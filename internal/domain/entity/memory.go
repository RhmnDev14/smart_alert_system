package entity

import (
	"time"

	"github.com/google/uuid"
)

// MemoryType represents the category of a user memory
type MemoryType string

const (
	MemoryTypePreference MemoryType = "preference" // e.g., "User prefers jogging at 5 AM"
	MemoryTypeFact       MemoryType = "fact"       // e.g., "User's wife name is Siti"
	MemoryTypeHabit      MemoryType = "habit"      // e.g., "User usually has meetings on Monday"
	MemoryTypePersonal   MemoryType = "personal"   // e.g., "User likes spicy food"
)

// UserMemory represents a piece of information the AI remembers about a user
type UserMemory struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	MemoryType MemoryType `json:"memory_type" db:"memory_type"`
	Content    string     `json:"content" db:"content"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

func NewUserMemory(userID uuid.UUID, memoryType MemoryType, content string) *UserMemory {
	now := time.Now()
	return &UserMemory{
		ID:         uuid.New(),
		UserID:     userID,
		MemoryType: memoryType,
		Content:    content,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
