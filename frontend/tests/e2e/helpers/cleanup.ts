/**
 * E2E Test Data Cleanup Utilities
 */

import type { Page } from '@playwright/test'

const cleanupMarkers = new Set<string>()

/**
 * Mark a resource for cleanup
 */
export function markForCleanup(
  id: string,
  type: 'project' | 'module' | 'plan' | 'testcase'
) {
  cleanupMarkers.add(`${type}:${id}`)
}

/**
 * Generate a cleanup ID
 */
export function generateCleanupId(
  type: 'project' | 'module' | 'plan' | 'testcase'
): string {
  const id = `${type}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
  markForCleanup(id, type)
  return id
}

/**
 * Cleanup test data via API calls
 * This should be called in test.afterAll hooks
 */
export async function cleanupTestData(page: Page) {
  // Get auth token from localStorage
  const token = await page.evaluate(() => localStorage.getItem('access_token'))

  if (!token) {
    console.warn('No auth token found, skipping cleanup')
    return
  }

  // Cleanup each marked resource
  for (const marker of cleanupMarkers) {
    const [type, id] = marker.split(':')

    try {
      await page.request.delete(`/api/v1/${type}s/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      console.log(`Cleaned up ${type}: ${id}`)
    } catch (error) {
      console.warn(`Failed to cleanup ${type}:${id}`, error)
    }
  }

  cleanupMarkers.clear()
}

/**
 * Cleanup all test data (for CI environment)
 */
export async function cleanupAllTestData() {
  // In CI, we might want to truncate all test data
  // This should be implemented based on your backend API
  console.log('Full cleanup not implemented - use test database reset in CI')
}
