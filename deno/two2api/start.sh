#!/bin/bash

# Start OpenAI-compatible API proxy
# Usage: ./start.sh

echo "🤖 Starting API Proxy..."
echo ""

# Load environment variables from .env file if it exists
if [ -f .env ]; then
    echo "Loading environment from .env"
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if port is already in use
PORT=${PORT:-9090}
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  Port $PORT is already in use!"
    echo "Please stop the existing process or use a different port:"
    echo "  PORT=9091 ./start.sh"
    exit 1
fi

# Display configuration
echo "Configuration:"
echo "  Port: ${PORT:-9090}"
echo "  API Key: ${DEFAULT_KEY:-sk-your-key}"
echo "  Model: ${MODEL_NAME:-default-model}"
echo "  Service: ${SERVICE_NAME:-AI2Api}"
echo "  Debug: ${DEBUG_MODE:-false}"
echo ""

# Start the server
deno run --allow-net --allow-env main.ts
