# 核心功能模块 - 原子化任务列表

**版本**: 1.0
**日期**: 2026-04-02
**状态**: 待执行

---

## 任务说明

- **[P]** 标记表示该任务与其他同级 [P] 任务可并行执行
- **依赖** 列表示必须先完成的前置任务
- **TDD**: 测试任务 (🔴) 必须在实现任务 (🟢) 之前完成
- **验收标准**: 每个任务的完成标准

---

## Phase 1: 基础设施 (P0)

### 1.1 配置管理

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-001 | 🔴 编写配置结构测试 | `internal/config/config_test.go` | - | 测试 YAML 解析、环境变量覆盖、默认值 |
| T-002 | 🟢 实现配置结构体 | `internal/config/config.go` | T-001 | 包含 Server/Database/Redis/LLM/Milvus/Storage/JWT/Log 配置 |
| T-003 | 🔴 编写配置加载器测试 | `internal/config/loader_test.go` | T-002 | 测试文件加载、环境变量替换、校验 |
| T-004 | 🟢 实现配置加载器 | `internal/config/loader.go` | T-003 | 使用 viper 加载配置，支持环境变量 |
| T-005 | 📝 创建配置示例文件 | `configs/config.example.yaml` | T-002 | 包含所有配置项的示例 |

### 1.2 统一错误码

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-006 | 🔴 编写错误码测试 | `internal/ierrors/codes_test.go` | - | 测试错误码常量定义完整性 |
| T-007 | 🟢 实现统一错误码 | `internal/ierrors/codes.go` | T-006 | 包含 1xxxx-6xxxx 所有错误码 |
| T-008 | 🔴 编写错误映射测试 | `internal/ierrors/mapping_test.go` | T-007 | 测试领域错误到错误码的映射 |
| T-009 | 🟢 实现错误映射 | `internal/ierrors/mapping.go` | T-008 | MapError 函数，支持 errors.Is |
| T-010 | 🔴 编写 HTTP 响应测试 | `internal/ierrors/response_test.go` | T-007 | 测试 ErrorResponse 序列化 |
| T-011 | 🟢 实现 HTTP 响应结构 | `internal/ierrors/response.go` | T-010 | ErrorResponse, NewErrorResponse, CodeToMessage |

### 1.3 数据库连接池

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-012 | 🔴 编写数据库连接测试 | `internal/app/database_test.go` | T-002 | 测试连接池配置、健康检查 |
| T-013 | 🟢 实现数据库连接池 | `internal/app/database.go` | T-012 | NewDB 函数，连接池参数配置 |

### 1.4 优雅停机

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-014 | 🔴 编写停机管理器测试 | `internal/app/shutdown_test.go` | - | 测试组件注册顺序、逆序关闭 |
| T-015 | 🟢 实现优雅停机管理器 | `internal/app/shutdown.go` | T-014 | ShutdownManager, Closer 接口 |

### 1.5 Wire 依赖注入

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-016 | 🟢 定义 Wire Provider Set | `internal/app/wire.go` | T-013, T-015 | InfrastructureSet, RepositorySet, ServiceSet, HTTPSet |
| T-017 | 🟢 实现 App 容器 | `internal/app/app.go` | T-016 | App 结构体，Run 方法含优雅停机 |

### 1.6 程序入口

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-018 | 🟢 实现程序入口 | `cmd/server/main.go` | T-017 | 加载配置，初始化 App，捕获 SIGTERM |

### Phase 1 完成标准
- [ ] `make build` 编译通过
- [ ] `make test` 单元测试通过
- [ ] `make lint` 无警告

---

## Phase 2: 领域模型 (P0)

### 2.1 Identity Context

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-019 | [P] 🔴 编写 UserRole 值对象测试 | `internal/domain/identity/user_role_test.go` | - | 测试角色解析、IsAdmin 方法 |
| T-020 | [P] 🟢 实现 UserRole 值对象 | `internal/domain/identity/user_role.go` | T-019 | RoleSuperAdmin, RoleAdmin, RoleNormal |
| T-021 | [P] 🔴 编写 User 聚合根测试 | `internal/domain/identity/user_test.go` | T-020 | 测试 NewUser, VerifyPassword, ChangePassword |
| T-022 | [P] 🟢 实现 User 聚合根 | `internal/domain/identity/user.go` | T-021 | 私有字段，只读访问器，领域行为 |
| T-023 | [P] 🟢 定义 Identity 领域错误 | `internal/domain/identity/errors.go` | - | ErrUserNotFound, ErrEmailDuplicate 等 |
| T-024 | [P] 🟢 定义 UserRepository 接口 | `internal/domain/identity/repository.go` | T-022 | Save, FindByID, FindByEmail, Update, Delete |

