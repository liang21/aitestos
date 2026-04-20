import { Button, Card, Typography } from '@arco-design/web-react'
import { IconPlus } from '@arco-design/web-react/icon'

const { Title, Paragraph } = Typography

/**
 * Project List Page - Placeholder
 * TODO: Implement in Phase 3
 */
export function ProjectListPage() {
  return (
    <div className="p-6">
      <Title heading={4}>项目管理</Title>
      <Paragraph className="text-gray-500 mb-6">
        管理您的测试项目、模块和配置
      </Paragraph>

      <Card className="max-w-md mx-auto text-center py-12">
        <div className="text-gray-400 mb-4">
          <IconPlus className="w-12 h-12 mx-auto mb-2" />
          <p>项目管理模块即将推出</p>
          <p className="text-sm">Phase 3: 项目管理模块</p>
        </div>
        <Button type="primary" disabled>
          新建项目
        </Button>
      </Card>
    </div>
  )
}
