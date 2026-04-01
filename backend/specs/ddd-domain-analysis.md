# 智能测试管理平台 - DDD 领域驱动设计分析

**版本**：1.0
**日期**：2026-04-01
**作者**：DDD 架构师

---

## 1. 核心业务流程梳理（Event Storming 视角）

### 关键领域事件

| 阶段 | 事件 | 触发者 |
|---|---|---|
| **项目初始化** | `ProjectCreated` | 超级管理员 |
| | `ModuleCreated` | 超级管理员 |
| | `ProjectConfigured` | 超级管理员 |
| **知识库构建** | `DocumentUploaded` | 超级管理员 |
| | `DocumentParsed` | 系统（异步） |
| | `ChunkEmbedded` | 系统（异步） |
| **AI 用例生成** | `GenerationTaskCreated` | 超级管理员 |
| | `RAGRetrieved` | 系统 |
| | `DraftGenerated` | 系统（LLM 调用） |
| | `DraftConfirmed` → `TestCaseCreated` | 测试工程师 |
| | `DraftRejected` | 测试工程师 |
| **测试执行** | `TestPlanCreated` | 超级管理员 |
| | `CaseAssignedToPlan` | 超级管理员 |
| | `ResultRecorded` | 测试工程师 |
| | `PlanCompleted` | 系统 |

### 核心命令（Command）

```
CreateProject → Project
UploadDocument → Document
CreateGenerationTask → GenerationTask
ConfirmDraft → TestCase
CreateTestPlan → TestPlan
RecordResult → TestResult
```

---

## 2. 子域识别

### 子域分类表

| 子域 | 类型 | 重要性 | 业务价值 | 竞争差异化 |
|---|---|---|---|---|
| **AI 用例生成** | 核心域 | ⭐⭐⭐⭐⭐ | 需求→用例自动化，核心竞争力 | 高 |
| **知识库构建** | 核心域 | ⭐⭐⭐⭐⭐ | RAG 检索基础，AI 能力支撑 | 高 |
| **测试用例管理** | 核心域 | ⭐⭐⭐⭐ | 用例全生命周期管理 | 中 |
| **测试计划执行** | 支撑域 | ⭐⭐⭐ | 执行流程编排 | 低 |
| **项目管理** | 通用域 | ⭐⭐ | 多项目隔离，行业标准 | 低 |
| **用户权限** | 通用域 | ⭐⭐ | 认证授权，可替换 | 低 |

### 划分依据

**核心域（Core Domain）**：
- **AI 用例生成**：业务核心，直接决定产品竞争力
- **知识库构建**：AI 生成的前置条件，RAG 技术壁垒
- **测试用例管理**：测试资产管理的核心价值

**支撑域（Supporting Domain）**：
- **测试计划执行**：支撑核心域的执行环节，无差异化

**通用域（Generic Domain）**：
- **项目管理**：通用能力，可购买或使用开源方案
- **用户权限**：标准功能，可集成 OAuth/LDAP

---

## 3. 限界上下文（Bounded Context）定义

### 上下文边界图

```
┌─────────────────────────────────────────────────────────────────┐
│                    智能测试管理平台                               │
├─────────────────┬─────────────────┬─────────────────────────────┤
│  Identity Context │  Project Context  │    Knowledge Context       │
│  (用户权限)        │  (项目管理)        │    (知识库)                 │
├─────────────────┼─────────────────┼─────────────────────────────┤
│  TestCase Context │  TestPlan Context │    Generation Context      │
│  (用例管理)        │  (计划执行)        │    (AI 生成)                │
└─────────────────┴─────────────────┴─────────────────────────────┘
```

### 3.1 Identity Context（身份认证上下文）

| 项目 | 描述 |
|---|---|
| **职责** | 用户注册、登录、角色权限管理 |
| **聚合根** | `User` |
| **实体** | `User` |
| **值对象** | `UserRole`, `Email`, `Password` |
| **领域服务** | `AuthService`, `PermissionService` |
| **边界内表** | `users` |

**不变量**：
- 用户邮箱全局唯一
- 密码必须加密存储

---

### 3.2 Project Context（项目管理上下文）

| 项目 | 描述 |
|---|---|
| **职责** | 项目 CRUD、模块管理、项目配置 |
| **聚合根** | `Project` |
| **实体** | `Project`, `Module`, `ProjectConfig` |
| **值对象** | `ProjectId`, `ModuleId`, `ConfigKey` |
| **领域服务** | `ModuleService`, `ConfigService` |
| **边界内表** | `project`, `module`, `project_config` |

