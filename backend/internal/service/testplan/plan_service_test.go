// Package testplan provides test plan management services
package testplan

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/testcase"
	"github.com/liang21/aitestos/internal/domain/testplan"
)

// MockTestPlanRepository implements testplan.TestPlanRepository for testing
type MockTestPlanRepository struct {
	plans   map[uuid.UUID]*testplan.TestPlan
	saveErr error
	findErr error
}

func NewMockTestPlanRepository() *MockTestPlanRepository {
	return &MockTestPlanRepository{
		plans: make(map[uuid.UUID]*testplan.TestPlan),
	}
}

func (m *MockTestPlanRepository) Save(ctx context.Context, p *testplan.TestPlan) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.plans[p.ID()] = p
	return nil
}

func (m *MockTestPlanRepository) FindByID(ctx context.Context, id uuid.UUID) (*testplan.TestPlan, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	p, ok := m.plans[id]
	if !ok {
		return nil, testplan.ErrPlanNotFound
	}
	return p, nil
}

func (m *MockTestPlanRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts testplan.QueryOptions) ([]*testplan.TestPlan, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	plans := make([]*testplan.TestPlan, 0)
	for _, p := range m.plans {
		if p.ProjectID() == projectID {
			plans = append(plans, p)
		}
	}
	return plans, nil
}

func (m *MockTestPlanRepository) FindAll(ctx context.Context, opts testplan.QueryOptions) ([]*testplan.TestPlan, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	plans := make([]*testplan.TestPlan, 0, len(m.plans))
	for _, p := range m.plans {
		plans = append(plans, p)
	}
	return plans, nil
}

func (m *MockTestPlanRepository) Update(ctx context.Context, p *testplan.TestPlan) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.plans[p.ID()] = p
	return nil
}

func (m *MockTestPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	if _, ok := m.plans[id]; !ok {
		return testplan.ErrPlanNotFound
	}
	delete(m.plans, id)
	return nil
}

func (m *MockTestPlanRepository) FindByStatus(ctx context.Context, status testplan.PlanStatus, opts testplan.QueryOptions) ([]*testplan.TestPlan, error) {
	return []*testplan.TestPlan{}, nil
}

func (m *MockTestPlanRepository) AddCase(ctx context.Context, planID, caseID uuid.UUID) error {
	return nil
}

func (m *MockTestPlanRepository) RemoveCase(ctx context.Context, planID, caseID uuid.UUID) error {
	return nil
}

func (m *MockTestPlanRepository) GetCaseIDs(ctx context.Context, planID uuid.UUID) ([]uuid.UUID, error) {
	return []uuid.UUID{}, nil
}

func (m *MockTestPlanRepository) UpdateStatus(ctx context.Context, planID uuid.UUID, status testplan.PlanStatus) error {
	return nil
}

// MockTestResultRepository implements testplan.TestResultRepository for testing
type MockTestResultRepository struct {
	results   map[uuid.UUID]*testplan.TestResult
	planIndex map[uuid.UUID][]*testplan.TestResult
	saveErr   error
	findErr   error
}

func NewMockTestResultRepository() *MockTestResultRepository {
	return &MockTestResultRepository{
		results:   make(map[uuid.UUID]*testplan.TestResult),
		planIndex: make(map[uuid.UUID][]*testplan.TestResult),
	}
}

func (m *MockTestResultRepository) Save(ctx context.Context, r *testplan.TestResult) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.results[r.ID()] = r
	m.planIndex[r.PlanID()] = append(m.planIndex[r.PlanID()], r)
	return nil
}

func (m *MockTestResultRepository) FindByID(ctx context.Context, id uuid.UUID) (*testplan.TestResult, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	r, ok := m.results[id]
	if !ok {
		return nil, testplan.ErrResultNotFound
	}
	return r, nil
}

func (m *MockTestResultRepository) FindByPlanID(ctx context.Context, planID uuid.UUID, opts testplan.QueryOptions) ([]*testplan.TestResult, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	results, ok := m.planIndex[planID]
	if !ok {
		return []*testplan.TestResult{}, nil
	}
	return results, nil
}

func (m *MockTestResultRepository) FindByPlanAndCase(ctx context.Context, planID uuid.UUID, caseID uuid.UUID) (*testplan.TestResult, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	for _, r := range m.planIndex[planID] {
		if r.CaseID() == caseID {
			return r, nil
		}
	}
	return nil, testplan.ErrResultNotFound
}

func (m *MockTestResultRepository) Update(ctx context.Context, r *testplan.TestResult) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.results[r.ID()] = r
	return nil
}

