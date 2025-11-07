# Dashboard Fix Guide - Real Data Display

## 🔍 Problem

API returns correct data at `/api/images`:
```json
[{
  "id": "sha256:4e2271e467fe7fd714a74fce3b6d05a21b81b7f958620beb477d8c49b257ee14",
  "name": "my-app",
  "tag": "latest",
  "size": "290.0 MB",
  "created": "2025-11-01T03:34:02Z"
}]
```

But dashboard shows mock data instead.

## ✅ Fixes Applied

1. **Removed Mock Data Fallback** ✅
   - Images page: No longer shows mock data on error
   - Containers page: No longer shows mock data on error
   - This helps identify real API errors

2. **Added Console Logging** ✅
   - Logs API responses to browser console
   - Makes debugging easier

## 🧪 Debugging Steps

### Step 1: Check Browser Console
1. Open Dashboard: http://localhost:3001
2. Open DevTools (F12)
3. Go to Console tab
4. Navigate to Images page
5. Look for:
   - `Loaded images from API: [...]` - Success
   - `Failed to load images: ...` - Error
   - CORS errors
   - Network errors

### Step 2: Check Network Tab
1. DevTools → Network tab
2. Refresh Images page
3. Find request to `/api/images`
4. Check:
   - Status: Should be `200 OK`
   - Response: Should show the image data
   - Headers: Check CORS headers

### Step 3: Test API Directly
```bash
# Test images API
curl http://localhost:8081/api/images | jq '.'

# Test containers API  
curl http://localhost:8081/api/containers | jq '.'

# Test mode
curl http://localhost:8081/api/mode | jq '.'
```

### Step 4: Verify Mode
The API should be in `test` mode (default):
```bash
curl http://localhost:8081/api/mode
```

Should return:
```json
{
  "mode": "test",
  "description": "Test mode: UAT sandbox - Built images and containers are stored in memory",
  "builtImages": 1,
  "totalImages": 4
}
```

## 🔧 Common Issues & Solutions

### Issue 1: CORS Error
**Symptom:** Console shows `CORS policy` error  
**Solution:** Already fixed in server.go, but verify:
```go
config.AllowOrigins = []string{
    "http://localhost:3000",
    "http://localhost:3001",
    ...
}
```

### Issue 2: API Not Running
**Symptom:** Network tab shows `Failed to fetch`  
**Solution:** Start API server:
```bash
./nebulabox-api
# or
cd cmd/api && go run main.go
```

### Issue 3: Wrong Port
**Symptom:** No response  
**Solution:** Check API_BASE_URL in frontend:
```typescript
// web/dashboard/src/lib/api.ts
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081/api'
```

### Issue 4: Response Format Mismatch
**Symptom:** API returns data but frontend doesn't display  
**Solution:** Check Image interface matches API response:
```typescript
export interface Image {
  id: string
  name: string
  tag: string
  size: string
  created: string
}
```

### Issue 5: Error Being Silently Caught
**Symptom:** Empty array displayed  
**Solution:** Check console for actual error message

## 🎯 Quick Test Script

```bash
#!/bin/bash

echo "Testing API Integration..."

# 1. Health check
echo "1. API Health:"
curl -s http://localhost:8081/api/health | jq '.status'

# 2. Mode check
echo "2. Mode:"
curl -s http://localhost:8081/api/mode | jq '.mode'

# 3. Images
echo "3. Images:"
curl -s http://localhost:8081/api/images | jq '.[] | {name, tag, size}'

# 4. Containers
echo "4. Containers:"
curl -s http://localhost:8081/api/containers | jq '.[] | {name, image, status}'
```

## 📋 Verification Checklist

- [ ] API server is running on port 8081
- [ ] Dashboard is running on port 3001
- [ ] API returns images at `/api/images`
- [ ] Browser console shows no errors
- [ ] Network tab shows successful API call
- [ ] Response matches Image interface
- [ ] Dashboard displays real data (not mock)
- [ ] Refresh button works
- [ ] After building image, it appears in dashboard

## 🔄 Next Steps After Fix

1. **Build Image:**
   ```bash
   curl -X POST http://localhost:8081/api/buildspec/build \
     -H "Content-Type: application/json" \
     -d @mvp/app/single-container/buildspec.json
   ```

2. **Create Container:**
   ```bash
   curl -X POST http://localhost:8081/api/containers/run \
     -H "Content-Type: application/json" \
     -d '{"image": "mern-mvp:latest", "name": "mern-test"}'
   ```

3. **Verify in Dashboard:**
   - Images page → Should see built image
   - Containers page → Should see created container
   - Refresh → Should persist

## 🐛 If Still Not Working

1. Share browser console errors
2. Share network tab screenshot
3. Share API response from curl
4. Check if API and Dashboard are on correct ports
5. Verify mode is "test" not "mock"

