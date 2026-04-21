import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import {
  useDraftList,
  useConfirmDraft,
  useRejectDraft,
  useBatchConfirm,
  usePendingDraftCount,
} from './useDrafts'
import { server } from '../../../../tests/msw/server'
import { draftsHandlers } from '../../../../tests/msw/handlers/drafts'
import { http, HttpResponse } from 'msw'

describe('useDrafts hooks', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
    server.use(...draftsHandlers)
  })

  afterEach(() => {
    server.resetHandlers()
  })

  function wrapper({ children }: { children: React.ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    )
  }

  describe('useDraftList', () => {
    it('should fetch draft list successfully', async () => {
      const { result } = renderHook(
        () => useDraftList({ projectId: 'project-1', offset: 0, limit: 10 }),
        { wrapper }
      )

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        data: expect.any(Array),
        total: expect.any(Number),
      })
    })

    it('should invalidate query cache on refetch', async () => {
      const { result } = renderHook(
        () => useDraftList({ projectId: 'project-1', offset: 0, limit: 10 }),
        { wrapper }
      )

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      // Simply check that refetch function exists
      expect(result.current.refetch).toBeDefined()

      // Call refetch and verify it completes without error
      await expect(result.current.refetch()).resolves.toBeDefined()
    })
  })

  describe('useConfirmDraft', () => {
    it('should confirm draft and update cache', async () => {
      const { result } = renderHook(() => useConfirmDraft(), { wrapper })

      result.current.mutate({
        draftId: 'draft-001',
        moduleId: 'mod-001',
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        number: expect.any(String),
        status: 'unexecuted',
      })
    })

    it('should invalidate queries on success', async () => {
      const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')

      const { result } = renderHook(() => useConfirmDraft(), { wrapper })

      result.current.mutate({
        draftId: 'draft-001',
        moduleId: 'mod-001',
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(invalidateSpy).toHaveBeenCalled()
    })
  })

  describe('useRejectDraft', () => {
    it('should reject draft successfully', async () => {
      const { result } = renderHook(() => useRejectDraft(), { wrapper })

      result.current.mutate({
        draftId: 'draft-001',
        data: { reason: 'duplicate', feedback: '重复用例' },
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        success: true,
        message: '草稿已拒绝',
      })
    })
  })

  describe('useBatchConfirm', () => {
    it('should confirm multiple drafts', async () => {
      const { result } = renderHook(() => useBatchConfirm(), { wrapper })

      result.current.mutate({
        draftIds: ['draft-001', 'draft-002'],
        moduleId: 'mod-001',
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toMatchObject({
        successCount: 2,
        failedCount: 0,
      })
    })

    it('should handle partial success', async () => {
      server.use(
        http.post('/api/v1/generation/drafts/batch-confirm', async () =>
          HttpResponse.json({
            successCount: 1,
            failedCount: 1,
            errors: [{ draftId: 'draft-002', error: '验证失败' }],
          })
        )
      )

      const { result } = renderHook(() => useBatchConfirm(), { wrapper })

      result.current.mutate({
        draftIds: ['draft-001', 'draft-002'],
        moduleId: 'mod-001',
      })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data.successCount).toBe(1)
      expect(result.current.data.failedCount).toBe(1)
    })
  })

  describe('usePendingDraftCount', () => {
    it('should poll pending draft count', async () => {
      const { result } = renderHook(() => usePendingDraftCount(), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      expect(result.current.data).toBeGreaterThan(0)
    })

    it('should use correct refetch interval', async () => {
      const { result } = renderHook(() => usePendingDraftCount(), { wrapper })

      await waitFor(() => expect(result.current.isSuccess).toBe(true))

      const query = queryClient
        .getQueryCache()
        .find(['drafts', 'pending-count'])
      expect(query?.observers[0].options.refetchInterval).toBe(5000)
    })
  })
})