### 2.2 Project Context

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-025 | [P] 🔴 编写 ProjectPrefix 值对象测试 | `internal/domain/project/prefix_test.go` | - | 测试 2-4 位大写字母校验 |
| T-026 | [P] 🟢 实现 ProjectPrefix 值对象 | `internal/domain/project/prefix.go` | T-025 | ParseProjectPrefix, String 方法 |
| T-027 | [P] 🔴 编写 ModuleAbbreviation 值对象测试 | `internal/domain/project/abbreviation_test.go` | - | 测试 2-4 位大写字母校验 |
| T-028 | [P] 🟢 实现 ModuleAbbreviation 值对象 | `internal/domain/project/abbreviation.go` | T-027 | ParseModuleAbbreviation |
| T-029 | 🔴 编写 Project 聚合根测试 | `internal/domain/project/project_test.go` | T-026 | 测试 NewProject, 访问器 |
| T-030 | 🟢 实现 Project 聚合根 | `internal/domain/project/project.go` | T-029 | 私有字段，工厂函数 |
| T-031 | 🔴 编写 Module 实体测试 | `internal/domain/project/module_test.go` | T-028, T-030 | 测试 NewModule |
| T-032 | 🟢 实现 Module 实体 | `internal/domain/project/module.go` | T-031 | 关联 ProjectID |
| T-033 | 🔴 编写 ProjectConfig 实体测试 | `internal/domain/project/project_config_test.go` | T-030 | 测试配置 CRUD |
| T-034 | 🟢 实现 ProjectConfig 实体 | `internal/domain/project/project_config.go` | T-033 | Key-Value 配置 |
| T-035 | 🟢 定义 Project 领域错误 | `internal/domain/project/errors.go` | - | ErrProjectNotFound, ErrPrefixDuplicate 等 |
| T-036 | 🟢 定义 Project Repository 接口 | `internal/domain/project/repository.go` | T-030, T-032, T-034 | ProjectRepository, ModuleRepository, ProjectConfigRepository |

### 2.3 TestCase Context

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-037 | [P] 🔴 编写 CaseNumber 值对象测试 | `internal/domain/testcase/case_number_test.go` | - | 测试格式校验、生成逻辑 |
| T-038 | [P] 🟢 实现 CaseNumber 值对象 | `internal/domain/testcase/case_number.go` | T-037 | ParseCaseNumber, GenerateCaseNumber |
| T-039 | [P] 🔴 编写 CaseStatus 值对象测试 | `internal/domain/testcase/case_status_test.go` | - | 测试状态枚举 |
| T-040 | [P] 🟢 实现 CaseStatus 值对象 | `internal/domain/testcase/case_status.go` | T-039 | StatusUnexecuted, StatusPass, StatusBlock, StatusFail |
| T-041 | [P] 🔴 编写 AiMetadata 值对象测试 | `internal/domain/testcase/ai_metadata_test.go` | - | 测试置信度计算逻辑 |
| T-042 | [P] 🟢 实现 AiMetadata 值对象 | `internal/domain/testcase/ai_metadata.go` | T-041 | Confidence, ReferencedChunk, CalculateConfidence |
| T-043 | 🔴 编写 TestCase 聚合根测试 | `internal/domain/testcase/test_case_test.go` | T-038, T-040, T-042 | 测试 NewTestCase, UpdateStatus |
| T-044 | 🟢 实现 TestCase 聚合根 | `internal/domain/testcase/test_case.go` | T-043 | 私有字段，Preconditions, Steps, ExpectedResult |
| T-045 | 🟢 定义 TestCase 领域错误 | `internal/domain/testcase/errors.go` | - | ErrCaseNotFound, ErrEmptySteps 等 |
| T-046 | 🟢 定义 TestCase Repository 接口 | `internal/domain/testcase/repository.go` | T-044 | Save, FindByID, FindByNumber, CountByDate |

