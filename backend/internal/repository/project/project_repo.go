// Package project provides project repository implementation
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

// ProjectRepository implements domainproject.ProjectRepository interface
type ProjectRepository struct {
	db *sqlx.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Save persists a new project
func (r *ProjectRepository) Save(ctx context.Context, project *domainproject.Project) error {
	query := `
		INSERT INTO project (id, name, prefix, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		project.ID(),
		project.Name(),
		project.Prefix().String(),
		project.Description(),
		project.CreatedAt(),
		project.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save project: %w", err)
	}
	return nil
}

// FindByID retrieves a project by ID
func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainproject.Project, error) {
	var row struct {
		ID          uuid.UUID `db:"id"`
		Name        string    `db:"name"`
		Prefix      string    `db:"prefix"`
		Description string    `db:"description"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	query := `
		SELECT id, name, prefix, description, created_at, updated_at
		FROM project
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainproject.ErrProjectNotFound
		}
		return nil, fmt.Errorf("find project by id: %w", err)
	}

	prefix, err := domainproject.ParseProjectPrefix(row.Prefix)
	if err != nil {
		return nil, fmt.Errorf("parse project prefix: %w", err)
	}

	// Reconstruct project using Reconstruct method
	return domainproject.Reconstruct(
		row.ID,
		row.Name,
		prefix,
		row.Description,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

// FindByName retrieves a project by name
func (r *ProjectRepository) FindByName(ctx context.Context, name string) (*domainproject.Project, error) {
	var row struct {
		ID          uuid.UUID `db:"id"`
		Name        string    `db:"name"`
		Prefix      string    `db:"prefix"`
		Description string    `db:"description"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	query := `
		SELECT id, name, prefix, description, created_at, updated_at
		FROM project
		WHERE name = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainproject.ErrProjectNotFound
		}
		return nil, fmt.Errorf("find project by name: %w", err)
	}

	prefix, err := domainproject.ParseProjectPrefix(row.Prefix)
	if err != nil {
		return nil, fmt.Errorf("parse project prefix: %w", err)
	}

	return domainproject.Reconstruct(
		row.ID,
		row.Name,
		prefix,
		row.Description,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

// FindByPrefix retrieves a project by prefix
func (r *ProjectRepository) FindByPrefix(ctx context.Context, prefix domainproject.ProjectPrefix) (*domainproject.Project, error) {
	var row struct {
		ID          uuid.UUID `db:"id"`
		Name        string    `db:"name"`
		Prefix      string    `db:"prefix"`
		Description string    `db:"description"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	query := `
		SELECT id, name, prefix, description, created_at, updated_at
		FROM project
		WHERE prefix = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, prefix.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainproject.ErrProjectNotFound
		}
		return nil, fmt.Errorf("find project by prefix: %w", err)
	}

	parsedPrefix, err := domainproject.ParseProjectPrefix(row.Prefix)
	if err != nil {
		return nil, fmt.Errorf("parse project prefix: %w", err)
	}

	return domainproject.Reconstruct(
		row.ID,
		row.Name,
		parsedPrefix,
		row.Description,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	), nil
}

// FindAll retrieves all projects with pagination
func (r *ProjectRepository) FindAll(ctx context.Context, opts domainproject.QueryOptions) ([]*domainproject.Project, error) {
	query := `
		SELECT id, name, prefix, description, created_at, updated_at
		FROM project
		WHERE deleted_at IS NULL
	`
	var args []interface{}
	argIdx := 1

	if opts.Keywords != "" {
		query += fmt.Sprintf(" AND (name LIKE '%%' || $%d || '%%' OR description LIKE '%%' || $%d || '%%')", argIdx, argIdx)
		args = append(args, opts.Keywords)
		argIdx++
	}

	if opts.OrderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", opts.OrderBy)
	} else {
		query += " ORDER BY created_at DESC"
	}

	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
		args = append(args, opts.Limit, opts.Offset)
	}

	var rows []struct {
		ID          uuid.UUID `db:"id"`
		Name        string    `db:"name"`
		Prefix      string    `db:"prefix"`
		Description string    `db:"description"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("find all projects: %w", err)
	}

	projects := make([]*domainproject.Project, 0, len(rows))
	for _, row := range rows {
		prefix, err := domainproject.ParseProjectPrefix(row.Prefix)
		if err != nil {
			return nil, fmt.Errorf("parse project prefix: %w", err)
		}

		project := domainproject.Reconstruct(
			row.ID,
			row.Name,
			prefix,
			row.Description,
			parseTime(row.CreatedAt),
			parseTime(row.UpdatedAt),
		)
		projects = append(projects, project)
	}

	return projects, nil
}

