// Package ierrors_test tests HTTP error response structures
package ierrors_test

import (
	"encoding/json"
	"testing"

	"github.com/liang21/aitestos/internal/ierrors"
)

func TestErrorResponse_ToJSON(t *testing.T) {
	tests := []struct {
		name     string
		response *ierrors.ErrorResponse
		wantErr  bool
	}{
		{
			name: "basic error response",
			response: &ierrors.ErrorResponse{
				Code:    ierrors.CodeUserNotFound,
				Message: "用户不存在",
				TraceID: "trace-123",
			},
			wantErr: false,
		},
		{
			name: "error without trace ID",
			response: &ierrors.ErrorResponse{
				Code:    ierrors.CodeInternalError,
				Message: "系统内部错误",
			},
			wantErr: false,
		},
		{
			name: "error with empty trace ID",
			response: &ierrors.ErrorResponse{
				Code:    ierrors.CodeValidationError,
				Message: "参数校验失败",
				TraceID: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.response.ToJSON()
			if len(got) == 0 {
				t.Error("ToJSON() returned empty bytes")
			}

			// Verify it's valid JSON
			var parsed ierrors.ErrorResponse
			if err := json.Unmarshal(got, &parsed); (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() produced invalid JSON: %v", err)
			}
		})
	}
}

func TestErrorResponse_JSONFields(t *testing.T) {
	resp := &ierrors.ErrorResponse{
		Code:    ierrors.CodeProjectNotFound,
		Message: "项目不存在",
		TraceID: "trace-abc-123",
	}

	data := resp.ToJSON()

	// Parse to verify field names
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check field names match expected JSON tags
	if _, ok := parsed["code"]; !ok {
		t.Error("Missing 'code' field in JSON")
	}
	if _, ok := parsed["message"]; !ok {
		t.Error("Missing 'message' field in JSON")
	}
	if _, ok := parsed["traceId"]; !ok {
		t.Error("Missing 'traceId' field in JSON")
	}
}

func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		traceID  string
		wantCode int
	}{
		{
			name:     "create user not found response",
			code:     ierrors.CodeUserNotFound,
			traceID:  "trace-1",
			wantCode: ierrors.CodeUserNotFound,
		},
		{
			name:     "create internal error response",
			code:     ierrors.CodeInternalError,
			traceID:  "trace-2",
			wantCode: ierrors.CodeInternalError,
		},
		{
			name:     "create validation error response",
			code:     ierrors.CodeValidationError,
			traceID:  "",
			wantCode: ierrors.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ierrors.NewErrorResponse(tt.code, tt.traceID)
			if got.Code != tt.wantCode {
				t.Errorf("NewErrorResponse().Code = %d, want %d", got.Code, tt.wantCode)
			}
			if got.TraceID != tt.traceID {
				t.Errorf("NewErrorResponse().TraceID = %s, want %s", got.TraceID, tt.traceID)
			}
			if got.Message == "" {
				t.Error("NewErrorResponse().Message should not be empty")
			}
		})
	}
}

