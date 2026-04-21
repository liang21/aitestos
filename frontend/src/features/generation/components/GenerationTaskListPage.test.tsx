import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { GenerationTaskListPage } from './GenerationTaskListPage'
import { server } from '../../../../tests/msw/server'
import { generationHandlers } from '../../../../tests/msw/handlers/generation'

describe('GenerationTaskListPage', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
    server.use(...generationHandlers)
  })

  afterEach(() => {
    server.resetHandlers()
  })

  function renderWithProviders(ui: React.ReactElement) {
    return render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>{ui}</BrowserRouter>
      </QueryClientProvider>
    )
  }

  it('should render task list with prompt summary, status tag, and created time', async () => {
    renderWithProviders(<GenerationTaskListPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/测试用户注册功能/i)).toBeInTheDocument()
    })

    // Should show status tags
    expect(screen.getByText(/pending/i)).toBeInTheDocument()
    expect(screen.getByText(/completed/i)).toBeInTheDocument()

    // Should show created time
    expect(screen.getByText(/2026-04-20/i)).toBeInTheDocument()
  })

  it('should filter tasks by status', async () => {
    const user = userEvent.setup()
    renderWithProviders(<GenerationTaskListPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/测试用户注册功能/i)).toBeInTheDocument()
    })

    // Select status filter
    const statusSelect = screen.getByRole('combobox', { name: /状态/i })
    await user.selectOptions(statusSelect, 'completed')

    await waitFor(() => {
      // Should only show completed tasks
      const tasks = screen.getAllByText(/completed/i)
      expect(tasks.length).toBeGreaterThan(0)
    })
  })

  it('should navigate to task detail on click', async () => {
    const user = userEvent.setup()
    renderWithProviders(<GenerationTaskListPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/测试用户注册功能/i)).toBeInTheDocument()
    })

    const firstTask = screen.getByText(/测试用户注册功能/i).closest('tr')
    expect(firstTask).toBeInTheDocument()

    await user.click(firstTask!)

    await waitFor(() => {
      expect(window.location.pathname).toMatch(/\/generation\/tasks\/[^/]+$/)
    })
  })
})
