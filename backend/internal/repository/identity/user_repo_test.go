// Package identity_test tests UserRepository implementation
package identity_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domainidentity "github.com/liang21/aitestos/internal/domain/identity"
)

// Test fixtures
func createTestUser(t *testing.T, username string) *domainidentity.User {
	t.Helper()
	user, err := domainidentity.NewUser(username, username+"@test.com", "password123", domainidentity.RoleNormal)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func TestUserRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires a database connection
	// In CI/CD, this would use testcontainers
	t.Run("save new user", func(t *testing.T) {
		// Placeholder for integration test
		// Real implementation would use testcontainers
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find existing user", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by email", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestUserRepository_FindByUsername(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("find by username", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestUserRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("update user", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestUserRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("soft delete user", func(t *testing.T) {
		// Placeholder for integration test
	})
}

func TestUserRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("list users with pagination", func(t *testing.T) {
		// Placeholder for integration test
	})
}

// Unit tests for repository helpers
func TestUserRow_Conversion(t *testing.T) {
	user := createTestUser(t, "testuser")

	// Test ToRow conversion
	row := user.ToRow()
	if row.ID != user.ID() {
		t.Errorf("ToRow().ID = %v, want %v", row.ID, user.ID())
	}
	if row.Username != user.Username() {
		t.Errorf("ToRow().Username = %v, want %v", row.Username, user.Username())
	}
	if row.Email != user.Email() {
		t.Errorf("ToRow().Email = %v, want %v", row.Email, user.Email())
	}
	if row.Role != string(user.Role()) {
		t.Errorf("ToRow().Role = %v, want %v", row.Role, user.Role())
	}

	// Test FromRow conversion
	reconstructed, err := domainidentity.FromRow(row)
	if err != nil {
		t.Fatalf("FromRow() error = %v", err)
	}
	if reconstructed.ID() != user.ID() {
		t.Errorf("FromRow().ID() = %v, want %v", reconstructed.ID(), user.ID())
	}
	if reconstructed.Username() != user.Username() {
		t.Errorf("FromRow().Username() = %v, want %v", reconstructed.Username(), user.Username())
	}
}

// MockUserRepository for testing without database
type MockUserRepository struct {
	users           map[uuid.UUID]*domainidentity.User
	usersByEmail    map[string]*domainidentity.User
	usersByUsername map[string]*domainidentity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:           make(map[uuid.UUID]*domainidentity.User),
		usersByEmail:    make(map[string]*domainidentity.User),
		usersByUsername: make(map[string]*domainidentity.User),
	}
}

func (m *MockUserRepository) Save(ctx context.Context, user *domainidentity.User) error {
	m.users[user.ID()] = user
	m.usersByEmail[user.Email()] = user
	m.usersByUsername[user.Username()] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainidentity.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, domainidentity.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domainidentity.User, error) {
	user, ok := m.usersByEmail[email]
	if !ok {
		return nil, domainidentity.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*domainidentity.User, error) {
	user, ok := m.usersByUsername[username]
	if !ok {
		return nil, domainidentity.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domainidentity.User) error {
	if _, ok := m.users[user.ID()]; !ok {
		return domainidentity.ErrUserNotFound
	}
	m.users[user.ID()] = user
	m.usersByEmail[user.Email()] = user
	m.usersByUsername[user.Username()] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	user, ok := m.users[id]
	if !ok {
		return domainidentity.ErrUserNotFound
	}
	delete(m.users, id)
	delete(m.usersByEmail, user.Email())
	delete(m.usersByUsername, user.Username())
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, opts domainidentity.QueryOptions) ([]*domainidentity.User, int64, error) {
	users := make([]*domainidentity.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, int64(len(users)), nil
}

func TestMockUserRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMockUserRepository()

	// Create
	user := createTestUser(t, "testuser")
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read by ID
	found, err := repo.FindByID(ctx, user.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.ID() != user.ID() {
		t.Errorf("FindByID().ID() = %v, want %v", found.ID(), user.ID())
	}

	// Read by Email
	found, err = repo.FindByEmail(ctx, user.Email())
	if err != nil {
		t.Fatalf("FindByEmail() error = %v", err)
	}
	if found.Email() != user.Email() {
		t.Errorf("FindByEmail().Email() = %v, want %v", found.Email(), user.Email())
	}

	// Read by Username
	found, err = repo.FindByUsername(ctx, user.Username())
	if err != nil {
		t.Fatalf("FindByUsername() error = %v", err)
	}
	if found.Username() != user.Username() {
		t.Errorf("FindByUsername().Username() = %v, want %v", found.Username(), user.Username())
	}

	// Update
	user.UpdateRole(domainidentity.RoleAdmin)
	err = repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	found, _ = repo.FindByID(ctx, user.ID())
	if found.Role() != domainidentity.RoleAdmin {
		t.Errorf("Update().Role() = %v, want %v", found.Role(), domainidentity.RoleAdmin)
	}

	// Delete
	err = repo.Delete(ctx, user.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err = repo.FindByID(ctx, user.ID())
	if err != domainidentity.ErrUserNotFound {
		t.Errorf("FindByID() after Delete() error = %v, want %v", err, domainidentity.ErrUserNotFound)
	}
}

func TestMockUserRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMockUserRepository()

	_, err := repo.FindByID(ctx, uuid.New())
	if err != domainidentity.ErrUserNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domainidentity.ErrUserNotFound)
	}

	_, err = repo.FindByEmail(ctx, "notfound@test.com")
	if err != domainidentity.ErrUserNotFound {
		t.Errorf("FindByEmail() error = %v, want %v", err, domainidentity.ErrUserNotFound)
	}

	_, err = repo.FindByUsername(ctx, "notfound")
	if err != domainidentity.ErrUserNotFound {
		t.Errorf("FindByUsername() error = %v, want %v", err, domainidentity.ErrUserNotFound)
	}

	err = repo.Update(ctx, createTestUser(t, "testuser"))
	if err != domainidentity.ErrUserNotFound {
		t.Errorf("Update() error = %v, want %v", err, domainidentity.ErrUserNotFound)
	}

	err = repo.Delete(ctx, uuid.New())
	if err != domainidentity.ErrUserNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domainidentity.ErrUserNotFound)
	}
}
