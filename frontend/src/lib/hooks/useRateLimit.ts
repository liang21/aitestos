import { useState, useCallback } from 'react'

/**
 * Rate limit configuration
 */
interface RateLimitConfig {
  maxAttempts: number
  windowMs: number
  cooldownMs: number
}

/**
 * Rate limit state
 */
interface RateLimitState {
  attempts: number
  isLocked: boolean
  remainingTime: number
}

/**
 * useRateLimit Hook
 *
 * Implements rate limiting for sensitive operations like login/register
 * - Max attempts per time window
 * - Cooldown period after max attempts reached
 * - Automatic reset after window expires
 */
export function useRateLimit(config: RateLimitConfig) {
  const { maxAttempts, windowMs, cooldownMs } = config

  const [state, setState] = useState<RateLimitState>({
    attempts: 0,
    isLocked: false,
    remainingTime: 0,
  })

  const [windowStart, setWindowStart] = useState<number>(0)
  const [lockUntil, setLockUntil] = useState<number>(0)

  /**
   * Check if rate limit allows the action
   */
  const canAttempt = useCallback((): boolean => {
    const now = Date.now()

    // Check if currently locked
    if (state.isLocked && now < lockUntil) {
      // Update remaining time
      setState((prev) => ({
        ...prev,
        remainingTime: Math.ceil((lockUntil - now) / 1000),
      }))
      return false
    }

    // Reset if lock period expired
    if (state.isLocked && now >= lockUntil) {
      setState({
        attempts: 0,
        isLocked: false,
        remainingTime: 0,
      })
      setWindowStart(0)
      setLockUntil(0)
      return true
    }

    // Reset window if expired
    if (windowStart > 0 && now - windowStart > windowMs) {
      setState((prev) => ({ ...prev, attempts: 0 }))
      setWindowStart(0)
      return true
    }

    // Check if max attempts reached
    if (state.attempts >= maxAttempts) {
      const lockTime = now + cooldownMs
      setLockUntil(lockTime)
      setState({
        ...state,
        isLocked: true,
        remainingTime: Math.ceil(cooldownMs / 1000),
      })
      return false
    }

    return true
  }, [state, windowStart, lockUntil, maxAttempts, windowMs, cooldownMs])

  /**
   * Record an attempt (should be called after action)
   */
  const recordAttempt = useCallback(
    (success: boolean): void => {
      if (success) {
        // Reset on successful attempt
        setState({
          attempts: 0,
          isLocked: false,
          remainingTime: 0,
        })
        setWindowStart(0)
      } else {
        // Increment attempts on failure
        const now = Date.now()
        if (windowStart === 0) {
          setWindowStart(now)
        }

        setState((prev) => {
          const newAttempts = prev.attempts + 1

          // Check if should lock
          if (newAttempts >= maxAttempts) {
            const lockTime = now + cooldownMs
            setLockUntil(lockTime)
            return {
              attempts: newAttempts,
              isLocked: true,
              remainingTime: Math.ceil(cooldownMs / 1000),
            }
          }

          return {
            ...prev,
            attempts: newAttempts,
          }
        })
      }
    },
    [windowStart, maxAttempts, cooldownMs]
  )

  /**
   * Reset the rate limit state
   */
  const reset = useCallback((): void => {
    setState({
      attempts: 0,
      isLocked: false,
      remainingTime: 0,
    })
    setWindowStart(0)
    setLockUntil(0)
  }, [])

  /**
   * Get remaining attempts
   */
  const getRemainingAttempts = useCallback((): number => {
    return Math.max(0, maxAttempts - state.attempts)
  }, [maxAttempts, state.attempts])

  return {
    canAttempt,
    recordAttempt,
    reset,
    isLocked: state.isLocked,
    attempts: state.attempts,
    remainingAttempts: getRemainingAttempts(),
    remainingTime: state.remainingTime,
  }
}

/**
 * Default rate limit configurations
 */
export const RateLimitConfig = {
  LOGIN: {
    maxAttempts: 5,
    windowMs: 15 * 60 * 1000, // 15 minutes
    cooldownMs: 15 * 60 * 1000, // 15 minutes
  },
  REGISTER: {
    maxAttempts: 3,
    windowMs: 60 * 60 * 1000, // 1 hour
    cooldownMs: 60 * 60 * 1000, // 1 hour
  },
} as const
