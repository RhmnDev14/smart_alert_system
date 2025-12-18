# Docker Setup Guide

Panduan untuk menjalankan Smart Alert System menggunakan Docker dan Docker Compose.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+

## Quick Start

### 1. Setup Environment

Copy dan edit file `.env`:

```bash
cp env.example .env
```

Edit `.env` dan sesuaikan konfigurasi:
- `WAHA_SERVER_URL`: URL Waha server Anda
- `WAHA_API_KEY`: API key Waha (jika diperlukan)
- `DB_PASSWORD`: Password database (default: postgres)
- `AI_PROVIDER`: `ollama` untuk AI gratis atau `openai` untuk OpenAI

### 2. Pull Ollama Model (Optional)

Jika menggunakan Ollama, pull model terlebih dahulu:

```bash
docker-compose up -d ollama
sleep 10
docker exec smart-alert-ollama ollama pull llama3.2
```

### 3. Start All Services

```bash
docker-compose up -d
```

Ini akan:
- Start PostgreSQL database
- Start Ollama (AI gratis)
- Build dan start Smart Alert System
- Run database migrations otomatis

### 4. Check Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f ollama
```

### 5. Stop Services

```bash
docker-compose down
```

Untuk menghapus volumes juga (data akan hilang):

```bash
docker-compose down -v
```

## Development Mode

Untuk development dengan hot reload, gunakan `docker-compose.dev.yml`:

```bash
# Start only database and Ollama
docker-compose -f docker-compose.dev.yml up -d

# Run app locally (outside Docker)
go run cmd/server/main.go
```

## Services

### PostgreSQL Database
- **Port**: 5432
- **User**: postgres (atau dari `.env`)
- **Password**: postgres (atau dari `.env`)
- **Database**: smart_alert_db
- **Volume**: `postgres_data`

### Ollama (AI)
- **Port**: 11434
- **API**: http://localhost:11434/v1
- **Volume**: `ollama_data`
- **Model**: llama3.2 (default)

### Smart Alert System
- **Port**: 8080
- **Health Check**: http://localhost:8080/health
- **Webhook**: http://localhost:8080/webhook

## Configuration

### Environment Variables

Semua konfigurasi melalui file `.env`:

```bash
# Database
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=smart_alert_db

# Waha Server
WAHA_SERVER_URL=http://host.docker.internal:3000
WAHA_API_KEY=your_api_key

# AI
AI_PROVIDER=ollama
AI_MODEL=llama3.2
AI_BASE_URL=http://ollama:11434/v1

# Application
APP_PORT=8080
TIMEZONE=Asia/Jakarta
```

### Network Configuration

Semua services terhubung melalui network `smart-alert-network`.

Untuk mengakses Waha server yang berjalan di host (bukan di Docker), gunakan:
- `host.docker.internal:3000` (Mac/Windows)
- `172.17.0.1:3000` (Linux)

## Troubleshooting

### Database Connection Error

```bash
# Check database logs
docker-compose logs postgres

# Check database is running
docker-compose ps postgres

# Restart database
docker-compose restart postgres
```

### Ollama Not Responding

```bash
# Check Ollama logs
docker-compose logs ollama

# Pull model manually
docker exec smart-alert-ollama ollama pull llama3.2

# Check if model is available
docker exec smart-alert-ollama ollama list
```

### Migration Errors

```bash
# Run migrations manually
docker exec smart-alert-app ./migrate

# Check migration logs
docker-compose logs app | grep -i migration
```

### Port Already in Use

Jika port sudah digunakan, ubah di `docker-compose.yml`:

```yaml
ports:
  - "8081:8080"  # Change host port
```

### Reset Everything

```bash
# Stop and remove all containers, networks, and volumes
docker-compose down -v

# Remove images
docker-compose down --rmi all

# Start fresh
docker-compose up -d --build
```

## Production Deployment

Untuk production:

1. **Use secrets management**: Jangan hardcode credentials di `.env`
2. **Use reverse proxy**: Nginx atau Traefik untuk HTTPS
3. **Use managed database**: AWS RDS, Google Cloud SQL, dll
4. **Monitor**: Setup logging dan monitoring
5. **Backup**: Setup automated database backups

### Example Production docker-compose

```yaml
version: '3.8'

services:
  app:
    image: your-registry/smart-alert-system:latest
    environment:
      DB_HOST: ${DB_HOST}
      DB_PASSWORD: ${DB_PASSWORD}
      # ... other env vars
    secrets:
      - db_password
      - waha_api_key

secrets:
  db_password:
    external: true
  waha_api_key:
    external: true
```

## Building Custom Image

```bash
# Build image
docker build -t smart-alert-system:latest .

# Tag for registry
docker tag smart-alert-system:latest your-registry/smart-alert-system:v1.0.0

# Push to registry
docker push your-registry/smart-alert-system:v1.0.0
```

## Health Checks

Semua services memiliki health checks:

```bash
# Check service health
docker-compose ps

# Manual health check
curl http://localhost:8080/health
```

## Volumes

Data persisten disimpan di volumes:

- `postgres_data`: Database data
- `ollama_data`: Ollama models

Backup volumes:

```bash
# Backup database
docker run --rm -v smart-alert-system_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz /data

# Restore database
docker run --rm -v smart-alert-system_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres_backup.tar.gz -C /
```

## Tips

1. **Development**: Gunakan `docker-compose.dev.yml` untuk development
2. **Production**: Gunakan `docker-compose.yml` dengan proper secrets
3. **Monitoring**: Setup logging dengan `docker-compose logs -f`
4. **Performance**: Adjust resources di `docker-compose.yml` sesuai kebutuhan

