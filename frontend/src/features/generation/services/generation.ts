import request from '@/lib/request'
import type {
  GenerationTask,
  CreateTaskRequest,
  TaskListParams,
  PaginatedResponse,
  CaseDraft,
} from '@/types/api'

/**
 * Generation API service
 * Handles AI test case generation tasks
 */
export const generationApi = {
  /**
   * Create a new generation task
   * POST /api/v1/generation/tasks
   */
  createTask: async (data: CreateTaskRequest): Promise<GenerationTask> => {
    return request.post<GenerationTask>('/generation/tasks', data)
  },

  /**
   * Get task details by ID
   * GET /api/v1/generation/tasks/:id
   */
  getTask: async (id: string): Promise<GenerationTask> => {
    return request.get<GenerationTask>(`/generation/tasks/${id}`)
  },

  /**
   * List generation tasks with filters
   * GET /api/v1/generation/tasks
   */
  listTasks: async (
    params: TaskListParams
  ): Promise<PaginatedResponse<GenerationTask>> => {
    return request.get<PaginatedResponse<GenerationTask>>('/generation/tasks', {
      params,
    })
  },

  /**
   * Get drafts for a specific task
   * GET /api/v1/generation/tasks/:id/drafts
   */
  getTaskDrafts: async (taskId: string): Promise<CaseDraft[]> => {
    return request.get<CaseDraft[]>(`/generation/tasks/${taskId}/drafts`)
  },
}
