import { QueryClient } from '@tanstack/react-query'

/**
 * React Query Client with global defaults
 *
 * - staleTime: 5 minutes (data stays fresh for 5min)
 * - retry: 1 (retry failed requests once)
 * - retryDelay: exponential backoff
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      retry: 1,
      retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 1,
    },
  },
})
