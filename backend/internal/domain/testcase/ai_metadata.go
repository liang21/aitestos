// Package testcase defines AiMetadata value object
package testcase

import (
	"time"

	"github.com/google/uuid"
)

// Confidence is a value object representing AI confidence level
type Confidence string

const (
	// ConfidenceHigh indicates high confidence
	ConfidenceHigh Confidence = "high"
	// ConfidenceMedium indicates medium confidence
	ConfidenceMedium Confidence = "medium"
	// ConfidenceLow indicates low confidence
	ConfidenceLow Confidence = "low"
)

// ReferencedChunk represents a document chunk referenced by AI
type ReferencedChunk struct {
	ChunkID         uuid.UUID `json:"chunk_id"`
	DocumentID      uuid.UUID `json:"document_id"`
	DocumentTitle   string    `json:"document_title"`
	SimilarityScore float64   `json:"similarity_score"`
}

// NewReferencedChunk creates a new ReferencedChunk
func NewReferencedChunk(chunkID, documentID uuid.UUID, title string, score float64) *ReferencedChunk {
	return &ReferencedChunk{
		ChunkID:         chunkID,
		DocumentID:      documentID,
		DocumentTitle:   title,
		SimilarityScore: score,
	}
}

// AiMetadata represents metadata about AI-generated test cases
type AiMetadata struct {
	GenerationTaskID uuid.UUID          `json:"generation_task_id"`
	Confidence       Confidence         `json:"confidence"`
	ReferencedChunks []*ReferencedChunk `json:"referenced_chunks"`
	ModelVersion     string             `json:"model_version"`
	GeneratedAt      time.Time          `json:"generated_at"`
}

// NewAiMetadata creates new AI metadata
func NewAiMetadata(taskID uuid.UUID, confidence Confidence, chunks []*ReferencedChunk, modelVersion string) *AiMetadata {
	return &AiMetadata{
		GenerationTaskID: taskID,
		Confidence:       confidence,
		ReferencedChunks: chunks,
		ModelVersion:     modelVersion,
		GeneratedAt:      time.Now(),
	}
}

// IsAIGenerated returns true if metadata represents AI-generated content
func (m *AiMetadata) IsAIGenerated() bool {
	return m != nil && m.GenerationTaskID != uuid.Nil
}

// CalculateConfidence determines confidence based on chunks
func CalculateConfidence(chunks []*ReferencedChunk) Confidence {
	if len(chunks) >= 2 && chunks[0].SimilarityScore > 0.8 {
		return ConfidenceHigh
	}
	if len(chunks) >= 1 && chunks[0].SimilarityScore >= 0.5 {
		return ConfidenceMedium
	}
	return ConfidenceLow
}
