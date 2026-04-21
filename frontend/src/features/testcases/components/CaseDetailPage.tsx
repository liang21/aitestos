/**
 * Case Detail Page
 * Displays detailed information about a test case
 */

import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Button, Card, Space, Message, Spin } from '@arco-design/web-react'
import { IconEdit, IconDelete } from '@arco-design/web-react/icon'
import { useCaseDetail, useDeleteTestCase } from '../hooks/useTestCases'
import { StatusTag } from '@/components/business/StatusTag'
import { ReferencePanel } from '@/components/business/ReferencePanel'
import type { ReferencedChunk } from '@/types/api'

interface CaseDetailPageProps {
  /** For testing purposes, caseId can be passed as prop */
  caseId?: string
}

export function CaseDetailPage({
  caseId: propCaseId,
}: CaseDetailPageProps = {}) {
  const { caseId: urlCaseId } = useParams<{ caseId: string }>()
  const caseId = propCaseId ?? urlCaseId ?? ''
  const navigate = useNavigate()
  const [deleteModalVisible, setDeleteModalVisible] = useState(false)

  const { data: testCase, isLoading, error } = useCaseDetail(caseId)
  const deleteMutation = useDeleteTestCase()

  // Handle delete
  const handleDelete = async () => {
    if (!testCase) return

    try {
      await deleteMutation.mutateAsync(testCase.id, {
        onSuccess: () => {
          Message.success('用例已删除')
          navigate('/testcases')
        },
      })
    } catch {
      Message.error('删除失败')
    }
  }

  // Handle edit
  const handleEdit = () => {
    navigate(`/testcases/${testCase?.id}/edit`)
  }

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin />
      </div>
    )
  }

  if (error || !testCase) {
    return (
      <div className="p-6 text-center text-gray-500">加载失败，用例不存在</div>
    )
  }

  return (
    <div className="p-6">
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">{testCase.title}</h1>
          <p className="text-gray-500 mt-1">{testCase.number}</p>
        </div>
        <Space>
          <Button icon={<IconEdit />} onClick={handleEdit}>
            编辑
          </Button>
          <Button
            status="danger"
            icon={<IconDelete />}
            onClick={() => setDeleteModalVisible(true)}
          >
            删除
          </Button>
        </Space>
      </div>

      {/* Type, Priority, Status Tags */}
      <div className="mb-4 flex gap-2">
        <StatusTag status={testCase.caseType} category="case_type" />
        <StatusTag status={testCase.priority} category="priority" />
        <StatusTag status={testCase.status} category="case_status" />
      </div>

      {/* Main Information Card */}
      <Card className="mb-4" title="用例信息">
        {/* Preconditions */}
        <div className="mb-4">
          <h3 className="font-medium mb-2">前置条件</h3>
          {testCase.preconditions && testCase.preconditions.length > 0 ? (
            <ul className="list-disc list-inside text-gray-700">
              {testCase.preconditions.map((condition, index) => (
                <li key={index}>{condition}</li>
              ))}
            </ul>
          ) : (
            <p className="text-gray-400">无前置条件</p>
          )}
        </div>

        {/* Steps */}
        <div className="mb-4">
          <h3 className="font-medium mb-2">测试步骤</h3>
          {testCase.steps && testCase.steps.length > 0 ? (
            <ol className="list-decimal list-inside text-gray-700">
              {testCase.steps.map((step, index) => (
                <li key={index} className="mb-1">
                  {step}
                </li>
              ))}
            </ol>
          ) : (
            <p className="text-gray-400">无测试步骤</p>
          )}
        </div>

        {/* Expected Results */}
        <div>
          <h3 className="font-medium mb-2">预期结果</h3>
          <pre className="bg-gray-50 p-3 rounded text-sm overflow-x-auto">
            {JSON.stringify(testCase.expected, null, 2)}
          </pre>
        </div>
      </Card>

      {/* AI Metadata Card */}
      {testCase.aiMetadata && (
        <Card className="mb-4" title="AI 来源">
          <div className="mb-3">
            <StatusTag
              status={testCase.aiMetadata.confidence}
              category="confidence"
            />
          </div>

          {testCase.aiMetadata.referencedChunks &&
          testCase.aiMetadata.referencedChunks.length > 0 ? (
            <ReferencePanel chunks={testCase.aiMetadata.referencedChunks} />
          ) : (
            <p className="text-gray-400">无引用来源</p>
          )}

          <div className="mt-3 text-sm text-gray-500">
            <p>模型版本: {testCase.aiMetadata.modelVersion}</p>
            <p>
              生成时间:{' '}
              {new Date(testCase.aiMetadata.generatedAt).toLocaleString(
                'zh-CN'
              )}
            </p>
          </div>
        </Card>
      )}

      {/* Metadata Card */}
      <Card title="元数据">
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-500">创建时间:</span>{' '}
            {new Date(testCase.createdAt).toLocaleString('zh-CN')}
          </div>
          <div>
            <span className="text-gray-500">更新时间:</span>{' '}
            {new Date(testCase.updatedAt).toLocaleString('zh-CN')}
          </div>
        </div>
      </Card>

      {/* Delete Confirmation Modal */}
      {deleteModalVisible && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <Card className="w-96">
            <h3 className="text-lg font-medium mb-4">确认删除</h3>
            <p className="mb-6">确定要删除此用例吗？此操作不可恢复。</p>
            <div className="flex justify-end gap-2">
              <Button onClick={() => setDeleteModalVisible(false)}>取消</Button>
              <Button
                status="danger"
                onClick={handleDelete}
                loading={deleteMutation.isPending}
              >
                确认删除
              </Button>
            </div>
          </Card>
        </div>
      )}
    </div>
  )
}
