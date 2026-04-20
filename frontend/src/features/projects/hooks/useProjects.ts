import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type {
  CreateProjectRequest,
  Project,
  ProjectDetail,
  ProjectStats,
  UpdateProjectRequest,
} from '@/types/api'
import { projectsApi } from '../services/projects'

/**
 * Query keys factory for projects
 */
export const projectKeys = {
  all: ['projects'] as const,
  lists: () => [...projectKeys.all, 'list'] as const,
  list: (params: Record<string, unknown>) =>
    [...projectKeys.lists(), params] as const,
  details: () => [...projectKeys.all, 'detail'] as const,
  detail: (id: string) => [...projectKeys.details(), id] as const,
  stats: (id: string) => [...projectKeys.all, 'stats', id] as const,
}

/**
 * Hook: List projects with pagination and filtering
 */
export function useProjectList(params?: {
  keywords?: string
  offset?: number
  limit?: number
}) {
  return useQuery({
    queryKey: projectKeys.list(params ?? {}),
    queryFn: () => projectsApi.list(params),
  })
}

/**
 * Hook: Get project detail by ID
 */
export function useProjectDetail(id: string) {
  return useQuery({
    queryKey: projectKeys.detail(id),
    queryFn: () => projectsApi.get(id),
    enabled: !!id,
  })
}

/**
 * Hook: Get project statistics
 */
export function useProjectStats(id: string) {
  return useQuery({
    queryKey: projectKeys.stats(id),
    queryFn: () => projectsApi.getStats(id),
    enabled: !!id,
  })
}

/**
 * Hook: Create new project
 */
export function useCreateProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateProjectRequest) => projectsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() })
    },
  })
}

/**
 * Hook: Update existing project
 */
export function useUpdateProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateProjectRequest }) =>
      projectsApi.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: projectKeys.detail(variables.id),
      })
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() })
    },
  })
}

/**
 * Hook: Delete project
 */
export function useDeleteProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => projectsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() })
    },
  })
}
