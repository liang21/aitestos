import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { LoginPage } from './LoginPage'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'

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
        <MemoryRouter initialEntries={['/login']}>{ui}</MemoryRouter>
      </QueryClientProvider>
    )
  }

  describe('rendering', () => {
    it('should render banner with slogan', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.getByAltText('Aitestos Platform')).toBeInTheDocument()
      expect(screen.getByText('因为热爱 快乐成长')).toBeInTheDocument()
    })

    it('should render form with branding', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.getByText('AI 测试管理平台')).toBeInTheDocument()
      expect(screen.getByText('账号登录')).toBeInTheDocument()
    })

    it('should render email input field', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.getByPlaceholderText('请输入邮箱')).toBeInTheDocument()
    })

    it('should render password input field', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.getByPlaceholderText('请输入密码')).toBeInTheDocument()
    })

    it('should render login button', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.getByRole('button', { name: '登录' })).toBeInTheDocument()
    })

    it('should render register link', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.getByText('还没有账号？')).toBeInTheDocument()
      expect(screen.getByRole('link', { name: '立即注册' })).toBeInTheDocument()
    })

    it('should not render "remember me" checkbox', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.queryByText('记住我')).not.toBeInTheDocument()
    })

    it('should not render "forgot password" link', () => {
      renderWithRouter(<LoginPage />)

      expect(screen.queryByText('忘记密码？')).not.toBeInTheDocument()
    })
  })

  describe('authentication check', () => {
    it('should redirect to /projects if already authenticated', () => {
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

      expect(screen.queryByText('账号登录')).not.toBeInTheDocument()
    })
  })
})
