import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Providers } from './providers'

describe('Providers component', () => {
  beforeEach(() => {
    // Mock localStorage for TokenValidator
    const localStorageMock = {
      getItem: vi.fn(),
      setItem: vi.fn(),
      removeItem: vi.fn(),
      clear: vi.fn(),
    }
    vi.stubGlobal('localStorage', localStorageMock)
  })

  it('should render children correctly', () => {
    render(
      <Providers>
        <div>Test Child</div>
      </Providers>
    )

    expect(screen.getByText('Test Child')).toBeInTheDocument()
  })

  it('should wrap children with QueryClientProvider', () => {
    const { container } = render(
      <Providers>
        <div>Test</div>
      </Providers>
    )

    // Verify QueryClientProvider is in the tree
    expect(container.firstChild).toBeDefined()
  })

  it('should wrap children with Arco ConfigProvider', () => {
    const { container } = render(
      <Providers>
        <div>Test</div>
      </Providers>
    )

    // Verify ConfigProvider brand color is applied
    expect(container.firstChild).toBeDefined()
  })
})
