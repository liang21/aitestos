/**
 * TestCases API Service
 * Handles test case related API calls
 */

import { get, post, put, del } from '@/lib/request'
import type {
  PaginatedResponse,
  TestCase,
  CreateTestCaseRequest,
  UpdateTestCaseRequest,
  TestCaseListParams,
} from '@/types/api'

/**
 * TestCases API
 */
export const testcasesApi = {
  /**
   * Get test cases list with optional filters
   */
  list: async (
    params?: TestCaseListParams
  ): Promise<PaginatedResponse<TestCase>> => {
    return get<PaginatedResponse<TestCase>>('/testcases', { params })
  },

  /**
   * Get test case detail by ID
   */
  get: async (id: string): Promise<TestCase> => {
    return get<TestCase>(`/testcases/${id}`)
  },

  /**
   * Create a new test case
   */
  create: async (data: CreateTestCaseRequest): Promise<TestCase> => {
    return post<CreateTestCaseRequest, TestCase>('/testcases', data)
  },

  /**
   * Update test case
   */
  update: async (
    id: string,
    data: UpdateTestCaseRequest
  ): Promise<TestCase> => {
    return put<UpdateTestCaseRequest, TestCase>(`/testcases/${id}`, data)
  },

  /**
   * Delete test case
   */
  delete: async (id: string): Promise<{ success: boolean }> => {
    return del<{ success: boolean }>(`/testcases/${id}`)
  },
}
