/**
 * Draft Confirm Page
 * Split panel layout for reviewing and confirming AI-generated drafts
 * Features: Navigation between drafts, confirmation with module selection, dirty state protection
 */

import { useState, useEffect, useCallback, useMemo } from 'react'
import { useParams, useNavigate, useBlocker } from 'react-router-dom'
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
import { IconCheck, IconClose, IconSave, IconLeft, IconRight } from '@arco-design/web-react/icon'
import {
  useDraftDetail,
  useConfirmDraft,
  useRejectDraft,
  useDraftList,
} from '../hooks/useDrafts'
import { useModuleList } from '@/features/modules/hooks/useModules'
import { SplitPanel } from '@/components/business/SplitPanel'
import { ArrayEditor } from '@/components/business/ArrayEditor'
import { ReferencePanel } from '@/components/business/ReferencePanel'
import { StatusTag } from '@/components/business/StatusTag'
import { buildProjectRoutes } from '@/lib/routes'
import type { ReferencedChunk } from '@/types/api'

const reasonOptions = [
  { label: '重复', value: 'duplicate' },
  { label: '无关', value: 'irrelevant' },
  { label: '低质量', value: 'low_quality' },
  { label: '其他', value: 'other' },
]

interface DraftNavigationState {
  draftIds: string[]
  currentIndex: number
}

