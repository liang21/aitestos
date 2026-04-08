package knowledge_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	"github.com/liang21/aitestos/internal/repository/knowledge"
	"github.com/liang21/aitestos/internal/repository/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	projectRepo := repository.NewProjectRepository(tc.DB)
	docRepo := repository.NewDocumentRepository(tc.DB)
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
			builder *testsetup.DocumentBuilder
			wantErr error
		}{
			{
				name:    "save PRD document",
				builder: testsetup.NewDocumentBuilder(project.ID()).WithType(knowledge.TypePRD),
				wantErr: nil,
			},
			{
				name:    "save Figma document",
				builder: testsetup.NewDocumentBuilder(project.ID()).WithType(knowledge.TypeFigma),
				wantErr: nil,
			},
			{
				name:    "save API spec document",
				builder: testsetup.NewDocumentBuilder(project.ID()).WithType(knowledge.TypeAPISpec),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				doc, err := tt.builder.Build()
				require.NoError(t, err, "build document should succeed")

				err = docRepo.Save(ctx, doc)
				testsetup.AssertErrorIs(t, tt.wantErr, err)

				if tt.wantErr == nil {
					found, err := docRepo.FindByID(ctx, doc.ID())
					require.NoError(t, err, "find document by ID should succeed")
					testsetup.AssertDocumentEqual(t, doc, found)
				}
			})
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		doc, err := testsetup.NewDocumentBuilder(project.ID()).Build()
		require.NoError(t, err, "build document should succeed")
		require.NoError(t, docRepo.Save(ctx, doc), "save document should succeed")

		found, err := docRepo.FindByID(ctx, doc.ID())
		require.NoError(t, err, "find document by ID should succeed")
		testsetup.AssertDocumentEqual(t, doc, found)

		// 测试不存在的 ID
		_, err = docRepo.FindByID(ctx, uuid.New())
		require.Error(t, err, "find non-existent document should fail")
		assert.ErrorIs(t, err, knowledge.ErrDocumentNotFound, "error should be ErrDocumentNotFound")
	})

	t.Run("FindByProjectID", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)

		// 创建多个文档
		for i := 0; i < 3; i++ {
			doc, err := testsetup.NewDocumentBuilder(project.ID()).Build()
			require.NoError(t, err, "build document should succeed")
			require.NoError(t, docRepo.Save(ctx, doc), "save document should succeed")
		}

		docs, err := docRepo.FindByProjectID(ctx, project.ID())
		require.NoError(t, err, "find documents by project ID should succeed")
		assert.Equal(t, 3, len(docs), "should return 3 documents")
	})

	t.Run("FindByType", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)

		// 创建不同类型的文档
		for i := 0; i < 2; i++ {
			doc, err := testsetup.NewDocumentBuilder(project.ID()).WithType(knowledge.TypePRD).Build()
			require.NoError(t, err, "build PRD document should succeed")
			require.NoError(t, docRepo.Save(ctx, doc), "save PRD document should succeed")
		}

		for i := 0; i < 3; i++ {
			doc, err := testsetup.NewDocumentBuilder(project.ID()).WithType(knowledge.TypeFigma).Build()
			require.NoError(t, err, "build Figma document should succeed")
			require.NoError(t, docRepo.Save(ctx, doc), "save Figma document should succeed")
		}

		// 查询 PRD 文档
		prdDocs, err := docRepo.FindByType(ctx, project.ID(), knowledge.TypePRD)
		require.NoError(t, err, "find PRD documents should succeed")
		assert.Equal(t, 2, len(prdDocs), "should return 2 PRD documents")

		// 查询 Figma 文档
		figmaDocs, err := docRepo.FindByType(ctx, project.ID(), knowledge.TypeFigma)
		require.NoError(t, err, "find Figma documents should succeed")
		assert.Equal(t, 3, len(figmaDocs), "should return 3 Figma documents")
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		doc, err := testsetup.NewDocumentBuilder(project.ID()).Build()
		require.NoError(t, err, "build document should succeed")
		require.NoError(t, docRepo.Save(ctx, doc), "save document should succeed")

		// 更新状态: pending -> processing
		err = docRepo.UpdateStatus(ctx, doc.ID(), knowledge.StatusProcessing)
		require.NoError(t, err, "update status to processing should succeed")

		found, err := docRepo.FindByID(ctx, doc.ID())
		require.NoError(t, err, "find document should succeed")
		assert.Equal(t, knowledge.StatusProcessing, found.Status(), "status should be processing")

		// 更新状态: processing -> completed
		err = docRepo.UpdateStatus(ctx, doc.ID(), knowledge.StatusCompleted)
		require.NoError(t, err, "update status to completed should succeed")

		found, err = docRepo.FindByID(ctx, doc.ID())
		require.NoError(t, err, "find document should succeed")
		assert.Equal(t, knowledge.StatusCompleted, found.Status(), "status should be completed")
	})

	t.Run("UpdateContentText", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		doc, err := testsetup.NewDocumentBuilder(project.ID()).Build()
		require.NoError(t, err, "build document should succeed")
		require.NoError(t, docRepo.Save(ctx, doc), "save document should succeed")

		content := "This is the extracted text content"
		err = docRepo.UpdateContentText(ctx, doc.ID(), content)
		require.NoError(t, err, "update content text should succeed")

		found, err := docRepo.FindByID(ctx, doc.ID())
		require.NoError(t, err, "find document should succeed")
		assert.Equal(t, content, found.ContentText(), "content text should be updated")
	})

	t.Run("Delete", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)
		doc, err := testsetup.NewDocumentBuilder(project.ID()).Build()
		require.NoError(t, err, "build document should succeed")
		require.NoError(t, docRepo.Save(ctx, doc), "save document should succeed")

		err = docRepo.Delete(ctx, doc.ID())
		require.NoError(t, err, "delete document should succeed")

		// 验证删除
		_, err = docRepo.FindByID(ctx, doc.ID())
		require.Error(t, err, "find deleted document should fail")
		assert.ErrorIs(t, err, knowledge.ErrDocumentNotFound, "error should be ErrDocumentNotFound")
	})

	t.Run("FindByStatus", func(t *testing.T) {
		tc.CleanupTest()

		project := createProject(t)

		// 创建不同状态的文档
		for i := 0; i < 2; i++ {
			doc, err := testsetup.NewDocumentBuilder(project.ID()).
				WithStatus(knowledge.StatusPending).Build()
			require.NoError(t, err, "build document should succeed")
			require.NoError(t, docRepo.Save(ctx, doc), "save document should succeed")
		}

		for i := 0; i < 3; i++ {
			doc, err := testsetup.NewDocumentBuilder(project.ID()).
				WithStatus(knowledge.StatusCompleted).Build()
			require.NoError(t, err, "build document should succeed")
			require.NoError(t, docRepo.Save(ctx, doc), "save document should succeed")
		}

		// 查询 pending 状态的文档
		pendingDocs, err := docRepo.FindByStatus(ctx, knowledge.StatusPending)
		require.NoError(t, err, "find pending documents should succeed")
		assert.GreaterOrEqual(t, len(pendingDocs), 2, "should return at least 2 pending documents")

		// 查询 completed 状态的文档
		completedDocs, err := docRepo.FindByStatus(ctx, knowledge.StatusCompleted)
		require.NoError(t, err, "find completed documents should succeed")
		assert.GreaterOrEqual(t, len(completedDocs), 3, "should return at least 3 completed documents")
	})
}
