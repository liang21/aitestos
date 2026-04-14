// Package ierrors_test tests error mapping functionality
package ierrors_test

import (
	"errors"
	"testing"

	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/identity"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	"github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/domain/testcase"
	"github.com/liang21/aitestos/internal/domain/testplan"
	"github.com/liang21/aitestos/internal/ierrors"
)

func TestMapError_IdentityContext(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{"ErrUserNotFound", identity.ErrUserNotFound, ierrors.CodeUserNotFound},
		{"ErrEmailDuplicate", identity.ErrEmailDuplicate, ierrors.CodeEmailDuplicate},
		{"ErrUsernameDuplicate", identity.ErrUsernameDuplicate, ierrors.CodeUsernameDuplicate},
		{"ErrPasswordMismatch", identity.ErrPasswordMismatch, ierrors.CodePasswordMismatch},
		{"ErrInvalidEmail", identity.ErrInvalidEmail, ierrors.CodeInvalidEmail},
		{"ErrPermissionDenied", identity.ErrPermissionDenied, ierrors.CodePermissionDenied},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ierrors.MapError(tt.err)
			if got != tt.expectedCode {
				t.Errorf("MapError(%v) = %d, want %d", tt.err, got, tt.expectedCode)
			}
		})
	}
}

func TestMapError_ProjectContext(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{"ErrProjectNotFound", project.ErrProjectNotFound, ierrors.CodeProjectNotFound},
		{"ErrProjectNameDuplicate", project.ErrProjectNameDuplicate, ierrors.CodeProjectNameDuplicate},
		{"ErrProjectPrefixDuplicate", project.ErrProjectPrefixDuplicate, ierrors.CodeProjectPrefixDuplicate},
		{"ErrInvalidProjectPrefix", project.ErrInvalidProjectPrefix, ierrors.CodeInvalidProjectPrefix},
		{"ErrModuleNotFound", project.ErrModuleNotFound, ierrors.CodeModuleNotFound},
		{"ErrModuleNameDuplicate", project.ErrModuleNameDuplicate, ierrors.CodeModuleNameDuplicate},
		{"ErrModuleAbbrevDuplicate", project.ErrModuleAbbrevDuplicate, ierrors.CodeModuleAbbrevDuplicate},
		{"ErrInvalidModuleAbbrev", project.ErrInvalidModuleAbbrev, ierrors.CodeInvalidModuleAbbrev},
		{"ErrConfigNotFound", project.ErrConfigNotFound, ierrors.CodeConfigNotFound},
		{"ErrConfigKeyDuplicate", project.ErrConfigKeyDuplicate, ierrors.CodeConfigKeyDuplicate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ierrors.MapError(tt.err)
			if got != tt.expectedCode {
				t.Errorf("MapError(%v) = %d, want %d", tt.err, got, tt.expectedCode)
			}
		})
	}
}

func TestMapError_KnowledgeContext(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{"ErrDocumentNotFound", knowledge.ErrDocumentNotFound, ierrors.CodeDocumentNotFound},
		{"ErrDocumentParseFailed", knowledge.ErrDocumentParseFailed, ierrors.CodeDocumentParseFailed},
		{"ErrUnsupportedDocumentType", knowledge.ErrUnsupportedDocumentType, ierrors.CodeUnsupportedDocumentType},
		{"ErrEmptyChunks", knowledge.ErrEmptyChunks, ierrors.CodeEmptyChunks},
		{"ErrEmbeddingFailed", knowledge.ErrEmbeddingFailed, ierrors.CodeEmbeddingFailed},
		{"ErrVectorSearchFailed", knowledge.ErrVectorSearchFailed, ierrors.CodeVectorSearchFailed},
		{"ErrKnowledgeBaseEmpty", knowledge.ErrKnowledgeBaseEmpty, ierrors.CodeKnowledgeBaseEmpty},
		{"ErrDocumentProcessing", knowledge.ErrDocumentProcessing, ierrors.CodeDocumentProcessing},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ierrors.MapError(tt.err)
			if got != tt.expectedCode {
				t.Errorf("MapError(%v) = %d, want %d", tt.err, got, tt.expectedCode)
			}
		})
	}
}

