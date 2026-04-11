package project_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/repository/project"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModuleRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	projectRepo := project.NewProjectRepository(tc.DB)
	moduleRepo := project.NewModuleRepository(tc.DB)
	ctx := context.Background()

	// 辅助函数：创建项目
	createProject := func(t *testing.T) *domainproject.Project {
		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")
		return project
	}

	t.Run("Save", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)

		tests := []struct {
			name    string
			builder *testsetup.ModuleBuilder
			wantErr error
		}{
			{
				name:    "save valid module",
				builder: testsetup.NewModuleBuilder(project.ID()),
				wantErr: nil,
			},
			{
				name:    "save module with 2-char abbreviation",
				builder: testsetup.NewModuleBuilder(project.ID()).WithAbbreviation("AB"),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				module, err := tt.builder.Build()
				require.NoError(t, err, "build module should succeed")

				err = moduleRepo.Save(ctx, module)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					found, err := moduleRepo.FindByID(ctx, module.ID())
					require.NoError(t, err, "find module by ID should succeed")
					testsetup.AssertModuleEqual(t, module, found)
				}
			})
		}
	})

	t.Run("Save duplicate name in same project", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		name := "duplicate_module"

		module1, err := testsetup.NewModuleBuilder(project.ID()).WithName(name).Build()
		require.NoError(t, err, "build module1 should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module1), "save module1 should succeed")

		module2, err := testsetup.NewModuleBuilder(project.ID()).WithName(name).Build()
		require.NoError(t, err, "build module2 should succeed")
		err = moduleRepo.Save(ctx, module2)
		require.Error(t, err, "save duplicate name should fail")
	})

	t.Run("Save duplicate abbreviation in same project", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		abbr := "DUP"

		module1, err := testsetup.NewModuleBuilder(project.ID()).WithAbbreviation(abbr).Build()
		require.NoError(t, err, "build module1 should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module1), "save module1 should succeed")

		module2, err := testsetup.NewModuleBuilder(project.ID()).WithAbbreviation(abbr).Build()
		require.NoError(t, err, "build module2 should succeed")
		err = moduleRepo.Save(ctx, module2)
		require.Error(t, err, "save duplicate abbreviation should fail")
	})

	t.Run("Save same name in different projects", func(t *testing.T) {
		tc.CleanupTest()

		project1 := createProject(t)
		project2 := createProject(t)
		name := "common_module"

		module1, err := testsetup.NewModuleBuilder(project1.ID()).WithName(name).Build()
		require.NoError(t, err, "build module1 should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module1), "save module1 should succeed")

		// 不同项目中可以有同名模块
		module2, err := testsetup.NewModuleBuilder(project2.ID()).WithName(name).Build()
		require.NoError(t, err, "build module2 should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module2), "save module2 should succeed")
	})

	t.Run("FindByID", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		module, err := testsetup.NewModuleBuilder(project.ID()).Build()
		require.NoError(t, err, "build module should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module), "save module should succeed")

		found, err := moduleRepo.FindByID(ctx, module.ID())
		require.NoError(t, err, "find module by ID should succeed")
		testsetup.AssertModuleEqual(t, module, found)

		// 测试不存在的 ID
		_, err = moduleRepo.FindByID(ctx, uuid.New())
		require.Error(t, err, "find non-existent module should fail")
		assert.ErrorIs(t, err, domainproject.ErrModuleNotFound, "error should be ErrModuleNotFound")
	})

	t.Run("FindByProjectID", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)

		// 创建多个模块
		for i := 0; i < 3; i++ {
			module, err := testsetup.NewModuleBuilder(project.ID()).Build()
			require.NoError(t, err, "build module should succeed")
			require.NoError(t, moduleRepo.Save(ctx, module), "save module should succeed")
		}

		modules, err := moduleRepo.FindByProjectID(ctx, project.ID())
		require.NoError(t, err, "find modules by project ID should succeed")
		assert.Equal(t, 3, len(modules), "should return 3 modules")
	})

	t.Run("FindByAbbreviation", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		abbrStr := "FIND"

		module, err := testsetup.NewModuleBuilder(project.ID()).WithAbbreviation(abbrStr).Build()
		require.NoError(t, err, "build module should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module), "save module should succeed")

		abbr, err := domainproject.ParseModuleAbbreviation(abbrStr)
		require.NoError(t, err, "parse abbreviation should succeed")

		found, err := moduleRepo.FindByAbbreviation(ctx, project.ID(), abbr)
		require.NoError(t, err, "find module by abbreviation should succeed")
		testsetup.AssertModuleEqual(t, module, found)

		// 测试不存在的缩写 - 使用无效缩写字符串
		/* 测试不存在的缩写 - 使用有效的缩写格式但不存在的值	*/
		invalidAbbr, _ := domainproject.ParseModuleAbbreviation("NF")
		_, err = moduleRepo.FindByAbbreviation(ctx, project.ID(), invalidAbbr)
		require.Error(t, err, "find non-existent abbreviation should fail")
	})

	t.Run("Update", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		module, err := testsetup.NewModuleBuilder(project.ID()).Build()
		require.NoError(t, err, "build module should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module), "save module should succeed")

		// 更新描述
		module.UpdateDescription("updated description")
		err = moduleRepo.Update(ctx, module)
		require.NoError(t, err, "update module should succeed")

		found, err := moduleRepo.FindByID(ctx, module.ID())
		require.NoError(t, err, "find module should succeed")
		assert.Equal(t, "updated description", found.Description(), "description should be updated")
	})

	t.Run("Delete", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		module, err := testsetup.NewModuleBuilder(project.ID()).Build()
		require.NoError(t, err, "build module should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module), "save module should succeed")

		err = moduleRepo.Delete(ctx, module.ID())
		require.NoError(t, err, "delete module should succeed")

		// 验证删除
		_, err = moduleRepo.FindByID(ctx, module.ID())
		require.Error(t, err, "find deleted module should fail")
		assert.ErrorIs(t, err, domainproject.ErrModuleNotFound, "error should be ErrModuleNotFound")
	})

	t.Run("Cascade delete on project deletion", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)

		// 创建模块
		module, err := testsetup.NewModuleBuilder(project.ID()).Build()
		require.NoError(t, err, "build module should succeed")
		require.NoError(t, moduleRepo.Save(ctx, module), "save module should succeed")

		// 删除项目
		err = projectRepo.Delete(ctx, project.ID())
		require.NoError(t, err, "delete project should succeed")

		// 验证模块也被删除（级联删除）
		_, err = moduleRepo.FindByID(ctx, module.ID())
		require.Error(t, err, "find module after project deletion should fail")
	})
}
