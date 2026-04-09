package testplan_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/identity"
	"github.com/liang21/aitestos/internal/domain/testcase"
	domaintestplan "github.com/liang21/aitestos/internal/domain/testplan"
	identityRepo "github.com/liang21/aitestos/internal/repository/identity"
	projectPackage "github.com/liang21/aitestos/internal/repository/project"
	testcaseRepo "github.com/liang21/aitestos/internal/repository/testcase"
	testplanRepo "github.com/liang21/aitestos/internal/repository/testplan"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResultRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	userRepo := identityRepo.NewUserRepository(tc.DB)
	projectRepo := projectPackage.NewProjectRepository(tc.DB)
	moduleRepo := projectPackage.NewModuleRepository(tc.DB)
	caseRepo := testcaseRepo.NewTestCaseRepository(tc.DB)
	planRepo := testplanRepo.NewTestPlanRepository(tc.DB)
	resultRepo := testplanRepo.NewTestResultRepository(tc.DB)
	ctx := context.Background()

	// 辅助函数：创建用户
	createUser := func(t *testing.T) *identity.User {
		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")
		require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")
		return user
	}

	// 辅助函数：创建项目、模块、用例、计划
	createTestData := func(t *testing.T) (*identity.User, *testcase.TestCase, *domaintestplan.TestPlan) {
		user := createUser(t)

		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")

		module, err := testsetup.NewModuleBuilder(project.ID()).Build()
		require.NoError(t, err, "build module should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module), "save module should succeed")

		tc, err := testsetup.NewTestCaseBuilder(module.ID(), user.ID()).Build()
		require.NoError(t, err, "build test case should succeed")
		require.NoError(t, caseRepo.Save(ctx, tc), "save test case should succeed")

		plan, err := testsetup.NewTestPlanBuilder(project.ID(), user.ID()).Build()
		require.NoError(t, err, "build plan should succeed")
		require.NoError(t, planRepo.Save(ctx, plan), "save plan should succeed")

		return user, tc, plan
	}

	t.Run("Save", func(t *testing.T) {
		tc.CleanupTest()

		user, testCase, plan := createTestData(t)

		tests := []struct {
			name    string
			builder *testsetup.TestResultBuilder
			wantErr error
		}{
			{
				name:    "save pass result",
				builder: testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).WithResult(domaintestplan.ResultPass),
				wantErr: nil,
			},
			{
				name:    "save fail result",
				builder: testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).WithResult(domaintestplan.ResultFail),
				wantErr: nil,
			},
			{
				name:    "save block result",
				builder: testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).WithResult(domaintestplan.ResultBlock),
				wantErr: nil,
			},
			{
				name:    "save skip result",
				builder: testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).WithResult(domaintestplan.ResultSkip),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := tt.builder.Build()
				require.NoError(t, err, "build result should succeed")

				err = resultRepo.Save(ctx, result)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					found, err := resultRepo.FindByID(ctx, result.ID())
					require.NoError(t, err, "find result by ID should succeed")
					testsetup.AssertTestResultEqual(t, result, found)
				}
			})
		}
	})

	t.Run("FindByPlanID", func(t *testing.T) {
		tc.CleanupTest()

		user, testCase, plan := createTestData(t)

		// 创建多个执行结果
		for i := 0; i < 3; i++ {
			result, err := testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).Build()
			require.NoError(t, err, "build result should succeed")
			require.NoError(t, resultRepo.Save(ctx, result), "save result should succeed")
		}

		results, err := resultRepo.FindByPlanID(ctx, plan.ID(), domaintestplan.QueryOptions{Limit: 10, Offset: 0})
		require.NoError(t, err, "find results by plan ID should succeed")
		assert.Equal(t, 3, len(results), "should return 3 results")
	})

	t.Run("FindByCaseID", func(t *testing.T) {
		tc.CleanupTest()

		user, testCase, plan := createTestData(t)

		// 创建多个执行结果（不同计划）
		for i := 0; i < 3; i++ {
			result, err := testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).Build()
			require.NoError(t, err, "build result should succeed")
			require.NoError(t, resultRepo.Save(ctx, result), "save result should succeed")
		}

		results, err := resultRepo.FindByCaseID(ctx, testCase.ID(), domaintestplan.QueryOptions{Limit: 10, Offset: 0})
		require.NoError(t, err, "find results by case ID should succeed")
		assert.GreaterOrEqual(t, len(results), 3, "should return at least 3 results")
	})

	t.Run("CountByStatus", func(t *testing.T) {
		tc.CleanupTest()

		user, testCase, plan := createTestData(t)

		// 创建不同状态的结果
		statuses := []domaintestplan.ResultStatus{
			domaintestplan.ResultPass,
			domaintestplan.ResultPass,
			domaintestplan.ResultFail,
			domaintestplan.ResultBlock,
			domaintestplan.ResultSkip,
		}

		for _, status := range statuses {
			result, err := testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).
				WithResult(status).Build()
			require.NoError(t, err, "build result should succeed")
			require.NoError(t, resultRepo.Save(ctx, result), "save result should succeed")
		}

		// 统计各状态数量
		counts, err := resultRepo.CountByStatus(ctx, plan.ID())
		require.NoError(t, err, "count by status should succeed")
		assert.Equal(t, 2, counts[domaintestplan.ResultPass], "should have 2 pass results")
		assert.Equal(t, 1, counts[domaintestplan.ResultFail], "should have 1 fail result")
		assert.Equal(t, 1, counts[domaintestplan.ResultBlock], "should have 1 block result")
		assert.Equal(t, 1, counts[domaintestplan.ResultSkip], "should have 1 skip result")
	})

	t.Run("FindLatestByCaseID", func(t *testing.T) {
		tc.CleanupTest()

		user, testCase, plan := createTestData(t)

		// 创建多个结果（不同时间）
		time.Sleep(10 * time.Millisecond)
		result1, err := testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).
			WithResult(domaintestplan.ResultFail).Build()
		require.NoError(t, err, "build result1 should succeed")
		require.NoError(t, resultRepo.Save(ctx, result1), "save result1 should succeed")

		time.Sleep(10 * time.Millisecond)
		result2, err := testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).
			WithResult(domaintestplan.ResultPass).Build()
		require.NoError(t, err, "build result2 should succeed")
		require.NoError(t, resultRepo.Save(ctx, result2), "save result2 should succeed")

		// 查询最新结果
		latest, err := resultRepo.FindLatestByCaseID(ctx, testCase.ID())
		require.NoError(t, err, "find latest result should succeed")
		assert.Equal(t, result2.ID(), latest.ID(), "should return the latest result")
		assert.Equal(t, domaintestplan.ResultPass, latest.Status(), "latest result should be pass")
	})

	t.Run("FindByPlanIDAndCaseID", func(t *testing.T) {
		tc.CleanupTest()

		user, testCase, plan := createTestData(t)

		result, err := testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).Build()
		require.NoError(t, err, "build result should succeed")
		require.NoError(t, resultRepo.Save(ctx, result), "save result should succeed")

		found, err := resultRepo.FindByPlanIDAndCaseID(ctx, plan.ID(), testCase.ID())
		require.NoError(t, err, "find result by plan and case ID should succeed")
		require.Len(t, found, 1, "should find 1 result")
		testsetup.AssertTestResultEqual(t, result, found[0])

		// 测试不存在的组合
		found, err = resultRepo.FindByPlanIDAndCaseID(ctx, uuid.New(), uuid.New())
		require.NoError(t, err, "find non-existent result should not error")
		assert.Len(t, found, 0, "should find no results for non-existent combination")
	})

	t.Run("DeleteByPlanID", func(t *testing.T) {
		tc.CleanupTest()

		user, testCase, plan := createTestData(t)

		// 创建多个结果
		for i := 0; i < 3; i++ {
			result, err := testsetup.NewTestResultBuilder(testCase.ID(), plan.ID(), user.ID()).Build()
			require.NoError(t, err, "build result should succeed")
			require.NoError(t, resultRepo.Save(ctx, result), "save result should succeed")
		}

		// 删除计划的所有结果
		err := resultRepo.DeleteByPlanID(ctx, plan.ID())
		require.NoError(t, err, "delete results by plan ID should succeed")

		// 验证删除
		results, err := resultRepo.FindByPlanID(ctx, plan.ID(), domaintestplan.QueryOptions{Limit: 10, Offset: 0})
		require.NoError(t, err, "find results should succeed")
		assert.Equal(t, 0, len(results), "should have no results after deletion")
	})
}
