// Package ierrors provides error code mapping from domain errors
package ierrors

import (
	"errors"

	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/identity"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	"github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/domain/testcase"
	"github.com/liang21/aitestos/internal/domain/testplan"
)

// MapError converts domain errors to unified error codes
func MapError(err error) int {
	if err == nil {
		return 0
	}

	// Identity Context
	switch {
	case errors.Is(err, identity.ErrUserNotFound):
		return CodeUserNotFound
	case errors.Is(err, identity.ErrEmailDuplicate):
		return CodeEmailDuplicate
	case errors.Is(err, identity.ErrUsernameDuplicate):
		return CodeUsernameDuplicate
	case errors.Is(err, identity.ErrPasswordMismatch):
		return CodePasswordMismatch
	case errors.Is(err, identity.ErrInvalidEmail):
		return CodeInvalidEmail
	case errors.Is(err, identity.ErrPermissionDenied):
		return CodePermissionDenied
	}

	// Project Context
	switch {
	case errors.Is(err, project.ErrProjectNotFound):
		return CodeProjectNotFound
	case errors.Is(err, project.ErrProjectNameDuplicate):
		return CodeProjectNameDuplicate
	case errors.Is(err, project.ErrProjectPrefixDuplicate):
		return CodeProjectPrefixDuplicate
	case errors.Is(err, project.ErrInvalidProjectPrefix):
		return CodeInvalidProjectPrefix
	case errors.Is(err, project.ErrModuleNotFound):
		return CodeModuleNotFound
	case errors.Is(err, project.ErrModuleNameDuplicate):
		return CodeModuleNameDuplicate
	case errors.Is(err, project.ErrModuleAbbrevDuplicate):
		return CodeModuleAbbrevDuplicate
	case errors.Is(err, project.ErrInvalidModuleAbbrev):
		return CodeInvalidModuleAbbrev
	case errors.Is(err, project.ErrConfigNotFound):
		return CodeConfigNotFound
	case errors.Is(err, project.ErrConfigKeyDuplicate):
		return CodeConfigKeyDuplicate
	}

	// Knowledge Context
	switch {
	case errors.Is(err, knowledge.ErrDocumentNotFound):
		return CodeDocumentNotFound
	case errors.Is(err, knowledge.ErrDocumentParseFailed):
		return CodeDocumentParseFailed
	case errors.Is(err, knowledge.ErrUnsupportedDocumentType):
		return CodeUnsupportedDocumentType
	case errors.Is(err, knowledge.ErrEmptyChunks):
		return CodeEmptyChunks
	case errors.Is(err, knowledge.ErrEmbeddingFailed):
		return CodeEmbeddingFailed
	case errors.Is(err, knowledge.ErrVectorSearchFailed):
		return CodeVectorSearchFailed
	case errors.Is(err, knowledge.ErrKnowledgeBaseEmpty):
		return CodeKnowledgeBaseEmpty
	case errors.Is(err, knowledge.ErrDocumentProcessing):
		return CodeDocumentProcessing
	}

	// TestCase Context
	switch {
	case errors.Is(err, testcase.ErrCaseNotFound):
		return CodeCaseNotFound
	case errors.Is(err, testcase.ErrCaseNumberDuplicate):
		return CodeCaseNumberDuplicate
	case errors.Is(err, testcase.ErrInvalidCaseNumber):
		return CodeInvalidCaseNumber
	case errors.Is(err, testcase.ErrEmptySteps):
		return CodeEmptySteps
	case errors.Is(err, testcase.ErrInvalidPriority):
		return CodeInvalidPriority
	}

	// TestPlan Context
	switch {
	case errors.Is(err, testplan.ErrPlanNotFound):
		return CodePlanNotFound
	case errors.Is(err, testplan.ErrPlanNameDuplicate):
		return CodePlanNameDuplicate
	case errors.Is(err, testplan.ErrPlanArchived):
		return CodePlanArchived
	case errors.Is(err, testplan.ErrResultNotFound):
		return CodeResultNotFound
	case errors.Is(err, testplan.ErrCaseNotInPlan):
		return CodeCaseNotInPlan
	case errors.Is(err, testplan.ErrDuplicateExecution):
		return CodeDuplicateExecution
	}

	// Generation Context
	switch {
	case errors.Is(err, generation.ErrTaskNotFound):
		return CodeTaskNotFound
	case errors.Is(err, generation.ErrTaskAlreadyProcessed):
		return CodeTaskAlreadyProcessed
	case errors.Is(err, generation.ErrDraftNotFound):
		return CodeDraftNotFound
	case errors.Is(err, generation.ErrDraftAlreadyConfirmed):
		return CodeDraftAlreadyConfirmed
	case errors.Is(err, generation.ErrDraftAlreadyRejected):
		return CodeDraftAlreadyRejected
	case errors.Is(err, generation.ErrInvalidDraftStatus):
		return CodeInvalidDraftStatus
	case errors.Is(err, generation.ErrLLMCallFailed):
		return CodeLLMCallFailed
	case errors.Is(err, generation.ErrRAGNoResult):
		return CodeRAGNoResult
	case errors.Is(err, generation.ErrConcurrentModification):
		return CodeConcurrentModification
	case errors.Is(err, generation.ErrLLMTimeout):
		return CodeLLMTimeout
	case errors.Is(err, generation.ErrGenerationQueueFull):
		return CodeGenerationQueueFull
	}

	// Unknown error
	return CodeInternalError
}
