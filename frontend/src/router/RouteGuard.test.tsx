import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import { RouteGuard } from './RouteGuard'
import { useAuthStore } from '../features/auth/hooks/useAuthStore'

describe('RouteGuard', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
    })
    vi.clearAllMocks()
    useAuthStore.getState().reset()
  })

  function renderWithRouter(
    ui: React.ReactElement,
    initialEntries = ['/']
  ) {
    return render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={initialEntries}>
          <Routes>
            <Route path="/login" element={<div>Login Page</div>} />
            <Route path="/projects" element={<div>Projects Page</div>} />
            <Route path="/admin" element={<div>Admin Page</div>} />
            <Route element={<RouteGuard requireAdmin={false} />}>
              <Route path="/" element={ui} />
            </Route>
            <Route
              path="/protected-admin"
              element={<RouteGuard requireAdmin={true}>{ui}</RouteGuard>}
            />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    )
  }

  describe('authentication check', () => {
    it('should redirect to /login when not authenticated', () => {
      renderWithRouter(<div>Protected Content</div>)

      expect(screen.getByText('Login Page')).toBeInTheDocument()
      expect(screen.queryByText('Protected Content')).not.toBeInTheDocument()
    })

    it('should render children when authenticated', () => {
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

      renderWithRouter(<div>Protected Content</div>)

      expect(screen.getByText('Protected Content')).toBeInTheDocument()
      expect(screen.queryByText('Login Page')).not.toBeInTheDocument()
    })

    it('should redirect to /login when token is expired', () => {
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

      // Mock JWT decode to return expired timestamp
      vi.spyOn(JSON, 'parse').mockReturnValue({
        exp: Math.floor(Date.now() / 1000) - 3600, // 1 hour ago
      })

      renderWithRouter(<div>Protected Content</div>)

      expect(screen.getByText('Login Page')).toBeInTheDocument()
    })
  })

  describe('admin check', () => {
    it('should redirect to /projects when requireAdmin=true and user is normal', () => {
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

      renderWithRouter(<div>Admin Content</div>, ['/protected-admin'])

      expect(screen.getByText('Projects Page')).toBeInTheDocument()
      expect(screen.queryByText('Admin Content')).not.toBeInTheDocument()
    })

    it('should render admin content when user is admin', () => {
      useAuthStore.setState({
        user: {
          id: 'admin-123',
          username: 'admin',
          email: 'admin@example.com',
          role: 'admin',
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        token: 'valid-token',
        refreshToken: 'valid-refresh',
        isAuthenticated: true,
      })

      renderWithRouter(<div>Admin Content</div>, ['/protected-admin'])

      expect(screen.getByText('Admin Content')).toBeInTheDocument()
    })
  })
})
