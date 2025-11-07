# Upload Binaries to GitHub Release - Step by Step

## Step 1: Build Binaries (Already Done)

Binaries should be in `dist/` directory:
- `nebulabox-linux-amd64`
- `nebulabox-linux-arm64`
- `nbx-linux-amd64`
- `nbx-linux-arm64`
- (and checksums)

## Step 2: Upload to GitHub Release

### Method 1: Via GitHub Web UI (Easiest)

1. **Go to Releases:**
   - Visit: `https://github.com/kunal1274/nebulabox/releases`

2. **Edit Existing Release:**
   - Find release `v0.1.0`
   - Click the **"Edit"** button (pencil icon) next to the release

3. **Upload Binaries:**
   - Scroll down to the **"Assets"** section
   - Click **"Attach binaries by selecting them"** or drag & drop
   - Select ALL files from `dist/` directory:
     ```
     nebulabox-linux-amd64
     nebulabox-linux-arm64
     nbx-linux-amd64
     nbx-linux-arm64
     nebulabox-linux-amd64.sha256
     nebulabox-linux-arm64.sha256
     ```
   - Wait for uploads to complete

4. **Update Release:**
   - Click **"Update release"** button at the bottom
   - Done! ✅

### Method 2: Via GitHub CLI (If Installed)

```bash
# Install GitHub CLI if not installed
# macOS: brew install gh
# Linux: sudo apt install gh

# Authenticate
gh auth login

# Upload binaries to release
cd dist/
gh release upload v0.1.0 nebulabox-linux-amd64
gh release upload v0.1.0 nebulabox-linux-arm64
gh release upload v0.1.0 nbx-linux-amd64
gh release upload v0.1.0 nbx-linux-arm64
gh release upload v0.1.0 *.sha256

# Or upload all at once
gh release upload v0.1.0 dist/*
```

### Method 3: Using curl (Advanced)

```bash
# Get GitHub token first (Settings > Developer settings > Personal access tokens)
GITHUB_TOKEN="your_token_here"
REPO="kunal1274/nebulabox"
TAG="v0.1.0"

# Upload each file
for file in dist/nebulabox-linux-amd64 dist/nebulabox-linux-arm64 dist/nbx-linux-amd64 dist/nbx-linux-arm64; do
    curl -X POST \
        -H "Authorization: token $GITHUB_TOKEN" \
        -H "Content-Type: application/octet-stream" \
        --data-binary @"$file" \
        "https://uploads.github.com/repos/$REPO/releases/$(gh api repos/$REPO/releases/tags/$TAG --jq .id)/assets?name=$(basename $file)"
done
```

## Step 3: Verify Upload

1. **Check Release Page:**
   - Go to: `https://github.com/kunal1274/nebulabox/releases/tag/v0.1.0`
   - You should see binaries listed under "Assets"

2. **Test Download:**
   ```bash
   # Test download
   curl -L https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64 -o /tmp/nbx-test
   chmod +x /tmp/nbx-test
   /tmp/nbx-test version
   ```

## Quick Upload Script

Save this as `upload-binaries.sh`:

```bash
#!/bin/bash
set -e

RELEASE="v0.1.0"
REPO="kunal1274/nebulabox"

echo "📦 Uploading binaries to release $RELEASE"
echo ""

# Check if gh CLI is installed
if command -v gh &> /dev/null; then
    echo "Using GitHub CLI..."
    cd dist/
    gh release upload "$RELEASE" nebulabox-linux-amd64 nebulabox-linux-arm64 nbx-linux-amd64 nbx-linux-arm64 *.sha256
    echo "✅ Upload complete!"
else
    echo "❌ GitHub CLI not installed"
    echo ""
    echo "Please upload manually:"
    echo "1. Go to: https://github.com/$REPO/releases/edit/$RELEASE"
    echo "2. Upload files from dist/ directory"
    echo ""
    echo "Or install GitHub CLI:"
    echo "  brew install gh  # macOS"
    echo "  sudo apt install gh  # Linux"
fi
```

## Troubleshooting

### Issue: "Edit" button not visible
- Make sure you're logged into GitHub
- Check you have write access to the repository

### Issue: Upload fails
- Check file size (GitHub has limits)
- Try uploading one file at a time
- Check internet connection

### Issue: Files not showing after upload
- Refresh the page
- Check browser console for errors
- Try uploading again

## After Upload

Once binaries are uploaded, users can:

```bash
# Download and install
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64
chmod +x nbx-linux-amd64
sudo mv nbx-linux-amd64 /usr/local/bin/nbx
nbx version
```

---

**Ready to upload!** 🚀

