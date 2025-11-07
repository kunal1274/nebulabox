# Mode Switcher Guide

## 🎯 What Was Added

A **Mode Switcher** component has been added to the dashboard header that allows you to switch between:

- **Mock Mode** - Static mock data only (no persistence)
- **Test Mode** - UAT sandbox with in-memory storage (default)
- **Live Mode** - Production mode with full persistence

## 📍 Location

The mode switcher appears in the **top-right corner** of every dashboard page, next to the header.

## 🎨 Visual Indicators

- **Current Mode Badge**: Shows current mode in colored badge (MOCK/TEST/LIVE)
- **Active Button**: The current mode button is highlighted
- **Image Counter**: Shows built images count vs total images

## 🚀 How to Use

### Switch Modes:

1. **Look for Mode Switcher** in top-right corner
2. **Click the mode button** you want (Mock/Test/Live)
3. **Dashboard reloads** automatically to apply new mode

### What Each Mode Does:

#### Mock Mode
- Shows only static mock data
- No built images or containers stored
- Perfect for UI demos
- All data cleared when switching to this mode

#### Test Mode (Default) ✅
- **Stores built images** in memory
- **Stores created containers** in memory
- Perfect for UAT sandbox testing
- Data persists during server session
- Data lost on server restart

#### Live Mode
- Full registry integration
- Persistent storage
- Production-ready
- Data survives restarts

## 📊 Current Status Display

The mode switcher shows:
- **Current mode** (colored badge)
- **Built Images**: Number of images built via buildspec
- **Total Images**: Total images available

Example:
```
Mode: TEST  [Mock] [Test] [Live]  Images: 2/5
```

## 🔄 Mode Switching Behavior

When switching modes:
1. Mode change request sent to API
2. Dashboard automatically reloads
3. All pages now use new mode
4. Data appropriate to new mode is displayed

### Important Notes:

- **Switching to Mock**: Clears all built images and containers
- **Switching to Test/Live**: Preserves existing data
- **Mode persists**: Until changed again or server restarts (for test mode)

## 🧪 Testing Flow

### 1. Verify Current Mode
- Check mode switcher shows "TEST" (default)
- Verify badge color is blue

### 2. Build Image in Test Mode
```bash
curl -X POST http://localhost:8081/api/buildspec/build \
  -H "Content-Type: application/json" \
  -d @mvp/app/single-container/buildspec.json
```

### 3. Check Image Appears
- Go to Images page
- Should see built image
- Mode switcher shows: `Images: 1/X`

### 4. Switch to Mock Mode
- Click "Mock" button
- Dashboard reloads
- Images page shows only mock data
- Built images cleared

### 5. Switch Back to Test Mode
- Click "Test" button
- Dashboard reloads
- Can build new images again

## 🎯 Quick Reference

| Action | Mock Mode | Test Mode | Live Mode |
|--------|-----------|-----------|-----------|
| Built Images Stored | ❌ | ✅ | ✅ |
| Containers Stored | ❌ | ✅ | ✅ |
| Persists After Restart | ❌ | ❌ | ✅ |
| Mock Data Shown | ✅ | ✅ (if no real data) | ❌ |
| Use Case | Demos | UAT Testing | Production |

## 🔍 Troubleshooting

### Mode Switcher Not Showing
- Rebuild dashboard: `cd web/dashboard && npm run build`
- Restart dashboard server
- Hard refresh browser (Ctrl+Shift+R)

### Mode Change Not Working
- Check browser console for errors
- Verify API is running on port 8081
- Check API logs for mode change requests

### Images Not Appearing After Switch
- Verify you're in Test or Live mode
- Build a new image after switching
- Check API: `curl http://localhost:8081/api/mode`

## 📝 API Commands

You can also change mode via API:

```bash
# Get current mode
curl http://localhost:8081/api/mode

# Set mode
curl -X PUT http://localhost:8081/api/mode \
  -H "Content-Type: application/json" \
  -d '{"mode": "test"}'
```

The dashboard mode switcher does this automatically!

