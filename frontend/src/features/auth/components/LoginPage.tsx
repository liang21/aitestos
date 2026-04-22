import { useEffect } from 'react'
import { Navigate } from 'react-router-dom'
import { LoginBanner } from './LoginBanner'
import { LoginForm } from './LoginForm'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'

/**
 * LoginPage Component
 *
 * Based on MeterSphere login page design:
 * - Left: Banner with illustration (50% width)
 * - Right: Login form with branding (50% width)
 *
 * Features:
 * - Redirects authenticated users to their intended destination
 * - Stores original destination in router state for post-login redirect
 */
export function LoginPage() {
  const { isAuthenticated } = useAuthStore()

  // If already authenticated, the LoginForm will handle the redirect
  // This component just prevents rendering the login form
  if (isAuthenticated) {
    return <Navigate to="/projects" replace />
  }

  return (
    <div className="flex h-screen w-screen overflow-hidden overflow-x-hidden m-0 p-0">
      <LoginBanner />
      <LoginForm />
    </div>
  )
}
