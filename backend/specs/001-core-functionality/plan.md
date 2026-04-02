# 智能测试管理平台 - 核心功能技术实现方案

**版本**: 1.0
**日期**: 2026-04-02
**状态**: 已批准

---

## 1. 技术上下文总结

### 1.1 技术栈选型

| 层级 | 技术选型 | 说明 |
|------|----------|------|
| **语言** | Go 1.24+ | 标准库优先，iter 包处理集合 |
| **Web 框架** | net/http + chi | 轻量路由，标准库兼容 |
| **数据库** | PostgreSQL 16+ | 关系数据存储，JSONB 支持 |
| **缓存** | Redis 7+ | 缓存、分布式锁 |
| **向量数据库** | Milvus | RAG 检索 |
| **消息队列** | RabbitMQ | 异步任务处理 |
| **对象存储** | MinIO/S3 | 文件存储 |
| **LLM** | DeepSeek API | 用例生成、Embedding |
| **日志** | rs/zerolog | 结构化日志 |
| **指标** | Prometheus | 监控埋点 |
| **DI** | Wire | 依赖注入 |

### 1.2 架构风格

采用 **分层架构 + DDD 限界上下文**：

```
cmd/server/main.go (入口)
    ↓
internal/app (DI + 生命周期)
    ↓
internal/transport/http (HTTP 处理)
    ↓
internal/service (应用服务)
    ↓
internal/domain (领域模型)
    ↓
internal/repository (数据持久化)
```

---

## 2. 合宪性审查

### 2.1 第一条：简单性原则 (Simplicity First)

| 条款 | 审查结果 | 说明 |
|------|----------|------|
| 1.1 YAGNI | ✅ 通过 | 仅实现 spec.md 定义的功能，无预测性功能 |
| 1.2 标准库优先 | ✅ 通过 | 核心逻辑使用 net/http, encoding/json, context |
| 1.3 拒绝过度抽象 | ✅ 通过 | 接口仅在 3+ 异构实现时提取，由消费者定义 |

### 2.2 第二条：测试先行铁律 (Test-First Imperative)

| 条款 | 审查结果 | 说明 |
|------|----------|------|
| 2.1 红绿循环 | ✅ 通过 | 每个功能先写失败测试，再实现 |
| 2.2 表格驱动 | ✅ 通过 | 单元测试采用 `tt := []struct{...}` 模式 |
| 2.3 真实依赖 | ✅ 通过 | 使用 testcontainers 进行集成测试 |

### 2.3 第三条：明确性原则 (Clarity & Explicitness)

| 条款 | 审查结果 | 说明 |
|------|----------|------|
| 3.1 错误必理 | ✅ 通过 | 所有错误使用 `fmt.Errorf("context: %w", err)` 包装 |
| 3.2 零全局依赖 | ✅ 通过 | 无 init() 修改全局状态，依赖通过构造函数注入 |
| 3.3 并发安全 | ✅ 通过 | Goroutine 通过 Context 控制生命周期 |

### 2.4 第四条：分布式健壮性 (Distributed Robustness)

| 条款 | 审查结果 | 说明 |
|------|----------|------|
| 4.1 契约先行 | ✅ 通过 | API 定义在 api/openapi.yaml，修改先更新契约 |
| 4.2 失败透明 | ✅ 通过 | 所有外部调用含 Context 超时，有降级逻辑 |
| 4.3 幂等性设计 | ✅ 通过 | 写操作设计幂等键，防止重复提交 |

### 2.5 第五条：代码质量 (Code Quality)

| 条款 | 审查结果 | 说明 |
|------|----------|------|
| 5.1 静态检查 | ✅ 通过 | 通过 golangci-lint run 检查 |
| 5.2 格式化 | ✅ 通过 | 提交前运行 go fmt 和 go vet |
| 5.3 依赖整洁 | ✅ 通过 | go mod tidy 清理依赖 |

---

## 3. 项目结构细化

### 3.1 完整目录结构

