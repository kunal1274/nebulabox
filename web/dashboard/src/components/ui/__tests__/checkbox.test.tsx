import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Checkbox } from '../checkbox'

describe('Checkbox', () => {
  it('renders checkbox element', () => {
    render(<Checkbox id="test-checkbox" />)
    const checkbox = screen.getByRole('checkbox')
    expect(checkbox).toBeInTheDocument()
  })

  it('is checked when checked prop is true', () => {
    render(<Checkbox id="test-checkbox" checked />)
    const checkbox = screen.getByRole('checkbox')
    expect(checkbox).toBeChecked()
  })

  it('is unchecked by default', () => {
    render(<Checkbox id="test-checkbox" />)
    const checkbox = screen.getByRole('checkbox')
    expect(checkbox).not.toBeChecked()
  })

  it('handles click events', async () => {
    const handleCheckedChange = vi.fn()
    const user = userEvent.setup()
    
    render(
      <Checkbox
        id="test-checkbox"
        onCheckedChange={handleCheckedChange}
      />
    )
    const checkbox = screen.getByRole('checkbox')
    
    await user.click(checkbox)
    expect(handleCheckedChange).toHaveBeenCalled()
  })

  it('is disabled when disabled prop is true', () => {
    render(<Checkbox id="test-checkbox" disabled />)
    const checkbox = screen.getByRole('checkbox')
    expect(checkbox).toBeDisabled()
  })

  it('forwards additional props', () => {
    render(<Checkbox id="test-checkbox" aria-label="Test checkbox" />)
    const checkbox = screen.getByLabelText('Test checkbox')
    expect(checkbox).toBeInTheDocument()
  })
})

