/**
 * useTestCases Hook Tests
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import {
  useCaseList,
  useCaseDetail,
  useCreateTestCase,
  useUpdateTestCase,
  useDeleteTestCase,
} from './useTestCases'
import type {
  TestCase,
  CreateTestCaseRequest,
  UpdateTestCaseRequest,
} from '@/types/api'

describe('useTestCases hooks', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
  })

  afterEach(() => {
    server.resetHandlers()
  })

  function wrapper({ children }: { children: React.ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    )
  }

  const mockCase: TestCase = {
    id: 'case-001',
    moduleId: 'mod-001',
    userId: 'user-001',
    number: 'TEST-USR-20260421-001',
    title: '用户登录成功',
    preconditions: ['用户已注册'],
    steps: ['打开登录页', '输入用户名密码', '点击登录'],
    expected: { step_3: '登录成功，跳转首页' },
    caseType: 'functionality',
    priority: 'high',
    status: 'unexecuted',
    createdAt: '2026-04-21T00:00:00Z',
    updatedAt: '2026-04-21T00:00:00Z',
  }

  describe('useCaseList', () => {
    it('should fetch cases by project_id', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: [mockCase],
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      const { result } = renderHook(
        () => useCaseList({ projectId: 'project-001', offset: 0, limit: 10 }),
        { wrapper }
      )

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.data).toHaveLength(1)
      expect(result.current.data?.data[0].title).toBe('用户登录成功')
    })

    it('should pass filter params to API', async () => {
      const { result } = renderHook(
        () =>
          useCaseList({
            projectId: 'project-001',
            status: 'unexecuted',
            caseType: 'functionality',
          }),
        { wrapper }
      )

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.data).toBeDefined()
    })
  })

  describe('useCaseDetail', () => {
    it('should fetch case detail by id', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      const { result } = renderHook(() => useCaseDetail('case-001'), {
        wrapper,
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.id).toBe('case-001')
      expect(result.current.data?.title).toBe('用户登录成功')
    })

    it('should not fetch when id is empty', async () => {
      const { result } = renderHook(() => useCaseDetail(''), { wrapper })

      expect(result.current.fetchStatus).toBe('idle')
    })
  })

  describe('useCreateTestCase', () => {
    it('should create case and invalidate cache', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      const createData: CreateTestCaseRequest = {
        moduleId: 'mod-001',
        title: '新用例',
        preconditions: [],
        steps: ['步骤1'],
        expected: {},
        caseType: 'functionality',
        priority: 'medium',
      }

      server.use(
        http.post('/api/v1/testcases', () =>
          HttpResponse.json(
            {
              ...mockCase,
              id: 'case-002',
              title: '新用例',
            },
            { status: 201 }
          )
        )
      )

      const { result } = renderHook(() => useCreateTestCase(), { wrapper })

      result.current.mutate(createData)

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.id).toBe('case-002')
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['testcases', 'list'],
      })
    })
  })

  describe('useUpdateTestCase', () => {
    it('should update case and invalidate cache', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      const updateData: UpdateTestCaseRequest = {
        title: '更新标题',
        status: 'passed',
      }

      server.use(
        http.put('/api/v1/testcases/:id', () =>
          HttpResponse.json({
            ...mockCase,
            title: '更新标题',
            status: 'passed',
          })
        )
      )

      const { result } = renderHook(() => useUpdateTestCase(), { wrapper })

      result.current.mutate({ id: 'case-001', data: updateData })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.title).toBe('更新标题')
      expect(result.current.data?.status).toBe('passed')
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['testcases'],
      })
    })
  })

  describe('useDeleteTestCase', () => {
    it('should delete case and invalidate cache', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      server.use(
        http.delete('/api/v1/testcases/:id', () =>
          HttpResponse.json({ success: true })
        )
      )

      const { result } = renderHook(() => useDeleteTestCase(), { wrapper })

      result.current.mutate('case-001')

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toEqual({ success: true })
      expect(invalidateSpy).toHaveBeenCalled()
    })
  })
})
