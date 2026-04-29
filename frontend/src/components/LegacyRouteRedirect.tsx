/**
 * Legacy Route Redirect Component
 *
 * Handles redirects from old route patterns to new project-scoped routes.
 * Extracts projectId from the most recent project in localStorage or redirects to project list.
 */

import { useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

interface LegacyRouteRedirectProps {
  to: (projectId: string) => string
}

interface LegacyRouteRedirectWithParamsProps {
  to: (projectId: string, param: string) => string
  paramKey: string
}

/**
 * Redirects legacy routes to new project-scoped routes.
 * Uses the most recently accessed project from localStorage, or redirects to /projects.
 */
export function LegacyRouteRedirect({ to }: LegacyRouteRedirectProps) {
  const navigate = useNavigate()

  useEffect(() => {
    // Try to get the most recent project from localStorage
    const recentProjectId = localStorage.getItem('lastAccessedProjectId')

    if (recentProjectId) {
      navigate(to(recentProjectId), { replace: true })
    } else {
      // No recent project, redirect to project list
      navigate('/projects', { replace: true })
    }
  }, [navigate, to])

  // Return null while redirecting
  return null
}

/**
 * Redirects legacy routes with URL params (e.g., /testcases/:caseId)
 * Extracts both projectId and the route param for the redirect.
 */
export function LegacyRouteRedirectWithParams({
  to,
  paramKey,
}: LegacyRouteRedirectWithParamsProps) {
  const navigate = useNavigate()
  const params = useParams<{ [key: string]: string }>()
  const paramValue = params[paramKey]

  useEffect(() => {
    const recentProjectId = localStorage.getItem('lastAccessedProjectId')

    if (recentProjectId && paramValue) {
      navigate(to(recentProjectId, paramValue), { replace: true })
    } else {
      navigate('/projects', { replace: true })
    }
  }, [navigate, to, paramValue])

  return null
}
