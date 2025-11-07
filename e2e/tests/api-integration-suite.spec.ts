import { test, expect } from '@playwright/test'
import { ApiHelper } from '../utils/api-helpers'

/**
 * API Integration Test Suite
 * Tests all API endpoints and their integration with UI
 */
test.describe('API Integration Suite', () => {
  let apiHelper: ApiHelper

  test.beforeEach(async ({ page }) => {
    apiHelper = new ApiHelper(page)
    await apiHelper.setMode('test')
  })

  test('@mvp API health check', async ({ page }) => {
    const response = await page.request.get('http://localhost:8081/api/health')
    expect(response.ok()).toBeTruthy()
    const data = await response.json()
    expect(data.status).toBe('healthy')
  })

  test('@mvp API mode endpoints', async ({ page }) => {
    // Get mode
    const modeData = await apiHelper.getMode()
    expect(['mock', 'test', 'live']).toContain(modeData.mode)
    
    // Set mode
    await apiHelper.setMode('test')
    const newModeData = await apiHelper.getMode()
    expect(newModeData.mode).toBe('test')
  })

  test('@mvp API images endpoints', async ({ page }) => {
    // List images
    const images = await apiHelper.listImages()
    expect(Array.isArray(images)).toBe(true)
    
    // Build image
    const buildspecPath = require.resolve('../../mvp/app/single-container/buildspec.json')
    const buildspecContent = require(buildspecPath)
    buildspecContent.tag = 'api-test:latest'
    
    const buildResult = await apiHelper.buildImage(buildspecContent, 'api-test:latest')
    expect(buildResult.valid).toBe(true)
    
    // Find image
    const image = await apiHelper.waitForImage('api-test', 10000)
    expect(image).toBeTruthy()
  })

  test('@mvp API containers endpoints', async ({ page }) => {
    // List containers
    const containers = await apiHelper.listContainers()
    expect(Array.isArray(containers)).toBe(true)
    
    // Create container
    const container = await apiHelper.createContainer({
      image: 'mern-mvp:latest',
      name: 'api-container-test',
      detach: true
    })
    expect(container.name).toBe('api-container-test')
    
    // Find container
    const foundContainer = await apiHelper.findContainer('api-container-test')
    expect(foundContainer).toBeTruthy()
    
    // Stop container
    await apiHelper.stopContainer(container.id)
    await apiHelper.waitForContainerStatus('api-container-test', 'stopped', 10000)
    
    // Start container
    await apiHelper.startContainer(container.id)
    await apiHelper.waitForContainerStatus('api-container-test', 'running', 10000)
    
    // Get logs
    const logs = await apiHelper.getContainerLogs(container.id)
    expect(logs).toBeTruthy()
  })

  test('@mvp API system stats', async ({ page }) => {
    const response = await page.request.get('http://localhost:8081/api/system/stats')
    expect(response.ok()).toBeTruthy()
    const stats = await response.json()
    expect(stats).toHaveProperty('cpuUsage')
    expect(stats).toHaveProperty('memoryUsage')
    expect(stats).toHaveProperty('containersRunning')
  })
})

