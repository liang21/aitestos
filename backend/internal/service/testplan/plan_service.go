// Package testplan provides test plan management services
package testplan

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/testcase"
	"github.com/liang21/aitestos/internal/domain/testplan"
)

// CreatePlanRequest contains test plan creation data
type CreatePlanRequest struct {
	ProjectID   uuid.UUID   `json:"project_id" validate:"required"`
	Name        string      `json:"name" validate:"required,min=2,max=255"`
	Description string      `json:"description"`
	CaseIDs     []uuid.UUID `json:"case_ids"`
}

// UpdatePlanRequest contains test plan update data
type UpdatePlanRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty"`
}

// AddCaseRequest contains case addition data
type AddCaseRequest struct {
	CaseIDs []uuid.UUID `json:"case_ids" validate:"required,min=1"`
}

// RecordResultRequest contains result recording data
type RecordResultRequest struct {
	PlanID uuid.UUID `json:"plan_id" validate:"required"`
	CaseID uuid.UUID `json:"case_id" validate:"required"`
	Status string    `json:"status" validate:"required,oneof=pass fail block skip"`
	Note   string    `json:"note"`
}

// PlanDetail contains test plan info with related data
type PlanDetail struct {
	*testplan.TestPlan
	Cases   []*testcase.TestCase   `json:"cases"`
	Results []*testplan.TestResult `json:"results"`
	Stats   *PlanStatistics        `json:"stats"`
}

// PlanStatistics contains test plan execution statistics
type PlanStatistics struct {
	TotalCases   int64 `json:"total_cases"`
	PassedCases  int64 `json:"passed_cases"`
	FailedCases  int64 `json:"failed_cases"`
	BlockedCases int64 `json:"blocked_cases"`
	SkippedCases int64 `json:"skipped_cases"`
	Unexecuted   int64 `json:"unexecuted"`
}

// PlanListOptions contains pagination and filtering options
type PlanListOptions struct {
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
	Status   string `json:"status,omitempty"`
	Keywords string `json:"keywords,omitempty"`
}

// PlanService provides test plan management operations
type PlanService interface {
	// Plan management
	CreatePlan(ctx context.Context, req *CreatePlanRequest, userID uuid.UUID) (*testplan.TestPlan, error)
	GetPlan(ctx context.Context, id uuid.UUID) (*PlanDetail, error)
	ListPlans(ctx context.Context, projectID uuid.UUID, opts PlanListOptions) ([]*testplan.TestPlan, int64, error)
	UpdatePlan(ctx context.Context, id uuid.UUID, req *UpdatePlanRequest) (*testplan.TestPlan, error)
	DeletePlan(ctx context.Context, id uuid.UUID) error
	UpdatePlanStatus(ctx context.Context, id uuid.UUID, status string) error

	// Case management within plan
	AddCases(ctx context.Context, planID uuid.UUID, caseIDs []uuid.UUID) error
	RemoveCase(ctx context.Context, planID uuid.UUID, caseID uuid.UUID) error

	// Result management
	RecordResult(ctx context.Context, req *RecordResultRequest, userID uuid.UUID) (*testplan.TestResult, error)
	GetResults(ctx context.Context, planID uuid.UUID) ([]*testplan.TestResult, error)
	GetResultByCase(ctx context.Context, planID uuid.UUID, caseID uuid.UUID) (*testplan.TestResult, error)
}

// PlanServiceImpl implements PlanService
type PlanServiceImpl struct {
	planRepo   testplan.TestPlanRepository
	resultRepo testplan.TestResultRepository
	caseRepo   testcase.TestCaseRepository
}

// NewPlanService creates a new PlanService instance
func NewPlanService(
	planRepo testplan.TestPlanRepository,
	resultRepo testplan.TestResultRepository,
	caseRepo testcase.TestCaseRepository,
) PlanService {
	return &PlanServiceImpl{
		planRepo:   planRepo,
		resultRepo: resultRepo,
		caseRepo:   caseRepo,
	}
}

// CreatePlan creates a new test plan
func (s *PlanServiceImpl) CreatePlan(ctx context.Context, req *CreatePlanRequest, userID uuid.UUID) (*testplan.TestPlan, error) {
	plan, err := testplan.NewTestPlan(req.ProjectID, req.Name, req.Description, userID)
	if err != nil {
		return nil, fmt.Errorf("create plan: %w", err)
	}

	// Add initial cases if provided
	for _, caseID := range req.CaseIDs {
		if err := plan.AddCase(caseID); err != nil {
			return nil, fmt.Errorf("add case %s: %w", caseID, err)
		}
	}

	if err := s.planRepo.Save(ctx, plan); err != nil {
		return nil, fmt.Errorf("save plan: %w", err)
	}

	return plan, nil
}

