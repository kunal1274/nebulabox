# Test NebulaBox Installation from Remote System

## Prerequisites

1. ✅ Code pushed to GitHub
2. ✅ Tag pushed to GitHub
3. ✅ GitHub Actions workflow should create release automatically

## Check Release Status

### Option 1: Check GitHub Web UI
1. Go to: `https://github.com/kunal1274/nebulabox/releases`
2. Check if release was created automatically
3. Verify binaries are uploaded

### Option 2: Check via GitHub CLI (if installed)
```bash
gh release list
```

### Option 3: Check GitHub Actions
1. Go to: `https://github.com/kunal1274/nebulabox/actions`
2. Check if "Release" workflow ran successfully
3. If it failed, check the logs

## Installation Methods to Test

### Method 1: Go Install (Recommended)

On a remote system with Go installed:

```bash
# Install latest version
go install github.com/kunal1274/nebulabox/cmd/nebulabox@latest

# Or install specific version
go install github.com/kunal1274/nebulabox/cmd/nebulabox@v0.1.0

# Verify installation
nbx version
# or
nebulabox version
```

**Note**: For `go install` to work, the repository must be public or you need authentication.

### Method 2: Download Binary from GitHub Release

On any system (Linux/macOS/Windows):

```bash
# Linux (amd64)
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64
chmod +x nbx-linux-amd64
sudo mv nbx-linux-amd64 /usr/local/bin/nbx

# macOS (amd64)
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-darwin-amd64
chmod +x nbx-darwin-amd64
sudo mv nbx-darwin-amd64 /usr/local/bin/nbx

# macOS (Apple Silicon)
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-darwin-arm64
chmod +x nbx-darwin-arm64
sudo mv nbx-darwin-arm64 /usr/local/bin/nbx

# Windows
# Download nbx-windows-amd64.exe and add to PATH
```

### Method 3: Using curl

```bash
# Linux
curl -L https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64 -o nbx
chmod +x nbx
sudo mv nbx /usr/local/bin/

# Verify
nbx version
```

## Testing Checklist

### On Remote System

- [ ] **Test Go Install**
  ```bash
  go install github.com/kunal1274/nebulabox/cmd/nebulabox@latest
  nbx version
  ```

- [ ] **Test Binary Download**
  ```bash
  wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64
  chmod +x nbx-linux-amd64
  ./nbx-linux-amd64 version
  ```

- [ ] **Test CLI Commands**
  ```bash
  nbx --help
  nbx version
  nbx ps
  ```

- [ ] **Test Different Platforms**
  - [ ] Linux (amd64)
  - [ ] Linux (arm64)
  - [ ] macOS (amd64)
  - [ ] macOS (arm64)
  - [ ] Windows (amd64)

## Troubleshooting

### Issue: Release Not Created Automatically

**Solution**: Manually create release on GitHub:
1. Go to: `https://github.com/kunal1274/nebulabox/releases/new`
2. Select tag: `v0.1.0`
3. Upload binaries from `dist/` directory
4. Publish release

### Issue: `go install` Fails

**Possible causes**:
- Repository is private (needs authentication)
- Go module path is incorrect
- Network issues

**Solutions**:
- Make repository public, OR
- Use GitHub token: `GOPRIVATE=github.com/kunal1274/nebulabox go install ...`
- Use binary download instead

### Issue: Binary Not Found

**Check**:
- Tag name matches release
- Release is published (not draft)
- Binary filename is correct

### Issue: Permission Denied

**Solution**:
```bash
chmod +x nbx
sudo mv nbx /usr/local/bin/
```

## Quick Test Script

Save this as `test-install.sh` on remote system:

```bash
#!/bin/bash
set -e

echo "🧪 Testing NebulaBox Installation"
echo ""

# Test 1: Go Install
echo "1️⃣ Testing Go Install..."
if command -v go &> /dev/null; then
    go install github.com/kunal1274/nebulabox/cmd/nebulabox@latest
    if command -v nbx &> /dev/null || command -v nebulabox &> /dev/null; then
        echo "✅ Go install successful"
        nbx version || nebulabox version
    else
        echo "⚠️  Binary installed but not in PATH"
        echo "   Add ~/go/bin to PATH"
    fi
else
    echo "⏭️  Go not installed, skipping"
fi

echo ""

# Test 2: Binary Download
echo "2️⃣ Testing Binary Download..."
VERSION="v0.1.0"
ARCH="linux-amd64"  # Change for your system

if [ "$(uname)" == "Darwin" ]; then
    ARCH="darwin-$(uname -m | sed 's/x86_64/amd64/;s/arm64/arm64/')"
fi

URL="https://github.com/kunal1274/nebulabox/releases/download/${VERSION}/nbx-${ARCH}"
echo "Downloading: $URL"

if curl -L -f -o /tmp/nbx-test "$URL" 2>/dev/null; then
    chmod +x /tmp/nbx-test
    echo "✅ Download successful"
    /tmp/nbx-test version
    rm /tmp/nbx-test
else
    echo "❌ Download failed - check release exists"
fi

echo ""
echo "✅ Installation test complete!"
```

## Next Steps After Successful Test

1. ✅ Update `INSTALL.md` with actual installation commands
2. ✅ Update `README.md` with installation instructions
3. ✅ Test on multiple platforms
4. ✅ Document any issues found

---

**Ready to test!** 🚀

