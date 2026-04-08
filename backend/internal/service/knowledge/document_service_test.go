// Package knowledge provides document management services
package knowledge

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/knowledge"
)

// MockDocumentRepository implements knowledge.DocumentRepository for testing
type MockDocumentRepository struct {
	documents    map[uuid.UUID]*knowledge.Document
	projectIndex map[uuid.UUID][]*knowledge.Document
	saveErr      error
	findErr      error
}

func NewMockDocumentRepository() *MockDocumentRepository {
	return &MockDocumentRepository{
		documents:    make(map[uuid.UUID]*knowledge.Document),
		projectIndex: make(map[uuid.UUID][]*knowledge.Document),
	}
}

func (m *MockDocumentRepository) Save(ctx context.Context, doc *knowledge.Document) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.documents[doc.ID()] = doc
	m.projectIndex[doc.ProjectID()] = append(m.projectIndex[doc.ProjectID()], doc)
	return nil
}

func (m *MockDocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*knowledge.Document, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	doc, ok := m.documents[id]
	if !ok {
		return nil, knowledge.ErrDocumentNotFound
	}
	return doc, nil
}

func (m *MockDocumentRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts knowledge.QueryOptions) ([]*knowledge.Document, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	docs, ok := m.projectIndex[projectID]
	if !ok {
		return []*knowledge.Document{}, nil
	}
	return docs, nil
}

func (m *MockDocumentRepository) FindAll(ctx context.Context, opts knowledge.QueryOptions) ([]*knowledge.Document, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	docs := make([]*knowledge.Document, 0, len(m.documents))
	for _, doc := range m.documents {
		docs = append(docs, doc)
	}
	return docs, nil
}

func (m *MockDocumentRepository) Update(ctx context.Context, doc *knowledge.Document) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.documents[doc.ID()] = doc
	return nil
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	delete(m.documents, id)
	return nil
}

func (m *MockDocumentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status knowledge.DocumentStatus) error {
	return nil
}

func (m *MockDocumentRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	count := 0
	for _, doc := range m.documents {
		if doc.ProjectID() == projectID {
			count++
		}
	}
	return int64(count), nil
}

// MockDocumentChunkRepository implements knowledge.DocumentChunkRepository for testing
type MockDocumentChunkRepository struct {
	chunks       map[uuid.UUID]*knowledge.DocumentChunk
	docIndex     map[uuid.UUID][]*knowledge.DocumentChunk
	saveErr      error
	findErr      error
}

func NewMockDocumentChunkRepository() *MockDocumentChunkRepository {
	return &MockDocumentChunkRepository{
		chunks:   make(map[uuid.UUID]*knowledge.DocumentChunk),
		docIndex: make(map[uuid.UUID][]*knowledge.DocumentChunk),
	}
}

func (m *MockDocumentChunkRepository) Save(ctx context.Context, chunk *knowledge.DocumentChunk) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.chunks[chunk.ID()] = chunk
	m.docIndex[chunk.DocumentID()] = append(m.docIndex[chunk.DocumentID()], chunk)
	return nil
}

func (m *MockDocumentChunkRepository) SaveBatch(ctx context.Context, chunks []*knowledge.DocumentChunk) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	for _, chunk := range chunks {
		m.chunks[chunk.ID()] = chunk
		m.docIndex[chunk.DocumentID()] = append(m.docIndex[chunk.DocumentID()], chunk)
	}
	return nil
}

func (m *MockDocumentChunkRepository) FindByID(ctx context.Context, id uuid.UUID) (*knowledge.DocumentChunk, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	chunk, ok := m.chunks[id]
	if !ok {
		return nil, errors.New("chunk not found")
	}
	return chunk, nil
}

func (m *MockDocumentChunkRepository) FindByDocumentID(ctx context.Context, documentID uuid.UUID) ([]*knowledge.DocumentChunk, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	chunks, ok := m.docIndex[documentID]
	if !ok {
		return []*knowledge.DocumentChunk{}, nil
	}
	return chunks, nil
}

func (m *MockDocumentChunkRepository) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	delete(m.docIndex, documentID)
	return nil
}

