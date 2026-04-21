import { test, expect } from '@playwright/test'
import { LoginPage } from './pages/LoginPage'
import { ProjectListPage } from './pages/ProjectListPage'
import { CreateProjectModal } from './pages/CreateProjectModal'
import {
  createProject,
  createModule,
  generateCleanupId,
} from './helpers/factories'
import { cleanupTestData } from './helpers/cleanup'
import { E2E_CONFIG } from './helpers/config'

/**
 * E2E Tests: Project Management Flow
 * Covers TC-03 (create project) and project dashboard
 *
 * Improved with:
 * - Page Object Model
 * - Test data factories with cleanup IDs
 * - Better test isolation
 */

test.describe('Project Management Flow', () => {
  let loginPage: LoginPage
  let projectListPage: ProjectListPage
  let createProjectModal: CreateProjectModal

  test.beforeEach(async ({ page }) => {
    loginPage = new LoginPage(page)
    projectListPage = new ProjectListPage(page)
    createProjectModal = new CreateProjectModal(page)

    // Login before each test
    await loginPage.goto()
    await loginPage.loginWithTestUser()
  })

  test.afterAll(async ({ page }) => {
    await cleanupTestData(page)
  })

  test('should create project successfully (TC-03)', async ({ page }) => {
    // Act: Click "New Project" button
    await projectListPage.clickNewProject()

    // Arrange: Create project with cleanup ID for later cleanup
    const cleanupId = generateCleanupId('project')
    const projectData = createProject({
      name: `E2E Project ${cleanupId.split('-')[2]}`,
    })

    // Fill and submit form
    await page.fill('[data-testid="project-name-input"]', projectData.name)
    await page.fill('[data-testid="project-prefix-input"]', projectData.prefix)
    await page.fill(
      '[data-testid="project-description-input"]',
      projectData.description
    )
    await page.click('[data-testid="project-submit-button"]')

    // Assert: Should show success message
    await expect(page.locator('text=创建成功')).toBeVisible()

    // Assert: New project should appear in list
    await projectListPage.waitForProject(projectData.name, 5000)
    expect(await projectListPage.hasProject(projectData.name)).toBe(true)
  })

  test('should show validation error for invalid project prefix', async ({
    page,
  }) => {
    // Act: Click "New Project" button
    await projectListPage.clickNewProject()

    // Fill with invalid prefix (not uppercase letters)
    await page.fill('[data-testid="project-name-input"]', 'Test Project')
    await page.fill('[data-testid="project-prefix-input"]', 'invalid')

    // Try to submit
    await page.click('[data-testid="project-submit-button"]')

    // Assert: Should show validation error
    await expect(page.locator('text=前缀必须是 2-4 位大写字母')).toBeVisible()

    // Cancel the modal
    await createProjectModal.cancel()
  })

  test('should view project dashboard', async ({ page }) => {
    // Arrange: Click on first project in list
    const projectCount = await projectListPage.getProjectCount()
    if (projectCount === 0) {
      test.skip('No projects available')
    }

    await page.click('table tbody tr:first-child')

    // Assert: Should navigate to project detail
    await expect(
      page.locator('h1, h2').filter({ hasText: /仪表盘/ })
    ).toBeVisible()

    // Assert: Should show stats cards
    await expect(page.locator('.arco-statistic').first()).toBeVisible()

    // Assert: Should show at least 4 stats
    const statsCount = await page.locator('.arco-statistic').count()
    expect(statsCount).toBeGreaterThanOrEqual(4)
  })

  test('should create and delete module in project', async ({ page }) => {
    // Arrange: Navigate to first project
    const projectCount = await projectListPage.getProjectCount()
    if (projectCount === 0) {
      test.skip('No projects available')
    }

    await page.click('table tbody tr:first-child')
    await page.click('text=项目设置')
    await page.click('text=模块管理')

    // Get initial module count
    const initialCount = await page.locator('table tbody tr').count()

    // Act: Create module
    await page.click('button:has-text("新建模块")')

    const moduleData = createModule()
    await page.fill('[data-testid="module-name-input"]', moduleData.name)
    await page.fill(
      '[data-testid="module-abbreviation-input"]',
      moduleData.abbreviation
    )
    await page.click('button:has-text("确定")')

    // Assert: Module should appear
    await expect(page.locator(`text=${moduleData.name}`)).toBeVisible()

    // Act: Delete the module
    await page.hover('table tbody tr:last-child')
    await page.click('button:has-text("删除")')
    await page.click('.arco-popconfirm-popup button:has-text("确定")')

    // Assert: Count should decrease
    const finalCount = await page.locator('table tbody tr').count()
    expect(finalCount).toBe(initialCount)
  })

  test('should search for projects', async ({ page }) => {
    // Arrange: Create a test project first
    await projectListPage.clickNewProject()
    const searchKeyword = 'SEARCH_TEST'
    await createProjectModal.createProject({
      name: `${searchKeyword} Project`,
      prefix: 'SRH',
    })

    // Act: Search for the project
    await projectListPage.searchProject(searchKeyword)

    // Assert: Should show only matching projects
    const visibleProjects = await page.locator('table tbody tr').count()
    expect(visibleProjects).toBeGreaterThan(0)

    // Verify the search keyword is in the results
    await expect(page.locator(`text=${searchKeyword}`)).toBeVisible()
  })
})