### 2.4 TestPlan Context

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-047 | [P] 🔴 编写 PlanStatus 值对象测试 | `internal/domain/testplan/plan_status_test.go` | - | 测试状态枚举和转换 |
| T-048 | [P] 🟢 实现 PlanStatus 值对象 | `internal/domain/testplan/plan_status.go` | T-047 | Draft, Active, Completed, Archived |
| T-049 | 🔴 编写 TestPlan 聚合根测试 | `internal/domain/testplan/test_plan_test.go` | T-048 | 测试 NewTestPlan, AddCase, UpdateStatus |
| T-050 | 🟢 实现 TestPlan 聚合根 | `internal/domain/testplan/test_plan.go` | T-049 | 关联 ProjectID, CaseIDs |
| T-051 | 🔴 编写 TestResult 实体测试 | `internal/domain/testplan/test_result_test.go` | T-050 | 测试执行结果录入 |
| T-052 | 🟢 实现 TestResult 实体 | `internal/domain/testplan/test_result.go` | T-051 | Pass/Fail/Block/Skip 状态 |
| T-053 | 🟢 定义 TestPlan 领域错误 | `internal/domain/testplan/errors.go` | - | ErrPlanNotFound, ErrPlanArchived 等 |
| T-054 | 🟢 定义 TestPlan Repository 接口 | `internal/domain/testplan/repository.go` | T-050, T-052 | TestPlanRepository, TestResultRepository |

### 2.5 Knowledge Context

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-055 | [P] 🔴 编写 DocumentType 值对象测试 | `internal/domain/knowledge/document_type_test.go` | - | 测试 PRD/Figma/APISpec 类型 |
| T-056 | [P] 🟢 实现 DocumentType 值对象 | `internal/domain/knowledge/document_type.go` | T-055 | TypePRD, TypeFigma, TypeAPISpec |
| T-057 | [P] 🔴 编写 DocumentStatus 值对象测试 | `internal/domain/knowledge/document_status_test.go` | - | 测试状态流转 |
| T-058 | [P] 🟢 实现 DocumentStatus 值对象 | `internal/domain/knowledge/document_status.go` | T-057 | Pending, Processing, Completed, Failed |
| T-059 | 🔴 编写 Document 聚合根测试 | `internal/domain/knowledge/document_test.go` | T-056, T-058 | 测试 NewDocument, UpdateStatus |
| T-060 | 🟢 实现 Document 聚合根 | `internal/domain/knowledge/document.go` | T-059 | 关联 ProjectID, URL, Metadata |
| T-061 | 🔴 编写 DocumentChunk 实体测试 | `internal/domain/knowledge/document_chunk_test.go` | T-060 | 测试分块创建 |
| T-062 | 🟢 实现 DocumentChunk 实体 | `internal/domain/knowledge/document_chunk.go` | T-061 | ChunkIndex, Content, Metadata |
| T-063 | 🟢 定义 Knowledge 领域错误 | `internal/domain/knowledge/errors.go` | - | ErrDocumentNotFound, ErrKnowledgeBaseEmpty 等 |
| T-064 | 🟢 定义 Knowledge Repository 接口 | `internal/domain/knowledge/repository.go` | T-060, T-062 | DocumentRepository, DocumentChunkRepository, VectorRepository |

### 2.6 Generation Context

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-065 | [P] 🔴 编写 TaskStatus 值对象测试 | `internal/domain/generation/task_status_test.go` | - | 测试状态枚举 |
| T-066 | [P] 🟢 实现 TaskStatus 值对象 | `internal/domain/generation/task_status.go` | T-065 | Pending, Processing, Completed, Failed |
| T-067 | [P] 🔴 编写 DraftStatus 值对象测试 | `internal/domain/generation/draft_status_test.go` | - | 测试状态枚举 |
| T-068 | [P] 🟢 实现 DraftStatus 值对象 | `internal/domain/generation/draft_status.go` | T-067 | Pending, Confirmed, Rejected |
| T-069 | 🔴 编写 GenerationTask 聚合根测试 | `internal/domain/generation/generation_task_test.go` | T-066 | 测试 NewGenerationTask, StartProcessing, Complete, Fail |
| T-070 | 🟢 实现 GenerationTask 聚合根 | `internal/domain/generation/generation_task.go` | T-069 | 关联 ProjectID, Prompt, ResultSummary |
| T-071 | 🔴 编写 GeneratedCaseDraft 实体测试 | `internal/domain/generation/case_draft_test.go` | T-068, T-042 | 测试 Confirm, Reject 方法 |
| T-072 | 🟢 实现 GeneratedCaseDraft 实体 | `internal/domain/generation/case_draft.go` | T-071 | 关联 TaskID, 引用 testcase 值对象 |
| T-073 | 🟢 定义 Generation 领域错误 | `internal/domain/generation/errors.go` | - | ErrTaskNotFound, ErrDraftAlreadyConfirmed 等 |
| T-074 | 🟢 定义 Generation Repository 接口 | `internal/domain/generation/repository.go` | T-070, T-072 | GenerationTaskRepository, CaseDraftRepository |

