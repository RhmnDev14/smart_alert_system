# Database Migrations

File-file migration database PostgreSQL untuk Smart Alert System.

## Urutan Migration

1. `001_create_extension_uuid.sql` - Enable UUID extension
2. `002_create_users_table.sql` - Tabel users
3. `003_create_activity_categories_table.sql` - Tabel activity_categories
4. `004_create_activities_table.sql` - Tabel activities
5. `005_create_activity_completions_table.sql` - Tabel activity_completions
6. `006_create_user_health_profiles_table.sql` - Tabel user_health_profiles
7. `007_create_recommendation_types_table.sql` - Tabel recommendation_types
8. `008_create_health_recommendations_table.sql` - Tabel health_recommendations
9. `009_create_message_history_table.sql` - Tabel message_history
10. `010_create_alert_logs_table.sql` - Tabel alert_logs
11. `011_create_scheduled_alerts_table.sql` - Tabel scheduled_alerts
12. `012_seed_initial_data.sql` - Seed data awal (categories dan recommendation types)

## Cara Menjalankan Migration

### Setup Environment

Buat file `.env` di root project (copy dari `env.example`) dengan konfigurasi database:

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
DATABASE_URL=postgresql://user:password@localhost:5432/smart_alert_db?sslmode=disable
```

### Menggunakan Makefile (Recommended - Tidak perlu psql)

```bash
# Jalankan semua migration (otomatis membaca dari .env, menggunakan Go tool)
make migrate-up

# Buat database baru (masih memerlukan psql)
make migrate-create-db

# Drop semua tabel (development only)
make migrate-down
```

### Menggunakan Go Tool (Tidak perlu psql)

```bash
# Tool akan otomatis membaca dari .env
go run cmd/migrate/main.go migrations

# Atau build dulu
go build -o bin/migrate cmd/migrate/main.go
./bin/migrate migrations
```

### Menggunakan Script

```bash
# Script akan otomatis menggunakan Go tool jika psql tidak tersedia
./migrations/run_migrations.sh

# Atau override dengan parameter (jika menggunakan psql)
./migrations/run_migrations.sh "postgresql://user:pass@host:port/db"
```

### Menggunakan psql

```bash
# Set environment variable untuk database connection
export DATABASE_URL="postgresql://user:password@localhost:5432/smart_alert_db"

# Jalankan semua migration
for file in migrations/*.sql; do
    psql $DATABASE_URL -f "$file"
done
```

### Menggunakan golang-migrate

Jika menggunakan library golang-migrate:

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Jalankan migration
migrate -path migrations -database "postgresql://user:password@localhost:5432/smart_alert_db?sslmode=disable" up
```

### Menggunakan psql langsung

```bash
psql -U username -d smart_alert_db -f migrations/001_create_extension_uuid.sql
psql -U username -d smart_alert_db -f migrations/002_create_users_table.sql
# ... dan seterusnya
```

## Rollback Migration

Untuk rollback (jika menggunakan golang-migrate):

```bash
migrate -path migrations -database "postgresql://user:password@localhost:5432/smart_alert_db?sslmode=disable" down
```

## Catatan

- Semua tabel menggunakan UUID sebagai primary key
- Semua timestamp menggunakan `TIMESTAMP WITH TIME ZONE`
- Foreign key constraints menggunakan `ON DELETE CASCADE` atau `ON DELETE SET NULL` sesuai kebutuhan
- Index sudah dibuat untuk kolom-kolom yang sering digunakan dalam query
- Trigger untuk auto-update `updated_at` sudah dibuat untuk tabel yang memerlukan

