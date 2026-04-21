import { http, HttpResponse } from 'msw'
import { afterEach, describe, expect, it } from 'vitest'
import { server } from '../../../../tests/msw/server'
import { configsApi } from './configs'

describe('configsApi', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  describe('list', () => {
    it('should call GET /projects/:id/configs', async () => {
      const mockConfigs = [
        {
          key: 'llm_model',
          value: 'deepseek-chat',
          description: 'LLM model selection',
        },
        {
          key: 'llm_temperature',
          value: 0.7,
          description: 'Generation temperature',
        },
      ]

      server.use(
        http.get('/api/v1/projects/123/configs', () =>
          HttpResponse.json({ data: mockConfigs })
        )
      )

      const result = await configsApi.list('123')
      expect(result).toEqual({ data: mockConfigs })
    })

    it('should return empty array when no configs exist', async () => {
      server.use(
        http.get('/api/v1/projects/123/configs', () =>
          HttpResponse.json({ data: [] })
        )
      )

      const result = await configsApi.list('123')
      expect(result.data).toHaveLength(0)
    })
  })

  describe('set', () => {
    it('should call PUT /projects/:id/configs/:key with value', async () => {
      const configData = {
        key: 'llm_model',
        value: 'gpt-4',
        description: 'Updated model',
      }

      server.use(
        http.put(
          '/api/v1/projects/123/configs/llm_model',
          async ({ request }) => {
            const body = await request.json()
            expect(body).toEqual({
              value: 'gpt-4',
              description: 'Updated model',
            })
            return HttpResponse.json(configData)
          }
        )
      )

      const result = await configsApi.set('123', 'llm_model', {
        value: 'gpt-4',
        description: 'Updated model',
      })
      expect(result.key).toBe('llm_model')
      expect(result.value).toBe('gpt-4')
    })

    it('should set config with JSON object value', async () => {
      const jsonValue = { model: 'gpt-4', temperature: 0.7 }
      const configData = {
        key: 'llm_config',
        value: jsonValue,
        description: 'LLM configuration',
      }

      server.use(
        http.put(
          '/api/v1/projects/123/configs/llm_config',
          async ({ request }) => {
            const body = await request.json()
            expect(body.value).toEqual(jsonValue)
            return HttpResponse.json(configData)
          }
        )
      )

      const result = await configsApi.set('123', 'llm_config', {
        value: jsonValue,
        description: 'LLM configuration',
      })
      expect(result.value).toEqual(jsonValue)
    })
  })

  describe('delete', () => {
    it('should call DELETE /projects/:id/configs/:key', async () => {
      server.use(
        http.delete(
          '/api/v1/projects/123/configs/llm_model',
          () => new HttpResponse(null, { status: 204 })
        )
      )

      await expect(configsApi.delete('123', 'llm_model')).resolves.toBe('')
    })
  })

  describe('import', () => {
    it('should call POST /projects/:id/configs/import with configs array', async () => {
      const configsToImport = [
        { key: 'model', value: 'gpt-4' },
        { key: 'temperature', value: 0.7 },
      ]

      server.use(
        http.post(
          '/api/v1/projects/123/configs/import',
          async ({ request }) => {
            const body = await request.json()
            expect(body).toEqual({ configs: configsToImport })
            return HttpResponse.json({
              successCount: 2,
              failedCount: 0,
            })
          }
        )
      )

      const result = await configsApi.import('123', configsToImport)
      expect(result.successCount).toBe(2)
      expect(result.failedCount).toBe(0)
    })

    it('should handle partial import failure', async () => {
      const configsToImport = [
        { key: 'model', value: 'gpt-4' },
        { key: 'temperature', value: 0.7 },
      ]

      server.use(
        http.post(
          '/api/v1/projects/123/configs/import',
          async ({ request }) => {
            const body = await request.json()
            expect(body).toEqual({ configs: configsToImport })
            return HttpResponse.json({
              successCount: 1,
              failedCount: 1,
              errors: [
                {
                  key: 'temperature',
                  error: 'Invalid value type',
                },
              ],
            })
          }
        )
      )

      const result = await configsApi.import('123', configsToImport)
      expect(result.successCount).toBe(1)
      expect(result.failedCount).toBe(1)
      expect(result.errors).toHaveLength(1)
    })
  })

  describe('export', () => {
    it('should call GET /projects/:id/configs/export', async () => {
      const mockConfigs = [
        { key: 'llm_model', value: 'deepseek-chat' },
        { key: 'llm_temperature', value: 0.7 },
      ]

      server.use(
        http.get('/api/v1/projects/123/configs/export', () =>
          HttpResponse.json({ configs: mockConfigs })
        )
      )

      const result = await configsApi.export('123')
      expect(result.configs).toEqual(mockConfigs)
      expect(result.configs).toHaveLength(2)
    })
  })
})
