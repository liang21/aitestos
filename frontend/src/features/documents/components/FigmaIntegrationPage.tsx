/**
 * Figma Integration Page (Placeholder)
 *
 * TODO: Implement Figma integration functionality
 * Phase 4 will complete this component with:
 * - Connection configuration (Personal Access Token / OAuth)
 * - File import from Figma URL
 * - Node selection tree
 *
 * @see plan.md §5.4
 */

import { Card, Empty } from '@arco-design/web-react'
import { Figma } from 'lucide-react'

export default function FigmaIntegrationPage() {
  return (
    <div className="p-6">
      <Card>
        <Empty
          icon={<Figma size={64} />}
          title="Figma 集成功能开发中"
          description="此功能将在后续版本中完成实现。"
        />
      </Card>
    </div>
  )
}
