/**
 * Test-Aware Development Utilities
 * Provides helpers for test-aware development and data-test-id management
 */

/**
 * Test ID Generator
 * Generates consistent test IDs following the pattern: feature.action.element
 */
export class TestIDGenerator {
  /**
   * Generate test ID
   * @param feature Feature name (e.g., 'containers', 'images', 'buildspec')
   * @param action Action name (e.g., 'create', 'delete', 'refresh')
   * @param element Element type (e.g., 'button', 'input', 'card')
   */
  static generate(feature: string, action: string, element: string = 'button'): string {
    return `${feature}.${action}.${element}`
  }

  /**
   * Generate test ID for nested components
   * @param parts Array of test ID parts
   */
  static fromParts(parts: string[]): string {
    return parts.join('.')
  }
}

/**
 * Test Mode Detection
 */
export class TestMode {
  /**
   * Check if running in test mode
   */
  static isTestMode(): boolean {
    if (typeof window === 'undefined') return false
    return (
      window.location.search.includes('test=true') ||
      (typeof process !== 'undefined' && process.env?.NODE_ENV === 'test') ||
      localStorage.getItem('testMode') === 'true'
    )
  }

  /**
   * Disable animations in test mode
   */
  static disableAnimations(): void {
    if (this.isTestMode()) {
      document.documentElement.classList.add('no-animations')
    }
  }

  /**
   * Get test session ID
   */
  static getSessionId(): string {
    if (typeof window === 'undefined') return ''
    
    let sessionId = sessionStorage.getItem('test-session-id')
    if (!sessionId) {
      sessionId = `session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
      sessionStorage.setItem('test-session-id', sessionId)
    }
    return sessionId
  }

  /**
   * Get test run ID
   */
  static getRunId(): string {
    if (typeof window === 'undefined') return ''
    
    let runId = sessionStorage.getItem('test-run-id')
    if (!runId) {
      runId = `run-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
      sessionStorage.setItem('test-run-id', runId)
    }
    return runId
  }
}

/**
 * Test Recorder API
 * Provides window.__RECORDER__ API for manual test annotation
 */
export class TestRecorder {
  private static instance: TestRecorder | null = null
  private annotations: Array<{
    step: number
    action: string
    timestamp: number
    selector?: string
    metadata?: any
  }> = []

  static getInstance(): TestRecorder {
    if (!TestRecorder.instance) {
      TestRecorder.instance = new TestRecorder()
    }
    return TestRecorder.instance
  }

  /**
   * Annotate a test step
   */
  annotate(action: string, metadata?: any): void {
    const step = this.annotations.length + 1
    this.annotations.push({
      step,
      action,
      timestamp: Date.now(),
      metadata,
    })
    console.log(`[Test Recorder] Step ${step}: ${action}`, metadata)
  }

  /**
   * Get all annotations
   */
  getAnnotations() {
    return [...this.annotations]
  }

  /**
   * Clear annotations
   */
  clear(): void {
    this.annotations = []
  }
}

/**
 * Initialize test-aware utilities
 */
export function initTestAware(): void {
  if (typeof window === 'undefined') return

  // Disable animations in test mode
  TestMode.disableAnimations()

  // Expose recorder API
  const recorder = TestRecorder.getInstance()
  ;(window as any).__RECORDER__ = {
    annotate: (action: string, metadata?: any) => recorder.annotate(action, metadata),
    getAnnotations: () => recorder.getAnnotations(),
    clear: () => recorder.clear(),
    getSessionId: () => TestMode.getSessionId(),
    getRunId: () => TestMode.getRunId(),
  }

  console.log('[Test-Aware] Initialized', {
    sessionId: TestMode.getSessionId(),
    runId: TestMode.getRunId(),
    testMode: TestMode.isTestMode(),
  })
}

// Auto-initialize
if (typeof window !== 'undefined') {
  initTestAware()
}