func (m *MockTestResultRepository) FindByCaseID(ctx context.Context, caseID uuid.UUID, opts testplan.QueryOptions) ([]*testplan.TestResult, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	results := make([]*testplan.TestResult, 0)
	for _, r := range m.results {
		if r.CaseID() == caseID {
			results = append(results, r)
		}
	}
	return results, nil
}

func (m *MockTestResultRepository) FindByExecutorID(ctx context.Context, executorID uuid.UUID, opts testplan.QueryOptions) ([]*testplan.TestResult, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	results := make([]*testplan.TestResult, 0)
	for _, r := range m.results {
		if r.ExecutedBy() == executorID {
			results = append(results, r)
		}
	}
	return results, nil
}

func (m *MockTestResultRepository) DeleteByPlanID(ctx context.Context, planID uuid.UUID) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	delete(m.planIndex, planID)
	return nil
}

func (m *MockTestResultRepository) CountByPlanID(ctx context.Context, planID uuid.UUID) (int64, error) {
	if m.findErr != nil {
		return 0, m.findErr
	}
	return int64(len(m.planIndex[planID])), nil
}

func (m *MockTestResultRepository) CountByStatus(ctx context.Context, planID uuid.UUID) (map[testplan.ResultStatus]int, error) {
	return make(map[testplan.ResultStatus]int), nil
}

func (m *MockTestResultRepository) FindLatestByCaseID(ctx context.Context, caseID uuid.UUID) (*testplan.TestResult, error) {
	return nil, nil
}

func (m *MockTestResultRepository) FindByPlanIDAndCaseID(ctx context.Context, planID, caseID uuid.UUID) ([]*testplan.TestResult, error) {
	return nil, nil
}

// MockTestCaseRepository for testing
type MockTestCaseRepository struct {
	cases map[uuid.UUID]*testcase.TestCase
}

func NewMockTestCaseRepoForPlan() *MockTestCaseRepository {
	return &MockTestCaseRepository{
		cases: make(map[uuid.UUID]*testcase.TestCase),
	}
}

func (m *MockTestCaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*testcase.TestCase, error) {
	tc, ok := m.cases[id]
	if !ok {
		return nil, testcase.ErrCaseNotFound
	}
	return tc, nil
}

func (m *MockTestCaseRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*testcase.TestCase, error) {
	cases := make([]*testcase.TestCase, 0)
	for _, id := range ids {
		if tc, ok := m.cases[id]; ok {
			cases = append(cases, tc)
		}
	}
	return cases, nil
}

func (m *MockTestCaseRepository) AddCase(tc *testcase.TestCase) {
	m.cases[tc.ID()] = tc
}

func (m *MockTestCaseRepository) Save(ctx context.Context, tc *testcase.TestCase) error {
	m.cases[tc.ID()] = tc
	return nil
}

func (m *MockTestCaseRepository) FindByNumber(ctx context.Context, number testcase.CaseNumber) (*testcase.TestCase, error) {
	for _, tc := range m.cases {
		if tc.Number() == number {
			return tc, nil
		}
	}
	return nil, testcase.ErrCaseNotFound
}

func (m *MockTestCaseRepository) FindByModuleID(ctx context.Context, moduleID uuid.UUID, opts testcase.QueryOptions) ([]*testcase.TestCase, error) {
	result := make([]*testcase.TestCase, 0)
	for _, tc := range m.cases {
		if tc.ModuleID() == moduleID {
			result = append(result, tc)
		}
	}
	return result, nil
}

func (m *MockTestCaseRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts testcase.QueryOptions) ([]*testcase.TestCase, error) {
	return []*testcase.TestCase{}, nil
}

func (m *MockTestCaseRepository) Update(ctx context.Context, tc *testcase.TestCase) error {
	m.cases[tc.ID()] = tc
	return nil
}

func (m *MockTestCaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.cases, id)
	return nil
}

func (m *MockTestCaseRepository) CountByDate(ctx context.Context, moduleID uuid.UUID, date time.Time) (int64, error) {
	return 0, nil
}

func (m *MockTestCaseRepository) CountByModuleID(ctx context.Context, moduleID uuid.UUID) (int64, error) {
	count := int64(0)
	for _, tc := range m.cases {
		if tc.ModuleID() == moduleID {
			count++
		}
	}
	return count, nil
}

func (m *MockTestCaseRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	return 0, nil
}

