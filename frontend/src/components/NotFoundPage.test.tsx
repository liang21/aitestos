import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { NotFoundPage } from './NotFoundPage'

// Mock useNavigate
const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  }
})

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

describe('NotFoundPage', () => {
  beforeEach(() => {
    mockNavigate.mockClear()
  })

  it('should render 404 message in Chinese', () => {
    render(<NotFoundPage />, { wrapper })

    // Check for 404 status (may be in SVG or other element)
    const container = screen.getByText('抱歉，您访问的页面不存在。')
    expect(container).toBeInTheDocument()

    // Verify Result component is rendered
    expect(document.querySelector('.arco-result-is-404')).toBeInTheDocument()
  })

  it('should navigate to /projects on button click', async () => {
    const user = userEvent.setup()
    render(<NotFoundPage />, { wrapper })

    const backButton = screen.getByRole('button', { name: /返回首页/ })
    expect(backButton).toBeInTheDocument()

    await user.click(backButton)

    expect(mockNavigate).toHaveBeenCalledTimes(1)
    expect(mockNavigate).toHaveBeenCalledWith('/projects')
  })
})
