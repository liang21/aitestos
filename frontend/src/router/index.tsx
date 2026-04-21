import { createBrowserRouter, Navigate } from 'react-router-dom'
import { RouteGuard } from '@/router/RouteGuard'
import { AuthErrorBoundary } from '@/components/ErrorBoundary'
import { App } from '@/app/App'
import { AppLayout } from '@/components/layout/AppLayout'

/**
 * Application Router Configuration
 *
 * Defines all application routes with lazy loading
 * Public routes: /login, /register
 * Protected routes: wrapped with RouteGuard and AuthErrorBoundary
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
      // Public routes (also wrapped with error boundary)
      {
        path: '/login',
        lazy: () =>
          import('../features/auth/components/LoginPage').then((m) => ({
            Component: () => (
              <AuthErrorBoundary>
                <m.LoginPage />
              </AuthErrorBoundary>
            ),
          })),
      },
      {
        path: '/register',
        lazy: () =>
          import('../features/auth/components/RegisterPage').then((m) => ({
            Component: () => (
              <AuthErrorBoundary>
                <m.RegisterPage />
              </AuthErrorBoundary>
            ),
          })),
      },
      // Protected routes (require authentication)
      {
        path: '/',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <AppLayout />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
        children: [
          {
            path: 'projects',
            lazy: () =>
              import('../features/projects/components/ProjectListPage').then(
                (m) => ({ Component: m.ProjectListPage })
              ),
          },
          {
            path: 'testcases',
            lazy: () =>
              import('../features/testcases/components/CaseListPage').then(
                (m) => ({ Component: m.CaseListPage })
              ),
          },
          {
            path: 'plans',
            lazy: () =>
              import('../features/plans/components/PlanListPage').then((m) => ({
                Component: m.PlanListPage,
              })),
          },
          {
            path: 'generation',
            lazy: () =>
              import('../features/generation/components/GenerationTaskListPage').then(
                (m) => ({ Component: m.GenerationTaskListPage })
              ),
          },
          {
            path: 'drafts',
            lazy: () =>
              import('../features/drafts/components/DraftListPage').then(
                (m) => ({ Component: m.DraftListPage })
              ),
          },
          {
            path: 'documents',
            lazy: () =>
              import('../features/documents/components/KnowledgeListPage').then(
                (m) => ({ Component: m.KnowledgeListPage })
              ),
          },
          {
            path: '*',
            element: <Navigate to="/projects" replace />,
          },
        ],
      },
      // 404 fallback
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
