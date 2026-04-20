import type { FieldError } from 'react-hook-form'

/**
 * Convert react-hook-form error to Arco Design validateStatus
 */
export function useFormError(error: FieldError | undefined): 'error' | undefined {
  return error ? 'error' : undefined
}
