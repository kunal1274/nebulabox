# Test-Driven Development Guide - NebulaBox

## 🎯 Philosophy: Testing as First-Class Citizen

**Every feature must have tests before it's considered complete.**

## 📋 Test Coverage Strategy

### 1. **Smoke Tests** (`@smoke`)
Quick validation that pages load without errors
- All pages render
- No critical JavaScript errors
- Basic navigation works

### 2. **MVP Tests** (`@mvp`)
Core functionality tests for MVP features
- Buildspec → Image → Container flow
- API integration
- Data persistence
- Mode switching

### 3. **Integration Tests**
Tests that verify multiple components work together
- Full workflows
- API ↔ UI synchronization
- Data flow end-to-end

### 4. **E2E Tests**
Complete user journey tests
- Real user scenarios
- Complete workflows
- Error handling

## 🏗️ Test Structure

```
e2e/
├── tests/
│   ├── 00-setup.spec.ts          # Health checks
│   ├── all-pages.spec.ts          # Smoke tests for all pages
│   ├── mode-switcher.spec.ts      # Mode switching tests
│   ├── buildspec.spec.ts          # Build spec tests
│   ├── images-api-integration.spec.ts  # Images API tests
│   ├── containers-api-integration.spec.ts  # Containers API tests
│   ├── comprehensive-integration.spec.ts  # Full integration
│   ├── complete-mvp-flow.spec.ts  # Complete MVP workflow
│   └── data-persistence.spec.ts  # Data persistence tests
├── utils/
│   ├── page-objects.ts           # Page object models
│   ├── api-helpers.ts             # API helper functions
│   └── test-helpers.ts            # Visual helpers, health checks
└── playwright.config.ts           # Playwright configuration
```

## 🔄 Development Workflow

### When Adding a New Feature:

1. **Write Test First** (TDD Approach)
   ```typescript
   test('@mvp should create new feature', async ({ page }) => {
     // Test the feature before implementing
   })
   ```

2. **Run Test** (Should fail initially)
   ```bash
   npm run test:mvp
   ```

3. **Implement Feature**
   - Make test pass
   - Ensure API integration works
   - Verify UI updates correctly

4. **Run Full Test Suite**
   ```bash
   npm run test
   ```

5. **Add to CI/CD**
   - Tests run automatically on PR
   - Tests run on commit (optional)

## 📊 Test Categories

### @smoke Tests
- Fast (< 5 seconds each)
- Check page loads
- Check no critical errors
- Should run on every commit

### @mvp Tests
- Moderate speed (< 30 seconds each)
- Test core functionality
- Verify API integration
- Should run on every PR

### Full Suite
- Comprehensive coverage
- May take longer
- Run on CI/CD pipeline
- Run before releases

## 🧪 Running Tests

### Development Mode
```bash
# Run all tests
npm run test

# Run with UI (recommended for development)
npm run test:ui

# Run only smoke tests (fast)
npm run test:smoke

# Run only MVP tests
npm run test:mvp

# Run in headed mode (see browser)
npm run test:headed

# Debug mode
npm run test:debug
```

### CI/CD Mode
```bash
# All tests in CI
npm run test:all

# Specific browser
npm run test:chromium
```

## 🎯 Test Coverage Goals

### Current Coverage
- ✅ Page smoke tests
- ✅ Mode switcher
- ✅ Build spec validation
- ✅ Images API integration
- ✅ Containers API integration
- ✅ Complete MVP flow
- ✅ Data persistence

### Target Coverage
- [ ] All API endpoints
- [ ] All UI interactions
- [ ] Error scenarios
- [ ] Edge cases
- [ ] Performance tests
- [ ] Accessibility tests

## 🔧 Test Helper Functions

### API Helper
```typescript
const apiHelper = new ApiHelper(page)
await apiHelper.setMode('test')
await apiHelper.buildImage(spec, tag)
await apiHelper.createContainer(options)
await apiHelper.waitForImage(name)
await apiHelper.waitForContainerStatus(name, status)
```

### Page Objects
```typescript
const buildSpecPage = new BuildSpecPage(page)
await buildSpecPage.setSpecJson(json)
await buildSpecPage.clickBuild()

const imagesPage = new ImagesPage(page)
await imagesPage.findImage('my-app')

const containersPage = new ContainersPage(page)
await containersPage.createContainer()
```

## 📝 Writing New Tests

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
1. **Arrange-Act-Assert** pattern
2. **Use page objects** for UI interactions
3. **Use API helpers** for API calls
4. **Wait appropriately** (don't use fixed timeouts)
5. **Clean up** test data in afterEach if needed
6. **Tag tests** with @smoke or @mvp
7. **Make tests independent** (no test order dependencies)

## 🚀 Continuous Testing

### Pre-commit Hooks
```bash
# Run smoke tests before commit
npm run test:smoke
```

### PR Checks
- All @smoke tests must pass
- All @mvp tests must pass
- No new test failures

### Release Checks
- Full test suite passes
- All browsers tested
- Performance benchmarks met

## 📈 Test Metrics

### Coverage Goals
- **Smoke Tests**: 100% of pages
- **MVP Tests**: 100% of MVP features
- **Integration Tests**: All critical workflows
- **E2E Tests**: All user journeys

### Quality Metrics
- **Test Execution Time**: < 5 minutes for full suite
- **Test Reliability**: > 95% pass rate
- **Test Coverage**: > 80% of critical paths

## 🐛 Debugging Tests

### View Test Results
```bash
# HTML report
npm run test:report

# Watch mode
npm run test:ui
```

### Common Issues
1. **Test Flakiness**: Increase wait times, use better selectors
2. **API Not Ready**: Check webServer config in playwright.config.ts
3. **Selector Issues**: Use page objects, use getByRole when possible
4. **Timing Issues**: Wait for networkidle, wait for specific elements

## 🎓 Learning Resources

- Playwright Docs: https://playwright.dev
- Best Practices: https://playwright.dev/docs/best-practices
- Page Object Model: https://playwright.dev/docs/pom

## 📋 Test Checklist for New Features

When adding a new feature, ensure:

- [ ] Smoke test for new page/component
- [ ] MVP test for core functionality
- [ ] API integration test
- [ ] Error handling test
- [ ] UI interaction test
- [ ] Data persistence test (if applicable)
- [ ] Mode compatibility test (if applicable)

## 🔄 Maintenance

- Review test failures weekly
- Update tests when UI changes
- Add tests for bug fixes
- Remove obsolete tests
- Keep test helpers updated

