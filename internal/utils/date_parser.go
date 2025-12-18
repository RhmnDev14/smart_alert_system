package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseTimeFromText parses time from natural language text in Indonesian
// Examples: "besok jam 6 pagi", "hari ini jam 2 siang", "lusa jam 8 malam"
func ParseTimeFromText(text string, baseTime time.Time) (*time.Time, error) {
	if text == "" {
		return nil, nil
	}

	text = strings.ToLower(strings.TrimSpace(text))
	loc := baseTime.Location()

	// Parse relative dates
	var targetDate time.Time
	if strings.Contains(text, "besok") || strings.Contains(text, "tomorrow") {
		targetDate = baseTime.AddDate(0, 0, 1)
	} else if strings.Contains(text, "lusa") || strings.Contains(text, "day after tomorrow") {
		targetDate = baseTime.AddDate(0, 0, 2)
	} else if strings.Contains(text, "hari ini") || strings.Contains(text, "today") {
		targetDate = baseTime
	} else {
		// Default to today if no date specified
		targetDate = baseTime
	}

	// Extract hour and minute
	hour := 9 // Default 9 AM
	minute := 0

	// Pattern untuk jam
	hourPattern := regexp.MustCompile(`(\d{1,2})\s*(pagi|siang|sore|malam|am|pm|:)?`)
	matches := hourPattern.FindStringSubmatch(text)

	if len(matches) >= 2 {
		if h, err := strconv.Atoi(matches[1]); err == nil {
			hour = h
		}

		// Adjust for AM/PM or Indonesian time indicators
		if len(matches) >= 3 {
			indicator := strings.ToLower(matches[2])
			if strings.Contains(indicator, "siang") || strings.Contains(indicator, "sore") {
				if hour < 12 {
					hour += 12
				}
			} else if strings.Contains(indicator, "malam") {
				if hour < 12 {
					hour += 12
				}
			} else if strings.Contains(indicator, "pm") {
				if hour < 12 {
					hour += 12
				}
			} else if strings.Contains(indicator, "am") && hour == 12 {
				hour = 0
			}
		}

		// Check for minute pattern (HH:MM)
		minutePattern := regexp.MustCompile(`(\d{1,2}):(\d{2})`)
		minMatches := minutePattern.FindStringSubmatch(text)
		if len(minMatches) >= 3 {
			if h, err := strconv.Atoi(minMatches[1]); err == nil {
				hour = h
			}
			if m, err := strconv.Atoi(minMatches[2]); err == nil {
				minute = m
			}
		}
	}

	// Create time
	result := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), hour, minute, 0, 0, loc)
	
	// If the time is in the past and it's today, assume it's for tomorrow
	if result.Before(baseTime) && targetDate.Equal(baseTime) {
		result = result.AddDate(0, 0, 1)
	}

	return &result, nil
}

// ParseISO8601Time parses ISO 8601 format time string
func ParseISO8601Time(timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("unable to parse time: %s", timeStr)
}

