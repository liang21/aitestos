// Package identity_test tests User aggregate
package identity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/identity"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		email    string
		password string
		role     identity.UserRole
		wantErr  bool
	}{
		{
			name:     "valid user",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			role:     identity.RoleNormal,
			wantErr:  false,
		},
		{
			name:     "invalid email",
			username: "testuser",
			email:    "",
			password: "password123",
			role:     identity.RoleNormal,
			wantErr:  true,
		},
		{
			name:     "invalid username - too short",
			username: "ab",
			email:    "test@example.com",
			password: "password123",
			role:     identity.RoleNormal,
			wantErr:  true,
		},
		{
			name:     "invalid password - too short",
			username: "testuser",
			email:    "test@example.com",
			password: "short",
			role:     identity.RoleNormal,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := identity.NewUser(tt.username, tt.email, tt.password, tt.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewUser() returned nil user")
					return
				}
				if got.Username() != tt.username {
					t.Errorf("User.Username() = %v, want %v", got.Username(), tt.username)
				}
				if got.Email() != tt.email {
					t.Errorf("User.Email() = %v, want %v", got.Email(), tt.email)
				}
				if got.Role() != tt.role {
					t.Errorf("User.Role() = %v, want %v", got.Role(), tt.role)
				}
				if got.ID() == uuid.Nil {
					t.Error("User.ID() should not be nil UUID")
				}
			}
		})
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	user, err := identity.NewUser("testuser", "test@example.com", "password123", identity.RoleNormal)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "correct password",
			password: "password123",
			want:     true,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := user.VerifyPassword(tt.password); got != tt.want {
				t.Errorf("User.VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_ChangePassword(t *testing.T) {
	user, err := identity.NewUser("testuser", "test@example.com", "password123", identity.RoleNormal)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name        string
		oldPassword string
		newPassword string
		wantErr     bool
	}{
		{
			name:        "valid password change",
			oldPassword: "password123",
			newPassword: "newpassword123",
			wantErr:     false,
		},
		{
			name:        "incorrect old password",
			oldPassword: "wrongpassword",
			newPassword: "newpassword123",
			wantErr:     true,
		},
		{
			name:        "new password too short",
			oldPassword: "password123",
			newPassword: "short",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.ChangePassword(tt.oldPassword, tt.newPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify new password works
				if !user.VerifyPassword(tt.newPassword) {
					t.Error("New password verification failed")
				}
			}
		})
	}
}

func TestUser_Accessors(t *testing.T) {
	user, err := identity.NewUser("testuser", "test@example.com", "password123", identity.RoleAdmin)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test all accessors
	if user.ID() == uuid.Nil {
		t.Error("User.ID() should not be nil")
	}
	if user.Username() != "testuser" {
		t.Errorf("User.Username() = %v, want testuser", user.Username())
	}
	if user.Email() != "test@example.com" {
		t.Errorf("User.Email() = %v, want test@example.com", user.Email())
	}
	if user.Role() != identity.RoleAdmin {
		t.Errorf("User.Role() = %v, want admin", user.Role())
	}
	if user.CreatedAt().IsZero() {
		t.Error("User.CreatedAt() should not be zero")
	}
	if user.UpdatedAt().IsZero() {
		t.Error("User.UpdatedAt() should not be zero")
	}
}
