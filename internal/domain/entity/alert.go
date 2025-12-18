package entity

import (
	"time"

	"github.com/google/uuid"
)

type AlertType string

const (
	AlertTypeMorning      AlertType = "morning_alert"
	AlertTypeEvening      AlertType = "evening_summary"
	AlertTypeActivityReminder AlertType = "activity_reminder"
)

type AlertStatus string

const (
	AlertStatusPending AlertStatus = "pending"
	AlertStatusSent    AlertStatus = "sent"
	AlertStatusFailed  AlertStatus = "failed"
)

type AlertLog struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	AlertType    AlertType  `json:"alert_type" db:"alert_type"`
	AlertContent string     `json:"alert_content" db:"alert_content"`
	ScheduledTime time.Time `json:"scheduled_time" db:"scheduled_time"`
	SentAt       *time.Time `json:"sent_at" db:"sent_at"`
	IsSent       bool       `json:"is_sent" db:"is_sent"`
	Status       AlertStatus `json:"status" db:"status"`
	ErrorMessage string     `json:"error_message" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

func NewAlertLog(userID uuid.UUID, alertType AlertType, content string, scheduledTime time.Time) *AlertLog {
	return &AlertLog{
		ID:           uuid.New(),
		UserID:       userID,
		AlertType:    alertType,
		AlertContent: content,
		ScheduledTime: scheduledTime,
		IsSent:       false,
		Status:       AlertStatusPending,
		CreatedAt:    time.Now(),
	}
}

func (a *AlertLog) MarkSent() {
	now := time.Now()
	a.IsSent = true
	a.Status = AlertStatusSent
	a.SentAt = &now
}

func (a *AlertLog) MarkFailed(err error) {
	a.Status = AlertStatusFailed
	if err != nil {
		a.ErrorMessage = err.Error()
	}
}

