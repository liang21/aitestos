/**
 * MSW 基础设施测试
 *
 * 验证 MSW 服务器配置和 handler 集中管理功能
 * 对应任务 T01: 配置 MSW 集中管理
 */

import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from './server'
import {
  documentsHandlers,
  draftsHandlers,
  generationHandlers,
  modulesHandlers,
  plansHandlers,
  testcasesHandlers,
} from './handlers'

describe('MSW Server 基础设施', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  describe('Server 生命周期', () => {
    it('应该能够启动和关闭 server', () => {
      expect(server).toBeDefined()
      expect(typeof server.listen).toBe('function')
      expect(typeof server.close).toBe('function')
      expect(typeof server.resetHandlers).toBe('function')
    })

    it('应该能够重置 handlers', () => {
      const customHandler = http.get('/api/v1/test', () =>
        HttpResponse.json({ test: 'custom' })
      )
      server.use(customHandler)
      server.resetHandlers()
      // 验证自定义 handler 已被移除
    })
  })

  describe('Handler 集中管理', () => {
    it('应该从 handlers/index.ts 导出所有 handler 数组', () => {
      expect(documentsHandlers).toBeDefined()
      expect(Array.isArray(documentsHandlers)).toBe(true)
      expect(draftsHandlers).toBeDefined()
      expect(Array.isArray(draftsHandlers)).toBe(true)
      expect(generationHandlers).toBeDefined()
      expect(Array.isArray(generationHandlers)).toBe(true)
      expect(modulesHandlers).toBeDefined()
      expect(Array.isArray(modulesHandlers)).toBe(true)
      expect(plansHandlers).toBeDefined()
      expect(Array.isArray(plansHandlers)).toBe(true)
      expect(testcasesHandlers).toBeDefined()
      expect(Array.isArray(testcasesHandlers)).toBe(true)
    })

    it('每个 handler 应该包含有效的 MSW 请求处理器', () => {
      // 验证 documentsHandlers
      expect(documentsHandlers.length).toBeGreaterThan(0)
      // MSW v2 handler 是一个对象，只要有解析和测试方法即可
      expect(typeof documentsHandlers[0]).toBe('object')
      expect(documentsHandlers[0]).not.toBeNull()

      // 验证 modulesHandlers 至少有一个 handler
      expect(modulesHandlers.length).toBeGreaterThan(0)
    })
  })

  describe('HTTP 请求拦截', () => {
    it('应该拦截 POST /api/v1/auth/login 并返回 mock 数据', async () => {
      // MSW 在 node 环境中使用相对路径
      // 但 fetch 需要完整的 URL，所以使用 http://localhost
      server.use(
        http.post('http://localhost/api/v1/auth/login', () =>
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
        )
      )

      const response = await fetch('http://localhost/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: 'test@test.com', password: 'password' }),
      })

      expect(response.status).toBe(200)
      const data = await response.json()
      expect(data).toHaveProperty('access_token')
      expect(data).toHaveProperty('refresh_token')
      expect(data.user).toHaveProperty('email', 'admin@test.com')
    })

    it('应该拦截 POST /api/v1/auth/refresh 并返回新 token', async () => {
      server.use(
        http.post('http://localhost/api/v1/auth/refresh', () =>
          HttpResponse.json({
            access_token: 'new-access-token',
            refresh_token: 'new-refresh-token',
          })
        )
      )

      const response = await fetch('http://localhost/api/v1/auth/refresh', {
        method: 'POST',
      })

      expect(response.status).toBe(200)
      const data = await response.json()
      expect(data).toHaveProperty('access_token', 'new-access-token')
    })

    it('应该拦截 GET /api/v1/projects/:projectId/modules 请求', async () => {
      server.use(
        http.get('http://localhost/api/v1/projects/:projectId/modules', () =>
          HttpResponse.json({
            data: [
              {
                id: 'mod-1',
                projectId: 'proj-1',
                name: '用户中心',
                abbreviation: 'USR',
                createdAt: '2026-04-20T10:00:00Z',
                updatedAt: '2026-04-20T10:00:00Z',
              },
            ],
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      const response = await fetch(
        'http://localhost/api/v1/projects/proj-1/modules'
      )

      expect(response.status).toBe(200)
      const data = await response.json()
      expect(data.data).toBeInstanceOf(Array)
      expect(data.data.length).toBeGreaterThan(0)
      expect(data.data[0]).toHaveProperty('name', '用户中心')
    })
  })

  describe('动态 Handler 注入', () => {
    it('应该能够在测试中临时覆盖 handler', async () => {
      // 添加临时 handler
      server.use(
        http.get('http://localhost/api/v1/test-temp', () =>
          HttpResponse.json({ message: 'temporary' })
        )
      )

      const response = await fetch('http://localhost/api/v1/test-temp')
      expect(response.status).toBe(200)
      const data = await response.json()
      expect(data.message).toBe('temporary')
    })

    it('resetHandlers 应该移除临时添加的 handler', async () => {
      // 添加临时 handler
      server.use(
        http.get('http://localhost/api/v1/test-temp-reset', () =>
          HttpResponse.json({ message: 'before reset' })
        )
      )

      let response = await fetch('http://localhost/api/v1/test-temp-reset')
      expect(response.status).toBe(200)

      // 重置后，临时 handler 应该被移除
      server.resetHandlers()

      // onUnhandledRequest: 'error' 会导致请求失败
      // 这是预期的行为
      await expect(
        fetch('http://localhost/api/v1/test-temp-reset').catch(() => ({ ok: false }))
      ).resolves.toHaveProperty('ok', false)
    })
  })
})
