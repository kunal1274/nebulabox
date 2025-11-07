# NebulaBox E2E Test Suite

Comprehensive end-to-end testing for NebulaBox using Playwright.

## 🎯 Test Categories

### @smoke - Quick Smoke Tests
Fast tests (< 5s each) that verify pages load without errors.
```bash
npm run test:smoke
```

### @mvp - MVP Feature Tests
Core functionality tests for MVP features.
```bash
npm run test:mvp
```

## 🚀 Quick Start

### Install Dependencies
```bash
cd e2e
npm install
```

### Run All Tests
```bash
npm run test
```

### Run Tests with UI (Recommended)
```bash
npm run test:ui
```

### Run Specific Test File
```bash
npx playwright test tests/complete-mvp-flow.spec.ts
```

## 📊 Test Coverage

### Current Test Files

1. **00-setup.spec.ts** - Health checks and setup
2. **all-pages.spec.ts** - Smoke tests for all pages
3. **mode-switcher.spec.ts** - Mode switching tests
4. **buildspec.spec.ts** - Build specification tests
5. **images-api-integration.spec.ts** - Images API integration
6. **containers-api-integration.spec.ts** - Containers API integration
7. **comprehensive-integration.spec.ts** - Full integration tests
8. **complete-mvp-flow.spec.ts** - Complete MVP workflow
9. **data-persistence.spec.ts** - Data persistence tests
10. **mvp-flow.spec.ts** - Original MVP flow test

## 🔧 Test Helpers

### API Helper
```typescript
import { ApiHelper } from './utils/api-helpers'

const apiHelper = new ApiHelper(page)
await apiHelper.setMode('test')
await apiHelper.buildImage(spec, tag)
await apiHelper.createContainer(options)
```

### Page Objects
```typescript
import { BuildSpecPage, ImagesPage, ContainersPage } from './utils/page-objects'

const buildSpecPage = new BuildSpecPage(page)
await buildSpecPage.navigate()
await buildSpecPage.setSpecJson(json)
await buildSpecPage.clickBuild()
```

## 📝 Writing New Tests

See `TEST_DRIVEN_DEVELOPMENT.md` for complete guide.

### Quick Template
```typescript
import { test, expect } from '@playwright/test'
import { ApiHelper } from '../utils/api-helpers'

test.describe('Feature Name', () => {
  test('@mvp should do something', async ({ page }) => {
    const apiHelper = new ApiHelper(page)
    // Test implementation
  })
})
```

## 🎨 Test Execution

### Development
```bash
# Interactive UI mode
npm run test:ui

# Headed mode (see browser)
npm run test:headed

# Debug mode
npm run test:debug
```

### CI/CD
```bash
# All tests
npm run test:all

# Specific browser
npm run test:chromium
```

## 📈 Test Reports

### View HTML Report
```bash
npm run test:report
```

### Test Results Location
- HTML Report: `test-results/`
- Videos: `test-results/**/video.webm`
- Screenshots: `test-results/**/screenshot.png`

## 🐛 Debugging

### Common Issues

1. **API Not Running**
   - Check `playwright.config.ts` webServer config
   - Verify ports 8081 and 3001 are available

2. **Tests Failing Intermittently**
   - Increase wait times
   - Use better selectors (getByRole, getByText)
   - Wait for networkidle

3. **Selector Issues**
   - Use page objects
   - Prefer semantic selectors (getByRole, getByLabel)
   - Avoid class-based selectors

## 🔄 Continuous Testing

Tests run automatically:
- On every commit (smoke tests)
- On every PR (all tests)
- Before releases (full suite)

## 📋 Test Checklist

For each new feature:
- [ ] Smoke test added
- [ ] MVP test added (if MVP feature)
- [ ] API integration test
- [ ] Error handling test
- [ ] Page object created (if new page)