export function DraftConfirmPage() {
  const { draftId } = useParams<{ draftId: string }>()
  const navigate = useNavigate()
  const [form] = Form.useForm()

  // State for modals and navigation
  const [rejectModal, setRejectModal] = useState(false)
  const [confirmModal, setConfirmModal] = useState(false)
  const [selectedReason, setSelectedReason] = useState('')
  const [feedback, setFeedback] = useState('')
  const [selectedModuleId, setSelectedModuleId] = useState('')
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false)

  // Draft navigation state
  const [navState, setNavState] = useState<DraftNavigationState>({
    draftIds: [],
    currentIndex: 0,
  })

  // Local edits state (for auto-save between drafts)
  const [localEdits, setLocalEdits] = useState<Record<string, any>>({})

  // Fetch draft list for navigation
  const { data: draftList } = useDraftList({ status: 'pending' })

  // Fetch current draft detail
  const { data: draft, isLoading } = useDraftDetail(draftId ?? '')

  // Fetch modules for current project
  const { data: modules } = useModuleList(draft?.projectId ?? '', {
    enabled: !!draft?.projectId,
  })

  // Mutations
  const confirmDraft = useConfirmDraft()
  const rejectDraft = useRejectDraft()

  // Initialize navigation state when draft list loads
  useEffect(() => {
    if (draftList?.data) {
      const ids = draftList.data.map(d => d.id)
      setNavState({
        draftIds: ids,
        currentIndex: ids.indexOf(draftId ?? ''),
      })
    }
  }, [draftList, draftId])

  // Navigation blocker for unsaved changes
  const shouldBlock = useMemo(
    () =>
      hasUnsavedChanges &&
      document.activeElement?.tagName !== 'INPUT' &&
      document.activeElement?.tagName !== 'TEXTAREA',
    [hasUnsavedChanges]
  )
  const blocker = useBlocker(shouldBlock)

  // Auto-save local edits before switching drafts
  const saveLocalEdits = useCallback(() => {
    if (draftId && hasUnsavedChanges) {
      const values = form.getFieldsValue()
      setLocalEdits(prev => ({
        ...prev,
        [draftId]: values,
      }))
    }
  }, [draftId, hasUnsavedChanges, form])

  // Handle draft navigation
  const navigateToDraft = useCallback((direction: 'prev' | 'next') => {
    if (!draftId) return

    // Save current edits before navigating
    saveLocalEdits()

    const newIndex = direction === 'next'
      ? navState.currentIndex + 1
      : navState.currentIndex - 1

    if (newIndex >= 0 && newIndex < navState.draftIds.length) {
      setHasUnsavedChanges(false)
      navigate(`/drafts/${navState.draftIds[newIndex]}`)
    }
  }, [draftId, navState, saveLocalEdits, navigate])

  // Handle keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Only handle arrow keys when not in input/textarea
      const target = e.target as HTMLElement
      if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') return

      if (e.key === 'ArrowLeft') {
        e.preventDefault()
        navigateToDraft('prev')
      } else if (e.key === 'ArrowRight') {
        e.preventDefault()
        navigateToDraft('next')
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [navigateToDraft])

  // Restore local edits when draft changes
  useEffect(() => {
    if (draftId && localEdits[draftId]) {
      form.setFieldsValue(localEdits[draftId])
      setHasUnsavedChanges(false)
    }
  }, [draftId, localEdits, form])

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

  // Handle confirm button click - show module selection modal
  const handleConfirmClick = () => {
    // Set default module if available
    if (modules && modules.length > 0) {
      setSelectedModuleId(modules[0].id)
    }
    setConfirmModal(true)
  }

  // Handle confirm submission
  const handleConfirmSubmit = async () => {
    if (!selectedModuleId) {
      Message.warning('请选择目标模块')
      return
    }

    try {
      const values = await form.validate()

      const result = await confirmDraft.mutateAsync({
        draftId: draft.id,
        data: {
          moduleId: selectedModuleId,
          ...values,
        },
      })

      Message.success(`用例 ${result.number} 已创建`)
      setHasUnsavedChanges(false)
      setConfirmModal(false)

      // Navigate to case detail page - use project-scoped route
      const routes = draft.projectId ? buildProjectRoutes(draft.projectId) : null
      navigate(routes?.cases.detail(result.id) ?? `/testcases/${result.id}`)
    } catch (error) {
      if (error instanceof Error) {
        Message.error(`确认失败：${error.message}`)
      } else {
        Message.error('确认草稿失败，请稍后重试')
      }
      console.error('Confirm failed:', error)
    }
  }

  // Handle save (temporary save to local state)
  const handleSave = () => {
    const values = form.getFieldsValue()
    setLocalEdits(prev => ({
      ...prev,
      [draftId]: values,
    }))
    setHasUnsavedChanges(false)
    Message.info('草稿已保存到本地')
  }

  // Handle reject
  const handleReject = () => {
    setRejectModal(true)
  }

  const handleRejectSubmit = async () => {
    try {
      await rejectDraft.mutateAsync({
        draftId: draft.id,
        data: {
          reason: selectedReason as any,
          feedback: feedback || undefined,
        },
      })

      Message.success('草稿已拒绝')
      setHasUnsavedChanges(false)
      setRejectModal(false)
      navigate(-1)
    } catch (error) {
      Message.error('拒绝草稿失败，请稍后重试')
      console.error('Reject failed:', error)
    }
  }

  // Handle form value changes
  const handleFormChange = () => {
    setHasUnsavedChanges(true)
  }

  const canGoPrev = navState.currentIndex > 0
  const canGoNext = navState.currentIndex < navState.draftIds.length - 1
  const currentNumber = navState.currentIndex + 1
  const totalDrafts = navState.draftIds.length

  return (
    <div className="h-screen flex flex-col">
      {/* Header with Navigation */}
      <div className="flex items-center justify-between px-6 py-4 border-b">
        <div className="flex items-center gap-4">
          <Button
            icon={<IconLeft />}
            onClick={() => navigate(-1)}
            type="text"
          >
            返回
          </Button>
          <div>
            <h1 className="text-xl font-semibold">草稿确认</h1>
            <div className="flex items-center gap-3 text-sm text-gray-500 mt-1">
              <span>来源：{draft.projectName} / {draft.moduleName}</span>
              {totalDrafts > 0 && (
                <>
                  <span>•</span>
                  <span className="font-medium">
                    第 {currentNumber} / {totalDrafts} 条
                  </span>
                </>
              )}
            </div>
          </div>
        </div>

        {/* Navigation Controls */}
        {totalDrafts > 1 && (
          <div className="flex items-center gap-2">
            {navState.draftIds.map((id, index) => (
              <button
                key={id}
                onClick={() => {
                  saveLocalEdits()
                  setHasUnsavedChanges(false)
                  navigate(`/drafts/${id}`)
                }}
                className={`w-2 h-2 rounded-full transition-colors ${
                  index === navState.currentIndex
                    ? 'bg-primary-600 scale-125'
                    : 'bg-gray-300 hover:bg-gray-400'
                }`}
                title={`跳转到第 ${index + 1} 条`}
              />
            ))}
          </div>
        )}

        <Space>
          <Button
            onClick={handleSave}
            icon={<IconSave />}
            disabled={!hasUnsavedChanges}
          >
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
            onClick={handleConfirmClick}
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
                onValuesChange={handleFormChange}
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
                  <ArrayEditor placeholder="请输入前置条件" minRows={0} />
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
                  <ArrayEditor placeholder="请输入测试步骤" minRows={1} />
                </Form.Item>

                {/* Expected Results */}
                <Form.Item
                  field="expected"
                  label="预期结果"
                  rules={[{ required: true, message: '请输入预期结果' }]}
                >
                  <Input.TextArea
                    placeholder='请输入预期结果（JSON 格式，例如：{"step_1": "操作成功"}）'
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
                      <Select.Option value="functionality">
                        功能测试
                      </Select.Option>
                      <Select.Option value="performance">
                        性能测试
                      </Select.Option>
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
            <div className="p-6 h-full overflow-auto bg-gradient-to-br from-purple-50 to-white">
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
        onOk={handleRejectSubmit}
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

      {/* Confirm Modal with Module Selection */}
      <Modal
        title="确认并转为正式用例"
        visible={confirmModal}
        onCancel={() => setConfirmModal(false)}
        onOk={handleConfirmSubmit}
        okText="确认"
        cancelText="取消"
        okButtonProps={{
          loading: confirmDraft.isPending,
          disabled: !selectedModuleId,
        }}
        style={{ width: 500 }}
      >
        <div className="py-4">
          <p className="mb-4 text-gray-600">
            请选择要将此草稿确认到哪个模块：
          </p>

          {modules && modules.length > 0 ? (
            <Select
              value={selectedModuleId}
              onChange={setSelectedModuleId}
              placeholder="请选择目标模块"
              className="w-full"
              options={modules.map(m => ({
                label: `${m.name} (${m.abbreviation})`,
                value: m.id,
              }))}
            />
          ) : (
            <div className="text-center py-8 text-gray-500">
              <p>当前项目暂无模块</p>
              <p className="text-sm mt-2">请先在项目设置中创建模块</p>
            </div>
          )}

          {selectedModuleId && (
            <div className="mt-4 p-3 bg-blue-50 rounded text-sm text-blue-700">
              确认后，草稿将转为正式用例并生成用例编号
            </div>
          )}
        </div>
      </Modal>

      {/* Navigation Blocker Confirmation Dialog */}
      {blocker.state === 'blocked' && (
        <Modal
          title="有未保存的修改"
          visible
          onOk={() => {
            blocker.proceed?.()
            setHasUnsavedChanges(false)
          }}
          onCancel={() => blocker.reset?.()}
          okText="放弃修改"
          cancelText="取消"
          focusOk={false}
        >
          <p>您有未保存的修改，确定要离开吗？离开后将丢失这些修改。</p>
        </Modal>
      )}
    </div>
  )
}