```
aitestos/
├── cmd/
│   └── server/
│       └── main.go                    # 程序入口 + 优雅停机
├── internal/
│   ├── app/
│   │   ├── app.go                     # 应用容器
│   │   ├── wire.go                    # Wire 依赖注入
│   │   ├── database.go                # DB 初始化
│   │   └── shutdown.go                # 优雅停机管理
│   ├── config/
│   │   ├── config.go                  # 配置结构
│   │   └── loader.go                  # 配置加载
│   ├── domain/                        # 领域层（纯 Go）
│   │   ├── identity/
│   │   │   ├── user.go                # 聚合根
│   │   │   ├── user_role.go           # 值对象
│   │   │   ├── errors.go              # 领域错误
│   │   │   └── repository.go          # Repository 接口
│   │   ├── project/
│   │   │   ├── project.go             # 聚合根
│   │   │   ├── module.go              # 实体
│   │   │   ├── project_config.go      # 实体
│   │   │   ├── prefix.go              # 值对象
│   │   │   ├── errors.go
│   │   │   └── repository.go
│   │   ├── knowledge/
│   │   │   ├── document.go            # 聚合根
│   │   │   ├── document_chunk.go      # 实体
│   │   │   ├── document_type.go       # 值对象
│   │   │   ├── errors.go
│   │   │   └── repository.go
│   │   ├── testcase/
│   │   │   ├── test_case.go           # 聚合根
│   │   │   ├── case_number.go         # 值对象
│   │   │   ├── case_status.go         # 值对象
│   │   │   ├── ai_metadata.go         # 值对象
│   │   │   ├── errors.go
│   │   │   └── repository.go
│   │   ├── testplan/
│   │   │   ├── test_plan.go           # 聚合根
│   │   │   ├── test_result.go         # 实体
│   │   │   ├── errors.go
│   │   │   └── repository.go
│   │   └── generation/
│   │       ├── generation_task.go     # 聚合根
│   │       ├── case_draft.go          # 实体
│   │       ├── confidence.go          # 值对象
│   │       ├── errors.go
│   │       └── repository.go
│   ├── ierrors/
│   │   ├── codes.go                   # 统一错误码
│   │   ├── mapping.go                 # 错误映射
│   │   └── response.go                # HTTP 响应
│   ├── repository/                    # 数据持久化实现
│   │   ├── transaction.go             # 事务管理器
│   │   ├── identity/
│   │   │   └── user_repo.go
│   │   ├── project/
│   │   │   ├── project_repo.go
│   │   │   ├── module_repo.go
│   │   │   └── config_repo.go
│   │   ├── knowledge/
│   │   │   ├── document_repo.go
│   │   │   └── chunk_repo.go
│   │   ├── testcase/
│   │   │   └── case_repo.go
│   │   ├── testplan/
│   │   │   ├── plan_repo.go
│   │   │   └── result_repo.go
│   │   └── generation/
│   │       ├── task_repo.go
│   │       └── draft_repo.go
│   ├── service/                       # 应用服务
│   │   ├── identity/
│   │   │   └── auth_service.go
│   │   ├── project/
│   │   │   └── project_service.go
│   │   ├── knowledge/
│   │   │   └── document_service.go
│   │   ├── testcase/
│   │   │   └── case_service.go
│   │   ├── testplan/
│   │   │   └── plan_service.go
│   │   └── generation/
│   │       ├── generation_service.go
│   │       ├── rag_service.go
│   │       └── llm_service.go
│   └── transport/
│       └── http/
│           ├── server.go              # HTTP 服务器
│           ├── router.go              # 路由定义
│           ├── middleware/
│           │   ├── auth.go            # JWT 认证
│           │   ├── logging.go         # 日志
│           │   ├── recovery.go        # Panic 恢复
│           │   ├── metrics.go         # Prometheus 指标
│           │   └── trace.go           # 链路追踪
│           └── handler/
│               ├── identity.go
│               ├── project.go
│               ├── knowledge.go
│               ├── testcase.go
│               ├── testplan.go
│               └── generation.go
├── pkg/                               # 公共工具
│   ├── validator/
│   │   └── validator.go
│   └── uuidx/
│       └── uuid.go
├── api/
│   └── openapi/
│       └── openapi.yaml               # API 契约
├── configs/
│   └── config.example.yaml            # 配置示例
├── scripts/
│   └── migrate.sh                     # 数据库迁移
├── tests/
│   └── integration/                   # 集成测试
├── Makefile
├── go.mod
└── go.sum
```

### 3.2 包依赖关系

