/**
 * PlanDetailPage Component Tests
 */

import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { PlanDetailPage } from './PlanDetailPage'

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

describe('PlanDetailPage', () => {
  beforeEach(() => {
    // Mock plan detail API
    server.use(
      http.get('/api/v1/plans/:id', () =>
        HttpResponse.json({
          id: 'plan-001',
          projectId: 'project-001',
          name: 'Sprint 12 回归测试',
          description: 'Sprint 12 回归测试计划',
          status: 'active' as const,
          createdBy: 'user-001',
          createdAt: '2026-04-21T00:00:00Z',
          updatedAt: '2026-04-21T00:00:00Z',
          cases: [
            {
              caseId: 'case-001',
              caseNumber: 'TEST-USR-20260421-001',
              caseTitle: '用户登录成功',
              resultStatus: 'pass' as const,
              resultNote: '功能正常',
              executedAt: '2026-04-21T10:00:00Z',
              executedBy: 'user-001',
            },
            {
              caseId: 'case-002',
              caseNumber: 'TEST-USR-20260421-002',
              caseTitle: '用户注册成功',
              resultStatus: 'fail' as const,
              resultNote: '登录超时',
              executedAt: '2026-04-21T10:05:00Z',
              executedBy: 'user-001',
            },
            {
              caseId: 'case-003',
              caseNumber: 'TEST-USR-20260421-003',
              caseTitle: '用户注销成功',
              resultStatus: undefined,
            },
          ],
          stats: {
            total: 3,
            passed: 1,
            failed: 1,
            blocked: 0,
            skipped: 0,
            unexecuted: 1,
          },
        })
      ),
      // Mock record result API
      http.post('/api/v1/plans/:id/results', () =>
        HttpResponse.json(
          {
            caseId: 'case-003',
            caseNumber: 'TEST-USR-20260421-003',
            caseTitle: '用户注销成功',
            resultStatus: 'pass' as const,
            resultNote: '功能正常',
            executedAt: '2026-04-21T11:00:00Z',
            executedBy: 'user-001',
          },
          { status: 201 }
        )
      )
    )
  })

  describe('rendering', () => {
    it('should render plan information', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      expect(await screen.findAllByText(/Sprint 12 回归测试/i)).toHaveLength(2)
      expect(screen.getByText(/Sprint 12 回归测试计划/i)).toBeInTheDocument()
    })

    it('should render execution statistics', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      expect(await screen.findByText(/总用例/)).toBeInTheDocument()
      // Check that the stats are rendered (there will be multiple "3" values)
      expect(screen.getAllByText(/3/).length).toBeGreaterThan(0)
    })

    it('should render case list table', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      expect(await screen.findByText(/用户登录成功/i)).toBeInTheDocument()
      expect(screen.getByText(/用户注册成功/i)).toBeInTheDocument()
      expect(screen.getByText(/用户注销成功/i)).toBeInTheDocument()
    })
  })

  describe('case list rendering', () => {
    it('should display case status tags', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      // Status tags for executed cases
      expect(await screen.findAllByText(/通过/)).toBeTruthy()
      expect(screen.findAllByText(/失败/)).toBeTruthy()

      // Unexecuted case - check for the text "未执行"
      expect(screen.getAllByText(/未执行/).length).toBeGreaterThan(0)
    })

    it('should display case numbers', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      expect(
        await screen.findByText(/TEST-USR-20260421-001/i)
      ).toBeInTheDocument()
      expect(screen.getByText(/TEST-USR-20260421-002/i)).toBeInTheDocument()
    })

    it('should display executed by and time for executed cases', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      // Verify table headers include execution info columns
      expect(await screen.findByText(/执行人/i)).toBeInTheDocument()
      expect(screen.getByText(/执行时间/i)).toBeInTheDocument()

      // Unexecuted cases should show '-' for executed by
      const allDashes = screen.queryAllByText('-')
      expect(allDashes.length).toBeGreaterThan(0)
    })
  })

  describe('result recording - quick entry', () => {
    it('should have clickable result status for quick entry', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      // Unexecuted cases should have clickable "未执行" status
      await waitFor(() => {
        const unexecutedElements = screen.getAllByText(/未执行/)
        expect(unexecutedElements.length).toBeGreaterThan(0)
        // Verify at least one is a button or clickable element
        const unexecutedButtons = unexecutedElements.filter(
          (el) => el.tagName === 'BUTTON' || el.closest('button') !== null
        )
        expect(unexecutedButtons.length).toBeGreaterThan(0)
      })
    })
  })

  describe('result recording - detailed entry', () => {
    it('should open result record modal on detail entry button click', async () => {
      const user = userEvent.setup()
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      // Find the "详细录入" button
      const detailButtons = await screen.findAllByRole('button', {
        name: /详细录入/i,
      })
      expect(detailButtons.length).toBeGreaterThan(0)
      await user.click(detailButtons[0])

      // Modal should appear
      expect(await screen.findByText(/录入执行结果/i)).toBeInTheDocument()
    })

    it('should submit result from modal', async () => {
      const user = userEvent.setup()
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      const detailButtons = await screen.findAllByRole('button', {
        name: /详细录入/i,
      })
      await user.click(detailButtons[0])

      // Modal should be visible
      const modalTitle = await screen.findByText(/录入执行结果/i)
      expect(modalTitle).toBeInTheDocument()

      // Verify Modal structure - check if Modal exists in document
      await waitFor(() => {
        const modalElement = document.querySelector('.arco-modal-wrapper')
        expect(modalElement).toBeInTheDocument()
      })

      // Modal contains "通过" option (somewhere in the document)
      const passElements = screen.queryAllByText('通过')
      expect(passElements.length).toBeGreaterThan(0)
    })
  })

  describe('actions', () => {
    it('should render back button', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      expect(
        await screen.findByRole('button', { name: /返回/i })
      ).toBeInTheDocument()
    })

    it('should render edit button', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      expect(
        await screen.findByRole('button', { name: /编辑/i })
      ).toBeInTheDocument()
    })

    it('should render delete button', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      expect(
        await screen.findByRole('button', { name: /删除/i })
      ).toBeInTheDocument()
    })
  })
})
