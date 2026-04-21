import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { configsApi } from '../services/configs'

/**
 * Query keys for configs
 */
export const configsKeys = {
  all: ['configs'] as const,
  lists: () => [...configsKeys.all, 'list'] as const,
  list: (projectId: string) => [...configsKeys.lists(), projectId] as const,
}

/**
 * Get project configuration list
 */
export function useConfigList(projectId: string) {
  return useQuery({
    queryKey: configsKeys.list(projectId),
    queryFn: () => configsApi.list(projectId),
    enabled: !!projectId,
  })
}

/**
 * Set a configuration item
 */
export function useSetConfig(projectId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      key,
      value,
      description,
    }: {
      key: string
      value: unknown
      description?: string
    }) => configsApi.set(projectId, key, { value, description }),
    onSuccess: () => {
      // Invalidate config list cache
      queryClient.invalidateQueries({ queryKey: configsKeys.lists() })
    },
  })
}

/**
 * Delete a configuration item
 */
export function useDeleteConfig(projectId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (key: string) => configsApi.delete(projectId, key),
    onSuccess: () => {
      // Invalidate config list cache
      queryClient.invalidateQueries({ queryKey: configsKeys.lists() })
    },
  })
}

/**
 * Batch import configurations
 */
export function useImportConfigs(projectId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (configs: Array<{ key: string; value: unknown }>) =>
      configsApi.import(projectId, configs),
    onSuccess: () => {
      // Invalidate config list cache
      queryClient.invalidateQueries({ queryKey: configsKeys.lists() })
    },
  })
}

/**
 * Export all project configurations
 */
export function useExportConfigs(projectId: string) {
  return useMutation({
    mutationFn: () => configsApi.export(projectId),
  })
}