// Update updates an existing project
func (r *ProjectRepository) Update(ctx context.Context, project *domainproject.Project) error {
	query := `
		UPDATE project
		SET name = $2, description = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
		project.ID(),
		project.Name(),
		project.Description(),
		project.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainproject.ErrProjectNotFound
	}
	return nil
}

// Delete removes a project (soft delete)
func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE project
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainproject.ErrProjectNotFound
	}
	return nil
}

// GetStatistics retrieves aggregated project statistics
func (r *ProjectRepository) GetStatistics(ctx context.Context, id uuid.UUID) (*domainproject.ProjectStatistics, error) {
	var stats domainproject.ProjectStatistics

	query := `
		SELECT
			(SELECT COUNT(*) FROM modules WHERE project_id = $1 AND deleted_at IS NULL) AS module_count,
			(SELECT COUNT(*) FROM test_cases tc
			 JOIN modules m ON tc.module_id = m.id
				WHERE m.project_id = $1 AND tc.deleted_at IS NULL) AS case_count,
			(SELECT COUNT(*) FROM documents WHERE project_id = $1 AND deleted_at IS NULL) AS document_count,
			COALESCE(
				(SELECT COUNT(*) * 100.0 / NULLIF(
					(SELECT COUNT(*) FROM test_cases tc
					 JOIN modules m ON tc.module_id = m.id
						WHERE m.project_id = $1 AND tc.deleted_at IS NULL AND tc.status != 'unexecuted'), 0)
					FROM test_cases tc
					JOIN modules m ON tc.module_id = m.id
					WHERE m.project_id = $1 AND tc.status = 'pass' AND tc.deleted_at IS NULL), 0) AS pass_rate,
			(SELECT COUNT(*) FROM test_cases tc
				JOIN modules m ON tc.module_id = m.id
				WHERE m.project_id = $1
					AND tc.ai_metadata->>'generation_task_id' IS NOT NULL
					AND tc.deleted_at IS NULL) AS ai_generated_count
	`

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&stats.ModuleCount,
		&stats.CaseCount,
		&stats.DocumentCount,
		&stats.PassRate,
		&stats.AIGeneratedCount,
	)
	if err != nil {
		return nil, fmt.Errorf("get project statistics: %w", err)
	}

	// Calculate coverage rate (cases per module)
	stats.CoverageRate = calculateCoverageRate(stats.CaseCount, stats.ModuleCount)
	stats.UpdatedAt = time.Now()

	return &stats, nil
}

// SetStatistics stores statistics in cache (no-op for base repository)
// This method is implemented by CachedProjectRepository for cache warming
func (r *ProjectRepository) SetStatistics(ctx context.Context, id uuid.UUID, stats *domainproject.ProjectStatistics) error {
	// Base repository doesn't support caching, return nil (no-op)
	return nil
}

// calculateCoverageRate calculates coverage rate as cases per module
func calculateCoverageRate(caseCount, moduleCount int64) float64 {
	if moduleCount == 0 {
		return 0
	}
	// Average cases per module as a simple coverage metric
	return float64(caseCount) / float64(moduleCount)
}
