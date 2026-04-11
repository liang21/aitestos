package testsetup

import (
	"time"

	"github.com/google/uuid"
	"github.com/liang21/aitestos/internal/domain/generation"
	"github.com/liang21/aitestos/internal/domain/identity"
	"github.com/liang21/aitestos/internal/domain/knowledge"
	domainproject "github.com/liang21/aitestos/internal/domain/project"
	"github.com/liang21/aitestos/internal/domain/testcase"
	"github.com/liang21/aitestos/internal/domain/testplan"
)

// UserBuilder 用户构建器
type UserBuilder struct {
	username string
	email    string
	password string
	role     identity.UserRole
}

// NewUserBuilder 创建用户构建器
func NewUserBuilder() *UserBuilder {
	uid := uuid.New().String()[:8]
	return &UserBuilder{
		username: "testuser_" + uid,
		email:    "test_" + uid + "@example.com",
		password: "TestPass123!",
		role:     identity.RoleNormal,
	}
}

// WithUsername 设置用户名
func (b *UserBuilder) WithUsername(username string) *UserBuilder {
	b.username = username
	return b
}

// WithEmail 设置邮箱
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.email = email
	return b
}

// WithPassword 设置密码
func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.password = password
	return b
}

// WithRole 设置角色
func (b *UserBuilder) WithRole(role identity.UserRole) *UserBuilder {
	b.role = role
	return b
}

// Build 构建用户
func (b *UserBuilder) Build() (*identity.User, error) {
	return identity.NewUser(b.username, b.email, b.password, b.role)
}

// ProjectBuilder 项目构建器
type ProjectBuilder struct {
	name        string
	prefix      string
	description string
}

// NewProjectBuilder 创建项目构建器
func NewProjectBuilder() *ProjectBuilder {
	uid := uuid.New().String()[:6]
	return &ProjectBuilder{
		name:        "Test Project " + uid,
		prefix:      "TP" + uid[:2],
		description: "Test description",
	}
}

// WithName 设置项目名称
func (b *ProjectBuilder) WithName(name string) *ProjectBuilder {
	b.name = name
	return b
}

// WithPrefix 设置项目前缀
func (b *ProjectBuilder) WithPrefix(prefix string) *ProjectBuilder {
	b.prefix = prefix
	return b
}

// WithDescription 设置项目描述
func (b *ProjectBuilder) WithDescription(desc string) *ProjectBuilder {
	b.description = desc
	return b
}

// Build 构建项目
func (b *ProjectBuilder) Build() (*domainproject.Project, error) {
	return domainproject.NewProject(b.name, b.prefix, b.description)
}

// ModuleBuilder 模块构建器
type ModuleBuilder struct {
	projectID    uuid.UUID
	name         string
	abbreviation string
	description  string
}

// NewModuleBuilder 创建模块构建器
func NewModuleBuilder(projectID uuid.UUID) *ModuleBuilder {
	uid := uuid.New().String()[:4]
	return &ModuleBuilder{
		projectID:    projectID,
		name:         "Test Module " + uid,
		abbreviation: "TM" + uid[:2],
		description:  "Test module description",
	}
}

// WithName 设置模块名称
func (b *ModuleBuilder) WithName(name string) *ModuleBuilder {
	b.name = name
	return b
}

// WithAbbreviation 设置模块缩写
func (b *ModuleBuilder) WithAbbreviation(abbr string) *ModuleBuilder {
	b.abbreviation = abbr
	return b
}

// WithDescription 设置模块描述
func (b *ModuleBuilder) WithDescription(desc string) *ModuleBuilder {
	b.description = desc
	return b
}

// Build 构建模块
func (b *ModuleBuilder) Build() (*domainproject.Module, error) {
	userID := uuid.New() // 测试时使用随机用户ID
	return domainproject.NewModule(b.projectID, b.name, b.abbreviation, b.description, userID)
}

// ProjectConfigBuilder 项目配置构建器
type ProjectConfigBuilder struct {
	projectID uuid.UUID
	key       string
	value     map[string]any
}

// NewProjectConfigBuilder 创建项目配置构建器
func NewProjectConfigBuilder(projectID uuid.UUID) *ProjectConfigBuilder {
	uid := uuid.New().String()[:8]
	return &ProjectConfigBuilder{
		projectID: projectID,
		key:       "config_key_" + uid,
		value:     map[string]any{"setting": "value"},
	}
}

