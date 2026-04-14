// Package ierrors_test tests error code definitions
package ierrors_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/ierrors"
)

func TestErrorCodeRanges(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		minRange int
		maxRange int
	}{
		// Identity Context (1xxxx)
		{"UserNotFound in Identity range", ierrors.CodeUserNotFound, 10000, 19999},
		{"EmailDuplicate in Identity range", ierrors.CodeEmailDuplicate, 10000, 19999},
		{"UsernameDuplicate in Identity range", ierrors.CodeUsernameDuplicate, 10000, 19999},
		{"PasswordMismatch in Identity range", ierrors.CodePasswordMismatch, 10000, 19999},
		{"InvalidEmail in Identity range", ierrors.CodeInvalidEmail, 10000, 19999},
		{"PermissionDenied in Identity range", ierrors.CodePermissionDenied, 10000, 19999},

		// Project Context (2xxxx)
		{"ProjectNotFound in Project range", ierrors.CodeProjectNotFound, 20000, 29999},
		{"ProjectNameDuplicate in Project range", ierrors.CodeProjectNameDuplicate, 20000, 29999},
		{"ModuleNotFound in Project range", ierrors.CodeModuleNotFound, 20000, 29999},
		{"ModuleNameDuplicate in Project range", ierrors.CodeModuleNameDuplicate, 20000, 29999},
		{"ConfigNotFound in Project range", ierrors.CodeConfigNotFound, 20000, 29999},
		{"ConfigKeyDuplicate in Project range", ierrors.CodeConfigKeyDuplicate, 20000, 29999},
		{"ProjectPrefixDuplicate in Project range", ierrors.CodeProjectPrefixDuplicate, 20000, 29999},
		{"InvalidProjectPrefix in Project range", ierrors.CodeInvalidProjectPrefix, 20000, 29999},
		{"ModuleAbbrevDuplicate in Project range", ierrors.CodeModuleAbbrevDuplicate, 20000, 29999},
		{"InvalidModuleAbbrev in Project range", ierrors.CodeInvalidModuleAbbrev, 20000, 29999},

		// Knowledge Context (3xxxx)
		{"DocumentNotFound in Knowledge range", ierrors.CodeDocumentNotFound, 30000, 39999},
		{"DocumentParseFailed in Knowledge range", ierrors.CodeDocumentParseFailed, 30000, 39999},
		{"UnsupportedDocumentType in Knowledge range", ierrors.CodeUnsupportedDocumentType, 30000, 39999},
		{"EmptyChunks in Knowledge range", ierrors.CodeEmptyChunks, 30000, 39999},
		{"EmbeddingFailed in Knowledge range", ierrors.CodeEmbeddingFailed, 30000, 39999},
		{"VectorSearchFailed in Knowledge range", ierrors.CodeVectorSearchFailed, 30000, 39999},
		{"KnowledgeBaseEmpty in Knowledge range", ierrors.CodeKnowledgeBaseEmpty, 30000, 39999},
		{"DocumentProcessing in Knowledge range", ierrors.CodeDocumentProcessing, 30000, 39999},

		// TestCase Context (4xxxx)
		{"CaseNotFound in TestCase range", ierrors.CodeCaseNotFound, 40000, 49999},
		{"CaseNumberDuplicate in TestCase range", ierrors.CodeCaseNumberDuplicate, 40000, 49999},
		{"InvalidCaseNumber in TestCase range", ierrors.CodeInvalidCaseNumber, 40000, 49999},
		{"EmptySteps in TestCase range", ierrors.CodeEmptySteps, 40000, 49999},
		{"InvalidPriority in TestCase range", ierrors.CodeInvalidPriority, 40000, 49999},
		{"InvalidCaseType in TestCase range", ierrors.CodeInvalidCaseType, 40000, 49999},

		// TestPlan Context (5xxxx)
		{"PlanNotFound in TestPlan range", ierrors.CodePlanNotFound, 50000, 59999},
		{"PlanNameDuplicate in TestPlan range", ierrors.CodePlanNameDuplicate, 50000, 59999},
		{"PlanArchived in TestPlan range", ierrors.CodePlanArchived, 50000, 59999},
		{"ResultNotFound in TestPlan range", ierrors.CodeResultNotFound, 50000, 59999},
		{"CaseNotInPlan in TestPlan range", ierrors.CodeCaseNotInPlan, 50000, 59999},
		{"DuplicateExecution in TestPlan range", ierrors.CodeDuplicateExecution, 50000, 59999},

		// Generation Context (6xxxx)
		{"TaskNotFound in Generation range", ierrors.CodeTaskNotFound, 60000, 69999},
		{"TaskAlreadyProcessed in Generation range", ierrors.CodeTaskAlreadyProcessed, 60000, 69999},
		{"DraftNotFound in Generation range", ierrors.CodeDraftNotFound, 60000, 69999},
		{"DraftAlreadyConfirmed in Generation range", ierrors.CodeDraftAlreadyConfirmed, 60000, 69999},
		{"DraftAlreadyRejected in Generation range", ierrors.CodeDraftAlreadyRejected, 60000, 69999},
		{"InvalidDraftStatus in Generation range", ierrors.CodeInvalidDraftStatus, 60000, 69999},
		{"LLMCallFailed in Generation range", ierrors.CodeLLMCallFailed, 60000, 69999},
		{"RAGNoResult in Generation range", ierrors.CodeRAGNoResult, 60000, 69999},
		{"ConcurrentModification in Generation range", ierrors.CodeConcurrentModification, 60000, 69999},
		{"LLMTimeout in Generation range", ierrors.CodeLLMTimeout, 60000, 69999},
		{"GenerationQueueFull in Generation range", ierrors.CodeGenerationQueueFull, 60000, 69999},

		// System Errors (9xxxx)
		{"InternalError in System range", ierrors.CodeInternalError, 90000, 99999},
		{"DatabaseError in System range", ierrors.CodeDatabaseError, 90000, 99999},
		{"ValidationError in System range", ierrors.CodeValidationError, 90000, 99999},
		{"Unauthorized in System range", ierrors.CodeUnauthorized, 90000, 99999},
		{"RateLimited in System range", ierrors.CodeRateLimited, 90000, 99999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code < tt.minRange || tt.code > tt.maxRange {
				t.Errorf("code %d not in expected range [%d, %d]", tt.code, tt.minRange, tt.maxRange)
			}
		})
	}
}

