// Package knowledge defines Document aggregate
package knowledge

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Document is the aggregate root for knowledge context
type Document struct {
	id          uuid.UUID
	projectID   uuid.UUID
	name        string
	docType     DocumentType
	url         string
	contentText string
	status      DocumentStatus
	createdBy   uuid.UUID
	createdAt   time.Time
	updatedAt   time.Time
}

// NewDocument creates a new document
func NewDocument(projectID uuid.UUID, name string, docType DocumentType, url string, userID uuid.UUID) (*Document, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}
	if name == "" {
		return nil, errors.New("document name cannot be empty")
	}
	if url == "" {
		return nil, errors.New("document URL cannot be empty")
	}
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be nil")
	}

	now := time.Now()
	return &Document{
		id:        uuid.New(),
		projectID: projectID,
		name:      name,
		docType:   docType,
		url:       url,
		status:    StatusPending,
		createdBy: userID,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns the document's unique identifier
func (d *Document) ID() uuid.UUID {
	return d.id
}

// ProjectID returns the associated project's ID
func (d *Document) ProjectID() uuid.UUID {
	return d.projectID
}

// Name returns the document's name
func (d *Document) Name() string {
	return d.name
}

// Type returns the document's type
func (d *Document) Type() DocumentType {
	return d.docType
}

// URL returns the document's URL
func (d *Document) URL() string {
	return d.url
}

// ContentText returns the extracted text content
func (d *Document) ContentText() string {
	return d.contentText
}

// UpdateContentText updates the extracted text content
func (d *Document) UpdateContentText(text string) {
	d.contentText = text
	d.updatedAt = time.Now()
}

// Status returns the document's processing status
func (d *Document) Status() DocumentStatus {
	return d.status
}

// CreatedBy returns the user who created this document
func (d *Document) CreatedBy() uuid.UUID {
	return d.createdBy
}

// CreatedAt returns the creation timestamp
func (d *Document) CreatedAt() time.Time {
	return d.createdAt
}

// UpdatedAt returns the last update timestamp
func (d *Document) UpdatedAt() time.Time {
	return d.updatedAt
}

// UpdateStatus updates the document processing status
func (d *Document) UpdateStatus(status DocumentStatus) error {
	if !d.status.CanTransitionTo(status) {
		return errors.New("invalid status transition")
	}
	d.status = status
	d.updatedAt = time.Now()
	return nil
}

// ReconstructDocument reconstructs a Document from stored data
func ReconstructDocument(
	id uuid.UUID,
	projectID uuid.UUID,
	name string,
	docType DocumentType,
	url string,
	contentText string,
	status DocumentStatus,
	createdBy uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
) *Document {
	return &Document{
		id:          id,
		projectID:   projectID,
		name:        name,
		docType:     docType,
		url:         url,
		contentText: contentText,
		status:      status,
		createdBy:   createdBy,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// DocumentChunk represents a chunk of a document for vector storage
type DocumentChunk struct {
	id         uuid.UUID
	documentID uuid.UUID
	projectID  uuid.UUID
	chunkIndex int
	content    string
	embedding  []byte
	createdAt  time.Time
}

// NewDocumentChunk creates a new document chunk
func NewDocumentChunk(documentID uuid.UUID, projectID uuid.UUID, chunkIndex int, content string) (*DocumentChunk, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID cannot be nil")
	}
	return &DocumentChunk{
		id:         uuid.New(),
		documentID: documentID,
		projectID:  projectID,
		chunkIndex: chunkIndex,
		content:    content,
		createdAt:  time.Now(),
	}, nil
}

// ID returns the chunk's unique identifier
func (c *DocumentChunk) ID() uuid.UUID {
	return c.id
}

// DocumentID returns the associated document's ID
func (c *DocumentChunk) DocumentID() uuid.UUID {
	return c.documentID
}

// ProjectID returns the associated project's ID
func (c *DocumentChunk) ProjectID() uuid.UUID {
	return c.projectID
}

// ChunkIndex returns the chunk's index in the document
func (c *DocumentChunk) ChunkIndex() int {
	return c.chunkIndex
}

// Content returns the chunk's content
func (c *DocumentChunk) Content() string {
	return c.content
}

// Embedding returns the chunk's embedding vector
func (c *DocumentChunk) Embedding() []byte {
	return c.embedding
}

// SetEmbedding sets the chunk's embedding vector
func (c *DocumentChunk) SetEmbedding(embedding []byte) {
	c.embedding = embedding
}

// UpdateContent updates the chunk's content
func (c *DocumentChunk) UpdateContent(content string) {
	c.content = content
}

// CreatedAt returns the creation timestamp
func (c *DocumentChunk) CreatedAt() time.Time {
	return c.createdAt
}

// ReconstructDocumentChunk reconstructs a DocumentChunk from stored data
func ReconstructDocumentChunk(
	id uuid.UUID,
	documentID uuid.UUID,
	projectID uuid.UUID,
	chunkIndex int,
	content string,
	embedding []byte,
	createdAt time.Time,
) *DocumentChunk {
	return &DocumentChunk{
		id:         id,
		documentID: documentID,
		projectID:  projectID,
		chunkIndex: chunkIndex,
		content:    content,
		embedding:  embedding,
		createdAt:  createdAt,
	}
}
