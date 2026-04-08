// Package knowledge_test tests DocumentChunkRepository implementation
package knowledge_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domainknowledge "github.com/liang21/aitestos/internal/domain/knowledge"
)

// Test fixtures
func createTestChunk(t *testing.T, documentID uuid.UUID, index int, content string) *domainknowledge.DocumentChunk {
	t.Helper()
	return domainknowledge.NewDocumentChunk(documentID, index, content)
}

func TestDocumentChunkRepository_SaveBatch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("save multiple chunks", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("save empty batch", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestDocumentChunkRepository_FindByDocumentID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find chunks by document id", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("find chunks for non-existent document", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestDocumentChunkRepository_DeleteByDocumentID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("delete all chunks for a document", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestDocumentChunkRepository_CountByDocumentID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("count chunks for a document", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// MockDocumentChunkRepository for testing without database
type MockDocumentChunkRepository struct {
	chunks       map[uuid.UUID]*domainknowledge.DocumentChunk
	chunksByDoc map[uuid.UUID][]*domainknowledge.DocumentChunk
}

func NewMockDocumentChunkRepository() *MockDocumentChunkRepository {
	return &MockDocumentChunkRepository{
		chunks:    make(map[uuid.UUID]*domainknowledge.DocumentChunk),
		chunksByDoc: make(map[uuid.UUID][]*domainknowledge.DocumentChunk),
	}
}

func (m *MockDocumentChunkRepository) SaveBatch(ctx context.Context, chunks []*domainknowledge.DocumentChunk) error {
	if len(chunks) == 0 {
		return nil
	}
	for _, chunk := range chunks {
		m.chunks[chunk.ID()] = chunk
		m.chunksByDoc[chunk.DocumentID()] = append(m.chunksByDoc[chunk.DocumentID()], chunk)
	}
	return nil
}

func (m *MockDocumentChunkRepository) FindByDocumentID(ctx context.Context, documentID uuid.UUID) ([]*domainknowledge.DocumentChunk, error) {
	chunks, ok := m.chunksByDoc[documentID]
	if !ok {
		return []*domainknowledge.DocumentChunk{}, nil
	}
	return chunks, nil
}

func (m *MockDocumentChunkRepository) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	chunks := m.chunksByDoc[documentID]
	for _, chunk := range chunks {
		delete(m.chunks, chunk.ID())
	}
	delete(m.chunksByDoc, documentID)
	return nil
}

func (m *MockDocumentChunkRepository) CountByDocumentID(ctx context.Context, documentID uuid.UUID) (int64, error) {
	return int64(len(m.chunksByDoc[documentID])), nil
}

func TestMockDocumentChunkRepository_BatchOperations(t *testing.T) {
	ctx := context.Background()
	repo := NewMockDocumentChunkRepository()
	documentID := uuid.New()

	// Create batch
	chunks := []*domainknowledge.DocumentChunk{
		createTestChunk(t, documentID, 0, "Chunk 0 content"),
		createTestChunk(t, documentID, 1, "Chunk 1 content"),
		createTestChunk(t, documentID, 2, "Chunk 2 content"),
	}

	// Save batch
	err := repo.SaveBatch(ctx, chunks)
	if err != nil {
		t.Fatalf("SaveBatch() error = %v", err)
	}

	// Find by DocumentID
	found, err := repo.FindByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("FindByDocumentID() error = %v", err)
	}
	if len(found) != 3 {
		t.Errorf("FindByDocumentID() returned %d chunks, want 3", len(found))
	}

	// Count by DocumentID
	count, err := repo.CountByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("CountByDocumentID() error = %v", err)
	}
	if count != 3 {
		t.Errorf("CountByDocumentID() = %d, want 3", count)
	}

	// Delete by DocumentID
	err = repo.DeleteByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("DeleteByDocumentID() error = %v", err)
	}

	// Verify deletion
	found, err = repo.FindByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("FindByDocumentID() error = %v", err)
	}
	if len(found) != 0 {
		t.Errorf("FindByDocumentID() after DeleteByDocumentID() returned %d chunks, want 0", len(found))
	}
}

func TestMockDocumentChunkRepository_EmptyBatch(t *testing.T) {
	ctx := context.Background()
	repo := NewMockDocumentChunkRepository()

	err := repo.SaveBatch(ctx, []*domainknowledge.DocumentChunk{})
	if err != nil {
		t.Fatalf("SaveBatch() with empty slice error = %v", err)
	}
}

func TestMockDocumentChunkRepository_NonExistentDocument(t *testing.T) {
	ctx := context.Background()
	repo := NewMockDocumentChunkRepository()

	chunks, err := repo.FindByDocumentID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("FindByDocumentID() error = %v", err)
	}
	if len(chunks) != 0 {
		t.Errorf("FindByDocumentID() for non-existent document returned %d chunks, want 0", len(chunks))
	}

	count, err := repo.CountByDocumentID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("CountByDocumentID() error = %v", err)
	}
	if count != 0 {
		t.Errorf("CountByDocumentID() for non-existent document = %d, want 0", count)
	}
}
