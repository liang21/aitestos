import { beforeEach, describe, expect, it } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { createElement } from 'react'
import { server } from '../../../../tests/msw/server'
import {
  useConfigList,
  useSetConfig,
  useDeleteConfig,
  useImportConfigs,
  useExportConfigs,
} from './useConfigs'

// Helper: 创建测试用 QueryClient
function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
}

function wrapper({ children }: { children: any }) {
  return createElement(
    QueryClientProvider,
    { client: createTestQueryClient() },
    children
  )
}

describe('useConfigs hooks', () => {
  beforeEach(() => {
    server.resetHandlers()
  })

  describe('useConfigList', () => {
    it('should return config list data', async () => {
      // Arrange
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

      // Act
      const { result } = renderHook(() => useConfigList('123'), { wrapper })

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.data).toHaveLength(2)
      expect(result.current.data?.data[0].key).toBe('llm_model')
    })

    it('should handle empty config list', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/projects/123/configs', () =>
          HttpResponse.json({ data: [] })
        )
      )

      // Act
      const { result } = renderHook(() => useConfigList('123'), { wrapper })

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.data).toHaveLength(0)
    })

    it('should handle fetch error', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/projects/123/configs', () =>
          HttpResponse.json({ error: 'Unauthorized' }, { status: 401 })
        )
      )

      // Act
      const { result } = renderHook(() => useConfigList('123'), { wrapper })

      // Assert
      await waitFor(() => expect(result.current.isError).toBe(true))
    })
  })

  describe('useSetConfig', () => {
    it('should set config successfully and invalidate cache', async () => {
      // Arrange
      const mockConfig = {
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
            return HttpResponse.json(mockConfig)
          }
        )
      )

      // Act
      const { result } = renderHook(() => useSetConfig('123'), { wrapper })

      result.current.mutate({
        key: 'llm_model',
        value: 'gpt-4',
        description: 'Updated model',
      })

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.key).toBe('llm_model')
    })

    it('should handle set config error', async () => {
      // Arrange
      server.use(
        http.put('/api/v1/projects/123/configs/llm_model', () =>
          HttpResponse.json({ error: 'Invalid config' }, { status: 400 })
        )
      )

      // Act
      const { result } = renderHook(() => useSetConfig('123'), { wrapper })

      result.current.mutate({
        key: 'llm_model',
        value: 'invalid-model',
        description: 'Invalid',
      })

      // Assert
      await waitFor(() => expect(result.current.isError).toBe(true))
    })
  })

  describe('useDeleteConfig', () => {
    it('should delete config successfully and invalidate cache', async () => {
      // Arrange
      server.use(
        http.delete(
          '/api/v1/projects/123/configs/llm_model',
          () => new HttpResponse(null, { status: 204 })
        )
      )

      // Act
      const { result } = renderHook(() => useDeleteConfig('123'), { wrapper })

      result.current.mutate('llm_model')

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
    })

    it('should handle delete config error', async () => {
      // Arrange
      server.use(
        http.delete('/api/v1/projects/123/configs/llm_model', () =>
          HttpResponse.json({ error: 'Config not found' }, { status: 404 })
        )
      )

      // Act
      const { result } = renderHook(() => useDeleteConfig('123'), { wrapper })

      result.current.mutate('llm_model')

      // Assert
      await waitFor(() => expect(result.current.isError).toBe(true))
    })
  })

  describe('useImportConfigs', () => {
    it('should import configs successfully', async () => {
      // Arrange
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

      // Act
      const { result } = renderHook(() => useImportConfigs('123'), { wrapper })

      result.current.mutate(configsToImport)

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.successCount).toBe(2)
    })

    it('should handle partial import failure', async () => {
      // Arrange
      const configsToImport = [
        { key: 'model', value: 'gpt-4' },
        { key: 'temperature', value: 0.7 },
      ]

      server.use(
        http.post('/api/v1/projects/123/configs/import', () =>
          HttpResponse.json({
            successCount: 1,
            failedCount: 1,
            errors: [{ key: 'temperature', error: 'Invalid value type' }],
          })
        )
      )

      // Act
      const { result } = renderHook(() => useImportConfigs('123'), { wrapper })

      result.current.mutate(configsToImport)

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.failedCount).toBe(1)
      expect(result.current.data?.errors).toHaveLength(1)
    })

    it('should handle import error', async () => {
      // Arrange
      server.use(
        http.post('/api/v1/projects/123/configs/import', () =>
          HttpResponse.json({ error: 'Import failed' }, { status: 500 })
        )
      )

      // Act
      const { result } = renderHook(() => useImportConfigs('123'), { wrapper })

      result.current.mutate([{ key: 'model', value: 'gpt-4' }])

      // Assert
      await waitFor(() => expect(result.current.isError).toBe(true))
    })
  })

  describe('useExportConfigs', () => {
    it('should export configs successfully', async () => {
      // Arrange
      const mockConfigs = [
        { key: 'llm_model', value: 'deepseek-chat' },
        { key: 'llm_temperature', value: 0.7 },
      ]

      server.use(
        http.get('/api/v1/projects/123/configs/export', () =>
          HttpResponse.json({ configs: mockConfigs })
        )
      )

      // Act
      const { result } = renderHook(() => useExportConfigs('123'), { wrapper })

      result.current.mutate()

      // Assert
      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.configs).toHaveLength(2)
      expect(result.current.data?.configs[0].key).toBe('llm_model')
    })

    it('should handle export error', async () => {
      // Arrange
      server.use(
        http.get('/api/v1/projects/123/configs/export', () =>
          HttpResponse.json({ error: 'Export failed' }, { status: 500 })
        )
      )

      // Act
      const { result } = renderHook(() => useExportConfigs('123'), { wrapper })

      result.current.mutate()

      // Assert
      await waitFor(() => expect(result.current.isError).toBe(true))
    })
  })
})
