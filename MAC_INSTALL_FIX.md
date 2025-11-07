# Mac Installation Fix

## Problem

`npm install -g nebulabox` fails on Mac because:
1. **Mac binaries don't exist yet** - Only Linux binaries are in the release
2. **Code is Linux-specific** - Uses Linux syscalls that don't work on Mac
3. **Install script tries to download** - But binary doesn't exist

## Current Status

✅ **Linux**: Works (binaries available)  
❌ **Mac**: Fails (no binaries, code is Linux-only)  
❌ **Windows**: Fails (no binaries, code is Linux-only)

## Solutions

### Option 1: Build from Source (CLI Only)

Mac can build the CLI structure, but container features won't work:

```bash
# Clone repository
git clone https://github.com/kunal1274/nebulabox.git
cd nebulabox

# Build (will work for CLI commands)
go build -o nbx ./cmd/nebulabox

# Install
sudo mv nbx /usr/local/bin/

# Test CLI
nbx version
nbx --help
```

**Note**: Container operations won't work because code uses Linux-specific syscalls.

### Option 2: Use Linux VM/Container

Run NebulaBox in a Linux environment on Mac:

```bash
# Using Docker
docker run -it --rm -v $(pwd):/workspace ubuntu:22.04 bash

# Inside container
apt-get update && apt-get install -y wget
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64
chmod +x nbx-linux-amd64
./nbx-linux-amd64 version
```

### Option 3: Fix Install Script (Better Error Messages)

The install script now:
- ✅ Detects Mac platform correctly
- ✅ Shows helpful error message
- ✅ Suggests alternatives
- ✅ Explains why it fails

### Option 4: Create Mac Binaries (Future)

To support Mac properly, need to:
1. Port Linux syscalls to Mac equivalents (or use compatibility layer)
2. Build Mac binaries: `GOOS=darwin go build ...`
3. Upload to GitHub release
4. Install script will auto-download

## Updated Install Script

The install script now:
- ✅ Better error handling
- ✅ Shows platform-specific messages
- ✅ Suggests alternatives for Mac
- ✅ Handles 404 errors gracefully

## Testing on Mac

After updating the install script:

```bash
# Try install (will show helpful error)
npm install -g nebulabox

# Should see:
# ⚠️  Note: Mac binaries may not be available yet
# ❌ Installation failed: Binary not found for platform: darwin-amd64
# 💡 Mac Installation Alternatives:
#    1. Build from source (CLI only)
#    2. Use Linux VM/Container
#    3. Wait for Mac binaries
```

## Next Steps

1. ✅ **Fixed**: Install script shows better errors
2. ⏳ **Future**: Port code to support Mac (or use compatibility)
3. ⏳ **Future**: Build and upload Mac binaries
4. ⏳ **Future**: Test Mac installation

## Quick Fix for Now

For Mac users, recommend:

```bash
# Option A: Build from source
git clone https://github.com/kunal1274/nebulabox.git
cd nebulabox && go build -o nbx ./cmd/nebulabox
sudo mv nbx /usr/local/bin/

# Option B: Use Linux
docker run -it ubuntu:22.04
# Then install inside Linux
```

---

**Install script updated with better Mac error handling!** 🍎

