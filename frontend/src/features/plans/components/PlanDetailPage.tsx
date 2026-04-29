/**
 * Plan Detail Page
 * Shows plan information, execution stats, and test results
 * Features: Quick entry (inline select), Batch entry, Undo functionality
 */

import React, { useState, useCallback, useEffect, useRef } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import {
  Button,
  Card,
  Table,
  Progress,
  Space,
  Modal,
  Message,
  Popconfirm,
  Select,
  Checkbox,
} from '@arco-design/web-react'
import {
  IconArrowLeft,
  IconEdit,
  IconDelete,
} from '@arco-design/web-react/icon'
import {
  usePlanDetail,
  useRecordResult,
  useDeleteResult,
  useDeletePlan,
} from '../hooks/usePlans'
import { StatusTag } from '@/components/business/StatusTag'
import { ResultRecordModal } from './ResultRecordModal'
import { buildProjectRoutes } from '@/lib/routes'
import type { ResultStatus } from '@/types/enums'
import type { PlanCase } from '@/types/api'

const resultStatusOptions = [
  { label: '通过', value: 'pass' },
  { label: '失败', value: 'fail' },
  { label: '阻塞', value: 'block' },
  { label: '跳过', value: 'skip' },
]

const resultStatusTextMap: Record<ResultStatus, string> = {
  pass: '通过',
  fail: '失败',
  block: '阻塞',
  skip: '跳过',
}

const resultStatusColorMap: Record<ResultStatus, string> = {
  pass: 'rgb(var(--success-6))',
  fail: 'rgb(var(--danger-6))',
  block: 'rgb(var(--warning-6))',
  skip: 'rgb(var(--text-3))',
}

interface UndoState {
  caseId: string
  previousStatus?: ResultStatus
  timerId: number
  messageId: string
}

