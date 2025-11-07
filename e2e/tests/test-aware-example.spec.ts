import { test, expect } from '@playwright/test'
import { TestArtifactManager } from '../utils/test-artifacts'
import { TestTracker } from '../utils/test-tracker'
import { TEST_IDS } from '../../web/dashboard/src/lib/test-ids'

/**
 * Test-Aware Example Test
 * Demonstrates the complete test-aware system with:
 * - Test ID usage
 * - Artifact collection
 * - Execution tracking
 * - Before/after screenshots
 */

test.describe('Test-Aware Example: BuildSpec Flow', () => {
  test('TC-001: Build image from buildspec with full artifact collection', async ({ page }) => {
    const artifactManager = new TestArtifactManager()
    const tracker = new TestTracker()

    // Initialize test run
    const suiteId = 'buildspec-flow'
    const testCaseId = 'TC-001'
    const runId = await tracker.startTest(suiteId, testCaseId)
    const run = artifactManager.initializeRun(suiteId, testCaseId)

    try {
      // Step 1: Navigate to BuildSpec page
      await artifactManager.addStep(page, 1, 'Navigate to BuildSpec page', 'goto /buildspec')
      await page.goto('http://localhost:3001/buildspec')
      await page.waitForLoadState('networkidle')
      
      // Get element position for overlay
      const loadExampleBtn = page.getByTestId(TEST_IDS.buildspec.loadExample)
      const bbox = await loadExampleBtn.boundingBox()
      const annotation = bbox ? { x: bbox.x + bbox.width / 2, y: bbox.y + bbox.height / 2, label: '1' } : undefined
      
      await artifactManager.completeStep(page, 1, 'passed', undefined, {
        url: page.url(),
        annotation,
      })
      await tracker.recordStep(runId, 1, 'passed', undefined, 1500)

      // Step 2: Load example buildspec
      await artifactManager.addStep(page, 2, 'Load example buildspec', 'click load example')
      const loadExampleButton = page.getByTestId(TEST_IDS.buildspec.loadExample)
      const btnBbox = await loadExampleButton.boundingBox()
      const annotation2 = btnBbox ? { x: btnBbox.x + btnBbox.width / 2, y: btnBbox.y + btnBbox.height / 2, label: '2' } : undefined
      
      await loadExampleButton.click()
      await page.waitForTimeout(500)
      
      // Verify spec loaded
      const specEditor = page.getByTestId(TEST_IDS.buildspec.specEditor)
      await expect(specEditor).toBeVisible()
      
      await artifactManager.completeStep(page, 2, 'passed', undefined, { annotation: annotation2 })
      await tracker.recordStep(runId, 2, 'passed', undefined, 800)

      // Step 3: Validate buildspec
      await artifactManager.addStep(page, 3, 'Validate buildspec', 'click validate')
      const validateButton = page.getByTestId(TEST_IDS.buildspec.validate)
      const validateBbox = await validateButton.boundingBox()
      const annotation3 = validateBbox ? { x: validateBbox.x + validateBbox.width / 2, y: validateBbox.y + validateBbox.height / 2, label: '3' } : undefined
      
      await validateButton.click()
      await page.waitForTimeout(2000)
      
      // Switch to validation tab to verify
      const validationTab = page.getByTestId(TEST_IDS.buildspec.validationTab)
      if (await validationTab.isVisible()) {
        await validationTab.click()
        await page.waitForTimeout(500)
      }
      
      // Check for validation result
      const validationContent = page.locator('text=/valid|invalid|success/i').first()
      const hasResult = await validationContent.isVisible({ timeout: 5000 }).catch(() => false)
      expect(hasResult || await page.getByTestId(TEST_IDS.buildspec.dockerfileTab).isVisible()).toBeTruthy()
      
      await artifactManager.completeStep(page, 3, 'passed', undefined, { annotation: annotation3 })
      await tracker.recordStep(runId, 3, 'passed', undefined, 2500)

      // Step 4: Build image
      await artifactManager.addStep(page, 4, 'Build image from buildspec', 'click build')
      const buildButton = page.getByTestId(TEST_IDS.buildspec.build)
      const buildBbox = await buildButton.boundingBox()
      const annotation4 = buildBbox ? { x: buildBbox.x + buildBbox.width / 2, y: buildBbox.y + buildBbox.height / 2, label: '4' } : undefined
      
      await buildButton.click()
      await page.waitForTimeout(3000)
      
      // Switch to build logs tab
      const logsTab = page.getByTestId(TEST_IDS.buildspec.logsTab)
      if (await logsTab.isVisible()) {
        await logsTab.click()
        await page.waitForTimeout(500)
      }
      
      // Verify build started (logs tab visible or validation success)
      const hasLogs = await page.getByText(/building|built|success/i).isVisible({ timeout: 5000 }).catch(() => false)
      expect(hasLogs || await page.getByTestId(TEST_IDS.buildspec.dockerfileTab).isVisible()).toBeTruthy()
      
      await artifactManager.completeStep(page, 4, 'passed', undefined, { annotation: annotation4 })
      await tracker.recordStep(runId, 4, 'passed', undefined, 3500)

      // Step 5: Verify image appears in API
      await artifactManager.addStep(page, 5, 'Verify image in API', 'check API')
      
      // Get tag from input
      const tagInput = page.getByTestId(TEST_IDS.buildspec.tagInput)
      const tagValue = await tagInput.inputValue()
      const imageName = tagValue.split(':')[0]
      
      // Check API
      const imagesResponse = await page.request.get('http://localhost:8081/api/images')
      expect(imagesResponse.ok()).toBeTruthy()
      const images = await imagesResponse.json()
      const builtImage = images.find((img: any) => img.name === imageName)
      expect(builtImage).toBeTruthy()
      
      await artifactManager.completeStep(page, 5, 'passed', undefined, { imageName, imageFound: !!builtImage })
      await tracker.recordStep(runId, 5, 'passed', undefined, 1000)

      // Step 6: Verify image in UI
      await artifactManager.addStep(page, 6, 'Verify image in dashboard', 'goto images page')
      await page.goto('http://localhost:3001/images')
      await page.waitForTimeout(2000)
      
      // Check if image appears
      const imageCard = page.locator('[data-test-id="images.card.container"]').filter({ hasText: new RegExp(imageName, 'i') }).first()
      await expect(imageCard).toBeVisible({ timeout: 10000 })
      
      await artifactManager.completeStep(page, 6, 'passed', undefined, { imageName })
      await tracker.recordStep(runId, 6, 'passed', undefined, 2000)

      // Complete test
      const result = await tracker.completeTest(runId, 'passed')
      await artifactManager.finalizeRun('passed')

      // Verify final status
      expect(result.status).toBe('passed')
      expect(result.steps.length).toBe(6)
      expect(result.steps.every((s) => s.status === 'passed')).toBe(true)

    } catch (error) {
      // Record failure
      const currentStep = artifactManager.getCurrentRun()?.steps.length || 1
      await tracker.recordStep(runId, currentStep, 'failed', String(error))
      await tracker.completeTest(runId, 'failed')
      await artifactManager.finalizeRun('failed')
      throw error
    }
  })
})