func TestErrorCodeUniqueness(t *testing.T) {
	codes := map[int]string{
		// Identity
		ierrors.CodeUserNotFound:      "CodeUserNotFound",
		ierrors.CodeEmailDuplicate:    "CodeEmailDuplicate",
		ierrors.CodeUsernameDuplicate: "CodeUsernameDuplicate",
		ierrors.CodePasswordMismatch:  "CodePasswordMismatch",
		ierrors.CodeInvalidEmail:      "CodeInvalidEmail",
		ierrors.CodePermissionDenied:  "CodePermissionDenied",

		// Project
		ierrors.CodeProjectNotFound:        "CodeProjectNotFound",
		ierrors.CodeProjectNameDuplicate:   "CodeProjectNameDuplicate",
		ierrors.CodeModuleNotFound:         "CodeModuleNotFound",
		ierrors.CodeModuleNameDuplicate:    "CodeModuleNameDuplicate",
		ierrors.CodeConfigNotFound:         "CodeConfigNotFound",
		ierrors.CodeConfigKeyDuplicate:     "CodeConfigKeyDuplicate",
		ierrors.CodeProjectPrefixDuplicate: "CodeProjectPrefixDuplicate",
		ierrors.CodeInvalidProjectPrefix:   "CodeInvalidProjectPrefix",
		ierrors.CodeModuleAbbrevDuplicate:  "CodeModuleAbbrevDuplicate",
		ierrors.CodeInvalidModuleAbbrev:    "CodeInvalidModuleAbbrev",

		// Knowledge
		ierrors.CodeDocumentNotFound:        "CodeDocumentNotFound",
		ierrors.CodeDocumentParseFailed:     "CodeDocumentParseFailed",
		ierrors.CodeUnsupportedDocumentType: "CodeUnsupportedDocumentType",
		ierrors.CodeEmptyChunks:             "CodeEmptyChunks",
		ierrors.CodeEmbeddingFailed:         "CodeEmbeddingFailed",
		ierrors.CodeVectorSearchFailed:      "CodeVectorSearchFailed",
		ierrors.CodeKnowledgeBaseEmpty:      "CodeKnowledgeBaseEmpty",
		ierrors.CodeDocumentProcessing:      "CodeDocumentProcessing",

		// TestCase
		ierrors.CodeCaseNotFound:        "CodeCaseNotFound",
		ierrors.CodeCaseNumberDuplicate: "CodeCaseNumberDuplicate",
		ierrors.CodeInvalidCaseNumber:   "CodeInvalidCaseNumber",
		ierrors.CodeEmptySteps:          "CodeEmptySteps",
		ierrors.CodeInvalidPriority:     "CodeInvalidPriority",
		ierrors.CodeInvalidCaseType:     "CodeInvalidCaseType",

		// TestPlan
		ierrors.CodePlanNotFound:       "CodePlanNotFound",
		ierrors.CodePlanNameDuplicate:  "CodePlanNameDuplicate",
		ierrors.CodePlanArchived:       "CodePlanArchived",
		ierrors.CodeResultNotFound:     "CodeResultNotFound",
		ierrors.CodeCaseNotInPlan:      "CodeCaseNotInPlan",
		ierrors.CodeDuplicateExecution: "CodeDuplicateExecution",

		// Generation
		ierrors.CodeTaskNotFound:           "CodeTaskNotFound",
		ierrors.CodeTaskAlreadyProcessed:   "CodeTaskAlreadyProcessed",
		ierrors.CodeDraftNotFound:          "CodeDraftNotFound",
		ierrors.CodeDraftAlreadyConfirmed:  "CodeDraftAlreadyConfirmed",
		ierrors.CodeDraftAlreadyRejected:   "CodeDraftAlreadyRejected",
		ierrors.CodeInvalidDraftStatus:     "CodeInvalidDraftStatus",
		ierrors.CodeLLMCallFailed:          "CodeLLMCallFailed",
		ierrors.CodeRAGNoResult:            "CodeRAGNoResult",
		ierrors.CodeConcurrentModification: "CodeConcurrentModification",
		ierrors.CodeLLMTimeout:             "CodeLLMTimeout",
		ierrors.CodeGenerationQueueFull:    "CodeGenerationQueueFull",

		// System
		ierrors.CodeInternalError:   "CodeInternalError",
		ierrors.CodeDatabaseError:   "CodeDatabaseError",
		ierrors.CodeValidationError: "CodeValidationError",
		ierrors.CodeUnauthorized:    "CodeUnauthorized",
		ierrors.CodeRateLimited:     "CodeRateLimited",
	}

	// Check for duplicates
	seen := make(map[int]string)
	for code, name := range codes {
		if existing, exists := seen[code]; exists {
			t.Errorf("duplicate error code %d: %s and %s", code, existing, name)
		}
		seen[code] = name
	}
}

