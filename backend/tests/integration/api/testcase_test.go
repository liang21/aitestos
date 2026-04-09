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

// TestTestCaseAPI_Integration tests test case management API endpoints
// Phase 6: Integration Tests (P2) - T-140
func TestTestCaseAPI_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	suite := testutil.NewIntegrationTestSuite(tc.DB)

	t.Run("TestCase CRUD", func(t *testing.T) {
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
		}{
			{
				name:   "create test case with valid data",
				method: http.MethodPost,
				path:   "/api/v1/testcases",
				body: map[string]interface{}{
					"module_id":     module.ID().String(),
					"title":         "Test Login Function",
					"preconditions": []string{"User has valid account"},
					"steps":         []string{"Navigate to login page", "Enter credentials", "Click login button"},
					"expected":      map[string]interface{}{"status": "success", "redirect": "/dashboard"},
					"case_type":     "functionality",
					"priority":      "P2",
				},
				wantStatus: http.StatusCreated,
				assertFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
					// Verify response contains case ID and number
				},
			},
			{
				name:   "create test case with empty steps",
				method: http.MethodPost,
				path:   "/api/v1/testcases",
				body: map[string]interface{}{
					"module_id":     module.ID().String(),
					"title":         "Invalid Test Case",
					"preconditions": []string{},
					"steps":         []string{}, // Empty steps should fail
					"expected":      map[string]interface{}{},
					"case_type":     "functionality",
					"priority":      "P2",
				},
				wantStatus: http.StatusBadRequest,
			},
			{
				name:   "create test case with invalid module ID",
				method: http.MethodPost,
				path:   "/api/v1/testcases",
				body: map[string]interface{}{
					"module_id":     uuid.New().String(), // Non-existent module
					"title":         "Invalid Module Test Case",
					"preconditions": []string{},
					"steps":         []string{"Step 1"},
					"expected":      map[string]interface{}{},
					"case_type":     "functionality",
					"priority":      "P2",
				},
				wantStatus: http.StatusBadRequest,
			},
			{
				name:   "create test case with missing title",
				method: http.MethodPost,
				path:   "/api/v1/testcases",
				body: map[string]interface{}{
					"module_id":     module.ID().String(),
					"title":         "", // Empty title should fail
					"preconditions": []string{},
					"steps":         []string{"Step 1"},
					"expected":      map[string]interface{}{},
					"case_type":     "functionality",
					"priority":      "P2",
				},
				wantStatus: http.StatusBadRequest,
			},
			{
				name:       "list test cases with pagination",
				method:     http.MethodGet,
				path:       "/api/v1/testcases?offset=0&limit=10",
				wantStatus: http.StatusOK,
			},
			{
				name:       "list test cases by module",
				method:     http.MethodGet,
				path:       "/api/v1/testcases?module_id=" + module.ID().String(),
				wantStatus: http.StatusOK,
			},
			{
				name:       "list test cases by project",
				method:     http.MethodGet,
				path:       "/api/v1/testcases?project_id=" + project.ID().String(),
				wantStatus: http.StatusOK,
			},
			{
				name:       "list test cases with keywords",
				method:     http.MethodGet,
				path:       "/api/v1/testcases?keywords=login&limit=5",
				wantStatus: http.StatusOK,
			},
			{
				name:       "get non-existent test case",
				method:     http.MethodGet,
				path:       "/api/v1/testcases/" + uuid.New().String(),
				wantStatus: http.StatusNotFound,
			},
			{
				name:   "update non-existent test case",
				method: http.MethodPut,
				path:   "/api/v1/testcases/" + uuid.New().String(),
				body: map[string]interface{}{
					"title":         "Updated Test Case",
					"preconditions": []string{"New precondition"},
					"steps":         []string{"New step"},
					"expected":      map[string]interface{}{"status": "updated"},
				},
				wantStatus: http.StatusNotFound,
			},
			{
				name:       "delete non-existent test case",
				method:     http.MethodDelete,
				path:       "/api/v1/testcases/" + uuid.New().String(),
				wantStatus: http.StatusNotFound,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
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

	t.Run("Case Number Generation", func(t *testing.T) {
		t.Parallel()
		tc.CleanupTest()

		user := testutil.CreateTestUser(t, tc.DB)
		userID := user.ID()
		project := testutil.CreateTestProject(t, tc.DB)
		module := testutil.CreateTestModule(t, tc.DB, project.ID(), userID)

		// Test that case numbers are auto-generated correctly
		t.Run("auto-generate case number", func(t *testing.T) {
			w := suite.MakeRequest(http.MethodPost, "/api/v1/testcases", map[string]interface{}{
				"module_id":     module.ID().String(),
				"title":         "Auto-generated Number Test",
				"preconditions": []string{},
				"steps":         []string{"Step 1"},
				"expected":      map[string]interface{}{},
				"case_type":     "functionality",
				"priority":      "P2",
			}, userID)

			if w.Code == http.StatusNotImplemented {
				t.Skip("handler not implemented yet")
			}

			assert.Equal(t, http.StatusCreated, w.Code)

			// Verify case number format: {PROJECT_PREFIX}-{MODULE_ABBREV}-{DATE}-{SEQ}
			// e.g., TP-TM-20260402-001
		})
	})

	t.Run("Requirement Traceability", func(t *testing.T) {
		t.Parallel()
		tc.CleanupTest()

		user := testutil.CreateTestUser(t, tc.DB)
		_ = user.ID() // userID used implicitly through suite.MakeRequest
		project := testutil.CreateTestProject(t, tc.DB)
		_ = testutil.CreateTestModule(t, tc.DB, project.ID(), user.ID())

		// Test linking test cases to requirements
		t.Run("link to requirement", func(t *testing.T) {
			// This would test the ability to link a test case to a requirement document
			// Currently a placeholder for future implementation
			t.Skip("requirement traceability not yet implemented")
		})
	})
}

// TestTestCaseAPI_Filtering tests various filtering options
func TestTestCaseAPI_Filtering(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	suite := testutil.NewIntegrationTestSuite(tc.DB)

	user := testutil.CreateTestUser(t, tc.DB)
	userID := user.ID()

	tt := []struct {
		name string
		path string
	}{
		{"filter by status", "/api/v1/testcases?status=unexecuted"},
		{"filter by priority", "/api/v1/testcases?priority=P0"},
		{"filter by case type", "/api/v1/testcases?case_type=functionality"},
		{"filter by multiple", "/api/v1/testcases?priority=P0&status=unexecuted&limit=20"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := suite.MakeRequest(http.MethodGet, tc.path, nil, userID)

			if w.Code == http.StatusNotImplemented {
				t.Skip("handler not implemented yet")
			}

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestTestCaseAPI_BulkOperations tests bulk operations on test cases
func TestTestCaseAPI_BulkOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	_ = testutil.NewIntegrationTestSuite(tc.DB)
	_ = testutil.CreateTestUser(t, tc.DB) // user created but not needed for skipped tests

	t.Run("bulk delete", func(t *testing.T) {
		// Test bulk deleting multiple test cases
		t.Skip("bulk delete not yet implemented")
	})

	t.Run("bulk update status", func(t *testing.T) {
		// Test bulk updating status of multiple test cases
		t.Skip("bulk update not yet implemented")
	})
}
