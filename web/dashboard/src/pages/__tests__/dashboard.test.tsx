import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { Dashboard } from '../Dashboard'
import * as api from '@/lib/api'

// Mock the API client
vi.mock('@/lib/api', () => ({
  api: {
    getSystemStats: vi.fn(),
  },
}))

describe('Dashboard', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('renders dashboard with loading state initially', () => {
    vi.mocked(api.api.getSystemStats).mockImplementation(
      () => new Promise(() => {}) // Never resolves to simulate loading
    )
    
    render(<Dashboard />)
    // The component should render even while loading
    expect(screen.getByText(/running containers/i)).toBeInTheDocument()
  })

  it('displays system stats after loading', async () => {
    const mockStats = {
      cpuUsage: 45.2,
      memoryUsage: 62.8,
      diskUsage: 38.5,
      containersRunning: 3,
      containersTotal: 5,
    }

    vi.mocked(api.api.getSystemStats).mockResolvedValue(mockStats)

    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText(/3/i)).toBeInTheDocument()
      expect(screen.getByText(/5/i)).toBeInTheDocument()
    })
  })

  it('handles API errors gracefully with mock data', async () => {
    vi.mocked(api.api.getSystemStats).mockRejectedValue(new Error('API Error'))

    render(<Dashboard />)

    await waitFor(() => {
      // Should fallback to mock data on error
      expect(screen.getByText(/running containers/i)).toBeInTheDocument()
    })
  })

  it('displays CPU usage', async () => {
    const mockStats = {
      cpuUsage: 75.5,
      memoryUsage: 50.0,
      diskUsage: 30.0,
      containersRunning: 2,
      containersTotal: 4,
    }

    vi.mocked(api.api.getSystemStats).mockResolvedValue(mockStats)

    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText(/75.5%/i)).toBeInTheDocument()
    })
  })

  it('displays memory usage', async () => {
    const mockStats = {
      cpuUsage: 30.0,
      memoryUsage: 80.5,
      diskUsage: 40.0,
      containersRunning: 1,
      containersTotal: 3,
    }

    vi.mocked(api.api.getSystemStats).mockResolvedValue(mockStats)

    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText(/80.5%/i)).toBeInTheDocument()
    })
  })
})

