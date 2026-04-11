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

func TestProjectRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	projectRepo := project.NewProjectRepository(tc.DB)
	ctx := context.Background()

	t.Run("Save", func(t *testing.T) {
		tc.CleanupTest()

		tests := []struct {
			name    string
			builder *testsetup.ProjectBuilder
			wantErr error
		}{
			{
				name:    "save valid project",
				builder: testsetup.NewProjectBuilder(),
				wantErr: nil,
			},
			{
				name:    "save project with 2-char prefix",
				builder: testsetup.NewProjectBuilder().WithPrefix("AB"),
				wantErr: nil,
			},
			{
				name:    "save project with 4-char prefix",
				builder: testsetup.NewProjectBuilder().WithPrefix("ABCD"),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				project, err := tt.builder.Build()
				require.NoError(t, err, "build project should succeed")

				err = projectRepo.Save(ctx, project)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					found, err := projectRepo.FindByID(ctx, project.ID())
					require.NoError(t, err, "find project by ID should succeed")
					testsetup.AssertProjectEqual(t, project, found)
				}
			})
		}
	})

	t.Run("Save duplicate name", func(t *testing.T) {
		tc.CleanupTest()

		name := "duplicate_project"

		project1, err := testsetup.NewProjectBuilder().WithName(name).Build()
		require.NoError(t, err, "build project1 should succeed")
		require.NoError(t, projectRepo.Save(ctx, project1), "save project1 should succeed")

		project2, err := testsetup.NewProjectBuilder().WithName(name).Build()
		require.NoError(t, err, "build project2 should succeed")
		err = projectRepo.Save(ctx, project2)
		require.Error(t, err, "save duplicate name should fail")
	})

	t.Run("Save duplicate prefix", func(t *testing.T) {
		tc.CleanupTest()

		prefix := "DUP"

		project1, err := testsetup.NewProjectBuilder().WithPrefix(prefix).Build()
		require.NoError(t, err, "build project1 should succeed")
		require.NoError(t, projectRepo.Save(ctx, project1), "save project1 should succeed")

		project2, err := testsetup.NewProjectBuilder().WithPrefix(prefix).Build()
		require.NoError(t, err, "build project2 should succeed")
		err = projectRepo.Save(ctx, project2)
		require.Error(t, err, "save duplicate prefix should fail")
	})

	t.Run("FindByID", func(t *testing.T) {
		tc.CleanupTest()

		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")

		found, err := projectRepo.FindByID(ctx, project.ID())
		require.NoError(t, err, "find project by ID should succeed")
		testsetup.AssertProjectEqual(t, project, found)

		// 测试不存在的 ID
		_, err = projectRepo.FindByID(ctx, uuid.New())
		require.Error(t, err, "find non-existent project should fail")
		assert.ErrorIs(t, err, domainproject.ErrProjectNotFound, "error should be ErrProjectNotFound")
	})

	t.Run("FindByName", func(t *testing.T) {
		tc.CleanupTest()

		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")

		found, err := projectRepo.FindByName(ctx, project.Name())
		require.NoError(t, err, "find project by name should succeed")
		testsetup.AssertProjectEqual(t, project, found)

		// 测试不存在的名称
		_, err = projectRepo.FindByName(ctx, "notfound")
		require.Error(t, err, "find non-existent name should fail")
		assert.ErrorIs(t, err, domainproject.ErrProjectNotFound, "error should be ErrProjectNotFound")
	})

	t.Run("FindByPrefix", func(t *testing.T) {
		tc.CleanupTest()

		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")

		found, err := projectRepo.FindByPrefix(ctx, project.Prefix())
		require.NoError(t, err, "find project by prefix should succeed")
		testsetup.AssertProjectEqual(t, project, found)

		// 测试不存在的 prefix
		_, err = projectRepo.FindByPrefix(ctx, "NF")
		require.Error(t, err, "find non-existent prefix should fail")
		assert.ErrorIs(t, err, domainproject.ErrProjectNotFound, "error should be ErrProjectNotFound")
	})

	t.Run("Update", func(t *testing.T) {
		tc.CleanupTest()

		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")

		// 更新描述
		project.UpdateDescription("updated description")
		err = projectRepo.Update(ctx, project)
		require.NoError(t, err, "update project should succeed")

		found, err := projectRepo.FindByID(ctx, project.ID())
		require.NoError(t, err, "find project should succeed")
		assert.Equal(t, "updated description", found.Description(), "description should be updated")
	})

	t.Run("Delete", func(t *testing.T) {
		tc.CleanupTest()

		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")

		err = projectRepo.Delete(ctx, project.ID())
		require.NoError(t, err, "delete project should succeed")

		// 验证软删除
		_, err = projectRepo.FindByID(ctx, project.ID())
		require.Error(t, err, "find deleted project should fail")
		assert.ErrorIs(t, err, domainproject.ErrProjectNotFound, "error should be ErrProjectNotFound")
	})

	t.Run("FindAll", func(t *testing.T) {
		tc.CleanupTest()

		// 创建多个项目
		for i := 0; i < 5; i++ {
			project, err := testsetup.NewProjectBuilder().Build()
			require.NoError(t, err, "build project should succeed")
			require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")
		}

		// 测试分页
		projects, err := projectRepo.FindAll(ctx, domainproject.QueryOptions{
			Limit:  3,
			Offset: 0,
		})
		require.NoError(t, err, "find all projects should succeed")
		assert.Equal(t, 3, len(projects), "should return 3 projects")
	})
}
