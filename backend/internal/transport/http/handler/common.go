// Package handler provides HTTP handlers for the API
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	domainGeneration "github.com/liang21/aitestos/internal/domain/generation"
	domainIdentity "github.com/liang21/aitestos/internal/domain/identity"
	domainKnowledge "github.com/liang21/aitestos/internal/domain/knowledge"
	domainProject "github.com/liang21/aitestos/internal/domain/project"
	domainTestcase "github.com/liang21/aitestos/internal/domain/testcase"
	domainTestplan "github.com/liang21/aitestos/internal/domain/testplan"
)

// Context key types for type-safe context values
type contextKey string

const (
	userIDContextKey     contextKey = "user_id"
	projectIDContextKey  contextKey = "project_id"
	moduleIDContextKey   contextKey = "module_id"
	caseIDContextKey     contextKey = "case_id"
	planIDContextKey     contextKey = "plan_id"
	taskIDContextKey     contextKey = "task_id"
	draftIDContextKey    contextKey = "draft_id"
	documentIDContextKey contextKey = "document_id"
)

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to encode response"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(data)
}

// handleServiceError maps service errors to HTTP status codes using errors.Is()
func handleServiceError(w http.ResponseWriter, err error) {
	// 401 Unauthorized
	switch {
	case errors.Is(err, domainIdentity.ErrUserNotFound),
		errors.Is(err, domainIdentity.ErrPasswordMismatch),
		errors.Is(err, domainIdentity.ErrTokenInvalid),
		errors.Is(err, domainIdentity.ErrTokenExpired),
		errors.Is(err, domainIdentity.ErrTokenRevoked),
		errors.Is(err, domainIdentity.ErrTokenMissing):
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// 403 Forbidden
	if errors.Is(err, domainIdentity.ErrPermissionDenied) {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	// 404 Not Found
	switch {
	case errors.Is(err, domainProject.ErrProjectNotFound),
		errors.Is(err, domainProject.ErrModuleNotFound),
		errors.Is(err, domainProject.ErrConfigNotFound),
		errors.Is(err, domainTestcase.ErrCaseNotFound),
		errors.Is(err, domainTestplan.ErrPlanNotFound),
		errors.Is(err, domainTestplan.ErrResultNotFound),
		errors.Is(err, domainGeneration.ErrTaskNotFound),
		errors.Is(err, domainGeneration.ErrDraftNotFound),
		errors.Is(err, domainKnowledge.ErrDocumentNotFound):
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// 409 Conflict
	switch {
	case errors.Is(err, domainIdentity.ErrEmailDuplicate),
		errors.Is(err, domainIdentity.ErrUsernameDuplicate),
		errors.Is(err, domainProject.ErrProjectNameDuplicate),
		errors.Is(err, domainProject.ErrProjectPrefixDuplicate),
		errors.Is(err, domainProject.ErrModuleNameDuplicate),
		errors.Is(err, domainProject.ErrModuleAbbrevDuplicate),
		errors.Is(err, domainGeneration.ErrDraftAlreadyConfirmed),
		errors.Is(err, domainGeneration.ErrDraftAlreadyRejected),
		errors.Is(err, domainGeneration.ErrTaskAlreadyProcessed),
		errors.Is(err, domainKnowledge.ErrDocumentProcessing):
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}

	// 400 Bad Request
	switch {
	case errors.Is(err, domainIdentity.ErrInvalidEmail),
		errors.Is(err, domainIdentity.ErrPasswordTooShort),
		errors.Is(err, domainIdentity.ErrInvalidUsername),
		errors.Is(err, domainIdentity.ErrInvalidRole),
		errors.Is(err, domainProject.ErrInvalidProjectPrefix),
		errors.Is(err, domainProject.ErrInvalidModuleAbbrev),
		errors.Is(err, domainTestcase.ErrInvalidCaseNumber),
		errors.Is(err, domainTestcase.ErrEmptySteps),
		errors.Is(err, domainTestcase.ErrInvalidPriority),
		errors.Is(err, domainTestplan.ErrPlanNameDuplicate):
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 500 Internal Server Error (fallback)
	respondWithError(w, http.StatusInternalServerError, err.Error())
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
