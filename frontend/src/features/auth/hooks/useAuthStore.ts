import { create } from 'zustand'
import type { UserJSON } from '../../../types/api'

/**
 * Authentication state managed by Zustand
 * Handles user session, tokens, and authentication status
 */
interface AuthState {
  // State
  user: UserJSON | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean

  // Actions
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  refresh: (refreshToken: string) => Promise<void>
  setUser: (user: UserJSON) => void
  reset: () => void
}

export const useAuthStore = create<AuthState>((set, get) => ({
  // Initial state
  user: null,
  token: null,
  refreshToken: null,
  isAuthenticated: false,

  /**
   * Login - authenticate user and store tokens
   */
  login: async (email: string, password: string) => {
    const { authApi } = await import('../services/auth')
    const response = await authApi.login({ email, password })

    // Store tokens in localStorage
    const tokenStorage = {
      setItem: (key: string, value: string) => {
        try {
          localStorage.setItem(key, value)
        } catch {
          // Silently fail if localStorage unavailable
        }
      },
    }

    tokenStorage.setItem('access_token', response.access_token)
    tokenStorage.setItem('refresh_token', response.refresh_token)

    // Update state
    set({
      user: response.user,
      token: response.access_token,
      refreshToken: response.refresh_token,
      isAuthenticated: true,
    })
  },

  /**
   * Logout - clear tokens and state
   */
  logout: () => {
    const tokenStorage = {
      removeItem: (key: string) => {
        try {
          localStorage.removeItem(key)
        } catch {
          // Silently fail if localStorage unavailable
        }
      },
    }

    tokenStorage.removeItem('access_token')
    tokenStorage.removeItem('refresh_token')

    // Reset state
    set({
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
    })
  },

  /**
   * Refresh tokens - update with new tokens from backend
   */
  refresh: async (refreshToken: string) => {
    try {
      const { authApi } = await import('../services/auth')
      const response = await authApi.refresh(refreshToken)

      const tokenStorage = {
        setItem: (key: string, value: string) => {
          try {
            localStorage.setItem(key, value)
          } catch {
            // Silently fail if localStorage unavailable
          }
        },
      }

      tokenStorage.setItem('access_token', response.access_token)
      tokenStorage.setItem('refresh_token', response.refresh_token)

      set({
        token: response.access_token,
        refreshToken: response.refresh_token,
      })
    } catch (error) {
      // Refresh failed - logout user
      get().logout()
      throw error
    }
  },

  /**
   * Set user data (e.g., after profile update)
   */
  setUser: (user: UserJSON) => {
    set({ user })
  },

  /**
   * Reset state to initial values
   */
  reset: () => {
    set({
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
    })
  },
}))
