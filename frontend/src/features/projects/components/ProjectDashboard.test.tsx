import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../../../tests/msw/server'
import { afterEach, describe, expect, it } from 'vitest'
import { ProjectDashboard } from './ProjectDashboard'

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
}

function renderWithProviders(ui: React.ReactElement, projectId: string = 'proj1') {
  const queryClient = createTestQueryClient()
  return render(
    <QueryClientProvider client={queryClient}>
      {ui}
    </QueryClientProvider>
  )
}

describe('ProjectDashboard', () => {
  afterEach(() => {
    server.resetHandlers()
  })

  it('should render 4 stats cards', async () => {
    const mockStats = {
      totalCases: 150,
      passRate: 88.5,
      coverage: 94.2,
      aiGeneratedCount: 75,
      trend: [
        { date: '2024-01-01', passRate: 82 },
        { date: '2024-01-02', passRate: 85 },
        { date: '2024-01-03', passRate: 88.5 },
      ],
    }

    server.use(
      http.get('/api/v1/projects/proj1/stats', () => HttpResponse.json(mockStats))
    )

    renderWithProviders(<ProjectDashboard projectId="proj1" />)

    await waitFor(() => {
      expect(screen.getByText('总用例数')).toBeInTheDocument()
      expect(screen.getByText('通过率')).toBeInTheDocument()
      expect(screen.getByText('覆盖率')).toBeInTheDocument()
      expect(screen.getByText('AI生成')).toBeInTheDocument()
    })
  })

  it('should render trend chart area', async () => {
    const mockStats = {
      totalCases: 100,
      passRate: 85,
      coverage: 90,
      aiGeneratedCount: 50,
      trend: [
        { date: '2024-01-01', passRate: 80 },
        { date: '2024-01-02', passRate: 85 },
        { date: '2024-01-03', passRate: 88 },
      ],
    }

    server.use(
      http.get('/api/v1/projects/proj1/stats', () => HttpResponse.json(mockStats))
    )

    renderWithProviders(<ProjectDashboard projectId="proj1" />)

    await waitFor(() => {
      expect(screen.getByText('通过率趋势')).toBeInTheDocument()
    })
  })

  it('should render recent tasks list', async () => {
    const mockStats = {
      totalCases: 100,
      passRate: 85,
      coverage: 90,
      aiGeneratedCount: 50,
      trend: [],
    }

    server.use(
      http.get('/api/v1/projects/proj1/stats', () => HttpResponse.json(mockStats))
    )

    renderWithProviders(<ProjectDashboard projectId="proj1" />)

    await waitFor(() => {
      expect(screen.getByText('最近任务')).toBeInTheDocument()
      expect(screen.getByText('暂无最近任务')).toBeInTheDocument()
    })
  })
})
