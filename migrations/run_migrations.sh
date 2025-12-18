#!/bin/bash

# Script untuk menjalankan semua migration database
# Usage: ./run_migrations.sh [database_url]
# Script akan membaca dari .env jika ada, atau menggunakan parameter/default

set -e

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# Load .env file if exists
if [ -f "$PROJECT_ROOT/.env" ]; then
    echo "Loading environment from .env file..."
    export $(cat "$PROJECT_ROOT/.env" | grep -v '^#' | xargs)
fi

# Build DATABASE_URL from .env variables or use provided parameter
if [ -n "$1" ]; then
    # Use provided parameter
    DATABASE_URL="$1"
elif [ -n "$DATABASE_URL" ]; then
    # Use DATABASE_URL from .env
    DATABASE_URL="$DATABASE_URL"
elif [ -n "$DB_HOST" ] && [ -n "$DB_USER" ] && [ -n "$DB_PASSWORD" ] && [ -n "$DB_NAME" ]; then
    # Build from individual DB variables
    DB_PORT="${DB_PORT:-5432}"
    DB_SSLMODE="${DB_SSLMODE:-disable}"
    DATABASE_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
else
    # Default fallback
    DATABASE_URL="postgresql://postgres:postgres@localhost:5432/smart_alert_db?sslmode=disable"
fi

echo "Running migrations on: $DATABASE_URL"
echo ""

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Check if psql is available, otherwise use Go tool
if command -v psql &> /dev/null; then
    # Use psql if available
    for file in "$SCRIPT_DIR"/*.sql; do
        if [[ -f "$file" ]] && [[ "$(basename "$file")" != "README.md" ]] && [[ "$(basename "$file")" != "drop_all_tables.sql" ]]; then
            filename=$(basename "$file")
            echo "Running migration: $filename"
            psql "$DATABASE_URL" -f "$file"
            if [ $? -eq 0 ]; then
                echo "✓ Success: $filename"
            else
                echo "✗ Failed: $filename"
                exit 1
            fi
            echo ""
        fi
    done
    echo "All migrations completed successfully!"
else
    # Use Go tool if psql is not available
    echo "psql not found, using Go migration tool..."
    cd "$PROJECT_ROOT"
    go run cmd/migrate/main.go migrations
fi

