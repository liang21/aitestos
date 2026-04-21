import { http, HttpResponse } from 'msw'
import type { GenerationTask, CaseDraft } from '@/types/api'

export const generationHandlers = [
  // POST /api/v1/generation/tasks - 创建生成任务
  http.post('/api/v1/generation/tasks', () =>
    HttpResponse.json<GenerationTask>({
      id: '550e8400-e29b-41d4-a716-446655440001',
      projectId: '550e8400-e29b-41d4-a716-446655440002',
      moduleId: '550e8400-e29b-41d4-a716-446655440003',
      status: 'pending',
      prompt: '测试用户注册功能',
      result: null,
      createdAt: '2026-04-20T10:00:00Z',
      updatedAt: '2026-04-20T10:00:00Z',
    })
  ),

  // GET /api/v1/generation/tasks/:id - 获取任务详情
  http.get('/api/v1/generation/tasks/:id', ({ params }) => {
    const { id } = params
    const tasks: Record<string, GenerationTask> = {
      '550e8400-e29b-41d4-a716-446655440001': {
        id: '550e8400-e29b-41d4-a716-446655440001',
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        moduleId: '550e8400-e29b-41d4-a716-446655440003',
        status: 'pending',
        prompt: '测试用户注册功能',
        result: null,
        createdAt: '2026-04-20T10:00:00Z',
        updatedAt: '2026-04-20T10:00:00Z',
      },
      '550e8400-e29b-41d4-a716-446655440004': {
        id: '550e8400-e29b-41d4-a716-446655440004',
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        moduleId: '550e8400-e29b-41d4-a716-446655440003',
        status: 'completed',
        prompt: '测试用户登录功能',
        result: { draftCount: 5, confidence: 'high' },
        createdAt: '2026-04-20T09:00:00Z',
        updatedAt: '2026-04-20T09:05:00Z',
      },
      '550e8400-e29b-41d4-a716-446655440005': {
        id: '550e8400-e29b-41d4-a716-446655440005',
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        moduleId: '550e8400-e29b-41d4-a716-446655440003',
        status: 'processing',
        prompt: '测试用户登录功能',
        result: null,
        createdAt: '2026-04-20T09:00:00Z',
        updatedAt: '2026-04-20T09:01:00Z',
      },
    }
    const task = tasks[id as string]
    if (!task) {
      return HttpResponse.json({ error: '任务不存在' }, { status: 404 })
    }
    return HttpResponse.json<GenerationTask>(task)
  }),

  // GET /api/v1/generation/tasks - 任务列表（通过 query params）
  http.get('/api/v1/generation/tasks', ({ request }) => {
    const url = new URL(request.url)
    const status = url.searchParams.get('status')

    const allTasks: GenerationTask[] = [
      {
        id: '550e8400-e29b-41d4-a716-446655440001',
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        moduleId: '550e8400-e29b-41d4-a716-446655440003',
        status: 'pending',
        prompt: '测试用户注册功能',
        result: null,
        createdAt: '2026-04-20T10:00:00Z',
        updatedAt: '2026-04-20T10:00:00Z',
      },
      {
        id: '550e8400-e29b-41d4-a716-446655440004',
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        moduleId: '550e8400-e29b-41d4-a716-446655440003',
        status: 'completed',
        prompt: '测试用户登录功能',
        result: { draftCount: 5, confidence: 'high' },
        createdAt: '2026-04-20T09:00:00Z',
        updatedAt: '2026-04-20T09:05:00Z',
      },
    ]

    // Filter by status if provided
    const filteredTasks = status
      ? allTasks.filter((task) => task.status === status)
      : allTasks

    return HttpResponse.json<{
      data: GenerationTask[]
      total: number
      offset: number
      limit: number
    }>({
      data: filteredTasks,
      total: filteredTasks.length,
      offset: 0,
      limit: 10,
    })
  }),

  // GET /api/v1/generation/tasks/:id/drafts - 获取任务草稿列表
  http.get('/api/v1/generation/tasks/:id/drafts', () =>
    HttpResponse.json<CaseDraft[]>([
      {
        id: '550e8400-e29b-41d4-a716-446655440010',
        taskId: '550e8400-e29b-41d4-a716-446655440004',
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
        createdAt: '2026-04-20T09:02:00Z',
        updatedAt: '2026-04-20T09:02:00Z',
      },
      {
        id: '550e8400-e29b-41d4-a716-446655440011',
        taskId: '550e8400-e29b-41d4-a716-446655440004',
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
        createdAt: '2026-04-20T09:03:00Z',
        updatedAt: '2026-04-20T09:03:00Z',
      },
    ])
  ),
]
