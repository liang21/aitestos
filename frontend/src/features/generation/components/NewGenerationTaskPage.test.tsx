import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { NewGenerationTaskPage } from './NewGenerationTaskPage'
import { server } from '../../../../tests/msw/server'
import { generationHandlers } from '../../../../tests/msw/handlers/generation'
import { http, HttpResponse } from 'msw'
import type { Module, Document, PaginatedResponse } from '@/types/api'

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
        HttpResponse.json<PaginatedResponse<Module>>({
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
    // Mock documents list for knowledge readiness check (default: sufficient)
    server.use(
      http.get('/api/v1/knowledge/documents', () =>
        HttpResponse.json<PaginatedResponse<Document>>({
          data: [
            { id: 'doc-1', name: 'PRD', type: 'prd', status: 'completed' },
            { id: 'doc-2', name: 'API Spec', type: 'api_spec', status: 'completed' },
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

    // Check for Collapse component by its header text
    // Arco Collapse renders the header text, so we can find it directly
    const collapseHeader = screen.queryByText(/高级选项.*场景类型/)
    expect(collapseHeader).toBeInTheDocument()

    // Also verify that the Collapse component contains the advanced options
    expect(screen.getByText(/场景类型/i)).toBeInTheDocument()
    expect(screen.getByText(/优先级/i)).toBeInTheDocument()
    expect(screen.getByText(/用例类型/i)).toBeInTheDocument()
  })

  it('should display submit button', async () => {
    renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: /立即生成/i })
      ).toBeInTheDocument()
    })
  })

  describe('Knowledge Readiness Indicator', () => {
    it('should show sufficient status (green) when 2+ completed documents exist', async () => {
      // Mock documents API with 2+ completed documents
      server.use(
        http.get('/api/v1/knowledge/documents', () =>
          HttpResponse.json<PaginatedResponse<Document>>({
            data: [
              { id: 'doc-1', name: 'PRD', type: 'prd', status: 'completed' },
              { id: 'doc-2', name: 'API Spec', type: 'api_spec', status: 'completed' },
            ],
            total: 2,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

      await waitFor(() => {
        expect(screen.getByText(/🟢 就绪/i)).toBeInTheDocument()
      })

      // Submit button should be enabled
      const submitButton = screen.getByRole('button', { name: /立即生成/i })
      expect(submitButton).not.toBeDisabled()
    })

    it('should show insufficient status (yellow) with warning when <2 documents exist', async () => {
      // Mock documents API with 1 completed document
      server.use(
        http.get('/api/v1/knowledge/documents', () =>
          HttpResponse.json<PaginatedResponse<Document>>({
            data: [{ id: 'doc-1', name: 'PRD', type: 'prd', status: 'completed' }],
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

      await waitFor(() => {
        expect(screen.getByText(/🟡 内容有限/i)).toBeInTheDocument()
      })

      // Should show warning alert
      expect(
        screen.getByText(/知识库内容较少，生成质量可能较低/i)
      ).toBeInTheDocument()

      // Submit button should still be enabled
      const submitButton = screen.getByRole('button', { name: /立即生成/i })
      expect(submitButton).not.toBeDisabled()
    })

    it('should show empty status (red) and disable submit when no documents exist', async () => {
      // Mock documents API with no documents
      server.use(
        http.get('/api/v1/knowledge/documents', () =>
          HttpResponse.json<PaginatedResponse<Document>>({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<NewGenerationTaskPage projectId="proj-1" />)

      await waitFor(() => {
        expect(screen.getByText(/🔴 请先上传需求文档/i)).toBeInTheDocument()
      })

      // Should show error alert
      expect(
        screen.getByText(/暂无需求文档，请先上传 PRD/i)
      ).toBeInTheDocument()

      // Submit button should be disabled
      const submitButton = screen.getByRole('button', { name: /立即生成/i })
      expect(submitButton).toBeDisabled()

      // Should show helper text
      expect(
        screen.getByText(/请先上传需求文档后再创建生成任务/i)
      ).toBeInTheDocument()
    })
  })
})
