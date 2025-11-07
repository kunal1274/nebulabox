# Installing NebulaBox CLI (nbx)

## Installation Methods

### Method 1: Install from Source (Recommended)

Install directly from GitHub using Go:

```bash
# Install latest version
go install github.com/nebulabox/nebulabox/cmd/nebulabox@latest

# This installs as 'nebulabox', create alias for 'nbx'
alias nbx=nebulabox

# Or install to specific location with custom name
go build -o ~/.local/bin/nbx ./cmd/nebulabox
```

### Method 2: Install from Binary Release

Download pre-built binaries from GitHub Releases:

```bash
# Download latest release
VERSION="v0.1.0"
OS="linux"  # or "darwin", "windows"
ARCH="amd64"  # or "arm64"

# Download binary
wget https://github.com/nebulabox/nebulabox/releases/download/${VERSION}/nebulabox-${OS}-${ARCH}

# Make executable and move to PATH
chmod +x nebulabox-${OS}-${ARCH}
sudo mv nebulabox-${OS}-${ARCH} /usr/local/bin/nbx

# Or install to user directory
mkdir -p ~/.local/bin
mv nebulabox-${OS}-${ARCH} ~/.local/bin/nbx
chmod +x ~/.local/bin/nbx
export PATH="$HOME/.local/bin:$PATH"
```

### Method 3: Build from Source

```bash
# Clone repository
git clone https://github.com/nebulabox/nebulabox.git
cd nebulabox

# Build
make build-cli-test

# Install locally
./scripts/install-nbx.sh

# Or manually
cp bin/nbx ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"
```

### Method 4: Using Package Managers (Future)

#### Homebrew (macOS/Linux)
```bash
brew install nebulabox/tap/nbx
```

#### Snap (Linux)
```bash
sudo snap install nbx
```

#### AUR (Arch Linux)
```bash
yay -S nebulabox-bin
```

## Verify Installation

```bash
# Check if nbx is available
which nbx

# Check version
nbx version

# Test help
nbx --help
```

## Update Installation

### If installed via go install:
```bash
go install github.com/nebulabox/nebulabox/cmd/nebulabox@latest
```

### If installed via binary:
```bash
# Download new version and replace
wget https://github.com/nebulabox/nebulabox/releases/download/latest/nebulabox-linux-amd64
chmod +x nebulabox-linux-amd64
sudo mv nebulabox-linux-amd64 /usr/local/bin/nbx
```

## Uninstall

```bash
# Remove from PATH locations
rm ~/.local/bin/nbx
rm ~/.local/bin/nebulabox
# or
sudo rm /usr/local/bin/nbx
sudo rm /usr/local/bin/nebulabox

# Remove Go installation
go clean -i github.com/nebulabox/nebulabox/cmd/nebulabox
```

## Troubleshooting

### Command Not Found

1. **Check if installed:**
   ```bash
   which nbx
   which nebulabox
   ```

2. **Check PATH:**
   ```bash
   echo $PATH
   ```

3. **Add to PATH:**
   ```bash
   export PATH="$HOME/.local/bin:$PATH"
   # Or add to ~/.bashrc permanently
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   ```

### Permission Denied

```bash
chmod +x ~/.local/bin/nbx
# or
sudo chmod +x /usr/local/bin/nbx
```

## Requirements

- Go 1.22+ (for building from source)
- Linux/macOS/Windows
- 50MB+ disk space

## Installation Locations

Default installation locations:
- `go install`: `$GOPATH/bin` or `$HOME/go/bin`
- Manual: `~/.local/bin` or `/usr/local/bin`
- System: `/usr/bin` (requires sudo)

