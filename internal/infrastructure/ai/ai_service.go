package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"smart_alert_system/internal/domain/entity"

	"github.com/google/uuid"
)

type AIService interface {
	ParseIntent(ctx context.Context, message string) (*entity.ParsedIntent, error)
	GenerateHealthRecommendation(ctx context.Context, userID uuid.UUID, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error)
	GenerateMorningAlert(ctx context.Context, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error)
	GenerateEveningSummary(ctx context.Context, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error)
}

type OpenAIService struct {
	apiKey  string
	model   string
	client  *http.Client
	baseURL string
}

func NewOpenAIService(apiKey, model, baseURL string) *OpenAIService {
	// Default to OpenAI if baseURL is empty
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAIService{
		apiKey:  apiKey,
		model:   model,
		client:  &http.Client{Timeout: 120 * time.Second}, // Longer timeout for local Ollama
		baseURL: baseURL,
	}
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func (s *OpenAIService) callAPI(prompt string) (string, error) {
	// Note: API key validation removed - Ollama doesn't need API key
	// For OpenAI, API key should be set but we don't fail here to allow Ollama usage

	url := fmt.Sprintf("%s/chat/completions", s.baseURL)

	reqBody := OpenAIRequest{
		Model: s.model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// Only set Authorization header if API key is provided (Ollama doesn't need it)
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return strings.TrimSpace(openAIResp.Choices[0].Message.Content), nil
}

// callAPIWithSystem calls API with system message (for Ollama)
func (s *OpenAIService) callAPIWithSystem(systemPrompt, userPrompt string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", s.baseURL)

	reqBody := OpenAIRequest{
		Model: s.model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return strings.TrimSpace(openAIResp.Choices[0].Message.Content), nil
}

func (s *OpenAIService) ParseIntent(ctx context.Context, message string) (*entity.ParsedIntent, error) {
	// Use system message for better instruction following
	systemPrompt := `You are a JSON-only response bot. You MUST respond with ONLY valid JSON, no explanations, no markdown, no code blocks, no text before or after.

Your task: Analyze WhatsApp messages and extract intent and entities. Return ONLY a JSON object.

Valid intents: "add_activity", "delete_activity", "update_activity", "list_activities", "question", "greeting", "unknown"

JSON format:
{
  "intent": "one of the valid intents",
  "confidence": 0.0 to 1.0,
  "entities": {
    "title": "activity title if exists",
    "description": "description if exists",
    "scheduled_time": "time in natural language",
    "activity_id": "id if for update/delete",
    "priority": 1-5 if mentioned
  }
}

Examples:
Input: "Hi"
Output: {"intent":"greeting","confidence":0.9,"entities":{}}

Input: "Saya mau olahraga besok jam 6 pagi"
Output: {"intent":"add_activity","confidence":0.9,"entities":{"title":"olahraga","scheduled_time":"besok jam 6 pagi"}}

Input: "Lihat kegiatan hari ini"
Output: {"intent":"list_activities","confidence":0.9,"entities":{}}

REMEMBER: Return ONLY JSON, nothing else. Start with { and end with }.`

	userPrompt := fmt.Sprintf(`Analyze this WhatsApp message and return JSON:

Message: "%s"`, message)

	// Use system message if baseURL suggests Ollama (local)
	useSystemMessage := strings.Contains(s.baseURL, "localhost") || strings.Contains(s.baseURL, "127.0.0.1")

	var response string
	var err error

	if useSystemMessage {
		// For Ollama, use system message
		response, err = s.callAPIWithSystem(systemPrompt, userPrompt)
	} else {
		// For OpenAI, use combined prompt
		combinedPrompt := fmt.Sprintf(`%s

%s`, systemPrompt, userPrompt)
		response, err = s.callAPI(combinedPrompt)
	}

	if err != nil {
		return &entity.ParsedIntent{
			Type:       entity.IntentUnknown,
			Confidence: 0.0,
			Entities:   make(map[string]interface{}),
		}, err
	}

	// Clean response - remove markdown code blocks, whitespace, etc.
	cleanedResponse := cleanJSONResponse(response)

	// Parse JSON response
	var result struct {
		Intent     string                 `json:"intent"`
		Confidence float64                `json:"confidence"`
		Entities   map[string]interface{} `json:"entities"`
	}

	if err := json.Unmarshal([]byte(cleanedResponse), &result); err != nil {
		// Try to extract intent and entities from text response
		log.Printf("⚠️  Failed to parse JSON response. Trying to extract from text...")
		if extractedIntent := extractIntentFromText(response, message); extractedIntent != nil {
			log.Printf("✓ Successfully extracted intent from text: %s", extractedIntent.Type)
			return extractedIntent, nil
		}

		// If extraction fails, return error to trigger fallback parser
		log.Printf("⚠️  Could not extract intent from text. Will use fallback parser.")
		return nil, fmt.Errorf("ai_parse_failed")
	}

	// Validate and clean entities from JSON response
	// For add_activity, always validate entities match the original message
	if result.Intent == "add_activity" {
		originalMessageLower := strings.ToLower(message)

		// Validate title - must exist in original message
		if title, ok := result.Entities["title"].(string); ok && title != "" {
			titleLower := strings.ToLower(title)
			// Check if title actually appears in original message
			if !strings.Contains(originalMessageLower, titleLower) {
				// Title from AI doesn't match, extract from message instead
				extractedTitle := extractTitleFromMessage(message)
				if extractedTitle != "" {
					result.Entities["title"] = extractedTitle
					log.Printf("⚠️  AI title '%s' doesn't match message, using extracted: '%s'", title, extractedTitle)
				} else {
					delete(result.Entities, "title")
					log.Printf("⚠️  AI title '%s' doesn't match message and extraction failed, removing", title)
				}
			}
		} else {
			// No title from AI, try to extract from message
			extractedTitle := extractTitleFromMessage(message)
			if extractedTitle != "" {
				result.Entities["title"] = extractedTitle
				log.Printf("✓ Extracted title from message: '%s'", extractedTitle)
			}
		}

		// Validate scheduled_time - must exist in original message
		if timeStr, ok := result.Entities["scheduled_time"].(string); ok && timeStr != "" {
			timeStrLower := strings.ToLower(timeStr)
			// Check if time actually appears in original message
			if !strings.Contains(originalMessageLower, timeStrLower) {
				// Time from AI doesn't match, extract from message instead
				timePatterns := []*regexp.Regexp{
					regexp.MustCompile(`(?i)(jam\s+)?(\d{1,2})\s+(pagi|siang|sore|malam)`),
					regexp.MustCompile(`(?i)(\d{1,2}):(\d{2})`),
					regexp.MustCompile(`(?i)jam\s+(\d{1,2})`),
				}

				var timeMatch string
				for _, pattern := range timePatterns {
					matches := pattern.FindStringSubmatch(message)
					if len(matches) > 0 {
						timeMatch = matches[0]
						break
					}
				}

				// Add day keywords if present
				dayKeywords := []string{"besok", "lusa", "hari ini"}
				var dayMatch string
				for _, day := range dayKeywords {
					if strings.Contains(originalMessageLower, day) {
						dayMatch = day
						break
					}
				}

				// Combine day and time
				if dayMatch != "" && timeMatch != "" {
					result.Entities["scheduled_time"] = dayMatch + " " + timeMatch
				} else if timeMatch != "" {
					result.Entities["scheduled_time"] = timeMatch
				} else if dayMatch != "" {
					result.Entities["scheduled_time"] = dayMatch
				} else {
					delete(result.Entities, "scheduled_time")
				}

				if timeMatch != "" || dayMatch != "" {
					log.Printf("⚠️  AI time '%s' doesn't match message, using extracted: '%s'", timeStr, result.Entities["scheduled_time"])
				}
			}
		} else {
			// No time from AI, try to extract from message
			timePatterns := []*regexp.Regexp{
				regexp.MustCompile(`(?i)(jam\s+)?(\d{1,2})\s+(pagi|siang|sore|malam)`),
				regexp.MustCompile(`(?i)(\d{1,2}):(\d{2})`),
				regexp.MustCompile(`(?i)jam\s+(\d{1,2})`),
			}

			var timeMatch string
			for _, pattern := range timePatterns {
				matches := pattern.FindStringSubmatch(message)
				if len(matches) > 0 {
					timeMatch = matches[0]
					break
				}
			}

			dayKeywords := []string{"besok", "lusa", "hari ini"}
			var dayMatch string
			for _, day := range dayKeywords {
				if strings.Contains(originalMessageLower, day) {
					dayMatch = day
					break
				}
			}

			if dayMatch != "" && timeMatch != "" {
				result.Entities["scheduled_time"] = dayMatch + " " + timeMatch
			} else if timeMatch != "" {
				result.Entities["scheduled_time"] = timeMatch
			} else if dayMatch != "" {
				result.Entities["scheduled_time"] = dayMatch
			}

			if timeMatch != "" || dayMatch != "" {
				log.Printf("✓ Extracted time from message: '%s'", result.Entities["scheduled_time"])
			}
		}

		// Remove invalid entities for new activities
		delete(result.Entities, "activity_id")
		// Priority is optional, but if provided and not in message, remove it
		if priority, ok := result.Entities["priority"].(float64); ok {
			// Keep priority only if it's reasonable (1-5)
			if priority < 1 || priority > 5 {
				delete(result.Entities, "priority")
			}
		}
	}

	return &entity.ParsedIntent{
		Type:       entity.IntentType(result.Intent),
		Confidence: result.Confidence,
		Entities:   result.Entities,
	}, nil
}

func (s *OpenAIService) GenerateHealthRecommendation(ctx context.Context, userID uuid.UUID, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error) {
	activitiesStr := formatActivitiesForAI(activities)

	prompt := fmt.Sprintf(`Based on the user's activities and health profile, generate a personalized health recommendation.

Activities:
%s

Health Profile: %s

Generate a concise, helpful health recommendation in Indonesian.`, activitiesStr, formatHealthProfileForAI(healthProfile))

	return s.callAPI(prompt)
}

func (s *OpenAIService) GenerateMorningAlert(ctx context.Context, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error) {
	activitiesStr := formatActivitiesForAI(activities)

	prompt := fmt.Sprintf(`Generate a friendly morning alert message in Indonesian that:
1. Lists today's scheduled activities
2. Provides personalized health tips based on the activities

Activities today:
%s

Health Profile: %s

Make it warm, encouraging, and concise.`, activitiesStr, formatHealthProfileForAI(healthProfile))

	return s.callAPI(prompt)
}

func (s *OpenAIService) GenerateEveningSummary(ctx context.Context, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error) {
	activitiesStr := formatActivitiesForAI(activities)

	prompt := fmt.Sprintf(`Generate an evening summary message in Indonesian that:
1. Summarizes completed activities today
2. Analyzes activity patterns
3. Provides recommendations for tomorrow

Completed activities:
%s

Health Profile: %s

Make it reflective, encouraging, and actionable.`, activitiesStr, formatHealthProfileForAI(healthProfile))

	return s.callAPI(prompt)
}

func formatActivitiesForAI(activities []*entity.Activity) string {
	if len(activities) == 0 {
		return "No activities"
	}

	var sb strings.Builder
	for i, activity := range activities {
		sb.WriteString(fmt.Sprintf("%d. %s - %s (Status: %s)\n", i+1, activity.Title, activity.Description, activity.Status))
	}
	return sb.String()
}

func formatHealthProfileForAI(profile *entity.UserHealthProfile) string {
	if profile == nil {
		return "No health profile available"
	}
	return fmt.Sprintf("Age: %v, Gender: %s", profile.Age, profile.Gender)
}

// cleanJSONResponse extracts JSON from response, handling markdown code blocks and extra text
func cleanJSONResponse(response string) string {
	response = strings.TrimSpace(response)

	// Remove markdown code blocks (```json ... ``` or ``` ... ```)
	if strings.HasPrefix(response, "```") {
		// Find the first newline after ```
		firstNewline := strings.Index(response, "\n")
		if firstNewline > 0 {
			response = response[firstNewline+1:]
		}
		// Remove trailing ```
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSuffix(response, "`")
	}

	// Find JSON object boundaries
	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")

	if startIdx >= 0 && endIdx > startIdx {
		response = response[startIdx : endIdx+1]
	}

	// Remove any remaining whitespace
	response = strings.TrimSpace(response)

	return response
}

// extractIntentFromText tries to extract intent and entities from AI's text response
// This handles cases where AI returns explanations instead of JSON
func extractIntentFromText(aiResponse, originalMessage string) *entity.ParsedIntent {
	aiResponseLower := strings.ToLower(aiResponse)
	originalMessageLower := strings.ToLower(originalMessage)

	entities := make(map[string]interface{})
	var intentType entity.IntentType
	confidence := 0.6 // Lower confidence for extracted intents

	// Extract intent from text
	intentPatterns := map[entity.IntentType][]string{
		entity.IntentAddActivity:    {"add_activity", "adding activity", "intent is \"add_activity\"", "intent: \"add_activity\""},
		entity.IntentDeleteActivity: {"delete_activity", "deleting activity", "intent is \"delete_activity\""},
		entity.IntentUpdateActivity: {"update_activity", "updating activity", "intent is \"update_activity\""},
		entity.IntentListActivities: {"list_activities", "listing activities", "intent is \"list_activities\""},
		entity.IntentGreeting:       {"greeting", "intent is \"greeting\""},
		entity.IntentQuestion:       {"question", "intent is \"question\""},
	}

	for intent, patterns := range intentPatterns {
		for _, pattern := range patterns {
			if strings.Contains(aiResponseLower, pattern) {
				intentType = intent
				confidence = 0.7
				break
			}
		}
		if intentType != "" {
			break
		}
	}

	// If no intent found, try to infer from original message
	if intentType == "" || intentType == entity.IntentUnknown {
		// Check if message contains activity-related keywords
		activityKeywords := []string{"jam", "pagi", "siang", "sore", "malam", "besok", "lusa", "hari ini", "mau", "akan", "ingin"}
		hasActivityKeyword := false
		for _, keyword := range activityKeywords {
			if strings.Contains(originalMessageLower, keyword) {
				hasActivityKeyword = true
				break
			}
		}

		if hasActivityKeyword {
			intentType = entity.IntentAddActivity
			confidence = 0.6
		}
	}

	// Extract entities from text or original message
	if intentType == entity.IntentAddActivity {
		// ALWAYS extract from original message first (more reliable)
		// Extract title from original message (remove common words)
		title := extractTitleFromMessage(originalMessage)
		if title != "" {
			entities["title"] = title
		}

		// Extract scheduled_time from original message
		timePatterns := []*regexp.Regexp{
			regexp.MustCompile(`(?i)(jam\s+)?(\d{1,2})\s+(pagi|siang|sore|malam)`),
			regexp.MustCompile(`(?i)(\d{1,2}):(\d{2})`),
			regexp.MustCompile(`(?i)jam\s+(\d{1,2})`),
		}

		var timeMatch string
		for _, pattern := range timePatterns {
			matches := pattern.FindStringSubmatch(originalMessage)
			if len(matches) > 0 {
				timeMatch = matches[0]
				break
			}
		}

		// Add day keywords if present
		dayKeywords := []string{"besok", "lusa", "hari ini"}
		var dayMatch string
		for _, day := range dayKeywords {
			if strings.Contains(originalMessageLower, day) {
				dayMatch = day
				break
			}
		}

		// Combine day and time
		if dayMatch != "" && timeMatch != "" {
			entities["scheduled_time"] = dayMatch + " " + timeMatch
		} else if timeMatch != "" {
			entities["scheduled_time"] = timeMatch
		} else if dayMatch != "" {
			entities["scheduled_time"] = dayMatch
		}

		// Remove invalid entities that don't match original message
		// Don't trust AI response entities, only use what we extracted
		if title == "" {
			// If we can't extract title, don't add it
			delete(entities, "title")
		}
		// Remove activity_id and priority from AI response (they're not in the message)
		delete(entities, "activity_id")
		delete(entities, "priority")
	}

	// If we found an intent, return it
	if intentType != "" && intentType != entity.IntentUnknown {
		return &entity.ParsedIntent{
			Type:       intentType,
			Confidence: confidence,
			Entities:   entities,
		}
	}

	return nil
}

// extractTitleFromMessage extracts activity title from message
func extractTitleFromMessage(message string) string {
	messageLower := strings.ToLower(message)

	// Remove common prefixes
	prefixes := []string{"saya mau", "saya akan", "saya ingin", "aku mau", "aku akan", "aku ingin"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(messageLower, prefix) {
			messageLower = strings.TrimPrefix(messageLower, prefix)
			messageLower = strings.TrimSpace(messageLower)
			break
		}
	}

	// Remove common time-related words and patterns
	timePatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\s*jam\s+\d+.*`),
		regexp.MustCompile(`(?i)\s*\d+\s+(pagi|siang|sore|malam).*`),
		regexp.MustCompile(`(?i)\s*\d+:\d+.*`),
		regexp.MustCompile(`(?i)\s*(besok|lusa|hari ini).*`),
	}

	for _, pattern := range timePatterns {
		messageLower = pattern.ReplaceAllString(messageLower, "")
	}

	// Remove individual time words
	timeWords := []string{
		"jam", "pagi", "siang", "sore", "malam",
		"besok", "lusa", "hari ini",
		"mau", "akan", "ingin", "tambah", "add", "saya", "aku",
	}

	words := strings.Fields(messageLower)
	var titleWords []string

	for _, word := range words {
		isTimeWord := false
		// Check if word is a time word
		for _, timeWord := range timeWords {
			if word == timeWord {
				isTimeWord = true
				break
			}
		}
		// Skip numbers (likely time)
		if matched, _ := regexp.MatchString(`^\d+$`, word); matched {
			isTimeWord = true
		}
		// Skip empty strings
		if word == "" {
			isTimeWord = true
		}
		if !isTimeWord {
			titleWords = append(titleWords, word)
		}
	}

	title := strings.Join(titleWords, " ")
	title = strings.TrimSpace(title)

	return title
}

// extractTitleFromAIText tries to extract title from AI response text
func extractTitleFromAIText(aiResponse, originalMessage string) string {
	aiResponseLower := strings.ToLower(aiResponse)

	// Look for activity type mentions in AI response
	activityPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)activity type[:\s]+["']?([^"',\n]+)["']?`),
		regexp.MustCompile(`(?i)activity[:\s]+["']?([^"',\n]+)["']?`),
		regexp.MustCompile(`(?i)title[:\s]+["']?([^"',\n]+)["']?`),
		regexp.MustCompile(`(?i)(futsal|soccer|olahraga|meeting|makan|minum|tidur|bangun|kerja|belajar)`),
	}

	for _, pattern := range activityPatterns {
		matches := pattern.FindStringSubmatch(aiResponseLower)
		if len(matches) > 1 && matches[1] != "" {
			title := strings.TrimSpace(matches[1])
			// Filter out common words
			if title != "add_activity" && title != "activity" && len(title) > 2 {
				return title
			}
		}
	}

	// Fallback to extracting from original message
	return extractTitleFromMessage(originalMessage)
}
