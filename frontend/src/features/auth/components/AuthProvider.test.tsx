import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render } from '@testing-library/react'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'
import { AuthProvider } from './AuthProvider'

// Mock request handlers
const mockSetAuthExpiredHandler = vi.fn()
const mockSetTokenUpdatedHandler = vi.fn()

vi.mock('@/lib/request', () => ({
  setAuthExpiredHandler: (handler: () => void) => mockSetAuthExpiredHandler(handler),
  setTokenUpdatedHandler: (handler: (a: string, b: string) => void) =>
    mockSetTokenUpdatedHandler(handler),
}))

describe('AuthProvider', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    useAuthStore.getState().reset()
    // Mock localStorage
    const localStorageMock = {
      getItem: vi.fn(),
      setItem: vi.fn(),
      removeItem: vi.fn(),
      clear: vi.fn(),
    }
    vi.stubGlobal('localStorage', localStorageMock)
  })

  it('should initialize auth state on mount', () => {
    // Get the mocked localStorage
    const localStorageMock = global.localStorage as {
      getItem: ReturnType<typeof vi.fn>
      setItem: ReturnType<typeof vi.fn>
      removeItem: ReturnType<typeof vi.fn>
      clear: ReturnType<typeof vi.fn>
    }

    // Set tokens in localStorage
    localStorageMock.getItem.mockImplementation((key) => {
      if (key === 'access_token') return 'test-token'
      if (key === 'refresh_token') return 'test-refresh'
      return null
    })

    render(<AuthProvider>Test</AuthProvider>)

    // Verify handlers were registered
    expect(mockSetAuthExpiredHandler).toHaveBeenCalled()
    expect(mockSetTokenUpdatedHandler).toHaveBeenCalled()

    // Verify auth state was initialized
    expect(useAuthStore.getState().isInitialized).toBe(true)
  })

  it('should register logout handler with request interceptor', () => {
    render(<AuthProvider>Test</AuthProvider>)

    expect(mockSetAuthExpiredHandler).toHaveBeenCalledWith(
      expect.any(Function)
    )
  })

  it('should register setTokens handler with request interceptor', () => {
    render(<AuthProvider>Test</AuthProvider>)

    expect(mockSetTokenUpdatedHandler).toHaveBeenCalledWith(
      expect.any(Function)
    )
  })

  it('should handle empty localStorage gracefully', () => {
    const localStorageMock = global.localStorage as {
      getItem: ReturnType<typeof vi.fn>
      setItem: ReturnType<typeof vi.fn>
      removeItem: ReturnType<typeof vi.fn>
      clear: ReturnType<typeof vi.fn>
    }

    localStorageMock.getItem.mockReturnValue(null)

    render(<AuthProvider>Test</AuthProvider>)

    expect(useAuthStore.getState().isInitialized).toBe(true)
    expect(useAuthStore.getState().isAuthenticated).toBe(false)
  })
})
