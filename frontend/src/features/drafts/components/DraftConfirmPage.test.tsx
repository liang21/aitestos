import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { DraftConfirmPage } from './DraftConfirmPage'
import { server } from '../../../../tests/msw/server'
import { draftsHandlers } from '../../../../tests/msw/handlers/drafts'
import { http, HttpResponse } from 'msw'
import type { CaseDraft } from '@/types/api'

const mockDraft: CaseDraft = {
  id: 'draft-001',
  taskId: 'task-001',
  title: '验证用户登录功能',
  preconditions: ['用户已注册', '账号状态正常'],
  steps: ['打开登录页面', '输入账号密码', '点击登录'],
  expected: {
    step_1: '页面正常加载',
    step_2: '输入成功',
    step_3: '登录成功',
  },
  caseType: 'functionality',
  priority: 'P1',
  status: 'pending',
  createdAt: '2026-04-16T08:00:00Z',
  updatedAt: '2026-04-16T08:00:00Z',
  projectName: 'ECommerce',
  moduleName: '用户中心',
}

describe('DraftConfirmPage', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
    server.use(...draftsHandlers)
  })

  afterEach(() => {
    server.resetHandlers()
  })

  function renderWithProviders(ui: React.ReactElement) {
    return render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={['/drafts/draft-001']}>
          <Routes>
            <Route path="/drafts/:draftId" element={ui} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    )
  }

  it('should render draft confirm page', async () => {
    server.use(
      http.get('/api/v1/generation/drafts/draft-001', () =>
        HttpResponse.json(mockDraft)
      )
    )

    renderWithProviders(<DraftConfirmPage />)

    await waitFor(() => {
      expect(screen.getByText('草稿确认')).toBeInTheDocument()
    })
  })

  it('should render draft title and metadata', async () => {
    server.use(
      http.get('/api/v1/generation/drafts/draft-001', () =>
        HttpResponse.json(mockDraft)
      )
    )

    renderWithProviders(<DraftConfirmPage />)

    await waitFor(() => {
      expect(screen.getByDisplayValue('验证用户登录功能')).toBeInTheDocument()
    })

    expect(screen.getByText(/ECommerce/)).toBeInTheDocument()
    expect(screen.getByText(/用户中心/)).toBeInTheDocument()
  })

  it('should render form with draft data', async () => {
    server.use(
      http.get('/api/v1/generation/drafts/draft-001', () =>
        HttpResponse.json(mockDraft)
      )
    )

    renderWithProviders(<DraftConfirmPage />)

    await waitFor(() => {
      expect(screen.getByDisplayValue('验证用户登录功能')).toBeInTheDocument()
    })

    // Check that form is rendered (using more flexible selectors)
    expect(screen.getByText('用例标题')).toBeInTheDocument()
    expect(screen.getByText('测试步骤')).toBeInTheDocument()
  })

  it('should show confirm and reject buttons', async () => {
    server.use(
      http.get('/api/v1/generation/drafts/draft-001', () =>
        HttpResponse.json(mockDraft)
      )
    )

    renderWithProviders(<DraftConfirmPage />)

    await waitFor(() => {
      expect(screen.getByText('草稿确认')).toBeInTheDocument()
    })

    expect(
      screen.getByRole('button', { name: /确认并转为正式用例/i })
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /拒绝/i })).toBeInTheDocument()
  })

  it('should show reject modal when reject button clicked', async () => {
    const user = userEvent.setup()
    server.use(
      http.get('/api/v1/generation/drafts/draft-001', () =>
        HttpResponse.json(mockDraft)
      )
    )

    renderWithProviders(<DraftConfirmPage />)

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /拒绝/i })).toBeInTheDocument()
    })

    await user.click(screen.getByRole('button', { name: /拒绝/i }))

    // Wait for modal to appear
    await waitFor(() => {
      expect(screen.getByText('拒绝草稿')).toBeInTheDocument()
    })
  })
})
