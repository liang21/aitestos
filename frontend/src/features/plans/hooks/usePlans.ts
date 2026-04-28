/**
 * Plans React Query Hooks
 */

import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryResult,
} from '@tanstack/react-query'
import { plansApi } from '../services/plans'
import type {
  PaginatedResponse,
  TestPlan,
  PlanDetail,
  CreatePlanRequest,
  UpdatePlanRequest,
  RecordResultRequest,
} from '@/types/api'
import type { PlanCase, PlanStats } from '@/types/api'

// Query keys
export const planKeys = {
  all: ['plans'] as const,
  lists: () => [...planKeys.all, 'list'] as const,
  list: (projectId: string, params?: {
    status?: string
    keywords?: string
    offset?: number
    limit?: number
  }) => [...planKeys.lists(), projectId, params] as const,
  details: () => [...planKeys.all, 'detail'] as const,
  detail: (id: string) => [...planKeys.details(), id] as const,
  results: (id: string) => [...planKeys.all, 'results', id] as const,
}

/**
 * Get plans list
 * @param projectId - Project ID from route params
 * @param params - Query parameters (filters, pagination)
 */
export function usePlanList(
  projectId: string,
  params?: {
    status?: string
    keywords?: string
    offset?: number
    limit?: number
  }
): UseQueryResult<PaginatedResponse<TestPlan>> {
  return useQuery({
    queryKey: planKeys.list(projectId, params),
    queryFn: () => plansApi.list({ project_id: projectId, ...params }),
  })
}

/**
 * Get plan detail
 */
export function usePlanDetail(id: string): UseQueryResult<PlanDetail> {
  return useQuery({
    queryKey: planKeys.detail(id),
    queryFn: () => plansApi.get(id),
    enabled: !!id,
  })
}

/**
 * Create plan mutation
 */
export function useCreatePlan() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreatePlanRequest) => plansApi.create(data),
    onSuccess: () => {
      // Invalidate plan queries
      queryClient.invalidateQueries({ queryKey: planKeys.all })
    },
  })
}

/**
 * Update plan mutation
 */
export function useUpdatePlan() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePlanRequest }) =>
      plansApi.update(id, data),
    onSuccess: () => {
      // Invalidate plan queries
      queryClient.invalidateQueries({ queryKey: planKeys.all })
    },
  })
}

/**
 * Delete plan mutation
 */
export function useDeletePlan() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => plansApi.delete(id),
    onSuccess: () => {
      // Invalidate plan queries
      queryClient.invalidateQueries({ queryKey: planKeys.all })
    },
  })
}

/**
 * Add cases to plan mutation
 */
export function useAddCases() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ planId, caseIds }: { planId: string; caseIds: string[] }) =>
      plansApi.addCases(planId, caseIds),
    onMutate: async ({ planId }) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: planKeys.detail(planId) })

      // Snapshot previous value
      const previousPlan = queryClient.getQueryData(planKeys.detail(planId))

      // Optimistically update to the new value
      queryClient.setQueryData(
        planKeys.detail(planId),
        (old: PlanDetail | undefined) => {
          if (!old) return old
          // Add new cases to the existing list (they would be added server-side)
          return { ...old, stats: { ...old.stats, total: old.stats.total + 1 } }
        }
      )

      // Return context with previous value
      return { previousPlan }
    },
    onError: (err, variables, context) => {
      // Rollback to previous value
      if (context?.previousPlan) {
        queryClient.setQueryData(
          planKeys.detail(context?.variables.planId),
          context.previousPlan
        )
      }
    },
    onSuccess: () => {
      // Invalidate to ensure consistency with server
      queryClient.invalidateQueries({ queryKey: planKeys.all })
    },
  })
}

/**
 * Remove case from plan mutation
 */
export function useRemoveCase() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ planId, caseId }: { planId: string; caseId: string }) =>
      plansApi.removeCase(planId, caseId),
    onMutate: async ({ planId, caseId }) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: planKeys.detail(planId) })

      // Snapshot previous value
      const previousPlan = queryClient.getQueryData(planKeys.detail(planId))

      // Optimistically update to remove the case
      queryClient.setQueryData(
        planKeys.detail(planId),
        (old: PlanDetail | undefined) => {
          if (!old) return old
          return {
            ...old,
            cases: old.cases.filter((c) => c.caseId !== caseId),
            stats: {
              ...old.stats,
              total: old.stats.total - 1,
              // Adjust status counts based on removed case
              ...(old.cases.find((c) => c.caseId === caseId)?.resultStatus ===
                'pass' && {
                passed: old.stats.passed - 1,
              }),
              ...(old.cases.find((c) => c.caseId === caseId)?.resultStatus ===
                'fail' && {
                failed: old.stats.failed - 1,
              }),
            },
          }
        }
      )

      // Return context with previous value
      return { previousPlan }
    },
    onError: (err, variables, context) => {
      // Rollback to previous value
      if (context?.previousPlan) {
        queryClient.setQueryData(
          planKeys.detail(context?.variables.planId),
          context.previousPlan
        )
      }
    },
    onSuccess: () => {
      // Invalidate to ensure consistency with server
      queryClient.invalidateQueries({ queryKey: planKeys.all })
    },
  })
}

