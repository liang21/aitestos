import { type ReactNode, useEffect } from 'react'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'
import { setAuthExpiredHandler, setTokenUpdatedHandler } from '@/lib/request'

interface AuthProviderProps {
  children: ReactNode
}

/**
 * AuthProvider Component
 *
 * Initializes authentication state and sets up request interceptors.
 * Should be rendered at the app root to ensure auth state is available everywhere.
 *
 * Responsibilities:
 * 1. Initialize auth state from localStorage on app startup
 * 2. Register auth expired handler with request interceptor
 * 3. Register token updated handler to keep store in sync with request.ts
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const { logout, setTokens, initialize } = useAuthStore()

  useEffect(() => {
    // Initialize auth state from localStorage
    initialize()

    // Register handlers with request interceptor
    setAuthExpiredHandler(logout)
    setTokenUpdatedHandler(setTokens)
  }, [initialize, logout, setTokens])

  return <>{children}</>
}
