# Konfigurasi Waha di .env

Panduan untuk mengkonfigurasi Waha Server di file `.env` Smart Alert System.

## Format Konfigurasi

Di file `.env`, pastikan ada konfigurasi Waha seperti ini:

```bash
# Waha Server Configuration
WAHA_SERVER_URL=http://localhost:3000
WAHA_API_KEY=
```

## Penjelasan Variabel

### WAHA_SERVER_URL
URL base dari Waha Server yang sedang running.

**Contoh:**
- Local: `http://localhost:3000`
- Remote: `http://192.168.1.100:3000`
- Production: `https://waha.yourdomain.com`
- Docker: `http://waha-container:3000`

### WAHA_API_KEY
API Key untuk autentikasi ke Waha Server (opsional).

**Jika Waha Server menggunakan API Key:**
```bash
WAHA_API_KEY=your_secret_api_key_here
```

**Jika tidak menggunakan API Key:**
```bash
WAHA_API_KEY=
# atau bisa dikosongkan
```

## Contoh Konfigurasi Berdasarkan Setup

### 1. Waha di Localhost (Development)
```bash
WAHA_SERVER_URL=http://localhost:3000
WAHA_API_KEY=
```

### 2. Waha di Server Berbeda (Same Network)
```bash
WAHA_SERVER_URL=http://192.168.1.100:3000
WAHA_API_KEY=your_api_key
```

### 3. Waha dengan Docker
```bash
# Jika Waha di container yang sama
WAHA_SERVER_URL=http://localhost:3000

# Jika Waha di container berbeda (docker network)
WAHA_SERVER_URL=http://waha-container:3000
```

### 4. Waha dengan HTTPS (Production)
```bash
WAHA_SERVER_URL=https://waha.yourdomain.com
WAHA_API_KEY=your_production_api_key
```

## Setup Webhook Setelah Konfigurasi

Setelah mengatur `.env`, setup webhook dengan:

### Menggunakan Script
```bash
# Script akan membaca WAHA_SERVER_URL dan WAHA_API_KEY dari .env
./scripts/setup_waha_webhook.sh
```

### Manual dengan cURL
```bash
# Baca dari .env atau set manual
WAHA_URL="http://localhost:3000"
API_KEY=""  # atau API key Anda
WEBHOOK_URL="http://localhost:8080/webhook"

curl -X POST "${WAHA_URL}/api/webhook" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: ${API_KEY}" \
  -d "{
    \"url\": \"${WEBHOOK_URL}\",
    \"events\": [\"message\"]
  }"
```

## Verifikasi Konfigurasi

### 1. Test Koneksi ke Waha Server
```bash
# Test apakah Waha Server accessible
curl http://localhost:3000/api/health
# atau
curl http://localhost:3000/api/sessions
```

### 2. Test Send Message (jika sudah setup)
```bash
# Test mengirim pesan via Waha API
curl -X POST "http://localhost:3000/api/sendText" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{
    "chatId": "6281234567890",
    "text": "Test message"
  }'
```

## Troubleshooting

### Error: Connection Refused
**Masalah:** Smart Alert System tidak bisa connect ke Waha Server

**Solusi:**
1. Pastikan Waha Server sudah running
2. Cek `WAHA_SERVER_URL` di `.env` sudah benar
3. Test koneksi: `curl $WAHA_SERVER_URL/api/health`

### Error: 401 Unauthorized
**Masalah:** API Key tidak valid atau tidak diset

**Solusi:**
1. Cek `WAHA_API_KEY` di `.env`
2. Pastikan API key sesuai dengan yang dikonfigurasi di Waha Server
3. Jika Waha tidak menggunakan API key, biarkan kosong

### Error: Webhook tidak menerima pesan
**Masalah:** Webhook sudah diset tapi tidak menerima pesan

**Solusi:**
1. Pastikan webhook URL accessible dari Waha Server
2. Jika Smart Alert System di localhost, pastikan Waha juga di localhost
3. Jika di server berbeda, gunakan IP/domain yang bisa diakses
4. Cek firewall/network rules

## Catatan Penting

1. **URL harus accessible:**
   - Jika Waha dan Smart Alert System di server berbeda, pastikan network/firewall mengizinkan
   - Untuk production, gunakan HTTPS

2. **API Key (opsional):**
   - Tidak semua setup Waha memerlukan API key
   - Jika Waha Anda tidak menggunakan API key, biarkan `WAHA_API_KEY` kosong

3. **Port:**
   - Default Waha Server: `3000`
   - Default Smart Alert System: `8080`
   - Pastikan tidak ada konflik port

4. **Restart setelah perubahan:**
   - Setelah mengubah `.env`, restart Smart Alert System:
   ```bash
   # Stop server (Ctrl+C)
   # Start lagi
   go run cmd/server/main.go
   ```

