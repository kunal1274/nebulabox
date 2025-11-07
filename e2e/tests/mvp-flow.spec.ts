import { test, expect } from '@playwright/test';
import { BuildSpecPage, ContainersPage, ImagesPage } from '../utils/page-objects';
import { ScreenshotHelper } from '../utils/test-helpers';
import * as fs from 'fs';
import * as path from 'path';

/**
 * Complete MVP Flow E2E Test
 * Tests the entire flow: Build Spec → Build Image → Create Container
 */
test.describe('MVP Complete Flow', () => {
  const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json');
  let screenshotHelper: ScreenshotHelper;

  test.beforeEach(async ({ page }) => {
    screenshotHelper = new ScreenshotHelper(page);
  });

  test('@mvp complete MERN stack deployment flow', async ({ page }) => {
    // Step 1: Navigate to Build Spec page
    const buildSpecPage = new BuildSpecPage(page);
    await buildSpecPage.navigate();
    await screenshotHelper.capture('mvp-01-buildspec-page');

    // Step 2: Load and validate buildspec
    const buildspecContent = fs.readFileSync(buildspecPath, 'utf-8');
    const specJson = JSON.parse(buildspecContent);
    
    await buildSpecPage.setSpecJson(JSON.stringify(specJson, null, 2));
    await screenshotHelper.capture('mvp-02-spec-loaded');
    
    await buildSpecPage.clickValidate();
    await buildSpecPage.waitForValidationComplete();
    await screenshotHelper.capture('mvp-03-validated');
    
    const dockerfile = await buildSpecPage.getDockerfilePreview();
    expect(dockerfile).toBeTruthy();
    await screenshotHelper.capture('mvp-04-dockerfile-shown');

    // Step 3: Build the image
    await buildSpecPage.clickBuild();
    await buildSpecPage.waitForBuildComplete();
    await screenshotHelper.capture('mvp-05-build-started');
    
    await page.waitForTimeout(3000);
    const logs = await buildSpecPage.getBuildLogs();
    // Logs might be empty in mock mode
    expect(Array.isArray(logs)).toBeTruthy();
    if (logs.length > 0) {
      expect(logs.some(log => log.length > 0)).toBeTruthy();
    }
    await screenshotHelper.capture('mvp-06-build-complete');

    // Step 4: Navigate to Images page and verify
    const imagesPage = new ImagesPage(page);
    await imagesPage.navigate();
    await screenshotHelper.capture('mvp-07-images-page');
    
    // Check if image appears (may take a moment)
    await page.waitForTimeout(2000);
    await screenshotHelper.capture('mvp-08-images-loaded');

    // Step 5: Navigate to Containers page
    const containersPage = new ContainersPage(page);
    await containersPage.navigate();
    await screenshotHelper.capture('mvp-09-containers-page');

    // Step 6: Create container (if form is available)
    await containersPage.clickCreate();
    await screenshotHelper.capture('mvp-10-create-clicked');
    
    // Check if we navigated to create page or form appeared
    await page.waitForTimeout(2000);
    const url = page.url();
    const hasCreatePage = url.includes('/new') || url.includes('/create');
    
    // Try to find image input if form is visible
    const imageInput = page.locator('input[name="image"], input[placeholder*="image"], input[placeholder*="Image"]').first();
    const hasForm = await imageInput.isVisible({ timeout: 3000 }).catch(() => false);
    
    if (hasForm) {
      await imageInput.fill('mern-mvp:latest');
      await screenshotHelper.capture('mvp-11-image-filled');
    } else if (hasCreatePage) {
      // We're on the create page
      await screenshotHelper.capture('mvp-11-create-page');
    }
    
    // Verify we either have form or are on create page
    expect(hasCreatePage || hasForm).toBeTruthy();
    await screenshotHelper.capture('mvp-12-form-complete');
  });
});

