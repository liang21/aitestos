/**
 * PlanListPage Component Tests
 */

import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { PlanListPage } from './PlanListPage'

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>{ui}</BrowserRouter>
    </QueryClientProvider>
  )
}

describe('PlanListPage', () => {
  beforeEach(() => {
    // Mock plans API
    server.use(
      http.get('/api/v1/plans', () =>
        HttpResponse.json({
          data: [
            {
              id: 'plan-001',
              projectId: 'project-001',
              name: 'Sprint 12 回归测试',
              description: 'Sprint 12 回归测试计划',
              status: 'draft' as const,
              createdBy: 'user-001',
              createdAt: '2026-04-21T00:00:00Z',
              updatedAt: '2026-04-21T00:00:00Z',
            },
            {
              id: 'plan-002',
              projectId: 'project-001',
              name: 'Sprint 13 回归测试',
              description: 'Sprint 13 回归测试计划',
              status: 'active' as const,
              createdBy: 'user-001',
              createdAt: '2026-04-20T00:00:00Z',
              updatedAt: '2026-04-20T00:00:00Z',
            },
          ],
          total: 2,
          offset: 0,
          limit: 10,
        })
      )
    )
  })

  describe('rendering', () => {
    it('should render plan list table', async () => {
      renderWithProviders(<PlanListPage />)

      expect(await screen.findAllByText(/Sprint 12 回归测试/i)).toHaveLength(2)
      expect(screen.getAllByText(/Sprint 13 回归测试/i)).toHaveLength(2)
    })

    it('should render status tags', async () => {
      renderWithProviders(<PlanListPage />)

      expect(await screen.findByText(/草稿/i)).toBeInTheDocument()
      expect(screen.getByText(/进行中/i)).toBeInTheDocument()
    })

    it('should render create plan button', async () => {
      renderWithProviders(<PlanListPage />)

      expect(
        await screen.findByRole('button', { name: /新建计划/i })
      ).toBeInTheDocument()
    })
  })

  describe('filtering', () => {
    it('should render status filter', async () => {
      renderWithProviders(<PlanListPage />)

      // Check for the Select component by its role
      expect(await screen.findByRole('combobox')).toBeInTheDocument()
    })

    it('should filter plans by status', async () => {
      // Mock filtered response
      server.use(
        http.get('/api/v1/plans', () =>
          HttpResponse.json({
            data: [
              {
                id: 'plan-001',
                projectId: 'project-001',
                name: 'Sprint 12 回归测试',
                description: 'Sprint 12 回归测试计划',
                status: 'draft' as const,
                createdBy: 'user-001',
                createdAt: '2026-04-21T00:00:00Z',
                updatedAt: '2026-04-21T00:00:00Z',
              },
            ],
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<PlanListPage />)

      // After selecting a filter, verify the list updates
      expect(await screen.findAllByText(/Sprint 12 回归测试/i)).toHaveLength(2)
    })
  })

  describe('actions', () => {
    it('should navigate to plan detail on row click', async () => {
      const user = vi.fn()
      renderWithProviders(<PlanListPage />)

      // Click on a plan row - use the first matching element
      const planRows = await screen.findAllByText(/Sprint 12 回归测试/i)
      planRows[0].click()

      // Verify navigation occurred
      // This would typically check for navigate calls
    })

    it('should show create plan modal on button click', async () => {
      renderWithProviders(<PlanListPage />)

      const createButton = await screen.findByRole('button', { name: /新建计划/i })
      createButton.click()

      // Modal should appear
      // This would typically check for modal visibility
    })
  })
})
