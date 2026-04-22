import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse } from 'msw'
import { LoginPage } from '@/features/auth/components/LoginPage'
import { RegisterPage } from '@/features/auth/components/RegisterPage'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'
import { useLogin, useRegister } from '@/features/auth/hooks/useAuth'

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
    server.resetHandlers()
  })

  function renderWithProviders(ui: React.ReactElement) {
    return render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>{ui}</BrowserRouter>
      </QueryClientProvider>
    )
  }

  describe('useLogin Hook Integration', () => {
    it('should successfully login and update auth store', async () => {
      const mockUser = {
        id: 'user-123',
        username: 'testuser',
        email: 'test@example.com',
        role: 'normal',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      server.use(
        http.post('/api/v1/auth/login', () =>
          HttpResponse.json({
            access_token: 'test-token',
            refresh_token: 'test-refresh',
            user: mockUser,
          })
        )
      )

      function TestComponent() {
        const login = useLogin()
        return (
          <button
            onClick={() =>
              login.mutate({
                email: 'test@example.com',
                password: 'password123',
              })
            }
          >
            Login
          </button>
        )
      }

      renderWithProviders(<TestComponent />)

      // Initially not authenticated
      expect(useAuthStore.getState().isAuthenticated).toBe(false)

      // Trigger login
      await screen.getByRole('button').click()

      // Verify auth store was updated
      await waitFor(() => {
        expect(useAuthStore.getState().user).toEqual(mockUser)
        expect(useAuthStore.getState().token).toBe('test-token')
        expect(useAuthStore.getState().isAuthenticated).toBe(true)
      })
    })

    it('should handle login errors', async () => {
      server.use(
        http.post('/api/v1/auth/login', () =>
          HttpResponse.json({ error: '邮箱或密码错误' }, { status: 401 })
        )
      )

      function TestComponent() {
        const login = useLogin()
        return (
          <button
            onClick={() =>
              login.mutate({ email: 'test@example.com', password: 'wrong' })
            }
          >
            Login
          </button>
        )
      }

      renderWithProviders(<TestComponent />)

      // Trigger login
      await screen.getByRole('button').click()

      // Verify user is still unauthenticated
      await waitFor(() => {
        expect(useAuthStore.getState().isAuthenticated).toBe(false)
        expect(useAuthStore.getState().user).toBeNull()
      })
    })
  })

  describe('useRegister Hook Integration', () => {
    it('should successfully register without auto-login', async () => {
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

      function TestComponent() {
        const register = useRegister()
        return (
          <button
            onClick={() =>
              register.mutate({
                username: 'newuser',
                email: 'new@example.com',
                password: 'password123',
                role: 'normal',
              })
            }
          >
            Register
          </button>
        )
      }

      renderWithProviders(<TestComponent />)

      // Trigger registration
      await screen.getByRole('button').click()

      // Verify user is NOT logged in after registration (registration doesn't auto-login)
      await waitFor(() => {
        expect(useAuthStore.getState().isAuthenticated).toBe(false)
      })
    })

    it('should handle registration errors', async () => {
      server.use(
        http.post('/api/v1/auth/register', () =>
          HttpResponse.json({ error: '邮箱已存在' }, { status: 409 })
        )
      )

      function TestComponent() {
        const register = useRegister()
        return (
          <button
            onClick={() =>
              register.mutate({
                username: 'testuser',
                email: 'existing@example.com',
                password: 'password123',
                role: 'normal',
              })
            }
          >
            Register
          </button>
        )
      }

      renderWithProviders(<TestComponent />)

      // Trigger registration
      await screen.getByRole('button').click()

      // Verify user is still unauthenticated
      await waitFor(() => {
        expect(useAuthStore.getState().isAuthenticated).toBe(false)
      })
    })
  })

  describe('LoginPage Component Rendering', () => {
    it('should render login page correctly', () => {
      renderWithProviders(<LoginPage />)

      expect(screen.getByText('AI 测试管理平台')).toBeInTheDocument()
      expect(screen.getByText('账号登录')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('请输入邮箱')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('请输入密码')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: '登录' })).toBeInTheDocument()
    })

    it('should redirect if already authenticated', () => {
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

      renderWithProviders(<LoginPage />)

      // Login form should not be visible when authenticated
      expect(screen.queryByText('账号登录')).not.toBeInTheDocument()
    })
  })

  describe('RegisterPage Component Rendering', () => {
    it('should render register page correctly', () => {
      renderWithProviders(<RegisterPage />)

      expect(screen.getByRole('heading', { name: '注册' })).toBeInTheDocument()
      expect(screen.getByText('创建您的账号')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('请输入用户名')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('请输入邮箱')).toBeInTheDocument()
      expect(
        screen.getByPlaceholderText('请输入密码（至少 8 位字符）')
      ).toBeInTheDocument()
      expect(screen.getByRole('button', { name: '注册' })).toBeInTheDocument()
    })

    it('should have login link on register page', () => {
      renderWithProviders(<RegisterPage />)

      const loginLink = screen.getByText('立即登录')
      expect(loginLink).toBeInTheDocument()
      expect(loginLink.closest('a')).toHaveAttribute('href', '/login')
    })
  })
})
