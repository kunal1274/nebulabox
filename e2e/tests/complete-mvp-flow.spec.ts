import { test, expect } from '@playwright/test'
import { BuildSpecPage, ImagesPage, ContainersPage } from '../utils/page-objects'
import * as fs from 'fs'
import * as path from 'path'

/**
 * Complete MVP Flow E2E Test
 * End-to-end test of the full MERN stack deployment workflow
 */
test.describe('Complete MVP Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Ensure test mode
    await page.request.put('http://localhost:8081/api/mode', {
      data: { mode: 'test' }
    })
  })

  test('@mvp complete MERN stack deployment flow', async ({ page }) => {
    // Step 1: Navigate to Build Spec page
    await page.goto('http://localhost:3001/buildspec')
    const buildSpecPage = new BuildSpecPage(page)

    // Step 2: Load buildspec.json
    const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json')
    const buildspecContent = fs.readFileSync(buildspecPath, 'utf-8')
    const specJson = JSON.parse(buildspecContent)

    await buildSpecPage.setSpecJson(buildspecContent)
    await page.waitForTimeout(500)

    // Step 3: Validate buildspec
    await buildSpecPage.clickValidate()
    await page.waitForTimeout(2000)
    
    // Step 4: Build image
    await buildSpecPage.clickBuild()
    await page.waitForTimeout(3000)

    // Step 5: Verify image was built (check API)
    const imagesResponse = await page.request.get('http://localhost:8081/api/images')
    const imagesData = await imagesResponse.json()
    const builtImage = imagesData.find((img: any) => img.name === 'mern-mvp' || img.name === specJson.name)
    expect(builtImage).toBeTruthy()

    // Step 6: Navigate to Images page and verify
    await page.goto('http://localhost:3001/images')
    await page.waitForTimeout(1000)
    
    const imageName = builtImage?.name || 'mern-mvp'
    await expect(page.getByText(new RegExp(imageName, 'i'))).toBeVisible({ timeout: 5000 })

    // Step 7: Navigate to Containers page
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(1000)

    // Step 8: Create container
    await page.getByRole('button', { name: /New Container|Create Container/i }).click()
    await page.waitForTimeout(500)

    // Fill container form
    const imageInput = page.getByLabel(/Image/i).or(page.locator('input[placeholder*="image" i]')).first()
    await imageInput.fill(`${imageName}:latest`)
    await page.waitForTimeout(200)

    const nameInput = page.getByLabel(/Name/i).or(page.locator('input[placeholder*="name" i]')).first()
    if (await nameInput.isVisible()) {
      await nameInput.fill('mern-mvp-test')
    }

    // Add ports if form has port inputs
    const portInputs = page.locator('input[type="text"]').filter({ hasText: /port/i })
    if (await portInputs.count() > 0) {
      // Try to add ports if form supports it
    }

    // Submit form
    const createButton = page.getByRole('button', { name: /Create|Submit/i }).first()
    await createButton.click()
    await page.waitForTimeout(2000)

    // Step 9: Verify container appears in list
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(1000)
    
    // Should see the container
    await expect(
      page.getByText(/mern-mvp-test/i).or(page.locator('.card, [class*="Card"]').filter({ hasText: /mern-mvp/i })).first()
    ).toBeVisible({ timeout: 5000 })

    // Step 10: Verify container is running
    const containersResponse = await page.request.get('http://localhost:8081/api/containers')
    const containersData = await containersResponse.json()
    const container = containersData.find((c: any) => 
      c.name === 'mern-mvp-test' || c.image.includes('mern-mvp')
    )
    expect(container).toBeTruthy()
    expect(['running', 'created']).toContain(container.status)

    // Step 11: Test container lifecycle - Stop
    const containerCard = page.locator('.card, [class*="Card"]').filter({ hasText: /mern-mvp/i }).first()
    const stopButton = containerCard.getByRole('button', { name: /Stop/i }).first()
    await stopButton.click()
    await page.waitForTimeout(1000)

    // Step 12: Verify stopped
    const stoppedResponse = await page.request.get('http://localhost:8081/api/containers')
    const stoppedData = await stoppedResponse.json()
    const stoppedContainer = stoppedData.find((c: any) => 
      c.name === 'mern-mvp-test' || c.image.includes('mern-mvp')
    )
    expect(stoppedContainer.status).toBe('stopped')

    // Step 13: Test container lifecycle - Start
    const startButton = containerCard.getByRole('button', { name: /Start/i }).first()
    await startButton.click()
    await page.waitForTimeout(1000)

    // Step 14: Verify started
    const startedResponse = await page.request.get('http://localhost:8081/api/containers')
    const startedData = await startedResponse.json()
    const startedContainer = startedData.find((c: any) => 
      c.name === 'mern-mvp-test' || c.image.includes('mern-mvp')
    )
    expect(startedContainer.status).toBe('running')
  })
})

