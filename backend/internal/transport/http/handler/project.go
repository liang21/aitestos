// Package handler provides HTTP handlers for the API
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	projectservice "github.com/liang21/aitestos/internal/service/project"
)

// ProjectHandler handles project-related HTTP requests
type ProjectHandler struct {
	projectService projectservice.ProjectService
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler(projectService projectservice.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// CreateProject handles project creation
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req projectservice.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := h.projectService.CreateProject(r.Context(), &req, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, project)
}

// GetProject handles getting a single project
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	detail, err := h.projectService.GetProject(r.Context(), projectID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, detail)
}

// ListProjects handles listing projects
func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	opts := projectservice.ListOptions{
		Offset:   getIntQueryParam(r, "offset", 0),
		Limit:    getIntQueryParam(r, "limit", 10),
		Keywords: r.URL.Query().Get("keywords"),
	}

	projects, total, err := h.projectService.ListProjects(r.Context(), opts)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":   projects,
		"total":  total,
		"offset": opts.Offset,
		"limit":  opts.Limit,
	})
}

// UpdateProject handles project updates
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	var req projectservice.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := h.projectService.UpdateProject(r.Context(), projectID, &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, project)
}

// DeleteProject handles project deletion
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	if err := h.projectService.DeleteProject(r.Context(), projectID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateModule handles module creation
func (h *ProjectHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	var req projectservice.CreateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	module, err := h.projectService.CreateModule(r.Context(), projectID, &req, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, module)
}

// ListModules handles listing modules
func (h *ProjectHandler) ListModules(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	modules, err := h.projectService.ListModules(r.Context(), projectID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, modules)
}

// DeleteModule handles module deletion
func (h *ProjectHandler) DeleteModule(w http.ResponseWriter, r *http.Request) {
	moduleID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid module ID")
		return
	}

	if err := h.projectService.DeleteModule(r.Context(), moduleID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SetConfig handles setting project configuration
func (h *ProjectHandler) SetConfig(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	key := chi.URLParam(r, "key")
	if key == "" {
		respondWithError(w, http.StatusBadRequest, "missing config key")
		return
	}

	var value map[string]any
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.projectService.SetConfig(r.Context(), projectID, key, value); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ListConfigs handles listing project configurations
func (h *ProjectHandler) ListConfigs(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	configs, err := h.projectService.ListConfigs(r.Context(), projectID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, configs)
}
