/**
 * useProjectRoutes Hook
 *
 * Provides project-scoped routes based on current URL params.
 * Must be used within a project route context (under /projects/:projectId).
 */

import { useParams } from 'react-router-dom'
import { buildProjectRoutes } from '@/lib/routes'

/**
 * Get project-scoped routes for the current project
 * @returns Project routes object
 * @throws Error if used outside of a project route context
 *
 * @example
 * ```tsx
 * function MyComponent() {
 *   const routes = useProjectRoutes()
 *   return <Link to={routes.cases.list}>Test Cases</Link>
 * }
 * ```
 */
export function useProjectRoutes() {
  const { projectId } = useParams<{ projectId: string }>()

  if (!projectId) {
    throw new Error(
      'useProjectRoutes must be used within a project route context. ' +
      'Ensure the component is rendered under a /projects/:projectId route.'
    )
  }

  return buildProjectRoutes(projectId)
}
