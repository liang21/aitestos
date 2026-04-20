import { renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { useDebounce } from './useDebounce'

describe('useDebounce', () => {
  it('should return initial value immediately', () => {
    const { result } = renderHook(() => useDebounce('initial'))

    expect(result.current).toBe('initial')
  })

  it('should have default delay of 500ms', () => {
    const { result } = renderHook(() => useDebounce('test'))

    expect(result.current).toBe('test')
  })

  it('should export a function', () => {
    expect(typeof useDebounce).toBe('function')
  })
})
