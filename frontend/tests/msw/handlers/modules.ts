import { http, HttpResponse } from 'msw'
import type { Module } from '@/types/api'

export const modulesHandlers = [
  // GET /api/v1/projects/:projectId/modules - 获取项目的模块列表
  http.get('/api/v1/projects/:projectId/modules', () =>
    HttpResponse.json<{
      data: Module[]
      total: number
      offset: number
      limit: number
    }>({
      data: [
        {
          id: 'mod-1',
          projectId: 'proj-1',
          name: '用户中心',
          abbreviation: 'USR',
          createdAt: '2026-04-20T10:00:00Z',
          updatedAt: '2026-04-20T10:00:00Z',
        },
        {
          id: 'mod-2',
          projectId: 'proj-1',
          name: '订单管理',
          abbreviation: 'ORD',
          createdAt: '2026-04-20T10:00:00Z',
          updatedAt: '2026-04-20T10:00:00Z',
        },
      ],
      total: 2,
      offset: 0,
      limit: 10,
    })
  ),
]
