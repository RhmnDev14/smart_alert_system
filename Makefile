.PHONY: migrate-up migrate-down migrate-drop migrate-create-db help load-env

# Load .env file if exists
-include .env
export

# Build DATABASE_URL from .env variables or use default
ifdef DATABASE_URL
	DB_URL := $(DATABASE_URL)
else ifdef DB_HOST
	DB_PORT := $(or $(DB_PORT),5432)
	DB_SSLMODE := $(or $(DB_SSLMODE),disable)
	DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
else
	DB_URL := postgresql://postgres:postgres@localhost:5432/smart_alert_db?sslmode=disable
endif

help: ## Menampilkan help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

migrate-up: ## Menjalankan semua migration (menggunakan Go tool)
	@echo "Running migrations..."
	@go run cmd/migrate/main.go migrations

migrate-down: ## Drop semua tabel (HATI-HATI: hanya untuk development!)
	@echo "⚠️  WARNING: This will drop all tables!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		go run cmd/migrate/main.go migrations; \
		echo "✓ All tables dropped!"; \
	else \
		echo "Cancelled."; \
	fi

migrate-create-db: ## Membuat database baru
	@echo "Creating database..."
	@if [ -n "$(DB_NAME)" ]; then \
		psql -h $(or $(DB_HOST),localhost) -U $(or $(DB_USER),postgres) -c "CREATE DATABASE $(DB_NAME);" || echo "Database mungkin sudah ada"; \
	else \
		psql -U postgres -c "CREATE DATABASE smart_alert_db;" || echo "Database mungkin sudah ada"; \
	fi
	@echo "✓ Database created!"

migrate-drop-db: ## Drop database (HATI-HATI!)
	@echo "⚠️  WARNING: This will drop the entire database!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		if [ -n "$(DB_NAME)" ]; then \
			psql -h $(or $(DB_HOST),localhost) -U $(or $(DB_USER),postgres) -c "DROP DATABASE IF EXISTS $(DB_NAME);"; \
		else \
			psql -U postgres -c "DROP DATABASE IF EXISTS smart_alert_db;"; \
		fi \
		echo "✓ Database dropped!"; \
	else \
		echo "Cancelled."; \
	fi

