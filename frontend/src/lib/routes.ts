/**
 * Application Route Constants and Factory Functions
 *
 * Provides centralized route management with type safety.
 * All routes should be imported from this file to ensure consistency.
 */

/**
 * Build project-scoped routes dynamically
 * @param projectId - The project ID from URL params
 * @returns Object containing all project-relative routes
 */
export function buildProjectRoutes(projectId: string) {
  const base = `/projects/${projectId}`

  return {
    /** Project dashboard */
    dashboard: `${base}/dashboard`,

    /** Knowledge base routes */
    knowledge: {
      /** Knowledge document list */
      list: `${base}/knowledge`,
      /** Knowledge document detail */
      detail: (docId: string) => `${base}/knowledge/${docId}`,
      /** Figma integration page (admin only) */
      figma: `${base}/knowledge/figma`,
    },

    /** AI generation routes */
    generation: {
      /** Generation task list */
      list: `${base}/generation`,
      /** Create new generation task */
      new: `${base}/generation/new`,
      /** Generation task detail */
      detail: (taskId: string) => `${base}/generation/${taskId}`,
    },

    /** Test case management routes */
    cases: {
      /** Test case list */
      list: `${base}/cases`,
      /** Test case detail */
      detail: (caseId: string) => `${base}/cases/${caseId}`,
    },

    /** Test plan management routes */
    plans: {
      /** Test plan list */
      list: `${base}/plans`,
      /** Create new test plan */
      new: `${base}/plans/new`,
      /** Test plan detail */
      detail: (planId: string) => `${base}/plans/${planId}`,
    },

    /** Project settings routes */
    settings: {
      /** Module management (admin only) */
      modules: `${base}/settings/modules`,
      /** Configuration management (admin only) */
      configs: `${base}/configs`,
    },
  } as const
}

/**
 * Global application routes (not project-scoped)
 */
export const GLOBAL_ROUTES = {
  /** Project list page */
  projects: '/projects',

  /** Draft box (global across projects) */
  drafts: '/drafts',
  /** Draft confirmation page */
  draftDetail: (draftId: string) => `/drafts/${draftId}`,

  /** Authentication */
  login: '/login',
  register: '/register',
} as const

/**
 * Legacy route mappings for backward compatibility
 * Maps old route patterns to new route patterns
 */
export const LEGACY_ROUTE_REDIRECTS = {
  // /testcases → /projects/:projectId/cases
  testcases: (projectId: string) => `/projects/${projectId}/cases`,

  // /documents → /projects/:projectId/knowledge
  documents: (projectId: string) => `/projects/${projectId}/knowledge`,

  // /generation → /projects/:projectId/generation
  generation: (projectId: string) => `/projects/${projectId}/generation`,

  // /generation/tasks/new → /projects/:projectId/generation/new
  'generation/tasks/new': (projectId: string) => `/projects/${projectId}/generation/new`,

  // /generation/tasks/:taskId → /projects/:projectId/generation/:taskId
  'generation/tasks/:taskId': (projectId: string, taskId: string) =>
    `/projects/${projectId}/generation/${taskId}`,

  // /plans → /projects/:projectId/plans
  plans: (projectId: string) => `/projects/${projectId}/plans`,

  // /plans/new → /projects/:projectId/plans/new
  'plans/new': (projectId: string) => `/projects/${projectId}/plans/new`,

  // /plans/:planId → /projects/:projectId/plans/:planId
  'plans/:planId': (projectId: string, planId: string) => `/projects/${projectId}/plans/${planId}`,

  // /projects/:projectId/modules → /projects/:projectId/settings/modules
  'projects/:projectId/modules': (projectId: string) => `/projects/${projectId}/settings/modules`,
} as const
