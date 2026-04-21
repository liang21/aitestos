/**
 * Drafts API Service
 * Handles draft-related API calls
 */

import { get, post } from '@/lib/request'
import type {
  PaginatedResponse,
  CaseDraft,
  TestCase,
  DraftListParams,
  ConfirmDraftRequest,
  RejectDraftRequest,
  BatchConfirmRequest,
  BatchConfirmResponse,
} from '@/types/api'

/**
 * Drafts API
 */
export const draftsApi = {
  /**
   * Get draft list with optional filters
   */
  getDrafts: async (params?: DraftListParams): Promise<PaginatedResponse<CaseDraft>> => {
    return get<PaginatedResponse<CaseDraft>>('/generation/drafts', { params })
  },

  /**
   * Get pending draft count (for badge)
   */
  getPendingCount: async (): Promise<{ count: number }> => {
    return get<{ count: number }>('/generation/drafts/count')
  },

  /**
   * Get draft detail by ID
   */
  getDraft: async (id: string): Promise<CaseDraft> => {
    return get<CaseDraft>(`/generation/drafts/${id}`)
  },

  /**
   * Confirm draft and convert to test case
   */
  confirmDraft: async (
    draftId: string,
    data: ConfirmDraftRequest
  ): Promise<TestCase> => {
    return post<ConfirmDraftRequest, TestCase>(
      `/generation/drafts/${draftId}/confirm`,
      data
    )
  },

  /**
   * Reject draft with reason
   */
  rejectDraft: async (
    draftId: string,
    data: RejectDraftRequest
  ): Promise<{ success: boolean }> => {
    return post<RejectDraftRequest, { success: boolean }>(
      `/generation/drafts/${draftId}/reject`,
      data
    )
  },

  /**
   * Batch confirm multiple drafts
   */
  batchConfirm: async (data: BatchConfirmRequest): Promise<BatchConfirmResponse> => {
    return post<BatchConfirmRequest, BatchConfirmResponse>(
      '/generation/drafts/batch-confirm',
      data
    )
  },
}
