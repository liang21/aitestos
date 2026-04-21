import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { draftsApi } from './drafts'
import { server } from '../../../../tests/msw/server'
import { draftsHandlers } from '../../../../tests/msw/handlers/drafts'
import { http, HttpResponse } from 'msw'

describe('Drafts API service', () => {
  beforeEach(() => {
    server.use(...draftsHandlers)
  })

  afterEach(() => {
    server.resetHandlers()
  })

  describe('getDrafts', () => {
    it('should fetch drafts with pagination', async () => {
      const params = {
        projectId: 'project-1',
        offset: 0,
        limit: 10,
      }

      const result = await draftsApi.getDrafts(params)

      expect(result).toMatchObject({
        data: expect.any(Array),
        total: expect.any(Number),
        offset: 0,
        limit: 10,
      })

      expect(result.data.length).toBeGreaterThan(0)
      expect(result.data[0]).toMatchObject({
        id: expect.any(String),
        title: expect.any(String),
        status: 'pending',
      })
    })

    it('should filter drafts by status', async () => {
      const params = {
        projectId: 'project-1',
        status: 'confirmed' as const,
        offset: 0,
        limit: 10,
      }

      const result = await draftsApi.getDrafts(params)

      expect(result.data.every((d) => d.status === 'confirmed')).toBe(true)
    })

    it('should filter drafts by projectId', async () => {
      const params = {
        projectId: 'project-1',
        offset: 0,
        limit: 10,
      }

      const result = await draftsApi.getDrafts(params)

      expect(result.data[0].projectName).toBe('ECommerce')
    })
  })

  describe('confirmDraft', () => {
    it('should confirm draft and convert to test case', async () => {
      const draftId = 'draft-001'

      const result = await draftsApi.confirmDraft(draftId, {
        moduleId: 'mod-001',
      })

      expect(result).toMatchObject({
        id: expect.any(String),
        number: expect.any(String),
        title: '验证有效邮箱注册',
        status: 'unexecuted',
      })
    })
  })

  describe('rejectDraft', () => {
    it('should reject draft with reason', async () => {
      const draftId = 'draft-001'

      const result = await draftsApi.rejectDraft(draftId, {
        reason: 'duplicate',
        feedback: '与现有用例重复',
      })

      expect(result).toMatchObject({
        success: true,
        message: '草稿已拒绝',
      })
    })
  })

  describe('batchConfirm', () => {
    it('should confirm multiple drafts at once', async () => {
      const result = await draftsApi.batchConfirm({
        draftIds: ['draft-001', 'draft-002'],
        moduleId: 'mod-001',
      })

      expect(result).toMatchObject({
        successCount: 2,
        failedCount: 0,
      })
    })

    it('should return partial success if some confirmations fail', async () => {
      // Configure MSW to return partial success
      server.use(
        http.post('/api/v1/generation/drafts/batch-confirm', async () =>
          HttpResponse.json({
            successCount: 1,
            failedCount: 1,
            errors: [{ draftId: 'draft-002', error: '验证失败' }],
          })
        )
      )

      const result = await draftsApi.batchConfirm({
        draftIds: ['draft-001', 'draft-002'],
        moduleId: 'mod-001',
      })

      expect(result.successCount).toBe(1)
      expect(result.failedCount).toBe(1)
      expect(result.errors).toBeDefined()
    })
  })
})
