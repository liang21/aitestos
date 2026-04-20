import { describe, it, expect } from 'vitest'
import { cn } from './utils'

describe('cn', () => {
  it('should merge class names correctly', () => {
    expect(cn('px-4', 'py-2')).toBe('px-4 py-2')
  })

  it('should handle conflicting tailwind classes', () => {
    expect(cn('px-4', 'px-2')).toBe('px-2')
  })

  it('should handle conditional classes', () => {
    const isActive = true
    const isInactive = false
    expect(
      cn('base-class', isActive && 'active', isInactive && 'inactive')
    ).toBe('base-class active')
  })

  it('should handle empty inputs', () => {
    expect(cn()).toBe('')
  })

  it('should handle arrays and objects', () => {
    expect(cn(['class1', 'class2'], { class3: true, class4: false })).toBe(
      'class1 class2 class3'
    )
  })
})
