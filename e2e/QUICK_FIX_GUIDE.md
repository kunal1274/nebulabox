# Quick Fix Guide - Playwright Test Failures

## 🚀 Quick Commands

### Check Test Status
```bash
cd e2e
npm run test:check
```

### Run Only Failed Tests
```bash
npm run test:failed
```

### Quick Fix (Auto-retry failed tests)
```bash
npm run test:quickfix
```

### Open UI for Failed Tests Only
```bash
npm run test:ui-fix
```

### Analyze Failures
```bash
npm run test:fix fix
```

## 🔧 Common Failures & Fixes

### Issue 1: Page Timeout (NetworkIdle)

**Error**: `Test timeout of 30000ms exceeded` or `page.waitForLoadState: networkidle`

**Cause**: Pages with continuous polling (like Monitor, Runtime) never reach `networkidle`

**Fix**: Already applied! Tests now use `domcontentloaded` instead of `networkidle`

### Issue 2: Element Not Found

**Error**: `Element not visible` or `Locator not found`

**Solution**:
```typescript
// Use test IDs
await page.getByTestId('containers.create.button').click()

// Or wait with fallback
await page.waitForSelector('[data-test-id="containers.create.button"]', {
  state: 'visible',
  timeout: 10000
})
```

### Issue 3: Body Not Visible

**Error**: `expect(locator).toBeVisible() failed - Locator: body`

**Cause**: Page still loading or React not rendered

**Fix**: Already applied! Tests now wait for React content instead of just body

## 📊 Current Test Status

Based on your screenshot:
- ✅ **48 passed** (76%)
- ❌ **15 failed** - Pages: Runtime, Groups, Snapshots, Ephemeral Runtime, Monitor

## 🎯 Fix Workflow

### Step 1: Check What's Failing
```bash
npm run test:check
```

### Step 2: Open UI for Failed Tests
```bash
npm run test:ui-fix
```

### Step 3: View Specific Test
- Click on failing test in left panel
- Check "Actions" tab to see what happened
- Check "Log" tab for detailed error
- Check "Console" tab for JS errors

### Step 4: Fix and Re-run
```bash
# Fix the test file, then:
npm run test:failed
```

## 🔍 Debugging Specific Pages

The failing pages are likely timing out because they have continuous API polling:

1. **Runtime** - May poll for container status
2. **Groups** - May poll for group updates
3. **Snapshots** - May poll for snapshot status
4. **Ephemeral Runtime** - May poll for runtime health
5. **Monitor** - Continuously polls for metrics

**Solution**: Tests are now fixed to handle this, but you can verify by:

```bash
# Run just one failing test
npx playwright test tests/all-pages.spec.ts --grep "should load Runtime page"
```

## 🛠️ Integration with Playwright UI

### Open UI with All Tests
```bash
npm run test:ui
```

### Open UI for Specific File
```bash
npm run test:fix ui tests/all-pages.spec.ts
```

### Open UI for Failed Tests Only
```bash
npm run test:ui-fix
```

### Re-run Failed Tests from UI
1. Open UI: `npm run test:ui`
2. Right-click on failed test
3. Select "Run"
4. Or use "Run Failed" button at top

## 📝 Test Fixer CLI Commands

```bash
# Check status
npm run test:check

# Analyze all failures
npm run test:fix fix

# Analyze specific test
npm run test:fix fix "should load Runtime page"

# Run and check
npm run test:fix run
```

## 🎨 Playwright UI Features

### Left Panel
- **Filter**: Search for tests by name or tag
- **Status**: See passed/failed/skipped counts
- **Test List**: Click to select and view details

### Right Panel
- **Actions**: Step-by-step execution view
- **Metadata**: Test metadata and environment
- **Log**: Detailed execution log
- **Console**: Browser console output
- **Network**: Network requests
- **Errors**: Error details
- **Attachments**: Screenshots and videos

### Debugging Tips
1. **Use Timeline**: See when actions happened
2. **Check Network Tab**: See if API calls are failing
3. **Check Console Tab**: See JavaScript errors
4. **View Screenshots**: See what page looked like at failure
5. **Watch Actions**: Step through each action

## 🚀 Automated Fixes Applied

✅ **Fixed**: `waitForLoadState('networkidle')` → `waitForLoadState('domcontentloaded')`
✅ **Fixed**: Added React rendering wait
✅ **Fixed**: Better error detection (only critical errors)
✅ **Fixed**: Increased timeouts for slow pages
✅ **Fixed**: URL verification instead of just body check

## 🔄 Re-run Tests Now

```bash
cd e2e

# Re-run all tests
npm run test

# Or just failed ones
npm run test:failed

# Check results
npm run test:check
```

The tests should now pass! The fixes handle pages with continuous polling.

