import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Form,
  Input,
  Button,
  Message,
  Card,
  Radio,
} from '@arco-design/web-react'
import { useNavigate, Link } from 'react-router-dom'
import { useRegister } from '@/features/auth/hooks/useAuth'
import { useRateLimit, RateLimitConfig } from '@/lib/hooks/useRateLimit'
import { RateLimiter } from '@/components/RateLimiter'

const { Item: FormItem } = Form

/**
 * Register schema validation
 */
const registerSchema = z.object({
  username: z
    .string()
    .min(3, '用户名至少为 3 个字符')
    .max(32, '用户名最多为 32 个字符'),
  email: z.string().min(1, '请输入邮箱').email('请输入有效的邮箱地址'),
  password: z
    .string()
    .min(8, '密码至少为 8 位字符')
    .max(100, '密码最多为 100 个字符'),
  role: z.enum(['super_admin', 'admin', 'normal'], {
    errorMap: () => '请选择用户角色',
  }),
})

type RegisterFormData = z.infer<typeof registerSchema>

/**
 * RegisterPage Component
 *
 * Handles new user registration with form validation and rate limiting
 */
export function RegisterPage() {
  const navigate = useNavigate()
  const register = useRegister()

  // Rate limiting
  const rateLimit = useRateLimit(RateLimitConfig.REGISTER)

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      username: '',
      email: '',
      password: '',
      role: 'normal' as const,
    },
  })

  const onSubmit = async (data: RegisterFormData) => {
    // Check rate limit before attempting registration
    if (!rateLimit.canAttempt()) {
      return
    }

    try {
      await register.mutateAsync(data)
      // Reset rate limit on successful registration
      rateLimit.recordAttempt(true)

      Message.success('注册成功，请登录')
      navigate('/login')
    } catch (error) {
      // Record failed attempt
      rateLimit.recordAttempt(false)

      Message.error(error instanceof Error ? error.message : '注册失败')
    }
  }

  const isDisabled =
    register.isPending || register.isSuccess || rateLimit.isLocked

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md p-8">
        <div className="mb-6 text-center">
          <h1 className="text-2xl font-bold">注册</h1>
          <p className="text-gray-500">创建您的账号</p>
        </div>

        <RateLimiter
          isLocked={rateLimit.isLocked}
          remainingAttempts={rateLimit.remainingAttempts}
          maxAttempts={RateLimitConfig.REGISTER.maxAttempts}
          remainingTime={rateLimit.remainingTime}
        >
          <Form onSubmit={handleSubmit(onSubmit)} layout="vertical">
            <FormItem
              label="用户名"
              required
              validateStatus={errors.username ? 'error' : undefined}
              help={errors.username?.message}
            >
              <Controller
                name="username"
                control={control}
                render={({ field }) => (
                  <Input
                    {...field}
                    placeholder="请输入用户名"
                    size="large"
                    disabled={isDisabled}
                  />
                )}
              />
            </FormItem>

            <FormItem
              label="邮箱"
              required
              validateStatus={errors.email ? 'error' : undefined}
              help={errors.email?.message}
            >
              <Controller
                name="email"
                control={control}
                render={({ field }) => (
                  <Input
                    {...field}
                    placeholder="请输入邮箱"
                    size="large"
                    disabled={isDisabled}
                  />
                )}
              />
            </FormItem>

            <FormItem
              label="密码"
              required
              validateStatus={errors.password ? 'error' : undefined}
              help={errors.password?.message}
            >
              <Controller
                name="password"
                control={control}
                render={({ field }) => (
                  <Input.Password
                    {...field}
                    placeholder="请输入密码（至少 8 位字符）"
                    size="large"
                    disabled={isDisabled}
                  />
                )}
              />
            </FormItem>

            <FormItem
              label="角色"
              required
              validateStatus={errors.role ? 'error' : undefined}
              help={errors.role?.message}
            >
              <Controller
                name="role"
                control={control}
                render={({ field }) => (
                  <Radio.Group
                    {...field}
                    disabled={isDisabled}
                    options={[
                      { label: '普通用户', value: 'normal' },
                      { label: '管理员', value: 'admin' },
                      { label: '超级管理员', value: 'super_admin' },
                    ]}
                  />
                )}
              />
            </FormItem>

            <Button
              type="primary"
              size="large"
              long
              htmlType="submit"
              loading={register.isPending}
              disabled={isDisabled}
            >
              {rateLimit.isLocked
                ? `请等待 ${rateLimit.remainingTime} 秒`
                : '注册'}
            </Button>

            <div className="mt-4 text-center text-sm text-gray-500">
              已有账号？
              <Link to="/login" className="ml-1 text-blue-500 hover:underline">
                立即登录
              </Link>
            </div>
          </Form>
        </RateLimiter>
      </Card>
    </div>
  )
}
