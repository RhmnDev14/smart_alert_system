package entity

import (
	"time"

	"github.com/google/uuid"
)

type ActivityStatus string

const (
	ActivityStatusPending   ActivityStatus = "pending"
	ActivityStatusCompleted ActivityStatus = "completed"
	ActivityStatusCancelled ActivityStatus = "cancelled"
	ActivityStatusOverdue   ActivityStatus = "overdue"
)

type Activity struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	UserID        uuid.UUID      `json:"user_id" db:"user_id"`
	CategoryID    *uuid.UUID     `json:"category_id" db:"category_id"`
	Title         string         `json:"title" db:"title"`
	Description   string         `json:"description" db:"description"`
	ScheduledTime time.Time      `json:"scheduled_time" db:"scheduled_time"`
	ReminderTime  *time.Time     `json:"reminder_time" db:"reminder_time"`
	Status        ActivityStatus `json:"status" db:"status"`
	Priority      int            `json:"priority" db:"priority"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	CompletedAt   *time.Time     `json:"completed_at" db:"completed_at"`
}

func NewActivity(userID uuid.UUID, title, description string, scheduledTime time.Time, priority int) *Activity {
	now := time.Now()
	return &Activity{
		ID:            uuid.New(),
		UserID:        userID,
		Title:         title,
		Description:   description,
		ScheduledTime: scheduledTime,
		Status:        ActivityStatusPending,
		Priority:      priority,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (a *Activity) Complete() {
	now := time.Now()
	a.Status = ActivityStatusCompleted
	a.CompletedAt = &now
	a.UpdatedAt = now
}

func (a *Activity) Cancel() {
	a.Status = ActivityStatusCancelled
	a.UpdatedAt = time.Now()
}

func (a *Activity) MarkOverdue() {
	a.Status = ActivityStatusOverdue
	a.UpdatedAt = time.Now()
}

