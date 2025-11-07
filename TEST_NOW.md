# ✅ Ready to Test - Installation Instructions

## Confirmation

✅ **Release is ready!** All binaries are uploaded and available for download.

## Test on Remote Linux System

### Quick Install & Test

```bash
# 1. Download binary
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64

# 2. Make executable
chmod +x nbx-linux-amd64

# 3. Install globally
sudo mv nbx-linux-amd64 /usr/local/bin/nbx

# 4. Test it works
nbx version
nbx --help
```

### Test All Commands

```bash
# Version
nbx version

# Help
nbx --help

# List containers
nbx ps

# List images
nbx images

# Build an image (if you have a BuildSpec)
nbx build /path/to/buildspec.json

# Run a container
nbx run <image-name>
```

## Test on Mac (CLI Only)

**Note:** Container features won't work on Mac (Linux-specific), but CLI structure can be tested.

### Option 1: Download Binary (if available)
```bash
# For Intel Mac
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-darwin-amd64

# For Apple Silicon
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-darwin-arm64
```

**Note:** Mac binaries aren't built yet (Linux-only code). Use Option 2 instead.

### Option 2: Build from Source (CLI Only)
```bash
# Clone
git clone https://github.com/kunal1274/nebulabox.git
cd nebulabox

# Build (will work for CLI, not containers)
go build -o nbx ./cmd/nebulabox

# Install
sudo mv nbx /usr/local/bin/

# Test CLI
nbx version
nbx --help
```

## Test Checklist

### Basic Tests
- [ ] Download binary successfully
- [ ] Install to PATH
- [ ] Run `nbx version` - shows version
- [ ] Run `nbx --help` - shows help menu
- [ ] Run `nbx ps` - lists containers (empty initially)
- [ ] Run `nbx images` - lists images (empty initially)

### Container Tests (Linux Only)
- [ ] Build an image: `nbx build <path>`
- [ ] Run a container: `nbx run <image>`
- [ ] List containers: `nbx ps`
- [ ] View logs: `nbx logs <container>`
- [ ] Stop container: `nbx stop <container>`

### Advanced Tests (Linux Only)
- [ ] Create nested container: `nbx hierarchy create`
- [ ] List hierarchy: `nbx hierarchy list`
- [ ] Create group: `nbx group create`

## Quick Test Script

Save this as `test-install.sh` on remote system:

```bash
#!/bin/bash
set -e

echo "🧪 Testing NebulaBox Installation"
echo ""

# Download
echo "1️⃣ Downloading binary..."
wget -q https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64 -O /tmp/nbx
chmod +x /tmp/nbx

# Test
echo "2️⃣ Testing binary..."
/tmp/nbx version

if [ $? -eq 0 ]; then
    echo "✅ Binary works!"
    
    # Install
    echo ""
    echo "3️⃣ Installing..."
    sudo mv /tmp/nbx /usr/local/bin/nbx
    
    echo ""
    echo "✅ Installation successful!"
    echo ""
    echo "Test commands:"
    nbx version
    nbx --help
else
    echo "❌ Binary test failed"
    exit 1
fi
```

**Run it:**
```bash
chmod +x test-install.sh
./test-install.sh
```

## Verify Installation

```bash
# Check if installed
which nbx

# Check version
nbx version

# Check help
nbx --help

# List commands
nbx ps
nbx images
```

## Troubleshooting

### "Command not found"
```bash
# Check PATH
echo $PATH

# Add to PATH if needed
export PATH=$PATH:/usr/local/bin
```

### "Permission denied"
```bash
# Make executable
chmod +x nbx-linux-amd64

# Or install with sudo
sudo mv nbx-linux-amd64 /usr/local/bin/nbx
```

### Download fails
```bash
# Check URL
curl -I https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64

# Try with curl instead
curl -L https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64 -o nbx
```

---

**✅ Everything is ready! Start testing now!** 🚀

