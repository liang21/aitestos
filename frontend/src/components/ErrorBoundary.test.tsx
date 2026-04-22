import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { ErrorBoundary } from './ErrorBoundary'

// Component that throws an error
function ThrowError({ shouldThrow }: { shouldThrow: boolean }) {
  if (shouldThrow) {
    throw new Error('Test error')
  }
  return <div>No error</div>
}

describe('ErrorBoundary', () => {
  it('should render children when there is no error', () => {
    render(
      <ErrorBoundary>
        <ThrowError shouldThrow={false} />
      </ErrorBoundary>
    )

    expect(screen.getByText('No error')).toBeInTheDocument()
  })

  it('should catch errors and display fallback UI', () => {
    // Suppress console.error for this test
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    render(
      <ErrorBoundary>
        <ThrowError shouldThrow={true} />
      </ErrorBoundary>
    )

    expect(screen.getByText('出现了一些问题')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '刷新页面' })).toBeInTheDocument()

    consoleSpy.mockRestore()
  })

  it('should display error message when provided', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    render(
      <ErrorBoundary>
        <ThrowError shouldThrow={true} />
      </ErrorBoundary>
    )

    expect(screen.getByText('Test error')).toBeInTheDocument()

    consoleSpy.mockRestore()
  })

  it('should call onError prop when error occurs', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    const onError = vi.fn()

    render(
      <ErrorBoundary onError={onError}>
        <ThrowError shouldThrow={true} />
      </ErrorBoundary>
    )

    expect(onError).toHaveBeenCalled()
    expect(onError).toHaveBeenCalledWith(
      expect.any(Error),
      expect.objectContaining({
        componentStack: expect.any(String),
      })
    )

    consoleSpy.mockRestore()
  })

  it('should reset error state when reset button is clicked', async () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    // Use a wrapper component with state to trigger error reset
    function TestWrapper() {
      const [shouldThrow, setShouldThrow] = React.useState(true)
      const [key, setKey] = React.useState(0)

      const handleReset = () => {
        setKey((prev) => prev + 1)
        setShouldThrow(false)
      }

      return (
        <ErrorBoundary
          key={key}
          fallback={
            <div>
              <div>Error fallback</div>
              <button onClick={handleReset}>重试</button>
            </div>
          }
        >
          <ThrowError shouldThrow={shouldThrow} />
        </ErrorBoundary>
      )
    }

    const { rerender } = render(<TestWrapper />)

    // Initial render with error
    expect(screen.getByText('Error fallback')).toBeInTheDocument()

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: '重试' }))

    // After reset, should render children without error
    await waitFor(() => {
      expect(screen.getByText('No error')).toBeInTheDocument()
    })

    consoleSpy.mockRestore()
  })
})
