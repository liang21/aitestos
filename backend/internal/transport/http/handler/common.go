// Package handler provides HTTP handlers for the API
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Context key types for type-safe context values
type contextKey string

const (
	userIDContextKey    contextKey = "user_id"
	projectIDContextKey contextKey = "project_id"
	moduleIDContextKey  contextKey = "module_id"
	caseIDContextKey    contextKey = "case_id"
	planIDContextKey    contextKey = "plan_id"
	taskIDContextKey    contextKey = "task_id"
	draftIDContextKey   contextKey = "draft_id"
	documentIDContextKey contextKey = "document_id"
)

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// handleServiceError maps service errors to HTTP status codes
func handleServiceError(w http.ResponseWriter, err error) {
	// Default to internal server error
	code := http.StatusInternalServerError
	message := err.Error()

	// Map common errors to HTTP status codes
	switch err.Error() {
	case "user not found", "invalid credentials", "password mismatch", "unauthorized":
		code = http.StatusUnauthorized
	case "project not found", "module not found", "test case not found", "test plan not found",
		"generation task not found", "draft not found", "document not found", "config not found":
		code = http.StatusNotFound
	case "email already exists", "username already exists", "project name already exists",
		"project prefix already exists", "module name already exists", "module abbreviation already exists":
		code = http.StatusConflict
	case "invalid email format", "password too short", "invalid role",
		"invalid project prefix", "invalid module abbreviation", "invalid case number",
		"empty steps", "invalid priority", "invalid case type", "invalid document type",
		"invalid plan status", "invalid result status", "invalid request body":
		code = http.StatusBadRequest
	case "insufficient permissions", "permission denied":
		code = http.StatusForbidden
	case "draft already confirmed", "draft already rejected", "task already processed",
		"document is being processed":
		code = http.StatusConflict
	}

	respondWithError(w, code, message)
}

// getIDFromURL extracts a UUID from URL parameters
func getIDFromURL(r *http.Request, param string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, param)
	return uuid.Parse(idStr)
}

// getUserIDFromContext extracts user ID from context
func getUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return userID, ok
}

// getIntQueryParam extracts an integer from query parameters
func getIntQueryParam(r *http.Request, key string, defaultValue int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

// getStringQueryParam extracts a string from query parameters with default value
func getStringQueryParam(r *http.Request, key string, defaultValue string) string {
	if v := r.URL.Query().Get(key); v != "" {
		return v
	}
	return defaultValue
}
