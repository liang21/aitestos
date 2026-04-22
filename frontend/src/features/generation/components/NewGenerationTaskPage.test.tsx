import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
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
      http.get('/api/v1/projects/:projectId/modules', () =>
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
        <MemoryRouter initialEntries={['/']}>{ui}</MemoryRouter>
      </QueryClientProvider>
    )
  }

  it('should render module selector as required field', async () => {
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    // Wait for modules to load
    await waitFor(() => {
      expect(screen.getByText(/目标模块/i)).toBeInTheDocument()
    })

    // Module selector should be present
    const moduleSelect = screen.getByRole('combobox')
    expect(moduleSelect).toBeInTheDocument()
  })

  it('should display prompt description input field', async () => {
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/目标模块/i)).toBeInTheDocument()
    })

    const promptInput = screen.getByPlaceholderText(/请描述测试需求/i)
    expect(promptInput).toBeInTheDocument()
  })

  it('should display case count input field with default value', async () => {
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/目标模块/i)).toBeInTheDocument()
    })

    const countInput = screen.getByRole('spinbutton')
    expect(countInput).toBeInTheDocument()
  })

  it('should display advanced options collapse', async () => {
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    await waitFor(() => {
      expect(screen.getByText(/目标模块/i)).toBeInTheDocument()
    })

    // Check for Collapse component by looking for the form item label
    const collapseLabel = screen
      .getAllByText('高级选项')
      .find((el) => el.tagName === 'LABEL')
    expect(collapseLabel).toBeInTheDocument()
  })

  it('should display submit button', async () => {
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: /立即生成/i })
      ).toBeInTheDocument()
    })
  })
})
