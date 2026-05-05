import { type ReactNode } from 'react'
import { Alert, Progress } from '@arco-design/web-react'

interface RateLimiterProps {
  isLocked: boolean
  remainingAttempts: number
  maxAttempts: number
  remainingTime?: number
  children: ReactNode
}

/**
 * RateLimiter Component
 *
 * Displays rate limit status and locks UI when limit is reached
 */
export function RateLimiter({
  isLocked,
  remainingAttempts,
  maxAttempts,
  remainingTime = 0,
  children,
}: RateLimiterProps) {
  if (isLocked) {
    return (
      <Alert
        type="error"
        content={
          <div>
            <div className="font-medium">尝试次数过多</div>
            <div className="text-sm mt-1">{`请等待 ${remainingTime} 秒后再试`}</div>
          </div>
        }
        showIcon
      />
    )
  }

  if (remainingAttempts < maxAttempts) {
    const percentage = (remainingAttempts / maxAttempts) * 100

    return (
      <div className="mb-4">
        <Alert
          type="warning"
          content={
            <div>
              <div className="font-medium">{`剩余尝试次数: ${remainingAttempts}/${maxAttempts}`}</div>
              <div className="text-sm mt-1">请检查您的输入后重试</div>
            </div>
          }
          showIcon
        />
        <Progress
          percent={percentage}
          size="small"
          color={percentage > 50 ? 'green' : percentage > 20 ? 'orange' : 'red'}
        />
      </div>
    )
  }

  return <>{children}</>
}
