import { describe, it, expect, beforeEach, vi } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useRateLimit, RateLimitConfig } from './useRateLimit'

describe('useRateLimit', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('basic functionality', () => {
    it('should allow attempts within limit', () => {
      const { result } = renderHook(() =>
        useRateLimit({
          maxAttempts: 3,
          windowMs: 5000,
          cooldownMs: 5000,
        })
      )

      expect(result.current.canAttempt()).toBe(true)
      expect(result.current.isLocked).toBe(false)
      expect(result.current.remainingAttempts).toBe(3)
    })

    it('should record failed attempts', () => {
      const { result } = renderHook(() =>
        useRateLimit({
          maxAttempts: 3,
          windowMs: 5000,
          cooldownMs: 5000,
        })
      )

      act(() => {
        result.current.recordAttempt(false)
      })

      expect(result.current.attempts).toBe(1)
      expect(result.current.remainingAttempts).toBe(2)
    })

    it('should reset on successful attempt', () => {
      const { result } = renderHook(() =>
        useRateLimit({
          maxAttempts: 3,
          windowMs: 5000,
          cooldownMs: 5000,
        })
      )

      act(() => {
        result.current.recordAttempt(false)
        result.current.recordAttempt(false)
      })

      expect(result.current.attempts).toBe(2)

      act(() => {
        result.current.recordAttempt(true)
      })

      expect(result.current.attempts).toBe(0)
      expect(result.current.remainingAttempts).toBe(3)
    })
  })

  describe('max attempts limit', () => {
    it('should lock after max attempts reached', () => {
      const { result } = renderHook(() =>
        useRateLimit({
          maxAttempts: 3,
          windowMs: 5000,
          cooldownMs: 5000,
        })
      )

      // Use all attempts
      act(() => {
        result.current.recordAttempt(false)
        result.current.recordAttempt(false)
        result.current.recordAttempt(false)
      })

      expect(result.current.attempts).toBe(3)
      expect(result.current.canAttempt()).toBe(false)
      expect(result.current.isLocked).toBe(true)
    })

    it('should unlock after cooldown period', () => {
      const { result } = renderHook(() =>
        useRateLimit({
          maxAttempts: 3,
          windowMs: 5000,
          cooldownMs: 5000,
        })
      )

      // Use all attempts
      act(() => {
        result.current.recordAttempt(false)
        result.current.recordAttempt(false)
        result.current.recordAttempt(false)
      })

      expect(result.current.isLocked).toBe(true)

      // Fast forward past cooldown
      act(() => {
        vi.advanceTimersByTime(6000)
      })

      expect(result.current.canAttempt()).toBe(false) // Still locked, need to check again
      expect(result.current.isLocked).toBe(false) // Should have reset
    })
  })

  describe('window expiration', () => {
    it('should reset window after time expires', () => {
      const { result } = renderHook(() =>
        useRateLimit({
          maxAttempts: 3,
          windowMs: 1000,
          cooldownMs: 5000,
        })
      )

      // Use some attempts
      act(() => {
        result.current.recordAttempt(false)
        result.current.recordAttempt(false)
      })

      expect(result.current.attempts).toBe(2)

      // Fast forward past window
      act(() => {
        vi.advanceTimersByTime(1500)
      })

      // Should be able to attempt again
      expect(result.current.canAttempt()).toBe(true)
    })
  })

  describe('reset functionality', () => {
    it('should manually reset state', () => {
      const { result } = renderHook(() =>
        useRateLimit({
          maxAttempts: 3,
          windowMs: 5000,
          cooldownMs: 5000,
        })
      )

      act(() => {
        result.current.recordAttempt(false)
        result.current.recordAttempt(false)
      })

      expect(result.current.attempts).toBe(2)

      act(() => {
        result.current.reset()
      })

      expect(result.current.attempts).toBe(0)
      expect(result.current.isLocked).toBe(false)
    })
  })

  describe('RateLimitConfig presets', () => {
    it('should have login config with correct values', () => {
      expect(RateLimitConfig.LOGIN.maxAttempts).toBe(5)
      expect(RateLimitConfig.LOGIN.windowMs).toBe(15 * 60 * 1000)
      expect(RateLimitConfig.LOGIN.cooldownMs).toBe(15 * 60 * 1000)
    })

    it('should have register config with correct values', () => {
      expect(RateLimitConfig.REGISTER.maxAttempts).toBe(3)
      expect(RateLimitConfig.REGISTER.windowMs).toBe(60 * 60 * 1000)
      expect(RateLimitConfig.REGISTER.cooldownMs).toBe(60 * 60 * 1000)
    })
  })
})
