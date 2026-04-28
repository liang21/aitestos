/**
 * Generation module constants
 * Centralized configuration for AI generation features
 */

/**
 * Polling configuration
 */
export const POLLING = {
  /**
   * Interval in milliseconds for polling task status
   */
  INTERVAL_MS: 3000,

  /**
   * Task statuses that trigger active polling
   */
  ACTIVE_STATUSES: ['pending', 'processing'] as const,

  /**
   * Task statuses that stop polling
   */
  INACTIVE_STATUSES: ['completed', 'failed'] as const,
} as const

/**
 * Knowledge base readiness thresholds
 */
export const KNOWLEDGE_READINESS = {
  /**
   * Minimum number of completed documents for "sufficient" status
   */
  MIN_DOCUMENTS_FOR_SUFFICIENT: 2,

  /**
   * Threshold for "insufficient" status (less than this)
   */
  MIN_DOCUMENTS_FOR_INSUFFICIENT: 1,

  /**
   * Document status to count towards readiness
   */
  RELEVANT_DOCUMENT_STATUS: 'completed',
} as const

/**
 * Task list configuration
 */
export const TASK_LIST = {
  /**
   * Default page size for task list
   */
  DEFAULT_PAGE_SIZE: 20,

  /**
   * Available page size options
   */
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],

  /**
   * Query stale time in milliseconds (5 minutes)
   */
  STALE_TIME_MS: 5 * 60 * 1000,

  /**
   * Query retry count
   */
  RETRY_COUNT: 1,
} as const

/**
 * Task creation configuration
 */
export const TASK_CREATION = {
  /**
   * Minimum prompt length in characters
   */
  MIN_PROMPT_LENGTH: 10,

  /**
   * Maximum prompt length in characters
   */
  MAX_PROMPT_LENGTH: 500,

  /**
   * Default case count when not specified
   */
  DEFAULT_CASE_COUNT: 5,

  /**
   * Minimum case count
   */
  MIN_CASE_COUNT: 1,

  /**
   * Maximum case count
   */
  MAX_CASE_COUNT: 20,

  /**
   * Default priority when knowledge is insufficient
   */
  DEFAULT_PRIORITY_INSUFFICIENT: 'P2',
} as const

/**
 * UI configuration
 */
export const UI = {
  /**
   * Maximum characters to show before truncating prompt in list
   */
  PROMPT_TRUNCATE_LENGTH: 50,

  /**
   * Truncation suffix
   */
  TRUNCATE_SUFFIX: '...',
} as const

/**
 * Combined configuration object
 */
export const GENERATION_CONFIG = {
  POLLING,
  KNOWLEDGE_READINESS,
  TASK_LIST,
  TASK_CREATION,
  UI,
} as const

/**
 * Type exports
 */
export type KnowledgeReadiness = 'sufficient' | 'insufficient' | 'empty'
