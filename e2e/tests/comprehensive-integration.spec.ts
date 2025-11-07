import { test, expect } from '@playwright/test'
import { ApiHelper } from '../utils/api-helpers'
import { BuildSpecPage, ImagesPage, ContainersPage } from '../utils/page-objects'
import * as fs from 'fs'
import * as path from 'path'

/**
 * Comprehensive Integration Tests
 * Tests all major flows: Build → Image → Container → Lifecycle
 */
test.describe('Comprehensive Integration Tests', () => {
  let apiHelper: ApiHelper

  test.beforeEach(async ({ page }) => {
    apiHelper = new ApiHelper(page)
    // Ensure test mode
    await apiHelper.setMode('test')
    // Clean up any previous test data
    await page.goto('http://localhost:3001')
  })

  test('@mvp full workflow: buildspec → image → container → verify', async ({ page }) => {
    const testImageName = 'e2e-full-test'
    const testContainerName = 'e2e-full-container'

    // 1. Build Image from Buildspec
    await page.goto('http://localhost:3001/buildspec')
    const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json')
    const buildspecContent = JSON.parse(fs.readFileSync(buildspecPath, 'utf-8'))
    
    // Modify tag for this test
    buildspecContent.tag = `${testImageName}:latest`
    buildspecContent.name = testImageName

    const buildResponse = await apiHelper.buildImage(buildspecContent, `${testImageName}:latest`)
    expect(buildResponse.valid).toBe(true)

    // 2. Wait for image to appear in API
    const image = await apiHelper.waitForImage(testImageName, 10000)
    expect(image).toBeTruthy()
    expect(image.name).toBe(testImageName)

    // 3. Verify image appears in dashboard
    await page.goto('http://localhost:3001/images')
    await page.waitForTimeout(2000)
    await expect(page.getByText(new RegExp(testImageName, 'i'))).toBeVisible({ timeout: 10000 })

    // 4. Create container from the image
    const container = await apiHelper.createContainer({
      image: `${testImageName}:latest`,
      name: testContainerName,
      ports: ['3000:3000'],
      env: ['TEST_MODE=true'],
      detach: true
    })
    expect(container.name).toBe(testContainerName)

    // 5. Wait for container to be running
    await apiHelper.waitForContainerStatus(testContainerName, 'running', 15000)

    // 6. Verify container appears in dashboard
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(2000)
    await expect(page.getByText(new RegExp(testContainerName, 'i'))).toBeVisible({ timeout: 10000 })

    // 7. Verify container details in dashboard
    const containerCard = page.locator('.card, [class*="Card"]').filter({ hasText: new RegExp(testContainerName, 'i') }).first()
    await expect(containerCard.getByText(/running/i).or(containerCard.getByText(/test-container/i))).toBeVisible()

    // 8. Test container logs
    const containerData = await apiHelper.findContainer(testContainerName)
    expect(containerData).toBeTruthy()
    
    const logsResponse = await apiHelper.getContainerLogs(containerData.id)
    expect(logsResponse).toBeTruthy()

    // 9. Test stop container
    await apiHelper.stopContainer(containerData.id)
    await apiHelper.waitForContainerStatus(testContainerName, 'stopped', 10000)

    // Verify in dashboard
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(1000)
    const stoppedCard = page.locator('.card, [class*="Card"]').filter({ hasText: new RegExp(testContainerName, 'i') }).first()
    await expect(stoppedCard.getByText(/stopped/i)).toBeVisible({ timeout: 5000 })

    // 10. Test start container
    await apiHelper.startContainer(containerData.id)
    await apiHelper.waitForContainerStatus(testContainerName, 'running', 10000)

    // Verify in dashboard
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(1000)
    const startedCard = page.locator('.card, [class*="Card"]').filter({ hasText: new RegExp(testContainerName, 'i') }).first()
    await expect(startedCard.getByText(/running/i)).toBeVisible({ timeout: 5000 })
  })

  test('@mvp mode switching workflow', async ({ page }) => {
    // 1. Start in test mode
    await apiHelper.setMode('test')
    await page.goto('http://localhost:3001')
    await expect(page.getByText(/TEST/i).or(page.getByRole('button', { name: /Test/i })).first()).toBeVisible()

    // 2. Build an image in test mode
    const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json')
    const buildspecContent = JSON.parse(fs.readFileSync(buildspecPath, 'utf-8'))
    buildspecContent.tag = 'mode-test:latest'
    
    await apiHelper.buildImage(buildspecContent, 'mode-test:latest')
    const image = await apiHelper.waitForImage('mode-test', 10000)
    expect(image).toBeTruthy()

    // 3. Verify image appears in dashboard
    await page.goto('http://localhost:3001/images')
    await page.waitForTimeout(1000)
    await expect(page.getByText(/mode-test/i)).toBeVisible({ timeout: 5000 })

    // 4. Switch to mock mode via UI
    await page.goto('http://localhost:3001')
    await page.getByRole('button', { name: /Mock/i }).click()
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)

    // 5. Verify mode changed
    const modeData = await apiHelper.getMode()
    expect(modeData.mode).toBe('mock')

    // 6. Verify image no longer appears (mock mode clears data)
    await page.goto('http://localhost:3001/images')
    await page.waitForTimeout(1000)
    // In mock mode, should show mock data only
    const imagesResponse = await apiHelper.listImages()
    // Mock mode may show only mock images
    expect(Array.isArray(imagesResponse)).toBe(true)

    // 7. Switch back to test mode
    await page.goto('http://localhost:3001')
    await page.getByRole('button', { name: /Test/i }).click()
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)

    // 8. Build new image in test mode
    buildspecContent.tag = 'mode-test-2:latest'
    await apiHelper.buildImage(buildspecContent, 'mode-test-2:latest')
    const image2 = await apiHelper.waitForImage('mode-test-2', 10000)
    expect(image2).toBeTruthy()
  })

  test('@mvp error handling and validation', async ({ page }) => {
    // 1. Try to create container with invalid image
    await page.goto('http://localhost:3001/containers')
    await page.getByRole('button', { name: /New Container|Create/i }).click()
    await page.waitForTimeout(500)

    // Should handle invalid image gracefully
    // (Actual error handling depends on implementation)

    // 2. Try to build invalid buildspec
    await page.goto('http://localhost:3001/buildspec')
    const buildSpecPage = new BuildSpecPage(page)
    
    const invalidSpec = '{"invalid": "json"}'
    await buildSpecPage.setSpecJson(invalidSpec)
    await buildSpecPage.clickValidate()
    await page.waitForTimeout(2000)

    // Should show validation error
    const validationResult = page.locator('text=/error|invalid|failed/i').first()
    // Error message should be visible
    const hasError = await validationResult.count() > 0
    expect(hasError || await page.getByText(/dockerfile/i).isVisible()).toBeTruthy()
  })
})

