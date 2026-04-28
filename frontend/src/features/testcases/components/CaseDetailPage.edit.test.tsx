import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { CaseDetailPage } from './CaseDetailPage'

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
}

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = createTestQueryClient()
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>{ui}</BrowserRouter>
    </QueryClientProvider>
  )
}

describe('CaseDetailPage - Edit and Copy', () => {
  const mockCaseDetail = {
    id: '1',
    number: 'ECO-USR-20260428-001',
    title: '密码错误超过5次锁定账号',
    caseType: 'functionality',
    priority: 'P1',
    status: 'pass',
    moduleId: 'mod-1',
    project_id: 'proj-1',
    module: { id: 'mod-1', name: '登录模块', abbreviation: 'USR' },
    preconditions: ['用户已登录系统'],
    steps: ['输入错误密码', '重复至第5次', '验证锁定状态'],
    expected: { step_3: '第5次后锁定' },
    userId: 'u-1',
    createdByName: '张三',
    createdAt: '2026-04-28T10:00:00Z',
    updatedAt: '2026-04-28T11:00:00Z',
  }

  const mockModules = [
    { id: 'mod-1', name: '登录模块', abbreviation: 'USR' },
    { id: 'mod-2', name: '支付模块', abbreviation: 'PAY' },
  ]

  beforeEach(() => {
    server.resetHandlers()
    server.use(
      http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCaseDetail)),
      http.get('/api/v1/projects/:id/modules', () => HttpResponse.json(mockModules))
    )
  })

  describe('edit functionality', () => {
    it('should open CreateCaseDrawer with pre-filled data when clicking edit button', async () => {
      const user = userEvent.setup()
      renderWithProviders(<CaseDetailPage caseId="1" />)

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument()
      })

      // Click edit button
      await user.click(screen.getByText('编辑'))

      // Drawer should open with pre-filled data
      await waitFor(() => {
        expect(screen.getByText('编辑测试用例')).toBeInTheDocument()
      })

      // Check pre-filled values
      expect(screen.getByDisplayValue('密码错误超过5次锁定账号')).toBeInTheDocument()
    })

    it('should call PUT /testcases/{id} when submitting edit form', async () => {
      const user = userEvent.setup()
      let updateCalled = false

      server.use(
        http.put('/api/v1/testcases/:id', () => {
          updateCalled = true
          return HttpResponse.json({
            ...mockCaseDetail,
            title: '更新后的标题',
            updated_at: '2026-04-28T12:00:00Z',
          })
        })
      )

      renderWithProviders(<CaseDetailPage caseId="1" />)

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument()
      })

      await user.click(screen.getByText('编辑'))

      await waitFor(() => {
        expect(screen.getByText('编辑测试用例')).toBeInTheDocument()
      })

      // Wait for drawer to fully render and form to be ready
      await waitFor(() => {
        expect(screen.getByDisplayValue('密码错误超过5次锁定账号')).toBeInTheDocument()
      }, { timeout: 3000 })

      // Modify title
      const titleInput = screen.getByDisplayValue('密码错误超过5次锁定账号')
      await user.clear(titleInput)
      await user.type(titleInput, '更新后的标题')

      // Submit
      const submitButton = screen.getByText('保存')
      await user.click(submitButton)

      await waitFor(() => {
        expect(updateCalled).toBe(true)
      })
    })
  })

  describe('copy functionality', () => {
    it('should open CreateCaseDrawer with pre-filled data and [副本] prefix when clicking copy button', async () => {
      const user = userEvent.setup()
      renderWithProviders(<CaseDetailPage caseId="1" />)

      await waitFor(() => {
        expect(screen.getByText('复制')).toBeInTheDocument()
      })

      // Click copy button
      await user.click(screen.getByText('复制'))

      // Drawer should open with pre-filled data
      await waitFor(() => {
        expect(screen.getByText('新建测试用例')).toBeInTheDocument()
      })

      // Check title has [副本] prefix
      expect(screen.getByDisplayValue('[副本] 密码错误超过5次锁定账号')).toBeInTheDocument()
    })

    it('should call POST /testcases when submitting copy form', async () => {
      const user = userEvent.setup()
      let createCalled = false

      server.use(
        http.post('/api/v1/testcases', () => {
          createCalled = true
          return HttpResponse.json({
            id: '2',
            number: 'ECO-USR-20260428-002',
            ...mockCaseDetail,
            title: '[副本] 密码错误超过5次锁定账号',
            created_at: '2026-04-28T12:00:00Z',
            updated_at: '2026-04-28T12:00:00Z',
          })
        })
      )

      renderWithProviders(<CaseDetailPage caseId="1" />)

      await waitFor(() => {
        expect(screen.getByText('复制')).toBeInTheDocument()
      })

      await user.click(screen.getByText('复制'))

      await waitFor(() => {
        expect(screen.getByText('新建测试用例')).toBeInTheDocument()
      })

      // Submit copy
      const submitButton = screen.getByText('创建用例')
      await user.click(submitButton)

      await waitFor(() => {
        expect(createCalled).toBe(true)
      })
    })
  })

  describe('delete functionality', () => {
    it('should show Popconfirm before deleting', async () => {
      const user = userEvent.setup()
      renderWithProviders(<CaseDetailPage caseId="1" />)

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument()
      })

      // Click delete button should show confirmation
      await user.click(screen.getByText('删除'))

      // Popconfirm should appear
      await waitFor(() => {
        expect(screen.getByText(/确认删除/)).toBeInTheDocument()
      })
    })

    it('should call DELETE /testcases/{id} after confirming delete', async () => {
      const user = userEvent.setup()
      let deleteCalled = false

      server.use(
        http.delete('/api/v1/testcases/:id', () => {
          deleteCalled = true
          return new HttpResponse(null, { status: 204 })
        })
      )

      renderWithProviders(<CaseDetailPage caseId="1" />)

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument()
      })

      await user.click(screen.getByText('删除'))

      await waitFor(() => {
        expect(screen.getByText(/确认删除/)).toBeInTheDocument()
      })

      // Click confirm
      const confirmButton = screen.getByText('确认')
      await user.click(confirmButton)

      await waitFor(() => {
        expect(deleteCalled).toBe(true)
      })
    })
  })
})
