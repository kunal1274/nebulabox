#!/bin/bash

# Install NebulaBox CLI using go install
# This makes nbx available globally via Go's installation system

set -e

echo "🚀 Installing NebulaBox CLI via Go..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo "   Please install Go 1.22+ from https://go.dev"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.22"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Error: Go version $GO_VERSION is too old"
    echo "   Required: Go $REQUIRED_VERSION or higher"
    exit 1
fi

# Get module path from go.mod
if [ -f "go.mod" ]; then
    MODULE_PATH=$(grep "^module" go.mod | awk '{print $2}')
else
    echo "❌ Error: go.mod not found"
    echo "   Please run this script from the project root"
    exit 1
fi

echo "📦 Module: $MODULE_PATH"
echo ""

# Install the binary
echo "Installing nebulabox..."
go install -ldflags "-X main.Version=$(git describe --tags --always 2>/dev/null || echo 'dev')" \
    "$MODULE_PATH/cmd/nebulabox@latest"

# Get Go bin directory
GO_BIN=$(go env GOPATH)/bin
if [ -z "$GO_BIN" ] || [ "$GO_BIN" = "/bin" ]; then
    GO_BIN="$HOME/go/bin"
fi

echo ""
echo "✅ Installation complete!"
echo ""

# Check if Go bin is in PATH
if [[ ":$PATH:" != *":$GO_BIN:"* ]]; then
    echo "⚠️  Warning: $GO_BIN is not in your PATH"
    echo ""
    echo "Add this to your ~/.bashrc or ~/.zshrc:"
    echo "  export PATH=\"\$PATH:$GO_BIN\""
    echo ""
    echo "Then run:"
    echo "  source ~/.bashrc  # or source ~/.zshrc"
    echo ""
    echo "Or use the command directly:"
    echo "  $GO_BIN/nebulabox version"
else
    echo "✅ You can now use:"
    echo "   nebulabox version"
    echo "   nebulabox ps"
    echo "   nebulabox --help"
fi

# Create nbx alias
if [ -f "$GO_BIN/nebulabox" ]; then
    ln -sf nebulabox "$GO_BIN/nbx" 2>/dev/null || cp "$GO_BIN/nebulabox" "$GO_BIN/nbx"
    chmod +x "$GO_BIN/nbx"
    echo ""
    echo "✅ Shortcut 'nbx' created at $GO_BIN/nbx"
    echo "   You can now use: nbx version"
fi

echo ""
echo "📝 Note: To update, run:"
echo "   go install $MODULE_PATH/cmd/nebulabox@latest"

