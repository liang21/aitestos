import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { CreateProjectModal } from './CreateProjectModal'

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
}

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = createTestQueryClient()
  return render(
    <QueryClientProvider client={queryClient}>
      {ui}
    </QueryClientProvider>
  )
}

describe('CreateProjectModal', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  it('should render form fields (name, prefix, description)', () => {
    const mockOnCancel = vi.fn()
    renderWithProviders(<CreateProjectModal visible onCancel={mockOnCancel} />)

    expect(screen.getByLabelText(/项目名称/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/项目前缀/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/项目描述/i)).toBeInTheDocument()
  })

  it('should show validation error when name is empty', async () => {
    const user = userEvent.setup()
    const mockOnCancel = vi.fn()
    const mockOnOk = vi.fn()

    renderWithProviders(
      <CreateProjectModal visible onCancel={mockOnCancel} onOk={mockOnOk} />
    )

    const submitButton = screen.getByRole('button', { name: /确定/i })
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/项目名称不能为空/i)).toBeInTheDocument()
    })
  })

  it('should show validation error when prefix is invalid', async () => {
    const user = userEvent.setup()
    const mockOnCancel = vi.fn()
    const mockOnOk = vi.fn()

    renderWithProviders(
      <CreateProjectModal visible onCancel={mockOnCancel} onOk={mockOnOk} />
    )

    const nameInput = screen.getByLabelText(/项目名称/i)
    await user.type(nameInput, 'Test Project')

    const prefixInput = screen.getByLabelText(/项目前缀/i)
    await user.type(prefixInput, 'invalid')

    const submitButton = screen.getByRole('button', { name: /确定/i })
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/前缀必须是2-4位大写字母/i)).toBeInTheDocument()
    })
  })

  it('should close modal after successful submission', async () => {
    const user = userEvent.setup()
    const mockOnCancel = vi.fn()
    const mockOnOk = vi.fn()

    const mockResponse = {
      id: '456',
      name: 'NewProject',
      prefix: 'NEW',
      description: 'New project description',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    }

    server.use(
      http.post('/api/v1/projects', async () => HttpResponse.json(mockResponse, { status: 201 }))
    )

    renderWithProviders(
      <CreateProjectModal visible onCancel={mockOnCancel} onOk={mockOnOk} />
    )

    const nameInput = screen.getByLabelText(/项目名称/i)
    await user.type(nameInput, 'NewProject')

    const prefixInput = screen.getByLabelText(/项目前缀/i)
    await user.type(prefixInput, 'NEW')

    const descInput = screen.getByLabelText(/项目描述/i)
    await user.type(descInput, 'New project description')

    const submitButton = screen.getByRole('button', { name: /确定/i })
    await user.click(submitButton)

    await waitFor(() => {
      expect(mockOnOk).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'NewProject',
          prefix: 'NEW',
          description: 'New project description',
        })
      )
    })
  })
})
