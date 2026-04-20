import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { LoginPage } from './LoginPage'
import { useAuthStore } from '../hooks/useAuthStore'

describe('LoginPage', () => {
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
  })

  function renderWithRouter(ui: React.ReactElement) {
    return render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>{ui}</MemoryRouter>
      </QueryClientProvider>
    )
  }

  describe('rendering', () => {
    it('should render email/password input fields and login button', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.getByText('邮箱')).toBeInTheDocument()
      expect(screen.getByText('密码')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: '登录' })).toBeInTheDocument()
    })

    it('should show register link', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.getByText('还没有账号？')).toBeInTheDocument()
      expect(screen.getByRole('link', { name: '立即注册' })).toBeInTheDocument()
    })

    it('should have correct placeholders', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.getByPlaceholderText('请输入邮箱')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('请输入密码')).toBeInTheDocument()
    })
  })

  describe('authentication check', () => {
    it('should redirect to /projects if already authenticated', () => {
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

      renderWithRouter(<LoginPage />)

      // Component should return null (redirects internally)
      expect(screen.queryByText('登录')).not.toBeInTheDocument()
    })
  })
})
