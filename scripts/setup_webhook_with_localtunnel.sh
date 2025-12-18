#!/bin/bash

# Script untuk setup webhook dengan localtunnel
# Usage: ./scripts/setup_webhook_with_localtunnel.sh [localtunnel_url] [port]

set -e

# Load .env file if exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Get port from parameter or .env
PORT="${2:-${APP_PORT:-8080}}"
LOCALTUNNEL_URL="${1}"

if [ -z "$LOCALTUNNEL_URL" ]; then
    echo "Usage: ./scripts/setup_webhook_with_localtunnel.sh <localtunnel_url> [port]"
    echo ""
    echo "Example:"
    echo "  ./scripts/setup_webhook_with_localtunnel.sh https://giant-ties-rule.loca.lt 8080"
    echo ""
    echo "Or start localtunnel first:"
    echo "  lt --port $PORT"
    echo "  # Copy the URL, then run this script with the URL"
    exit 1
fi

# Remove trailing slash
LOCALTUNNEL_URL="${LOCALTUNNEL_URL%/}"

# Build webhook URL
WEBHOOK_URL="${LOCALTUNNEL_URL}/webhook"

# Get Waha config
WAHA_URL="${WAHA_SERVER_URL:-http://localhost:3000}"
WAHA_URL="${WAHA_URL%/}"
API_KEY="${WAHA_API_KEY:-}"

echo "Setting up Waha webhook with localtunnel..."
echo "Localtunnel URL: $LOCALTUNNEL_URL"
echo "Webhook URL: $WEBHOOK_URL"
echo "Waha URL: $WAHA_URL"
echo "Port: $PORT"
if [ -n "$API_KEY" ]; then
    echo "API Key: *** (hidden)"
else
    echo "API Key: (not set)"
fi
echo ""

# Build curl command
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
echo "Important:"
echo "  1. Make sure Smart Alert System is running on port $PORT"
echo "  2. Keep localtunnel running: lt --port $PORT"
echo "  3. Test by sending a WhatsApp message"
echo ""
echo "To verify webhook status:"
echo "curl -X GET \"${WAHA_URL}/api/webhook\""

