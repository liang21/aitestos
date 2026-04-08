// Package testplan defines TestPlan aggregate
package testplan

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// TestPlan is the aggregate root for testplan context
type TestPlan struct {
	id          uuid.UUID
	projectID   uuid.UUID
	name        string
	description string
	status      PlanStatus
	caseIDs     map[uuid.UUID]bool
	createdBy   uuid.UUID
	createdAt   time.Time
	updatedAt   time.Time
}

// NewTestPlan creates a new test plan
func NewTestPlan(projectID uuid.UUID, name, description string, userID uuid.UUID) (*TestPlan, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}
	if name == "" {
		return nil, errors.New("plan name cannot be empty")
	}
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be nil")
	}

	now := time.Now()
	return &TestPlan{
		id:          uuid.New(),
		projectID:   projectID,
		name:        name,
		description: description,
		status:      StatusDraft,
		caseIDs:     make(map[uuid.UUID]bool),
		createdBy:   userID,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ID returns the test plan's unique identifier
func (p *TestPlan) ID() uuid.UUID {
	return p.id
}

// ProjectID returns the associated project's ID
func (p *TestPlan) ProjectID() uuid.UUID {
	return p.projectID
}

// Name returns the test plan's name
func (p *TestPlan) Name() string {
	return p.name
}

// Description returns the test plan's description
func (p *TestPlan) Description() string {
	return p.description
}

// Status returns the test plan's status
func (p *TestPlan) Status() PlanStatus {
	return p.status
}

// CaseIDs returns the associated test case IDs
func (p *TestPlan) CaseIDs() []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(p.caseIDs))
	for id := range p.caseIDs {
		ids = append(ids, id)
	}
	return ids
}

// CreatedBy returns the user who created this plan
func (p *TestPlan) CreatedBy() uuid.UUID {
	return p.createdBy
}

// CreatedAt returns the creation timestamp
func (p *TestPlan) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the last update timestamp
func (p *TestPlan) UpdatedAt() time.Time {
	return p.updatedAt
}

// AddCase adds a test case to the plan
func (p *TestPlan) AddCase(caseID uuid.UUID) error {
	if p.status == StatusArchived {
		return errors.New("cannot add case to archived plan")
	}
	if p.caseIDs[caseID] {
		return errors.New("case already exists in plan")
	}
	p.caseIDs[caseID] = true
	p.updatedAt = time.Now()
	return nil
}

// RemoveCase removes a test case from the plan
func (p *TestPlan) RemoveCase(caseID uuid.UUID) {
	delete(p.caseIDs, caseID)
	p.updatedAt = time.Now()
}

// HasCase checks if a case is in the plan
func (p *TestPlan) HasCase(caseID uuid.UUID) bool {
	return p.caseIDs[caseID]
}

// UpdateStatus updates the test plan status
func (p *TestPlan) UpdateStatus(status PlanStatus) error {
	if !p.status.CanTransitionTo(status) {
		return errors.New("invalid status transition")
	}
	p.status = status
	p.updatedAt = time.Now()
	return nil
}

// Reconstruct reconstructs a TestPlan from stored data
func Reconstruct(
	id uuid.UUID,
	projectID uuid.UUID,
	name string,
	description string,
	status PlanStatus,
	caseIDs []uuid.UUID,
	createdBy uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
) *TestPlan {
	caseMap := make(map[uuid.UUID]bool)
	for _, caseID := range caseIDs {
		caseMap[caseID] = true
	}

	return &TestPlan{
		id:          id,
		projectID:   projectID,
		name:        name,
		description: description,
		status:      status,
		caseIDs:     caseMap,
		createdBy:   createdBy,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
