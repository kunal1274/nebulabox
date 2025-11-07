import { test, expect } from '@playwright/test'
import { ApiHelper } from '../utils/api-helpers'
import * as fs from 'fs'
import * as path from 'path'

/**
 * Data Persistence Tests
 * Tests that data persists correctly in test mode
 */
test.describe('Data Persistence in Test Mode', () => {
  let apiHelper: ApiHelper

  test.beforeEach(async ({ page }) => {
    apiHelper = new ApiHelper(page)
    await apiHelper.setMode('test')
  })

  test('@mvp should persist images across page refreshes', async ({ page }) => {
    const testImageName = 'persist-test-image'

    // Build image
    const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json')
    const buildspecContent = JSON.parse(fs.readFileSync(buildspecPath, 'utf-8'))
    buildspecContent.tag = `${testImageName}:latest`
    
    await apiHelper.buildImage(buildspecContent, `${testImageName}:latest`)
    const image = await apiHelper.waitForImage(testImageName, 10000)
    expect(image).toBeTruthy()

    // Refresh images page
    await page.goto('http://localhost:3001/images')
    await page.waitForTimeout(1000)
    await page.reload()
    await page.waitForTimeout(1000)

    // Image should still be there
    const images = await apiHelper.listImages()
    const persistedImage = images.find((img: any) => img.name === testImageName)
    expect(persistedImage).toBeTruthy()

    // Verify in dashboard
    await expect(page.getByText(new RegExp(testImageName, 'i'))).toBeVisible({ timeout: 5000 })
  })

  test('@mvp should persist containers across page refreshes', async ({ page }) => {
    const testContainerName = 'persist-test-container'

    // Create container
    const container = await apiHelper.createContainer({
      image: 'mern-mvp:latest',
      name: testContainerName,
      detach: true
    })
    expect(container).toBeTruthy()

    // Refresh containers page
    await page.goto('http://localhost:3001/containers')
    await page.waitForTimeout(1000)
    await page.reload()
    await page.waitForTimeout(1000)

    // Container should still be there
    const containers = await apiHelper.listContainers()
    const persistedContainer = containers.find((c: any) => c.name === testContainerName)
    expect(persistedContainer).toBeTruthy()

    // Verify in dashboard
    await expect(page.getByText(new RegExp(testContainerName, 'i'))).toBeVisible({ timeout: 5000 })
  })

  test('@mvp should clear data when switching to mock mode', async ({ page }) => {
    // Build image in test mode
    const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json')
    const buildspecContent = JSON.parse(fs.readFileSync(buildspecPath, 'utf-8'))
    buildspecContent.tag = 'clear-test:latest'
    
    await apiHelper.buildImage(buildspecContent, 'clear-test:latest')
    const image = await apiHelper.waitForImage('clear-test', 10000)
    expect(image).toBeTruthy()

    // Create container
    const container = await apiHelper.createContainer({
      image: 'clear-test:latest',
      name: 'clear-test-container',
      detach: true
    })
    expect(container).toBeTruthy()

    // Switch to mock mode
    await apiHelper.setMode('mock')

    // Images should be cleared
    const images = await apiHelper.listImages()
    const clearedImage = images.find((img: any) => img.name === 'clear-test')
    expect(clearedImage).toBeFalsy()

    // Containers should be cleared
    const containers = await apiHelper.listContainers()
    const clearedContainer = containers.find((c: any) => c.name === 'clear-test-container')
    expect(clearedContainer).toBeFalsy()
  })

  test('@mvp should preserve data when switching between test and live mode', async ({ page }) => {
    // Build image in test mode
    const buildspecPath = path.join(__dirname, '../../mvp/app/single-container/buildspec.json')
    const buildspecContent = JSON.parse(fs.readFileSync(buildspecPath, 'utf-8'))
    buildspecContent.tag = 'preserve-test:latest'
    
    await apiHelper.buildImage(buildspecContent, 'preserve-test:latest')
    const image = await apiHelper.waitForImage('preserve-test', 10000)
    expect(image).toBeTruthy()

    // Switch to live mode
    await apiHelper.setMode('live')

    // Image should still be there
    const images = await apiHelper.listImages()
    const preservedImage = images.find((img: any) => img.name === 'preserve-test')
    expect(preservedImage).toBeTruthy()

    // Switch back to test mode
    await apiHelper.setMode('test')

    // Image should still be there
    const images2 = await apiHelper.listImages()
    const preservedImage2 = images2.find((img: any) => img.name === 'preserve-test')
    expect(preservedImage2).toBeTruthy()
  })
})

