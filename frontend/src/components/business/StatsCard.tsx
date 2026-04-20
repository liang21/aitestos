import { Card, Statistic } from '@arco-design/web-react'
import type { ReactNode } from 'react'
import { TrendingDown, TrendingUp } from 'lucide-react'

export interface StatsCardProps {
  /** Card title */
  title: string
  /** Statistic value */
  value: number
  /** Trend direction */
  trend?: 'up' | 'down'
  /** Trend percentage value */
  trendValue?: number
  /** Custom icon component */
  icon?: ReactNode
  /** Additional class name */
  className?: string
  /** Custom style */
  style?: React.CSSProperties
  /** Value formatter (e.g., suffix like "%") */
  suffix?: string
  /** Value precision (decimal places) */
  precision?: number
}

/**
 * Statistics card component for displaying metrics
 * Shows title, value, optional trend arrow, and custom icon
 */
export function StatsCard({
  title,
  value,
  trend,
  trendValue,
  icon,
  className,
  style,
  suffix,
  precision = 0,
}: StatsCardProps) {
  const renderTrend = () => {
    if (trend === undefined || trendValue === undefined) return null

    const TrendIcon = trend === 'up' ? TrendingUp : TrendingDown
    const trendColor = trend === 'up' ? 'text-green-600' : 'text-red-600'

    return (
      <div className={`flex items-center gap-1 text-sm ${trendColor}`}>
        <TrendIcon className="w-4 h-4" />
        <span>{trendValue}%</span>
      </div>
    )
  }

  return (
    <Card className={className} style={style} bordered={false}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="text-gray-500 text-sm mb-2">{title}</div>
          <Statistic
            value={value}
            precision={precision}
            suffix={suffix}
            valueStyle={{
              fontSize: '28px',
              fontWeight: 600,
              color: '#1D2129',
            }}
          />
          {renderTrend()}
        </div>
        {icon && (
          <div className="ml-4 flex-shrink-0">
            <div className="w-12 h-12 rounded-lg bg-purple-50 flex items-center justify-center text-purple-600">
              {icon}
            </div>
          </div>
        )}
      </div>
    </Card>
  )
}
