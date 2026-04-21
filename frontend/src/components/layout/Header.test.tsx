import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { Header } from './Header'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'
import { useAppStore } from '@/store/useAppStore'
import { useLogout } from '@/features/auth/hooks/useAuth'

// Mock stores
vi.mock('@/features/auth/hooks/useAuthStore', () => ({
  useAuthStore: vi.fn(),
}))

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}))

vi.mock('@/features/auth/hooks/useAuth', () => ({
  useLogout: vi.fn(),
}))

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
}

function wrapper({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={createTestQueryClient()}>
      <BrowserRouter>{children}</BrowserRouter>
    </QueryClientProvider>
  )
}

describe('Header', () => {
  const mockToggleSidebar = vi.fn()
  const mockLogout = vi.fn()

  beforeEach(() => {
    // Reset mocks before each test
    vi.mocked(useAuthStore).mockReturnValue({
      user: {
        id: '1',
        username: 'admin',
        email: 'admin@test.com',
        role: 'admin',
      },
    } as any)

    vi.mocked(useAppStore).mockReturnValue({
      sidebarCollapsed: false,
      toggleSidebar: mockToggleSidebar,
    } as any)

    vi.mocked(useLogout).mockReturnValue(mockLogout)
  })

  it('should render collapse button', () => {
    render(<Header title="项目列表" />, { wrapper })

    const collapseButton = screen.getByRole('button', { name: '折叠侧边栏' })
    expect(collapseButton).toBeInTheDocument()
  })

  it('should render breadcrumb navigation', () => {
    const breadcrumbs = [{ title: '项目', path: '/projects' }]
    render(<Header title="项目列表" breadcrumbs={breadcrumbs} />, { wrapper })

    expect(screen.getByText('项目')).toBeInTheDocument()
    expect(screen.getByText('项目列表')).toBeInTheDocument()
  })

  it('should render user dropdown with username', () => {
    render(<Header title="项目列表" />, { wrapper })

    // Check that username is displayed
    expect(screen.getByText('admin')).toBeInTheDocument()
  })
})
