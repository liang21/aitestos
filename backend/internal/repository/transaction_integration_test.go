package repository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/liang21/aitestos/internal/domain/identity"
	identityrepo "github.com/liang21/aitestos/internal/repository/identity"
	"github.com/liang21/aitestos/internal/repository"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxManager_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	t.Run("commit transaction", func(t *testing.T) {
		tc.CleanupTest()
		ctx := context.Background()

		userRepo := identityrepo.NewUserRepository(tc.DB)
		txManager := repository.NewTxManager(func(ctx context.Context) (repository.Tx, error) {
			return tc.DB.BeginTxx(ctx, nil)
		})

		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")

		err = txManager.WithTransaction(ctx, func(txCtx context.Context) error {
			return userRepo.Save(txCtx, user)
		})
		require.NoError(t, err, "transaction should commit successfully")

		// 验证数据已提交
		found, err := userRepo.FindByEmail(ctx, user.Email())
		require.NoError(t, err, "find user by email should succeed")
		testsetup.AssertUserEqual(t, user, found)
	})

	t.Run("rollback on error", func(t *testing.T) {
		tc.CleanupTest()
		ctx := context.Background()

		userRepo := identityrepo.NewUserRepository(tc.DB)
		txManager := repository.NewTxManager(func(ctx context.Context) (repository.Tx, error) {
			return tc.DB.BeginTxx(ctx, nil)
		})

		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")

		// 在事务中保存用户，然后返回错误
		err = txManager.WithTransaction(ctx, func(txCtx context.Context) error {
			if err := userRepo.Save(txCtx, user); err != nil {
				return err
			}
			// 模拟业务错误
			return identity.ErrUsernameDuplicate
		})
		require.Error(t, err, "transaction should return error")
		assert.ErrorIs(t, err, identity.ErrUsernameDuplicate, "error should be ErrUsernameDuplicate")

		// 验证数据未提交
		_, err = userRepo.FindByEmail(ctx, user.Email())
		require.Error(t, err, "find user should fail")
		assert.ErrorIs(t, err, identity.ErrUserNotFound, "error should be ErrUserNotFound")
	})

	t.Run("rollback on panic", func(t *testing.T) {
		tc.CleanupTest()
		ctx := context.Background()

		userRepo := identityrepo.NewUserRepository(tc.DB)
		txManager := repository.NewTxManager(func(ctx context.Context) (repository.Tx, error) {
			return tc.DB.BeginTxx(ctx, nil)
		})

		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")

		// 在事务中保存用户，然后 panic
		// recovered tracks panic recovery
		func() {
			defer func() {
				_ = recover()
			}()

			_ = txManager.WithTransaction(ctx, func(txCtx context.Context) error {
				if err := userRepo.Save(txCtx, user); err != nil {
					return err
				}
				panic("test panic")
			})
		}()

		// 注意：testcontainers 的 TxManager 应该捕获 panic 并回滚
		// 但不会重新 panic，所以这里可能需要调整
		// 验证数据未提交
		_, err = userRepo.FindByEmail(ctx, user.Email())
		require.Error(t, err, "find user should fail after panic")
	})

	t.Run("nested transaction uses same connection", func(t *testing.T) {
		tc.CleanupTest()
		ctx := context.Background()

		userRepo := identityrepo.NewUserRepository(tc.DB)
		txManager := repository.NewTxManager(func(ctx context.Context) (repository.Tx, error) {
			return tc.DB.BeginTxx(ctx, nil)
		})

		user1, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user1 should succeed")

		user2, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user2 should succeed")

		// 嵌套事务
		err = txManager.WithTransaction(ctx, func(ctx1 context.Context) error {
			// 保存第一个用户
			if err := userRepo.Save(ctx1, user1); err != nil {
				return err
			}

			// 嵌套事务保存第二个用户
			return txManager.WithTransaction(ctx1, func(ctx2 context.Context) error {
				return userRepo.Save(ctx2, user2)
			})
		})
		require.NoError(t, err, "nested transaction should succeed")

		// 验证两个用户都已提交
		found1, err := userRepo.FindByEmail(ctx, user1.Email())
		require.NoError(t, err, "find user1 should succeed")
		testsetup.AssertUserEqual(t, user1, found1)

		found2, err := userRepo.FindByEmail(ctx, user2.Email())
		require.NoError(t, err, "find user2 should succeed")
		testsetup.AssertUserEqual(t, user2, found2)
	})

	t.Run("nested transaction rollback", func(t *testing.T) {
		tc.CleanupTest()
		ctx := context.Background()

		userRepo := identityrepo.NewUserRepository(tc.DB)
		txManager := repository.NewTxManager(func(ctx context.Context) (repository.Tx, error) {
			return tc.DB.BeginTxx(ctx, nil)
		})

		user1, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user1 should succeed")

		user2, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user2 should succeed")

		// 嵌套事务，内部事务返回错误
		err = txManager.WithTransaction(ctx, func(ctx1 context.Context) error {
			// 保存第一个用户
			if err := userRepo.Save(ctx1, user1); err != nil {
				return err
			}

			// 嵌套事务保存第二个用户，然后返回错误
			return txManager.WithTransaction(ctx1, func(ctx2 context.Context) error {
				if err := userRepo.Save(ctx2, user2); err != nil {
					return err
				}
				return errors.New("nested transaction error")
			})
		})
		require.Error(t, err, "nested transaction should return error")

		// 验证两个用户都未提交
		_, err = userRepo.FindByEmail(ctx, user1.Email())
		require.Error(t, err, "find user1 should fail")

		_, err = userRepo.FindByEmail(ctx, user2.Email())
		require.Error(t, err, "find user2 should fail")
	})

	t.Run("concurrent transactions", func(t *testing.T) {
		tc.CleanupTest()
		ctx := context.Background()

		userRepo := identityrepo.NewUserRepository(tc.DB)
		txManager := repository.NewTxManager(func(ctx context.Context) (repository.Tx, error) {
			return tc.DB.BeginTxx(ctx, nil)
		})

		// 并发执行多个事务
		const goroutines = 5
		errCh := make(chan error, goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				user, err := testsetup.NewUserBuilder().Build()
				if err != nil {
					errCh <- fmt.Errorf("build user: %w", err)
					return
				}

				err = txManager.WithTransaction(ctx, func(txCtx context.Context) error {
					return userRepo.Save(txCtx, user)
				})
				errCh <- err
			}()
		}

		// 等待所有事务完成
		for i := 0; i < goroutines; i++ {
			err := <-errCh
			assert.NoError(t, err, "concurrent transaction should succeed")
		}
	})
}
