/**
 * Test case execution status
 */
export type CaseStatus = 'unexecuted' | 'pass' | 'block' | 'fail'

/**
 * Test case type/category
 */
export type CaseType =
  | 'functionality'
  | 'performance'
  | 'api'
  | 'ui'
  | 'security'

/**
 * Test plan status
 */
export type PlanStatus = 'draft' | 'active' | 'completed' | 'archived'

/**
 * Priority levels
 */
export type Priority = 'P0' | 'P1' | 'P2' | 'P3'

/**
 * Test result execution status
 */
export type ResultStatus = 'pass' | 'fail' | 'block' | 'skip'

/**
 * AI generation task status
 */
export type TaskStatus = 'pending' | 'processing' | 'completed' | 'failed'

/**
 * Draft confirmation status
 */
export type DraftStatus = 'pending' | 'confirmed' | 'rejected'

/**
 * Document type in knowledge base
 */
export type DocumentType = 'prd' | 'figma' | 'api_spec' | 'swagger' | 'markdown'

/**
 * Document processing status
 */
export type DocumentStatus = 'pending' | 'processing' | 'completed' | 'failed'

/**
 * User role for access control
 */
export type UserRole = 'super_admin' | 'admin' | 'normal'

/**
 * AI generation confidence level
 */
export type Confidence = 'high' | 'medium' | 'low'

/**
 * Test scene type
 */
export type SceneType = 'positive' | 'negative' | 'boundary'
