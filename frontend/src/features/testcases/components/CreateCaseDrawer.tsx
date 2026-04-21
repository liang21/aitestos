/**
 * Create Case Drawer
 * Form for creating new test cases
 */

import { useEffect } from 'react'
import {
  Drawer,
  Form,
  Input,
  Select,
  Button,
  Message,
} from '@arco-design/web-react'
import { useModuleList } from '@/features/modules/hooks/useModules'
import { useCreateTestCase } from '../hooks/useTestCases'
import { ArrayEditor } from '@/components/business/ArrayEditor'
import type { CaseType, Priority } from '@/types/api'

const { TextArea } = Input
const { Option } = Select

interface CreateCaseDrawerProps {
  visible: boolean
  projectId: string
  onClose: () => void
}

const caseTypeOptions = [
  { label: '功能测试', value: 'functionality' },
  { label: '性能测试', value: 'performance' },
  { label: '接口测试', value: 'api' },
  { label: 'UI 测试', value: 'ui' },
  { label: '安全测试', value: 'security' },
]

const priorityOptions = [
  { label: 'P0 紧急', value: 'P0' },
  { label: 'P1 高', value: 'P1' },
  { label: 'P2 中', value: 'P2' },
  { label: 'P3 低', value: 'P3' },
]

const rejectionReasons = [
  { label: '重复', value: 'duplicate' },
  { label: '无关', value: 'irrelevant' },
  { label: '低质量', value: 'low_quality' },
  { label: '其他', value: 'other' },
]

export function CreateCaseDrawer({
  visible,
  projectId,
  onClose,
}: CreateCaseDrawerProps) {
  const [form] = Form.useForm()
  const createMutation = useCreateTestCase()

  // Fetch modules for selection
  const { data: modules } = useModuleList(projectId)

  // Reset form when drawer opens
  useEffect(() => {
    if (visible) {
      form.resetFields()
    }
  }, [visible, form])

  // Handle form submission
  const handleSubmit = async () => {
    try {
      const values = await form.validate()

      // Validate module is selected
      if (!values.moduleId) {
        Message.error('请选择所属模块')
        return
      }

      // Validate title
      if (
        !values.title ||
        values.title.length < 2 ||
        values.title.length > 500
      ) {
        Message.error('请输入用例标题，长度 2-500 字符')
        return
      }

      // Validate steps
      if (!values.steps || values.steps.length < 1) {
        Message.error('至少添加 1 条测试步骤')
        return
      }

      await createMutation.mutateAsync(
        {
          moduleId: values.moduleId,
          title: values.title,
          preconditions: values.preconditions || [],
          steps: values.steps,
          expected: values.expected ? JSON.parse(values.expected) : {},
          caseType: values.caseType || 'functionality',
          priority: values.priority || 'medium',
        },
        {
          onSuccess: () => {
            Message.success('用例创建成功')
            onClose()
          },
          onError: (error) => {
            Message.error(`创建失败：${error.message}`)
          },
        }
      )
    } catch (error) {
      // Form validation failed
    }
  }

  return (
    <Drawer
      title="新建测试用例"
      visible={visible}
      onClose={onClose}
      width={600}
      footer={
        <div className="flex justify-end gap-2">
          <Button onClick={onClose}>取消</Button>
          <Button
            type="primary"
            onClick={handleSubmit}
            loading={createMutation.isPending}
          >
            确认创建
          </Button>
        </div>
      }
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{ preconditions: [], steps: [], expected: '' }}
      >
        {/* Module Selection */}
        <Form.Item
          field="moduleId"
          label="所属模块"
          rules={[{ required: true, message: '请选择所属模块' }]}
        >
          <Select placeholder="请选择所属模块">
            {modules?.data?.map((module) => (
              <Option key={module.id} value={module.id}>
                {module.name} ({module.abbreviation})
              </Option>
            ))}
          </Select>
        </Form.Item>

        {/* Title */}
        <Form.Item
          field="title"
          label="用例标题"
          rules={[
            { required: true, message: '请输入用例标题' },
            {
              minLength: 2,
              maxLength: 500,
              message: '请输入用例标题，长度 2-500 字符',
            },
          ]}
        >
          <Input placeholder="请输入用例标题" />
        </Form.Item>

        {/* Preconditions */}
        <Form.Item field="preconditions" label="前置条件">
          <ArrayEditor placeholder="请输入前置条件" minRows={0} />
        </Form.Item>

        {/* Steps */}
        <Form.Item
          field="steps"
          label="测试步骤"
          rules={[{ required: true, message: '请添加测试步骤' }]}
        >
          <ArrayEditor placeholder="请输入测试步骤" minRows={1} />
        </Form.Item>

        {/* Expected Results */}
        <Form.Item
          field="expected"
          label="预期结果 (JSON 格式)"
          rules={[
            {
              validator: (value, callback) => {
                if (value && value.trim()) {
                  try {
                    JSON.parse(value)
                  } catch {
                    callback('请输入有效的 JSON 格式')
                  }
                }
                callback()
              },
            },
          ]}
        >
          <TextArea
            placeholder='请输入预期结果（JSON 格式，例如：{"step_1": "操作成功"}）'
            rows={4}
          />
        </Form.Item>

        {/* Case Type */}
        <Form.Item
          field="caseType"
          label="用例类型"
          initialValue="functionality"
        >
          <Select options={caseTypeOptions} />
        </Form.Item>

        {/* Priority */}
        <Form.Item field="priority" label="优先级" initialValue="P2">
          <Select options={priorityOptions} />
        </Form.Item>
      </Form>
    </Drawer>
  )
}
