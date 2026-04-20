import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type {
  CreateModuleRequest,
  Module,
  PaginatedResponse,
} from '@/types/api'
import { modulesApi } from '../services/modules'

/**
 * Query keys factory for modules
 */
export const moduleKeys = {
  all: ['modules'] as const,
  lists: () => [...moduleKeys.all, 'list'] as const,
  list: (projectId: string) => [...moduleKeys.lists(), projectId] as const,
}

/**
 * Hook: List modules for a project
 */
export function useModuleList(projectId: string) {
  return useQuery({
    queryKey: moduleKeys.list(projectId),
    queryFn: () => modulesApi.list(projectId),
    enabled: !!projectId,
  })
}

/**
 * Hook: Create new module
 */
export function useCreateModule() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      projectId,
      data,
    }: {
      projectId: string
      data: CreateModuleRequest
    }) => modulesApi.create(projectId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: moduleKeys.list(variables.projectId),
      })
    },
  })
}

/**
 * Hook: Delete module
 */
export function useDeleteModule() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ projectId, id }: { projectId: string; id: string }) =>
      modulesApi.delete(id),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: moduleKeys.list(variables.projectId),
      })
    },
  })
}
