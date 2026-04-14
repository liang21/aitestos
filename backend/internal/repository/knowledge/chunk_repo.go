// Package knowledge provides document chunk repository implementation
package knowledge

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domainknowledge "github.com/liang21/aitestos/internal/domain/knowledge"
)

// chunkRow is a common row struct for scanning document chunk data
type chunkRow struct {
	ID         uuid.UUID `db:"id"`
	DocumentID uuid.UUID `db:"document_id"`
	ProjectID  uuid.UUID `db:"project_id"`
	ChunkIndex int       `db:"chunk_index"`
	Content    string    `db:"content"`
	Embedding  []byte    `db:"embedding"`
	CreatedAt  string    `db:"created_at"`
}

// DocumentChunkRepository implements domainknowledge.DocumentChunkRepository interface
type DocumentChunkRepository struct {
	db *sqlx.DB
}

// NewDocumentChunkRepository creates a new document chunk repository
func NewDocumentChunkRepository(db *sqlx.DB) *DocumentChunkRepository {
	return &DocumentChunkRepository{db: db}
}

// Save persists a single document chunk
func (r *DocumentChunkRepository) Save(ctx context.Context, chunk *domainknowledge.DocumentChunk) error {
	query := `
		INSERT INTO document_chunks (id, document_id, project_id, chunk_index, content, embedding, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		chunk.ID(),
		chunk.DocumentID(),
		chunk.ProjectID(),
		chunk.ChunkIndex(),
		chunk.Content(),
		chunk.Embedding(),
		chunk.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save document chunk: %w", err)
	}
	return nil
}

// SaveBatch persists multiple document chunks
func (r *DocumentChunkRepository) SaveBatch(ctx context.Context, chunks []*domainknowledge.DocumentChunk) error {
	if len(chunks) == 0 {
		return nil
	}

	query := `
		INSERT INTO document_chunks (id, document_id, project_id, chunk_index, content, embedding, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, chunk := range chunks {
		_, err := r.db.ExecContext(ctx, query,
			chunk.ID(),
			chunk.DocumentID(),
			chunk.ProjectID(),
			chunk.ChunkIndex(),
			chunk.Content(),
			chunk.Embedding(),
			chunk.CreatedAt(),
		)
		if err != nil {
			return fmt.Errorf("save document chunk: %w", err)
		}
	}
	return nil
}

// FindByID retrieves a document chunk by ID
func (r *DocumentChunkRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainknowledge.DocumentChunk, error) {
	query := `
		SELECT id, document_id, project_id, chunk_index, content, embedding, created_at
		FROM document_chunks
		WHERE id = $1
	`
	var row chunkRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainknowledge.ErrChunkNotFound
		}
		return nil, fmt.Errorf("find document chunk by id: %w", err)
	}

	return domainknowledge.ReconstructDocumentChunk(
		row.ID,
		row.DocumentID,
		row.ProjectID,
		row.ChunkIndex,
		row.Content,
		row.Embedding,
		parseTime(row.CreatedAt),
	), nil
}

// FindByDocumentID retrieves all chunks for a document
func (r *DocumentChunkRepository) FindByDocumentID(ctx context.Context, documentID uuid.UUID) ([]*domainknowledge.DocumentChunk, error) {
	query := `
		SELECT id, document_id, project_id, chunk_index, content, embedding, created_at
		FROM document_chunks
		WHERE document_id = $1
		ORDER BY chunk_index ASC
	`

	var rows []chunkRow
	if err := r.db.SelectContext(ctx, &rows, query, documentID); err != nil {
		return nil, fmt.Errorf("find document chunks by document id: %w", err)
	}

	chunks := make([]*domainknowledge.DocumentChunk, 0, len(rows))
	for _, row := range rows {
		chunk := domainknowledge.ReconstructDocumentChunk(
			row.ID,
			row.DocumentID,
			row.ProjectID,
			row.ChunkIndex,
			row.Content,
			row.Embedding,
			parseTime(row.CreatedAt),
		)
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

// FindByChunkIndex retrieves a chunk by document ID and chunk index
func (r *DocumentChunkRepository) FindByChunkIndex(ctx context.Context, documentID uuid.UUID, chunkIndex int) (*domainknowledge.DocumentChunk, error) {
	query := `
		SELECT id, document_id, project_id, chunk_index, content, embedding, created_at
		FROM document_chunks
		WHERE document_id = $1 AND chunk_index = $2
	`
	var row chunkRow
	err := r.db.GetContext(ctx, &row, query, documentID, chunkIndex)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainknowledge.ErrChunkNotFound
		}
		return nil, fmt.Errorf("find document chunk by index: %w", err)
	}

	return domainknowledge.ReconstructDocumentChunk(
		row.ID,
		row.DocumentID,
		row.ProjectID,
		row.ChunkIndex,
		row.Content,
		row.Embedding,
		parseTime(row.CreatedAt),
	), nil
}

// Update updates an existing document chunk
func (r *DocumentChunkRepository) Update(ctx context.Context, chunk *domainknowledge.DocumentChunk) error {
	query := `
		UPDATE document_chunks
		SET content = $2, embedding = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		chunk.ID(),
		chunk.Content(),
		chunk.Embedding(),
	)
	if err != nil {
		return fmt.Errorf("update document chunk: %w", err)
	}
	return nil
}

// DeleteByDocumentID removes all chunks for a document
func (r *DocumentChunkRepository) DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error {
	query := `DELETE FROM document_chunks WHERE document_id = $1`
	_, err := r.db.ExecContext(ctx, query, documentID)
	if err != nil {
		return fmt.Errorf("delete document chunks by document id: %w", err)
	}
	return nil
}

// CountByDocumentID counts chunks for a document
func (r *DocumentChunkRepository) CountByDocumentID(ctx context.Context, documentID uuid.UUID) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM document_chunks WHERE document_id = $1`
	err := r.db.GetContext(ctx, &count, query, documentID)
	if err != nil {
		return 0, fmt.Errorf("count document chunks: %w", err)
	}
	return count, nil
}
