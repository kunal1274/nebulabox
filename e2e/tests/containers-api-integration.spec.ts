import { test, expect } from '@playwright/test'
import { ContainersPage } from '../utils/page-objects'

/**
 * Containers API Integration E2E Tests
 * Tests that created containers appear in dashboard after API creation
 */
test.describe('Containers API Integration', () => {
  test.beforeEach(async ({ page }) => {
    // Ensure we're in test mode
    await page.goto('http://localhost:3001')
    const modeResponse = await page.request.put('http://localhost:8081/api/mode', {
      data: { mode: 'test' }
    })
    expect(modeResponse.ok()).toBeTruthy()
  })

  test('@mvp should create container and see it in dashboard', async ({ page }) => {
    // Step 1: Create container via API
    const createResponse = await page.request.post('http://localhost:8081/api/containers/run', {
      data: {
        image: 'mern-mvp:latest',
        name: 'e2e-test-container',
        ports: ['3000:3000'],
        env: ['TEST_VAR=test_value'],
        detach: true
      }
    })

    expect(createResponse.ok()).toBeTruthy()
    const containerData = await createResponse.json()
    expect(containerData.name).toBe('e2e-test-container')
    const containerId = containerData.id

    // Step 2: Verify container in API
    const containersResponse = await page.request.get('http://localhost:8081/api/containers')
    const containersData = await containersResponse.json()
    const createdContainer = containersData.find((c: any) => c.name === 'e2e-test-container')
    expect(createdContainer).toBeTruthy()
    expect(createdContainer.image).toBe('mern-mvp:latest')

    // Step 3: Navigate to Containers page
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(1000)

    // Step 4: Verify container appears in dashboard
    await expect(page.getByText(/e2e-test-container/i)).toBeVisible()
    await expect(page.getByText(/mern-mvp:latest/i)).toBeVisible()
  })

  test('@mvp should stop container and update status', async ({ page }) => {
    // Create container first
    const createResponse = await page.request.post('http://localhost:8081/api/containers/run', {
      data: {
        image: 'mern-mvp:latest',
        name: 'e2e-stop-test',
        detach: true
      }
    })
    const containerData = await createResponse.json()
    const containerId = containerData.id

    // Navigate to containers page
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(1000)

    // Find container card
    const containerCard = page.locator('.card, [class*="Card"]').filter({ hasText: /e2e-stop-test/i }).first()
    await expect(containerCard).toBeVisible()

    // Find and click stop button
    const stopButton = containerCard.getByRole('button', { name: /Stop/i }).first()
    await stopButton.click()
    await page.waitForTimeout(1000)

    // Verify status changed via API
    const containersResponse = await page.request.get('http://localhost:8081/api/containers')
    const containersData = await containersResponse.json()
    const container = containersData.find((c: any) => c.name === 'e2e-stop-test')
    expect(container).toBeTruthy()
    expect(container.status).toBe('stopped')
  })

  test('@mvp should start container and update status', async ({ page }) => {
    // Create and stop container first
    const createResponse = await page.request.post('http://localhost:8081/api/containers/run', {
      data: {
        image: 'mern-mvp:latest',
        name: 'e2e-start-test',
        detach: true
      }
    })
    const containerData = await createResponse.json()
    const containerId = containerData.id

    // Stop it
    await page.request.post(`http://localhost:8081/api/containers/${containerId}/stop`)

    // Navigate to containers page
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(1000)

    // Find container and start it
    const containerCard = page.locator('.card, [class*="Card"]').filter({ hasText: /e2e-start-test/i }).first()
    const startButton = containerCard.getByRole('button', { name: /Start/i }).first()
    await startButton.click()
    await page.waitForTimeout(1000)

    // Verify status changed
    const containersResponse = await page.request.get('http://localhost:8081/api/containers')
    const containersData = await containersResponse.json()
    const container = containersData.find((c: any) => c.name === 'e2e-start-test')
    expect(container.status).toBe('running')
  })

  test('@mvp should display container list correctly', async ({ page }) => {
    // Create a test container
    await page.request.post('http://localhost:8081/api/containers/run', {
      data: {
        image: 'mern-mvp:latest',
        name: 'e2e-list-test',
        detach: true
      }
    })

    // Navigate to containers page
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(1000)

    // Should see containers list
    const containersList = page.locator('.card, [class*="Card"], table, [data-testid="container-list"]').first()
    await expect(containersList).toBeVisible()

    // Should see at least one container
    const containerCards = page.locator('.card, [class*="Card"], table tbody tr').filter({ hasNotText: 'Loading' })
    const count = await containerCards.count()
    expect(count).toBeGreaterThan(0)
  })
})

