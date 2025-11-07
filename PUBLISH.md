# Publishing NebulaBox CLI (nbx)

## Overview

This guide explains how to publish NebulaBox CLI so it can be installed from anywhere in the world.

## Prerequisites

1. **GitHub Repository**: Your code should be in a GitHub repository
2. **Go Module**: Project must be a valid Go module
3. **Git Tags**: Use semantic versioning (v0.1.0, v1.0.0, etc.)

## Publishing Methods

### Method 1: Go Install (Recommended for Developers)

This allows users to install directly from your GitHub repository:

```bash
go install github.com/yourusername/nebulabox/cmd/nebulabox@latest
```

**Setup Steps:**

1. **Ensure your module path is correct in go.mod:**
   ```bash
   module github.com/yourusername/nebulabox
   ```

2. **Push your code to GitHub:**
   ```bash
   git remote add origin https://github.com/yourusername/nebulabox.git
   git push -u origin main
   ```

3. **Create a release tag:**
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

4. **Users can now install:**
   ```bash
   go install github.com/yourusername/nebulabox/cmd/nebulabox@latest
   ```

### Method 2: GitHub Releases (Recommended for End Users)

This provides pre-built binaries for all platforms.

**Setup Steps:**

1. **Build release binaries:**
   ```bash
   ./scripts/build-release.sh v0.1.0
   ```

2. **Create a GitHub release:**
   - Go to GitHub → Releases → Draft a new release
   - Tag: `v0.1.0`
   - Title: `Release v0.1.0`
   - Upload all files from `dist/` directory
   - Publish release

3. **Users can download:**
   ```bash
   wget https://github.com/yourusername/nebulabox/releases/download/v0.1.0/nbx-linux-amd64
   chmod +x nbx-linux-amd64
   sudo mv nbx-linux-amd64 /usr/local/bin/nbx
   ```

### Method 3: Automated GitHub Actions

Use the provided workflow to automatically build and release:

1. **Push workflow file** (already created: `.github/workflows/release.yml`)

2. **Create a tag:**
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

3. **GitHub Actions will automatically:**
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload all binaries

## Version Management

### Semantic Versioning

Use semantic versioning: `MAJOR.MINOR.PATCH`

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### Creating a New Release

```bash
# 1. Update version in code/docs
# 2. Commit changes
git add .
git commit -m "Prepare release v0.1.0"

# 3. Create and push tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# 4. Build and upload (if manual)
./scripts/build-release.sh v0.1.0
# Then upload dist/* to GitHub Releases
```

## Installation Instructions for Users

### Via Go Install

```bash
# Install latest
go install github.com/yourusername/nebulabox/cmd/nebulabox@latest

# Install specific version
go install github.com/yourusername/nebulabox/cmd/nebulabox@v0.1.0

# Create nbx alias
alias nbx=nebulabox
# Or: ln -s $(go env GOPATH)/bin/nebulabox $(go env GOPATH)/bin/nbx
```

### Via GitHub Releases

```bash
# Linux
wget https://github.com/yourusername/nebulabox/releases/download/v0.1.0/nbx-linux-amd64
chmod +x nbx-linux-amd64
sudo mv nbx-linux-amd64 /usr/local/bin/nbx

# macOS (Intel)
wget https://github.com/yourusername/nebulabox/releases/download/v0.1.0/nbx-darwin-amd64
chmod +x nbx-darwin-amd64
sudo mv nbx-darwin-amd64 /usr/local/bin/nbx

# macOS (Apple Silicon)
wget https://github.com/yourusername/nebulabox/releases/download/v0.1.0/nbx-darwin-arm64
chmod +x nbx-darwin-arm64
sudo mv nbx-darwin-arm64 /usr/local/bin/nbx

# Windows
# Download nbx-windows-amd64.exe and add to PATH
```

## Testing Installation

Before publishing, test installation:

```bash
# Test go install locally
go install ./cmd/nebulabox
$(go env GOPATH)/bin/nebulabox version

# Test release build
./scripts/build-release.sh v0.1.0-test
# Test binaries in dist/
```

## Distribution Checklist

- [ ] Code pushed to GitHub
- [ ] go.mod has correct module path
- [ ] Version tags created
- [ ] Release binaries built
- [ ] GitHub release created
- [ ] Installation instructions documented
- [ ] README updated with installation steps
- [ ] Version information displays correctly

## Future: Package Managers

### Homebrew (macOS)

Create a Homebrew formula:
```ruby
class Nebulabox < Formula
  desc "NebulaBox container platform"
  homepage "https://github.com/yourusername/nebulabox"
  url "https://github.com/yourusername/nebulabox/archive/v0.1.0.tar.gz"
  sha256 "..."
  
  depends_on "go" => :build
  
  def install
    system "go", "build", "-o", bin/"nbx", "./cmd/nebulabox"
  end
end
```

### Snap (Linux)

Create `snap/snapcraft.yaml`:
```yaml
name: nbx
version: '0.1.0'
summary: NebulaBox container platform
description: |
  NebulaBox is a modern container platform.
  
grade: stable
confinement: strict

apps:
  nbx:
    command: nbx

parts:
  nbx:
    plugin: go
    source: .
```

## Quick Publish Commands

```bash
# Full release process
VERSION="v0.1.0"
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION
./scripts/build-release.sh $VERSION
# Then upload dist/* to GitHub Releases manually or via GitHub Actions
```

---

**Your CLI is now ready to be published and installed from anywhere!** 🚀

