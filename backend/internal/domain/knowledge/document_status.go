// Package knowledge defines DocumentStatus value object
package knowledge

import "errors"

// DocumentStatus is a value object representing the processing status of a document
type DocumentStatus string

const (
	// StatusPending means the document is pending processing
	StatusPending DocumentStatus = "pending"
	// StatusProcessing means the document is being processed
	StatusProcessing DocumentStatus = "processing"
	// StatusCompleted means the document has been processed
	StatusCompleted DocumentStatus = "completed"
	// StatusFailed means the document processing failed
	StatusFailed DocumentStatus = "failed"
)

// ParseDocumentStatus validates and creates a DocumentStatus
func ParseDocumentStatus(s string) (DocumentStatus, error) {
	switch DocumentStatus(s) {
	case StatusPending, StatusProcessing, StatusCompleted, StatusFailed:
		return DocumentStatus(s), nil
	default:
		return "", errors.New("invalid document status")
	}
}

// String returns the string representation
func (s DocumentStatus) String() string {
	return string(s)
}

// CanTransitionTo checks if transition to another status is valid
func (s DocumentStatus) CanTransitionTo(to DocumentStatus) bool {
	transitions := map[DocumentStatus]map[DocumentStatus]bool{
		StatusPending:    {StatusProcessing: true},
		StatusProcessing: {StatusCompleted: true, StatusFailed: true},
		StatusFailed:     {StatusPending: true}, // Allow retry
		StatusCompleted:  {},
	}

	allowed, exists := transitions[s][to]
	return exists && allowed
}

// IsFinal returns true if the status is a final state
func (s DocumentStatus) IsFinal() bool {
	return s == StatusCompleted || s == StatusFailed
}
