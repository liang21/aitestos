// Package http provides HTTP transport layer implementation
package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	router := NewRouter(nil, zerolog.Nop(), nil)
	assert.NotNil(t, router)
}

func TestRouterRoutes(t *testing.T) {
	t.Parallel()

	router := NewRouter(nil, zerolog.Nop(), nil)
	require.NotNil(t, router)

	tests := []struct {
		method string
		path   string
	}{
		// Auth routes
		{"POST", "/api/v1/auth/register"},
		{"POST", "/api/v1/auth/login"},
		{"POST", "/api/v1/auth/refresh"},

		// Project routes
		{"GET", "/api/v1/projects"},
		{"POST", "/api/v1/projects"},
		{"GET", "/api/v1/projects/00000000-0000-0000-0000-000000000001"},
		{"PUT", "/api/v1/projects/00000000-0000-0000-0000-000000000001"},
		{"DELETE", "/api/v1/projects/00000000-0000-0000-0000-000000000001"},

		// Module routes
		{"GET", "/api/v1/projects/00000000-0000-0000-0000-000000000001/modules"},
		{"POST", "/api/v1/projects/00000000-0000-0000-0000-000000000001/modules"},
		{"DELETE", "/api/v1/modules/00000000-0000-0000-0000-000000000001"},

		// Config routes
		{"GET", "/api/v1/projects/00000000-0000-0000-0000-000000000001/configs"},
		{"PUT", "/api/v1/projects/00000000-0000-0000-0000-000000000001/configs/test-key"},

		// TestCase routes
		{"GET", "/api/v1/testcases"},
		{"POST", "/api/v1/testcases"},
		{"GET", "/api/v1/testcases/00000000-0000-0000-0000-000000000001"},
		{"PUT", "/api/v1/testcases/00000000-0000-0000-0000-000000000001"},
		{"DELETE", "/api/v1/testcases/00000000-0000-0000-0000-000000000001"},

		// TestPlan routes
		{"GET", "/api/v1/plans"},
		{"POST", "/api/v1/plans"},
		{"GET", "/api/v1/plans/00000000-0000-0000-0000-000000000001"},
		{"PUT", "/api/v1/plans/00000000-0000-0000-0000-000000000001"},
		{"DELETE", "/api/v1/plans/00000000-0000-0000-0000-000000000001"},

		// TestPlan case management
		{"POST", "/api/v1/plans/00000000-0000-0000-0000-000000000001/cases"},
		{"DELETE", "/api/v1/plans/00000000-0000-0000-0000-000000000001/cases/00000000-0000-0000-0000-000000000001"},

		// TestPlan result routes
		{"POST", "/api/v1/plans/00000000-0000-0000-0000-000000000001/results"},
		{"GET", "/api/v1/plans/00000000-0000-0000-0000-000000000001/results"},

		// Generation routes
		{"POST", "/api/v1/generation/tasks"},
		{"GET", "/api/v1/generation/tasks/00000000-0000-0000-0000-000000000001"},
		{"GET", "/api/v1/generation/tasks/00000000-0000-0000-0000-000000000001/drafts"},
		{"POST", "/api/v1/generation/drafts/00000000-0000-0000-0000-000000000001/confirm"},
		{"POST", "/api/v1/generation/drafts/00000000-0000-0000-0000-000000000001/reject"},
		{"POST", "/api/v1/generation/drafts/batch-confirm"},

		// Knowledge routes
		{"GET", "/api/v1/knowledge/documents"},
		{"POST", "/api/v1/knowledge/documents"},
		{"GET", "/api/v1/knowledge/documents/00000000-0000-0000-0000-000000000001"},
		{"DELETE", "/api/v1/knowledge/documents/00000000-0000-0000-0000-000000000001"},
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.path, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			// Route should exist (not 404), though may return 401 if auth required
			assert.NotEqual(t, http.StatusNotFound, w.Code, "Route %s %s should exist", tt.method, tt.path)
		})
	}
}

func TestRouterHealthCheck(t *testing.T) {
	t.Parallel()

	router := NewRouter(nil, zerolog.Nop(), nil)
	require.NotNil(t, router)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterNotFound(t *testing.T) {
	t.Parallel()

	router := NewRouter(nil, zerolog.Nop(), nil)
	require.NotNil(t, router)

	req := httptest.NewRequest("GET", "/nonexistent-path", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterMethodNotAllowed(t *testing.T) {
	t.Parallel()

	router := NewRouter(nil, zerolog.Nop(), nil)
	require.NotNil(t, router)

	// Try GET on a POST-only route
	req := httptest.NewRequest("GET", "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
