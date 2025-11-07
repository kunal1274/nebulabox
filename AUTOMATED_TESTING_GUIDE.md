# Automated Testing Guide - Complete E2E Coverage

## 🎯 Overview

NebulaBox now has **comprehensive automated E2E testing** that covers:
- ✅ All pages (smoke tests)
- ✅ Mode switching
- ✅ Image building and management
- ✅ Container creation and lifecycle
- ✅ API integration
- ✅ Data persistence
- ✅ Complete workflows

## 🚀 Quick Start

### 1. Install Dependencies
```bash
cd e2e
npm install
```

### 2. Run All Tests
```bash
npm run test
```

### 3. Run with UI (Recommended for Development)
```bash
npm run test:ui
```

### 4. Run Specific Test Categories
```bash
# Smoke tests only (fast)
npm run test:smoke

# MVP tests only
npm run test:mvp

# Specific test file
npx playwright test tests/complete-mvp-flow.spec.ts
```

## 📊 Test Coverage

### Test Files Created

1. **00-setup.spec.ts** - API and dashboard health checks
2. **all-pages.spec.ts** - Smoke tests for all 21 pages
3. **mode-switcher.spec.ts** - Mode switching functionality
4. **buildspec.spec.ts** - Build specification validation and building
5. **images-api-integration.spec.ts** - Image building and dashboard display
6. **containers-api-integration.spec.ts** - Container creation and lifecycle
7. **comprehensive-integration.spec.ts** - Full integration workflows
8. **complete-mvp-flow.spec.ts** - Complete MERN stack deployment
9. **data-persistence.spec.ts** - Data persistence across refreshes
10. **api-integration-suite.spec.ts** - All API endpoints
11. **mvp-flow.spec.ts** - Original MVP workflow test

### Test Helpers Created

1. **api-helpers.ts** - API interaction utilities
2. **page-objects.ts** - Page object models for UI
3. **test-helpers.ts** - Visual helpers and health checks

## 🎨 Test Execution Modes

### Development Mode
```bash
# Interactive UI (best for development)
npm run test:ui

# See browser (headed mode)
npm run test:headed

# Debug mode (step through tests)
npm run test:debug

# Generate test code
npm run test:codegen
```

### CI/CD Mode
```bash
# All tests
npm run test:all

# Specific browser
npm run test:chromium
npm run test:firefox
npm run test:webkit
```

## 📋 Test Categories

### @smoke Tests
**Purpose**: Quick validation that nothing is broken
- All pages load
- No critical errors
- Basic navigation

**Run**: On every commit
**Time**: < 5 seconds total

### @mvp Tests
**Purpose**: Core functionality validation
- Buildspec → Image → Container flow
- API integration
- Data persistence
- Mode switching

**Run**: On every PR
**Time**: < 2 minutes total

### Full Suite
**Purpose**: Comprehensive coverage
- All workflows
- Edge cases
- Error handling
- Cross-browser

**Run**: On releases, nightly
**Time**: ~5-10 minutes

## 🔄 Automated Test Execution

### On Every Commit
```bash
# Pre-commit hook runs smoke tests
npm run test:smoke
```

### On PR
- All @smoke tests must pass
- All @mvp tests must pass
- No new failures

### On Release
- Full test suite
- All browsers
- Performance checks

## 📝 Test Structure

### Page Object Pattern
```typescript
// Reusable page objects
const buildSpecPage = new BuildSpecPage(page)
await buildSpecPage.navigate()
await buildSpecPage.setSpecJson(json)
await buildSpecPage.clickBuild()
```

### API Helper Pattern
```typescript
// API interaction helpers
const apiHelper = new ApiHelper(page)
await apiHelper.setMode('test')
await apiHelper.buildImage(spec, tag)
await apiHelper.waitForImage(name)
```

### Test Organization
```
tests/
├── 00-setup.spec.ts          # Setup and health
├── all-pages.spec.ts          # Smoke tests
├── mode-switcher.spec.ts      # Mode tests
├── buildspec.spec.ts          # Build spec tests
├── images-api-integration.spec.ts
├── containers-api-integration.spec.ts
├── comprehensive-integration.spec.ts
├── complete-mvp-flow.spec.ts
├── data-persistence.spec.ts
└── api-integration-suite.spec.ts
```

