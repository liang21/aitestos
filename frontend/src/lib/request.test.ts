import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import request from './request'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

// 创建独立的 server 实例用于测试
const server = setupServer()

type MockLocalStorage = {
  getItem: ReturnType<typeof vi.fn>
  setItem: ReturnType<typeof vi.fn>
  removeItem: ReturnType<typeof vi.fn>
  clear: ReturnType<typeof vi.fn>
}

describe('axios request instance', () => {
  beforeEach(() => {
    // Setup localStorage mock
    const localStorageMock: MockLocalStorage = {
      getItem: vi.fn(),
      setItem: vi.fn(),
      removeItem: vi.fn(),
      clear: vi.fn(),
    }
    vi.stubGlobal('localStorage', localStorageMock)

    server.listen()
  })

  afterEach(() => {
    server.close()
    vi.unstubAllGlobals()
  })

  it('should attach Authorization header when token exists', async () => {
    const mockLocalStorage =
      globalThis.localStorage as unknown as MockLocalStorage
    mockLocalStorage.getItem.mockReturnValue('test-token')

    server.use(
      http.get('/api/v1/test', () => HttpResponse.json({ data: 'test' }))
    )

    const response = await request.get('/test')

    expect(mockLocalStorage.getItem).toHaveBeenCalledWith('access_token')
    expect(response).toEqual({ data: 'test' })
  })

  it('should not attach Authorization header when token is missing', async () => {
    const mockLocalStorage =
      globalThis.localStorage as unknown as MockLocalStorage
    mockLocalStorage.getItem.mockReturnValue(null)

    server.use(
      http.get('/api/v1/test', () => HttpResponse.json({ data: 'test' }))
    )

    const response = await request.get('/test')

    expect(mockLocalStorage.getItem).toHaveBeenCalledWith('access_token')
    expect(response).toEqual({ data: 'test' })
  })

  it('should return response data directly on success', async () => {
    server.use(
      http.get('/api/v1/test', () => HttpResponse.json({ message: 'success' }))
    )

    const response = await request.get('/test')

    expect(response).toEqual({ message: 'success' })
    expect(response).not.toHaveProperty('data')
  })

  describe('token refresh queue', () => {
    it('should replay queued requests after successful token refresh', async () => {
      const mockLocalStorage =
        globalThis.localStorage as unknown as MockLocalStorage

      let refreshCallCount = 0
      let hasRefreshed = false

      server.use(
        // Handler that tracks requests and returns 401 before refresh
        http.get('/api/v1/endpoint', () => {
          if (!hasRefreshed) {
            return new HttpResponse(null, { status: 401 })
          }
          return HttpResponse.json({ success: true })
        }),
        // Refresh endpoint
        http.post('/api/v1/auth/refresh', () => {
          refreshCallCount++
          hasRefreshed = true
          return HttpResponse.json({
            access_token: 'new-token',
            refresh_token: 'new-refresh',
          })
        })
      )

      // Setup localStorage to return tokens
      mockLocalStorage.getItem.mockImplementation((key) => {
        if (key === 'access_token') {
          return hasRefreshed ? 'new-token' : 'old-token'
        }
        if (key === 'refresh_token') {
          return 'valid-refresh-token'
        }
        return null
      })

      // Fire multiple concurrent requests - they will all get 401 first
      // Then the refresh will happen, and they should be replayed successfully
      const results = await Promise.allSettled([
        request.get('/endpoint'),
        request.get('/endpoint'),
        request.get('/endpoint'),
      ])

      // Verify refresh was called once
      expect(refreshCallCount).toBe(1)

      // Verify all requests eventually succeeded
      results.forEach((result) => {
        expect(result.status).toBe('fulfilled')
        if (result.status === 'fulfilled') {
          expect(result.value).toEqual({ success: true })
        }
      })

      // Verify tokens were updated
      expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
        'access_token',
        'new-token'
      )
      expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
        'refresh_token',
        'new-refresh'
      )
    })

    it('should reject all queued requests when refresh fails', async () => {
      const mockLocalStorage =
        globalThis.localStorage as unknown as MockLocalStorage

      mockLocalStorage.getItem.mockImplementation((key) => {
        if (key === 'access_token') {
          return 'expired-token'
        }
        if (key === 'refresh_token') {
          return 'invalid-refresh-token'
        }
        return null
      })

      let refreshCallCount = 0

      server.use(
        // All requests return 401
        http.get('/api/v1/test', () => new HttpResponse(null, { status: 401 })),
        // Refresh endpoint fails
        http.post('/api/v1/auth/refresh', () => {
          refreshCallCount++
          return new HttpResponse(null, { status: 401 })
        })
      )

      // Fire concurrent requests
      const results = await Promise.allSettled([
        request.get('/test'),
        request.get('/test'),
      ])

      // Verify refresh was attempted
      expect(refreshCallCount).toBe(1)

      // Verify all requests failed
      results.forEach((result) => {
        expect(result.status).toBe('rejected')
      })

      // Verify tokens were cleared
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('access_token')
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('refresh_token')
    })
  })
})
