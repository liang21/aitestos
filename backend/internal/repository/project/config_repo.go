// Package project provides project config repository implementation
package project

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
)

// ProjectConfigRepository implements domainproject.ProjectConfigRepository interface
type ProjectConfigRepository struct {
	db *sqlx.DB
}

// NewProjectConfigRepository creates a new project config repository
func NewProjectConfigRepository(db *sqlx.DB) *ProjectConfigRepository {
	return &ProjectConfigRepository{db: db}
}

// Save persists a project configuration (upsert)
func (r *ProjectConfigRepository) Save(ctx context.Context, config *domainproject.ProjectConfig) error {
	query := `
		INSERT INTO project_configs (id, project_id, key, value, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (project_id, key) DO UPDATE SET
			value = EXCLUDED.value,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at
	`

	valueJSON, err := toJSON(config.Value())
	if err != nil {
		return fmt.Errorf("serialize config value: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		config.ID(),
		config.ProjectID(),
		config.Key(),
		valueJSON,
		config.Description(),
		config.CreatedAt(),
		config.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save project config: %w", err)
	}
	return nil
}

// FindByProjectID retrieves all configurations for a project
func (r *ProjectConfigRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*domainproject.ProjectConfig, error) {
	query := `
		SELECT id, project_id, key, value, description, created_at, updated_at
		FROM project_configs
		WHERE project_id = $1
		ORDER BY key ASC
	`

	var rows []struct {
		ID          uuid.UUID `db:"id"`
		ProjectID   uuid.UUID `db:"project_id"`
		Key         string    `db:"key"`
		Value       string    `db:"value"`
		Description string    `db:"description"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, projectID); err != nil {
		return nil, fmt.Errorf("find project configs by project id: %w", err)
	}

	configs := make([]*domainproject.ProjectConfig, 0, len(rows))
	for _, row := range rows {
		var value map[string]any
		if err := fromJSON(row.Value, &value); err != nil {
			return nil, fmt.Errorf("deserialize config value: %w", err)
		}

		config, err := domainproject.ReconstructProjectConfig(
			row.ID,
			row.ProjectID,
			row.Key,
			value,
			row.Description,
			parseTime(row.CreatedAt),
			parseTime(row.UpdatedAt),
		)
		if err != nil {
			return nil, fmt.Errorf("reconstruct project config: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// FindByKey retrieves a configuration by project ID and key
func (r *ProjectConfigRepository) FindByKey(ctx context.Context, projectID uuid.UUID, key string) (*domainproject.ProjectConfig, error) {
	var row struct {
		ID          uuid.UUID `db:"id"`
		ProjectID   uuid.UUID `db:"project_id"`
		Key         string    `db:"key"`
		Value       string    `db:"value"`
		Description string    `db:"description"`
		CreatedAt   string    `db:"created_at"`
		UpdatedAt   string    `db:"updated_at"`
	}

	query := `
		SELECT id, project_id, key, value, description, created_at, updated_at
		FROM project_configs
		WHERE project_id = $1 AND key = $2
	`
	err := r.db.GetContext(ctx, &row, query, projectID, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainproject.ErrConfigNotFound
		}
		return nil, fmt.Errorf("find project config by key: %w", err)
	}

	var value map[string]any
	if err := fromJSON(row.Value, &value); err != nil {
		return nil, fmt.Errorf("deserialize config value: %w", err)
	}

	return domainproject.ReconstructProjectConfig(
		row.ID,
		row.ProjectID,
		row.Key,
		value,
		row.Description,
		parseTime(row.CreatedAt),
		parseTime(row.UpdatedAt),
	)
}

// Delete removes a configuration
func (r *ProjectConfigRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM project_configs WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete project config: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainproject.ErrConfigNotFound
	}
	return nil
}

// Update updates an existing configuration
func (r *ProjectConfigRepository) Update(ctx context.Context, config *domainproject.ProjectConfig) error {
	valueJSON, err := toJSON(config.Value())
	if err != nil {
		return fmt.Errorf("serialize config value: %w", err)
	}

	query := `
		UPDATE project_configs
		SET value = $1, description = $2, updated_at = $3
		WHERE id = $4
	`
	result, err := r.db.ExecContext(ctx, query,
		valueJSON,
		config.Description(),
		config.UpdatedAt(),
		config.ID(),
	)
	if err != nil {
		return fmt.Errorf("update project config: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domainproject.ErrConfigNotFound
	}
	return nil
}

// BatchUpsert batch upserts configurations
func (r *ProjectConfigRepository) BatchUpsert(ctx context.Context, configs []*domainproject.ProjectConfig) error {
	if len(configs) == 0 {
		return nil
	}

	query := `
		INSERT INTO project_configs (id, project_id, key, value, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (project_id, key)
		DO UPDATE SET value = EXCLUDED.value, description = EXCLUDED.description, updated_at = NOW()
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, config := range configs {
		valueJSON, err := toJSON(config.Value())
		if err != nil {
			return fmt.Errorf("serialize config value: %w", err)
		}

		_, err = stmt.ExecContext(ctx, config.ID(), config.ProjectID(), config.Key(), valueJSON, config.Description())
		if err != nil {
			return fmt.Errorf("batch upsert config [%s]: %w", config.Key(), err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// ExportConfigs exports all configs for a project as JSON-compatible format
func (r *ProjectConfigRepository) ExportConfigs(ctx context.Context, projectID uuid.UUID) ([]map[string]any, error) {
	configs, err := r.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]any, len(configs))
	for i, config := range configs {
		result[i] = map[string]any{
			"key":         config.Key(),
			"value":       config.Value(),
			"description": config.Description(),
		}
	}

	return result, nil
}
