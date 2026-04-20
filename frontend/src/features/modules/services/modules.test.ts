import { afterAll, describe, expect, it } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { modulesApi } from './modules'

describe('modulesApi', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  describe('list', () => {
    it('should call GET /modules with projectId param', async () => {
      const mockData = {
        data: [
          { id: '1', projectId: 'proj1', name: 'User Module', abbreviation: 'USR', createdAt: '2024-01-01', updatedAt: '2024-01-01' },
          { id: '2', projectId: 'proj1', name: 'Order Module', abbreviation: 'ORD', createdAt: '2024-01-01', updatedAt: '2024-01-01' },
        ],
        total: 2,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get('/api/v1/projects/proj1/modules', () => HttpResponse.json(mockData))
      )

      const result = await modulesApi.list('proj1')
      expect(result.data).toHaveLength(2)
      expect(result.data[0].name).toBe('User Module')
    })

    it('should return modules list for project', async () => {
      const mockData = {
        data: [
          { id: '1', projectId: 'proj1', name: 'Module A', abbreviation: 'MDA', createdAt: '2024-01-01', updatedAt: '2024-01-01' },
        ],
        total: 1,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get('/api/v1/projects/proj1/modules', () => HttpResponse.json(mockData))
      )

      const result = await modulesApi.list('proj1')
      expect(result.data[0].abbreviation).toBe('MDA')
    })
  })

  describe('create', () => {
    it('should call POST /projects/:projectId/modules with data', async () => {
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
        http.post('/api/v1/projects/proj1/modules', async ({ request }) => {
          const body = await request.json()
          expect(body).toEqual(newModule)
          return HttpResponse.json(mockResponse, { status: 201 })
        })
      )

      const result = await modulesApi.create('proj1', newModule)
      expect(result.id).toBe('456')
      expect(result.name).toBe('PaymentModule')
    })
  })

  describe('delete', () => {
    it('should call DELETE /modules/:id', async () => {
      server.use(
        http.delete('/api/v1/modules/123', () => new HttpResponse(null, { status: 204 }))
      )

      await expect(modulesApi.delete('123')).resolves.toBe('')
    })
  })
})
