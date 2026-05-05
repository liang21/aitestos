import { useState } from 'react'
import { useForm, useController } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input, Modal, Form, Button } from '@arco-design/web-react'
import { useCreateProject } from '../hooks/useProjects'
import { messageSuccess, messageError } from '@/lib/notification'

const { TextArea } = Input

// Schema validation with Zod
const projectSchema = z.object({
  name: z
    .string({ required_error: '项目名称不能为空' })
    .min(1, '项目名称不能为空'),
  description: z.string().optional(),
})

type ProjectFormValues = z.infer<typeof projectSchema>

interface CreateProjectModalProps {
  visible: boolean
  onCancel: () => void
  onOk?: (data?: ProjectFormValues) => void
}

/**
 * Modal for creating a new project
 */
export function CreateProjectModal({
  visible,
  onCancel,
  onOk,
}: CreateProjectModalProps) {
  const createProject = useCreateProject()
  const [submitAttempted, setSubmitAttempted] = useState(false)

  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ProjectFormValues>({
    resolver: zodResolver(projectSchema),
    mode: 'onBlur',
    defaultValues: {
      name: '',
      description: '',
    },
  })

  const nameField = useController({ name: 'name', control })
  const descriptionField = useController({ name: 'description', control })

  const onSubmit = async (data: ProjectFormValues) => {
    try {
      await createProject.mutateAsync(data)
      messageSuccess('项目创建成功')
      reset()
      setSubmitAttempted(false)
      onOk?.(data)
    } catch {
      messageError('项目创建失败')
    }
  }

  const handleOk = async () => {
    setSubmitAttempted(true)
    await handleSubmit(onSubmit)()
  }

  const handleCancel = () => {
    reset()
    setSubmitAttempted(false)
    onCancel()
  }

  return (
    <Modal
      title="新建项目"
      visible={visible}
      onCancel={handleCancel}
      onOk={handleOk}
      confirmLoading={createProject.isPending}
    >
      <form onSubmit={(e) => e.preventDefault()}>
        <div className="mb-4">
          <label htmlFor="project-name" className="block mb-2">
            <span className="text-red-500 mr-1">*</span>项目名称
          </label>
          <Input
            id="project-name"
            {...nameField.field}
            placeholder="请输入项目名称"
            aria-label="项目名称"
          />
          {submitAttempted && errors.name && (
            <div className="text-red-500 text-sm mt-1">
              {errors.name.message}
            </div>
          )}
        </div>

        <div className="mb-4">
          <label htmlFor="project-description" className="block mb-2">
            项目描述
          </label>
          <TextArea
            id="project-description"
            {...descriptionField.field}
            placeholder="请输入项目描述（可选）"
            maxLength={200}
            autoSize={{ minRows: 3, maxRows: 6 }}
            aria-label="项目描述"
          />
          {submitAttempted && errors.description && (
            <div className="text-red-500 text-sm mt-1">
              {errors.description.message}
            </div>
          )}
        </div>
      </form>
    </Modal>
  )
}
