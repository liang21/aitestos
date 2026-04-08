// Package project provides module repository implementation
package project

import (
	"context"
	"database/sql"
	"fmt"

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
		INSERT INTO modules (id, project_id, name, abbreviation, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		module.ID(),
		module.ProjectID(),
		module.Name(),
		module.Abbreviation().String(),
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
		ID           uuid.UUID `db:"id"`
		ProjectID    uuid.UUID `db:"project_id"`
		Name         string    `db:"name"`
		Abbreviation string    `db:"abbreviation"`
		Description  string    `db:"description"`
		CreatedAt    string    `db:"created_at"`
		UpdatedAt    string    `db:"updated_at"`
	}

	query := `
		SELECT id, project_id, name, abbreviation, description, created_at, updated_at
		FROM modules
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainproject.ErrModuleNotFound
		}
		return nil, fmt.Errorf("find module by id: %w", err)
	}

	abbrev, err := domainproject.ParseModuleAbbreviation(row.Abbreviation)
	if err != nil {
		return nil, fmt.Errorf("parse module abbreviation: %w", err)
	}

	return domainproject.ReconstructModule(
		row.ID,
		row.ProjectID,
		row.Name,
		abbrev,
		row.Description,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

// FindByProjectID retrieves all modules for a project
func (r *ModuleRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*domainproject.Module, error) {
	query := `
		SELECT id, project_id, name, abbreviation, description, created_at, updated_at
		FROM modules
		WHERE project_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	var rows []struct {
		ID           uuid.UUID `db:"id"`
		ProjectID    uuid.UUID `db:"project_id"`
		Name         string    `db:"name"`
		Abbreviation string    `db:"abbreviation"`
		Description  string    `db:"description"`
		CreatedAt    string    `db:"created_at"`
		UpdatedAt    string    `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, projectID); err != nil {
		return nil, fmt.Errorf("find modules by project id: %w", err)
	}

	modules := make([]*domainproject.Module, 0, len(rows))
	for _, row := range rows {
		abbrev, err := domainproject.ParseModuleAbbreviation(row.Abbreviation)
		if err != nil {
			return nil, fmt.Errorf("parse module abbreviation: %w", err)
		}

		module := domainproject.ReconstructModule(
			row.ID,
			row.ProjectID,
			row.Name,
			abbrev,
			row.Description,
			parseTime(row.CreatedAt),
			parseTime(row.UpdatedAt),
		)
		modules = append(modules, module)
	}

	return modules, nil
}

// FindByAbbreviation retrieves a module by abbreviation within a project
func (r *ModuleRepository) FindByAbbreviation(ctx context.Context, projectID uuid.UUID, abbrev domainproject.ModuleAbbreviation) (*domainproject.Module, error) {
	var row struct {
		ID           uuid.UUID `db:"id"`
		ProjectID    uuid.UUID `db:"project_id"`
		Name         string    `db:"name"`
		Abbreviation string    `db:"abbreviation"`
		Description  string    `db:"description"`
		CreatedAt    string    `db:"created_at"`
		UpdatedAt    string    `db:"updated_at"`
	}

	query := `
		SELECT id, project_id, name, abbreviation, description, created_at, updated_at
		FROM modules
		WHERE project_id = $1 AND abbreviation = $2 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, projectID, abbrev.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainproject.ErrModuleNotFound
		}
		return nil, fmt.Errorf("find module by abbreviation: %w", err)
	}

	parsedAbbrev, err := domainproject.ParseModuleAbbreviation(row.Abbreviation)
	if err != nil {
		return nil, fmt.Errorf("parse module abbreviation: %w", err)
	}

	return domainproject.ReconstructModule(
		row.ID,
		row.ProjectID,
		row.Name,
		parsedAbbrev,
		row.Description,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

// Delete removes a module (soft delete)
func (r *ModuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE modules
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
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
