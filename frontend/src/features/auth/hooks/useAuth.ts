import { useMutation } from '@tanstack/react-query'
import { useAuthStore } from './useAuthStore'
import { authApi } from '../services/auth'

/**
 * Login mutation hook
 * Handles user authentication and stores session
 */
export function useLogin() {
  const login = useMutation({
    mutationFn: async ({ email, password }: { email: string; password: string }) => {
      return authApi.login({ email, password })
    },
    onSuccess: async (data) => {
      // Store will be updated by authApi.login internally
      // But we need to update Zustand store separately
      const { user, token, refreshToken } = useAuthStore.getState()

      // Check if store was already updated by the API call
      if (useAuthStore.getState().token !== data.access_token) {
        // Update store with response data
        useAuthStore.setState({
          user: data.user,
          token: data.access_token,
          refreshToken: data.refresh_token,
          isAuthenticated: true,
        })
      }
    },
    onError: (error: Error) => {
      // Error will be handled by the component
      console.error('Login failed:', error)
    },
  })

  return login
}

/**
 * Register mutation hook
 * Handles user registration
 */
export function useRegister() {
  const register = useMutation({
    mutationFn: async ({
      username,
      email,
      password,
      role,
    }: {
      username: string
      email: string
      password: string
      role: 'normal' | 'super_admin' | 'admin'
    }) => {
      return authApi.register({ username, email, password, role })
    },
    onSuccess: () => {
      // Registration successful - component will handle navigation
    },
    onError: (error: Error) => {
      console.error('Registration failed:', error)
    },
  })

  return register
}

/**
 * Refresh token hook
 * Refreshes the access token using the refresh token
 */
export function useRefresh() {
  const { refreshToken } = useAuthStore()

  const refresh = useMutation({
    mutationFn: async () => {
      if (!refreshToken) {
        throw new Error('No refresh token available')
      }
      return authApi.refresh(refreshToken)
    },
    onSuccess: async (data) => {
      // Store will be updated by authApi.refresh internally
      // Check if store was already updated
      if (useAuthStore.getState().token !== data.access_token) {
        useAuthStore.setState({
          token: data.access_token,
          refreshToken: data.refresh_token,
        })
      }
    },
    onError: () => {
      // Refresh failed - logout will be called by useAuthStore.refresh
      console.error('Token refresh failed')
    },
  })

  return refresh
}
