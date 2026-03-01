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
	ProcessMessage(ctx context.Context, message string, userName string, currentTime time.Time, chatHistory []*entity.MessageHistory, userMemories []*entity.UserMemory) (*entity.AIProcessResult, error)
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

func (s *OpenAIService) callAPIWithMessages(messages []Message) (string, error) {
	url := s.baseURL
	if !strings.HasSuffix(url, "/chat/completions") {
		url = fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(s.baseURL, "/"))
	}

	reqBody := OpenAIRequest{
		Model:    s.model,
		Messages: messages,
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

func (s *OpenAIService) ProcessMessage(ctx context.Context, message string, userName string, currentTime time.Time, chatHistory []*entity.MessageHistory, userMemories []*entity.UserMemory) (*entity.AIProcessResult, error) {
	// Build persistent memory context
	var memoryContext string
	if len(userMemories) > 0 {
		var sb strings.Builder
		sb.WriteString("\n\nHAL-HAL YANG KAMU INGAT TENTANG USER INI:\n")
		for _, mem := range userMemories {
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", mem.MemoryType, mem.Content))
		}
		memoryContext = sb.String()
	}

	systemPrompt := fmt.Sprintf(`Kamu adalah asisten pintar bernama "Smart Alert Bot". Tugasmu ada 3:

1. **Merespons pesan user secara natural** — bisa ngobrol biasa, menjawab pertanyaan, dll.
2. **Mendeteksi dan menjalankan operasi data kegiatan** — berdasarkan pesan user, deteksi apakah mereka ingin:
   - **create**: Membuat jadwal/kegiatan baru
   - **get**: Melihat/mencari jadwal kegiatan yang sudah ada
   - **update**: Mengubah/memperbarui jadwal kegiatan (reschedule, cancel, complete, ubah judul, dll)
   - **none**: Hanya obrolan biasa tanpa operasi data, ATAU informasi jadwal belum lengkap dan perlu ditanyakan dulu
3. **Mengingat informasi penting tentang user** — jika user menyebutkan preferensi, kebiasaan, fakta personal, atau info penting tentang diri mereka, simpan sebagai memory.

INFO USER SAAT INI:
- Nama lengkap user yang sedang bicara: %s
- Waktu saat ini: %s
%s
ATURAN PENTING:
- Selalu balas dalam Bahasa Indonesia yang natural dan friendly
- User yang sedang berbicara denganmu bernama "%s". INGAT NAMA INI BAIK-BAIK.
- Jika user bertanya "saya siapa", "siapa saya", "siapa nama saya", atau pertanyaan identitas serupa, JAWAB DENGAN NAMA USER yaitu "%s". JANGAN menjawab tentang identitas kamu (bot). Contoh jawaban yang benar: "Kamu adalah %s! Ada yang bisa saya bantu?"
- Deteksi intent user dengan cermat dari konteks percakapan
- Gunakan informasi dari memory (jika ada) untuk memberikan respons yang lebih personal

ATURAN MEMORY (PERSISTENT MEMORY):
- Jika user menyebutkan preferensi, kebiasaan, atau fakta tentang diri mereka, tambahkan ke memories_to_save
- Tipe memory: "preference" (kesukaan/preferensi), "fact" (fakta), "habit" (kebiasaan), "personal" (info personal)
- Contoh: "Saya biasa jogging jam 5 pagi" → habit: "User biasa jogging jam 5 pagi"
- Contoh: "Saya alergi kacang" → personal: "User alergi kacang"
- Contoh: "Saya lebih suka reminder 30 menit sebelum" → preference: "User suka reminder 30 menit sebelum jadwal"
- JANGAN simpan memory untuk hal yang sudah diingat (cek daftar memory di atas)
- Jika tidak ada info baru yang perlu diingat, kosongkan memories_to_save (array kosong [])

ATURAN PER ACTION:

**CREATE** — Jika user menyebutkan jadwal/kegiatan/rencana baru:
- Contoh LENGKAP: "besok meeting jam 10", "hari ini mau olahraga jam 6 sore", "tambahkan jadwal belajar jam 8 malam"
- PENTING: Jika informasi WAKTU/JAM belum disebutkan secara spesifik, JANGAN langsung buat jadwal! Tanyakan dulu jam/waktu spesifiknya.
  - Contoh kurang lengkap: "mau bukber jumat depan" → tanyakan "Jam berapa bukbernya?"
  - Contoh kurang lengkap: "ada rencana jalan-jalan besok" → tanyakan "Jam berapa rencananya?"
  - Contoh kurang lengkap: "mau meeting minggu depan" → tanyakan "Hari apa dan jam berapa meetingnya?"
  - Jika waktu kurang spesifik, set action = "none" dan tanyakan detailnya di response
- Hanya set action = "create" jika user sudah menyebutkan JAM/WAKTU secara jelas (misal: "jam 5", "pukul 10", "sore jam 4", "jam 8 malam")
- Jika waktu sudah jelas: Set action = "create", has_schedule = true, isi field "schedule"
- scheduled_time dalam format ISO 8601 (contoh: "2026-02-23T10:00:00+07:00")
- Jika user bilang "besok", hitung tanggal saat ini + 1 hari
- Jika user bilang "hari ini", gunakan tanggal saat ini
- Response harus mengkonfirmasi bahwa jadwal akan disimpan

**GET** — HANYA jika user minta lihat / cari KEGIATAN PRIBADI mereka yang tersimpan di sistem:
- Contoh: "lihat jadwal hari ini", "ada apa besok?", "jadwal apa saja yang pending?", "cari meeting saya", "tampilkan semua kegiatan"
- PENTING: Action "get" HANYA untuk kegiatan PRIBADI user di database, BUKAN untuk pertanyaan umum!
  - "jadwal saya hari ini" → action = "get" (kegiatan pribadi user)
  - "jadwal pertandingan liga Indonesia" → action = "none" (pertanyaan umum, jawab secara natural)
  - "ada pertandingan Persija hari ini?" → action = "none" (pertanyaan umum, jawab secara natural)
  - "berita hari ini apa?" → action = "none" (pertanyaan umum, jawab secara natural)
  - "kapan meeting saya?" → action = "get" (kegiatan pribadi user)
- Set action = "get", has_schedule = false, isi field "query"
- filter_type bisa: "today", "tomorrow", "date", "status", "search", "all"
- Untuk "date", isi field date dengan format "YYYY-MM-DD"
- Untuk "status", isi field status dengan: "pending", "completed", "cancelled", "overdue"
- Untuk "search", isi field keyword
- Response harus bilang bahwa kamu akan mencarikan datanya

**UPDATE** — Jika user minta ubah/update kegiatan yang sudah ada:
- Contoh: "ubah meeting besok jadi jam 2 siang", "cancel jadwal olahraga", "tandai meeting sebagai selesai", "reschedule rapat ke lusa", "hapus jadwal belajar"
- Set action = "update", has_schedule = false, isi field "update"
- search_title = kata kunci judul kegiatan yang akan diubah (cukup kata kuncinya saja, tidak perlu exact)
- Isi field yang berubah saja (new_title, new_description, new_scheduled_time, new_status, new_priority)
- Untuk cancel: new_status = "cancelled"
- Untuk selesai/complete: new_status = "completed"  
- Untuk reschedule: isi new_scheduled_time dengan format ISO 8601
- Response harus mengkonfirmasi perubahan yang akan dilakukan

**NONE** — Gunakan action "none" untuk:
- Obrolan biasa, sapaan, atau pertanyaan umum
- Informasi jadwal yang belum lengkap (perlu klarifikasi waktu dulu)
- Pertanyaan tentang pengetahuan umum, berita, olahraga, cuaca, dll — JAWAB secara natural dan sebaik mungkin seperti asisten cerdas
- Contoh: "jadwal pertandingan persija" → jawab dengan informasi seputar itu sebaik yang kamu bisa, action tetap "none"
- Set action = "none", has_schedule = false

KAMU HARUS MEMBALAS DENGAN FORMAT JSON BERIKUT (TANPA markdown code block, HANYA JSON murni):

{
  "response": "pesan balasan",
  "has_schedule": true/false,
  "action": "create/get/update/none",
  "schedule": { ... },
  "query": { ... },
  "update": { ... },
  "memories_to_save": [
    {"type": "preference/fact/habit/personal", "content": "deskripsi memory"}
  ]
}

Contoh lengkap CREATE + memory:
{
  "response": "Baik, jadwal jogging kamu jam 5 pagi besok sudah saya simpan!",
  "has_schedule": true,
  "action": "create",
  "schedule": {
    "title": "Jogging",
    "description": "Jogging pagi",
    "scheduled_time": "2026-02-23T05:00:00+07:00",
    "priority": 3
  },
  "memories_to_save": [
    {"type": "habit", "content": "User biasa jogging jam 5 pagi"}
  ]
}

Contoh NONE tanpa memory:
{
  "response": "pesan balasan",
  "has_schedule": false,
  "action": "none",
  "memories_to_save": []
}

Priority: 1(rendah) - 5(tinggi), default 3.
Untuk field yang tidak berubah pada UPDATE, set null.

INGAT: Balas HANYA JSON, tanpa backtick, tanpa markdown, tanpa penjelasan tambahan. Mulai dengan { dan akhiri dengan }.`, userName, currentTime.Format("2006-01-02T15:04:05-07:00"), memoryContext, userName, userName, userName)

	// Build multi-turn message array
	messages := []Message{
		{Role: "system", Content: systemPrompt},
	}

	// Add conversation history as proper user/assistant turns
	if len(chatHistory) > 0 {
		// chatHistory is DESC order, reverse for chronological
		for i := len(chatHistory) - 1; i >= 0; i-- {
			msg := chatHistory[i]
			if msg.MessageType == entity.MessageTypeIncoming {
				messages = append(messages, Message{Role: "user", Content: msg.MessageContent})
			} else if msg.MessageType == entity.MessageTypeOutgoing {
				messages = append(messages, Message{Role: "assistant", Content: msg.MessageContent})
			}
		}
	}

	// Add current message
	messages = append(messages, Message{Role: "user", Content: message})

	// Call API with multi-turn messages
	response, err := s.callAPIWithMessages(messages)
	if err != nil {
		log.Printf("❌ AI API call failed: %v", err)
		return &entity.AIProcessResult{
			Response:    "Maaf, saya sedang mengalami gangguan. Silakan coba lagi nanti.",
			HasSchedule: false,
			Action:      entity.ActionNone,
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
			Action:      entity.ActionNone,
		}, nil
	}

	// Backward compatibility: if action is empty but has_schedule is true, treat as create
	if result.Action == "" && result.HasSchedule {
		result.Action = entity.ActionCreate
	} else if result.Action == "" {
		result.Action = entity.ActionNone
	}

	return &result, nil
}

func (s *OpenAIService) GenerateMorningAlert(ctx context.Context, activities []*entity.Activity, healthProfile *entity.UserHealthProfile) (string, error) {
	activitiesStr := formatActivitiesForAI(activities)

	prompt := fmt.Sprintf(`Buatkan pesan pengingat pagi hari (Morning Alert) dalam bahasa Indonesia dengan kriteria berikut:
1. Sapa user dengan ramah dan hangat untuk memulai hari.
2. Jika ada kegiatan terjadwal hari ini, tampilkan daftar jadwal kegiatannya beserta jam-nya.
3. Jika tidak ada kegiatan, sampaikan bahwa belum ada jadwal hari ini.
4. Tanyakan apakah ada kegiatan lain yang ingin ditambahkan untuk hari ini.
5. JANGAN berikan tips kesehatan atau motivasi panjang. Buat pesan singkat dan to the point.
6. DI AKHIR PESAN, SELALU TAMBAHKAN TEKS PERSIS SEPERTI INI:
Smart Alert System
Develop by Rahman Umardi

Jadwal kegiatan hari ini:
%s

Buat pesan yang singkat, hangat, dan bersahabat.`, activitiesStr)

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
