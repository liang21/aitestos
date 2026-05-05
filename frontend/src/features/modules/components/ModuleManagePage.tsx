import { useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Button,
  Card,
  Table,
  Typography,
  Modal,
  Space,
  Form,
  Input,
} from '@arco-design/web-react'
import { IconPlus } from '@arco-design/web-react/icon'
import { messageSuccess, messageError } from '@/lib/notification'
import {
  useModuleList,
  useDeleteModule,
  useCreateModule,
  useUpdateModule,
} from '../hooks/useModules'
import type { Module } from '@/types/api'

const { Title } = Typography
const { Item: FormItem } = Form

/**
 * Module Management Page
 * Lists and manages modules for a project
 */
export function ModuleManagePage() {
  const { projectId = '' } = useParams<{ projectId: string }>()
  const [createModalVisible, setCreateModalVisible] = useState(false)
  const [deleteModalVisible, setDeleteModalVisible] = useState(false)
  const [editModalVisible, setEditModalVisible] = useState(false)
  const [moduleToDelete, setModuleToDelete] = useState<{ id: string; name: string } | null>(null)
  const [editingModule, setEditingModule] = useState<Module | null>(null)
  const [form] = Form.useForm()
  const [editForm] = Form.useForm()

  const { data, isLoading } = useModuleList(projectId)
  const deleteModule = useDeleteModule()
  const createModule = useCreateModule()
  const updateModule = useUpdateModule()

  const columns = [
    {
      title: '模块名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '缩写',
      dataIndex: 'abbreviation',
      key: 'abbreviation',
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (date: string) => new Date(date).toLocaleDateString('zh-CN'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: unknown, record: Module) => (
        <Space>
          <Button
            type="text"
            size="small"
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="text"
            status="danger"
            size="small"
            onClick={() => handleDelete(record.id, record.name)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

  const handleDelete = (id: string, name: string) => {
    setModuleToDelete({ id, name })
    setDeleteModalVisible(true)
  }

  const confirmDelete = () => {
    if (moduleToDelete) {
      deleteModule.mutate({ projectId, id: moduleToDelete.id })
      setDeleteModalVisible(false)
      setModuleToDelete(null)
    }
  }

  const cancelDelete = () => {
    setDeleteModalVisible(false)
    setModuleToDelete(null)
  }

  const handleEdit = (module: Module) => {
    setEditingModule(module)
    editForm.setFieldsValue({
      name: module.name,
      abbreviation: module.abbreviation,
    })
    setEditModalVisible(true)
  }

  const cancelEdit = () => {
    setEditModalVisible(false)
    setEditingModule(null)
    editForm.reset()
  }

  const handleEditSubmit = async () => {
    if (!editingModule) return

    try {
      const values = await editForm.validate()
      await updateModule.mutateAsync({
        id: editingModule.id,
        data: values,
      })
      messageSuccess('模块更新成功')
      setEditModalVisible(false)
      setEditingModule(null)
      editForm.reset()
    } catch {
      messageError('模块更新失败')
    }
  }

  const handleCreate = async () => {
    try {
      const values = await form.validate()
      await createModule.mutateAsync({ projectId, data: values })
      messageSuccess('模块创建成功')
      form.reset()
      setCreateModalVisible(false)
    } catch {
      messageError('模块创建失败')
    }
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <Title heading={4}>模块管理</Title>
        <Button
          type="primary"
          icon={<IconPlus />}
          onClick={() => setCreateModalVisible(true)}
        >
          新建模块
        </Button>
      </div>

      <Card>
        <Table
          columns={columns}
          data={data?.data ?? []}
          loading={isLoading}
          rowKey="id"
          pagination={false}
        />
      </Card>

      <Modal
        title="新建模块"
        visible={createModalVisible}
        onCancel={() => {
          form.reset()
          setCreateModalVisible(false)
        }}
        onOk={handleCreate}
        confirmLoading={createModule.isPending}
      >
        <Form form={form} layout="vertical">
          <FormItem
            label="模块名称"
            required
            field="name"
            rules={[{ required: true, message: '模块名称不能为空' }]}
          >
            <Input placeholder="请输入模块名称" />
          </FormItem>

          <FormItem
            label="缩写"
            required
            field="abbreviation"
            rules={[
              { required: true, message: '缩写不能为空' },
              {
                pattern: /^[A-Z]{2,4}$/,
                message: '缩写必须是2-4位大写字母',
              },
            ]}
          >
            <Input
              placeholder="2-4位大写字母"
              maxLength={4}
              style={{ textTransform: 'uppercase' }}
            />
          </FormItem>
        </Form>
      </Modal>

      {/* Edit Module Modal */}
      <Modal
        title="编辑模块"
        visible={editModalVisible}
        onCancel={cancelEdit}
        onOk={handleEditSubmit}
        okText="保存"
        cancelText="取消"
        focus={false}
      >
        <Form form={editForm} layout="vertical">
          <FormItem
            label="模块名称"
            required
            field="name"
            rules={[{ required: true, message: '模块名称不能为空' }]}
          >
            <Input placeholder="请输入模块名称" />
          </FormItem>

          <FormItem
            label="缩写"
            required
            field="abbreviation"
            rules={[
              { required: true, message: '缩写不能为空' },
              {
                pattern: /^[A-Z]{2,4}$/,
                message: '缩写必须是2-4位大写字母',
              },
            ]}
          >
            <Input
              placeholder="2-4位大写字母"
              maxLength={4}
              style={{ textTransform: 'uppercase' }}
            />
          </FormItem>
        </Form>
      </Modal>

      {/* Delete Confirmation Modal */}
      <Modal
        title="确认删除"
        visible={deleteModalVisible}
        onCancel={cancelDelete}
        onOk={confirmDelete}
        okText="确认"
        cancelText="取消"
        focus={false}
      >
        <p>确定要删除模块 <strong>{moduleToDelete?.name}</strong> 吗？此操作无法撤销。</p>
      </Modal>
    </div>
  )
}