// GetPlan retrieves test plan details with related data
func (s *PlanServiceImpl) GetPlan(ctx context.Context, id uuid.UUID) (*PlanDetail, error) {
	plan, err := s.planRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find plan: %w", err)
	}

	detail := &PlanDetail{
		TestPlan: plan,
		Cases:    make([]*testcase.TestCase, 0),
		Results:  make([]*testplan.TestResult, 0),
		Stats:    &PlanStatistics{},
	}

	// Get cases
	for _, caseID := range plan.CaseIDs() {
		tc, err := s.caseRepo.FindByID(ctx, caseID)
		if err == nil {
			detail.Cases = append(detail.Cases, tc)
		}
	}

	// Get results
	results, err := s.resultRepo.FindByPlanID(ctx, id, testplan.QueryOptions{})
	if err == nil {
		detail.Results = results
	}

	// Calculate stats
	detail.Stats.TotalCases = int64(len(detail.Cases))
	detail.Stats.Unexecuted = detail.Stats.TotalCases
	for _, result := range detail.Results {
		switch result.Status() {
		case testplan.ResultPass:
			detail.Stats.PassedCases++
			detail.Stats.Unexecuted--
		case testplan.ResultFail:
			detail.Stats.FailedCases++
			detail.Stats.Unexecuted--
		case testplan.ResultBlock:
			detail.Stats.BlockedCases++
			detail.Stats.Unexecuted--
		case testplan.ResultSkip:
			detail.Stats.SkippedCases++
			detail.Stats.Unexecuted--
		}
	}

	return detail, nil
}

// ListPlans lists test plans with pagination
func (s *PlanServiceImpl) ListPlans(ctx context.Context, projectID uuid.UUID, opts PlanListOptions) ([]*testplan.TestPlan, int64, error) {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}

	queryOpts := testplan.QueryOptions{
		Offset:   opts.Offset,
		Limit:    opts.Limit,
		Keywords: opts.Keywords,
	}

	plans, err := s.planRepo.FindByProjectID(ctx, projectID, queryOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("list plans: %w", err)
	}

	return plans, int64(len(plans)), nil
}

// UpdatePlan updates test plan information
func (s *PlanServiceImpl) UpdatePlan(ctx context.Context, id uuid.UUID, req *UpdatePlanRequest) (*testplan.TestPlan, error) {
	plan, err := s.planRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find plan: %w", err)
	}

	if req.Name != nil {
		plan.UpdateName(*req.Name)
	}

	if req.Description != nil {
		plan.UpdateDescription(*req.Description)
	}

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return nil, fmt.Errorf("update plan: %w", err)
	}

	return plan, nil
}

// DeletePlan deletes a test plan
func (s *PlanServiceImpl) DeletePlan(ctx context.Context, id uuid.UUID) error {
	if err := s.planRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete plan: %w", err)
	}
	return nil
}

// UpdatePlanStatus updates test plan status
func (s *PlanServiceImpl) UpdatePlanStatus(ctx context.Context, id uuid.UUID, status string) error {
	plan, err := s.planRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find plan: %w", err)
	}

	newStatus, err := testplan.ParsePlanStatus(status)
	if err != nil {
		return errors.New("invalid plan status")
	}

	if err := plan.UpdateStatus(newStatus); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return fmt.Errorf("save plan: %w", err)
	}

	return nil
}

// AddCases adds test cases to a plan
func (s *PlanServiceImpl) AddCases(ctx context.Context, planID uuid.UUID, caseIDs []uuid.UUID) error {
	plan, err := s.planRepo.FindByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("find plan: %w", err)
	}

	for _, caseID := range caseIDs {
		// Verify case exists
		_, err := s.caseRepo.FindByID(ctx, caseID)
		if err != nil {
			return testcase.ErrCaseNotFound
		}

		if err := plan.AddCase(caseID); err != nil {
			return fmt.Errorf("add case %s: %w", caseID, err)
		}
	}

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return fmt.Errorf("save plan: %w", err)
	}

	return nil
}

// RemoveCase removes a test case from a plan
func (s *PlanServiceImpl) RemoveCase(ctx context.Context, planID uuid.UUID, caseID uuid.UUID) error {
	plan, err := s.planRepo.FindByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("find plan: %w", err)
	}

	plan.RemoveCase(caseID)

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return fmt.Errorf("save plan: %w", err)
	}

	return nil
}

// RecordResult records a test execution result
func (s *PlanServiceImpl) RecordResult(ctx context.Context, req *RecordResultRequest, userID uuid.UUID) (*testplan.TestResult, error) {
	// Verify plan exists
	plan, err := s.planRepo.FindByID(ctx, req.PlanID)
	if err != nil {
		return nil, fmt.Errorf("find plan: %w", err)
	}

	// Verify case is in plan
	if !plan.HasCase(req.CaseID) {
		return nil, errors.New("case not in plan")
	}

	status, err := testplan.ParseResultStatus(req.Status)
	if err != nil {
		return nil, errors.New("invalid result status")
	}

	result, err := testplan.NewTestResult(req.PlanID, req.CaseID, userID, status, req.Note)
	if err != nil {
		return nil, fmt.Errorf("create result: %w", err)
	}

	if err := s.resultRepo.Save(ctx, result); err != nil {
		return nil, fmt.Errorf("save result: %w", err)
	}

	return result, nil
}

// GetResults retrieves all results for a plan
func (s *PlanServiceImpl) GetResults(ctx context.Context, planID uuid.UUID) ([]*testplan.TestResult, error) {
	results, err := s.resultRepo.FindByPlanID(ctx, planID, testplan.QueryOptions{})
	if err != nil {
		return nil, fmt.Errorf("get results: %w", err)
	}
	return results, nil
}

// GetResultByCase retrieves result for a specific case in a plan
func (s *PlanServiceImpl) GetResultByCase(ctx context.Context, planID uuid.UUID, caseID uuid.UUID) (*testplan.TestResult, error) {
	result, err := s.resultRepo.FindByID(ctx, caseID) // Use FindByID for now
	if err != nil {
		return nil, fmt.Errorf("get result by case: %w", err)
	}
	return result, nil
}
