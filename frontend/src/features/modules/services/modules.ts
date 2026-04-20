import type {
  CreateModuleRequest,
  Module,
  PaginatedResponse,
} from '@/types/api'
import request from '@/lib/request'

/**
 * Modules API service
 * Handles all module-related API calls
 */
export const modulesApi = {
  /**
   * List modules for a project
   */
  list: async (
    projectId: string
  ): Promise<PaginatedResponse<Module>> => {
    return request.get(`/projects/${projectId}/modules`)
  },

  /**
   * Create new module
   */
  create: async (
    projectId: string,
    data: CreateModuleRequest
  ): Promise<Module> => {
    return request.post(`/projects/${projectId}/modules`, data)
  },

  /**
   * Delete module (soft delete)
   */
  delete: async (id: string): Promise<void> => {
    return request.delete(`/modules/${id}`)
  },
}
