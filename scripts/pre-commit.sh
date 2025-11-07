#!/bin/bash

# NebulaBox Pre-commit Hook
# Run tests before allowing commits

set -e

echo "🔍 Running pre-commit checks..."

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Skipping Go checks."
    exit 0
fi

# Run linter (if available)
if command -v golangci-lint &> /dev/null; then
    echo "🔍 Running golangci-lint..."
    golangci-lint run --timeout=5m ./... || echo "⚠️  Linter warnings found (non-blocking)"
fi

# Run unit tests
echo "🧪 Running unit tests..."
if ! make test-unit > /tmp/test-output.log 2>&1; then
    echo "❌ Unit tests failed!"
    cat /tmp/test-output.log
    exit 1
fi

# Format check
echo "📝 Checking code formatting..."
if ! gofmt -l . | grep -q .; then
    echo "⚠️  Some files are not formatted. Run 'go fmt ./...' to fix."
    gofmt -l . | head -10
fi

# Check for common issues
echo "🔍 Checking for common issues..."

# Check for TODO/FIXME comments in committed code (non-blocking)
if git diff --cached | grep -qE "(TODO|FIXME|XXX)"; then
    echo "⚠️  Found TODO/FIXME comments in staged changes (non-blocking)"
fi

echo "✅ Pre-commit checks passed!"
exit 0

