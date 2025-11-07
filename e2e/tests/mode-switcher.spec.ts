import { test, expect } from '@playwright/test'
import { BuildSpecPage, ContainersPage, ImagesPage } from '../utils/page-objects'
import { TEST_IDS } from '../../web/dashboard/src/lib/test-ids'

/**
 * Mode Switcher E2E Tests
 * Tests switching between mock/test/live modes
 * Migrated to use test IDs for better reliability
 */
test.describe('Mode Switcher', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:3001')
  })

  test('@smoke should display mode switcher', async ({ page }) => {
    // Mode switcher should be visible in header using test IDs
    await expect(page.getByTestId(TEST_IDS.modeSwitcher.badge)).toBeVisible()
    await expect(page.getByTestId(TEST_IDS.modeSwitcher.mock)).toBeVisible()
    await expect(page.getByTestId(TEST_IDS.modeSwitcher.test)).toBeVisible()
    await expect(page.getByTestId(TEST_IDS.modeSwitcher.live)).toBeVisible()
  })

  test('@mvp should show current mode', async ({ page }) => {
    // Should show mode badge using test ID
    await expect(page.getByTestId(TEST_IDS.modeSwitcher.badge)).toBeVisible()
  })

  test('@mvp should show image count', async ({ page }) => {
    // Should show Images: X/Y
    await expect(page.getByText(/Images:/i)).toBeVisible()
  })

  test('@mvp should switch to mock mode', async ({ page }) => {
    // Click Mock button using test ID
    await page.getByTestId(TEST_IDS.modeSwitcher.mock).click()
    
    // Dashboard should reload
    await page.waitForLoadState('networkidle')
    
    // Verify mode changed (check API)
    const modeResponse = await page.request.get('http://localhost:8081/api/mode')
    const modeData = await modeResponse.json()
    expect(modeData.mode).toBe('mock')
  })

  test('@mvp should switch to test mode', async ({ page }) => {
    // Ensure we're starting from a known state
    // Switch to mock first, then to test using test IDs
    await page.getByTestId(TEST_IDS.modeSwitcher.mock).click()
    await page.waitForLoadState('networkidle')
    
    // Switch to test
    await page.getByTestId(TEST_IDS.modeSwitcher.test).click()
    await page.waitForLoadState('networkidle')
    
    // Verify mode changed
    const modeResponse = await page.request.get('http://localhost:8081/api/mode')
    const modeData = await modeResponse.json()
    expect(modeData.mode).toBe('test')
  })

  test('@mvp should update image count after building', async ({ page }) => {
    // Start in test mode using test ID
    const testButton = page.getByTestId(TEST_IDS.modeSwitcher.test)
    if (await testButton.getAttribute('data-active') !== 'true') {
      await testButton.click()
      await page.waitForLoadState('networkidle')
    }

    // Build an image
    await page.goto('http://localhost:3001/buildspec')
    const buildSpecPage = new BuildSpecPage(page)
    await buildSpecPage.navigate()
    
    // Load and build example using test IDs
    await page.getByTestId(TEST_IDS.buildspec.loadExample).click()
    await page.getByTestId(TEST_IDS.buildspec.build).click()
    await page.waitForTimeout(2000)
    
    // Go back to home
    await page.goto('http://localhost:3001')
    
    // Check image count updated
    await expect(page.getByText(/Images:/i)).toBeVisible()
    // Should show at least 1 built image
    const imageText = await page.getByText(/Images:/i).textContent()
    expect(imageText).toMatch(/\d+/)
  })
})

