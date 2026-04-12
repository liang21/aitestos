package generation_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/identity"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	repoGeneration "github.com/liang21/aitestos/internal/repository/generation"
	repoIdentity "github.com/liang21/aitestos/internal/repository/identity"
	repoProject "github.com/liang21/aitestos/internal/repository/project"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaseDraftRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	userRepo := repoIdentity.NewUserRepository(tc.DB)
	projectRepo := repoProject.NewProjectRepository(tc.DB)
	moduleRepo := repoProject.NewModuleRepository(tc.DB)
	taskRepo := repoGeneration.NewGenerationTaskRepository(tc.DB)
	draftRepo := repoGeneration.NewCaseDraftRepository(tc.DB)
	ctx := context.Background()

	// 辅助函数：创建用户
	createUser := func(t *testing.T) *identity.User {
		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")
		require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")
		return user
	}

	// 辅助函数：创建项目、模块和任务
	createProjectModuleTask := func(t *testing.T) (*domainproject.Project, *domainproject.Module, *generation.GenerationTask) {
		user := createUser(t)

		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")

		module, err := testsetup.NewModuleBuilder(project.ID()).Build()
		require.NoError(t, err, "build module should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module), "save module should succeed")

		task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, taskRepo.Save(ctx, task), "save task should succeed")

		return project, module, task
	}

	t.Run("Save", func(t *testing.T) {
		tc.CleanupTest()

		_, _, task := createProjectModuleTask(t)

		tests := []struct {
			name    string
			builder *testsetup.CaseDraftBuilder
			wantErr error
		}{
			{
				name:    "save pending draft",
				builder: testsetup.NewCaseDraftBuilder(task.ID()).WithStatus(generation.DraftPending),
				wantErr: nil,
			},
			{
				name:    "save confirmed draft",
				builder: testsetup.NewCaseDraftBuilder(task.ID()).WithStatus(generation.DraftConfirmed),
				wantErr: nil,
			},
			{
				name:    "save rejected draft",
				builder: testsetup.NewCaseDraftBuilder(task.ID()).WithStatus(generation.DraftRejected),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				draft, err := tt.builder.Build()
				require.NoError(t, err, "build draft should succeed")

				err = draftRepo.Save(ctx, draft)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					found, err := draftRepo.FindByID(ctx, draft.ID())
					require.NoError(t, err, "find draft by ID should succeed")
					testsetup.AssertCaseDraftEqual(t, draft, found)
				}
			})
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		tc.CleanupTest()

		_, _, task := createProjectModuleTask(t)

		draft, err := testsetup.NewCaseDraftBuilder(task.ID()).Build()
		require.NoError(t, err, "build draft should succeed")
		require.NoError(t, draftRepo.Save(ctx, draft), "save draft should succeed")

		found, err := draftRepo.FindByID(ctx, draft.ID())
		require.NoError(t, err, "find draft by ID should succeed")
		testsetup.AssertCaseDraftEqual(t, draft, found)

		// 测试不存在的 ID
		_, err = draftRepo.FindByID(ctx, uuid.New())
		require.Error(t, err, "find non-existent draft should fail")
		assert.ErrorIs(t, err, generation.ErrDraftNotFound, "error should be ErrDraftNotFound")
	})

	t.Run("FindByTaskID", func(t *testing.T) {
		tc.CleanupTest()

		_, _, task := createProjectModuleTask(t)

		// 创建多个草稿
		for i := 0; i < 3; i++ {
			draft, err := testsetup.NewCaseDraftBuilder(task.ID()).Build()
			require.NoError(t, err, "build draft should succeed")
			require.NoError(t, draftRepo.Save(ctx, draft), "save draft should succeed")
		}

		drafts, err := draftRepo.FindByTaskID(ctx, task.ID())
		require.NoError(t, err, "find drafts by task ID should succeed")
		assert.Equal(t, 3, len(drafts), "should return 3 drafts")
	})

	t.Run("FindByTaskIDAndStatus", func(t *testing.T) {
		tc.CleanupTest()

		_, _, task := createProjectModuleTask(t)

		// 创建不同状态的草稿
		for i := 0; i < 2; i++ {
			draft, err := testsetup.NewCaseDraftBuilder(task.ID()).
				WithStatus(generation.DraftPending).Build()
			require.NoError(t, err, "build pending draft should succeed")
			require.NoError(t, draftRepo.Save(ctx, draft), "save pending draft should succeed")
		}

		for i := 0; i < 3; i++ {
			draft, err := testsetup.NewCaseDraftBuilder(task.ID()).
				WithStatus(generation.DraftConfirmed).Build()
			require.NoError(t, err, "build confirmed draft should succeed")
			require.NoError(t, draftRepo.Save(ctx, draft), "save confirmed draft should succeed")
		}

		// 查询 pending 状态的草稿
		pendingDrafts, err := draftRepo.FindByTaskIDAndStatus(ctx, task.ID(), generation.DraftPending)
		require.NoError(t, err, "find pending drafts should succeed")
		assert.Equal(t, 2, len(pendingDrafts), "should return 2 pending drafts")

		// 查询 confirmed 状态的草稿
		confirmedDrafts, err := draftRepo.FindByTaskIDAndStatus(ctx, task.ID(), generation.DraftConfirmed)
		require.NoError(t, err, "find confirmed drafts should succeed")
		assert.Equal(t, 3, len(confirmedDrafts), "should return 3 confirmed drafts")
	})

	t.Run("Update", func(t *testing.T) {
		tc.CleanupTest()

		_, module, task := createProjectModuleTask(t)

		draft, err := testsetup.NewCaseDraftBuilder(task.ID()).Build()
		require.NoError(t, err, "build draft should succeed")
		require.NoError(t, draftRepo.Save(ctx, draft), "save draft should succeed")

		// 更新：确认草稿
		require.NoError(t, draft.Confirm(module.ID()), "confirm draft should succeed")
		err = draftRepo.Update(ctx, draft)
		require.NoError(t, err, "update draft should succeed")

		found, err := draftRepo.FindByID(ctx, draft.ID())
		require.NoError(t, err, "find draft should succeed")
		assert.Equal(t, generation.DraftConfirmed, found.Status(), "status should be confirmed")
		assert.Equal(t, module.ID(), *found.ModuleID(), "module ID should be set")
	})

	t.Run("Update with rejection", func(t *testing.T) {
		tc.CleanupTest()

		_, _, task := createProjectModuleTask(t)

		draft, err := testsetup.NewCaseDraftBuilder(task.ID()).Build()
		require.NoError(t, err, "build draft should succeed")
		require.NoError(t, draftRepo.Save(ctx, draft), "save draft should succeed")

		// 更新：拒绝草稿
		require.NoError(t, draft.Reject(generation.ReasonLowQuality, "Content is too simple"), "reject draft should succeed")
		err = draftRepo.Update(ctx, draft)
		require.NoError(t, err, "update draft should succeed")

		found, err := draftRepo.FindByID(ctx, draft.ID())
		require.NoError(t, err, "find draft should succeed")
		assert.Equal(t, generation.DraftRejected, found.Status(), "status should be rejected")
		assert.NotEmpty(t, found.Feedback(), "feedback should be set")
	})

	t.Run("BatchUpdateStatus", func(t *testing.T) {
		tc.CleanupTest()

		_, module, task := createProjectModuleTask(t)

		// 创建多个草稿
		draftIDs := make([]uuid.UUID, 3)
		for i := 0; i < 3; i++ {
			draft, err := testsetup.NewCaseDraftBuilder(task.ID()).Build()
			require.NoError(t, err, "build draft should succeed")
			require.NoError(t, draftRepo.Save(ctx, draft), "save draft should succeed")
			draftIDs[i] = draft.ID()
		}

		// 批量确认
		err := draftRepo.BatchUpdateStatus(ctx, draftIDs, generation.DraftConfirmed, module.ID())
		require.NoError(t, err, "batch update status should succeed")

		// 验证所有草稿都已确认
		for _, id := range draftIDs {
			found, err := draftRepo.FindByID(ctx, id)
			require.NoError(t, err, "find draft should succeed")
			assert.Equal(t, generation.DraftConfirmed, found.Status(), "status should be confirmed")
			assert.Equal(t, module.ID(), *found.ModuleID(), "module ID should be set")
		}
	})

	t.Run("CountByTaskIDAndStatus", func(t *testing.T) {
		tc.CleanupTest()

		_, _, task := createProjectModuleTask(t)

		// 创建不同状态的草稿
		for i := 0; i < 2; i++ {
			draft, err := testsetup.NewCaseDraftBuilder(task.ID()).
				WithStatus(generation.DraftPending).Build()
			require.NoError(t, err, "build pending draft should succeed")
			require.NoError(t, draftRepo.Save(ctx, draft), "save pending draft should succeed")
		}

		for i := 0; i < 3; i++ {
			draft, err := testsetup.NewCaseDraftBuilder(task.ID()).
				WithStatus(generation.DraftConfirmed).Build()
			require.NoError(t, err, "build confirmed draft should succeed")
			require.NoError(t, draftRepo.Save(ctx, draft), "save confirmed draft should succeed")
		}

		// 统计 pending 草稿
		pendingCount, err := draftRepo.CountByTaskIDAndStatus(ctx, task.ID(), generation.DraftPending)
		require.NoError(t, err, "count pending drafts should succeed")
		assert.Equal(t, int64(2), pendingCount, "should count 2 pending drafts")

		// 统计 confirmed 草稿
		confirmedCount, err := draftRepo.CountByTaskIDAndStatus(ctx, task.ID(), generation.DraftConfirmed)
		require.NoError(t, err, "count confirmed drafts should succeed")
		assert.Equal(t, int64(3), confirmedCount, "should count 3 confirmed drafts")
	})

	t.Run("Delete", func(t *testing.T) {
		tc.CleanupTest()

		_, _, task := createProjectModuleTask(t)

		draft, err := testsetup.NewCaseDraftBuilder(task.ID()).Build()
		require.NoError(t, err, "build draft should succeed")
		require.NoError(t, draftRepo.Save(ctx, draft), "save draft should succeed")

		err = draftRepo.Delete(ctx, draft.ID())
		require.NoError(t, err, "delete draft should succeed")

		// 验证删除
		_, err = draftRepo.FindByID(ctx, draft.ID())
		require.Error(t, err, "find deleted draft should fail")
		assert.ErrorIs(t, err, generation.ErrDraftNotFound, "error should be ErrDraftNotFound")
	})

	t.Run("DeleteByTaskID", func(t *testing.T) {
		tc.CleanupTest()

		_, _, task := createProjectModuleTask(t)

		// 创建多个草稿
		for i := 0; i < 3; i++ {
			draft, err := testsetup.NewCaseDraftBuilder(task.ID()).Build()
			require.NoError(t, err, "build draft should succeed")
			require.NoError(t, draftRepo.Save(ctx, draft), "save draft should succeed")
		}

		// 删除任务的所有草稿
		err := draftRepo.DeleteByTaskID(ctx, task.ID())
		require.NoError(t, err, "delete drafts by task ID should succeed")

		// 验证删除
		drafts, err := draftRepo.FindByTaskID(ctx, task.ID())
		require.NoError(t, err, "find drafts should succeed")
		assert.Equal(t, 0, len(drafts), "should have no drafts after deletion")
	})

	t.Run("Cascade delete on task deletion", func(t *testing.T) {
		tc.CleanupTest()

		_, _, task := createProjectModuleTask(t)

		draft, err := testsetup.NewCaseDraftBuilder(task.ID()).Build()
		require.NoError(t, err, "build draft should succeed")
		require.NoError(t, draftRepo.Save(ctx, draft), "save draft should succeed")

		// 删除任务
		err = taskRepo.Delete(ctx, task.ID())
		require.NoError(t, err, "delete task should succeed")

		// 验证草稿也被删除（级联删除）
		_, err = draftRepo.FindByID(ctx, draft.ID())
		require.Error(t, err, "find draft after task deletion should fail")
	})
}
