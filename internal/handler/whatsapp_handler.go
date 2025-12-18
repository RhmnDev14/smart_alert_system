package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/domain/repository"
	"smart_alert_system/internal/infrastructure/ai"
	"smart_alert_system/internal/infrastructure/whatsapp"
	"smart_alert_system/internal/usecase"
	"smart_alert_system/internal/utils"

	"github.com/google/uuid"
)

type WhatsAppHandler struct {
	userUseCase     *usecase.UserUseCase
	activityUseCase *usecase.ActivityUseCase
	aiService       ai.AIService
	wahaClient      *whatsapp.WahaClient
	messageRepo     repository.MessageRepository
	alertRepo       repository.AlertRepository
}

func NewWhatsAppHandler(
	userUseCase *usecase.UserUseCase,
	activityUseCase *usecase.ActivityUseCase,
	aiService ai.AIService,
	wahaClient *whatsapp.WahaClient,
	messageRepo repository.MessageRepository,
	alertRepo repository.AlertRepository,
) *WhatsAppHandler {
	return &WhatsAppHandler{
		userUseCase:     userUseCase,
		activityUseCase: activityUseCase,
		aiService:       aiService,
		wahaClient:      wahaClient,
		messageRepo:     messageRepo,
		alertRepo:       alertRepo,
	}
}

// Waha webhook payload structure
type WebhookPayload struct {
	ID        string      `json:"id"`
	Timestamp int64       `json:"timestamp"`
	Event     string      `json:"event"`
	Session   string      `json:"session"`
	Metadata  interface{} `json:"metadata"`
	Me        struct {
		ID       string `json:"id"`
		PushName string `json:"pushName"`
	} `json:"me"`
	Payload MessageData `json:"payload"`
}

// MessageData structure from Waha (actual format)
type MessageData struct {
	ID        string      `json:"id"` // String format: "false_25675515867262@lid_AC91F329..."
	Timestamp int64       `json:"timestamp"`
	From      string      `json:"from"` // Format: "25675515867262@lid" or "6281234567890@c.us"
	FromMe    bool        `json:"fromMe"`
	Source    string      `json:"source"`
	To        string      `json:"to"` // Format: "62881024952694@c.us"
	Body      string      `json:"body"`
	HasMedia  bool        `json:"hasMedia"`
	Media     interface{} `json:"media"`
	Ack       int         `json:"ack"`
	AckName   string      `json:"ackName"`
	Data      *struct {
		ID   *MessageID `json:"id"` // Nested ID structure in _data
		Type string     `json:"type"`
		Body string     `json:"body"`
		From string     `json:"from"`
		To   string     `json:"to"`
	} `json:"_data"`
}

