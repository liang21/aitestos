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
  Message,
} from '@arco-design/web-react'
import { IconPlus } from '@arco-design/web-react/icon'
import {
  useModuleList,
  useDeleteModule,
  useCreateModule,
} from '../hooks/useModules'

const { Title } = Typography
const { Item: FormItem } = Form

/**
 * Module Management Page
 * Lists and manages modules for a project
 */
export function ModuleManagePage() {
  const { projectId = '' } = useParams<{ projectId: string }>()
  const [createModalVisible, setCreateModalVisible] = useState(false)
  const [form] = Form.useForm()

  const { data, isLoading } = useModuleList(projectId)
  const deleteModule = useDeleteModule()
  const createModule = useCreateModule()

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
      render: (_: unknown, record: { id: string }) => (
        <Space>
          <Button type="text" size="small">
            编辑
          </Button>
          <Button
            type="text"
            status="danger"
            size="small"
            onClick={() => handleDelete(record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

  const handleDelete = (id: string) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除此模块吗？',
      onOk: () => deleteModule.mutate({ projectId, id }),
    })
  }

  const handleCreate = async () => {
    try {
      const values = await form.validate()
      await createModule.mutateAsync({ projectId, data: values })
      Message.success('模块创建成功')
      form.reset()
      setCreateModalVisible(false)
    } catch {
      Message.error('模块创建失败')
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
    </div>
  )
}
