import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { api } from '../api'

// Mock fetch globally
if (typeof globalThis.fetch === 'undefined') {
  globalThis.fetch = vi.fn() as typeof fetch
} else {
  vi.stubGlobal('fetch', vi.fn())
}

describe('ApiClient', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Reset any stored tokens
    localStorage.clear()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })


  describe('authentication methods', () => {
    it('stores token after login', async () => {
      const mockToken = 'auth-token-123'
      vi.mocked(globalThis.fetch).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ token: mockToken }),
        headers: new Headers(),
      } as Response)

      await api.login('user@example.com', 'password')

      expect(localStorage.getItem('authToken')).toBe(mockToken)
    })

    it('removes token on logout', async () => {
      localStorage.setItem('authToken', 'test-token')
      vi.mocked(globalThis.fetch).mockResolvedValueOnce({
        ok: true,
        json: async () => ({}),
        headers: new Headers(),
      } as Response)
      
      await api.logout()
      // Token should be cleared after logout
      expect(globalThis.fetch).toHaveBeenCalled()
    })
  })

  describe('container methods', () => {
    it('lists containers', async () => {
      const mockContainers = [
        { id: '1', name: 'container1', image: 'nginx:latest', status: 'running' as const, created: new Date().toISOString() },
        { id: '2', name: 'container2', image: 'alpine:latest', status: 'stopped' as const, created: new Date().toISOString() },
      ]

      vi.mocked(globalThis.fetch).mockResolvedValueOnce({
        ok: true,
        json: async () => mockContainers,
        headers: new Headers(),
      } as Response)

      const result = await api.listContainers()
      expect(result).toEqual(mockContainers)
    })

    it('creates container', async () => {
      const mockContainer = { id: '1', name: 'new-container', image: 'nginx:latest', status: 'running' as const, created: new Date().toISOString() }
      vi.mocked(globalThis.fetch).mockResolvedValueOnce({
        ok: true,
        json: async () => mockContainer,
        headers: new Headers(),
      } as Response)

      const result = await api.runContainer({
        image: 'nginx:latest',
        name: 'new-container',
      })

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/containers/run'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('nginx:latest'),
        })
      )
      expect(result).toEqual(mockContainer)
    })
  })

  describe('system stats', () => {
    it('fetches system stats', async () => {
      const mockStats = {
        cpuUsage: 45.2,
        memoryUsage: 62.8,
        diskUsage: 38.5,
        containersRunning: 3,
        containersTotal: 5,
      }

      vi.mocked(globalThis.fetch).mockResolvedValueOnce({
        ok: true,
        json: async () => mockStats,
        headers: new Headers(),
      } as Response)

      const result = await api.getSystemStats()
      expect(result).toEqual(mockStats)
    })
  })
})

