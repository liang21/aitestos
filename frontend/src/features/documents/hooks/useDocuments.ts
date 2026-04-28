import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryResult,
} from '@tanstack/react-query'
import { message } from '@arco-design/web-react'
import { documentsApi } from '../services/documents'
import type {
  Document,
  DocumentDetail,
  DocumentChunk,
  UploadDocumentRequest,
} from '@/types/api'

// ============================================================================
// Query Keys
// ============================================================================

const documentKeys = {
  all: ['documents'] as const,
  lists: () => [...documentKeys.all, 'list'] as const,
  list: (params: any) => [...documentKeys.lists(), params] as const,
  details: () => [...documentKeys.all, 'detail'] as const,
  detail: (id: string) => [...documentKeys.details(), id] as const,
  chunks: (id: string) => [...documentKeys.all, 'chunks', id] as const,
}

// ============================================================================
// Hooks
// ============================================================================

/**
 * Query document list
 */
export function useDocumentList(params: {
  projectId: string
  type?: string
  status?: string
  offset?: number
  limit?: number
}): UseQueryResult<{
  data: Document[]
  total: number
  offset: number
  limit: number
}> {
  return useQuery({
    queryKey: documentKeys.list(params),
    queryFn: () => documentsApi.list(params),
  })
}

/**
 * Query document detail
 */
export function useDocumentDetail(id: string): UseQueryResult<DocumentDetail> {
  return useQuery({
    queryKey: documentKeys.detail(id),
    queryFn: () => documentsApi.get(id),
    enabled: !!id,
  })
}

/**
 * Upload document mutation
 */
export function useUploadDocument() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: UploadDocumentRequest) => documentsApi.create(data),
    onSuccess: () => {
      if (message && typeof message.success === 'function') {
        message.success('文档上传成功')
      }
      // Invalidate document lists
      queryClient.invalidateQueries({ queryKey: documentKeys.lists() })
    },
    onError: (error: unknown) => {
      const errorMessage =
        error instanceof Error ? error.message : '文档上传失败'
      if (message && typeof message.error === 'function') {
        message.error(errorMessage)
      }
    },
  })
}

/**
 * Delete document mutation
 */
export function useDeleteDocument() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => documentsApi.delete(id),
    onSuccess: () => {
      if (message && typeof message.success === 'function') {
        message.success('文档删除成功')
      }
      // Invalidate document lists
      queryClient.invalidateQueries({ queryKey: documentKeys.lists() })
    },
    onError: (error: unknown) => {
      const errorMessage =
        error instanceof Error ? error.message : '文档删除失败'
      if (message && typeof message.error === 'function') {
        message.error(errorMessage)
      }
    },
  })
}

/**
 * Query document chunks
 */
export function useDocumentChunks(
  docId: string
): UseQueryResult<DocumentChunk[]> {
  return useQuery({
    queryKey: documentKeys.chunks(docId),
    queryFn: () => documentsApi.getChunks(docId),
    enabled: !!docId,
  })
}
