import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Form, Input, Button, Message } from '@arco-design/web-react'
import { IconEmail, IconLock } from '@arco-design/web-react/icon'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { z } from 'zod'
import { useLogin } from '@/features/auth/hooks/useAuth'
import { useRateLimit, RateLimitConfig } from '@/lib/hooks/useRateLimit'
import { RateLimiter } from '@/components/RateLimiter'

// Login schema with email validation
const loginSchema = z.object({
  email: z.string().min(1, '请输入邮箱').email('请输入有效的邮箱地址'),
  password: z.string().min(1, '请输入密码'),
})

type LoginInput = z.infer<typeof loginSchema>

interface LoginFormProps {
  onSuccess?: () => void
}

/**
 * Login form component - Email + Password authentication
 */
export function LoginForm({ onSuccess }: LoginFormProps) {
  const login = useLogin()
  const rateLimit = useRateLimit(RateLimitConfig.LOGIN)
  const location = useLocation()
  const navigate = useNavigate()

  // Get redirect destination from router state
  const from = (location.state as { from?: string })?.from || '/projects'

  const form = useForm<LoginInput>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  })

  const handleSubmit = async (data: LoginInput) => {
    if (!rateLimit.canAttempt()) {
      return
    }

    try {
      await login.mutateAsync(data)
      rateLimit.recordAttempt(true)
      Message.success('登录成功')

      // Call onSuccess callback or navigate to redirect destination
      if (onSuccess) {
        onSuccess()
      } else {
        navigate(from, { replace: true })
      }
    } catch (error) {
      rateLimit.recordAttempt(false)
      Message.error(error instanceof Error ? error.message : '登录失败')
    }
  }

  const isDisabled = login.isPending || login.isSuccess || rateLimit.isLocked

  return (
    <div className="login-form-wrap">
      {/* Logo & Title */}
      <div className="login-title">
        <div className="logo-container">
          <img src="/favicon.svg" alt="Aitestos" className="logo-image" />
        </div>
        <div className="title-container">
          <span className="title-welcome">AI 测试管理平台</span>
        </div>
      </div>

      {/* Form Card */}
      <RateLimiter
        isLocked={rateLimit.isLocked}
        remainingAttempts={rateLimit.remainingAttempts}
        maxAttempts={RateLimitConfig.LOGIN.maxAttempts}
        remainingTime={rateLimit.remainingTime}
      >
        <div className="form-card">
          <div className="form-header">账号登录</div>

          <Form
            onSubmit={form.handleSubmit(handleSubmit)}
            layout="vertical"
            className="login-form"
          >
            <Form.Item
              field="email"
              rules={[{ required: true, message: '请输入邮箱' }]}
              className="login-form-item"
            >
              <Input
                {...(form.register('email') as object)}
                prefix={<IconEmail />}
                placeholder="请输入邮箱"
                size="large"
                maxLength={64}
                disabled={isDisabled}
                className="login-input"
              />
            </Form.Item>

            <Form.Item
              field="password"
              rules={[{ required: true, message: '请输入密码' }]}
              className="login-form-item"
            >
              <Input.Password
                {...(form.register('password') as object)}
                prefix={<IconLock />}
                placeholder="请输入密码"
                size="large"
                maxLength={64}
                disabled={isDisabled}
                allowClear
                className="login-password-input"
              />
            </Form.Item>

            <div className="form-actions">
              <Button
                type="primary"
                size="large"
                long
                htmlType="submit"
                loading={login.isPending}
                disabled={isDisabled}
              >
                {rateLimit.isLocked
                  ? `请等待 ${rateLimit.remainingTime} 秒`
                  : '登录'}
              </Button>
            </div>

            {/* Register Link */}
            <div className="form-footer">
              <span className="footer-text">还没有账号？</span>
              <Link to="/register" className="register-link">
                立即注册
              </Link>
            </div>
          </Form>
        </div>
      </RateLimiter>
    </div>
  )
}
