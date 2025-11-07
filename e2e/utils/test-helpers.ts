import { Page, Locator, expect } from '@playwright/test';

/**
 * Visual action highlighting utilities
 * Adds sparkles-like effects to actions for better visibility in recordings
 */
export class VisualHelper {
  constructor(private page: Page) {}

  /**
   * Highlight an element with a sparkle effect before clicking
   */
  async highlightAndClick(locator: Locator, options?: { timeout?: number }) {
    const selector = await this.getSelectorFromLocator(locator);
    if (selector) {
      await this.highlightElement(selector, options);
    }
    await locator.click(options);
    if (selector) {
      await this.unhighlightElement(selector);
    }
  }

  /**
   * Highlight an element with a sparkle effect before typing
   */
  async highlightAndType(locator: Locator, text: string, options?: { timeout?: number }) {
    const selector = await this.getSelectorFromLocator(locator);
    if (selector) {
      await this.highlightElement(selector, options);
    }
    await locator.fill(text, options);
    if (selector) {
      await this.unhighlightElement(selector);
    }
  }

  /**
   * Get a CSS selector string from a Locator (if possible)
   */
  private async getSelectorFromLocator(locator: Locator): Promise<string | null> {
    try {
      // Try to get a unique selector
      return await locator.evaluate((el) => {
        // Generate a simple selector if possible
        if (el.id) return `#${el.id}`;
        if (el.className) {
          const classes = el.className.split(' ').filter(c => c).join('.');
          if (classes) return `.${classes}`;
        }
        return el.tagName.toLowerCase();
      }).catch(() => null);
    } catch {
      return null;
    }
  }

  /**
   * Create sparkle effect on element
   */
  private async highlightElement(selector: string, options?: { timeout?: number }) {
    await this.page.evaluate(
      ({ sel }) => {
        const element = document.querySelector(sel);
        if (element) {
          const rect = element.getBoundingClientRect();
          
          // Create sparkle overlay
          const overlay = document.createElement('div');
          overlay.id = 'pw-sparkle-overlay';
          overlay.style.cssText = `
            position: fixed;
            top: ${rect.top}px;
            left: ${rect.left}px;
            width: ${rect.width}px;
            height: ${rect.height}px;
            border: 3px solid #00ff00;
            border-radius: 8px;
            box-shadow: 0 0 20px #00ff00, 0 0 40px #00ff00, inset 0 0 20px rgba(0, 255, 0, 0.3);
            pointer-events: none;
            z-index: 99999;
            animation: sparkle 0.5s ease-in-out;
          `;
          
          // Add sparkle animation
          const style = document.createElement('style');
          style.textContent = `
            @keyframes sparkle {
              0% { transform: scale(0.8); opacity: 0; }
              50% { transform: scale(1.1); opacity: 1; }
              100% { transform: scale(1); opacity: 0.8; }
            }
          `;
          if (!document.getElementById('pw-sparkle-style')) {
            style.id = 'pw-sparkle-style';
            document.head.appendChild(style);
          }
          
          document.body.appendChild(overlay);
        }
      },
      { sel: selector }
    );
    
    // Wait for animation
    await this.page.waitForTimeout(300);
  }

  /**
   * Remove highlight
   */
  private async unhighlightElement(selector: string) {
    await this.page.evaluate(() => {
      const overlay = document.getElementById('pw-sparkle-overlay');
      if (overlay) {
        overlay.style.opacity = '0';
        overlay.style.transition = 'opacity 0.3s';
        setTimeout(() => overlay.remove(), 300);
      }
    });
    await this.page.waitForTimeout(100);
  }

  /**
   * Wait for element and scroll into view with highlight
   */
  async scrollToAndHighlight(selector: string) {
    await this.page.waitForSelector(selector, { state: 'visible' });
    await this.page.locator(selector).scrollIntoViewIfNeeded();
    await this.highlightElement(selector);
    await this.page.waitForTimeout(200);
    await this.unhighlightElement(selector);
  }
}

/**
 * API health check utilities
 */
export class HealthChecker {
  constructor(private baseURL: string = 'http://localhost:8081') {}

  /**
   * Check if API server is healthy
   */
  async checkAPIHealth(maxRetries = 10, delay = 1000): Promise<boolean> {
    for (let i = 0; i < maxRetries; i++) {
      try {
        const response = await fetch(`${this.baseURL}/api/health`);
        if (response.ok) {
          const data = await response.json();
          return data.status === 'healthy';
        }
      } catch (error) {
        // Server not ready yet
      }
      await new Promise(resolve => setTimeout(resolve, delay));
    }
    return false;
  }

  /**
   * Check if frontend is accessible
   */
  async checkFrontendHealth(maxRetries = 10, delay = 1000): Promise<boolean> {
    for (let i = 0; i < maxRetries; i++) {
      try {
        const response = await fetch('http://localhost:3001');
        return response.ok;
      } catch (error) {
        // Server not ready yet
      }
      await new Promise(resolve => setTimeout(resolve, delay));
    }
    return false;
  }

  /**
   * Wait for all services to be ready
   */
  async waitForServices(): Promise<void> {
    const [apiReady, frontendReady] = await Promise.all([
      this.checkAPIHealth(),
      this.checkFrontendHealth(),
    ]);

    if (!apiReady) {
      throw new Error('API server is not healthy');
    }
    if (!frontendReady) {
      throw new Error('Frontend server is not accessible');
    }
  }
}

/**
 * Screenshot utilities with timestamps
 */
export class ScreenshotHelper {
  constructor(private page: Page) {}

  /**
   * Take screenshot with descriptive name
   */
  async capture(name: string) {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    await this.page.screenshot({
      path: `test-results/screenshots/${name}-${timestamp}.png`,
      fullPage: true,
    });
  }

  /**
   * Take screenshot of specific element
   */
  async captureElement(selector: string, name: string) {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    await this.page.locator(selector).screenshot({
      path: `test-results/screenshots/${name}-${timestamp}.png`,
    });
  }
}

