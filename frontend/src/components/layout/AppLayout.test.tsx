import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { AppLayout } from './AppLayout'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'
import { useAppStore } from '@/store/useAppStore'
import { usePendingDraftCount } from '@/features/drafts/hooks/useDrafts'

// Mock stores
vi.mock('@/features/auth/hooks/useAuthStore', () => ({
  useAuthStore: vi.fn(),
}))

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}))

vi.mock('@/features/drafts/hooks/useDrafts', () => ({
  usePendingDraftCount: vi.fn(),
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

describe('AppLayout', () => {
  const mockToggleSidebar = vi.fn()

  beforeEach(() => {
    // Reset mocks before each test
    vi.mocked(useAuthStore).mockReturnValue({
      isAuthenticated: true,
      user: {
        id: '1',
        username: 'admin',
        email: 'admin@test.com',
        role: 'admin',
      },
    } as any)

    vi.mocked(usePendingDraftCount).mockReturnValue({
      data: 0,
      isLoading: false,
    } as any)
  })

  it('should render Sidebar, Header, and Content three-section layout', () => {
    vi.mocked(useAppStore).mockReturnValue({
      sidebarCollapsed: false,
      toggleSidebar: mockToggleSidebar,
    } as any)

    render(
      <AppLayout title="项目列表">
        <div data-testid="content">Page Content</div>
      </AppLayout>,
      { wrapper }
    )

    expect(screen.getByTestId('app-sidebar')).toBeInTheDocument()
    expect(screen.getByTestId('app-header')).toBeInTheDocument()
    expect(screen.getByTestId('content')).toBeInTheDocument()
  })

  it('should apply collapsed class when sidebarCollapsed is true', () => {
    vi.mocked(useAppStore).mockReturnValue({
      sidebarCollapsed: true,
      toggleSidebar: mockToggleSidebar,
    } as any)

    render(
      <AppLayout title="项目列表">
        <div>Content</div>
      </AppLayout>,
      { wrapper }
    )

    const sidebar = screen.getByTestId('app-sidebar')
    expect(sidebar).toHaveClass('collapsed')
  })

  it('should not apply collapsed class when sidebarCollapsed is false', () => {
    vi.mocked(useAppStore).mockReturnValue({
      sidebarCollapsed: false,
      toggleSidebar: mockToggleSidebar,
    } as any)

    render(
      <AppLayout title="项目列表">
        <div>Content</div>
      </AppLayout>,
      { wrapper }
    )

    const sidebar = screen.getByTestId('app-sidebar')
    expect(sidebar).not.toHaveClass('collapsed')
  })
})
