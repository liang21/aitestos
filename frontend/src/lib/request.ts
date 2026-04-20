import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import type { RefreshResponse } from '@/types/api'

// ============================================================================
// Token Refresh State
// ============================================================================

interface QueuedRequest {
  request: InternalAxiosRequestConfig
  resolve: (value: unknown) => void
  reject: (reason: unknown) => void
}

let isRefreshing = false
const pendingRequests: QueuedRequest[] = []

// Auth expired callback
let authExpiredHandler: (() => void) | null = null

/**
 * Register callback for auth expiration (e.g., redirect to login)
 */
export function setAuthExpiredHandler(handler: () => void) {
  authExpiredHandler = handler
}

// ============================================================================
// Token Storage Interface (for easier testing and SSR compatibility)
// ============================================================================

const tokenStorage = {
  getItem: (key: string): string | null => {
    try {
      return localStorage.getItem(key)
    } catch {
      return null
    }
  },
  setItem: (key: string, value: string): void => {
    try {
      localStorage.setItem(key, value)
    } catch {
      // Silently fail if localStorage is unavailable
    }
  },
  removeItem: (key: string): void => {
    try {
      localStorage.removeItem(key)
    } catch {
      // Silently fail if localStorage is unavailable
    }
  },
}

// ============================================================================
// Axios Instance
// ============================================================================

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
})

// ============================================================================
// Request Interceptor
// ============================================================================

request.interceptors.request.use(
  (config) => {
    const token = tokenStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// ============================================================================
// Response Interceptor
// ============================================================================

request.interceptors.response.use(
  (response) => {
    return response.data
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean
      _queued?: boolean
    }

    // Handle 401 Unauthorized - Token refresh
    if (
      error.response?.status === 401 &&
      originalRequest &&
      !originalRequest._retry
    ) {
      // If already refreshing, queue this request
      if (isRefreshing) {
        // Prevent the same request from being queued multiple times
        if (originalRequest._queued) {
          return Promise.reject(error)
        }
        originalRequest._queued = true

        return new Promise((resolve, reject) => {
          pendingRequests.push({
            request: originalRequest,
            resolve,
            reject,
          })
        })
      }

      // Start refresh process
      originalRequest._retry = true
      isRefreshing = true

      const refreshToken = tokenStorage.getItem('refresh_token')

      if (!refreshToken) {
        handleAuthExpired()
        flushQueue(error)
        return Promise.reject(error)
      }

      try {
        // Call refresh endpoint
        const response = await axios.post<RefreshResponse>(
          `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/auth/refresh`,
          { refresh_token: refreshToken },
          { headers: { 'Content-Type': 'application/json' } }
        )

        const { access_token, refresh_token: newRefreshToken } = response.data

        // Update storage
        tokenStorage.setItem('access_token', access_token)
        tokenStorage.setItem('refresh_token', newRefreshToken)

        // Process queued requests: replay each with new token
        const queuedResults = await Promise.allSettled(
          pendingRequests.map(
            async ({ request: queuedReq, resolve, reject }) => {
              try {
                // Update Authorization header for queued request
                if (queuedReq.headers) {
                  queuedReq.headers.Authorization = `Bearer ${access_token}`
                }
                // Replay the request and return its result
                const result = await request(queuedReq)
                resolve(result)
                return result
              } catch (err) {
                reject(err)
                throw err
              }
            }
          )
        )

        // Clear queue
        pendingRequests.length = 0

        // Check if any queued requests failed
        const failedRequests = queuedResults.filter(
          (result): result is PromiseRejectedResult =>
            result.status === 'rejected'
        )

        if (failedRequests.length > 0) {
          console.error(
            `${failedRequests.length} queued request(s) failed after token refresh`
          )
        }

        // Retry original request
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${access_token}`
        }
        return request(originalRequest)
      } catch (refreshError) {
        // Refresh failed - clear tokens and trigger auth expired
        tokenStorage.removeItem('access_token')
        tokenStorage.removeItem('refresh_token')

        handleAuthExpired()

        // Reject all queued requests
        flushQueue(refreshError)

        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    // Handle other errors
    return Promise.reject(error)
  }
)

// ============================================================================
// Helper Functions
// ============================================================================

function handleAuthExpired() {
  if (authExpiredHandler) {
    authExpiredHandler()
  }
}

/**
 * Flush the pending request queue with a rejection reason
 */
function flushQueue(reason: unknown) {
  pendingRequests.forEach(({ reject }) => {
    reject(reason)
  })
  pendingRequests.length = 0
}

// ============================================================================
// Typed API Wrappers
// ============================================================================

/**
 * Typed GET request
 */
export function get<TResponse>(
  url: string,
  config?: InternalAxiosRequestConfig
) {
  return request.get<never, TResponse>(url, config)
}

/**
 * Typed POST request
 */
export function post<TRequest, TResponse>(
  url: string,
  data?: TRequest,
  config?: InternalAxiosRequestConfig
) {
  return request.post<never, TResponse>(url, data, config)
}

/**
 * Typed PUT request
 */
export function put<TRequest, TResponse>(
  url: string,
  data?: TRequest,
  config?: InternalAxiosRequestConfig
) {
  return request.put<never, TResponse>(url, data, config)
}

/**
 * Typed PATCH request
 */
export function patch<TRequest, TResponse>(
  url: string,
  data?: TRequest,
  config?: InternalAxiosRequestConfig
) {
  return request.patch<never, TResponse>(url, data, config)
}

/**
 * Typed DELETE request
 */
export function del<TResponse>(
  url: string,
  config?: InternalAxiosRequestConfig
) {
  return request.delete<never, TResponse>(url, config)
}

export default request
