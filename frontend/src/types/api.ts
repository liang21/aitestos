import type {
  CaseStatus,
  CaseType,
  Confidence,
  DocumentStatus,
  DocumentType,
  PlanStatus,
  Priority,
  ResultStatus,
  SceneType,
  TaskStatus,
  UserRole,
  DraftStatus,
} from './enums'

// ============================================================================
// Common Types
// ============================================================================

/**
 * Standard paginated response wrapper
 */
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  offset: number
  limit: number
}

// ============================================================================
// Auth Types
// ============================================================================

/**
 * Login request payload
 */
export interface LoginRequest {
  email: string
  password: string
}

/**
 * Register request payload
 */
export interface RegisterRequest {
  username: string
  email: string
  password: string
  role: UserRole
}

/**
 * Login/Register response
 */
export interface AuthResponse {
  access_token: string
  refresh_token: string
  user: UserJSON
}

/**
 * Refresh token response
 */
export interface RefreshResponse {
  access_token: string
  refresh_token: string
}

// ============================================================================
// User Types
// ============================================================================

/**
 * User profile data
 */
export interface UserJSON {
  id: string
  username: string
  email: string
  role: UserRole
  createdAt: string
  updatedAt: string
}

// ============================================================================
// Project Types
// ============================================================================

/**
 * Project base fields
 */
export interface Project {
  id: string
  name: string
  prefix: string
  description: string
  createdAt: string
  updatedAt: string
}

/**
 * Create project request
 */
export interface CreateProjectRequest {
  name: string
  prefix: string
  description?: string
}

/**
 * Update project request
 */
export interface UpdateProjectRequest {
  name?: string
  description?: string
}

/**
 * Project detail with extended info
 */
export interface ProjectDetail extends Project {
  moduleCount: number
  caseCount: number
  draftCount: number
}

/**
 * Project statistics
 */
export interface ProjectStats {
  totalCases: number
  passRate: number
  coverage: number
  aiGeneratedCount: number
  trend: Array<{
    date: string
    passRate: number
  }>
}

// ============================================================================
// Module Types
// ============================================================================

/**
 * Module data
 */
export interface Module {
  id: string
  projectId: string
  name: string
  abbreviation: string
  createdAt: string
  updatedAt: string
  caseCount?: number
}

/**
 * Create module request
 */
export interface CreateModuleRequest {
  name: string
  abbreviation: string
}

// ============================================================================
// Test Case Types
// ============================================================================

/**
 * AI metadata for generated test cases
 */
export interface AiMetadata {
  generationTaskId: string
  confidence: Confidence
  referencedChunks: ReferencedChunk[]
  modelVersion: string
  generatedAt: string
}

/**
 * Referenced document chunk
 */
export interface ReferencedChunk {
  chunkId: string
  documentId: string
  documentTitle: string
  similarityScore: number
}

/**
 * Test case data
 */
export interface TestCase {
  id: string
  moduleId: string
  userId: string
  number: string // Format: {prefix}-{abbreviation}-{YYYYMMDD}-{001}
  title: string
  preconditions: string[]
  steps: string[]
  expected: Record<string, unknown>
  caseType: CaseType
  priority: Priority
  status: CaseStatus
  aiMetadata?: AiMetadata
  createdAt: string
  updatedAt: string
  // Joined fields
  moduleName?: string
  projectName?: string
  projectPrefix?: string
  createdByName?: string
}

/**
 * Create test case request
 */
export interface CreateTestCaseRequest {
  moduleId: string
  title: string
  preconditions: string[]
  steps: string[]
  expected: Record<string, unknown>
  caseType: CaseType
  priority: Priority
}

/**
 * Update test case request
 */
export interface UpdateTestCaseRequest {
  title?: string
  preconditions?: string[]
  steps?: string[]
  expected?: Record<string, unknown>
  caseType?: CaseType
  priority?: Priority
  status?: CaseStatus
}

/**
 * Test case list query parameters
 */
export interface TestCaseListParams {
  projectId?: string
  moduleId?: string
  status?: CaseStatus
  caseType?: CaseType
  priority?: Priority
  keywords?: string
  offset?: number
  limit?: number
}

// ============================================================================
// Draft Types
// ============================================================================

/**
 * Draft test case data
 */
export interface CaseDraft {
  id: string
  taskId: string
  projectId: string
  title: string
  preconditions: string[]
  steps: string[]
  expected: Record<string, unknown>
  caseType: CaseType
  priority: Priority
  status: DraftStatus
  feedback?: string
  aiMetadata?: {
    confidence: Confidence
    referencedChunks: ReferencedChunk[]
    modelVersion: string
  }
  createdAt: string
  updatedAt: string
  // Joined fields
  projectName?: string
  moduleName?: string
}

