import { Page, Locator } from '@playwright/test'

/**
 * Page Object Models for E2E Tests
 * Centralized selectors and actions for UI components
 */

export class BasePage {
  constructor(protected page: Page) {}

  async navigate() {
    // Override in subclasses
  }

  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle')
    await this.page.waitForTimeout(500)
  }
}

export class BuildSpecPage extends BasePage {
  async navigate() {
    await this.page.goto('http://localhost:3001/buildspec')
    await this.waitForPageLoad()
  }

  async setSpecJson(json: string) {
    const editor = this.page.locator('textarea').or(this.page.locator('input[type="text"]').filter({ hasText: /json/i })).first()
    await editor.fill(json)
  }

  async clickValidate() {
    const validateButton = this.page.getByRole('button', { name: /Validate/i }).first()
    await validateButton.click()
  }

  async clickConvert() {
    const convertButton = this.page.getByRole('button', { name: /Convert/i }).first()
    await convertButton.click()
  }

  async clickBuild() {
    const buildButton = this.page.getByRole('button', { name: /Build/i }).first()
    await buildButton.click()
  }

  async clickLoadExample() {
    const loadButton = this.page.getByRole('button', { name: /Load Example/i }).first()
    await loadButton.click()
  }

  async waitForValidationComplete() {
    // Switch to Validation tab if exists
    const validationTab = this.page.getByRole('tab', { name: /Validation/i })
    if (await validationTab.isVisible()) {
      await validationTab.click()
    }
    await this.page.waitForTimeout(2000)
  }

  async waitForBuildComplete() {
    // Switch to Build Logs tab if exists
    const logsTab = this.page.getByRole('tab', { name: /Build Logs|Logs/i })
    if (await logsTab.isVisible()) {
      await logsTab.click()
    }
    await this.page.waitForTimeout(3000)
  }

  async getValidationResult(): Promise<string | null> {
    // Try to get validation result
    const result = this.page.locator('text=/valid|invalid|error/i').first()
    if (await result.isVisible()) {
      return await result.textContent()
    }
    return null
  }

  async getDockerfilePreview(): Promise<string | null> {
    // Switch to Dockerfile tab
    const dockerfileTab = this.page.getByRole('tab', { name: /Dockerfile/i })
    if (await dockerfileTab.isVisible()) {
      await dockerfileTab.click()
      await this.page.waitForTimeout(500)
    }
    
    const dockerfile = this.page.locator('textarea, pre, code').filter({ hasText: /FROM/i }).first()
    if (await dockerfile.isVisible()) {
      return await dockerfile.textContent()
    }
    return null
  }

  async getBuildLogs(): Promise<string[]> {
    // Switch to Build Logs tab
    const logsTab = this.page.getByRole('tab', { name: /Build Logs|Logs/i })
    if (await logsTab.isVisible()) {
      await logsTab.click()
      await this.page.waitForTimeout(500)
    }

    const logsContainer = this.page.locator('pre, code, [class*="log"]').filter({ hasText: /Building|built|success/i })
    if (await logsContainer.isVisible()) {
      const text = await logsContainer.textContent()
      return text ? text.split('\n') : []
    }
    return []
  }
}

export class ImagesPage extends BasePage {
  async navigate() {
    await this.page.goto('http://localhost:3001/images')
    await this.waitForPageLoad()
  }

  async getImageCount(): Promise<number> {
    const images = this.page.locator('.card, [class*="Card"], table tbody tr').filter({ hasNotText: 'Loading' })
    return await images.count()
  }

  async findImage(name: string): Promise<Locator | null> {
    const imageCard = this.page.locator('.card, [class*="Card"]').filter({ hasText: new RegExp(name, 'i') }).first()
    if (await imageCard.isVisible()) {
      return imageCard
    }
    return null
  }

  async clickRefresh() {
    const refreshButton = this.page.getByRole('button', { name: /Refresh/i }).first()
    await refreshButton.click()
    await this.page.waitForTimeout(1000)
  }

  async clickPullImage() {
    const pullButton = this.page.getByRole('button', { name: /Pull/i }).first()
    await pullButton.click()
  }

  async setPullImageName(imageName: string) {
    const input = this.page.getByPlaceholder(/nginx:latest|image/i).first()
    await input.fill(imageName)
  }
}

export class ContainersPage extends BasePage {
  async navigate() {
    await this.page.goto('http://localhost:3001/containers')
    await this.waitForPageLoad()
  }

  async clickCreate() {
    const createButton = this.page.getByRole('button', { name: /New Container|Create Container/i }).first()
    await createButton.click()
  }

  async getContainerCount(): Promise<number> {
    const containers = this.page.locator('.card, [class*="Card"], table tbody tr').filter({ hasNotText: /Loading/i })
    return await containers.count()
  }

  async findContainer(name: string): Promise<Locator | null> {
    const containerCard = this.page.locator('.card, [class*="Card"]').filter({ hasText: new RegExp(name, 'i') }).first()
    if (await containerCard.isVisible()) {
      return containerCard
    }
    return null
  }

  async stopContainer(name: string) {
    const container = await this.findContainer(name)
    if (container) {
      const stopButton = container.getByRole('button', { name: /Stop/i }).first()
      await stopButton.click()
      await this.page.waitForTimeout(1000)
    }
  }

  async startContainer(name: string) {
    const container = await this.findContainer(name)
    if (container) {
      const startButton = container.getByRole('button', { name: /Start/i }).first()
      await startButton.click()
      await this.page.waitForTimeout(1000)
    }
  }

  async clickRefresh() {
    const refreshButton = this.page.getByRole('button', { name: /Refresh|Show All/i }).first()
    await refreshButton.click()
    await this.page.waitForTimeout(1000)
  }
}
