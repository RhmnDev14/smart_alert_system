package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserHealthProfile struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	Age                *int      `json:"age" db:"age"`
	Gender             string    `json:"gender" db:"gender"`
	MedicalConditions  string    `json:"medical_conditions" db:"medical_conditions"` // JSONB
	Allergies          string    `json:"allergies" db:"allergies"`                   // JSONB
	Medications        string    `json:"medications" db:"medications"`               // JSONB
	ActivityPreferences string   `json:"activity_preferences" db:"activity_preferences"` // JSONB
	HealthGoals        string    `json:"health_goals" db:"health_goals"`              // JSONB
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type RecommendationType struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	TriggerCondition string    `json:"trigger_condition" db:"trigger_condition"`
}

type HealthRecommendation struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	UserID              uuid.UUID  `json:"user_id" db:"user_id"`
	RecommendationTypeID *uuid.UUID `json:"recommendation_type_id" db:"recommendation_type_id"`
	RecommendationText  string     `json:"recommendation_text" db:"recommendation_text"`
	ActivityID          *uuid.UUID `json:"activity_id" db:"activity_id"`
	GeneratedAt         time.Time  `json:"generated_at" db:"generated_at"`
	SentAt              *time.Time `json:"sent_at" db:"sent_at"`
	IsRead              bool       `json:"is_read" db:"is_read"`
	Priority            int        `json:"priority" db:"priority"`
}

