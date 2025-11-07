# Quick Steps to Upload Binaries

## ✅ What's Ready

Binaries are built and ready in `dist/` directory:
- `nebulabox-linux-amd64` (11MB)
- `nebulabox-linux-arm64` (11MB)
- `nbx-linux-amd64` (11MB) - alias
- `nbx-linux-arm64` (11MB) - alias
- Checksum files (.sha256)

## 🚀 Upload Steps (5 minutes)

### Step 1: Open GitHub Release
Go to: **https://github.com/kunal1274/nebulabox/releases**

### Step 2: Edit Release
- Find release **"v0.1.0"**
- Click the **"Edit"** button (pencil icon ✏️) on the right side

### Step 3: Upload Files
- Scroll down to **"Assets"** section
- Click **"Attach binaries by selecting them"** or drag & drop
- Navigate to: `nebulabox/dist/` directory
- Select **ALL** these files:
  ```
  ✅ nebulabox-linux-amd64
  ✅ nebulabox-linux-arm64
  ✅ nbx-linux-amd64
  ✅ nbx-linux-arm64
  ✅ nebulabox-linux-amd64.sha256
  ✅ nebulabox-linux-arm64.sha256
  ```
- Wait for uploads to complete (shows progress)

### Step 4: Save
- Scroll to bottom
- Click **"Update release"** button
- Done! ✅

## ✅ Verify Upload

After uploading, refresh the page. You should see:
- All 6 files listed under "Assets"
- Download buttons for each file

## 🧪 Test Download

After upload, test from another system:

```bash
# Download
wget https://github.com/kunal1274/nebulabox/releases/download/v0.1.0/nbx-linux-amd64

# Install
chmod +x nbx-linux-amd64
sudo mv nbx-linux-amd64 /usr/local/bin/nbx

# Test
nbx version
```

## 📝 Alternative: GitHub CLI

If you have GitHub CLI installed:

```bash
cd dist/
gh release upload v0.1.0 nebulabox-linux-amd64 nebulabox-linux-arm64 nbx-linux-amd64 nbx-linux-arm64 *.sha256
```

---

**That's it!** After upload, users can download and install. 🎉

