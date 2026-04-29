/**
 * New Plan Page
 * Form for creating new test plans with case selection
 */

import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import {
  Button,
  Card,
  Form,
  Input,
  Checkbox,
  Message,
  Space,
} from '@arco-design/web-react'
import { useCreatePlan } from '../hooks/usePlans'
import { useCaseList } from '@/features/testcases/hooks/useTestCases'
import type { CreatePlanRequest } from '@/types/api'

const { TextArea } = Input

export function NewPlanPage() {
  const navigate = useNavigate()
  const { projectId } = useParams<{ projectId: string }>()
  const [form] = Form.useForm()
  const createMutation = useCreatePlan()

  // Get available test cases
  const { data: casesData, isLoading: casesLoading } = useCaseList(
    projectId || '',
    { offset: 0, limit: 100 }
  )

  const [selectedCaseIds, setSelectedCaseIds] = useState<string[]>([])

  // Handle form submission
  const handleSubmit = async (values: {
    name: string
    description: string
  }) => {
    try {
      // Validate name
      if (!values.name || values.name.length < 3) {
        Message.error('请输入计划名称，至少 3 个字符')
        return
      }

      // Validate at least one case is selected
      if (selectedCaseIds.length === 0) {
        Message.error('请至少选择一个用例')
        return
      }

      // Create plan
      const planData: CreatePlanRequest = {
        projectId: effectiveProjectId,
        name: values.name,
        description: values.description,
      }

      const plan = await createMutation.mutateAsync(planData)

      // Add cases to plan (this would be done in the backend or via a separate API call)
      // For now, we'll assume the plan is created and navigate to detail page

      Message.success('计划创建成功')
      navigate(`/plans/${plan.id}`)
    } catch (error) {
      Message.error(
        `创建失败：${error instanceof Error ? error.message : '未知错误'}`
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

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <div className="mb-4">
        <Button onClick={() => navigate('/plans')} type="text">
          ← 返回
        </Button>
      </div>

      <h1 className="text-2xl font-semibold mb-6">新建测试计划</h1>

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
          <Button onClick={() => navigate('/plans')}>取消</Button>
          <Button
            type="primary"
            htmlType="submit"
            loading={createMutation.isPending}
          >
            创建
          </Button>
        </div>
      </Form>
    </div>
  )
}
