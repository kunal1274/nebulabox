/**
 * API Test-Aware Utilities
 * Provides test-aware tracking for API calls
 */

import { TestMode } from './test-aware'

/**
 * Test-Aware API Client Wrapper
 * Tracks all API calls for testing purposes
 */

interface ApiCall {
  timestamp: number
  method: string
  endpoint: string
  request: any
  response?: any
  error?: string
  duration?: number
}

class TestAwareApiClient {
  private calls: ApiCall[] = []
  private maxCalls = 1000 // Limit stored calls

  /**
   * Record API call
   */
  recordCall(
    method: string,
    endpoint: string,
    request: any,
    response?: any,
    error?: string,
    duration?: number
  ): void {
    // Always record calls, but only if tracking is enabled
    // This allows tests to force enable tracking if needed
    const shouldTrack = TestMode.isTestMode() || 
      (typeof window !== 'undefined' && (window as any).__FORCE_API_TRACKING === true)
    
    if (!shouldTrack) {
      return
    }

    const call: ApiCall = {
      timestamp: Date.now(),
      method,
      endpoint,
      request,
      response,
      error,
      duration,
    }

    this.calls.push(call)

    // Maintain max size
    if (this.calls.length > this.maxCalls) {
      this.calls.shift()
    }

    // Emit event for test recording
    if (typeof window !== 'undefined' && (window as any).__RECORDER__) {
      ;(window as any).__RECORDER__.annotate(`API Call: ${method} ${endpoint}`, {
        type: 'api',
        call,
      })
    }
  }

  /**
   * Get all recorded calls
   */
  getCalls(): ApiCall[] {
    return [...this.calls]
  }

  /**
   * Get calls for specific endpoint
   */
  getCallsForEndpoint(endpoint: string): ApiCall[] {
    return this.calls.filter((call) => call.endpoint === endpoint)
  }

  /**
   * Get last call
   */
  getLastCall(): ApiCall | null {
    return this.calls.length > 0 ? this.calls[this.calls.length - 1] : null
  }

  /**
   * Clear calls
   */
  clear(): void {
    this.calls = []
  }

  /**
   * Export calls as JSON
   */
  export(): string {
    return JSON.stringify(this.calls, null, 2)
  }
}

export const apiTracker = new TestAwareApiClient()

/**
 * Wrap API calls with tracking
 */
export function trackApiCall<T>(
  method: string,
  endpoint: string,
  request: any,
  apiCall: () => Promise<T>
): Promise<T> {
  const startTime = Date.now()

  return apiCall()
    .then((response) => {
      const duration = Date.now() - startTime
      apiTracker.recordCall(method, endpoint, request, response, undefined, duration)
      return response
    })
    .catch((error) => {
      const duration = Date.now() - startTime
      apiTracker.recordCall(method, endpoint, request, undefined, String(error), duration)
      throw error
    })
}

/**
 * Get API call summary for test artifacts
 */
export function getApiCallSummary(): {
  total: number
  successful: number
  failed: number
  calls: ApiCall[]
} {
  const calls = apiTracker.getCalls()
  return {
    total: calls.length,
    successful: calls.filter((c) => !c.error).length,
    failed: calls.filter((c) => c.error).length,
    calls,
  }
}

/**
 * Initialize API tracking in test mode
 */
export function initApiTracking(): void {
  if (typeof window === 'undefined') return

  // Always expose API tracker (but only tracks in test mode)
  // This allows tests to access it even if test mode isn't detected
  ;(window as any).__API_TRACKER__ = {
    getCalls: () => apiTracker.getCalls(),
    getCallsForEndpoint: (endpoint: string) => apiTracker.getCallsForEndpoint(endpoint),
    getLastCall: () => apiTracker.getLastCall(),
    clear: () => apiTracker.clear(),
    export: () => apiTracker.export(),
    getSummary: () => getApiCallSummary(),
    isEnabled: () => TestMode.isTestMode(),
  }

  if (TestMode.isTestMode()) {
    console.log('[API Tracker] Initialized and tracking enabled')
  } else {
    console.log('[API Tracker] Initialized (tracking disabled - not in test mode)')
  }
}

// Auto-initialize
if (typeof window !== 'undefined') {
  initApiTracking()
}

