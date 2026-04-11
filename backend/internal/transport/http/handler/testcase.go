// Package handler provides HTTP handlers for the API
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	domaintestcase "github.com/liang21/aitestos/internal/domain/testcase"
	caseservice "github.com/liang21/aitestos/internal/service/testcase"
)

// TestCaseHandler handles test case-related HTTP requests
type TestCaseHandler struct {
	caseService caseservice.CaseService
}

// NewTestCaseHandler creates a new TestCaseHandler
func NewTestCaseHandler(caseService caseservice.CaseService) *TestCaseHandler {
	return &TestCaseHandler{
		caseService: caseService,
	}
}

// CreateCase handles test case creation
func (h *TestCaseHandler) CreateCase(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req caseservice.CreateCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tc, err := h.caseService.CreateCase(r.Context(), &req, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, tc)
}

// GetCase handles getting a single test case
func (h *TestCaseHandler) GetCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid case ID")
		return
	}

	detail, err := h.caseService.GetCaseDetail(r.Context(), caseID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, detail)
}

// ListCases handles listing test cases
func (h *TestCaseHandler) ListCases(w http.ResponseWriter, r *http.Request) {
	opts := caseservice.CaseListOptions{
		Offset:   getIntQueryParam(r, "offset", 0),
		Limit:    getIntQueryParam(r, "limit", 10),
		Keywords: r.URL.Query().Get("keywords"),
	}

	// Parse optional filters
	if projectID := r.URL.Query().Get("project_id"); projectID != "" {
		if id, err := uuid.Parse(projectID); err == nil {
			opts.ProjectID = id
		}
	}

	if moduleID := r.URL.Query().Get("module_id"); moduleID != "" {
		if id, err := uuid.Parse(moduleID); err == nil {
			opts.ModuleID = id
		}
	}

	opts.Status = r.URL.Query().Get("status")
	opts.CaseType = r.URL.Query().Get("case_type")
	opts.Priority = r.URL.Query().Get("priority")

	var cases []*domaintestcase.TestCase
	var total int64
	var err error

	if opts.ModuleID != uuid.Nil {
		cases, total, err = h.caseService.ListByModule(r.Context(), opts.ModuleID, opts)
	} else if opts.ProjectID != uuid.Nil {
		cases, total, err = h.caseService.ListByProject(r.Context(), opts.ProjectID, opts)
	} else {
		respondWithError(w, http.StatusBadRequest, "project_id or module_id is required")
		return
	}

	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"data":   cases,
		"total":  total,
		"offset": opts.Offset,
		"limit":  opts.Limit,
	})
}

// UpdateCase handles test case updates
func (h *TestCaseHandler) UpdateCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid case ID")
		return
	}

	var req caseservice.UpdateCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tc, err := h.caseService.UpdateCase(r.Context(), caseID, &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, tc)
}

// DeleteCase handles test case deletion
func (h *TestCaseHandler) DeleteCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid case ID")
		return
	}

	if err := h.caseService.DeleteCase(r.Context(), caseID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