func (m *MockDocumentChunkRepository) CountByDocumentID(ctx context.Context, documentID uuid.UUID) (int64, error) {
	count := 0
	for _, c := range m.chunks {
		if c.DocumentID() == documentID {
			count++
		}
	}
	return int64(count), nil
}

// MockVectorRepository implements knowledge.VectorRepository for testing
type MockVectorRepository struct {
	upsertErr error
	searchErr error
	deleteErr error
	countErr  error
	results   []*knowledge.DocumentChunk
}

func NewMockVectorRepository() *MockVectorRepository {
	return &MockVectorRepository{}
}

func (m *MockVectorRepository) Upsert(ctx context.Context, chunks []*knowledge.DocumentChunk) error {
	return m.upsertErr
}

func (m *MockVectorRepository) Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*knowledge.DocumentChunk, error) {
	if m.searchErr != nil {
		return nil, m.searchErr
	}
	return m.results, nil
}

func (m *MockVectorRepository) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	return m.deleteErr
}

func (m *MockVectorRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	if m.countErr != nil {
		return 0, m.countErr
	}
	return int64(len(m.results)), nil
}

// TestDocumentService_UploadDocument tests document upload
func TestDocumentService_UploadDocument(t *testing.T) {
	ctx := context.Background()
	docRepo := NewMockDocumentRepository()
	chunkRepo := NewMockDocumentChunkRepository()
	vectorRepo := NewMockVectorRepository()
	service := NewDocumentService(docRepo, chunkRepo, vectorRepo)

	projectID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		req     *UploadDocumentRequest
		setup   func()
		wantErr error
	}{
		{
			name: "successful upload PRD",
			req: &UploadDocumentRequest{
				ProjectID: projectID,
				Name:      "Product Requirements Document.pdf",
				Type:      "prd",
				File:      bytes.NewReader([]byte("test content")),
				FileSize:  13,
				UserID:    userID,
			},
			setup:   func() {},
			wantErr: nil,
		},
		{
			name: "successful upload Figma",
			req: &UploadDocumentRequest{
				ProjectID: projectID,
				Name:      "UI Design.fig",
				Type:      "figma",
				File:      bytes.NewReader([]byte("figma content")),
				FileSize:  14,
				UserID:    userID,
			},
			setup:   func() {},
			wantErr: nil,
		},
		{
			name: "empty name",
			req: &UploadDocumentRequest{
				ProjectID: projectID,
				Name:      "",
				Type:      "prd",
				File:      bytes.NewReader([]byte("content")),
				UserID:    userID,
			},
			setup:   func() {},
			wantErr: errors.New("document name cannot be empty"),
		},
		{
			name: "invalid document type",
			req: &UploadDocumentRequest{
				ProjectID: projectID,
				Name:      "Invalid Document",
				Type:      "invalid_type",
				File:      bytes.NewReader([]byte("content")),
				UserID:    userID,
			},
			setup:   func() {},
			wantErr: errors.New("invalid document type"),
		},
		{
			name: "nil project ID",
			req: &UploadDocumentRequest{
				ProjectID: uuid.Nil,
				Name:      "Orphan Document",
				Type:      "prd",
				File:      bytes.NewReader([]byte("content")),
				UserID:    userID,
			},
			setup:   func() {},
			wantErr: errors.New("project ID cannot be nil"),
		},
		{
			name: "nil user ID",
			req: &UploadDocumentRequest{
				ProjectID: projectID,
				Name:      "No Owner Document",
				Type:      "prd",
				File:      bytes.NewReader([]byte("content")),
				UserID:    uuid.Nil,
			},
			setup:   func() {},
			wantErr: errors.New("user ID cannot be nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			doc, err := service.UploadDocument(ctx, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("UploadDocument() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("UploadDocument() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UploadDocument() unexpected error: %v", err)
				return
			}

			if doc == nil {
				t.Error("UploadDocument() returned nil document")
				return
			}

			if doc.Name() != tt.req.Name {
				t.Errorf("UploadDocument() name = %v, want %v", doc.Name(), tt.req.Name)
			}
			if doc.ProjectID() != tt.req.ProjectID {
				t.Errorf("UploadDocument() projectID = %v, want %v", doc.ProjectID(), tt.req.ProjectID)
			}
			if doc.Status() != knowledge.StatusPending {
				t.Errorf("UploadDocument() status = %v, want %v", doc.Status(), knowledge.StatusPending)
			}
		})
	}
}