// WithKey 设置配置键
func (b *ProjectConfigBuilder) WithKey(key string) *ProjectConfigBuilder {
	b.key = key
	return b
}

// WithValue 设置配置值
func (b *ProjectConfigBuilder) WithValue(value map[string]any) *ProjectConfigBuilder {
	b.value = value
	return b
}

// Build 构建项目配置
func (b *ProjectConfigBuilder) Build() (*domainproject.ProjectConfig, error) {
	userID := uuid.New() // 测试时使用随机用户ID
	return domainproject.NewProjectConfig(b.projectID, b.key, b.value, userID.String())
}

// TestCaseBuilder 测试用例构建器
type TestCaseBuilder struct {
	moduleID      uuid.UUID
	userID        uuid.UUID
	number        testcase.CaseNumber
	title         string
	preconditions testcase.Preconditions
	steps         testcase.Steps
	expected      testcase.ExpectedResult
	caseType      testcase.CaseType
	priority      testcase.Priority
	aiMetadata    *testcase.AiMetadata
}

// NewTestCaseBuilder 创建测试用例构建器
func NewTestCaseBuilder(moduleID, userID uuid.UUID) *TestCaseBuilder {
	uid := uuid.New().String()[:8]
	return &TestCaseBuilder{
		moduleID: moduleID,
		userID:   userID,
		title:    "Test Case " + uid,
		preconditions: testcase.Preconditions{
			"用户已登录",
			"项目已创建",
		},
		steps: testcase.Steps{
			"步骤1：打开页面",
			"步骤2：点击按钮",
			"步骤3：验证结果",
		},
		expected: testcase.ExpectedResult{
			"status":  "success",
			"message": "操作成功",
		},
		caseType:   testcase.CaseTypeFunctionality,
		priority:   testcase.PriorityP2,
		aiMetadata: nil,
	}
}

// WithTitle 设置标题
func (b *TestCaseBuilder) WithTitle(title string) *TestCaseBuilder {
	b.title = title
	return b
}

// WithNumber 设置编号
func (b *TestCaseBuilder) WithNumber(number testcase.CaseNumber) *TestCaseBuilder {
	b.number = number
	return b
}

// WithSteps 设置步骤
func (b *TestCaseBuilder) WithSteps(steps testcase.Steps) *TestCaseBuilder {
	b.steps = steps
	return b
}

// WithPreconditions 设置前置条件
func (b *TestCaseBuilder) WithPreconditions(preconditions testcase.Preconditions) *TestCaseBuilder {
	b.preconditions = preconditions
	return b
}

// WithExpected 设置预期结果
func (b *TestCaseBuilder) WithExpected(expected testcase.ExpectedResult) *TestCaseBuilder {
	b.expected = expected
	return b
}

// WithCaseType 设置用例类型
func (b *TestCaseBuilder) WithCaseType(caseType testcase.CaseType) *TestCaseBuilder {
	b.caseType = caseType
	return b
}

// WithPriority 设置优先级
func (b *TestCaseBuilder) WithPriority(priority testcase.Priority) *TestCaseBuilder {
	b.priority = priority
	return b
}

// WithAIMetadata 设置 AI 元数据
func (b *TestCaseBuilder) WithAIMetadata(metadata *testcase.AiMetadata) *TestCaseBuilder {
	b.aiMetadata = metadata
	return b
}

// Build 构建测试用例
func (b *TestCaseBuilder) Build() (*testcase.TestCase, error) {
	return testcase.NewTestCase(
		b.moduleID,
		b.userID,
		b.number,
		b.title,
		b.preconditions,
		b.steps,
		b.expected,
		b.caseType,
		b.priority,
	)
}

// TestPlanBuilder 测试计划构建器
type TestPlanBuilder struct {
	projectID uuid.UUID
	userID    uuid.UUID
	name      string
	status    testplan.PlanStatus
}

// NewTestPlanBuilder 创建测试计划构建器
func NewTestPlanBuilder(projectID, userID uuid.UUID) *TestPlanBuilder {
	uid := uuid.New().String()[:8]
	return &TestPlanBuilder{
		projectID: projectID,
		userID:    userID,
		name:      "Test Plan " + uid,
		status:    testplan.StatusDraft,
	}
}

