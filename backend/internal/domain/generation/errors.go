// Package generation defines generation domain errors
package generation

import "errors"

var (
	// ErrTaskNotFound indicates generation task does not exist
	ErrTaskNotFound = errors.New("generation task not found")
	// ErrTaskAlreadyProcessed indicates task has already been processed
	ErrTaskAlreadyProcessed = errors.New("task already processed")
	// ErrDraftNotFound indicates draft does not exist
	ErrDraftNotFound = errors.New("case draft not found")
	// ErrDraftAlreadyConfirmed indicates draft has already been confirmed
	ErrDraftAlreadyConfirmed = errors.New("draft already confirmed")
	// ErrDraftAlreadyRejected indicates draft has already been rejected
	ErrDraftAlreadyRejected = errors.New("draft already rejected")
	// ErrInvalidDraftStatus indicates draft status is invalid for operation
	ErrInvalidDraftStatus = errors.New("invalid draft status for operation")
	// ErrLLMCallFailed indicates LLM call failed
	ErrLLMCallFailed = errors.New("LLM call failed")
	// ErrRAGNoResult indicates RAG retrieval returned no results
	ErrRAGNoResult = errors.New("RAG retrieval returned no results")
	// ErrConcurrentModification indicates concurrent modification detected
	ErrConcurrentModification = errors.New("concurrent modification detected")
	// ErrLLMTimeout indicates LLM call timed out
	ErrLLMTimeout = errors.New("LLM call timed out")
	// ErrGenerationQueueFull indicates generation queue is full
	ErrGenerationQueueFull = errors.New("generation queue is full")
)
