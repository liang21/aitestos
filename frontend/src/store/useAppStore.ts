import { create } from 'zustand'
import { useEffect } from 'react'

interface AppState {
  sidebarCollapsed: boolean
  toggleSidebar: () => void
}

export const useAppStore = create<AppState>((set, get) => ({
  // Initial state: collapsed on small screens, expanded otherwise
  // Uses CSS variable breakpoint for consistency
  sidebarCollapsed: typeof window !== 'undefined'
    ? window.innerWidth < parseInt(getComputedStyle(document.documentElement).getPropertyValue('--responsive-breakpoint') || '768', 10)
    : false,

  toggleSidebar: () =>
    set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
}))

// Hook to initialize sidebar state based on window width
export function useSidebarResponsive() {
  useEffect(() => {
    const getBreakpoint = () => {
      const cssVar = getComputedStyle(document.documentElement).getPropertyValue('--responsive-breakpoint')
      return parseInt(cssVar || '768', 10)
    }

    const handleResize = () => {
      const breakpoint = getBreakpoint()
      const shouldBeCollapsed = window.innerWidth < breakpoint
      const currentState = useAppStore.getState().sidebarCollapsed
      if (shouldBeCollapsed !== currentState) {
        useAppStore.setState({ sidebarCollapsed: shouldBeCollapsed })
      }
    }

    // Initial check
    handleResize()

    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, []) // Empty deps - only register once

  return useAppStore()
}
