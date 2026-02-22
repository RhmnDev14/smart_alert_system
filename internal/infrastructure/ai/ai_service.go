package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"smart_alert_system/internal/domain/entity"

	"github.com/google/uuid"
)

// AIService defines the interface for AI operations
type AIService interface {
	// ProcessMessage processes a user message and returns a response + optional schedule data
	ProcessMessage(ctx context.Context, message string, currentTime time.Time) (*entity.AIProcessResult, error)
	// GenerateMorningAlert generates morning alert content
	GenerateMorningAlert(ctx context.Context, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error)
	// GenerateEveningSummary generates evening summary content
	GenerateEveningSummary(ctx context.Context, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error)
	// GenerateActivityReminder generates dynamic engaging reminder for a specific activity
	GenerateActivityReminder(ctx context.Context, title, description, timeStr string) (string, error)
}

type OpenAIService struct {
	apiKey  string
	model   string
	client  *http.Client
	baseURL string
}

func NewOpenAIService(apiKey, model, baseURL string) *OpenAIService {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAIService{
		apiKey:  apiKey,
		model:   model,
		client:  &http.Client{Timeout: 120 * time.Second},
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

func (s *OpenAIService) callAPIWithSystem(systemPrompt, userPrompt string) (string, error) {
	url := s.baseURL
	if !strings.HasSuffix(url, "/chat/completions") {
		url = fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(s.baseURL, "/"))
	}

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
	// Important for OpenRouter routing
	req.Header.Set("HTTP-Referer", "https://smart-alert-system.local")
	req.Header.Set("X-Title", "Smart Alert System")

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

func (s *OpenAIService) callAPI(prompt string) (string, error) {
	url := s.baseURL
	if !strings.HasSuffix(url, "/chat/completions") {
		url = fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(s.baseURL, "/"))
	}

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
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	// Important for OpenRouter routing
	req.Header.Set("HTTP-Referer", "https://smart-alert-system.local")
	req.Header.Set("X-Title", "Smart Alert System")

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

// ProcessMessage is the main gateway function.
// AI will:
// 1. Respond naturally to the user
// 2. Detect if the message contains schedule/activity information
// 3. If yes, extract the schedule data for saving to DB
func (s *OpenAIService) ProcessMessage(ctx context.Context, message string, currentTime time.Time) (*entity.AIProcessResult, error) {
	systemPrompt := fmt.Sprintf(`Kamu adalah asisten pintar bernama "Smart Alert Bot". Tugasmu ada 2:

1. **Merespons pesan user secara natural** — bisa ngobrol biasa, menjawab pertanyaan, dll.
2. **Mendeteksi jadwal kegiatan** — jika pesan user mengandung informasi tentang jadwal/kegiatan/rencana yang perlu diingat, kamu harus mengekstrak data jadwal tersebut.

WAKTU SAAT INI: %s

ATURAN PENTING:
- Selalu balas dalam Bahasa Indonesia
- Jika user menyebutkan jadwal/kegiatan/rencana (contoh: "besok meeting jam 10", "hari ini mau olahraga jam 6 pagi", "minggu depan ada acara", dll), set has_schedule = true dan isi data schedule
- Jika pesan hanya obrolan biasa, sapaan, atau pertanyaan umum, set has_schedule = false
- Untuk scheduled_time, SELALU gunakan format ISO 8601 (contoh: "2026-02-23T10:00:00+07:00"), hitung berdasarkan waktu saat ini
- Jika user bilang "besok", hitung dari tanggal saat ini + 1 hari
- Jika user bilang "hari ini", gunakan tanggal saat ini
- Jika waktu tidak spesifik, gunakan estimasi yang masuk akal (pagi=08:00, siang=12:00, sore=16:00, malam=20:00)
- Response harus friendly dan mengkonfirmasi jika ada jadwal yang disimpan

KAMU HARUS MEMBALAS DENGAN FORMAT JSON BERIKUT (TANPA markdown code block, HANYA JSON murni):
{
  "response": "pesan balasan untuk user dalam bahasa Indonesia",
  "has_schedule": true/false,
  "schedule": {
    "title": "judul kegiatan",
    "description": "deskripsi singkat (opsional, boleh kosong)",
    "scheduled_time": "2026-02-23T10:00:00+07:00",
    "priority": 3
  }
}

Jika has_schedule = false, field "schedule" boleh null atau tidak ada.
Priority: 1(rendah) - 5(tinggi), default 3.

INGAT: Balas HANYA JSON, tanpa backtick, tanpa markdown, tanpa penjelasan tambahan. Mulai dengan { dan akhiri dengan }.`, currentTime.Format("2006-01-02T15:04:05-07:00"))

	userPrompt := fmt.Sprintf(`Pesan user: "%s"`, message)

	response, err := s.callAPIWithSystem(systemPrompt, userPrompt)
	if err != nil {
		log.Printf("❌ AI API call failed: %v", err)
		return &entity.AIProcessResult{
			Response:    "Maaf, saya sedang mengalami gangguan. Silakan coba lagi nanti.",
			HasSchedule: false,
		}, err
	}

	// Clean and parse the response
	cleanedResponse := cleanJSONResponse(response)
	log.Printf("🤖 AI Raw Response: %s", cleanedResponse)

	var result entity.AIProcessResult
	if err := json.Unmarshal([]byte(cleanedResponse), &result); err != nil {
		log.Printf("⚠️ Failed to parse AI JSON response: %v", err)
		log.Printf("⚠️ Raw response: %s", response)
		// If JSON parsing fails, use the raw response as the message
		return &entity.AIProcessResult{
			Response:    response,
			HasSchedule: false,
		}, nil
	}

	return &result, nil
}

func (s *OpenAIService) GenerateMorningAlert(ctx context.Context, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error) {
	activitiesStr := formatActivitiesForAI(activities)

	prompt := fmt.Sprintf(`Buatkan pesan pengingat pagi hari (Morning Alert) dalam bahasa Indonesia dengan kriteria berikut:
1. Sapa dengan penuh semangat dan berikan ucapan motivasi atau kutipan penyemangat untuk memulai hari.
2. Jika ada kegiatan terjadwal hari ini, sebutkan dan ingatkan kegiatan-kegiatan tersebut.
3. Tanyakan kepada user apa saja rencana atau kegiatan tambahan yang akan dilakukan hari ini, agar user bisa menambahkannya.
4. Jangan lupa berikan tips kesehatan personal yang relevan jika ada (opsional).
5. DI AKHIR PESAN, SELALU TAMBAHKAN TEKS PERSIS SEPERTI INI:
Smart Alert System
Develop by Rahman Umardi

Jadwal kegiatan hari ini:
%s

Profil Kesehatan User (jadikan referensi jika perlu): %s

Buat pesan yang ringkas, hangat, bersahabat, dan memotivasi.`, activitiesStr, formatHealthProfileForAI(healthProfile))

	return s.callAPI(prompt)
}

func (s *OpenAIService) GenerateEveningSummary(ctx context.Context, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error) {
	activitiesStr := formatActivitiesForAI(activities)

	prompt := fmt.Sprintf(`Buatkan pesan ringkasan malam hari (Evening Summary) dalam bahasa Indonesia dengan kriteria berikut:
1. Sapa user dengan ramah dan apresiasi usaha mereka hari ini.
2. Berikan ringkasan kegiatan yang telah dilakukan (completed) hari ini.
3. Berikan analisis singkat terkait pola kegiatan mereka hari ini (misalnya terlalu sibuk, cukup seimbang, atau kurang aktivitas fisik).
4. Berikan rekomendasi kesehatan yang relevan untuk dipersiapkan menghadap hari esok.
5. DI AKHIR PESAN, SELALU TAMBAHKAN TEKS PERSIS SEPERTI INI (jangan gunakan format markdown bold):
Smart Alert System
Develop by Rahman Umardi

Kegiatan yang telah Selesai hari ini:
%s

Profil Kesehatan User (jadikan referensi jika perlu): %s

Buat pesan yang reflektif, mendukung, dan memotivasi untuk beristirahat.`, activitiesStr, formatHealthProfileForAI(healthProfile))

	return s.callAPI(prompt)
}

func (s *OpenAIService) GenerateActivityReminder(ctx context.Context, title, description, timeStr string) (string, error) {
	prompt := fmt.Sprintf(`Buatkan pesan pengingat (Reminder) singkat, natural, asyik, dan tidak kaku dalam bahasa Indonesia untuk kegiatan berikut:
Judul Kegiatan: %s
Jam: %s
Catatan/Deskripsi: %s

Kriteria:
1. Sapa user dengan ceria layaknya asisten pribadi yang ramah, dan ingatkan bahwa sudah saatnya untuk kegiatan tersebut.
2. Gunakan emoji yang relevan dengan kegiatannya.
3. Pesannya harus ringkas (cukup 1-2 paragraf pendek atau kalimat singkat yang memotivasi).
`, title, timeStr, description)

	return s.callAPI(prompt)
}

func formatActivitiesForAI(activities []*entity.Activity) string {
	if len(activities) == 0 {
		return "No activities"
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	var sb strings.Builder
	for i, activity := range activities {
		localTime := activity.ScheduledTime.In(loc)
		descStr := ""
		if activity.Description != "" {
			descStr = " - " + activity.Description
		}
		sb.WriteString(fmt.Sprintf("%d. %s%s (Waktu: %s)\n",
			i+1, activity.Title, descStr, localTime.Format("15:04")))
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
		firstNewline := strings.Index(response, "\n")
		if firstNewline > 0 {
			response = response[firstNewline+1:]
		}
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSuffix(response, "`")
	}

	// Find JSON object boundaries
	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")

	if startIdx >= 0 && endIdx > startIdx {
		response = response[startIdx : endIdx+1]
	}

	response = strings.TrimSpace(response)
	return response
}

// GenerateHealthRecommendation generates health recommendations (kept for compatibility)
func (s *OpenAIService) GenerateHealthRecommendation(ctx context.Context, userID uuid.UUID, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error) {
	activitiesStr := formatActivitiesForAI(activities)

	prompt := fmt.Sprintf(`Based on the user's activities and health profile, generate a personalized health recommendation.

Activities:
%s

Health Profile: %s

Generate a concise, helpful health recommendation in Indonesian.`, activitiesStr, formatHealthProfileForAI(healthProfile))

	return s.callAPI(prompt)
}
