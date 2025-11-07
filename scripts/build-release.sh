#!/bin/bash

# Build release binaries for multiple platforms
# This creates distributable binaries for GitHub Releases

set -e

# Check Git configuration
if ! git config --get user.name > /dev/null 2>&1 || ! git config --get user.email > /dev/null 2>&1; then
    echo "⚠️  Warning: Git user.name and user.email not configured"
    echo ""
    echo "Please configure Git:"
    echo "  git config --global user.name 'Your Name'"
    echo "  git config --global user.email 'your.email@example.com'"
    echo ""
    echo "Or configure for this repository only:"
    echo "  git config user.name 'Your Name'"
    echo "  git config user.email 'your.email@example.com'"
    echo ""
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

VERSION="${1:-$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')}"
BUILD_DIR="dist"
BINARY_NAME="nebulabox"

echo "🚀 Building NebulaBox release binaries"
echo "Version: $VERSION"
echo ""

# Clean previous builds
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Platforms to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Build for each platform
for PLATFORM in "${PLATFORMS[@]}"; do
    PLATFORM_SPLIT=(${PLATFORM//\// })
    GOOS=${PLATFORM_SPLIT[0]}
    GOARCH=${PLATFORM_SPLIT[1]}
    
    OUTPUT_NAME="$BINARY_NAME-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="$OUTPUT_NAME.exe"
    fi
    
    echo "📦 Building for $GOOS/$GOARCH..."
    
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
        -o "$BUILD_DIR/$OUTPUT_NAME" \
        ./cmd/nebulabox
    
    # Create checksum
    if [ "$GOOS" = "windows" ]; then
        sha256sum "$BUILD_DIR/$OUTPUT_NAME" > "$BUILD_DIR/$OUTPUT_NAME.sha256"
    else
        shasum -a 256 "$BUILD_DIR/$OUTPUT_NAME" > "$BUILD_DIR/$OUTPUT_NAME.sha256"
    fi
    
    echo "✅ Built: $OUTPUT_NAME"
done

# Create nbx symlinks/aliases for Unix systems
for PLATFORM in "linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64"; do
    PLATFORM_SPLIT=(${PLATFORM//\// })
    GOOS=${PLATFORM_SPLIT[0]}
    GOARCH=${PLATFORM_SPLIT[1]}
    
    BINARY="$BUILD_DIR/$BINARY_NAME-$GOOS-$GOARCH"
    if [ -f "$BINARY" ]; then
        # Create nbx as a copy (since symlinks don't work in releases)
        cp "$BINARY" "$BUILD_DIR/nbx-$GOOS-$GOARCH"
        chmod +x "$BUILD_DIR/nbx-$GOOS-$GOARCH"
    fi
done

# Create release notes
cat > "$BUILD_DIR/RELEASE_NOTES.md" <<EOF
# NebulaBox CLI Release $VERSION

## Installation

### Linux (amd64)
\`\`\`bash
wget https://github.com/nebulabox/nebulabox/releases/download/$VERSION/nbx-linux-amd64
chmod +x nbx-linux-amd64
sudo mv nbx-linux-amd64 /usr/local/bin/nbx
\`\`\`

### macOS (amd64)
\`\`\`bash
wget https://github.com/nebulabox/nebulabox/releases/download/$VERSION/nbx-darwin-amd64
chmod +x nbx-darwin-amd64
sudo mv nbx-darwin-amd64 /usr/local/bin/nbx
\`\`\`

### macOS (Apple Silicon)
\`\`\`bash
wget https://github.com/nebulabox/nebulabox/releases/download/$VERSION/nbx-darwin-arm64
chmod +x nbx-darwin-arm64
sudo mv nbx-darwin-arm64 /usr/local/bin/nbx
\`\`\`

### Windows
Download \`nbx-windows-amd64.exe\` and add to PATH.

## Verify Installation

\`\`\`bash
nbx version
nbx --help
\`\`\`

## Checksums

All binaries include SHA256 checksums for verification.
EOF

echo ""
echo "✅ Release build complete!"
echo "📦 Binaries in: $BUILD_DIR/"
echo ""
echo "To create a GitHub release:"
echo "  1. Create a tag: git tag -a v$VERSION -m 'Release $VERSION'"
echo "  2. Push tag: git push origin v$VERSION"
echo "  3. Upload files from $BUILD_DIR/ to GitHub Releases"
echo ""

