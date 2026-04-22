import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { authApi } from '@/features/auth/services/auth'

vi.mock('@/features/auth/services/auth')

// Set up localStorage mock BEFORE importing useAuthStore
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}

vi.stubGlobal('localStorage', localStorageMock)

// Import after localStorage is mocked
import { useAuthStore } from './useAuthStore'

describe('useAuthStore', () => {
  beforeEach(() => {
    // Reset mocks
    localStorageMock.getItem.mockReturnValue(null)
    localStorageMock.setItem.mockImplementation(() => {})
    localStorageMock.removeItem.mockImplementation(() => {})
    localStorageMock.clear.mockImplementation(() => {})

    // Reset store before each test
    useAuthStore.getState().reset()
    vi.clearAllMocks()
  })

  afterEach(() => {
    // Clean up
    useAuthStore.getState().reset()
  })

  describe('initialize', () => {
    it('should mark store as initialized', () => {
      useAuthStore.getState().initialize()

      expect(useAuthStore.getState().isInitialized).toBe(true)
    })

    it('should handle missing tokens gracefully', () => {
      localStorageMock.getItem.mockReturnValue(null)

      useAuthStore.getState().initialize()

      expect(useAuthStore.getState().isAuthenticated).toBe(false)
      expect(useAuthStore.getState().isInitialized).toBe(true)
    })

    it('should clear expired tokens on initialization', () => {
      // Create an expired token (exp in the past)
      const expiredPayload = btoa(
        JSON.stringify({ exp: Math.floor(Date.now() / 1000) - 3600 })
      )
      const expiredToken = `header.${expiredPayload}.signature`

      localStorageMock.getItem.mockImplementation((key) => {
        if (key === 'access_token') return expiredToken
        if (key === 'refresh_token') return 'refresh'
        return null
      })

      useAuthStore.getState().initialize()

      expect(localStorageMock.removeItem).toHaveBeenCalledWith('access_token')
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')
      expect(useAuthStore.getState().isAuthenticated).toBe(false)
    })
  })

  describe('setTokens', () => {
    it('should update tokens in localStorage and state', () => {
      useAuthStore.getState().setTokens('new-access', 'new-refresh')

      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'access_token',
        'new-access'
      )
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'refresh_token',
        'new-refresh'
      )
      expect(useAuthStore.getState().token).toBe('new-access')
      expect(useAuthStore.getState().refreshToken).toBe('new-refresh')
    })
  })

  describe('login', () => {
    it('should write tokens to localStorage and update state on success', async () => {
      const mockUser = {
        id: 'user-123',
        username: 'testuser',
        email: 'test@example.com',
        role: 'normal' as const,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      vi.mocked(authApi.login).mockResolvedValue({
        access_token: 'test-token',
        refresh_token: 'test-refresh',
        user: mockUser,
      })

      await useAuthStore.getState().login('test@example.com', 'password123')

      // Verify state updated
      expect(useAuthStore.getState().user).toEqual(mockUser)
      expect(useAuthStore.getState().token).toBe('test-token')
      expect(useAuthStore.getState().refreshToken).toBe('test-refresh')
      expect(useAuthStore.getState().isAuthenticated).toBe(true)
      expect(useAuthStore.getState().isInitialized).toBe(true)

      // Verify localStorage called
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'access_token',
        'test-token'
      )
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'refresh_token',
        'test-refresh'
      )
    })

    it('should call authApi.login with correct parameters', async () => {
      vi.mocked(authApi.login).mockResolvedValue({
        access_token: 'token',
        refresh_token: 'refresh',
        user: {
          id: 'user-1',
          username: 'user',
          email: 'user@example.com',
          role: 'normal',
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      })

      await useAuthStore.getState().login('user@example.com', 'password')

      expect(authApi.login).toHaveBeenCalledWith({
        email: 'user@example.com',
        password: 'password',
      })
    })
  })

  describe('logout', () => {
    it('should clear tokens and state', () => {
      // Set initial state
      useAuthStore.setState({
        user: {
          id: 'user-123',
          username: 'testuser',
          email: 'test@example.com',
          role: 'normal',
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        token: 'test-token',
        refreshToken: 'test-refresh',
        isAuthenticated: true,
      })

      useAuthStore.getState().logout()

      // Verify state cleared
      expect(useAuthStore.getState().user).toBeNull()
      expect(useAuthStore.getState().token).toBeNull()
      expect(useAuthStore.getState().refreshToken).toBeNull()
      expect(useAuthStore.getState().isAuthenticated).toBe(false)

      // Verify localStorage cleared
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('access_token')
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')
    })
  })

  describe('refresh', () => {
    it('should update tokens on success', async () => {
      // Set initial state
      useAuthStore.setState({
        user: {
          id: 'user-123',
          username: 'testuser',
          email: 'test@example.com',
          role: 'normal',
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        token: 'old-token',
        refreshToken: 'old-refresh',
        isAuthenticated: true,
      })

      vi.mocked(authApi.refresh).mockResolvedValue({
        access_token: 'new-token',
        refresh_token: 'new-refresh',
      })

      await useAuthStore.getState().refresh('old-refresh')

      // Verify state updated
      expect(useAuthStore.getState().token).toBe('new-token')
      expect(useAuthStore.getState().refreshToken).toBe('new-refresh')

      // Verify localStorage updated
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'access_token',
        'new-token'
      )
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'refresh_token',
        'new-refresh'
      )
    })

    it('should call logout on refresh failure', async () => {
      // Set initial state
      useAuthStore.setState({
        user: {
          id: 'user-123',
          username: 'testuser',
          email: 'test@example.com',
          role: 'normal',
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        token: 'test-token',
        refreshToken: 'test-refresh',
        isAuthenticated: true,
      })

      const error = new Error('无效或过期的刷新令牌')
      vi.mocked(authApi.refresh).mockRejectedValue(error)

      await expect(
        useAuthStore.getState().refresh('invalid-token')
      ).rejects.toThrow('无效或过期的刷新令牌')

      // Verify logout was called (state cleared)
      expect(useAuthStore.getState().user).toBeNull()
      expect(useAuthStore.getState().isAuthenticated).toBe(false)
    })
  })

  describe('setUser', () => {
    it('should update user data', () => {
      const newUser = {
        id: 'user-456',
        username: 'newuser',
        email: 'new@example.com',
        role: 'admin',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      useAuthStore.getState().setUser(newUser)

      expect(useAuthStore.getState().user).toEqual(newUser)
    })
  })

  describe('reset', () => {
    it('should reset state to initial values', () => {
      // Set some state
      useAuthStore.setState({
        user: {
          id: 'user-123',
          username: 'testuser',
          email: 'test@example.com',
          role: 'normal',
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        token: 'test-token',
        refreshToken: 'test-refresh',
        isAuthenticated: true,
      })

      useAuthStore.getState().reset()

      expect(useAuthStore.getState().user).toBeNull()
      expect(useAuthStore.getState().token).toBeNull()
      expect(useAuthStore.getState().refreshToken).toBeNull()
      expect(useAuthStore.getState().isAuthenticated).toBe(false)
    })
  })
})
