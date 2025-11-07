import { test, expect } from '@playwright/test'
import { ImagesPage } from '../utils/page-objects'
import * as fs from 'fs'
import * as path from 'path'

/**
 * Images API Integration E2E Tests
 * Tests that built images appear in dashboard after API build
 */
test.describe('Images API Integration', () => {
  test.beforeEach(async ({ page }) => {
    // Ensure we're in test mode
    await page.goto('http://localhost:3001')
    const modeResponse = await page.request.put('http://localhost:8081/api/mode', {
      data: { mode: 'test' }
    })
    expect(modeResponse.ok()).toBeTruthy()
  })

  test('@mvp should build image and see it in dashboard', async ({ page }) => {
    const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json')
    const buildspecContent = JSON.parse(fs.readFileSync(buildspecPath, 'utf-8'))

    // Step 1: Build image via API
    const buildResponse = await page.request.post('http://localhost:8081/api/buildspec/build', {
      data: {
        spec: buildspecContent,
        tag: 'e2e-test-image:latest'
      }
    })
    
    expect(buildResponse.ok()).toBeTruthy()
    const buildData = await buildResponse.json()
    expect(buildData.valid).toBe(true)
    expect(buildData.tag).toBe('e2e-test-image:latest')

    // Step 2: Verify image in API
    const imagesResponse = await page.request.get('http://localhost:8081/api/images')
    const imagesData = await imagesResponse.json()
    const builtImage = imagesData.find((img: any) => img.name === 'e2e-test-image')
    expect(builtImage).toBeTruthy()
    expect(builtImage.tag).toBe('latest')

    // Step 3: Navigate to Images page
    await page.goto('http://localhost:3001/images')
    const imagesPage = new ImagesPage(page)
    await imagesPage.navigate()

    // Step 4: Wait for images to load
    await page.waitForTimeout(1000)

    // Step 5: Verify image appears in dashboard
    await expect(page.getByText(/e2e-test-image/i)).toBeVisible()
    await expect(page.getByText(/latest/i).first()).toBeVisible()
  })

  test('@mvp should refresh images list after building', async ({ page }) => {
    // Build image
    const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json')
    const buildspecContent = JSON.parse(fs.readFileSync(buildspecPath, 'utf-8'))
    
    await page.request.post('http://localhost:8081/api/buildspec/build', {
      data: {
        spec: buildspecContent,
        tag: 'refresh-test:latest'
      }
    })

    // Navigate to images page
    await page.goto('http://localhost:3001/images')
    await page.waitForTimeout(500)

    // Click refresh
    await page.getByRole('button', { name: /Refresh/i }).click()
    await page.waitForTimeout(1000)

    // Should see the new image
    await expect(page.getByText(/refresh-test/i).or(page.getByText(/mern-mvp/i)).or(page.getByText(/my-app/i))).toBeVisible()
  })

  test('@mvp should show image details', async ({ page }) => {
    // Build image first
    const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json')
    const buildspecContent = JSON.parse(fs.readFileSync(buildspecPath, 'utf-8'))
    
    await page.request.post('http://localhost:8081/api/buildspec/build', {
      data: {
        spec: buildspecContent,
        tag: 'detail-test:latest'
      }
    })

    // Navigate to images
    await page.goto('http://localhost:3001/images')
    await page.waitForTimeout(1000)

    // Find image card
    const imageCard = page.locator('.card, [class*="Card"]').filter({ hasText: /detail-test|mern-mvp|my-app/i }).first()
    
    // Should show image details
    await expect(imageCard.getByText(/latest/i)).toBeVisible()
    await expect(imageCard.getByText(/MB|GB/i)).toBeVisible() // Size
  })
})

