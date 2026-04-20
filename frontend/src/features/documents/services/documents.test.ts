import { describe, it, expect, beforeEach } from 'vitest'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse } from 'msw'
import { documentsApi } from './documents'
import type { Document, DocumentDetail, DocumentChunk } from '@/types/api'

describe('Documents API service', () => {
  beforeEach(() => {
    server.resetHandlers()
  })

  describe('list(projectId)', () => {
    it('should call GET /knowledge/documents and return document list', async () => {
      // Arrange
      const mockResponse = {
        data: [
          {
            id: 'doc-1',
            projectId: 'proj-1',
            name: '用户注册模块 PRD v2.0',
            type: 'prd' as const,
            status: 'completed' as const,
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
      }

      server.use(
        http.get('/api/v1/knowledge/documents', () =>
          HttpResponse.json(mockResponse)
        )
      )

      // Act
      const result = await documentsApi.list({ projectId: 'proj-1' })

      // Assert
      expect(result).toEqual(mockResponse)
      expect(result.data[0].id).toBe('doc-1')
      expect(result.data[0].type).toBe('prd')
    })

    it('should support type filter parameter', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents', ({ request }) => {
          const url = new URL(request.url)
          expect(url.searchParams.get('type')).toBe('api_spec')
          return HttpResponse.json({ data: [], total: 0, offset: 0, limit: 10 })
        })
      )

      // Act
      await documentsApi.list({ projectId: 'proj-1', type: 'api_spec' })

      // Assert - handled in MSW handler
      expect(true).toBe(true)
    })

    it('should support status filter parameter', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/knowledge/documents', ({ request }) => {
          const url = new URL(request.url)
          expect(url.searchParams.get('status')).toBe('processing')
          return HttpResponse.json({ data: [], total: 0, offset: 0, limit: 10 })
        })
      )

      // Act
      await documentsApi.list({ projectId: 'proj-1', status: 'processing' })

      // Assert - handled in MSW handler
      expect(true).toBe(true)
    })
  })

  describe('get(id)', () => {
    it('should call GET /knowledge/documents/{id} and return document detail with chunk_count', async () => {
      // Arrange
      const mockDetail: DocumentDetail = {
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
        chunks: [
          {
            id: 'chunk-1',
            documentId: 'doc-1',
            chunkIndex: 0,
            content: '用户注册功能概述...',
          },
        ],
      }

      server.use(
        http.get('/api/v1/knowledge/documents/:id', () =>
          HttpResponse.json(mockDetail)
        )
      )

      // Act
      const result = await documentsApi.get('doc-1')

      // Assert
      expect(result).toEqual(mockDetail)
      expect(result.chunkCount).toBe(42)
      expect(result.chunks).toHaveLength(1)
    })
  })

  describe('create(data)', () => {
    it('should call POST /knowledge/documents and return created document', async () => {
      // Arrange
      const uploadData = {
        projectId: 'proj-1',
        name: '新文档.pdf',
        type: 'prd' as const,
        file: new File(['test'], 'test.pdf'),
      }

      const mockDocument: Document = {
        id: 'doc-new',
        projectId: 'proj-1',
        name: '新文档.pdf',
        type: 'prd',
        status: 'pending',
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
      const result = await documentsApi.create(uploadData)

      // Assert
      expect(result).toEqual(mockDocument)
      expect(result.status).toBe('pending')
    })
  })

  describe('delete(id)', () => {
    it('should call DELETE /knowledge/documents/{id} and return 204', async () => {
      // Arrange
      server.use(
        http.delete('/api/v1/knowledge/documents/:id', () =>
          new HttpResponse(null, { status: 204 })
        )
      )

      // Act
      await documentsApi.delete('doc-1')

      // Assert - no exception thrown means success
      expect(true).toBe(true)
    })
  })

  describe('getChunks(docId)', () => {
    it('should call GET /knowledge/documents/{id}/chunks and return chunk list', async () => {
      // Arrange
      const mockChunks: DocumentChunk[] = [
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
      ]

      server.use(
        http.get('/api/v1/knowledge/documents/:id/chunks', () =>
          HttpResponse.json(mockChunks)
        )
      )

      // Act
      const result = await documentsApi.getChunks('doc-1')

      // Assert
      expect(result).toEqual(mockChunks)
      expect(result).toHaveLength(2)
      expect(result[0].chunkIndex).toBe(0)
      expect(result[1].chunkIndex).toBe(1)
    })
  })
})
