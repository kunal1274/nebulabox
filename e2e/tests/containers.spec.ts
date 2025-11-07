import { test, expect } from '@playwright/test';
import { ContainersPage } from '../utils/page-objects';
import { ScreenshotHelper } from '../utils/test-helpers';
import { TEST_IDS } from '../../web/dashboard/src/lib/test-ids';

/**
 * Container Management E2E Tests
 * Migrated to use test IDs for better reliability
 */
test.describe('Container Lifecycle', () => {
  let containersPage: ContainersPage;
  let screenshotHelper: ScreenshotHelper;

  test.beforeEach(async ({ page }) => {
    containersPage = new ContainersPage(page);
    screenshotHelper = new ScreenshotHelper(page);
    await containersPage.navigate();
    await screenshotHelper.capture('containers-01-page-loaded');
  });

  test('@smoke should load containers page', async ({ page }) => {
    await expect(page).toHaveURL(/.*containers/);
    // Use test ID to verify containers list is visible
    await expect(page.getByTestId(TEST_IDS.containers.list)).toBeVisible({ timeout: 5000 });
  });

  test('@mvp should display container list', async ({ page }) => {
    // Wait for page to load and verify container list using test ID
    await page.waitForTimeout(1000);
    await expect(page.getByTestId(TEST_IDS.containers.list)).toBeVisible();
    await screenshotHelper.capture('containers-02-list-visible');
  });

  test('@mvp should navigate to create container', async ({ page }) => {
    // Use test ID to click create button
    await page.getByTestId(TEST_IDS.containers.create).click();
    await screenshotHelper.capture('containers-03-create-clicked');
    
    // Should navigate to create page or show form/modal
    // Check if URL changed or form appeared using test IDs
    await page.waitForTimeout(1000);
    const url = page.url();
    const hasCreatePage = url.includes('/new') || url.includes('/create');
    const hasForm = await page.getByTestId(TEST_IDS.createContainer.imageInput).isVisible({ timeout: 3000 }).catch(() => false);
    
    expect(hasCreatePage || hasForm).toBeTruthy();
    await screenshotHelper.capture('containers-04-create-form');
  });
});

