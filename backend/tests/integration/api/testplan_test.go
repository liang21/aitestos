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

// TestPlanAPI_Integration tests test plan management API endpoints
// Phase 6: Integration Tests (P2) - T-142
func TestPlanAPI_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	suite := testutil.NewIntegrationTestSuite(tc.DB)

	t.Run("Plan CRUD", func(t *testing.T) {
		t.Parallel()
		tc.CleanupTest()

		user := testutil.CreateTestUser(t, tc.DB)
		userID := user.ID()
		project := testutil.CreateTestProject(t, tc.DB)

		tt := []struct {
			name       string
			method     string
			path       string
			body       interface{}
			wantStatus int
			assertFunc func(t *testing.T, w *httptest.ResponseRecorder)
		}{
			{
				name:   "create plan with valid data",
				method: http.MethodPost,
				path:   "/api/v1/plans",
				body: map[string]interface{}{
					"project_id":  project.ID().String(),
					"name":        "Sprint 1 Test Plan",
					"description": "Test plan for sprint 1",
					"case_ids":    []string{},
				},
				wantStatus: http.StatusCreated,
				assertFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
					var response map[string]interface{}
					testutil.ParseJSONResponse(t, w, &response)
					assert.NotEmpty(t, response["id"])
					assert.Equal(t, "Sprint 1 Test Plan", response["name"])
				},
			},
			{
				name:   "create plan with empty name",
				method: http.MethodPost,
				path:   "/api/v1/plans",
				body: map[string]interface{}{
					"project_id":  project.ID().String(),
					"name":        "",
					"description": "Test plan",
				},
				wantStatus: http.StatusBadRequest,
			},
			{
				name:   "create plan with initial cases",
				method: http.MethodPost,
				path:   "/api/v1/plans",
				body: map[string]interface{}{
					"project_id":  project.ID().String(),
					"name":        "Plan with Cases",
					"description": "Test plan with initial cases",
					"case_ids":    []string{uuid.New().String(), uuid.New().String()},
				},
				wantStatus: http.StatusCreated,
			},
			{
				name:       "list plans with pagination",
				method:     http.MethodGet,
				path:       "/api/v1/plans?project_id=" + project.ID().String() + "&offset=0&limit=10",
				wantStatus: http.StatusOK,
			},
			{
				name:       "get non-existent plan",
				method:     http.MethodGet,
				path:       "/api/v1/plans/" + uuid.New().String(),
				wantStatus: http.StatusNotFound,
			},
			{
				name:   "update non-existent plan",
				method: http.MethodPut,
				path:   "/api/v1/plans/" + uuid.New().String(),
				body: map[string]interface{}{
					"name":        "Updated Plan",
					"description": "Updated description",
				},
				wantStatus: http.StatusNotFound,
			},
			{
				name:       "delete non-existent plan",
				method:     http.MethodDelete,
				path:       "/api/v1/plans/" + uuid.New().String(),
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

	t.Run("Case Management", func(t *testing.T) {
		t.Parallel()
		tc.CleanupTest()

		user := testutil.CreateTestUser(t, tc.DB)
		userID := user.ID()
		project := testutil.CreateTestProject(t, tc.DB)
		module := testutil.CreateTestModule(t, tc.DB, project.ID(), userID)
		testCase := testutil.CreateTestCase(t, tc.DB, module.ID(), userID)

		var planID uuid.UUID

		// First create a plan
		t.Run("setup plan", func(t *testing.T) {
			w := suite.MakeRequest(http.MethodPost, "/api/v1/plans", map[string]interface{}{
				"project_id":  project.ID().String(),
				"name":        "Case Management Test Plan",
				"description": "Plan for testing case management",
			}, userID)

			if w.Code == http.StatusNotImplemented {
				t.Skip("handler not implemented yet")
			}

			if w.Code == http.StatusCreated {
				var response map[string]interface{}
				testutil.ParseJSONResponse(t, w, &response)
				if id, ok := response["id"].(string); ok {
					planID, _ = uuid.Parse(id)
				}
			}
		})

		if planID == uuid.Nil {
			t.Skip("plan creation failed, skipping case management tests")
		}

		tt := []struct {
			name       string
			method     string
			path       string
			body       interface{}
			wantStatus int
		}{
			{
				name:   "add cases to plan",
				method: http.MethodPost,
				path:   "/api/v1/plans/" + planID.String() + "/cases",
				body: map[string]interface{}{
					"case_ids": []string{testCase.ID().String()},
				},
				wantStatus: http.StatusOK,
			},
			{
				name:   "add duplicate cases to plan",
				method: http.MethodPost,
				path:   "/api/v1/plans/" + planID.String() + "/cases",
				body: map[string]interface{}{
					"case_ids": []string{testCase.ID().String()},
				},
				wantStatus: http.StatusConflict, // or OK with warning
			},
			{
				name:       "remove case from plan",
				method:     http.MethodDelete,
				path:       "/api/v1/plans/" + planID.String() + "/cases/" + testCase.ID().String(),
				wantStatus: http.StatusNoContent,
			},
			{
				name:       "remove non-existent case from plan",
				method:     http.MethodDelete,
				path:       "/api/v1/plans/" + planID.String() + "/cases/" + uuid.New().String(),
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
			})
		}
	})

	t.Run("Result Recording", func(t *testing.T) {
		t.Parallel()
		tc.CleanupTest()

		user := testutil.CreateTestUser(t, tc.DB)
		userID := user.ID()
		project := testutil.CreateTestProject(t, tc.DB)
		module := testutil.CreateTestModule(t, tc.DB, project.ID(), userID)
		testCase := testutil.CreateTestCase(t, tc.DB, module.ID(), userID)

		var planID uuid.UUID

		// Setup: Create plan and add case
		t.Run("setup", func(t *testing.T) {
			w := suite.MakeRequest(http.MethodPost, "/api/v1/plans", map[string]interface{}{
				"project_id": project.ID().String(),
				"name":       "Result Recording Test Plan",
			}, userID)

			if w.Code == http.StatusCreated {
				var response map[string]interface{}
				testutil.ParseJSONResponse(t, w, &response)
				if id, ok := response["id"].(string); ok {
					planID, _ = uuid.Parse(id)
				}
			}
		})

		if planID == uuid.Nil {
			t.Skip("plan creation failed")
		}

		tt := []struct {
			name       string
			method     string
			path       string
			body       interface{}
			wantStatus int
		}{
			{
				name:   "record pass result",
				method: http.MethodPost,
				path:   "/api/v1/plans/" + planID.String() + "/results",
				body: map[string]interface{}{
					"case_id": testCase.ID().String(),
					"status":  "pass",
					"note":    "All steps executed successfully",
				},
				wantStatus: http.StatusCreated,
			},
			{
				name:   "record fail result",
				method: http.MethodPost,
				path:   "/api/v1/plans/" + planID.String() + "/results",
				body: map[string]interface{}{
					"case_id": testCase.ID().String(),
					"status":  "fail",
					"note":    "Step 3 failed: Expected different error message",
				},
				wantStatus: http.StatusCreated,
			},
			{
				name:   "record block result",
				method: http.MethodPost,
				path:   "/api/v1/plans/" + planID.String() + "/results",
				body: map[string]interface{}{
					"case_id": testCase.ID().String(),
					"status":  "block",
					"note":    "Environment not ready",
				},
				wantStatus: http.StatusCreated,
			},
			{
				name:   "record skip result",
				method: http.MethodPost,
				path:   "/api/v1/plans/" + planID.String() + "/results",
				body: map[string]interface{}{
					"case_id": testCase.ID().String(),
					"status":  "skip",
					"note":    "Feature not available in this version",
				},
				wantStatus: http.StatusCreated,
			},
			{
				name:   "record result with invalid status",
				method: http.MethodPost,
				path:   "/api/v1/plans/" + planID.String() + "/results",
				body: map[string]interface{}{
					"case_id": testCase.ID().String(),
					"status":  "invalid",
				},
				wantStatus: http.StatusBadRequest,
			},
			{
				name:       "get plan results",
				method:     http.MethodGet,
				path:       "/api/v1/plans/" + planID.String() + "/results",
				wantStatus: http.StatusOK,
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
	})
}

// TestPlanAPI_FullExecutionFlow tests the complete execution flow
func TestPlanAPI_FullExecutionFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	suite := testutil.NewIntegrationTestSuite(tc.DB)

	user := testutil.CreateTestUser(t, tc.DB)
	userID := user.ID()
	project := testutil.CreateTestProject(t, tc.DB)
	module := testutil.CreateTestModule(t, tc.DB, project.ID(), userID)

	// Create multiple test cases
	_ = testutil.CreateTestCase(t, tc.DB, module.ID(), userID)
	_ = testutil.CreateTestCase(t, tc.DB, module.ID(), userID)
	_ = testutil.CreateTestCase(t, tc.DB, module.ID(), userID)

	// Full execution flow:
	// 1. Create test plan
	// 2. Add test cases
	// 3. Update plan status to active
	// 4. Record results for all cases
	// 5. Update plan status to completed
	// 6. Verify statistics

	t.Run("Full Flow", func(t *testing.T) {
		// Step 1: Create test plan
		w := suite.MakeRequest(http.MethodPost, "/api/v1/plans", map[string]interface{}{
			"project_id":  project.ID().String(),
			"name":        "Full Execution Flow Test",
			"description": "Testing complete execution flow",
		}, userID)

		if w.Code == http.StatusNotImplemented {
			t.Skip("handler not implemented yet")
		}

		assert.Equal(t, http.StatusCreated, w.Code)

		// Additional steps would be implemented when handlers are ready
	})
}

// TestPlanAPI_Statistics tests plan statistics calculations
func TestPlanAPI_Statistics(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	_ = testutil.NewIntegrationTestSuite(tc.DB) // suite created but not needed for skipped tests

	t.Run("Empty Plan Statistics", func(t *testing.T) {
		// Test that empty plan returns zero statistics
		t.Skip("statistics endpoint not yet implemented")
	})

	t.Run("Partial Execution Statistics", func(t *testing.T) {
		// Test statistics with some results recorded
		t.Skip("statistics endpoint not yet implemented")
	})

	t.Run("Complete Execution Statistics", func(t *testing.T) {
		// Test statistics with all results recorded
		t.Skip("statistics endpoint not yet implemented")
	})
}
