// Package identity_test tests UserRole value object
package identity_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/domain/identity"
)

func TestParseUserRole(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    identity.UserRole
		wantErr bool
	}{
		{
			name:  "super admin role",
			input: "super_admin",
			want:  identity.RoleSuperAdmin,
		},
		{
			name:  "admin role",
			input: "admin",
			want:  identity.RoleAdmin,
		},
		{
			name:  "normal role",
			input: "normal",
			want:  identity.RoleNormal,
		},
		{
			name:    "invalid role",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty role",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := identity.ParseUserRole(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUserRole(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseUserRole(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUserRole_IsAdmin(t *testing.T) {
	tests := []struct {
		name string
		role identity.UserRole
		want bool
	}{
		{
			name: "super admin is admin",
			role: identity.RoleSuperAdmin,
			want: true,
		},
		{
			name: "admin is admin",
			role: identity.RoleAdmin,
			want: true,
		},
		{
			name: "normal is not admin",
			role: identity.RoleNormal,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsAdmin(); got != tt.want {
				t.Errorf("UserRole.IsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRole_IsSuperAdmin(t *testing.T) {
	tests := []struct {
		name string
		role identity.UserRole
		want bool
	}{
		{
			name: "super admin is super admin",
			role: identity.RoleSuperAdmin,
			want: true,
		},
		{
			name: "admin is not super admin",
			role: identity.RoleAdmin,
			want: false,
		},
		{
			name: "normal is not super admin",
			role: identity.RoleNormal,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsSuperAdmin(); got != tt.want {
				t.Errorf("UserRole.IsSuperAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}
