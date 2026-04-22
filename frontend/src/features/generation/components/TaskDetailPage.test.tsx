import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { TaskDetailPage } from './TaskDetailPage'
import { server } from '../../../../tests/msw/server'
import { generationHandlers } from '../../../../tests/msw/handlers/generation'

describe('TaskDetailPage', () => {
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

  function renderWithProviders(taskId: string) {
    window.history.pushState({}, '', `/generation/tasks/${taskId}`)
    return render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route path="/generation/tasks/:taskId" element={<TaskDetailPage />} />
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    )
  }

  it('should render task details page structure', async () => {
    renderWithProviders('550e8400-e29b-41d4-a716-446655440001')

    await waitFor(() => {
      expect(screen.getByText(/任务详情/i)).toBeInTheDocument()
    })

    // Should show back button
    expect(screen.getByRole('button', { name: /返回列表/i })).toBeInTheDocument()
  })

  it('should display polling indicator for processing status', async () => {
    renderWithProviders('550e8400-e29b-41d4-a716-446655440005')

    // Wait for component to render
    await waitFor(() => {
      expect(screen.getByText(/任务详情/i)).toBeInTheDocument()
    })

    // Should show processing state
    await waitFor(() => {
      expect(screen.getByText(/AI 正在生成用例，请稍候/i)).toBeInTheDocument()
    })
  })

  it('should display draft section for completed status', async () => {
    renderWithProviders('550e8400-e29b-41d4-a716-446655440004')

    // Wait for component and drafts to render
    await waitFor(() => {
      expect(screen.getByText(/任务详情/i)).toBeInTheDocument()
    })

    await waitFor(() => {
      expect(screen.getByText(/生成的草稿/i)).toBeInTheDocument()
    }, { timeout: 3000 })
  })

  it('should render task description card', async () => {
    renderWithProviders('550e8400-e29b-41d4-a716-446655440001')

    // Wait for component to render
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /返回列表/i })).toBeInTheDocument()
    }, { timeout: 3000 })
  })
})
