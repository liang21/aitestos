import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ModuleManagePage } from './ModuleManagePage'

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
}

function renderWithProviders(ui: React.ReactElement, projectId = 'proj1') {
  const queryClient = createTestQueryClient()
  const path = `/projects/${projectId}/modules`

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route path="/projects/:projectId/modules" element={ui} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  )
}

describe('ModuleManagePage', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  it('should render modules table', async () => {
    const mockModules = {
      data: [
        { id: '1', projectId: 'proj1', name: '用户模块', abbreviation: 'USR', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
        { id: '2', projectId: 'proj1', name: '订单模块', abbreviation: 'ORD', createdAt: '2024-01-02T00:00:00Z', updatedAt: '2024-01-02T00:00:00Z' },
      ],
      total: 2,
      offset: 0,
      limit: 10,
    }

    server.use(
      http.get('/api/v1/projects/proj1/modules', () => HttpResponse.json(mockModules))
    )

    renderWithProviders(<ModuleManagePage />)

    await waitFor(() => {
      expect(screen.getByText('用户模块')).toBeInTheDocument()
      expect(screen.getByText('订单模块')).toBeInTheDocument()
      expect(screen.getByText('USR')).toBeInTheDocument()
      expect(screen.getByText('ORD')).toBeInTheDocument()
    })
  })

  it('should open create module modal', async () => {
    const user = userEvent.setup()
    const mockModules = { data: [], total: 0, offset: 0, limit: 10 }

    server.use(
      http.get('/api/v1/projects/proj1/modules', () => HttpResponse.json(mockModules))
    )

    renderWithProviders(<ModuleManagePage />)

    const createButton = screen.getByRole('button', { name: /新建模块/i })
    await user.click(createButton)

    await waitFor(() => {
      expect(screen.getByPlaceholderText(/请输入模块名称/i)).toBeInTheDocument()
    })
  })

  it('should show delete confirmation', async () => {
    const mockModules = {
      data: [
        { id: '1', projectId: 'proj1', name: 'TestModule', abbreviation: 'TST', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
      ],
      total: 1,
      offset: 0,
      limit: 10,
    }

    server.use(
      http.get('/api/v1/projects/proj1/modules', () => HttpResponse.json(mockModules))
    )

    renderWithProviders(<ModuleManagePage />)

    await waitFor(() => {
      const deleteButtons = screen.getAllByRole('button', { name: /删除/i })
      expect(deleteButtons.length).toBeGreaterThan(0)
    })
  })
})
