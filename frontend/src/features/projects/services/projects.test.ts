import { http, HttpResponse } from 'msw'
import { afterAll, afterEach, beforeEach, describe, expect, it } from 'vitest'
import { server } from '../../../../tests/msw/server'
import { projectsApi } from './projects'

describe('projectsApi', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  describe('list', () => {
    it('should call GET /projects with params', async () => {
      const mockData = {
        data: [
          { id: '1', name: 'ECommerce', prefix: 'ECO', description: 'Test project', createdAt: '2024-01-01', updatedAt: '2024-01-01' },
        ],
        total: 1,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get('/api/v1/projects', ({ request }) => {
          const url = new URL(request.url)
          expect(url.searchParams.get('keywords')).toBe('test')
          expect(url.searchParams.get('offset')).toBe('0')
          expect(url.searchParams.get('limit')).toBe('10')
          return HttpResponse.json(mockData)
        })
      )

      const result = await projectsApi.list({ keywords: 'test', offset: 0, limit: 10 })
      expect(result).toEqual(mockData)
    })

    it('should return projects list', async () => {
      const mockData = {
        data: [
          { id: '1', name: 'Project A', prefix: 'PA', description: 'Desc A', createdAt: '2024-01-01', updatedAt: '2024-01-01' },
          { id: '2', name: 'Project B', prefix: 'PB', description: 'Desc B', createdAt: '2024-01-02', updatedAt: '2024-01-02' },
        ],
        total: 2,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get('/api/v1/projects', () => HttpResponse.json(mockData))
      )

      const result = await projectsApi.list()
      expect(result.data).toHaveLength(2)
      expect(result.data[0].name).toBe('Project A')
    })
  })

  describe('get', () => {
    it('should call GET /projects/:id', async () => {
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

      const result = await projectsApi.get('123')
      expect(result).toEqual(mockProject)
    })
  })

  describe('create', () => {
    it('should call POST /projects with data', async () => {
      const newProject = {
        name: 'NewProject',
        prefix: 'NEW',
        description: 'New project description',
      }

      const mockResponse = {
        id: '456',
        ...newProject,
        createdAt: '2024-01-01',
        updatedAt: '2024-01-01',
      }

      server.use(
        http.post('/api/v1/projects', async ({ request }) => {
          const body = await request.json()
          expect(body).toEqual(newProject)
          return HttpResponse.json(mockResponse, { status: 201 })
        })
      )

      const result = await projectsApi.create(newProject)
      expect(result.id).toBe('456')
      expect(result.name).toBe('NewProject')
    })
  })

  describe('update', () => {
    it('should call PUT /projects/:id with data', async () => {
      const updateData = {
        name: 'UpdatedProject',
        description: 'Updated description',
      }

      const mockResponse = {
        id: '123',
        prefix: 'ECO',
        ...updateData,
        createdAt: '2024-01-01',
        updatedAt: '2024-01-02',
      }

      server.use(
        http.put('/api/v1/projects/123', async ({ request }) => {
          const body = await request.json()
          expect(body).toEqual(updateData)
          return HttpResponse.json(mockResponse)
        })
      )

      const result = await projectsApi.update('123', updateData)
      expect(result.name).toBe('UpdatedProject')
    })
  })

  describe('delete', () => {
    it('should call DELETE /projects/:id', async () => {
      server.use(
        http.delete('/api/v1/projects/123', () => new HttpResponse(null, { status: 204 }))
      )

      // DELETE returns 204, axios returns empty string
      await expect(projectsApi.delete('123')).resolves.toBe('')
    })
  })

  describe('getStats', () => {
    it('should call GET /projects/:id/stats', async () => {
      const mockStats = {
        totalCases: 100,
        passRate: 85.5,
        coverage: 92.3,
        aiGeneratedCount: 45,
        trend: [
          { date: '2024-01-01', passRate: 80 },
          { date: '2024-01-02', passRate: 85 },
        ],
      }

      server.use(
        http.get('/api/v1/projects/123/stats', () => HttpResponse.json(mockStats))
      )

      const result = await projectsApi.getStats('123')
      expect(result.totalCases).toBe(100)
      expect(result.passRate).toBe(85.5)
    })
  })
})
