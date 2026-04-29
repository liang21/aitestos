/**
 * Virtual List Component
 * Provides efficient rendering for large lists using windowing technique
 * Based on @tanstack/react-virtual pattern
 */

import { useRef, useMemo, useState } from 'react'

interface VirtualListProps<T> {
  items: T[]
  itemHeight: number // Height of each item in pixels
  containerHeight: number // Height of the visible container
  renderItem: (item: T, index: number) => React.ReactNode
  overscan?: number // Number of items to render outside visible area
  className?: string
}

export function VirtualList<T>({
  items,
  itemHeight,
  containerHeight,
  renderItem,
  overscan = 3,
  className = '',
}: VirtualListProps<T>) {
  const containerRef = useRef<HTMLDivElement>(null)
  const [scrollTop, setScrollTop] = useState(0)

  // Calculate visible range based on scroll position
  const { visibleRange, totalHeight, offsetY } = useMemo(() => {
    const startIndex = Math.max(0, Math.floor(scrollTop / itemHeight) - overscan)
    const endIndex = Math.min(
      items.length - 1,
      Math.ceil((scrollTop + containerHeight) / itemHeight) + overscan
    )

    return {
      visibleRange: { start: startIndex, end: endIndex },
      totalHeight: items.length * itemHeight,
      offsetY: startIndex * itemHeight,
    }
  }, [items.length, itemHeight, containerHeight, overscan, scrollTop])

  const visibleItems = items.slice(visibleRange.start, visibleRange.end + 1)

  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    setScrollTop(e.currentTarget.scrollTop)
  }

  return (
    <div
      ref={containerRef}
      className={`overflow-auto ${className}`}
      style={{ height: containerHeight }}
      onScroll={handleScroll}
    >
      <div style={{ height: totalHeight, position: 'relative' }}>
        <div
          style={{
            transform: `translateY(${offsetY}px)`,
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
          }}
        >
          {visibleItems.map((item, index) => (
            <div
              key={(item as { id?: string }).id || index}
              style={{ height: itemHeight }}
            >
              {renderItem(item, visibleRange.start + index)}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

/**
 * Hook for managing virtual list state
 * Alternative simpler implementation for small to medium lists
 */
export function useVirtualList<T>(options: {
  items: T[]
  itemHeight: number
  containerHeight: number
  overscan?: number
}) {
  const { items, itemHeight, containerHeight, overscan = 3 } = options

  const [scrollTop, setScrollTop] = useState(0)

  const visibleRange = useMemo(() => {
    const startIndex = Math.max(0, Math.floor(scrollTop / itemHeight) - overscan)
    const endIndex = Math.min(
      items.length - 1,
      Math.ceil((scrollTop + containerHeight) / itemHeight) + overscan
    )
    return { start: startIndex, end: endIndex }
  }, [items.length, itemHeight, containerHeight, overscan, scrollTop])

  const totalHeight = items.length * itemHeight
  const offsetY = visibleRange.start * itemHeight

  const visibleItems = items.slice(visibleRange.start, visibleRange.end + 1)

  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    setScrollTop(e.currentTarget.scrollTop)
  }

  return {
    visibleItems,
    totalHeight,
    offsetY,
    handleScroll,
    visibleRange,
  }
}
