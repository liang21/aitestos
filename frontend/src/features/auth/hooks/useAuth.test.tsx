import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useLogin, useRegister } from './useAuth'
import { useAuthStore } from './useAuthStore'
import { authApi } from '@/features/auth/services/auth'
import { server } from '@/tests/msw/server'

vi.mock('@/features/auth/services/auth')

describe('useAuth hooks', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
    vi.clearAllMocks()
    // Reset store
    useAuthStore.getState().reset()
    server.listen()
  })

  afterEach(() => {
    server.close()
  })

  function wrapper({ children }: { children: React.ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    )
  }

  describe('useLogin', () => {
    it('should call store.login on successful mutation', async () => {
      const mockUser = {
        id: 'user-123',
        username: 'testuser',
        email: 'test@example.com',
        role: 'normal' as const,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      vi.mocked(authApi.login).mockResolvedValue({
        access_token: 'test-token',
        refresh_token: 'test-refresh',
        user: mockUser,
      })

      const { result } = renderHook(() => useLogin(), { wrapper })

      result.current.mutate({ email: 'test@example.com', password: 'password' })

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      // Verify store.login was called
      expect(useAuthStore.getState().user).toEqual(mockUser)
      expect(useAuthStore.getState().token).toBe('test-token')
    })

    it('should throw error on login failure', async () => {
      vi.mocked(authApi.login).mockRejectedValue(new Error('邮箱或密码错误'))

      const { result } = renderHook(() => useLogin(), { wrapper })

      try {
        await result.current.mutateAsync({
          email: 'test@example.com',
          password: 'wrong',
        })
        expect.fail('Should have thrown error')
      } catch (error) {
        expect(error).toEqual(new Error('邮箱或密码错误'))
      }

      await waitFor(() => {
        expect(result.current.isError).toBe(true)
      })
    })
  })

  describe('useRegister', () => {
    it('should call authApi.register on mutation', async () => {
      const mockUser = {
        id: 'user-456',
        username: 'newuser',
        email: 'new@example.com',
        role: 'normal' as const,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      vi.mocked(authApi.register).mockResolvedValue(mockUser)

      const { result } = renderHook(() => useRegister(), { wrapper })

      result.current.mutate({
        username: 'newuser',
        email: 'new@example.com',
        password: 'password123',
        role: 'normal',
      })

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      expect(authApi.register).toHaveBeenCalledWith({
        username: 'newuser',
        email: 'new@example.com',
        password: 'password123',
        role: 'normal',
      })
    })

    it('should throw error on registration failure', async () => {
      vi.mocked(authApi.register).mockRejectedValue(new Error('邮箱已存在'))

      const { result } = renderHook(() => useRegister(), { wrapper })

      try {
        await result.current.mutateAsync({
          username: 'testuser',
          email: 'existing@example.com',
          password: 'password123',
          role: 'normal',
        })
        expect.fail('Should have thrown error')
      } catch (error) {
        expect(error).toEqual(new Error('邮箱已存在'))
      }

      await waitFor(() => {
        expect(result.current.isError).toBe(true)
      })
    })
  })
})