// WithName 设置名称
func (b *TestPlanBuilder) WithName(name string) *TestPlanBuilder {
	b.name = name
	return b
}

// WithStatus 设置状态
func (b *TestPlanBuilder) WithStatus(status testplan.PlanStatus) *TestPlanBuilder {
	b.status = status
	return b
}

// Build 构建测试计划
func (b *TestPlanBuilder) Build() (*testplan.TestPlan, error) {
	description := "Test plan description"
	return testplan.NewTestPlan(b.projectID, b.name, description, b.userID)
}

// TestResultBuilder 测试结果构建器
type TestResultBuilder struct {
	caseID      uuid.UUID
	planID      uuid.UUID
	executorID  uuid.UUID
	result      testplan.ResultStatus
	details     map[string]any
	executedAt  time.Time
}

// NewTestResultBuilder 创建测试结果构建器
func NewTestResultBuilder(caseID, planID, executorID uuid.UUID) *TestResultBuilder {
	return &TestResultBuilder{
		caseID:     caseID,
		planID:     planID,
		executorID: executorID,
		result:     testplan.ResultPass,
		details:    map[string]any{},
		executedAt: time.Now(),
	}
}

// WithResult 设置结果状态
func (b *TestResultBuilder) WithResult(result testplan.ResultStatus) *TestResultBuilder {
	b.result = result
	return b
}

// WithDetails 设置结果详情
func (b *TestResultBuilder) WithDetails(details map[string]any) *TestResultBuilder {
	b.details = details
	return b
}

// Build 构建测试结果
func (b *TestResultBuilder) Build() (*testplan.TestResult, error) {
	note := "Test execution note"
	return testplan.NewTestResult(b.planID, b.caseID, b.executorID, b.result, note)
}

// DocumentBuilder 文档构建器
type DocumentBuilder struct {
	projectID   uuid.UUID
	name        string
	docType     knowledge.DocumentType
	url         string
	contentText string
	metadata    map[string]any
	status      knowledge.DocumentStatus
}

// NewDocumentBuilder 创建文档构建器
func NewDocumentBuilder(projectID uuid.UUID) *DocumentBuilder {
	uid := uuid.New().String()[:8]
	return &DocumentBuilder{
		projectID:   projectID,
		name:        "Test Document " + uid,
		docType:     knowledge.TypePRD,
		url:         "https://example.com/doc.pdf",
		contentText: "Document content",
		metadata:    map[string]any{},
		status:      knowledge.StatusPending,
	}
}

// WithName 设置名称
func (b *DocumentBuilder) WithName(name string) *DocumentBuilder {
	b.name = name
	return b
}

// WithType 设置类型
func (b *DocumentBuilder) WithType(docType knowledge.DocumentType) *DocumentBuilder {
	b.docType = docType
	return b
}

// WithURL 设置 URL
func (b *DocumentBuilder) WithURL(url string) *DocumentBuilder {
	b.url = url
	return b
}

// WithContentText 设置文本内容
func (b *DocumentBuilder) WithContentText(text string) *DocumentBuilder {
	b.contentText = text
	return b
}

// WithStatus 设置状态
func (b *DocumentBuilder) WithStatus(status knowledge.DocumentStatus) *DocumentBuilder {
	b.status = status
	return b
}

// Build 构建文档
func (b *DocumentBuilder) Build() (*knowledge.Document, error) {
	userID := uuid.New() // 测试时使用随机用户ID
	return knowledge.NewDocument(b.projectID, b.name, b.docType, b.url, userID)
}

// DocumentChunkBuilder 文档块构建器
type DocumentChunkBuilder struct {
	documentID uuid.UUID
	projectID  uuid.UUID
	chunkIndex int
	content    string
	metadata   map[string]any
}

// NewDocumentChunkBuilder 创建文档块构建器
func NewDocumentChunkBuilder(documentID uuid.UUID, projectID uuid.UUID, chunkIndex int) *DocumentChunkBuilder {
	return &DocumentChunkBuilder{
		documentID: documentID,
		projectID:  projectID,
		chunkIndex: chunkIndex,
		content:    "Chunk content " + uuid.New().String()[:8],
		metadata:   map[string]any{},
	}
}

