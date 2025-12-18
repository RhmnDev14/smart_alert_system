# Setup Ollama (Free AI)

Ollama adalah AI lokal yang 100% gratis dan tidak memerlukan API key. Ini adalah alternatif yang sempurna untuk OpenAI jika Anda ingin menggunakan AI tanpa biaya.

## Instalasi Ollama

### macOS
```bash
# Download dan install dari https://ollama.ai
# Atau menggunakan Homebrew:
brew install ollama
```

### Linux (Ubuntu/Debian)
```bash
# Install Ollama dengan script resmi (paling mudah)
curl -fsSL https://ollama.ai/install.sh | sh

# Atau jika script tidak bekerja, install manual:
# 1. Download binary
curl -L https://ollama.ai/download/ollama-linux-amd64 -o /usr/local/bin/ollama
chmod +x /usr/local/bin/ollama

# 2. Buat systemd service (opsional, untuk auto-start)
sudo tee /etc/systemd/system/ollama.service > /dev/null <<EOF
[Unit]
Description=Ollama Service
After=network-online.target

[Service]
ExecStart=/usr/local/bin/ollama serve
User=$USER
Group=$USER
Restart=always
RestartSec=3

[Install]
WantedBy=default.target
EOF

# 3. Enable dan start service
sudo systemctl enable ollama
sudo systemctl start ollama
```

### Windows
Download installer dari https://ollama.ai/download

## Menjalankan Ollama

Setelah instalasi, jalankan Ollama server:
```bash
ollama serve
```

Server akan berjalan di `http://localhost:11434`

## Download Model

Download model yang akan digunakan (pilih salah satu):

```bash
# Model kecil dan cepat (recommended untuk development)
ollama pull llama3.2

# Model lebih besar dan lebih pintar
ollama pull llama3.1
ollama pull mistral
ollama pull qwen2.5
```

**Rekomendasi:** Gunakan `llama3.2` untuk development karena lebih cepat dan menggunakan lebih sedikit RAM.

## Konfigurasi .env

Update file `.env` Anda:

```bash
# AI Configuration
AI_PROVIDER=ollama
AI_API_KEY=                    # Kosongkan untuk Ollama
AI_MODEL=llama3.2             # Sesuaikan dengan model yang Anda download
AI_BASE_URL=http://localhost:11434/v1
```

## Testing

1. Pastikan Ollama server berjalan:
   ```bash
   ollama serve
   ```

2. Test model:
   ```bash
   ollama run llama3.2 "Halo, apa kabar?"
   ```

3. Restart Smart Alert System:
   ```bash
   ./bin/server
   ```

4. Cek log saat startup, seharusnya muncul:
   ```
   ✓ AI Service: Ollama (Free)
     Base URL: http://localhost:11434/v1
     Model: llama3.2
   ```

## Model yang Tersedia

| Model | Size | Speed | Quality | Recommended For |
|-------|------|--------|----------|-----------------|
| llama3.2 | ~2GB | ⚡⚡⚡ | ⭐⭐⭐ | Development, Testing |
| llama3.1 | ~4.7GB | ⚡⚡ | ⭐⭐⭐⭐ | Production |
| mistral | ~4.1GB | ⚡⚡ | ⭐⭐⭐⭐ | Production |
| qwen2.5 | ~4.4GB | ⚡⚡ | ⭐⭐⭐⭐ | Production |

## Troubleshooting

### Error: "connection refused"
- Pastikan Ollama server berjalan: `ollama serve`
- Cek apakah port 11434 terbuka: `curl http://localhost:11434/api/tags`

### Error: "model not found"
- Download model terlebih dahulu: `ollama pull llama3.2`
- Pastikan nama model di `.env` sesuai dengan model yang sudah di-download

### Response lambat
- Gunakan model yang lebih kecil (llama3.2)
- Pastikan RAM cukup (minimal 8GB untuk llama3.2)
- Tutup aplikasi lain yang menggunakan banyak RAM

## Keuntungan Ollama

✅ **100% Gratis** - Tidak ada biaya API  
✅ **Privasi** - Data tidak dikirim ke server eksternal  
✅ **Offline** - Bisa digunakan tanpa internet  
✅ **Tidak ada Rate Limit** - Tidak ada batasan request  
✅ **OpenAI-Compatible** - Menggunakan API yang sama dengan OpenAI  

## Kekurangan Ollama

⚠️ **Membutuhkan RAM** - Minimal 8GB untuk model kecil  
⚠️ **Lebih Lambat** - Dibandingkan dengan OpenAI API (tergantung hardware)  
⚠️ **Perlu Setup** - Harus install dan download model  

## Tips

1. Untuk development, gunakan `llama3.2` (kecil dan cepat)
2. Untuk production, pertimbangkan `llama3.1` atau `mistral` (lebih pintar)
3. Jika RAM terbatas, gunakan model yang lebih kecil
4. Pastikan Ollama server selalu berjalan saat menggunakan Smart Alert System

