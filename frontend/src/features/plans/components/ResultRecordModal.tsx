/**
 * Result Record Modal
 * Modal for recording test execution results
 */

import { useState, useEffect } from 'react'
import {
  Modal,
  Radio,
  Input,
  Button,
  Space,
  Message,
} from '@arco-design/web-react'
import type { ResultStatus, PlanCase } from '@/types/api'
import { plansApi } from '../services/plans'

const resultStatusTextMap: Record<ResultStatus, string> = {
  pass: '通过',
  fail: '失败',
  block: '阻塞',
  skip: '跳过',
}

interface ResultRecordModalProps {
  visible: boolean
  planId: string
  caseId: string
  caseTitle: string
  existingResult?: PlanCase
  onClose: () => void
  onSubmit?: (data: { status: ResultStatus; note?: string }) => void
}

const { TextArea } = Input

export function ResultRecordModal({
  visible,
  planId,
  caseId,
  caseTitle,
  existingResult,
  onClose,
  onSubmit,
}: ResultRecordModalProps) {
  const [status, setStatus] = useState<ResultStatus>(
    existingResult?.resultStatus || 'pass'
  )
  const [note, setNote] = useState(existingResult?.resultNote || '')
  const [existingCase, setExistingCase] = useState<PlanCase | undefined>(
    undefined
  )

  // Fetch existing plan detail to check if case is already executed
  useEffect(() => {
    if (visible && planId && !existingResult) {
      plansApi
        .get(planId)
        .then((planDetail) => {
          const existing = planDetail.cases.find((c) => c.caseId === caseId)
          setExistingCase(existing)
          if (existing && existing.resultStatus) {
            setStatus(existing.resultStatus)
            setNote(existing.resultNote || '')
          }
        })
        .catch(() => {
          // Ignore error, will proceed without existing data
        })
    }
  }, [visible, planId, caseId, existingResult])

  // Handle modal close
  const handleClose = () => {
    setStatus('pass')
    setNote('')
    setExistingCase(undefined)
    onClose()
  }

  // Handle submit
  const handleSubmit = () => {
    if (!onSubmit) {
      handleClose()
      return
    }

    // Validate note length if provided
    if (note && note.length > 500) {
      Message.error('备注不能超过 500 字符')
      return
    }

    // Check if result already exists and warn user
    if (existingCase && existingCase.resultStatus) {
      Modal.confirm({
        title: '确认覆盖？',
        content: `该用例已有执行结果：${resultStatusTextMap[existingCase.resultStatus]}，确定要覆盖吗？`,
        onOk: () => {
          onSubmit({ status, note: note || undefined })
        },
      })
    } else {
      onSubmit({ status, note: note || undefined })
    }
  }

  // Render warning if result exists
  const hasExistingResult = existingCase && existingCase.resultStatus

  return (
    <Modal
      title="录入执行结果"
      visible={visible}
      onCancel={handleClose}
      footer={
        <Space>
          <Button onClick={handleClose}>取消</Button>
          <Button type="primary" onClick={handleSubmit}>
            提交
          </Button>
        </Space>
      }
      width={600}
    >
      <div className="mb-4">
        <div className="text-gray-500 text-sm mb-1">用例</div>
        <div>{caseTitle}</div>
        {hasExistingResult && (
          <div className="mt-2 p-2 bg-orange-50 border border-orange-200 rounded text-orange-700 text-sm">
            ⚠️ 该用例已有执行结果：
            {resultStatusTextMap[existingCase.resultStatus]}
          </div>
        )}
      </div>

      <div className="mb-4">
        <div className="mb-2 font-medium">执行状态</div>
        <Radio.Group
          value={status}
          onChange={(value) => setStatus(value as ResultStatus)}
          options={[
            { label: '通过', value: 'pass' },
            { label: '失败', value: 'fail' },
            { label: '阻塞', value: 'block' },
            { label: '跳过', value: 'skip' },
          ]}
        />
      </div>

      <div>
        <div className="mb-2 font-medium">
          备注{' '}
          <span className="text-gray-400 font-normal">{note.length}/500</span>
        </div>
        <TextArea
          placeholder="请输入备注（可选）"
          rows={4}
          value={note}
          onChange={setNote}
          maxLength={500}
          showWordLimit
        />
      </div>
    </Modal>
  )
}