```
                    ┌─────────────┐
                    │ cmd/server  │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ internal/app│
                    └──────┬──────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
   ┌──────▼──────┐  ┌──────▼──────┐  ┌──────▼──────┐
   │  transport  │  │   config    │  │   service   │
   └──────┬──────┘  └─────────────┘  └──────┬──────┘
          │                                 │
          │         ┌───────────────────────┤
          │         │                       │
   ┌──────▼──────┐  │                ┌──────▼──────┐
   │   service   │◄─┘                │  repository │
   └──────┬──────┘                   └──────┬──────┘
          │                                 │
   ┌──────▼──────┐                   ┌──────▼──────┐
   │   domain    │◄──────────────────│  database   │
   └─────────────┘                   └─────────────┘
```

**依赖规则**：
- `transport` → `service` → `domain` ← `repository`
- `domain` 层无外部依赖，纯 Go 代码
- `repository` 实现 `domain` 定义的接口

---

## 4. 核心数据结构

### 4.1 Project Context

```go
// internal/domain/project/project.go
package project

import (
    "time"
    "github.com/google/uuid"
)

// Project 聚合根
type Project struct {
    id          uuid.UUID
    name        string
    prefix      ProjectPrefix    // 值对象
    description string
    config      map[string]any   // 简单配置
    createdAt   time.Time
    updatedAt   time.Time
}

// NewProject 创建项目（工厂函数）
func NewProject(name, prefixStr, description string) (*Project, error) {
    prefix, err := ParseProjectPrefix(prefixStr)
    if err != nil {
        return nil, err
    }

    now := time.Now()
    return &Project{
        id:          uuid.New(),
        name:        name,
        prefix:      prefix,
        description: description,
        config:      make(map[string]any),
        createdAt:   now,
        updatedAt:   now,
    }, nil
}

// 只读访问器
func (p *Project) ID() uuid.UUID         { return p.id }
func (p *Project) Name() string          { return p.name }
func (p *Project) Prefix() ProjectPrefix { return p.prefix }
func (p *Project) Description() string   { return p.description }
func (p *Project) CreatedAt() time.Time  { return p.createdAt }
func (p *Project) UpdatedAt() time.Time  { return p.updatedAt }
```

```go
// internal/domain/project/prefix.go
package project

import (
    "errors"
    "regexp"
)

// ProjectPrefix 项目前缀值对象（2-4位大写字母）
type ProjectPrefix string

var prefixRegex = regexp.MustCompile(`^[A-Z]{2,4}$`)

var (
    ErrInvalidProjectPrefix  = errors.New("invalid project prefix: must be 2-4 uppercase letters")
    ErrProjectPrefixDuplicate = errors.New("project prefix already exists")
)

func ParseProjectPrefix(s string) (ProjectPrefix, error) {
    if !prefixRegex.MatchString(s) {
        return "", ErrInvalidProjectPrefix
    }
    return ProjectPrefix(s), nil
}

func (p ProjectPrefix) String() string { return string(p) }
```

```go
// internal/domain/project/module.go
package project

import (
    "time"
    "github.com/google/uuid"
)

// Module 实体
type Module struct {
    id            uuid.UUID
    projectID     uuid.UUID
    name          string
    abbreviation  ModuleAbbreviation  // 值对象
    description   string
    createdAt     time.Time
    updatedAt     time.Time
}

// ModuleAbbreviation 模块缩写值对象
type ModuleAbbreviation string

var abbrevRegex = regexp.MustCompile(`^[A-Z]{2,4}$`)

var (
    ErrInvalidModuleAbbrev   = errors.New("invalid module abbreviation: must be 2-4 uppercase letters")
    ErrModuleAbbrevDuplicate = errors.New("module abbreviation already exists in project")
)

func ParseModuleAbbreviation(s string) (ModuleAbbreviation, error) {
    if !abbrevRegex.MatchString(s) {
        return "", ErrInvalidModuleAbbrev
    }
    return ModuleAbbreviation(s), nil
}
```

### 4.2 TestCase Context

