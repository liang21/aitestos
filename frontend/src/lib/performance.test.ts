import { describe, expect, it, vi, beforeEach } from 'vitest'
import { debounce, throttle, propsAreEqual, perfMonitor } from './performance'

describe('Performance Utilities', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('debounce', () => {
    it('should delay function execution', () => {
      const fn = vi.fn()
      const debouncedFn = debounce(fn, 300)

      debouncedFn()
      expect(fn).not.toHaveBeenCalled()

      vi.advanceTimersByTime(300)
      expect(fn).toHaveBeenCalledTimes(1)
    })

    it('should reset timer on subsequent calls', () => {
      const fn = vi.fn()
      const debouncedFn = debounce(fn, 300)

      debouncedFn()
      vi.advanceTimersByTime(100)
      debouncedFn()
      vi.advanceTimersByTime(100)
      debouncedFn()
      vi.advanceTimersByTime(100)

      expect(fn).not.toHaveBeenCalled()

      vi.advanceTimersByTime(200)
      expect(fn).toHaveBeenCalledTimes(1)
    })
  })

  describe('throttle', () => {
    it('should limit function execution rate', () => {
      const fn = vi.fn().mockReturnValue('result')
      const throttledFn = throttle(fn, 300)

      const result1 = throttledFn()
      const result2 = throttledFn()
      const result3 = throttledFn()

      expect(fn).toHaveBeenCalledTimes(1)
      expect(result1).toBe('result')
      expect(result2).toBe('result')
      expect(result3).toBe('result')

      vi.advanceTimersByTime(300)

      const result4 = throttledFn()
      expect(fn).toHaveBeenCalledTimes(2)
      expect(result4).toBe('result')
    })
  })

  describe('propsAreEqual', () => {
    it('should return true for identical objects', () => {
      const obj1 = { a: 1, b: 2 }
      const obj2 = { a: 1, b: 2 }

      expect(propsAreEqual(obj1, obj2)).toBe(true)
    })

    it('should return false for different objects', () => {
      const obj1 = { a: 1, b: 2 }
      const obj2 = { a: 1, b: 3 }

      expect(propsAreEqual(obj1, obj2)).toBe(false)
    })

    it('should handle arrays correctly', () => {
      const obj1 = { items: [1, 2, 3] }
      const obj2 = { items: [1, 2, 3] }
      const obj3 = { items: [1, 2, 3, 4] }

      expect(propsAreEqual(obj1, obj2)).toBe(true)
      expect(propsAreEqual(obj1, obj3)).toBe(false)
    })

    it('should handle nested objects', () => {
      const obj1 = { nested: { a: 1 } }
      const obj2 = { nested: { a: 1 } }

      expect(propsAreEqual(obj1, obj2)).toBe(true)
    })
  })

  describe('PerformanceMonitor', () => {
    it('should mark and measure performance', () => {
      const consoleSpy = vi.spyOn(console, 'log').mockReturnValue(undefined)

      perfMonitor.mark('test-start')
      vi.advanceTimersByTime(100)
      perfMonitor.measure('test-operation', 'test-start')

      expect(consoleSpy).toHaveBeenCalled()
      const logCall = consoleSpy.mock.calls[0]
      expect(logCall[0]).toContain('test-operation')
      expect(logCall[0]).toContain('ms')

      consoleSpy.mockRestore()
    })

    it('should warn when measuring non-existent mark', () => {
      const consoleWarnSpy = vi
        .spyOn(console, 'warn')
        .mockReturnValue(undefined)

      perfMonitor.measure('non-existent', 'missing-mark')

      expect(consoleWarnSpy).toHaveBeenCalledWith(
        expect.stringContaining('missing-mark')
      )

      consoleWarnSpy.mockRestore()
    })

    it('should measure async operations', async () => {
      const consoleSpy = vi.spyOn(console, 'log').mockReturnValue(undefined)

      const asyncFn = async () => {
        return 'result'
      }

      const result = await perfMonitor.measureAsync('async-op', asyncFn)

      expect(result).toBe('result')
      expect(consoleSpy).toHaveBeenCalled()
      const logMessage = consoleSpy.mock.calls[0][0] as string
      expect(logMessage).toContain('async-op')
      expect(logMessage).toContain('ms')

      consoleSpy.mockRestore()
    })
  })
})
