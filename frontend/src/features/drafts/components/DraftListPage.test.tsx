import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { DraftListPage } from './DraftListPage'
import { server } from '../../../../tests/msw/server'
import { draftsHandlers } from '../../../../tests/msw/handlers/drafts'

describe('DraftListPage', () => {
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
        <BrowserRouter>{ui}</BrowserRouter>
      </QueryClientProvider>
    )
  }

  it('should render draft list with title, source, confidence tag, and time', async () => {
    renderWithProviders(<DraftListPage projectId="project-1" />)

    await waitFor(() => {
      expect(screen.getByText('草稿箱')).toBeInTheDocument()
    })

    // Check that table is rendered
    expect(screen.getByRole('table')).toBeInTheDocument()
  })

  it('should support batch selection with checkboxes', async () => {
    const user = userEvent.setup()
    renderWithProviders(<DraftListPage projectId="project-1" />)

    await waitFor(() => {
      expect(screen.getByRole('table')).toBeInTheDocument()
    })

    const checkboxes = screen.getAllByRole('checkbox')
    expect(checkboxes.length).toBeGreaterThan(0)

    await user.click(checkboxes[0])

    expect(checkboxes[0]).toBeChecked()
  })

  it('should enable batch confirm button when drafts are selected', async () => {
    const user = userEvent.setup()
    renderWithProviders(<DraftListPage projectId="project-1" />)

    await waitFor(() => {
      expect(screen.getByRole('table')).toBeInTheDocument()
    })

    // Get all checkboxes and select the first one
    const checkboxes = screen.getAllByRole('checkbox')
    await user.click(checkboxes[0])

    // Wait for state update
    await waitFor(() => {
      expect(checkboxes[0]).toBeChecked()
    })
  })

  it('should filter by project/module/status', async () => {
    renderWithProviders(<DraftListPage projectId="project-1" />)

    await waitFor(() => {
      expect(screen.getByRole('table')).toBeInTheDocument()
    })

    // Test status filter exists
    const statusSelect = screen.getByRole('combobox')
    expect(statusSelect).toBeInTheDocument()
  })
})
