package testsetup

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/identity"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/domain/testcase"
	"github.com/liang21/aitestos/internal/domain/testplan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertUserEqual 断言用户相等
func AssertUserEqual(t *testing.T, expected, actual *identity.User) {
	t.Helper()

	require.NotNil(t, actual, "user should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "user ID should match")
	assert.Equal(t, expected.Username(), actual.Username(), "username should match")
	assert.Equal(t, expected.Email(), actual.Email(), "email should match")
	assert.Equal(t, expected.Role(), actual.Role(), "role should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
	AssertTimeEqual(t, expected.UpdatedAt(), actual.UpdatedAt(), "updated_at")
}

// AssertProjectEqual 断言项目相等
func AssertProjectEqual(t *testing.T, expected, actual *domainproject.Project) {
	t.Helper()

	require.NotNil(t, actual, "project should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "project ID should match")
	assert.Equal(t, expected.Name(), actual.Name(), "project name should match")
	assert.Equal(t, expected.Prefix(), actual.Prefix(), "project prefix should match")
	assert.Equal(t, expected.Description(), actual.Description(), "project description should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
	AssertTimeEqual(t, expected.UpdatedAt(), actual.UpdatedAt(), "updated_at")
}

// AssertModuleEqual 断言模块相等
func AssertModuleEqual(t *testing.T, expected, actual *domainproject.Module) {
	t.Helper()

	require.NotNil(t, actual, "module should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "module ID should match")
	assert.Equal(t, expected.ProjectID(), actual.ProjectID(), "module project ID should match")
	assert.Equal(t, expected.Name(), actual.Name(), "module name should match")
	assert.Equal(t, expected.Abbreviation(), actual.Abbreviation(), "module abbreviation should match")
	assert.Equal(t, expected.Description(), actual.Description(), "module description should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
	AssertTimeEqual(t, expected.UpdatedAt(), actual.UpdatedAt(), "updated_at")
}

// AssertProjectConfigEqual 断言项目配置相等
func AssertProjectConfigEqual(t *testing.T, expected, actual *domainproject.ProjectConfig) {
	t.Helper()

	require.NotNil(t, actual, "project config should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "config ID should match")
	assert.Equal(t, expected.ProjectID(), actual.ProjectID(), "config project ID should match")
	assert.Equal(t, expected.Key(), actual.Key(), "config key should match")
	assert.Equal(t, expected.Value(), actual.Value(), "config value should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
	AssertTimeEqual(t, expected.UpdatedAt(), actual.UpdatedAt(), "updated_at")
}

// AssertTestCaseEqual 断言测试用例相等
func AssertTestCaseEqual(t *testing.T, expected, actual *testcase.TestCase) {
	t.Helper()

	require.NotNil(t, actual, "test case should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "case ID should match")
	assert.Equal(t, expected.ModuleID(), actual.ModuleID(), "case module ID should match")
	assert.Equal(t, expected.UserID(), actual.UserID(), "case user ID should match")
	assert.Equal(t, expected.Number(), actual.Number(), "case number should match")
	assert.Equal(t, expected.Title(), actual.Title(), "case title should match")
	assert.Equal(t, expected.Preconditions(), actual.Preconditions(), "case preconditions should match")
	assert.Equal(t, expected.Steps(), actual.Steps(), "case steps should match")
	assert.Equal(t, expected.ExpectedResult(), actual.ExpectedResult(), "case expected result should match")
	assert.Equal(t, expected.CaseType(), actual.CaseType(), "case type should match")
	assert.Equal(t, expected.Priority(), actual.Priority(), "case priority should match")
	assert.Equal(t, expected.Status(), actual.Status(), "case status should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
	AssertTimeEqual(t, expected.UpdatedAt(), actual.UpdatedAt(), "updated_at")
}

// AssertTestPlanEqual 断言测试计划相等
func AssertTestPlanEqual(t *testing.T, expected, actual *testplan.TestPlan) {
	t.Helper()

	require.NotNil(t, actual, "test plan should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "plan ID should match")
	assert.Equal(t, expected.ProjectID(), actual.ProjectID(), "plan project ID should match")
	assert.Equal(t, expected.CreatedBy(), actual.CreatedBy(), "plan created by should match")
	assert.Equal(t, expected.Name(), actual.Name(), "plan name should match")
	assert.Equal(t, expected.Status(), actual.Status(), "plan status should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
	AssertTimeEqual(t, expected.UpdatedAt(), actual.UpdatedAt(), "updated_at")
}

// AssertTestResultEqual 断言测试结果相等
func AssertTestResultEqual(t *testing.T, expected, actual *testplan.TestResult) {
	t.Helper()

	require.NotNil(t, actual, "test result should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "result ID should match")
	assert.Equal(t, expected.CaseID(), actual.CaseID(), "result case ID should match")
	assert.Equal(t, expected.PlanID(), actual.PlanID(), "result plan ID should match")
	assert.Equal(t, expected.ExecutedBy(), actual.ExecutedBy(), "result executor ID should match")
	assert.Equal(t, expected.Status(), actual.Status(), "result status should match")
	AssertTimeEqual(t, expected.ExecutedAt(), actual.ExecutedAt(), "executed_at")
}

// AssertDocumentEqual 断言文档相等
func AssertDocumentEqual(t *testing.T, expected, actual *knowledge.Document) {
	t.Helper()

	require.NotNil(t, actual, "document should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "document ID should match")
	assert.Equal(t, expected.ProjectID(), actual.ProjectID(), "document project ID should match")
	assert.Equal(t, expected.Name(), actual.Name(), "document name should match")
	assert.Equal(t, expected.Type(), actual.Type(), "document type should match")
	assert.Equal(t, expected.URL(), actual.URL(), "document URL should match")
	assert.Equal(t, expected.Status(), actual.Status(), "document status should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
	AssertTimeEqual(t, expected.UpdatedAt(), actual.UpdatedAt(), "updated_at")
}

// AssertDocumentChunkEqual 断言文档块相等
func AssertDocumentChunkEqual(t *testing.T, expected, actual *knowledge.DocumentChunk) {
	t.Helper()

	require.NotNil(t, actual, "document chunk should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "chunk ID should match")
	assert.Equal(t, expected.DocumentID(), actual.DocumentID(), "chunk document ID should match")
	assert.Equal(t, expected.ChunkIndex(), actual.ChunkIndex(), "chunk index should match")
	assert.Equal(t, expected.Content(), actual.Content(), "chunk content should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
}

// AssertGenerationTaskEqual 断言生成任务相等
func AssertGenerationTaskEqual(t *testing.T, expected, actual *generation.GenerationTask) {
	t.Helper()

	require.NotNil(t, actual, "generation task should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "task ID should match")
	assert.Equal(t, expected.ProjectID(), actual.ProjectID(), "task project ID should match")
	assert.Equal(t, expected.CreatedBy(), actual.CreatedBy(), "task created by should match")
	assert.Equal(t, expected.Status(), actual.Status(), "task status should match")
	assert.Equal(t, expected.Prompt(), actual.Prompt(), "task prompt should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
	AssertTimeEqual(t, expected.UpdatedAt(), actual.UpdatedAt(), "updated_at")
}

// AssertCaseDraftEqual 断言用例草稿相等
func AssertCaseDraftEqual(t *testing.T, expected, actual *generation.GeneratedCaseDraft) {
	t.Helper()

	require.NotNil(t, actual, "case draft should not be nil")
	assert.Equal(t, expected.ID(), actual.ID(), "draft ID should match")
	assert.Equal(t, expected.TaskID(), actual.TaskID(), "draft task ID should match")
	assert.Equal(t, expected.Title(), actual.Title(), "draft title should match")
	assert.Equal(t, expected.Status(), actual.Status(), "draft status should match")
	AssertTimeEqual(t, expected.CreatedAt(), actual.CreatedAt(), "created_at")
	AssertTimeEqual(t, expected.UpdatedAt(), actual.UpdatedAt(), "updated_at")
}

// AssertIDEqual 断言 UUID 相等
func AssertIDEqual(t *testing.T, expected, actual uuid.UUID, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, expected, actual, msgAndArgs...)
}

// AssertErrorIs 断言错误类型
func AssertErrorIs(t *testing.T, expected, actual error, msgAndArgs ...interface{}) {
	t.Helper()
	if expected != nil {
		require.Error(t, actual, msgAndArgs...)
		assert.ErrorIs(t, actual, expected, msgAndArgs...)
	} else {
		assert.NoError(t, actual, msgAndArgs...)
	}
}

// AssertTimeEqual 断言时间相等（允许 1 秒误差）
func AssertTimeEqual(t *testing.T, expected, actual time.Time, fieldName string) {
	t.Helper()
	if expected.IsZero() && actual.IsZero() {
		return
	}
	assert.WithinDuration(t, expected, actual, time.Second, "%s should be within 1 second", fieldName)
}

// AssertIDsEqual 断言 UUID 列表相等
func AssertIDsEqual(t *testing.T, expected, actual []uuid.UUID, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, len(expected), len(actual), "length should match")
	for i, id := range expected {
		assert.Equal(t, id, actual[i], msgAndArgs...)
	}
}

// AssertSliceLen 断言切片长度
func AssertSliceLen(t *testing.T, slice interface{}, expectedLen int, msgAndArgs ...interface{}) {
	t.Helper()
	switch v := slice.(type) {
	case []uuid.UUID:
		assert.Equal(t, expectedLen, len(v), msgAndArgs...)
	case []*identity.User:
		assert.Equal(t, expectedLen, len(v), msgAndArgs...)
	case []*domainproject.Project:
		assert.Equal(t, expectedLen, len(v), msgAndArgs...)
	case []*domainproject.Module:
		assert.Equal(t, expectedLen, len(v), msgAndArgs...)
	case []*testcase.TestCase:
		assert.Equal(t, expectedLen, len(v), msgAndArgs...)
	case []*testplan.TestPlan:
		assert.Equal(t, expectedLen, len(v), msgAndArgs...)
	default:
		t.Fatalf("unsupported slice type: %T", slice)
	}
}
