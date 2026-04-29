import { Button, Input, Space } from '@arco-design/web-react'
import { IconMinus, IconPlus } from '@arco-design/web-react/icon'
import { useMemo } from 'react'

export interface ArrayEditorProps {
  /** Current array value */
  value: string[]
  /** Change callback with updated array */
  onChange: (value: string[]) => void
  /** Input placeholder text */
  placeholder?: string
  /** Minimum number of rows (default 1) */
  minRows?: number
  /** Maximum number of rows (default no limit) */
  maxRows?: number
  /** Additional class name */
  className?: string
  /** Show add button (default true) */
  showAddButton?: boolean
  /** Show remove button (default true) */
  showRemoveButton?: boolean
}

/**
 * Dynamic array editor component
 * Used for editing test case steps, preconditions, etc.
 */
export function ArrayEditor({
  value,
  onChange,
  placeholder = '请输入内容',
  minRows = 1,
  maxRows,
  showAddButton = true,
  showRemoveButton = true,
  className,
}: ArrayEditorProps) {
  // Use value directly as source of truth, default to one empty row
  const displayRows = useMemo(() => {
    return value.length > 0 ? value : ['']
  }, [value])

  const handleInputChange = (index: number, newValue: string) => {
    const newRows = [...displayRows]
    newRows[index] = newValue
    onChange(newRows)
  }

  const handleAddRow = () => {
    if (maxRows && displayRows.length >= maxRows) return
    const newRows = [...displayRows, '']
    onChange(newRows)
  }

  const handleRemoveRow = (index: number) => {
    if (displayRows.length <= minRows) return
    const newRows = displayRows.filter((_, i) => i !== index)
    onChange(newRows)
  }

  return (
    <Space
      className={className}
      direction="vertical"
      size="small"
      style={{ width: '100%' }}
    >
      {displayRows.map((row, index) => (
        <div key={index} className="flex gap-2 items-center">
          <Input
            value={row}
            onChange={(val) => handleInputChange(index, val)}
            placeholder={placeholder}
            className="flex-1"
          />
          {showRemoveButton && displayRows.length > minRows && (
            <Button
              type="text"
              status="danger"
              icon={<IconMinus />}
              aria-label="删除"
              onClick={() => handleRemoveRow(index)}
            />
          )}
        </div>
      ))}
      {showAddButton && (!maxRows || displayRows.length < maxRows) && (
        <Button
          type="outline"
          icon={<IconPlus />}
          onClick={handleAddRow}
          className="w-full"
        >
          添加
        </Button>
      )}
    </Space>
  )
}
