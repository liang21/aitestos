import type { ReactNode } from 'react'
import { QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { ConfigProvider } from '@arco-design/web-react'
import { queryClient } from '../lib/query-client'

interface ProvidersProps {
  children: ReactNode
}

/**
 * Application Providers
 *
 * Wraps the app with:
 * - QueryClientProvider (React Query)
 * - ConfigProvider (Arco Design theme)
 * - ReactQueryDevtools (development only)
 */
export function Providers({ children }: ProvidersProps) {
  return (
    <QueryClientProvider client={queryClient}>
      <ConfigProvider>{children}</ConfigProvider>
      {import.meta.env.DEV && <ReactQueryDevtools initialIsOpen={false} />}
    </QueryClientProvider>
  )
}
