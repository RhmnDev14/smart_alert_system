package entity

import "time"

// AIProcessResult represents the result from AI processing a user message.
// The AI acts as a gateway: it responds to the user AND detects schedules.
type AIProcessResult struct {
	// Response is the natural language response to send back to the user
	Response string `json:"response"`

	// HasSchedule indicates whether the message contains schedule/activity information
	HasSchedule bool `json:"has_schedule"`

	// Schedule contains the extracted schedule data (only if HasSchedule is true)
	Schedule *ScheduleData `json:"schedule,omitempty"`
}

// ScheduleData represents extracted schedule information from a user message
type ScheduleData struct {
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	ScheduledTime string `json:"scheduled_time"` // ISO 8601 format from AI
	Priority      int    `json:"priority,omitempty"`
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
