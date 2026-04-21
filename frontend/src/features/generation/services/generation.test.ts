import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { generationApi } from './generation'
import { server } from '../../../../tests/msw/server'
import { generationHandlers } from '../../../../tests/msw/handlers/generation'

describe('Generation API service', () => {
  beforeEach(() => {
    server.use(...generationHandlers)
  })

  afterEach(() => {
    server.resetHandlers()
  })

  describe('createTask', () => {
    it('should create a new generation task', async () => {
      const requestData = {
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        moduleId: '550e8400-e29b-41d4-a716-446655440003',
        prompt: '测试用户注册功能',
        count: 5,
        caseType: 'functionality' as const,
        priority: 'P1' as const,
        sceneType: 'positive' as const,
      }

      const result = await generationApi.createTask(requestData)

      expect(result).toMatchObject({
        id: expect.any(String),
        projectId: requestData.projectId,
        moduleId: requestData.moduleId,
        status: 'pending',
        prompt: requestData.prompt,
        result: null,
        createdAt: expect.any(String),
        updatedAt: expect.any(String),
      })
    })
  })

  describe('getTask', () => {
    it('should fetch task details by id', async () => {
      const taskId = '550e8400-e29b-41d4-a716-446655440001'

      const result = await generationApi.getTask(taskId)

      expect(result).toMatchObject({
        id: taskId,
        status: 'pending',
        prompt: '测试用户注册功能',
        result: null,
      })
    })

    it('should return completed task with result', async () => {
      const taskId = '550e8400-e29b-41d4-a716-446655440004'

      const result = await generationApi.getTask(taskId)

      expect(result).toMatchObject({
        id: taskId,
        status: 'completed',
        result: {
          draftCount: 5,
        },
      })
    })
  })

  describe('listTasks', () => {
    it('should fetch task list with pagination', async () => {
      const params = {
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        status: 'pending' as const,
        offset: 0,
        limit: 10,
      }

      const result = await generationApi.listTasks(params)

      expect(result).toMatchObject({
        data: expect.any(Array),
        total: expect.any(Number),
        offset: 0,
        limit: 10,
      })

      expect(result.data.length).toBeGreaterThan(0)
      expect(result.data[0]).toMatchObject({
        id: expect.any(String),
        status: 'pending',
      })
    })

    it('should filter tasks by status', async () => {
      const params = {
        projectId: '550e8400-e29b-41d4-a716-446655440002',
        status: 'completed' as const,
        offset: 0,
        limit: 10,
      }

      const result = await generationApi.listTasks(params)

      expect(result.data.every((task) => task.status === 'completed')).toBe(
        true
      )
    })
  })

  describe('getTaskDrafts', () => {
    it('should fetch drafts for a completed task', async () => {
      const taskId = '550e8400-e29b-41d4-a716-446655440004'

      const result = await generationApi.getTaskDrafts(taskId)

      expect(result).toEqual(expect.any(Array))
      expect(result.length).toBeGreaterThan(0)
      expect(result[0]).toMatchObject({
        id: expect.any(String),
        taskId: taskId,
        title: expect.any(String),
        preconditions: expect.any(Array),
        steps: expect.any(Array),
        expected: expect.any(Object),
        caseType: expect.any(String),
        priority: expect.any(String),
        status: 'pending',
      })
    })

    it('should return drafts with correct structure', async () => {
      const taskId = '550e8400-e29b-41d4-a716-446655440004'

      const result = await generationApi.getTaskDrafts(taskId)

      expect(result[0]).toMatchObject({
        id: '550e8400-e29b-41d4-a716-446655440010',
        title: '验证有效邮箱注册',
        caseType: 'functionality',
        priority: 'P1',
        steps: expect.arrayContaining(['打开注册页面', '输入有效邮箱']),
      })
    })
  })
})
