package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/domain/repository"
	"smart_alert_system/internal/infrastructure/ai"
	"smart_alert_system/internal/infrastructure/telegram"
	"smart_alert_system/internal/usecase"
	"smart_alert_system/internal/utils"

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
	log.Printf("Method: %s", r.Method)
	log.Printf("Remote Addr: %s", r.RemoteAddr)

	if r.Method != http.MethodPost {
		log.Printf("❌ Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Error reading request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("📥 Raw Payload (length: %d bytes):", len(bodyBytes))
	log.Printf("%s", string(bodyBytes))

	var update telegram.Update
	if err := json.Unmarshal(bodyBytes, &update); err != nil {
		log.Printf("❌ Error decoding update: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Only process text messages
	if update.Message == nil || update.Message.Text == "" {
		log.Printf("⚠️ Ignoring non-text update")
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("📨 Message Details:")
	log.Printf("  From: %s (ID: %d)", update.Message.From.FirstName, update.Message.From.ID)
	log.Printf("  Chat ID: %d", update.Message.Chat.ID)
	log.Printf("  Text: %s", update.Message.Text)

	// Process message asynchronously
	go h.processMessage(context.Background(), update.Message)

	w.WriteHeader(http.StatusOK)
	log.Printf("✓ Response sent (200 OK)")
	log.Printf("=== End Telegram Webhook ===\n")
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

				// Skip bot commands that start with /start (handled separately)
				log.Printf("📨 New message from %s: %s", update.Message.From.FirstName, update.Message.Text)

				go h.processMessage(ctx, update.Message)
			}
		}
	}
}

func (h *TelegramHandler) processMessage(ctx context.Context, message *telegram.Message) {
	log.Printf("🔄 Processing message...")

	// Use chat ID as the unique identifier (like WhatsApp number)
	chatIDStr := fmt.Sprintf("%d", message.Chat.ID)
	senderName := message.From.FirstName
	if message.From.LastName != "" {
		senderName += " " + message.From.LastName
	}
	log.Printf("  Chat ID: %s, Sender: %s", chatIDStr, senderName)

	// Only process if message is not from a bot
	if message.From.IsBot {
		log.Printf("⚠️ Ignoring message from bot")
		return
	}

	messageContent := message.Text
	log.Printf("  Message content: %s", messageContent)

	// Get or create user (using chat ID as the identifier, stored in whatsapp_number field)
	log.Printf("  Getting or creating user: %s", chatIDStr)
	user, err := h.userUseCase.GetOrCreateUser(ctx, chatIDStr, senderName, "Asia/Jakarta")
	if err != nil {
		log.Printf("❌ Error getting/creating user: %v", err)
		return
	}
	log.Printf("  ✓ User ID: %s, IsFirstTime: %v", user.ID, user.IsFirstTime)

	// Save incoming message
	now := time.Now()
	messageHistory := entity.NewMessageHistory(user.ID, messageContent, entity.MessageTypeIncoming)
	messageHistory.ReceivedAt = &now
	if err := h.messageRepo.Create(ctx, messageHistory); err != nil {
		log.Printf("Error saving message: %v", err)
	}

	// Check if first time user
	if user.IsFirstTime {
		welcomeMsg := "Halo! Selamat datang di Smart Alert System 🤖\n\nSaya akan membantu Anda mengelola kegiatan dan memberikan rekomendasi kesehatan.\n\nAnda bisa menambahkan kegiatan dengan format:\n• \"Besok saya akan olahraga jam 6 pagi\"\n• \"Hari ini ada meeting jam 2 siang\"\n• \"Tambah kegiatan [nama kegiatan] [waktu]\"\n\nSilakan coba kirim pesan untuk menambahkan kegiatan!"
		log.Printf("  Sending welcome message to chat: %d", message.Chat.ID)

		if err := h.telegramClient.SendMessage(message.Chat.ID, welcomeMsg); err != nil {
			log.Printf("❌ Error sending welcome message: %v", err)
		} else {
			log.Printf("✓ Welcome message sent successfully")
			h.userUseCase.MarkAsNotFirstTime(ctx, user.ID)
		}
	}

	// Parse intent with AI
	log.Printf("  Parsing intent with AI...")
	parsedIntent, err := h.aiService.ParseIntent(ctx, messageContent)
	if err != nil {
		log.Printf("⚠️ AI parsing failed, using fallback parser: %v", err)
		parsedIntent = utils.FallbackIntentParser(messageContent, time.Now())
		log.Printf("  ✓ Fallback intent detected: %s (confidence: %.2f)", parsedIntent.Type, parsedIntent.Confidence)
		if len(parsedIntent.Entities) > 0 {
			log.Printf("  Entities: %+v", parsedIntent.Entities)
		}
	} else {
		log.Printf("  ✓ Intent detected: %s (confidence: %.2f)", parsedIntent.Type, parsedIntent.Confidence)
		if len(parsedIntent.Entities) > 0 {
			log.Printf("  Entities: %+v", parsedIntent.Entities)
		}
	}

	messageHistory.IntentDetected = string(parsedIntent.Type)
	messageHistory.IsProcessed = true
	h.messageRepo.Update(ctx, messageHistory)

	// Handle intent
	log.Printf("  Handling intent: %s", parsedIntent.Type)
	response, err := h.handleIntent(ctx, user.ID, parsedIntent, messageContent)
	if err != nil {
		log.Printf("❌ Error handling intent: %v", err)
		response = "Maaf, terjadi kesalahan. Silakan coba lagi."
	} else {
		log.Printf("  ✓ Response generated: %s", response)
	}

	// Send response
	log.Printf("  Sending response to chat: %d", message.Chat.ID)
	if err := h.telegramClient.SendMessage(message.Chat.ID, response); err != nil {
		log.Printf("❌ Error sending response: %v", err)
	} else {
		log.Printf("✓ Response sent successfully")
		// Save outgoing message
		outgoingMsg := entity.NewMessageHistory(user.ID, response, entity.MessageTypeOutgoing)
		sentAt := time.Now()
		outgoingMsg.SentAt = &sentAt
		outgoingMsg.AIResponse = response
		h.messageRepo.Create(ctx, outgoingMsg)
	}
}

func (h *TelegramHandler) handleIntent(ctx context.Context, userID uuid.UUID, intent *entity.ParsedIntent, originalMessage string) (string, error) {
	switch intent.Type {
	case entity.IntentAddActivity:
		return h.handleAddActivity(ctx, userID, intent)
	case entity.IntentDeleteActivity:
		return h.handleDeleteActivity(ctx, userID, intent)
	case entity.IntentUpdateActivity:
		return h.handleUpdateActivity(ctx, userID, intent)
	case entity.IntentListActivities:
		return h.handleListActivities(ctx, userID)
	case entity.IntentQuestion:
		return h.handleQuestion(ctx, userID, originalMessage)
	case entity.IntentGreeting:
		return "Halo! Ada yang bisa saya bantu hari ini? 😊", nil
	default:
		return "Maaf, saya belum memahami pesan Anda. Silakan coba lagi dengan format yang lebih jelas.", nil
	}
}

func (h *TelegramHandler) handleAddActivity(ctx context.Context, userID uuid.UUID, intent *entity.ParsedIntent) (string, error) {
	log.Printf("  📝 Processing add activity intent...")
	data := extractActivityData(intent.Entities, time.Now())

	if data.Title == "" {
		if data.Description != "" {
			data.Title = data.Description
			data.Description = ""
		} else {
			return "Maaf, saya tidak dapat menemukan judul kegiatan. Silakan coba lagi dengan format: 'Tambah kegiatan [judul] [waktu]'", nil
		}
	}

	log.Printf("  Activity data: Title=%s, Description=%s, ScheduledTime=%v, Priority=%d",
		data.Title, data.Description, data.ScheduledTime, data.Priority)

	activity, err := h.activityUseCase.CreateActivity(ctx, userID, data)
	if err != nil {
		log.Printf("❌ Failed to create activity: %v", err)
		return "", fmt.Errorf("failed to create activity: %w", err)
	}

	log.Printf("✓ Activity created successfully: ID=%s, Title=%s, ScheduledTime=%s",
		activity.ID, activity.Title, activity.ScheduledTime.Format("02 Jan 2006 15:04"))

	return fmt.Sprintf("✅ Kegiatan '%s' berhasil ditambahkan untuk %s",
		activity.Title, activity.ScheduledTime.Format("02 Jan 2006 15:04")), nil
}

func (h *TelegramHandler) handleDeleteActivity(ctx context.Context, userID uuid.UUID, intent *entity.ParsedIntent) (string, error) {
	activityIDStr, ok := intent.Entities["activity_id"].(string)
	if !ok {
		return "Maaf, ID kegiatan tidak ditemukan. Silakan coba lagi.", nil
	}

	activityID, err := uuid.Parse(activityIDStr)
	if err != nil {
		return "Maaf, ID kegiatan tidak valid.", nil
	}

	if err := h.activityUseCase.DeleteActivity(ctx, activityID); err != nil {
		return "", fmt.Errorf("failed to delete activity: %w", err)
	}

	return "✅ Kegiatan berhasil dihapus.", nil
}

func (h *TelegramHandler) handleUpdateActivity(ctx context.Context, userID uuid.UUID, intent *entity.ParsedIntent) (string, error) {
	activityIDStr, ok := intent.Entities["activity_id"].(string)
	if !ok {
		return "Maaf, ID kegiatan tidak ditemukan.", nil
	}

	activityID, err := uuid.Parse(activityIDStr)
	if err != nil {
		return "Maaf, ID kegiatan tidak valid.", nil
	}

	data := extractUpdateActivityData(intent.Entities)
	if err := h.activityUseCase.UpdateActivity(ctx, activityID, data); err != nil {
		return "", fmt.Errorf("failed to update activity: %w", err)
	}

	return "✅ Kegiatan berhasil diupdate.", nil
}

func (h *TelegramHandler) handleListActivities(ctx context.Context, userID uuid.UUID) (string, error) {
	activities, err := h.activityUseCase.GetTodayActivities(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get activities: %w", err)
	}

	if len(activities) == 0 {
		return "Anda tidak memiliki kegiatan untuk hari ini.", nil
	}

	response := "📋 Kegiatan Hari Ini:\n\n"
	for i, activity := range activities {
		response += fmt.Sprintf("%d. %s - %s\n   Waktu: %s\n   Status: %s\n\n",
			i+1, activity.Title, activity.Description,
			activity.ScheduledTime.Format("15:04"), activity.Status)
	}

	return response, nil
}

func (h *TelegramHandler) handleQuestion(ctx context.Context, userID uuid.UUID, question string) (string, error) {
	_, err := h.aiService.ParseIntent(ctx, question)
	if err != nil {
		return "Maaf, saya tidak dapat menjawab pertanyaan tersebut saat ini.", nil
	}

	return "Terima kasih atas pertanyaannya. Fitur ini sedang dalam pengembangan.", nil
}

func extractActivityData(entities map[string]interface{}, baseTime time.Time) entity.ActivityIntentData {
	data := entity.ActivityIntentData{}

	if title, ok := entities["title"].(string); ok && title != "" {
		data.Title = title
	}

	if desc, ok := entities["description"].(string); ok && desc != "" {
		data.Description = desc
	}

	if timeStr, ok := entities["scheduled_time"].(string); ok && timeStr != "" {
		if parsedTime, err := utils.ParseISO8601Time(timeStr); err == nil {
			data.ScheduledTime = parsedTime
		} else {
			if parsedTime, err := utils.ParseTimeFromText(timeStr, baseTime); err == nil && parsedTime != nil {
				data.ScheduledTime = parsedTime
			}
		}
	}

	if priority, ok := entities["priority"].(float64); ok {
		data.Priority = int(priority)
	} else if priorityStr, ok := entities["priority"].(string); ok {
		if p, err := strconv.Atoi(priorityStr); err == nil {
			data.Priority = p
		}
	}

	if catIDStr, ok := entities["category_id"].(string); ok {
		if catID, err := uuid.Parse(catIDStr); err == nil {
			data.CategoryID = &catID
		}
	}

	return data
}

func extractUpdateActivityData(entities map[string]interface{}) entity.UpdateActivityIntentData {
	data := entity.UpdateActivityIntentData{}

	if idStr, ok := entities["activity_id"].(string); ok {
		if id, err := uuid.Parse(idStr); err == nil {
			data.ActivityID = id
		}
	}
	if title, ok := entities["title"].(string); ok {
		data.Title = &title
	}
	if desc, ok := entities["description"].(string); ok {
		data.Description = &desc
	}

	return data
}
