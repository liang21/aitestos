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
	chunkID         uuid.UUID
	documentID      uuid.UUID
	documentTitle   string
	similarityScore float64
}

// NewReferencedChunk creates a new ReferencedChunk
func NewReferencedChunk(chunkID, documentID uuid.UUID, title string, score float64) *ReferencedChunk {
	return &ReferencedChunk{
		chunkID:         chunkID,
		documentID:      documentID,
		documentTitle:   title,
		similarityScore: score,
	}
}

// ChunkID returns the chunk's ID
func (r *ReferencedChunk) ChunkID() uuid.UUID {
	return r.chunkID
}

// DocumentID returns the document's ID
func (r *ReferencedChunk) DocumentID() uuid.UUID {
	return r.documentID
}

// DocumentTitle returns the document title
func (r *ReferencedChunk) DocumentTitle() string {
	return r.documentTitle
}

// SimilarityScore returns the similarity score
func (r *ReferencedChunk) SimilarityScore() float64 {
	return r.similarityScore
}

// AiMetadata represents metadata about AI-generated test cases
type AiMetadata struct {
	generationTaskID uuid.UUID
	confidence       Confidence
	referencedChunks []*ReferencedChunk
	modelVersion     string
	generatedAt      time.Time
}

// NewAiMetadata creates new AI metadata
func NewAiMetadata(taskID uuid.UUID, confidence Confidence, chunks []*ReferencedChunk, modelVersion string) *AiMetadata {
	return &AiMetadata{
		generationTaskID: taskID,
		confidence:       confidence,
		referencedChunks: chunks,
		modelVersion:     modelVersion,
		generatedAt:      time.Now(),
	}
}

// GenerationTaskID returns the generation task ID
func (m *AiMetadata) GenerationTaskID() uuid.UUID {
	return m.generationTaskID
}

// Confidence returns the confidence level
func (m *AiMetadata) Confidence() Confidence {
	return m.confidence
}

// ReferencedChunks returns the referenced chunks
func (m *AiMetadata) ReferencedChunks() []*ReferencedChunk {
	return m.referencedChunks
}

// ModelVersion returns the model version
func (m *AiMetadata) ModelVersion() string {
	return m.modelVersion
}

// GeneratedAt returns the generation timestamp
func (m *AiMetadata) GeneratedAt() time.Time {
	return m.generatedAt
}

// IsAIGenerated returns true if metadata represents AI-generated content
func (m *AiMetadata) IsAIGenerated() bool {
	return m != nil && m.generationTaskID != uuid.Nil
}

// CalculateConfidence determines confidence based on chunks
func CalculateConfidence(chunks []*ReferencedChunk) Confidence {
	if len(chunks) >= 2 && chunks[0].SimilarityScore() > 0.8 {
		return ConfidenceHigh
	}
	if len(chunks) >= 1 && chunks[0].SimilarityScore() >= 0.5 {
		return ConfidenceMedium
	}
	return ConfidenceLow
}
