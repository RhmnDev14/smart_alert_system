package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type TelegramClient struct {
	botToken string
	baseURL  string
	client   *http.Client
}

type SendMessageRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type SendMessageResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
}

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      *User  `json:"from"`
	Chat      *Chat  `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
}

type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type GetUpdatesRequest struct {
	Offset  int `json:"offset,omitempty"`
	Limit   int `json:"limit,omitempty"`
	Timeout int `json:"timeout,omitempty"`
}

type GetUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func NewTelegramClient(botToken string) *TelegramClient {
	return &TelegramClient{
		botToken: botToken,
		baseURL:  fmt.Sprintf("https://api.telegram.org/bot%s", botToken),
		client: &http.Client{
			Timeout: 60 * time.Second, // Long polling timeout
		},
	}
}

// SendMessage sends a text message to a specific chat ID
func (c *TelegramClient) SendMessage(chatID int64, message string) error {
	url := fmt.Sprintf("%s/sendMessage", c.baseURL)

	reqBody := SendMessageRequest{
		ChatID: chatID,
		Text:   message,
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

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var sendResp SendMessageResponse
	if err := json.Unmarshal(body, &sendResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !sendResp.OK {
		return fmt.Errorf("telegram API returned not OK: %s", string(body))
	}

	return nil
}

// SendMessageByStringID sends a message using a string chat ID (for compatibility)
// Telegram chat IDs are int64, but we store them as strings in the DB
func (c *TelegramClient) SendMessageByStringID(chatIDStr string, message string) error {
	var chatID int64
	_, err := fmt.Sscanf(chatIDStr, "%d", &chatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID '%s': %w", chatIDStr, err)
	}
	return c.SendMessage(chatID, message)
}

// GetUpdates fetches new updates from Telegram using long polling
func (c *TelegramClient) GetUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=30", c.baseURL, offset)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var updatesResp GetUpdatesResponse
	if err := json.Unmarshal(body, &updatesResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !updatesResp.OK {
		return nil, fmt.Errorf("telegram API returned not OK: %s", string(body))
	}

	return updatesResp.Result, nil
}

// GetMe verifies the bot token and returns bot info
func (c *TelegramClient) GetMe() (*User, error) {
	url := fmt.Sprintf("%s/getMe", c.baseURL)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call getMe: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		OK     bool `json:"ok"`
		Result User `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !result.OK {
		return nil, fmt.Errorf("telegram API returned not OK: %s", string(body))
	}

	log.Printf("✓ Telegram Bot: @%s (%s)", result.Result.Username, result.Result.FirstName)
	return &result.Result, nil
}
