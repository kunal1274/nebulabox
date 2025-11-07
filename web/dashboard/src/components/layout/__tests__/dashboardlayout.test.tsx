import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { DashboardLayout } from '../DashboardLayout'

const renderLayout = () => {
  return render(
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<DashboardLayout />}>
          <Route index element={<div>Test Content</div>} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

describe('DashboardLayout', () => {
  it('renders layout with sidebar navigation', () => {
    renderLayout()
    // Sidebar should be rendered (check for NebulaBox brand)
    expect(screen.getByText(/NebulaBox/i)).toBeInTheDocument()
  })

  it('renders outlet content', () => {
    renderLayout()
    // Content from Route should be rendered
    expect(screen.getByText('Test Content')).toBeInTheDocument()
  })
})