func TestMapError_TestCaseContext(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{"ErrCaseNotFound", testcase.ErrCaseNotFound, ierrors.CodeCaseNotFound},
		{"ErrCaseNumberDuplicate", testcase.ErrCaseNumberDuplicate, ierrors.CodeCaseNumberDuplicate},
		{"ErrInvalidCaseNumber", testcase.ErrInvalidCaseNumber, ierrors.CodeInvalidCaseNumber},
		{"ErrEmptySteps", testcase.ErrEmptySteps, ierrors.CodeEmptySteps},
		{"ErrInvalidPriority", testcase.ErrInvalidPriority, ierrors.CodeInvalidPriority},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ierrors.MapError(tt.err)
			if got != tt.expectedCode {
				t.Errorf("MapError(%v) = %d, want %d", tt.err, got, tt.expectedCode)
			}
		})
	}
}

func TestMapError_TestPlanContext(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{"ErrPlanNotFound", testplan.ErrPlanNotFound, ierrors.CodePlanNotFound},
		{"ErrPlanNameDuplicate", testplan.ErrPlanNameDuplicate, ierrors.CodePlanNameDuplicate},
		{"ErrPlanArchived", testplan.ErrPlanArchived, ierrors.CodePlanArchived},
		{"ErrResultNotFound", testplan.ErrResultNotFound, ierrors.CodeResultNotFound},
		{"ErrCaseNotInPlan", testplan.ErrCaseNotInPlan, ierrors.CodeCaseNotInPlan},
		{"ErrDuplicateExecution", testplan.ErrDuplicateExecution, ierrors.CodeDuplicateExecution},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ierrors.MapError(tt.err)
			if got != tt.expectedCode {
				t.Errorf("MapError(%v) = %d, want %d", tt.err, got, tt.expectedCode)
			}
		})
	}
}

func TestMapError_GenerationContext(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{"ErrTaskNotFound", generation.ErrTaskNotFound, ierrors.CodeTaskNotFound},
		{"ErrTaskAlreadyProcessed", generation.ErrTaskAlreadyProcessed, ierrors.CodeTaskAlreadyProcessed},
		{"ErrDraftNotFound", generation.ErrDraftNotFound, ierrors.CodeDraftNotFound},
		{"ErrDraftAlreadyConfirmed", generation.ErrDraftAlreadyConfirmed, ierrors.CodeDraftAlreadyConfirmed},
		{"ErrDraftAlreadyRejected", generation.ErrDraftAlreadyRejected, ierrors.CodeDraftAlreadyRejected},
		{"ErrInvalidDraftStatus", generation.ErrInvalidDraftStatus, ierrors.CodeInvalidDraftStatus},
		{"ErrLLMCallFailed", generation.ErrLLMCallFailed, ierrors.CodeLLMCallFailed},
		{"ErrRAGNoResult", generation.ErrRAGNoResult, ierrors.CodeRAGNoResult},
		{"ErrConcurrentModification", generation.ErrConcurrentModification, ierrors.CodeConcurrentModification},
		{"ErrLLMTimeout", generation.ErrLLMTimeout, ierrors.CodeLLMTimeout},
		{"ErrGenerationQueueFull", generation.ErrGenerationQueueFull, ierrors.CodeGenerationQueueFull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ierrors.MapError(tt.err)
			if got != tt.expectedCode {
				t.Errorf("MapError(%v) = %d, want %d", tt.err, got, tt.expectedCode)
			}
		})
	}
}

func TestMapError_NilError(t *testing.T) {
	got := ierrors.MapError(nil)
	if got != 0 {
		t.Errorf("MapError(nil) = %d, want 0", got)
	}
}

func TestMapError_UnknownError(t *testing.T) {
	unknownErr := errors.New("some unknown error")
	got := ierrors.MapError(unknownErr)
	if got != ierrors.CodeInternalError {
		t.Errorf("MapError(unknown) = %d, want %d", got, ierrors.CodeInternalError)
	}
}

func TestMapError_WrappedError(t *testing.T) {
	// Test that wrapped errors are correctly identified using errors.Is
	wrappedErr := errors.Join(identity.ErrUserNotFound, errors.New("additional context"))
	got := ierrors.MapError(wrappedErr)
	if got != ierrors.CodeUserNotFound {
		t.Errorf("MapError(wrapped) = %d, want %d", got, ierrors.CodeUserNotFound)
	}
}

func TestMapError_ErrorChain(t *testing.T) {
	// Test error chain with fmt.Errorf %w
	baseErr := project.ErrProjectNotFound
	wrappedErr := errors.Join(baseErr, errors.New("context"))
	got := ierrors.MapError(wrappedErr)
	if got != ierrors.CodeProjectNotFound {
		t.Errorf("MapError(wrapped chain) = %d, want %d", got, ierrors.CodeProjectNotFound)
	}
}
