import { describe, it, expect } from 'vitest'
import { queryClient } from './query-client'

describe('QueryClient configuration', () => {
  it('should have default staleTime of 5 minutes', () => {
    expect(queryClient.getDefaultOptions().queries?.staleTime).toBe(
      5 * 60 * 1000
    )
  })

  it('should have retry count of 1', () => {
    expect(queryClient.getDefaultOptions().queries?.retry).toBe(1)
  })

  it('should have queryClient instance', () => {
    expect(queryClient).toBeDefined()
    expect(queryClient.constructor.name).toBe('QueryClient')
  })
})
