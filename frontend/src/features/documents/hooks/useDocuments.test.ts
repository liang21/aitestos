import { describe, it, expect, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse } from 'msw'
import { useDocumentList, useDocumentDetail, useUploadDocument, useDeleteDocument, useDocumentChunks } from './useDocuments'
import { createElement } from 'react'

// Helper: 创建测试用 QueryClient
function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
}

function wrapper({ children }: { children: any }) {
  return createElement(QueryClientProvider, { client: createTestQueryClient() }, children)
}

describe('useDocuments hooks', () => {
  beforeEach(() => {
    server.resetHandlers()
  })

  describe('useDocumentList', () => {
    it('should return document list data', async () => {
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
            ],
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      // Act
      const { result } = renderHook(() => useDocumentList({ projectId: 'proj-1' }), { wrapper })

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.data).toHaveLength(1)
      expect(result.current.data?.data[0].name).toBe('用户注册模块 PRD v2.0')
    })

    it('should support filtering by type and status', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents', ({ request }) => {
          const url = new URL(request.url)
          expect(url.searchParams.get('type')).toBe('api_spec')
          expect(url.searchParams.get('status')).toBe('processing')
          return HttpResponse.json({ data: [], total: 0, offset: 0, limit: 10 })
        })
      )

      // Act
      const { result } = renderHook(
        () => useDocumentList({ projectId: 'proj-1', type: 'api_spec', status: 'processing' }),
        { wrapper }
      )

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
    })
  })

  describe('useDocumentDetail', () => {
    it('should fetch document detail by id', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents/:id', () =>
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
      const { result } = renderHook(() => useDocumentDetail('doc-1'), { wrapper })

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.name).toBe('用户注册模块 PRD v2.0')
      expect(result.current.data?.chunkCount).toBe(42)
    })

    it('should not fetch when id is empty', async () => {
      // Act
      const { result } = renderHook(() => useDocumentDetail(''), { wrapper })

      // Assert
      expect(result.current.fetchStatus).toBe('idle')
    })

    it('should handle fetch error', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents/:id', () =>
          HttpResponse.json({ error: 'Not found' }, { status: 404 })
        )
      )

      // Act
      const { result } = renderHook(() => useDocumentDetail('doc-1'), { wrapper })

      // Assert
      await waitFor(() => expect(result.current.isError).toBe(true))
    })
  })

  describe('useUploadDocument', () => {
    it('should upload document successfully', async () => {
      // Arrange
      const mockDocument = {
        id: 'doc-new',
        projectId: 'proj-1',
        name: '新文档.pdf',
        type: 'prd' as const,
        status: 'pending' as const,
        chunkCount: 0,
        uploadedBy: 'user-1',
        uploadedByName: '张三',
        createdAt: '2026-04-16T10:00:00Z',
        updatedAt: '2026-04-16T10:00:00Z',
      }

      server.use(
        http.post('/api/v1/knowledge/documents', () =>
          HttpResponse.json(mockDocument, { status: 201 })
        )
      )

      // Act
      const { result } = renderHook(() => useUploadDocument(), { wrapper })

      const uploadData = {
        projectId: 'proj-1',
        name: '新文档.pdf',
        type: 'prd' as const,
        file: new File(['test'], 'test.pdf'),
      }

      result.current.mutate(uploadData)

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
    })

    it('should handle upload error', async () => {
      // Arrange
      server.use(
        http.post('/api/v1/knowledge/documents', () =>
          HttpResponse.json({ error: 'Upload failed' }, { status: 500 })
        )
      )

      // Act
      const { result } = renderHook(() => useUploadDocument(), { wrapper })

      const uploadData = {
        projectId: 'proj-1',
        name: '新文档.pdf',
        type: 'prd' as const,
        file: new File(['test'], 'test.pdf'),
      }

      result.current.mutate(uploadData)

      // Assert
      await waitFor(() => expect(result.current.isError).toBe(true))
    })
  })

  describe('useDeleteDocument', () => {
    it('should delete document successfully', async () => {
      // Arrange
      server.use(
        http.delete('/api/v1/knowledge/documents/:id', () =>
          new HttpResponse(null, { status: 204 })
        )
      )

      // Act
      const { result } = renderHook(() => useDeleteDocument(), { wrapper })

      result.current.mutate('doc-1')

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
    })

    it('should handle delete error', async () => {
      // Arrange
      server.use(
        http.delete('/api/v1/knowledge/documents/:id', () =>
          HttpResponse.json({ error: 'Delete failed' }, { status: 500 })
        )
      )

      // Act
      const { result } = renderHook(() => useDeleteDocument(), { wrapper })

      result.current.mutate('doc-1')

      // Assert
      await waitFor(() => expect(result.current.isError).toBe(true))
    })
  })

  describe('useDocumentChunks', () => {
    it('should fetch document chunks', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents/:id/chunks', () =>
          HttpResponse.json([
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
          ])
        )
      )

      // Act
      const { result } = renderHook(() => useDocumentChunks('doc-1'), { wrapper })

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data).toHaveLength(2)
      expect(result.current.data?.[0].chunkIndex).toBe(0)
    })
  })
})