**聚合设计**：
```
Project (Aggregate Root)
├── Module[] (Entity, 项目内唯一)
└── ProjectConfig[] (Entity, KV 配置)
```

**不变量**：
- 模块名称在项目内唯一
- 配置键在项目内唯一

---

### 3.3 Knowledge Context（知识库上下文）

| 项目 | 描述 |
|---|---|
| **职责** | 文档上传、解析、分块、向量化 |
| **聚合根** | `Document` |
| **实体** | `Document`, `DocumentChunk` |
| **值对象** | `DocumentType`, `ChunkMetadata`, `EmbeddingVector` |
| **领域服务** | `DocumentParser`, `ChunkService`, `EmbeddingService` |
| **边界内表** | `document`, `document_chunk` |
| **外部依赖** | Milvus（向量存储）, OSS（文件存储） |

**聚合设计**：
```
Document (Aggregate Root)
└── DocumentChunk[] (Entity, 生命周期跟随 Document)
```

**不变量**：
- 删除文档时级联删除所有分块
- 分块索引在文档内有序

---

### 3.4 TestCase Context（测试用例上下文）

| 项目 | 描述 |
|---|---|
| **职责** | 用例 CRUD、版本管理、需求追溯 |
| **聚合根** | `TestCase` |
| **实体** | `TestCase` |
| **值对象** | `CaseNumber`, `CaseStatus`, `CaseType`, `Priority`, `Preconditions`, `Steps`, `ExpectedResult`, `AiMetadata` |
| **领域服务** | `CaseNumberGenerator`, `TraceabilityService` |
| **边界内表** | `test_case` |

**不变量**：
- 用例编号全局唯一（业务标识）
- AI 元数据记录来源追溯

---

### 3.5 TestPlan Context（测试计划上下文）

| 项目 | 描述 |
|---|---|
| **职责** | 计划编排、用例关联、执行结果记录 |
| **聚合根** | `TestPlan` |
| **实体** | `TestPlan`, `TestResult` |
| **值对象** | `PlanStatus`, `ResultStatus`, `ResultDetails` |
| **领域服务** | `PlanExecutionService`, `StatisticsService` |
| **边界内表** | `test_plan`, `test_result` |

**聚合设计**：
```
TestPlan (Aggregate Root)
└── TestResult[] (Entity, 执行记录)
```

**不变量**：
- 计划删除时级联删除执行结果
- 执行结果关联的用例不应随计划删除

---

### 3.6 Generation Context（AI 生成上下文）

| 项目 | 描述 |
|---|---|
| **职责** | 生成任务管理、RAG 检索、用例草稿生成与确认 |
| **聚合根** | `GenerationTask` |
| **实体** | `GenerationTask`, `GeneratedCaseDraft` |
| **值对象** | `TaskStatus`, `DraftStatus`, `Prompt`, `Feedback` |
| **领域服务** | `RAGService`, `LLMService`, `DraftConfirmationService` |
| **边界内表** | `generation_task`, `generated_case_draft` |
| **外部依赖** | LLM API（DeepSeek）, Knowledge Context |

**聚合设计**：
```
GenerationTask (Aggregate Root)
└── GeneratedCaseDraft[] (Entity, 草稿列表)
```

**不变量**：
- 草稿确认后转为正式用例
- 草稿拒绝时记录反馈

---

## 4. 上下文映射（Context Mapping）

### 上下文交互关系图

```
                    ┌──────────────────┐
                    │ Identity Context │
                    └────────┬─────────┘
                             │ Shared Kernel (User VO)
                             ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│ Knowledge Context│◄─│  Project Context │─►│ Generation Context│
└────────┬─────────┘  └────────┬─────────┘  └────────┬─────────┘
         │                     │                     │
         │ OHS (Document DTO)  │ OHS (Project VO)    │ OHS (Draft DTO)
         ▼                     ▼                     ▼
┌─────────────────────────────────────────────────────────────────┐
│                      TestCase Context                           │
└─────────────────────────────┬───────────────────────────────────┘
                              │ ACL (Case Reference)
                              ▼
                    ┌──────────────────┐
                    │ TestPlan Context │
                    └──────────────────┘
```

### 上下文映射表

