import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import {
  afterAll,
  afterEach,
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from 'vitest'
import {
  useProjectList,
  useProjectDetail,
  useProjectStats,
  useCreateProject,
  useUpdateProject,
  useDeleteProject,
} from './useProjects'
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

describe('useProjects hooks', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  describe('useProjectList', () => {
    it('should return projects list', async () => {
      const mockData = {
        data: [
          {
            id: '1',
            name: 'Project A',
            prefix: 'PA',
            description: 'Desc A',
            createdAt: '2024-01-01',
            updatedAt: '2024-01-01',
          },
        ],
        total: 1,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get('/api/v1/projects', () => HttpResponse.json(mockData))
      )

      const { result } = renderHook(() => useProjectList(), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data).toEqual(mockData)
    })

    it('should support keyword search', async () => {
      const mockData = {
        data: [
          {
            id: '1',
            name: 'TestProject',
            prefix: 'TST',
            description: 'Test',
            createdAt: '2024-01-01',
            updatedAt: '2024-01-01',
          },
        ],
        total: 1,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get('/api/v1/projects', ({ request }) => {
          const url = new URL(request.url)
          expect(url.searchParams.get('keywords')).toBe('TestProject')
          return HttpResponse.json(mockData)
        })
      )

      const { result } = renderHook(
        () => useProjectList({ keywords: 'TestProject' }),
        { wrapper }
      )

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
    })
  })

  describe('useProjectDetail', () => {
    it('should fetch project by id', async () => {
      const mockProject = {
        id: '123',
        name: 'ECommerce',
        prefix: 'ECO',
        description: 'Test project',
        createdAt: '2024-01-01',
        updatedAt: '2024-01-01',
      }

      server.use(
        http.get('/api/v1/projects/123', () => HttpResponse.json(mockProject))
      )

      const { result } = renderHook(() => useProjectDetail('123'), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data).toEqual(mockProject)
    })

    it('should not fetch when id is empty', () => {
      const { result } = renderHook(() => useProjectDetail(''), { wrapper })

      expect(result.current.fetchStatus).toBe('idle')
    })
  })

  describe('useProjectStats', () => {
    it('should fetch project statistics', async () => {
      const mockStats = {
        totalCases: 100,
        passRate: 85.5,
        coverage: 92.3,
        aiGeneratedCount: 45,
        trend: [{ date: '2024-01-01', passRate: 80 }],
      }

      server.use(
        http.get('/api/v1/projects/123/stats', () =>
          HttpResponse.json(mockStats)
        )
      )

      const { result } = renderHook(() => useProjectStats('123'), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(result.current.data?.totalCases).toBe(100)
    })
  })

  describe('useCreateProject', () => {
    it('should create project and invalidate cache', async () => {
      const queryClient = createTestQueryClient()
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      const newProject = {
        name: 'NewProject',
        prefix: 'NEW',
        description: 'New project',
      }

      const mockResponse = {
        id: '456',
        ...newProject,
        createdAt: '2024-01-01',
        updatedAt: '2024-01-01',
      }

      server.use(
        http.post('/api/v1/projects', async () =>
          HttpResponse.json(mockResponse, { status: 201 })
        )
      )

      const { result } = renderHook(() => useCreateProject(), {
        wrapper: ({ children }) => (
          <QueryClientProvider client={queryClient}>
            {children}
          </QueryClientProvider>
        ),
      })

      await result.current.mutateAsync(newProject)

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['projects', 'list'],
      })
    })
  })

  describe('useUpdateProject', () => {
    it('should update project and invalidate detail cache', async () => {
      const queryClient = createTestQueryClient()
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      const updateData = { name: 'UpdatedProject' }

      const mockResponse = {
        id: '123',
        prefix: 'ECO',
        ...updateData,
        description: 'Desc',
        createdAt: '2024-01-01',
        updatedAt: '2024-01-02',
      }

      server.use(
        http.put('/api/v1/projects/123', async () =>
          HttpResponse.json(mockResponse)
        )
      )

      const { result } = renderHook(() => useUpdateProject(), {
        wrapper: ({ children }) => (
          <QueryClientProvider client={queryClient}>
            {children}
          </QueryClientProvider>
        ),
      })

      await result.current.mutateAsync({ id: '123', data: updateData })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['projects', 'detail', '123'],
      })
    })
  })

  describe('useDeleteProject', () => {
    it('should delete project and invalidate list cache', async () => {
      const queryClient = createTestQueryClient()
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      server.use(
        http.delete(
          '/api/v1/projects/123',
          () => new HttpResponse(null, { status: 204 })
        )
      )

      const { result } = renderHook(() => useDeleteProject(), {
        wrapper: ({ children }) => (
          <QueryClientProvider client={queryClient}>
            {children}
          </QueryClientProvider>
        ),
      })

      await result.current.mutateAsync('123')

      await waitFor(() => expect(result.current.isSuccess).toBe(true))
      expect(invalidateSpy).toHaveBeenCalledWith({
        queryKey: ['projects', 'list'],
      })
    })
  })
})
