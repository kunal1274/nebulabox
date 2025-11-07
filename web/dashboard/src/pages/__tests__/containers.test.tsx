import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { Containers } from '../Containers'
import * as api from '@/lib/api'

// Mock the API client
vi.mock('@/lib/api', () => ({
  api: {
    listContainers: vi.fn(),
    stopContainer: vi.fn(),
    deleteContainer: vi.fn(),
  },
}))

const renderContainers = () => {
  return render(
    <BrowserRouter>
      <Containers />
    </BrowserRouter>
  )
}

describe('Containers Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders containers page', () => {
    vi.mocked(api.api.listContainers).mockResolvedValue([])
    renderContainers()
    expect(screen.getByText(/containers/i)).toBeInTheDocument()
  })

  it('displays list of containers', async () => {
    const mockContainers = [
      {
        id: '1',
        name: 'container1',
        image: 'nginx:latest',
        status: 'running' as const,
        created: new Date().toISOString(),
      },
      {
        id: '2',
        name: 'container2',
        image: 'alpine:latest',
        status: 'stopped' as const,
        created: new Date().toISOString(),
      },
    ]

    vi.mocked(api.api.listContainers).mockResolvedValue(mockContainers)
    renderContainers()

    await waitFor(() => {
      expect(screen.getByText('container1')).toBeInTheDocument()
      expect(screen.getByText('container2')).toBeInTheDocument()
    })
  })

  it('displays empty state when no containers', async () => {
    vi.mocked(api.api.listContainers).mockResolvedValue([])
    renderContainers()

    await waitFor(() => {
      expect(api.api.listContainers).toHaveBeenCalled()
    })
  })

  it('handles API errors gracefully', async () => {
    vi.mocked(api.api.listContainers).mockRejectedValue(new Error('API Error'))
    renderContainers()

    await waitFor(() => {
      // Component should still render even on error
      expect(screen.getByText(/containers/i)).toBeInTheDocument()
    })
  })
})

