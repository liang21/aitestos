/**
 * Plans API Service Tests
 */

import { describe, it, expect, afterEach } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { plansApi } from './plans'
import type {
  TestPlan,
  PlanDetail,
  CreatePlanRequest,
  UpdatePlanRequest,
  PlanCase,
  PlanStats,
  RecordResultRequest,
} from '@/types/api'

describe('plansApi', () => {
  afterEach(() => {
    server.resetHandlers()
  })

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

  const mockPlanCases: PlanCase[] = [
    {
      caseId: 'case-001',
      caseNumber: 'TEST-USR-20260421-001',
      caseTitle: '用户登录成功',
      resultStatus: 'pass',
      resultNote: '功能正常',
      executedAt: '2026-04-21T10:00:00Z',
      executedBy: 'user-001',
    },
    {
      caseId: 'case-002',
      caseNumber: 'TEST-USR-20260421-002',
      caseTitle: '用户注册成功',
      resultStatus: 'fail',
      resultNote: '登录超时',
      executedAt: '2026-04-21T10:05:00Z',
      executedBy: 'user-001',
    },
  ]

  const mockPlanStats: PlanStats = {
    total: 2,
    passed: 1,
    failed: 1,
    blocked: 0,
    skipped: 0,
    unexecuted: 0,
  }

  describe('list', () => {
    it('should fetch plans list with params', async () => {
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

      const result = await plansApi.list({
        projectId: 'project-001',
        offset: 0,
        limit: 10,
      })

      expect(result.data).toEqual([mockPlan])
      expect(result.total).toBe(1)
    })

    it('should filter by status', async () => {
      server.use(
        http.get('/api/v1/plans', () =>
          HttpResponse.json({
            data: [{ ...mockPlan, status: 'active' }],
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      const result = await plansApi.list({
        projectId: 'project-001',
        status: 'active',
      })

      expect(result.data).toHaveLength(1)
      expect(result.data[0].status).toBe('active')
    })

    it('should search by keywords', async () => {
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

      const result = await plansApi.list({
        projectId: 'project-001',
        keywords: 'Sprint 12',
      })

      expect(result.data).toBeDefined()
    })
  })

  describe('get', () => {
    it('should fetch plan detail by id', async () => {
      const mockDetail: PlanDetail = {
        ...mockPlan,
        cases: mockPlanCases,
        stats: mockPlanStats,
      }

      server.use(
        http.get('/api/v1/plans/:id', () => HttpResponse.json(mockDetail))
      )

      const result = await plansApi.get('plan-001')

      expect(result.id).toBe('plan-001')
      expect(result.cases).toHaveLength(2)
      expect(result.stats.total).toBe(2)
    })
  })

  describe('create', () => {
    it('should create a new plan', async () => {
      const createData: CreatePlanRequest = {
        projectId: 'project-001',
        name: 'Sprint 13 回归测试',
        description: 'Sprint 13 回归测试计划',
      }

      server.use(
        http.post('/api/v1/plans', () =>
          HttpResponse.json({ ...mockPlan, id: 'plan-002', ...createData }, {
            status: 201,
          })
        )
      )

      const result = await plansApi.create(createData)

      expect(result.id).toBe('plan-002')
      expect(result.name).toBe('Sprint 13 回归测试')
    })
  })

  describe('update', () => {
    it('should update plan', async () => {
      const updateData: UpdatePlanRequest = {
        name: '更新后的计划名称',
        status: 'active',
      }

      server.use(
        http.put('/api/v1/plans/:id', () =>
          HttpResponse.json({ ...mockPlan, ...updateData })
        )
      )

      const result = await plansApi.update('plan-001', updateData)

      expect(result.name).toBe('更新后的计划名称')
      expect(result.status).toBe('active')
    })
  })

  describe('delete', () => {
    it('should delete plan', async () => {
      server.use(
        http.delete('/api/v1/plans/:id', () =>
          new HttpResponse(null, { status: 204 })
        )
      )

      const result = await plansApi.delete('plan-001')

      // 204 responses return empty string from axios
      expect(result).toBe('')
    })
  })

  describe('addCases', () => {
    it('should add cases to plan', async () => {
      server.use(
        http.post('/api/v1/plans/:id/cases', () =>
          HttpResponse.json({ success: true })
        )
      )

      const result = await plansApi.addCases('plan-001', ['case-001', 'case-002'])

      expect(result).toEqual({ success: true })
    })
  })

  describe('removeCase', () => {
    it('should remove case from plan', async () => {
      server.use(
        http.delete('/api/v1/plans/:id/cases/:caseId', () =>
          new HttpResponse(null, { status: 204 })
        )
      )

      const result = await plansApi.removeCase('plan-001', 'case-001')

      // 204 responses return empty string from axios
      expect(result).toBe('')
    })
  })

  describe('getResults', () => {
    it('should fetch plan results', async () => {
      server.use(
        http.get('/api/v1/plans/:id/results', () =>
          HttpResponse.json(mockPlanCases)
        )
      )

      const result = await plansApi.getResults('plan-001')

      expect(result).toHaveLength(2)
      expect(result[0].caseId).toBe('case-001')
    })
  })

  describe('recordResult', () => {
    it('should record test result', async () => {
      const recordData: RecordResultRequest = {
        caseId: 'case-001',
        status: 'pass',
        note: '功能正常',
      }

      server.use(
        http.post('/api/v1/plans/:id/results', () =>
          HttpResponse.json({
            caseId: 'case-001',
            status: 'pass',
            note: '功能正常',
          }, { status: 201 })
        )
      )

      const result = await plansApi.recordResult('plan-001', recordData)

      expect(result.caseId).toBe('case-001')
      expect(result.status).toBe('pass')
    })
  })
})