| 上游上下文 | 下游上下文 | 关系模式 | 交互方式 | 说明 |
|---|---|---|---|---|
| Identity | All | **Shared Kernel** | 共享 `User` 值对象 | 用户基本信息作为共享概念 |
| Project | Knowledge | **OHS/PL** | 同步调用 | 知识库按项目隔离 |
| Project | TestCase | **OHS/PL** | 同步调用 | 用例归属于项目模块 |
| Project | Generation | **OHS/PL** | 同步调用 | 生成任务按项目发起 |
| Project | TestPlan | **OHS/PL** | 同步调用 | 计划归属于项目 |
| Knowledge | Generation | **OHS/PL** | gRPC 同步调用 | RAG 检索文档块 |
| Generation | TestCase | **U/D** | 领域事件 | 草稿确认后创建用例 |
| TestCase | TestPlan | **ACL** | 用例引用 | 计划引用用例，非聚合内部 |

### 关系模式说明

| 模式 | 符号 | 适用场景 |
|---|---|---|
| **Shared Kernel** | SK | 共享核心概念（如 User） |
| **Open Host Service** | OHS | 上游提供标准 API |
| **Published Language** | PL | 使用共享的领域语言 |
| **Upstream/Downstream** | U/D | 上游发布事件，下游订阅 |
| **Anti-Corruption Layer** | ACL | 下游隔离上游模型变化 |

---

## 5. 聚合与事务边界

### 聚合根一览

| 上下文 | 聚合根 | 边界内实体 | 事务一致性范围 |
|---|---|---|---|
| Identity | `User` | - | 单用户操作 |
| Project | `Project` | `Module`, `ProjectConfig` | 项目+模块+配置 |
| Knowledge | `Document` | `DocumentChunk` | 文档+分块 |
| TestCase | `TestCase` | - | 单用例操作 |
| TestPlan | `TestPlan` | `TestResult` | 计划+执行结果 |
| Generation | `GenerationTask` | `GeneratedCaseDraft` | 任务+草稿 |

### 跨聚合一致性策略

| 场景 | 策略 | 实现方式 |
|---|---|---|
| 草稿确认 → 创建用例 | **最终一致性** | 领域事件 + 消息队列 |
| 用例删除 → 计划引用 | 保留历史 | 标记失效 |
| 文档更新 → 向量重建 | **最终一致性** | 异步任务处理 |

---

## 6. 划分逻辑依据

### 解决的耦合问题

| 问题 | DDD 解决方案 |
|---|---|
| **AI 生成与用例管理耦合** | 拆分为 `Generation Context` 和 `TestCase Context`，通过领域事件解耦 |
| **知识库与业务逻辑混合** | 独立 `Knowledge Context`，通过 OHS 提供检索服务 |
| **计划执行与用例定义混淆** | 分离 `TestPlan Context` 和 `TestCase Context`，用例是稳定资产，计划是执行流程 |
| **项目配置散落各处** | 集中在 `Project Context` 的 `ProjectConfig` 实体 |

### 架构优势

1. **独立演进**：AI 生成逻辑可独立迭代，不影响用例管理
2. **技术异构**：Knowledge Context 可使用专门的向量数据库
3. **水平扩展**：Generation Context 可独立扩容应对 LLM 调用
4. **团队协作**：不同上下文可由不同团队负责

---

## 7. 目录结构建议

```
internal/
├── domain/                    # 领域模型（纯 Go，无外部依赖）
│   ├── identity/
│   │   ├── user.go
│   │   ├── user_role.go
│   │   └── repository.go
│   ├── project/
│   │   ├── project.go         # 聚合根
│   │   ├── module.go
│   │   ├── project_config.go
│   │   └── repository.go
│   ├── knowledge/
│   │   ├── document.go        # 聚合根
│   │   ├── document_chunk.go
│   │   └── repository.go
│   ├── testcase/
│   │   ├── test_case.go       # 聚合根
│   │   ├── case_number.go
│   │   └── repository.go
│   ├── testplan/
│   │   ├── test_plan.go       # 聚合根
│   │   ├── test_result.go
│   │   └── repository.go
│   └── generation/
│       ├── generation_task.go # 聚合根
│       ├── case_draft.go
│       └── repository.go
├── service/                   # 应用服务（编排领域对象）
│   ├── identity/
│   ├── project/
│   ├── knowledge/
│   ├── testcase/
│   ├── testplan/
│   └── generation/
└── transport/                 # 传输层（HTTP/gRPC）
    └── http/
```

---

## 8. Repository 接口定义

> **设计原则**: Repository 接口定义在 `domain` 层，实现在 `repository` 层，遵循依赖倒置原则。

### 8.1 Identity Context