/**
 * Confirm draft request
 */
export interface ConfirmDraftRequest {
  moduleId: string
  title?: string
  preconditions?: string[]
  steps?: string[]
  expected?: Record<string, unknown>
  caseType?: CaseType
  priority?: Priority
}

/**
 * Reject draft request
 */
export interface RejectDraftRequest {
  reason: string
  feedback?: string
}

/**
 * Batch confirm request
 */
export interface BatchConfirmRequest {
  draftIds: string[]
  moduleId: string
}

/**
 * Batch confirm response
 */
export interface BatchConfirmResponse {
  successCount: number
  failedCount: number
  errors?: Array<{ draftId: string; error: string }>
}

/**
 * Draft list parameters
 */
export interface DraftListParams {
  projectId?: string
  status?: DraftStatus
  keywords?: string
  offset?: number
  limit?: number
}

// ============================================================================
// Test Plan Types
// ============================================================================

/**
 * Test plan data
 */
export interface TestPlan {
  id: string
  projectId: string
  name: string
  description: string
  status: PlanStatus
  createdBy: string
  createdAt: string
  updatedAt: string
  // Joined fields
  caseCount?: number
  createdByName?: string
}

/**
 * Create plan request
 */
export interface CreatePlanRequest {
  projectId: string
  name: string
  description?: string
}

/**
 * Update plan request
 */
export interface UpdatePlanRequest {
  name?: string
  description?: string
  status?: PlanStatus
}

/**
 * Plan detail with execution stats
 */
export interface PlanDetail extends TestPlan {
  cases: PlanCase[]
  stats: PlanStats
}

/**
 * Test case in plan
 */
export interface PlanCase {
  caseId: string
  caseNumber: string
  caseTitle: string
  resultStatus?: ResultStatus
  resultNote?: string
  executedAt?: string
  executedBy?: string
}

/**
 * Plan execution statistics
 */
export interface PlanStats {
  total: number
  passed: number
  failed: number
  blocked: number
  skipped: number
  unexecuted: number
}

/**
 * Record test result request
 */
export interface RecordResultRequest {
  caseId: string
  status: ResultStatus
  note?: string
}

// ============================================================================
// AI Generation Types
// ============================================================================

/**
 * Generation task data
 */
export interface GenerationTask {
  id: string
  projectId: string
  moduleId: string
  status: TaskStatus
  prompt: string
  result: GenerationTaskResult | null
  createdAt: string
  updatedAt: string
  // Joined fields
  projectName?: string
  moduleName?: string
  createdByName?: string
}

/**
 * Generation task result
 */
export interface GenerationTaskResult {
  draftCount: number
  confidence: Confidence
  completedAt?: string
  error?: string
}

/**
 * Create generation task request
 */
export interface CreateTaskRequest {
  projectId: string
  moduleId: string
  prompt: string
  count?: number
  caseType?: CaseType
  priority?: Priority
  sceneType?: SceneType
}

/**
 * Generation task list parameters
 */
export interface TaskListParams {
  projectId?: string
  status?: TaskStatus
  offset?: number
  limit?: number
}

// ============================================================================
// Document Types
// ============================================================================

/**
 * Document data
 */
export interface Document {
  id: string
  projectId: string
  name: string
  type: DocumentType
  status: DocumentStatus
  chunkCount: number
  uploadedBy: string
  createdAt: string
  updatedAt: string
  // Joined fields
  uploadedByName?: string
}

/**
 * Upload document request
 */
export interface UploadDocumentRequest {
  projectId: string
  name: string
  type: DocumentType
  file?: File
  url?: string
}

/**
 * Document detail with chunks
 */
export interface DocumentDetail extends Document {
  chunks: DocumentChunk[]
}

/**
 * Document chunk data
 */
export interface DocumentChunk {
  id: string
  documentId: string
  chunkIndex: number
  content: string
  metadata?: Record<string, unknown>
}

// ============================================================================
// Config Types
// ============================================================================

/**
 * Project configuration item
 */
export interface ConfigItem {
  key: string
  value: unknown
  description?: string
}

/**
 * Set config request
 */
export interface SetConfigRequest {
  key: string
  value: unknown
  description?: string
}

/**
 * Import configs request
 */
export interface ImportConfigsRequest {
  configs: Array<{ key: string; value: unknown }>
}

/**
 * Export configs response
 */
export interface ExportConfigsResponse {
  configs: Array<{ key: string; value: unknown }>
}
