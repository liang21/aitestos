import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useGenerationTasks, useGenerationTask, useCreateGenerationTask } from './useGeneration'
import { server } from '../../../../tests/msw/server'
import { generationHandlers } from '../../../../tests/msw/handlers/generation'
import { http, HttpResponse } from 'msw'

describe('useGeneration hooks', () => {
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

  describe('useGenerationTasks', () => {
    it('should fetch task list successfully', async () => {
      const { result } = renderHook(
        () =>
          useGenerationTasks({
            projectId: '550e8400-e29b-41d4-a716-446655440002',
            offset: 0,
            limit: 10,
          }),
        { wrapper }
      )

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        data: expect.any(Array),
        total: expect.any(Number),
        offset: 0,
        limit: 10,
      })
    })

    it('should filter tasks by status param', async () => {
      const { result } = renderHook(
        () =>
          useGenerationTasks({
            projectId: '550e8400-e29b-41d4-a716-446655440002',
            status: 'completed',
            offset: 0,
            limit: 10,
          }),
        { wrapper }
      )

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.data.every((task) => task.status === 'completed')).toBe(
        true
      )
    })

    it('should handle loading state', () => {
      const { result } = renderHook(
        () =>
          useGenerationTasks({
            projectId: '550e8400-e29b-41d4-a716-446655440002',
            offset: 0,
            limit: 10,
          }),
        { wrapper }
      )

      expect(result.current.isLoading).toBe(true)
    })
  })

  describe('useGenerationTask', () => {
    it('should fetch single task details', async () => {
      const taskId = '550e8400-e29b-41d4-a716-446655440001'

      const { result } = renderHook(() => useGenerationTask(taskId), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        id: taskId,
        status: 'pending',
        prompt: '测试用户注册功能',
      })
    })

    it('should not fetch when taskId is empty', async () => {
      const { result } = renderHook(() => useGenerationTask(''), { wrapper })

      // When enabled is false, query is disabled
      expect(result.current.fetchStatus).toBe('idle')
      expect(result.current.data).toBeUndefined()
    })
  })

  describe('useCreateGenerationTask', () => {
    it('should create task successfully', async () => {
      const { result } = renderHook(() => useCreateGenerationTask(), { wrapper })

      result.current.mutate({
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        moduleId: '550e8400-e29b-41d4-a716-446655440003',
        prompt: '测试用户注册功能',
        count: 5,
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        id: expect.any(String),
        status: 'pending',
        prompt: '测试用户注册功能',
      })
    })

    it('should handle API errors', async () => {
      // Configure MSW to return error
      server.use(
        http.post('/api/v1/generation/tasks', () =>
          HttpResponse.json({ error: '创建任务失败' }, { status: 400 })
        )
      )

      const { result } = renderHook(() => useCreateGenerationTask(), { wrapper })

      result.current.mutate({
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        moduleId: '550e8400-e29b-41d4-a716-446655440003',
        prompt: '测试',
      })

      await waitFor(() => expect(result.current.isError).toBe(true))
    })

    it('should invalidate queries on success', async () => {
      // First, populate the cache with a tasks query
      queryClient.setQueryData(
        ['generation', 'tasks', 'list', { projectId: 'test', offset: 0, limit: 10 }],
        { data: [], total: 0, offset: 0, limit: 10 }
      )

      const { result } = renderHook(() => useCreateGenerationTask(), { wrapper })

      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      result.current.mutate({
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        moduleId: '550e8400-e29b-41d4-a716-446655440003',
        prompt: '测试用户注册功能',
        count: 5,
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      // Verify invalidateQueries was called
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['generation', 'tasks'],
      })
    })
  })
})
