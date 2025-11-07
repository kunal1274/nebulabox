# Install NebulaBox on macOS (Remote System)

## Quick Installation for Mac

### Method 1: Go Install (Recommended if Go is installed)

```bash
# Install latest version
go install github.com/kunal1274/nebulabox/cmd/nebulabox@latest

# Add to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin

# Verify installation
nbx version
```

**If `nbx` command not found:**
```bash
# Add Go bin to PATH permanently
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# Or for bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile
source ~/.bash_profile
```

### Method 2: Download Binary from GitHub Release

#### For Intel Mac (amd64):
```bash
# Download binary
curl -L https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-darwin-amd64 -o nbx

# Make executable
chmod +x nbx

# Move to system path
sudo mv nbx /usr/local/bin/

# Verify
nbx version
```

#### For Apple Silicon Mac (arm64/M1/M2/M3):
```bash
# Download binary
curl -L https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-darwin-arm64 -o nbx

# Make executable
chmod +x nbx

# Move to system path
sudo mv nbx /usr/local/bin/

# Verify
nbx version
```

### Method 3: Install via Homebrew (if formula exists)

```bash
# If you create a Homebrew formula later
brew install nebulabox
```

### Method 4: Build from Source

If binaries aren't available:

```bash
# Clone repository
git clone https://github.com/kunal1274/nebulabox.git
cd nebulabox

# Build
go build -o nbx ./cmd/nebulabox

# Install
sudo mv nbx /usr/local/bin/

# Verify
nbx version
```

## Detect Your Mac Architecture

```bash
# Check architecture
uname -m

# Output:
# x86_64 = Intel Mac (use nbx-darwin-amd64)
# arm64 = Apple Silicon (use nbx-darwin-arm64)
```

## Quick Test Script for Mac

Save this as `install-nbx-mac.sh`:

```bash
#!/bin/bash
set -e

echo "🍎 Installing NebulaBox on macOS"
echo ""

# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" == "x86_64" ]; then
    BINARY="nbx-darwin-amd64"
    echo "Detected: Intel Mac (amd64)"
elif [ "$ARCH" == "arm64" ]; then
    BINARY="nbx-darwin-arm64"
    echo "Detected: Apple Silicon (arm64)"
else
    echo "❌ Unknown architecture: $ARCH"
    exit 1
fi

VERSION="v0.1.0"
URL="https://github.com/kunal1274/nebulabox/releases/download/${VERSION}/${BINARY}"

echo "Downloading: $URL"
echo ""

# Download
curl -L -f -o /tmp/nbx "$URL"

if [ $? -eq 0 ]; then
    chmod +x /tmp/nbx
    
    # Test first
    echo "Testing binary..."
    /tmp/nbx version
    
    if [ $? -eq 0 ]; then
        # Install
        echo ""
        echo "Installing to /usr/local/bin/nbx..."
        sudo mv /tmp/nbx /usr/local/bin/nbx
        
        echo ""
        echo "✅ Installation successful!"
        echo ""
        echo "Verify:"
        nbx version
        nbx --help
    else
        echo "❌ Binary test failed"
        exit 1
    fi
else
    echo "❌ Download failed"
    echo ""
    echo "Trying Go install instead..."
    if command -v go &> /dev/null; then
        go install github.com/kunal1274/nebulabox/cmd/nebulabox@latest
        echo "✅ Installed via Go"
    else
        echo "❌ Go not installed. Please install Go first."
        exit 1
    fi
fi
```

**Usage:**
```bash
chmod +x install-nbx-mac.sh
./install-nbx-mac.sh
```

## Troubleshooting

### Issue: "Command not found: nbx"

**Solution:**
```bash
# Check if binary exists
which nbx
ls -la /usr/local/bin/nbx

# If not found, check PATH
echo $PATH

# Add to PATH if needed
export PATH=$PATH:/usr/local/bin
```

### Issue: "Permission denied"

**Solution:**
```bash
# Make executable
chmod +x nbx

# Or install with sudo
sudo mv nbx /usr/local/bin/
```

### Issue: "Binary not found on GitHub"

**Possible causes:**
- Release doesn't have binaries yet
- GitHub Actions workflow didn't run
- Wrong version tag

**Solutions:**
1. Check release: `https://github.com/kunal1274/nebulabox/releases`
2. Build from source (Method 4 above)
3. Use Go install: `go install github.com/kunal1274/nebulabox/cmd/nebulabox@latest`

### Issue: macOS Security Warning

If you see "nbx cannot be opened because it is from an unidentified developer":

**Solution:**
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/nbx

# Or allow in System Preferences:
# System Preferences > Security & Privacy > Allow
```

## Verify Installation

```bash
# Check version
nbx version

# Check help
nbx --help

# List commands
nbx ps
nbx images
```

## Next Steps After Installation

1. ✅ Test basic commands
2. ✅ Try building an image
3. ✅ Run a container
4. ✅ Test hierarchical containers

---

**Ready to test on Mac!** 🍎

