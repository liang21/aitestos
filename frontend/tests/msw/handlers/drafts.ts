import { http, HttpResponse } from 'msw'
import type { CaseDraft, TestCase } from '@/types/api'

/**
 * MSW handlers for Draft API endpoints
 */
export const draftsHandlers = [
  // GET /api/v1/generation/drafts/count - 获取待处理草稿数量
  http.get('/api/v1/generation/drafts/count', () => {
    return HttpResponse.json({ count: 2 })
  }),

  // GET /api/v1/generation/drafts - 获取草稿列表
  http.get('/api/v1/generation/drafts', ({ request }) => {
    const url = new URL(request.url)
    const projectId = url.searchParams.get('projectId')
    const status = url.searchParams.get('status')

    const allDrafts: CaseDraft[] = [
      {
        id: 'draft-001',
        taskId: 'task-001',
        title: '验证有效邮箱注册',
        preconditions: ['用户未注册过该邮箱'],
        steps: ['打开注册页面', '输入有效邮箱', '输入密码', '点击注册'],
        expected: {
          step_1: '页面正常加载',
          step_2: '邮箱格式正确',
          step_3: '密码通过校验',
          step_4: '注册成功',
        },
        caseType: 'functionality',
        priority: 'P1',
        status: 'pending',
        createdAt: '2026-04-20T10:00:00Z',
        updatedAt: '2026-04-20T10:00:00Z',
        projectName: 'ECommerce',
        moduleName: '用户中心',
      },
      {
        id: 'draft-002',
        taskId: 'task-001',
        title: '验证邮箱格式校验',
        preconditions: ['用户在注册页面'],
        steps: ['打开注册页面', '输入无效邮箱格式', '点击注册'],
        expected: {
          step_1: '页面正常加载',
          step_2: '显示邮箱格式错误',
          step_3: '注册失败',
        },
        caseType: 'functionality',
        priority: 'P2',
        status: 'pending',
        createdAt: '2026-04-20T10:01:00Z',
        updatedAt: '2026-04-20T10:01:00Z',
        projectName: 'ECommerce',
        moduleName: '用户中心',
      },
      {
        id: 'draft-003',
        taskId: 'task-002',
        title: '验证密码强度校验',
        preconditions: ['用户在注册页面'],
        steps: ['打开注册页面', '输入弱密码', '点击注册'],
        expected: {
          step_1: '页面正常加载',
          step_2: '显示密码强度不足',
          step_3: '注册失败',
        },
        caseType: 'functionality',
        priority: 'P1',
        status: 'confirmed',
        createdAt: '2026-04-20T10:02:00Z',
        updatedAt: '2026-04-20T10:05:00Z',
        projectName: 'ECommerce',
        moduleName: '用户中心',
      },
    ]

    // Filter by projectId and status
    let filtered = allDrafts
    if (projectId) {
      filtered = filtered.filter((d) => d.projectName === 'ECommerce')
    }
    if (status) {
      filtered = filtered.filter((d) => d.status === status)
    }

    return HttpResponse.json({
      data: filtered,
      total: filtered.length,
      offset: 0,
      limit: 10,
    })
  }),

  // POST /api/v1/generation/drafts/:id/confirm - 确认草稿
  http.post('/api/v1/generation/drafts/:id/confirm', () => {
    return HttpResponse.json<TestCase>({
      id: 'tc-001',
      moduleId: 'mod-001',
      userId: 'user-001',
      number: 'ECO-USR-20260421-001',
      title: '验证有效邮箱注册',
      preconditions: ['用户未注册过该邮箱'],
      steps: ['打开注册页面', '输入有效邮箱', '输入密码', '点击注册'],
      expected: {
        step_1: '页面正常加载',
        step_2: '邮箱格式正确',
        step_3: '密码通过校验',
        step_4: '注册成功',
      },
      caseType: 'functionality',
      priority: 'P1',
      status: 'unexecuted',
      createdAt: '2026-04-20T10:00:00Z',
      updatedAt: '2026-04-20T10:00:00Z',
    })
  }),

  // POST /api/v1/generation/drafts/:id/reject - 拒绝草稿
  http.post('/api/v1/generation/drafts/:id/reject', () => {
    return HttpResponse.json({
      success: true,
      message: '草稿已拒绝',
    })
  }),

  // POST /api/v1/generation/drafts/batch-confirm - 批量确认
  http.post('/api/v1/generation/drafts/batch-confirm', async () => {
    return HttpResponse.json({
      successCount: 2,
      failedCount: 0,
    })
  }),
]
