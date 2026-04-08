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

func TestDocumentChunkRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tc := testsetup.SetupTest(t)
	defer tc.CleanupTest()

	projectRepo := repository.NewProjectRepository(tc.DB)
	docRepo := repository.NewDocumentRepository(tc.DB)
	chunkRepo := repository.NewDocumentChunkRepository(tc.DB)
	ctx := context.Background()

	// 辅助函数：创建项目和文档
	createProjectAndDocument := func(t *testing.T) (*domainproject.Project, *knowledge.Document) {
		project, err := testsetup.NewProjectBuilder().Build()
		require.NoError(t, err, "build project should succeed")
		require.NoError(t, projectRepo.Save(ctx, project), "save project should succeed")

		doc, err := testsetup.NewDocumentBuilder(project.ID()).Build()
		require.NoError(t, err, "build document should succeed")
		require.NoError(t, docRepo.Save(ctx, doc), "save document should succeed")

		return project, doc
	}

	t.Run("Save", func(t *testing.T) {
		tc.CleanupTest()

		_, doc := createProjectAndDocument(t)

		chunk, err := testsetup.NewDocumentChunkBuilder(doc.ID(), 0).Build()
		require.NoError(t, err, "build chunk should succeed")

		err = chunkRepo.Save(ctx, chunk)
		require.NoError(t, err, "save chunk should succeed")

		found, err := chunkRepo.FindByID(ctx, chunk.ID())
		require.NoError(t, err, "find chunk by ID should succeed")
		testsetup.AssertDocumentChunkEqual(t, chunk, found)
	})

	t.Run("BatchSave", func(t *testing.T) {
		tc.CleanupTest()

		_, doc := createProjectAndDocument(t)

		// 批量创建块
		chunks := make([]*knowledge.DocumentChunk, 5)
		for i := 0; i < 5; i++ {
			chunk, err := testsetup.NewDocumentChunkBuilder(doc.ID(), i).Build()
			require.NoError(t, err, "build chunk should succeed")
			chunks[i] = chunk
		}

		err := chunkRepo.BatchSave(ctx, chunks)
		require.NoError(t, err, "batch save chunks should succeed")

		// 验证所有块都已保存
		found, err := chunkRepo.FindByDocumentID(ctx, doc.ID())
		require.NoError(t, err, "find chunks by document ID should succeed")
		assert.Equal(t, 5, len(found), "should return 5 chunks")
	})

	t.Run("FindByDocumentID", func(t *testing.T) {
		tc.CleanupTest()

		_, doc := createProjectAndDocument(t)

		// 创建多个块
		for i := 0; i < 3; i++ {
			chunk, err := testsetup.NewDocumentChunkBuilder(doc.ID(), i).Build()
			require.NoError(t, err, "build chunk should succeed")
			require.NoError(t, chunkRepo.Save(ctx, chunk), "save chunk should succeed")
		}

		chunks, err := chunkRepo.FindByDocumentID(ctx, doc.ID())
		require.NoError(t, err, "find chunks by document ID should succeed")
		assert.Equal(t, 3, len(chunks), "should return 3 chunks")

		// 验证按 chunk_index 排序
		for i, chunk := range chunks {
			assert.Equal(t, i, chunk.ChunkIndex(), "chunk should be ordered by index")
		}
	})

	t.Run("FindByChunkIndex", func(t *testing.T) {
		tc.CleanupTest()

		_, doc := createProjectAndDocument(t)

		chunk, err := testsetup.NewDocumentChunkBuilder(doc.ID(), 5).Build()
		require.NoError(t, err, "build chunk should succeed")
		require.NoError(t, chunkRepo.Save(ctx, chunk), "save chunk should succeed")

		found, err := chunkRepo.FindByChunkIndex(ctx, doc.ID(), 5)
		require.NoError(t, err, "find chunk by index should succeed")
		testsetup.AssertDocumentChunkEqual(t, chunk, found)

		// 测试不存在的索引
		_, err = chunkRepo.FindByChunkIndex(ctx, doc.ID(), 999)
		require.Error(t, err, "find non-existent chunk index should fail")
		assert.ErrorIs(t, err, knowledge.ErrChunkNotFound, "error should be ErrChunkNotFound")
	})

	t.Run("FindByID", func(t *testing.T) {
		tc.CleanupTest()

		_, doc := createProjectAndDocument(t)

		chunk, err := testsetup.NewDocumentChunkBuilder(doc.ID(), 0).Build()
		require.NoError(t, err, "build chunk should succeed")
		require.NoError(t, chunkRepo.Save(ctx, chunk), "save chunk should succeed")

		found, err := chunkRepo.FindByID(ctx, chunk.ID())
		require.NoError(t, err, "find chunk by ID should succeed")
		testsetup.AssertDocumentChunkEqual(t, chunk, found)

		// 测试不存在的 ID
		_, err = chunkRepo.FindByID(ctx, uuid.New())
		require.Error(t, err, "find non-existent chunk should fail")
		assert.ErrorIs(t, err, knowledge.ErrChunkNotFound, "error should be ErrChunkNotFound")
	})

	t.Run("DeleteByDocumentID", func(t *testing.T) {
		tc.CleanupTest()

		_, doc := createProjectAndDocument(t)

		// 创建多个块
		for i := 0; i < 3; i++ {
			chunk, err := testsetup.NewDocumentChunkBuilder(doc.ID(), i).Build()
			require.NoError(t, err, "build chunk should succeed")
			require.NoError(t, chunkRepo.Save(ctx, chunk), "save chunk should succeed")
		}

		// 删除文档的所有块
		err := chunkRepo.DeleteByDocumentID(ctx, doc.ID())
		require.NoError(t, err, "delete chunks by document ID should succeed")

		// 验证删除
		chunks, err := chunkRepo.FindByDocumentID(ctx, doc.ID())
		require.NoError(t, err, "find chunks should succeed")
		assert.Equal(t, 0, len(chunks), "should have no chunks after deletion")
	})

	t.Run("Cascade delete on document deletion", func(t *testing.T) {
		tc.CleanupTest()

		_, doc := createProjectAndDocument(t)

		// 创建块
		chunk, err := testsetup.NewDocumentChunkBuilder(doc.ID(), 0).Build()
		require.NoError(t, err, "build chunk should succeed")
		require.NoError(t, chunkRepo.Save(ctx, chunk), "save chunk should succeed")

		// 删除文档
		err = docRepo.Delete(ctx, doc.ID())
		require.NoError(t, err, "delete document should succeed")

		// 验证块也被删除（级联删除）
		_, err = chunkRepo.FindByID(ctx, chunk.ID())
		require.Error(t, err, "find chunk after document deletion should fail")
	})

	t.Run("Update", func(t *testing.T) {
		tc.CleanupTest()

		_, doc := createProjectAndDocument(t)

		chunk, err := testsetup.NewDocumentChunkBuilder(doc.ID(), 0).Build()
		require.NoError(t, err, "build chunk should succeed")
		require.NoError(t, chunkRepo.Save(ctx, chunk), "save chunk should succeed")

		// 更新内容
		newContent := "Updated chunk content"
		chunk.UpdateContent(newContent)
		err = chunkRepo.Update(ctx, chunk)
		require.NoError(t, err, "update chunk should succeed")

		found, err := chunkRepo.FindByID(ctx, chunk.ID())
		require.NoError(t, err, "find chunk should succeed")
		assert.Equal(t, newContent, found.Content(), "content should be updated")
	})

	t.Run("CountByDocumentID", func(t *testing.T) {
		tc.CleanupTest()

		_, doc := createProjectAndDocument(t)

		// 创建多个块
		for i := 0; i < 7; i++ {
			chunk, err := testsetup.NewDocumentChunkBuilder(doc.ID(), i).Build()
			require.NoError(t, err, "build chunk should succeed")
			require.NoError(t, chunkRepo.Save(ctx, chunk), "save chunk should succeed")
		}

		count, err := chunkRepo.CountByDocumentID(ctx, doc.ID())
		require.NoError(t, err, "count chunks by document ID should succeed")
		assert.Equal(t, int64(7), count, "should count 7 chunks")
	})
}
