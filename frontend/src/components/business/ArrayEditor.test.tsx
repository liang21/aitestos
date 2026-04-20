import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, vi } from 'vitest'
import { ArrayEditor } from './ArrayEditor'

describe('ArrayEditor', () => {
  it('should render initial single row when value is empty', () => {
    const onChange = vi.fn()
    render(
      <ArrayEditor value={[]} onChange={onChange} placeholder="请输入内容" />
    )

    const inputs = screen.getAllByRole('textbox')
    expect(inputs.length).toBe(1)
  })

  it('should render all rows from initial value', () => {
    const onChange = vi.fn()
    const initialValue = ['步骤1', '步骤2', '步骤3']
    render(
      <ArrayEditor value={initialValue} onChange={onChange} placeholder="请输入" />
    )

    const inputs = screen.getAllByRole('textbox')
    expect(inputs.length).toBe(3)
    expect(inputs[0]).toHaveValue('步骤1')
    expect(inputs[1]).toHaveValue('步骤2')
    expect(inputs[2]).toHaveValue('步骤3')
  })

  it('should add new row when add button is clicked', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()
    render(
      <ArrayEditor value={['步骤1']} onChange={onChange} placeholder="请输入" />
    )

    const addButton = screen.getByRole('button', { name: /添加|add/i })
    await user.click(addButton)

    const inputs = screen.getAllByRole('textbox')
    expect(inputs.length).toBe(2)
    expect(onChange).toHaveBeenCalled()
  })

  it('should remove row when delete button is clicked', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()
    render(
      <ArrayEditor value={['步骤1', '步骤2', '步骤3']} onChange={onChange} placeholder="请输入" />
    )

    const deleteButtons = screen.getAllByRole('button', { name: /删除/i })
    await user.click(deleteButtons[1]) // Delete second row

    const inputs = screen.getAllByRole('textbox')
    expect(inputs.length).toBe(2)
    expect(onChange).toHaveBeenCalled()
  })

  it('should call onChange with latest value on input change', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()
    render(
      <ArrayEditor value={['步骤1']} onChange={onChange} placeholder="请输入" />
    )

    const input = screen.getByRole('textbox')
    await user.clear(input)
    await user.type(input, '新步骤')

    expect(onChange).toHaveBeenCalledWith(['新步骤'])
  })

  it('should always keep at least one row when trying to delete last row', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()
    render(
      <ArrayEditor value={['步骤1']} onChange={onChange} placeholder="请输入" />
    )

    // When only one row, delete button should not be visible
    const deleteButtons = screen.queryAllByRole('button', { name: /删除/i })
    expect(deleteButtons.length).toBe(0)

    // Should still have one input
    const inputs = screen.getAllByRole('textbox')
    expect(inputs.length).toBe(1)
  })
})
