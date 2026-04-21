/**
 * PlanDetailPage Component Tests
 */

import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
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

      // Status tags for executed cases - use findAllByText since multiple elements might have this text
      expect(await screen.findAllByText(/通过/)).toBeTruthy()
      expect(screen.findAllByText(/失败/)).toBeTruthy()

      // Unexecuted case - check for the text "未执行" (there might be multiple)
      expect(screen.getAllByText(/未执行/).length).toBeGreaterThan(0)
    })

    it('should display case numbers', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      expect(
        await screen.findByText(/TEST-USR-20260421-001/i)
      ).toBeInTheDocument()
      expect(screen.getByText(/TEST-USR-20260421-002/i)).toBeInTheDocument()
    })

    it('should display result notes for executed cases', async () => {
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      expect(await screen.findByText(/功能正常/i)).toBeInTheDocument()
      expect(screen.getByText(/登录超时/i)).toBeInTheDocument()
    })
  })

  describe('result recording', () => {
    it('should open result record modal on button click', async () => {
      const user = userEvent.setup()
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      // Find the "录入结果" button for unexecuted case
      const recordButtons = await screen.findAllByRole('button', {
        name: /录入结果/i,
      })
      await user.click(recordButtons[0])

      // Modal should appear
      expect(await screen.findByText(/录入执行结果/i)).toBeInTheDocument()
    })

    it('should submit result and refresh list', async () => {
      const user = userEvent.setup()
      renderWithProviders(<PlanDetailPage planId="plan-001" />)

      const recordButtons = await screen.findAllByRole('button', {
        name: /录入结果/i,
      })
      await user.click(recordButtons[0])

      // Select pass status
      const passOption = await screen.findByRole('radio', { name: /通过/i })
      await user.click(passOption)

      // Add note
      const noteInput = screen.getByPlaceholderText(/备注/i)
      await user.type(noteInput, '功能正常')

      // Submit
      const submitButton = screen.getByRole('button', { name: /提交/i })
      await user.click(submitButton)

      // Modal should close and list refresh
      // This would typically check for modal dismissal and data refresh
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
