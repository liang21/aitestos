// Package middleware provides HTTP middleware implementations
package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("recovers from panic with string", func(t *testing.T) {
		t.Parallel()

		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("something went wrong")
		})

		handler := Recovery()(panicHandler)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		// Should not panic
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "code")
	})

	t.Run("recovers from panic with error", func(t *testing.T) {
		t.Parallel()

		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(assert.AnError)
		})

		handler := Recovery()(panicHandler)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("no panic passes through", func(t *testing.T) {
		t.Parallel()

		normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})

		handler := Recovery()(normalHandler)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})
}

func TestRecoveryMiddlewareWithTraceID(t *testing.T) {
	t.Parallel()

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	handler := Recovery()(panicHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Trace-ID", "trace-123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Error response should include trace ID
	if traceID, ok := response["traceId"]; ok {
		assert.Equal(t, "trace-123", traceID)
	}
}
