package entity

import "time"

// ActionType represents the type of action detected by AI
type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionGet    ActionType = "get"
	ActionUpdate ActionType = "update"
	ActionNone   ActionType = "none" // pure conversation, no data operation
)

// AIProcessResult represents the result from AI processing a user message.
// The AI acts as a gateway: it responds to the user AND detects data operations (create/get/update).
type AIProcessResult struct {
	// Response is the natural language response to send back to the user
	Response string `json:"response"`

	// HasSchedule indicates whether the message contains schedule/activity information (kept for backward compat)
	HasSchedule bool `json:"has_schedule"`

	// Action indicates what operation the user wants: create, get, update, or none
	Action ActionType `json:"action"`

	// Schedule contains the extracted schedule data (only if Action is "create")
	Schedule *ScheduleData `json:"schedule,omitempty"`

	// Query contains query/filter parameters (only if Action is "get")
	Query *QueryData `json:"query,omitempty"`

	// Update contains update parameters (only if Action is "update")
	Update *UpdateData `json:"update,omitempty"`

	// MemoriesToSave contains facts/preferences to remember about the user (persistent memory)
	MemoriesToSave []MemoryItem `json:"memories_to_save,omitempty"`
}

// MemoryItem represents a piece of information the AI wants to remember about the user
type MemoryItem struct {
	Type    string `json:"type"`    // "preference", "fact", "habit", "personal"
	Content string `json:"content"` // The actual memory content
}

// ScheduleData represents extracted schedule information from a user message
type ScheduleData struct {
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	ScheduledTime string `json:"scheduled_time"` // ISO 8601 format from AI
	Priority      int    `json:"priority,omitempty"`
}

// QueryData represents query/filter parameters for GET operations
type QueryData struct {
	// FilterType: "today", "tomorrow", "date", "status", "all", "search"
	FilterType string `json:"filter_type"`
	// Date filter (ISO 8601 date, e.g. "2026-02-28")
	Date string `json:"date,omitempty"`
	// Status filter: "pending", "completed", "cancelled", "overdue"
	Status string `json:"status,omitempty"`
	// Keyword for search
	Keyword string `json:"keyword,omitempty"`
}

// UpdateData represents update parameters for UPDATE operations
type UpdateData struct {
	// Identifier to find the activity — could be title (fuzzy) or keyword
	SearchTitle string `json:"search_title"`
	// Fields to update (nil = no change)
	NewTitle         *string `json:"new_title,omitempty"`
	NewDescription   *string `json:"new_description,omitempty"`
	NewScheduledTime *string `json:"new_scheduled_time,omitempty"` // ISO 8601
	NewStatus        *string `json:"new_status,omitempty"`         // "completed", "cancelled"
	NewPriority      *int    `json:"new_priority,omitempty"`
}

// ParseScheduledTime parses the scheduled_time string from AI into time.Time
func (s *ScheduleData) ParseScheduledTime() (*time.Time, error) {
	if s.ScheduledTime == "" {
		return nil, nil
	}

	// Try ISO 8601 format first
	t, err := time.Parse(time.RFC3339, s.ScheduledTime)
	if err == nil {
		return &t, nil
	}

	// Try other common formats
	formats := []string{
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, format := range formats {
		t, err := time.Parse(format, s.ScheduledTime)
		if err == nil {
			return &t, nil
		}
	}

	return nil, err
}
