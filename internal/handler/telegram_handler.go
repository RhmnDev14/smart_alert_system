package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/domain/repository"
	"smart_alert_system/internal/infrastructure/ai"
	"smart_alert_system/internal/infrastructure/telegram"
	"smart_alert_system/internal/usecase"

	"github.com/google/uuid"
)

type TelegramHandler struct {
	userUseCase     *usecase.UserUseCase
	activityUseCase *usecase.ActivityUseCase
	aiService       ai.AIService
	telegramClient  *telegram.TelegramClient
	messageRepo     repository.MessageRepository
	alertRepo       repository.AlertRepository
	memoryRepo      repository.MemoryRepository
}

func NewTelegramHandler(
	userUseCase *usecase.UserUseCase,
	activityUseCase *usecase.ActivityUseCase,
	aiService ai.AIService,
	telegramClient *telegram.TelegramClient,
	messageRepo repository.MessageRepository,
	alertRepo repository.AlertRepository,
	memoryRepo repository.MemoryRepository,
) *TelegramHandler {
	return &TelegramHandler{
		userUseCase:     userUseCase,
		activityUseCase: activityUseCase,
		aiService:       aiService,
		telegramClient:  telegramClient,
		messageRepo:     messageRepo,
		alertRepo:       alertRepo,
		memoryRepo:      memoryRepo,
	}
}

// HandleWebhook handles incoming Telegram webhook updates (POST /webhook)
func (h *TelegramHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Telegram Webhook Received ===")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var update telegram.Update
	if err := json.Unmarshal(bodyBytes, &update); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if update.Message == nil || update.Message.Text == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	go h.processMessage(context.Background(), update.Message)

	w.WriteHeader(http.StatusOK)
}

// StartLongPolling starts the Telegram long-polling loop to receive messages
func (h *TelegramHandler) StartLongPolling(ctx context.Context) {
	log.Println("🤖 Starting Telegram long-polling...")
	offset := 0

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Telegram long-polling...")
			return
		default:
			updates, err := h.telegramClient.GetUpdates(offset)
			if err != nil {
				log.Printf("❌ Error getting updates: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, update := range updates {
				offset = update.UpdateID + 1

				if update.Message == nil || update.Message.Text == "" {
					continue
				}

				log.Printf("📨 New message from %s: %s", update.Message.From.FirstName, update.Message.Text)
				go h.processMessage(ctx, update.Message)
			}
		}
	}
}

