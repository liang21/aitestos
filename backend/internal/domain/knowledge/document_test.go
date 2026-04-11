// Package knowledge_test tests Document aggregate
package knowledge_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/knowledge"
)

func TestNewDocument(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		projectID uuid.UUID
		docName   string
		docType   knowledge.DocumentType
		url       string
		userID    uuid.UUID
		wantErr   bool
	}{
		{
			name:      "valid PRD document",
			projectID: projectID,
			docName:   "Product Requirements",
			docType:   knowledge.TypePRD,
			url:       "https://example.com/prd.pdf",
			userID:    userID,
			wantErr:   false,
		},
		{
			name:      "valid Figma document",
			projectID: projectID,
			docName:   "UI Design",
			docType:   knowledge.TypeFigma,
			url:       "https://figma.com/design/abc",
			userID:    userID,
			wantErr:   false,
		},
		{
			name:      "empty name",
			projectID: projectID,
			docName:   "",
			docType:   knowledge.TypePRD,
			url:       "https://example.com/doc.pdf",
			userID:    userID,
			wantErr:   true,
		},
		{
			name:      "nil project ID",
			projectID: uuid.Nil,
			docName:   "Test Document",
			docType:   knowledge.TypePRD,
			url:       "https://example.com/doc.pdf",
			userID:    userID,
			wantErr:   true,
		},
		{
			name:      "empty url",
			projectID: projectID,
			docName:   "Test Document",
			docType:   knowledge.TypePRD,
			url:       "",
			userID:    userID,
			wantErr:   true,
		},
		{
			name:      "nil user ID",
			projectID: projectID,
			docName:   "Test Document",
			docType:   knowledge.TypePRD,
			url:       "https://example.com/doc.pdf",
			userID:    uuid.Nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := knowledge.NewDocument(tt.projectID, tt.docName, tt.docType, tt.url, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewDocument() returned nil document")
					return
				}
				if got.Name() != tt.docName {
					t.Errorf("Document.Name() = %v, want %v", got.Name(), tt.docName)
				}
				if got.Type() != tt.docType {
					t.Errorf("Document.Type() = %v, want %v", got.Type(), tt.docType)
				}
				if got.Status() != knowledge.StatusPending {
					t.Errorf("Document.Status() = %v, want pending", got.Status())
				}
			}
		})
	}
}

func TestDocument_Accessors(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()
	doc, err := knowledge.NewDocument(projectID, "Test Doc", knowledge.TypePRD, "https://example.com/doc.pdf", userID)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	if doc.ID() == uuid.Nil {
		t.Error("Document.ID() should not be nil")
	}
	if doc.ProjectID() != projectID {
		t.Errorf("Document.ProjectID() = %v, want %v", doc.ProjectID(), projectID)
	}
	if doc.Name() != "Test Doc" {
		t.Errorf("Document.Name() = %v, want Test Doc", doc.Name())
	}
	if doc.Type() != knowledge.TypePRD {
		t.Errorf("Document.Type() = %v, want prd", doc.Type())
	}
	if doc.URL() != "https://example.com/doc.pdf" {
		t.Errorf("Document.URL() = %v, want https://example.com/doc.pdf", doc.URL())
	}
	if doc.Status() != knowledge.StatusPending {
		t.Errorf("Document.Status() = %v, want pending", doc.Status())
	}
	if doc.CreatedAt().IsZero() {
		t.Error("Document.CreatedAt() should not be zero")
	}
	if doc.UpdatedAt().IsZero() {
		t.Error("Document.UpdatedAt() should not be zero")
	}
	if doc.CreatedBy() != userID {
		t.Errorf("Document.CreatedBy() = %v, want %v", doc.CreatedBy(), userID)
	}
}

func TestDocument_UpdateStatus(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()
	doc, err := knowledge.NewDocument(projectID, "Test Doc", knowledge.TypePRD, "https://example.com/doc.pdf", userID)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	originalUpdatedAt := doc.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	// Test valid transition: pending -> processing
	err = doc.UpdateStatus(knowledge.StatusProcessing)
	if err != nil {
		t.Errorf("Document.UpdateStatus() error = %v", err)
	}
	if doc.Status() != knowledge.StatusProcessing {
		t.Errorf("Document.Status() = %v, want processing", doc.Status())
	}
	if !doc.UpdatedAt().After(originalUpdatedAt) {
		t.Error("Document.UpdatedAt() should be updated")
	}

	// Test invalid transition: processing -> pending
	err = doc.UpdateStatus(knowledge.StatusPending)
	if err == nil {
		t.Error("Document.UpdateStatus() should return error for invalid transition")
	}
}

func TestDocumentChunk(t *testing.T) {
	docID := uuid.New()
	projectID := uuid.New()
	content := "This is a test chunk content"

	chunk, err := knowledge.NewDocumentChunk(docID, projectID, 0, content)
	if err != nil {
		t.Fatalf("NewDocumentChunk() error = %v", err)
	}

	if chunk.ID() == uuid.Nil {
		t.Error("DocumentChunk.ID() should not be nil")
	}
	if chunk.DocumentID() != docID {
		t.Errorf("DocumentChunk.DocumentID() = %v, want %v", chunk.DocumentID(), docID)
	}
	if chunk.ChunkIndex() != 0 {
		t.Errorf("DocumentChunk.ChunkIndex() = %v, want 0", chunk.ChunkIndex())
	}
	if chunk.Content() != content {
		t.Errorf("DocumentChunk.Content() = %v, want %v", chunk.Content(), content)
	}
	if chunk.CreatedAt().IsZero() {
		t.Error("DocumentChunk.CreatedAt() should not be zero")
	}
}
