// Package handler provides HTTP handlers for the API
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/generation"
	genservice "github.com/liang21/aitestos/internal/service/generation"
)

// GenerationHandler handles generation-related HTTP requests
type GenerationHandler struct {
	genService genservice.GenerationService
}

// NewGenerationHandler creates a new GenerationHandler
func NewGenerationHandler(genService genservice.GenerationService) *GenerationHandler {
	return &GenerationHandler{
		genService: genService,
	}
}

// CreateTask handles generation task creation
func (h *GenerationHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req genservice.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.genService.CreateTask(r.Context(), &req, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

// GetTask handles getting a single generation task
func (h *GenerationHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromURL(r, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	task, err := h.genService.GetTask(r.Context(), taskID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// ListTasks handles listing generation tasks
func (h *GenerationHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "project_id is required")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}

	opts := genservice.ListTaskOptions{
		Offset:   getIntQueryParam(r, "offset", 0),
		Limit:    getIntQueryParam(r, "limit", 10),
		Status:   r.URL.Query().Get("status"),
		Keywords: r.URL.Query().Get("keywords"),
	}

	// Support module_id filter
	if moduleIDStr := r.URL.Query().Get("module_id"); moduleIDStr != "" {
		moduleID, err := uuid.Parse(moduleIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid module ID format")
			return
		}
		opts.ModuleID = moduleID
	}

	tasks, total, err := h.genService.ListTasks(r.Context(), projectID, opts)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":   tasks,
		"total":  total,
		"offset": opts.Offset,
		"limit":  opts.Limit,
	})
}

// GetDrafts handles getting drafts for a task
func (h *GenerationHandler) GetDrafts(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromURL(r, "taskID")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	drafts, err := h.genService.GetDrafts(r.Context(), taskID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, drafts)
}

// ConfirmDraft handles confirming a single draft
func (h *GenerationHandler) ConfirmDraft(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	draftID, err := getIDFromURL(r, "draftID")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid draft ID")
		return
	}

	var req struct {
		ModuleID uuid.UUID `json:"module_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	confirmReq := &genservice.ConfirmDraftRequest{
		DraftID:  draftID,
		ModuleID: req.ModuleID,
	}

	tc, err := h.genService.ConfirmDraft(r.Context(), confirmReq, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, tc)
}

// RejectDraft handles rejecting a single draft
func (h *GenerationHandler) RejectDraft(w http.ResponseWriter, r *http.Request) {
	draftID, err := getIDFromURL(r, "draftID")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid draft ID")
		return
	}

	var req struct {
		Reason   string `json:"reason"`
		Feedback string `json:"feedback"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	rejectReq := &genservice.RejectDraftRequest{
		DraftID:  draftID,
		Reason:   generation.RejectionReason(req.Reason),
		Feedback: req.Feedback,
	}

	if err := h.genService.RejectDraft(r.Context(), rejectReq); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// BatchConfirm handles confirming multiple drafts at once
func (h *GenerationHandler) BatchConfirm(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req struct {
		DraftIDs []uuid.UUID `json:"draft_ids"`
		ModuleID uuid.UUID   `json:"module_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	batchReq := &genservice.BatchConfirmRequest{
		DraftIDs: req.DraftIDs,
		ModuleID: req.ModuleID,
	}

	result, err := h.genService.BatchConfirm(r.Context(), batchReq, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, result)
}

// ListAllDrafts handles listing all drafts with filters
func (h *GenerationHandler) ListAllDrafts(w http.ResponseWriter, r *http.Request) {
	opts := genservice.ListAllDraftsOptions{
		Offset: getIntQueryParam(r, "offset", 0),
		Limit:  getIntQueryParam(r, "limit", 10),
		Status: r.URL.Query().Get("status"),
	}

	// Parse project_id if provided
	if projectIDStr := r.URL.Query().Get("project_id"); projectIDStr != "" {
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid project ID format")
			return
		}
		opts.ProjectID = projectID
	}

	// Parse task_id if provided
	if taskIDStr := r.URL.Query().Get("task_id"); taskIDStr != "" {
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid task ID format")
			return
		}
		opts.TaskID = taskID
	}

	drafts, total, err := h.genService.ListAllDrafts(r.Context(), opts)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":   drafts,
		"total":  total,
		"offset": opts.Offset,
		"limit":  opts.Limit,
	})
}