// processMessage is the main gateway: processes every message through AI
func (h *TelegramHandler) processMessage(ctx context.Context, message *telegram.Message) {
	log.Printf("🔄 Processing message from %s (Chat: %d)...", message.From.FirstName, message.Chat.ID)

	// Build user identifier from chat ID
	chatIDStr := fmt.Sprintf("%d", message.Chat.ID)
	senderName := message.From.FirstName
	if message.From.LastName != "" {
		senderName += " " + message.From.LastName
	}

	// Skip bot messages
	if message.From.IsBot {
		log.Printf("⚠️ Ignoring message from bot")
		return
	}

	messageContent := message.Text
	log.Printf("  📝 Message: %s", messageContent)

	// Step 1: Get or create user
	user, err := h.userUseCase.GetOrCreateUser(ctx, chatIDStr, senderName, "Asia/Jakarta")
	if err != nil {
		log.Printf("❌ Error getting/creating user: %v", err)
		return
	}
	log.Printf("  ✓ User: %s (ID: %s, FirstTime: %v)", user.Name, user.ID, user.IsFirstTime)

	// Step 2: Save incoming message to history
	now := time.Now()
	messageHistory := entity.NewMessageHistory(user.ID, messageContent, entity.MessageTypeIncoming)
	messageHistory.ReceivedAt = &now
	if err := h.messageRepo.Create(ctx, messageHistory); err != nil {
		log.Printf("⚠️ Error saving message history: %v", err)
	}

	// Step 3: Handle first-time user with welcome message
	if user.IsFirstTime {
		welcomeMsg := "Halo! Selamat datang di Smart Alert Bot 🤖\n\n" +
			"Saya adalah asisten pintar yang bisa:\n" +
			"• 💬 Ngobrol dan menjawab pertanyaan\n" +
			"• 📅 Otomatis mendeteksi dan menyimpan jadwal kegiatan\n" +
			"• ⏰ Mengingatkan kegiatan yang sudah dijadwalkan\n" +
			"• 🔍 Menampilkan dan mencari jadwal kegiatan\n" +
			"• ✏️ Mengubah, membatalkan, atau menandai kegiatan selesai\n\n" +
			"Cukup kirim pesan apa saja! Contoh:\n" +
			"• \"Besok ada meeting jam 10 pagi\"\n" +
			"• \"Lihat jadwal hari ini\"\n" +
			"• \"Cancel meeting besok\"\n" +
			"• \"Tandai olahraga sebagai selesai\"\n" +
			"• Atau ngobrol biasa juga boleh 😊\n\n" +
			"Silakan mulai!\n\n" +
			"Smart Alert System\n" +
			"Develop by Rahman Umardi"

		if err := h.telegramClient.SendMessage(message.Chat.ID, welcomeMsg); err != nil {
			log.Printf("❌ Error sending welcome message: %v", err)
		} else {
			log.Printf("✓ Welcome message sent")
			h.userUseCase.MarkAsNotFirstTime(ctx, user.ID)
		}
	}

	// Step 4: Fetch recent conversation history and user memories for AI context
	log.Printf("  🤖 Processing through AI gateway...")
	currentTime := time.Now()
	chatHistory, err := h.messageRepo.GetByUserID(ctx, user.ID, 10)
	if err != nil {
		log.Printf("⚠️ Error fetching chat history: %v", err)
		chatHistory = nil
	}
	userMemories, err := h.memoryRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Printf("⚠️ Error fetching user memories: %v", err)
		userMemories = nil
	}
	if len(userMemories) > 0 {
		log.Printf("  🧠 Loaded %d memories for user", len(userMemories))
	}
	aiResult, err := h.aiService.ProcessMessage(ctx, messageContent, user.Name, currentTime, chatHistory, userMemories)
	if err != nil {
		log.Printf("⚠️ AI processing error: %v", err)
	}

	if aiResult == nil {
		aiResult = &entity.AIProcessResult{
			Response:    "Maaf, saya sedang mengalami gangguan. Silakan coba lagi nanti.",
			HasSchedule: false,
			Action:      entity.ActionNone,
		}
	}

	// Step 4b: Save any new memories extracted by AI
	if len(aiResult.MemoriesToSave) > 0 {
		for _, mem := range aiResult.MemoriesToSave {
			if mem.Content == "" {
				continue
			}
			memType := entity.MemoryType(mem.Type)
			if memType == "" {
				memType = entity.MemoryTypeFact
			}
			newMemory := entity.NewUserMemory(user.ID, memType, mem.Content)
			if err := h.memoryRepo.Create(ctx, newMemory); err != nil {
				log.Printf("⚠️ Error saving memory: %v", err)
			} else {
				log.Printf("  🧠 Memory saved: [%s] %s", mem.Type, mem.Content)
			}
		}
	}

	// Step 5: Route based on detected action
	var finalResponse string
	switch aiResult.Action {
	case entity.ActionCreate:
		log.Printf("  📅 CREATE action detected")
		finalResponse = h.handleCreateAction(ctx, user.ID, message.Chat.ID, aiResult)

	case entity.ActionGet:
		log.Printf("  🔍 GET action detected")
		finalResponse = h.handleGetAction(ctx, user.ID, aiResult)

	case entity.ActionUpdate:
		log.Printf("  ✏️ UPDATE action detected")
		finalResponse = h.handleUpdateAction(ctx, user.ID, aiResult)

	default:
		log.Printf("  💬 No data operation, regular conversation")
		finalResponse = aiResult.Response
	}

	// Step 6: Update message history with AI response
	messageHistory.AIResponse = finalResponse
	messageHistory.IsProcessed = true
	switch aiResult.Action {
	case entity.ActionCreate:
		messageHistory.IntentDetected = "schedule_create"
	case entity.ActionGet:
		messageHistory.IntentDetected = "schedule_get"
	case entity.ActionUpdate:
		messageHistory.IntentDetected = "schedule_update"
	default:
		messageHistory.IntentDetected = "conversation"
	}
	h.messageRepo.Update(ctx, messageHistory)

	// Step 7: Send response to user
	log.Printf("  📤 Sending response to chat %d...", message.Chat.ID)
	if err := h.telegramClient.SendMessage(message.Chat.ID, finalResponse); err != nil {
		log.Printf("❌ Error sending response: %v", err)
	} else {
		log.Printf("✓ Response sent successfully")

		// Save outgoing message
		outgoingMsg := entity.NewMessageHistory(user.ID, finalResponse, entity.MessageTypeOutgoing)
		sentAt := time.Now()
		outgoingMsg.SentAt = &sentAt
		outgoingMsg.AIResponse = finalResponse
		h.messageRepo.Create(ctx, outgoingMsg)
	}
}

// handleCreateAction saves a new schedule from AI extraction
func (h *TelegramHandler) handleCreateAction(ctx context.Context, userID uuid.UUID, chatID int64, aiResult *entity.AIProcessResult) string {
	if aiResult.Schedule == nil {
		log.Printf("⚠️ CREATE action but no schedule data")
		return aiResult.Response
	}

	h.saveScheduleFromAI(ctx, userID, aiResult.Schedule)
	return aiResult.Response
}