/**
 * Record test result mutation
 */
export function useRecordResult() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      planId,
      data,
    }: {
      planId: string
      data: RecordResultRequest
    }) => plansApi.recordResult(planId, data),
    onMutate: async ({ planId, data }) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: planKeys.detail(planId) })

      // Snapshot previous value
      const previousPlan = queryClient.getQueryData(planKeys.detail(planId))

      // Optimistically update the case result
      queryClient.setQueryData(
        planKeys.detail(planId),
        (old: PlanDetail | undefined) => {
          if (!old) return old

          // Find and update the specific case
          const updatedCases = old.cases.map((c) =>
            c.caseId === data.caseId
              ? {
                  ...c,
                  resultStatus: data.status,
                  resultNote: data.note,
                  executedAt: new Date().toISOString(),
                }
              : c
          )

          // Recalculate stats
          const oldResult = old.cases.find(
            (c) => c.caseId === data.caseId
          )?.resultStatus
          const stats = { ...old.stats }

          // Decrement old status count if case was previously executed
          if (oldResult === 'pass') stats.passed = Math.max(0, stats.passed - 1)
          if (oldResult === 'fail') stats.failed = Math.max(0, stats.failed - 1)
          if (oldResult === 'block')
            stats.blocked = Math.max(0, stats.blocked - 1)
          if (oldResult === 'skip')
            stats.skipped = Math.max(0, stats.skipped - 1)
          if (oldResult === undefined)
            stats.unexecuted = Math.max(0, stats.unexecuted - 1)

          // Increment new status count
          if (data.status === 'pass') stats.passed += 1
          if (data.status === 'fail') stats.failed += 1
          if (data.status === 'block') stats.blocked += 1
          if (data.status === 'skip') stats.skipped += 1

          return {
            ...old,
            cases: updatedCases,
            stats,
          }
        }
      )

      // Return context with previous value
      return { previousPlan }
    },
    onError: (err, variables, context) => {
      // Rollback to previous value
      if (context?.previousPlan) {
        queryClient.setQueryData(
          planKeys.detail(context?.variables.planId),
          context.previousPlan
        )
      }
    },
    onSuccess: () => {
      // Invalidate to ensure consistency with server
      queryClient.invalidateQueries({ queryKey: planKeys.all })
    },
  })
}

/**
 * Update plan status mutation
 */
export function useUpdatePlanStatus() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ planId, status }: { planId: string; status: string }) =>
      plansApi.updateStatus(planId, status),
    onSuccess: () => {
      // Invalidate plan queries
      queryClient.invalidateQueries({ queryKey: planKeys.all })
    },
  })
}

/**
 * Delete result mutation (for undo functionality)
 */
export function useDeleteResult() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ planId, caseId }: { planId: string; caseId: string }) =>
      plansApi.deleteResult(planId, caseId),
    onMutate: async ({ planId, caseId }) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: planKeys.detail(planId) })

      // Snapshot previous value
      const previousPlan = queryClient.getQueryData(planKeys.detail(planId))

      // Optimistically remove the result
      queryClient.setQueryData(
        planKeys.detail(planId),
        (old: PlanDetail | undefined) => {
          if (!old) return old

          // Find the case and remove its result
          const updatedCases = old.cases.map((c) =>
            c.caseId === caseId
              ? {
                  ...c,
                  resultStatus: undefined,
                  resultNote: undefined,
                  executedAt: undefined,
                  executedBy: undefined,
                }
              : c
          )

          // Recalculate stats
          const targetCase = old.cases.find((c) => c.caseId === caseId)
          const oldResult = targetCase?.resultStatus
          const stats = { ...old.stats }

          // Decrement the old status count
          if (oldResult === 'pass') stats.passed = Math.max(0, stats.passed - 1)
          if (oldResult === 'fail') stats.failed = Math.max(0, stats.failed - 1)
          if (oldResult === 'block')
            stats.blocked = Math.max(0, stats.blocked - 1)
          if (oldResult === 'skip')
            stats.skipped = Math.max(0, stats.skipped - 1)

          // Increment unexecuted
          stats.unexecuted += 1

          return {
            ...old,
            cases: updatedCases,
            stats,
          }
        }
      )

      // Return context with previous value
      return { previousPlan }
    },
    onError: (err, variables, context) => {
      // Rollback to previous value
      if (context?.previousPlan) {
        queryClient.setQueryData(
          planKeys.detail(context?.variables.planId),
          context.previousPlan
        )
      }
    },
    onSuccess: () => {
      // Invalidate to ensure consistency with server
      queryClient.invalidateQueries({ queryKey: planKeys.all })
    },
  })
}
