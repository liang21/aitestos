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
	"github.com/liang21/aitestos/internal/domain/project"
	projectservice "github.com/liang21/aitestos/internal/service/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockProjectService implements projectservice.ProjectService for testing
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) CreateProject(ctx context.Context, req *projectservice.CreateProjectRequest, userID uuid.UUID) (*project.Project, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*project.Project), args.Error(1)
}

func (m *MockProjectService) GetProject(ctx context.Context, id uuid.UUID) (*projectservice.ProjectDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectservice.ProjectDetail), args.Error(1)
}

func (m *MockProjectService) ListProjects(ctx context.Context, opts projectservice.ListOptions) ([]*project.Project, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]*project.Project), args.Get(1).(int64), args.Error(2)
}

func (m *MockProjectService) UpdateProject(ctx context.Context, id uuid.UUID, req *projectservice.UpdateProjectRequest) (*project.Project, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*project.Project), args.Error(1)
}

func (m *MockProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectService) CreateModule(ctx context.Context, projectID uuid.UUID, req *projectservice.CreateModuleRequest, userID uuid.UUID) (*project.Module, error) {
	args := m.Called(ctx, projectID, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*project.Module), args.Error(1)
}

func (m *MockProjectService) ListModules(ctx context.Context, projectID uuid.UUID) ([]*project.Module, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]*project.Module), args.Error(1)
}

func (m *MockProjectService) GetModule(ctx context.Context, id uuid.UUID) (*project.Module, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*project.Module), args.Error(1)
}

func (m *MockProjectService) DeleteModule(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectService) SetConfig(ctx context.Context, projectID uuid.UUID, key string, value map[string]any) error {
	args := m.Called(ctx, projectID, key, value)
	return args.Error(0)
}

func (m *MockProjectService) GetConfig(ctx context.Context, projectID uuid.UUID, key string) (*project.ProjectConfig, error) {
	args := m.Called(ctx, projectID, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*project.ProjectConfig), args.Error(1)
}

func (m *MockProjectService) ListConfigs(ctx context.Context, projectID uuid.UUID) ([]*project.ProjectConfig, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]*project.ProjectConfig), args.Error(1)
}

func (m *MockProjectService) GetProjectStatistics(ctx context.Context, id uuid.UUID) (*project.ProjectStatistics, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*project.ProjectStatistics), args.Error(1)
}

func (m *MockProjectService) ImportConfigs(ctx context.Context, projectID uuid.UUID, req *projectservice.ImportConfigsRequest) (*projectservice.ImportConfigsResult, error) {
	args := m.Called(ctx, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectservice.ImportConfigsResult), args.Error(1)
}

func (m *MockProjectService) ExportConfigs(ctx context.Context, projectID uuid.UUID) ([]map[string]any, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]any), args.Error(1)
}

func TestCreateProjectHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful creation", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		mockSvc.On("CreateProject", mock.Anything, mock.Anything, mock.Anything).Return(&project.Project{}, nil)

		handler := NewProjectHandler(mockSvc)
		require.NotNil(t, handler)

		body := map[string]string{
			"name":        "Test Project",
			"prefix":      "TEST",
			"description": "A test project",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateProject(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("missing user context", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		handler := NewProjectHandler(mockSvc)

		body := map[string]string{
			"name":   "Test Project",
			"prefix": "TEST",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateProject(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDContextKey, uuid.New())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateProject(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListProjectsHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful list", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		mockSvc.On("ListProjects", mock.Anything, mock.Anything).Return([]*project.Project{}, int64(0), nil)

		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects?offset=0&limit=10", nil)
		w := httptest.NewRecorder()

		handler.ListProjects(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("with pagination", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		mockSvc.On("ListProjects", mock.Anything, mock.Anything).Return([]*project.Project{}, int64(0), nil)

		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects?offset=10&limit=20", nil)
		w := httptest.NewRecorder()

		handler.ListProjects(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetProjectHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		mockSvc.On("GetProject", mock.Anything, projectID).Return(&projectservice.ProjectDetail{}, nil)

		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects/"+projectID.String(), nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", projectID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetProject(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("project not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		mockSvc.On("GetProject", mock.Anything, projectID).Return(nil, project.ErrProjectNotFound)

		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects/"+projectID.String(), nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", projectID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetProject(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCreateModuleHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful creation", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		mockSvc.On("CreateModule", mock.Anything, projectID, mock.Anything, mock.Anything).Return(&project.Module{}, nil)

		handler := NewProjectHandler(mockSvc)

		body := map[string]string{
			"name":         "Test Module",
			"abbreviation": "TMOD",
			"description":  "A test module",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/projects/"+projectID.String()+"/modules", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("projectID", projectID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		// Set user context
		req = req.WithContext(context.WithValue(req.Context(), userIDContextKey, uuid.New()))
		w := httptest.NewRecorder()

		handler.CreateModule(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestSetConfigHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful set", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		mockSvc.On("SetConfig", mock.Anything, projectID, "test-key", mock.Anything).Return(nil)

		handler := NewProjectHandler(mockSvc)

		body := map[string]interface{}{
			"value": map[string]interface{}{"setting": "value"},
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/api/v1/projects/"+projectID.String()+"/configs/test-key", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("projectID", projectID.String())
		rctx.URLParams.Add("key", "test-key")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.SetConfig(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetProjectStatisticsHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		mockSvc.On("GetProjectStatistics", mock.Anything, projectID).Return(&project.ProjectStatistics{}, nil)

		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects/"+projectID.String()+"/stats", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", projectID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetProjectStatistics(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify response is valid JSON
		var stats project.ProjectStatistics
		err := json.Unmarshal(w.Body.Bytes(), &stats)
		assert.NoError(t, err)
	})

	t.Run("project not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		mockSvc.On("GetProjectStatistics", mock.Anything, projectID).Return(nil, project.ErrProjectNotFound)

		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects/"+projectID.String()+"/stats", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", projectID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetProjectStatistics(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid project ID", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects/invalid-uuid/stats", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetProjectStatistics(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestImportConfigsHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful import", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		result := &projectservice.ImportConfigsResult{Imported: 2, Failed: 0}
		mockSvc.On("ImportConfigs", mock.Anything, projectID, mock.Anything).Return(result, nil)

		handler := NewProjectHandler(mockSvc)

		body := map[string]interface{}{
			"configs": []map[string]interface{}{
				{
					"key":         "llm_config",
					"value":       map[string]interface{}{"model": "deepseek-chat"},
					"description": "LLM configuration",
				},
				{
					"key":         "rag_config",
					"value":       map[string]interface{}{"enabled": true},
					"description": "RAG configuration",
				},
			},
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/projects/"+projectID.String()+"/configs/import", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", projectID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.ImportConfigs(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify response is valid JSON
		var resp projectservice.ImportConfigsResult
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 2, resp.Imported)
	})

	t.Run("project not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		mockSvc.On("ImportConfigs", mock.Anything, projectID, mock.Anything).Return(nil, project.ErrProjectNotFound)

		handler := NewProjectHandler(mockSvc)

		body := map[string]interface{}{
			"configs": []map[string]interface{}{
				{
					"key":   "test_config",
					"value": map[string]interface{}{"test": true},
				},
			},
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/projects/"+projectID.String()+"/configs/import", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", projectID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.ImportConfigs(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("POST", "/api/v1/projects/"+uuid.New().String()+"/configs/import", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", uuid.New().String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.ImportConfigs(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestExportConfigsHandler(t *testing.T) {
	t.Parallel()

	t.Run("successful export", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		configs := []map[string]any{
			{
				"key":         "llm_config",
				"value":       map[string]any{"model": "deepseek-chat"},
				"description": "LLM configuration",
			},
			{
				"key":         "rag_config",
				"value":       map[string]any{"enabled": true},
				"description": "RAG configuration",
			},
		}
		mockSvc.On("ExportConfigs", mock.Anything, projectID).Return(configs, nil)

		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects/"+projectID.String()+"/configs/export", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", projectID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.ExportConfigs(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify response is valid JSON array
		var resp []map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
	})

	t.Run("project not found", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		projectID := uuid.New()
		mockSvc.On("ExportConfigs", mock.Anything, projectID).Return(nil, project.ErrProjectNotFound)

		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects/"+projectID.String()+"/configs/export", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", projectID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.ExportConfigs(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid project ID", func(t *testing.T) {
		t.Parallel()

		mockSvc := new(MockProjectService)
		handler := NewProjectHandler(mockSvc)

		req := httptest.NewRequest("GET", "/api/v1/projects/invalid-uuid/configs/export", nil)
		// Set chi URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.ExportConfigs(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
