/**
 * E2E Test Configuration
 */

export const E2E_CONFIG = {
  /** Timeout for page navigation */
  NAVIGATION_TIMEOUT: 5000,

  /** Timeout for element visibility */
  ELEMENT_TIMEOUT: 10000,

  /** Timeout for async operations (AI generation, etc) */
  ASYNC_OPERATION_TIMEOUT: 60000,

  /** Base URL for the application */
  BASE_URL: process.env.BASE_URL || 'http://localhost:5174',

  /** Test user credentials */
  TEST_USER: {
    email: 'admin@aitestos.com',
    password: 'Test1234',
    username: 'admin',
  },

  /** Project prefix for E2E tests */
  PROJECT_PREFIX: 'E2E',
} as const

/**
 * Reusable wait strategies
 */
export const waitStrategies = {
  /** Wait for navigation to complete */
  async forNavigation(page: import('@playwright/test').Page) {
    await page.waitForLoadState('networkidle', {
      timeout: E2E_CONFIG.NAVIGATION_TIMEOUT,
    })
  },

  /** Wait for element to be visible */
  async forElement(page: import('@playwright/test').Page, selector: string) {
    await page.waitForSelector(selector, {
      state: 'visible',
      timeout: E2E_CONFIG.ELEMENT_TIMEOUT,
    })
  },

  /** Wait for task completion (polling) */
  async forTaskCompletion(
    page: import('@playwright/test').Page,
    timeout = E2E_CONFIG.ASYNC_OPERATION_TIMEOUT
  ) {
    await page.waitForSelector(
      '[data-status="completed"], [data-status="failed"]',
      {
        timeout,
      }
    )
  },
}
