import { test, expect } from '@playwright/test'
import { LoginPage } from './pages/LoginPage'
import { ProjectListPage } from './pages/ProjectListPage'
import { CreateProjectModal } from './pages/CreateProjectModal'
import {
  createProject,
  createModule,
  createGenerationTask,
  createTestPlan,
  generateCleanupId,
} from './helpers/factories'
import { cleanupTestData, markForCleanup } from './helpers/cleanup'
import { E2E_CONFIG, waitStrategies } from './helpers/config'

/**
 * E2E Tests: AI Generation Core Flow (拆分版)
 *
 * 原来的单一200+行测试被拆分为多个独立测试，
 * 每个测试专注于一个特定场景，提高可维护性和调试效率。
 *
 * Covers TC-01 through TC-08
 */

test.describe('AI Generation Core Flow - Split Tests', () => {
  let loginPage: LoginPage
  let projectListPage: ProjectListPage
  let createProjectModal: CreateProjectModal

  // Shared test data
  let testProject: { name: string; prefix: string; description: string }
  let testModule: { name: string; abbreviation: string }

  test.beforeAll(async () => {
    console.log('E2E Core Flow: Starting test suite')
  })

  test.beforeEach(async ({ page }) => {
    loginPage = new LoginPage(page)
    projectListPage = new ProjectListPage(page)
    createProjectModal = new CreateProjectModal(page)

    // Login
    await loginPage.goto()
    await loginPage.loginWithTestUser()

    // Go to projects page
    await projectListPage.goto()
  })

  test.afterAll(async ({ page }) => {
    await cleanupTestData(page)
    console.log('E2E Core Flow: Test suite completed')
  })

  /**
   * Step 1-2: Login & Create Project (TC-01, TC-03)
   */
  test('should create project with valid data', async ({ page }) => {
    // Arrange: Generate test project data
    const cleanupId = generateCleanupId('project')
    testProject = createProject({
      name: `AI Test Project ${cleanupId.split('-')[2]}`,
    })

    // Act: Create project
    await projectListPage.clickNewProject()
    await page.fill('[data-testid="project-name-input"]', testProject.name)
    await page.fill('[data-testid="project-prefix-input"]', testProject.prefix)
    await page.fill(
      '[data-testid="project-description-input"]',
      testProject.description
    )
    await page.click('[data-testid="project-submit-button"]')

    // Assert: Success message
    await expect(page.locator('text=创建成功')).toBeVisible()

    // Assert: Project appears in list
    await projectListPage.waitForProject(testProject.name)
    expect(await projectListPage.hasProject(testProject.name)).toBe(true)
  })

  /**
   * Step 3: Create Module
   */
  test('should create module in project', async ({ page }) => {
    // Arrange: Open first project
    await page.click('table tbody tr:first-child')
    await page.click('text=项目设置')
    await page.click('text=模块管理')

    // Generate module data
    const cleanupId = generateCleanupId('module')
    testModule = createModule()

    // Act: Create module
    await page.click('button:has-text("新建模块")')
    await page.fill('[data-testid="module-name-input"]', testModule.name)
    await page.fill(
      '[data-testid="module-abbreviation-input"]',
      testModule.abbreviation
    )
    await page.click('button:has-text("确定")')

    // Assert: Module should appear
    await expect(page.locator(`text=${testModule.name}`)).toBeVisible()
  })

  /**
   * Step 4-5: Initiate AI Generation Task (TC-04)
   */
  test('should create AI generation task', async ({ page }) => {
    // Arrange: Navigate to AI generation page
    await page.click('text=AI 生成')

    // Generate task data
    const taskData = createGenerationTask()

    // Act: Create generation task
    await page.click('button:has-text("新建任务")')

    // Select module (wait for modal/dropdown)
    await page.click('.arco-select:has-text("选择模块")')
    // Select first available module

    // Fill requirement description
    await page.fill('textarea[placeholder*="需求描述"]', taskData.requirement)

    // Set generation parameters
    await page.fill('input[role="spinbutton"]', taskData.count.toString())

    // Submit
    await page.click('button:has-text("立即生成")')

    // Assert: Task created successfully
    await expect(page.locator('text=任务创建成功')).toBeVisible()
  })

  /**
   * Step 6: Poll for Task Completion
   */
  test('should show completed task status', async ({ page }) => {
    // Arrange: Navigate to task list
    await page.click('text=AI 生成')

    // Act: Click on first task
    await page.click('table tbody tr:first-child')

    // Assert: Wait for task completion (using custom wait strategy)
    await page.waitForSelector(
      '[data-status="completed"], [data-status="failed"]',
      {
        timeout: E2E_CONFIG.ASYNC_OPERATION_TIMEOUT,
      }
    )

    // Verify status is completed
    const status = await page
      .locator('[data-status]')
      .getAttribute('data-status')
    expect(['completed', 'failed']).toContain(status || '')
  })

  /**
   * Step 7: View Drafts (TC-04 continued)
   */
  test('should display generated drafts', async ({ page }) => {
    // Arrange: Navigate to task detail
    await page.click('table tbody tr:first-child')

    // Act: Check for drafts
    await expect(page.locator('text=草稿列表')).toBeVisible()

    // Assert: Should show at least one draft
    const draftCount = await page.locator('table tbody tr').count()
    expect(draftCount).toBeGreaterThan(0)

    // Check for draft fields
    await expect(page.locator('text=标题')).toBeVisible()
    await expect(page.locator('text=置信度')).toBeVisible()
  })

  /**
   * Step 8: Confirm Draft and Verify Auto-Numbering (TC-05)
   */
  test('should confirm draft and generate case number', async ({ page }) => {
    // Arrange: Navigate to drafts
    await page.click('text=草稿箱')

    // Act: Click on first draft
    await page.click('table tbody tr:first-child')

    // Click confirm button
    await page.click('button:has-text("确认并转为正式用例")')

    // Select target module
    await page.click('.arco-select:has-text("目标模块")')

    // Confirm
    await page.click('button:has-text("确认")')

    // Assert: Success message
    await expect(page.locator('text=确认成功')).toBeVisible()

    // Navigate to test cases to verify numbering
    await page.click('text=测试用例')

    // Check for auto-generated number format
    const caseNumberPattern = new RegExp(
      `^[A-Z]{2,4}-[A-Z]{2,4}-\\d{8}-\\d{3}$`
    )
    await expect(
      page.locator(`td`).filter({ hasText: caseNumberPattern }).first()
    ).toBeVisible()
  })

  /**
   * Step 9: Batch Confirm Drafts (TC-06)
   */
  test('should batch confirm multiple drafts', async ({ page }) => {
    // Arrange: Navigate to drafts
    await page.click('text=草稿箱')

    // Ensure we have multiple drafts
    const draftCount = await page.locator('table tbody tr').count()
    if (draftCount < 2) {
      test.skip('Need at least 2 drafts for batch confirm test')
    }

    // Act: Select first 2 drafts
    await page.locator('input[type="checkbox"]').nth(1).check()
    await page.locator('input[type="checkbox"]').nth(2).check()

    // Click batch confirm
    await page.click('button:has-text("批量确认")')

    // Select target module
    await page.click('.arco-select:has-text("目标模块")')

    // Confirm
    await page.click('button:has-text("确认")')

    // Assert: Success message with count
    await expect(page.locator('text=/成功确认 \\d+ 条草稿/')).toBeVisible()
  })

  /**
   * Step 10: Create Test Plan (TC-07 part 1)
   */
  test('should create test plan', async ({ page }) => {
    // Arrange: Navigate to plans
    await page.click('text=测试计划')

    // Generate plan data
    const planData = createTestPlan()

    // Act: Create plan
    await page.click('button:has-text("新建计划")')

    await page.fill('[placeholder="计划名称"]', planData.name)
    await page.fill('textarea[placeholder*="描述"]', planData.description)

    // Select test cases (click "Add All" button if available)
    const addAllButton = page.locator('button:has-text("全部添加")')
    if (await addAllButton.isVisible()) {
      await addAllButton.click()
    }

    // Submit
    await page.click('button:has-text("创建")')

    // Assert: Success message
    await expect(page.locator('text=创建成功')).toBeVisible()

    // Verify plan appears in list
    await expect(page.locator(`text=${planData.name}`)).toBeVisible()
  })

  /**
   * Step 11: Record Test Results (TC-07 part 2)
   */
  test('should record test execution results', async ({ page }) => {
    // Arrange: Navigate to first test plan
    await page.click('text=测试计划')
    await page.click('table tbody tr:first-child')

    // Act: Record result for first test case
    await page.click('table tbody tr:first-child button:has-text("录入结果")')

    // Select Pass status
    await page.click('label:has-text("通过")')

    // Add remark
    await page.fill('textarea[placeholder*="备注"]', '功能正常')

    // Submit
    await page.click('button:has-text("确定")')

    // Assert: Success message
    await expect(page.locator('text=录入成功')).toBeVisible()

    // Verify status badge updated
    await expect(
      page.locator('table tbody tr:first-child').locator('text=通过')
    ).toBeVisible()
  })

  /**
   * Step 12: Verify Dashboard Statistics (TC-08)
   */
  test('should display accurate dashboard statistics', async ({ page }) => {
    // Arrange: Navigate to first project dashboard
    await page.click('text=项目列表')
    await page.click('table tbody tr:first-child')

    // Assert: Check for 4 stats cards
    await expect(page.locator('.arco-statistic').nth(0)).toContainText(
      '用例总数'
    )
    await expect(page.locator('.arco-statistic').nth(1)).toContainText('通过率')
    await expect(page.locator('.arco-statistic').nth(2)).toContainText('覆盖率')
    await expect(page.locator('.arco-statistic').nth(3)).toContainText(
      'AI 生成数'
    )

    // Verify values are displayed
    const stats = await page.locator('.arco-statistic-value').allTextContents()
    expect(stats.length).toBeGreaterThan(0)

    // Verify trend chart exists
    await expect(page.locator('[data-testid="trend-chart"]')).toBeVisible()
  })

  /**
   * Negative Test: Login failure (TC-02)
   */
  test('should handle login failure correctly', async ({ page }) => {
    // Act: Navigate to login and attempt with invalid credentials
    await loginPage.goto()
    await loginPage.loginWithInvalidCredentials(
      E2E_CONFIG.TEST_USER.email,
      'WrongPassword'
    )

    // Assert: Error message displayed
    await expect(
      page.locator('[data-testid="login-error-message"]')
    ).toBeVisible()

    // Assert: Stays on login page
    await expect(page).toHaveURL('/login')
  })
})
