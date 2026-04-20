import { Button, Input, Space } from '@arco-design/web-react'
import { IconMinus, IconPlus } from '@arco-design/web-react/icon'
import { useEffect, useRef, useState } from 'react'

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
 * Check if two string arrays are deeply equal
 */
function arraysEqual(a: string[], b: string[]): boolean {
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false
  }
  return true
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
  // Sync internal state with external value prop
  const [rows, setRows] = useState<string[]>(value.length > 0 ? value : [''])
  const prevValueRef = useRef<string[]>()

  useEffect(() => {
    // Only update if value actually changed (deep comparison)
    if (!arraysEqual(value, prevValueRef.current ?? [])) {
      if (value.length > 0) {
        setRows(value)
      }
      prevValueRef.current = value
    }
  }, [value])

  const handleInputChange = (index: number, newValue: string) => {
    const newRows = [...rows]
    newRows[index] = newValue
    setRows(newRows)
    onChange(newRows)
  }

  const handleAddRow = () => {
    if (maxRows && rows.length >= maxRows) return
    const newRows = [...rows, '']
    setRows(newRows)
    onChange(newRows)
  }

  const handleRemoveRow = (index: number) => {
    if (rows.length <= minRows) return
    const newRows = rows.filter((_, i) => i !== index)
    setRows(newRows)
    onChange(newRows)
  }

  return (
    <Space
      className={className}
      direction="vertical"
      size="small"
      style={{ width: '100%' }}
    >
      {rows.map((row, index) => (
        <div key={index} className="flex gap-2 items-center">
          <Input
            value={row}
            onChange={(val) => handleInputChange(index, val)}
            placeholder={placeholder}
            className="flex-1"
          />
          {showRemoveButton && rows.length > minRows && (
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
      {showAddButton && (!maxRows || rows.length < maxRows) && (
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
