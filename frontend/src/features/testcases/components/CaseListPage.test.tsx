/**
 * CaseListPage Component Tests
 */

import { describe, it, expect, beforeEach } from 'vitest'
import { render, screen, within } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { CaseListPage } from './CaseListPage'

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

describe('CaseListPage', () => {
  beforeEach(() => {
    // Setup before each test
  })

  const mockCases = [
    {
      id: 'case-001',
      moduleId: 'mod-001',
      userId: 'user-001',
      number: 'TEST-USR-20260421-001',
      title: '用户登录成功',
      preconditions: ['用户已注册'],
      steps: ['打开登录页', '输入用户名密码', '点击登录'],
      expected: { step_3: '登录成功，跳转首页' },
      caseType: 'functionality' as const,
      priority: 'P1' as const,
      status: 'unexecuted' as const,
      createdAt: '2026-04-21T00:00:00Z',
      updatedAt: '2026-04-21T00:00:00Z',
    },
    {
      id: 'case-002',
      moduleId: 'mod-001',
      userId: 'user-001',
      number: 'TEST-USR-20260421-002',
      title: '用户登录失败',
      preconditions: ['用户已注册'],
      steps: ['打开登录页', '输入错误密码', '点击登录'],
      expected: { step_3: '显示错误提示' },
      caseType: 'functionality' as const,
      priority: 'P2' as const,
      status: 'pass' as const,
      createdAt: '2026-04-21T00:00:00Z',
      updatedAt: '2026-04-21T00:00:00Z',
    },
  ]

  describe('rendering', () => {
    it('should render case table with correct columns', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: mockCases,
            total: 2,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      // Check column headers
      expect(screen.getByText('编号')).toBeInTheDocument()
      expect(screen.getByText('标题')).toBeInTheDocument()
      expect(screen.getByText('类型')).toBeInTheDocument()
      expect(screen.getByText('优先级')).toBeInTheDocument()
      expect(screen.getByText('状态')).toBeInTheDocument()
    })

    it('should render case data rows', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: mockCases,
            total: 2,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      expect(
        await screen.findByText('TEST-USR-20260421-001')
      ).toBeInTheDocument()
      expect(screen.getByText('用户登录成功')).toBeInTheDocument()
      expect(screen.getByText('TEST-USR-20260421-002')).toBeInTheDocument()
      expect(screen.getByText('用户登录失败')).toBeInTheDocument()
    })

    it('should display type tags', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: mockCases,
            total: 2,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      // Functionality type should be displayed
      const typeTags = await screen.findAllByText('功能测试')
      expect(typeTags.length).toBeGreaterThan(0)
    })

    it('should display priority tags', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: mockCases,
            total: 2,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      expect(await screen.findByText('P1 高')).toBeInTheDocument()
      expect(screen.getByText('P2 中')).toBeInTheDocument()
    })

    it('should display status tags', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: mockCases,
            total: 2,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      expect(await screen.findByText('未执行')).toBeInTheDocument()
      expect(screen.getByText('通过')).toBeInTheDocument()
    })

    it('should show "新建用例" button', () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      expect(
        screen.getByRole('button', { name: /新建用例/i })
      ).toBeInTheDocument()
    })
  })

  describe('filtering', () => {
    it('should render filter bar with type select', () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      // Type filter select should exist
      const typeSelect = screen.getByPlaceholderText(/类型/i)
      expect(typeSelect).toBeInTheDocument()
    })

    it('should render filter bar with priority select', () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      // Priority filter select should exist
      const prioritySelect = screen.getByPlaceholderText(/优先级/i)
      expect(prioritySelect).toBeInTheDocument()
    })

    it('should render filter bar with status select', () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      // Status filter select should exist
      const statusSelect = screen.getByPlaceholderText(/状态/i)
      expect(statusSelect).toBeInTheDocument()
    })
  })

  describe('pagination', () => {
    it('should render pagination with total count', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: mockCases,
            total: 25,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      expect(await screen.findByText(/共 25 条/i)).toBeInTheDocument()
    })
  })

  describe('row click', () => {
    it('should navigate to detail page on row click', async () => {
      server.use(
        http.get('/api/v1/testcases', () =>
          HttpResponse.json({
            data: [mockCases[0]],
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(<CaseListPage projectId="project-001" />)

      const caseLink = await screen.findByText('TEST-USR-20260421-001')
      expect(caseLink.closest('a')).toHaveAttribute(
        'href',
        '/testcases/case-001'
      )
    })
  })
})
