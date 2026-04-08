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

func TestProjectConfigRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	projectRepo := repository.NewProjectRepository(tc.DB)
	configRepo := repository.NewProjectConfigRepository(tc.DB)
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
			builder *testsetup.ProjectConfigBuilder
			wantErr error
		}{
			{
				name:    "save valid config",
				builder: testsetup.NewProjectConfigBuilder(project.ID()),
				wantErr: nil,
			},
			{
				name: "save config with complex value",
				builder: testsetup.NewProjectConfigBuilder(project.ID()).
					WithValue(map[string]any{"nested": map[string]any{"key": "value"}}),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				config, err := tt.builder.Build()
				require.NoError(t, err, "build config should succeed")

				err = configRepo.Save(ctx, config)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					found, err := configRepo.FindByKey(ctx, project.ID(), config.Key())
					require.NoError(t, err, "find config by key should succeed")
					testsetup.AssertProjectConfigEqual(t, config, found)
				}
			})
		}
	})

	t.Run("Save duplicate key", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		key := "duplicate_key"

		config1, err := testsetup.NewProjectConfigBuilder(project.ID()).WithKey(key).Build()
		require.NoError(t, err, "build config1 should succeed")
		require.NoError(t, configRepo.Save(ctx, config1), "save config1 should succeed")

		config2, err := testsetup.NewProjectConfigBuilder(project.ID()).WithKey(key).Build()
		require.NoError(t, err, "build config2 should succeed")
		err = configRepo.Save(ctx, config2)
		require.Error(t, err, "save duplicate key should fail")
	})

	t.Run("Same key in different projects", func(t *testing.T) {
		tc.CleanupTest()

		project1 := createProject(t)
		project2 := createProject(t)
		key := "common_key"

		config1, err := testsetup.NewProjectConfigBuilder(project1.ID()).WithKey(key).Build()
		require.NoError(t, err, "build config1 should succeed")
		require.NoError(t, configRepo.Save(ctx, config1), "save config1 should succeed")

		// 不同项目可以有相同的 key
		config2, err := testsetup.NewProjectConfigBuilder(project2.ID()).WithKey(key).Build()
		require.NoError(t, err, "build config2 should succeed")
		require.NoError(t, configRepo.Save(ctx, config2), "save config2 should succeed")
	})

	t.Run("FindByProjectID", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)

		// 创建多个配置
		for i := 0; i < 3; i++ {
			config, err := testsetup.NewProjectConfigBuilder(project.ID()).Build()
			require.NoError(t, err, "build config should succeed")
			require.NoError(t, configRepo.Save(ctx, config), "save config should succeed")
		}

		configs, err := configRepo.FindByProjectID(ctx, project.ID())
		require.NoError(t, err, "find configs by project ID should succeed")
		assert.Equal(t, 3, len(configs), "should return 3 configs")
	})

	t.Run("FindByKey", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		key := "find_key"

		config, err := testsetup.NewProjectConfigBuilder(project.ID()).WithKey(key).Build()
		require.NoError(t, err, "build config should succeed")
		require.NoError(t, configRepo.Save(ctx, config), "save config should succeed")

		found, err := configRepo.FindByKey(ctx, project.ID(), key)
		require.NoError(t, err, "find config by key should succeed")
		testsetup.AssertProjectConfigEqual(t, config, found)

		// 测试不存在的 key
		_, err = configRepo.FindByKey(ctx, project.ID(), "notfound")
		require.Error(t, err, "find non-existent key should fail")
		assert.ErrorIs(t, err, domainproject.ErrConfigNotFound, "error should be ErrConfigNotFound")
	})

	t.Run("Update", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		config, err := testsetup.NewProjectConfigBuilder(project.ID()).Build()
		require.NoError(t, err, "build config should succeed")
		require.NoError(t, configRepo.Save(ctx, config), "save config should succeed")

		// 更新值
		newValue := map[string]any{"updated": true}
		config.UpdateValue(newValue)
		err = configRepo.Update(ctx, config)
		require.NoError(t, err, "update config should succeed")

		found, err := configRepo.FindByKey(ctx, project.ID(), config.Key())
		require.NoError(t, err, "find config should succeed")
		assert.Equal(t, newValue, found.Value(), "value should be updated")
	})

	t.Run("Delete", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		config, err := testsetup.NewProjectConfigBuilder(project.ID()).Build()
		require.NoError(t, err, "build config should succeed")
		require.NoError(t, configRepo.Save(ctx, config), "save config should succeed")

		err = configRepo.Delete(ctx, config.ID())
		require.NoError(t, err, "delete config should succeed")

		// 验证删除
		_, err = configRepo.FindByKey(ctx, project.ID(), config.Key())
		require.Error(t, err, "find deleted config should fail")
		assert.ErrorIs(t, err, domainproject.ErrConfigNotFound, "error should be ErrConfigNotFound")
	})

	t.Run("Cascade delete on project deletion", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)

		// 创建配置
		config, err := testsetup.NewProjectConfigBuilder(project.ID()).Build()
		require.NoError(t, err, "build config should succeed")
		require.NoError(t, configRepo.Save(ctx, config), "save config should succeed")

		// 删除项目
		err = projectRepo.Delete(ctx, project.ID())
		require.NoError(t, err, "delete project should succeed")

		// 验证配置也被删除（级联删除）
		_, err = configRepo.FindByKey(ctx, project.ID(), config.Key())
		require.Error(t, err, "find config after project deletion should fail")
	})
}
