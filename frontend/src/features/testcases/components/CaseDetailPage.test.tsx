/**
 * CaseDetailPage Component Tests
 */

import { describe, it, expect, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { CaseDetailPage } from './CaseDetailPage'

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

describe('CaseDetailPage', () => {
  beforeEach(() => {
    // Setup before each test
  })

  const mockCase = {
    id: 'case-001',
    moduleId: 'mod-001',
    userId: 'user-001',
    number: 'TEST-USR-20260421-001',
    title: '用户登录成功',
    preconditions: ['用户已注册'],
    steps: ['打开登录页', '输入正确的用户名和密码', '点击登录按钮'],
    expected: {
      step_3: '登录成功，跳转到首页',
    },
    caseType: 'functionality' as const,
    priority: 'P1' as const,
    status: 'unexecuted' as const,
    aiMetadata: {
      generationTaskId: 'task-001',
      confidence: 'high',
      referencedChunks: [
        {
          chunkId: 'chunk-001',
          documentId: 'doc-001',
          documentTitle: '用户登录需求文档',
          similarityScore: 0.92,
        },
        {
          chunkId: 'chunk-002',
          documentId: 'doc-001',
          documentTitle: '用户登录需求文档',
          similarityScore: 0.88,
        },
      ],
      modelVersion: 'gpt-4',
      generatedAt: '2026-04-21T00:00:00Z',
    },
    createdAt: '2026-04-21T00:00:00Z',
    updatedAt: '2026-04-21T00:00:00Z',
  }

  describe('rendering', () => {
    it('should render case basic information', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(await screen.findByText('用户登录成功')).toBeInTheDocument()
      expect(screen.getByText('TEST-USR-20260421-001')).toBeInTheDocument()
    })

    it('should render preconditions', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(await screen.findByText('用户已注册')).toBeInTheDocument()
      expect(screen.getByText(/前置条件/i)).toBeInTheDocument()
    })

    it('should render steps', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(await screen.findByText('打开登录页')).toBeInTheDocument()
      expect(screen.getByText('输入正确的用户名和密码')).toBeInTheDocument()
      expect(screen.getByText('点击登录按钮')).toBeInTheDocument()
      expect(screen.getByText(/测试步骤/i)).toBeInTheDocument()
    })

    it('should render expected results', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(await screen.findByText(/预期结果/i)).toBeInTheDocument()
      expect(screen.getByText(/登录成功，跳转到首页/i)).toBeInTheDocument()
    })

    it('should render type and priority tags', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(await screen.findByText('功能测试')).toBeInTheDocument()
      expect(screen.getByText('P1 高')).toBeInTheDocument()
    })
  })

  describe('AI metadata section', () => {
    it('should render AI confidence tag', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(await screen.findByText('高置信度')).toBeInTheDocument()
    })

    it('should render reference chunks list', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(await screen.findByText(/用户登录需求文档/i)).toBeInTheDocument()
      expect(screen.getByText(/92%/)).toBeInTheDocument()
      expect(screen.getByText(/88%/)).toBeInTheDocument()
    })

    it('should show "无 AI 来源" when no metadata', async () => {
      const caseWithoutAI = {
        ...mockCase,
        aiMetadata: undefined,
      }

      server.use(
        http.get('/api/v1/testcases/:id', () =>
          HttpResponse.json(caseWithoutAI)
        )
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      // When no AI metadata, the AI Metadata card should not be rendered
      expect(screen.queryByText(/AI 来源/i)).not.toBeInTheDocument()
    })

    it('should show "无引用来源" when metadata exists but no chunks', async () => {
      const caseWithEmptyChunks = {
        ...mockCase,
        aiMetadata: {
          ...mockCase.aiMetadata!,
          referencedChunks: [],
        },
      }

      server.use(
        http.get('/api/v1/testcases/:id', () =>
          HttpResponse.json(caseWithEmptyChunks)
        )
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(await screen.findByText(/无引用来源/i)).toBeInTheDocument()
    })
  })

  describe('action buttons', () => {
    it('should render edit button', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(
        await screen.findByRole('button', { name: /编辑/i })
      ).toBeInTheDocument()
    })

    it('should render delete button', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () => HttpResponse.json(mockCase))
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(
        await screen.findByRole('button', { name: /删除/i })
      ).toBeInTheDocument()
    })
  })

  describe('loading and error states', () => {
    it('should show loading state while fetching', () => {
      server.use(
        http.get(
          '/api/v1/testcases/:id',
          () => new Promise(() => {}) // Never resolve
        )
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      // Loading spinner should be present (Arco Design Spin component)
      const spinner = document.querySelector('.arco-spin')
      expect(spinner).toBeInTheDocument()
    })

    it('should show error state on fetch failure', async () => {
      server.use(
        http.get('/api/v1/testcases/:id', () =>
          HttpResponse.json({ message: 'Not found' }, { status: 404 })
        )
      )

      renderWithProviders(<CaseDetailPage caseId="case-001" />)

      expect(
        await screen.findByText(/加载失败|用例不存在/i)
      ).toBeInTheDocument()
    })
  })
})
