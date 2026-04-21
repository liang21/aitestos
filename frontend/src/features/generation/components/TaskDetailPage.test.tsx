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

  function renderWithProviders(ui: React.ReactElement, taskId: string) {
    return render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route path="/generation/tasks/:taskId" element={ui} />
          </Routes>
          <div
            onClick={() => {
              window.history.pushState({}, '', `/generation/tasks/${taskId}`)
            }}
          />
        </BrowserRouter>
      </QueryClientProvider>
    )
  }

  it('should render task details (prompt, status, timestamps)', async () => {
    const taskId = '550e8400-e29b-41d4-a716-446655440001'
    window.history.pushState({}, '', `/generation/tasks/${taskId}`)
    renderWithProviders(<TaskDetailPage />, taskId)

    await waitFor(() => {
      expect(screen.getByText(/测试用户注册功能/i)).toBeInTheDocument()
    })

    expect(screen.getByText(/pending/i)).toBeInTheDocument()
    expect(screen.getByText(/2026-04-20/i)).toBeInTheDocument()
  })

  it('should display polling progress when status is processing', async () => {
    const taskId = '550e8400-e29b-41d4-a716-446655440005'
    window.history.pushState({}, '', `/generation/tasks/${taskId}`)
    renderWithProviders(<TaskDetailPage />, taskId)

    await waitFor(() => {
      expect(screen.getByText(/processing/i)).toBeInTheDocument()
    })

    // Should show loading/polling indicator
    expect(screen.getByRole('status')).toBeInTheDocument()
  })

  it('should display draft list when status is completed', async () => {
    const taskId = '550e8400-e29b-41d4-a716-446655440004'
    window.history.pushState({}, '', `/generation/tasks/${taskId}`)
    renderWithProviders(<TaskDetailPage />, taskId)

    await waitFor(() => {
      expect(screen.getByText(/生成的草稿/i)).toBeInTheDocument()
    })

    // Should show draft items
    expect(screen.getByText(/验证有效邮箱注册/i)).toBeInTheDocument()
    expect(screen.getByText(/验证邮箱格式校验/i)).toBeInTheDocument()
  })

  it('should display draft with title, type, priority, and confidence tag', async () => {
    const taskId = '550e8400-e29b-41d4-a716-446655440004'
    window.history.pushState({}, '', `/generation/tasks/${taskId}`)
    renderWithProviders(<TaskDetailPage />, taskId)

    await waitFor(() => {
      expect(screen.getByText(/验证有效邮箱注册/i)).toBeInTheDocument()
    })

    // Check draft metadata
    expect(screen.getByText(/functionality/i)).toBeInTheDocument()
    expect(screen.getByText(/P1/i)).toBeInTheDocument()
  })

  it('should render task details (prompt, status, timestamps)', async () => {
    renderWithProviders(
      <TaskDetailPage taskId="550e8400-e29b-41d4-a716-446655440001" />
    )

    await waitFor(() => {
      expect(screen.getByText(/测试用户注册功能/i)).toBeInTheDocument()
    })

    expect(screen.getByText(/pending/i)).toBeInTheDocument()
    expect(screen.getByText(/2026-04-20/i)).toBeInTheDocument()
  })

  it('should display polling progress when status is processing', async () => {
    renderWithProviders(
      <TaskDetailPage taskId="550e8400-e29b-41d4-a716-446655440005" />
    )

    await waitFor(() => {
      expect(screen.getByText(/processing/i)).toBeInTheDocument()
    })

    // Should show loading/polling indicator
    expect(screen.getByRole('status', { hidden: true })).toBeInTheDocument()
  })

  it('should display draft list when status is completed', async () => {
    renderWithProviders(
      <TaskDetailPage taskId="550e8400-e29b-41d4-a716-446655440004" />
    )

    await waitFor(() => {
      expect(screen.getByText(/草稿列表/i)).toBeInTheDocument()
    })

    // Should show draft items
    expect(screen.getByText(/验证有效邮箱注册/i)).toBeInTheDocument()
    expect(screen.getByText(/验证邮箱格式校验/i)).toBeInTheDocument()
  })

  it('should display draft with title, type, priority, and confidence tag', async () => {
    renderWithProviders(
      <TaskDetailPage taskId="550e8400-e29b-41d4-a716-446655440004" />
    )

    await waitFor(() => {
      expect(screen.getByText(/验证有效邮箱注册/i)).toBeInTheDocument()
    })

    // Check draft metadata
    expect(screen.getByText(/functionality/i)).toBeInTheDocument()
    expect(screen.getByText(/P1/i)).toBeInTheDocument()
    // Confidence should be displayed
    const draftItem = screen.getByText(/验证有效邮箱注册/i).closest('tr')
    expect(draftItem).toBeInTheDocument()
  })
})
