/**
 * Figma Node Preview Component
 * Displays Figma node previews with lazy-loaded images
 */

import { useState } from 'react'
import { Image } from 'lucide-react'
import { LazyImage } from '@/components/business/LazyImage'

interface FigmaNode {
  id: string
  name: string
  type: string
  thumbnailUrl?: string
  children?: FigmaNode[]
}

interface FigmaNodePreviewProps {
  nodeName: string
  nodeType?: string
  thumbnailUrl?: string
  children?: FigmaNodePreviewProps[]
  expanded?: boolean
  onToggle?: () => void
}

/**
 * Single Figma node with preview thumbnail
 */
export function FigmaNodePreview({
  nodeName,
  nodeType = 'COMPONENT',
  thumbnailUrl,
  children,
  expanded = false,
  onToggle,
}: FigmaNodePreviewProps) {
  const [imageError, setImageError] = useState(false)
  const hasChildren = children && children.length > 0

  const handleImageError = () => {
    setImageError(true)
  }

  const handleImageLoad = () => {
    setImageError(false)
  }

  return (
    <div className="figma-node mb-2">
      <div
        className={`flex items-center gap-2 p-2 rounded border hover:bg-gray-50 transition-colors cursor-pointer ${
          hasChildren ? 'cursor-pointer' : ''
        }`}
        onClick={hasChildren ? onToggle : undefined}
      >
        {/* Expand/Collapse Icon */}
        {hasChildren && (
          <span className="text-gray-400">
            {expanded ? '▼' : '▶'}
          </span>
        )}

        {/* Node Type Icon */}
        <span className="text-gray-400">
          {nodeType === 'FRAME' ? '📄' : nodeType === 'COMPONENT' ? '🧩' : '📝'}
        </span>

        {/* Node Name */}
        <span className="flex-1">{nodeName}</span>

        {/* Thumbnail Preview */}
        {thumbnailUrl && !imageError && (
          <div className="w-16 h-16 rounded border overflow-hidden">
            <LazyImage
              src={thumbnailUrl}
              alt={`${nodeName} 预览`}
              className="w-full h-full object-cover"
              onError={handleImageError}
              onLoad={handleImageLoad}
            />
          </div>
        )}

        {/* Error State */}
        {imageError && (
          <div className="w-16 h-16 rounded border bg-gray-100 flex items-center justify-center text-gray-400">
            <Image size={14} />
          </div>
        )}
      </div>

      {/* Child Nodes */}
      {hasChildren && expanded && (
        <div className="ml-6 pl-2 border-l border-gray-200">
          {children?.map((child, index) => (
            <FigmaNodePreview
              key={`${child.nodeName}-${index}`}
              {...child}
            />
          ))}
        </div>
      )}
    </div>
  )
}

/**
 * Figma Node Preview List
 * Displays a list of Figma nodes with lazy-loaded thumbnails
 */
interface FigmaNodePreviewListProps {
  nodes: Array<{
    id: string
    name: string
    type: string
    thumbnailUrl?: string
    children?: Array<{
      id: string
      name: string
      type: string
      thumbnailUrl?: string
      children?: unknown[]
    }>
  }>
  maxVisible?: number
}

export function FigmaNodePreviewList({
  nodes,
  maxVisible = 10,
}: FigmaNodePreviewListProps) {
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set())

  const toggleNode = (nodeId: string) => {
    setExpandedNodes((prev) => {
      const newSet = new Set(prev)
      if (newSet.has(nodeId)) {
        newSet.delete(nodeId)
      } else {
        newSet.add(nodeId)
      }
      return newSet
    })
  }

  // Limit visible nodes for performance
  const visibleNodes = nodes.slice(0, maxVisible)

  return (
    <div className="figma-node-preview-list">
      {visibleNodes.map((node) => (
        <FigmaNodePreview
          key={node.id}
          nodeName={node.name}
          nodeType={node.type}
          thumbnailUrl={node.thumbnailUrl}
          expanded={expandedNodes.has(node.id)}
          onToggle={() => toggleNode(node.id)}
          children={node.children?.map((child) => ({
            nodeName: child.name,
            nodeType: child.type,
            thumbnailUrl: child.thumbnailUrl,
            children: child.children as FigmaNodePreviewProps['children'],
          }))}
        />
      ))}

      {nodes.length > maxVisible && (
        <div className="text-center text-gray-500 text-sm mt-4">
          还有 {nodes.length - maxVisible} 个节点未显示
        </div>
      )}
    </div>
  )
}

/**
 * Compact Figma Node Tree for selection
 * Optimized for large number of nodes
 */
interface CompactFigmaTreeProps {
  nodes: FigmaNode[]
  onNodeSelect?: (nodeIds: string[]) => void
  maxNodes?: number
}

export function CompactFigmaTree({
  nodes,
  onNodeSelect,
  maxNodes = 100,
}: CompactFigmaTreeProps) {
  const [selectedNodes, setSelectedNodes] = useState<Set<string>>(new Set())

  const toggleNodeSelection = (nodeId: string) => {
    setSelectedNodes((prev) => {
      const newSet = new Set(prev)
      if (newSet.has(nodeId)) {
        newSet.delete(nodeId)
      } else {
        newSet.add(nodeId)
      }
      onNodeSelect?.(Array.from(newSet))
      return newSet
    })
  }

  // Memoized node renderer for performance
  const renderNode = (node: FigmaNode, level: number = 0): React.ReactNode => {
    const isSelected = selectedNodes.has(node.id)

    return (
      <div key={node.id} style={{ marginLeft: `${level * 16}px` }} className="mb-1">
        <label className="flex items-center gap-2 cursor-pointer hover:bg-gray-50 p-1 rounded">
          <input
            type="checkbox"
            checked={isSelected}
            onChange={() => toggleNodeSelection(node.id)}
            className="rounded"
          />
          <span className="text-sm">
            {node.type === 'FRAME' ? '📄' : '🧩'} {node.name}
          </span>
          {node.thumbnailUrl && (
            <LazyImage
              src={node.thumbnailUrl}
              alt={`${node.name} 预览`}
              className="w-8 h-8 rounded"
              style={{ display: 'inline-block' }}
            />
          )}
        </label>
        {node.children?.map((child) => renderNode(child, level + 1))}
      </div>
    )
  }

  // Limit total nodes for performance
  const limitNodes = (nodes: FigmaNode[], maxCount: number): FigmaNode[] => {
    const result: FigmaNode[] = []
    let count = 0

    const traverse = (nodeList: FigmaNode[]) => {
      for (const node of nodeList) {
        if (count >= maxCount) break
        result.push(node)
        count++
        if (node.children) {
          traverse(node.children)
        }
      }
    }

    traverse(nodes)
    return result
  }

  const limitedNodes = limitNodes(nodes, maxNodes)

  return (
    <div className="max-h-96 overflow-auto p-2 border rounded">
      {limitedNodes.map((node) => renderNode(node))}
      {nodes.length > maxNodes && (
        <div className="text-center text-gray-500 text-sm mt-2">
          显示前 {maxNodes} 个节点，共 {nodes.length} 个
        </div>
      )}
    </div>
  )
}
