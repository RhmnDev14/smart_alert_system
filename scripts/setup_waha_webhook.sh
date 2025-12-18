#!/bin/bash

# Script untuk setup webhook Waha Server
# Usage: ./scripts/setup_waha_webhook.sh [webhook_url] [waha_url] [api_key]
# Jika tidak ada parameter, akan membaca dari .env

set -e

# Load .env file if exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Get values from parameters or .env or defaults
WEBHOOK_URL="${1:-${WEBHOOK_URL:-http://localhost:8080/webhook}}"
WAHA_URL="${2:-${WAHA_SERVER_URL:-http://localhost:3000}}"
API_KEY="${3:-${WAHA_API_KEY:-}}"

# Get APP_PORT from .env if available
APP_PORT="${APP_PORT:-8080}"

# Remove trailing slash from WAHA_URL
WAHA_URL="${WAHA_URL%/}"

# If webhook URL is default and APP_PORT is set, use it
if [ "$WEBHOOK_URL" = "http://localhost:8080/webhook" ] && [ "$APP_PORT" != "8080" ]; then
    WEBHOOK_URL="http://localhost:${APP_PORT}/webhook"
fi

# Check if Waha is remote but webhook is localhost
if [[ "$WAHA_URL" == http*://*.*.* ]] && [[ "$WEBHOOK_URL" == http://localhost* ]]; then
    echo "⚠️  WARNING: Waha Server is remote but webhook URL is localhost!"
    echo "   Waha Server: $WAHA_URL"
    echo "   Webhook URL: $WEBHOOK_URL"
    echo ""
    echo "   This won't work! You need to use a public URL for webhook."
    echo "   Options:"
    echo "   1. Use ngrok: ngrok http $APP_PORT"
    echo "   2. Use your public IP/domain"
    echo "   3. Deploy Smart Alert System to a public server"
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Cancelled."
        exit 1
    fi
    echo ""
fi

echo "Setting up Waha webhook..."
echo "Webhook URL: $WEBHOOK_URL"
echo "Waha URL: $WAHA_URL"
if [ -n "$API_KEY" ]; then
    echo "API Key: *** (hidden)"
else
    echo "API Key: (not set)"
fi
echo ""

# Check if using localtunnel URL
if [[ "$WEBHOOK_URL" == *".loca.lt"* ]]; then
    echo "✓ Detected localtunnel URL"
    echo "  Make sure localtunnel is running: lt --port ${APP_PORT:-8080}"
    echo ""
fi

# Build curl command (ensure no double slashes)
WEBHOOK_API_URL="${WAHA_URL}/api/webhook"
CURL_CMD="curl -X POST \"${WEBHOOK_API_URL}\" \
  -H \"Content-Type: application/json\""

if [ -n "$API_KEY" ]; then
  CURL_CMD="$CURL_CMD -H \"X-Api-Key: ${API_KEY}\""
fi

CURL_CMD="$CURL_CMD -d '{
  \"url\": \"${WEBHOOK_URL}\",
  \"events\": [\"message\"]
}'"

echo "Executing: $CURL_CMD"
echo ""

# Execute
eval $CURL_CMD

echo ""
echo "✓ Webhook setup completed!"
echo ""
echo "To verify, check webhook status:"
echo "curl -X GET \"${WAHA_URL}/api/webhook\""

