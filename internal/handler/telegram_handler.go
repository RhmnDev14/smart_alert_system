package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
}

func NewTelegramHandler(
	userUseCase *usecase.UserUseCase,
	activityUseCase *usecase.ActivityUseCase,
	aiService ai.AIService,
	telegramClient *telegram.TelegramClient,
	messageRepo repository.MessageRepository,
	alertRepo repository.AlertRepository,
) *TelegramHandler {
	return &TelegramHandler{
		userUseCase:     userUseCase,
		activityUseCase: activityUseCase,
		aiService:       aiService,
		telegramClient:  telegramClient,
		messageRepo:     messageRepo,
		alertRepo:       alertRepo,
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
			"• ⏰ Mengingatkan kegiatan yang sudah dijadwalkan\n\n" +
			"Cukup kirim pesan apa saja! Contoh:\n" +
			"• \"Besok ada meeting jam 10 pagi\"\n" +
			"• \"Hari ini mau olahraga jam 6 sore\"\n" +
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

	// Step 4: Process message through AI Gateway
	log.Printf("  🤖 Processing through AI gateway...")
	currentTime := time.Now()
	aiResult, err := h.aiService.ProcessMessage(ctx, messageContent, currentTime)
	if err != nil {
		log.Printf("⚠️ AI processing error: %v", err)
		// AI returned a fallback response
	}

	if aiResult == nil {
		aiResult = &entity.AIProcessResult{
			Response:    "Maaf, saya sedang mengalami gangguan. Silakan coba lagi nanti.",
			HasSchedule: false,
		}
	}

	// Step 5: If AI detected a schedule, save it to database
	if aiResult.HasSchedule && aiResult.Schedule != nil {
		log.Printf("  📅 Schedule detected! Title: %s", aiResult.Schedule.Title)
		h.saveScheduleFromAI(ctx, user.ID, aiResult.Schedule)
	} else {
		log.Printf("  💬 No schedule detected, regular conversation")
	}

	// Step 6: Update message history with AI response
	messageHistory.AIResponse = aiResult.Response
	messageHistory.IsProcessed = true
	if aiResult.HasSchedule {
		messageHistory.IntentDetected = "schedule_detected"
	} else {
		messageHistory.IntentDetected = "conversation"
	}
	h.messageRepo.Update(ctx, messageHistory)

	// Step 7: Send response to user
	log.Printf("  📤 Sending response to chat %d...", message.Chat.ID)
	if err := h.telegramClient.SendMessage(message.Chat.ID, aiResult.Response); err != nil {
		log.Printf("❌ Error sending response: %v", err)
	} else {
		log.Printf("✓ Response sent successfully")

		// Save outgoing message
		outgoingMsg := entity.NewMessageHistory(user.ID, aiResult.Response, entity.MessageTypeOutgoing)
		sentAt := time.Now()
		outgoingMsg.SentAt = &sentAt
		outgoingMsg.AIResponse = aiResult.Response
		h.messageRepo.Create(ctx, outgoingMsg)
	}
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
