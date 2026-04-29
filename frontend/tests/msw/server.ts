import { setupServer } from 'msw/node'
import { HttpResponse, http } from 'msw'
import { documentsHandlers } from './handlers/documents'
import { draftsHandlers } from './handlers/drafts'
import { testcasesHandlers } from './handlers/testcases'
import { plansHandlers } from './handlers/plans'

export const server = setupServer(
  // Default handlers for common endpoints (with /api/v1 prefix)
  http.post('/api/v1/auth/login', () =>
    HttpResponse.json({
      access_token: 'mock-access-token',
      refresh_token: 'mock-refresh-token',
      user: {
        id: '1',
        username: 'admin',
        email: 'admin@test.com',
        role: 'admin',
      },
    })
  ),
  // Also handle without prefix for tests
  http.post('/auth/login', () =>
    HttpResponse.json({
      access_token: 'mock-access-token',
      refresh_token: 'mock-refresh-token',
      user: {
        id: '1',
        username: 'admin',
        email: 'admin@test.com',
        role: 'admin',
      },
    })
  ),
  http.post('/api/v1/auth/refresh', () =>
    HttpResponse.json({
      access_token: 'new-access-token',
      refresh_token: 'new-refresh-token',
    })
  ),
  http.post('/auth/refresh', () =>
    HttpResponse.json({
      access_token: 'new-access-token',
      refresh_token: 'new-refresh-token',
    })
  ),
  // Modules handler (with and without prefix)
  http.get('/api/v1/projects/:projectId/modules', () =>
    HttpResponse.json({
      data: [
        {
          id: 'mod-1',
          projectId: 'project-1',
          name: '用户中心',
          abbreviation: 'USR',
          createdAt: '2026-04-16T00:00:00Z',
          updatedAt: '2026-04-16T00:00:00Z',
        },
        {
          id: 'mod-2',
          projectId: 'project-1',
          name: '订单管理',
          abbreviation: 'ORD',
          createdAt: '2026-04-16T00:00:00Z',
          updatedAt: '2026-04-16T00:00:00Z',
        },
      ],
      total: 2,
      offset: 0,
      limit: 10,
    })
  ),
  http.get('/projects/:projectId/modules', () =>
    HttpResponse.json({
      data: [
        {
          id: 'mod-1',
          projectId: 'project-1',
          name: '用户中心',
          abbreviation: 'USR',
          createdAt: '2026-04-16T00:00:00Z',
          updatedAt: '2026-04-16T00:00:00Z',
        },
        {
          id: 'mod-2',
          projectId: 'project-1',
          name: '订单管理',
          abbreviation: 'ORD',
          createdAt: '2026-04-16T00:00:00Z',
          updatedAt: '2026-04-16T00:00:00Z',
        },
      ],
      total: 2,
      offset: 0,
      limit: 10,
    })
  ),
  // Documents handlers
  ...documentsHandlers,
  // Draft handlers (also add without prefix versions)
  ...draftsHandlers,
  http.get('/generation/drafts/count', ({ request }) => {
    console.log('[MSW] GET /generation/drafts/count called', request.url)
    const response = HttpResponse.json({ count: 2 })
    console.log('[MSW] Response:', response)
    return response
  }),
  http.get('/generation/drafts', ({ request }) => {
    console.log('[MSW] GET /generation/drafts called', request.url)
    return HttpResponse.json({
      data: [],
      total: 0,
      offset: 0,
      limit: 10,
    })
  }),
  http.get('/generation/drafts/:id', ({ params }) => {
    console.log('[MSW] GET /generation/drafts/:id called', params.id)
    return HttpResponse.json({
      id: 'draft-001',
      taskId: 'task-001',
      projectId: 'project-1',
      title: '验证有效邮箱注册',
      preconditions: [],
      steps: [],
      expected: {},
      caseType: 'functionality' as const,
      priority: 'P1' as const,
      status: 'pending' as const,
      createdAt: '2026-04-20T10:00:00Z',
      updatedAt: '2026-04-20T10:00:00Z',
    })
  }),
  http.post('/generation/drafts/:id/confirm', () => {
    return HttpResponse.json({
      id: 'tc-001',
      number: 'ECO-USR-20260421-001',
      status: 'unexecuted' as const,
    })
  }),
  http.post('/generation/drafts/:id/reject', () => {
    return HttpResponse.json({
      success: true,
      message: '草稿已拒绝',
    })
  }),
  http.post('/generation/drafts/batch-confirm', () => {
    return HttpResponse.json({
      successCount: 2,
      failedCount: 0,
    })
  }),
  // Testcases handlers
  ...testcasesHandlers,
  // Plans handlers
  ...plansHandlers
)
