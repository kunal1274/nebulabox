# Release Guide - NebulaBox CLI

## Prerequisites

### 1. Configure Git (Required)

Before creating a release, you must configure Git:

```bash
# Quick setup
./scripts/setup-git.sh

# Or manually
git config user.name "Your Name"
git config user.email "your.email@example.com"
```

**Verify configuration:**
```bash
git config --get user.name
git config --get user.email
```

### 2. Ensure Code is Committed

```bash
# Check status
git status

# Commit any changes
git add .
git commit -m "Prepare release v0.1.0"
```

## Creating a Release

### Method 1: Automated Release Script (Recommended)

```bash
# Create release (handles Git config check, building, tagging)
./scripts/create-release.sh v0.1.0
```

This script will:
1. ✅ Check Git configuration
2. ✅ Build binaries for all platforms
3. ✅ Create Git tag
4. ✅ Optionally push tag to GitHub

### Method 2: Manual Release

```bash
# 1. Build release binaries
./scripts/build-release.sh v0.1.0

# 2. Create Git tag
git tag -a v0.1.0 -m "Release v0.1.0"

# 3. Push tag
git push origin v0.1.0

# 4. Create GitHub release (upload dist/* files)
```

## Release Process

### Step-by-Step

1. **Setup Git (if not done):**
   ```bash
   ./scripts/setup-git.sh
   ```

2. **Update Version:**
   - Update version in code/docs if needed
   - Commit changes

3. **Create Release:**
   ```bash
   ./scripts/create-release.sh v0.1.0
   ```

4. **Push Tag:**
   ```bash
   git push origin v0.1.0
   ```

5. **GitHub Actions** (if configured):
   - Automatically builds and creates release
   - Uploads binaries to GitHub Releases

6. **Manual GitHub Release** (if needed):
   - Go to: https://github.com/nebulabox/nebulabox/releases/new
   - Select tag: v0.1.0
   - Upload files from `dist/` directory

## Version Format

Use semantic versioning: `vMAJOR.MINOR.PATCH`

- **v0.1.0** - First release
- **v0.2.0** - New features
- **v1.0.0** - Stable release
- **v1.0.1** - Bug fixes

## Troubleshooting

### Git Configuration Error

**Error:** "Make sure you configure your 'user.name' and 'user.email' in git"

**Solution:**
```bash
# Quick fix
./scripts/setup-git.sh

# Or manual
git config user.name "Your Name"
git config user.email "your.email@example.com"
```

### Tag Already Exists

**Error:** "Tag v0.1.0 already exists"

**Solution:**
```bash
# Delete local tag
git tag -d v0.1.0

# Delete remote tag
git push origin :refs/tags/v0.1.0

# Create new tag
./scripts/create-release.sh v0.1.0
```

### Uncommitted Changes

**Warning:** "Working directory has uncommitted changes"

**Solution:**
```bash
# Commit changes
git add .
git commit -m "Your commit message"

# Or stash
git stash
```

## After Release

### Users Can Install

Once the release is published, users can install:

```bash
# Install latest
go install github.com/nebulabox/nebulabox/cmd/nebulabox@latest

# Install specific version
go install github.com/nebulabox/nebulabox/cmd/nebulabox@v0.1.0
```

### Verify Release

```bash
# Check tag exists
git tag -l

# Check release on GitHub
# Visit: https://github.com/nebulabox/nebulabox/releases
```

## Quick Commands

```bash
# Setup Git
./scripts/setup-git.sh

# Create release
./scripts/create-release.sh v0.1.0

# Build only (no tag)
./scripts/build-release.sh v0.1.0

# Check Git config
git config --get user.name
git config --get user.email
```

---

**Ready to release!** 🚀

