import { useQuery } from '@tanstack/react-query'
import { generationApi } from '@/features/generation/services/generation'
import { generationKeys } from './useGeneration'
import { GENERATION_CONFIG } from '@/features/generation/constants'

/**
 * Poll a generation task by ID
 * - Polls at configured interval when status is pending or processing
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
      // Poll when task is in active status
      if (
        data?.status === 'pending' ||
        data?.status === 'processing'
      ) {
        return GENERATION_CONFIG.POLLING.INTERVAL_MS
      }
      // Stop polling for completed or failed tasks
      return false
    },
    refetchIntervalInBackground: false,
  })
}
