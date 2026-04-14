package testplan_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/identity"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	domaintestplan "github.com/liang21/aitestos/internal/domain/testplan"
	identityRepo "github.com/liang21/aitestos/internal/repository/identity"
	projectPackage "github.com/liang21/aitestos/internal/repository/project"
	testplanRepo "github.com/liang21/aitestos/internal/repository/testplan"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	userRepo := identityRepo.NewUserRepository(tc.DB)
	projectRepo := projectPackage.NewProjectRepository(tc.DB)
	planRepo := testplanRepo.NewTestPlanRepository(tc.DB)
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
			builder *testsetup.TestPlanBuilder
			wantErr error
		}{
			{
				name:    "save draft plan",
				builder: testsetup.NewTestPlanBuilder(project.ID(), user.ID()).WithStatus(domaintestplan.StatusDraft),
				wantErr: nil,
			},
			{
				name:    "save active plan",
				builder: testsetup.NewTestPlanBuilder(project.ID(), user.ID()).WithStatus(domaintestplan.StatusActive),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				plan, err := tt.builder.Build()
				require.NoError(t, err, "build plan should succeed")

				err = planRepo.Save(ctx, plan)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					found, err := planRepo.FindByID(ctx, plan.ID())
					require.NoError(t, err, "find plan by ID should succeed")
					testsetup.AssertTestPlanEqual(t, plan, found)
				}
			})
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		plan, err := testsetup.NewTestPlanBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, err, "build plan should succeed")
		require.NoError(t, planRepo.Save(ctx, plan), "save plan should succeed")

		found, err := planRepo.FindByID(ctx, plan.ID())
		require.NoError(t, err, "find plan by ID should succeed")
		testsetup.AssertTestPlanEqual(t, plan, found)

		// 测试不存在的 ID
		_, err = planRepo.FindByID(ctx, uuid.New())
		require.Error(t, err, "find non-existent plan should fail")
		assert.ErrorIs(t, err, domaintestplan.ErrPlanNotFound, "error should be ErrPlanNotFound")
	})

	t.Run("FindByProjectID", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		// 创建多个计划
		for i := 0; i < 3; i++ {
			plan, err := testsetup.NewTestPlanBuilder(project.ID(), user.ID()).Build()
			require.NoError(t, err, "build plan should succeed")
			require.NoError(t, planRepo.Save(ctx, plan), "save plan should succeed")
		}

		plans, err := planRepo.FindByProjectID(ctx, project.ID(), domaintestplan.QueryOptions{
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "find plans by project ID should succeed")
		assert.Equal(t, 3, len(plans), "should return 3 plans")
	})

	t.Run("AddCase", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		plan, err := testsetup.NewTestPlanBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, err, "build plan should succeed")
		require.NoError(t, planRepo.Save(ctx, plan), "save plan should succeed")

		caseID := uuid.New()
		err = planRepo.AddCase(ctx, plan.ID(), caseID)
		require.NoError(t, err, "add case to plan should succeed")

		// 验证关联
		caseIDs, err := planRepo.GetCaseIDs(ctx, plan.ID())
		require.NoError(t, err, "get case IDs should succeed")
		assert.Equal(t, 1, len(caseIDs), "should have 1 case")
		assert.Equal(t, caseID, caseIDs[0], "case ID should match")
	})

	t.Run("RemoveCase", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		plan, err := testsetup.NewTestPlanBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, err, "build plan should succeed")
		require.NoError(t, planRepo.Save(ctx, plan), "save plan should succeed")

		caseID := uuid.New()
		require.NoError(t, planRepo.AddCase(ctx, plan.ID(), caseID), "add case should succeed")

		err = planRepo.RemoveCase(ctx, plan.ID(), caseID)
		require.NoError(t, err, "remove case from plan should succeed")

		// 验证关联已删除
		caseIDs, err := planRepo.GetCaseIDs(ctx, plan.ID())
		require.NoError(t, err, "get case IDs should succeed")
		assert.Equal(t, 0, len(caseIDs), "should have 0 cases")
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		plan, err := testsetup.NewTestPlanBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, err, "build plan should succeed")
		require.NoError(t, planRepo.Save(ctx, plan), "save plan should succeed")

		// 状态流转: draft -> active -> completed -> archived
		err = planRepo.UpdateStatus(ctx, plan.ID(), domaintestplan.StatusActive)
		require.NoError(t, err, "update status to active should succeed")

		found, err := planRepo.FindByID(ctx, plan.ID())
		require.NoError(t, err, "find plan should succeed")
		assert.Equal(t, domaintestplan.StatusActive, found.Status(), "status should be active")
	})

	t.Run("Delete", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		plan, err := testsetup.NewTestPlanBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, err, "build plan should succeed")
		require.NoError(t, planRepo.Save(ctx, plan), "save plan should succeed")

		err = planRepo.Delete(ctx, plan.ID())
		require.NoError(t, err, "delete plan should succeed")

		// 验证软删除
		_, err = planRepo.FindByID(ctx, plan.ID())
		require.Error(t, err, "find deleted plan should fail")
		assert.ErrorIs(t, err, domaintestplan.ErrPlanNotFound, "error should be ErrPlanNotFound")
	})

	t.Run("FindByStatus", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project := createProject(t)

		// 创建不同状态的计划
		for i := 0; i < 2; i++ {
			plan, err := testsetup.NewTestPlanBuilder(project.ID(), user.ID()).
				WithStatus(domaintestplan.StatusDraft).Build()
			require.NoError(t, err, "build plan should succeed")
			require.NoError(t, planRepo.Save(ctx, plan), "save plan should succeed")
		}

		for i := 0; i < 3; i++ {
			plan, err := testsetup.NewTestPlanBuilder(project.ID(), user.ID()).
				WithStatus(domaintestplan.StatusActive).Build()
			require.NoError(t, err, "build plan should succeed")
			require.NoError(t, planRepo.Save(ctx, plan), "save plan should succeed")
		}

		// 查询 draft 状态的计划
		plans, err := planRepo.FindByStatus(ctx, domaintestplan.StatusDraft, domaintestplan.QueryOptions{
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "find plans by status should succeed")
		assert.GreaterOrEqual(t, len(plans), 2, "should return at least 2 draft plans")

		// 查询 active 状态的计划
		plans, err = planRepo.FindByStatus(ctx, domaintestplan.StatusActive, domaintestplan.QueryOptions{
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "find active plans should succeed")
		assert.GreaterOrEqual(t, len(plans), 3, "should return at least 3 active plans")
	})
}
