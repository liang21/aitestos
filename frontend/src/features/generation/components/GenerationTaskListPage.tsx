import { useState, useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Card, Select, Button, Tag } from '@arco-design/web-react'
import { IconPlus } from '@arco-design/web-react/icon'
import { useGenerationTasks } from '@/features/generation/hooks/useGeneration'
import { StatusTag } from '@/components/business/StatusTag'
import type { TaskStatus } from '@/types/enums'

const { Option: SelectOption } = Select

interface GenerationTaskListPageProps {
  projectId: string
}

export function GenerationTaskListPage({
  projectId,
}: GenerationTaskListPageProps) {
  const navigate = useNavigate()
  const [statusFilter, setStatusFilter] = useState<TaskStatus | ''>('')

  const { data: tasksData, isLoading } = useGenerationTasks({
    projectId,
    ...(statusFilter ? { status: statusFilter } : {}),
    offset: 0,
    limit: 20,
  })

  const tasks = tasksData?.data ?? []

  // Memoize columns to prevent re-creation on every render
  const columns = useMemo(
    () => [
      {
        title: '需求描述',
        dataIndex: 'prompt',
        key: 'prompt',
        ellipsis: true,
        render: (prompt: string, record: { id: string }) => (
          <a
            onClick={() => navigate(`/generation/tasks/${record.id}`)}
            className="text-blue-500 hover:underline cursor-pointer"
          >
            {prompt.length > 50 ? `${prompt.slice(0, 50)}...` : prompt}
          </a>
        ),
      },
      {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
        width: 100,
        render: (status: TaskStatus) => (
          <StatusTag status={status} category="task_status" size="small" />
        ),
      },
      {
        title: '草稿数',
        dataIndex: 'result',
        key: 'draftCount',
        width: 80,
        render: (result: { draftCount?: number } | null) => (
          <span>{result?.draftCount ?? '-'}</span>
        ),
      },
      {
        title: '创建时间',
        dataIndex: 'createdAt',
        key: 'createdAt',
        width: 180,
        render: (date: string) => new Date(date).toLocaleString('zh-CN'),
      },
    ],
    [navigate]
  )

  return (
    <div className="p-6">
      <Card>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">AI 生成任务</h2>
          <div className="flex gap-3">
            <Select
              placeholder="筛选状态"
              style={{ width: 150 }}
              value={statusFilter}
              onChange={setStatusFilter}
              allowClear
            >
              <SelectOption value="pending">待处理</SelectOption>
              <SelectOption value="processing">生成中</SelectOption>
              <SelectOption value="completed">已完成</SelectOption>
              <SelectOption value="failed">失败</SelectOption>
            </Select>
            <Button
              type="primary"
              icon={<IconPlus />}
              onClick={() =>
                navigate('/generation/tasks/new', { state: { projectId } })
              }
            >
              新建任务
            </Button>
          </div>
        </div>

        <Table
          columns={columns}
          data={tasks}
          loading={isLoading}
          rowKey="id"
          pagination={{
            total: tasksData?.total ?? 0,
            pageSize: 20,
            showTotal: (total) => `共 ${total} 条`,
          }}
          onRow={(record) => ({
            onClick: () => navigate(`/generation/tasks/${record.id}`),
            style: { cursor: 'pointer' },
          })}
        />
      </Card>
    </div>
  )
}
