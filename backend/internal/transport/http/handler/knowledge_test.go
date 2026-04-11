// Package handler provides HTTP handlers for the API
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	docservice "github.com/liang21/aitestos/internal/service/knowledge"
)

// MockDocumentService implements docservice.DocumentService for testing
type MockDocumentService struct {
	mock.Mock
}

func (m *MockDocumentService) UploadDocument(ctx context.Context, req *docservice.UploadDocumentRequest) (*knowledge.Document, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*knowledge.Document), args.Error(1)
}

func (m *MockDocumentService) GetDocument(ctx context.Context, id uuid.UUID) (*docservice.DocumentDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*docservice.DocumentDetail), args.Error(1)
}

func (m *MockDocumentService) ListDocuments(ctx context.Context, opts docservice.DocumentListOptions) ([]*knowledge.Document, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]*knowledge.Document), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentService) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockDocumentService) ProcessDocument(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentService) GetChunks(ctx context.Context, documentID uuid.UUID) ([]*docservice.ChunkInfo, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).([]*docservice.ChunkInfo), args.Error(1)
}

func TestUploadDocumentHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful upload", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		projectID := uuid.New()
		userID := uuid.New()
		mockSvc.On("UploadDocument", mock.Anything, mock.Anything).Return(&knowledge.Document{}, nil)

		handler := NewKnowledgeHandler(mockSvc)
		require.NotNil(t, handler)

		body := map[string]interface{}{
			"project_id": projectID.String(),
			"name":       "Test Document",
			"type":       "prd",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/knowledge/documents", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, userID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.UploadDocument(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("missing user context", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		handler := NewKnowledgeHandler(mockSvc)

		body := map[string]interface{}{
			"project_id": uuid.New().String(),
			"name":       "Test Document",
			"type":       "prd",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/knowledge/documents", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UploadDocument(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid document type", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		handler := NewKnowledgeHandler(mockSvc)

		body := map[string]interface{}{
			"project_id": uuid.New().String(),
			"name":       "Test Document",
			"type":       "invalid_type",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/knowledge/documents", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.UploadDocument(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetDocumentHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		docID := uuid.New()
		mockSvc.On("GetDocument", mock.Anything, docID).Return(&docservice.DocumentDetail{
			Document:   &knowledge.Document{},
			ChunkCount: 0,
		}, nil)

		handler := NewKnowledgeHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/knowledge/documents/"+docID.String(), nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", docID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetDocument(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("document not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		docID := uuid.New()
		mockSvc.On("GetDocument", mock.Anything, docID).Return(nil, knowledge.ErrDocumentNotFound)

		handler := NewKnowledgeHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/knowledge/documents/"+docID.String(), nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", docID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetDocument(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestListDocumentsHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful list", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		projectID := uuid.New()
		mockSvc.On("ListDocuments", mock.Anything, mock.MatchedBy(func(opts docservice.DocumentListOptions) bool {
			return opts.ProjectID == projectID
		})).Return([]*knowledge.Document{}, int64(0), nil)

		handler := NewKnowledgeHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/knowledge/documents?project_id="+projectID.String(), nil)
		w := httptest.NewRecorder()

		handler.ListDocuments(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteDocumentHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful delete", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		docID := uuid.New()
		mockSvc.On("DeleteDocument", mock.Anything, docID).Return(nil)

		handler := NewKnowledgeHandler(mockSvc)

		req := httptest.NewRequest("DELETE", "/api/v1/knowledge/documents/"+docID.String(), nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", docID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.DeleteDocument(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("document not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		docID := uuid.New()
		mockSvc.On("DeleteDocument", mock.Anything, docID).Return(knowledge.ErrDocumentNotFound)

		handler := NewKnowledgeHandler(mockSvc)

		req := httptest.NewRequest("DELETE", "/api/v1/knowledge/documents/"+docID.String(), nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", docID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.DeleteDocument(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetChunksHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get chunks", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		docID := uuid.New()
		mockSvc.On("GetChunks", mock.Anything, docID).Return([]*docservice.ChunkInfo{}, nil)

		handler := NewKnowledgeHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/knowledge/documents/"+docID.String()+"/chunks", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", docID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetChunks(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestProcessDocumentHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful process", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		docID := uuid.New()
		mockSvc.On("ProcessDocument", mock.Anything, docID).Return(nil)

		handler := NewKnowledgeHandler(mockSvc)

		req := httptest.NewRequest("POST", "/api/v1/knowledge/documents/"+docID.String()+"/process", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", docID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.ProcessDocument(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)
	})

	t.Run("document processing already in progress", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockDocumentService)
		docID := uuid.New()
		mockSvc.On("ProcessDocument", mock.Anything, docID).Return(errors.New("document is being processed"))

		handler := NewKnowledgeHandler(mockSvc)

		req := httptest.NewRequest("POST", "/api/v1/knowledge/documents/"+docID.String()+"/process", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", docID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.ProcessDocument(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
