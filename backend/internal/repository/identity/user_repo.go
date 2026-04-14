// Package identity provides user repository implementation
package identity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/liang21/aitestos/internal/domain/identity"
)

// UserRepository implements identity.UserRepository interface
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save persists a new user
func (r *UserRepository) Save(ctx context.Context, user *identity.User) error {
	row := user.ToRow()
	query := `
		INSERT INTO users (id, username, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		row.ID,
		row.Username,
		row.Email,
		row.Password,
		row.Role,
		row.CreatedAt,
		row.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	var row identity.UserRow
	query := `
		SELECT id, username, email, password, role, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, identity.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return identity.FromRow(&row)
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*identity.User, error) {
	var row identity.UserRow
	query := `
		SELECT id, username, email, password, role, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, identity.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return identity.FromRow(&row)
}

// FindByUsername retrieves a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*identity.User, error) {
	var row identity.UserRow
	query := `
		SELECT id, username, email, password, role, created_at, updated_at, deleted_at
		FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &row, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, identity.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return identity.FromRow(&row)
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *identity.User) error {
	row := user.ToRow()
	query := `
		UPDATE users
		SET username = $2, email = $3, password = $4, role = $5, updated_at = $6
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
		row.ID,
		row.Username,
		row.Email,
		row.Password,
		row.Role,
		row.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return identity.ErrUserNotFound
	}
	return nil
}

// Delete removes a user (soft delete)
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return identity.ErrUserNotFound
	}
	return nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, opts identity.QueryOptions) ([]*identity.User, int64, error) {
	// Count total
	var total int64
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE deleted_at IS NULL
	`
	if opts.Keywords != "" {
		countQuery += " AND (username LIKE '%' || $1 || '%' OR email LIKE '%' || $1 || '%')"
	}
	if opts.Role != "" {
		if opts.Keywords != "" {
			countQuery += " AND role = $2"
		} else {
			countQuery += " AND role = $1"
		}
	}

	var countArgs []interface{}
	if opts.Keywords != "" {
		countArgs = append(countArgs, opts.Keywords)
	}
	if opts.Role != "" {
		countArgs = append(countArgs, string(opts.Role))
	}

	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// Query users
	query := `
		SELECT id, username, email, password, role, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
	`
	var args []interface{}
	argIdx := 1

	if opts.Keywords != "" {
		query += fmt.Sprintf(" AND (username LIKE '%%' || $%d || '%%' OR email LIKE '%%' || $%d || '%%')", argIdx, argIdx)
		args = append(args, opts.Keywords)
		argIdx++
	}
	if opts.Role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, string(opts.Role))
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

	var rows []identity.UserRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	users := make([]*identity.User, 0, len(rows))
	for _, row := range rows {
		user, err := identity.FromRow(&row)
		if err != nil {
			return nil, 0, fmt.Errorf("convert row to user: %w", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}
