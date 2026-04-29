/**
 * Figma Integration Page
 *
 * Allows importing Figma designs as knowledge base documents.
 * Requires admin role.
 *
 * @see plan.md §5.4
 * @see US-2.2 关联 Figma 设计稿
 */

import { useState } from 'react'
import {
  Card,
  Form,
  Input,
  Radio,
  Button,
  Tree,
  Space,
  Message,
  Steps,
} from '@arco-design/web-react'
import { Figma, Link as LinkIcon, Loader2 } from 'lucide-react'

const { TextArea } = Input

interface FigmaNode {
  id: string
  name: string
  type: string
  children?: FigmaNode[]
}

export default function FigmaIntegrationPage() {
  const [form] = Form.useForm()
  const [step, setStep] = useState(1)
  const [loading, setLoading] = useState(false)
  const [figmaNodes, setFigmaNodes] = useState<FigmaNode[]>([])
  const [selectedKeys, setSelectedKeys] = useState<string[]>([])

  // Handle connection test
  const handleTestConnection = async () => {
    setLoading(true)
    try {
      // TODO: Implement actual Figma API connection test
      await new Promise((resolve) => setTimeout(resolve, 1000))
      Message.success('连接成功')
      setStep(2)
    } catch {
      Message.error('连接失败，请检查令牌是否正确')
    } finally {
      setLoading(false)
    }
  }

  // Handle file parsing
  const handleParseFile = async () => {
    const url = form.getFieldValue('fileUrl')
    if (!url) {
      Message.error('请输入 Figma 文件 URL')
      return
    }

    setLoading(true)
    try {
      // TODO: Implement actual Figma API file parsing
      await new Promise((resolve) => setTimeout(resolve, 1500))

      // Mock nodes data
      const mockNodes: FigmaNode[] = [
        {
          id: '1',
          name: '登录页面',
          type: 'FRAME',
          children: [
            { id: '1-1', name: '用户名输入框', type: 'COMPONENT' },
            { id: '1-2', name: '密码输入框', type: 'COMPONENT' },
            { id: '1-3', name: '登录按钮', type: 'COMPONENT' },
          ],
        },
        {
          id: '2',
          name: '主页',
          type: 'FRAME',
          children: [
            { id: '2-1', name: '导航栏', type: 'COMPONENT' },
            { id: '2-2', name: '内容区', type: 'COMPONENT' },
          ],
        },
      ]

      setFigmaNodes(mockNodes)
      Message.success('解析成功')
      setStep(3)
    } catch {
      Message.error('解析失败，请检查 URL 是否正确')
    } finally {
      setLoading(false)
    }
  }

  // Handle import confirmation
  const handleConfirmImport = async () => {
    if (selectedKeys.length === 0) {
      Message.warning('请至少选择一个节点')
      return
    }

    setLoading(true)
    try {
      // TODO: Implement actual Figma node import
      await new Promise((resolve) => setTimeout(resolve, 1500))
      Message.success(`成功导入 ${selectedKeys.length} 个节点`)
      setStep(1)
      form.resetFields()
      setFigmaNodes([])
      setSelectedKeys([])
    } catch {
      Message.error('导入失败')
    } finally {
      setLoading(false)
    }
  }

  // Build tree data structure
  const buildTreeData = (nodes: FigmaNode[]) => {
    return nodes.map((node) => ({
      key: node.id,
      title: node.name,
      children: node.children ? buildTreeData(node.children) : undefined,
    }))
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <Card>
        <div className="mb-6">
          <h2 className="text-xl font-semibold mb-2">Figma 设计稿集成</h2>
          <p className="text-gray-500">
            从 Figma 导入设计稿作为知识库文档，AI 将基于设计生成测试用例
          </p>
        </div>

        <Steps current={step - 1} className="mb-8">
          <Steps.Item title="连接配置" />
          <Steps.Item title="导入文件" />
          <Steps.Item title="选择节点" />
        </Steps>

        {/* Step 1: Connection Configuration */}
        {step === 1 && (
          <Form form={form} layout="vertical">
            <Form.Item
              label="认证方式"
              field="authType"
              initialValue="token"
              rules={[{ required: true, message: '请选择认证方式' }]}
            >
              <Radio.Group>
                <Radio value="token">个人访问令牌</Radio>
                <Radio value="oauth">OAuth 2.0</Radio>
              </Radio.Group>
            </Form.Item>

            <Form.Item
              label="访问令牌"
              field="token"
              rules={[{ required: true, message: '请输入 Figma 访问令牌' }]}
            >
              <Input.Password
                placeholder="figma_..."
                style={{ width: '100%' }}
              />
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                icon={<LinkIcon size={16} />}
                loading={loading}
                onClick={handleTestConnection}
              >
                测试连接
              </Button>
            </Form.Item>
          </Form>
        )}

        {/* Step 2: File Import */}
        {step === 2 && (
          <Form form={form} layout="vertical">
            <Form.Item
              label="Figma 文件 URL"
              field="fileUrl"
              rules={[
                { required: true, message: '请输入 Figma 文件 URL' },
                {
                  pattern: /^https:\/\/www\.figma\.com\/file\/.+/,
                  message: '请输入有效的 Figma 文件 URL',
                },
              ]}
            >
              <Input
                placeholder="https://www.figma.com/file/..."
                prefix={<Figma size={16} />}
              />
            </Form.Item>

            <Form.Item>
              <Space>
                <Button
                  type="primary"
                  loading={loading}
                  onClick={handleParseFile}
                >
                  解析文件
                </Button>
                <Button onClick={() => setStep(1)}>上一步</Button>
              </Space>
            </Form.Item>
          </Form>
        )}

        {/* Step 3: Node Selection */}
        {step === 3 && (
          <>
            <div className="mb-4">
              <p className="text-gray-500 mb-4">选择要导入的页面和组件</p>
              {figmaNodes.length > 0 ? (
                <Tree
                  checkable
                  treeData={buildTreeData(figmaNodes)}
                  defaultExpandAll
                  checkedKeys={selectedKeys}
                  onCheck={(keys) => setSelectedKeys(keys as string[])}
                  style={{ maxHeight: 400, overflow: 'auto' }}
                />
              ) : (
                <div className="text-center text-gray-400 py-8">
                  暂无节点数据
                </div>
              )}
            </div>

            <div className="flex gap-3">
              <Button
                type="primary"
                loading={loading}
                onClick={handleConfirmImport}
              >
                确认导入
              </Button>
              <Button onClick={() => setStep(2)}>上一步</Button>
              <Button
                onClick={() => {
                  setStep(1)
                  form.resetFields()
                  setFigmaNodes([])
                  setSelectedKeys([])
                }}
              >
                取消
              </Button>
            </div>
          </>
        )}
      </Card>
    </div>
  )
}
