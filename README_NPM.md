# NebulaBox - npm Package

Install NebulaBox via npm, pnpm, or yarn!

## Quick Install

```bash
# npm
npm install -g nebulabox

# pnpm
pnpm add -g nebulabox

# yarn
yarn global add nebulabox

# npx (no install needed)
npx nebulabox version
```

## Usage

After installation:

```bash
# Check version
nbx version

# Get help
nbx --help

# List containers
nbx ps

# List images
nbx images

# Build an image
nbx build /path/to/buildspec.json

# Run a container
nbx run <image-name>
```

## How It Works

The npm package automatically:
1. Detects your platform (Linux/macOS/Windows)
2. Downloads the appropriate binary from GitHub Releases
3. Creates `nbx` and `nebulabox` commands

## Platform Support

- ✅ Linux (amd64, arm64)
- ✅ macOS (amd64, arm64) - CLI only
- ✅ Windows (amd64, arm64) - CLI only

**Note:** Container features require Linux. On macOS/Windows, only CLI structure works.

## Troubleshooting

### Binary not found

```bash
# Re-run install script
cd node_modules/nebulabox
node scripts/install.js
```

### Permission denied

```bash
# Use sudo for global install
sudo npm install -g nebulabox

# Or fix npm permissions
npm config set prefix ~/.npm-global
export PATH=~/.npm-global/bin:$PATH
```

### Wrong platform

The package auto-detects your platform. If you need a specific platform:

```bash
# Download manually
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64
chmod +x nbx-linux-amd64
sudo mv nbx-linux-amd64 /usr/local/bin/nbx
```

## Development

```bash
# Clone repository
git clone https://github.com/kunal1274/nebulabox.git
cd nebulabox

# Install dependencies (if any)
npm install

# Test install script
node scripts/install.js

# Build from source (Go required)
go build -o nbx ./cmd/nebulabox
```

## Links

- **GitHub**: https://github.com/kunal1274/nebulabox
- **npm**: https://www.npmjs.com/package/nebulabox
- **Releases**: https://github.com/kunal1274/nebulabox/releases

## License

MIT

