package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type WahaClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type WahaMessage struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Body      string `json:"body"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
}

type SendMessageRequest struct {
	Session string `json:"session,omitempty"` // Optional session name
	ChatID  string `json:"chatId"`
	Text    string `json:"text"`
}

type SendMessageResponse struct {
	Sent bool   `json:"sent"`
	ID   string `json:"id"`
}

func NewWahaClient(baseURL, apiKey string) *WahaClient {
	return &WahaClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *WahaClient) SendMessage(chatID, message string) error {
	// Remove trailing slash from baseURL
	baseURL := strings.TrimSuffix(c.baseURL, "/")

	// Format chatID - ensure it has proper format
	formattedChatID := chatID
	if !strings.Contains(chatID, "@") {
		// If no @, assume it's a phone number, add @c.us
		formattedChatID = chatID + "@c.us"
	}

	// Use the correct format based on Waha API documentation
	// Format: POST /api/sendText with {"session":"default","chatId":"...","text":"..."}
	url := fmt.Sprintf("%s/api/sendText", baseURL)

	reqBody := SendMessageRequest{
		Session: "default", // Session name is required
		ChatID:  formattedChatID,
		Text:    message,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-Api-Key", c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Accept both 200 OK and 201 Created as success
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	// Log success (optional, can be removed in production)
	if resp.StatusCode == http.StatusCreated {
		// 201 Created is also valid for POST requests
	}

	return nil
}

func (c *WahaClient) GetWebhookURL() string {
	return fmt.Sprintf("%s/api/webhook", c.baseURL)
}
