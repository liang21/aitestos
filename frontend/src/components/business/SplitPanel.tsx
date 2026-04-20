import type { ReactNode } from 'react'
import { useState } from 'react'

export interface SplitPanelProps {
  /** Children must be SplitPanel.Left and SplitPanel.Right */
  children: ReactNode
  /** Initial split ratio (0-1), default 0.5 */
  defaultSplit?: number
  /** Minimum width for left panel (pixels) */
  minLeftWidth?: number
  /** Minimum width for right panel (pixels) */
  minRightWidth?: number
  /** Additional class name */
  className?: string
  /** Custom style */
  style?: React.CSSProperties
}

interface PanelProps {
  children: ReactNode
  className?: string
  style?: React.CSSProperties
}

const splitPanelContext = {
  splitRatio: 0.5,
  setSplitRatio: (_ratio: number) => {},
}

/**
 * Left panel component
 */
function LeftPanel({ children, className, style }: PanelProps) {
  return (
    <div
      className={className}
      style={{
        ...style,
        flex: '0 0 50%',
        overflow: 'auto',
      }}
    >
      {children}
    </div>
  )
}

/**
 * Right panel component
 */
function RightPanel({ children, className, style }: PanelProps) {
  return (
    <div
      className={className}
      style={{
        ...style,
        flex: '1 1 auto',
        overflow: 'auto',
        minWidth: 0,
      }}
    >
      {children}
    </div>
  )
}

/**
 * Split panel component with resizable divider
 * Used for side-by-side layouts like draft confirmation page
 */
export function SplitPanel({
  children,
  defaultSplit = 0.5,
  minLeftWidth = 200,
  minRightWidth = 200,
  className,
  style,
}: SplitPanelProps) {
  const [splitRatio, setSplitRatio] = useState(defaultSplit)
  const [isDragging, setIsDragging] = useState(false)

  const handleMouseDown = () => {
    setIsDragging(true)
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging) return

    const container = e.currentTarget as HTMLElement
    if (!container) return

    const containerRect = container.getBoundingClientRect()
    const containerWidth = containerRect.width

    // Guard against invalid container dimensions
    if (containerWidth <= 0) return

    // Clamp mouse position to container bounds
    const clientX = Math.max(containerRect.left, Math.min(containerRect.right, e.clientX))
    const newLeftWidth = clientX - containerRect.left

    // Calculate new split ratio with boundaries
    const minRatio = minLeftWidth / containerWidth
    const maxRatio = 1 - (minRightWidth / containerWidth)

    // Ensure minRatio < maxRatio
    const effectiveMinRatio = Math.max(0, Math.min(minRatio, 0.9))
    const effectiveMaxRatio = Math.max(0.1, Math.min(maxRatio, 1))

    const newRatio = Math.max(
      effectiveMinRatio,
      Math.min(effectiveMaxRatio, newLeftWidth / containerWidth)
    )

    setSplitRatio(newRatio)
  }

  const handleMouseUp = () => {
    setIsDragging(false)
  }

  // Find Left and Right panel children
  const childArray = Array.isArray(children) ? children : [children]
  const leftChild = childArray.find(
    (c) => (c as any)?.type === LeftPanel
  )
  const rightChild = childArray.find(
    (c) => (c as any)?.type === RightPanel
  )

  return (
    <div
      className={className}
      style={{
        ...style,
        display: 'flex',
        flexDirection: 'row',
        height: '100%',
        overflow: 'hidden',
      }}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseUp}
    >
      <div
        style={{
          flex: `0 0 ${splitRatio * 100}%`,
          overflow: 'auto',
        }}
      >
        {leftChild}
      </div>

      {/* Resizable divider */}
      <div
        className="resizer-divider"
        style={{
          width: '4px',
          cursor: 'col-resize',
          background: isDragging ? '#7B61FF' : '#E5E8EF',
          flex: '0 0 4px',
          transition: isDragging ? 'none' : 'background 0.2s',
          userSelect: 'none',
        }}
        onMouseDown={handleMouseDown}
      />

      <div
        style={{
          flex: '1 1 auto',
          overflow: 'auto',
          minWidth: minRightWidth,
        }}
      >
        {rightChild}
      </div>
    </div>
  )
}

SplitPanel.Left = LeftPanel
SplitPanel.Right = RightPanel
