/**
 * Plans API MSW Handlers
 */

import { http, HttpResponse } from 'msw'
import type { TestPlan, PlanDetail, PlanCase, PlanStats } from '@/types/api'

export const plansHandlers = [
  // GET /api/v1/plans - List plans
  http.get('/api/v1/plans', () => {
    const mockPlans: TestPlan[] = [
      {
        id: 'plan-001',
        projectId: 'project-001',
        name: 'Sprint 12 回归测试',
        description: 'Sprint 12 回归测试计划',
        status: 'draft',
        createdBy: 'user-001',
        createdAt: '2026-04-21T00:00:00Z',
        updatedAt: '2026-04-21T00:00:00Z',
      },
    ]
    return HttpResponse.json({
      data: mockPlans,
      total: 1,
      offset: 0,
      limit: 10,
    })
  }),

  // GET /api/v1/plans/:id - Get plan detail
  http.get('/api/v1/plans/:id', () => {
    const mockCases: PlanCase[] = [
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
    const mockStats: PlanStats = {
      total: 2,
      passed: 1,
      failed: 1,
      blocked: 0,
      skipped: 0,
      unexecuted: 0,
    }
    const mockDetail: PlanDetail = {
      id: 'plan-001',
      projectId: 'project-001',
      name: 'Sprint 12 回归测试',
      description: 'Sprint 12 回归测试计划',
      status: 'draft',
      createdBy: 'user-001',
      createdAt: '2026-04-21T00:00:00Z',
      updatedAt: '2026-04-21T00:00:00Z',
      cases: mockCases,
      stats: mockStats,
    }
    return HttpResponse.json(mockDetail)
  }),

  // POST /api/v1/plans - Create plan
  http.post('/api/v1/plans', () => {
    const newPlan: TestPlan = {
      id: 'plan-002',
      projectId: 'project-001',
      name: 'Sprint 13 回归测试',
      description: 'Sprint 13 回归测试计划',
      status: 'draft',
      createdBy: 'user-001',
      createdAt: '2026-04-21T00:00:00Z',
      updatedAt: '2026-04-21T00:00:00Z',
    }
    return HttpResponse.json(newPlan, { status: 201 })
  }),

  // PUT /api/v1/plans/:id - Update plan
  http.put('/api/v1/plans/:id', () => {
    const updatedPlan: TestPlan = {
      id: 'plan-001',
      projectId: 'project-001',
      name: '更新后的计划名称',
      description: 'Sprint 12 回归测试计划',
      status: 'active',
      createdBy: 'user-001',
      createdAt: '2026-04-21T00:00:00Z',
      updatedAt: '2026-04-21T00:00:00Z',
    }
    return HttpResponse.json(updatedPlan)
  }),

  // DELETE /api/v1/plans/:id - Delete plan
  http.delete('/api/v1/plans/:id', () => {
    return new HttpResponse(null, { status: 204 })
  }),

  // POST /api/v1/plans/:id/cases - Add cases to plan
  http.post('/api/v1/plans/:id/cases', () => {
    return HttpResponse.json({ success: true })
  }),

  // DELETE /api/v1/plans/:id/cases/:caseId - Remove case from plan
  http.delete('/api/v1/plans/:id/cases/:caseId', () => {
    return new HttpResponse(null, { status: 204 })
  }),

  // GET /api/v1/plans/:id/results - Get plan results
  http.get('/api/v1/plans/:id/results', () => {
    const mockResults: PlanCase[] = [
      {
        caseId: 'case-001',
        caseNumber: 'TEST-USR-20260421-001',
        caseTitle: '用户登录成功',
        resultStatus: 'pass',
        resultNote: '功能正常',
        executedAt: '2026-04-21T10:00:00Z',
        executedBy: 'user-001',
      },
    ]
    return HttpResponse.json(mockResults)
  }),

  // POST /api/v1/plans/:id/results - Record result
  http.post('/api/v1/plans/:id/results', () => {
    const mockResult: PlanCase = {
      caseId: 'case-001',
      caseNumber: 'TEST-USR-20260421-001',
      caseTitle: '用户登录成功',
      resultStatus: 'pass',
      resultNote: '功能正常',
      executedAt: '2026-04-21T10:00:00Z',
      executedBy: 'user-001',
    }
    return HttpResponse.json(mockResult, { status: 201 })
  }),
]
