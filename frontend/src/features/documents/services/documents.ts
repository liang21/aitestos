import { get, post, del } from '@/lib/request'
import type {
  Document,
  DocumentDetail,
  DocumentChunk,
  PaginatedResponse,
  UploadDocumentRequest,
} from '@/types/api'

// ============================================================================
// Documents API
// ============================================================================

export const documentsApi = {
  /**
   * Get document list
   * GET /knowledge/documents
   */
  list: (params?: {
    projectId: string
    type?: string
    status?: string
    offset?: number
    limit?: number
  }): Promise<PaginatedResponse<Document>> => {
    return get<PaginatedResponse<Document>>('/knowledge/documents', {
      params,
    })
  },

  /**
   * Get document detail
   * GET /knowledge/documents/{id}
   */
  get: (id: string): Promise<DocumentDetail> => {
    return get<DocumentDetail>(`/knowledge/documents/${id}`)
  },

  /**
   * Upload document
   * POST /knowledge/documents
   */
  create: (data: UploadDocumentRequest): Promise<Document> => {
    return post<UploadDocumentRequest, Document>('/knowledge/documents', data)
  },

  /**
   * Delete document
   * DELETE /knowledge/documents/{id}
   */
  delete: (id: string): Promise<void> => {
    return del<void>(`/knowledge/documents/${id}`)
  },

  /**
   * Get document chunks
   * GET /knowledge/documents/{id}/chunks
   */
  getChunks: (docId: string): Promise<DocumentChunk[]> => {
    return get<DocumentChunk[]>(`/knowledge/documents/${docId}/chunks`)
  },
}
