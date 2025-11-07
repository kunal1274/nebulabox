# Publish NebulaBox to npm/pnpm/yarn

## Overview

NebulaBox can be installed via:
- ✅ **npm** - `npm install -g nebulabox`
- ✅ **pnpm** - `pnpm add -g nebulabox`
- ✅ **yarn** - `yarn global add nebulabox`
- ✅ **npx** - `npx nebulabox version` (no install needed)

## How It Works

The npm package:
1. Downloads the appropriate binary for your platform during `postinstall`
2. Creates symlinks (`nbx` and `nebulabox`) in `node_modules/.bin`
3. Makes binaries available globally or via `npx`

## Prerequisites

1. **npm account** - Create at https://www.npmjs.com/signup
2. **Login to npm** - `npm login`
3. **Package.json** - Already created ✅
4. **Install scripts** - Already created ✅

## Publishing Steps

### Step 1: Prepare Package

```bash
# Ensure package.json is correct
cat package.json

# Test install script locally
node scripts/install.js
```

### Step 2: Test Locally

```bash
# Create a test package
npm pack

# Install locally
npm install -g ./nebulabox-0.1.0.tgz

# Test
nbx version
```

### Step 3: Publish to npm

```bash
# Make sure you're logged in
npm whoami

# If not logged in
npm login

# Publish
npm publish

# Or publish as public (if private by default)
npm publish --access public
```

### Step 4: Verify Publication

```bash
# Check package on npm
npm view nebulabox

# Test install from npm
npm install -g nebulabox@0.1.0
nbx version
```

## Installation Methods

### Via npm

```bash
# Global install
npm install -g nebulabox

# Local install
npm install nebulabox

# Use with npx (no install)
npx nebulabox version
```

### Via pnpm

```bash
# Global install
pnpm add -g nebulabox

# Local install
pnpm add nebulabox

# Use with pnpm exec
pnpm exec nbx version
```

### Via yarn

```bash
# Global install
yarn global add nebulabox

# Local install
yarn add nebulabox

# Use with yarn
yarn nbx version
```

## Package Structure

```
nebulabox/
├── package.json          # Package metadata
├── bin/                  # Binaries (created during install)
│   ├── nbx              # Symlink to binary
│   └── nebulabox        # Symlink to binary
├── scripts/
│   ├── install.js       # Downloads binary for platform
│   └── uninstall.js     # Cleans up binaries
└── README.md            # Package documentation
```

## Version Management

### Update Version

```bash
# Update package.json version
npm version patch  # 0.1.0 -> 0.1.1
npm version minor  # 0.1.0 -> 0.2.0
npm version major  # 0.1.0 -> 1.0.0

# This also creates a git tag
```

### Publish New Version

```bash
# 1. Update version
npm version patch

# 2. Publish
npm publish

# 3. Create GitHub release (if needed)
git push --tags
```

## Platform Support

The package automatically downloads the correct binary:
- **Linux**: `nbx-linux-amd64` or `nbx-linux-arm64`
- **macOS**: `nbx-darwin-amd64` or `nbx-darwin-arm64`
- **Windows**: `nbx-windows-amd64.exe` or `nbx-windows-arm64.exe`

## Troubleshooting

### Issue: Binary not found after install

**Solution:**
```bash
# Check if binary was downloaded
ls node_modules/nebulabox/bin/

# Re-run install script
cd node_modules/nebulabox
node scripts/install.js
```

### Issue: Permission denied

**Solution:**
```bash
# Use sudo for global install
sudo npm install -g nebulabox

# Or fix npm permissions
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
export PATH=~/.npm-global/bin:$PATH
```

### Issue: Wrong platform binary

**Solution:**
The install script auto-detects platform. If wrong:
1. Check `process.platform` and `process.arch`
2. Update `scripts/install.js` if needed

### Issue: Download fails

**Solution:**
- Check GitHub release exists
- Verify version matches package.json
- Check network connectivity
- Try manual download to verify URL

## Alternative: Using npx (No Install)

Users can use without installing:

```bash
# Run directly
npx nebulabox version

# Or with alias
npx nebulabox@latest ps
```

## Publishing Checklist

- [ ] Update `package.json` version
- [ ] Test `npm pack` locally
- [ ] Test install script: `node scripts/install.js`
- [ ] Ensure GitHub release exists with binaries
- [ ] Login to npm: `npm login`
- [ ] Publish: `npm publish`
- [ ] Verify: `npm view nebulabox`
- [ ] Test install: `npm install -g nebulabox`
- [ ] Update README with npm install instructions

## CI/CD Integration

Add to `.github/workflows/publish-npm.yml`:

```yaml
name: Publish to npm

on:
  release:
    types: [created]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
          registry-url: 'https://registry.npmjs.org'
      - run: npm publish
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
```

## Next Steps

1. ✅ Create npm account
2. ✅ Test package locally
3. ✅ Publish to npm
4. ✅ Update documentation
5. ✅ Announce on social media/forums

---

**Ready to publish!** 🚀

