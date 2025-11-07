# Quick Mac Testing Guide

## Current Situation

✅ Release `v0.1.0` exists on GitHub  
❌ Binaries not uploaded to release  
✅ You can still test on Mac!

## Fastest Way to Test on Mac (Right Now)

### Step 1: On Your Mac, Clone and Build

```bash
# Clone the repository
git clone https://github.com/kunal1274/nebulabox.git
cd nebulabox

# Build for Mac (auto-detects your architecture)
go build -o nbx ./cmd/nebulabox

# Install globally
sudo mv nbx /usr/local/bin/

# Test it works
nbx version
```

**That's it!** You can now test all NebulaBox commands.

### Step 2: Test Basic Commands

```bash
# Check version
nbx version

# See help
nbx --help

# List containers
nbx ps

# List images
nbx images
```

## Alternative: Quick Install Script

On your Mac, create `install-nbx.sh`:

```bash
#!/bin/bash
set -e

echo "🍎 Installing NebulaBox on Mac..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    echo "Install Go: https://golang.org/dl/"
    exit 1
fi

# Clone and build
git clone https://github.com/kunal1274/nebulabox.git /tmp/nebulabox
cd /tmp/nebulabox

# Build
go build -o nbx ./cmd/nebulabox

# Install
sudo mv nbx /usr/local/bin/

# Cleanup
rm -rf /tmp/nebulabox

echo "✅ NebulaBox installed!"
echo "Test: nbx version"
```

**Run it:**
```bash
chmod +x install-nbx.sh
./install-nbx.sh
```

## What to Test on Mac

### 1. Basic Commands
- [ ] `nbx version` - Shows version
- [ ] `nbx --help` - Shows help
- [ ] `nbx ps` - Lists containers
- [ ] `nbx images` - Lists images

### 2. Container Operations
- [ ] `nbx run <image>` - Run a container
- [ ] `nbx ps` - See running containers
- [ ] `nbx logs <container>` - View logs
- [ ] `nbx stop <container>` - Stop container

### 3. Image Operations
- [ ] `nbx build <path>` - Build an image
- [ ] `nbx images` - List images
- [ ] `nbx rmi <image>` - Remove image

### 4. Advanced Features
- [ ] `nbx hierarchy create` - Create nested container
- [ ] `nbx hierarchy list` - List hierarchy
- [ ] `nbx group create` - Create container group

## Troubleshooting

### "Command not found: nbx"

```bash
# Check if installed
which nbx
ls -la /usr/local/bin/nbx

# If not found, rebuild and install
cd nebulabox
go build -o nbx ./cmd/nebulabox
sudo mv nbx /usr/local/bin/
```

### "Permission denied"

```bash
# Make executable
chmod +x nbx
sudo mv nbx /usr/local/bin/
```

### macOS Security Warning

If macOS blocks the binary:

```bash
# Remove quarantine
xattr -d com.apple.quarantine /usr/local/bin/nbx
```

## Next: Fix Release Binaries

After testing, you can:
1. Build all platform binaries: `./scripts/build-release.sh v0.1.0`
2. Upload to GitHub release manually
3. Then others can download binaries instead of building

---

**Start testing on Mac now!** 🍎✨