```go
// internal/domain/identity/repository.go
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByUsername(ctx context.Context, username string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### 8.2 Project Context

```go
// internal/domain/project/repository.go
type ProjectRepository interface {
    Save(ctx context.Context, project *Project) error
    FindByID(ctx context.Context, id uuid.UUID) (*Project, error)
    FindAll(ctx context.Context) ([]*Project, error)
    Update(ctx context.Context, project *Project) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type ModuleRepository interface {
    Save(ctx context.Context, module *Module) error
    FindByID(ctx context.Context, id uuid.UUID) (*Module, error)
    FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*Module, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectConfigRepository interface {
    Save(ctx context.Context, config *ProjectConfig) error
    FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*ProjectConfig, error)
    FindByKey(ctx context.Context, projectID uuid.UUID, key string) (*ProjectConfig, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### 8.3 Knowledge Context

```go
// internal/domain/knowledge/repository.go
type DocumentRepository interface {
    Save(ctx context.Context, doc *Document) error
    FindByID(ctx context.Context, id uuid.UUID) (*Document, error)
    FindByProjectID(ctx context.Context, projectID uuid.UUID, opts QueryOptions) ([]*Document, error)
    Update(ctx context.Context, doc *Document) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type DocumentChunkRepository interface {
    SaveBatch(ctx context.Context, chunks []*DocumentChunk) error
    FindByDocumentID(ctx context.Context, documentID uuid.UUID) ([]*DocumentChunk, error)
    DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error
}

// 向量检索接口（可由 Milvus 实现）
type VectorRepository interface {
    Upsert(ctx context.Context, chunks []*DocumentChunk) error
    Search(ctx context.Context, queryVector []float32, topK int, filter map[string]any) ([]*DocumentChunk, error)
    DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error
}
```

### 8.4 TestCase Context

```go
// internal/domain/testcase/repository.go
type TestCaseRepository interface {
    Save(ctx context.Context, tc *TestCase) error
    FindByID(ctx context.Context, id uuid.UUID) (*TestCase, error)
    FindByNumber(ctx context.Context, number string) (*TestCase, error)
    FindByModuleID(ctx context.Context, moduleID uuid.UUID, opts QueryOptions) ([]*TestCase, error)
    FindByStatus(ctx context.Context, status CaseStatus, opts QueryOptions) ([]*TestCase, error)
    Update(ctx context.Context, tc *TestCase) error
    Delete(ctx context.Context, id uuid.UUID) error
    CountByModuleID(ctx context.Context, moduleID uuid.UUID) (int64, error)
}

type QueryOptions struct {
    Offset int
    Limit  int
    OrderBy string
}
```

### 8.5 TestPlan Context

```go
// internal/domain/testplan/repository.go
type TestPlanRepository interface {
    Save(ctx context.Context, plan *TestPlan) error
    FindByID(ctx context.Context, id uuid.UUID) (*TestPlan, error)
    FindByProjectID(ctx context.Context, projectID uuid.UUID, opts QueryOptions) ([]*TestPlan, error)
    Update(ctx context.Context, plan *TestPlan) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type TestResultRepository interface {
    Save(ctx context.Context, result *TestResult) error
    FindByID(ctx context.Context, id uuid.UUID) (*TestResult, error)
    FindByPlanID(ctx context.Context, planID uuid.UUID, opts QueryOptions) ([]*TestResult, error)
    FindByCaseID(ctx context.Context, caseID uuid.UUID, opts QueryOptions) ([]*TestResult, error)
    FindByExecutorID(ctx context.Context, executorID uuid.UUID, opts QueryOptions) ([]*TestResult, error)
    DeleteByPlanID(ctx context.Context, planID uuid.UUID) error
    CountByPlanID(ctx context.Context, planID uuid.UUID) (int64, error)
}
```

### 8.6 Generation Context

```go
// internal/domain/generation/repository.go
type GenerationTaskRepository interface {
    Save(ctx context.Context, task *GenerationTask) error
    FindByID(ctx context.Context, id uuid.UUID) (*GenerationTask, error)
    FindByProjectID(ctx context.Context, projectID uuid.UUID, opts QueryOptions) ([]*GenerationTask, error)
    FindByStatus(ctx context.Context, status TaskStatus, opts QueryOptions) ([]*GenerationTask, error)
    Update(ctx context.Context, task *GenerationTask) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type CaseDraftRepository interface {
    Save(ctx context.Context, draft *GeneratedCaseDraft) error
    FindByID(ctx context.Context, id uuid.UUID) (*GeneratedCaseDraft, error)
    FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]*GeneratedCaseDraft, error)
    FindByStatus(ctx context.Context, status DraftStatus, opts QueryOptions) ([]*GeneratedCaseDraft, error)
    Update(ctx context.Context, draft *GeneratedCaseDraft) error
    Delete(ctx context.Context, id uuid.UUID) error
    DeleteByTaskID(ctx context.Context, taskID uuid.UUID) error
}
```

---

## 9. 领域错误定义

> **设计原则**: 每个上下文定义独立的领域错误，通过 `ierrors` 映射为统一错误码。

### 9.1 Identity Context

```go
// internal/domain/identity/errors.go
package identity

import "errors"

var (
    // 用户不存在
    ErrUserNotFound = errors.New("user not found")
    // 邮箱已被注册
    ErrEmailDuplicate = errors.New("email already exists")
    // 用户名已被占用
    ErrUsernameDuplicate = errors.New("username already exists")
    // 密码不匹配
    ErrPasswordMismatch = errors.New("password mismatch")
    // 无效的邮箱格式
    ErrInvalidEmail = errors.New("invalid email format")
    // 权限不足
    ErrPermissionDenied = errors.New("permission denied")
)
```

### 9.2 Project Context

```go
// internal/domain/project/errors.go
package project

import "errors"

var (
    // 项目不存在
    ErrProjectNotFound = errors.New("project not found")
    // 项目名称已存在
    ErrProjectNameDuplicate = errors.New("project name already exists")
    // 模块不存在
    ErrModuleNotFound = errors.New("module not found")
    // 模块名称在项目内重复
    ErrModuleNameDuplicate = errors.New("module name duplicate in project")
    // 配置项不存在
    ErrConfigNotFound = errors.New("config not found")
    // 配置键重复
    ErrConfigKeyDuplicate = errors.New("config key duplicate in project")
)
```

### 9.3 Knowledge Context

```go
// internal/domain/knowledge/errors.go
package knowledge

import "errors"

var (
    // 文档不存在
    ErrDocumentNotFound = errors.New("document not found")
    // 文档解析失败
    ErrDocumentParseFailed = errors.New("document parse failed")
    // 不支持的文档类型
    ErrUnsupportedDocumentType = errors.New("unsupported document type")
    // 文档分块为空
    ErrEmptyChunks = errors.New("document chunks is empty")
    // 向量化失败
    ErrEmbeddingFailed = errors.New("embedding failed")
    // 向量检索失败
    ErrVectorSearchFailed = errors.New("vector search failed")
)
```

### 9.4 TestCase Context

```go
// internal/domain/testcase/errors.go
package testcase

import "errors"

var (
    // 用例不存在
    ErrCaseNotFound = errors.New("test case not found")
    // 用例编号已存在
    ErrCaseNumberDuplicate = errors.New("case number already exists")
    // 无效的用例编号格式
    ErrInvalidCaseNumber = errors.New("invalid case number format")
    // 用例步骤为空
    ErrEmptySteps = errors.New("test case steps cannot be empty")
    // 无效的优先级
    ErrInvalidPriority = errors.New("invalid priority")
    // 无效的用例类型
    ErrInvalidCaseType = errors.New("invalid case type")
)
```

### 9.5 TestPlan Context

```go
// internal/domain/testplan/errors.go
package testplan

import "errors"

var (
    // 计划不存在
    ErrPlanNotFound = errors.New("test plan not found")
    // 计划名称重复
    ErrPlanNameDuplicate = errors.New("plan name already exists")
    // 计划已归档，无法修改
    ErrPlanArchived = errors.New("plan is archived")
    // 执行结果不存在
    ErrResultNotFound = errors.New("test result not found")
    // 用例未关联到计划
    ErrCaseNotInPlan = errors.New("case not assigned to plan")
    // 重复执行
    ErrDuplicateExecution = errors.New("duplicate execution for same case")
)
```

### 9.6 Generation Context

```go
// internal/domain/generation/errors.go
package generation

import "errors"

var (
    // 任务不存在
    ErrTaskNotFound = errors.New("generation task not found")
    // 任务已处理
    ErrTaskAlreadyProcessed = errors.New("task already processed")
    // 草稿不存在
    ErrDraftNotFound = errors.New("case draft not found")
    // 草稿已确认
    ErrDraftAlreadyConfirmed = errors.New("draft already confirmed")
    // 草稿已拒绝
    ErrDraftAlreadyRejected = errors.New("draft already rejected")
    // 草稿状态无效
    ErrInvalidDraftStatus = errors.New("invalid draft status for operation")
    // LLM 调用失败
    ErrLLMCallFailed = errors.New("LLM call failed")
    // RAG 检索无结果
    ErrRAGNoResult = errors.New("RAG retrieval returned no results")
    // 并发修改冲突
    ErrConcurrentModification = errors.New("concurrent modification detected")
)
```

---

## 10. 统一错误码定义 (ierrors)

> **宪法要求**: 所有业务错误必须映射至此处的 Code，用于 HTTP 响应、日志追踪、前端国际化。

### 10.1 错误码结构

```go
// internal/ierrors/codes.go
package ierrors

// 错误码规范:
// - 1xxxx: Identity Context
// - 2xxxx: Project Context
// - 3xxxx: Knowledge Context
// - 4xxxx: TestCase Context
// - 5xxxx: TestPlan Context
// - 6xxxx: Generation Context
// - 9xxxx: 系统级错误

const (
    // ============ Identity Context (1xxxx) ============
    CodeUserNotFound      = 10001
    CodeEmailDuplicate    = 10002
    CodeUsernameDuplicate = 10003
    CodePasswordMismatch  = 10004
    CodeInvalidEmail      = 10005
    CodePermissionDenied  = 10006

    // ============ Project Context (2xxxx) ============
    CodeProjectNotFound      = 20001
    CodeProjectNameDuplicate = 20002
    CodeModuleNotFound       = 20003
    CodeModuleNameDuplicate  = 20004
    CodeConfigNotFound       = 20005
    CodeConfigKeyDuplicate   = 20006

    // ============ Knowledge Context (3xxxx) ============
    CodeDocumentNotFound        = 30001
    CodeDocumentParseFailed     = 30002
    CodeUnsupportedDocumentType = 30003
    CodeEmptyChunks             = 30004
    CodeEmbeddingFailed         = 30005
    CodeVectorSearchFailed      = 30006

    // ============ TestCase Context (4xxxx) ============
    CodeCaseNotFound       = 40001
    CodeCaseNumberDuplicate = 40002
    CodeInvalidCaseNumber  = 40003
    CodeEmptySteps         = 40004
    CodeInvalidPriority    = 40005
    CodeInvalidCaseType    = 40006

    // ============ TestPlan Context (5xxxx) ============
    CodePlanNotFound        = 50001
    CodePlanNameDuplicate   = 50002
    CodePlanArchived        = 50003
    CodeResultNotFound      = 50004
    CodeCaseNotInPlan       = 50005
    CodeDuplicateExecution  = 50006

    // ============ Generation Context (6xxxx) ============
    CodeTaskNotFound          = 60001
    CodeTaskAlreadyProcessed  = 60002
    CodeDraftNotFound         = 60003
    CodeDraftAlreadyConfirmed = 60004
    CodeDraftAlreadyRejected  = 60005
    CodeInvalidDraftStatus    = 60006
    CodeLLMCallFailed         = 60007
    CodeRAGNoResult           = 60008
    CodeConcurrentModification = 60009

    // ============ System Errors (9xxxx) ============
    CodeInternalError    = 90001
    CodeDatabaseError    = 90002
    CodeValidationError  = 90003
    CodeUnauthorized     = 90004
    CodeRateLimited      = 90005
)
```

### 10.2 错误映射

```go
// internal/ierrors/mapping.go
package ierrors

import (
    "errors"
    "github.com/yourorg/aitestos/internal/domain/identity"
    "github.com/yourorg/aitestos/internal/domain/project"
    "github.com/yourorg/aitestos/internal/domain/knowledge"
    "github.com/yourorg/aitestos/internal/domain/testcase"
    "github.com/yourorg/aitestos/internal/domain/testplan"
    "github.com/yourorg/aitestos/internal/domain/generation"
)

// MapError 将领域错误映射为统一错误码
func MapError(err error) int {
    if err == nil {
        return 0
    }

    // Identity Context
    switch {
    case errors.Is(err, identity.ErrUserNotFound):
        return CodeUserNotFound
    case errors.Is(err, identity.ErrEmailDuplicate):
        return CodeEmailDuplicate
    case errors.Is(err, identity.ErrUsernameDuplicate):
        return CodeUsernameDuplicate
    case errors.Is(err, identity.ErrPasswordMismatch):
        return CodePasswordMismatch
    case errors.Is(err, identity.ErrInvalidEmail):
        return CodeInvalidEmail
    case errors.Is(err, identity.ErrPermissionDenied):
        return CodePermissionDenied
    }

    // Project Context
    switch {
    case errors.Is(err, project.ErrProjectNotFound):
        return CodeProjectNotFound
    case errors.Is(err, project.ErrProjectNameDuplicate):
        return CodeProjectNameDuplicate
    case errors.Is(err, project.ErrModuleNotFound):
        return CodeModuleNotFound
    case errors.Is(err, project.ErrModuleNameDuplicate):
        return CodeModuleNameDuplicate
    case errors.Is(err, project.ErrConfigNotFound):
        return CodeConfigNotFound
    case errors.Is(err, project.ErrConfigKeyDuplicate):
        return CodeConfigKeyDuplicate
    }

    // Knowledge Context
    switch {
    case errors.Is(err, knowledge.ErrDocumentNotFound):
        return CodeDocumentNotFound
    case errors.Is(err, knowledge.ErrDocumentParseFailed):
        return CodeDocumentParseFailed
    case errors.Is(err, knowledge.ErrUnsupportedDocumentType):
        return CodeUnsupportedDocumentType
    case errors.Is(err, knowledge.ErrEmptyChunks):
        return CodeEmptyChunks
    case errors.Is(err, knowledge.ErrEmbeddingFailed):
        return CodeEmbeddingFailed
    case errors.Is(err, knowledge.ErrVectorSearchFailed):
        return CodeVectorSearchFailed
    }

    // TestCase Context
    switch {
    case errors.Is(err, testcase.ErrCaseNotFound):
        return CodeCaseNotFound
    case errors.Is(err, testcase.ErrCaseNumberDuplicate):
        return CodeCaseNumberDuplicate
    case errors.Is(err, testcase.ErrInvalidCaseNumber):
        return CodeInvalidCaseNumber
    case errors.Is(err, testcase.ErrEmptySteps):
        return CodeEmptySteps
    case errors.Is(err, testcase.ErrInvalidPriority):
        return CodeInvalidPriority
    case errors.Is(err, testcase.ErrInvalidCaseType):
        return CodeInvalidCaseType
    }

    // TestPlan Context
    switch {
    case errors.Is(err, testplan.ErrPlanNotFound):
        return CodePlanNotFound
    case errors.Is(err, testplan.ErrPlanNameDuplicate):
        return CodePlanNameDuplicate
    case errors.Is(err, testplan.ErrPlanArchived):
        return CodePlanArchived
    case errors.Is(err, testplan.ErrResultNotFound):
        return CodeResultNotFound
    case errors.Is(err, testplan.ErrCaseNotInPlan):
        return CodeCaseNotInPlan
    case errors.Is(err, testplan.ErrDuplicateExecution):
        return CodeDuplicateExecution
    }

    // Generation Context
    switch {
    case errors.Is(err, generation.ErrTaskNotFound):
        return CodeTaskNotFound
    case errors.Is(err, generation.ErrTaskAlreadyProcessed):
        return CodeTaskAlreadyProcessed
    case errors.Is(err, generation.ErrDraftNotFound):
        return CodeDraftNotFound
    case errors.Is(err, generation.ErrDraftAlreadyConfirmed):
        return CodeDraftAlreadyConfirmed
    case errors.Is(err, generation.ErrDraftAlreadyRejected):
        return CodeDraftAlreadyRejected
    case errors.Is(err, generation.ErrInvalidDraftStatus):
        return CodeInvalidDraftStatus
    case errors.Is(err, generation.ErrLLMCallFailed):
        return CodeLLMCallFailed
    case errors.Is(err, generation.ErrRAGNoResult):
        return CodeRAGNoResult
    case errors.Is(err, generation.ErrConcurrentModification):
        return CodeConcurrentModification
    }

    // 未知错误
    return CodeInternalError
}
```

### 10.3 HTTP 响应结构

```go
// internal/ierrors/response.go
package ierrors

import "encoding/json"

// ErrorResponse 统一错误响应结构
type ErrorResponse struct {
    Code    int    `json:"code"`    // 业务错误码
    Message string `json:"message"` // 错误信息
    TraceID string `json:"trace_id,omitempty"` // 追踪ID
}

// ToJSON 转换为 JSON
func (e *ErrorResponse) ToJSON() []byte {
    data, _ := json.Marshal(e)
    return data
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, traceID string) *ErrorResponse {
    return &ErrorResponse{
        Code:    code,
        Message: CodeToMessage(code),
        TraceID: traceID,
    }
}

// CodeToMessage 错误码转消息（可扩展为 i18n）
func CodeToMessage(code int) string {
    messages := map[int]string{
        // Identity
        CodeUserNotFound:      "用户不存在",
        CodeEmailDuplicate:    "邮箱已被注册",
        CodeUsernameDuplicate: "用户名已被占用",
        CodePasswordMismatch:  "密码错误",
        CodeInvalidEmail:      "邮箱格式无效",
        CodePermissionDenied:  "权限不足",

        // Project
        CodeProjectNotFound:      "项目不存在",
        CodeProjectNameDuplicate: "项目名称已存在",
        CodeModuleNotFound:       "模块不存在",
        CodeModuleNameDuplicate:  "模块名称重复",
        CodeConfigNotFound:       "配置项不存在",
        CodeConfigKeyDuplicate:   "配置键重复",

        // Knowledge
        CodeDocumentNotFound:        "文档不存在",
        CodeDocumentParseFailed:     "文档解析失败",
        CodeUnsupportedDocumentType: "不支持的文档类型",
        CodeEmptyChunks:             "文档分块为空",
        CodeEmbeddingFailed:         "向量化失败",
        CodeVectorSearchFailed:      "向量检索失败",

        // TestCase
        CodeCaseNotFound:       "测试用例不存在",
        CodeCaseNumberDuplicate: "用例编号已存在",
        CodeInvalidCaseNumber:  "用例编号格式无效",
        CodeEmptySteps:         "用例步骤不能为空",
        CodeInvalidPriority:    "优先级无效",
        CodeInvalidCaseType:    "用例类型无效",

        // TestPlan
        CodePlanNotFound:       "测试计划不存在",
        CodePlanNameDuplicate:  "计划名称已存在",
        CodePlanArchived:       "计划已归档",
        CodeResultNotFound:     "执行结果不存在",
        CodeCaseNotInPlan:      "用例未关联到计划",
        CodeDuplicateExecution: "重复执行",

        // Generation
        CodeTaskNotFound:           "生成任务不存在",
        CodeTaskAlreadyProcessed:   "任务已处理",
        CodeDraftNotFound:          "草稿不存在",
        CodeDraftAlreadyConfirmed:  "草稿已确认",
        CodeDraftAlreadyRejected:   "草稿已拒绝",
        CodeInvalidDraftStatus:     "草稿状态无效",
        CodeLLMCallFailed:          "LLM 调用失败",
        CodeRAGNoResult:            "RAG 检索无结果",
        CodeConcurrentModification: "并发修改冲突",

        // System
        CodeInternalError:   "系统内部错误",
        CodeDatabaseError:   "数据库错误",
        CodeValidationError: "参数校验失败",
        CodeUnauthorized:    "未授权",
        CodeRateLimited:     "请求过于频繁",
    }

    if msg, ok := messages[code]; ok {
        return msg
    }
    return "未知错误"
}
```

---

## 11. 更新后的目录结构

```
internal/
├── domain/                    # 领域模型（纯 Go，无外部依赖）
│   ├── identity/
│   │   ├── user.go
│   │   ├── user_role.go
│   │   ├── errors.go          # ✅ 领域错误
│   │   └── repository.go      # ✅ Repository 接口
│   ├── project/
│   │   ├── project.go
│   │   ├── module.go
│   │   ├── project_config.go
│   │   ├── errors.go          # ✅ 领域错误
│   │   └── repository.go      # ✅ Repository 接口
│   ├── knowledge/
│   │   ├── document.go
│   │   ├── document_chunk.go
│   │   ├── errors.go          # ✅ 领域错误
│   │   └── repository.go      # ✅ Repository 接口
│   ├── testcase/
│   │   ├── test_case.go
│   │   ├── case_number.go
│   │   ├── errors.go          # ✅ 领域错误
│   │   └── repository.go      # ✅ Repository 接口
│   ├── testplan/
│   │   ├── test_plan.go
│   │   ├── test_result.go
│   │   ├── errors.go          # ✅ 领域错误
│   │   └── repository.go      # ✅ Repository 接口
│   └── generation/
│       ├── generation_task.go
│       ├── case_draft.go
│       ├── errors.go          # ✅ 领域错误
│       └── repository.go      # ✅ Repository 接口
├── ierrors/                   # ✅ 统一错误码 (宪法要求)
│   ├── codes.go
│   ├── mapping.go
│   └── response.go
├── repository/                # ✅ 数据持久化实现
│   ├── identity/
│   ├── project/
│   ├── knowledge/
│   ├── testcase/
│   ├── testplan/
│   └── generation/
├── service/                   # 应用服务（编排领域对象）
│   ├── identity/
│   ├── project/
│   ├── knowledge/
│   ├── testcase/
│   ├── testplan/
│   └── generation/
└── transport/                 # 传输层（HTTP/gRPC）
    └── http/
```

---

## 12. 验证方式

1. 检查每个上下文是否可独立编译
2. 验证跨上下文调用仅通过定义的接口
3. 确认聚合根是唯一事务入口
4. 测试领域事件发布/订阅机制
5. **验证所有领域错误已映射到 ierrors**
6. **确认 HTTP 响应使用统一 ErrorResponse**
