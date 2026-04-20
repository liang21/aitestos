import { type ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '../features/auth/hooks/useAuthStore'

interface RouteGuardProps {
  children: ReactNode
  requireAdmin?: boolean
}

/**
 * RouteGuard Component
 *
 * Protects routes by checking authentication status and user role
 * - Redirects to /login if not authenticated
 * - Checks token expiration via JWT exp claim
 * - Enforces admin role if requireAdmin is true
 */
export function RouteGuard({ children, requireAdmin = false }: RouteGuardProps) {
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

  // Check token expiration (decode JWT exp claim)
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    const exp = payload.exp
    const now = Math.floor(Date.now() / 1000)

    if (exp < now) {
      // Token expired
      logout()
      return (
        <Navigate
          to="/login"
          state={{ from: location.pathname + location.search }}
          replace
        />
      )
    }
  } catch {
    // Invalid token format
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
