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
          // Projects
          {
            path: 'projects',
            lazy: () =>
              import('../features/projects/components/ProjectListPage').then(
                (m) => ({ Component: m.ProjectListPage })
              ),
          },
          {
            path: 'projects/:projectId',
            lazy: () =>
              import('../features/projects/components/ProjectDashboard').then(
                (m) => ({ Component: m.ProjectDashboard })
              ),
          },
          {
            path: 'projects/:projectId/modules',
            lazy: () =>
              import('../features/modules/components/ModuleManagePage').then(
                (m) => ({ Component: m.ModuleManagePage })
              ),
          },
          // Test Cases
          {
            path: 'testcases',
            lazy: () =>
              import('../features/testcases/components/CaseListPage').then(
                (m) => ({ Component: m.CaseListPage })
              ),
          },
          {
            path: 'testcases/:caseId',
            lazy: () =>
              import('../features/testcases/components/CaseDetailPage').then(
                (m) => ({ Component: m.CaseDetailPage })
              ),
          },
          // Test Plans
          {
            path: 'plans',
            lazy: () =>
              import('../features/plans/components/PlanListPage').then((m) => ({
                Component: m.PlanListPage,
              })),
          },
          {
            path: 'plans/new',
            lazy: () =>
              import('../features/plans/components/NewPlanPage').then((m) => ({
                Component: m.NewPlanPage,
              })),
          },
          {
            path: 'plans/:planId',
            lazy: () =>
              import('../features/plans/components/PlanDetailPage').then((m) => ({
                Component: m.PlanDetailPage,
              })),
          },
          // AI Generation
          {
            path: 'generation',
            lazy: () =>
              import('../features/generation/components/GenerationTaskListPage').then(
                (m) => ({ Component: m.GenerationTaskListPage })
              ),
          },
          {
            path: 'generation/tasks/new',
            lazy: () =>
              import('../features/generation/components/NewGenerationTaskPage').then(
                (m) => ({ Component: m.NewGenerationTaskPage })
              ),
          },
          {
            path: 'generation/tasks/:taskId',
            lazy: () =>
              import('../features/generation/components/TaskDetailPage').then(
                (m) => ({ Component: m.TaskDetailPage })
              ),
          },
          // Drafts
          {
            path: 'drafts',
            lazy: () =>
              import('../features/drafts/components/DraftListPage').then(
                (m) => ({ Component: m.DraftListPage })
              ),
          },
          {
            path: 'drafts/:draftId',
            lazy: () =>
              import('../features/drafts/components/DraftConfirmPage').then(
                (m) => ({ Component: m.DraftConfirmPage })
              ),
          },
          // Documents/Knowledge Base
          {
            path: 'documents',
            lazy: () =>
              import('../features/documents/components/KnowledgeListPage').then(
                (m) => ({ Component: m.KnowledgeListPage })
              ),
          },
          {
            path: 'documents/:documentId',
            lazy: () =>
              import('../features/documents/components/DocumentDetailPage').then(
                (m) => ({ Component: m.DocumentDetailPage })
              ),
          },
          // Configs
          {
            path: 'projects/:projectId/configs',
            lazy: () =>
              import('../features/configs/components/ConfigManagePage').then(
                (m) => ({ Component: m.ConfigManagePage })
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
