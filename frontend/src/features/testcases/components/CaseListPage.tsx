/**
 * Case List Page
 * Displays list of test cases with filtering
 */

import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Table, Button, Select, Space, Message } from '@arco-design/web-react'
import { IconPlus } from '@arco-design/web-react/icon'
import { useCaseList } from '../hooks/useTestCases'
import { StatusTag } from '@/components/business/StatusTag'
import { VirtualTable } from '@/components/business/VirtualTable'
import { CreateCaseDrawer } from './CreateCaseDrawer'
import type { TestCase, CaseStatus, CaseType, Priority } from '@/types/api'

const caseTypeOptions = [
  { label: '全部', value: '' },
  { label: '功能测试', value: 'functionality' },
  { label: '性能测试', value: 'performance' },
  { label: '接口测试', value: 'api' },
  { label: 'UI 测试', value: 'ui' },
  { label: '安全测试', value: 'security' },
]

const priorityOptions = [
  { label: '全部', value: '' },
  { label: '紧急', value: 'critical' },
  { label: '高', value: 'high' },
  { label: '中', value: 'medium' },
  { label: '低', value: 'low' },
]

const statusOptions: Array<{ label: string; value: CaseStatus | '' }> = [
  { label: '全部', value: '' },
  { label: '未执行', value: 'unexecuted' },
  { label: '通过', value: 'passed' },
  { label: '失败', value: 'failed' },
  { label: '阻塞', value: 'blocked' },
  { label: '跳过', value: 'skipped' },
]

export function CaseListPage() {
  const navigate = useNavigate()
  const { projectId } = useParams<{ projectId: string }>()
  const [caseTypeFilter, setCaseTypeFilter] = useState<CaseType | ''>('')
  const [priorityFilter, setPriorityFilter] = useState<Priority | ''>('')
  const [statusFilter, setStatusFilter] = useState<CaseStatus | ''>('')
  const [createDrawerVisible, setCreateDrawerVisible] = useState(false)

  // Fetch test cases list
  const { data, isLoading, refetch } = useCaseList(
    projectId || '',
    {
      caseType: caseTypeFilter || undefined,
      priority: priorityFilter || undefined,
      status: statusFilter || undefined,
      offset: 0,
      limit: 100,
    }
  )

  // Handle row click - navigate to detail
  const handleRowClick = (record: TestCase) => {
    navigate(`/projects/${projectId}/cases/${record.id}`)
  }

  // Table columns
  const columns = [
    {
      title: '编号',
      dataIndex: 'number',
      width: 180,
      render: (number: string, record: TestCase) => (
        <a
          href={`/projects/${projectId}/cases/${record.id}`}
          onClick={(e) => e.preventDefault()}
        >
          {number}
        </a>
      ),
    },
    {
      title: '标题',
      dataIndex: 'title',
      width: 300,
    },
    {
      title: '类型',
      dataIndex: 'caseType',
      width: 100,
      render: (caseType: CaseType) => (
        <StatusTag status={caseType} category="case_type" />
      ),
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      width: 100,
      render: (priority: Priority) => (
        <StatusTag status={priority} category="priority" />
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: CaseStatus) => (
        <StatusTag status={status} category="case_status" />
      ),
    },
  ]

  return (
    <div className="p-6">
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-2xl font-semibold">测试用例</h1>
        <Button
          type="primary"
          icon={<IconPlus />}
          onClick={() => setCreateDrawerVisible(true)}
        >
          新建用例
        </Button>
      </div>

      {/* Filter Bar */}
      <div className="mb-4 flex gap-4">
        <Select
          value={caseTypeFilter}
          onChange={(value) => setCaseTypeFilter(value as CaseType | '')}
          options={caseTypeOptions}
          placeholder="类型"
          style={{ width: 120 }}
        />
        <Select
          value={priorityFilter}
          onChange={(value) => setPriorityFilter(value as Priority | '')}
          options={priorityOptions}
          placeholder="优先级"
          style={{ width: 120 }}
        />
        <Select
          value={statusFilter}
          onChange={(value) => setStatusFilter(value as CaseStatus | '')}
          options={statusOptions}
          placeholder="状态"
          style={{ width: 120 }}
        />
      </div>

      {/* Table with virtual scrolling */}
      <VirtualTable
        loading={isLoading}
        data={data?.data ?? []}
        columns={columns}
        rowKey="id"
        enableVirtual={(data?.data ?? []).length > 500}
        pagination={{
          total: data?.total ?? 0,
          pageSize: 20,
          showTotal: (total) => `共 ${total} 条`,
        }}
        onRow={(record) => ({
          onClick: () => handleRowClick(record),
          style: { cursor: 'pointer' },
        })}
      />

      {/* Create Case Drawer */}
      <CreateCaseDrawer
        visible={createDrawerVisible}
        projectId={projectId ?? ''}
        onClose={() => setCreateDrawerVisible(false)}
      />
    </div>
  )
}
