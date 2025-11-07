#!/bin/bash

# Install NebulaBox CLI (nbx) to system PATH
# This script creates symlinks for 'nbx' and 'nebulabox' commands

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_DIR="$PROJECT_ROOT/bin"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

echo "🚀 Installing NebulaBox CLI..."

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Check if binary exists
if [ ! -f "$BIN_DIR/nebulabox" ]; then
    echo "❌ Error: Binary not found at $BIN_DIR/nebulabox"
    echo "   Please run 'make build-cli-test' first"
    exit 1
fi

# Create symlinks
echo "📦 Creating symlinks..."
ln -sf "$BIN_DIR/nebulabox" "$INSTALL_DIR/nbx"
ln -sf "$BIN_DIR/nebulabox" "$INSTALL_DIR/nebulabox"

# Make executable
chmod +x "$INSTALL_DIR/nbx"
chmod +x "$INSTALL_DIR/nebulabox"

echo "✅ Symlinks created:"
echo "   $INSTALL_DIR/nbx"
echo "   $INSTALL_DIR/nebulabox"

# Check if install directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "⚠️  Warning: $INSTALL_DIR is not in your PATH"
    echo ""
    echo "Add this to your ~/.bashrc or ~/.zshrc:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    echo "Then run:"
    echo "  source ~/.bashrc  # or source ~/.zshrc"
    echo ""
    echo "Or use the command directly:"
    echo "  $INSTALL_DIR/nbx version"
else
    echo ""
    echo "✅ Installation complete! You can now use:"
    echo "   nbx version"
    echo "   nbx ps"
    echo "   nbx --help"
fi

