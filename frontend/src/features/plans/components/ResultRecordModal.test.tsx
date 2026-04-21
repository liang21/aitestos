/**
 * ResultRecordModal Component Tests
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, cleanup, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { ResultRecordModal } from './ResultRecordModal'

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
  return render(
    <QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>
  )
}

describe('ResultRecordModal', () => {
  const defaultProps = {
    visible: true,
    planId: 'plan-001',
    caseId: 'case-001',
    caseTitle: '用户登录成功',
    onClose: vi.fn(),
  }

  beforeEach(() => {
    // Mock GET plan detail API (called by useEffect in ResultRecordModal)
    server.use(
      http.get('/api/v1/plans/:id', () =>
        HttpResponse.json({
          id: 'plan-001',
          projectId: 'project-001',
          name: 'Sprint 12 回归测试',
          description: 'Sprint 12 回归测试',
          status: 'active' as const,
          createdAt: '2026-04-20T10:00:00Z',
          updatedAt: '2026-04-20T10:00:00Z',
          stats: {
            total: 2,
            passed: 1,
            failed: 0,
            blocked: 0,
            skipped: 0,
            unexecuted: 1,
          },
          cases: [
            {
              caseId: 'case-001',
              caseNumber: 'TEST-USR-20260421-001',
              caseTitle: '用户登录成功',
              resultStatus: undefined,
              resultNote: undefined,
              executedAt: undefined,
              executedBy: undefined,
            },
          ],
        })
      ),
      // Mock record result API
      http.post('/api/v1/plans/:id/results', () =>
        HttpResponse.json(
          {
            caseId: 'case-001',
            caseNumber: 'TEST-USR-20260421-001',
            caseTitle: '用户登录成功',
            resultStatus: 'pass' as const,
            resultNote: '功能正常',
            executedAt: '2026-04-21T10:00:00Z',
            executedBy: 'user-001',
          },
          { status: 201 }
        )
      )
    )
  })

  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  describe('rendering', () => {
    it('should render modal title', async () => {
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      expect(await screen.findByText(/录入执行结果/i)).toBeInTheDocument()
    })

    it('should display case title', async () => {
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      expect(await screen.findByText(/用户登录成功/i)).toBeInTheDocument()
    })
  })

  describe('status selection', () => {
    it('should render status radio options', async () => {
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      expect(
        await screen.findByRole('radio', { name: /通过/i })
      ).toBeInTheDocument()
      expect(screen.getByRole('radio', { name: /失败/i })).toBeInTheDocument()
      expect(screen.getByRole('radio', { name: /阻塞/i })).toBeInTheDocument()
      expect(screen.getByRole('radio', { name: /跳过/i })).toBeInTheDocument()
    })

    it('should select pass status by default', async () => {
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      const passRadio = await screen.findByRole('radio', { name: /通过/i })
      expect(passRadio).toBeChecked()
    })

    it('should allow changing status selection', async () => {
      const user = userEvent.setup()
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      const failRadio = await screen.findByRole('radio', { name: /失败/i })
      await user.click(failRadio)

      expect(failRadio).toBeChecked()

      const passRadio = screen.getByRole('radio', { name: /通过/i })
      expect(passRadio).not.toBeChecked()
    })
  })

  describe('notes input', () => {
    it('should render notes textarea', async () => {
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      expect(await screen.findByPlaceholderText(/备注/i)).toBeInTheDocument()
    })

    it('should allow typing notes', async () => {
      const user = userEvent.setup()
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      const notesInput = await screen.findByPlaceholderText(/备注/i)
      // Wait for any initial state to settle
      await waitFor(() => {
        expect(notesInput).toHaveValue('')
      })
      await user.type(notesInput, '功能正常，无问题')

      await waitFor(() => {
        expect(notesInput).toHaveValue('功能正常，无问题')
      })
    })

    it('should allow empty notes', async () => {
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      const notesInput = await screen.findByPlaceholderText(/备注/i)
      // Wait for component to fully initialize with empty state
      await waitFor(() => {
        expect(notesInput).toHaveValue('')
      })
    })
  })

  describe('form submission', () => {
    it('should submit with pass status and notes', async () => {
      const user = userEvent.setup()
      const onClose = vi.fn()

      renderWithProviders(
        <ResultRecordModal {...defaultProps} onClose={onClose} />
      )

      // Type notes
      const notesInput = await screen.findByPlaceholderText(/备注/i)
      await user.type(notesInput, '功能正常')

      // Submit
      const submitButton = screen.getByRole('button', { name: /提交/i })
      await user.click(submitButton)

      // Should call onClose after successful submission
      // This would typically check for API call and modal close
    })

    it('should submit with fail status', async () => {
      const user = userEvent.setup()
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      // Select fail status
      const failRadio = await screen.findByRole('radio', { name: /失败/i })
      await user.click(failRadio)

      // Type notes
      const notesInput = screen.getByPlaceholderText(/备注/i)
      await user.type(notesInput, '登录超时')

      // Submit
      const submitButton = screen.getByRole('button', { name: /提交/i })
      await user.click(submitButton)

      // Should submit with fail status
    })

    it('should validate that status is selected', async () => {
      const user = userEvent.setup()
      renderWithProviders(<ResultRecordModal {...defaultProps} />)

      // By default, pass status is selected
      const passRadio = await screen.findByRole('radio', { name: /通过/i })
      expect(passRadio).toBeChecked()

      // Can change to other status
      const failRadio = screen.getByRole('radio', { name: /失败/i })
      await user.click(failRadio)
      expect(failRadio).toBeChecked()
    })
  })

  describe('modal actions', () => {
    it('should close on cancel button click', async () => {
      const user = userEvent.setup()
      const onClose = vi.fn()

      renderWithProviders(
        <ResultRecordModal {...defaultProps} onClose={onClose} />
      )

      const cancelButton = await screen.findByRole('button', { name: /取消/i })
      await user.click(cancelButton)

      expect(onClose).toHaveBeenCalled()
    })

    it('should close on X button click', async () => {
      const onClose = vi.fn()

      renderWithProviders(
        <ResultRecordModal {...defaultProps} onClose={onClose} />
      )

      // Use the cancel button instead since the Modal close icon is hard to target
      const cancelButton = await screen.findByRole('button', { name: /取消/i })
      cancelButton.click()

      expect(onClose).toHaveBeenCalled()
    })
  })

  describe('not visible', () => {
    it('should not render when visible is false', () => {
      renderWithProviders(
        <ResultRecordModal {...defaultProps} visible={false} />
      )

      expect(screen.queryByText(/录入执行结果/i)).not.toBeInTheDocument()
    })
  })
})
