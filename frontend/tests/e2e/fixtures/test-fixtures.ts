/**
 * E2E Test Fixtures
 *
 * Extends Playwright test with custom fixtures
 */
import { test as base } from '@playwright/test'
import { LoginPage } from '../pages/LoginPage'
import { ProjectListPage } from '../pages/ProjectListPage'
import { CreateProjectModal } from '../pages/CreateProjectModal'
import { cleanupTestData, generateCleanupId } from '../helpers/cleanup'

type Fixtures = {
  loginPage: LoginPage
  projectListPage: ProjectListPage
  createProjectModal: CreateProjectModal
}

// Extend base test with custom fixtures
export const test = base.extend<Fixtures>({
  loginPage: async ({ page }, use) => {
    const loginPage = new LoginPage(page)
    await use(loginPage)
  },

  projectListPage: async ({ page }, use) => {
    const projectListPage = new ProjectListPage(page)
    await use(projectListPage)
  },

  createProjectModal: async ({ page }, use) => {
    const createProjectModal = new CreateProjectModal(page)
    await use(createProjectModal)
  },
})

// Re-export everything from Playwright
export { expect } from '@playwright/test'

/**
 * Global setup for E2E tests
 */
export async function globalSetup() {
  console.log('E2E Test Suite: Global Setup')
  // Could reset test database here
}

/**
 * Global teardown for E2E tests
 */
export async function globalTeardown() {
  console.log('E2E Test Suite: Global Teardown')
  // Could cleanup test artifacts here
}
