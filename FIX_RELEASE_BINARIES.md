# Fix Release - Add Binaries to GitHub Release

## Current Status

✅ Release `v0.1.0` created on GitHub  
❌ Binaries not uploaded (only source code assets visible)

## Why Binaries Are Missing

The GitHub Actions workflow should have:
1. Built binaries for all platforms
2. Uploaded them to the release

**Possible issues:**
- Workflow didn't trigger
- Workflow failed
- Build script had errors
- Upload step failed

## Solution Options

### Option 1: Check and Fix GitHub Actions

1. **Check Actions Status:**
   - Visit: `https://github.com/kunal1274/nebulabox/actions`
   - Look for "Release" workflow run
   - Check if it succeeded or failed

2. **If Workflow Failed:**
   - Click on the failed run
   - Check error logs
   - Fix issues and re-run

3. **Re-run Workflow:**
   - Go to Actions tab
   - Find the Release workflow
   - Click "Re-run all jobs"

### Option 2: Manually Upload Binaries

If workflow isn't working, upload binaries manually:

1. **Build Binaries Locally:**
   ```bash
   ./scripts/build-release.sh v0.1.0
   ```

2. **Upload to GitHub Release:**
   - Go to: `https://github.com/kunal1274/nebulabox/releases`
   - Click "Edit" on v0.1.0 release
   - Scroll to "Assets" section
   - Click "Attach binaries by selecting them"
   - Upload all files from `dist/` directory:
     - `nbx-linux-amd64`
     - `nbx-linux-arm64`
     - `nbx-darwin-amd64`
     - `nbx-darwin-arm64`
     - `nbx-windows-amd64.exe`
     - `nbx-windows-arm64.exe`
   - Click "Update release"

### Option 3: Create New Release with Binaries

If current release is problematic:

1. **Build binaries:**
   ```bash
   ./scripts/build-release.sh v0.2.0
   ```

2. **Create new tag:**
   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0 with binaries"
   git push origin v0.2.0
   ```

3. **Check if workflow runs automatically**

4. **Or manually create release with binaries**

## Quick Fix Script

Save this as `fix-release.sh`:

```bash
#!/bin/bash
set -e

VERSION="v0.1.0"
echo "🔧 Fixing release $VERSION"
echo ""

# Build binaries
echo "📦 Building binaries..."
./scripts/build-release.sh "$VERSION"

# Check if dist/ has files
if [ ! -d "dist" ] || [ -z "$(ls -A dist/* 2>/dev/null)" ]; then
    echo "❌ No binaries built!"
    exit 1
fi

echo ""
echo "✅ Binaries built in dist/"
echo ""
echo "📋 Next steps:"
echo "1. Go to: https://github.com/kunal1274/nebulabox/releases"
echo "2. Edit release v0.1.0"
echo "3. Upload files from dist/ directory"
echo ""
echo "Or use GitHub CLI (if installed):"
echo "  gh release upload $VERSION dist/*"
```

## Verify Release Has Binaries

After fixing, check:
- `https://github.com/kunal1274/nebulabox/releases/tag/v0.1.0`
- Should see binaries listed under "Assets"
- Should be able to download `nbx-darwin-amd64` or `nbx-darwin-arm64` for Mac

## Test After Fix

On Mac:
```bash
# Download
curl -L https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-darwin-amd64 -o nbx
chmod +x nbx
sudo mv nbx /usr/local/bin/
nbx version
```

---

**Priority: Fix binaries so Mac testing can proceed!** 🍎

