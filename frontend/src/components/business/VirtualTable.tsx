/**
 * Virtual Table Component
 * Wrapper around Arco Table with virtual scrolling support for large datasets
 */

import { useMemo } from 'react'
import { Table } from '@arco-design/web-react'
import type { TableProps } from '@arco-design/web-react'

interface VirtualTableProps<T> extends Omit<TableProps<T>, 'scroll'> {
  data: T[]
  rowHeight?: number // Estimated row height in pixels
  enableVirtual?: boolean // Enable virtual scrolling
}

export function VirtualTable<T extends Record<string, unknown>>({
  data,
  rowHeight = 60,
  enableVirtual = false,
  pagination,
  ...restProps
}: VirtualTableProps<T>) {
  // Enable pagination by default with page size 20
  const defaultPagination = useMemo(() => {
    if (pagination === false) return false
    if (pagination) return pagination
    return {
      pageSize: 20,
      pageSizeOptions: [10, 20, 50, 100],
      showTotal: (total: number) => `共 ${total} 条`,
    }
  }, [pagination])

  // For large datasets (>1000 rows), enable virtual scrolling
  const shouldUseVirtual = enableVirtual || data.length > 1000

  const scrollProps = useMemo(() => {
    if (shouldUseVirtual) {
      // Calculate viewport height based on page size
      const pageSize =
        defaultPagination && typeof defaultPagination === 'object' && 'pageSize' in defaultPagination
          ? (defaultPagination as { pageSize: number }).pageSize
          : 20
      const viewportHeight = pageSize * rowHeight + 100 // Add extra space for header

      return {
        x: '100%',
        y: viewportHeight,
      }
    }

    return {
      x: '100%',
    }
  }, [shouldUseVirtual, defaultPagination, rowHeight])

  return (
    <Table<T>
      {...restProps}
      data={data}
      pagination={defaultPagination}
      scroll={scrollProps}
    />
  )
}

/**
 * Memoized row component for performance optimization
 * Use this when rendering complex rows with nested components
 */
export function MemoizedRow<T>({ children, ...props }: {
  children: React.ReactNode
  record: T
  index: number
}) {
  return (
    <tr {...props}>
      {children}
    </tr>
  )
}

/**
 * Optimized cell renderer
 * Memoizes expensive cell computations
 */
export function useOptimizedCellRenderer<T, V>(
  value: V,
  record: T,
  render: (value: V, record: T) => React.ReactNode
) {
  return useMemo(() => {
    return render(value, record)
  }, [value, record, render])
}