## 🎯 What Gets Tested Automatically

### 1. Page Loading
- ✅ All 21 pages load without errors
- ✅ No 404s or 500s
- ✅ JavaScript errors caught

### 2. Mode Switching
- ✅ Mock/Test/Live mode switching
- ✅ Data clearing on mode switch
- ✅ UI updates correctly

### 3. Image Management
- ✅ Build image from buildspec
- ✅ Image appears in API
- ✅ Image appears in dashboard
- ✅ Image details displayed
- ✅ Image persistence

### 4. Container Management
- ✅ Create container
- ✅ Container appears in dashboard
- ✅ Start/Stop container
- ✅ Status updates
- ✅ Container persistence

### 5. Complete Workflows
- ✅ Buildspec → Image → Container → Verify
- ✅ Full MERN stack deployment
- ✅ Error handling
- ✅ Data persistence across refreshes

### 6. API Integration
- ✅ All API endpoints tested
- ✅ API ↔ UI synchronization
- ✅ Response format validation

## 🔧 Writing New Tests

### Template
```typescript
import { test, expect } from '@playwright/test'
import { ApiHelper } from '../utils/api-helpers'

test.describe('Feature Name', () => {
  let apiHelper: ApiHelper

  test.beforeEach(async ({ page }) => {
    apiHelper = new ApiHelper(page)
    await apiHelper.setMode('test')
  })

  test('@mvp should do something', async ({ page }) => {
    // Arrange
    // Act  
    // Assert
  })
})
```

### Best Practices
1. Use `@smoke` or `@mvp` tags
2. Use page objects for UI
3. Use API helpers for API calls
4. Wait appropriately (not fixed timeouts)
5. Clean up test data
6. Make tests independent

## 📊 Test Reports

### HTML Report
```bash
npm run test:report
```

### Test Results Location
- HTML: `test-results/playwright-report/index.html`
- Videos: `test-results/**/video.webm`
- Screenshots: `test-results/**/screenshot.png`
- JSON: `test-results/results.json`

## 🚦 Continuous Integration

### GitHub Actions
Tests run automatically:
- On push to main/develop
- On pull requests
- On manual trigger

### Local CI
```bash
# Run full suite locally
npm run test:all

# Run in CI mode (headless)
CI=true npm run test
```

## 📈 Coverage Metrics

### Current Coverage
- ✅ 21 pages (smoke tests)
- ✅ Mode switching
- ✅ Build spec workflow
- ✅ Images API integration
- ✅ Containers API integration
- ✅ Complete MVP flow
- ✅ Data persistence
- ✅ API endpoints

### Coverage Goals
- [ ] 100% of MVP features
- [ ] 100% of critical workflows
- [ ] All error scenarios
- [ ] Performance benchmarks
- [ ] Accessibility tests

## 🎓 Next Steps

1. **Run Tests Now**:
   ```bash
   cd e2e
   npm install
   npm run test:ui
   ```

2. **Watch Tests Run**:
   - See browser automation
   - Watch API calls
   - See test results

3. **Add More Tests**:
   - Follow test templates
   - Use page objects
   - Use API helpers

4. **Integrate with CI/CD**:
   - Tests run on every PR
   - Tests block merges if failing
   - Test reports available

## 🔄 Development Workflow

### When Adding Features:
1. Write test first (TDD)
2. Run test (should fail)
3. Implement feature
4. Run test (should pass)
5. Run full suite
6. Commit

### When Fixing Bugs:
1. Write test that reproduces bug
2. Fix bug
3. Test should pass
4. Run full suite
5. Commit

## 📄 Documentation

- **TEST_DRIVEN_DEVELOPMENT.md** - Complete TDD guide
- **e2e/README.md** - E2E test suite documentation
- **This file** - Automated testing guide

## 🎉 Benefits

✅ **Automated Testing**: No manual checking needed
✅ **Fast Feedback**: Tests run in minutes
✅ **Comprehensive**: Covers all major flows
✅ **Reliable**: Tests run consistently
✅ **Documentation**: Tests serve as documentation
✅ **Confidence**: Know what works and what doesn't

## 🚀 Run Tests Now!

```bash
cd e2e
npm install
npm run test:ui
```

Watch the magic happen! ✨

