#!/bin/bash

# Create a new release with proper Git configuration checks
# This script ensures Git is configured before creating a release

set -e

# Check Git configuration
if ! git config --get user.name > /dev/null 2>&1; then
    echo "❌ Error: Git user.name not configured"
    echo ""
    echo "Please run:"
    echo "  ./scripts/setup-git.sh"
    echo ""
    echo "Or configure manually:"
    echo "  git config user.name 'Your Name'"
    echo "  git config user.email 'your.email@example.com'"
    exit 1
fi

if ! git config --get user.email > /dev/null 2>&1; then
    echo "❌ Error: Git user.email not configured"
    echo ""
    echo "Please run:"
    echo "  ./scripts/setup-git.sh"
    echo ""
    echo "Or configure manually:"
    echo "  git config user.name 'Your Name'"
    echo "  git config user.email 'your.email@example.com'"
    exit 1
fi

# Get version
if [ -z "$1" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v0.1.0"
    exit 1
fi

VERSION="$1"

# Validate version format (should start with 'v')
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+ ]]; then
    echo "⚠️  Warning: Version should follow semantic versioning (e.g., v0.1.0)"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo "🚀 Creating release: $VERSION"
echo ""

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo "❌ Error: Tag $VERSION already exists"
    echo ""
    echo "To update the tag:"
    echo "  git tag -d $VERSION"
    echo "  git push origin :refs/tags/$VERSION"
    exit 1
fi

# Check if working directory is clean
if ! git diff-index --quiet HEAD --; then
    echo "⚠️  Warning: Working directory has uncommitted changes"
    echo ""
    git status --short
    echo ""
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Build release binaries
echo "📦 Building release binaries..."
./scripts/build-release.sh "$VERSION"

# Create annotated tag
echo ""
echo "🏷️  Creating Git tag: $VERSION"
git tag -a "$VERSION" -m "Release $VERSION"

# Show tag info
echo ""
echo "✅ Tag created:"
git show "$VERSION" --no-patch --format="%H%n%an <%ae>%n%s" | head -3

echo ""
echo "📋 Next steps:"
echo "  1. Push the tag:"
echo "     git push origin $VERSION"
echo ""
echo "  2. If using GitHub Actions, the release will be created automatically"
echo ""
echo "  3. Or manually create GitHub release:"
echo "     - Go to: https://github.com/nebulabox/nebulabox/releases/new"
echo "     - Select tag: $VERSION"
echo "     - Upload files from: dist/"
echo ""
read -p "Push tag now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Pushing tag..."
    git push origin "$VERSION"
    echo "✅ Tag pushed!"
    echo ""
    echo "🎉 Release $VERSION is now available!"
    echo ""
    echo "Users can install with:"
    echo "  go install github.com/nebulabox/nebulabox/cmd/nebulabox@$VERSION"
else
    echo ""
    echo "Tag created locally. Push when ready:"
    echo "  git push origin $VERSION"
fi

