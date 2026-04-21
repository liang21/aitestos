/**
 * Drafts module constants
 */

/**
 * Polling interval for pending draft count (milliseconds)
 */
export const DRAFT_POLL_INTERVAL = 5000

/**
 * Maximum number of drafts that can be batch confirmed at once
 */
export const MAX_BATCH_CONFIRM_COUNT = 50

/**
 * Draft rejection reasons
 */
export const DRAFT_REJECTION_REASONS = [
  { label: '重复', value: 'duplicate' },
  { label: '无关', value: 'irrelevant' },
  { label: '低质量', value: 'low_quality' },
  { label: '其他', value: 'other' },
] as const

/**
 * Draft rejection reason type
 */
export type DraftRejectionReason =
  (typeof DRAFT_REJECTION_REASONS)[number]['value']
