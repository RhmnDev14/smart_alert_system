package entity

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeIncoming MessageType = "incoming"
	MessageTypeOutgoing MessageType = "outgoing"
)

type MessageHistory struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	MessageContent string    `json:"message_content" db:"message_content"`
	MessageType   MessageType `json:"message_type" db:"message_type"`
	IntentDetected string    `json:"intent_detected" db:"intent_detected"`
	AIResponse     string    `json:"ai_response" db:"ai_response"`
	ReceivedAt     *time.Time `json:"received_at" db:"received_at"`
	SentAt         *time.Time `json:"sent_at" db:"sent_at"`
	IsProcessed    bool       `json:"is_processed" db:"is_processed"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

func NewMessageHistory(userID uuid.UUID, content string, msgType MessageType) *MessageHistory {
	now := time.Now()
	return &MessageHistory{
		ID:            uuid.New(),
		UserID:        userID,
		MessageContent: content,
		MessageType:   msgType,
		IsProcessed:   false,
		CreatedAt:     now,
	}
}

