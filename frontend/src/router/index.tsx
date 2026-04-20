import { createBrowserRouter, Navigate, Route } from 'react-router-dom'
import { RouteGuard } from './RouteGuard'
import { App } from '../app/App'

/**
 * Application Router Configuration
 *
 * Defines all application routes with lazy loading
 * Public routes: /login, /register
 * Protected routes: wrapped with RouteGuard
 */
export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {
        index: true,
        element: <Navigate to="/projects" replace />,
      },
      // Public routes
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
      // Protected routes (require authentication)
      {
        path: '/',
        element: <RouteGuard />,
        children: [
          {
            path: 'projects',
            lazy: () =>
              import('../features/projects/components/ProjectListPage').then(
                (m) => ({ Component: m.ProjectListPage })
              ),
          },
          // More protected routes will be added here
          {
            path: '*',
            element: <Navigate to="/projects" replace />,
          },
        ],
      },
      // 404 fallback
      {
        path: '*',
        element: () => import('../components/NotFoundPage').then((m) => ({
          Component: m.NotFoundPage,
        })),
      },
    ],
  },
])

