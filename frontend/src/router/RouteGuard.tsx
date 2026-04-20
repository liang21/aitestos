import { type ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'

interface RouteGuardProps {
  children: ReactNode
  requireAdmin?: boolean
}

/**
 * Check if JWT token is expired by decoding the exp claim
 * @param token - JWT access token
 * @returns true if token is expired or invalid, false otherwise
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

/**
 * RouteGuard Component
 *
 * Protects routes by checking authentication status and user role
 * - Redirects to /login if not authenticated
 * - Checks token expiration via JWT exp claim
 * - Enforces admin role if requireAdmin is true
 */
export function RouteGuard({
  children,
  requireAdmin = false,
}: RouteGuardProps) {
  const location = useLocation()
  const { user, token, isAuthenticated, logout } = useAuthStore()

  // Check if authenticated
  if (!isAuthenticated || !user || !token) {
    // Store the intended destination for post-login redirect
    return (
      <Navigate
        to="/login"
        state={{ from: location.pathname + location.search }}
        replace
      />
    )
  }

  // Check token expiration
  if (isTokenExpired(token)) {
    logout()
    return (
      <Navigate
        to="/login"
        state={{ from: location.pathname + location.search }}
        replace
      />
    )
  }

  // Check admin role if required
  if (requireAdmin && user.role !== 'super_admin' && user.role !== 'admin') {
    // User is not admin, redirect to projects
    return <Navigate to="/projects" replace />
  }

  return <>{children}</>
}
