/**
 * usePlans Hook Tests
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import {
  usePlanList,
  usePlanDetail,
  useCreatePlan,
  useUpdatePlan,
  useDeletePlan,
  useAddCases,
  useRemoveCase,
  useRecordResult,
} from './usePlans'
import type {
  TestPlan,
  CreatePlanRequest,
  UpdatePlanRequest,
  RecordResultRequest,
} from '@/types/api'
import React from 'react'

describe('usePlans hooks', () => {
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

  const mockPlan: TestPlan = {
    id: 'plan-001',
    projectId: 'project-001',
    name: 'Sprint 12 回归测试',
    description: 'Sprint 12 回归测试计划',
    status: 'draft',
    createdBy: 'user-001',
    createdAt: '2026-04-21T00:00:00Z',
    updatedAt: '2026-04-21T00:00:00Z',
  }

  describe('usePlanList', () => {
    it('should fetch plans by project_id', async () => {
      server.use(
        http.get('/api/v1/plans', () =>
          HttpResponse.json({
            data: [mockPlan],
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      const { result } = renderHook(
        () => usePlanList({ projectId: 'project-001', offset: 0, limit: 10 }),
        { wrapper }
      )

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.data).toHaveLength(1)
      expect(result.current.data?.data[0].name).toBe('Sprint 12 回归测试')
    })

    it('should pass filter params to API', async () => {
      const { result } = renderHook(
        () => usePlanList({ projectId: 'project-001', status: 'active' }),
        { wrapper }
      )

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.data).toBeDefined()
    })
  })

  describe('usePlanDetail', () => {
    it('should fetch plan detail by id', async () => {
      server.use(
        http.get('/api/v1/plans/:id', () =>
          HttpResponse.json({
            ...mockPlan,
            cases: [],
            stats: {
              total: 0,
              passed: 0,
              failed: 0,
              blocked: 0,
              skipped: 0,
              unexecuted: 0,
            },
          })
        )
      )

      const { result } = renderHook(() => usePlanDetail('plan-001'), {
        wrapper,
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.id).toBe('plan-001')
      expect(result.current.data?.name).toBe('Sprint 12 回归测试')
    })

    it('should not fetch when id is empty', async () => {
      const { result } = renderHook(() => usePlanDetail(''), { wrapper })

      expect(result.current.fetchStatus).toBe('idle')
    })
  })

  describe('useCreatePlan', () => {
    it('should create plan and invalidate cache', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      const createData: CreatePlanRequest = {
        projectId: 'project-001',
        name: 'Sprint 13 回归测试',
        description: 'Sprint 13 回归测试计划',
      }

      server.use(
        http.post('/api/v1/plans', () =>
          HttpResponse.json(
            {
              ...mockPlan,
              id: 'plan-002',
              ...createData,
            },
            { status: 201 }
          )
        )
      )

      const { result } = renderHook(() => useCreatePlan(), { wrapper })

      result.current.mutate(createData)

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.id).toBe('plan-002')
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['plans'],
      })
    })
  })

  describe('useUpdatePlan', () => {
    it('should update plan and invalidate cache', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      const updateData: UpdatePlanRequest = {
        name: '更新后的计划名称',
        status: 'active',
      }

      server.use(
        http.put('/api/v1/plans/:id', () =>
          HttpResponse.json({
            ...mockPlan,
            ...updateData,
          })
        )
      )

      const { result } = renderHook(() => useUpdatePlan(), { wrapper })

      result.current.mutate({ id: 'plan-001', data: updateData })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.name).toBe('更新后的计划名称')
      expect(result.current.data?.status).toBe('active')
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['plans'],
      })
    })
  })

  describe('useDeletePlan', () => {
    it('should delete plan and invalidate cache', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      server.use(
        http.delete(
          '/api/v1/plans/:id',
          () => new HttpResponse(null, { status: 204 })
        )
      )

      const { result } = renderHook(() => useDeletePlan(), { wrapper })

      result.current.mutate('plan-001')

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(invalidateSpy).toHaveBeenCalled()
    })
  })

  describe('useAddCases', () => {
    it('should add cases to plan and invalidate cache', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      server.use(
        http.post('/api/v1/plans/:id/cases', () =>
          HttpResponse.json({ success: true })
        )
      )

      const { result } = renderHook(() => useAddCases(), { wrapper })

      result.current.mutate({
        planId: 'plan-001',
        caseIds: ['case-001', 'case-002'],
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toEqual({ success: true })
      expect(invalidateSpy).toHaveBeenCalled()
    })
  })

  describe('useRemoveCase', () => {
    it('should remove case from plan and invalidate cache', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      server.use(
        http.delete(
          '/api/v1/plans/:id/cases/:caseId',
          () => new HttpResponse(null, { status: 204 })
        )
      )

      const { result } = renderHook(() => useRemoveCase(), { wrapper })

      result.current.mutate({ planId: 'plan-001', caseId: 'case-001' })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(invalidateSpy).toHaveBeenCalled()
    })
  })

  describe('useRecordResult', () => {
    it('should record result and invalidate cache', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      const recordData: RecordResultRequest = {
        caseId: 'case-001',
        status: 'pass',
        note: '功能正常',
      }

      server.use(
        http.post('/api/v1/plans/:id/results', () =>
          HttpResponse.json(
            {
              caseId: 'case-001',
              caseNumber: 'TEST-USR-20260421-001',
              caseTitle: '用户登录成功',
              resultStatus: 'pass',
              resultNote: '功能正常',
              executedAt: '2026-04-21T10:00:00Z',
              executedBy: 'user-001',
            },
            { status: 201 }
          )
        )
      )

      const { result } = renderHook(() => useRecordResult(), { wrapper })

      result.current.mutate({ planId: 'plan-001', data: recordData })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data?.caseId).toBe('case-001')
      expect(result.current.data?.resultStatus).toBe('pass')
      expect(invalidateSpy).toHaveBeenCalled()
    })
  })
})
