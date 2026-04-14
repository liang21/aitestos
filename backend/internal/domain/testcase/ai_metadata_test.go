// Package testcase_test tests AiMetadata value object
package testcase_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/testcase"
)

func TestConfidence_String(t *testing.T) {
	tests := []struct {
		name       string
		confidence testcase.Confidence
		want       string
	}{
		{
			name:       "high confidence",
			confidence: testcase.ConfidenceHigh,
			want:       "high",
		},
		{
			name:       "medium confidence",
			confidence: testcase.ConfidenceMedium,
			want:       "medium",
		},
		{
			name:       "low confidence",
			confidence: testcase.ConfidenceLow,
			want:       "low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.confidence); got != tt.want {
				t.Errorf("Confidence = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReferencedChunk(t *testing.T) {
	chunkID := uuid.New()
	docID := uuid.New()

	chunk := testcase.NewReferencedChunk(chunkID, docID, "Test Document", 0.85)

	if chunk.ChunkID != chunkID {
		t.Errorf("ChunkID = %v, want %v", chunk.ChunkID, chunkID)
	}
	if chunk.DocumentID != docID {
		t.Errorf("DocumentID = %v, want %v", chunk.DocumentID, docID)
	}
	if chunk.DocumentTitle != "Test Document" {
		t.Errorf("DocumentTitle = %v, want Test Document", chunk.DocumentTitle)
	}
	if chunk.SimilarityScore != 0.85 {
		t.Errorf("SimilarityScore = %v, want 0.85", chunk.SimilarityScore)
	}
}

func TestCalculateConfidence(t *testing.T) {
	tests := []struct {
		name      string
		numChunks int
		score     float64
		want      testcase.Confidence
	}{
		{
			name:      "high confidence - 2+ chunks with score > 0.8",
			numChunks: 2,
			score:     0.85,
			want:      testcase.ConfidenceHigh,
		},
		{
			name:      "medium confidence - 1 chunk with score >= 0.5",
			numChunks: 1,
			score:     0.65,
			want:      testcase.ConfidenceMedium,
		},
		{
			name:      "low confidence - score < 0.5",
			numChunks: 1,
			score:     0.3,
			want:      testcase.ConfidenceLow,
		},
		{
			name:      "medium confidence - 2 chunks but score <= 0.8",
			numChunks: 2,
			score:     0.7,
			want:      testcase.ConfidenceMedium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := make([]*testcase.ReferencedChunk, tt.numChunks)
			for i := 0; i < tt.numChunks; i++ {
				chunks[i] = testcase.NewReferencedChunk(uuid.New(), uuid.New(), "Doc", tt.score)
			}
			if got := testcase.CalculateConfidence(chunks); got != tt.want {
				t.Errorf("CalculateConfidence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAiMetadata(t *testing.T) {
	taskID := uuid.New()
	chunk1 := testcase.NewReferencedChunk(uuid.New(), uuid.New(), "Doc1", 0.9)
	chunk2 := testcase.NewReferencedChunk(uuid.New(), uuid.New(), "Doc2", 0.7)
	chunks := []*testcase.ReferencedChunk{chunk1, chunk2}

	metadata := testcase.NewAiMetadata(taskID, testcase.ConfidenceHigh, chunks, "deepseek-v3")

	if metadata.GenerationTaskID != taskID {
		t.Errorf("GenerationTaskID = %v, want %v", metadata.GenerationTaskID, taskID)
	}
	if metadata.Confidence != testcase.ConfidenceHigh {
		t.Errorf("Confidence = %v, want high", metadata.Confidence)
	}
	if metadata.ModelVersion != "deepseek-v3" {
		t.Errorf("ModelVersion = %v, want deepseek-v3", metadata.ModelVersion)
	}
	if len(metadata.ReferencedChunks) != 2 {
		t.Errorf("ReferencedChunks length = %v, want 2", len(metadata.ReferencedChunks))
	}
	if metadata.GeneratedAt.IsZero() {
		t.Error("GeneratedAt should not be zero")
	}
}

func TestAiMetadata_IsAIGenerated(t *testing.T) {
	taskID := uuid.New()
	metadata := testcase.NewAiMetadata(taskID, testcase.ConfidenceHigh, nil, "v1")

	if !metadata.IsAIGenerated() {
		t.Error("IsAIGenerated() should return true for non-nil metadata")
	}

	var nilMetadata *testcase.AiMetadata
	if nilMetadata.IsAIGenerated() {
		t.Error("IsAIGenerated() should return false for nil metadata")
	}
}
