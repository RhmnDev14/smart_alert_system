# Setup Webhook untuk Waha Server Remote

Jika Waha Server Anda berada di server remote (bukan localhost), webhook URL harus bisa diakses dari internet.

## Masalah

Jika Waha Server di remote (`https://waha.rahmandigital.my.id/`) tapi webhook URL masih `http://localhost:8080/webhook`, Waha Server tidak akan bisa mengirim webhook ke Smart Alert System yang running di localhost.

## Solusi

### Opsi 1: Menggunakan ngrok (Development/Testing)

**ngrok** membuat tunnel dari localhost ke public URL.

1. **Install ngrok:**
```bash
# macOS
brew install ngrok

# atau download dari https://ngrok.com/download
```

2. **Start ngrok tunnel:**
```bash
# Tunnel port 8080 (atau port Smart Alert System Anda)
ngrok http 8080
```

3. **Copy HTTPS URL dari ngrok:**
```
Forwarding: https://abc123.ngrok.io -> http://localhost:8080
```

4. **Setup webhook dengan URL ngrok:**
```bash
./scripts/setup_waha_webhook.sh "https://abc123.ngrok.io/webhook"
```

**Catatan:** URL ngrok berubah setiap kali restart (kecuali pakai plan berbayar).

### Opsi 2: Deploy Smart Alert System ke Server Public

Deploy Smart Alert System ke server yang bisa diakses dari internet.

1. **Setup di server (VPS/Cloud):**
```bash
# Clone dan setup di server
git clone your-repo
cd smart-alert-system
cp env.example .env
# Edit .env dengan konfigurasi server
make migrate-up
go build -o bin/server cmd/server/main.go
./bin/server
```

2. **Setup webhook dengan URL server:**
```bash
./scripts/setup_waha_webhook.sh "https://your-server.com/webhook"
```

### Opsi 3: Gunakan Public IP dengan Port Forwarding

Jika Smart Alert System di localhost tapi punya public IP:

1. **Setup port forwarding di router:**
   - Forward port 8080 ke IP localhost Anda

2. **Setup webhook dengan public IP:**
```bash
./scripts/setup_waha_webhook.sh "http://YOUR_PUBLIC_IP:8080/webhook"
```

**Catatan:** Ini kurang aman untuk production, gunakan HTTPS.

### Opsi 4: Gunakan Domain dengan Reverse Proxy

Jika punya domain dan server:

1. **Setup Nginx reverse proxy:**
```nginx
server {
    listen 80;
    server_name smart-alert.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

2. **Setup SSL dengan Let's Encrypt:**
```bash
sudo certbot --nginx -d smart-alert.yourdomain.com
```

3. **Setup webhook:**
```bash
./scripts/setup_waha_webhook.sh "https://smart-alert.yourdomain.com/webhook"
```

## Verifikasi Setup

### 1. Test Webhook URL Accessible

```bash
# Test dari command line (harus bisa diakses)
curl https://your-webhook-url.com/webhook/health

# Atau test dari browser
# Buka: https://your-webhook-url.com/health
```

### 2. Test Webhook Manual

```bash
curl -X POST "https://your-webhook-url.com/webhook" \
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

### 3. Cek Log Smart Alert System

Setelah mengirim test webhook, cek log aplikasi:
```
Processing message from: 6281234567890
```

## Troubleshooting

### Webhook tidak menerima pesan dari Waha

1. **Cek apakah webhook URL accessible:**
   ```bash
   curl https://your-webhook-url.com/health
   ```

2. **Cek firewall:**
   - Pastikan port terbuka
   - Cek security group (jika di cloud)

3. **Cek log Waha Server:**
   - Lihat apakah ada error saat mengirim webhook
   - Cek response dari webhook

4. **Test dengan ngrok dulu:**
   - Gunakan ngrok untuk memastikan setup benar
   - Jika ngrok bekerja, masalahnya di network/firewall

### Error: Connection Refused

**Masalah:** Waha Server tidak bisa connect ke webhook URL

**Solusi:**
- Pastikan Smart Alert System running
- Pastikan webhook URL benar (bukan localhost jika Waha remote)
- Cek firewall/network rules

### Error: Timeout

**Masalah:** Webhook request timeout

**Solusi:**
- Pastikan Smart Alert System merespons dengan cepat
- Cek network latency
- Pastikan tidak ada blocking di firewall

## Rekomendasi

- **Development:** Gunakan ngrok
- **Staging:** Deploy ke staging server dengan domain
- **Production:** Deploy ke production server dengan HTTPS dan domain

## Contoh Setup Lengkap dengan ngrok

```bash
# Terminal 1: Start Smart Alert System
go run cmd/server/main.go

# Terminal 2: Start ngrok
ngrok http 8080

# Terminal 3: Setup webhook dengan URL ngrok
./scripts/setup_waha_webhook.sh "https://abc123.ngrok.io/webhook"

# Test: Kirim pesan ke WhatsApp
# Cek log di Terminal 1 untuk melihat pesan masuk
```

