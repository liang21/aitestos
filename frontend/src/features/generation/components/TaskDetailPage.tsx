import { useMemo } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Tag, Table, Button, Spin } from '@arco-design/web-react'
import { IconArrowLeft } from '@arco-design/web-react/icon'
import { useQuery } from '@tanstack/react-query'
import { generationApi } from '@/features/generation/services/generation'
import { StatusTag } from '@/components/business/StatusTag'
import { generationKeys } from '@/features/generation/hooks/useGeneration'

export function TaskDetailPage() {
  const { taskId } = useParams<{ taskId: string }>()
  const navigate = useNavigate()

  // Fetch task details with polling
  const { data: task, isLoading: taskLoading } = useQuery({
    queryKey: generationKeys.task(taskId ?? ''),
    queryFn: () => generationApi.getTask(taskId!),
    enabled: !!taskId,
    refetchInterval: (data) => {
      if (data?.status === 'pending' || data?.status === 'processing') {
        return 3000
      }
      return false
    },
  })

  // Fetch drafts when task is completed
  const { data: drafts = [], isLoading: draftsLoading } = useQuery({
    queryKey: ['generation-task-drafts', taskId],
    queryFn: () => generationApi.getTaskDrafts(taskId!),
    enabled: !!taskId && task?.status === 'completed',
  })

  // Memoize draft columns
  const draftColumns = useMemo(
    () => [
      {
        title: '标题',
        dataIndex: 'title',
        key: 'title',
        ellipsis: true,
      },
      {
        title: '类型',
        dataIndex: 'caseType',
        key: 'caseType',
        width: 100,
        render: (type: string) => <Tag color="arcoblue">{type}</Tag>,
      },
      {
        title: '优先级',
        dataIndex: 'priority',
        key: 'priority',
        width: 80,
        render: (priority: string) => (
          <StatusTag status={priority} category="priority" size="small" />
        ),
      },
      {
        title: '置信度',
        dataIndex: 'aiMetadata',
        key: 'confidence',
        width: 100,
        render: (metadata: { confidence?: string } | undefined) =>
          metadata?.confidence ? (
            <Tag
              color={
                metadata.confidence === 'high'
                  ? 'green'
                  : metadata.confidence === 'medium'
                  ? 'orange'
                  : 'red'
              }
            >
              {metadata.confidence === 'high' ? '高' : metadata.confidence === 'medium' ? '中' : '低'}
            </Tag>
          ) : null,
      },
      {
        title: '操作',
        key: 'action',
        width: 100,
        render: (_: unknown, record: { id: string }) => (
          <Button
            type="text"
            size="small"
            onClick={() => navigate(`/drafts/${record.id}`)}
          >
            查看详情
          </Button>
        ),
      },
    ],
    [navigate]
  )

  if (taskLoading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <Spin size="large" />
      </div>
    )
  }

  if (!task) {
    return (
      <div className="p-6">
        <div className="text-red-500">任务不存在</div>
      </div>
    )
  }

  const isProcessing = task.status === 'pending' || task.status === 'processing'

  return (
    <div className="p-6">
      <div className="mb-4">
        <Button
          icon={<IconArrowLeft />}
          onClick={() => navigate(-1)}
          type="text"
        >
          返回列表
        </Button>
      </div>

      <Card className="mb-4">
        <Descriptions title="任务详情" column={2}>
          <Descriptions.Item label="任务ID">{task.id}</Descriptions.Item>
          <Descriptions.Item label="状态">
            <StatusTag status={task.status} category="task_status" />
          </Descriptions.Item>
          <Descriptions.Item label="需求描述" span={2}>
            {task.prompt}
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">
            {new Date(task.createdAt).toLocaleString('zh-CN')}
          </Descriptions.Item>
          <Descriptions.Item label="更新时间">
            {new Date(task.updatedAt).toLocaleString('zh-CN')}
          </Descriptions.Item>
          {task.result && (
            <>
              <Descriptions.Item label="生成草稿数">
                {task.result.draftCount}
              </Descriptions.Item>
              {task.result.error && (
                <Descriptions.Item label="错误信息">
                  <Tag color="red">{task.result.error}</Tag>
                </Descriptions.Item>
              )}
            </>
          )}
        </Descriptions>

        {isProcessing && (
          <div className="mt-4 flex items-center gap-2">
            <Spin />
            <span className="text-gray-500">
              {task.status === 'pending' ? '任务等待中...' : 'AI 正在生成用例，请稍候...'}
            </span>
          </div>
        )}
      </Card>

      {task.status === 'completed' && (
        <Card title="生成的草稿" loading={draftsLoading}>
          <Table
            columns={draftColumns}
            data={drafts}
            rowKey="id"
            pagination={false}
            onRow={(record) => ({
              onClick: () => navigate(`/drafts/${record.id}`),
              style: { cursor: 'pointer' },
            })}
          />
        </Card>
      )}
    </div>
  )
}
