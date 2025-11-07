# Playwright Integration Guide

## 🎯 Overview

This guide explains how to integrate with Playwright UI, fix test failures, and monitor tests from the console.

## 🛠️ Available Commands

### Quick Commands

```bash
# Check test status (after running tests)
npm run test:check

# Show detailed status
npm run test:status

# Fix and analyze failures
npm run test:fix fix

# Run tests and auto-check
npm run test:fix run

# Open Playwright UI
npm run test:ui

# Open UI for specific file
npm run test:fix ui tests/all-pages.spec.ts

# Monitor tests continuously
npm run test:monitor watch
```

### Using the Test Fixer CLI

```bash
cd e2e

# Check for failures
npm run test:fix check

# Analyze specific test failure
npm run test:fix fix "should load"

# Analyze failures in specific file
npm run test:fix fix "" tests/all-pages.spec.ts

# Run tests and check
npm run test:fix run
```

### Using Monitor Script

```bash
cd e2e

# Check existing results
./scripts/monitor-tests.sh check

# Run tests and check
./scripts/monitor-tests.sh run

# Watch mode (auto-re-run)
./scripts/monitor-tests.sh watch
```

## 🔧 Fixing Test Failures

### Common Issues and Fixes

#### 1. Page Not Loading (404/Timeout)

**Error**: `Navigation timeout` or `404 Not Found`

**Fix**:
- Check if route exists in `web/dashboard/src/App.tsx`
- Verify page component is imported
- Check if API server is running on port 8081
- Verify dashboard is running on port 3001

**Example Fix**:
```typescript
// Wait longer for page to load
await page.goto('http://localhost:3001/some-page', { 
  waitUntil: 'networkidle', 
  timeout: 30000 
})
```

#### 2. Element Not Found

**Error**: `Element not visible` or `Element not found`

**Fix**:
- Add proper wait conditions
- Use test IDs instead of CSS selectors
- Check if element is rendered conditionally

**Example Fix**:
```typescript
// Wait for element with test ID
await page.getByTestId('containers.create.button').waitFor({ state: 'visible' })
```

#### 3. Network Errors

**Error**: `Network request failed` or `CORS error`

**Fix**:
- Check API server is running
- Verify CORS configuration
- Check network requests in Playwright UI

#### 4. Test Flakiness

**Error**: Intermittent failures

**Fix**:
- Add explicit waits instead of fixed timeouts
- Wait for network to be idle
- Use `page.waitForSelector` with proper conditions

## 📊 Monitoring Tests

### Real-time Monitoring

```bash
# Watch mode - automatically re-runs tests
npm run test:watch

# Or using monitor script
./scripts/monitor-tests.sh watch
```

### Check Status After Run

```bash
# After running tests in Playwright UI
npm run test:check

# Get detailed failure information
npm run test:status
```

## 🎨 Playwright UI Integration

### Opening UI for Specific Tests

```bash
# Open UI for all tests
npm run test:ui

# Open UI for specific file
npm run test:fix ui tests/all-pages.spec.ts

# Open UI for specific test pattern
npx playwright test --ui --grep "should load"
```

### Using UI to Fix Tests

1. **Open UI**: `npm run test:ui`
2. **Select failing test** in left panel
3. **View error details** in right panel
4. **Check "Actions" tab** to see what happened
5. **Check "Log" tab** for detailed logs
6. **Check "Console" tab** for JavaScript errors
7. **Check "Network" tab** for API issues

### Re-running Failed Tests

```bash
# Re-run only failed tests
npx playwright test --last-failed

# Re-run specific file
npx playwright test tests/all-pages.spec.ts

# Re-run with UI
npx playwright test --ui tests/all-pages.spec.ts
```

## 🔍 Debugging Workflow

### Step 1: Identify Failures

```bash
npm run test:check
```

### Step 2: Analyze Specific Failure

```bash
npm run test:fix fix "should load Images page"
```

### Step 3: Open UI for Detailed Debugging

```bash
npm run test:fix ui tests/all-pages.spec.ts
```

### Step 4: Fix and Re-run

```bash
# Fix the test file, then:
npx playwright test tests/all-pages.spec.ts --reporter=list
```

## 📝 Test Results Structure

Test results are stored in:
```
e2e/test-results/
├── results.json          # JSON results (parsed by fixer)
├── playwright-report/    # HTML report
└── [test-files]/         # Individual test artifacts
```

## 🚀 Quick Fix Workflow

```bash
# 1. Run tests
npm run test

# 2. Check for failures
npm run test:check

# 3. If failures found, analyze them
npm run test:fix fix

# 4. Open UI to debug specific test
npm run test:fix ui tests/[failing-file].spec.ts

# 5. Fix the test, then re-run
npm run test

# 6. Verify fix
npm run test:check
```

## 🎯 Best Practices

1. **Always check status after test run**
   ```bash
   npm run test:check
   ```

2. **Use test IDs for reliable selectors**
   ```typescript
   await page.getByTestId('containers.create.button').click()
   ```

3. **Wait for network idle when navigating**
   ```typescript
   await page.goto('http://localhost:3001/page', { 
     waitUntil: 'networkidle' 
   })
   ```

4. **Add explicit waits instead of fixed timeouts**
   ```typescript
   await expect(page.getByText('Content')).toBeVisible({ timeout: 10000 })
   ```

5. **Use Playwright UI for visual debugging**
   ```bash
   npm run test:ui
   ```

## 📚 Additional Resources

- [Playwright Documentation](https://playwright.dev)
- [Test-Aware System Guide](../TEST_AWARE_SYSTEM_GUIDE.md)
- [Automated Testing Guide](../AUTOMATED_TESTING_GUIDE.md)

