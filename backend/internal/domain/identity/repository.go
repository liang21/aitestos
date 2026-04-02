// Package identity defines repository interfaces
package identity

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	// Save persists a new user
	Save(ctx context.Context, user *User) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByUsername retrieves a user by username
	FindByUsername(ctx context.Context, username string) (*User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *User) error

	// Delete removes a user (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves users with pagination
	List(ctx context.Context, opts QueryOptions) ([]*User, int64, error)
}

// QueryOptions holds pagination and filtering options
type QueryOptions struct {
	Offset   int
	Limit    int
	OrderBy  string
	Keywords string
	Role     UserRole
}
