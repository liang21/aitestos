// Package knowledge_test tests DocumentRepository implementation
package knowledge_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domainknowledge "github.com/liang21/aitestos/internal/domain/knowledge"
)

// Test fixtures
func createTestDocument(t *testing.T, projectID uuid.UUID, name string) *domainknowledge.Document {
	t.Helper()
	doc, err := domainknowledge.NewDocument(projectID, name, domainknowledge.TypePRD, "https://example.com/doc.pdf", uuid.New())
	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}
	return doc
}

func TestDocumentRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("save new document", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestDocumentRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find existing document", func(t *testing.T) {
		// Placeholder for integration test
	})

	t.Run("find non-existent document", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestDocumentRepository_FindByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by project id with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestDocumentRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("update document", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestDocumentRepository_UpdateStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("update document status", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestDocumentRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("soft delete document", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestDocumentRepository_CountByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("count documents by project id", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// MockDocumentRepository for testing without database
type MockDocumentRepository struct {
	documents     map[uuid.UUID]*domainknowledge.Document
	docsByProject map[uuid.UUID][]*domainknowledge.Document
}

func NewMockDocumentRepository() *MockDocumentRepository {
	return &MockDocumentRepository{
		documents:     make(map[uuid.UUID]*domainknowledge.Document),
		docsByProject: make(map[uuid.UUID][]*domainknowledge.Document),
	}
}

func (m *MockDocumentRepository) Save(ctx context.Context, doc *domainknowledge.Document) error {
	m.documents[doc.ID()] = doc
	m.docsByProject[doc.ProjectID()] = append(m.docsByProject[doc.ProjectID()], doc)
	return nil
}

func (m *MockDocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainknowledge.Document, error) {
	doc, ok := m.documents[id]
	if !ok {
		return nil, domainknowledge.ErrDocumentNotFound
	}
	return doc, nil
}

func (m *MockDocumentRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts domainknowledge.QueryOptions) ([]*domainknowledge.Document, error) {
	docs, ok := m.docsByProject[projectID]
	if !ok {
		return []*domainknowledge.Document{}, nil
	}
	return docs, nil
}

func (m *MockDocumentRepository) Update(ctx context.Context, doc *domainknowledge.Document) error {
	if _, ok := m.documents[doc.ID()]; !ok {
		return domainknowledge.ErrDocumentNotFound
	}
	m.documents[doc.ID()] = doc
	return nil
}

func (m *MockDocumentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domainknowledge.DocumentStatus) error {
	doc, ok := m.documents[id]
	if !ok {
		return domainknowledge.ErrDocumentNotFound
	}
	_ = doc.UpdateStatus(status)
	return nil
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	doc, ok := m.documents[id]
	if !ok {
		return domainknowledge.ErrDocumentNotFound
	}
	delete(m.documents, id)

	// Remove from project list
	projectDocs := m.docsByProject[doc.ProjectID()]
	for i, d := range projectDocs {
		if d.ID() == id {
			m.docsByProject[doc.ProjectID()] = append(projectDocs[:i], projectDocs[i+1:]...)
			break
		}
	}
	return nil
}

func (m *MockDocumentRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	return int64(len(m.docsByProject[projectID])), nil
}

func TestMockDocumentRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMockDocumentRepository()
	projectID := uuid.New()

	// Create
	doc := createTestDocument(t, projectID, "PRD Document")
	err := repo.Save(ctx, doc)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read by ID
	found, err := repo.FindByID(ctx, doc.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.ID() != doc.ID() {
		t.Errorf("FindByID().ID() = %v, want %v", found.ID(), doc.ID())
	}

	// Read by ProjectID
	docs, err := repo.FindByProjectID(ctx, projectID, domainknowledge.QueryOptions{})
	if err != nil {
		t.Fatalf("FindByProjectID() error = %v", err)
	}
	if len(docs) != 1 {
		t.Errorf("FindByProjectID() returned %d docs, want 1", len(docs))
	}

	// Update Status
	err = repo.UpdateStatus(ctx, doc.ID(), domainknowledge.StatusProcessing)
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
	found, _ = repo.FindByID(ctx, doc.ID())
	if found.Status() != domainknowledge.StatusProcessing {
		t.Errorf("UpdateStatus().Status() = %v, want processing", found.Status())
	}

	// Count by ProjectID
	count, err := repo.CountByProjectID(ctx, projectID)
	if err != nil {
		t.Fatalf("CountByProjectID() error = %v", err)
	}
	if count != 1 {
		t.Errorf("CountByProjectID() = %d, want 1", count)
	}

	// Delete
	err = repo.Delete(ctx, doc.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err = repo.FindByID(ctx, doc.ID())
	if err != domainknowledge.ErrDocumentNotFound {
		t.Errorf("FindByID() after Delete() error = %v, want %v", err, domainknowledge.ErrDocumentNotFound)
	}
}

func TestMockDocumentRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMockDocumentRepository()

	_, err := repo.FindByID(ctx, uuid.New())
	if err != domainknowledge.ErrDocumentNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domainknowledge.ErrDocumentNotFound)
	}

	doc := createTestDocument(t, uuid.New(), "test")
	err = repo.Update(ctx, doc)
	if err != domainknowledge.ErrDocumentNotFound {
		t.Errorf("Update() error = %v, want %v", err, domainknowledge.ErrDocumentNotFound)
	}

	err = repo.UpdateStatus(ctx, uuid.New(), domainknowledge.StatusProcessing)
	if err != domainknowledge.ErrDocumentNotFound {
		t.Errorf("UpdateStatus() error = %v, want %v", err, domainknowledge.ErrDocumentNotFound)
	}

	err = repo.Delete(ctx, uuid.New())
	if err != domainknowledge.ErrDocumentNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domainknowledge.ErrDocumentNotFound)
	}
}
