import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { useModuleList, useCreateModule, useDeleteModule } from './useModules'
import React from 'react'

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
}

function wrapper({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={createTestQueryClient()}>
      {children}
    </QueryClientProvider>
  )
}

describe('useModules hooks', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  describe('useModuleList', () => {
    it('should fetch modules by projectId', async () => {
      const mockData = {
        data: [
          {
            id: '1',
            projectId: 'proj1',
            name: 'User Module',
            abbreviation: 'USR',
            createdAt: '2024-01-01',
            updatedAt: '2024-01-01',
          },
          {
            id: '2',
            projectId: 'proj1',
            name: 'Order Module',
            abbreviation: 'ORD',
            createdAt: '2024-01-01',
            updatedAt: '2024-01-01',
          },
        ],
        total: 2,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get('/api/v1/projects/proj1/modules', () =>
          HttpResponse.json(mockData)
        )
      )

      const { result } = renderHook(() => useModuleList('proj1'), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.data).toHaveLength(2)
      expect(result.current.data?.data[0].name).toBe('User Module')
    })

    it('should not fetch when projectId is empty', () => {
      const { result } = renderHook(() => useModuleList(''), { wrapper })

      expect(result.current.fetchStatus).toBe('idle')
    })
  })

  describe('useCreateModule', () => {
    it('should create module and invalidate cache', async () => {
      const queryClient = createTestQueryClient()
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      const newModule = {
        name: 'PaymentModule',
        abbreviation: 'PAY',
      }

      const mockResponse = {
        id: '456',
        projectId: 'proj1',
        ...newModule,
        createdAt: '2024-01-01',
        updatedAt: '2024-01-01',
      }

      server.use(
        http.post('/api/v1/projects/proj1/modules', async () =>
          HttpResponse.json(mockResponse, { status: 201 })
        )
      )

      const { result } = renderHook(() => useCreateModule(), {
        wrapper: ({ children }) => (
          <QueryClientProvider client={queryClient}>
            {children}
          </QueryClientProvider>
        ),
      })

      await result.current.mutateAsync({ projectId: 'proj1', data: newModule })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['modules', 'list', 'proj1'],
      })
    })
  })

  describe('useDeleteModule', () => {
    it('should delete module and invalidate cache', async () => {
      const queryClient = createTestQueryClient()
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      server.use(
        http.delete(
          '/api/v1/modules/123',
          () => new HttpResponse(null, { status: 204 })
        )
      )

      const { result } = renderHook(() => useDeleteModule(), {
        wrapper: ({ children }) => (
          <QueryClientProvider client={queryClient}>
            {children}
          </QueryClientProvider>
        ),
      })

      await result.current.mutateAsync({ projectId: 'proj1', id: '123' })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['modules', 'list', 'proj1'],
      })
    })
  })
})
