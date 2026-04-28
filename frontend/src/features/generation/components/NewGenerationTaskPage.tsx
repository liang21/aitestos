import { useState, useMemo } from 'react'
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
  Card,
  Alert,
  Modal,
} from '@arco-design/web-react'
import { IconPlus, IconCheckCircle, IconExclamationCircle, IconCloseCircle } from '@arco-design/web-react/icon'
import { useCreateGenerationTask } from '@/features/generation/hooks/useGeneration'
import { useModuleList } from '@/features/modules/hooks/useModules'
import { useDocumentList } from '@/features/documents/hooks/useDocuments'
import { GENERATION_CONFIG } from '@/features/generation/constants'

const { Item: FormItem } = Form
const { Option: SelectOption } = Select
const { TextArea } = Input

/* eslint-disable @typescript-eslint/no-unused-vars */
/**
 * Zod schema for generation task form
 */
const createTaskSchema = z.object({
  moduleId: z.string().min(1, '请选择模块'),
  prompt: z
    .string()
    .min(GENERATION_CONFIG.TASK_CREATION.MIN_PROMPT_LENGTH, `需求描述至少${GENERATION_CONFIG.TASK_CREATION.MIN_PROMPT_LENGTH}个字`)
    .max(GENERATION_CONFIG.TASK_CREATION.MAX_PROMPT_LENGTH, `需求描述最多${GENERATION_CONFIG.TASK_CREATION.MAX_PROMPT_LENGTH}个字`),
  count: z
    .number()
    .int()
    .min(GENERATION_CONFIG.TASK_CREATION.MIN_CASE_COUNT, `用例数量至少${GENERATION_CONFIG.TASK_CREATION.MIN_CASE_COUNT}`)
    .max(GENERATION_CONFIG.TASK_CREATION.MAX_CASE_COUNT, `用例数量最多${GENERATION_CONFIG.TASK_CREATION.MAX_CASE_COUNT}`),
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

type KnowledgeReadiness = 'sufficient' | 'insufficient' | 'empty'

export function NewGenerationTaskPage({
  projectId: propProjectId,
}: NewGenerationTaskPageProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const [advancedExpanded, setAdvancedExpanded] = useState(false)
  const [showInsufficientWarning, setShowInsufficientWarning] = useState(false)
  const createTask = useCreateGenerationTask()

  // Get projectId from prop or router state
  const projectId =
    propProjectId || (location.state as { projectId?: string })?.projectId || ''

  // Fetch modules for selection
  const { data: modulesData } = useModuleList(projectId)
  const modules = modulesData?.data ?? []

  // Fetch documents for knowledge readiness check
  const { data: documentsData } = useDocumentList({
    projectId,
    status: GENERATION_CONFIG.KNOWLEDGE_READINESS.RELEVANT_DOCUMENT_STATUS,
    limit: 100,
  })
  const completedDocuments = documentsData?.data ?? []

  // Calculate knowledge readiness
  const readiness = useMemo((): KnowledgeReadiness => {
    if (completedDocuments.length === 0) return 'empty'
    if (
      completedDocuments.length <
      GENERATION_CONFIG.KNOWLEDGE_READINESS.MIN_DOCUMENTS_FOR_SUFFICIENT
    )
      return 'insufficient'
    return 'sufficient'
  }, [completedDocuments.length])

  const isSubmitDisabled = readiness === 'empty'

  const {
    control,
    handleSubmit,
    getValues,
    formState: { errors, isSubmitting },
  } = useForm<CreateTaskFormData>({
    resolver: zodResolver(createTaskSchema),
    defaultValues: {
      count: GENERATION_CONFIG.TASK_CREATION.DEFAULT_CASE_COUNT,
      caseType: 'functionality',
      priority: 'P1',
      sceneType: 'positive',
    },
  })

  const onSubmit = async (data: CreateTaskFormData) => {
    // Check if knowledge base is insufficient before proceeding
    if (readiness === 'insufficient') {
      setShowInsufficientWarning(true)
      return
    }

    try {
      const result = await createTask.mutateAsync({
        projectId,
        moduleId: data.moduleId,
        prompt: data.prompt,
        count: data.count,
        caseType: data.caseType,
        priority:
          data.priority ||
          GENERATION_CONFIG.TASK_CREATION.DEFAULT_PRIORITY_INSUFFICIENT, // Force lower priority when knowledge is limited
        sceneType: data.sceneType,
      })

      Message.success('生成任务创建成功')
      navigate(`/generation/tasks/${result.id}`)
    } catch (error) {
      Message.error('创建任务失败，请重试')
    }
  }

  const handleConfirmInsufficient = async () => {
    // Get form data and proceed with submission
    const formData = getValues()
    setShowInsufficientWarning(false)

    try {
      const result = await createTask.mutateAsync({
        projectId,
        moduleId: formData.moduleId,
        prompt: formData.prompt,
        count: formData.count,
        caseType: formData.caseType,
        priority: GENERATION_CONFIG.TASK_CREATION.DEFAULT_PRIORITY_INSUFFICIENT, // Force lower priority for insufficient knowledge
        sceneType: formData.sceneType,
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

      {/* Knowledge Readiness Indicator */}
      <Card
        className="mb-6"
        style={{
          borderLeft: `4px solid ${
            readiness === 'sufficient'
              ? '#00b42a'
              : readiness === 'insufficient'
              ? '#ff7d00'
              : '#f53f3f'
          }`,
        }}
      >
        <div className="flex items-center gap-3">
          {readiness === 'sufficient' && (
            <IconCheckCircle style={{ fontSize: 20, color: '#00b42a' }} />
          )}
          {readiness === 'insufficient' && (
            <IconExclamationCircle style={{ fontSize: 20, color: '#ff7d00' }} />
          )}
          {readiness === 'empty' && (
            <IconCloseCircle style={{ fontSize: 20, color: '#f53f3f' }} />
          )}
          <span className="text-sm">
            📄 {completedDocuments.length} 份文档 ·{' '}
            {readiness === 'sufficient' && (
              <span style={{ color: '#00b42a', fontWeight: 500 }}>🟢 就绪</span>
            )}
            {readiness === 'insufficient' && (
              <span style={{ color: '#ff7d00', fontWeight: 500 }}>🟡 内容有限</span>
            )}
            {readiness === 'empty' && (
              <span style={{ color: '#f53f3f', fontWeight: 500 }}>
                🔴 请先上传需求文档
              </span>
            )}
          </span>
        </div>
        {readiness === 'insufficient' && (
          <Alert
            type="warning"
            style={{ marginTop: 12 }}
            content="知识库内容较少，生成质量可能较低。建议上传更多需求文档后再试。"
          />
        )}
        {readiness === 'empty' && (
          <Alert
            type="error"
            style={{ marginTop: 12 }}
            content="暂无需求文档，请先上传 PRD、API 文档或 Figma 设计稿。"
          />
        )}
      </Card>

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

          <div className="mt-4">
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
          </div>

          <FormItem>
            <div className="flex gap-3">
              <Button
                type="primary"
                htmlType="submit"
                loading={isSubmitting || createTask.isPending}
                icon={<IconPlus />}
                disabled={isSubmitting || createTask.isPending || isSubmitDisabled}
              >
                立即生成
              </Button>
              <Button onClick={() => navigate(-1)}>取消</Button>
            </div>
            {isSubmitDisabled && (
              <p className="text-sm text-red-500 mt-2">
                请先上传需求文档后再创建生成任务
              </p>
            )}
          </FormItem>
        </Form>
      </form>

      {/* RAG Degradation Warning Modal */}
      <Modal
        title="知识库内容不足"
        visible={showInsufficientWarning}
        onOk={handleConfirmInsufficient}
        onCancel={() => setShowInsufficientWarning(false)}
        okText="继续生成"
        cancelText="取消"
        focusOk
        style={{ width: 480 }}
      >
        <Alert
          type="warning"
          content="知识库内容不足，生成质量可能较低。建议上传更多需求文档后再试。"
          style={{ marginBottom: 16 }}
        />
        <p className="text-sm text-gray-600">
          继续生成时，系统将使用较低的置信度设置，可能导致生成的测试用例质量下降。
        </p>
      </Modal>
    </div>
  )
}
