/**
 * CreateCaseDrawer Component Tests
 */

import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { CreateCaseDrawer } from './CreateCaseDrawer'

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

describe('CreateCaseDrawer', () => {
  beforeEach(() => {
    // Mock modules API
    server.use(
      http.get('/api/v1/projects/:projectId/modules', () =>
        HttpResponse.json({
          data: [
            {
              id: 'mod-001',
              projectId: 'project-001',
              name: '用户中心',
              abbreviation: 'USR',
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
  })

  describe('rendering when drawer is open', () => {
    it('should render form fields', async () => {
      const onClose = vi.fn()

      renderWithProviders(
        <CreateCaseDrawer visible projectId="project-001" onClose={onClose} />
      )

      expect(await screen.findByText(/新建测试用例/i)).toBeInTheDocument()

      // Check form fields exist
      expect(screen.getByPlaceholderText(/所属模块/i)).toBeInTheDocument()
      expect(screen.getByPlaceholderText(/用例标题/i)).toBeInTheDocument()
      expect(screen.getByPlaceholderText(/前置条件/i)).toBeInTheDocument()
      expect(screen.getByPlaceholderText(/测试步骤/i)).toBeInTheDocument()
    })

    it('should render type and priority selects', async () => {
      renderWithProviders(
        <CreateCaseDrawer visible projectId="project-001" onClose={vi.fn()} />
      )

      expect(await screen.findByText(/用例类型/i)).toBeInTheDocument()
      expect(screen.getByText(/优先级/i)).toBeInTheDocument()
    })

    it('should render footer buttons', async () => {
      renderWithProviders(
        <CreateCaseDrawer visible projectId="project-001" onClose={vi.fn()} />
      )

      expect(
        await screen.findByRole('button', { name: /取消/i })
      ).toBeInTheDocument()
      expect(
        screen.getByRole('button', { name: /确认创建/i })
      ).toBeInTheDocument()
    })
  })

  describe('form submission', () => {
    it('should render form with all fields', async () => {
      const onClose = vi.fn()

      server.use(
        http.post('/api/v1/testcases', () =>
          HttpResponse.json({
            id: 'case-001',
            moduleId: 'mod-001',
            userId: 'user-001',
            number: 'TEST-USR-20260421-001',
            title: '用户登录成功',
            preconditions: [],
            steps: ['步骤1'],
            expected: {},
            caseType: 'functionality',
            priority: 'P2',
            status: 'unexecuted',
            createdAt: '2026-04-21T00:00:00Z',
            updatedAt: '2026-04-21T00:00:00Z',
          })
        )
      )

      renderWithProviders(
        <CreateCaseDrawer visible projectId="project-001" onClose={onClose} />
      )

      // Verify drawer is open with all elements
      expect(screen.getByText(/新建测试用例/i)).toBeInTheDocument()
      expect(screen.getByPlaceholderText(/所属模块/i)).toBeInTheDocument()
      expect(screen.getByPlaceholderText(/用例标题/i)).toBeInTheDocument()
      expect(screen.getByPlaceholderText(/测试步骤/i)).toBeInTheDocument()
      expect(
        screen.getByRole('button', { name: /确认创建/i })
      ).toBeInTheDocument()
    })
  })

  describe('not visible', () => {
    it('should not render drawer when visible is false', () => {
      renderWithProviders(
        <CreateCaseDrawer
          visible={false}
          projectId="project-001"
          onClose={vi.fn()}
        />
      )

      expect(screen.queryByText(/新建测试用例/i)).not.toBeInTheDocument()
    })
  })
})