// TestPlanService_CreatePlan tests plan creation
func TestPlanService_CreatePlan(t *testing.T) {
	ctx := context.Background()
	planRepo := NewMockTestPlanRepository()
	resultRepo := NewMockTestResultRepository()
	caseRepo := NewMockTestCaseRepoForPlan()
	service := NewPlanService(planRepo, resultRepo, caseRepo)

	projectID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		req     *CreatePlanRequest
		wantErr error
	}{
		{
			name: "successful creation",
			req: &CreatePlanRequest{
				ProjectID:   projectID,
				Name:        "Sprint 1 Test Plan",
				Description: "Test plan for first sprint",
				CaseIDs:     []uuid.UUID{},
			},
			wantErr: nil,
		},
		{
			name: "successful creation with cases",
			req: &CreatePlanRequest{
				ProjectID:   projectID,
				Name:        "Sprint 2 Test Plan",
				Description: "Test plan with initial cases",
				CaseIDs:     []uuid.UUID{uuid.New(), uuid.New()},
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			req: &CreatePlanRequest{
				ProjectID:   projectID,
				Name:        "",
				Description: "Invalid plan",
				CaseIDs:     []uuid.UUID{},
			},
			wantErr: errors.New("create plan: plan name cannot be empty"),
		},
		{
			name: "nil project ID",
			req: &CreatePlanRequest{
				ProjectID:   uuid.Nil,
				Name:        "Orphan Plan",
				Description: "No project",
				CaseIDs:     []uuid.UUID{},
			},
			wantErr: errors.New("create plan: project ID cannot be nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := service.CreatePlan(ctx, tt.req, userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CreatePlan() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) {
					t.Errorf("CreatePlan() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreatePlan() unexpected error: %v", err)
				return
			}

			if plan == nil {
				t.Error("CreatePlan() returned nil plan")
				return
			}

			if plan.Name() != tt.req.Name {
				t.Errorf("CreatePlan() name = %v, want %v", plan.Name(), tt.req.Name)
			}
			if plan.ProjectID() != tt.req.ProjectID {
				t.Errorf("CreatePlan() projectID = %v, want %v", plan.ProjectID(), tt.req.ProjectID)
			}
			if plan.Status() != testplan.StatusDraft {
				t.Errorf("CreatePlan() initial status = %v, want %v", plan.Status(), testplan.StatusDraft)
			}
		})
	}
}

// TestPlanService_GetPlan tests plan retrieval
func TestPlanService_GetPlan(t *testing.T) {
	ctx := context.Background()
	planRepo := NewMockTestPlanRepository()
	resultRepo := NewMockTestResultRepository()
	caseRepo := NewMockTestCaseRepoForPlan()
	service := NewPlanService(planRepo, resultRepo, caseRepo)

	// Create test plan
	projectID := uuid.New()
	userID := uuid.New()
	testPlan, _ := testplan.NewTestPlan(projectID, "Test Plan", "Description", userID)
	planRepo.plans[testPlan.ID()] = testPlan

	tests := []struct {
		name    string
		planID  uuid.UUID
		wantErr error
	}{
		{
			name:    "successful retrieval",
			planID:  testPlan.ID(),
			wantErr: nil,
		},
		{
			name:    "plan not found",
			planID:  uuid.New(),
			wantErr: testplan.ErrPlanNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := service.GetPlan(ctx, tt.planID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetPlan() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetPlan() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetPlan() unexpected error: %v", err)
				return
			}

			if detail == nil {
				t.Error("GetPlan() returned nil detail")
				return
			}

			if detail.ID() != tt.planID {
				t.Errorf("GetPlan() ID = %v, want %v", detail.ID(), tt.planID)
			}
		})
	}
}

