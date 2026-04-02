// Package identity defines user role value object
package identity

import "errors"

// UserRole is a value object representing user permissions
type UserRole string

const (
	// RoleSuperAdmin has full system access
	RoleSuperAdmin UserRole = "super_admin"
	// RoleAdmin has project-level management access
	RoleAdmin UserRole = "admin"
	// RoleNormal has basic test execution access
	RoleNormal UserRole = "normal"
)

// ParseUserRole converts string to UserRole
func ParseUserRole(s string) (UserRole, error) {
	switch UserRole(s) {
	case RoleSuperAdmin, RoleAdmin, RoleNormal:
		return UserRole(s), nil
	default:
		return "", errors.New("invalid user role")
	}
}

// IsAdmin returns true if user has admin privileges
func (r UserRole) IsAdmin() bool {
	return r == RoleSuperAdmin || r == RoleAdmin
}

// IsSuperAdmin returns true if user is super admin
func (r UserRole) IsSuperAdmin() bool {
	return r == RoleSuperAdmin
}

// CanManageProject returns true if user can manage projects
func (r UserRole) CanManageProject() bool {
	return r.IsAdmin()
}

// CanUploadDocument returns true if user can upload documents
func (r UserRole) CanUploadDocument() bool {
	return true // All roles can upload after permission update
}

// CanGenerateTask returns true if user can create generation tasks
func (r UserRole) CanGenerateTask() bool {
	return true // All roles can generate after permission update
}

// CanConfirmDraft returns true if user can confirm drafts
func (r UserRole) CanConfirmDraft() bool {
	return true // All roles can confirm
}
