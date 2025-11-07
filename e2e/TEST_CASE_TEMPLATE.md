# Test Case Template

## Test Case Structure

Use this template for documenting test cases with full artifact tracking.

### Metadata

```yaml
suiteId: "buildspec-flow"
testCaseId: "TC-001"
name: "Build image from buildspec.json"
description: "Validates complete buildspec → image → dashboard flow"
priority: "high"
tags: ["@mvp", "@buildspec", "@smoke"]
```

### Preconditions

- [ ] API server running on port 8081
- [ ] Dashboard running on port 3001
- [ ] Test mode enabled
- [ ] No existing test images

### Test Steps

#### Step 1: Navigate to BuildSpec Page

- **Action**: Navigate to `/buildspec`
- **Expected**: BuildSpec page loads with editor visible
- **Test ID**: `sidebar.buildspec.link`
- **Before Screenshot**: `step-01_before.png`
- **After Screenshot**: `step-01_after.png`
- **Expected Duration**: < 2s
- **Status**: `passed` | `failed` | `skipped`

#### Step 2: Load Example Buildspec

- **Action**: Click "Load Example" button
- **Expected**: Editor populated with example buildspec JSON
- **Test ID**: `buildspec.loadExample.button`
- **Before Screenshot**: `step-02_before.png`
- **After Screenshot**: `step-02_after.png`
- **Expected Duration**: < 1s
- **Status**: `passed` | `failed` | `skipped`

#### Step 3: Validate Buildspec

- **Action**: Click "Validate" button
- **Expected**: Validation result shown, Dockerfile generated
- **Test ID**: `buildspec.validate.button`
- **Before Screenshot**: `step-03_before.png`
- **After Screenshot**: `step-03_after.png`
- **Expected Duration**: < 3s
- **API Call**: `POST /api/buildspec/validate`
- **Status**: `passed` | `failed` | `skipped`

#### Step 4: Build Image

- **Action**: Click "Build" button
- **Expected**: Build completes, logs shown
- **Test ID**: `buildspec.build.button`
- **Before Screenshot**: `step-04_before.png`
- **After Screenshot**: `step-04_after.png`
- **Expected Duration**: < 5s
- **API Call**: `POST /api/buildspec/build`
- **Status**: `passed` | `failed` | `skipped`

#### Step 5: Verify Image in API

- **Action**: Check API for built image
- **Expected**: Image exists in `/api/images`
- **API Call**: `GET /api/images`
- **Expected Response**: Image with matching name/tag
- **Status**: `passed` | `failed` | `skipped`

#### Step 6: Verify Image in Dashboard

- **Action**: Navigate to Images page
- **Expected**: Built image appears in list
- **Test ID**: `images.card.container`
- **Before Screenshot**: `step-06_before.png`
- **After Screenshot**: `step-06_after.png`
- **Expected Duration**: < 2s
- **Status**: `passed` | `failed` | `skipped`

### Test Execution Result

```json
{
  "runId": "run-20251031-001",
  "testCaseId": "TC-001",
  "suiteId": "buildspec-flow",
  "status": "passed",
  "startTime": 1698768000000,
  "endTime": 1698768015000,
  "duration": 15000,
  "steps": [
    {
      "stepNumber": 1,
      "status": "passed",
      "duration": 1500
    },
    {
      "stepNumber": 2,
      "status": "passed",
      "duration": 800
    },
    {
      "stepNumber": 3,
      "status": "passed",
      "duration": 2500
    },
    {
      "stepNumber": 4,
      "status": "passed",
      "duration": 3500
    },
    {
      "stepNumber": 5,
      "status": "passed",
      "duration": 1000
    },
    {
      "stepNumber": 6,
      "status": "passed",
      "duration": 2000
    }
  ]
}
```

### Artifacts

- **Screenshots**: `test-results/buildspec-flow/TC-001/run-20251031-001/screenshots/`
- **Videos**: `test-results/buildspec-flow/TC-001/run-20251031-001/videos/`
- **Logs**: `test-results/buildspec-flow/TC-001/run-20251031-001/logs/`
- **Network**: `test-results/buildspec-flow/TC-001/run-20251031-001/network/`
- **Report**: `test-results/buildspec-flow/TC-001/run-20251031-001/reports/report.html`

### Notes

- All screenshots include step annotations (numbered circles)
- Network requests are captured with timing
- Console logs are saved for debugging
- Test can be replayed using artifacts

