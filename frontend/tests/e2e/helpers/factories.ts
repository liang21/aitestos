/**
 * E2E Test Data Factories
 *
 * Generates consistent test data with optional overrides
 */

/**
 * Generate a test project
 */
export function createProject(
  overrides?: Partial<{ name: string; prefix: string; description: string }>
) {
  const timestamp = Date.now().toString().slice(-6)
  return {
    name: overrides?.name || `E2E Test Project ${timestamp}`,
    prefix: overrides?.prefix || 'E2E',
    description: overrides?.description || 'E2E automated test project',
  }
}

/**
 * Generate a test module
 */
export function createModule(
  overrides?: Partial<{ name: string; abbreviation: string }>
) {
  return {
    name: overrides?.name || '用户中心',
    abbreviation: overrides?.abbreviation || 'USR',
  }
}

/**
 * Generate test credentials
 */
export function createTestCredentials() {
  return {
    email: 'admin@aitestos.com',
    password: 'Test1234',
  }
}

/**
 * Generate AI generation task data
 */
export function createGenerationTask(
  overrides?: Partial<{ requirement: string; count: number }>
) {
  return {
    requirement:
      overrides?.requirement || '测试用户注册功能，包括邮箱验证和密码强度校验',
    count: overrides?.count || 5,
  }
}

/**
 * Generate test plan data
 */
export function createTestPlan(
  overrides?: Partial<{ name: string; description: string }>
) {
  const timestamp = Date.now().toString().slice(-6)
  return {
    name: overrides?.name || `Sprint E2E ${timestamp}`,
    description: overrides?.description || 'E2E automated test plan',
  }
}
