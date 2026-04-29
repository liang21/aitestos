import { afterAll, afterEach, beforeAll, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import '@testing-library/jest-dom/vitest'
import { server } from './msw/server'
import { setupTestLocalStorage, clearTestLocalStorage } from './utils/test-localStorage'

// Start MSW server before all tests
beforeAll(() => {
  // Setup test localStorage
  setupTestLocalStorage()

  // Set environment variable for API base URL
  process.env.VITE_API_BASE_URL = '/api/v1'
  server.listen({ onUnhandledRequest: 'error' })
})

// Reset handlers after each test
afterEach(() => {
  server.resetHandlers()
  cleanup()
  clearTestLocalStorage()
})

// Close MSW server after all tests
afterAll(() => server.close())

// Mock window.matchMedia for Arco Design responsive components
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Note: Arco Design components are not mocked here to preserve all functionality
// Message side effects in tests can be suppressed by using vi.spyOn if needed
// localStorage is provided by jsdom, no need to mock it
