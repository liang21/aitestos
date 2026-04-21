/**
 * Login Page Object
 *
 * Encapsulates login page interactions and selectors
 */
import { expect } from '@playwright/test'
import { BasePage } from './BasePage'
import { E2E_CONFIG } from '../helpers/config'

export class LoginPage extends BasePage {
  // Selectors
  readonly emailInput = '[data-testid="login-email-input"]'
  readonly passwordInput = '[data-testid="login-password-input"]'
  readonly loginButton = '[data-testid="login-submit-button"]'
  readonly registerLink = '[data-testid="register-link"]'
  readonly errorMessage = '[data-testid="login-error-message"]'

  /**
   * Navigate to login page
   */
  async goto() {
    await super.goto('/login')
  }

  /**
   * Login with credentials
   */
  async login(email: string, password: string) {
    await this.page.fill(this.emailInput, email)
    await this.page.fill(this.passwordInput, password)
    await this.page.click(this.loginButton)

    // Wait for navigation to complete after successful login
    await this.page.waitForURL('**/projects', {
      timeout: E2E_CONFIG.NAVIGATION_TIMEOUT,
    })
  }

  /**
   * Login with default test credentials
   */
  async loginWithTestUser() {
    await this.login(E2E_CONFIG.TEST_USER.email, E2E_CONFIG.TEST_USER.password)
  }

  /**
   * Attempt login with invalid credentials (expecting failure)
   */
  async loginWithInvalidCredentials(email: string, password: string) {
    await this.page.fill(this.emailInput, email)
    await this.page.fill(this.passwordInput, password)
    await this.page.click(this.loginButton)

    // Should stay on login page
    await expect(this.page).toHaveURL('/login')

    // Error message should be visible
    await expect(this.page.locator(this.errorMessage)).toBeVisible()
  }

  /**
   * Navigate to registration page
   */
  async gotoRegister() {
    await this.page.click(this.registerLink)
    await expect(this.page).toHaveURL('/register')
  }

  /**
   * Check if currently on login page
   */
  async isOnLoginPage(): Promise<boolean> {
    return (
      (await this.exists(this.emailInput)) &&
      (await this.exists(this.passwordInput))
    )
  }
}
