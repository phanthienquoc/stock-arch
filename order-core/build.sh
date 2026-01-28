#!/bin/bash
# Build script for order-core (Go)

set -e

echo "🔨 Building order-core (Go)..."
echo "================================"

cd "$(dirname "$0")"

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

# Build
echo "🏗️  Building executable..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-s -w" -o main .

# Verify
if [ -f main ]; then
    SIZE=$(ls -lh main | awk '{print $5}')
    echo "✅ Build successful! Binary size: $SIZE"
else
    echo "❌ Build failed"
    exit 1
fi

# Test
echo ""
echo "🧪 Testing binary..."
echo "./main" 
echo ""
echo "To run: docker compose -f docker-compose.yml build order-core"
