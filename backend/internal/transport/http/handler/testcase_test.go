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
	"github.com/liang21/aitestos/internal/domain/testcase"
	caseservice "github.com/liang21/aitestos/internal/service/testcase"
)

// MockCaseService implements caseservice.CaseService for testing
type MockCaseService struct {
	mock.Mock
}

func (m *MockCaseService) CreateCase(ctx context.Context, req *caseservice.CreateCaseRequest, userID uuid.UUID) (*testcase.TestCase, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*testcase.TestCase), args.Error(1)
}

func (m *MockCaseService) UpdateCase(ctx context.Context, id uuid.UUID, req *caseservice.UpdateCaseRequest) (*testcase.TestCase, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*testcase.TestCase), args.Error(1)
}

func (m *MockCaseService) GetCaseDetail(ctx context.Context, id uuid.UUID) (*caseservice.CaseDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*caseservice.CaseDetail), args.Error(1)
}

func (m *MockCaseService) GetCaseByNumber(ctx context.Context, number testcase.CaseNumber) (*caseservice.CaseDetail, error) {
	args := m.Called(ctx, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*caseservice.CaseDetail), args.Error(1)
}

func (m *MockCaseService) ListByModule(ctx context.Context, moduleID uuid.UUID, opts caseservice.CaseListOptions) ([]*testcase.TestCase, int64, error) {
	args := m.Called(ctx, moduleID, opts)
	return args.Get(0).([]*testcase.TestCase), args.Get(1).(int64), args.Error(2)
}

func (m *MockCaseService) ListByProject(ctx context.Context, projectID uuid.UUID, opts caseservice.CaseListOptions) ([]*testcase.TestCase, int64, error) {
	args := m.Called(ctx, projectID, opts)
	return args.Get(0).([]*testcase.TestCase), args.Get(1).(int64), args.Error(2)
}

func (m *MockCaseService) DeleteCase(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCaseService) GenerateCaseNumber(ctx context.Context, moduleID uuid.UUID) (testcase.CaseNumber, error) {
	args := m.Called(ctx, moduleID)
	return args.Get(0).(testcase.CaseNumber), args.Error(1)
}

func TestCreateCaseHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful creation", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		mockSvc.On("CreateCase", mock.Anything, mock.Anything, mock.Anything).Return(&testcase.TestCase{}, nil)

		handler := NewTestCaseHandler(mockSvc)
		require.NotNil(t, handler)

		body := map[string]interface{}{
			"module_id":      uuid.New().String(),
			"title":          "Test Case Title",
			"preconditions":  []string{"Precondition 1"},
			"steps":          []string{"Step 1", "Step 2"},
			"expected":       map[string]interface{}{"result": "success"},
			"case_type":      "functionality",
			"priority":       "P1",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/testcases", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateCase(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("missing user context", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		handler := NewTestCaseHandler(mockSvc)

		body := map[string]interface{}{
			"module_id": uuid.New().String(),
			"title":     "Test Case",
			"steps":     []string{"Step 1"},
			"expected":  map[string]interface{}{},
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/testcases", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateCase(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		handler := NewTestCaseHandler(mockSvc)

		req := httptest.NewRequest("POST", "/api/v1/testcases", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateCase(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("empty steps", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		mockSvc.On("CreateCase", mock.Anything, mock.Anything, mock.Anything).Return(nil, testcase.ErrEmptySteps)

		handler := NewTestCaseHandler(mockSvc)

		body := map[string]interface{}{
			"module_id": uuid.New().String(),
			"title":     "Test Case",
			"steps":     []string{},
			"expected":  map[string]interface{}{},
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/testcases", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateCase(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetCaseHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		caseID := uuid.New()
		mockSvc.On("GetCaseDetail", mock.Anything, caseID).Return(&caseservice.CaseDetail{}, nil)

		handler := NewTestCaseHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/testcases/"+caseID.String(), nil)
		ctx := context.WithValue(req.Context(), caseIDContextKey, caseID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetCase(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("case not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		caseID := uuid.New()
		mockSvc.On("GetCaseDetail", mock.Anything, caseID).Return(nil, testcase.ErrCaseNotFound)

		handler := NewTestCaseHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/testcases/"+caseID.String(), nil)
		ctx := context.WithValue(req.Context(), caseIDContextKey, caseID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetCase(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestListCasesHandler(t *testing.T) {
	t.Parallel()

	t.Run("list by project", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		projectID := uuid.New()
		mockSvc.On("ListByProject", mock.Anything, projectID, mock.Anything).Return([]*testcase.TestCase{}, int64(0), nil)

		handler := NewTestCaseHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/testcases?project_id="+projectID.String(), nil)
		w := httptest.NewRecorder()

		handler.ListCases(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("list by module", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		moduleID := uuid.New()
		mockSvc.On("ListByModule", mock.Anything, moduleID, mock.Anything).Return([]*testcase.TestCase{}, int64(0), nil)

		handler := NewTestCaseHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/testcases?module_id="+moduleID.String(), nil)
		w := httptest.NewRecorder()

		handler.ListCases(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUpdateCaseHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful update", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		caseID := uuid.New()
		mockSvc.On("UpdateCase", mock.Anything, caseID, mock.Anything).Return(&testcase.TestCase{}, nil)

		handler := NewTestCaseHandler(mockSvc)

		body := map[string]interface{}{
			"title": "Updated Title",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/api/v1/testcases/"+caseID.String(), bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), caseIDContextKey, caseID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.UpdateCase(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteCaseHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful delete", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockCaseService)
		caseID := uuid.New()
		mockSvc.On("DeleteCase", mock.Anything, caseID).Return(nil)

		handler := NewTestCaseHandler(mockSvc)

		req := httptest.NewRequest("DELETE", "/api/v1/testcases/"+caseID.String(), nil)
		ctx := context.WithValue(req.Context(), caseIDContextKey, caseID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.DeleteCase(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

const caseIDContextKey contextKey = "case_id"