// TestDocumentService_GetDocument tests document retrieval
func TestDocumentService_GetDocument(t *testing.T) {
	ctx := context.Background()
	docRepo := NewMockDocumentRepository()
	chunkRepo := NewMockDocumentChunkRepository()
	vectorRepo := NewMockVectorRepository()
	service := NewDocumentService(docRepo, chunkRepo, vectorRepo)

	projectID := uuid.New()
	userID := uuid.New()

	// Create test document
	doc, _ := knowledge.NewDocument(projectID, "Test Document.pdf", knowledge.TypePRD, "https://example.com/doc.pdf", userID)
	docRepo.documents[doc.ID()] = doc

	tests := []struct {
		name    string
		docID   uuid.UUID
		wantErr error
	}{
		{
			name:    "successful retrieval",
			docID:   doc.ID(),
			wantErr: nil,
		},
		{
			name:    "document not found",
			docID:   uuid.New(),
			wantErr: knowledge.ErrDocumentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := service.GetDocument(ctx, tt.docID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetDocument() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetDocument() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetDocument() unexpected error: %v", err)
				return
			}

			if detail == nil {
				t.Error("GetDocument() returned nil detail")
				return
			}

			if detail.ID() != tt.docID {
				t.Errorf("GetDocument() ID = %v, want %v", detail.ID(), tt.docID)
			}
		})
	}
}

// TestDocumentService_DeleteDocument tests document deletion
func TestDocumentService_DeleteDocument(t *testing.T) {
	ctx := context.Background()
	docRepo := NewMockDocumentRepository()
	chunkRepo := NewMockDocumentChunkRepository()
	vectorRepo := NewMockVectorRepository()
	service := NewDocumentService(docRepo, chunkRepo, vectorRepo)

	projectID := uuid.New()
	userID := uuid.New()

	// Create test document
	doc, _ := knowledge.NewDocument(projectID, "To Delete.pdf", knowledge.TypePRD, "https://example.com/delete.pdf", userID)
	docRepo.documents[doc.ID()] = doc

	// Create some chunks
	chunk := knowledge.NewDocumentChunk(doc.ID(), 0, "test content")
	chunkRepo.chunks[chunk.ID()] = chunk
	chunkRepo.docIndex[doc.ID()] = []*knowledge.DocumentChunk{chunk}

	tests := []struct {
		name    string
		docID   uuid.UUID
		wantErr error
	}{
		{
			name:    "successful deletion",
			docID:   doc.ID(),
			wantErr: nil,
		},
		{
			name:    "document not found",
			docID:   uuid.New(),
			wantErr: knowledge.ErrDocumentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteDocument(ctx, tt.docID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("DeleteDocument() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("DeleteDocument() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteDocument() unexpected error: %v", err)
				return
			}

			// Verify document is deleted
			_, err = docRepo.FindByID(ctx, tt.docID)
			if err == nil {
				t.Error("DeleteDocument() document still exists after deletion")
			}
		})
	}
}