```go
// internal/domain/testcase/test_case.go
package testcase

import (
    "time"
    "github.com/google/uuid"
)

// TestCase 聚合根
type TestCase struct {
    id            uuid.UUID
    moduleID      uuid.UUID
    userID        uuid.UUID
    number        CaseNumber          // 值对象
    title         string
    preconditions Preconditions       // 值对象
    steps         Steps               // 值对象
    expected      ExpectedResult      // 值对象
    aiMetadata    *AiMetadata         // 值对象，AI生成用例才有
    caseType      CaseType            // 值对象
    priority      Priority            // 值对象
    status        CaseStatus          // 值对象
    createdAt     time.Time
    updatedAt     time.Time
}

// Preconditions 前置条件
type Preconditions []string

// Steps 测试步骤
type Steps []string

// ExpectedResult 预期结果
type ExpectedResult map[string]any

// CaseType 用例类型
type CaseType string

const (
    CaseTypeFunctionality CaseType = "functionality"
    CaseTypePerformance   CaseType = "performance"
    CaseTypeAPI           CaseType = "api"
    CaseTypeUI            CaseType = "ui"
    CaseTypeSecurity      CaseType = "security"
)

// Priority 优先级
type Priority string

const (
    PriorityP0 Priority = "P0"
    PriorityP1 Priority = "P1"
    PriorityP2 Priority = "P2"
    PriorityP3 Priority = "P3"
)

// CaseStatus 用例状态
type CaseStatus string

const (
    StatusUnexecuted CaseStatus = "unexecuted"
    StatusPass       CaseStatus = "pass"
    StatusBlock      CaseStatus = "block"
    StatusFail       CaseStatus = "fail"
)

// NewTestCase 创建测试用例
func NewTestCase(
    moduleID, userID uuid.UUID,
    number CaseNumber,
    title string,
    preconditions Preconditions,
    steps Steps,
    expected ExpectedResult,
    caseType CaseType,
    priority Priority,
) (*TestCase, error) {
    if len(steps) == 0 {
        return nil, ErrEmptySteps
    }
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }

    now := time.Now()
    return &TestCase{
        id:            uuid.New(),
        moduleID:      moduleID,
        userID:        userID,
        number:        number,
        title:         title,
        preconditions: preconditions,
        steps:         steps,
        expected:      expected,
        caseType:      caseType,
        priority:      priority,
        status:        StatusUnexecuted,
        createdAt:     now,
        updatedAt:     now,
    }, nil
}

// UpdateStatus 更新状态
func (tc *TestCase) UpdateStatus(status CaseStatus) {
    tc.status = status
    tc.updatedAt = time.Now()
}
```

```go
// internal/domain/testcase/case_number.go
package testcase

import (
    "fmt"
    "regexp"
    "time"
)

// CaseNumber 用例编号值对象
// 格式: {项目前缀}-{模块缩写}-{日期}-{序号}
// 示例: ECO-USR-20260402-001
type CaseNumber string

var caseNumberRegex = regexp.MustCompile(`^[A-Z]{2,4}-[A-Z]{2,4}-\d{8}-\d{3}$`)

func ParseCaseNumber(s string) (CaseNumber, error) {
    if !caseNumberRegex.MatchString(s) {
        return "", ErrInvalidCaseNumber
    }
    return CaseNumber(s), nil
}

// GenerateCaseNumber 生成用例编号
func GenerateCaseNumber(projectPrefix, moduleAbbrev string, seq int) CaseNumber {
    date := time.Now().Format("20060102")
    return CaseNumber(fmt.Sprintf("%s-%s-%s-%03d",
        projectPrefix, moduleAbbrev, date, seq))
}

func (n CaseNumber) String() string { return string(n) }
```

```go
// internal/domain/testcase/ai_metadata.go
package testcase

import (
    "time"
    "github.com/google/uuid"
)

// Confidence AI 置信度
type Confidence string

const (
    ConfidenceHigh   Confidence = "high"
    ConfidenceMedium Confidence = "medium"
    ConfidenceLow    Confidence = "low"
)

// ReferencedChunk 引用的文档块
type ReferencedChunk struct {
    ChunkID         uuid.UUID `json:"chunk_id"`
    DocumentID      uuid.UUID `json:"document_id"`
    DocumentTitle   string    `json:"document_title"`
    SimilarityScore float64   `json:"similarity_score"`
}

// AiMetadata AI 元数据
type AiMetadata struct {
    GenerationTaskID  uuid.UUID         `json:"generation_task_id"`
    Confidence        Confidence        `json:"confidence"`
    ReferencedChunks  []ReferencedChunk `json:"referenced_chunks"`
    ModelVersion      string            `json:"model_version"`
    GeneratedAt       time.Time         `json:"generated_at"`
}

// CalculateConfidence 根据检索结果计算置信度
func CalculateConfidence(chunks []ReferencedChunk) Confidence {
    if len(chunks) >= 2 && chunks[0].SimilarityScore > 0.8 {
        return ConfidenceHigh
    }
    if len(chunks) >= 1 && chunks[0].SimilarityScore >= 0.5 {
        return ConfidenceMedium
    }
    return ConfidenceLow
}
```

