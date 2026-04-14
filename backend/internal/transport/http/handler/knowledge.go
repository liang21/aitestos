// Package handler provides HTTP handlers for the API
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	docservice "github.com/liang21/aitestos/internal/service/knowledge"
)

// KnowledgeHandler handles knowledge document-related HTTP requests
type KnowledgeHandler struct {
	docService docservice.DocumentService
}

// NewKnowledgeHandler creates a new KnowledgeHandler
func NewKnowledgeHandler(docService docservice.DocumentService) *KnowledgeHandler {
	return &KnowledgeHandler{
		docService: docService,
	}
}

// UploadDocument handles document upload
func (h *KnowledgeHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req struct {
		ProjectID uuid.UUID `json:"project_id"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate document type
	docType, err := knowledge.ParseDocumentType(req.Type)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid document type")
		return
	}

	uploadReq := &docservice.UploadDocumentRequest{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Type:      docType.String(),
		UserID:    userID,
	}

	doc, err := h.docService.UploadDocument(r.Context(), uploadReq)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, doc)
}

// GetDocument handles getting a single document
func (h *KnowledgeHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	docID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	detail, err := h.docService.GetDocument(r.Context(), docID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, detail)
}

// ListDocuments handles listing documents
func (h *KnowledgeHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(r.URL.Query().Get("project_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	opts := docservice.DocumentListOptions{
		Offset:    getIntQueryParam(r, "offset", 0),
		Limit:     getIntQueryParam(r, "limit", 10),
		ProjectID: projectID,
		Type:      r.URL.Query().Get("type"),
		Status:    r.URL.Query().Get("status"),
	}

	docs, total, err := h.docService.ListDocuments(r.Context(), opts)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":   docs,
		"total":  total,
		"offset": opts.Offset,
		"limit":  opts.Limit,
	})
}

// DeleteDocument handles document deletion
func (h *KnowledgeHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	docID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	if err := h.docService.DeleteDocument(r.Context(), docID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetChunks handles getting document chunks
func (h *KnowledgeHandler) GetChunks(w http.ResponseWriter, r *http.Request) {
	docID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	chunks, err := h.docService.GetChunks(r.Context(), docID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, chunks)
}

// ProcessDocument handles triggering document processing
func (h *KnowledgeHandler) ProcessDocument(w http.ResponseWriter, r *http.Request) {
	docID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	if err := h.docService.ProcessDocument(r.Context(), docID); err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusAccepted, map[string]string{
		"message": "document processing started",
	})
}