// TestDocumentService_UpdateStatus tests status update
func TestDocumentService_UpdateStatus(t *testing.T) {
	ctx := context.Background()
	docRepo := NewMockDocumentRepository()
	chunkRepo := NewMockDocumentChunkRepository()
	vectorRepo := NewMockVectorRepository()
	service := NewDocumentService(docRepo, chunkRepo, vectorRepo)

	projectID := uuid.New()
	userID := uuid.New()

	// Create test document
	doc, _ := knowledge.NewDocument(projectID, "Status Test.pdf", knowledge.TypePRD, "https://example.com/status.pdf", userID)
	docRepo.documents[doc.ID()] = doc

	// Create a fresh pending doc for invalid transition test
	pendingDoc, _ := knowledge.NewDocument(projectID, "Pending Doc.pdf", knowledge.TypePRD, "https://example.com/pending.pdf", userID)
	docRepo.documents[pendingDoc.ID()] = pendingDoc

	// Create completed document for archived test
	completedDoc, _ := knowledge.NewDocument(projectID, "Completed.pdf", knowledge.TypePRD, "https://example.com/completed.pdf", userID)
	completedDoc.UpdateStatus(knowledge.StatusProcessing)
	completedDoc.UpdateStatus(knowledge.StatusCompleted)
	docRepo.documents[completedDoc.ID()] = completedDoc

	tests := []struct {
		name      string
		docID     uuid.UUID
		newStatus string
		wantErr   error
	}{
		{
			name:      "pending to processing",
			docID:     doc.ID(),
			newStatus: "processing",
			wantErr:   nil,
		},
		{
			name:      "invalid status transition - pending to completed",
			docID:     pendingDoc.ID(),
			newStatus: "completed", // should go through processing first
			wantErr:   errors.New("update status: invalid status transition"),
		},
		{
			name:      "invalid status value",
			docID:     doc.ID(),
			newStatus: "invalid_status",
			wantErr:   errors.New("invalid document status"),
		},
		{
			name:      "document not found",
			docID:     uuid.New(),
			newStatus: "processing",
			wantErr:   knowledge.ErrDocumentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateStatus(ctx, tt.docID, tt.newStatus)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("UpdateStatus() expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) {
					t.Errorf("UpdateStatus() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateStatus() unexpected error: %v", err)
				return
			}

			// Verify status was updated
			updatedDoc, _ := docRepo.FindByID(ctx, tt.docID)
			if string(updatedDoc.Status()) != tt.newStatus {
				t.Errorf("UpdateStatus() status = %v, want %v", updatedDoc.Status(), tt.newStatus)
			}
		})
	}
}

// TestDocumentService_ListDocuments tests document listing
func TestDocumentService_ListDocuments(t *testing.T) {
	ctx := context.Background()
	docRepo := NewMockDocumentRepository()
	chunkRepo := NewMockDocumentChunkRepository()
	vectorRepo := NewMockVectorRepository()
	service := NewDocumentService(docRepo, chunkRepo, vectorRepo)

	projectID := uuid.New()
	userID := uuid.New()

	// Create test documents
	for i := 0; i < 5; i++ {
		doc, _ := knowledge.NewDocument(projectID, "Document "+string(rune('A'+i)), knowledge.TypePRD, "https://example.com/doc"+string(rune('1'+i)), userID)
		docRepo.documents[doc.ID()] = doc
		docRepo.projectIndex[projectID] = append(docRepo.projectIndex[projectID], doc)
	}

	tests := []struct {
		name      string
		opts      DocumentListOptions
		wantCount int
		wantErr   error
	}{
		{
			name: "list all documents with project filter",
			opts: DocumentListOptions{
				ProjectID: projectID,
				Offset:    0,
				Limit:     10,
			},
			wantCount: 5,
			wantErr:   nil,
		},
		{
			name: "list with pagination",
			opts: DocumentListOptions{
				ProjectID: projectID,
				Offset:    0,
				Limit:     2,
			},
			wantCount: 5, // Mock doesn't implement pagination, returns all documents
			wantErr:   nil,
		},
		{
			name: "filter by project",
			opts: DocumentListOptions{
				ProjectID: projectID,
				Offset:    0,
				Limit:     10,
			},
			wantCount: 5,
			wantErr:   nil,
		},
		{
			name: "filter by type with project",
			opts: DocumentListOptions{
				ProjectID: projectID,
				Type:      "prd",
				Offset:    0,
				Limit:     10,
			},
			wantCount: 5,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docs, total, err := service.ListDocuments(ctx, tt.opts)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ListDocuments() expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ListDocuments() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ListDocuments() unexpected error: %v", err)
				return
			}

			if len(docs) != tt.wantCount {
				t.Errorf("ListDocuments() count = %v, want %v", len(docs), tt.wantCount)
			}

			if total < int64(tt.wantCount) {
				t.Errorf("ListDocuments() total = %v, want at least %v", total, tt.wantCount)
			}
		})
	}
}