export function PlanDetailPage({ planId: propPlanId }: { planId?: string }) {
  const navigate = useNavigate()
  const { planId: urlPlanId, projectId } = useParams<{
    planId?: string
    projectId?: string
  }>()
  const planId = propPlanId || urlPlanId || ''

  // Get project-scoped routes
  const routes = projectId ? buildProjectRoutes(projectId) : null

  const { data: planDetail, isLoading, error, refetch } = usePlanDetail(planId)
  const recordResultMutation = useRecordResult()
  const deleteResultMutation = useDeleteResult()
  const deletePlanMutation = useDeletePlan()

  // Result Modal state
  const [resultModalVisible, setResultModalVisible] = useState(false)
  const [selectedCaseId, setSelectedCaseId] = useState<string>('')
  const [selectedCaseTitle, setSelectedCaseTitle] = useState<string>('')

  // Quick entry state
  const [editingCaseId, setEditingCaseId] = useState<string | null>(null)
  const [flashRows, setFlashRows] = useState<Set<string>>(new Set())
  const [undoStates, setUndoStates] = useState<Map<string, UndoState>>(new Map())

  // Batch entry state
  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([])
  const [batchModalVisible, setBatchModalVisible] = useState(false)
  const [batchResult, setBatchResult] = useState<ResultStatus | ''>('')

  // Refs for timers
  const undoTimersRef = useRef<Map<string, number>>(new Map())

  // Cleanup timers on unmount
  useEffect(() => {
    return () => {
      undoTimersRef.current.forEach((timerId) => clearTimeout(timerId))
      undoTimersRef.current.clear()
    }
  }, [])

  // Handle back button
  const handleBack = () => {
    navigate(routes?.plans.list ?? '/plans')
  }

  // Handle edit button
  const handleEdit = () => {
    // TODO: Implement edit route when available
    Message.info('编辑功能待实现')
    // navigate(routes?.plans.detail(planId) + '/edit')
  }

  // Handle delete button
  const handleDelete = async () => {
    try {
      await deletePlanMutation.mutateAsync(planId)
      Message.success('计划已删除')
      navigate(routes?.plans.list ?? '/plans')
    } catch (err) {
      Message.error(
        `删除失败：${err instanceof Error ? err.message : '未知错误'}`
      )
    }
  }

  // Handle open result record modal (detailed entry)
  const handleOpenResultModal = (caseId: string, caseTitle: string) => {
    setSelectedCaseId(caseId)
    setSelectedCaseTitle(caseTitle)
    setResultModalVisible(true)
  }

  // Handle close result record modal
  const handleCloseResultModal = () => {
    setResultModalVisible(false)
    setSelectedCaseId('')
    setSelectedCaseTitle('')
  }

  // Handle submit result from modal
  const handleSubmitResult = async (data: {
    status: ResultStatus
    note?: string
  }) => {
    try {
      await recordResultMutation.mutateAsync({
        planId,
        data: {
          caseId: selectedCaseId,
          status: data.status,
          note: data.note,
        },
      })
      Message.success('结果录入成功')
      handleCloseResultModal()
      refetch()
    } catch (err) {
      Message.error(
        `录入失败：${err instanceof Error ? err.message : '未知错误'}`
      )
    }
  }

  // Quick entry: handle inline select change
  const handleQuickEntryChange = async (
    caseId: string,
    status: ResultStatus
  ) => {
    const targetCase = planDetail?.cases.find((c) => c.caseId === caseId)
    const previousStatus = targetCase?.resultStatus

    try {
      await recordResultMutation.mutateAsync({
        planId,
        data: { caseId, status },
      })

      // Flash the row
      setFlashRows((prev) => new Set(prev).add(caseId))
      setTimeout(() => {
        setFlashRows((prev) => {
          const next = new Set(prev)
          next.delete(caseId)
          return next
        })
      }, 500)

      // Show toast with undo link
      const messageId = `result-${caseId}-${Date.now()}`
      Message.success({
        id: messageId,
        content: (
          <span className="flex items-center gap-2">
            已录入：{resultStatusTextMap[status]}
            <button
              className="text-white underline hover:no-underline"
              onClick={() => handleUndo(caseId, previousStatus, messageId)}
            >
              撤销
            </button>
          </span>
        ),
        duration: 3000,
        closable: true,
      })

      // Set up undo state with timer
      const timerId = window.setTimeout(() => {
        setUndoStates((prev) => {
          const next = new Map(prev)
          next.delete(caseId)
          return next
        })
        undoTimersRef.current.delete(caseId)
      }, 3000)

      undoTimersRef.current.set(caseId, timerId)
      setUndoStates((prev) =>
        new Map(prev).set(caseId, { caseId, previousStatus, timerId, messageId })
      )

      setEditingCaseId(null)
    } catch (err) {
      Message.error(
        `录入失败：${err instanceof Error ? err.message : '未知错误'}`
      )
    }
  }

  // Undo quick entry result
  const handleUndo = async (
    caseId: string,
    previousStatus: ResultStatus | undefined,
    messageId: string
  ) => {
    // Clear existing timer
    const existingUndo = undoStates.get(caseId)
    if (existingUndo) {
      clearTimeout(existingUndo.timerId)
      undoTimersRef.current.delete(caseId)
    }

    // Close the toast
    Message.clear(messageId)

    try {
      await deleteResultMutation.mutateAsync({ planId, caseId })
      Message.success('已撤销录入')

      setUndoStates((prev) => {
        const next = new Map(prev)
        next.delete(caseId)
        return next
      })
    } catch (err) {
      Message.error(
        `撤销失败：${err instanceof Error ? err.message : '未知错误'}`
      )
    }
  }

  // Batch entry: handle checkbox selection
  const handleRowSelectionChange = (selectedKeys: string[]) => {
    setSelectedRowKeys(selectedKeys)
  }

  // Batch entry: open modal
  const handleOpenBatchModal = () => {
    setBatchResult('')
    setBatchModalVisible(true)
  }

  // Batch entry: close modal
  const handleCloseBatchModal = () => {
    setBatchModalVisible(false)
    setBatchResult('')
  }

  // Batch entry: submit
  const handleSubmitBatchEntry = async () => {
    if (!batchResult) {
      Message.warning('请选择执行结果')
      return
    }

    try {
      // Record result for each selected case
      await Promise.all(
        selectedRowKeys.map((caseId) =>
          recordResultMutation.mutateAsync({
            planId,
            data: { caseId, status: batchResult as ResultStatus },
          })
        )
      )

      Message.success(`已批量录入 ${selectedRowKeys.length} 条结果`)
      handleCloseBatchModal()
      setSelectedRowKeys([])
    } catch (err) {
      Message.error(
        `批量录入失败：${err instanceof Error ? err.message : '未知错误'}`
      )
    }
  }

  if (error) {
    return (
      <div className="p-6">
        <Button onClick={handleBack} type="text">
          ← 返回
        </Button>
        <div className="mt-4 text-red-500">
          加载失败：{error instanceof Error ? error.message : '未知错误'}
        </div>
      </div>
    )
  }

  if (isLoading || !planDetail) {
    return (
      <div className="p-6">
        <Button onClick={handleBack} type="text">
          ← 返回
        </Button>
        <div className="mt-4">加载中...</div>
      </div>
    )
  }

  const { stats, cases, status } = planDetail

  // Calculate pass rate
  const passRate = stats.total > 0 ? (stats.passed / stats.total) * 100 : 0

  // Table columns with quick entry
  const columns = [
    {
      title: '用例编号',
      dataIndex: 'caseNumber',
      width: 180,
    },
    {
      title: '用例标题',
      dataIndex: 'caseTitle',
      ellipsis: true,
    },
    {
      title: '执行状态',
      dataIndex: 'resultStatus',
      width: 120,
      render: (status: ResultStatus | undefined, record: PlanCase) => {
        const isEditing = editingCaseId === record.caseId

        if (isEditing) {
          return (
            <Select
              size="small"
              options={resultStatusOptions}
              value={status}
              placeholder="选择结果"
              onChange={(value: ResultStatus) =>
                handleQuickEntryChange(record.caseId, value)
              }
              onBlur={() => setEditingCaseId(null)}
              autoFocus
              style={{ width: '100%' }}
            />
          )
        }

        return status ? (
          <button
            className="cursor-pointer hover:opacity-80"
            onClick={() => setEditingCaseId(record.caseId)}
          >
            <StatusTag status={status} category="case_status" />
          </button>
        ) : (
          <button
            className="text-gray-400 cursor-pointer hover:text-gray-600"
            onClick={() => setEditingCaseId(record.caseId)}
          >
            未执行
          </button>
        )
      },
    },
    {
      title: '执行人',
      dataIndex: 'executedBy',
      width: 100,
      render: (executedBy: string | undefined) => executedBy || '-',
    },
    {
      title: '执行时间',
      dataIndex: 'executedAt',
      width: 160,
      render: (executedAt: string | undefined) =>
        executedAt
          ? new Date(executedAt).toLocaleString('zh-CN')
          : '-',
    },
    {
      title: '操作',
      width: 100,
      render: (_: unknown, record: { caseId: string; caseTitle: string }) => (
        <Button
          type="text"
          size="small"
          onClick={() => handleOpenResultModal(record.caseId, record.caseTitle)}
        >
          详细录入
        </Button>
      ),
    },
  ]

  // Row class for flash animation
  const getRowClassName = (record: PlanCase) => {
    return flashRows.has(record.caseId) ? 'animate-flash' : ''
  }

  return (
    <div className="p-6">
      <style>{`
        @keyframes flash {
          0% { background-color: transparent; }
          50% { background-color: ${
            batchResult ? resultStatusColorMap[batchResult as ResultStatus] : 'rgba(var(--primary-6), 0.1)'
          }; }
          100% { background-color: transparent; }
        }
        .animate-flash {
          animation: flash 0.5s ease-in-out;
        }
      `}</style>

      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-2">
          <Button onClick={handleBack} type="text" icon={<IconArrowLeft />}>
            返回
          </Button>
          <h1 className="text-2xl font-semibold">{planDetail.name}</h1>
        </div>
        <Space>
          <Button onClick={handleEdit} icon={<IconEdit />}>
            编辑
          </Button>
          <Popconfirm
            title="确认删除此测试计划？"
            onOk={handleDelete}
            okText="确认"
            cancelText="取消"
          >
            <Button status="danger" icon={<IconDelete />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      </div>

      {/* Plan Info */}
      <Card className="mb-6">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <div className="text-gray-500 text-sm">状态</div>
            <StatusTag status={status} category="plan_status" />
          </div>
          <div>
            <div className="text-gray-500 text-sm">创建时间</div>
            <div>{new Date(planDetail.createdAt).toLocaleString('zh-CN')}</div>
          </div>
        </div>
        {planDetail.description && (
          <div className="mt-4">
            <div className="text-gray-500 text-sm mb-1">描述</div>
            <div>{planDetail.description}</div>
          </div>
        )}
      </Card>

      {/* Execution Stats */}
      <Card className="mb-6" title="执行统计">
        <div className="grid grid-cols-6 gap-4 mb-6">
          <div className="text-center">
            <div className="text-2xl font-semibold">{stats.total}</div>
            <div className="text-gray-500 text-sm">总用例</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-semibold text-green-600">
              {stats.passed}
            </div>
            <div className="text-gray-500 text-sm">通过</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-semibold text-red-600">
              {stats.failed}
            </div>
            <div className="text-gray-500 text-sm">失败</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-semibold text-orange-600">
              {stats.blocked}
            </div>
            <div className="text-gray-500 text-sm">阻塞</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-semibold text-gray-600">
              {stats.skipped}
            </div>
            <div className="text-gray-500 text-sm">跳过</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-semibold text-gray-400">
              {stats.unexecuted}
            </div>
            <div className="text-gray-500 text-sm">未执行</div>
          </div>
        </div>

        <div>
          <div className="flex justify-between mb-2">
            <span className="text-gray-500">通过率</span>
            <span className="font-medium">{passRate.toFixed(1)}%</span>
          </div>
          <Progress percent={passRate} color="#00B42A" animation />
        </div>
      </Card>

      {/* Test Cases Table with Batch Entry */}
      <Card
        title="用例列表"
        extra={
          selectedRowKeys.length > 0 && (
            <Button type="primary" onClick={handleOpenBatchModal}>
              批量录入 ({selectedRowKeys.length})
            </Button>
          )
        }
      >
        <Table
          columns={columns}
          data={cases}
          rowKey="caseId"
          pagination={false}
          loading={isLoading}
          rowSelection={
            status === 'active'
              ? {
                  type: 'checkbox',
                  selectedRowKeys,
                  onChange: handleRowSelectionChange,
                }
              : undefined
          }
          rowClassName={getRowClassName}
        />
      </Card>

      {/* Result Record Modal (Detailed Entry) */}
      <ResultRecordModal
        visible={resultModalVisible}
        planId={planId}
        caseId={selectedCaseId}
        caseTitle={selectedCaseTitle}
        onClose={handleCloseResultModal}
        onSubmit={handleSubmitResult}
      />

      {/* Batch Entry Modal */}
      <Modal
        title="批量录入结果"
        visible={batchModalVisible}
        onCancel={handleCloseBatchModal}
        onOk={handleSubmitBatchEntry}
        okText="确认录入"
        cancelText="取消"
      >
        <div className="py-4">
          <div className="mb-4">
            已选择 <span className="font-medium">{selectedRowKeys.length}</span>{' '}
            条用例
          </div>
          <div>
            <label className="block mb-2">执行结果</label>
            <Select
              options={resultStatusOptions}
              value={batchResult}
              onChange={setBatchResult}
              placeholder="请选择执行结果"
              style={{ width: '100%' }}
            />
          </div>
        </div>
      </Modal>
    </div>
  )
}
