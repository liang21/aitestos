/**
 * JWT Token Refresh Integration Tests
 * Tests automatic token refresh on 401 responses and concurrent request handling
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse, delay } from 'msw'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'
import request from '@/lib/request'
import { setAuthExpiredHandler, setTokenUpdatedHandler } from '@/lib/request'
import type { UserRole } from '@/types/enums'

// Mock logger
vi.mock('@/lib/logger', () => ({
  logger: {
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  },
  logAuthError: vi.fn(),
}))

describe('JWT Token Refresh Integration Tests', () => {
  const mockUser = {
    id: 'user-123',
    username: 'testuser',
    email: 'test@example.com',
    role: 'admin' as UserRole,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  }

  beforeEach(() => {
    vi.clearAllMocks()
    useAuthStore.getState().reset()
    server.resetHandlers()

    // Setup initial localStorage state
    localStorage.setItem('access_token', 'expired-token')
    localStorage.setItem('refresh_token', 'valid-refresh-token')

    // Setup token update handler
    setTokenUpdatedHandler((access, refresh) => {
      useAuthStore.getState().setTokens(access, refresh)
    })
  })

  afterEach(() => {
    // Clean up localStorage
    try {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    } catch (e) {
      // Ignore if localStorage is not available
    }
  })

  describe('Automatic Token Refresh', () => {
    it('should automatically refresh token on 401 response', async () => {
      let refreshCallCount = 0
      let apiCallCount = 0

      // Mock refresh endpoint
      server.use(
        http.post('/api/v1/auth/refresh', async () => {
          refreshCallCount++
          await delay(50)
          return HttpResponse.json({
            access_token: 'new-access-token',
            refresh_token: 'new-refresh-token',
          })
        })
      )

      // Mock API that returns 401 first, then 200 after retry
      server.use(
        http.get('/api/v1/projects', async ({ request }) => {
          apiCallCount++
          const authHeader = request.headers.get('Authorization')

          // First call with expired token returns 401
          if (authHeader === 'Bearer expired-token') {
            return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 })
          }

          // Subsequent calls with new token return 200
          if (authHeader === 'Bearer new-access-token') {
            return HttpResponse.json({
              data: [
                {
                  id: 'proj-1',
                  name: 'Test Project',
                  prefix: 'TEST',
                  description: '',
                  createdAt: '2024-01-01T00:00:00Z',
                  updatedAt: '2024-01-01T00:00:00Z',
                },
              ],
              total: 1,
              offset: 0,
              limit: 10,
            })
          }

          return HttpResponse.json({ error: 'Invalid token' }, { status: 401 })
        })
      )

      // Make API call that will trigger refresh
      const result = await request.get('/api/v1/projects')

      // Verify refresh was called once
      expect(refreshCallCount).toBe(1)

      // Verify API was called twice (original + retry)
      expect(apiCallCount).toBe(2)

      // Verify result is correct
      expect(result.data).toHaveLength(1)
      expect(result.data[0].name).toBe('Test Project')

      // Verify tokens were updated
      expect(localStorage.getItem('access_token')).toBe('new-access-token')
      expect(localStorage.getItem('refresh_token')).toBe('new-refresh-token')
    })

    it('should handle concurrent requests during token refresh', async () => {
      let refreshCallCount = 0

      // Mock refresh endpoint with delay
      server.use(
        http.post('/api/v1/auth/refresh', async () => {
          refreshCallCount++
          await delay(100)
          return HttpResponse.json({
            access_token: 'new-access-token',
            refresh_token: 'new-refresh-token',
          })
        })
      )

      // Mock API that returns 401
      server.use(
        http.get('/api/v1/testcases', async ({ request }) => {
          const authHeader = request.headers.get('Authorization')

          if (authHeader === 'Bearer expired-token') {
            return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 })
          }

          return HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        })
      )

      // Make concurrent requests
      const promises = [
        request.get('/api/v1/testcases?offset=0'),
        request.get('/api/v1/testcases?offset=10'),
        request.get('/api/v1/testcases?offset=20'),
      ]

      await Promise.all(promises)

      // Verify refresh was called only once (not 3 times)
      expect(refreshCallCount).toBe(1)

      // Verify all requests succeeded
      const results = await Promise.all(promises)
      results.forEach((result: { data: unknown[] }) => {
        expect(result.data).toEqual([])
      })
    })

    it('should logout on refresh failure', async () => {
      let authExpiredCalled = false

      // Setup auth expired handler
      setAuthExpiredHandler(() => {
        authExpiredCalled = true
        useAuthStore.getState().logout()
      })

      // Mock failing refresh endpoint
      server.use(
        http.post('/api/v1/auth/refresh', () =>
          HttpResponse.json({ error: 'Invalid refresh token' }, { status: 401 })
        )
      )

      // Mock API that returns 401
      server.use(
        http.get('/api/v1/projects', () =>
          HttpResponse.json({ error: 'Unauthorized' }, { status: 401 })
        )
      )

      // Make API call
      await expect(request.get('/api/v1/projects')).rejects.toThrow()

      // Verify auth expired handler was called
      expect(authExpiredCalled).toBe(true)

      // Verify tokens were cleared
      expect(localStorage.getItem('access_token')).toBeNull()
      expect(localStorage.getItem('refresh_token')).toBeNull()
    })
  })

  describe('Token Refresh with Auth Store Integration', () => {
    it('should update auth store after successful refresh', async () => {
      // Set initial auth state
      useAuthStore.setState({
        user: mockUser,
        token: 'expired-token',
        refreshToken: 'valid-refresh-token',
        isAuthenticated: true,
      })

      expect(useAuthStore.getState().token).toBe('expired-token')

      // Mock refresh endpoint
      server.use(
        http.post('/api/v1/auth/refresh', async () => {
          await delay(50)
          return HttpResponse.json({
            access_token: 'new-access-token',
            refresh_token: 'new-refresh-token',
          })
        })
      )

      // Mock API that triggers refresh
      server.use(
        http.get('/api/v1/projects', async ({ request }) => {
          const authHeader = request.headers.get('Authorization')

          if (authHeader === 'Bearer expired-token') {
            return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 })
          }

          return HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        })
      )

      // Trigger refresh by making API call
      await request.get('/api/v1/projects')

      // Verify auth store was updated
      expect(useAuthStore.getState().token).toBe('new-access-token')
      expect(useAuthStore.getState().refreshToken).toBe('new-refresh-token')
      expect(useAuthStore.getState().isAuthenticated).toBe(true)
    })

    it('should clear auth state on refresh failure', async () => {
      // Set initial auth state
      useAuthStore.setState({
        user: mockUser,
        token: 'expired-token',
        refreshToken: 'invalid-refresh-token',
        isAuthenticated: true,
      })

      // Setup auth expired handler
      setAuthExpiredHandler(() => {
        useAuthStore.getState().logout()
      })

      // Mock failing refresh
      server.use(
        http.post('/api/v1/auth/refresh', () =>
          HttpResponse.json({ error: 'Invalid refresh token' }, { status: 401 })
        )
      )

      // Mock API that returns 401
      server.use(
        http.get('/api/v1/projects', () =>
          HttpResponse.json({ error: 'Unauthorized' }, { status: 401 })
        )
      )

      // Make API call that will fail
      await expect(request.get('/api/v1/projects')).rejects.toThrow()

      // Verify auth store was cleared
      expect(useAuthStore.getState().isAuthenticated).toBe(false)
      expect(useAuthStore.getState().token).toBeNull()
      expect(useAuthStore.getState().user).toBeNull()
    })
  })

  describe('Edge Cases', () => {
    it('should skip refresh for auth endpoints', async () => {
      let refreshCallCount = 0

      // Mock refresh endpoint (should not be called)
      server.use(
        http.post('/api/v1/auth/refresh', () => {
          refreshCallCount++
          return HttpResponse.json({
            access_token: 'new-token',
            refresh_token: 'new-refresh',
          })
        })
      )

      // Mock login endpoint returning 401 (wrong password)
      server.use(
        http.post('/api/v1/auth/login', () =>
          HttpResponse.json({ error: '邮箱或密码错误' }, { status: 401 })
        )
      )

      // Make login request with wrong credentials
      await expect(
        request.post('/api/v1/auth/login', {
          email: 'test@example.com',
          password: 'wrong',
        })
      ).rejects.toThrow()

      // Verify refresh was NOT called
      expect(refreshCallCount).toBe(0)
    })

    it('should handle rapid successive 401 responses', async () => {
      let refreshCallCount = 0

      // Mock refresh endpoint with delay
      server.use(
        http.post('/api/v1/auth/refresh', async () => {
          refreshCallCount++
          await delay(200)
          return HttpResponse.json({
            access_token: 'new-access-token',
            refresh_token: 'new-refresh-token',
          })
        })
      )

      // Mock API endpoints
      const endpointIds = ['projects', 'testcases', 'plans']

      endpointIds.forEach((id) => {
        server.use(
          http.get(`/api/v1/${id}`, async ({ request }) => {
            const authHeader = request.headers.get('Authorization')

            if (authHeader === 'Bearer expired-token') {
              return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 })
            }

            return HttpResponse.json({
              data: [],
              total: 0,
              offset: 0,
              limit: 10,
            })
          })
        )
      })

      // Make rapid successive requests
      const promises = endpointIds.map((id) => request.get(`/api/v1/${id}`))

      await Promise.all(promises)

      // Verify refresh was called only once
      expect(refreshCallCount).toBe(1)
    })
  })
})
