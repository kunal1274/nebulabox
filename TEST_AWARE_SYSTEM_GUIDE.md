# Test-Aware System Guide

## 🎯 Overview

NebulaBox now has a **comprehensive test-aware system** that enables:
- ✅ Automatic test ID generation and management
- ✅ Test artifact collection (screenshots, videos, logs)
- ✅ Test execution tracking (pass/fail/unexecuted)
- ✅ Hierarchical test structure (Suite > Test Case > Test Run > Artifacts)
- ✅ Data flow testing
- ✅ Visual test annotations

## 📋 System Architecture

### Hierarchy Structure

```
Test Suite (e.g., "Authentication", "Container Lifecycle")
  └── Test Case (e.g., "TC-001: Login Flow", "TC-002: Container Creation")
      └── Test Run (e.g., "run-20251031-001")
          ├── screenshots/
          │   ├── step-01_before.png
          │   ├── step-01_after.png
          │   ├── step-02_before.png
          │   └── step-02_after.png
          ├── videos/
          │   └── run.mp4
          ├── logs/
          │   ├── step-01.json
          │   ├── step-02.json
          │   └── console.json
          ├── network/
          │   └── requests.json
          ├── reports/
          │   └── report.html
          ├── run-metadata.json
          └── execution-result.json
```

## 🔧 Components

### 1. Test ID System (`lib/test-ids.ts`)

Centralized test IDs for all interactive elements:

```typescript
import { TEST_IDS } from '@/lib/test-ids'

// Usage in components
<Button data-test-id={TEST_IDS.containers.create}>
  Create Container
</Button>

<Input data-test-id={TEST_IDS.buildspec.tagInput} />
```

### 2. Test-Aware Utilities (`lib/test-aware.ts`)

**Test Mode Detection**:
```typescript
import { TestMode } from '@/lib/test-aware'

// Check if in test mode
if (TestMode.isTestMode()) {
  // Disable animations
  TestMode.disableAnimations()
}

// Get session/run IDs
const sessionId = TestMode.getSessionId()
const runId = TestMode.getRunId()
```

**Test Recorder API**:
```typescript
// Manual annotation in browser console
window.__RECORDER__.annotate('Clicked create button', {
  selector: 'containers.create.button',
  timestamp: Date.now()
})
```

### 3. Test Artifact Manager (`e2e/utils/test-artifacts.ts`)

Manages all test artifacts:

```typescript
import { TestArtifactManager } from '../utils/test-artifacts'

const manager = new TestArtifactManager()

// Initialize run
const run = manager.initializeRun('suite-1', 'TC-001')

// Capture step with before/after screenshots
const step = await manager.addStep(page, 1, 'Navigate to page', 'goto /containers')
// ... perform action ...
await manager.completeStep(page, 1, 'passed')

// Finalize run
await manager.finalizeRun('passed')
```

### 4. Test Tracker (`e2e/utils/test-tracker.ts`)

Tracks test execution status:

```typescript
import { TestTracker } from '../utils/test-tracker'

const tracker = new TestTracker()

// Start test
const runId = await tracker.startTest('suite-1', 'TC-001')

// Record step
await tracker.recordStep(runId, 1, 'passed', undefined, 1500)

// Complete test
const result = await tracker.completeTest(runId, 'passed')
```

## 📝 Test ID Naming Convention

Pattern: `feature.action.element`

Examples:
- `containers.create.button` - Create container button
- `images.refresh.button` - Refresh images button
- `buildspec.validate.button` - Validate buildspec button
- `modeSwitcher.test.button` - Test mode button

## 🎨 Adding Test IDs to Components

### 1. Update Component Props

```typescript
// Button component already supports data-test-id
<Button data-test-id="containers.create.button">
  Create
</Button>
```

### 2. Use Test ID Constants

```typescript
import { TEST_IDS } from '@/lib/test-ids'

<Button data-test-id={TEST_IDS.containers.create}>
  Create
</Button>
```

### 3. Add to Inputs

```typescript
<Input 
  data-test-id={TEST_IDS.buildspec.tagInput}
  value={tag}
  onChange={(e) => setTag(e.target.value)}
/>
```

## 🧪 Writing Test-Aware Tests

### Example: Container Creation Test

