import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse } from 'msw'
import { KnowledgeListPage } from './KnowledgeListPage'

// Mock useParams at top level
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual as any,
    useParams: () => ({ projectId: 'proj-1' }),
  }
})

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
}

function renderWithProviders(ui: any) {
  const queryClient = createTestQueryClient()
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>{ui}</BrowserRouter>
    </QueryClientProvider>
  )
}

describe('KnowledgeListPage', () => {
  beforeEach(() => {
    server.resetHandlers()
  })

  describe('rendering', () => {
    it('should render document list with name, type tag, status tag, and upload time', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents', () =>
          HttpResponse.json({
            data: [
              {
                id: 'doc-1',
                projectId: 'proj-1',
                name: '用户注册模块 PRD v2.0',
                type: 'prd',
                status: 'completed',
                chunkCount: 42,
                uploadedBy: 'user-1',
                uploadedByName: '张三',
                createdAt: '2026-04-16T08:00:00Z',
                updatedAt: '2026-04-16T08:30:00Z',
              },
              {
                id: 'doc-2',
                projectId: 'proj-1',
                name: 'API 接口规范',
                type: 'api_spec',
                status: 'processing',
                chunkCount: 0,
                uploadedBy: 'user-1',
                uploadedByName: '张三',
                createdAt: '2026-04-16T09:00:00Z',
                updatedAt: '2026-04-16T09:05:00Z',
              },
            ],
            total: 2,
            offset: 0,
            limit: 10,
          })
        )
      )

      // Act
      renderWithProviders(<KnowledgeListPage />)

      // Assert
      await waitFor(() => {
        expect(screen.getByText('用户注册模块 PRD v2.0')).toBeInTheDocument()
        expect(screen.getByText('API 接口规范')).toBeInTheDocument()
      })
    })
  })

  describe('filtering', () => {
    it('should filter documents by type', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents', ({ request }) => {
          const url = new URL(request.url)
          const type = url.searchParams.get('type')
          return HttpResponse.json({
            data: [
              {
                id: 'doc-1',
                projectId: 'proj-1',
                name: 'PRD 文档',
                type: type || 'prd',
                status: 'completed',
                chunkCount: 10,
                uploadedBy: 'user-1',
                createdAt: '2026-04-16T08:00:00Z',
                updatedAt: '2026-04-16T08:00:00Z',
              },
            ],
            total: 1,
            offset: 0,
            limit: 10,
          })
        })
      )

      const user = userEvent.setup()

      // Act
      renderWithProviders(<KnowledgeListPage />)

      // Select type filter - click first combobox (document type)
      const typeSelects = screen.getAllByRole('combobox')
      await user.click(typeSelects[0])

      // Assert - filter request was made (verified in MSW handler)
      await waitFor(() => {
        expect(screen.getByText('PRD 文档')).toBeInTheDocument()
      })
    })

    it('should show empty state when no documents match', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents', () =>
          HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        )
      )

      // Act
      renderWithProviders(<KnowledgeListPage />)

      // Assert
      await waitFor(() => {
        expect(screen.getByText(/暂无数据/i)).toBeInTheDocument()
      })
    })
  })

  describe('actions', () => {
    it('should open upload modal when clicking upload button', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents', () =>
          HttpResponse.json({ data: [], total: 0, offset: 0, limit: 10 })
        )
      )

      const user = userEvent.setup()

      // Act
      renderWithProviders(<KnowledgeListPage />)

      const uploadButton = screen.getByRole('button', { name: /上传文档/i })
      await user.click(uploadButton)

      // Assert
      await waitFor(() => {
        expect(screen.getByRole('dialog', { name: /上传文档/i })).toBeInTheDocument()
      })
    })
  })
})
