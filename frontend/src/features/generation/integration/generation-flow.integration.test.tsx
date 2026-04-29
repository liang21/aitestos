/**
 * AI Generation Flow Integration Tests
 * Tests the complete flow: Document Upload → Create Task → View Results → Confirm Draft
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse, delay } from 'msw'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'

// Import components to test
import { KnowledgeListPage } from '@/features/documents/components/KnowledgeListPage'
import { NewGenerationTaskPage } from '@/features/generation/components/NewGenerationTaskPage'
import { TaskDetailPage } from '@/features/generation/components/TaskDetailPage'
import { DraftConfirmPage } from '@/features/drafts/components/DraftConfirmPage'

describe('AI Generation Flow Integration Tests', () => {
  let queryClient: QueryClient
  const mockProjectId = 'project-123'
  const mockModuleId = 'module-123'

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

  function RouterWrapper({ children }: { children: React.ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route path="/projects/:projectId/knowledge" element={children} />
            <Route path="/projects/:projectId/generation/new" element={children} />
            <Route path="/projects/:projectId/generation/:taskId" element={children} />
            <Route path="/drafts/:draftId" element={children} />
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    )
  }

  describe('Complete AI Generation Flow', () => {
    it('should complete full flow: upload document → create task → view results → confirm draft', async () => {
      const user = userEvent.setup()

      // Step 1: Mock document list (simulate uploaded documents)
      const mockDocuments = {
        data: [
          {
            id: 'doc-1',
            projectId: mockProjectId,
            name: 'Login Module PRD',
            type: 'prd' as const,
            status: 'completed' as const,
            chunkCount: 12,
            uploadedBy: 'user-123',
            createdAt: '2024-01-01T00:00:00Z',
            updatedAt: '2024-01-01T00:00:00Z',
          },
        ],
        total: 1,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get(`/api/v1/knowledge/documents`, () =>
          HttpResponse.json(mockDocuments)
        )
      )

      // Mock module list
      const mockModules = {
        data: [
          {
            id: mockModuleId,
            projectId: mockProjectId,
            name: '登录模块',
            abbreviation: 'LOGIN',
            createdAt: '2024-01-01T00:00:00Z',
            updatedAt: '2024-01-01T00:00:00Z',
          },
        ],
        total: 1,
        offset: 0,
        limit: 10,
      }

      server.use(
        http.get(`/api/v1/projects/${mockProjectId}/modules`, () =>
          HttpResponse.json(mockModules)
        )
      )

      // Mock task creation
      const mockTask = {
        id: 'task-123',
        projectId: mockProjectId,
        moduleId: mockModuleId,
        status: 'processing' as const,
        prompt: '测试用户登录功能',
        result: null,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      let taskStatus: 'processing' | 'completed' = 'processing'

      server.use(
        http.post(`/api/v1/generation/tasks`, async () => {
          await delay(100)
          return HttpResponse.json(mockTask, { status: 201 })
        }),
        http.get(`/api/v1/generation/tasks/:taskId`, async ({ params }) => {
          await delay(50)
          if (params.taskId === 'task-123') {
            // Simulate status change on second poll
            if (taskStatus === 'processing') {
              taskStatus = 'completed'
              return HttpResponse.json({
                ...mockTask,
                status: 'completed',
                result: { draftCount: 3, confidence: 'high' },
              })
            }
          }
          return HttpResponse.json(mockTask)
        })
      )

      // Mock draft list
      const mockDrafts = [
        {
          id: 'draft-1',
          taskId: 'task-123',
          projectId: mockProjectId,
          title: '用户名密码登录',
          preconditions: ['用户已注册'],
          steps: ['打开登录页', '输入用户名密码', '点击登录'],
          expected: { result: '登录成功' },
          caseType: 'functionality' as const,
          priority: 'P1' as const,
          status: 'pending' as const,
          aiMetadata: {
            confidence: 'high' as const,
            referencedChunks: [],
            modelVersion: 'deepseek-chat-v3',
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
          projectName: '测试项目',
          moduleName: '登录模块',
        },
      ]

      server.use(
        http.get(`/api/v1/generation/tasks/task-123/drafts`, () =>
          HttpResponse.json(mockDrafts)
        )
      )

      // Mock draft detail
      server.use(
        http.get(`/api/v1/generation/drafts/draft-1`, () =>
          HttpResponse.json(mockDrafts[0])
        )
      )

      // Mock draft confirmation
      const mockTestCase = {
        id: 'case-123',
        moduleId: mockModuleId,
        userId: 'user-123',
        number: 'TEST-LOGIN-20240429-001',
        title: '用户名密码登录',
        preconditions: ['用户已注册'],
        steps: ['打开登录页', '输入用户名密码', '点击登录'],
        expected: { result: '登录成功' },
        caseType: 'functionality' as const,
        priority: 'P1' as const,
        status: 'unexecuted' as const,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      server.use(
        http.post(`/api/v1/generation/drafts/draft-1/confirm`, () =>
          HttpResponse.json(mockTestCase, { status: 201 })
        )
      )

      // Verify Step 1: Knowledge base has documents (knowledge readiness)
      const { unmount: unmountKnowledge } = renderWithProviders(
        <KnowledgeListPage />,
        `/projects/${mockProjectId}/knowledge`
      )

      await waitFor(() => {
        expect(screen.getByText('Login Module PRD')).toBeInTheDocument()
      })

      unmountKnowledge()

      // Verify Step 2: Create generation task with sufficient knowledge
      const { unmount: unmountNewTask } = renderWithProviders(
        <NewGenerationTaskPage projectId={mockProjectId} />,
        `/projects/${mockProjectId}/generation/new`
      )

      // Check knowledge readiness indicator shows sufficient
      await waitFor(() => {
        expect(screen.getByText(/🟢 就绪/)).toBeInTheDocument()
      })

      // Fill form and submit
      const promptInput = screen.getByPlaceholderText(/请描述测试需求/)
      await user.click(promptInput)
      await user.keyboard('测试用户登录功能，包括正常登录和错误处理')

      const submitButton = screen.getByRole('button', { name: /立即生成/ })
      await user.click(submitButton)

      // Verify success message
      await waitFor(() => {
        expect(screen.getByText('生成任务创建成功')).toBeInTheDocument()
      })

      unmountNewTask()

      // Verify Step 3: View task results
      const { unmount: unmountTaskDetail } = renderWithProviders(
        <TaskDetailPage />,
        `/projects/${mockProjectId}/generation/task-123`
      )

      // Wait for task status to change to completed
      await waitFor(
        () => {
          expect(screen.getByText('已完成')).toBeInTheDocument()
        },
        { timeout: 5000 }
      )

      // Verify drafts are displayed
      await waitFor(() => {
        expect(screen.getByText('生成的草稿')).toBeInTheDocument()
        expect(screen.getByText('用户名密码登录')).toBeInTheDocument()
      })

      unmountTaskDetail()

      // Verify Step 4: Confirm draft
      const { unmount: unmountDraftConfirm } = renderWithProviders(
        <DraftConfirmPage />,
        `/drafts/draft-1`
      )

      // Check draft content is displayed
      await waitFor(() => {
        expect(screen.getByDisplayValue('用户名密码登录')).toBeInTheDocument()
      })

      // Confirm draft
      const confirmButton = screen.getByRole('button', { name: /确认并转为正式用例/ })
      await user.click(confirmButton)

      // Select module in modal
      await waitFor(() => {
        expect(screen.getByText(/请选择要将此草稿确认到哪个模块/)).toBeInTheDocument()
      })

      const moduleSelect = screen.getByRole('combobox')
      await user.click(moduleSelect)
      await user.click(screen.getByText(/登录模块/))

      const modalConfirmButton = screen.getByRole('button', { name: '确认' })
      await user.click(modalConfirmButton)

      // Verify success message with case number
      await waitFor(() => {
        expect(
          screen.getByText(/用例 TEST-LOGIN-20240429-001 已创建/)
        ).toBeInTheDocument()
      })

      unmountDraftConfirm()
    })

    it('should show insufficient knowledge warning when documents are low', async () => {
      // Mock empty document list
      server.use(
        http.get(`/api/v1/knowledge/documents`, () =>
          HttpResponse.json({
            data: [],
            total: 0,
            offset: 0,
            limit: 10,
          })
        )
      )

      // Mock module list
      server.use(
        http.get(`/api/v1/projects/${mockProjectId}/modules`, () =>
          HttpResponse.json({
            data: [
              {
                id: mockModuleId,
                projectId: mockProjectId,
                name: '测试模块',
                abbreviation: 'TEST',
                createdAt: '2024-01-01T00:00:00Z',
                updatedAt: '2024-01-01T00:00:00Z',
              },
            ],
            total: 1,
            offset: 0,
            limit: 10,
          })
        )
      )

      renderWithProviders(
        <NewGenerationTaskPage projectId={mockProjectId} />,
        `/projects/${mockProjectId}/generation/new`
      )

      // Should show empty knowledge base warning
      await waitFor(() => {
        expect(screen.getByText(/🔴 请先上传需求文档/)).toBeInTheDocument()
      })

      // Submit button should be disabled
      const submitButton = screen.getByRole('button', { name: /立即生成/ })
      expect(submitButton).toBeDisabled()
    })
  })

  describe('Draft Confirmation Flow', () => {
    it('should support navigation between drafts', async () => {
      const mockDrafts = [
        {
          id: 'draft-1',
          taskId: 'task-123',
          projectId: mockProjectId,
          title: '草稿 1',
          preconditions: [],
          steps: ['步骤 1'],
          expected: { result: '成功' },
          caseType: 'functionality' as const,
          priority: 'P1' as const,
          status: 'pending' as const,
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        {
          id: 'draft-2',
          taskId: 'task-123',
          projectId: mockProjectId,
          title: '草稿 2',
          preconditions: [],
          steps: ['步骤 1'],
          expected: { result: '成功' },
          caseType: 'functionality' as const,
          priority: 'P2' as const,
          status: 'pending' as const,
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ]

      server.use(
        http.get('/api/v1/generation/drafts', () =>
          HttpResponse.json({
            data: mockDrafts,
            total: 2,
            offset: 0,
            limit: 10,
          })
        ),
        http.get('/api/v1/generation/drafts/:draftId', ({ params }) => {
          const draft = mockDrafts.find(d => d.id === params.draftId)
          return HttpResponse.json(draft)
        })
      )

      renderWithProviders(
        <DraftConfirmPage />,
        '/drafts/draft-1'
      )

      // Should show navigation indicator
      await waitFor(() => {
        expect(screen.getByText(/第 1 \/ 2 条/)).toBeInTheDocument()
      })

      // Should have navigation dots
      const dots = screen.getAllByRole('button').filter(
        btn => btn.className.includes('rounded-full')
      )
      expect(dots).toHaveLength(2)
    })

    it('should reject draft with reason', async () => {
      const mockDraft = {
        id: 'draft-1',
        taskId: 'task-123',
        projectId: mockProjectId,
        title: '低质量草稿',
        preconditions: [],
        steps: ['步骤'],
        expected: { result: '成功' },
        caseType: 'functionality' as const,
        priority: 'P1' as const,
        status: 'pending' as const,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      server.use(
        http.get('/api/v1/generation/drafts', () =>
          HttpResponse.json({
            data: [mockDraft],
            total: 1,
            offset: 0,
            limit: 10,
          })
        ),
        http.get('/api/v1/generation/drafts/draft-1', () =>
          HttpResponse.json(mockDraft)
        ),
        http.post('/api/v1/generation/drafts/draft-1/reject', () =>
          HttpResponse.json({ success: true })
        )
      )

      const user = userEvent.setup()

      renderWithProviders(
        <DraftConfirmPage />,
        '/drafts/draft-1'
      )

      await waitFor(() => {
        expect(screen.getByDisplayValue('低质量草稿')).toBeInTheDocument()
      })

      // Click reject button
      const rejectButton = screen.getByRole('button', { name: /拒绝/ })
      await user.click(rejectButton)

      // Should show reject modal
      await waitFor(() => {
        expect(screen.getByText('拒绝草稿')).toBeInTheDocument()
      })

      // Select reason
      const reasonSelect = screen.getByRole('combobox')
      await user.click(reasonSelect)
      await user.click(screen.getByText('低质量'))

      // Submit rejection
      const modalConfirmButton = screen.getAllByRole('button').find(
        btn => btn.textContent === '确认拒绝'
      )
      await user.click(modalConfirmButton!)

      // Verify success message
      await waitFor(() => {
        expect(screen.getByText('草稿已拒绝')).toBeInTheDocument()
      })
    })
  })
})
