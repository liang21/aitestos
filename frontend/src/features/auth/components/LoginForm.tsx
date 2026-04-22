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
    <div className="flex flex-col items-center justify-center h-screen flex-1 py-15 px-10 bg-gradient-to-br from-gray-50 to-gray-100">
      {/* Logo & Title */}
      <div className="text-center mb-10 animate-fade-in-up">
        <div className="flex justify-center mb-5">
          <img
            src="/favicon.svg"
            alt="Aitestos"
            className="w-35 h-auto object-contain drop-shadow-[0_4px_12px_rgba(139,92,246,0.2)]"
          />
        </div>
        <div className="">
          <span className="text-purple-700 text-xl font-semibold tracking-wide">
            AI 测试管理平台
          </span>
        </div>
      </div>

      {/* Form Card */}
      <RateLimiter
        isLocked={rateLimit.isLocked}
        remainingAttempts={rateLimit.remainingAttempts}
        maxAttempts={RateLimitConfig.LOGIN.maxAttempts}
        remainingTime={rateLimit.remainingTime}
      >
        <div className="min-w-120 px-12 py-12 rounded-2xl bg-white shadow-[0_20px_40px_rgba(0,0,0,0.08),0_8px_16px_rgba(0,0,0,0.04)] border border-purple-500/10 relative animate-slide-up hover:shadow-[0_25px_50px_rgba(0,0,0,0.12),0_12px_24px_rgba(0,0,0,0.08)] transition-shadow duration-300">
          <div className="text-xl font-semibold text-purple-700 text-center mb-8 tracking-wide">
            账号登录
          </div>

          <Form
            onSubmit={form.handleSubmit(handleSubmit)}
            layout="vertical"
            className="login-form"
          >
            <Form.Item
              field="email"
              rules={[{ required: true, message: '请输入邮箱' }]}
              className="mb-6"
            >
              <Input
                {...(form.register('email') as object)}
                prefix={<IconEmail />}
                placeholder="请输入邮箱"
                size="large"
                maxLength={64}
                disabled={isDisabled}
                className="login-input-rounded"
              />
            </Form.Item>

            <Form.Item
              field="password"
              rules={[{ required: true, message: '请输入密码' }]}
              className="mb-6"
            >
              <Input.Password
                {...(form.register('password') as object)}
                prefix={<IconLock />}
                placeholder="请输入密码"
                size="large"
                maxLength={64}
                disabled={isDisabled}
                allowClear
                className="login-input-rounded"
              />
            </Form.Item>

            <div className="mt-8 mb-0">
              <Button
                type="primary"
                size="large"
                long
                htmlType="submit"
                loading={login.isPending}
                disabled={isDisabled}
                className="!h-12 !text-base !font-semibold !bg-gradient-to-r !from-purple-500 !to-purple-700 !border-0 !rounded-3xl !shadow-[0_4px_12px_rgba(139,92,246,0.3)] hover:!bg-gradient-to-r hover:!from-purple-600 hover:!to-purple-800 hover:!shadow-[0_8px_20px_rgba(139,92,246,0.4)] hover:!-translate-y-[-2px] active:!translate-y-0 active:!shadow-[0_2px_8px_rgba(139,92,246,0.3)]"
              >
                {rateLimit.isLocked
                  ? `请等待 ${rateLimit.remainingTime} 秒`
                  : '登录'}
              </Button>
            </div>

            {/* Register Link */}
            <div className="flex justify-center items-center gap-2 mt-6 pt-6 border-t border-gray-100">
              <span className="text-sm text-gray-500">还没有账号？</span>
              <Link
                to="/register"
                className="text-sm text-purple-500 font-medium no-underline hover:text-purple-600 hover:underline"
              >
                立即注册
              </Link>
            </div>
          </Form>
        </div>
      </RateLimiter>
    </div>
  )
}
