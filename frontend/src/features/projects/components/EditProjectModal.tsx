import { useState, useEffect, useId } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm, useController } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Modal, Input, Message, Popconfirm } from '@arco-design/web-react'
import { useUpdateProject, useDeleteProject } from '../hooks/useProjects'
import type { Project } from '@/types/api'

const { TextArea } = Input

// Schema validation with Zod
const projectSchema = z.object({
  name: z
    .string({ required_error: '项目名称不能为空' })
    .min(1, '项目名称不能为空')
    .max(255, '项目名称最多255个字符'),
  prefix: z
    .string({ required_error: '前缀不能为空' })
    .min(2, '前缀至少2位')
    .max(4, '前缀最多4位')
    .regex(/^[A-Z]+$/, '前缀必须是2-4位大写字母'),
  description: z.string().optional(),
})

type ProjectFormValues = z.infer<typeof projectSchema>

interface EditProjectModalProps {
  visible: boolean
  project: Project
  onClose: () => void
  onSuccess: () => void
}

/**
 * Modal for editing or deleting an existing project
 * Reuses the form structure from CreateProjectModal
 */
export function EditProjectModal({
  visible,
  project,
  onClose,
  onSuccess,
}: EditProjectModalProps) {
  const navigate = useNavigate()
  const updateProject = useUpdateProject()
  const deleteProject = useDeleteProject()
  const [submitAttempted, setSubmitAttempted] = useState(false)

  // 生成唯一 ID，避免多实例冲突
  const nameId = useId()
  const prefixId = useId()
  const descriptionId = useId()

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
      prefix: '',
      description: '',
    },
  })

  const nameField = useController({ name: 'name', control })
  const prefixField = useController({ name: 'prefix', control })
  const descriptionField = useController({ name: 'description', control })

  // 同步 project prop 变化到表单
  useEffect(() => {
    if (visible && project) {
      reset({
        name: project.name,
        prefix: project.prefix,
        description: project.description || '',
      })
    }
  }, [visible, project, reset])

  const onSubmit = async (data: ProjectFormValues) => {
    try {
      await updateProject.mutateAsync({ id: project.id, data })
      Message.success('项目更新成功')
      reset()
      setSubmitAttempted(false)
      onSuccess()
      onClose()
    } catch {
      Message.error('项目更新失败')
    }
  }

  const handleOk = async () => {
    setSubmitAttempted(true)
    await handleSubmit(onSubmit)()
  }

  const handleCancel = () => {
    reset()
    setSubmitAttempted(false)
    onClose()
  }

  const handleDelete = async () => {
    try {
      await deleteProject.mutateAsync(project.id)
      Message.success('项目已删除')
      onSuccess()
      // 使用 React Router 导航到项目列表
      navigate('/projects')
    } catch {
      Message.error('项目删除失败')
    }
  }

  return (
    <Modal
      title="编辑项目"
      visible={visible}
      onCancel={handleCancel}
      onOk={handleOk}
      confirmLoading={updateProject.isPending}
      footer={
        <div className="flex items-center justify-between">
          <Popconfirm
            focusLock
            title="确认删除项目"
            content="删除项目后，所有关联的模块和用例将被级联删除，此操作不可恢复。"
            onOk={handleDelete}
            okText="确认删除"
            cancelText="取消"
          >
            <button
              type="button"
              className="arco-btn arco-btn-danger arco-btn-secondary"
              style={{ marginRight: 'auto' }}
            >
              删除项目
            </button>
          </Popconfirm>
          <div className="flex gap-2">
            <button
              type="button"
              className="arco-btn arco-btn-secondary"
              onClick={handleCancel}
            >
              取消
            </button>
            <button
              type="button"
              className="arco-btn arco-btn-primary"
              onClick={handleOk}
              disabled={updateProject.isPending}
            >
              保存
            </button>
          </div>
        </div>
      }
    >
      <form onSubmit={(e) => e.preventDefault()}>
        <div className="mb-4">
          <label htmlFor={`edit-project-name-${nameId}`} className="block mb-2">
            <span className="text-red-500 mr-1">*</span>项目名称
          </label>
          <Input
            id={`edit-project-name-${nameId}`}
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
          <label htmlFor={`edit-project-prefix-${prefixId}`} className="block mb-2">
            <span className="text-red-500 mr-1">*</span>项目前缀
          </label>
          <Input
            id={`edit-project-prefix-${prefixId}`}
            {...prefixField.field}
            placeholder="2-4位大写字母，如：ECO"
            maxLength={4}
            style={{ textTransform: 'uppercase' }}
            aria-label="项目前缀"
          />
          {submitAttempted && errors.prefix && (
            <div className="text-red-500 text-sm mt-1">
              {errors.prefix.message}
            </div>
          )}
        </div>

        <div className="mb-4">
          <label htmlFor={`edit-project-description-${descriptionId}`} className="block mb-2">
            项目描述
          </label>
          <TextArea
            id={`edit-project-description-${descriptionId}`}
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
