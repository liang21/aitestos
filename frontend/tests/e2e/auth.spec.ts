import { test, expect } from '@playwright/test'
import { LoginPage } from './pages/LoginPage'
import { E2E_CONFIG, waitStrategies } from './helpers/config'
import { cleanupTestData } from './helpers/cleanup'

/**
 * E2E Tests: Authentication Flow
 * Covers TC-01 (successful login) and TC-02 (login failure)
 *
 * Improved with:
 * - Page Object Model
 * - Test data factories
 * - Centralized config
 * - Cleanup hooks
 */

test.describe('Authentication Flow', () => {
  let loginPage: LoginPage

  test.beforeEach(async ({ page }) => {
    loginPage = new LoginPage(page)
    await loginPage.goto()
  })

  test.afterAll(async ({ page }) => {
    await cleanupTestData(page)
  })

  test('should redirect to /projects on successful login (TC-01)', async ({
    page,
  }) => {
    // Act: Login with test credentials
    await loginPage.loginWithTestUser()

    // Assert: Should redirect to /projects
    await expect(page).toHaveURL('/projects')

    // Assert: Should show user avatar in sidebar
    await expect(page.locator('[data-testid="app-sidebar"]')).toBeVisible()
    await expect(
      page
        .locator('.username')
        .filter({ hasText: E2E_CONFIG.TEST_USER.username })
    ).toBeVisible()
  })

  test('should show error message on login failure (TC-02)', async ({
    page,
  }) => {
    // Act: Login with invalid password
    await loginPage.loginWithInvalidCredentials(
      'admin@aitestos.com',
      'WrongPassword'
    )

    // Assert: Error message should be visible
    await expect(
      page.locator('[data-testid="login-error-message"]')
    ).toBeVisible()

    // Assert: Should stay on login page
    await expect(page).toHaveURL('/login')
  })

  test('should redirect to /login on successful registration', async ({
    page,
  }) => {
    // Arrange: Navigate to registration
    await loginPage.gotoRegister()

    // Fill registration form
    const timestamp = Date.now()
    await page.fill(
      '[data-testid="register-username-input"]',
      `testuser_${timestamp}`
    )
    await page.fill(
      '[data-testid="register-email-input"]',
      `test_${timestamp}@example.com`
    )
    await page.fill('[data-testid="register-password-input"]', 'Test1234')
    await page.fill(
      '[data-testid="register-confirm-password-input"]',
      'Test1234'
    )

    // Select role
    await page.click('label:has-text("普通用户")')

    // Act: Submit registration
    await page.click('[data-testid="register-submit-button"]')

    // Assert: Should redirect to /login
    await expect(page).toHaveURL('/login')
  })

  test('should redirect to /login when accessing protected route without auth', async ({
    page,
  }) => {
    // Act: Try to access protected route directly
    await page.goto('/projects')

    // Assert: Should redirect to /login
    await expect(page).toHaveURL('/login')
  })

  test('should handle email validation error', async ({ page }) => {
    // Arrange: Enter invalid email
    await page.fill('[data-testid="login-email-input"]', 'invalid-email')

    // Act: Try to submit
    await page.click('[data-testid="login-submit-button"]')

    // Assert: Should show validation error
    await expect(page.locator('text=请输入有效的邮箱地址')).toBeVisible()
  })

  test('should handle password validation error', async ({ page }) => {
    // Arrange: Enter short password
    await page.fill('[data-testid="login-password-input"]', '123')

    // Act: Try to submit
    await page.click('[data-testid="login-submit-button"]')

    // Assert: Should show validation error
    await expect(page.locator('text=密码至少需要 8 个字符')).toBeVisible()
  })
})
