/**
 * E2E Test Global Setup
 *
 * This file is run once before all E2E tests
 */

import { chromium, FullConfig } from '@playwright/test'

async function globalSetup(config: FullConfig) {
  console.log('🚀 E2E Test Suite: Starting global setup')

  // Could perform the following:
  // 1. Start/ensure test database is ready
  // 2. Seed initial test data
  // 3. Ensure test services are running

  const browser = await chromium.launch()
  const page = await browser.newPage()

  // Check if application is accessible
  try {
    await page.goto(process.env.BASE_URL || 'http://localhost:5174')
    const title = await page.title()
    console.log(`✅ Application is accessible: ${title}`)
  } catch (error) {
    console.error(
      '❌ Application is not accessible. Make sure dev server is running.'
    )
    throw new Error('E2E tests require application to be running')
  } finally {
    await browser.close()
  }

  console.log('✅ E2E Test Suite: Global setup completed')
}

export default globalSetup