func (h *WhatsAppHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Log incoming request
	log.Printf("=== Webhook Request Received ===")
	log.Printf("Method: %s", r.Method)
	log.Printf("URL: %s", r.URL.String())
	log.Printf("Remote Addr: %s", r.RemoteAddr)
	log.Printf("User-Agent: %s", r.UserAgent())
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

	if r.Method != http.MethodPost {
		log.Printf("❌ Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read body for debugging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Error reading request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Log raw payload
	log.Printf("📥 Raw Payload (length: %d bytes):", len(bodyBytes))
	log.Printf("%s", string(bodyBytes))

	var payload WebhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		log.Printf("❌ Error decoding webhook payload: %v", err)
		log.Printf("Attempting to parse as legacy format...")

		// Try alternative format (legacy format)
		var altPayload struct {
			Event string          `json:"event"`
			Data  json.RawMessage `json:"data"`
		}
		if err2 := json.Unmarshal(bodyBytes, &altPayload); err2 == nil && altPayload.Event == "message" {
			log.Printf("✓ Parsed as legacy format, event: %s", altPayload.Event)
			var messageData MessageData
			if err3 := json.Unmarshal(altPayload.Data, &messageData); err3 == nil {
				log.Printf("✓ Message data parsed successfully")
				log.Printf("  From: %s, Body: %s", messageData.From, messageData.Body)
				go h.processMessage(context.Background(), messageData)
				w.WriteHeader(http.StatusOK)
				return
			} else {
				log.Printf("❌ Error unmarshaling legacy message data: %v", err3)
			}
		}
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	log.Printf("✓ Payload parsed successfully")
	log.Printf("  Event: %s", payload.Event)
	log.Printf("  Session: %s", payload.Session)

	if payload.Event != "message" {
		log.Printf("⚠️  Ignoring non-message event: %s", payload.Event)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Log message details
	log.Printf("📨 Message Details:")
	log.Printf("  From: %s", payload.Payload.From)
	log.Printf("  To: %s", payload.Payload.To)
	log.Printf("  Body: %s", payload.Payload.Body)
	log.Printf("  Type: %s", getMessageType(payload.Payload))
	log.Printf("  Timestamp: %d", payload.Payload.Timestamp)
	log.Printf("  FromMe: %v", payload.Payload.FromMe)
	log.Printf("  ID: %s", payload.Payload.ID)

	// Process message asynchronously
	log.Printf("🚀 Processing message asynchronously...")
	go h.processMessage(context.Background(), payload.Payload)

	w.WriteHeader(http.StatusOK)
	log.Printf("✓ Response sent (200 OK)")
	log.Printf("=== End Webhook Request ===\n")
}

// Helper function to get message type
func getMessageType(msg MessageData) string {
	if msg.Data != nil && msg.Data.Type != "" {
		return msg.Data.Type
	}
	return "chat" // Default type
}

// Message ID structure (nested in _data)
type MessageID struct {
	FromMe     bool   `json:"fromMe"`
	Remote     string `json:"remote"`
	ID         string `json:"id"`
	Serialized string `json:"_serialized"`
}

func (h *WhatsAppHandler) processMessage(ctx context.Context, messageData MessageData) {
	log.Printf("🔄 Processing message...")

	// Extract WhatsApp number from "from" field
	// Format from Waha: "6281234567890@c.us" or "25675515867262@lid" or just number
	whatsappNumber := messageData.From
	if strings.Contains(whatsappNumber, "@") {
		whatsappNumber = strings.Split(whatsappNumber, "@")[0]
	}
	// Remove "lid" prefix if exists (WhatsApp Business)
	whatsappNumber = strings.TrimPrefix(whatsappNumber, "lid")
	log.Printf("  Extracted WhatsApp number: %s", whatsappNumber)

	// Only process if message is not from us
	if messageData.FromMe {
		log.Printf("⚠️  Ignoring message from self (fromMe: true)")
		return
	}

	messageContent := messageData.Body
	log.Printf("  Message content: %s", messageContent)

	// Get or create user
	log.Printf("  Getting or creating user: %s", whatsappNumber)
	user, err := h.userUseCase.GetOrCreateUser(ctx, whatsappNumber, "", "Asia/Jakarta")
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
		welcomeMsg := "Halo! Selamat datang di Smart Alert System. Saya akan membantu Anda mengelola kegiatan dan memberikan rekomendasi kesehatan.\n\nAnda bisa menambahkan kegiatan dengan format:\n• \"Besok saya akan olahraga jam 6 pagi\"\n• \"Hari ini ada meeting jam 2 siang\"\n• \"Tambah kegiatan [nama kegiatan] [waktu]\"\n\nSilakan coba kirim pesan untuk menambahkan kegiatan!"
		log.Printf("  Sending welcome message to: %s", whatsappNumber)

		// Use original 'from' format for sending message (with @lid or @c.us)
		sendTo := messageData.From
		if err := h.wahaClient.SendMessage(sendTo, welcomeMsg); err != nil {
			log.Printf("❌ Error sending welcome message: %v", err)
			log.Printf("  Tried sending to: %s", sendTo)
		} else {
			log.Printf("✓ Welcome message sent successfully")
			h.userUseCase.MarkAsNotFirstTime(ctx, user.ID)
			// After welcome message, also try to process the current message if it contains activity
			// This allows user to add activity in the first message
			log.Printf("  Processing first message for activity detection...")
			// Continue to process the message for activity
		}
		// Don't return early - continue processing the message to detect activity
	}

	// Parse intent with AI
	log.Printf("  Parsing intent with AI...")
	parsedIntent, err := h.aiService.ParseIntent(ctx, messageContent)
	if err != nil {
		log.Printf("⚠️  AI parsing failed, using fallback parser: %v", err)
		// Use fallback parser when AI fails
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
	// Use original 'from' format for sending message (with @lid or @c.us)
	sendTo := messageData.From
	log.Printf("  Sending response to: %s", sendTo)
	if err := h.wahaClient.SendMessage(sendTo, response); err != nil {
		log.Printf("❌ Error sending response: %v", err)
		log.Printf("  Tried sending to: %s", sendTo)
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

func (h *WhatsAppHandler) handleIntent(ctx context.Context, userID uuid.UUID, intent *entity.ParsedIntent, originalMessage string) (string, error) {
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
		return "Halo! Ada yang bisa saya bantu hari ini?", nil
	default:
		return "Maaf, saya belum memahami pesan Anda. Silakan coba lagi dengan format yang lebih jelas.", nil
	}
}

func (h *WhatsAppHandler) handleAddActivity(ctx context.Context, userID uuid.UUID, intent *entity.ParsedIntent) (string, error) {
	log.Printf("  📝 Processing add activity intent...")
	data := extractActivityData(intent.Entities, time.Now())

	// If title is empty, use description or ask user
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

	return fmt.Sprintf("✓ Kegiatan '%s' berhasil ditambahkan untuk %s",
		activity.Title, activity.ScheduledTime.Format("02 Jan 2006 15:04")), nil
}

func (h *WhatsAppHandler) handleDeleteActivity(ctx context.Context, userID uuid.UUID, intent *entity.ParsedIntent) (string, error) {
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

	return "✓ Kegiatan berhasil dihapus.", nil
}

func (h *WhatsAppHandler) handleUpdateActivity(ctx context.Context, userID uuid.UUID, intent *entity.ParsedIntent) (string, error) {
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

	return "✓ Kegiatan berhasil diupdate.", nil
}

func (h *WhatsAppHandler) handleListActivities(ctx context.Context, userID uuid.UUID) (string, error) {
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

func (h *WhatsAppHandler) handleQuestion(ctx context.Context, userID uuid.UUID, question string) (string, error) {
	// Use AI to answer general questions
	_, err := h.aiService.ParseIntent(ctx, question)
	if err != nil {
		return "Maaf, saya tidak dapat menjawab pertanyaan tersebut saat ini.", nil
	}

	// For now, return a simple response
	// In production, you'd want a dedicated Q&A endpoint
	return "Terima kasih atas pertanyaannya. Fitur ini sedang dalam pengembangan.", nil
}

func extractActivityData(entities map[string]interface{}, baseTime time.Time) entity.ActivityIntentData {
	data := entity.ActivityIntentData{}

	// Extract title
	if title, ok := entities["title"].(string); ok && title != "" {
		data.Title = title
	}

	// Extract description
	if desc, ok := entities["description"].(string); ok && desc != "" {
		data.Description = desc
	}

	// Extract scheduled_time
	if timeStr, ok := entities["scheduled_time"].(string); ok && timeStr != "" {
		// Try parsing as ISO 8601 first
		if parsedTime, err := utils.ParseISO8601Time(timeStr); err == nil {
			data.ScheduledTime = parsedTime
		} else {
			// Try parsing as natural language (Indonesian)
			if parsedTime, err := utils.ParseTimeFromText(timeStr, baseTime); err == nil && parsedTime != nil {
				data.ScheduledTime = parsedTime
			}
		}
	}

	// Extract priority
	if priority, ok := entities["priority"].(float64); ok {
		data.Priority = int(priority)
	} else if priorityStr, ok := entities["priority"].(string); ok {
		if p, err := strconv.Atoi(priorityStr); err == nil {
			data.Priority = p
		}
	}

	// Extract category_id if provided
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
