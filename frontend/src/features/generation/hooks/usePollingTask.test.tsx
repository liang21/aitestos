import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { usePollingTask } from './usePollingTask'
import { server } from '../../../../tests/msw/server'
import { generationHandlers } from '../../../../tests/msw/handlers/generation'
import { http, HttpResponse } from 'msw'

describe('usePollingTask', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
    server.use(...generationHandlers)
  })

  afterEach(() => {
    server.resetHandlers()
  })

  function wrapper({ children }: { children: React.ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    )
  }

  describe('pending state', () => {
    it('should poll when status is pending', async () => {
      const taskId = '550e8400-e29b-41d4-a716-446655440001'

      const { result } = renderHook(() => usePollingTask(taskId), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        id: taskId,
        status: 'pending',
      })

      // Verify refetchInterval function is set and returns 3000 for pending status
      const query = queryClient.getQueryCache().find(['generation-task', taskId])
      const refetchIntervalFn = query?.observers[0].options.refetchInterval
      expect(typeof refetchIntervalFn).toBe('function')
      if (typeof refetchIntervalFn === 'function') {
        expect(refetchIntervalFn(result.current.data)).toBe(3000)
      }
    })
  })

  describe('processing state', () => {
    it('should poll when status is processing', async () => {
      const taskId = '550e8400-e29b-41d4-a716-446655440005'

      const { result } = renderHook(() => usePollingTask(taskId), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        id: taskId,
        status: 'processing',
      })

      // Verify refetchInterval function returns 3000 for processing status
      const query = queryClient.getQueryCache().find(['generation-task', taskId])
      const refetchIntervalFn = query?.observers[0].options.refetchInterval
      if (typeof refetchIntervalFn === 'function') {
        expect(refetchIntervalFn(result.current.data)).toBe(3000)
      }
    })
  })

  describe('completed state', () => {
    it('should stop polling when status is completed', async () => {
      const taskId = '550e8400-e29b-41d4-a716-446655440004'

      const { result } = renderHook(() => usePollingTask(taskId), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        id: taskId,
        status: 'completed',
        result: { draftCount: 5 },
      })

      // Verify refetchInterval function returns false for completed status
      const query = queryClient.getQueryCache().find(['generation-task', taskId])
      const refetchIntervalFn = query?.observers[0].options.refetchInterval
      if (typeof refetchIntervalFn === 'function') {
        expect(refetchIntervalFn(result.current.data)).toBe(false)
      }
    })
  })

  describe('failed state', () => {
    it('should stop polling when status is failed', async () => {
      // Mock a failed task
      server.use(
        http.get('/api/v1/generation/tasks/:id', () =>
          HttpResponse.json({
            id: 'failed-task-id',
            status: 'failed',
            prompt: '测试失败',
            result: { error: 'Generation failed' },
            createdAt: '2026-04-20T10:00:00Z',
            updatedAt: '2026-04-20T10:05:00Z',
          })
        )
      )

      const taskId = 'failed-task-id'

      const { result } = renderHook(() => usePollingTask(taskId), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        status: 'failed',
      })

      // Verify refetchInterval function returns false for failed status
      const query = queryClient.getQueryCache().find(['generation-task', taskId])
      const refetchIntervalFn = query?.observers[0].options.refetchInterval
      if (typeof refetchIntervalFn === 'function') {
        expect(refetchIntervalFn(result.current.data)).toBe(false)
      }
    })
  })

  describe('empty taskId', () => {
    it('should not fetch when taskId is empty string', async () => {
      const { result } = renderHook(() => usePollingTask(''), { wrapper })

      // Query should be disabled
      expect(result.current.fetchStatus).toBe('idle')
      expect(result.current.data).toBeUndefined()
    })

    it('should not fetch when taskId is undefined', async () => {
      const { result } = renderHook(() => usePollingTask(undefined as unknown as string), {
        wrapper,
      })

      // Query should be disabled
      expect(result.current.fetchStatus).toBe('idle')
      expect(result.current.data).toBeUndefined()
    })
  })
})
