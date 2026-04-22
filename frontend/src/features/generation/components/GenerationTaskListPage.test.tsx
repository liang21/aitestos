import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
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
        <MemoryRouter>{ui}</MemoryRouter>
      </QueryClientProvider>
    )
  }

  it('should render task list page structure', async () => {
    renderWithProviders(<GenerationTaskListPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/AI 生成任务/i)).toBeInTheDocument()
    })

    // Should show new task button
    expect(
      screen.getByRole('button', { name: /新建任务/i })
    ).toBeInTheDocument()
  })

  it('should render status filter dropdown', async () => {
    renderWithProviders(<GenerationTaskListPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/AI 生成任务/i)).toBeInTheDocument()
    })

    // Should have a select for filtering
    const filterSelect = screen.getByPlaceholderText(/筛选状态/i)
    expect(filterSelect).toBeInTheDocument()
  })

  it('should display table component', async () => {
    renderWithProviders(<GenerationTaskListPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/AI 生成任务/i)).toBeInTheDocument()
    })

    // Should have a table for displaying tasks
    expect(screen.getByRole('table')).toBeInTheDocument()
  })
})
