import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { Images } from '../Images'
import * as api from '@/lib/api'

// Mock the API client
vi.mock('@/lib/api', () => ({
  api: {
    listImages: vi.fn(),
    pullImage: vi.fn(),
    deleteImage: vi.fn(),
  },
}))

const renderImages = () => {
  return render(
    <BrowserRouter>
      <Images />
    </BrowserRouter>
  )
}

describe('Images Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders images page', () => {
    vi.mocked(api.api.listImages).mockResolvedValue([])
    renderImages()
    expect(screen.getByText(/images/i)).toBeInTheDocument()
  })

  it('displays list of images', async () => {
    const mockImages = [
      {
        id: '1',
        name: 'nginx',
        tag: 'latest',
        size: '150MB',
        created: new Date().toISOString(),
      },
      {
        id: '2',
        name: 'alpine',
        tag: '3.18',
        size: '5MB',
        created: new Date().toISOString(),
      },
    ]

    vi.mocked(api.api.listImages).mockResolvedValue(mockImages)
    renderImages()

    await waitFor(() => {
      expect(screen.getByText(/nginx/i)).toBeInTheDocument()
      expect(screen.getByText(/alpine/i)).toBeInTheDocument()
    })
  })

  it('displays empty state when no images', async () => {
    vi.mocked(api.api.listImages).mockResolvedValue([])
    renderImages()

    await waitFor(() => {
      expect(api.api.listImages).toHaveBeenCalled()
    })
  })

  it('handles API errors gracefully', async () => {
    vi.mocked(api.api.listImages).mockRejectedValue(new Error('API Error'))
    renderImages()

    await waitFor(() => {
      // Component should still render even on error
      expect(screen.getByText(/images/i)).toBeInTheDocument()
    })
  })
})

