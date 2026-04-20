import { Component, type ReactNode } from 'react'
import { Button, Result, Card } from '@arco-design/web-react'

interface ErrorBoundaryState {
  hasError: boolean
  error?: Error
}

interface ErrorBoundaryProps {
  children: ReactNode
  fallback?: ReactNode
  onError?: (error: Error, errorInfo: React.ErrorInfo) => void
}

/**
 * Error Boundary Component
 *
 * Catches JavaScript errors in component tree,
 * logs them, and displays a fallback UI
 */
export class ErrorBoundary extends Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    // Log error to error reporting service
    console.error('ErrorBoundary caught an error:', error, errorInfo)

    // Call custom error handler if provided
    this.props.onError?.(error, errorInfo)
  }

  handleReset = () => {
    this.setState({ hasError: false, error: undefined })
  }

  render() {
    if (this.state.hasError) {
      // Use custom fallback if provided
      if (this.props.fallback) {
        return this.props.fallback
      }

      // Default error UI
      return (
        <div className="flex min-h-screen items-center justify-center bg-gray-50 p-4">
          <Card className="w-full max-w-md">
            <Result
              status="error"
              title="出现了一些问题"
              subTitle={
                this.state.error?.message ||
                '应用程序遇到了意外错误，请稍后重试。'
              }
              extra={[
                <Button
                  type="primary"
                  key="reload"
                  onClick={() => window.location.reload()}
                >
                  刷新页面
                </Button>,
                <Button key="reset" onClick={this.handleReset}>
                  重试
                </Button>,
              ]}
            />
          </Card>
        </div>
      )
    }

    return this.props.children
  }
}

/**
 * AuthErrorBoundary - Specialized error boundary for auth components
 */
export function AuthErrorBoundary({ children }: { children: ReactNode }) {
  return (
    <ErrorBoundary
      onError={(error) => {
        // Log auth-specific errors
        console.error('Auth error:', error)

        // In production, send to error tracking service
        if (import.meta.env.PROD) {
          // TODO: Integrate with error tracking service (e.g., Sentry)
        }
      }}
    >
      {children}
    </ErrorBoundary>
  )
}
