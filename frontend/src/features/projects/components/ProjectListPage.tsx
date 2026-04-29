import { useState, useEffect } from 'react'
import {
  Button,
  Card,
  Input,
  Typography,
  Modal,
  Space,
} from '@arco-design/web-react'
import { IconPlus } from '@arco-design/web-react/icon'
import { useProjectList, useDeleteProject } from '../hooks/useProjects'
import { CreateProjectModal } from './CreateProjectModal'
import { SearchTable } from '@/components/business/SearchTable'
import { useDebounce } from '@/hooks/useDebounce'

const { Title } = Typography

/**
 * Project List Page
 * Lists all projects with search and pagination
 */
export function ProjectListPage() {
  const [searchInput, setSearchInput] = useState('')
  const debouncedKeywords = useDebounce(searchInput, 300)

  const [searchParams, setSearchParams] = useState({
    keywords: '',
    offset: 0,
    limit: 10,
  })
  const [createModalVisible, setCreateModalVisible] = useState(false)

  // Update search params when debounced keywords change
  useEffect(() => {
    setSearchParams((prev) => ({
      ...prev,
      keywords: debouncedKeywords,
      offset: debouncedKeywords === '' ? 0 : prev.offset, // Reset offset when searching
    }))
  }, [debouncedKeywords])

  const { data, isLoading } = useProjectList(searchParams)
  const deleteProject = useDeleteProject()

  const columns = [
    {
      title: '项目名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '前缀',
      dataIndex: 'prefix',
      key: 'prefix',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
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
      content: '确定要删除此项目吗？',
      onOk: () => deleteProject.mutate(id),
    })
  }

  const handlePageChange = (page: number) => {
    setSearchParams((prev) => ({
      ...prev,
      offset: (page - 1) * prev.limit,
    }))
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <Title heading={4}>项目管理</Title>
        <Button
          type="primary"
          icon={<IconPlus />}
          onClick={() => setCreateModalVisible(true)}
        >
          新建项目
        </Button>
      </div>

      <Card>
        <div className="mb-4">
          <Input.Search
            placeholder="搜索项目"
            value={searchInput}
            onChange={setSearchInput}
            allowClear
          />
        </div>

        <SearchTable
          columns={columns}
          data={data?.data ?? []}
          total={data?.total ?? 0}
          loading={isLoading}
          current={Math.floor(searchParams.offset / searchParams.limit) + 1}
          pageSize={searchParams.limit}
          onPageChange={(page) => {
            setSearchParams((prev) => ({
              ...prev,
              offset: (page - 1) * prev.limit,
            }))
          }}
          emptyText="暂无项目"
        />
      </Card>

      <CreateProjectModal
        visible={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onOk={() => setCreateModalVisible(false)}
      />
    </div>
  )
}
