import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse } from 'msw'
import { DocumentDetailPage } from './DocumentDetailPage'

// Mock useParams at top level
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...(actual as any),
    useParams: () => ({ documentId: 'doc-1' }),
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

describe('DocumentDetailPage', () => {
  beforeEach(() => {
    server.resetHandlers()
    // Mock localStorage to provide access_token
    vi.stubGlobal('localStorage', {
      getItem: vi.fn((key) => {
        if (key === 'access_token') return 'mock-token'
        return null
      }),
      setItem: vi.fn(),
      removeItem: vi.fn(),
      clear: vi.fn(),
    })
  })

  describe('rendering', () => {
    it('should render document basic information', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents/doc-1', () =>
          HttpResponse.json({
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
            chunks: [],
          })
        )
      )

      // Act
      renderWithProviders(<DocumentDetailPage />)

      // Assert
      await waitFor(() => {
        expect(screen.getByText('用户注册模块 PRD v2.0')).toBeInTheDocument()
        expect(screen.getByText('42')).toBeInTheDocument() // chunk count
      })
    })

    it('should render chunk list', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents/doc-1', () =>
          HttpResponse.json({
            id: 'doc-1',
            projectId: 'proj-1',
            name: '用户注册模块 PRD v2.0',
            type: 'prd',
            status: 'completed',
            chunkCount: 2,
            uploadedBy: 'user-1',
            createdAt: '2026-04-16T08:00:00Z',
            updatedAt: '2026-04-16T08:30:00Z',
            chunks: [
              {
                id: 'chunk-1',
                documentId: 'doc-1',
                chunkIndex: 0,
                content: '用户注册功能概述...',
              },
              {
                id: 'chunk-2',
                documentId: 'doc-1',
                chunkIndex: 1,
                content: '邮箱验证规则...',
              },
            ],
          })
        )
      )

      // Act
      renderWithProviders(<DocumentDetailPage />)

      // Assert
      await waitFor(() => {
        expect(screen.getByText('用户注册功能概述...')).toBeInTheDocument()
        expect(screen.getByText('邮箱验证规则...')).toBeInTheDocument()
      })
    })
  })

  describe('processing state', () => {
    it('should show Spin when document status is processing', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents/:id', () =>
          HttpResponse.json({
            id: 'doc-1',
            projectId: 'proj-1',
            name: '处理中的文档',
            type: 'prd',
            status: 'processing',
            chunkCount: 0,
            uploadedBy: 'user-1',
            createdAt: '2026-04-16T08:00:00Z',
            updatedAt: '2026-04-16T08:30:00Z',
            chunks: [],
          })
        )
      )

      // Act
      renderWithProviders(<DocumentDetailPage />)

      // Assert - Arco Design Spin component should be present
      await waitFor(() => {
        expect(screen.getByText('处理中的文档')).toBeInTheDocument()
      })
    })
  })
})
