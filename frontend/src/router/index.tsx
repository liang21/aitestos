import { createBrowserRouter, Navigate } from 'react-router-dom'
import { App } from '../app/App'

/**
 * Application Router Configuration
 *
 * This is a minimal router setup. Routes will be added as features are implemented.
 * For now, it redirects root to a placeholder path.
 */
export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {
        index: true,
        element: <Navigate to="/login" replace />,
      },
      {
        path: '/login',
        lazy: () =>
          import('../features/auth/components/LoginPage').then((m) => ({
            Component: m.LoginPage,
          })),
      },
      {
        path: '/register',
        lazy: () =>
          import('../features/auth/components/RegisterPage').then((m) => ({
            Component: m.RegisterPage,
          })),
      },
      {
        path: '*',
        lazy: () =>
          import('../components/NotFoundPage').then((m) => ({
            Component: m.NotFoundPage,
          })),
      },
    ],
  },
])
