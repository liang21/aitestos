// Package knowledge provides document management services
package knowledge

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/knowledge"
)

// UploadDocumentRequest contains document upload data
type UploadDocumentRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=2,max=255"`
	Type      string    `json:"type" validate:"required,oneof=prd figma api_spec swagger markdown"`
	File      io.Reader `json:"-"`
	FileSize  int64     `json:"file_size"`
	UserID    uuid.UUID `json:"-"`
}

// DocumentDetail contains document info with chunk statistics
type DocumentDetail struct {
	*knowledge.Document
	ChunkCount int64 `json:"chunk_count"`
}

// DocumentListOptions contains pagination and filtering options
type DocumentListOptions struct {
	Offset    int       `json:"offset"`
	Limit     int       `json:"limit"`
	ProjectID uuid.UUID `json:"project_id,omitempty"`
	Type      string    `json:"type,omitempty"`
	Status    string    `json:"status,omitempty"`
}

// ChunkInfo contains document chunk information
type ChunkInfo struct {
	ID         uuid.UUID `json:"id"`
	DocumentID uuid.UUID `json:"document_id"`
	ChunkIndex int       `json:"chunk_index"`
	Content    string    `json:"content"`
	CreatedAt  int64     `json:"created_at"`
}

// DocumentService provides document management operations
type DocumentService interface {
	// Document management
	UploadDocument(ctx context.Context, req *UploadDocumentRequest) (*knowledge.Document, error)
	GetDocument(ctx context.Context, id uuid.UUID) (*DocumentDetail, error)
	ListDocuments(ctx context.Context, opts DocumentListOptions) ([]*knowledge.Document, int64, error)
	DeleteDocument(ctx context.Context, id uuid.UUID) error

	// Status management
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	ProcessDocument(ctx context.Context, id uuid.UUID) error

	// Chunk management
	GetChunks(ctx context.Context, documentID uuid.UUID) ([]*ChunkInfo, error)
}

// DocumentServiceImpl implements DocumentService
type DocumentServiceImpl struct {
	docRepo    knowledge.DocumentRepository
	chunkRepo  knowledge.DocumentChunkRepository
	vectorRepo knowledge.VectorRepository
}

// NewDocumentService creates a new DocumentService instance
func NewDocumentService(
	docRepo knowledge.DocumentRepository,
	chunkRepo knowledge.DocumentChunkRepository,
	vectorRepo knowledge.VectorRepository,
) DocumentService {
	return &DocumentServiceImpl{
		docRepo:    docRepo,
		chunkRepo:  chunkRepo,
		vectorRepo: vectorRepo,
	}
}

// UploadDocument uploads a new document
func (s *DocumentServiceImpl) UploadDocument(ctx context.Context, req *UploadDocumentRequest) (*knowledge.Document, error) {
	// Validate document type
	docType, err := knowledge.ParseDocumentType(req.Type)
	if err != nil {
		return nil, errors.New("invalid document type")
	}

	// Validate project ID
	if req.ProjectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}

	// Validate user ID
	if req.UserID == uuid.Nil {
		return nil, errors.New("user ID cannot be nil")
	}

	// Validate name
	if req.Name == "" {
		return nil, errors.New("document name cannot be empty")
	}

	// Generate URL (in production, this would be from storage service)
	url := fmt.Sprintf("https://storage.example.com/documents/%s/%s", req.ProjectID, req.Name)

	// Create document
	doc, err := knowledge.NewDocument(req.ProjectID, req.Name, docType, url, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	// Save document
	if err := s.docRepo.Save(ctx, doc); err != nil {
		return nil, fmt.Errorf("save document: %w", err)
	}

	// TODO: Upload file to storage and process chunks
	// This would involve:
	// 1. Uploading file to object storage (MinIO/S3)
	// 2. Extracting text content
	// 3. Splitting into chunks
	// 4. Generating embeddings
	// 5. Storing chunks and vectors

	return doc, nil
}

