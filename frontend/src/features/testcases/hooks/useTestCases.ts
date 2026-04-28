/**
 * TestCases React Query Hooks
 */

import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryResult,
} from '@tanstack/react-query'
import { testcasesApi } from '../services/testcases'
import type {
  TestCase,
  PaginatedResponse,
  TestCaseListParams,
  CreateTestCaseRequest,
  UpdateTestCaseRequest,
} from '@/types/api'

// Query keys
export const testcaseKeys = {
  all: ['testcases'] as const,
  lists: () => [...testcaseKeys.all, 'list'] as const,
  list: (projectId: string, params?: TestCaseListParams) =>
    [...testcaseKeys.lists(), projectId, params] as const,
  details: () => [...testcaseKeys.all, 'detail'] as const,
  detail: (id: string) => [...testcaseKeys.details(), id] as const,
}

/**
 * Get test case list
 * @param projectId - Project ID from route params
 * @param params - Query parameters (filters, pagination)
 */
export function useCaseList(
  projectId: string,
  params?: TestCaseListParams
): UseQueryResult<PaginatedResponse<TestCase>> {
  return useQuery({
    queryKey: testcaseKeys.list(projectId, params),
    queryFn: () => testcasesApi.list({ project_id: projectId, ...params }),
  })
}

/**
 * Get test case detail
 */
export function useCaseDetail(id: string): UseQueryResult<TestCase> {
  return useQuery({
    queryKey: testcaseKeys.detail(id),
    queryFn: () => testcasesApi.get(id),
    enabled: !!id,
  })
}

/**
 * Create test case mutation
 */
export function useCreateTestCase() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateTestCaseRequest) => testcasesApi.create(data),
    onSuccess: () => {
      // Invalidate test case list queries
      queryClient.invalidateQueries({ queryKey: testcaseKeys.lists() })
    },
  })
}

/**
 * Update test case mutation
 */
export function useUpdateTestCase() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTestCaseRequest }) =>
      testcasesApi.update(id, data),
    onSuccess: () => {
      // Invalidate test case queries
      queryClient.invalidateQueries({ queryKey: testcaseKeys.all })
    },
  })
}

/**
 * Delete test case mutation
 */
export function useDeleteTestCase() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => testcasesApi.delete(id),
    onSuccess: () => {
      // Invalidate test case queries
      queryClient.invalidateQueries({ queryKey: testcaseKeys.all })
    },
  })
}