### Phase 2 完成标准
- [ ] 所有领域模型单元测试通过
- [ ] `go vet ./internal/domain/...` 无警告
- [ ] 领域层无外部依赖（仅标准库 + google/uuid）

---

## Phase 3: Repository 层 (P0)


### 3.1 事务管理器

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-075 | 🔴 编写事务管理器测试 | `internal/repository/transaction_test.go` | T-013 | 测试事务提交、回滚、嵌套事务 |
| T-076 | 🟢 实现事务管理器 | `internal/repository/transaction.go` | T-075 | TxManager, WithTransaction, TxFromContext |

### 3.2 Identity Repository

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-077 | 🔴 编写 UserRepository 测试 | `internal/repository/identity/user_repo_test.go` | T-024, T-076 | 使用 testcontainers，测试 CRUD |
| T-078 | 🟢 实现 UserRepository | `internal/repository/identity/user_repo.go` | T-077 | 实现 domain.UserRepository 接口 |

### 3.3 Project Repository

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-079 | [P] 🔴 编写 ProjectRepository 测试 | `internal/repository/project/project_repo_test.go` | T-036, T-076 | 测试 CRUD、唯一性校验 |
| T-080 | [P] 🟢 实现 ProjectRepository | `internal/repository/project/project_repo.go` | T-079 | 实现 domain.ProjectRepository |
| T-081 | [P] 🔴 编写 ModuleRepository 测试 | `internal/repository/project/module_repo_test.go` | T-036, T-080 | 测试 CRUD、级联删除 |
| T-082 | [P] 🟢 实现 ModuleRepository | `internal/repository/project/module_repo.go` | T-081 | 实现 domain.ModuleRepository |
| T-083 | [P] 🔴 编写 ProjectConfigRepository 测试 | `internal/repository/project/config_repo_test.go` | T-036 | 测试 KV 存储 |
| T-084 | [P] 🟢 实现 ProjectConfigRepository | `internal/repository/project/config_repo.go` | T-083 | 实现 domain.ProjectConfigRepository |

### 3.4 TestCase Repository

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-085 | 🔴 编写 TestCaseRepository 测试 | `internal/repository/testcase/case_repo_test.go` | T-046, T-076 | 测试 CRUD、编号唯一性、日期计数 |
| T-086 | 🟢 实现 TestCaseRepository | `internal/repository/testcase/case_repo.go` | T-085 | 实现 domain.TestCaseRepository |

### 3.5 TestPlan Repository

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-087 | [P] 🔴 编写 TestPlanRepository 测试 | `internal/repository/testplan/plan_repo_test.go` | T-054, T-076 | 测试 CRUD、状态转换 |
| T-088 | [P] 🟢 实现 TestPlanRepository | `internal/repository/testplan/plan_repo.go` | T-087 | 实现 domain.TestPlanRepository |
| T-089 | [P] 🔴 编写 TestResultRepository 测试 | `internal/repository/testplan/result_repo_test.go` | T-054 | 测试执行结果录入 |
| T-090 | [P] 🟢 实现 TestResultRepository | `internal/repository/testplan/result_repo.go` | T-089 | 实现 domain.TestResultRepository |

### 3.6 Knowledge Repository

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-091 | [P] 🔴 编写 DocumentRepository 测试 | `internal/repository/knowledge/document_repo_test.go` | T-064, T-076 | 测试 CRUD、状态更新 |
| T-092 | [P] 🟢 实现 DocumentRepository | `internal/repository/knowledge/document_repo.go` | T-091 | 实现 domain.DocumentRepository |
| T-093 | [P] 🔴 编写 DocumentChunkRepository 测试 | `internal/repository/knowledge/chunk_repo_test.go` | T-064 | 测试批量保存、按文档查询 |
| T-094 | [P] 🟢 实现 DocumentChunkRepository | `internal/repository/knowledge/chunk_repo.go` | T-093 | 实现 domain.DocumentChunkRepository |

