/**
 * MSW Handlers for TestCases API
 */

import { http, HttpResponse } from 'msw'

export const testcasesHandlers = [
  // GET /testcases - List test cases
  http.get('/api/v1/testcases', () =>
    HttpResponse.json({
      data: [
        {
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
        },
        {
          id: 'case-002',
          moduleId: 'mod-001',
          userId: 'user-001',
          number: 'TEST-USR-20260421-002',
          title: '用户登录失败',
          preconditions: ['用户已注册'],
          steps: ['打开登录页', '输入错误密码', '点击登录'],
          expected: { step_3: '显示错误提示' },
          caseType: 'functionality',
          priority: 'medium',
          status: 'passed',
          createdAt: '2026-04-21T00:00:00Z',
          updatedAt: '2026-04-21T00:00:00Z',
        },
      ],
      total: 2,
      offset: 0,
      limit: 10,
    })
  ),

  // GET /testcases/:id - Get case detail
  http.get('/api/v1/testcases/:id', () =>
    HttpResponse.json({
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
      aiMetadata: {
        generationTaskId: 'task-001',
        confidence: 'high',
        referencedChunks: [
          {
            chunkId: 'chunk-001',
            documentId: 'doc-001',
            documentTitle: '用户登录需求文档',
            similarityScore: 0.92,
          },
          {
            chunkId: 'chunk-002',
            documentId: 'doc-001',
            documentTitle: '用户登录需求文档',
            similarityScore: 0.88,
          },
        ],
        modelVersion: 'gpt-4',
        generatedAt: '2026-04-21T00:00:00Z',
      },
      createdAt: '2026-04-21T00:00:00Z',
      updatedAt: '2026-04-21T00:00:00Z',
    })
  ),

  // POST /testcases - Create case
  http.post('/api/v1/testcases', () =>
    HttpResponse.json(
      {
        id: 'case-001',
        moduleId: 'mod-001',
        userId: 'user-001',
        number: 'TEST-USR-20260421-001',
        title: '用户登录成功',
        preconditions: [],
        steps: ['步骤1'],
        expected: {},
        caseType: 'functionality',
        priority: 'high',
        status: 'unexecuted',
        createdAt: '2026-04-21T00:00:00Z',
        updatedAt: '2026-04-21T00:00:00Z',
      },
      { status: 201 }
    )
  ),

  // PUT /testcases/:id - Update case
  http.put('/api/v1/testcases/:id', () =>
    HttpResponse.json({
      id: 'case-001',
      moduleId: 'mod-001',
      userId: 'user-001',
      number: 'TEST-USR-20260421-001',
      title: '更新后的标题',
      preconditions: [],
      steps: [],
      expected: {},
      caseType: 'functionality',
      priority: 'high',
      status: 'passed',
      createdAt: '2026-04-21T00:00:00Z',
      updatedAt: '2026-04-21T00:00:00Z',
    })
  ),

  // DELETE /testcases/:id - Delete case
  http.delete('/api/v1/testcases/:id', () =>
    HttpResponse.json({ success: true })
  ),
]
