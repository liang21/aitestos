import { type ReactNode, useEffect, useRef } from 'react'
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
  const initializedRef = useRef(false)

  useEffect(() => {
    // Only run once
    if (initializedRef.current) return
    initializedRef.current = true

    // Initialize auth state from localStorage
    useAuthStore.getState().initialize()

    // Register handlers with request interceptor - use store methods directly
    const store = useAuthStore.getState()
    setAuthExpiredHandler(store.logout)
    setTokenUpdatedHandler(store.setTokens)
  }, [])

  return <>{children}</>
}
