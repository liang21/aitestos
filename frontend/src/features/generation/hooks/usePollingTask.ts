import { useQuery } from '@tanstack/react-query'
import { generationApi } from '@/features/generation/services/generation'
import type { GenerationTask } from '@/types/api'
import { generationKeys } from './useGeneration'

/**
 * Poll a generation task by ID
 * - Polls every 3 seconds when status is pending or processing
 * - Stops polling when status is completed or failed
 */
export function usePollingTask(taskId: string) {
  return useQuery({
    queryKey: generationKeys.task(taskId),
    queryFn: async () => {
      const task = await generationApi.getTask(taskId)
      return task
    },
    enabled: !!taskId && taskId.length > 0,
    refetchInterval: (data) => {
      // Poll every 3 seconds if task is pending or processing
      if (data?.status === 'pending' || data?.status === 'processing') {
        return 3000
      }
      // Stop polling for completed or failed tasks
      return false
    },
  })
}
