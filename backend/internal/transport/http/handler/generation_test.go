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
	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/testcase"
	genservice "github.com/liang21/aitestos/internal/service/generation"
)

// MockGenerationService implements genservice.GenerationService for testing
type MockGenerationService struct {
	mock.Mock
}

func (m *MockGenerationService) CreateTask(ctx context.Context, req *genservice.CreateTaskRequest, userID uuid.UUID) (*generation.GenerationTask, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*generation.GenerationTask), args.Error(1)
}

func (m *MockGenerationService) GetTask(ctx context.Context, id uuid.UUID) (*generation.GenerationTask, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*generation.GenerationTask), args.Error(1)
}

func (m *MockGenerationService) ListTasks(ctx context.Context, projectID uuid.UUID, opts genservice.ListTaskOptions) ([]*generation.GenerationTask, int64, error) {
	args := m.Called(ctx, projectID, opts)
	return args.Get(0).([]*generation.GenerationTask), args.Get(1).(int64), args.Error(2)
}

func (m *MockGenerationService) GetDrafts(ctx context.Context, taskID uuid.UUID) ([]*generation.GeneratedCaseDraft, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).([]*generation.GeneratedCaseDraft), args.Error(1)
}

func (m *MockGenerationService) ConfirmDraft(ctx context.Context, req *genservice.ConfirmDraftRequest, userID uuid.UUID) (*testcase.TestCase, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*testcase.TestCase), args.Error(1)
}

func (m *MockGenerationService) RejectDraft(ctx context.Context, req *genservice.RejectDraftRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockGenerationService) BatchConfirm(ctx context.Context, req *genservice.BatchConfirmRequest, userID uuid.UUID) (*genservice.BatchConfirmResult, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*genservice.BatchConfirmResult), args.Error(1)
}

func (m *MockGenerationService) ProcessTask(ctx context.Context, taskID uuid.UUID) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

// Helper function to create a request with chi URL parameters set
func createRequestWithChiParams(method, url string, body []byte, urlParams map[string]string) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	// Set chi URL parameters
	if len(urlParams) > 0 {
		rctx := chi.NewRouteContext()
		for key, value := range urlParams {
			rctx.URLParams.Add(key, value)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}

	return req
}

func TestCreateTaskHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful creation", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		mockSvc.On("CreateTask", mock.Anything, mock.Anything, mock.Anything).Return(&generation.GenerationTask{}, nil)

		handler := NewGenerationHandler(mockSvc)
		require.NotNil(t, handler)

		body := map[string]interface{}{
			"project_id": uuid.New().String(),
			"module_id":  uuid.New().String(),
			"prompt":     "Generate test cases for user login functionality",
			"case_count": 5,
			"scene_types": []string{"positive", "negative"},
		}
		jsonBody, _ := json.Marshal(body)

		req := createRequestWithChiParams("POST", "/api/v1/generation/tasks", jsonBody, nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateTask(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("missing user context", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"project_id": uuid.New().String(),
			"module_id":  uuid.New().String(),
			"prompt":     "Generate test cases",
		}
		jsonBody, _ := json.Marshal(body)

		req := createRequestWithChiParams("POST", "/api/v1/generation/tasks", jsonBody, nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateTask(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid prompt too short", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		// Set mock expectation - service should return validation error for short prompt
		mockSvc.On("CreateTask", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("invalid request: prompt too short"))

		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"project_id": uuid.New().String(),
			"module_id":  uuid.New().String(),
			"prompt":     "short",
		}
		jsonBody, _ := json.Marshal(body)

		req := createRequestWithChiParams("POST", "/api/v1/generation/tasks", jsonBody, nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateTask(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetTaskHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		taskID := uuid.New()
		mockSvc.On("GetTask", mock.Anything, taskID).Return(&generation.GenerationTask{}, nil)

		handler := NewGenerationHandler(mockSvc)

		req := createRequestWithChiParams("GET", "/api/v1/generation/tasks/"+taskID.String(), nil, map[string]string{"id": taskID.String()})
		w := httptest.NewRecorder()

		handler.GetTask(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("task not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		taskID := uuid.New()
		mockSvc.On("GetTask", mock.Anything, taskID).Return(nil, generation.ErrTaskNotFound)

		handler := NewGenerationHandler(mockSvc)

		req := createRequestWithChiParams("GET", "/api/v1/generation/tasks/"+taskID.String(), nil, map[string]string{"id": taskID.String()})
		w := httptest.NewRecorder()

		handler.GetTask(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid task ID", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		handler := NewGenerationHandler(mockSvc)

		req := createRequestWithChiParams("GET", "/api/v1/generation/tasks/invalid-uuid", nil, map[string]string{"id": "invalid-uuid"})
		w := httptest.NewRecorder()

		handler.GetTask(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetDraftsHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get drafts", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		taskID := uuid.New()
		mockSvc.On("GetDrafts", mock.Anything, taskID).Return([]*generation.GeneratedCaseDraft{}, nil)

		handler := NewGenerationHandler(mockSvc)

		req := createRequestWithChiParams("GET", "/api/v1/generation/tasks/"+taskID.String()+"/drafts", nil, map[string]string{"taskID": taskID.String()})
		w := httptest.NewRecorder()

		handler.GetDrafts(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("task not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		taskID := uuid.New()
		mockSvc.On("GetDrafts", mock.Anything, taskID).Return([]*generation.GeneratedCaseDraft{}, generation.ErrTaskNotFound)

		handler := NewGenerationHandler(mockSvc)

		req := createRequestWithChiParams("GET", "/api/v1/generation/tasks/"+taskID.String()+"/drafts", nil, map[string]string{"taskID": taskID.String()})
		w := httptest.NewRecorder()

		handler.GetDrafts(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestConfirmDraftHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful confirm", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		draftID := uuid.New()
		moduleID := uuid.New()
		mockSvc.On("ConfirmDraft", mock.Anything, mock.Anything, mock.Anything).Return(&testcase.TestCase{}, nil)

		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"module_id": moduleID.String(),
		}
		jsonBody, _ := json.Marshal(body)

		req := createRequestWithChiParams("POST", "/api/v1/generation/drafts/"+draftID.String()+"/confirm", jsonBody, map[string]string{"draftID": draftID.String()})
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ConfirmDraft(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("missing user context", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"module_id": uuid.New().String(),
		}
		jsonBody, _ := json.Marshal(body)

		req := createRequestWithChiParams("POST", "/api/v1/generation/drafts/confirm", jsonBody, nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ConfirmDraft(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRejectDraftHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful reject", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		draftID := uuid.New()
		mockSvc.On("RejectDraft", mock.Anything, mock.Anything).Return(nil)

		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"reason":  "duplicate",
			"feedback": "Already exists",
		}
		jsonBody, _ := json.Marshal(body)

		req := createRequestWithChiParams("POST", "/api/v1/generation/drafts/"+draftID.String()+"/reject", jsonBody, map[string]string{"draftID": draftID.String()})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.RejectDraft(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid draft ID", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"reason":  "duplicate",
			"feedback": "Already exists",
		}
		jsonBody, _ := json.Marshal(body)

		req := createRequestWithChiParams("POST", "/api/v1/generation/drafts/invalid-uuid/reject", jsonBody, map[string]string{"id": "invalid-uuid"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.RejectDraft(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestBatchConfirmHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful batch confirm", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		draftIDs := []uuid.UUID{uuid.New(), uuid.New()}
		moduleID := uuid.New()
		mockSvc.On("BatchConfirm", mock.Anything, mock.Anything, mock.Anything).Return(&genservice.BatchConfirmResult{
			SuccessCount: 2,
			FailedCount:  0,
		}, nil)

		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"draft_ids": draftIDs,
			"module_id": moduleID.String(),
		}
		jsonBody, _ := json.Marshal(body)

		req := createRequestWithChiParams("POST", "/api/v1/generation/drafts/batch-confirm", jsonBody, nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.BatchConfirm(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing user context", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"draft_ids": []string{uuid.New().String()},
			"module_id": uuid.New().String(),
		}
		jsonBody, _ := json.Marshal(body)

		req := createRequestWithChiParams("POST", "/api/v1/generation/drafts/batch-confirm", jsonBody, nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.BatchConfirm(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
