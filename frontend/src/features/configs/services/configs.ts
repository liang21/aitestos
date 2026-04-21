import type {
  ConfigItem,
  ExportConfigsResponse,
  ImportConfigsRequest,
  SetConfigRequest,
} from '@/types/api'
import request from '@/lib/request'

/**
 * Configs API service
 * Handles all project configuration-related API calls
 */
export const configsApi = {
  /**
   * Get project configuration list
   */
  list: async (projectId: string): Promise<{ data: ConfigItem[] }> => {
    return request.get(`/projects/${projectId}/configs`)
  },

  /**
   * Set a configuration item
   */
  set: async (
    projectId: string,
    key: string,
    data: SetConfigRequest
  ): Promise<ConfigItem> => {
    return request.put(`/projects/${projectId}/configs/${key}`, data)
  },

  /**
   * Delete a configuration item
   */
  delete: async (projectId: string, key: string): Promise<void> => {
    return request.delete(`/projects/${projectId}/configs/${key}`)
  },

  /**
   * Batch import configurations
   */
  import: async (
    projectId: string,
    configs: ImportConfigsRequest['configs']
  ): Promise<{
    successCount: number
    failedCount: number
    errors?: Array<{ key: string; error: string }>
  }> => {
    return request.post(`/projects/${projectId}/configs/import`, { configs })
  },

  /**
   * Export all project configurations
   */
  export: async (projectId: string): Promise<ExportConfigsResponse> => {
    return request.get(`/projects/${projectId}/configs/export`)
  },
}
