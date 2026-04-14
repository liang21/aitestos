// Package handler provides HTTP handlers for the API
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/testcase"
	"github.com/liang21/aitestos/internal/domain/testplan"
	planservice "github.com/liang21/aitestos/internal/service/testplan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockPlanService implements planservice.PlanService for testing
type MockPlanService struct {
	mock.Mock
}

func (m *MockPlanService) CreatePlan(ctx context.Context, req *planservice.CreatePlanRequest, userID uuid.UUID) (*testplan.TestPlan, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*testplan.TestPlan), args.Error(1)
}

func (m *MockPlanService) GetPlan(ctx context.Context, id uuid.UUID) (*planservice.PlanDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*planservice.PlanDetail), args.Error(1)
}

func (m *MockPlanService) ListPlans(ctx context.Context, projectID uuid.UUID, opts planservice.PlanListOptions) ([]*testplan.TestPlan, int64, error) {
	args := m.Called(ctx, projectID, opts)
	return args.Get(0).([]*testplan.TestPlan), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlanService) UpdatePlan(ctx context.Context, id uuid.UUID, req *planservice.UpdatePlanRequest) (*testplan.TestPlan, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*testplan.TestPlan), args.Error(1)
}

func (m *MockPlanService) DeletePlan(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPlanService) UpdatePlanStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPlanService) AddCases(ctx context.Context, planID uuid.UUID, caseIDs []uuid.UUID) error {
	args := m.Called(ctx, planID, caseIDs)
	return args.Error(0)
}

func (m *MockPlanService) RemoveCase(ctx context.Context, planID uuid.UUID, caseID uuid.UUID) error {
	args := m.Called(ctx, planID, caseID)
	return args.Error(0)
}

func (m *MockPlanService) RecordResult(ctx context.Context, req *planservice.RecordResultRequest, userID uuid.UUID) (*testplan.TestResult, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*testplan.TestResult), args.Error(1)
}

func (m *MockPlanService) GetResults(ctx context.Context, planID uuid.UUID) ([]*testplan.TestResult, error) {
	args := m.Called(ctx, planID)
	return args.Get(0).([]*testplan.TestResult), args.Error(1)
}

func (m *MockPlanService) GetResultByCase(ctx context.Context, planID uuid.UUID, caseID uuid.UUID) (*testplan.TestResult, error) {
	args := m.Called(ctx, planID, caseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*testplan.TestResult), args.Error(1)
}

func TestCreatePlanHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful creation", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockPlanService)
		mockSvc.On("CreatePlan", mock.Anything, mock.Anything, mock.Anything).Return(&testplan.TestPlan{}, nil)

		handler := NewTestPlanHandler(mockSvc)
		require.NotNil(t, handler)

		body := map[string]interface{}{
			"project_id": uuid.New().String(),
			"name":       "Test Plan",
			"case_ids":   []string{uuid.New().String()},
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/plans", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreatePlan(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("missing user context", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockPlanService)
		handler := NewTestPlanHandler(mockSvc)

		body := map[string]interface{}{
			"project_id": uuid.New().String(),
			"name":       "Test Plan",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/plans", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreatePlan(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGetPlanHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockPlanService)
		planID := uuid.New()
		mockSvc.On("GetPlan", mock.Anything, planID).Return(&planservice.PlanDetail{}, nil)

		handler := NewTestPlanHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/plans/"+planID.String(), nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", planID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetPlan(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("plan not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockPlanService)
		planID := uuid.New()
		mockSvc.On("GetPlan", mock.Anything, planID).Return(nil, testplan.ErrPlanNotFound)

		handler := NewTestPlanHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/plans/"+planID.String(), nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", planID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetPlan(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAddCasesToPlanHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful add", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockPlanService)
		planID := uuid.New()
		caseIDs := []uuid.UUID{uuid.New(), uuid.New()}
		mockSvc.On("AddCases", mock.Anything, planID, caseIDs).Return(nil)

		handler := NewTestPlanHandler(mockSvc)

		body := map[string]interface{}{
			"case_ids": []string{caseIDs[0].String(), caseIDs[1].String()},
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/plans/"+planID.String()+"/cases", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", planID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.AddCases(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRecordResultHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful record", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockPlanService)
		planID := uuid.New()
		caseID := uuid.New()
		mockSvc.On("RecordResult", mock.Anything, mock.Anything, mock.Anything).Return(&testplan.TestResult{}, nil)

		handler := NewTestPlanHandler(mockSvc)

		body := map[string]interface{}{
			"case_id": caseID.String(),
			"status":  "pass",
			"note":    "Test passed",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/plans/"+planID.String()+"/results", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", planID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		// Set user context
		req = req.WithContext(context.WithValue(req.Context(), userIDContextKey, uuid.New()))
		w := httptest.NewRecorder()

		handler.RecordResult(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid status", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockPlanService)
		planID := uuid.New()
		// Set mock expectation - service should return validation error
		mockSvc.On("RecordResult", mock.Anything, mock.Anything, mock.Anything).Return(nil, testcase.ErrInvalidPriority)

		handler := NewTestPlanHandler(mockSvc)

		body := map[string]interface{}{
			"case_id": uuid.New().String(),
			"status":  "invalid_status",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/plans/"+planID.String()+"/results", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", planID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		// Set user context
		req = req.WithContext(context.WithValue(req.Context(), userIDContextKey, uuid.New()))
		w := httptest.NewRecorder()

		handler.RecordResult(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetResultsHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockPlanService)
		planID := uuid.New()
		mockSvc.On("GetResults", mock.Anything, planID).Return([]*testplan.TestResult{}, nil)

		handler := NewTestPlanHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/plans/"+planID.String()+"/results", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", planID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetResults(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
