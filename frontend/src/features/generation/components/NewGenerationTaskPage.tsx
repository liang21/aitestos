import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Button,
  Form,
  Select,
  Input,
  Collapse,
  InputNumber,
  Message,
} from '@arco-design/web-react'
import { IconPlus } from '@arco-design/web-react/icon'
import { useCreateGenerationTask } from '@/features/generation/hooks/useGeneration'
import { useModuleList } from '@/features/modules/hooks/useModules'
import type { CaseType, Priority, SceneType } from '@/types/enums'

const { Item: FormItem } = Form
const { Option: SelectOption } = Select
const { TextArea } = Input

/**
 * Zod schema for generation task form
 */
const createTaskSchema = z.object({
  moduleId: z.string().min(1, '请选择模块'),
  prompt: z.string().min(10, '需求描述至少10个字'),
  count: z.number().int().min(1).max(20, '用例数量范围为1-20'),
  caseType: z
    .enum(['functionality', 'performance', 'api', 'ui', 'security'] as const)
    .optional(),
  priority: z.enum(['P0', 'P1', 'P2', 'P3'] as const).optional(),
  sceneType: z.enum(['positive', 'negative', 'boundary'] as const).optional(),
})

type CreateTaskFormData = z.infer<typeof createTaskSchema>

interface NewGenerationTaskPageProps {
  projectId?: string
}

export function NewGenerationTaskPage({
  projectId: propProjectId,
}: NewGenerationTaskPageProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const [advancedExpanded, setAdvancedExpanded] = useState(false)
  const createTask = useCreateGenerationTask()

  // Get projectId from prop or router state
  const projectId =
    propProjectId || (location.state as { projectId?: string })?.projectId || ''

  // Fetch modules for selection
  const { data: modulesData } = useModuleList(projectId)
  const modules = modulesData?.data ?? []

  const {
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<CreateTaskFormData>({
    resolver: zodResolver(createTaskSchema),
    defaultValues: {
      count: 5,
      caseType: 'functionality',
      priority: 'P1',
      sceneType: 'positive',
    },
  })

  const onSubmit = async (data: CreateTaskFormData) => {
    try {
      const result = await createTask.mutateAsync({
        projectId,
        moduleId: data.moduleId,
        prompt: data.prompt,
        count: data.count,
        caseType: data.caseType,
        priority: data.priority,
        sceneType: data.sceneType,
      })

      Message.success('生成任务创建成功')
      navigate(`/generation/tasks/${result.id}`)
    } catch (error) {
      Message.error('创建任务失败，请重试')
    }
  }

  return (
    <div className="max-w-3xl mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold mb-2">新建 AI 生成任务</h1>
        <p className="text-gray-500">描述测试需求，AI 将自动生成测试用例草稿</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)}>
        <Form layout="vertical">
          <FormItem
            label="目标模块"
            required
            validateStatus={errors.moduleId ? 'error' : undefined}
            help={errors.moduleId?.message}
          >
            <Controller
              name="moduleId"
              control={control}
              render={({ field }) => (
                <Select placeholder="请选择模块" {...field}>
                  {modules.map((mod) => (
                    <SelectOption key={mod.id} value={mod.id}>
                      {mod.name}
                    </SelectOption>
                  ))}
                </Select>
              )}
            />
          </FormItem>

          <FormItem
            label="需求描述"
            required
            validateStatus={errors.prompt ? 'error' : undefined}
            help={errors.prompt?.message}
          >
            <Controller
              name="prompt"
              control={control}
              render={({ field }) => (
                <TextArea
                  placeholder="请描述测试需求，例如：测试用户注册功能，包括邮箱验证和密码强度校验"
                  {...field}
                  rows={4}
                  maxLength={500}
                  showWordLimit
                />
              )}
            />
          </FormItem>

          <FormItem
            label="用例数量"
            required
            validateStatus={errors.count ? 'error' : undefined}
            help={errors.count?.message || '生成 1-20 个测试用例'}
          >
            <Controller
              name="count"
              control={control}
              render={({ field }) => (
                <InputNumber
                  {...field}
                  min={1}
                  max={20}
                  precision={0}
                  style={{ width: '100%' }}
                />
              )}
            />
          </FormItem>

          <FormItem label="高级选项">
            <Collapse
              activeKey={advancedExpanded ? ['advanced'] : []}
              onChange={(keys) => setAdvancedExpanded(keys.length > 0)}
              style={{ borderWidth: 0, background: 'transparent' }}
            >
              <Collapse.Item
                key="advanced"
                header="高级选项（场景类型、优先级、用例类型）"
                style={{ background: 'transparent' }}
              >
                <div className="space-y-4">
                  <FormItem label="场景类型" help={errors.sceneType?.message}>
                    <Controller
                      name="sceneType"
                      control={control}
                      render={({ field }) => (
                        <Select {...field} placeholder="选择场景类型">
                          <SelectOption value="positive">正向测试</SelectOption>
                          <SelectOption value="negative">负向测试</SelectOption>
                          <SelectOption value="boundary">边界测试</SelectOption>
                        </Select>
                      )}
                    />
                  </FormItem>

                  <FormItem label="优先级" help={errors.priority?.message}>
                    <Controller
                      name="priority"
                      control={control}
                      render={({ field }) => (
                        <Select {...field} placeholder="选择优先级">
                          <SelectOption value="P0">P0 紧急</SelectOption>
                          <SelectOption value="P1">P1 高</SelectOption>
                          <SelectOption value="P2">P2 中</SelectOption>
                          <SelectOption value="P3">P3 低</SelectOption>
                        </Select>
                      )}
                    />
                  </FormItem>

                  <FormItem label="用例类型" help={errors.caseType?.message}>
                    <Controller
                      name="caseType"
                      control={control}
                      render={({ field }) => (
                        <Select {...field} placeholder="选择用例类型">
                          <SelectOption value="functionality">
                            功能测试
                          </SelectOption>
                          <SelectOption value="performance">
                            性能测试
                          </SelectOption>
                          <SelectOption value="api">API 测试</SelectOption>
                          <SelectOption value="ui">UI 测试</SelectOption>
                          <SelectOption value="security">安全测试</SelectOption>
                        </Select>
                      )}
                    />
                  </FormItem>
                </div>
              </Collapse.Item>
            </Collapse>
          </FormItem>

          <FormItem>
            <div className="flex gap-3">
              <Button
                type="primary"
                htmlType="submit"
                loading={isSubmitting || createTask.isPending}
                icon={<IconPlus />}
                disabled={isSubmitting || createTask.isPending}
              >
                立即生成
              </Button>
              <Button onClick={() => navigate(-1)}>取消</Button>
            </div>
          </FormItem>
        </Form>
      </form>
    </div>
  )
}
