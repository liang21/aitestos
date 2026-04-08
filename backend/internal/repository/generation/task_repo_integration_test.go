package generation_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/identity"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/repository/generation"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerationTaskRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	userRepo := identityrepo.NewUserRepository(tc.DB)
	projectRepo := repository.NewProjectRepository(tc.DB)
	taskRepo := repository.NewGenerationTaskRepository(tc.DB)
	ctx := context.Background()

	// 辅助函数：创建用户
	createUser := func(t *testing.T) *identity.User {
		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")
		require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")
		return user
	}

	// 辅助函数：创建项目
	createProject := func(t *testing.T) *domainproject.Project {
		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")
		return project
	}

	t.Run("Save", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		tests := []struct {
			name    string
			builder *testsetup.GenerationTaskBuilder
			wantErr error
		}{
			{
				name:    "save pending task",
				builder: testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).WithStatus(generation.TaskPending),
				wantErr: nil,
			},
			{
				name:    "save processing task",
				builder: testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).WithStatus(generation.TaskProcessing),
				wantErr: nil,
			},
			{
				name:    "save completed task",
				builder: testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).WithStatus(generation.TaskCompleted),
				wantErr: nil,
			},
			{
				name:    "save failed task",
				builder: testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).WithStatus(generation.TaskFailed),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				task := tt.builder.Build()

				err := taskRepo.Save(ctx, task)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					found, err := taskRepo.FindByID(ctx, task.ID())
					require.NoError(t, err, "find task by ID should succeed")
					testsetup.AssertGenerationTaskEqual(t, task, found)
				}
			})
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, taskRepo.Save(ctx, task), "save task should succeed")

		found, err := taskRepo.FindByID(ctx, task.ID())
		require.NoError(t, err, "find task by ID should succeed")
		testsetup.AssertGenerationTaskEqual(t, task, found)

		// 测试不存在的 ID
		_, err = taskRepo.FindByID(ctx, uuid.New())
		require.Error(t, err, "find non-existent task should fail")
		assert.ErrorIs(t, err, generation.ErrTaskNotFound, "error should be ErrTaskNotFound")
	})

	t.Run("FindByProjectID", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		// 创建多个任务
		for i := 0; i < 3; i++ {
			task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).Build()
			require.NoError(t, taskRepo.Save(ctx, task), "save task should succeed")
		}

		tasks, err := taskRepo.FindByProjectID(ctx, project.ID(), generation.QueryOptions{
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "find tasks by project ID should succeed")
		assert.Equal(t, 3, len(tasks), "should return 3 tasks")
	})

	t.Run("FindByStatus", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		// 创建不同状态的任务
		for i := 0; i < 2; i++ {
			task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).
				WithStatus(generation.TaskPending).Build()
			require.NoError(t, taskRepo.Save(ctx, task), "save pending task should succeed")
		}

		for i := 0; i < 3; i++ {
			task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).
				WithStatus(generation.TaskCompleted).Build()
			require.NoError(t, taskRepo.Save(ctx, task), "save completed task should succeed")
		}

		// 查询 pending 状态的任务
		pendingTasks, err := taskRepo.FindByStatus(ctx, generation.TaskPending, generation.QueryOptions{
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "find pending tasks should succeed")
		assert.GreaterOrEqual(t, len(pendingTasks), 2, "should return at least 2 pending tasks")

		// 查询 completed 状态的任务
		completedTasks, err := taskRepo.FindByStatus(ctx, generation.TaskCompleted, generation.QueryOptions{
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "find completed tasks should succeed")
		assert.GreaterOrEqual(t, len(completedTasks), 3, "should return at least 3 completed tasks")
	})

	t.Run("Update", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, taskRepo.Save(ctx, task), "save task should succeed")

		// 更新状态：pending -> processing
		task.StartProcessing()
		err := taskRepo.Update(ctx, task)
		require.NoError(t, err, "update task should succeed")

		found, err := taskRepo.FindByID(ctx, task.ID())
		require.NoError(t, err, "find task should succeed")
		assert.Equal(t, generation.TaskProcessing, found.Status(), "status should be processing")

		// 更新状态：processing -> completed
		summary := map[string]any{"total": 10, "generated": 10}
		task.Complete(summary)
		err = taskRepo.Update(ctx, task)
		require.NoError(t, err, "update task should succeed")

		found, err = taskRepo.FindByID(ctx, task.ID())
		require.NoError(t, err, "find task should succeed")
		assert.Equal(t, generation.TaskCompleted, found.Status(), "status should be completed")
		assert.Equal(t, summary, found.ResultSummary(), "result summary should be set")
	})

	t.Run("Update with failure", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, taskRepo.Save(ctx, task), "save task should succeed")

		// 更新状态：pending -> processing -> failed
		task.StartProcessing()
		task.Fail("API timeout error")
		err := taskRepo.Update(ctx, task)
		require.NoError(t, err, "update task should succeed")

		found, err := taskRepo.FindByID(ctx, task.ID())
		require.NoError(t, err, "find task should succeed")
		assert.Equal(t, generation.TaskFailed, found.Status(), "status should be failed")
		assert.Equal(t, "API timeout error", found.ErrorMsg(), "error message should be set")
	})

	t.Run("FindByUserID", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		// 创建多个任务
		for i := 0; i < 3; i++ {
			task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).Build()
			require.NoError(t, taskRepo.Save(ctx, task), "save task should succeed")
		}

		tasks, err := taskRepo.FindByUserID(ctx, user.ID(), generation.QueryOptions{
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "find tasks by user ID should succeed")
		assert.GreaterOrEqual(t, len(tasks), 3, "should return at least 3 tasks")
	})

	t.Run("Delete", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, taskRepo.Save(ctx, task), "save task should succeed")

		err := taskRepo.Delete(ctx, task.ID())
		require.NoError(t, err, "delete task should succeed")

		// 验证删除
		_, err = taskRepo.FindByID(ctx, task.ID())
		require.Error(t, err, "find deleted task should fail")
		assert.ErrorIs(t, err, generation.ErrTaskNotFound, "error should be ErrTaskNotFound")
	})

	t.Run("Cascade delete on project deletion", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		task := testsetup.NewGenerationTaskBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, taskRepo.Save(ctx, task), "save task should succeed")

		// 删除项目
		err := projectRepo.Delete(ctx, project.ID())
		require.NoError(t, err, "delete project should succeed")

		// 验证任务也被删除（级联删除）
		_, err = taskRepo.FindByID(ctx, task.ID())
		require.Error(t, err, "find task after project deletion should fail")
	})
}