### 4.3 Generation Context

```go
// internal/domain/generation/generation_task.go
package generation

import (
    "time"
    "github.com/google/uuid"
)

// TaskStatus 任务状态
type TaskStatus string

const (
    TaskPending    TaskStatus = "pending"
    TaskProcessing TaskStatus = "processing"
    TaskCompleted  TaskStatus = "completed"
    TaskFailed     TaskStatus = "failed"
)

// GenerationTask 聚合根
type GenerationTask struct {
    id            uuid.UUID
    projectID     uuid.UUID
    userID        uuid.UUID
    status        TaskStatus
    prompt        string
    resultSummary map[string]any
    errorMsg      string
    createdAt     time.Time
    updatedAt     time.Time
}

// NewGenerationTask 创建生成任务
func NewGenerationTask(projectID, userID uuid.UUID, prompt string) *GenerationTask {
    now := time.Now()
    return &GenerationTask{
        id:        uuid.New(),
        projectID: projectID,
        userID:    userID,
        status:    TaskPending,
        prompt:    prompt,
        createdAt: now,
        updatedAt: now,
    }
}

// StartProcessing 开始处理
func (t *GenerationTask) StartProcessing() {
    t.status = TaskProcessing
    t.updatedAt = time.Now()
}

// Complete 完成
func (t *GenerationTask) Complete(summary map[string]any) {
    t.status = TaskCompleted
    t.resultSummary = summary
    t.updatedAt = time.Now()
}

// Fail 失败
func (t *GenerationTask) Fail(errMsg string) {
    t.status = TaskFailed
    t.errorMsg = errMsg
    t.updatedAt = time.Now()
}
```

```go
// internal/domain/generation/case_draft.go
package generation

import (
    "time"
    "github.com/google/uuid"
    "github.com/liang21/aitestos/internal/domain/testcase"
)

// DraftStatus 草稿状态
type DraftStatus string

const (
    DraftPending   DraftStatus = "pending"
    DraftConfirmed DraftStatus = "confirmed"
    DraftRejected  DraftStatus = "rejected"
)

// RejectionReason 拒绝原因
type RejectionReason string

const (
    RejectionDuplicate   RejectionReason = "duplicate"
    RejectionIrrelevant  RejectionReason = "irrelevant"
    RejectionLowQuality  RejectionReason = "low_quality"
    RejectionOther       RejectionReason = "other"
)

// GeneratedCaseDraft 实体
type GeneratedCaseDraft struct {
    id            uuid.UUID
    taskID        uuid.UUID
    moduleID      *uuid.UUID
    title         string
    preconditions testcase.Preconditions
    steps         testcase.Steps
    expected      testcase.ExpectedResult
    caseType      testcase.CaseType
    priority      testcase.Priority
    aiMetadata    *testcase.AiMetadata
    status        DraftStatus
    feedback      string
    createdAt     time.Time
    updatedAt     time.Time
}

// Confirm 确认草稿
func (d *GeneratedCaseDraft) Confirm(moduleID uuid.UUID) {
    d.moduleID = &moduleID
    d.status = DraftConfirmed
    d.updatedAt = time.Now()
}

// Reject 拒绝草稿
func (d *GeneratedCaseDraft) Reject(reason RejectionReason, detail string) {
    d.status = DraftRejected
    d.feedback = string(reason) + ": " + detail
    d.updatedAt = time.Now()
}
```

---

## 5. 接口设计

### 5.1 Domain Repository 接口

```go
// internal/domain/project/repository.go
package project

import (
    "context"
    "github.com/google/uuid"
)

type ProjectRepository interface {
    Save(ctx context.Context, project *Project) error
    FindByID(ctx context.Context, id uuid.UUID) (*Project, error)
    FindByName(ctx context.Context, name string) (*Project, error)
    FindByPrefix(ctx context.Context, prefix ProjectPrefix) (*Project, error)
    FindAll(ctx context.Context, opts QueryOptions) ([]*Project, error)
    Update(ctx context.Context, project *Project) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type ModuleRepository interface {
    Save(ctx context.Context, module *Module) error
    FindByID(ctx context.Context, id uuid.UUID) (*Module, error)
    FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*Module, error)
    FindByAbbreviation(ctx context.Context, projectID uuid.UUID, abbrev ModuleAbbreviation) (*Module, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectConfigRepository interface {
    Save(ctx context.Context, config *ProjectConfig) error
    FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*ProjectConfig, error)
    FindByKey(ctx context.Context, projectID uuid.UUID, key string) (*ProjectConfig, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type QueryOptions struct {
    Offset   int
    Limit    int
    OrderBy  string
    Keywords string
}
```

