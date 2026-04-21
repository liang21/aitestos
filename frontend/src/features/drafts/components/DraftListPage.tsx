/**
 * Draft List Page
 * Displays list of draft test cases with batch operations
 */

import { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import {
  Table,
  Button,
  Select,
  Space,
  Modal,
  Message,
} from '@arco-design/web-react'
import { IconCheck, IconClose, IconSync } from '@arco-design/web-react/icon'
import {
  useDraftList,
  useBatchConfirm,
  useRejectDraft,
} from '../hooks/useDrafts'
import { useModuleList } from '@/features/modules/hooks/useModules'
import { StatusTag } from '@/components/business/StatusTag'
import type { CaseDraft, DraftStatus, Module } from '@/types/api'

interface DraftListPageProps {
  projectId?: string
}

const reasonOptions = [
  { label: '重复', value: 'duplicate' },
  { label: '无关', value: 'irrelevant' },
  { label: '低质量', value: 'low_quality' },
  { label: '其他', value: 'other' },
]

export function DraftListPage({ projectId }: DraftListPageProps) {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([])
  const [statusFilter, setStatusFilter] = useState<DraftStatus | ''>(
    (searchParams.get('status') as DraftStatus) || ''
  )
  const [rejectModal, setRejectModal] = useState<{
    visible: boolean
    draftIds: string[]
  }>({ visible: false, draftIds: [] })
  const [batchConfirmModal, setBatchConfirmModal] = useState({
    visible: false,
    targetModuleId: '',
  })

  // Fetch modules for selection
  const { data: modules } = useModuleList(projectId ?? '')

  // Fetch draft list
  const { data, isLoading, refetch } = useDraftList({
    projectId,
    status: statusFilter || undefined,
    offset: 0,
    limit: 100,
  })

  // Mutations
  const batchConfirm = useBatchConfirm()
  const rejectDraft = useRejectDraft()

  // Handle row click - navigate to detail
  const handleRowClick = (record: CaseDraft) => {
    navigate(`/drafts/${record.id}`)
  }

  // Handle batch confirm - open module selection modal
  const handleBatchConfirm = () => {
    if (selectedRowKeys.length === 0) return
    setBatchConfirmModal({ visible: true, targetModuleId: '' })
  }

  // Execute batch confirm with selected module
  const executeBatchConfirm = async () => {
    if (!batchConfirmModal.targetModuleId) {
      Message.warning('请选择目标模块')
      return
    }

    const result = await batchConfirm.mutateAsync({
      draftIds: selectedRowKeys,
      moduleId: batchConfirmModal.targetModuleId,
    })

    if (result.failedCount > 0) {
      Message.warning(
        `成功确认 ${result.successCount} 条，失败 ${result.failedCount} 条`
      )
    } else {
      Message.success(`成功确认 ${result.successCount} 条草稿`)
    }

    setBatchConfirmModal({ visible: false, targetModuleId: '' })
    setSelectedRowKeys([])
  }

  // Handle reject
  const handleReject = async (reason: string, feedback?: string) => {
    const results = await Promise.allSettled(
      rejectModal.draftIds.map((draftId) =>
        rejectDraft.mutateAsync({
          draftId,
          data: { reason: reason as any, feedback },
        })
      )
    )

    const successCount = results.filter((r) => r.status === 'fulfilled').length
    Message.success(`已拒绝 ${successCount} 条草稿`)

    setRejectModal({ visible: false, draftIds: [] })
    setSelectedRowKeys([])
  }

  // Status filter options
  const statusOptions = [
    { label: '全部', value: '' },
    { label: '待处理', value: 'pending' },
    { label: '已确认', value: 'confirmed' },
    { label: '已拒绝', value: 'rejected' },
  ]

  // Table columns
  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      width: 300,
      render: (title: string, record: CaseDraft) => (
        <a onClick={() => handleRowClick(record)}>{title}</a>
      ),
    },
    {
      title: '来源',
      dataIndex: 'projectName',
      width: 120,
      render: (_: unknown, record: CaseDraft) => (
        <span>
          {record.projectName} / {record.moduleName}
        </span>
      ),
    },
    {
      title: '置信度',
      dataIndex: 'aiMetadata',
      width: 100,
      render: (_: unknown, record: CaseDraft) => {
        // In real implementation, confidence would come from aiMetadata
        return <StatusTag status="high" category="confidence" />
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: DraftStatus) => (
        <StatusTag status={status} category="draft_status" />
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      width: 120,
      render: (_: unknown, record: CaseDraft) => (
        <Space>
          <Button
            type="text"
            size="small"
            icon={<IconCheck />}
            onClick={() => navigate(`/drafts/${record.id}`)}
          >
            确认
          </Button>
          <Button
            type="text"
            size="small"
            status="danger"
            icon={<IconClose />}
            onClick={() =>
              setRejectModal({ visible: true, draftIds: [record.id] })
            }
          >
            拒绝
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div className="p-6">
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-2xl font-semibold">草稿箱</h1>
        <Space>
          <Select
            value={statusFilter}
            onChange={(value) => {
              setStatusFilter(value)
              setSearchParams(value ? { status: value } : {})
            }}
            options={statusOptions}
            style={{ width: 120 }}
          />
          <Button icon={<IconSync />} onClick={() => refetch()}>
            刷新
          </Button>
        </Space>
      </div>

      {/* Batch actions */}
      {selectedRowKeys.length > 0 && (
        <div className="mb-4">
          <Space>
            <Button
              type="primary"
              icon={<IconCheck />}
              onClick={handleBatchConfirm}
              loading={batchConfirm.isPending}
            >
              批量确认 ({selectedRowKeys.length})
            </Button>
            <Button
              status="danger"
              icon={<IconClose />}
              onClick={() =>
                setRejectModal({ visible: true, draftIds: selectedRowKeys })
              }
            >
              批量拒绝
            </Button>
          </Space>
        </div>
      )}

      {/* Table */}
      <Table
        loading={isLoading}
        data={data?.data ?? []}
        columns={columns}
        rowKey="id"
        pagination={{
          total: data?.total ?? 0,
          pageSize: 20,
          showTotal: (total) => `共 ${total} 条`,
        }}
        rowSelection={{
          type: 'checkbox',
          selectedRowKeys,
          onChange: (keys) => setSelectedRowKeys(keys as string[]),
        }}
        onRow={(record) => ({
          onClick: () => handleRowClick(record),
          style: { cursor: 'pointer' },
        })}
      />

      {/* Batch confirm module selection modal */}
      <Modal
        title={`批量确认草稿 (${selectedRowKeys.length} 条)`}
        visible={batchConfirmModal.visible}
        onCancel={() =>
          setBatchConfirmModal({ visible: false, targetModuleId: '' })
        }
        onOk={executeBatchConfirm}
        okText="确认"
        cancelText="取消"
        okButtonProps={{ disabled: !batchConfirmModal.targetModuleId }}
        confirmLoading={batchConfirm.isPending}
      >
        <p className="mb-4">请选择目标模块：</p>
        <Select
          value={batchConfirmModal.targetModuleId}
          onChange={(value) =>
            setBatchConfirmModal({
              ...batchConfirmModal,
              targetModuleId: value,
            })
          }
          placeholder="请选择目标模块"
          className="w-full"
          options={modules?.data?.map((m: Module) => ({
            label: `${m.name} (${m.abbreviation})`,
            value: m.id,
          }))}
          notFoundContent="暂无可用模块"
        />
      </Modal>

      {/* Reject modal */}
      <Modal
        title="拒绝草稿"
        visible={rejectModal.visible}
        onCancel={() => setRejectModal({ visible: false, draftIds: [] })}
        onOk={() => {
          const reason = (
            document.querySelector(
              '[name="reject-reason"]'
            ) as HTMLSelectElement
          )?.value
          const feedback = (
            document.querySelector(
              '[name="reject-feedback"]'
            ) as HTMLTextAreaElement
          )?.value
          handleReject(reason, feedback)
        }}
        okText="确认拒绝"
        cancelText="取消"
      >
        <p className="mb-4">请选择拒绝原因：</p>
        <Select
          name="reject-reason"
          options={reasonOptions}
          className="w-full mb-4"
        />
        <p className="mb-2">反馈意见（可选）：</p>
        <textarea
          name="reject-feedback"
          className="w-full p-2 border border-gray-300 rounded"
          rows={3}
          placeholder="请输入具体的反馈意见..."
        />
      </Modal>
    </div>
  )
}
