import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { server } from '@/tests/msw/server'
import { http, HttpResponse } from 'msw'
import { LoginPage } from '@/features/auth/components/LoginPage'
import { RegisterPage } from '@/features/auth/components/RegisterPage'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'
import { AuthErrorBoundary } from '@/components/ErrorBoundary'

vi.mock('@/features/auth/services/auth')

describe('Authentication Integration Tests', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
    vi.clearAllMocks()
    useAuthStore.getState().reset()
    server.listen()
  })

  afterEach(() => {
    server.close()
  })

  function renderWithProviders(ui: React.ReactElement) {
    return render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <AuthErrorBoundary>{ui}</AuthErrorBoundary>
        </BrowserRouter>
      </QueryClientProvider>
    )
  }

  describe('Complete Registration Flow', () => {
    it('should successfully register a new user and redirect to login', async () => {
      const { authApi } = await import('@/features/auth/services/auth')

      // Mock successful registration
      vi.mocked(authApi).register = vi.fn().mockResolvedValue({
        id: 'user-456',
        username: 'newuser',
        email: 'new@example.com',
        role: 'normal',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      })

      renderWithProviders(<RegisterPage />)

      const user = userEvent.setup()

      // Fill out registration form
      await user.type(screen.getByPlaceholderText('请输入用户名'), 'newuser')
      await user.type(
        screen.getByPlaceholderText('请输入邮箱'),
        'new@example.com'
      )
      await user.type(
        screen.getByPlaceholderText('请输入密码（至少 8 位字符）'),
        'password123'
      )

      // Select role
      await user.click(screen.getByText('普通用户'))

      // Submit form
      await user.click(screen.getByRole('button', { name: '注册' }))

      // Verify API was called
      await waitFor(() => {
        expect(authApi.register).toHaveBeenCalledWith({
          username: 'newuser',
          email: 'new@example.com',
          password: 'password123',
          role: 'normal',
        })
      })

      // Verify user is not logged in after registration
      expect(useAuthStore.getState().isAuthenticated).toBe(false)
    })

    it('should handle registration errors gracefully', async () => {
      const { authApi } = await import('@/features/auth/services/auth')

      // Mock registration failure
      vi.mocked(authApi).register = vi
        .fn()
        .mockRejectedValue(new Error('邮箱已存在'))

      renderWithProviders(<RegisterPage />)

      const user = userEvent.setup()

      // Fill out form with existing email
      await user.type(screen.getByPlaceholderText('请输入用户名'), 'testuser')
      await user.type(
        screen.getByPlaceholderText('请输入邮箱'),
        'existing@example.com'
      )
      await user.type(
        screen.getByPlaceholderText('请输入密码（至少 8 位字符）'),
        'password123'
      )
      await user.click(screen.getByText('普通用户'))

      // Submit form
      await user.click(screen.getByRole('button', { name: '注册' }))

      // Verify error was handled (button should still be enabled)
      await waitFor(() => {
        expect(authApi.register).toHaveBeenCalled()
        expect(screen.getByRole('button', { name: '注册' })).toBeEnabled()
      })
    })
  })

  describe('Complete Login Flow', () => {
    it('should successfully login and update auth store', async () => {
      const { authApi } = await import('@/features/auth/services/auth')

      const mockUser = {
        id: 'user-123',
        username: 'testuser',
        email: 'test@example.com',
        role: 'normal',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      // Mock successful login
      vi.mocked(authApi).login = vi.fn().mockResolvedValue({
        access_token: 'test-token',
        refresh_token: 'test-refresh',
        user: mockUser,
      })

      renderWithProviders(<LoginPage />)

      const user = userEvent.setup()

      // Fill out login form
      await user.type(
        screen.getByPlaceholderText('请输入邮箱'),
        'test@example.com'
      )
      await user.type(screen.getByPlaceholderText('请输入密码'), 'password123')

      // Submit form
      await user.click(screen.getByRole('button', { name: '登录' }))

      // Verify API was called
      await waitFor(() => {
        expect(authApi.login).toHaveBeenCalledWith({
          email: 'test@example.com',
          password: 'password123',
        })
      })

      // Verify auth store was updated
      expect(useAuthStore.getState().user).toEqual(mockUser)
      expect(useAuthStore.getState().token).toBe('test-token')
      expect(useAuthStore.getState().isAuthenticated).toBe(true)
    })

    it('should handle login errors and keep user unauthenticated', async () => {
      const { authApi } = await import('@/features/auth/services/auth')

      // Mock login failure
      vi.mocked(authApi).login = vi
        .fn()
        .mockRejectedValue(new Error('邮箱或密码错误'))

      renderWithProviders(<LoginPage />)

      const user = userEvent.setup()

      // Fill out login form with wrong credentials
      await user.type(
        screen.getByPlaceholderText('请输入邮箱'),
        'test@example.com'
      )
      await user.type(
        screen.getByPlaceholderText('请输入密码'),
        'wrongpassword'
      )

      // Submit form
      await user.click(screen.getByRole('button', { name: '登录' }))

      // Verify API was called
      await waitFor(() => {
        expect(authApi.login).toHaveBeenCalled()
      })

      // Verify user is still unauthenticated
      expect(useAuthStore.getState().isAuthenticated).toBe(false)
      expect(useAuthStore.getState().user).toBeNull()
    })
  })

  describe('Protected Routes Behavior', () => {
    it('should redirect unauthenticated users to login', async () => {
      const { RouteGuard } = await import('@/router/RouteGuard')

      renderWithProviders(
        <Routes>
          <Route path="/login" element={<div>Login Page</div>} />
          <Route
            path="/protected"
            element={
              <RouteGuard>
                <div>Protected Content</div>
              </RouteGuard>
            }
          />
        </Routes>
      )

      // Navigate to protected route
      window.history.pushState({}, '', '/protected')

      // Should redirect to login
      await waitFor(() => {
        expect(screen.getByText('Login Page')).toBeInTheDocument()
        expect(screen.queryByText('Protected Content')).not.toBeInTheDocument()
      })
    })

    it('should allow authenticated users to access protected routes', async () => {
      const { RouteGuard } = await import('@/router/RouteGuard')

      // Set authenticated state
      useAuthStore.setState({
        user: {
          id: 'user-123',
          username: 'testuser',
          email: 'test@example.com',
          role: 'normal',
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        token: 'valid-token',
        refreshToken: 'valid-refresh',
        isAuthenticated: true,
      })

      renderWithProviders(
        <Routes>
          <Route path="/login" element={<div>Login Page</div>} />
          <Route
            path="/protected"
            element={
              <RouteGuard>
                <div>Protected Content</div>
              </RouteGuard>
            }
          />
        </Routes>
      )

      // Navigate to protected route
      window.history.pushState({}, '', '/protected')

      // Should show protected content
      await waitFor(() => {
        expect(screen.getByText('Protected Content')).toBeInTheDocument()
        expect(screen.queryByText('Login Page')).not.toBeInTheDocument()
      })
    })

    it('should redirect to login when token is expired', async () => {
      const { RouteGuard } = await import('@/router/RouteGuard')

      // Set authenticated state with expired token
      useAuthStore.setState({
        user: {
          id: 'user-123',
          username: 'testuser',
          email: 'test@example.com',
          role: 'normal',
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        token: 'expired-token',
        refreshToken: 'valid-refresh',
        isAuthenticated: true,
      })

      renderWithProviders(
        <Routes>
          <Route path="/login" element={<div>Login Page</div>} />
          <Route
            path="/protected"
            element={
              <RouteGuard>
                <div>Protected Content</div>
              </RouteGuard>
            }
          />
        </Routes>
      )

      // Navigate to protected route
      window.history.pushState({}, '', '/protected')

      // Should redirect to login due to expired token
      await waitFor(() => {
        expect(screen.getByText('Login Page')).toBeInTheDocument()
      })

      // Store should be cleared
      expect(useAuthStore.getState().isAuthenticated).toBe(false)
    })
  })

  describe('Error Boundary', () => {
    it('should catch errors in auth components and show fallback UI', () => {
      // This test verifies the ErrorBoundary works
      const ThrowError = () => {
        throw new Error('Component error')
      }

      renderWithProviders(
        <Routes>
          <Route
            path="/"
            element={
              <AuthErrorBoundary>
                <ThrowError />
              </AuthErrorBoundary>
            }
          />
        </Routes>
      )

      // Should show error fallback
      expect(screen.getByText('出现了一些问题')).toBeInTheDocument()
    })
  })
})
