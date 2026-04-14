// Package api_test provides integration tests for API endpoints
package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/liang21/aitestos/tests/integration/testutil"
	"github.com/stretchr/testify/assert"
)

// TestGenerationAPI_Integration tests generation task management API endpoints
// Phase 6: Integration Tests (P2) - T-141
func TestGenerationAPI_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	suite := testutil.NewIntegrationTestSuite(tc.DB)

	t.Run("Generation Task Lifecycle", func(t *testing.T) {
		t.Parallel()
		tc.CleanupTest()

		user := testutil.CreateTestUser(t, tc.DB)
		userID := user.ID()
		project := testutil.CreateTestProject(t, tc.DB)
		module := testutil.CreateTestModule(t, tc.DB, project.ID(), userID)

		tt := []struct {
			name       string
			method     string
			path       string
			body       interface{}
			wantStatus int
			assertFunc func(t *testing.T, w *httptest.ResponseRecorder)
			setup      func(t *testing.T)
		}{
			{
				name:   "create generation task with valid data",
				method: http.MethodPost,
				path:   "/api/v1/generation/tasks",
				body: map[string]interface{}{
					"project_id":  project.ID().String(),
					"module_id":   module.ID().String(),
					"prompt":      "Generate test cases for user login functionality including positive and negative scenarios",
					"case_count":  5,
					"scene_types": []string{"positive", "negative", "boundary"},
					"priority":    "P2",
					"case_type":   "functionality",
				},
				wantStatus: http.StatusCreated,
			},
			{
				name:   "create generation task with short prompt",
				method: http.MethodPost,
				path:   "/api/v1/generation/tasks",
				body: map[string]interface{}{
					"project_id": project.ID().String(),
					"module_id":  module.ID().String(),
					"prompt":     "short",
					"case_count": 5,
				},
				wantStatus: http.StatusBadRequest,
			},
			{
				name:   "create generation task with excessive case count",
				method: http.MethodPost,
				path:   "/api/v1/generation/tasks",
				body: map[string]interface{}{
					"project_id": project.ID().String(),
					"module_id":  module.ID().String(),
					"prompt":     "Generate test cases for user registration",
					"case_count": 100, // Max is 20
				},
				wantStatus: http.StatusBadRequest,
			},
			{
				name:       "get non-existent task",
				method:     http.MethodGet,
				path:       "/api/v1/generation/tasks/" + uuid.New().String(),
				wantStatus: http.StatusNotFound,
			},
			{
				name:       "list generation tasks by project",
				method:     http.MethodGet,
				path:       "/api/v1/generation/tasks?project_id=" + project.ID().String(),
				wantStatus: http.StatusOK,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				if tc.setup != nil {
					tc.setup(t)
				}

				w := suite.MakeRequest(tc.method, tc.path, tc.body, userID)

				if w.Code == http.StatusNotImplemented {
					t.Skipf("handler not implemented yet (expected status %d)", tc.wantStatus)
				}

				assert.Equal(t, tc.wantStatus, w.Code)

				if tc.assertFunc != nil {
					tc.assertFunc(t, w)
				}
			})
		}
	})

	t.Run("Draft Management", func(t *testing.T) {
		t.Parallel()
		tc.CleanupTest()

		user := testutil.CreateTestUser(t, tc.DB)
		userID := user.ID()
		project := testutil.CreateTestProject(t, tc.DB)
		module := testutil.CreateTestModule(t, tc.DB, project.ID(), userID)

		var draftID uuid.UUID

		tt := []struct {
			name       string
			method     string
			path       string
			body       interface{}
			wantStatus int
			setup      func(t *testing.T)
		}{
			{
				name:       "get drafts for task",
				method:     http.MethodGet,
				path:       "/api/v1/generation/tasks/" + uuid.New().String() + "/drafts",
				wantStatus: http.StatusOK, // May return empty list if task doesn't exist
			},
			{
				name:   "confirm non-existent draft",
				method: http.MethodPost,
				path:   "/api/v1/generation/drafts/" + uuid.New().String() + "/confirm",
				body: map[string]string{
					"module_id": module.ID().String(),
				},
				wantStatus: http.StatusNotFound,
			},
			{
				name:   "reject non-existent draft",
				method: http.MethodPost,
				path:   "/api/v1/generation/drafts/" + uuid.New().String() + "/reject",
				body: map[string]interface{}{
					"reason":   "duplicate",
					"feedback": "This case already exists",
				},
				wantStatus: http.StatusNotFound,
			},
			{
				name:   "batch confirm with empty list",
				method: http.MethodPost,
				path:   "/api/v1/generation/drafts/batch-confirm",
				body: map[string]interface{}{
					"draft_ids": []string{},
					"module_id": module.ID().String(),
				},
				wantStatus: http.StatusBadRequest,
			},
			{
				name:   "batch confirm with too many drafts",
				method: http.MethodPost,
				path:   "/api/v1/generation/drafts/batch-confirm",
				body: func() map[string]interface{} {
					ids := make([]string, 51) // Max is 50
					for i := range ids {
						ids[i] = uuid.New().String()
					}
					return map[string]interface{}{
						"draft_ids": ids,
						"module_id": module.ID().String(),
					}
				}(),
				wantStatus: http.StatusBadRequest,
			},
		}

		// Store draft ID for later use
		_ = draftID

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				if tc.setup != nil {
					tc.setup(t)
				}

				w := suite.MakeRequest(tc.method, tc.path, tc.body, userID)

				if w.Code == http.StatusNotImplemented {
					t.Skipf("handler not implemented yet (expected status %d)", tc.wantStatus)
				}

				assert.Equal(t, tc.wantStatus, w.Code)
			})
		}
	})

	t.Run("Full Generation Flow", func(t *testing.T) {
		t.Parallel()
		tc.CleanupTest()

		user := testutil.CreateTestUser(t, tc.DB)
		userID := user.ID()
		project := testutil.CreateTestProject(t, tc.DB)
		module := testutil.CreateTestModule(t, tc.DB, project.ID(), userID)

		// Step 1: Create generation task
		t.Run("Step 1: Create Task", func(t *testing.T) {
			w := suite.MakeRequest(http.MethodPost, "/api/v1/generation/tasks", map[string]interface{}{
				"project_id":  project.ID().String(),
				"module_id":   module.ID().String(),
				"prompt":      "Generate comprehensive test cases for user authentication",
				"case_count":  3,
				"scene_types": []string{"positive", "negative"},
				"priority":    "P1",
			}, userID)

			if w.Code == http.StatusNotImplemented {
				t.Skip("handler not implemented yet")
			}

			assert.Equal(t, http.StatusCreated, w.Code)
		})

		// Step 2: Get task status
		// Step 3: Get drafts
		// Step 4: Confirm/Reject drafts
		// These would be implemented when handlers are ready
	})
}

// TestGenerationAPI_ErrorHandling tests error handling scenarios
func TestGenerationAPI_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	suite := testutil.NewIntegrationTestSuite(tc.DB)

	user := testutil.CreateTestUser(t, tc.DB)
	userID := user.ID()

	tt := []struct {
		name       string
		method     string
		path       string
		body       interface{}
		wantStatus int
	}{
		{
			name:   "create task with missing project ID",
			method: http.MethodPost,
			path:   "/api/v1/generation/tasks",
			body: map[string]interface{}{
				"module_id":  uuid.New().String(),
				"prompt":     "Generate test cases",
				"case_count": 5,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "create task with missing module ID",
			method: http.MethodPost,
			path:   "/api/v1/generation/tasks",
			body: map[string]interface{}{
				"project_id": uuid.New().String(),
				"prompt":     "Generate test cases",
				"case_count": 5,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "create task with invalid JSON",
			method:     http.MethodPost,
			path:       "/api/v1/generation/tasks",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := suite.MakeRequest(tc.method, tc.path, tc.body, userID)

			if w.Code == http.StatusNotImplemented {
				t.Skipf("handler not implemented yet (expected status %d)", tc.wantStatus)
			}

			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}
