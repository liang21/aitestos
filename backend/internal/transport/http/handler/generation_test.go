// Package handler provides HTTP handlers for the API
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

		req := httptest.NewRequest("POST", "/api/v1/generation/tasks", bytes.NewReader(jsonBody))
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

		req := httptest.NewRequest("POST", "/api/v1/generation/tasks", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateTask(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid prompt too short", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"project_id": uuid.New().String(),
			"module_id":  uuid.New().String(),
			"prompt":     "short",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/generation/tasks", bytes.NewReader(jsonBody))
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

		req := httptest.NewRequest("GET", "/api/v1/generation/tasks/"+taskID.String(), nil)
		ctx := context.WithValue(req.Context(), taskIDContextKey, taskID)
		req = req.WithContext(ctx)
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

		req := httptest.NewRequest("GET", "/api/v1/generation/tasks/"+taskID.String(), nil)
		ctx := context.WithValue(req.Context(), taskIDContextKey, taskID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetTask(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetDraftsHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		taskID := uuid.New()
		mockSvc.On("GetDrafts", mock.Anything, taskID).Return([]*generation.GeneratedCaseDraft{}, nil)

		handler := NewGenerationHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/generation/tasks/"+taskID.String()+"/drafts", nil)
		ctx := context.WithValue(req.Context(), taskIDContextKey, taskID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetDrafts(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
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

		req := httptest.NewRequest("POST", "/api/v1/generation/drafts/"+draftID.String()+"/confirm", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), draftIDContextKey, draftID)
		ctx = context.WithValue(ctx, userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ConfirmDraft(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("draft already confirmed", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		draftID := uuid.New()
		mockSvc.On("ConfirmDraft", mock.Anything, mock.Anything, mock.Anything).Return(nil, generation.ErrDraftAlreadyConfirmed)

		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"module_id": uuid.New().String(),
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/generation/drafts/"+draftID.String()+"/confirm", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), draftIDContextKey, draftID)
		ctx = context.WithValue(ctx, userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ConfirmDraft(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
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
			"reason":   "duplicate",
			"feedback": "This case already exists",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/generation/drafts/"+draftID.String()+"/reject", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), draftIDContextKey, draftID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.RejectDraft(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("draft already rejected", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		draftID := uuid.New()
		mockSvc.On("RejectDraft", mock.Anything, mock.Anything).Return(generation.ErrDraftAlreadyRejected)

		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"reason": "irrelevant",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/generation/drafts/"+draftID.String()+"/reject", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), draftIDContextKey, draftID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.RejectDraft(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestBatchConfirmHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful batch confirm", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		mockSvc.On("BatchConfirm", mock.Anything, mock.Anything, mock.Anything).Return(&genservice.BatchConfirmResult{
			SuccessCount: 3,
			FailedCount:  0,
		}, nil)

		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"draft_ids": []string{
				uuid.New().String(),
				uuid.New().String(),
				uuid.New().String(),
			},
			"module_id": uuid.New().String(),
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/generation/drafts/batch-confirm", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.BatchConfirm(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "success_count")
	})

	t.Run("partial success", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockGenerationService)
		mockSvc.On("BatchConfirm", mock.Anything, mock.Anything, mock.Anything).Return(&genservice.BatchConfirmResult{
			SuccessCount: 2,
			FailedCount:  1,
		}, nil)

		handler := NewGenerationHandler(mockSvc)

		body := map[string]interface{}{
			"draft_ids": []string{
				uuid.New().String(),
				uuid.New().String(),
				uuid.New().String(),
			},
			"module_id": uuid.New().String(),
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/generation/drafts/batch-confirm", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.BatchConfirm(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
