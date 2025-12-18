# Smart Alert System

Sistem alert pintar yang mengintegrasikan AI dan WhatsApp untuk mengingatkan kegiatan user yang diagendakan serta memberikan rekomendasi kesehatan terkait kegiatan yang dilakukan.

## Fitur Utama

1. **Alert Pagi (05:00)**
   - Mengingatkan kegiatan yang diagendakan hari ini
   - Memberikan tips kesehatan personalisasi berdasarkan kegiatan

2. **Summary Malam (22:00)**
   - Ringkasan kegiatan yang telah dilakukan hari ini
   - Analisis pola kegiatan
   - Rekomendasi kesehatan untuk hari berikutnya

3. **Input Kegiatan Fleksibel**
   - User dapat menambahkan kegiatan kapan saja
   - Sistem menerima format pesan apa saja (natural language)
   - AI akan memparse dan mengekstrak informasi kegiatan

4. **Welcome Message**
   - Pesan default untuk user baru yang pertama kali mengirim pesan

5. **AI-Powered**
   - Parsing pesan natural language
   - Rekomendasi kesehatan kontekstual
   - Analisis pola kegiatan

## Teknologi

- **Backend**: Go
- **WhatsApp API**: Waha Server
- **AI**: (TBD - bisa menggunakan OpenAI, Gemini, atau LLM lokal)
- **Database**: PostgreSQL

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
│   ├── database/
│   ├── models/
│   ├── handlers/
│   ├── services/
│   │   ├── whatsapp/
│   │   ├── ai/
│   │   ├── scheduler/
│   │   └── activity/
│   └── utils/
└── migrations/
```

## Alur Kerja

1. **User mengirim pesan** → Sistem menerima via Waha
2. **AI memparse pesan** → Ekstrak intent dan data
3. **Proses intent** → Simpan/update/hapus kegiatan atau jawab pertanyaan
4. **Scheduler berjalan** → Alert pagi (05:00) dan summary malam (22:00)
5. **AI generate rekomendasi** → Berdasarkan kegiatan dan profil kesehatan user

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
```

Atau gunakan `DATABASE_URL` lengkap:
```bash
DATABASE_URL=postgresql://user:password@host:port/dbname?sslmode=disable
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
- Go 1.25.1 atau lebih baru
- PostgreSQL 12+
- Waha Server yang sudah running
- API Key untuk AI service (OpenAI, dll)

### Setup

1. **Setup Environment**
```bash
cp env.example .env
# Edit .env dengan konfigurasi Anda
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

# Atau langsung run
go run cmd/server/main.go
```

### Konfigurasi Waha Webhook

Setelah server running, konfigurasi webhook di Waha Server untuk mengarahkan ke:
```
http://your-server:8080/webhook
```

**Cara Setup Webhook:**

1. **Menggunakan Script (Paling Mudah):**
```bash
# Setup webhook dengan default values
./scripts/setup_waha_webhook.sh

# Atau dengan custom URL
./scripts/setup_waha_webhook.sh "http://your-server:8080/webhook" "http://waha-server:3000" "your_api_key"
```

2. **Menggunakan cURL:**
```bash
curl -X POST "http://localhost:3000/api/webhook" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_WAHA_API_KEY" \
  -d '{
    "url": "http://your-server:8080/webhook",
    "events": ["message"]
  }'
```

3. **Dokumentasi Lengkap:**
Lihat [docs/WAHA_WEBHOOK_SETUP.md](./docs/WAHA_WEBHOOK_SETUP.md) untuk panduan lengkap.

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
│       ├── whatsapp/     # Waha client
│       ├── ai/           # AI service
│       └── scheduler/    # Cron scheduler
└── migrations/           # Database migrations
```

## Fitur yang Sudah Diimplementasikan

1. ✅ Clean Architecture dengan separation of concerns
2. ✅ Domain entities sesuai ERD
3. ✅ Repository pattern dengan PostgreSQL
4. ✅ Use cases untuk business logic
5. ✅ WhatsApp webhook handler
6. ✅ AI service untuk parsing intent dan generate rekomendasi
7. ✅ Scheduler untuk alert pagi (05:00) dan summary malam (22:00)
8. ✅ Welcome message untuk user baru
9. ✅ Natural language processing untuk input kegiatan

## Next Steps

1. ✅ Setup database schema berdasarkan ERD
2. ✅ Implementasi koneksi Waha Server
3. ✅ Implementasi AI service untuk parsing dan rekomendasi
4. ✅ Implementasi scheduler untuk alert pagi dan malam
5. ✅ Implementasi handler untuk pesan WhatsApp
6. Testing dan deployment
7. Improve AI prompt untuk parsing yang lebih akurat
8. Tambahkan unit tests
9. Tambahkan error handling yang lebih robust

# smart_alert_system
