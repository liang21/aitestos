import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { NewGenerationTaskPage } from './NewGenerationTaskPage'
import { server } from '../../../../tests/msw/server'
import { generationHandlers } from '../../../../tests/msw/handlers/generation'
import { http, HttpResponse } from 'msw'

describe('NewGenerationTaskPage', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
    server.use(...generationHandlers)
    // Mock modules list for module selection
    server.use(
      http.get('/api/v1/projects/:id/modules', () =>
        HttpResponse.json({
          data: [
            { id: 'mod-1', name: '用户中心', abbreviation: 'USR' },
            { id: 'mod-2', name: '订单管理', abbreviation: 'ORD' },
          ],
          total: 2,
          offset: 0,
          limit: 10,
        })
      )
    )
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

  it('should render module selector as required field', async () => {
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/模块/i)).toBeInTheDocument()
    })

    // Module selector should be present
    const moduleSelect = screen.getByRole('combobox', { name: /模块/i })
    expect(moduleSelect).toBeInTheDocument()
    expect(moduleSelect).toBeRequired()
  })

  it('should validate prompt description (min 10 characters)', async () => {
    const user = userEvent.setup()
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    const promptInput = screen.getByPlaceholderText(/请描述测试需求/i)

    // Type less than 10 characters
    await user.type(promptInput, 'short')

    const submitButton = screen.getByRole('button', { name: /立即生成/i })
    await user.click(submitButton)

    // Should show validation error
    await waitFor(() => {
      expect(screen.getByText(/需求描述至少10个字/i)).toBeInTheDocument()
    })
  })

  it('should validate case count (1-20 range)', async () => {
    const user = userEvent.setup()
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    const countInput = screen.getByRole('spinbutton', { name: /用例数量/i })

    // Test upper bound
    await user.clear(countInput)
    await user.type(countInput, '25')

    const submitButton = screen.getByRole('button', { name: /立即生成/i })
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/用例数量范围为1-20/i)).toBeInTheDocument()
    })
  })

  it('should toggle advanced options collapse', async () => {
    const user = userEvent.setup()
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    const collapseButton = screen.getByText(/高级选项/i)

    // Initially collapsed
    expect(screen.queryByText(/场景类型/i)).not.toBeInTheDocument()

    // Click to expand
    await user.click(collapseButton)

    await waitFor(() => {
      expect(screen.getByText(/场景类型/i)).toBeInTheDocument()
    })
  })

  it('should navigate to task detail on successful submission', async () => {
    const user = userEvent.setup()
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    // Fill required fields
    await user.selectOptions(screen.getByRole('combobox', { name: /模块/i }), 'mod-1')
    await user.type(screen.getByPlaceholderText(/请描述测试需求/i), '测试用户注册功能，包括邮箱验证和密码强度校验')
    await user.type(screen.getByRole('spinbutton', { name: /用例数量/i }), '5')

    const submitButton = screen.getByRole('button', { name: /立即生成/i })
    await user.click(submitButton)

    await waitFor(() => {
      expect(submitButton).toBeDisabled()
    })

    // Should navigate to task detail page
    await waitFor(() => {
      expect(window.location.pathname).toMatch(/\/generation\/tasks\/[^/]+$/)
    })
  })
})
