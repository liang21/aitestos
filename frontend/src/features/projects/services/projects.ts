import type {
  CreateProjectRequest,
  Project,
  ProjectDetail,
  ProjectStats,
  UpdateProjectRequest,
} from '@/types/api'
import request from '@/lib/request'

/**
 * Projects API service
 * Handles all project-related API calls
 */
export const projectsApi = {
  /**
   * Get projects list with pagination and filtering
   */
  list: async (params?: {
    keywords?: string
    offset?: number
    limit?: number
  }): Promise<{
    data: Project[]
    total: number
    offset: number
    limit: number
  }> => {
    return request.get('/projects', { params })
  },

  /**
   * Get project detail by ID
   */
  get: async (id: string): Promise<ProjectDetail> => {
    return request.get(`/projects/${id}`)
  },

  /**
   * Get project statistics
   */
  getStats: async (id: string): Promise<ProjectStats> => {
    return request.get(`/projects/${id}/stats`)
  },

  /**
   * Create new project
   */
  create: async (data: CreateProjectRequest): Promise<Project> => {
    return request.post('/projects', data)
  },

  /**
   * Update existing project
   */
  update: async (id: string, data: UpdateProjectRequest): Promise<Project> => {
    return request.put(`/projects/${id}`, data)
  },

  /**
   * Delete project (soft delete)
   */
  delete: async (id: string): Promise<void> => {
    return request.delete(`/projects/${id}`)
  },
}
