import { setupServer } from 'msw/node'
import { HttpResponse, http } from 'msw'
import { draftsHandlers } from './handlers/drafts'

export const server = setupServer(
  // Default handlers for common endpoints
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
  http.post('/api/v1/auth/refresh', () =>
    HttpResponse.json({
      access_token: 'new-access-token',
      refresh_token: 'new-refresh-token',
    })
  ),
  // Modules handler for drafts module tests
  http.get('/api/v1/projects/:projectId/modules', () =>
    HttpResponse.json({
      data: [
        { id: 'mod-1', projectId: 'project-1', name: '用户中心', abbreviation: 'USR', createdAt: '2026-04-16T00:00:00Z', updatedAt: '2026-04-16T00:00:00Z' },
        { id: 'mod-2', projectId: 'project-1', name: '订单管理', abbreviation: 'ORD', createdAt: '2026-04-16T00:00:00Z', updatedAt: '2026-04-16T00:00:00Z' },
      ],
      total: 2,
      offset: 0,
      limit: 10,
    })
  ),
  // Draft handlers
  ...draftsHandlers
)