func TestCodeToMessage(t *testing.T) {
	tests := []struct {
		name         string
		code         int
		wantNotEmpty bool
		expectedMsg  string
	}{
		// Identity
		{name: "CodeUserNotFound", code: ierrors.CodeUserNotFound, expectedMsg: "用户不存在"},
		{name: "CodeEmailDuplicate", code: ierrors.CodeEmailDuplicate, expectedMsg: "邮箱已被注册"},
		{name: "CodePasswordMismatch", code: ierrors.CodePasswordMismatch, expectedMsg: "密码错误"},

		// Project
		{name: "CodeProjectNotFound", code: ierrors.CodeProjectNotFound, expectedMsg: "项目不存在"},
		{name: "CodeProjectNameDuplicate", code: ierrors.CodeProjectNameDuplicate, expectedMsg: "项目名称已存在"},
		{name: "CodeProjectPrefixDuplicate", code: ierrors.CodeProjectPrefixDuplicate, expectedMsg: "项目前缀已存在"},

		// Knowledge
		{name: "CodeKnowledgeBaseEmpty", code: ierrors.CodeKnowledgeBaseEmpty, expectedMsg: "知识库为空，请先上传文档"},
		{name: "CodeDocumentParseFailed", code: ierrors.CodeDocumentParseFailed, expectedMsg: "文档解析失败"},

		// TestCase
		{name: "CodeCaseNotFound", code: ierrors.CodeCaseNotFound, expectedMsg: "测试用例不存在"},
		{name: "CodeEmptySteps", code: ierrors.CodeEmptySteps, expectedMsg: "用例步骤不能为空"},

		// TestPlan
		{name: "CodePlanNotFound", code: ierrors.CodePlanNotFound, expectedMsg: "测试计划不存在"},
		{name: "CodePlanArchived", code: ierrors.CodePlanArchived, expectedMsg: "计划已归档"},

		// Generation
		{name: "CodeLLMTimeout", code: ierrors.CodeLLMTimeout, expectedMsg: "LLM 调用超时，请重试"},
		{name: "CodeGenerationQueueFull", code: ierrors.CodeGenerationQueueFull, expectedMsg: "生成任务队列已满，请稍后重试"},

		// System
		{name: "CodeInternalError", code: ierrors.CodeInternalError, expectedMsg: "系统内部错误"},
		{name: "CodeUnauthorized", code: ierrors.CodeUnauthorized, expectedMsg: "未授权"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ierrors.CodeToMessage(tt.code)
			if got != tt.expectedMsg {
				t.Errorf("CodeToMessage(%d) = %q, want %q", tt.code, got, tt.expectedMsg)
			}
		})
	}
}

func TestCodeToMessage_UnknownCode(t *testing.T) {
	got := ierrors.CodeToMessage(99999)
	if got != "未知错误" {
		t.Errorf("CodeToMessage(99999) = %q, want %q", got, "未知错误")
	}
}

func TestErrorResponse_AllCodesHaveMessages(t *testing.T) {
	// Test that all defined error codes have corresponding messages
	codes := []int{
		ierrors.CodeUserNotFound,
		ierrors.CodeEmailDuplicate,
		ierrors.CodeUsernameDuplicate,
		ierrors.CodePasswordMismatch,
		ierrors.CodeInvalidEmail,
		ierrors.CodePermissionDenied,
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
		ierrors.CodeDocumentNotFound,
		ierrors.CodeDocumentParseFailed,
		ierrors.CodeUnsupportedDocumentType,
		ierrors.CodeEmptyChunks,
		ierrors.CodeEmbeddingFailed,
		ierrors.CodeVectorSearchFailed,
		ierrors.CodeKnowledgeBaseEmpty,
		ierrors.CodeDocumentProcessing,
		ierrors.CodeCaseNotFound,
		ierrors.CodeCaseNumberDuplicate,
		ierrors.CodeInvalidCaseNumber,
		ierrors.CodeEmptySteps,
		ierrors.CodeInvalidPriority,
		ierrors.CodeInvalidCaseType,
		ierrors.CodePlanNotFound,
		ierrors.CodePlanNameDuplicate,
		ierrors.CodePlanArchived,
		ierrors.CodeResultNotFound,
		ierrors.CodeCaseNotInPlan,
		ierrors.CodeDuplicateExecution,
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
		ierrors.CodeInternalError,
		ierrors.CodeDatabaseError,
		ierrors.CodeValidationError,
		ierrors.CodeUnauthorized,
		ierrors.CodeRateLimited,
	}

	for _, code := range codes {
		msg := ierrors.CodeToMessage(code)
		if msg == "" || msg == "未知错误" {
			t.Errorf("Code %d has no message defined", code)
		}
	}
}

func TestErrorResponse_SerializationRoundTrip(t *testing.T) {
	original := &ierrors.ErrorResponse{
		Code:    ierrors.CodeUserNotFound,
		Message: "用户不存在",
		TraceID: "trace-xyz-789",
	}

	// Serialize
	data := original.ToJSON()

	// Deserialize
	var parsed ierrors.ErrorResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Compare
	if parsed.Code != original.Code {
		t.Errorf("Code mismatch: got %d, want %d", parsed.Code, original.Code)
	}
	if parsed.Message != original.Message {
		t.Errorf("Message mismatch: got %q, want %q", parsed.Message, original.Message)
	}
	if parsed.TraceID != original.TraceID {
		t.Errorf("TraceID mismatch: got %q, want %q", parsed.TraceID, original.TraceID)
	}
}
