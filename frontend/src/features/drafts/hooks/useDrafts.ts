/**
 * Drafts React Query Hooks
 */

import { useQuery, useMutation, useQueryClient, type UseQueryResult } from '@tanstack/react-query'
import { draftsApi } from '../services/drafts'
import { DRAFT_POLL_INTERVAL } from '../constants'
import type {
  CaseDraft,
  TestCase,
  DraftListParams,
  ConfirmDraftRequest,
  RejectDraftRequest,
  BatchConfirmRequest,
  BatchConfirmResponse,
} from '@/types/api'

// Query keys
export const draftKeys = {
  all: ['drafts'] as const,
  lists: () => [...draftKeys.all, 'list'] as const,
  list: (params?: DraftListParams) => [...draftKeys.lists(), params] as const,
  details: () => [...draftKeys.all, 'detail'] as const,
  detail: (id: string) => [...draftKeys.details(), id] as const,
  pendingCount: () => [...draftKeys.all, 'pending-count'] as const,
}

/**
 * Get draft list
 */
export function useDraftList(params?: DraftListParams): UseQueryResult<PaginatedResponse<CaseDraft>> {
  return useQuery({
    queryKey: draftKeys.list(params),
    queryFn: () => draftsApi.getDrafts(params),
  })
}

/**
 * Get draft detail
 */
export function useDraftDetail(id: string): UseQueryResult<CaseDraft> {
  return useQuery({
    queryKey: draftKeys.detail(id),
    queryFn: () => draftsApi.getDraft(id),
    enabled: !!id,
  })
}

/**
 * Confirm draft mutation
 */
export function useConfirmDraft() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ draftId, data }: { draftId: string; data: ConfirmDraftRequest }) =>
      draftsApi.confirmDraft(draftId, data),
    onSuccess: () => {
      // Invalidate draft list queries
      queryClient.invalidateQueries({ queryKey: draftKeys.lists() })
      // Invalidate pending count
      queryClient.invalidateQueries({ queryKey: draftKeys.pendingCount() })
    },
  })
}

/**
 * Reject draft mutation
 */
export function useRejectDraft() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ draftId, data }: { draftId: string; data: RejectDraftRequest }) =>
      draftsApi.rejectDraft(draftId, data),
    onSuccess: () => {
      // Invalidate draft list queries
      queryClient.invalidateQueries({ queryKey: draftKeys.lists() })
      // Invalidate pending count
      queryClient.invalidateQueries({ queryKey: draftKeys.pendingCount() })
    },
  })
}

/**
 * Batch confirm mutation
 */
export function useBatchConfirm() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: BatchConfirmRequest) => draftsApi.batchConfirm(data),
    onSuccess: () => {
      // Invalidate draft list queries
      queryClient.invalidateQueries({ queryKey: draftKeys.lists() })
      // Invalidate pending count
      queryClient.invalidateQueries({ queryKey: draftKeys.pendingCount() })
    },
  })
}

/**
 * Get pending draft count with polling (for badge)
 * Polls every 5 seconds when enabled
 */
export function usePendingDraftCount(enabled: boolean = true): UseQueryResult<number> {
  return useQuery({
    queryKey: draftKeys.pendingCount(),
    queryFn: async () => {
      const result = await draftsApi.getPendingCount()
      return result.count
    },
    refetchInterval: enabled ? DRAFT_POLL_INTERVAL : false,
    enabled,
  })
}
