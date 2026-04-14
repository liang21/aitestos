package testcase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/identity"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/domain/testcase"
	identityrepo "github.com/liang21/aitestos/internal/repository/identity"
	repository "github.com/liang21/aitestos/internal/repository/project"
	testcaserepo "github.com/liang21/aitestos/internal/repository/testcase"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaseRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	userRepo := identityrepo.NewUserRepository(tc.DB)
	projectRepo := repository.NewProjectRepository(tc.DB)
	moduleRepo := repository.NewModuleRepository(tc.DB)
	caseRepo := testcaserepo.NewTestCaseRepository(tc.DB)
	ctx := context.Background()

	// 辅助函数：创建用户
	createUser := func(t *testing.T) *identity.User {
		user, err := testsetup.NewUserBuilder().Build()
		require.NoError(t, err, "build user should succeed")
		require.NoError(t, userRepo.Save(ctx, user), "save user should succeed")
		return user
	}

	// 辅助函数：创建项目和模块
	createProjectAndModule := func(t *testing.T) (*domainproject.Project, *domainproject.Module) {
		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")

		module, err := testsetup.NewModuleBuilder(project.ID()).Build()
		require.NoError(t, err, "build module should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module), "save module should succeed")

		return project, module
	}

	t.Run("Save", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		_, module := createProjectAndModule(t)

		tests := []struct {
			name    string
			builder *testsetup.TestCaseBuilder
			wantErr error
		}{
			{
				name:    "save valid test case",
				builder: testsetup.NewTestCaseBuilder(module.ID(), user.ID()),
				wantErr: nil,
			},
			{
				name: "save test case with AI metadata",
				builder: testsetup.NewTestCaseBuilder(module.ID(), user.ID()).
					WithAIMetadata(testcase.NewAiMetadata(
						uuid.New(),
						testcase.ConfidenceHigh,
						[]*testcase.ReferencedChunk{},
						"deepseek-v3",
					)),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tc, err := tt.builder.Build()
				require.NoError(t, err, "build test case should succeed")

				err = caseRepo.Save(ctx, tc)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					found, err := caseRepo.FindByID(ctx, tc.ID())
					require.NoError(t, err, "find test case by ID should succeed")
					testsetup.AssertTestCaseEqual(t, tc, found)
				}
			})
		}
	})

	t.Run("Save duplicate number", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		_, module := createProjectAndModule(t)
		number := testcase.CaseNumber("TP-TM-20260403-001")

		case1, err := testsetup.NewTestCaseBuilder(module.ID(), user.ID()).
			WithNumber(number).Build()
		require.NoError(t, err, "build case1 should succeed")
		require.NoError(t, caseRepo.Save(ctx, case1), "save case1 should succeed")

		case2, err := testsetup.NewTestCaseBuilder(module.ID(), user.ID()).
			WithNumber(number).Build()
		require.NoError(t, err, "build case2 should succeed")
		err = caseRepo.Save(ctx, case2)
		require.Error(t, err, "save duplicate number should fail")
	})

	t.Run("FindByNumber", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		_, module := createProjectAndModule(t)
		number := testcase.CaseNumber("TP-TM-20260403-002")

		tc, err := testsetup.NewTestCaseBuilder(module.ID(), user.ID()).
			WithNumber(number).Build()
		require.NoError(t, err, "build test case should succeed")
		require.NoError(t, caseRepo.Save(ctx, tc), "save test case should succeed")

		found, err := caseRepo.FindByNumber(ctx, number)
		require.NoError(t, err, "find test case by number should succeed")
		testsetup.AssertTestCaseEqual(t, tc, found)

		// 测试不存在的编号
		_, err = caseRepo.FindByNumber(ctx, "NF-NF-20260403-999")
		require.Error(t, err, "find non-existent number should fail")
		assert.ErrorIs(t, err, testcase.ErrCaseNotFound, "error should be ErrCaseNotFound")
	})

	t.Run("FindByModuleID", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		_, module := createProjectAndModule(t)

		// 创建多个用例
		for i := 0; i < 3; i++ {
			tc, err := testsetup.NewTestCaseBuilder(module.ID(), user.ID()).Build()
			require.NoError(t, err, "build test case should succeed")
			require.NoError(t, caseRepo.Save(ctx, tc), "save test case should succeed")
		}

		cases, err := caseRepo.FindByModuleID(ctx, module.ID(), testcase.QueryOptions{
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "find cases by module ID should succeed")
		assert.Equal(t, 3, len(cases), "should return 3 cases")
	})

	t.Run("FindByProjectID", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		project, module := createProjectAndModule(t)

		// 创建多个用例
		for i := 0; i < 3; i++ {
			tc, err := testsetup.NewTestCaseBuilder(module.ID(), user.ID()).Build()
			require.NoError(t, err, "build test case should succeed")
			require.NoError(t, caseRepo.Save(ctx, tc), "save test case should succeed")
		}

		cases, err := caseRepo.FindByProjectID(ctx, project.ID(), testcase.QueryOptions{
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err, "find cases by project ID should succeed")
		assert.Equal(t, 3, len(cases), "should return 3 cases")
	})

	t.Run("CountByDate", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		_, module := createProjectAndModule(t)
		today := time.Now()

		// 创建 3 个今天的用例
		for i := 0; i < 3; i++ {
			tc, err := testsetup.NewTestCaseBuilder(module.ID(), user.ID()).Build()
			require.NoError(t, err, "build test case should succeed")
			require.NoError(t, caseRepo.Save(ctx, tc), "save test case should succeed")
		}

		count, err := caseRepo.CountByDate(ctx, module.ID(), today)
		require.NoError(t, err, "count cases by date should succeed")
		assert.Equal(t, int64(3), count, "should count 3 cases for today")

		// 测试昨天的计数（应该是 0）
		count, err = caseRepo.CountByDate(ctx, module.ID(), today.AddDate(0, 0, -1))
		require.NoError(t, err, "count cases for yesterday should succeed")
		assert.Equal(t, int64(0), count, "should count 0 cases for yesterday")
	})

	t.Run("Update", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		_, module := createProjectAndModule(t)

		tc, err := testsetup.NewTestCaseBuilder(module.ID(), user.ID()).Build()
		require.NoError(t, err, "build test case should succeed")
		require.NoError(t, caseRepo.Save(ctx, tc), "save test case should succeed")

		// 更新状态
		tc.UpdateStatus(testcase.StatusPass)
		err = caseRepo.Update(ctx, tc)
		require.NoError(t, err, "update test case should succeed")

		found, err := caseRepo.FindByID(ctx, tc.ID())
		require.NoError(t, err, "find test case should succeed")
		assert.Equal(t, testcase.StatusPass, found.Status(), "status should be updated to pass")
	})

	t.Run("Delete", func(t *testing.T) {
		tc.CleanupTest()

		user := createUser(t)
		_, module := createProjectAndModule(t)

		tc, err := testsetup.NewTestCaseBuilder(module.ID(), user.ID()).Build()
		require.NoError(t, err, "build test case should succeed")
		require.NoError(t, caseRepo.Save(ctx, tc), "save test case should succeed")

		err = caseRepo.Delete(ctx, tc.ID())
		require.NoError(t, err, "delete test case should succeed")

		// 验证软删除
		_, err = caseRepo.FindByID(ctx, tc.ID())
		require.Error(t, err, "find deleted test case should fail")
		assert.ErrorIs(t, err, testcase.ErrCaseNotFound, "error should be ErrCaseNotFound")
	})
}
