import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ProjectListPage } from './ProjectListPage'

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

describe('ProjectListPage', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  it('should render project table with columns', async () => {
    const mockProjects = {
      data: [
        { id: '1', name: 'ECommerce', prefix: 'ECO', description: '电商平台', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
        { id: '2', name: 'CRMSystem', prefix: 'CRM', description: '客户管理系统', createdAt: '2024-01-02T00:00:00Z', updatedAt: '2024-01-02T00:00:00Z' },
      ],
      total: 2,
      offset: 0,
      limit: 10,
    }

    server.use(
      http.get('/api/v1/projects', () => HttpResponse.json(mockProjects))
    )

    renderWithProviders(<ProjectListPage />)

    expect(screen.getByText('项目管理')).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByText('ECommerce')).toBeInTheDocument()
      expect(screen.getByText('CRMSystem')).toBeInTheDocument()
      expect(screen.getByText('电商平台')).toBeInTheDocument()
    })
  })

  it('should trigger search on input', async () => {
    const user = userEvent.setup()
    const mockProjects = {
      data: [
        { id: '1', name: 'TestProject', prefix: 'TST', description: 'Test', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
      ],
      total: 1,
      offset: 0,
      limit: 10,
    }

    server.use(
      http.get('/api/v1/projects', () => HttpResponse.json(mockProjects))
    )

    renderWithProviders(<ProjectListPage />)

    const searchInput = screen.getByPlaceholderText(/搜索项目/i)
    await user.clear(searchInput)
    await user.type(searchInput, 'TestProject')
    await user.keyboard('{Enter}')

    await waitFor(() => {
      expect(screen.getByText('TestProject')).toBeInTheDocument()
    })
  })

  it('should handle pagination change', async () => {
    const user = userEvent.setup()
    const mockProjectsPage1 = {
      data: Array.from({ length: 10 }, (_, i) => ({
        id: String(i + 1),
        name: `Project ${i + 1}`,
        prefix: `P${String(i + 1).padStart(2, '0')}`,
        description: `Desc ${i + 1}`,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      })),
      total: 25,
      offset: 0,
      limit: 10,
    }

    server.use(
      http.get('/api/v1/projects', ({ request }) => {
        const url = new URL(request.url)
        const offset = url.searchParams.get('offset')
        if (offset === '10') {
          return HttpResponse.json({
            data: [
              { id: '11', name: 'Project 11', prefix: 'P11', description: 'Desc 11', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
            ],
            total: 25,
            offset: 10,
            limit: 10,
          })
        }
        return HttpResponse.json(mockProjectsPage1)
      })
    )

    renderWithProviders(<ProjectListPage />)

    await waitFor(() => {
      expect(screen.getByText('Project 1')).toBeInTheDocument()
    })

    // Click next page button (page 2)
    const page2Button = screen.getByText('2')
    await user.click(page2Button)

    await waitFor(() => {
      expect(screen.getByText('Project 11')).toBeInTheDocument()
    })
  })

  it('should open CreateProjectModal on button click', async () => {
    const user = userEvent.setup()
    const mockProjects = { data: [], total: 0, offset: 0, limit: 10 }

    server.use(
      http.get('/api/v1/projects', () => HttpResponse.json(mockProjects))
    )

    renderWithProviders(<ProjectListPage />)

    const createButton = screen.getByRole('button', { name: /新建项目/i })
    await user.click(createButton)

    // Modal should open (we can check if modal title appears)
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/请输入项目名称/i)).toBeInTheDocument()
    })
  })
})
