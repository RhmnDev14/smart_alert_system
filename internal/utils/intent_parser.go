package utils

import (
	"regexp"
	"strings"
	"time"

	"smart_alert_system/internal/domain/entity"
)

// FallbackIntentParser provides rule-based intent detection when AI fails
func FallbackIntentParser(message string, baseTime time.Time) *entity.ParsedIntent {
	message = strings.ToLower(strings.TrimSpace(message))
	
	// Check for greeting
	if isGreeting(message) {
		return &entity.ParsedIntent{
			Type:       entity.IntentGreeting,
			Confidence: 0.8,
			Entities:   make(map[string]interface{}),
		}
	}
	
	// Check for list activities
	if isListActivities(message) {
		return &entity.ParsedIntent{
			Type:       entity.IntentListActivities,
			Confidence: 0.8,
			Entities:   make(map[string]interface{}),
		}
	}
	
	// Check for add activity (most common case)
	if intent := detectAddActivity(message, baseTime); intent != nil {
		return intent
	}
	
	// Default to question/unknown
	return &entity.ParsedIntent{
		Type:       entity.IntentQuestion,
		Confidence: 0.5,
		Entities:   make(map[string]interface{}),
	}
}

func isGreeting(message string) bool {
	greetings := []string{
		"halo", "hai", "hi", "hello", "hey",
		"selamat pagi", "selamat siang", "selamat sore", "selamat malam",
		"pagi", "siang", "sore", "malam",
		"apa kabar", "kabar", "gimana",
	}
	
	for _, greeting := range greetings {
		if strings.Contains(message, greeting) {
			return true
		}
	}
	return false
}

func isListActivities(message string) bool {
	patterns := []string{
		"lihat kegiatan", "list kegiatan", "daftar kegiatan",
		"kegiatan hari ini", "agenda hari ini", "jadwal hari ini",
		"apa kegiatan", "kegiatan apa", "agenda apa",
		"show activities", "list activities",
	}
	
	messageLower := strings.ToLower(message)
	for _, pattern := range patterns {
		if strings.Contains(messageLower, pattern) {
			return true
		}
	}
	return false
}

func detectAddActivity(message string, baseTime time.Time) *entity.ParsedIntent {
	// Patterns that suggest adding activity
	activityKeywords := []string{
		"jam", "pagi", "siang", "sore", "malam",
		"besok", "lusa", "hari ini", "hari ini",
		"mau", "akan", "ingin", "rencana", "agenda",
		"tambah", "add", "buat", "jadwalkan",
	}
	
	hasTimeKeyword := false
	for _, keyword := range activityKeywords {
		if strings.Contains(message, keyword) {
			hasTimeKeyword = true
			break
		}
	}
	
	if !hasTimeKeyword {
		return nil
	}
	
	entities := make(map[string]interface{})
	
	// Extract time information
	timePatterns := []struct {
		pattern *regexp.Regexp
		format  string
	}{
		{regexp.MustCompile(`jam\s+(\d{1,2})\s+(pagi|siang|sore|malam)`), "jam %d %s"},
		{regexp.MustCompile(`(\d{1,2})\s+(pagi|siang|sore|malam)`), "%d %s"},
		{regexp.MustCompile(`jam\s+(\d{1,2})`), "jam %d"},
		{regexp.MustCompile(`(\d{1,2}):(\d{2})`), "%d:%s"},
	}
	
	var scheduledTime string
	for _, tp := range timePatterns {
		matches := tp.pattern.FindStringSubmatch(message)
		if len(matches) > 0 {
			scheduledTime = matches[0]
			break
		}
	}
	
	// Add day keywords
	dayKeywords := []string{"besok", "lusa", "hari ini", "hari ini"}
	for _, day := range dayKeywords {
		if strings.Contains(message, day) {
			if scheduledTime != "" {
				scheduledTime = day + " " + scheduledTime
			} else {
				scheduledTime = day
			}
			break
		}
	}
	
	if scheduledTime != "" {
		entities["scheduled_time"] = scheduledTime
	}
	
	// Extract activity title (everything before time keywords)
	title := extractActivityTitle(message)
	if title != "" {
		entities["title"] = title
	}
	
	// If we found time or title, it's likely an add_activity intent
	if scheduledTime != "" || title != "" {
		return &entity.ParsedIntent{
			Type:       entity.IntentAddActivity,
			Confidence: 0.7,
			Entities:   entities,
		}
	}
	
	return nil
}

func extractActivityTitle(message string) string {
	// Remove common time-related words to extract activity
	timeWords := []string{
		"jam", "pagi", "siang", "sore", "malam",
		"besok", "lusa", "hari ini",
		"mau", "akan", "ingin", "tambah", "add",
	}
	
	words := strings.Fields(message)
	var titleWords []string
	
	for _, word := range words {
		isTimeWord := false
		for _, timeWord := range timeWords {
			if strings.Contains(strings.ToLower(word), timeWord) {
				isTimeWord = true
				break
			}
		}
		if !isTimeWord {
			titleWords = append(titleWords, word)
		}
	}
	
	title := strings.Join(titleWords, " ")
	title = strings.TrimSpace(title)
	
	// Remove common prefixes
	prefixes := []string{"saya", "aku", "saya mau", "saya akan", "saya ingin"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(title), prefix) {
			title = strings.TrimPrefix(strings.ToLower(title), prefix)
			title = strings.TrimSpace(title)
			break
		}
	}
	
	return title
}

