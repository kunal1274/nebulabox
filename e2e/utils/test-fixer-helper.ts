import { Page } from '@playwright/test'

/**
 * Test Fixer Helper Utilities
 * Provides helper functions for fixing common test issues
 */

/**
 * Wait for page to be ready (handles pages with continuous polling)
 */
export async function waitForPageReady(page: Page, timeout: number = 30000): Promise<void> {
  // Wait for DOM content loaded
  await page.waitForLoadState('domcontentloaded', { timeout })
  
  // Wait for React to render - check for common React app indicators
  await Promise.race([
    page.waitForSelector('nav, [role="navigation"], aside, main, [class*="layout"]', { 
      state: 'visible', 
      timeout: 10000 
    }).catch(() => null),
    page.waitForTimeout(2000), // Fallback timeout
  ])
}

/**
 * Check if page loaded successfully
 */
export async function isPageLoaded(page: Page): Promise<boolean> {
  // Check for critical errors
  const criticalErrors = await page.locator('text=/404|500|Not Found|Server Error/i').count()
  if (criticalErrors > 0) {
    return false
  }
  
  // Check for main content
  const hasContent = await Promise.race([
    page.locator('nav, h1, h2, [role="main"], main, body').first().isVisible().then(() => true),
    page.waitForTimeout(1000).then(() => false),
  ])
  
  return hasContent
}

/**
 * Safe navigation with retry
 */
export async function safeNavigate(
  page: Page,
  url: string,
  retries: number = 3
): Promise<void> {
  for (let i = 0; i < retries; i++) {
    try {
      await page.goto(url, {
        waitUntil: 'domcontentloaded',
        timeout: 30000,
      })
      await waitForPageReady(page, 10000)
      
      // Verify page loaded
      const loaded = await isPageLoaded(page)
      if (loaded) {
        return
      }
      
      if (i < retries - 1) {
        console.log(`Navigation attempt ${i + 1} failed, retrying...`)
        await page.waitForTimeout(1000)
      }
    } catch (error) {
      if (i === retries - 1) {
        throw error
      }
      await page.waitForTimeout(1000)
    }
  }
}

/**
 * Wait for element with multiple fallback strategies
 */
export async function waitForElementWithFallback(
  page: Page,
  selectors: string[],
  timeout: number = 10000
): Promise<boolean> {
  for (const selector of selectors) {
    try {
      await page.waitForSelector(selector, { state: 'visible', timeout })
      return true
    } catch {
      continue
    }
  }
  return false
}

/**
 * Check for network errors
 */
export async function checkNetworkErrors(page: Page): Promise<string[]> {
  const errors: string[] = []
  
  page.on('requestfailed', (request) => {
    if (request.url().includes('localhost')) {
      errors.push(`Failed request: ${request.url()} - ${request.failure()?.errorText}`)
    }
  })
  
  return errors
}