### 3.7 Generation Repository

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-095 | [P] 🔴 编写 GenerationTaskRepository 测试 | `internal/repository/generation/task_repo_test.go` | T-074, T-076 | 测试 CRUD、状态查询 |
| T-096 | [P] 🟢 实现 GenerationTaskRepository | `internal/repository/generation/task_repo.go` | T-095 | 实现 domain.GenerationTaskRepository |
| T-097 | [P] 🔴 编写 CaseDraftRepository 测试 | `internal/repository/generation/draft_repo_test.go` | T-074 | 测试 CRUD、按任务查询 |
| T-098 | [P] 🟢 实现 CaseDraftRepository | `internal/repository/generation/draft_repo.go` | T-097 | 实现 domain.CaseDraftRepository |

### Phase 3 完成标准
- [ ] 所有 Repository 单元测试通过
- [ ] 使用 testcontainers 进行集成测试
- [ ] 事务管理器支持跨 Repository 事务

---

## Phase 4: Service 层 (P1)

### 4.1 Identity Service

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-099 | 🔴 编写 AuthService 测试 | `internal/service/identity/auth_service_test.go` | T-078 | 测试注册、登录、Token 生成 |
| T-100 | 🟢 实现 AuthService | `internal/service/identity/auth_service.go` | T-099 | Register, Login, ValidateToken |

### 4.2 Project Service

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-101 | 🔴 编写 ProjectService 测试 | `internal/service/project/project_service_test.go` | T-080, T-082, T-084 | 测试项目 CRUD、模块管理、配置管理 |
| T-102 | 🟢 实现 ProjectService | `internal/service/project/project_service.go` | T-101 | CreateProject, CreateModule, SetConfig |

### 4.3 TestCase Service

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-103 | 🔴 编写 TestCaseService 测试 | `internal/service/testcase/case_service_test.go` | T-086, T-082 | 测试用例 CRUD、编号生成、需求追溯 |
| T-104 | 🟢 实现 TestCaseService | `internal/service/testcase/case_service.go` | T-103 | CreateCase, UpdateCase, GetCaseDetail |

### 4.4 TestPlan Service

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-105 | 🔴 编写 PlanService 测试 | `internal/service/testplan/plan_service_test.go` | T-088, T-090, T-086 | 测试计划 CRUD、用例关联、结果录入 |
| T-106 | 🟢 实现 PlanService | `internal/service/testplan/plan_service.go` | T-105 | CreatePlan, AddCase, RecordResult |

### 4.5 Knowledge Service

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-107 | 🔴 编写 DocumentService 测试 | `internal/service/knowledge/document_service_test.go` | T-092, T-094 | 测试文档上传、状态更新 |
| T-108 | 🟢 实现 DocumentService | `internal/service/knowledge/document_service.go` | T-107 | UploadDocument, GetDocument, DeleteDocument |

### 4.6 Generation Service

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-109 | 🔴 编写 RAGService 测试 | `internal/service/generation/rag_service_test.go` | T-064 | Mock 向量检索，测试 Top-K |
| T-110 | 🟢 实现 RAGService | `internal/service/generation/rag_service.go` | T-109 | Retrieve, CalculateConfidence |
| T-111 | 🔴 编写 LLMService 测试 | `internal/service/generation/llm_service_test.go` | - | Mock DeepSeek API，测试生成 |
| T-112 | 🟢 实现 LLMService | `internal/service/generation/llm_service.go` | T-111 | GenerateCases, 超时控制 |
| T-113 | 🔴 编写 GenerationService 测试 | `internal/service/generation/generation_service_test.go` | T-096, T-098, T-110, T-112 | 测试任务创建、草稿确认/拒绝 |
| T-114 | 🟢 实现 GenerationService | `internal/service/generation/generation_service.go` | T-113 | CreateTask, ConfirmDraft, RejectDraft, BatchConfirm |

### Phase 4 完成标准
- [ ] 所有 Service 单元测试通过
- [ ] Mock Repository 进行隔离测试
- [ ] 错误处理使用 ierrors 映射

