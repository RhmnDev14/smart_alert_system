# Smart Alert System

Sistem alert pintar yang mengintegrasikan AI dan Telegram untuk mengingatkan serta mengelola kegiatan user yang diagendakan — dilengkapi fitur **Persistent Memory** dan **Multi-turn Conversation**.

🤖 **Cobalah Bot Ini di Telegram:** [@SmartAlertActivity_bot](https://t.me/SmartAlertActivity_bot)

## Fitur Utama

1. **Alert Pagi (05:00 - Dinamis via .env)**
   - Menyapa user untuk memulai hari
   - Mengingatkan jadwal kegiatan yang diagendakan hari ini
   - Menanyakan apakah ada tambahan rencana/kegiatan lain untuk hari ini

2. **Summary Malam (22:00 - Dinamis via .env)**
   - Merangkum daftar kegiatan yang telah diselesaikan hari ini
   - Memberikan pesan penutup untuk beristirahat

3. **Pengingat Real-Time (Activity Reminder)**
   - Sistem aktif mengecek setiap menit (cron job)
   - Notifikasi otomatis dikirim tepat pada waktu kegiatan yang telah dijadwalkan

4. **Gateway AI Dinamis dengan CRUD Flexible**
   - User dapat mengirim pesan obrolan biasa atau menjadwalkan kegiatan
   - AI berfungsi sebagai asisten yang mengobrol secara natural (bahasa Indonesia)
   - AI secara otomatis mendeteksi intent: **create**, **get**, **update**, atau **none**
   - **Klarifikasi otomatis** — jika jadwal kurang spesifik (misal tanpa jam), AI bertanya dulu sebelum menyimpan
   - **Fuzzy search** dengan PostgreSQL Trigram (`pg_trgm`) — typo seperti "joging" tetap cocok dengan "Jogging"

5. **🧠 Persistent Memory (OpenClaw-style)**
   - Bot mengingat preferensi, kebiasaan, dan fakta personal user
   - Tipe memory: `preference`, `fact`, `habit`, `personal`
   - Memory disimpan di tabel `user_memories` dan di-load setiap ada pesan masuk

6. **💬 Multi-turn Conversation**
   - Menggunakan format multi-turn message array (system + user/assistant history)
   - 10 pesan terakhir di-load sebagai konteks percakapan
   - AI bisa memahami jawaban follow-up

7. **Welcome Message**
   - Pesan dan panduan interaksi untuk user baru yang mengontak Bot Telegram

## Teknologi

- **Backend**: Go (Golang)
- **Messaging API**: Telegram Bot API / Long-polling
- **AI Gateway**: Model LLM via API (contoh: OpenAI, SumoPod, OpenRouter, Google Gemini, atau local LLM)
- **Database**: PostgreSQL (Data persistence & Transactional Manager)
- **Database Extensions**: `pg_trgm` (Fuzzy text matching / trigram similarity)
- **Message Queue / Job Manager**: Redis & Asynq (Background Workers Pipeline)
- **Persistent Memory**: PostgreSQL-backed user memory system (preferences, facts, habits)

## Dokumentasi

- [Flowchart Sistem](./FLOWCHART.md) - Diagram alur proses sistem
- [ERD](./ERD.md) - Entity Relationship Diagram database
- [Migrations](./migrations/README.md) - Dokumentasi database migrations

## Struktur Proyek

```
smart-alert-system/
├── go.mod
├── README.md
├── FLOWCHART.md
├── ERD.md
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── domain/
│   ├── handler/
│   ├── usecase/
│   ├── infrastructure/
│   │   ├── database/
│   │   ├── telegram/
│   │   ├── ai/
│   │   ├── scheduler/
│   │   └── repository/
└── migrations/
```

## Alur Kerja

1. **User mengirim pesan** → Diterima lewat polling Telegram Bot
2. **Load Konteks** → Sistem mengambil 10 pesan terakhir (chat history) dan persistent memories user dari database
3. **Gateway AI (Multi-turn)** → Membaca pesan dengan konteks percakapan penuh dan memories, menyiapkan respons obrolan natural, sekaligus mendeteksi apakah ada intent CRUD.
4. **Klarifikasi Otomatis** → Jika jadwal kurang spesifik (misal tanpa jam), AI bertanya dulu sebelum menyimpan.
5. **Fuzzy Search** → Pencarian kegiatan menggunakan trigram similarity (`pg_trgm`), sehingga typo seperti "joging" tetap cocok dengan "Jogging".
6. **Persistent Memory** → Jika AI mendeteksi informasi penting tentang user (preferensi, kebiasaan, dll), otomatis disimpan ke database `user_memories` untuk digunakan di percakapan mendatang.
7. **Penyimpanan Otomatis** → Bila ada jadwal lengkap, sistem otomatis merekam detail nama, deskripsi, dan waktu jadwal ke _database PostgreSQL_.
8. **Scheduler & Message Queue Berjalan (Asynq + Redis)** →
   - Scheduler menjalankan tugas rutin berdasarkan jam (`MORNING_ALERT_TIME` dan `EVENING_SUMMARY_TIME`) lalu menembakkan event ke **Redis Message Queue** (Asynq).
   - _Background Job Worker_ menarik tugas dari antrean Redis dan memproses pemanggilan API AI + Pengiriman Telegram secara paralel non-blocking, dilengkapi fitur _Auto Retry on Failure_ & _Transactional Rollback_ tingkat database.
   - Demikian pula Pengingat Real-time menit-ke-menit bekerja mengecek aktivitas _pending_, mendorongnya ke antrean, lalu merubah statusnya menjadi _completed_ segera setelah notifikasi berhasil terkirim.

## Setup Database

### Prerequisites

- PostgreSQL 12+ terinstall
- Akses ke database PostgreSQL

### Konfigurasi Environment

1. Copy file `env.example` ke `.env`:

```bash
cp env.example .env
```

2. Edit file `.env` dan sesuaikan dengan konfigurasi database Anda:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=smart_alert_db
DB_SSLMODE=disable

# Redis (Untuk Background Job Asynq)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Telegram
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here

# AI & LLM Service (Google Gemini via OpenAI-compatible endpoint)
AI_PROVIDER=gemini
AI_BASE_URL=https://generativelanguage.googleapis.com/v1beta/openai
AI_MODEL=gemini-2.0-flash
AI_API_KEY=your_gemini_api_key_here
# Dapatkan API Key dari: https://aistudio.google.com/apikey

# Konfigurasi Jam
MORNING_ALERT_TIME=05:00
EVENING_SUMMARY_TIME=22:00
```

### Menjalankan Migration

#### Menggunakan Makefile (Recommended - Tidak perlu psql)

```bash
# Jalankan semua migration (akan membaca dari .env, menggunakan Go tool)
make migrate-up
```

#### Menggunakan Go Tool Langsung

```bash
# Tool akan otomatis membaca dari .env
go run cmd/migrate/main.go migrations

# Atau build dulu untuk performa lebih baik
go build -o bin/migrate cmd/migrate/main.go
./bin/migrate migrations
```

#### Menggunakan Script (Fallback ke Go tool jika psql tidak ada)

```bash
# Script akan otomatis menggunakan Go tool jika psql tidak tersedia
./migrations/run_migrations.sh
```

#### Menggunakan psql (Jika terinstall)

```bash
# Set database URL
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/smart_alert_db?sslmode=disable"

# Jalankan migration dengan psql
for file in migrations/*.sql; do
    psql "$DATABASE_URL" -f "$file"
done
```

## Menjalankan Aplikasi

### Prerequisites

- Go 1.23+
- PostgreSQL 12+
- Telegram Bot Token (Dapatkan dari [@BotFather](https://t.me/botfather))
- API Key untuk AI service (OpenAI / kompatibel OpenAI)

### Setup

1. **Setup Environment**

```bash
cp env.example .env
# Edit .env dengan konfigurasi Database, Telegram Token, dan API AI
```

2. **Jalankan Migration**

```bash
make migrate-up
```

3. **Build dan Run**

```bash
# Build
go build -o bin/server cmd/server/main.go

# Run
./bin/server

# Atau langsung run tanpa build
go run cmd/server/main.go
```

_Bot akan otomatis menggunakan Long-polling secara background untuk mengambil pesan dari Telegram._

## Struktur Clean Architecture

```
smart-alert-system/
├── cmd/
│   ├── migrate/          # Migration tool
│   └── server/           # Main application
├── internal/
│   ├── config/           # Configuration
│   ├── domain/           # Domain layer
│   │   ├── entity/       # Domain entities
│   │   └── repository/  # Repository interfaces
│   ├── usecase/          # Business logic
│   ├── handler/          # HTTP handlers
│   └── infrastructure/   # External services
│       ├── database/     # Database connection
│       ├── repository/   # Repository implementations
│       ├── telegram/     # Telegram Bot Client
│       ├── ai/           # AI service gateway
│       └── scheduler/    # Cron scheduler
└── migrations/           # Database migrations
```
