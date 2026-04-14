// Package identity provides token storage interfaces for refresh tokens
package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TokenStore stores and retrieves refresh token information
type TokenStore interface {
	// Store saves a refresh token with its expiration time
	Store(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error

	// Get retrieves the user ID and expiration for a token
	// Returns (userID, expiresAt, found, error)
	Get(ctx context.Context, token string) (uuid.UUID, time.Time, bool, error)

	// Delete removes a token from the store
	Delete(ctx context.Context, token string) error

	// DeleteAllByUser removes all tokens for a specific user
	DeleteAllByUser(ctx context.Context, userID uuid.UUID) error
}
