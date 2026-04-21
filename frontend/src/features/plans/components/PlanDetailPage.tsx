/**
 * Plan Detail Page
 * Shows plan information, execution stats, and test results
 */

import React from 'react'
import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { usePlanList } from '@/features/testcases/hooks/useTestCases'
import {
  Button,
  Card,
  Table,
  Progress,
  Space,
  Modal,
  Message,
  Popconfirm,
} from '@arco-design/web-react'
import { IconArrowLeft, IconEdit, IconDelete } from '@arco-design/web-react/icon'
import { usePlanDetail, useRecordResult, useDeletePlan } from '../hooks/usePlans'
import { StatusTag } from '@/components/business/StatusTag'
import { ResultRecordModal } from './ResultRecordModal'
import type { ResultStatus } from '@/types/enums'

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

export function PlanDetailPage({ planId: propPlanId }: { planId?: string }) {
  const navigate = useNavigate()
  const { planId: urlPlanId, projectId } = useParams<{ planId?: string; projectId?: string }>()
  const planId = propPlanId || urlPlanId || ''

  // If we have a projectId param, we might need to fetch the plan differently
  // For now, use the planId directly
  const effectiveProjectId = projectId || ''

  const { data: planDetail, isLoading, error, refetch } = usePlanDetail(planId)
  const recordResultMutation = useRecordResult()
  const deletePlanMutation = useDeletePlan()

  const [resultModalVisible, setResultModalVisible] = useState(false)
  const [selectedCaseId, setSelectedCaseId] = useState<string>('')
  const [selectedCaseTitle, setSelectedCaseTitle] = useState<string>('')

  // Handle back button
  const handleBack = () => {
    navigate('/plans')
  }

  // Handle edit button
  const handleEdit = () => {
    navigate(`/plans/${planId}/edit`)
  }

  // Handle delete button
  const handleDelete = async () => {
    try {
      await deletePlanMutation.mutateAsync(planId)
      Message.success('计划已删除')
      navigate('/plans')
    } catch (err) {
      Message.error(`删除失败：${err instanceof Error ? err.message : '未知错误'}`)
    }
  }

  // Handle open result record modal
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

  // Handle submit result
  const handleSubmitResult = async (data: { status: ResultStatus; note?: string }) => {
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
      Message.error(`录入失败：${err instanceof Error ? err.message : '未知错误'}`)
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

  const { stats, cases } = planDetail

  // Calculate pass rate
  const passRate = stats.total > 0 ? (stats.passed / stats.total) * 100 : 0

  // Table columns
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
      width: 100,
      render: (status: ResultStatus | undefined) =>
        status ? (
          <StatusTag status={status} category="case_status" />
        ) : (
          <span className="text-gray-400">未执行</span>
        ),
    },
    {
      title: '备注',
      dataIndex: 'resultNote',
      ellipsis: true,
      render: (note: string | undefined) => note || '-',
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
          录入结果
        </Button>
      ),
    },
  ]

  return (
    <div className="p-6">
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
            <StatusTag status={planDetail.status} category="plan_status" />
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
            <div className="text-2xl font-semibold text-green-600">{stats.passed}</div>
            <div className="text-gray-500 text-sm">通过</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-semibold text-red-600">{stats.failed}</div>
            <div className="text-gray-500 text-sm">失败</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-semibold text-orange-600">{stats.blocked}</div>
            <div className="text-gray-500 text-sm">阻塞</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-semibold text-gray-600">{stats.skipped}</div>
            <div className="text-gray-500 text-sm">跳过</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-semibold text-gray-400">{stats.unexecuted}</div>
            <div className="text-gray-500 text-sm">未执行</div>
          </div>
        </div>

        <div>
          <div className="flex justify-between mb-2">
            <span className="text-gray-500">通过率</span>
            <span className="font-medium">{passRate.toFixed(1)}%</span>
          </div>
          <Progress
            percent={passRate}
            color="#00B42A"
            animation
          />
        </div>
      </Card>

      {/* Test Cases Table */}
      <Card title="用例列表">
        <Table
          columns={columns}
          data={cases}
          rowKey="caseId"
          pagination={false}
          loading={isLoading}
        />
      </Card>

      {/* Result Record Modal */}
      <ResultRecordModal
        visible={resultModalVisible}
        planId={planId}
        caseId={selectedCaseId}
        caseTitle={selectedCaseTitle}
        onClose={handleCloseResultModal}
        onSubmit={handleSubmitResult}
      />
    </div>
  )
}
