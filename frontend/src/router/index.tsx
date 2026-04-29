/**
 * Application Router Configuration
 *
 * Defines all application routes with lazy loading and proper nesting.
 * Supports both new project-scoped routes and legacy routes for backward compatibility.
 */

import { createBrowserRouter, Navigate, useParams } from 'react-router-dom'
import { RouteGuard } from '@/router/RouteGuard'
import { AuthErrorBoundary } from '@/components/ErrorBoundary'
import { App } from '@/app/App'
import { AppLayout } from '@/components/layout/AppLayout'
import { LegacyRouteRedirect, LegacyRouteRedirectWithParams } from '@/components/LegacyRouteRedirect'

/**
 * Application Router
 *
 * Route structure:
 * - Public routes: /login, /register
 * - Project-scoped routes: /projects/:projectId/*
 * - Global routes: /drafts, /drafts/:draftId
 * - Legacy redirects: /testcases, /documents, etc. → new routes
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
      // ============================================
      // Public Routes (Authentication)
      // ============================================
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

      // ============================================
      // Global Routes (Not Project-Scoped)
      // ============================================
      {
        path: '/drafts',
        lazy: () =>
          import('../features/drafts/components/DraftListPage').then((m) => ({
            Component: () => (
              <AuthErrorBoundary>
                <RouteGuard>
                  <AppLayout>
                    <m.DraftListPage />
                  </AppLayout>
                </RouteGuard>
              </AuthErrorBoundary>
            ),
          })),
      },
      {
        path: '/drafts/:draftId',
        lazy: () =>
          import('../features/drafts/components/DraftConfirmPage').then((m) => ({
            Component: () => (
              <AuthErrorBoundary>
                <RouteGuard>
                  <AppLayout>
                    <m.DraftConfirmPage />
                  </AppLayout>
                </RouteGuard>
              </AuthErrorBoundary>
            ),
          })),
      },

      // ============================================
      // Legacy Route Redirects (Backward Compatibility)
      // ============================================
      {
        path: '/testcases',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <LegacyRouteRedirect to={(pid) => `/projects/${pid}/cases`} />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
      },
      {
        path: '/testcases/:caseId',
        lazy: () =>
          import('../features/testcases/components/CaseListPage').then(() => ({
            Component: () => (
              <AuthErrorBoundary>
                <RouteGuard>
                  <LegacyRouteRedirectWithParams
                    to={(pid, caseId) => `/projects/${pid}/cases/${caseId}`}
                    paramKey="caseId"
                  />
                </RouteGuard>
              </AuthErrorBoundary>
            ),
          })),
      },
      {
        path: '/documents',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <LegacyRouteRedirect to={(pid) => `/projects/${pid}/knowledge`} />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
      },
      {
        path: '/documents/:documentId',
        lazy: () =>
          import('../features/documents/components/KnowledgeListPage').then(() => ({
            Component: () => (
              <AuthErrorBoundary>
                <RouteGuard>
                  <LegacyRouteRedirectWithParams
                    to={(pid, docId) => `/projects/${pid}/knowledge/${docId}`}
                    paramKey="documentId"
                  />
                </RouteGuard>
              </AuthErrorBoundary>
            ),
          })),
      },
      {
        path: '/generation',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <LegacyRouteRedirect to={(pid) => `/projects/${pid}/generation`} />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
      },
      {
        path: '/generation/tasks/new',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <LegacyRouteRedirect to={(pid) => `/projects/${pid}/generation/new`} />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
      },
      {
        path: '/generation/tasks/:taskId',
        lazy: () =>
          import('../features/generation/components/GenerationTaskListPage').then(() => ({
            Component: () => (
              <AuthErrorBoundary>
                <RouteGuard>
                  <LegacyRouteRedirectWithParams
                    to={(pid, taskId) => `/projects/${pid}/generation/${taskId}`}
                    paramKey="taskId"
                  />
                </RouteGuard>
              </AuthErrorBoundary>
            ),
          })),
      },
      {
        path: '/plans',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <LegacyRouteRedirect to={(pid) => `/projects/${pid}/plans`} />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
      },
      {
        path: '/plans/new',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <LegacyRouteRedirect to={(pid) => `/projects/${pid}/plans/new`} />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
      },
      {
        path: '/plans/:planId',
        lazy: () =>
          import('../features/plans/components/PlanListPage').then(() => ({
            Component: () => (
              <AuthErrorBoundary>
                <RouteGuard>
                  <LegacyRouteRedirectWithParams
                    to={(pid, planId) => `/projects/${pid}/plans/${planId}`}
                    paramKey="planId"
                  />
                </RouteGuard>
              </AuthErrorBoundary>
            ),
          })),
      },
      {
        path: '/projects/:projectId/modules',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <Navigate to="./settings/modules" replace />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
      },

      // ============================================
      // Project-Scoped Routes (New Structure)
      // ============================================
      {
        path: '/projects',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <AppLayout />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
        children: [
          // Project list (index)
          {
            index: true,
            lazy: () =>
              import('../features/projects/components/ProjectListPage').then((m) => ({
                Component: m.ProjectListPage,
              })),
          },

          // Project-scoped routes
          {
            path: ':projectId',
            children: [
              // Dashboard redirect
              {
                index: true,
                element: <Navigate to="dashboard" replace />,
              },

              // Dashboard
              {
                path: 'dashboard',
                lazy: () =>
                  import('../features/projects/components/ProjectDashboard').then((m) => ({
                    Component: () => {
                      const { projectId } = useParams<{ projectId: string }>()
                      return projectId ? <m.ProjectDashboard projectId={projectId} /> : null
                    },
                  })),
              },

              // Knowledge Base
              {
                path: 'knowledge',
                children: [
                  {
                    index: true,
                    lazy: () =>
                      import('../features/documents/components/KnowledgeListPage').then(
                        (m) => ({
                          Component: m.KnowledgeListPage,
                        })
                      ),
                  },
                  {
                    path: ':docId',
                    lazy: () =>
                      import('../features/documents/components/DocumentDetailPage').then(
                        (m) => ({
                          Component: m.DocumentDetailPage,
                        })
                      ),
                  },
                  {
                    path: 'figma',
                    lazy: () =>
                      import('../features/documents/components/FigmaIntegrationPage').then((m) => ({
                        Component: () => (
                          <RouteGuard requireAdmin>
                            <m.default />
                          </RouteGuard>
                        ),
                      })),
                  },
                ],
              },

              // AI Generation
              {
                path: 'generation',
                children: [
                  {
                    index: true,
                    lazy: () =>
                      import('../features/generation/components/GenerationTaskListPage').then(
                        (m) => ({
                          Component: m.GenerationTaskListPage,
                        })
                      ),
                  },
                  {
                    path: 'new',
                    lazy: () =>
                      import('../features/generation/components/NewGenerationTaskPage').then(
                        (m) => ({
                          Component: m.NewGenerationTaskPage,
                        })
                      ),
                  },
                  {
                    path: ':taskId',
                    lazy: () =>
                      import('../features/generation/components/TaskDetailPage').then((m) => ({
                        Component: m.TaskDetailPage,
                      })),
                  },
                ],
              },

              // Test Cases
              {
                path: 'cases',
                children: [
                  {
                    index: true,
                    lazy: () =>
                      import('../features/testcases/components/CaseListPage').then((m) => ({
                        Component: m.CaseListPage,
                      })),
                  },
                  {
                    path: ':caseId',
                    lazy: () =>
                      import('../features/testcases/components/CaseDetailPage').then((m) => ({
                        Component: m.CaseDetailPage,
                      })),
                  },
                ],
              },

              // Test Plans
              {
                path: 'plans',
                children: [
                  {
                    index: true,
                    lazy: () =>
                      import('../features/plans/components/PlanListPage').then((m) => ({
                        Component: m.PlanListPage,
                      })),
                  },
                  {
                    path: 'new',
                    lazy: () =>
                      import('../features/plans/components/NewPlanPage').then((m) => ({
                        Component: m.NewPlanPage,
                      })),
                  },
                  {
                    path: ':planId',
                    lazy: () =>
                      import('../features/plans/components/PlanDetailPage').then((m) => ({
                        Component: m.PlanDetailPage,
                      })),
                  },
                ],
              },

              // Settings (Admin Only)
              {
                path: 'settings/modules',
                lazy: () =>
                  import('../features/modules/components/ModuleManagePage').then((m) => ({
                    Component: () => (
                      <RouteGuard requireAdmin>
                        <m.ModuleManagePage />
                      </RouteGuard>
                    ),
                  })),
              },
              {
                path: 'configs',
                lazy: () =>
                  import('../features/configs/components/ConfigManagePage').then((m) => ({
                    Component: () => (
                      <RouteGuard requireAdmin>
                        <m.ConfigManagePage />
                      </RouteGuard>
                    ),
                  })),
              },
            ],
          },
        ],
      },

      // ============================================
      // 404 Fallback
      // ============================================
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
