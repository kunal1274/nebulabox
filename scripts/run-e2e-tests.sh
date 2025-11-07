#!/bin/bash
# E2E Test Runner Script for NebulaBox

set -e

API_URL="${API_URL:-http://localhost:8081}"
START_SERVER="${START_SERVER:-false}"
SERVER_BINARY="${SERVER_BINARY:-./nebulabox-api}"

echo "🧪 NebulaBox E2E Test Runner"
echo "============================"
echo "API URL: $API_URL"
echo "Start Server: $START_SERVER"
echo ""

# Check if API server is already running
if [ "$START_SERVER" = "false" ]; then
    echo "📡 Checking API server availability..."
    if ! curl -s -f "$API_URL/api/auth/me" > /dev/null 2>&1; then
        echo "⚠️  API server not available at $API_URL"
        echo "💡 Start the server first or set START_SERVER=true"
        exit 1
    fi
    echo "✅ API server is available"
fi

# Run E2E tests
echo ""
echo "🚀 Running E2E tests..."
echo ""

cd "$(dirname "$0")/.."

if [ "$START_SERVER" = "true" ]; then
    # Build server if needed
    if [ ! -f "$SERVER_BINARY" ]; then
        echo "🔨 Building API server..."
        make build-api
    fi
    
    # Run tests with server startup
    go test -v -timeout 5m ./tests/... \
        -api-url="$API_URL" \
        -start-server=true \
        -server-binary="$SERVER_BINARY" \
        -wait-timeout=30s \
        -run "TestE2E"
else
    # Run tests against existing server
    go test -v -timeout 5m ./tests/... \
        -api-url="$API_URL" \
        -run "TestE2E"
fi

TEST_EXIT_CODE=$?

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo ""
    echo "✅ All E2E tests passed!"
else
    echo ""
    echo "❌ Some E2E tests failed"
fi

exit $TEST_EXIT_CODE