```typescript
import { test, expect } from '@playwright/test'
import { TestArtifactManager } from '../utils/test-artifacts'
import { TestTracker } from '../utils/test-tracker'
import { TEST_IDS } from '../../../web/dashboard/src/lib/test-ids'

test('TC-001: Create container', async ({ page }) => {
  const artifactManager = new TestArtifactManager()
  const tracker = new TestTracker()

  // Initialize test run
  const runId = await tracker.startTest('container-lifecycle', 'TC-001')
  const run = artifactManager.initializeRun('container-lifecycle', 'TC-001')

  try {
    // Step 1: Navigate to containers page
    await artifactManager.addStep(page, 1, 'Navigate to containers', 'goto /containers')
    await page.goto('http://localhost:3001/containers')
    await page.waitForLoadState('networkidle')
    await artifactManager.completeStep(page, 1, 'passed')
    await tracker.recordStep(runId, 1, 'passed', undefined, 1000)

    // Step 2: Click create button
    await artifactManager.addStep(page, 2, 'Click create button', 'click create')
    const createButton = page.getByTestId(TEST_IDS.containers.create)
    await createButton.click()
    await page.waitForLoadState('networkidle')
    await artifactManager.completeStep(page, 2, 'passed')
    await tracker.recordStep(runId, 2, 'passed', undefined, 500)

    // Step 3: Fill form
    await artifactManager.addStep(page, 3, 'Fill container form', 'fill inputs')
    await page.getByTestId(TEST_IDS.createContainer.imageInput).fill('nginx:latest')
    await page.getByTestId(TEST_IDS.createContainer.nameInput).fill('test-container')
    await artifactManager.completeStep(page, 3, 'passed')
    await tracker.recordStep(runId, 3, 'passed', undefined, 800)

    // Step 4: Submit form
    await artifactManager.addStep(page, 4, 'Submit form', 'click create')
    await page.getByTestId(TEST_IDS.createContainer.createButton).click()
    await page.waitForTimeout(2000)
    
    // Verify container appears
    const container = page.getByTestId(TEST_IDS.containers.card).first()
    await expect(container).toBeVisible({ timeout: 10000 })
    
    await artifactManager.completeStep(page, 4, 'passed')
    await tracker.recordStep(runId, 4, 'passed', undefined, 2500)

    // Complete test
    const result = await tracker.completeTest(runId, 'passed')
    await artifactManager.finalizeRun('passed')

    expect(result.status).toBe('passed')
  } catch (error) {
    await tracker.recordStep(runId, 4, 'failed', String(error))
    await tracker.completeTest(runId, 'failed')
    await artifactManager.finalizeRun('failed')
    throw error
  }
})
```

## 📊 Test Execution Status

### Status Types

- **passed**: All steps completed successfully
- **failed**: One or more steps failed
- **skipped**: Test was skipped (not executed)
- **unexecuted**: Test initialized but not executed

### Viewing Results

```bash
# View HTML report
cat test-results/suite-1/TC-001/run-123/reports/report.html

# View execution result
cat test-results/suite-1/TC-001/run-123/execution-result.json

# View test case summary
cat test-results/suite-1/TC-001/test-case-report.json
```

## 🔄 Data Flow Testing

### API Integration Test Pattern

```typescript
test('Data flow: Build image → Create container → Verify', async ({ page }) => {
  // 1. Build image via API
  const buildResponse = await page.request.post('http://localhost:8081/api/buildspec/build', {
    data: { spec: buildspecContent, tag: 'test-image:latest' }
  })
  expect(buildResponse.ok()).toBeTruthy()

  // 2. Verify image in API
  const imagesResponse = await page.request.get('http://localhost:8081/api/images')
  const images = await imagesResponse.json()
  expect(images.find((img: any) => img.name === 'test-image')).toBeTruthy()

  // 3. Verify image in UI
  await page.goto('http://localhost:3001/images')
  await expect(page.getByText('test-image')).toBeVisible()

  // 4. Create container from image via API
  const containerResponse = await page.request.post('http://localhost:8081/api/containers/run', {
    data: { image: 'test-image:latest', name: 'test-container', detach: true }
  })
  expect(containerResponse.ok()).toBeTruthy()

  // 5. Verify container in UI
  await page.goto('http://localhost:3001/containers')
  await expect(page.getByText('test-container')).toBeVisible()
})
```

## 📁 Artifact Organization

### Directory Structure

```
test-results/
├── suite-1/
│   ├── TC-001/
│   │   ├── run-20251031-001/
│   │   │   ├── screenshots/
│   │   │   ├── videos/
│   │   │   ├── logs/
│   │   │   ├── network/
│   │   │   ├── reports/
│   │   │   ├── run-metadata.json
│   │   │   └── execution-result.json
│   │   ├── run-20251031-002/
│   │   └── test-case-report.json
│   ├── summaries/
│   └── suite-summary.json
```

## 🚀 Quick Start

### 1. Add Test IDs to Components

```typescript
import { TEST_IDS } from '@/lib/test-ids'

<Button data-test-id={TEST_IDS.containers.create}>
  Create
</Button>
```

### 2. Write Test-Aware Tests

```typescript
const createButton = page.getByTestId(TEST_IDS.containers.create)
await createButton.click()
```

### 3. Run Tests

```bash
cd e2e
npm run test:ui
```

### 4. View Artifacts

```bash
# Open HTML report
open test-results/suite-1/TC-001/run-123/reports/report.html
```

## 📋 Checklist

When adding new features:

- [ ] Add test IDs to all interactive elements
- [ ] Use `TEST_IDS` constants
- [ ] Add test IDs to buttons, inputs, forms
- [ ] Write test-aware E2E tests
- [ ] Capture before/after screenshots for each step
- [ ] Track test execution status
- [ ] Verify data flow (API ↔ UI)

## 🎓 Best Practices

1. **Always use test IDs**: Never rely on CSS classes or text content
2. **Use constants**: Import from `TEST_IDS` for consistency
3. **Capture artifacts**: Screenshots, videos, logs for every step
4. **Track status**: Record pass/fail for each step
5. **Test data flow**: Verify API → UI synchronization
6. **Organize artifacts**: Follow hierarchical structure

## 📄 Related Documentation

- `TEST_DRIVEN_DEVELOPMENT.md` - TDD approach
- `AUTOMATED_TESTING_GUIDE.md` - Automated testing guide
- `e2e/README.md` - E2E test suite docs

