/**
 * Performance optimization utilities
 */

/**
 * Debounce function execution
 * Useful for search inputs, resize handlers, etc.
 */
export function debounce<T extends (...args: unknown[]) => unknown>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: NodeJS.Timeout

  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn(...args), delay)
  }
}

/**
 * Throttle function execution
 * Useful for scroll handlers, resize handlers, etc.
 */
export function throttle<T extends (...args: unknown[]) => unknown>(
  fn: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle = false
  let lastResult: ReturnType<T>

  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      inThrottle = true
      lastResult = fn(...args) as ReturnType<T>
      setTimeout(() => (inThrottle = false), limit)
    }
    return lastResult
  }
}

/**
 * Batch state updates to reduce re-renders
 * Use this when multiple state updates happen in quick succession
 */
export function batchUpdates<T>(updates: Array<() => T>): T[] {
  return updates.map(update => update())
}

/**
 * Memoize component props comparison
 * Use with React.memo for fine-grained re-render control
 */
export function propsAreEqual(
  prevProps: Record<string, unknown>,
  nextProps: Record<string, unknown>
): boolean {
  const keys = Object.keys(nextProps)

  for (const key of keys) {
    if (prevProps[key] !== nextProps[key]) {
      // Special handling for arrays and objects
      if (Array.isArray(nextProps[key])) {
        if (
          !Array.isArray(prevProps[key]) ||
          nextProps[key].length !== prevProps[key].length
        ) {
          return false
        }
      } else if (typeof nextProps[key] === 'object' && nextProps[key] !== null) {
        if (typeof prevProps[key] !== 'object' || prevProps[key] === null) {
          return false
        }
      } else {
        return false
      }
    }
  }

  return true
}

/**
 * Lazy load component
 * Delays component loading until it's needed
 */
export function lazyLoad<T extends React.ComponentType<any>>(
  componentLoader: () => Promise<{ default: T }>
): React.LazyExoticComponent<T> {
  return React.lazy(componentLoader)
}

/**
 * Create a memoized component with custom comparison
 */
export function memoWithComparison<T extends Record<string, unknown>>(
  Component: React.ComponentType<T>,
  areEqual?: (prevProps: T, nextProps: T) => boolean
): React.MemoExoticComponent<React.ComponentType<T>> {
  return React.memo(Component, areEqual)
}

/**
 * Performance monitoring utilities
 */
export class PerformanceMonitor {
  private marks = new Map<string, number>()

  mark(name: string): void {
    this.marks.set(name, performance.now())
  }

  measure(name: string, startMark?: string): number {
    const end = performance.now()
    const start = startMark
      ? this.marks.get(startMark)
      : this.marks.get(name)

    if (start === undefined) {
      console.warn(`Mark "${startMark || name}" not found`)
      return 0
    }

    const duration = end - start
    console.log(`[Performance] ${name}: ${duration.toFixed(2)}ms`)

    return duration
  }

  measureAsync<T>(
    name: string,
    fn: () => Promise<T>
  ): Promise<T> {
    const start = performance.now()
    return fn().finally(() => {
      const duration = performance.now() - start
      console.log(`[Performance] ${name}: ${duration.toFixed(2)}ms`)
    })
  }
}

/**
 * Create a performance monitor instance
 */
export const perfMonitor = new PerformanceMonitor()

/**
 * Hook for measuring component render time
 */
export function useRenderPerf(componentName: string) {
  if (process.env.NODE_ENV === 'development') {
    useEffect(() => {
      perfMonitor.mark(`${componentName}-render-start`)
      return () => {
        perfMonitor.measure(`${componentName}-render`, `${componentName}-render-start`)
      }
    })
  }
}

/**
 * Optimize list rendering by ensuring stable keys
 * Warns if unstable keys are detected
 */
export function useStableKeys<T>(
  items: T[],
  keyExtractor: (item: T) => string | number
): void {
  if (process.env.NODE_ENV === 'development') {
    const prevKeysRef = useRef<Set<string | number>>(new Set())

    useEffect(() => {
      const currentKeys = new Set(items.map(keyExtractor))
      const prevKeys = prevKeysRef.current

      // Check for key stability
      items.forEach((item, index) => {
        const key = keyExtractor(item)
        if (
          prevKeys.has(key) &&
          Array.from(prevKeys).indexOf(key) !== index
        ) {
          console.warn(
            `[Performance Warning] Unstable key detected at index ${index}: ${key}. This may cause unnecessary re-renders.`
          )
        }
      })

      prevKeysRef.current = currentKeys
    })
  }
}

// Add missing imports
import React, { useRef, useEffect, useMemo, useState } from 'react'