// handleGetAction queries activities for the user and formats them into the response
func (h *TelegramHandler) handleGetAction(ctx context.Context, userID uuid.UUID, aiResult *entity.AIProcessResult) string {
	activities, err := h.activityUseCase.GetActivities(ctx, userID, aiResult.Query)
	if err != nil {
		log.Printf("❌ Error getting activities: %v", err)
		return aiResult.Response + "\n\n⚠️ Maaf, terjadi kesalahan saat mengambil data kegiatan."
	}

	if len(activities) == 0 {
		return aiResult.Response + "\n\n📭 Tidak ada kegiatan yang ditemukan."
	}

	// Format activities into a readable list
	var sb strings.Builder
	sb.WriteString(aiResult.Response)
	sb.WriteString("\n\n📋 Daftar Kegiatan:\n")
	sb.WriteString("────────────────\n")

	loc, _ := time.LoadLocation("Asia/Jakarta")
	for i, act := range activities {
		localTime := act.ScheduledTime.In(loc)
		statusEmoji := getStatusEmoji(act.Status)
		priorityStr := getPriorityLabel(act.Priority)

		sb.WriteString(fmt.Sprintf("%d. %s %s\n", i+1, statusEmoji, act.Title))
		sb.WriteString(fmt.Sprintf("   🕐 %s\n", localTime.Format("02 Jan 2006, 15:04")))
		sb.WriteString(fmt.Sprintf("   📊 Status: %s | Prioritas: %s\n", act.Status, priorityStr))
		if act.Description != "" {
			sb.WriteString(fmt.Sprintf("   📝 %s\n", act.Description))
		}
		if i < len(activities)-1 {
			sb.WriteString("\n")
		}
	}

	sb.WriteString(fmt.Sprintf("\n📈 Total: %d kegiatan", len(activities)))

	return sb.String()
}

// handleUpdateAction finds and updates an activity based on AI-extracted update data
func (h *TelegramHandler) handleUpdateAction(ctx context.Context, userID uuid.UUID, aiResult *entity.AIProcessResult) string {
	if aiResult.Update == nil {
		log.Printf("⚠️ UPDATE action but no update data")
		return aiResult.Response
	}

	updatedActivity, err := h.activityUseCase.UpdateActivity(ctx, userID, aiResult.Update)
	if err != nil {
		log.Printf("❌ Error updating activity: %v", err)
		return aiResult.Response + fmt.Sprintf("\n\n⚠️ Gagal mengupdate kegiatan: %v", err)
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	localTime := updatedActivity.ScheduledTime.In(loc)
	statusEmoji := getStatusEmoji(updatedActivity.Status)

	return aiResult.Response + fmt.Sprintf("\n\n✅ Kegiatan berhasil diupdate:\n%s %s\n🕐 %s\n📊 Status: %s",
		statusEmoji, updatedActivity.Title, localTime.Format("02 Jan 2006, 15:04"), updatedActivity.Status)
}

// saveScheduleFromAI saves schedule data extracted by AI to the database
func (h *TelegramHandler) saveScheduleFromAI(ctx context.Context, userID uuid.UUID, scheduleData *entity.ScheduleData) {
	if scheduleData.Title == "" {
		log.Printf("⚠️ Schedule title is empty, skipping save")
		return
	}

	// Parse the scheduled time
	scheduledTime, err := scheduleData.ParseScheduledTime()
	if err != nil {
		log.Printf("⚠️ Failed to parse scheduled time '%s': %v", scheduleData.ScheduledTime, err)
		// Use default: 1 hour from now
		defaultTime := time.Now().Add(1 * time.Hour)
		scheduledTime = &defaultTime
	}
	if scheduledTime == nil {
		defaultTime := time.Now().Add(1 * time.Hour)
		scheduledTime = &defaultTime
	}

	// Set default priority
	priority := scheduleData.Priority
	if priority == 0 {
		priority = 3
	}

	// Create activity data
	data := entity.ActivityIntentData{
		Title:         scheduleData.Title,
		Description:   scheduleData.Description,
		ScheduledTime: scheduledTime,
		Priority:      priority,
	}

	// Save to database
	activity, err := h.activityUseCase.CreateActivity(ctx, userID, data)
	if err != nil {
		log.Printf("❌ Failed to save schedule: %v", err)
		return
	}

	log.Printf("✅ Schedule saved: ID=%s, Title='%s', Time=%s",
		activity.ID, activity.Title, activity.ScheduledTime.Format("02 Jan 2006 15:04"))
}

// getStatusEmoji returns an emoji for the activity status
func getStatusEmoji(status entity.ActivityStatus) string {
	switch status {
	case entity.ActivityStatusPending:
		return "⏳"
	case entity.ActivityStatusCompleted:
		return "✅"
	case entity.ActivityStatusCancelled:
		return "❌"
	case entity.ActivityStatusOverdue:
		return "⚠️"
	default:
		return "📌"
	}
}

// getPriorityLabel returns a human-readable label for priority
func getPriorityLabel(priority int) string {
	switch {
	case priority <= 1:
		return "Rendah"
	case priority == 2:
		return "Sedang-Rendah"
	case priority == 3:
		return "Sedang"
	case priority == 4:
		return "Tinggi"
	case priority >= 5:
		return "Sangat Tinggi"
	default:
		return "Sedang"
	}
}
