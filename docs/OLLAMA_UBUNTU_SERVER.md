# Setup Ollama di Ubuntu Server

Panduan khusus untuk install Ollama di Ubuntu Server (seperti Tencent Cloud, AWS, dll).

## Install Ollama (Cara Paling Mudah)

**JANGAN build dari source!** Gunakan binary yang sudah jadi:

```bash
# Install Ollama dengan script resmi
curl -fsSL https://ollama.ai/install.sh | sh
```

Setelah install, Ollama akan otomatis berjalan sebagai service.

## Verifikasi Install

```bash
# Cek apakah Ollama berjalan
ollama --version

# Cek service status
systemctl status ollama

# Test API
curl http://localhost:11434/api/tags
```

## Download Model

```bash
# Model kecil dan cepat (recommended untuk server)
ollama pull llama3.2

# Atau model lebih pintar (butuh lebih banyak RAM)
ollama pull llama3.1
ollama pull mistral
```

## Konfigurasi Firewall (Jika Perlu)

Jika server Anda memiliki firewall, buka port 11434:

```bash
# UFW (Ubuntu Firewall)
sudo ufw allow 11434/tcp

# Atau untuk firewall lain, pastikan port 11434 terbuka
```

## Konfigurasi Smart Alert System

Update file `.env` di server:

```bash
# AI Configuration
AI_PROVIDER=ollama
AI_API_KEY=
AI_MODEL=llama3.2
AI_BASE_URL=http://localhost:11434/v1
```

Jika Smart Alert System berjalan di server yang sama dengan Ollama, gunakan `localhost`.
Jika berbeda, ganti dengan IP server Ollama: `http://YOUR_SERVER_IP:11434/v1`

## Troubleshooting

### Error: "connection refused"
```bash
# Pastikan Ollama service berjalan
sudo systemctl status ollama

# Jika tidak berjalan, start service
sudo systemctl start ollama

# Enable auto-start saat boot
sudo systemctl enable ollama
```

### Error: "port already in use"
```bash
# Cek apakah ada proses lain yang menggunakan port 11434
sudo lsof -i :11434

# Atau ubah port Ollama (edit service file)
sudo systemctl edit ollama
# Tambahkan: Environment="OLLAMA_HOST=0.0.0.0:11435"
# Lalu update AI_BASE_URL di .env
```

### Model tidak ditemukan
```bash
# List model yang sudah di-download
ollama list

# Download model yang diperlukan
ollama pull llama3.2
```

### RAM tidak cukup
```bash
# Cek penggunaan RAM
free -h

# Gunakan model yang lebih kecil
ollama pull llama3.2  # ~2GB RAM
# Atau
ollama pull phi3      # ~2GB RAM, lebih kecil
```

## Menjalankan Ollama sebagai Service

Ollama biasanya sudah otomatis terinstall sebagai service. Jika tidak:

```bash
# Cek service
systemctl status ollama

# Start service
sudo systemctl start ollama

# Enable auto-start
sudo systemctl enable ollama

# Restart service
sudo systemctl restart ollama

# Lihat logs
sudo journalctl -u ollama -f
```

## Test Ollama

```bash
# Test dengan curl
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "Halo, apa kabar?",
  "stream": false
}'

# Test dengan ollama CLI
ollama run llama3.2 "Halo, apa kabar?"
```

## Tips untuk Server

1. **Gunakan model kecil** untuk menghemat RAM (llama3.2 recommended)
2. **Monitor RAM usage** - pastikan server punya cukup RAM
3. **Enable swap** jika RAM terbatas:
   ```bash
   sudo fallocate -l 4G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```
4. **Auto-restart Ollama** jika crash (sudah otomatis dengan systemd)
5. **Monitor logs** untuk debugging:
   ```bash
   sudo journalctl -u ollama -f
   ```

## Perbandingan Model untuk Server

| Model | Size | RAM Needed | Speed | Quality |
|-------|------|------------|-------|---------|
| llama3.2 | ~2GB | 4-6GB | ⚡⚡⚡ | ⭐⭐⭐ |
| llama3.1 | ~4.7GB | 8-10GB | ⚡⚡ | ⭐⭐⭐⭐ |
| mistral | ~4.1GB | 8-10GB | ⚡⚡ | ⭐⭐⭐⭐ |
| phi3 | ~2.3GB | 4-6GB | ⚡⚡⚡ | ⭐⭐⭐ |

**Rekomendasi untuk server:** `llama3.2` atau `phi3` (kecil, cepat, cukup pintar)

