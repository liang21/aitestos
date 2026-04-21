import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { generationApi } from '@/features/generation/services/generation'
import type {
  PaginatedResponse,
  GenerationTask,
  CreateTaskRequest,
  TaskListParams,
} from '@/types/api'
import { Message } from '@arco-design/web-react'

/**
 * Query keys for generation tasks
 * Exported for use in other hooks and components
 */
export const generationKeys = {
  all: ['generation'] as const,
  tasks: () => [...generationKeys.all, 'tasks'] as const,
  task: (id: string) => [...generationKeys.tasks(), id] as const,
  tasksList: (params: TaskListParams) =>
    [...generationKeys.tasks(), 'list', params] as const,
}

/**
 * Fetch generation tasks list
 */
export function useGenerationTasks(params: TaskListParams) {
  return useQuery({
    queryKey: generationKeys.tasksList(params),
    queryFn: () => generationApi.listTasks(params),
  })
}

/**
 * Fetch single generation task by ID
 */
export function useGenerationTask(taskId: string) {
  return useQuery({
    queryKey: generationKeys.task(taskId),
    queryFn: () => generationApi.getTask(taskId),
    enabled: !!taskId && taskId.length > 0,
  })
}

/**
 * Create a new generation task
 */
export function useCreateGenerationTask() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateTaskRequest) => generationApi.createTask(data),
    onSuccess: () => {
      // Invalidate all tasks queries to refetch
      queryClient.invalidateQueries({ queryKey: generationKeys.tasks() })
    },
    onError: (error: unknown) => {
      const errorMessage =
        error instanceof Error ? error.message : '创建任务失败'
      Message.error(errorMessage)
      console.error('Generation task creation error:', error)
    },
  })
}
