import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { authApi } from './auth'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse } from 'msw'
import * as requestModule from '../../../lib/request'

describe('authApi', () => {
  beforeEach(() => {
    server.listen()
  })

  afterEach(() => {
    server.close()
  })

  describe('login', () => {
    it('should call POST /auth/login and return tokens', async () => {
      server.use(
        http.post('/api/v1/auth/login', () =>
          HttpResponse.json({
            access_token: 'test-access-token',
            refresh_token: 'test-refresh-token',
            user: {
              id: 'user-123',
              username: 'testuser',
              email: 'test@example.com',
              role: 'normal',
              createdAt: '2024-01-01T00:00:00Z',
              updatedAt: '2024-01-01T00:00:00Z',
            },
          })
        )
      )

      const result = await authApi.login({
        email: 'test@example.com',
        password: 'password123',
      })

      expect(result).toEqual({
        access_token: 'test-access-token',
        refresh_token: 'test-refresh-token',
        user: {
          id: 'user-123',
          username: 'testuser',
          email: 'test@example.com',
          role: 'normal',
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      })
    })

    it('should throw error on 401 unauthorized', async () => {
      server.use(
        http.post('/api/v1/auth/login', () =>
          HttpResponse.json({ error: '邮箱或密码错误' }, { status: 401 })
        )
      )

      await expect(
        authApi.login({ email: 'test@example.com', password: 'wrong' })
      ).rejects.toThrow()
    })

    it('should throw error on network failure', async () => {
      server.use(http.post('/api/v1/auth/login', () => HttpResponse.error()))

      await expect(
        authApi.login({ email: 'test@example.com', password: 'password' })
      ).rejects.toThrow()
    })
  })

  describe('register', () => {
    it('should call POST /auth/register and return user', async () => {
      server.use(
        http.post('/api/v1/auth/register', () =>
          HttpResponse.json({
            id: 'user-456',
            username: 'newuser',
            email: 'new@example.com',
            role: 'normal',
            createdAt: '2024-01-01T00:00:00Z',
            updatedAt: '2024-01-01T00:00:00Z',
          })
        )
      )

      const result = await authApi.register({
        username: 'newuser',
        email: 'new@example.com',
        password: 'password123',
        role: 'normal',
      })

      expect(result).toEqual({
        id: 'user-456',
        username: 'newuser',
        email: 'new@example.com',
        role: 'normal',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      })
    })

    it('should throw error on 409 conflict (email already exists)', async () => {
      server.use(
        http.post('/api/v1/auth/register', () =>
          HttpResponse.json({ error: '邮箱已存在' }, { status: 409 })
        )
      )

      await expect(
        authApi.register({
          username: 'testuser',
          email: 'existing@example.com',
          password: 'password123',
          role: 'normal',
        })
      ).rejects.toThrow()
    })
  })

  describe('refresh', () => {
    it('should call POST /auth/refresh and return new tokens', async () => {
      server.use(
        http.post('/api/v1/auth/refresh', () =>
          HttpResponse.json({
            access_token: 'new-access-token',
            refresh_token: 'new-refresh-token',
          })
        )
      )

      const result = await authApi.refresh('valid-refresh-token')

      expect(result).toEqual({
        access_token: 'new-access-token',
        refresh_token: 'new-refresh-token',
      })
    })

    it('should throw error on 401 (invalid refresh token)', async () => {
      // Mock post to bypass the token refresh interceptor
      const mockPost = vi
        .spyOn(requestModule, 'post')
        .mockRejectedValue(new Error('Request failed with status code 401'))

      await expect(authApi.refresh('invalid-token')).rejects.toThrow()

      mockPost.mockRestore()
    })
  })
})
