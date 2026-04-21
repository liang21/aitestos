/**
 * Base Page Object
 *
 * Provides common functionality for all page objects
 */
import type { Page, Locator } from '@playwright/test'

export class BasePage {
  readonly page: Page

  constructor(page: Page) {
    this.page = page
  }

  /**
   * Navigate to a URL
   */
  async goto(path: string) {
    await this.page.goto(path)
    await this.page.waitForLoadState('networkidle')
  }

  /**
   * Wait for element to be visible
   */
  async waitForVisible(selector: string, timeout = 5000) {
    await this.page.waitForSelector(selector, { state: 'visible', timeout })
  }

  /**
   * Wait for element to be hidden
   */
  async waitForHidden(selector: string, timeout = 5000) {
    await this.page.waitForSelector(selector, { state: 'hidden', timeout })
  }

  /**
   * Click element by data-testid
   */
  async clickByTestId(testId: string) {
    await this.page.click(`[data-testid="${testId}"]`)
  }

  /**
   * Fill input by data-testid
   */
  async fillByTestId(testId: string, value: string) {
    await this.page.fill(`[data-testid="${testId}"]`, value)
  }

  /**
   * Get text content by data-testid
   */
  async getTextByTestId(testId: string): Promise<string> {
    return (
      (await this.page.locator(`[data-testid="${testId}"]`).textContent()) || ''
    )
  }

  /**
   * Check if element exists
   */
  async exists(selector: string): Promise<boolean> {
    return (await this.page.locator(selector).count()) > 0
  }

  /**
   * Take screenshot on error
   */
  async screenshotOnFailure(testName: string) {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-')
    await this.page.screenshot({
      path: `test-results/screenshots/${testName}-${timestamp}.png`,
      fullPage: true,
    })
  }
}
