/**
 * Plans API Service
 */

import { get, post, put, del } from '@/lib/request'
import type {
  PaginatedResponse,
  TestPlan,
  PlanDetail,
  CreatePlanRequest,
  UpdatePlanRequest,
  PlanCase,
  RecordResultRequest,
} from '@/types/api'

/**
 * Plans API
 */
export const plansApi = {
  /**
   * Get plans list with pagination and filters
   */
  list: async (params?: {
    projectId?: string
    status?: string
    keywords?: string
    offset?: number
    limit?: number
  }): Promise<PaginatedResponse<TestPlan>> => {
    return get<PaginatedResponse<TestPlan>>('/plans', { params })
  },

  /**
   * Get plan detail by ID
   */
  get: async (id: string): Promise<PlanDetail> => {
    return get<PlanDetail>(`/plans/${id}`)
  },

  /**
   * Create a new plan
   */
  create: async (data: CreatePlanRequest): Promise<TestPlan> => {
    return post<CreatePlanRequest, TestPlan>('/plans', data)
  },

  /**
   * Update plan
   */
  update: async (id: string, data: UpdatePlanRequest): Promise<TestPlan> => {
    return put<UpdatePlanRequest, TestPlan>(`/plans/${id}`, data)
  },

  /**
   * Delete plan
   */
  delete: async (id: string): Promise<void> => {
    return del<void>(`/plans/${id}`)
  },

  /**
   * Add cases to plan
   */
  addCases: async (planId: string, caseIds: string[]): Promise<{ success: boolean }> => {
    return post<{ case_ids: string[] }, { success: boolean }>(
      `/plans/${planId}/cases`,
      { case_ids: caseIds }
    )
  },

  /**
   * Remove case from plan
   */
  removeCase: async (planId: string, caseId: string): Promise<void> => {
    return del<void>(`/plans/${planId}/cases/${caseId}`)
  },

  /**
   * Get plan execution results
   */
  getResults: async (planId: string): Promise<PlanCase[]> => {
    return get<PlanCase[]>(`/plans/${planId}/results`)
  },

  /**
   * Record test result
   */
  recordResult: async (
    planId: string,
    data: RecordResultRequest
  ): Promise<PlanCase> => {
    return post<RecordResultRequest, PlanCase>(
      `/plans/${planId}/results`,
      data
    )
  },
}
