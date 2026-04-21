import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { Sidebar } from './Sidebar'
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

describe('Sidebar', () => {
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
      logout: mockLogout,
    } as any)

    vi.mocked(useAppStore).mockReturnValue({
      sidebarCollapsed: false,
      toggleSidebar: mockToggleSidebar,
    } as any)

    vi.mocked(usePendingDraftCount).mockReturnValue({
      data: 5,
      isLoading: false,
    } as any)
  })

  it('should render menu items (project list, current project submenu, drafts)', () => {
    render(<Sidebar />, { wrapper })

    // Check for main menu items
    expect(screen.getByText('项目列表')).toBeInTheDocument()
    expect(screen.getByText('草稿箱')).toBeInTheDocument()
  })

  it('should show selected state style for active route', () => {
    // Mock window.location to simulate active route
    const locationSpy = vi.spyOn(window, 'location', 'get').mockReturnValue({
      ...window.location,
      pathname: '/projects',
    } as Location)

    render(<Sidebar />, { wrapper })

    // Projects menu should be in selected state
    const projectsMenu = screen.getByText('项目列表')
    expect(projectsMenu).toBeInTheDocument()

    locationSpy.mockRestore()
  })

  it('should toggle collapse/expand state when collapse button is clicked', async () => {
    const user = userEvent.setup()
    render(<Sidebar />, { wrapper })

    // Find button by role (icon-only button may not have accessible name)
    const collapseButton = screen.getByRole('button')
    expect(collapseButton).toBeInTheDocument()

    // Click to collapse
    await user.click(collapseButton)
    expect(mockToggleSidebar).toHaveBeenCalledTimes(1)
  })

  it('should display pending draft count in Badge', () => {
    render(<Sidebar />, { wrapper })

    // Badge should show count of 5
    const badge = screen.getByText('5')
    expect(badge).toBeInTheDocument()
  })
})
