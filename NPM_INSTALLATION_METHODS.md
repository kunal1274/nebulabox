# NebulaBox Installation via npm/pnpm/yarn

## Overview

NebulaBox can now be installed via popular Node.js package managers!

## Installation Methods

### 1. npm (Node Package Manager)

```bash
# Global install (recommended)
npm install -g nebulabox

# Local install (project-specific)
npm install nebulabox

# Use without installing (npx)
npx nebulabox version
```

### 2. pnpm (Fast, disk space efficient)

```bash
# Global install
pnpm add -g nebulabox

# Local install
pnpm add nebulabox

# Execute without install
pnpm exec nbx version
```

### 3. yarn (Classic or v2+)

```bash
# Global install
yarn global add nebulabox

# Local install
yarn add nebulabox

# Execute
yarn nbx version
```

### 4. npx (No installation needed)

```bash
# Run directly without installing
npx nebulabox version
npx nebulabox ps
npx nebulabox --help
```

## How It Works

1. **Package Installation**: Installs the npm package
2. **Post-Install Script**: Automatically downloads the correct binary for your platform
3. **Binary Location**: Binaries are stored in `node_modules/nebulabox/bin/`
4. **Commands Available**: `nbx` and `nebulabox` commands become available

## Platform Support

The package automatically detects and downloads the correct binary:

| Platform | Architecture | Binary |
|----------|-------------|--------|
| Linux | amd64 | `nbx-linux-amd64` |
| Linux | arm64 | `nbx-linux-arm64` |
| macOS | amd64 | `nbx-darwin-amd64` |
| macOS | arm64 | `nbx-darwin-arm64` |
| Windows | amd64 | `nbx-windows-amd64.exe` |
| Windows | arm64 | `nbx-windows-arm64.exe` |

## Usage After Installation

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

# Hierarchy commands
nbx hierarchy create
nbx hierarchy list
```

## Advantages of npm Installation

✅ **Easy Installation**: One command to install  
✅ **Automatic Updates**: `npm update -g nebulabox`  
✅ **Version Management**: Install specific versions  
✅ **No Manual Downloads**: Binary downloaded automatically  
✅ **Cross-Platform**: Works on Linux, macOS, Windows  
✅ **npx Support**: Use without installing  

## Comparison with Other Methods

| Method | Pros | Cons |
|--------|------|------|
| **npm/pnpm/yarn** | Easy, auto-updates, version management | Requires Node.js |
| **Go Install** | No Node.js needed, always latest | Requires Go |
| **Binary Download** | Simple, no dependencies | Manual updates |
| **Build from Source** | Full control | Requires Go, build time |

## Troubleshooting

### Issue: Binary not found after install

```bash
# Check if binary was downloaded
ls node_modules/nebulabox/bin/

# Re-run install script
cd node_modules/nebulabox
node scripts/install.js
```

### Issue: Command not found

```bash
# Check if in PATH
which nbx

# For global install, check npm global bin
npm config get prefix
# Add to PATH: export PATH=$PATH:$(npm config get prefix)/bin
```

### Issue: Permission denied

```bash
# Use sudo for global install
sudo npm install -g nebulabox

# Or fix npm permissions
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
export PATH=~/.npm-global/bin:$PATH
```

### Issue: Wrong platform binary

The install script auto-detects platform. If wrong:
1. Check `process.platform` and `process.arch`
2. Verify GitHub release has binary for your platform
3. Manually download if needed

## Publishing to npm

See `npm-publish-guide.md` for complete publishing instructions.

Quick steps:
```bash
# 1. Login to npm
npm login

# 2. Publish
npm publish

# 3. Verify
npm view nebulabox
```

## Version Management

```bash
# Install specific version
npm install -g nebulabox@0.1.0

# Update to latest
npm update -g nebulabox

# Check installed version
npm list -g nebulabox
```

## CI/CD Integration

Use in GitHub Actions, GitLab CI, etc.:

```yaml
# Example: GitHub Actions
- name: Install NebulaBox
  run: npm install -g nebulabox

- name: Use NebulaBox
  run: nbx version
```

## Next Steps

1. ✅ Test locally: `npm pack` and `npm install -g ./nebulabox-0.1.0.tgz`
2. ✅ Create npm account: https://www.npmjs.com/signup
3. ✅ Publish: `npm publish`
4. ✅ Test install: `npm install -g nebulabox`
5. ✅ Update documentation

---

**Ready to publish to npm!** 🚀

