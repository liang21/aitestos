import { Card, Grid, Typography, Empty } from '@arco-design/web-react'
import { useProjectStats } from '../hooks/useProjects'
import { StatsCard } from '@/components/business/StatsCard'

const { Row, Col } = Grid
const { Title } = Typography

interface ProjectDashboardProps {
  projectId: string
}

/**
 * Project Dashboard Component
 * Shows project statistics and recent activity
 */
export function ProjectDashboard({ projectId }: ProjectDashboardProps) {
  const { data: stats, isLoading } = useProjectStats(projectId)

  if (!stats && !isLoading) {
    return <Empty description="暂无统计数据" />
  }

  return (
    <div>
      <Title heading={5} className="mb-4">
        项目概览
      </Title>

      <Row gutter={16} className="mb-6">
        <Col span={6}>
          <StatsCard
            title="总用例数"
            value={stats?.totalCases ?? 0}
            loading={isLoading}
          />
        </Col>
        <Col span={6}>
          <StatsCard
            title="通过率"
            value={stats?.passRate ?? 0}
            suffix="%"
            loading={isLoading}
          />
        </Col>
        <Col span={6}>
          <StatsCard
            title="覆盖率"
            value={stats?.coverage ?? 0}
            suffix="%"
            loading={isLoading}
          />
        </Col>
        <Col span={6}>
          <StatsCard
            title="AI生成"
            value={stats?.aiGeneratedCount ?? 0}
            loading={isLoading}
          />
        </Col>
      </Row>

      {stats?.trend && stats.trend.length > 0 && (
        <Card title="通过率趋势" className="mb-6">
          <div className="h-48 flex items-end justify-between gap-2">
            {stats.trend.map((item) => (
              <div key={item.date} className="flex-1 text-center">
                <div
                  className="bg-blue-500 rounded-t mx-auto"
                  style={{
                    height: `${(item.passRate / 100) * 100}%`,
                    maxWidth: '40px',
                  }}
                />
                <div className="text-xs mt-2 text-gray-500">
                  {new Date(item.date).toLocaleDateString('zh-CN', {
                    month: 'numeric',
                    day: 'numeric',
                  })}
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}

      <Card title="最近任务">
        <Empty description="暂无最近任务" />
      </Card>
    </div>
  )
}
