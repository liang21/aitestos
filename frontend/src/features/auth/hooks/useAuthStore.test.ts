import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuthStore } from './useAuthStore'
import { authApi } from '@/features/auth/services/auth'

vi.mock('@/features/auth/services/auth')

type MockLocalStorage = {
  getItem: ReturnType<typeof vi.fn>
  setItem: ReturnType<typeof vi.fn>
  removeItem: ReturnType<typeof vi.fn>
  clear: ReturnType<typeof vi.fn>
}

describe('useAuthStore', () => {
  let mockLocalStorage: MockLocalStorage

  beforeEach(() => {
    // Create localStorage mock before store is used
    mockLocalStorage = {
      getItem: vi.fn(),
      setItem: vi.fn(),
      removeItem: vi.fn(),
      clear: vi.fn(),
    }
    vi.stubGlobal('localStorage', mockLocalStorage)

    // Reset store before each test
    useAuthStore.getState().reset()
    vi.clearAllMocks()
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

      // Verify localStorage called
      expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
        'access_token',
        'test-token'
      )
      expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
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
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('access_token')
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('refresh_token')
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
      expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
        'access_token',
        'new-token'
      )
      expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
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
