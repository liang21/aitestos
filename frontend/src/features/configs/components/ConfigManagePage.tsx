import { useState, useEffect } from 'react'
import {
  Button,
  Card,
  Table,
  Typography,
  Modal,
  Space,
  Message,
  Input,
  Popconfirm,
} from '@arco-design/web-react'
import {
  IconPlus,
  IconExport,
  IconImport,
  IconEdit,
} from '@arco-design/web-react/icon'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  useConfigList,
  useSetConfig,
  useDeleteConfig,
  useExportConfigs,
  useImportConfigs,
} from '../hooks/useConfigs'
import { useParams, Navigate } from 'react-router-dom'
import { configSchema, type ConfigInput } from '../schema/configSchema'

const { Title, Text } = Typography
const { TextArea } = Input

/**
 * Format config value for display
 */
function formatConfigValue(value: unknown): string {
  if (typeof value === 'object') {
    return JSON.stringify(value, null, 2)
  }
  return String(value)
}

/**
 * Parse JSON value safely with error feedback
 */
function parseJsonValue(value: string): {
  success: boolean
  data: unknown
  error?: string
} {
  try {
    const parsed = JSON.parse(value)
    return { success: true, data: parsed }
  } catch {
    return { success: false, data: value, error: 'JSON 格式错误' }
  }
}

/**
 * Config Modal State
 */
interface ModalState {
  visible: boolean
  editing: { key: string; value: unknown; description?: string } | null
}

/**
 * Config Manage Page
 * Displays and manages project configurations
 */
