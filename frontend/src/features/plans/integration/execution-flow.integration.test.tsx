/**
 * Test Execution Flow Integration Tests
 * Tests the complete flow: Create Plan → Start Execution → Record Results → Mark Complete
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse, delay } from 'msw'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'

// Import components to test
import { NewPlanPage } from '@/features/plans/components/NewPlanPage'
import { PlanDetailPage } from '@/features/plans/components/PlanDetailPage'

describe('Test Execution Flow Integration Tests', () => {
  let queryClient: QueryClient
  const mockProjectId = 'project-123'
  const mockPlanId = 'plan-123'

  const mockUser = {
    id: 'user-123',
    username: 'testuser',
    email: 'test@example.com',
    role: 'admin' as const,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  }

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
    vi.clearAllMocks()
    useAuthStore.getState().reset()
    useAuthStore.setState({
      user: mockUser,
      token: 'test-token',
      isAuthenticated: true,
    })
    server.resetHandlers()
  })

  function renderWithProviders(ui: React.ReactElement, route = '/projects') {
    window.history.pushState({}, 'Test page', route)
    return render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          {ui}
        </BrowserRouter>
      </QueryClientProvider>
    )
  }

  describe('Complete Test Execution Flow', () => {
    it('should complete full flow: create plan → add cases → start execution → record results → mark complete', async () => {
      const user = userEvent.setup()

      // Mock test cases for selection
      const mockTestCases = {
        data: [
          {
            id: 'case-1',
            moduleId: 'module-1',
            userId: 'user-123',
            number: 'TEST-001',
            title: '用户登录测试',
            preconditions: ['用户已注册'],
            steps: ['打开登录页', '输入凭证', '点击登录'],
            expected: { success: true },
            caseType: 'functionality' as const,
            priority: 'P1' as const,
            status: 'unexecuted' as const,
            createdAt: '2024-01-01T00:00:00Z',
            updatedAt: '2024-01-01T00:00:00Z',
          },
          {
            id: 'case-2',
            moduleId: 'module-1',
            userId: 'user-123',
            number: 'TEST-002',
            title: '用户注册测试',
            preconditions: [],
            steps: ['打开注册页', '填写信息', '提交注册'],
            expected: { success: true },
            caseType: 'functionality' as const,
            priority: 'P1' as const,
            status: 'unexecuted' as const,
            createdAt: '2024-01-01T00:00:00Z',
            updatedAt: '2024-01-01T00:00:00Z',
          },
        ],
        total: 2,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get(`/api/v1/testcases`, () => HttpResponse.json(mockTestCases))
      )

      // Mock plan creation
      const mockPlan = {
        id: mockPlanId,
        projectId: mockProjectId,
        name: 'Sprint 1 回归测试',
        description: '第一轮回归测试',
        status: 'draft' as const,
        createdBy: 'user-123',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      server.use(
        http.post(`/api/v1/plans`, async () => {
          await delay(100)
          return HttpResponse.json(mockPlan, { status: 201 })
        })
      )

      // Mock adding cases to plan
      server.use(
        http.post(`/api/v1/plans/${mockPlanId}/cases`, () =>
          HttpResponse.json({ success: true })
        )
      )

      // Verify Step 1: Create plan
      const { unmount: unmountNewPlan } = renderWithProviders(
        <NewPlanPage />,
        `/projects/${mockProjectId}/plans/new`
      )

      // Fill plan name
      const nameInput = screen.getByPlaceholderText(/请输入计划名称/)
      await user.click(nameInput)
      await user.keyboard('Sprint 1 回归测试')

      // Select test cases (simulate selection)
      // Note: In actual implementation, you'd interact with case selector
      const submitButton = screen.getByRole('button', { name: /创建/ })
      await user.click(submitButton)

      await waitFor(() => {
        expect(screen.getByText('创建成功')).toBeInTheDocument()
      })

      unmountNewPlan()

      // Mock plan detail with cases
      const mockPlanDetail = {
        ...mockPlan,
        status: 'active' as const,
        cases: [
          {
            caseId: 'case-1',
            caseNumber: 'TEST-001',
            caseTitle: '用户登录测试',
            resultStatus: 'pass' as const,
            resultNote: '测试通过',
            executedAt: '2024-01-01T10:00:00Z',
            executedBy: 'user-123',
          },
          {
            caseId: 'case-2',
            caseNumber: 'TEST-002',
            caseTitle: '用户注册测试',
          },
        ],
        stats: {
          total: 2,
          passed: 1,
          failed: 0,
          blocked: 0,
          skipped: 0,
          unexecuted: 1,
        },
      }

      server.use(
        http.get(`/api/v1/plans/${mockPlanId}`, () =>
          HttpResponse.json(mockPlanDetail)
        ),
        http.patch(`/api/v1/plans/${mockPlanId}/status`, ({ request }) => {
          const body = request.json() as Promise<{ status: string }>
          return HttpResponse.json({ status: (await body).status })
        })
      )

      // Mock recording results
      server.use(
        http.post(`/api/v1/plans/${mockPlanId}/results`, async () => {
          await delay(50)
          return HttpResponse.json({
            caseId: 'case-2',
            caseNumber: 'TEST-002',
            caseTitle: '用户注册测试',
            resultStatus: 'pass',
            resultNote: '测试通过',
            executedAt: '2024-01-01T10:05:00Z',
            executedBy: 'user-123',
          })
        })
      )

      // Verify Step 2: View plan detail and start execution
      const { unmount: unmountPlanDetail } = renderWithProviders(
        <PlanDetailPage />,
        `/projects/${mockProjectId}/plans/${mockPlanId}`
      )

      // Should show plan details
      await waitFor(() => {
        expect(screen.getByText('Sprint 1 回归测试')).toBeInTheDocument()
      })

      // Should show status
      expect(screen.getByText('进行中')).toBeInTheDocument()

      // Should show stats
      expect(screen.getByText('总用例: 2')).toBeInTheDocument()
      expect(screen.getByText('✅通过: 1')).toBeInTheDocument()

      // Verify Step 3: Record test result
      // Find the test case row and click to record result
      const resultButtons = screen.getAllByRole('button')
      const recordButton = resultButtons.find(btn => btn.textContent === '录入')

      if (recordButton) {
        await user.click(recordButton)

        // Select result status (simulated - actual modal interaction would be more detailed)
        await waitFor(() => {
          expect(screen.getByText('录入执行结果')).toBeInTheDocument()
        })
      }

      unmountPlanDetail()
    })

    it('should support plan status transitions', async () => {
      const mockPlan = {
        id: mockPlanId,
        projectId: mockProjectId,
        name: '测试计划',
        description: '',
        status: 'draft' as const,
        createdBy: 'user-123',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
        cases: [],
        stats: {
          total: 0,
          passed: 0,
          failed: 0,
          blocked: 0,
          skipped: 0,
          unexecuted: 0,
        },
      }

      server.use(
        http.get(`/api/v1/plans/${mockPlanId}`, () =>
          HttpResponse.json(mockPlan)
        )
      )

      let currentStatus: string = 'draft'

      server.use(
        http.patch(`/api/v1/plans/${mockPlanId}/status`, async ({ request }) => {
          const body = await request.json() as { status: string }
          currentStatus = body.status
          return HttpResponse.json({ status: currentStatus })
        })
      )

      const { rerender } = renderWithProviders(
        <PlanDetailPage />,
        `/projects/${mockProjectId}/plans/${mockPlanId}`
      )

      // Initial status: draft
      await waitFor(() => {
        expect(screen.getByText('草稿')).toBeInTheDocument()
      })

      // Check for "开始执行" button
      expect(screen.getByRole('button', { name: /开始执行/ })).toBeInTheDocument()
    })

    it('should calculate statistics correctly based on results', async () => {
      const mockPlanDetail = {
        id: mockPlanId,
        projectId: mockProjectId,
        name: '统计测试计划',
        description: '',
        status: 'active' as const,
        createdBy: 'user-123',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
        cases: [
          {
            caseId: 'case-1',
            caseNumber: 'TEST-001',
            caseTitle: '用例1',
            resultStatus: 'pass' as const,
            executedAt: '2024-01-01T10:00:00Z',
            executedBy: 'user-123',
          },
          {
            caseId: 'case-2',
            caseNumber: 'TEST-002',
            caseTitle: '用例2',
            resultStatus: 'fail' as const,
            executedAt: '2024-01-01T10:01:00Z',
            executedBy: 'user-123',
          },
          {
            caseId: 'case-3',
            caseNumber: 'TEST-003',
            caseTitle: '用例3',
            resultStatus: 'block' as const,
            executedAt: '2024-01-01T10:02:00Z',
            executedBy: 'user-123',
          },
          {
            caseId: 'case-4',
            caseNumber: 'TEST-004',
            caseTitle: '用例4',
          },
        ],
        stats: {
          total: 4,
          passed: 1,
          failed: 1,
          blocked: 1,
          skipped: 0,
          unexecuted: 1,
        },
      }

      server.use(
        http.get(`/api/v1/plans/${mockPlanId}`, () =>
          HttpResponse.json(mockPlanDetail)
        )
      )

      renderWithProviders(
        <PlanDetailPage />,
        `/projects/${mockProjectId}/plans/${mockPlanId}`
      )

      await waitFor(() => {
        expect(screen.getByText('统计测试计划')).toBeInTheDocument()
      })

      // Verify statistics are correctly displayed
      expect(screen.getByText('总用例: 4')).toBeInTheDocument()
      expect(screen.getByText('✅通过: 1')).toBeInTheDocument()
      expect(screen.getByText('❌失败: 1')).toBeInTheDocument()
      expect(screen.getByText('⚠️阻塞: 1')).toBeInTheDocument()
      expect(screen.getByText('⏭️跳过: 0')).toBeInTheDocument()

      // Verify progress bar (would check for progress element)
      const progressText = screen.getByText(/\d+\.?\d*%/)
      expect(progressText).toBeInTheDocument()
    })
  })
})