// GetDocument retrieves document details with statistics
func (s *DocumentServiceImpl) GetDocument(ctx context.Context, id uuid.UUID) (*DocumentDetail, error) {
	doc, err := s.docRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find document: %w", err)
	}

	// Get chunk count
	chunks, err := s.chunkRepo.FindByDocumentID(ctx, id)
	chunkCount := int64(0)
	if err == nil {
		chunkCount = int64(len(chunks))
	}

	return &DocumentDetail{
		Document:   doc,
		ChunkCount: chunkCount,
	}, nil
}

// ListDocuments lists documents with pagination
func (s *DocumentServiceImpl) ListDocuments(ctx context.Context, opts DocumentListOptions) ([]*knowledge.Document, int64, error) {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Limit > 100 {
		opts.Limit = 100
	}

	queryOpts := knowledge.QueryOptions{
		Offset: opts.Offset,
		Limit:  opts.Limit,
	}

	if opts.ProjectID != uuid.Nil {
		docs, err := s.docRepo.FindByProjectID(ctx, opts.ProjectID, queryOpts)
		if err != nil {
			return nil, 0, fmt.Errorf("list documents: %w", err)
		}
		return docs, int64(len(docs)), nil
	}

	// If no project ID, return empty list
	return []*knowledge.Document{}, 0, nil
}

// DeleteDocument deletes a document and its chunks
func (s *DocumentServiceImpl) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	// Check if document exists
	_, err := s.docRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find document: %w", err)
	}

	// Delete chunks first
	if err := s.chunkRepo.DeleteByDocumentID(ctx, id); err != nil {
		return fmt.Errorf("delete chunks: %w", err)
	}

	// Delete document
	if err := s.docRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete document: %w", err)
	}

	return nil
}

// UpdateStatus updates document processing status
func (s *DocumentServiceImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	doc, err := s.docRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find document: %w", err)
	}

	newStatus, err := knowledge.ParseDocumentStatus(status)
	if err != nil {
		return errors.New("invalid document status")
	}

	if err := doc.UpdateStatus(newStatus); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	if err := s.docRepo.Update(ctx, doc); err != nil {
		return fmt.Errorf("save document: %w", err)
	}

	return nil
}

// ProcessDocument starts document processing (chunking and embedding)
func (s *DocumentServiceImpl) ProcessDocument(ctx context.Context, id uuid.UUID) error {
	doc, err := s.docRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find document: %w", err)
	}

	// Update status to processing
	if err := doc.UpdateStatus(knowledge.StatusProcessing); err != nil {
		return fmt.Errorf("update status to processing: %w", err)
	}
	if err := s.docRepo.Update(ctx, doc); err != nil {
		return fmt.Errorf("save document: %w", err)
	}

	// TODO: Implement actual processing
	// 1. Download document from storage
	// 2. Extract text content
	// 3. Split into chunks
	// 4. Generate embeddings for each chunk
	// 5. Store chunks in database
	// 6. Store vectors in Milvus
	// 7. Update status to completed

	// For now, just mark as completed
	if err := doc.UpdateStatus(knowledge.StatusCompleted); err != nil {
		return fmt.Errorf("update status to completed: %w", err)
	}
	if err := s.docRepo.Update(ctx, doc); err != nil {
		return fmt.Errorf("save document: %w", err)
	}

	return nil
}

// GetChunks retrieves chunks for a document
func (s *DocumentServiceImpl) GetChunks(ctx context.Context, documentID uuid.UUID) ([]*ChunkInfo, error) {
	chunks, err := s.chunkRepo.FindByDocumentID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("get chunks: %w", err)
	}

	result := make([]*ChunkInfo, 0, len(chunks))
	for _, chunk := range chunks {
		result = append(result, &ChunkInfo{
			ID:         chunk.ID(),
			DocumentID: chunk.DocumentID(),
			ChunkIndex: chunk.ChunkIndex(),
			Content:    chunk.Content(),
			CreatedAt:  chunk.CreatedAt().Unix(),
		})
	}

	return result, nil
}