// TestPlanService_AddCases tests adding cases to plan
func TestPlanService_AddCases(t *testing.T) {
	ctx := context.Background()
	planRepo := NewMockTestPlanRepository()
	resultRepo := NewMockTestResultRepository()
	caseRepo := NewMockTestCaseRepoForPlan()
	service := NewPlanService(planRepo, resultRepo, caseRepo)

	// Create test project, module and cases
	projectID := uuid.New()
	userID := uuid.New()
	moduleID := uuid.New()

	// Create test cases
	case1 := createTestCase(t, moduleID, userID, "Case 1")
	case2 := createTestCase(t, moduleID, userID, "Case 2")
	caseRepo.AddCase(case1)
	caseRepo.AddCase(case2)

	// Create test plan
	testPlan, _ := testplan.NewTestPlan(projectID, "Test Plan", "Description", userID)
	planRepo.plans[testPlan.ID()] = testPlan

	// Create archived plan
	archivedPlan, _ := testplan.NewTestPlan(projectID, "Archived Plan", "Description", userID)
	_ = archivedPlan.UpdateStatus(testplan.StatusArchived)
	planRepo.plans[archivedPlan.ID()] = archivedPlan

	tests := []struct {
		name    string
		planID  uuid.UUID
		caseIDs []uuid.UUID
		wantErr error
	}{
		{
			name:    "successful add single case",
			planID:  testPlan.ID(),
			caseIDs: []uuid.UUID{case1.ID()},
			wantErr: nil,
		},
		{
			name:    "successful add multiple cases",
			planID:  testPlan.ID(),
			caseIDs: []uuid.UUID{case2.ID()},
			wantErr: nil,
		},
		{
			name:    "plan not found",
			planID:  uuid.New(),
			caseIDs: []uuid.UUID{case1.ID()},
			wantErr: testplan.ErrPlanNotFound,
		},
		{
			name:    "archived plan",
			planID:  archivedPlan.ID(),
			caseIDs: []uuid.UUID{case1.ID()},
			wantErr: errors.New("cannot add case to archived plan"),
		},
		{
			name:    "case not found",
			planID:  testPlan.ID(),
			caseIDs: []uuid.UUID{uuid.New()},
			wantErr: testcase.ErrCaseNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AddCases(ctx, tt.planID, tt.caseIDs)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("AddCases() expected error %v, got nil", tt.wantErr)
					return
				}
				// Support both exact match and substring match for wrapped errors
				if err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) && !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Errorf("AddCases() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("AddCases() unexpected error: %v", err)
				return
			}

			// Verify cases were added
			plan, _ := planRepo.FindByID(ctx, tt.planID)
			for _, caseID := range tt.caseIDs {
				if !plan.HasCase(caseID) {
					t.Errorf("AddCases() case %v not in plan", caseID)
				}
			}
		})
	}
}

// TestPlanService_RecordResult tests result recording
func TestPlanService_RecordResult(t *testing.T) {
	ctx := context.Background()
	planRepo := NewMockTestPlanRepository()
	resultRepo := NewMockTestResultRepository()
	caseRepo := NewMockTestCaseRepoForPlan()
	service := NewPlanService(planRepo, resultRepo, caseRepo)

	// Create test project, module and case
	projectID := uuid.New()
	userID := uuid.New()
	moduleID := uuid.New()

	testCase := createTestCase(t, moduleID, userID, "Test Case")
	caseRepo.AddCase(testCase)

	// Create test plan with the case
	testPlan, _ := testplan.NewTestPlan(projectID, "Test Plan", "Description", userID)
	_ = testPlan.AddCase(testCase.ID())
	planRepo.plans[testPlan.ID()] = testPlan

	tests := []struct {
		name    string
		req     *RecordResultRequest
		userID  uuid.UUID
		wantErr error
	}{
		{
			name: "successful pass result",
			req: &RecordResultRequest{
				PlanID: testPlan.ID(),
				CaseID: testCase.ID(),
				Status: "pass",
				Note:   "Test passed successfully",
			},
			userID:  userID,
			wantErr: nil,
		},
		{
			name: "successful fail result",
			req: &RecordResultRequest{
				PlanID: testPlan.ID(),
				CaseID: testCase.ID(),
				Status: "fail",
				Note:   "Found a bug in login flow",
			},
			userID:  userID,
			wantErr: nil,
		},
		{
			name: "invalid status",
			req: &RecordResultRequest{
				PlanID: testPlan.ID(),
				CaseID: testCase.ID(),
				Status: "invalid",
				Note:   "",
			},
			userID:  userID,
			wantErr: errors.New("invalid result status"),
		},
		{
			name: "plan not found",
			req: &RecordResultRequest{
				PlanID: uuid.New(),
				CaseID: testCase.ID(),
				Status: "pass",
				Note:   "",
			},
			userID:  userID,
			wantErr: testplan.ErrPlanNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.RecordResult(ctx, tt.req, tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("RecordResult() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) {
					t.Errorf("RecordResult() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("RecordResult() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("RecordResult() returned nil result")
				return
			}

			if result.PlanID() != tt.req.PlanID {
				t.Errorf("RecordResult() planID = %v, want %v", result.PlanID(), tt.req.PlanID)
			}
			if result.CaseID() != tt.req.CaseID {
				t.Errorf("RecordResult() caseID = %v, want %v", result.CaseID(), tt.req.CaseID)
			}
			if string(result.Status()) != tt.req.Status {
				t.Errorf("RecordResult() status = %v, want %v", result.Status(), tt.req.Status)
			}
		})
	}
}

// TestPlanService_UpdatePlanStatus tests plan status update
func TestPlanService_UpdatePlanStatus(t *testing.T) {
	ctx := context.Background()
	planRepo := NewMockTestPlanRepository()
	resultRepo := NewMockTestResultRepository()
	caseRepo := NewMockTestCaseRepoForPlan()
	service := NewPlanService(planRepo, resultRepo, caseRepo)

	// Create test plan
	projectID := uuid.New()
	userID := uuid.New()
	testPlan, _ := testplan.NewTestPlan(projectID, "Test Plan", "Description", userID)
	planRepo.plans[testPlan.ID()] = testPlan

	// Create completed plan
	completedPlan, _ := testplan.NewTestPlan(projectID, "Completed Plan", "Description", userID)
	_ = completedPlan.UpdateStatus(testplan.StatusActive)
	_ = completedPlan.UpdateStatus(testplan.StatusCompleted)
	planRepo.plans[completedPlan.ID()] = completedPlan

	// Create archived plan
	archivedPlan, _ := testplan.NewTestPlan(projectID, "Archived Plan", "Description", userID)
	_ = archivedPlan.UpdateStatus(testplan.StatusActive)
	_ = archivedPlan.UpdateStatus(testplan.StatusCompleted)
	_ = archivedPlan.UpdateStatus(testplan.StatusArchived)
	planRepo.plans[archivedPlan.ID()] = archivedPlan

	tests := []struct {
		name      string
		planID    uuid.UUID
		newStatus string
		wantErr   error
	}{
		{
			name:      "draft to active",
			planID:    testPlan.ID(),
			newStatus: "active",
			wantErr:   nil,
		},
		{
			name:      "active to completed",
			planID:    testPlan.ID(),
			newStatus: "completed",
			wantErr:   nil,
		},
		{
			name:      "completed to archived",
			planID:    completedPlan.ID(),
			newStatus: "archived",
			wantErr:   nil,
		},
		{
			name:      "invalid transition - archived to active",
			planID:    archivedPlan.ID(),
			newStatus: "active",
			wantErr:   errors.New("update status: invalid status transition"),
		},
		{
			name:      "invalid status value",
			planID:    testPlan.ID(),
			newStatus: "invalid",
			wantErr:   errors.New("invalid plan status"),
		},
		{
			name:      "plan not found",
			planID:    uuid.New(),
			newStatus: "active",
			wantErr:   testplan.ErrPlanNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdatePlanStatus(ctx, tt.planID, tt.newStatus)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("UpdatePlanStatus() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) {
					t.Errorf("UpdatePlanStatus() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdatePlanStatus() unexpected error: %v", err)
				return
			}

			// Verify status was updated
			plan, _ := planRepo.FindByID(ctx, tt.planID)
			if string(plan.Status()) != tt.newStatus {
				t.Errorf("UpdatePlanStatus() status = %v, want %v", plan.Status(), tt.newStatus)
			}
		})
	}
}