---

## Phase 5: Transport 层 (P1)

### 5.1 HTTP Server & Router

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-115 | 🔴 编写 HTTP Server 测试 | `internal/transport/http/server_test.go` | - | 测试服务器启动、关闭 |
| T-116 | 🟢 实现 HTTP Server | `internal/transport/http/server.go` | T-115 | NewServer, Run, Shutdown |
| T-117 | 🔴 编写 Router 测试 | `internal/transport/http/router_test.go` | T-116 | 测试路由注册、路径匹配 |
| T-118 | 🟢 实现 Router | `internal/transport/http/router.go` | T-117 | 使用 chi 路由，注册所有 Handler |

### 5.2 Middleware

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-119 | [P] 🔴 编写 Auth Middleware 测试 | `internal/transport/http/middleware/auth_test.go` | - | 测试 JWT 验证、权限校验 |
| T-120 | [P] 🟢 实现 Auth Middleware | `internal/transport/http/middleware/auth.go` | T-119 | JWT 解析，用户信息注入 Context |
| T-121 | [P] 🔴 编写 Logging Middleware 测试 | `internal/transport/http/middleware/logging_test.go` | - | 测试日志输出格式 |
| T-122 | [P] 🟢 实现 Logging Middleware | `internal/transport/http/middleware/logging.go` | T-121 | zerolog 结构化日志 |
| T-123 | [P] 🔴 编写 Recovery Middleware 测试 | `internal/transport/http/middleware/recovery_test.go` | - | 测试 Panic 恢复 |
| T-124 | [P] 🟢 实现 Recovery Middleware | `internal/transport/http/middleware/recovery.go` | T-123 | Panic 捕获，返回 500 |
| T-125 | [P] 🔴 编写 Metrics Middleware 测试 | `internal/transport/http/middleware/metrics_test.go` | - | 测试 Prometheus 指标 |
| T-126 | [P] 🟢 实现 Metrics Middleware | `internal/transport/http/middleware/metrics.go` | T-125 | 请求计数、延迟直方图 |

### 5.3 Handlers

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-127 | 🔴 编写 Identity Handler 测试 | `internal/transport/http/handler/identity_test.go` | T-100 | 测试注册、登录 API |
| T-128 | 🟢 实现 Identity Handler | `internal/transport/http/handler/identity.go` | T-127 | Register, Login Handler |
| T-129 | 🔴 编写 Project Handler 测试 | `internal/transport/http/handler/project_test.go` | T-102 | 测试项目、模块、配置 API |
| T-130 | 🟢 实现 Project Handler | `internal/transport/http/handler/project.go` | T-129 | 项目 CRUD、模块 CRUD、配置 CRUD |
| T-131 | 🔴 编写 Knowledge Handler 测试 | `internal/transport/http/handler/knowledge_test.go` | T-108 | 测试文档上传、查询 API |
| T-132 | 🟢 实现 Knowledge Handler | `internal/transport/http/handler/knowledge.go` | T-131 | 文档上传、分块查询 |
| T-133 | 🔴 编写 TestCase Handler 测试 | `internal/transport/http/handler/testcase_test.go` | T-104 | 测试用例 CRUD API |
| T-134 | 🟢 实现 TestCase Handler | `internal/transport/http/handler/testcase.go` | T-133 | 用例 CRUD、需求追溯 |
| T-135 | 🔴 编写 TestPlan Handler 测试 | `internal/transport/http/handler/testplan_test.go` | T-106 | 测试计划、执行结果 API |
| T-136 | 🟢 实现 TestPlan Handler | `internal/transport/http/handler/testplan.go` | T-135 | 计划 CRUD、结果录入 |
| T-137 | 🔴 编写 Generation Handler 测试 | `internal/transport/http/handler/generation_test.go` | T-114 | 测试生成任务、草稿操作 API |
| T-138 | 🟢 实现 Generation Handler | `internal/transport/http/handler/generation.go` | T-137 | 任务创建、草稿确认/拒绝 |

### Phase 5 完成标准
- [ ] 所有 Handler 测试通过
- [ ] API 响应符合 openapi.yaml 定义
- [ ] 中间件链正确执行

---

## Phase 6: 集成测试 (P2)

