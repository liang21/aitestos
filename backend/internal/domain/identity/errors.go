// Package identity defines domain errors
package identity

import "errors"

var (
	// ErrUserNotFound indicates user does not exist
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailDuplicate indicates email is already registered
	ErrEmailDuplicate = errors.New("email already exists")
	// ErrUsernameDuplicate indicates username is already taken
	ErrUsernameDuplicate = errors.New("username already exists")
	// ErrPasswordMismatch indicates password verification failed
	ErrPasswordMismatch = errors.New("password mismatch")
	// ErrInvalidEmail indicates email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")
	// ErrInvalidUsername indicates username format is invalid
	ErrInvalidUsername = errors.New("invalid username: must be 3-32 characters")
	// ErrPasswordTooShort indicates password is too short
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	// ErrPermissionDenied indicates user lacks required permissions
	ErrPermissionDenied = errors.New("permission denied")

	// Token/Auth errors
	// ErrTokenMissing indicates authorization header is missing
	ErrTokenMissing = errors.New("missing authorization header")
	// ErrTokenInvalid indicates token is invalid or malformed
	ErrTokenInvalid = errors.New("invalid token")
	// ErrTokenExpired indicates token has expired
	ErrTokenExpired = errors.New("token expired")
	// ErrTokenRevoked indicates refresh token has been revoked
	ErrTokenRevoked = errors.New("refresh token has been revoked")
	// ErrInvalidRole indicates user role is invalid
	ErrInvalidRole = errors.New("invalid user role")
)
