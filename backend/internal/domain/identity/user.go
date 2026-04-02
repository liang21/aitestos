// Package identity defines user domain model
package identity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User is the aggregate root for identity context
type User struct {
	id        uuid.UUID
	username  string
	email     string
	password  string
	role      UserRole
	createdAt time.Time
	updatedAt time.Time
}

// NewUser creates a new user with validated fields
func NewUser(username, email, rawPassword string, role UserRole) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if err := validatePassword(rawPassword); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(rawPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		id:        uuid.New(),
		username:  username,
		email:     email,
		password:  string(hashedPassword),
		role:      role,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns the user's unique identifier
func (u *User) ID() uuid.UUID { return u.id }

// Username returns the user's username
func (u *User) Username() string { return u.username }

// Email returns the user's email
func (u *User) Email() string { return u.email }

// Role returns the user's role
func (u *User) Role() UserRole { return u.role }

// CreatedAt returns the creation timestamp
func (u *User) CreatedAt() time.Time { return u.createdAt }

// UpdatedAt returns the last update timestamp
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// VerifyPassword checks if the provided password matches
func (u *User) VerifyPassword(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(u.password),
		[]byte(rawPassword),
	)
	return err == nil
}

// ChangePassword updates the user's password after verification
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if !u.VerifyPassword(oldPassword) {
		return ErrPasswordMismatch
	}

	if err := validatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	u.password = string(hashedPassword)
	u.updatedAt = time.Now()
	return nil
}

// UpdateRole changes the user's role
func (u *User) UpdateRole(role UserRole) {
	u.role = role
	u.updatedAt = time.Now()
}

// UserJSON is the JSON representation for API responses
type UserJSON struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToJSON converts User to JSON format
func (u *User) ToJSON() *UserJSON {
	return &UserJSON{
		ID:        u.id,
		Username:  u.username,
		Email:     u.email,
		Role:      string(u.role),
		CreatedAt: u.createdAt,
		UpdatedAt: u.updatedAt,
	}
}

// UserRow is the database row structure
type UserRow struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ToRow converts User to database row
func (u *User) ToRow() *UserRow {
	return &UserRow{
		ID:        u.id,
		Username:  u.username,
		Email:     u.email,
		Password:  u.password,
		Role:      string(u.role),
		CreatedAt: u.createdAt,
		UpdatedAt: u.updatedAt,
	}
}

// FromRow converts database row to User
func FromRow(row *UserRow) (*User, error) {
	role, err := ParseUserRole(row.Role)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        row.ID,
		username:  row.Username,
		email:     row.Email,
		password:  row.Password,
		role:      role,
		createdAt: row.CreatedAt,
		updatedAt: row.UpdatedAt,
	}, nil
}

// Validation functions
func validateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}
	// Basic email validation
	if len(email) < 5 || len(email) > 255 {
		return ErrInvalidEmail
	}
	return nil
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return ErrInvalidUsername
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}
