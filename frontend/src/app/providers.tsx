import { useEffect } from 'react'
import type { ReactNode } from 'react'
import { QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { ConfigProvider } from '@arco-design/web-react'
import { queryClient } from '@/lib/query-client'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'

interface ProvidersProps {
  children: ReactNode
}

/**
 * Validate and clean up stored tokens on app startup
 */
function TokenValidator() {
  const reset = useAuthStore((state) => state.reset)

  useEffect(() => {
    const validateAndCleanupTokens = () => {
      const token = localStorage.getItem('access_token')

      if (!token) {
        return
      }

      try {
        // Check if token is valid JWT format
        const parts = token.split('.')
        if (parts.length !== 3) {
          throw new Error('Invalid JWT format')
        }

        // Decode and check expiration
        const payload = JSON.parse(atob(parts[1]))
        const exp = payload.exp
        const now = Math.floor(Date.now() / 1000)

        if (exp && exp < now) {
          // Token expired - clean up
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          reset()
        }
      } catch {
        // Invalid token - clean up
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        reset()
      }
    }

    validateAndCleanupTokens()
  }, [reset])

  return null
}

/**
 * Application Providers
 *
 * Wraps the app with:
 * - QueryClientProvider (React Query)
 * - ConfigProvider (Arco Design theme)
 * - ReactQueryDevtools (development only)
 * - TokenValidator (cleanup expired tokens on startup)
 */
export function Providers({ children }: ProvidersProps) {
  return (
    <QueryClientProvider client={queryClient}>
      <ConfigProvider>
        <TokenValidator />
        {children}
      </ConfigProvider>
      {import.meta.env.DEV && <ReactQueryDevtools initialIsOpen={false} />}
    </QueryClientProvider>
  )
}
