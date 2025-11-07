import { Page } from '@playwright/test'

/**
 * API Helper Utilities for E2E Tests
 * Provides helper functions for API interactions during tests
 */

export class ApiHelper {
  constructor(private page: Page) {}

  /**
   * Set operating mode
   */
  async setMode(mode: 'mock' | 'test' | 'live'): Promise<any> {
    const response = await this.page.request.put('http://localhost:8081/api/mode', {
      data: { mode }
    })
    expect(response.ok()).toBeTruthy()
    return response.json()
  }

  /**
   * Get current mode
   */
  async getMode(): Promise<any> {
    const response = await this.page.request.get('http://localhost:8081/api/mode')
    expect(response.ok()).toBeTruthy()
    return response.json()
  }

  /**
   * Build image from buildspec
   */
  async buildImage(spec: any, tag: string): Promise<any> {
    const response = await this.page.request.post('http://localhost:8081/api/buildspec/build', {
      data: { spec, tag }
    })
    expect(response.ok()).toBeTruthy()
    return response.json()
  }

  /**
   * List images
   */
  async listImages(): Promise<any[]> {
    const response = await this.page.request.get('http://localhost:8081/api/images')
    expect(response.ok()).toBeTruthy()
    return response.json()
  }

  /**
   * Find image by name
   */
  async findImage(name: string): Promise<any | null> {
    const images = await this.listImages()
    return images.find((img: any) => img.name === name) || null
  }

  /**
   * Create container
   */
  async createContainer(options: {
    image: string
    name: string
    ports?: string[]
    env?: string[]
    detach?: boolean
  }): Promise<any> {
    const response = await this.page.request.post('http://localhost:8081/api/containers/run', {
      data: options
    })
    expect(response.ok()).toBeTruthy()
    return response.json()
  }

  /**
   * List containers
   */
  async listContainers(): Promise<any[]> {
    const response = await this.page.request.get('http://localhost:8081/api/containers')
    expect(response.ok()).toBeTruthy()
    return response.json()
  }

  /**
   * Find container by name
   */
  async findContainer(name: string): Promise<any | null> {
    const containers = await this.listContainers()
    return containers.find((c: any) => c.name === name) || null
  }

  /**
   * Stop container
   */
  async stopContainer(containerId: string): Promise<any> {
    const response = await this.page.request.post(`http://localhost:8081/api/containers/${containerId}/stop`)
    expect(response.ok()).toBeTruthy()
    return response.json()
  }

  /**
   * Start container
   */
  async startContainer(containerId: string): Promise<any> {
    const response = await this.page.request.post(`http://localhost:8081/api/containers/${containerId}/start`)
    expect(response.ok()).toBeTruthy()
    return response.json()
  }

  /**
   * Get container logs
   */
  async getContainerLogs(containerId: string): Promise<any> {
    const response = await this.page.request.get(`http://localhost:8081/api/containers/${containerId}/logs`)
    expect(response.ok()).toBeTruthy()
    return response.json()
  }

  /**
   * Wait for container status
   */
  async waitForContainerStatus(containerName: string, status: string, timeout: number = 10000): Promise<void> {
    const startTime = Date.now()
    while (Date.now() - startTime < timeout) {
      const container = await this.findContainer(containerName)
      if (container && container.status === status) {
        return
      }
      await this.page.waitForTimeout(500)
    }
    throw new Error(`Container ${containerName} did not reach status ${status} within ${timeout}ms`)
  }

  /**
   * Wait for image to appear
   */
  async waitForImage(imageName: string, timeout: number = 10000): Promise<any> {
    const startTime = Date.now()
    while (Date.now() - startTime < timeout) {
      const image = await this.findImage(imageName)
      if (image) {
        return image
      }
      await this.page.waitForTimeout(500)
    }
    throw new Error(`Image ${imageName} did not appear within ${timeout}ms`)
  }

  /**
   * Clean up test data
   */
  async cleanup(): Promise<void> {
    // This could be extended to delete test containers/images
    // For now, just ensure we're in test mode
    await this.setMode('test')
  }
}

