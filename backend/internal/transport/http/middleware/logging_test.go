// Package middleware provides HTTP middleware implementations
package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		status     int
		wantFields map[string]interface{}
	}{
		{
			name:   "GET request",
			method: "GET",
			path:   "/api/v1/test",
			status: http.StatusOK,
			wantFields: map[string]interface{}{
				"method": "GET",
				"path":   "/api/v1/test",
				"status": float64(http.StatusOK),
			},
		},
		{
			name:   "POST request",
			method: "POST",
			path:   "/api/v1/users",
			status: http.StatusCreated,
			wantFields: map[string]interface{}{
				"method": "POST",
				"path":   "/api/v1/users",
				"status": float64(http.StatusCreated),
			},
		},
		{
			name:   "Not Found",
			method: "GET",
			path:   "/nonexistent",
			status: http.StatusNotFound,
			wantFields: map[string]interface{}{
				"method": "GET",
				"path":   "/nonexistent",
				"status": float64(http.StatusNotFound),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := zerolog.New(&buf)

			handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			}))

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			// Split log output into lines and find the completion log
			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			var completionLog map[string]interface{}
			for _, line := range lines {
				if line == "" {
					continue
				}
				var logEntry map[string]interface{}
				err := json.Unmarshal([]byte(line), &logEntry)
				assert.NoError(t, err)

				// Check if this is the completion log (has status field)
				if _, hasStatus := logEntry["status"]; hasStatus {
					completionLog = logEntry
					break
				}
			}

			// Verify we found the completion log
			assert.NotNil(t, completionLog, "should have completion log")

			// Verify expected fields
			for key, expected := range tt.wantFields {
				if actual, ok := completionLog[key]; ok {
					assert.Equal(t, expected, actual, "log field %s mismatch", key)
				}
			}

			// Should have duration
			_, hasDuration := completionLog["duration"]
			assert.True(t, hasDuration, "log should contain duration")
		})
	}
}

func TestLoggingMiddlewareWithTraceID(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	traceID := "test-trace-123"
	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Trace-ID", traceID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Parse all log lines
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	var foundTraceID bool
	for _, line := range lines {
		if line == "" {
			continue
		}
		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(line), &logEntry)
		assert.NoError(t, err)

		// Check if this log has trace_id
		if traceIDVal, ok := logEntry["trace_id"]; ok {
			if traceIDVal == traceID {
				foundTraceID = true
				break
			}
		}
	}

	assert.True(t, foundTraceID, "should find trace ID in logs")
}

func TestLoggingMiddlewareGeneratesTraceID(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that trace ID is in response header
		traceID := w.Header().Get("X-Trace-ID")
		assert.NotEmpty(t, traceID, "should generate trace ID")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Response should have trace ID header
	traceID := w.Header().Get("X-Trace-ID")
	assert.NotEmpty(t, traceID, "response should contain trace ID")
}
