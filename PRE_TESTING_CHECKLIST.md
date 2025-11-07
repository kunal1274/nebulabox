# Pre-Testing Checklist

**Date:** $(date)
**Status:** ⚠️ Review Required

## ✅ Test-Aware Infrastructure

### 1. Test ID System
- ✅ `test-ids.ts` - All constants defined (260+ test IDs)
- ✅ Pattern: `feature.action.element`
- ✅ Centralized constants for maintainability

### 2. UI Components Support
- ✅ `Button` - Supports `data-test-id` prop
- ✅ `Input` - Supports `data-test-id` prop
- ✅ `Textarea` - Supports `data-test-id` prop
- ✅ `TabsTrigger` - Supports `data-test-id` (via Radix props)
- ⚠️ Other components - May need verification

### 3. API Test-Aware Integration
- ✅ `api-test-aware.ts` - Created and functional
- ✅ Integrated into `ApiClient.request()`
- ✅ Auto-initialized in `main.tsx`
- ✅ Tracks: method, endpoint, request, response, errors, duration
- ✅ Exposes `window.__API_TRACKER__` in test mode
- ✅ Exposes `window.__RECORDER__` for test recording

### 4. Test-Aware Utilities
- ✅ `test-aware.ts` - Test mode detection
- ✅ Session/Run ID generation
- ✅ Animation disabling
- ✅ Initialized in `main.tsx`

## ✅ Test Infrastructure

### 1. Playwright Setup
- ✅ `playwright.config.ts` - Configured
- ✅ Web server auto-start (API + Dashboard)
- ✅ Video recording on failure
- ✅ Screenshots on failure
- ✅ JSON reporter for test results

### 2. E2E Tests
- ✅ `00-setup.spec.ts` - Health checks
- ✅ `buildspec.spec.ts` - Build spec flow
- ✅ `containers.spec.ts` - Container lifecycle
- ✅ `mvp-flow.spec.ts` - Full MVP deployment
- ✅ `all-pages.spec.ts` - Smoke tests (fixed timeouts)

### 3. Test Utilities
- ✅ `test-helpers.ts` - Visual helpers, health checkers
- ✅ `page-objects.ts` - Page object models
- ⚠️ Uses `getByRole` - May need migration to `getByTestId`

### 4. Test Fixer Tools
- ✅ `test-fixer.ts` - CLI tool for test management
- ✅ `monitor-tests.sh` - Shell script for monitoring
- ✅ `quick-fix.sh` - Quick fix script
- ⚠️ Requires `tsx` package (installed)

## ⚠️ Page Coverage

### ✅ Pages WITH Test IDs (5/25 - 20%)
1. ✅ Containers.tsx
2. ✅ Images.tsx
3. ✅ BuildSpec.tsx
4. ✅ CreateContainer.tsx
5. ✅ Registry.tsx
6. ✅ ModeSwitcher.tsx (component)

### ❌ Pages WITHOUT Test IDs (20/25 - 80%)

**HIGH PRIORITY:**
- Security.tsx
- Orchestrator.tsx
- Runtime.tsx
- AIOps.tsx
- ContainerGroups.tsx
- Composition.tsx
- Templates.tsx
- SharedRuntime.tsx
- Snapshots.tsx
- EphemeralRuntime.tsx

**MEDIUM PRIORITY:**
- Dashboard.tsx
- Monitor.tsx
- Networks.tsx
- Services.tsx
- Teams.tsx
- Tenants.tsx
- Settings.tsx

**LOW PRIORITY:**
- ContainerLogs.tsx
- ContainerEnv.tsx
- Logs.tsx
- Performance.tsx

## ⚠️ Potential Issues

### 1. Test ID Usage in E2E Tests
- ⚠️ Current tests use `getByRole`, `getByText`, `getByPlaceholderText`
- ⚠️ Should migrate to `getByTestId` for stability
- ⚠️ Page objects may need updates

### 2. Test Fixer Commands
- ⚠️ `npm run test:check` - Requires test results JSON
- ⚠️ `npm run test:fix` - Requires test results JSON
- ⚠️ Will work after first test run

### 3. Incomplete Page Coverage
- ⚠️ 80% of pages still need test IDs
- ⚠️ Can still test with current 5 pages
- ⚠️ But full coverage requires all pages

## ✅ Ready for Testing

### What Works Now:
1. ✅ Test infrastructure is complete
2. ✅ 5 pages have test IDs (core workflows)
3. ✅ API tracking is active
4. ✅ Playwright config is ready
5. ✅ Health check tests exist
6. ✅ MVP flow test exists

### What Needs Work:
1. ⚠️ Add test IDs to remaining 20 pages
2. ⚠️ Update E2E tests to use `getByTestId` instead of `getByRole`
3. ⚠️ Update page objects to use test IDs
4. ⚠️ Fix test fixer tools (needs test results)

## 🎯 Recommendation

### Option 1: Test Now (Partial Coverage)
- ✅ Can test with 5 pages that have test IDs
- ✅ Core workflows (Containers, Images, BuildSpec, CreateContainer, Registry)
- ✅ Good for MVP validation
- ⚠️ Limited to tested pages

### Option 2: Complete First (Full Coverage)
- ⚠️ Add test IDs to remaining 20 pages first
- ✅ Then run full test suite
- ✅ Better for comprehensive testing
- ⚠️ Takes more time upfront

### Option 3: Hybrid Approach
- ✅ Test current 5 pages now
- ✅ Gradually add test IDs to other pages
- ✅ Expand test coverage incrementally

## 📋 Action Items Before Full Testing

1. **Decide on testing scope:**
   - [ ] Test with 5 pages now (partial)
   - [ ] Add test IDs to all pages first (complete)
   - [ ] Hybrid approach

2. **Update E2E tests (if testing now):**
   - [ ] Migrate to `getByTestId` where available
   - [ ] Keep `getByRole` as fallback for untested pages
   - [ ] Update page objects

3. **Verify test fixer:**
   - [ ] Run a test first to generate results JSON
   - [ ] Then test fixer commands will work
   - [ ] Or fix test fixer to handle missing results

4. **Check dependencies:**
   - [ ] `tsx` installed in e2e/package.json ✅
   - [ ] Playwright installed ✅
   - [ ] All scripts defined ✅

## 🚀 Ready Status

**Current Status:** ⚠️ **PARTIALLY READY**

- **Infrastructure:** ✅ 100% Complete
- **Test IDs:** ⚠️ 20% Complete (5/25 pages)
- **E2E Tests:** ✅ Exists but may need updates
- **Test Tools:** ⚠️ Need test results to work

**Recommendation:** You can start testing with the 5 pages that have test IDs, or complete all pages first for full coverage.

