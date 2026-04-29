import { Empty, Pagination, Spin, Table } from '@arco-design/web-react'
import type { PaginationProps, TableProps } from '@arco-design/web-react'
import { IconLoading } from '@arco-design/web-react/icon'
import { useMemo } from 'react'

export interface SearchTableProps<T = unknown> extends Omit<
  TableProps,
  'pagination'
> {
  /** Table columns definition */
  columns: TableProps<T>['columns']
  /** Table data */
  data: T[]
  /** Total count for pagination */
  total: number
  /** Loading state */
  loading?: boolean
  /** Current page number (1-indexed) */
  current?: number
  /** Page size */
  pageSize?: number
  /** Pagination change callback */
  onPageChange?: (page: number, pageSize: number) => void
  /** Empty description */
  emptyText?: string
  /** Additional class name */
  className?: string
  /** Custom style */
  style?: React.CSSProperties
}

/**
 * Unified search table component with pagination and loading states
 * Wraps Arco Design Table with consistent behavior
 */
export function SearchTable<T extends Record<string, unknown>>({
  columns,
  data,
  total,
  loading = false,
  current = 1,
  pageSize = 10,
  onPageChange,
  emptyText = '暂无数据',
  className,
  style,
  ...restProps
}: SearchTableProps<T>) {
  // Calculate pagination props
  const paginationProps = useMemo((): PaginationProps | false => {
    if (!onPageChange && total <= pageSize) {
      return false
    }

    return {
      current,
      pageSize,
      total,
      sizeCanChange: false,
      showTotal: true,
      onChange: onPageChange,
      className: 'flex justify-end mt-4',
    }
  }, [current, pageSize, total, onPageChange])

  return (
    <div className={className} style={style}>
      <Spin
        loading={loading}
        icon={<IconLoading spin />}
        style={{ width: '100%', display: 'block' }}
      >
        <Table<T>
          columns={columns}
          data={data}
          pagination={paginationProps}
          rowKey={(record, index) => {
            if (typeof record === 'object' && record !== null) {
              const id = (record as Record<string, unknown>).id ?? (record as Record<string, unknown>).key
              if (typeof id === 'string' || typeof id === 'number') return String(id)
            }
            // Use index as fallback - warn in development
            if (process.env.NODE_ENV === 'development') {
              console.warn(
                'SearchTable: record missing unique "id" or "key" property, using index as fallback. This may cause rendering issues.'
              )
            }
            return `row-${index}`
          }}
          noDataElement={<Empty description={emptyText} />}
          border
          stripe
          {...restProps}
        />
      </Spin>
    </div>
  )
}
