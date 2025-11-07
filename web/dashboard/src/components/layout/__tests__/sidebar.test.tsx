import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { Sidebar } from '../Sidebar'

const renderSidebar = () => {
  // Mock useLocation hook
  const MockSidebar = () => {
    // Since we can't easily mock useLocation in a simple way, we'll render with BrowserRouter
    return (
      <BrowserRouter>
        <Sidebar />
      </BrowserRouter>
    )
  }
  return render(<MockSidebar />)
}

describe('Sidebar', () => {
  it('renders navigation items', () => {
    renderSidebar()
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText('Containers')).toBeInTheDocument()
    expect(screen.getByText('Images')).toBeInTheDocument()
    expect(screen.getByText('Registry')).toBeInTheDocument()
    expect(screen.getByText('Monitor')).toBeInTheDocument()
  })

  it('renders NebulaBox brand', () => {
    renderSidebar()
    expect(screen.getByText(/NebulaBox/i)).toBeInTheDocument()
  })

  it('has links to all main sections', () => {
    renderSidebar()
    const links = [
      'Dashboard',
      'Containers',
      'Images',
      'Registry',
      'Build Spec',
      'Security',
      'Orchestrator',
      'Runtime',
      'AI Ops',
      'Monitor',
      'Networks',
      'Services',
      'Teams',
      'Settings',
    ]
    
    links.forEach((link) => {
      const linkElement = screen.getByRole('link', { name: new RegExp(link, 'i') })
      expect(linkElement).toBeInTheDocument()
    })
  })

  it('links have correct href attributes', () => {
    renderSidebar()
    expect(screen.getByRole('link', { name: /dashboard/i })).toHaveAttribute('href', '/')
    expect(screen.getByRole('link', { name: /containers/i })).toHaveAttribute('href', '/containers')
    expect(screen.getByRole('link', { name: /images/i })).toHaveAttribute('href', '/images')
  })
})