### 6.1 API 集成测试

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-139 | [P] 🔴 编写项目管理集成测试 | `tests/integration/api/project_test.go` | T-130 | 端到端测试项目 CRUD |
| T-140 | [P] 🔴 编写用例管理集成测试 | `tests/integration/api/testcase_test.go` | T-134 | 端到端测试用例 CRUD |
| T-141 | [P] 🔴 编写生成任务集成测试 | `tests/integration/api/generation_test.go` | T-138 | 端到端测试生成流程 |
| T-142 | [P] 🔴 编写测试计划集成测试 | `tests/integration/api/testplan_test.go` | T-136 | 端到端测试执行流程 |

### 6.2 性能测试

| ID | 任务 | 文件 | 依赖 | 验收标准 |
|----|------|------|------|----------|
| T-143 | 📝 编写 k6 性能测试脚本 | `tests/performance/load_test.js` | T-118 | P99 < 500ms |

### Phase 6 完成标准
- [ ] 所有集成测试通过
- [ ] k6 性能测试达标

---

## 任务依赖图

```
Phase 1 (基础设施)
T-001 → T-002 → T-003 → T-004
                  ↓
T-006 → T-007 → T-008 → T-009
      → T-010 → T-011
                  ↓
T-012 → T-013 → T-016 → T-017 → T-018
      → T-014 → T-015 ↗

Phase 2 (领域模型) - 可并行
┌─ T-019 → T-020 → T-021 → T-022 → T-023
│                                   → T-024
├─ T-025 → T-026 → T-029 → T-030 → T-035
│   → T-027 → T-028 → T-031 → T-032   → T-036
│                   → T-033 → T-034
├─ T-037 → T-038 → T-043 → T-044 → T-045
│   → T-039 → T-040                  → T-046
│   → T-041 → T-042
├─ T-047 → T-048 → T-049 → T-050 → T-053
│                   → T-051 → T-052   → T-054
├─ T-055 → T-056 → T-059 → T-060 → T-063
│   → T-057 → T-058                  → T-064
│                   → T-061 → T-062
└─ T-065 → T-066 → T-069 → T-070 → T-073
    → T-067 → T-068                  → T-074
                    → T-071 → T-072

Phase 3 (Repository)
T-075 → T-076
    ↓
┌─ T-077 → T-078
├─ T-079 → T-080 → T-081 → T-082
│               → T-083 → T-084
├─ T-085 → T-086
├─ T-087 → T-088
│   → T-089 → T-090
├─ T-091 → T-092
│   → T-093 → T-094
└─ T-095 → T-096
    → T-097 → T-098

Phase 4 (Service)
┌─ T-099 → T-100
├─ T-101 → T-102
├─ T-103 → T-104
├─ T-105 → T-106
├─ T-107 → T-108
└─ T-109 → T-110 → T-113 → T-114
    → T-111 → T-112 ↗

Phase 5 (Transport)
T-115 → T-116 → T-117 → T-118
    ↓
┌─ T-119 → T-120
├─ T-121 → T-122
├─ T-123 → T-124
├─ T-125 → T-126
├─ T-127 → T-128
├─ T-129 → T-130
├─ T-131 → T-132
├─ T-133 → T-134
├─ T-135 → T-136
└─ T-137 → T-138

Phase 6 (集成测试)
┌─ T-139
├─ T-140
├─ T-141
├─ T-142
└─ T-143
```

---

## 统计信息

| 阶段 | 测试任务 | 实现任务 | 文档任务 | 总计 |
|------|----------|----------|----------|------|
| Phase 1 | 7 | 10 | 1 | 18 |
| Phase 2 | 24 | 24 | 6 | 54 |
| Phase 3 | 12 | 12 | 0 | 24 |
| Phase 4 | 8 | 8 | 0 | 16 |
| Phase 5 | 12 | 12 | 0 | 24 |
| Phase 6 | 4 | 0 | 1 | 5 |
| **总计** | **67** | **66** | **8** | **141** |

---

## 执行建议

1. **严格按顺序执行**：遵循 TDD 原则，先写测试再实现
2. **并行执行**：标记 [P] 的任务可由不同 Agent 并行执行
3. **每个任务完成后**：
   - 运行 `go test ./...` 确保测试通过
   - 运行 `go fmt ./...` 和 `go vet ./...`
4. **每个阶段完成后**：
   - 运行 `golangci-lint run ./...`
   - 检查测试覆盖率