```go
// internal/domain/testcase/repository.go
package testcase

import (
    "context"
    "time"
    "github.com/google/uuid"
)

type TestCaseRepository interface {
    Save(ctx context.Context, tc *TestCase) error
    FindByID(ctx context.Context, id uuid.UUID) (*TestCase, error)
    FindByNumber(ctx context.Context, number CaseNumber) (*TestCase, error)
    FindByModuleID(ctx context.Context, moduleID uuid.UUID, opts QueryOptions) ([]*TestCase, error)
    FindByProjectID(ctx context.Context, projectID uuid.UUID, opts QueryOptions) ([]*TestCase, error)
    Update(ctx context.Context, tc *TestCase) error
    Delete(ctx context.Context, id uuid.UUID) error
    CountByDate(ctx context.Context, moduleID uuid.UUID, date time.Time) (int64, error)
}
```

```go
// internal/domain/generation/repository.go
package generation

import (
    "context"
    "github.com/google/uuid"
)

type GenerationTaskRepository interface {
    Save(ctx context.Context, task *GenerationTask) error
    FindByID(ctx context.Context, id uuid.UUID) (*GenerationTask, error)
    FindByProjectID(ctx context.Context, projectID uuid.UUID, opts QueryOptions) ([]*GenerationTask, error)
    FindByStatus(ctx context.Context, status TaskStatus, opts QueryOptions) ([]*GenerationTask, error)
    Update(ctx context.Context, task *GenerationTask) error
}

type CaseDraftRepository interface {
    Save(ctx context.Context, draft *GeneratedCaseDraft) error
    FindByID(ctx context.Context, id uuid.UUID) (*GeneratedCaseDraft, error)
    FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]*GeneratedCaseDraft, error)
    Update(ctx context.Context, draft *GeneratedCaseDraft) error
}
```

### 5.2 Service 接口

```go
// internal/service/project/project_service.go
package project

import (
    "context"
    "github.com/google/uuid"
    "github.com/liang21/aitestos/internal/domain/project"
)

type ProjectService interface {
    // 项目管理
    CreateProject(ctx context.Context, req *CreateProjectRequest) (*project.Project, error)
    GetProject(ctx context.Context, id uuid.UUID) (*ProjectDetail, error)
    ListProjects(ctx context.Context, opts ListOptions) ([]*project.Project, int64, error)
    UpdateProject(ctx context.Context, id uuid.UUID, req *UpdateProjectRequest) (*project.Project, error)
    DeleteProject(ctx context.Context, id uuid.UUID) error

    // 模块管理
    CreateModule(ctx context.Context, projectID uuid.UUID, req *CreateModuleRequest) (*project.Module, error)
    ListModules(ctx context.Context, projectID uuid.UUID) ([]*project.Module, error)
    DeleteModule(ctx context.Context, id uuid.UUID) error

    // 配置管理
    SetConfig(ctx context.Context, projectID uuid.UUID, key string, value any) error
    GetConfig(ctx context.Context, projectID uuid.UUID, key string) (any, error)
}

type CreateProjectRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Prefix      string `json:"prefix" validate:"required,min=2,max=4"`
    Description string `json:"description"`
}

type ProjectDetail struct {
    *project.Project
    ModuleCount   int64 `json:"module_count"`
    CaseCount     int64 `json:"case_count"`
    DocumentCount int64 `json:"document_count"`
}
```

