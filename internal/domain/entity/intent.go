package entity

import (
	"time"

	"github.com/google/uuid"
)

type IntentType string

const (
	IntentAddActivity    IntentType = "add_activity"
	IntentDeleteActivity IntentType = "delete_activity"
	IntentUpdateActivity IntentType = "update_activity"
	IntentListActivities IntentType = "list_activities"
	IntentQuestion       IntentType = "question"
	IntentGreeting       IntentType = "greeting"
	IntentUnknown        IntentType = "unknown"
)

type ParsedIntent struct {
	Type       IntentType
	Confidence float64
	Entities   map[string]interface{}
}

type ActivityIntentData struct {
	Title         string
	Description   string
	ScheduledTime *time.Time
	CategoryID    *uuid.UUID
	Priority      int
}

type UpdateActivityIntentData struct {
	ActivityID    uuid.UUID
	Title         *string
	Description   *string
	ScheduledTime *time.Time
	Status        *string
	Priority      *int
}

