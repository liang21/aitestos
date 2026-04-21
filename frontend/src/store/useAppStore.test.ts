import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useAppStore, useSidebarResponsive } from './useAppStore'

describe('useAppStore', () => {
  // Reset store state before each test
  beforeEach(() => {
    // Reset to initial state
    useAppStore.setState({ sidebarCollapsed: false })
  })

  afterEach(() => {
    // Cleanup after each test
    useAppStore.setState({ sidebarCollapsed: false })
  })

  it('should have sidebarCollapsed default value as false', () => {
    const { result } = renderHook(() => useAppStore())
    expect(result.current.sidebarCollapsed).toBe(false)
  })

  it('should toggle sidebarCollapsed when toggleSidebar is called', () => {
    const { result } = renderHook(() => useAppStore())

    expect(result.current.sidebarCollapsed).toBe(false)

    act(() => {
      result.current.toggleSidebar()
    })

    expect(result.current.sidebarCollapsed).toBe(true)

    act(() => {
      result.current.toggleSidebar()
    })

    expect(result.current.sidebarCollapsed).toBe(false)
  })

  it('should respond to window width changes when using useSidebarResponsive', () => {
    const { result } = renderHook(() => useSidebarResponsive())

    // Initial state
    expect(result.current.sidebarCollapsed).toBe(false)

    // Simulate window resize to small screen
    act(() => {
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 500,
      })
      window.dispatchEvent(new Event('resize'))
    })

    // Should collapse on small screen
    expect(result.current.sidebarCollapsed).toBe(true)

    // Simulate window resize to large screen
    act(() => {
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 1024,
      })
      window.dispatchEvent(new Event('resize'))
    })

    // Should expand on large screen
    expect(result.current.sidebarCollapsed).toBe(false)
  })
})
