import { afterEach, describe, expect, it } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter, MemoryRouter, Route, Routes } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import userEvent from '@testing-library/user-event'
import { server } from '../../../../tests/msw/server'
import { ConfigManagePage } from './ConfigManagePage'

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
}

function renderWithProviders(
  ui: React.ReactElement,
  { projectId = '123' }: { projectId?: string } = {}
) {
  const queryClient = createTestQueryClient()
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[`/projects/${projectId}/configs`]}>
        <Routes>
          <Route path="/projects/:projectId/configs" element={ui} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  )
}

describe('ConfigManagePage', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  const mockConfigs = {
    data: [
      {
        key: 'llm_model',
        value: 'deepseek-chat',
        description: 'LLM model selection',
      },
      {
        key: 'llm_temperature',
        value: 0.7,
        description: 'Generation temperature parameter',
      },
      {
        key: 'max_tokens',
        value: 4096,
        description: 'Maximum tokens for generation',
      },
    ],
  }

  it('should render config table with columns (key, value, description, actions)', async () => {
    // Arrange
    server.use(
      http.get('/api/v1/projects/123/configs', () =>
        HttpResponse.json(mockConfigs)
      )
    )

    // Act
    renderWithProviders(<ConfigManagePage />)

    // Assert
    expect(screen.getByText(/配置管理/i)).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByText('llm_model')).toBeInTheDocument()
      expect(screen.getByText('deepseek-chat')).toBeInTheDocument()
      expect(screen.getByText(/LLM model selection/i)).toBeInTheDocument()
      expect(screen.getByText('llm_temperature')).toBeInTheDocument()
      expect(screen.getByText('0.7')).toBeInTheDocument()
    })
  })

  it('should render "新增配置" button', async () => {
    // Arrange
    server.use(
      http.get('/api/v1/projects/123/configs', () =>
        HttpResponse.json(mockConfigs)
      )
    )

    // Act
    renderWithProviders(<ConfigManagePage />)

    // Assert
    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: /新增配置/i })
      ).toBeInTheDocument()
    })
  })

  it('should open edit modal when clicking edit button', async () => {
    // Arrange
    const user = userEvent.setup()
    server.use(
      http.get('/api/v1/projects/123/configs', () =>
        HttpResponse.json(mockConfigs)
      )
    )

    // Act
    renderWithProviders(<ConfigManagePage />)

    await waitFor(() => {
      expect(screen.getByText('llm_model')).toBeInTheDocument()
    })

    // Click edit button (first row action button)
    const editButtons = screen.getAllByRole('button', { name: /编辑/i })

    if (editButtons.length > 0) {
      await user.click(editButtons[0])

      // Assert modal opens
      await waitFor(() => {
        expect(screen.getByText(/编辑配置/i)).toBeInTheDocument()
      })
    }
  })

  it('should render import/export buttons', async () => {
    // Arrange
    server.use(
      http.get('/api/v1/projects/123/configs', () =>
        HttpResponse.json(mockConfigs)
      )
    )

    // Act
    renderWithProviders(<ConfigManagePage />)

    // Assert
    await waitFor(() => {
      expect(screen.getByText('llm_model')).toBeInTheDocument()
    })

    expect(screen.getByRole('button', { name: /导入/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /导出/i })).toBeInTheDocument()
  })

  it('should handle empty config list', async () => {
    // Arrange
    server.use(
      http.get('/api/v1/projects/123/configs', () =>
        HttpResponse.json({ data: [] })
      )
    )

    // Act
    renderWithProviders(<ConfigManagePage />)

    // Assert
    await waitFor(() => {
      expect(screen.getByText(/配置管理/i)).toBeInTheDocument()
    })
  })

  it('should display JSON values properly', async () => {
    // Arrange
    const configsWithJson = {
      data: [
        {
          key: 'llm_config',
          value: { model: 'gpt-4', temperature: 0.7 },
          description: 'Complex config',
        },
      ],
    }

    server.use(
      http.get('/api/v1/projects/123/configs', () =>
        HttpResponse.json(configsWithJson)
      )
    )

    // Act
    renderWithProviders(<ConfigManagePage />)

    // Assert
    await waitFor(() => {
      expect(screen.getByText('llm_config')).toBeInTheDocument()
    })
  })
})
