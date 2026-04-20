import { Form, Input, Select, Modal } from '@arco-design/web-react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useUploadDocument } from '../hooks/useDocuments'
import type { DocumentType } from '@/types/api'

const { Option } = Select

// Document type options
const DOCUMENT_TYPE_OPTIONS = [
  { label: 'PRD', value: 'prd' },
  { label: 'API Spec', value: 'api_spec' },
  { label: 'Swagger', value: 'swagger' },
  { label: 'Figma', value: 'figma' },
  { label: 'Markdown', value: 'markdown' },
]

// Validation schema
const uploadDocumentSchema = z.object({
  name: z.string().min(2, '请输入文档名称').max(255, '文档名称最多 255 个字符'),
  type: z.enum(['prd', 'api_spec', 'swagger', 'figma', 'markdown'], {
    errorMap: () => ({ message: '请选择文档类型' }),
  }),
})

export type UploadDocumentFormData = z.infer<typeof uploadDocumentSchema>

export interface UploadDocumentModalProps {
  visible: boolean
  projectId: string
  onCancel: () => void
  onSuccess?: () => void
}

export function UploadDocumentModal({
  visible,
  projectId,
  onCancel,
  onSuccess,
}: UploadDocumentModalProps) {
  const uploadDocument = useUploadDocument()
  const form = useForm<UploadDocumentFormData>({
    resolver: zodResolver(uploadDocumentSchema),
    defaultValues: {
      name: '',
      type: 'prd',
    },
  })

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = form

  // Handle submit
  const onSubmit = async (data: UploadDocumentFormData) => {
    await uploadDocument.mutateAsync({
      projectId,
      name: data.name,
      type: data.type,
    })
    reset()
    onSuccess?.()
  }

  // Handle modal close
  const handleCancel = () => {
    reset()
    onCancel()
  }

  return (
    <Modal
      title="上传文档"
      visible={visible}
      onCancel={handleCancel}
      onOk={handleSubmit(onSubmit)}
      okText="确定"
      cancelText="取消"
      confirmLoading={uploadDocument.isPending}
      style={{ width: 500 }}
    >
      <Form layout="vertical" autoComplete="off">
        <Form.Item
          label="文档名称"
          required
          validateStatus={errors.name ? 'error' : undefined}
          help={errors.name?.message}
        >
          <Controller
            name="name"
            control={control}
            render={({ field }) => (
              <Input
                placeholder="请输入文档名称"
                aria-label="文档名称"
                {...field}
              />
            )}
          />
        </Form.Item>

        <Form.Item
          label="文档类型"
          required
          validateStatus={errors.type ? 'error' : undefined}
          help={errors.type?.message}
        >
          <Controller
            name="type"
            control={control}
            render={({ field }) => (
              <Select
                placeholder="请选择文档类型"
                aria-label="文档类型"
                {...field}
              >
                {DOCUMENT_TYPE_OPTIONS.map((opt) => (
                  <Option key={opt.value} value={opt.value}>
                    {opt.label}
                  </Option>
                ))}
              </Select>
            )}
          />
        </Form.Item>
      </Form>
    </Modal>
  )
}
