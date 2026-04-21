/**
 * TestCases API Service Tests
 */

import { describe, it, expect, afterEach } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { testcasesApi } from './testcases'
import type {
  TestCase,
  CreateTestCaseRequest,
  UpdateTestCaseRequest,
} from '@/types/api'

describe('testcasesApi', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  const mockData: TestCase[] = [
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
  ]

  describe('list', () => {
    it('should fetch test cases list with params', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: mockData,
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      const result = await testcasesApi.list({
        projectId: 'project-001',
        offset: 0,
        limit: 10,
      })

      expect(result.data).toEqual(mockData)
      expect(result.total).toBe(1)
    })

    it('should filter by status', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        )
      )

      const result = await testcasesApi.list({
        projectId: 'project-001',
        status: 'unexecuted',
      })

      expect(result.data).toBeDefined()
    })

    it('should filter by caseType and priority', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: mockData,
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      const result = await testcasesApi.list({
        projectId: 'project-001',
        caseType: 'functionality',
        priority: 'high',
      })

      expect(result.data).toBeDefined()
    })

    it('should search by keywords', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        )
      )

      const result = await testcasesApi.list({
        projectId: 'project-001',
        keywords: '登录',
      })

      expect(result.data).toBeDefined()
    })
  })

  describe('get', () => {
    it('should fetch test case detail by id', async () => {
      const mockCase = {
        id: 'case-001',
        moduleId: 'mod-001',
        userId: 'user-001',
        number: 'TEST-USR-20260421-001',
        title: '用户登录成功',
        preconditions: ['用户已注册'],
        steps: ['打开登录页', '输入用户名密码', '点击登录'],
        expected: { step_3: '登录成功，跳转首页' },
        caseType: 'functionality' as const,
        priority: 'high' as const,
        status: 'unexecuted' as const,
        aiMetadata: {
          generationTaskId: 'task-001',
          confidence: 'high',
          referencedChunks: [],
          modelVersion: 'gpt-4',
          generatedAt: '2026-04-21T00:00:00Z',
        },
        createdAt: '2026-04-21T00:00:00Z',
        updatedAt: '2026-04-21T00:00:00Z',
      }

      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      const result = await testcasesApi.get('case-001')

      expect(result.id).toBe('case-001')
      expect(result.aiMetadata).toBeDefined()
      expect(result.aiMetadata?.confidence).toBe('high')
    })
  })

  describe('create', () => {
    it('should create a new test case', async () => {
      const createData: CreateTestCaseRequest = {
        moduleId: 'mod-001',
        title: '用户登录成功',
        preconditions: ['用户已注册'],
        steps: ['打开登录页', '输入用户名密码', '点击登录'],
        expected: { step_3: '登录成功，跳转首页' },
        caseType: 'functionality',
        priority: 'high',
      }

      const mockResponse: TestCase = {
        id: 'case-001',
        moduleId: 'mod-001',
        userId: 'user-001',
        number: 'TEST-USR-20260421-001',
        ...createData,
        status: 'unexecuted',
        createdAt: '2026-04-21T00:00:00Z',
        updatedAt: '2026-04-21T00:00:00Z',
      }

      server.use(
        http.post('/api/v1/testcases', () =>
          HttpResponse.json(mockResponse, { status: 201 })
        )
      )

      const result = await testcasesApi.create(createData)

      expect(result.id).toBe('case-001')
      expect(result.number).toMatch(/^TEST-USR-\d{8}-\d{3}$/)
    })
  })

  describe('update', () => {
    it('should update test case', async () => {
      const updateData: UpdateTestCaseRequest = {
        title: '更新后的标题',
        status: 'passed',
      }

      const mockResponse: TestCase = {
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
      }

      server.use(
        http.put('/api/v1/testcases/:id', () => HttpResponse.json(mockResponse))
      )

      const result = await testcasesApi.update('case-001', updateData)

      expect(result.title).toBe('更新后的标题')
      expect(result.status).toBe('passed')
    })
  })

  describe('delete', () => {
    it('should delete test case', async () => {
      server.use(
        http.delete('/api/v1/testcases/:id', () =>
          HttpResponse.json({ success: true })
        )
      )

      const result = await testcasesApi.delete('case-001')

      expect(result).toEqual({ success: true })
    })
  })
})
