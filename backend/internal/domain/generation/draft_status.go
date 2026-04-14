// Package generation defines DraftStatus value object
package generation

import "errors"

// DraftStatus is a value object representing the status of a generated case draft
type DraftStatus string

const (
	// DraftPending means the draft is pending review
	DraftPending DraftStatus = "pending"
	// DraftConfirmed means the draft has been confirmed
	DraftConfirmed DraftStatus = "confirmed"
	// DraftRejected means the draft has been rejected
	DraftRejected DraftStatus = "rejected"
)

// ParseDraftStatus parses a string into DraftStatus
func ParseDraftStatus(s string) (DraftStatus, error) {
	switch s {
	case string(DraftPending):
		return DraftPending, nil
	case string(DraftConfirmed):
		return DraftConfirmed, nil
	case string(DraftRejected):
		return DraftRejected, nil
	default:
		return "", errors.New("invalid draft status")
	}
}

// String returns the string representation of DraftStatus
func (s DraftStatus) String() string {
	return string(s)
}

// CanTransitionTo checks if the status can transition to the target status
func (s DraftStatus) CanTransitionTo(target DraftStatus) bool {
	transitions := map[DraftStatus][]DraftStatus{
		DraftPending:   {DraftConfirmed, DraftRejected},
		DraftConfirmed: {},
		DraftRejected:  {},
	}

	allowed, exists := transitions[s]
	if !exists {
		return false
	}

	for _, t := range allowed {
		if t == target {
			return true
		}
	}
	return false
}

// IsFinal returns true if the status is a final state
func (s DraftStatus) IsFinal() bool {
	return s == DraftConfirmed || s == DraftRejected
}

// RejectionReason represents why a draft was rejected
type RejectionReason string

const (
	// ReasonDuplicate means the draft is a duplicate
	ReasonDuplicate RejectionReason = "duplicate"
	// ReasonIrrelevant means the draft is irrelevant
	ReasonIrrelevant RejectionReason = "irrelevant"
	// ReasonLowQuality means the draft has low quality
	ReasonLowQuality RejectionReason = "low_quality"
	// ReasonOther means other reasons
	ReasonOther RejectionReason = "other"
)

// ParseRejectionReason parses a string into RejectionReason
func ParseRejectionReason(s string) (RejectionReason, error) {
	switch s {
	case string(ReasonDuplicate):
		return ReasonDuplicate, nil
	case string(ReasonIrrelevant):
		return ReasonIrrelevant, nil
	case string(ReasonLowQuality):
		return ReasonLowQuality, nil
	case string(ReasonOther):
		return ReasonOther, nil
	default:
		return "", errors.New("invalid rejection reason")
	}
}

// String returns the string representation of RejectionReason
func (r RejectionReason) String() string {
	return string(r)
}
