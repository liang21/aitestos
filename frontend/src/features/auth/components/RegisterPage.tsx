import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Form, Input, Button, Message, Card, Radio } from '@arco-design/web-react'
import { useNavigate, Link } from 'react-router-dom'
import { useRegister } from '../hooks/useAuth'

/**
 * Register schema validation
 */
const registerSchema = z.object({
  username: z
    .string()
    .min(3, '用户名至少为 3 个字符')
    .max(32, '用户名最多为 32 个字符'),
  email: z
    .string()
    .min(1, '请输入邮箱')
    .email('请输入有效的邮箱地址'),
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
 * Handles new user registration with form validation
 */
export function RegisterPage() {
  const navigate = useNavigate()
  const register = useRegister()

  const form = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      username: '',
      email: '',
      password: '',
      role: 'normal' as const,
    },
  })

  const handleSubmit = async (data: RegisterFormData) => {
    try {
      await register.mutateAsync(data)
      Message.success('注册成功，请登录')
      navigate('/login')
    } catch (error) {
      Message.error(error instanceof Error ? error.message : '注册失败')
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md p-8">
        <div className="mb-6 text-center">
          <h1 className="text-2xl font-bold">注册</h1>
          <p className="text-gray-500">创建您的账号</p>
        </div>

        <Form onSubmit={form.handleSubmit(handleSubmit)} layout="vertical">
          <Form.Item
            field="username"
            label="用户名"
            required
            rules={[
              { required: true, message: '请输入用户名' },
              {
                minLength: 3,
                message: '用户名至少为 3 个字符',
              },
              {
                maxLength: 32,
                message: '用户名最多为 32 个字符',
              },
            ]}
          >
            <Input
              placeholder="请输入用户名"
              size="large"
              {...form.register('username')}
            />
          </Form.Item>

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
              { maxLength: 100, message: '密码最多为 100 个字符' },
            ]}
          >
            <Input.Password
              placeholder="请输入密码（至少 8 位字符）"
              size="large"
              {...form.register('password')}
            />
          </Form.Item>

          <Form.Item
            field="role"
            label="角色"
            required
            rules={[{ required: true, message: '请选择用户角色' }]}
          >
            <Radio.Group
              {...form.register('role')}
              options={[
                { label: '普通用户', value: 'normal' },
                { label: '管理员', value: 'admin' },
                { label: '超级管理员', value: 'super_admin' },
              ]}
            />
          </Form.Item>

          <Button
            type="primary"
            size="large"
            long
            htmlType="submit"
            loading={register.isPending}
            disabled={register.isSuccess}
          >
            注册
          </Button>

          <div className="mt-4 text-center text-sm text-gray-500">
            已有账号？
            <Link to="/login" className="ml-1 text-blue-500 hover:underline">
              立即登录
            </Link>
          </div>
        </Form>
      </Card>
    </div>
  )
}
