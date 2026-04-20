import { useCallback } from 'react'
import type { AxiosError } from 'axios'

interface UseMutationErrorHandlerOptions {
  onError?: (error: Error) => void
}

/**
 * Hook to handle mutation errors consistently
 * Provides error parsing and user-friendly messages
 */
export function useMutationErrorHandler(
  options: UseMutationErrorHandlerOptions = {}
) {
  const handleError = useCallback(
    (error: unknown) => {
      if (options.onError) {
        options.onError(error as Error)
        return
      }

      // Default error handling
      let errorMessage = '操作失败，请稍后重试'

      if (isAxiosError(error)) {
        if (error.response?.data?.message) {
          errorMessage = error.response.data.message
        } else if (error.response?.status === 409) {
          errorMessage = '数据冲突，请刷新后重试'
        } else if (error.response?.status === 403) {
          errorMessage = '没有权限执行此操作'
        } else if (error.response?.status === 404) {
          errorMessage = '请求的资源不存在'
        } else if (error.code === 'ERR_NETWORK') {
          errorMessage = '网络连接失败，请检查网络设置'
        }
      } else if (error instanceof Error) {
        errorMessage = error.message
      }

      console.error('Mutation error:', error)
      return errorMessage
    },
    [options.onError]
  )

  return { handleError }
}

/**
 * Type guard for Axios errors
 */
function isAxiosError(error: unknown): error is AxiosError {
  return (
    typeof error === 'object' &&
    error !== null &&
    'isAxiosError' in error &&
    (error as AxiosError).isAxiosError === true
  )
}
