// Package ierrors defines unified error codes
package ierrors

// Error code ranges:
// - 1xxxx: Identity Context
// - 2xxxx: Project Context
// - 3xxxx: Knowledge Context
// - 4xxxx: TestCase Context
// - 5xxxx: TestPlan Context
// - 6xxxx: Generation Context
// - 9xxxx: System Errors

const (
	// ============ Identity Context (1xxxx) ============
	CodeUserNotFound      = 10001
	CodeEmailDuplicate    = 10002
	CodeUsernameDuplicate = 10003
	CodePasswordMismatch  = 10004
	CodeInvalidEmail      = 10005
	CodePermissionDenied  = 10006

	// ============ Project Context (2xxxx) ============
	CodeProjectNotFound        = 20001
	CodeProjectNameDuplicate   = 20002
	CodeModuleNotFound         = 20003
	CodeModuleNameDuplicate    = 20004
	CodeConfigNotFound         = 20005
	CodeConfigKeyDuplicate     = 20006
	CodeProjectPrefixDuplicate = 20007
	CodeInvalidProjectPrefix   = 20008
	CodeModuleAbbrevDuplicate  = 20009
	CodeInvalidModuleAbbrev    = 20010

	// ============ Knowledge Context (3xxxx) ============
	CodeDocumentNotFound        = 30001
	CodeDocumentParseFailed     = 30002
	CodeUnsupportedDocumentType = 30003
	CodeEmptyChunks             = 30004
	CodeEmbeddingFailed         = 30005
	CodeVectorSearchFailed      = 30006
	CodeKnowledgeBaseEmpty      = 30007
	CodeDocumentProcessing      = 30008

	// ============ TestCase Context (4xxxx) ============
	CodeCaseNotFound        = 40001
	CodeCaseNumberDuplicate = 40002
	CodeInvalidCaseNumber   = 40003
	CodeEmptySteps          = 40004
	CodeInvalidPriority     = 40005
	CodeInvalidCaseType     = 40006

	// ============ TestPlan Context (5xxxx) ============
	CodePlanNotFound       = 50001
	CodePlanNameDuplicate  = 50002
	CodePlanArchived       = 50003
	CodeResultNotFound     = 50004
	CodeCaseNotInPlan      = 50005
	CodeDuplicateExecution = 50006

	// ============ Generation Context (6xxxx) ============
	CodeTaskNotFound           = 60001
	CodeTaskAlreadyProcessed   = 60002
	CodeDraftNotFound          = 60003
	CodeDraftAlreadyConfirmed  = 60004
	CodeDraftAlreadyRejected   = 60005
	CodeInvalidDraftStatus     = 60006
	CodeLLMCallFailed          = 60007
	CodeRAGNoResult            = 60008
	CodeConcurrentModification = 60009
	CodeLLMTimeout             = 60010
	CodeGenerationQueueFull    = 60011

	// ============ System Errors (9xxxx) ============
	CodeInternalError   = 90001
	CodeDatabaseError   = 90002
	CodeValidationError = 90003
	CodeUnauthorized    = 90004
	CodeRateLimited     = 90005
)
