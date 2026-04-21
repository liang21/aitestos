/**
 * Create Project Modal Page Object
 */
import { BasePage } from './BasePage'
import { createProject } from '../helpers/factories'

export class CreateProjectModal extends BasePage {
  // Selectors
  readonly modal = '[data-testid="create-project-modal"]'
  readonly nameInput = '[data-testid="project-name-input"]'
  readonly prefixInput = '[data-testid="project-prefix-input"]'
  readonly descriptionInput = '[data-testid="project-description-input"]'
  readonly submitButton = '[data-testid="project-submit-button"]'
  readonly cancelButton = '[data-testid="project-cancel-button"]'

  /**
   * Create a project with given data
   */
  async createProject(
    data?: Partial<{ name: string; prefix: string; description: string }>
  ) {
    const projectData = createProject(data)

    await this.page.fill(this.nameInput, projectData.name)
    await this.page.fill(this.prefixInput, projectData.prefix)
    await this.page.fill(this.descriptionInput, projectData.description)

    await this.page.click(this.submitButton)

    // Wait for modal to close
    await this.page.waitForSelector(this.modal, {
      state: 'hidden',
      timeout: 5000,
    })

    return projectData
  }

  /**
   * Create project with default factory data
   */
  async createWithDefaults() {
    return await this.createProject()
  }

  /**
   * Cancel the modal
   */
  async cancel() {
    await this.page.click(this.cancelButton)
    await this.page.waitForSelector(this.modal, { state: 'hidden' })
  }

  /**
   * Check if modal is visible
   */
  async isVisible(): Promise<boolean> {
    return await this.page.locator(this.modal).isVisible()
  }
}
