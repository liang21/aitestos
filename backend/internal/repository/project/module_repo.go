// Package project provides module repository implementation
package project

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
)

// ModuleRepository implements domainproject.ModuleRepository interface
type ModuleRepository struct {
	db *sqlx.DB
}

// NewModuleRepository creates a new module repository
func NewModuleRepository(db *sqlx.DB) *ModuleRepository {
	return &ModuleRepository{db: db}
}

// Save persists a new module
func (r *ModuleRepository) Save(ctx context.Context, module *domainproject.Module) error {
	query := `
		INSERT INTO modules (id, project_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (project_id, name) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.ExecContext(ctx, query,
		module.ID(),
		module.ProjectID(),
		module.Name(),
		module.Description(),
		module.CreatedAt(),
		module.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save module: %w", err)
	}
	return nil
}

// FindByID retrieves a module by ID
func (r *ModuleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainproject.Module, error) {
	var row struct {
		ID          uuid.UUID `db:"id"`
		ProjectID   uuid.UUID `db:"project_id"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		CreatedBy   string    `db:"created_by"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	query := `
		SELECT id, project_id, name, description, created_by, created_at, updated_at
		FROM modules
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainproject.ErrModuleNotFound
		}
		return nil, fmt.Errorf("find module by id: %w", err)
	}

	return domainproject.ReconstructModule(
		row.ID,
		row.ProjectID,
		row.Name,
		row.Description,
		parseUUID(row.CreatedBy),
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

// FindByProjectID retrieves all modules for a project
func (r *ModuleRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*domainproject.Module, error) {
	query := `
		SELECT id, project_id, name, description, created_by, created_at, updated_at
		FROM modules
		WHERE project_id = $1
		ORDER BY created_at ASC
	`

	var rows []struct {
		ID          uuid.UUID `db:"id"`
		ProjectID   uuid.UUID `db:"project_id"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		CreatedBy   string    `db:"created_by"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, projectID); err != nil {
		return nil, fmt.Errorf("find modules by project id: %w", err)
	}

	modules := make([]*domainproject.Module, 0, len(rows))
	for _, row := range rows {
		module := domainproject.ReconstructModule(
			row.ID,
			row.ProjectID,
			row.Name,
			row.Description,
			parseUUID(row.CreatedBy),
			parseTime(row.CreatedAt),
			parseTime(row.UpdatedAt),
		)
		modules = append(modules, module)
	}

	return modules, nil
}

// FindByName retrieves a module by name within a project
func (r *ModuleRepository) FindByName(ctx context.Context, projectID uuid.UUID, name string) (*domainproject.Module, error) {
	var row struct {
		ID          uuid.UUID `db:"id"`
		ProjectID   uuid.UUID `db:"project_id"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		CreatedBy   string    `db:"created_by"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	query := `
		SELECT id, project_id, name, description, created_by, created_at, updated_at
		FROM modules
		WHERE project_id = $1 AND name = $2
	`
	err := r.db.GetContext(ctx, &row, query, projectID, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainproject.ErrModuleNotFound
		}
		return nil, fmt.Errorf("find module by name: %w", err)
	}

	return domainproject.ReconstructModule(
		row.ID,
		row.ProjectID,
		row.Name,
		row.Description,
		parseUUID(row.CreatedBy),
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

// Delete removes a module (hard delete)
func (r *ModuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM modules WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete module: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainproject.ErrModuleNotFound
	}
	return nil
}

// Update updates an existing module
func (r *ModuleRepository) Update(ctx context.Context, module *domainproject.Module) error {
	query := `
		UPDATE modules
		SET name = $2, description = $3, updated_at = $4
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		module.ID(),
		module.Name(),
		module.Description(),
		module.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("update module: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainproject.ErrModuleNotFound
	}
	return nil
}

// parseTime parses a time string from the database
func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339Nano, s)
	return t
}

// parseUUID parses a UUID string from the database
func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
