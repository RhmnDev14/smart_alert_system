package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	WhatsAppNumber     string     `json:"whatsapp_number" db:"whatsapp_number"`
	Name               string     `json:"name" db:"name"`
	Timezone           string     `json:"timezone" db:"timezone"`
	IsActive           bool       `json:"is_active" db:"is_active"`
	IsFirstTime        bool        `json:"is_first_time" db:"is_first_time"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	LastInteractionAt  *time.Time `json:"last_interaction_at" db:"last_interaction_at"`
}

func NewUser(whatsappNumber, name, timezone string) *User {
	now := time.Now()
	return &User{
		ID:            uuid.New(),
		WhatsAppNumber: whatsappNumber,
		Name:          name,
		Timezone:      timezone,
		IsActive:      true,
		IsFirstTime:   true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

