# Smart Alert System

Sistem alert pintar yang mengintegrasikan AI dan WhatsApp untuk mengingatkan kegiatan user yang diagendakan serta memberikan rekomendasi kesehatan terkait kegiatan yang dilakukan.

## Fitur Utama

1. **Alert Pagi (05:00 - Dinamis via .env)**
   - Menyapa user dengan penuh semangat dan motivasi
   - Mengingatkan jadwal kegiatan yang diagendakan hari ini
   - Menanyakan apakah ada tambahan rencana/kegiatan lain hari ini
   - Memberikan tips kesehatan personalisasi berdasarkan kegiatan

2. **Summary Malam (22:00 - Dinamis via .env)**
   - Merangkum kegiatan yang telah diselesaikan hari ini
   - Menganalisis pola kegiatan (kesibukan, aktivitas fisik, dll)
   - Memberikan rekomendasi kesehatan sebagai persiapan hari esok

3. **Pengingat Real-Time (Activity Reminder)**
   - Sistem aktif mengecek setiap menit (cron job)
   - Notifikasi otomatis dikirim tepat pada waktu kegiatan yang telah dijadwalkan

4. **Gateway AI Dinamis**
   - User dapat mengirim pesan obrolan biasa atau menjadwalkan kegiatan
   - AI berfungsi sebagai asisten yang mengobrol secara natural (bahasa Indonesia)
   - AI secara otomatis di background mendeteksi, mengekstrak, dan menyimpan jadwal ke database jika ada konteks kegiatan

5. **Welcome Message**
   - Pesan dan panduan interaksi untuk user baru yang mengontak Bot Telegram

## Teknologi

- **Backend**: Go (Golang)
- **Messaging API**: Telegram Bot API / Long-polling
- **AI Gateway**: Model LLM via API (contoh: OpenAI, SumoPod, OpenRouter, atau local LLM)
- **Database**: PostgreSQL (Data persistence & Transactional Manager)
- **Message Queue / Job Manager**: Redis & Asynq (Background Workers Pipeline)

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
2. **Gateway AI** → Membaca pesan, menyiapkan respons obrolan natural, sekaligus mendeteksi kalau ada informasi jadwal.
3. **Penyimpanan Otomatis** → Bila ada jadwal, sistem otomatis merekam detail nama, deskripsi, dan waktu jadwal ke _database PostgreSQL_.
4. **Scheduler & Message Queue Berjalan (Asynq + Redis)** →
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

# AI & LLM Service
AI_PROVIDER=openai
AI_BASE_URL=https://openrouter.ai/api/v1
AI_MODEL=google/gemini-2.5-flash-free
AI_API_KEY=your_api_key_here

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

### Drop Database (Development Only)

```bash
# Hati-hati: ini akan menghapus semua tabel!
make migrate-down
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

## Fitur yang Sudah Diimplementasikan

1. ✅ Clean Architecture dengan separation of concerns
2. ✅ Domain entities sesuai ERD
3. ✅ Repository pattern dengan PostgreSQL
4. ✅ Use cases untuk business logic
5. ✅ Telegram webhook / long-polling handler
6. ✅ AI dinamis Gateway untuk chat sekaligus parsing objek JSON
7. ✅ Redis Message Queue / Asynq Job Workers
8. ✅ Scheduler dinamis (Morning Alert, Evening Summary, dan Pengingat Menit-ke-Menit)
9. ✅ Ekstraksi otomatis natural language ke database time
10. ✅ Transaksional Rollback & Graceful Error Handling dengan Queue Retry System

## Next Steps

1. ✅ Konfigurasi dan testing Gateway LLM OpenRouter
2. ✅ Penyempurnaan Asynq (Message Broker) + Redis Integration
3. ✅ Bungkus logic database ke dalam _Transaction Manager_ untuk konsistensi data
4. Tingkatkan prompt _fallback validation_ untuk sistem error AI
5. Tambahkan unit tests lengkap untuk infrastruktur

# smart_alert_system