// WithContent 设置内容
func (b *DocumentChunkBuilder) WithContent(content string) *DocumentChunkBuilder {
	b.content = content
	return b
}

// WithMetadata 设置元数据
func (b *DocumentChunkBuilder) WithMetadata(metadata map[string]any) *DocumentChunkBuilder {
	b.metadata = metadata
	return b
}

// Build 构建文档块
func (b *DocumentChunkBuilder) Build() (*knowledge.DocumentChunk, error) {
	return knowledge.NewDocumentChunk(b.documentID, b.projectID, b.chunkIndex, b.content)
}

// GenerationTaskBuilder 生成任务构建器
type GenerationTaskBuilder struct {
	projectID uuid.UUID
	userID    uuid.UUID
	prompt    string
	status    generation.TaskStatus
}

// NewGenerationTaskBuilder 创建生成任务构建器
func NewGenerationTaskBuilder(projectID, userID uuid.UUID) *GenerationTaskBuilder {
	return &GenerationTaskBuilder{
		projectID: projectID,
		userID:    userID,
		prompt:    "Generate test cases for user login functionality",
		status:    generation.TaskPending,
	}
}

// WithPrompt 设置提示词
func (b *GenerationTaskBuilder) WithPrompt(prompt string) *GenerationTaskBuilder {
	b.prompt = prompt
	return b
}

// WithStatus 设置状态
func (b *GenerationTaskBuilder) WithStatus(status generation.TaskStatus) *GenerationTaskBuilder {
	b.status = status
	return b
}

// Build 构建生成任务
func (b *GenerationTaskBuilder) Build() *generation.GenerationTask {
	moduleID := uuid.New() // 测试时使用随机模块ID
	userID := uuid.New()   // 测试时使用随机用户ID
	task, _ := generation.NewGenerationTask(b.projectID, moduleID, b.prompt, userID)
	return task
}

// CaseDraftBuilder 用例草稿构建器
type CaseDraftBuilder struct {
	taskID        uuid.UUID
	moduleID      *uuid.UUID
	title         string
	preconditions testcase.Preconditions
	steps         testcase.Steps
	expected      testcase.ExpectedResult
	caseType      testcase.CaseType
	priority      testcase.Priority
	aiMetadata    *testcase.AiMetadata
	status        generation.DraftStatus
}

// NewCaseDraftBuilder 创建用例草稿构建器
func NewCaseDraftBuilder(taskID uuid.UUID) *CaseDraftBuilder {
	uid := uuid.New().String()[:8]
	return &CaseDraftBuilder{
		taskID:   taskID,
		title:    "Generated Case " + uid,
		preconditions: testcase.Preconditions{
			"用户已登录",
		},
		steps: testcase.Steps{
			"步骤1：输入用户名",
			"步骤2：输入密码",
			"步骤3：点击登录",
		},
		expected: testcase.ExpectedResult{
			"status": "success",
		},
		caseType:   testcase.CaseTypeFunctionality,
		priority:   testcase.PriorityP2,
		aiMetadata: nil,
		status:     generation.DraftPending,
	}
}

// WithModuleID 设置模块 ID
func (b *CaseDraftBuilder) WithModuleID(moduleID uuid.UUID) *CaseDraftBuilder {
	b.moduleID = &moduleID
	return b
}

// WithTitle 设置标题
func (b *CaseDraftBuilder) WithTitle(title string) *CaseDraftBuilder {
	b.title = title
	return b
}

// WithStatus 设置状态
func (b *CaseDraftBuilder) WithStatus(status generation.DraftStatus) *CaseDraftBuilder {
	b.status = status
	return b
}

// WithAIMetadata 设置 AI 元数据
func (b *CaseDraftBuilder) WithAIMetadata(metadata *testcase.AiMetadata) *CaseDraftBuilder {
	b.aiMetadata = metadata
	return b
}

// Build 构建用例草稿
func (b *CaseDraftBuilder) Build() (*generation.GeneratedCaseDraft, error) {
	return generation.NewGeneratedCaseDraft(
		b.taskID,
		b.title,
		b.preconditions,
		b.steps,
		b.expected,
		b.caseType,
		b.priority,
	)
}