func TestErrorCodeCount(t *testing.T) {
	// Ensure we have expected number of error codes per context
	tests := []struct {
		name          string
		expectedCount int
		codes         []int
	}{
		{"Identity codes", 6, []int{
			ierrors.CodeUserNotFound,
			ierrors.CodeEmailDuplicate,
			ierrors.CodeUsernameDuplicate,
			ierrors.CodePasswordMismatch,
			ierrors.CodeInvalidEmail,
			ierrors.CodePermissionDenied,
		}},
		{"Project codes", 10, []int{
			ierrors.CodeProjectNotFound,
			ierrors.CodeProjectNameDuplicate,
			ierrors.CodeModuleNotFound,
			ierrors.CodeModuleNameDuplicate,
			ierrors.CodeConfigNotFound,
			ierrors.CodeConfigKeyDuplicate,
			ierrors.CodeProjectPrefixDuplicate,
			ierrors.CodeInvalidProjectPrefix,
			ierrors.CodeModuleAbbrevDuplicate,
			ierrors.CodeInvalidModuleAbbrev,
		}},
		{"Knowledge codes", 8, []int{
			ierrors.CodeDocumentNotFound,
			ierrors.CodeDocumentParseFailed,
			ierrors.CodeUnsupportedDocumentType,
			ierrors.CodeEmptyChunks,
			ierrors.CodeEmbeddingFailed,
			ierrors.CodeVectorSearchFailed,
			ierrors.CodeKnowledgeBaseEmpty,
			ierrors.CodeDocumentProcessing,
		}},
		{"TestCase codes", 6, []int{
			ierrors.CodeCaseNotFound,
			ierrors.CodeCaseNumberDuplicate,
			ierrors.CodeInvalidCaseNumber,
			ierrors.CodeEmptySteps,
			ierrors.CodeInvalidPriority,
			ierrors.CodeInvalidCaseType,
		}},
		{"TestPlan codes", 6, []int{
			ierrors.CodePlanNotFound,
			ierrors.CodePlanNameDuplicate,
			ierrors.CodePlanArchived,
			ierrors.CodeResultNotFound,
			ierrors.CodeCaseNotInPlan,
			ierrors.CodeDuplicateExecution,
		}},
		{"Generation codes", 11, []int{
			ierrors.CodeTaskNotFound,
			ierrors.CodeTaskAlreadyProcessed,
			ierrors.CodeDraftNotFound,
			ierrors.CodeDraftAlreadyConfirmed,
			ierrors.CodeDraftAlreadyRejected,
			ierrors.CodeInvalidDraftStatus,
			ierrors.CodeLLMCallFailed,
			ierrors.CodeRAGNoResult,
			ierrors.CodeConcurrentModification,
			ierrors.CodeLLMTimeout,
			ierrors.CodeGenerationQueueFull,
		}},
		{"System codes", 5, []int{
			ierrors.CodeInternalError,
			ierrors.CodeDatabaseError,
			ierrors.CodeValidationError,
			ierrors.CodeUnauthorized,
			ierrors.CodeRateLimited,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.codes) != tt.expectedCount {
				t.Errorf("%s: expected %d codes, got %d", tt.name, tt.expectedCount, len(tt.codes))
			}
		})
	}
}
