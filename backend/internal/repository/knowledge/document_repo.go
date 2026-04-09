// Package knowledge provides document repository implementation
package knowledge

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domainknowledge "github.com/liang21/aitestos/internal/domain/knowledge"
)

// documentRow is a common row struct for scanning document data
type documentRow struct {
	ID          uuid.UUID `db:"id"`
	ProjectID   uuid.UUID `db:"project_id"`
	Name        string    `db:"name"`
	Type        string    `db:"type"`
	URL         string    `db:"url"`
	ContentText string    `db:"content_text"`
	Status      string    `db:"status"`
	CreatedBy   uuid.UUID `db:"created_by"`
	CreatedAt   string    `db:"created_at"`
	UpdatedAt   string    `db:"updated_at"`
}

// DocumentRepository implements domainknowledge.DocumentRepository interface
type DocumentRepository struct {
	db *sqlx.DB
}

// NewDocumentRepository creates a new document repository
func NewDocumentRepository(db *sqlx.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Save persists a new document
func (r *DocumentRepository) Save(ctx context.Context, doc *domainknowledge.Document) error {
	query := `
		INSERT INTO documents (id, project_id, name, type, url, content_text, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		doc.ID(),
		doc.ProjectID(),
		doc.Name(),
		string(doc.Type()),
		doc.URL(),
		doc.ContentText(),
		string(doc.Status()),
		doc.CreatedBy(),
		doc.CreatedAt(),
		doc.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save document: %w", err)
	}
	return nil
}

// FindByID retrieves a document by ID
func (r *DocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainknowledge.Document, error) {
	query := `
		SELECT id, project_id, name, type, url, content_text, status, created_by, created_at, updated_at
		FROM documents
		WHERE id = $1 AND deleted_at IS NULL
	`
	var row documentRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainknowledge.ErrDocumentNotFound
		}
		return nil, fmt.Errorf("find document by id: %w", err)
	}

	return rowToDocument(row)
}

// FindByProjectID retrieves all documents for a project with pagination
func (r *DocumentRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, opts domainknowledge.QueryOptions) ([]*domainknowledge.Document, error) {
	query := `
		SELECT id, project_id, name, type, url, content_text, status, created_by, created_at, updated_at
		FROM documents
		WHERE project_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var rows []documentRow
	if err := r.db.SelectContext(ctx, &rows, query, projectID, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find documents by project id: %w", err)
	}

	return rowsToDocuments(rows)
}

// FindByType retrieves documents by type for a project with pagination
func (r *DocumentRepository) FindByType(ctx context.Context, projectID uuid.UUID, docType domainknowledge.DocumentType, opts domainknowledge.QueryOptions) ([]*domainknowledge.Document, error) {
	query := `
		SELECT id, project_id, name, type, url, content_text, status, created_by, created_at, updated_at
		FROM documents
		WHERE project_id = $1 AND type = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	var rows []documentRow
	if err := r.db.SelectContext(ctx, &rows, query, projectID, string(docType), opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find documents by type: %w", err)
	}

	return rowsToDocuments(rows)
}

// FindByStatus retrieves documents by status with pagination
func (r *DocumentRepository) FindByStatus(ctx context.Context, status domainknowledge.DocumentStatus, opts domainknowledge.QueryOptions) ([]*domainknowledge.Document, error) {
	query := `
		SELECT id, project_id, name, type, url, content_text, status, created_by, created_at, updated_at
		FROM documents
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var rows []documentRow
	if err := r.db.SelectContext(ctx, &rows, query, string(status), opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("find documents by status: %w", err)
	}

	return rowsToDocuments(rows)
}

// Update updates an existing document
func (r *DocumentRepository) Update(ctx context.Context, doc *domainknowledge.Document) error {
	query := `
		UPDATE documents
		SET name = $2, status = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
		doc.ID(),
		doc.Name(),
		string(doc.Status()),
		doc.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("update document: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainknowledge.ErrDocumentNotFound
	}
	return nil
}

// UpdateStatus updates the document processing status
func (r *DocumentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domainknowledge.DocumentStatus) error {
	query := `UPDATE documents SET status = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id, string(status))
	if err != nil {
		return fmt.Errorf("update document status: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainknowledge.ErrDocumentNotFound
	}
	return nil
}

// UpdateContentText updates the document's extracted text content
func (r *DocumentRepository) UpdateContentText(ctx context.Context, id uuid.UUID, contentText string) error {
	query := `UPDATE documents SET content_text = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id, contentText)
	if err != nil {
		return fmt.Errorf("update document content text: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainknowledge.ErrDocumentNotFound
	}
	return nil
}

// Delete removes a document (soft delete)
func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE documents SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainknowledge.ErrDocumentNotFound
	}
	return nil
}

// CountByProjectID counts documents for a project
func (r *DocumentRepository) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM documents WHERE project_id = $1 AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &count, query, projectID)
	if err != nil {
		return 0, fmt.Errorf("count documents by project id: %w", err)
	}
	return count, nil
}

// Helper functions
func rowToDocument(row documentRow) (*domainknowledge.Document, error) {
	docType, err := domainknowledge.ParseDocumentType(row.Type)
	if err != nil {
		return nil, fmt.Errorf("parse document type: %w", err)
	}

	status, err := domainknowledge.ParseDocumentStatus(row.Status)
	if err != nil {
		return nil, fmt.Errorf("parse document status: %w", err)
	}

	return domainknowledge.ReconstructDocument(
		row.ID,
		row.ProjectID,
		row.Name,
		docType,
		row.URL,
		row.ContentText,
		status,
		row.CreatedBy,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

func rowsToDocuments(rows []documentRow) ([]*domainknowledge.Document, error) {
	docs := make([]*domainknowledge.Document, 0, len(rows))
	for _, row := range rows {
		doc, err := rowToDocument(row)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}
