/**
 * Test utilities for localStorage
 * Provides safe localStorage operations for test environments
 */

export function createTestLocalStorage() {
  const store = new Map<string, string>()

  return {
    getItem: (key: string) => store.get(key) ?? null,
    setItem: (key: string, value: string) => store.set(key, value),
    removeItem: (key: string) => store.delete(key),
    clear: () => store.clear(),
    get length() {
      return store.size
    },
    key: (index: number) => {
      const keys = Array.from(store.keys())
      return keys[index] ?? null
    },
  }
}

export function setupTestLocalStorage() {
  // Check if localStorage is available
  const testLocalStorage = () => {
    try {
      localStorage.setItem('test', 'test')
      localStorage.removeItem('test')
      return true
    } catch {
      return false
    }
  }

  if (!testLocalStorage()) {
    // Create mock localStorage
    const mockLocalStorage = createTestLocalStorage()
    Object.defineProperty(window, 'localStorage', {
      value: mockLocalStorage,
      writable: true,
    })
  }
}

export function clearTestLocalStorage() {
  try {
    if (localStorage.clear) {
      localStorage.clear()
    } else {
      // Fallback: remove all items
      const keys = Object.keys(localStorage)
      keys.forEach(key => localStorage.removeItem(key))
    }
  } catch {
    // Ignore errors
  }
}