// TestPlanService_DeletePlan tests plan deletion
func TestPlanService_DeletePlan(t *testing.T) {
	ctx := context.Background()
	planRepo := NewMockTestPlanRepository()
	resultRepo := NewMockTestResultRepository()
	caseRepo := NewMockTestCaseRepoForPlan()
	service := NewPlanService(planRepo, resultRepo, caseRepo)

	// Create test plan
	projectID := uuid.New()
	userID := uuid.New()
	testPlan, _ := testplan.NewTestPlan(projectID, "To Delete", "Description", userID)
	planRepo.plans[testPlan.ID()] = testPlan

	tests := []struct {
		name    string
		planID  uuid.UUID
		wantErr error
	}{
		{
			name:    "successful deletion",
			planID:  testPlan.ID(),
			wantErr: nil,
		},
		{
			name:    "plan not found",
			planID:  uuid.New(),
			wantErr: testplan.ErrPlanNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeletePlan(ctx, tt.planID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("DeletePlan() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("DeletePlan() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeletePlan() unexpected error: %v", err)
				return
			}

			// Verify plan is deleted
			_, err = planRepo.FindByID(ctx, tt.planID)
			if err == nil {
				t.Error("DeletePlan() plan still exists after deletion")
			}
		})
	}
}

// Helper function to create test cases
func createTestCase(t *testing.T, moduleID, userID uuid.UUID, title string) *testcase.TestCase {
	caseNumber := testcase.GenerateCaseNumber("TEST", "USER", 1)
	tc, err := testcase.NewTestCase(
		moduleID,
		userID,
		caseNumber,
		title,
		[]string{"Precondition"},
		[]string{"Step 1", "Step 2"},
		map[string]any{"status": "success"},
		testcase.CaseTypeFunctionality,
		testcase.PriorityP1,
	)
	if err != nil {
		t.Fatalf("Failed to create test case: %v", err)
	}
	return tc
}
