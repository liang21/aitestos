/**
 * Edit Plan Page
 * Form for editing existing test plans with case selection
 */

import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import {
  Button,
  Card,
  Form,
  Input,
  Checkbox,
  Space,
} from '@arco-design/web-react'
import { usePlanDetail, useUpdatePlan } from '../hooks/usePlans'
import { useCaseList } from '@/features/testcases/hooks/useTestCases'
import { buildProjectRoutes } from '@/lib/routes'
import type { UpdatePlanRequest } from '@/types/api'
import { messageSuccess, messageError } from '@/lib/notification'

const { TextArea } = Input

export function EditPlanPage() {
  const navigate = useNavigate()
  const { projectId, planId } = useParams<{ projectId: string; planId: string }>()
  const [form] = Form.useForm()
  const updateMutation = useUpdatePlan()

  // Get project-scoped routes
  const routes = projectId ? buildProjectRoutes(projectId) : null

  // Load existing plan data
  const { data: planDetail, isLoading: planLoading } = usePlanDetail(planId ?? '')

  // Get available test cases
  const { data: casesData, isLoading: casesLoading } = useCaseList(
    projectId ?? '',
    { offset: 0, limit: 100 }
  )

  const [selectedCaseIds, setSelectedCaseIds] = useState<string[]>([])

  // Load existing plan data into form when available
  useEffect(() => {
    if (planDetail && form) {
      form.setFieldsValue({
        name: planDetail.name,
        description: planDetail.description || '',
      })
      // Load selected cases from plan
      const caseIds = planDetail.cases.map((c) => c.caseId)
      setSelectedCaseIds(caseIds)
    }
  }, [planDetail, form])

  // Handle form submission
  const handleSubmit = async (values: {
    name: string
    description: string
  }) => {
    if (!planId) return

    try {
      // Validate name
      if (!values.name || values.name.length < 3) {
        messageError('请输入计划名称，至少 3 个字符')
        return
      }

      // Validate at least one case is selected
      if (selectedCaseIds.length === 0) {
        messageError('请至少选择一个用例')
        return
      }

      // Update plan
      const planData: UpdatePlanRequest = {
        name: values.name,
        description: values.description,
      }

      await updateMutation.mutateAsync({
        id: planId,
        data: planData,
      })

      messageSuccess('计划更新成功')
      navigate(routes?.plans.detail(planId) ?? `/plans/${planId}`)
    } catch (error) {
      messageError(
        `更新失败：${error instanceof Error ? error.message : '未知错误'}`
      )
    }
  }

  // Handle case selection
  const handleCaseSelection = (caseId: string, checked: boolean) => {
    if (checked) {
      setSelectedCaseIds([...selectedCaseIds, caseId])
    } else {
      setSelectedCaseIds(selectedCaseIds.filter((id) => id !== caseId))
    }
  }

  // Handle select all
  const handleSelectAll = (checked: boolean) => {
    if (checked && casesData?.data) {
      setSelectedCaseIds(casesData.data.map((c) => c.id))
    } else {
      setSelectedCaseIds([])
    }
  }

  if (planLoading) {
    return (
      <div className="p-6 max-w-4xl mx-auto">
        <div className="text-center py-8">加载中...</div>
      </div>
    )
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <div className="mb-4">
        <Button onClick={() => navigate(routes?.plans.detail(planId ?? '') ?? `/plans/${planId}`)} type="text">
          ← 返回
        </Button>
      </div>

      <h1 className="text-2xl font-semibold mb-6">编辑测试计划</h1>

      <Form
        form={form}
        layout="vertical"
        onSubmit={handleSubmit}
        initialValues={{ name: '', description: '' }}
      >
        <Card className="mb-6">
          <Form.Item
            field="name"
            label="计划名称"
            rules={[
              { required: true, message: '请输入计划名称' },
              { minLength: 3, message: '计划名称至少 3 个字符' },
            ]}
          >
            <Input placeholder="请输入计划名称" />
          </Form.Item>

          <Form.Item field="description" label="计划描述">
            <TextArea placeholder="请输入计划描述" rows={3} />
          </Form.Item>
        </Card>

        <Card
          title={`选择用例 ${selectedCaseIds.length > 0 ? `(已选择 ${selectedCaseIds.length} 个)` : ''}`}
          className="mb-6"
          extra={
            <Checkbox
              checked={
                casesData?.data &&
                selectedCaseIds.length === casesData.data.length &&
                selectedCaseIds.length > 0
              }
              onChange={handleSelectAll}
            >
              全选
            </Checkbox>
          }
        >
          {casesLoading ? (
            <div className="text-center py-8">加载中...</div>
          ) : (
            <div className="space-y-2">
              {casesData?.data?.map((testCase) => (
                <div
                  key={testCase.id}
                  className="flex items-center p-3 border rounded hover:bg-gray-50"
                >
                  <Checkbox
                    checked={selectedCaseIds.includes(testCase.id)}
                    onChange={(checked) =>
                      handleCaseSelection(testCase.id, checked)
                    }
                  >
                    <span className="ml-2 font-medium">{testCase.number}</span>
                    <span className="ml-2 text-gray-600">{testCase.title}</span>
                  </Checkbox>
                </div>
              ))}
            </div>
          )}
        </Card>

        <div className="flex justify-end gap-2">
          <Button onClick={() => navigate(routes?.plans.detail(planId ?? '') ?? `/plans/${planId}`)}>取消</Button>
          <Button
            type="primary"
            htmlType="submit"
            loading={updateMutation.isPending}
          >
            保存
          </Button>
        </div>
      </Form>
    </div>
  )
}
