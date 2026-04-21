/**
 * Draft Confirm Page
 * Split panel layout for reviewing and confirming AI-generated drafts
 */

import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Form,
  Input,
  Select,
  Button,
  Space,
  Message,
  Modal,
  Spin,
} from '@arco-design/web-react'
import { IconCheck, IconClose, IconSave } from '@arco-design/web-react/icon'
import { useDraftDetail, useConfirmDraft, useRejectDraft } from '../hooks/useDrafts'
import { SplitPanel } from '@/components/business/SplitPanel'
import { ArrayEditor } from '@/components/business/ArrayEditor'
import { ReferencePanel } from '@/components/business/ReferencePanel'
import { StatusTag } from '@/components/business/StatusTag'
import type { ReferencedChunk } from '@/types/api'

const reasonOptions = [
  { label: '重复', value: 'duplicate' },
  { label: '无关', value: 'irrelevant' },
  { label: '低质量', value: 'low_quality' },
  { label: '其他', value: 'other' },
]

export function DraftConfirmPage() {
  const { draftId } = useParams<{ draftId: string }>()
  const navigate = useNavigate()
  const [form] = Form.useForm()

  const [rejectModal, setRejectModal] = useState(false)
  const [selectedReason, setSelectedReason] = useState('')
  const [feedback, setFeedback] = useState('')

  // Fetch draft detail
  const { data: draft, isLoading } = useDraftDetail(draftId ?? '')

  // Mutations
  const confirmDraft = useConfirmDraft()
  const rejectDraft = useRejectDraft()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Spin size="large" />
      </div>
    )
  }

  if (!draft) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p>草稿不存在</p>
      </div>
    )
  }

  // Get AI metadata for reference panel
  const referencedChunks: ReferencedChunk[] =
    (draft as any).aiMetadata?.referencedChunks ?? []

  // Handle confirm
  const handleConfirm = async () => {
    try {
      const values = await form.validate()

      const result = await confirmDraft.mutateAsync({
        draftId: draft.id,
        data: {
          moduleId: 'default-module', // TODO: user should select target module
          ...values,
        },
      })

      Message.success(`草稿已确认，用例编号：${result.number}`)
      navigate(-1)
    } catch (error) {
      // Validation error or API error
      if (error instanceof Error) {
        Message.error(`确认失败：${error.message}`)
      } else {
        Message.error('确认草稿失败，请稍后重试')
      }
      console.error('Confirm failed:', error)
    }
  }

  // Handle save (temporary save, not confirm)
  const handleSave = () => {
    Message.info('草稿已保存到本地')
    // TODO: implement local storage save
  }

  // Handle reject
  const handleReject = () => {
    setRejectModal(true)
  }

  const confirmReject = async () => {
    try {
      await rejectDraft.mutateAsync({
        draftId: draft.id,
        data: {
          reason: selectedReason as any,
          feedback: feedback || undefined,
        },
      })

      Message.success('草稿已拒绝')
      setRejectModal(false)
      navigate(-1)
    } catch (error) {
      Message.error('拒绝草稿失败，请稍后重试')
      console.error('Reject failed:', error)
    }
  }

  return (
    <div className="h-screen flex flex-col">
      {/* Header */}
      <div className="flex items-center justify-between px-6 py-4 border-b">
        <div>
          <h1 className="text-xl font-semibold">草稿确认</h1>
          <p className="text-gray-500 text-sm mt-1">
            来源：{draft.projectName} / {draft.moduleName}
          </p>
        </div>
        <Space>
          <Button onClick={handleSave} icon={<IconSave />}>
            保存
          </Button>
          <Button
            status="danger"
            onClick={handleReject}
            icon={<IconClose />}
          >
            拒绝
          </Button>
          <Button
            type="primary"
            onClick={handleConfirm}
            loading={confirmDraft.isPending}
            icon={<IconCheck />}
          >
            确认并转为正式用例
          </Button>
        </Space>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-hidden">
        <SplitPanel defaultSplit={0.6} minLeftWidth={400} minRightWidth={300}>
          <SplitPanel.Left>
            <div className="p-6 h-full overflow-auto">
              <Form
                form={form}
                layout="vertical"
                initialValues={{
                  title: draft.title,
                  preconditions: draft.preconditions,
                  steps: draft.steps,
                  expected: draft.expected,
                  caseType: draft.caseType,
                  priority: draft.priority,
                }}
              >
                {/* Title */}
                <Form.Item
                  field="title"
                  label="用例标题"
                  rules={[{ required: true, message: '请输入用例标题' }]}
                >
                  <Input placeholder="请输入用例标题" />
                </Form.Item>

                {/* Preconditions */}
                <Form.Item field="preconditions" label="前置条件">
                  <ArrayEditor
                    placeholder="请输入前置条件"
                    minRows={0}
                  />
                </Form.Item>

                {/* Steps */}
                <Form.Item
                  field="steps"
                  label="测试步骤"
                  rules={[
                    {
                      validator: (value: string[]) =>
                        value && value.length >= 1
                          ? undefined
                          : '请至少输入一个测试步骤',
                    },
                  ]}
                >
                  <ArrayEditor
                    placeholder="请输入测试步骤"
                    minRows={1}
                  />
                </Form.Item>

                {/* Expected Results */}
                <Form.Item
                  field="expected"
                  label="预期结果"
                  rules={[{ required: true, message: '请输入预期结果' }]}
                >
                  <Input.TextArea
                    placeholder="请输入预期结果（JSON 格式，例如：{&quot;step_1&quot;: &quot;操作成功&quot;}）"
                    rows={4}
                  />
                </Form.Item>

                {/* Case Type & Priority */}
                <Space size="large">
                  <Form.Item
                    field="caseType"
                    label="用例类型"
                    rules={[{ required: true }]}
                  >
                    <Select>
                      <Select.Option value="functionality">功能测试</Select.Option>
                      <Select.Option value="performance">性能测试</Select.Option>
                      <Select.Option value="api">接口测试</Select.Option>
                      <Select.Option value="ui">UI 测试</Select.Option>
                      <Select.Option value="security">安全测试</Select.Option>
                    </Select>
                  </Form.Item>

                  <Form.Item
                    field="priority"
                    label="优先级"
                    rules={[{ required: true }]}
                  >
                    <Select>
                      <Select.Option value="P0">P0 紧急</Select.Option>
                      <Select.Option value="P1">P1 高</Select.Option>
                      <Select.Option value="P2">P2 中</Select.Option>
                      <Select.Option value="P3">P3 低</Select.Option>
                    </Select>
                  </Form.Item>
                </Space>

                {/* AI Metadata Display */}
                {(draft as any).aiMetadata && (
                  <div className="mt-4 pt-4 border-t">
                    <p className="text-sm text-gray-500 mb-2">AI 生成信息</p>
                    <Space>
                      <StatusTag
                        status={(draft as any).aiMetadata?.confidence}
                        category="confidence"
                      />
                      <span className="text-sm text-gray-400">
                        模型：{(draft as any).aiMetadata?.modelVersion}
                      </span>
                    </Space>
                  </div>
                )}
              </Form>
            </div>
          </SplitPanel.Left>

          <SplitPanel.Right>
            <div className="p-6 h-full overflow-auto">
              <ReferencePanel chunks={referencedChunks} />
            </div>
          </SplitPanel.Right>
        </SplitPanel>
      </div>

      {/* Reject Modal */}
      <Modal
        title="拒绝草稿"
        visible={rejectModal}
        onCancel={() => setRejectModal(false)}
        onOk={confirmReject}
        okText="确认拒绝"
        cancelText="取消"
        okButtonProps={{
          status: 'danger',
          disabled: !selectedReason,
        }}
      >
        <p className="mb-4">请选择拒绝原因：</p>
        <Select
          value={selectedReason}
          onChange={setSelectedReason}
          options={reasonOptions}
          placeholder="请选择拒绝原因"
          className="w-full mb-4"
        />
        <p className="mb-2">反馈意见（可选）：</p>
        <Input.TextArea
          value={feedback}
          onChange={setFeedback}
          placeholder="请输入具体的反馈意见，帮助改进生成质量..."
          rows={4}
        />
      </Modal>
    </div>
  )
}
