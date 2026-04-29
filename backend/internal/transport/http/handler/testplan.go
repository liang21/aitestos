// Package handler provides HTTP handlers for the API
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	planservice "github.com/liang21/aitestos/internal/service/testplan"
)

// TestPlanHandler handles test plan-related HTTP requests
type TestPlanHandler struct {
	planService planservice.PlanService
}

// NewTestPlanHandler creates a new TestPlanHandler
func NewTestPlanHandler(planService planservice.PlanService) *TestPlanHandler {
	return &TestPlanHandler{
		planService: planService,
	}
}

// CreatePlan handles test plan creation
func (h *TestPlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req planservice.CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	plan, err := h.planService.CreatePlan(r.Context(), &req, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, plan)
}

// GetPlan handles getting a single test plan
func (h *TestPlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

	detail, err := h.planService.GetPlan(r.Context(), planID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, detail)
}

// ListPlans handles listing test plans
func (h *TestPlanHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	projectID, err := getIDFromURL(r, "projectID")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	opts := planservice.PlanListOptions{
		Offset:   getIntQueryParam(r, "offset", 0),
		Limit:    getIntQueryParam(r, "limit", 10),
		Status:   r.URL.Query().Get("status"),
		Keywords: r.URL.Query().Get("keywords"),
	}

	plans, total, err := h.planService.ListPlans(r.Context(), projectID, opts)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":   plans,
		"total":  total,
		"offset": opts.Offset,
		"limit":  opts.Limit,
	})
}

// UpdatePlan handles test plan updates
func (h *TestPlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	planID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

	var req planservice.UpdatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	plan, err := h.planService.UpdatePlan(r.Context(), planID, &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, plan)
}

// DeletePlan handles test plan deletion
func (h *TestPlanHandler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	planID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

	if err := h.planService.DeletePlan(r.Context(), planID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddCases handles adding test cases to a plan
func (h *TestPlanHandler) AddCases(w http.ResponseWriter, r *http.Request) {
	planID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

	var req struct {
		CaseIDs []string `json:"case_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	caseIDs := make([]uuid.UUID, 0, len(req.CaseIDs))
	for _, idStr := range req.CaseIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid case ID: "+idStr)
			return
		}
		caseIDs = append(caseIDs, id)
	}

	if err := h.planService.AddCases(r.Context(), planID, caseIDs); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RecordResult handles recording a test execution result
func (h *TestPlanHandler) RecordResult(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req planservice.RecordResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.planService.RecordResult(r.Context(), &req, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}

// GetResults handles getting test results for a plan
func (h *TestPlanHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	planID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

	results, err := h.planService.GetResults(r.Context(), planID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}

// RemoveCase handles removing a test case from a plan
func (h *TestPlanHandler) RemoveCase(w http.ResponseWriter, r *http.Request) {
	planID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

	caseID, err := getIDFromURL(r, "caseId")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid case ID")
		return
	}

	if err := h.planService.RemoveCase(r.Context(), planID, caseID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePlanStatus handles test plan status updates
func (h *TestPlanHandler) UpdatePlanStatus(w http.ResponseWriter, r *http.Request) {
	planID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Status == "" {
		respondWithError(w, http.StatusBadRequest, "status is required")
		return
	}

	if err := h.planService.UpdatePlanStatus(r.Context(), planID, req.Status); err != nil {
		handleServiceError(w, err)
		return
	}

	// Fetch updated plan and return it
	detail, err := h.planService.GetPlan(r.Context(), planID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, detail.TestPlan)
}
