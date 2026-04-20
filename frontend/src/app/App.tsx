import { Suspense } from 'react'
import { Outlet } from 'react-router-dom'

/**
 * Application Root Component
 *
 * Serves as the layout shell for the entire application.
 * The actual routing is defined in src/router/index.tsx
 */
export function App() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Outlet />
    </Suspense>
  )
}
