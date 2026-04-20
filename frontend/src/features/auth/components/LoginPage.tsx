import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Form, Input, Button, Message, Card } from '@arco-design/web-react'
import { useNavigate, Link } from 'react-router-dom'
import { useLogin } from '../hooks/useAuth'
import { useAuthStore } from '../hooks/useAuthStore'

/**
 * Login schema validation
 */
const loginSchema = z.object({
  email: z
    .string()
    .min(1, '请输入邮箱')
    .email('请输入有效的邮箱地址'),
  password: z
    .string()
    .min(8, '密码至少为 8 位字符')
})

type LoginFormData = z.infer<typeof loginSchema>

/**
 * LoginPage Component
 *
 * Handles user authentication with form validation
 */
export function LoginPage() {
  const navigate = useNavigate()
  const login = useLogin()
  const { isAuthenticated } = useAuthStore()

  // Redirect if already authenticated
  if (isAuthenticated) {
    navigate('/projects')
    return null
  }

  const form = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  })

  const handleSubmit = async (data: LoginFormData) => {
    try {
      await login.mutateAsync(data)
      // Store will be updated by useLogin onSuccess
      Message.success('登录成功')
      navigate('/projects')
    } catch (error) {
      // Error is handled by useLogin onError
      Message.error(error instanceof Error ? error.message : '登录失败')
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md p-8">
        <div className="mb-6 text-center">
          <h1 className="text-2xl font-bold">登录</h1>
          <p className="text-gray-500">欢迎回到 Aitestos</p>
        </div>

        <Form onSubmit={form.handleSubmit(handleSubmit)} layout="vertical">
          <Form.Item
            field="email"
            label="邮箱"
            required
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input
              placeholder="请输入邮箱"
              size="large"
              {...form.register('email')}
            />
          </Form.Item>

          <Form.Item
            field="password"
            label="密码"
            required
            rules={[
              { required: true, message: '请输入密码' },
              { minLength: 8, message: '密码至少为 8 位字符' },
            ]}
          >
            <Input.Password
              placeholder="请输入密码"
              size="large"
              {...form.register('password')}
            />
          </Form.Item>

          <Button
            type="primary"
            size="large"
            long
            htmlType="submit"
            loading={login.isPending}
            disabled={login.isSuccess}
          >
            登录
          </Button>

          <div className="mt-4 text-center text-sm text-gray-500">
            还没有账号？
            <Link to="/register" className="ml-1 text-blue-500 hover:underline">
              立即注册
            </Link>
          </div>
        </Form>
      </Card>
    </div>
  )
}
