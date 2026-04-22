import { type ReactNode } from 'react'
import { QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { ConfigProvider } from '@arco-design/web-react'
import { queryClient } from '@/lib/query-client'
import { AuthProvider } from '@/features/auth/components/AuthProvider'

interface ProvidersProps {
  children: ReactNode
}

/**
 * Application Providers
 *
 * Wraps the app with:
 * - QueryClientProvider (React Query)
 * - ConfigProvider (Arco Design theme)
 * - AuthProvider (authentication state initialization)
 * - ReactQueryDevtools (development only)
 */
export function Providers({ children }: ProvidersProps) {
  return (
    <QueryClientProvider client={queryClient}>
      <ConfigProvider>
        <AuthProvider>{children}</AuthProvider>
      </ConfigProvider>
      {import.meta.env.DEV && <ReactQueryDevtools initialIsOpen={false} />}
    </QueryClientProvider>
  )
}
