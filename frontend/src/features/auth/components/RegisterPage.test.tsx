import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { RegisterPage } from './RegisterPage'

describe('RegisterPage', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
    })
    vi.clearAllMocks()
  })

  function renderWithRouter(ui: React.ReactElement) {
    return render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
          {ui}
        </MemoryRouter>
      </QueryClientProvider>
    )
  }

  describe('rendering', () => {
    it('should render username/email/password/role form fields', () => {
      renderWithRouter(<RegisterPage />)

      expect(screen.getByText('用户名')).toBeInTheDocument()
      expect(screen.getByText('邮箱')).toBeInTheDocument()
      expect(screen.getByText('密码')).toBeInTheDocument()
      expect(screen.getByText('角色')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: '注册' })).toBeInTheDocument()
    })

    it('should show login link', () => {
      renderWithRouter(<RegisterPage />)

      expect(screen.getByText('已有账号？')).toBeInTheDocument()
      expect(screen.getByRole('link', { name: '立即登录' })).toBeInTheDocument()
    })

    it('should have correct placeholders', () => {
      renderWithRouter(<RegisterPage />)

      expect(screen.getByPlaceholderText('请输入用户名')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('请输入邮箱')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('请输入密码（至少 8 位字符）')).toBeInTheDocument()
    })

    it('should display role options', () => {
      renderWithRouter(<RegisterPage />)

      expect(screen.getByText('普通用户')).toBeInTheDocument()
      expect(screen.getByText('管理员')).toBeInTheDocument()
      expect(screen.getByText('超级管理员')).toBeInTheDocument()
    })
  })
})
