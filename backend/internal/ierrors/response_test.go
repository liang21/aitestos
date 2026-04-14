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
