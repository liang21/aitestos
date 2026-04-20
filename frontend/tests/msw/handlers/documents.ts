import { http, HttpResponse } from 'msw'

export const documentsHandlers = [
  // GET /knowledge/documents - 获取文档列表
  http.get('/api/v1/knowledge/documents', ({ request }) => {
    const url = new URL(request.url)
    const projectId = url.searchParams.get('project_id')
    const type = url.searchParams.get('type')
    const status = url.searchParams.get('status')

    return HttpResponse.json({
      data: [
        {
          id: 'doc-1',
          projectId: projectId || 'proj-1',
          name: '用户注册模块 PRD v2.0',
          type: type || 'prd',
          status: status || 'completed',
          chunkCount: 42,
          uploadedBy: 'user-1',
          uploadedByName: '张三',
          createdAt: '2026-04-16T08:00:00Z',
          updatedAt: '2026-04-16T08:30:00Z',
        },
        {
          id: 'doc-2',
          projectId: projectId || 'proj-1',
          name: 'API 接口规范',
          type: type || 'api_spec',
          status: status || 'processing',
          chunkCount: 0,
          uploadedBy: 'user-1',
          uploadedByName: '张三',
          createdAt: '2026-04-16T09:00:00Z',
          updatedAt: '2026-04-16T09:05:00Z',
        },
      ],
      total: 2,
      offset: 0,
      limit: 10,
    })
  }),

  // GET /knowledge/documents/{id} - 获取文档详情
  http.get('/api/v1/knowledge/documents/:id', ({ params }) => {
    const hasChunks = params.id === 'doc-1'
    return HttpResponse.json({
      id: params.id,
      projectId: 'proj-1',
      name: '用户注册模块 PRD v2.0',
      type: 'prd',
      status: 'completed',
      chunkCount: hasChunks ? 2 : 0,
      uploadedBy: 'user-1',
      uploadedByName: '张三',
      createdAt: '2026-04-16T08:00:00Z',
      updatedAt: '2026-04-16T08:30:00Z',
      chunks: hasChunks
        ? [
            {
              id: 'chunk-1',
              documentId: params.id,
              chunkIndex: 0,
              content: '用户注册功能概述...',
              metadata: { page: 1 },
            },
            {
              id: 'chunk-2',
              documentId: params.id,
              chunkIndex: 1,
              content: '邮箱验证规则...',
              metadata: { page: 2 },
            },
          ]
        : [],
    })
  }),

  // POST /knowledge/documents - 上传文档
  http.post('/api/v1/knowledge/documents', async ({ request }) => {
    const body = await request.json()
    return HttpResponse.json(
      {
        id: 'doc-new',
        projectId: (body as any).project_id,
        name: (body as any).name,
        type: (body as any).type,
        status: 'pending',
        chunkCount: 0,
        uploadedBy: 'user-1',
        uploadedByName: '张三',
        createdAt: '2026-04-16T10:00:00Z',
        updatedAt: '2026-04-16T10:00:00Z',
      },
      { status: 201 }
    )
  }),

  // DELETE /knowledge/documents/{id} - 删除文档
  http.delete('/api/v1/knowledge/documents/:id', () => {
    return new HttpResponse(null, { status: 204 })
  }),

  // GET /knowledge/documents/{id}/chunks - 获取文档分块列表
  http.get('/api/v1/knowledge/documents/:id/chunks', ({ params }) => {
    return HttpResponse.json([
      {
        id: 'chunk-1',
        documentId: params.id,
        chunkIndex: 0,
        content: '用户注册功能概述...',
      },
      {
        id: 'chunk-2',
        documentId: params.id,
        chunkIndex: 1,
        content: '邮箱验证规则...',
      },
    ])
  }),
]
