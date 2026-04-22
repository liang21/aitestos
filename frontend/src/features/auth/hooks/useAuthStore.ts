import { create } from 'zustand'
import type { UserJSON } from '@/types/api'

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
  isInitialized: boolean

  // Actions
  initialize: () => void
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  refresh: (refreshToken: string) => Promise<void>
  setTokens: (accessToken: string, refreshToken: string) => void
  setUser: (user: UserJSON) => void
  reset: () => void
}

/**
 * Token storage interface for easier testing and SSR compatibility
 */
const tokenStorage = {
  getItem: (key: string): string | null => {
    try {
      return localStorage.getItem(key)
    } catch {
      return null
    }
  },
  setItem: (key: string, value: string): void => {
    try {
      localStorage.setItem(key, value)
    } catch {
      // Silently fail if localStorage unavailable
    }
  },
  removeItem: (key: string): void => {
    try {
      localStorage.removeItem(key)
    } catch {
      // Silently fail if localStorage unavailable
    }
  },
}

/**
 * Decode JWT and check expiration
 */
function isTokenExpired(token: string): boolean {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) {
      return true // Invalid JWT format
    }

    const payload = JSON.parse(atob(parts[1]))
    const exp = payload.exp
    if (!exp) {
      return true // No exp claim
    }

    const now = Math.floor(Date.now() / 1000)
    return exp < now
  } catch {
    return true // Treat invalid tokens as expired
  }
}

export const useAuthStore = create<AuthState>((set, get) => ({
  // Initial state
  user: null,
  token: null,
  refreshToken: null,
  isAuthenticated: false,
  isInitialized: false,

  /**
   * Initialize auth state from localStorage
   * Called on app startup to restore user session
   */
  initialize: () => {
    const accessToken = tokenStorage.getItem('access_token')
    const refreshTokenValue = tokenStorage.getItem('refresh_token')

    if (!accessToken || !refreshTokenValue) {
      set({ isInitialized: true })
      return
    }

    // Check if access token is expired
    if (isTokenExpired(accessToken)) {
      // Clear expired tokens
      tokenStorage.removeItem('access_token')
      tokenStorage.removeItem('refresh_token')
      set({ isInitialized: true })
      return
    }

    // We have valid tokens but user data might be lost
    // Mark as initialized but don't set user (will need to fetch if needed)
    set({
      token: accessToken,
      refreshToken: refreshTokenValue,
      isAuthenticated: true,
      isInitialized: true,
    })
  },

  /**
   * Login - authenticate user and store tokens
   */
  login: async (email: string, password: string) => {
    const { authApi } = await import('../services/auth')
    const response = await authApi.login({ email, password })

    // Store tokens in localStorage
    tokenStorage.setItem('access_token', response.access_token)
    tokenStorage.setItem('refresh_token', response.refresh_token)

    // Update state
    set({
      user: response.user,
      token: response.access_token,
      refreshToken: response.refresh_token,
      isAuthenticated: true,
      isInitialized: true,
    })
  },

  /**
   * Logout - clear tokens and state
   */
  logout: () => {
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
  refresh: async (refreshTokenValue: string) => {
    try {
      const { authApi } = await import('../services/auth')
      const response = await authApi.refresh(refreshTokenValue)

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
   * Set tokens directly (called by request.ts interceptor after refresh)
   * This keeps the store in sync with localStorage
   */
  setTokens: (accessToken: string, newRefreshToken: string) => {
    tokenStorage.setItem('access_token', accessToken)
    tokenStorage.setItem('refresh_token', newRefreshToken)

    set({
      token: accessToken,
      refreshToken: newRefreshToken,
    })
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
