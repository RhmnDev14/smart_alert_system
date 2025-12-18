# Setup Webhook Waha Server

Panduan lengkap untuk mengkonfigurasi webhook di Waha Server agar terhubung dengan Smart Alert System.

## Prerequisites

1. Waha Server sudah running dan terhubung dengan WhatsApp
2. Smart Alert System sudah running
3. Pastikan kedua server bisa saling mengakses (network/firewall)

## Langkah-langkah Setup

### 1. Pastikan Smart Alert System Running

Jalankan aplikasi Smart Alert System:

```bash
# Build
go build -o bin/server cmd/server/main.go

# Run
./bin/server
```

Server akan berjalan di port yang dikonfigurasi di `.env` (default: `8080`).

Webhook endpoint: `http://your-server:8080/webhook`

### 2. Setup Webhook via Waha API

#### Metode 1: Menggunakan cURL

```bash
# Set webhook URL
curl -X POST "http://localhost:3000/api/webhook" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_WAHA_API_KEY" \
  -d '{
    "url": "http://your-server:8080/webhook",
    "events": ["message"]
  }'
```

#### Metode 2: Menggunakan HTTP Client (Postman/Insomnia)

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:3000/api/webhook`
- **Headers:**
  - `Content-Type: application/json`
  - `X-Api-Key: YOUR_WAHA_API_KEY` (jika menggunakan API key)
- **Body (JSON):**
```json
{
  "url": "http://your-server:8080/webhook",
  "events": ["message"]
}
```

#### Metode 3: Menggunakan Waha Dashboard (jika tersedia)

Jika Waha Server memiliki dashboard web:
1. Buka dashboard Waha
2. Navigasi ke Settings > Webhooks
3. Masukkan URL: `http://your-server:8080/webhook`
4. Pilih events: `message`
5. Save

### 3. Verifikasi Webhook

#### Test dengan mengirim pesan ke WhatsApp

1. Kirim pesan ke nomor WhatsApp yang terhubung dengan Waha
2. Cek log Smart Alert System, seharusnya ada log:
   ```
   Processing message from: 6281234567890
   ```

#### Cek Webhook Status via API

```bash
# Get webhook configuration
curl -X GET "http://localhost:3000/api/webhook" \
  -H "X-Api-Key: YOUR_WAHA_API_KEY"
```

### 4. Format Webhook Payload

Smart Alert System mengharapkan payload dalam format:

```json
{
  "event": "message",
  "data": {
    "id": "message_id",
    "from": "6281234567890",
    "to": "6289876543210",
    "body": "Pesan dari user",
    "timestamp": 1234567890,
    "type": "text"
  }
}
```

## Konfigurasi untuk Production

### 1. Gunakan HTTPS

Untuk production, gunakan HTTPS:

```bash
# Update webhook URL di Waha
curl -X POST "http://your-waha-server:3000/api/webhook" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_WAHA_API_KEY" \
  -d '{
    "url": "https://your-domain.com/webhook",
    "events": ["message"]
  }'
```

### 2. Setup Reverse Proxy (Nginx)

Jika menggunakan Nginx sebagai reverse proxy:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /webhook {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. Environment Variables

Pastikan `.env` sudah dikonfigurasi dengan benar:

```bash
# Waha Server Configuration
WAHA_SERVER_URL=http://localhost:3000
WAHA_API_KEY=your_api_key_if_needed

# Application Configuration
APP_PORT=8080
```

## Troubleshooting

### Webhook tidak menerima pesan

1. **Cek apakah server Smart Alert System running:**
   ```bash
   curl http://localhost:8080/health
   # Should return: OK
   ```

2. **Cek log Waha Server:**
   - Pastikan webhook URL sudah terdaftar
   - Cek apakah ada error saat mengirim webhook

3. **Cek firewall/network:**
   - Pastikan port 8080 terbuka
   - Pastikan Waha Server bisa mengakses Smart Alert System

4. **Test webhook manual:**
   ```bash
   curl -X POST "http://localhost:8080/webhook" \
     -H "Content-Type: application/json" \
     -d '{
       "event": "message",
       "data": {
         "id": "test123",
         "from": "6281234567890",
         "to": "6289876543210",
         "body": "Test message",
         "timestamp": 1234567890,
         "type": "text"
       }
     }'
   ```

### Webhook menerima pesan tapi tidak ada response

1. **Cek log aplikasi:**
   - Pastikan tidak ada error saat processing message
   - Cek apakah AI service berjalan dengan baik

2. **Cek database connection:**
   - Pastikan database bisa diakses
   - Cek apakah migration sudah dijalankan

3. **Cek Waha API key:**
   - Pastikan `WAHA_API_KEY` di `.env` sudah benar
   - Pastikan Waha Server bisa mengirim pesan

## Testing Webhook

### Test dengan script sederhana

Buat file `test_webhook.sh`:

```bash
#!/bin/bash

WEBHOOK_URL="http://localhost:8080/webhook"

curl -X POST "$WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "data": {
      "id": "test_'$(date +%s)'",
      "from": "6281234567890",
      "to": "6289876543210",
      "body": "Tambah kegiatan olahraga jam 6 pagi",
      "timestamp": '$(date +%s)',
      "type": "text"
    }
  }'

echo ""
```

Jalankan:
```bash
chmod +x test_webhook.sh
./test_webhook.sh
```

## Catatan Penting

1. **Webhook URL harus accessible dari Waha Server**
   - Jika Waha di server berbeda, gunakan IP public atau domain
   - Jika localhost, pastikan Waha dan Smart Alert System di server yang sama

2. **Event yang didukung:**
   - `message`: Pesan masuk dari WhatsApp

3. **Rate Limiting:**
   - Waha Server mungkin memiliki rate limiting
   - Smart Alert System akan memproses pesan secara asynchronous

4. **Security:**
   - Pertimbangkan untuk menambahkan authentication/authorization
   - Gunakan HTTPS untuk production
   - Validasi webhook signature jika Waha mendukung

## Referensi

- [Waha Documentation](https://github.com/devlikeapro/waha)
- [Waha API Reference](https://waha.devlike.pro/)

