/**
 * Project List Page Object
 */
import { expect } from '@playwright/test'
import { BasePage } from './BasePage'

export class ProjectListPage extends BasePage {
  // Selectors
  readonly newProjectButton = '[data-testid="new-project-button"]'
  readonly projectTable = '[data-testid="project-list-table"]'
  readonly projectSearchInput = '[data-testid="project-search-input"]'
  readonly projectNameCell = '[data-testid="project-name-cell"]'

  /**
   * Navigate to projects page
   */
  async goto() {
    await super.goto('/projects')
  }

  /**
   * Click "New Project" button
   */
  async clickNewProject() {
    await this.page.click(this.newProjectButton)
    // Wait for modal to appear
    await this.page.waitForSelector('[data-testid="create-project-modal"]', {
      state: 'visible',
    })
  }

  /**
   * Search for a project by name
   */
  async searchProject(keyword: string) {
    await this.page.fill(this.projectSearchInput, keyword)
    // Wait for search to complete
    await this.page.waitForTimeout(300)
  }

  /**
   * Click on a project by name
   */
  async openProject(projectName: string) {
    await this.page.click(`[data-project-name="${projectName}"]`)
    await this.page.waitForURL(`/projects/*`)
  }

  /**
   * Get project count
   */
  async getProjectCount(): Promise<number> {
    return await this.page
      .locator(this.projectTable)
      .locator('tbody tr')
      .count()
  }

  /**
   * Check if project exists in list
   */
  async hasProject(projectName: string): Promise<boolean> {
    return (
      (await this.page
        .locator(`[data-project-name="${projectName}"]`)
        .count()) > 0
    )
  }

  /**
   * Wait for project to appear in list
   */
  async waitForProject(projectName: string, timeout = 5000) {
    await this.page.waitForSelector(`[data-project-name="${projectName}"]`, {
      state: 'visible',
      timeout,
    })
  }
}
