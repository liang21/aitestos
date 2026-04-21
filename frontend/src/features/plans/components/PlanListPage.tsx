/**
 * Plan List Page
 * Lists all test plans with filtering and pagination
 */

import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Button, Card, Select, Table, Space, Message } from '@arco-design/web-react'
import { IconPlus } from '@arco-design/web-react/icon'
import { usePlanList } from '../hooks/usePlans'
import { SearchTable } from '@/components/business/SearchTable'
import { StatusTag } from '@/components/business/StatusTag'
import type { PlanStatus } from '@/types/enums'

const { Option } = Select

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '进行中', value: 'active' },
  { label: '已完成', value: 'completed' },
  { label: '已归档', value: 'archived' },
]

const statusTextMap: Record<PlanStatus, string> = {
  draft: '草稿',
  active: '进行中',
  completed: '已完成',
  archived: '已归档',
}

export function PlanListPage() {
  const navigate = useNavigate()
  const { projectId } = useParams<{ projectId: string }>()
  const [filters, setFilters] = useState<{
    status?: PlanStatus
    keywords?: string
  }>({})

  // Use provided projectId from route params
  const effectiveProjectId = projectId || 'project-001'

  const { data, isLoading, error } = usePlanList({
    projectId: effectiveProjectId,
    status: filters.status,
    keywords: filters.keywords,
    offset: 0,
    limit: 100,
  })

  // Handle row click
  const handleRowClick = (record: { id: string }) => {
    navigate(`/plans/${record.id}`)
  }

  // Handle create button click
  const handleCreateClick = () => {
    navigate('/plans/new')
  }

  // Handle status filter change
  const handleStatusChange = (value: PlanStatus | undefined) => {
    setFilters({ ...filters, status: value })
  }

  // Handle search
  const handleSearch = (keywords: string) => {
    setFilters({ ...filters, keywords })
  }

  // Table columns
  const columns = [
    {
      title: '计划名称',
      dataIndex: 'name',
      width: 300,
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: PlanStatus) => (
        <StatusTag status={status} category="plan_status" />
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      width: 180,
      render: (date: string) =>
        new Date(date).toLocaleString('zh-CN'),
    },
  ]

  if (error) {
    Message.error(`加载失败：${error.message}`)
    return null
  }

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-semibold">测试计划</h1>
        <Button type="primary" icon={<IconPlus />} onClick={handleCreateClick}>
          新建计划
        </Button>
      </div>

      <Card>
        <Space className="mb-4" size="medium">
          <Select
            placeholder="筛选状态"
            style={{ width: 150 }}
            allowClear
            onChange={handleStatusChange}
            value={filters.status}
          >
            {statusOptions.map((option) => (
              <Option key={option.value} value={option.value}>
                {option.label}
              </Option>
            ))}
          </Select>
        </Space>

        <SearchTable
          loading={isLoading}
          data={data?.data ?? []}
          total={data?.total ?? 0}
          columns={columns}
          onRow={(record) => ({
            onClick: () => handleRowClick(record),
            style: { cursor: 'pointer' },
          })}
          onSearch={handleSearch}
          rowKey="id"
        />
      </Card>
    </div>
  )
}