```go
// internal/service/generation/generation_service.go
package generation

import (
    "context"
    "github.com/google/uuid"
)

type GenerationService interface {
    // 创建生成任务
    CreateTask(ctx context.Context, req *CreateTaskRequest) (*GenerationTask, error)

    // 获取任务状态
    GetTask(ctx context.Context, taskID uuid.UUID) (*TaskDetail, error)

    // 获取草稿列表
    ListDrafts(ctx context.Context, opts DraftListOptions) ([]*CaseDraft, int64, error)

    // 确认草稿
    ConfirmDraft(ctx context.Context, draftID uuid.UUID, moduleID uuid.UUID) (*TestCase, error)

    // 拒绝草稿
    RejectDraft(ctx context.Context, draftID uuid.UUID, reason string, feedback string) error

    // 批量确认
    BatchConfirm(ctx context.Context, draftIDs []uuid.UUID, moduleID uuid.UUID) (*BatchResult, error)
}

type CreateTaskRequest struct {
    ProjectID          uuid.UUID `json:"project_id" validate:"required"`
    ModuleID           uuid.UUID `json:"module_id" validate:"required"`
    Prompt             string    `json:"prompt" validate:"required,min=20"`
    DocumentScope      string    `json:"document_scope"`       // all/prd_only/figma_only
    CaseCount          int       `json:"case_count"`           // 1-20
    SceneTypes         []string  `json:"scene_types"`          // positive/negative/boundary
    PriorityPreference string    `json:"priority_preference"`  // P0-P3
    GenerationMode     string    `json:"generation_mode"`      // normal/deep
}
```

---

## 6. 实施阶段

### Phase 1: 基础设施 (P0)

| 任务 | 文件 | 预计时间 |
|------|------|----------|
| 配置管理完善 | internal/config/* | 2h |
| 数据库连接池 | internal/app/database.go | 2h |
| Wire 依赖注入 | internal/app/wire.go | 3h |
| 优雅停机 | internal/app/shutdown.go | 2h |
| 统一错误码 | internal/ierrors/* | 2h |

### Phase 2: 领域模型 (P0)

| 任务 | 文件 | 预计时间 |
|------|------|----------|
| Identity Context | internal/domain/identity/* | 2h |
| Project Context | internal/domain/project/* | 4h |
| TestCase Context | internal/domain/testcase/* | 4h |
| TestPlan Context | internal/domain/testplan/* | 3h |
| Knowledge Context | internal/domain/knowledge/* | 3h |
| Generation Context | internal/domain/generation/* | 4h |

### Phase 3: Repository 层 (P0)

| 任务 | 文件 | 预计时间 |
|------|------|----------|
| 事务管理器 | internal/repository/transaction.go | 2h |
| Project Repository | internal/repository/project/* | 4h |
| TestCase Repository | internal/repository/testcase/* | 3h |
| TestPlan Repository | internal/repository/testplan/* | 3h |
| Knowledge Repository | internal/repository/knowledge/* | 4h |
| Generation Repository | internal/repository/generation/* | 3h |

### Phase 4: Service 层 (P1)

| 任务 | 文件 | 预计时间 |
|------|------|----------|
| Auth Service | internal/service/identity/* | 3h |
| Project Service | internal/service/project/* | 4h |
| TestCase Service | internal/service/testcase/* | 4h |
| TestPlan Service | internal/service/testplan/* | 4h |
| Document Service | internal/service/knowledge/* | 4h |
| Generation Service | internal/service/generation/* | 6h |

### Phase 5: Transport 层 (P1)

| 任务 | 文件 | 预计时间 |
|------|------|----------|
| HTTP Server & Router | internal/transport/http/* | 3h |
| Middleware | internal/transport/http/middleware/* | 4h |
| Handlers | internal/transport/http/handler/* | 6h |

### Phase 6: 集成测试 (P2)

| 任务 | 文件 | 预计时间 |
|------|------|----------|
| Repository 测试 | tests/integration/repository/* | 4h |
| Service 测试 | tests/integration/service/* | 4h |
| API 测试 | tests/integration/api/* | 4h |

---

## 7. 验证方式

1. **编译检查**: `make build` 无错误
2. **静态检查**: `make lint` 无警告
3. **单元测试**: `make test` 覆盖率 > 80%
4. **集成测试**: `make test-integration` 通过
5. **API 测试**: 使用 Postman/curl 验证 HTTP 接口
6. **性能测试**: k6 负载测试，P99 < 500ms

---

## 8. 参考文档

- [spec.md](./spec.md) - 核心功能规范
- [../ddd-domain-analysis.md](../ddd-domain-analysis.md) - DDD 领域分析
- [../go-engineering-design.md](../go-engineering-design.md) - Go 工程设计
- [../aitestos_optimized.sql](../aitestos_optimized.sql) - 数据库 Schema
- [../openapi.yaml](../openapi.yaml) - API 契约
