package identity_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/identity"
	identityrepo "github.com/liang21/aitestos/internal/repository/identity"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	userRepo := identityrepo.NewUserRepository(tc.DB)
	ctx := context.Background()

	t.Run("Save", func(t *testing.T) {
		tc.CleanupTest()

		tests := []struct {
			name    string
			builder *testsetup.UserBuilder
			wantErr error
		}{
			{
				name:    "save normal user",
				builder: testsetup.NewUserBuilder().WithRole(identity.RoleNormal),
				wantErr: nil,
			},
			{
				name:    "save admin user",
				builder: testsetup.NewUserBuilder().WithRole(identity.RoleAdmin),
				wantErr: nil,
			},
			{
				name:    "save super admin user",
				builder: testsetup.NewUserBuilder().WithRole(identity.RoleSuperAdmin),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				user, err := tt.builder.Build()
				require.NoError(t, err, "build user should succeed")

				err = userRepo.Save(ctx, user)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					// 验证可以找到保存的用户
					found, err := userRepo.FindByID(ctx, user.ID())
					require.NoError(t, err, "find user by ID should succeed")
					testsetup.AssertUserEqual(t, user, found)
				}
			})
		}
	})

	t.Run("Save duplicate username", func(t *testing.T) {
		tc.CleanupTest()

		username := "duplicate_user"

		user1, err := testsetup.NewUserBuilder().WithUsername(username).Build()
		require.NoError(t, err, "build user1 should succeed")
		require.NoError(t, userRepo.Save(ctx, user1), "save user1 should succeed")

		user2, err := testsetup.NewUserBuilder().WithUsername(username).Build()
		require.NoError(t, err, "build user2 should succeed")
		err = userRepo.Save(ctx, user2)
		require.Error(t, err, "save duplicate username should fail")
	})

	t.Run("Save duplicate email", func(t *testing.T) {
		tc.CleanupTest()

		email := "duplicate@example.com"

		user1, err := testsetup.NewUserBuilder().WithEmail(email).Build()
		require.NoError(t, err, "build user1 should succeed")
		require.NoError(t, userRepo.Save(ctx, user1), "save user1 should succeed")

		user2, err := testsetup.NewUserBuilder().WithEmail(email).Build()
		require.NoError(t, err, "build user2 should succeed")
		err = userRepo.Save(ctx, user2)
		require.Error(t, err, "save duplicate email should fail")
	})

	t.Run("FindByID", func(t *testing.T) {
		tc.CleanupTest()

		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")
		require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")

		found, err := userRepo.FindByID(ctx, user.ID())
		require.NoError(t, err, "find user by ID should succeed")
		testsetup.AssertUserEqual(t, user, found)

		// 测试不存在的 ID
		_, err = userRepo.FindByID(ctx, randomUUID())
		require.Error(t, err, "find non-existent user should fail")
		assert.ErrorIs(t, err, identity.ErrUserNotFound, "error should be ErrUserNotFound")
	})

	t.Run("FindByEmail", func(t *testing.T) {
		tc.CleanupTest()

		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")
		require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")

		found, err := userRepo.FindByEmail(ctx, user.Email())
		require.NoError(t, err, "find user by email should succeed")
		testsetup.AssertUserEqual(t, user, found)

		// 测试不存在的邮箱
		_, err = userRepo.FindByEmail(ctx, "notfound@example.com")
		require.Error(t, err, "find non-existent email should fail")
		assert.ErrorIs(t, err, identity.ErrUserNotFound, "error should be ErrUserNotFound")
	})

	t.Run("FindByUsername", func(t *testing.T) {
		tc.CleanupTest()

		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")
		require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")

		found, err := userRepo.FindByUsername(ctx, user.Username())
		require.NoError(t, err, "find user by username should succeed")
		testsetup.AssertUserEqual(t, user, found)

		// 测试不存在的用户名
		_, err = userRepo.FindByUsername(ctx, "notfound")
		require.Error(t, err, "find non-existent username should fail")
		assert.ErrorIs(t, err, identity.ErrUserNotFound, "error should be ErrUserNotFound")
	})

	t.Run("Update", func(t *testing.T) {
		tc.CleanupTest()

		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")
		require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")

		// 更新角色
		user.UpdateRole(identity.RoleAdmin)
		err = userRepo.Update(ctx, user)
		require.NoError(t, err, "update user should succeed")

		found, err := userRepo.FindByID(ctx, user.ID())
		require.NoError(t, err, "find user should succeed")
		assert.Equal(t, identity.RoleAdmin, found.Role(), "role should be updated to admin")
	})

	t.Run("Delete", func(t *testing.T) {
		tc.CleanupTest()

		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")
		require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")

		err = userRepo.Delete(ctx, user.ID())
		require.NoError(t, err, "delete user should succeed")

		// 验证软删除
		_, err = userRepo.FindByID(ctx, user.ID())
		require.Error(t, err, "find deleted user should fail")
		assert.ErrorIs(t, err, identity.ErrUserNotFound, "error should be ErrUserNotFound")
	})

	t.Run("List", func(t *testing.T) {
		tc.CleanupTest()

		// 创建多个用户
		for i := 0; i < 5; i++ {
			user, err := testsetup.NewUserBuilder().Build()
			require.NoError(t, err, "build user should succeed")
			require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")
		}

		// 测试分页
		users, total, err := userRepo.List(ctx, identity.QueryOptions{
			Limit:  3,
			Offset: 0,
		})
		require.NoError(t, err, "list users should succeed")
		assert.Equal(t, 3, len(users), "should return 3 users")
		assert.GreaterOrEqual(t, total, int64(5), "total count should be at least 5")

		// 测试第二页
		users, _, err = userRepo.List(ctx, identity.QueryOptions{
			Limit:  3,
			Offset: 3,
		})
		require.NoError(t, err, "list users second page should succeed")
		assert.GreaterOrEqual(t, len(users), 2, "should return at least 2 users")
	})

	t.Run("List with filter", func(t *testing.T) {
		tc.CleanupTest()

		// 创建不同角色的用户
		adminUser, err := testsetup.NewUserBuilder().WithRole(identity.RoleAdmin).Build()
		require.NoError(t, err, "build admin user should succeed")
		require.NoError(t, userRepo.Save(ctx, adminUser), "save admin user should succeed")

		normalUser, err := testsetup.NewUserBuilder().WithRole(identity.RoleNormal).Build()
		require.NoError(t, err, "build normal user should succeed")
		require.NoError(t, userRepo.Save(ctx, normalUser), "save normal user should succeed")

		// 测试按角色过滤
		users, _, err := userRepo.List(ctx, identity.QueryOptions{
			Role:   identity.RoleAdmin,
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "list users with filter should succeed")
		assert.GreaterOrEqual(t, len(users), 1, "should return at least 1 admin user")
		for _, u := range users {
			assert.Equal(t, identity.RoleAdmin, u.Role(), "all users should be admins")
		}
	})
}

func randomUUID() uuid.UUID {
	return uuid.New()
}
