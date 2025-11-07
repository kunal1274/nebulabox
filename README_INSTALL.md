# Install NebulaBox CLI (nbx) - Quick Guide

## 🚀 Quick Install (Recommended)

### Option 1: Install from GitHub (Go Install)

```bash
# Install latest version
go install github.com/nebulabox/nebulabox/cmd/nebulabox@latest

# Create nbx alias
alias nbx=nebulabox
# Or: ln -s $(go env GOPATH)/bin/nebulabox $(go env GOPATH)/bin/nbx

# Verify
nbx version
```

### Option 2: Download Binary Release

```bash
# Linux (amd64)
wget https://github.com/nebulabox/nebulabox/releases/download/v0.1.0/nbx-linux-amd64
chmod +x nbx-linux-amd64
sudo mv nbx-linux-amd64 /usr/local/bin/nbx

# macOS (Intel)
wget https://github.com/nebulabox/nebulabox/releases/download/v0.1.0/nbx-darwin-amd64
chmod +x nbx-darwin-amd64
sudo mv nbx-darwin-amd64 /usr/local/bin/nbx

# macOS (Apple Silicon)
wget https://github.com/nebulabox/nebulabox/releases/download/v0.1.0/nbx-darwin-arm64
chmod +x nbx-darwin-arm64
sudo mv nbx-darwin-arm64 /usr/local/bin/nbx
```

### Option 3: Build from Source

```bash
# Clone repository
git clone https://github.com/nebulabox/nebulabox.git
cd nebulabox

# Build
make build-cli-test

# Install locally
./scripts/install-nbx.sh

# Or use directly
./bin/nbx version
```

## ✅ Verify Installation

```bash
nbx version
nbx --help
nbx ps
```

## 📚 Full Documentation

- **Installation Guide**: See `INSTALL.md`
- **Publishing Guide**: See `PUBLISH.md`
- **Quick Start**: See `QUICK_START.md`

## 🔄 Update Installation

```bash
# If installed via go install
go install github.com/nebulabox/nebulabox/cmd/nebulabox@latest

# If installed via binary
# Download new version and replace old binary
```

---

**That's it! You can now use `nbx` from anywhere!** 🎉