export function ConfigManagePage() {
  const { projectId } = useParams<{ projectId: string }>()

  const [modalState, setModalState] = useState<ModalState>({
    visible: false,
    editing: null,
  })

  const { data, isLoading } = useConfigList(projectId || '')
  const setConfig = useSetConfig(projectId || '')
  const deleteConfig = useDeleteConfig(projectId || '')
  const exportConfigs = useExportConfigs(projectId || '')
  const importConfigs = useImportConfigs(projectId || '')

  // Conditional render instead of early return (React Hooks rules)
  if (!projectId) {
    return <Navigate to="/projects" replace />
  }

  const columns = [
    {
      title: '配置键',
      dataIndex: 'key',
      key: 'key',
      width: 200,
    },
    {
      title: '配置值',
      dataIndex: 'value',
      key: 'value',
      render: (value: unknown) => (
        <Text ellipsis={{ showTooltip: true }} style={{ maxWidth: 300 }}>
          {formatConfigValue(value)}
        </Text>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      render: (description?: string) => description || '-',
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (
        _: unknown,
        record: { key: string; value: unknown; description?: string }
      ) => (
        <Space>
          <Button
            type="text"
            size="small"
            icon={<IconEdit />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            content="确定要删除此配置吗？"
            onOk={() => handleDelete(record.key)}
          >
            <Button
              type="text"
              status="danger"
              size="small"
              loading={deleteConfig.isPending}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  const handleEdit = (config: {
    key: string
    value: unknown
    description?: string
  }) => {
    setModalState({ visible: true, editing: config })
  }

  const handleDelete = async (key: string) => {
    try {
      await deleteConfig.mutateAsync(key)
      Message.success('配置删除成功')
    } catch {
      Message.error('配置删除失败')
    }
  }

  const handleAdd = () => {
    setModalState({ visible: true, editing: null })
  }

  const handleModalClose = () => {
    setModalState({ visible: false, editing: null })
  }

  // P5: JSON 解析错误提示
  const handleModalOk = async (data: ConfigInput) => {
    const parseResult = parseJsonValue(data.value)

    // 如果值看起来应该是 JSON 对象但解析失败
    if (
      data.value.trim().startsWith('{') &&
      !parseResult.success &&
      parseResult.error
    ) {
      Message.error(`配置值格式有误：${parseResult.error}`)
      return
    }

    try {
      await setConfig.mutateAsync({
        key: data.key,
        value: parseResult.data,
        description: data.description || undefined,
      })
      Message.success(modalState.editing ? '配置更新成功' : '配置创建成功')
      handleModalClose()
    } catch {
      Message.error(modalState.editing ? '配置更新失败' : '配置创建失败')
    }
  }

  const handleExport = async () => {
    try {
      const result = await exportConfigs.mutateAsync()
      const jsonStr = JSON.stringify(result.configs, null, 2)
      const blob = new Blob([jsonStr], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `configs-${projectId}-${Date.now()}.json`
      a.click()
      URL.revokeObjectURL(url)
      Message.success('配置导出成功')
    } catch {
      Message.error('配置导出失败')
    }
  }

  // P1/P0: 完整实现导入功能
  const handleImport = () => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = 'application/json'
    input.onchange = async (e) => {
      const file = (e.target as HTMLInputElement).files?.[0]
      if (!file) return

      try {
        const text = await file.text()
        const configs = JSON.parse(text)

        // 验证格式
        if (!Array.isArray(configs)) {
          throw new Error('配置文件格式错误：必须是数组格式')
        }

        // 验证每个配置项的结构
        if (!configs.every((c) => c.key && typeof c.value !== 'undefined')) {
          throw new Error('配置项结构错误：每项必须包含 key 和 value')
        }

        // 调用导入 API
        const result = await importConfigs.mutateAsync(configs)

        // 显示结果
        Message.success(
          `导入成功：${result.successCount} 个配置${
            result.failedCount > 0 ? `，${result.failedCount} 个失败` : ''
          }`
        )

        // 如果有失败的，显示详情
        if (result.failedCount > 0 && result.errors) {
          result.errors.forEach((err) => {
            Message.error(`配置 ${err.key} 导入失败：${err.error}`)
          })
        }
      } catch (error) {
        Message.error(
          error instanceof Error ? error.message : '配置文件格式错误'
        )
      }
    }
    input.click()
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <Title heading={4}>配置管理</Title>
        <Space>
          <Button
            icon={<IconImport />}
            onClick={handleImport}
            loading={importConfigs.isPending}
          >
            导入
          </Button>
          <Button
            icon={<IconExport />}
            onClick={handleExport}
            loading={exportConfigs.isPending}
          >
            导出
          </Button>
          <Button type="primary" icon={<IconPlus />} onClick={handleAdd}>
            新增配置
          </Button>
        </Space>
      </div>

      <Card>
        <Table
          columns={columns}
          data={data?.data || []}
          loading={isLoading}
          pagination={false}
          rowKey="key"
        />
      </Card>

      <ConfigModal
        visible={modalState.visible}
        config={modalState.editing}
        onCancel={handleModalClose}
        onOk={handleModalOk}
      />
    </div>
  )
}

/**
 * Config Modal Component
 * Uses React Hook Form + Zod for form validation
 */
interface ConfigModalProps {
  visible: boolean
  config: { key: string; value: unknown; description?: string } | null
  onCancel: () => void
  onOk: (data: ConfigInput) => void
}

function ConfigModal({ visible, config, onCancel, onOk }: ConfigModalProps) {
  const isEdit = !!config

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<ConfigInput>({
    resolver: zodResolver(configSchema),
    defaultValues: {
      key: '',
      value: '',
      description: '',
    },
  })

  // Reset form when config changes or modal opens/closes
  useEffect(() => {
    if (visible) {
      reset({
        key: config?.key || '',
        value: formatConfigValue(config?.value || ''),
        description: config?.description || '',
      })
    }
  }, [visible, config, reset])

  const onSubmit = (data: ConfigInput) => {
    onOk(data)
  }

  const handleCancel = () => {
    reset()
    onCancel()
  }

  return (
    <Modal
      title={isEdit ? '编辑配置' : '新增配置'}
      visible={visible}
      onOk={handleSubmit(onSubmit)}
      onCancel={handleCancel}
      okText="确定"
      cancelText="取消"
    >
      <div className="space-y-4">
        <div>
          <div className="mb-2">
            <Text bold>配置键</Text>
          </div>
          <Controller
            name="key"
            control={control}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="请输入配置键，如: llm_model"
                disabled={isEdit}
                maxLength={100}
                status={errors.key ? 'error' : undefined}
              />
            )}
          />
          {errors.key && (
            <Text type="error" className="mt-1">
              {errors.key.message}
            </Text>
          )}
        </div>

        <div>
          <div className="mb-2">
            <Text bold>配置值</Text>
            <Text type="secondary" className="ml-2">
              (支持 JSON 格式)
            </Text>
          </div>
          <Controller
            name="value"
            control={control}
            render={({ field }) => (
              <TextArea
                {...field}
                placeholder='请输入配置值，如: {"model": "gpt-4", "temperature": 0.7}'
                autoSize={{ minRows: 3, maxRows: 6 }}
                status={errors.value ? 'error' : undefined}
              />
            )}
          />
          {errors.value && (
            <Text type="error" className="mt-1">
              {errors.value.message}
            </Text>
          )}
        </div>

        <div>
          <div className="mb-2">
            <Text bold>描述</Text>
            <Text type="secondary" className="ml-2">
              (可选)
            </Text>
          </div>
          <Controller
            name="description"
            control={control}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="请输入配置描述"
                maxLength={200}
                status={errors.description ? 'error' : undefined}
              />
            )}
          />
          {errors.description && (
            <Text type="error" className="mt-1">
              {errors.description.message}
            </Text>
          )}
        </div>
      </div>
    </Modal>
  )
}
