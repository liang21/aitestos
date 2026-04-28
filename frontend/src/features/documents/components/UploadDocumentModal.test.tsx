import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { server } from '../../../../tests/msw/server'
import { http, HttpResponse } from 'msw'
import { UploadDocumentModal } from './UploadDocumentModal'

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
}

function renderWithProviders(ui: any) {
  const queryClient = createTestQueryClient()
  return render(
    <QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>
  )
}

describe('UploadDocumentModal', () => {
  const mockOnCancel = vi.fn()
  const mockOnSuccess = vi.fn()

  const defaultProps = {
    visible: true,
    projectId: 'proj-1',
    onCancel: mockOnCancel,
    onSuccess: mockOnSuccess,
  }

  beforeEach(() => {
    server.resetHandlers()
    vi.clearAllMocks()
    // Mock localStorage to provide access_token
    vi.stubGlobal('localStorage', {
      getItem: vi.fn((key) => {
        if (key === 'access_token') return 'mock-token'
        return null
      }),
      setItem: vi.fn(),
      removeItem: vi.fn(),
      clear: vi.fn(),
    })
  })

  describe('form rendering', () => {
    it('should render project id, name input, and type selector', () => {
      // Act
      renderWithProviders(<UploadDocumentModal {...defaultProps} />)

      // Assert
      expect(screen.getByLabelText(/文档名称/i)).toBeInTheDocument()
      expect(
        screen.getByRole('combobox', { name: /文档类型/i })
      ).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /确定/i })).toBeInTheDocument()
    })
  })

  describe('validation', () => {
    it('should show validation error when name is empty', async () => {
      const user = userEvent.setup()

      // Act
      renderWithProviders(<UploadDocumentModal {...defaultProps} />)

      const submitButton = screen.getByRole('button', { name: /确定/i })
      await user.click(submitButton)

      // Assert
      await waitFor(() => {
        expect(screen.getByText(/请输入文档名称/i)).toBeInTheDocument()
      })
    })
  })

  describe('submission', () => {
    it('should close modal on successful upload', async () => {
      // Arrange
      server.use(
        http.post('/api/v1/knowledge/documents', () =>
          HttpResponse.json(
            {
              id: 'doc-new',
              projectId: 'proj-1',
              name: '新文档.pdf',
              type: 'prd',
              status: 'pending',
              chunkCount: 0,
              uploadedBy: 'user-1',
              createdAt: '2026-04-16T10:00:00Z',
              updatedAt: '2026-04-16T10:00:00Z',
            },
            { status: 201 }
          )
        )
      )

      const user = userEvent.setup()

      // Act
      renderWithProviders(<UploadDocumentModal {...defaultProps} />)

      const nameInput = screen.getByLabelText(/文档名称/i)
      await user.clear(nameInput)
      await user.type(nameInput, '新文档.pdf')

      // Select option is already selected by default (prd)
      const submitButton = screen.getByRole('button', { name: /确定/i })
      await user.click(submitButton)

      // Assert
      await waitFor(() => {
        expect(mockOnSuccess).toHaveBeenCalled()
      })
    })

    it('should show loading state during upload', async () => {
      // Arrange
      server.use(
        http.post('/api/v1/knowledge/documents', async () => {
          // Delay response to test loading state
          await new Promise((resolve) => setTimeout(resolve, 100))
          return HttpResponse.json({}, { status: 201 })
        })
      )

      const user = userEvent.setup()

      // Act
      renderWithProviders(<UploadDocumentModal {...defaultProps} />)

      const nameInput = screen.getByLabelText(/文档名称/i)
      await user.clear(nameInput)
      await user.type(nameInput, '测试文档.pdf')

      const submitButton = screen.getByRole('button', { name: /确定/i })
      await user.click(submitButton)

      // Assert - button should have loading class
      expect(submitButton).toHaveClass('arco-btn-loading')
    })
  })
})
