# 智能测试管理平台 - Go 工程设计文档

**版本**：2.0
**日期**：2026-04-01
**架构师**：Go 资深架构师
**Go 版本**：1.24+

---

## 1. 模块拓扑图

### 1.1 分层架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        cmd/server/main.go                        │
│                      (程序入口 + 优雅停机)                         │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                      internal/app/app.go                         │
│                    (依赖注入 + 生命周期编排)                       │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                    internal/transport/http/                      │
│              (HTTP Handler + 路由 + 中间件)                       │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                      internal/service/                           │
│                   (应用服务，编排领域对象)                         │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                       internal/domain/                           │
│            (聚合根 + 实体 + 值对象 + Repository 接口)              │
│                     ⚡ 纯 Go，无外部依赖 ⚡                        │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                     internal/repository/                         │
│                  (Repository 实现 + DB/Cache)                    │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 限界上下文模块矩阵

| 上下文 | domain/ | service/ | repository/ | transport/ |
|---|---|---|---|---|
| Identity | `identity/` | `identity/` | `identity/` | `identity.go` |
| Project | `project/` | `project/` | `project/` | `project.go` |
| Knowledge | `knowledge/` | `knowledge/` | `knowledge/` | `knowledge.go` |
| TestCase | `testcase/` | `testcase/` | `testcase/` | `testcase.go` |
| TestPlan | `testplan/` | `testplan/` | `testplan/` | `testplan.go` |
| Generation | `generation/` | `generation/` | `generation/` | `generation.go` |

---

## 2. 目录结构

```
aitestos/
├── cmd/
│   └── server/
│       └── main.go                    # 程序入口
├── internal/
│   ├── app/
│   │   ├── app.go                     # 应用容器
│   │   └── wire.go                    # Wire 依赖注入
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
│   │       ├── errors.go
│   │       └── repository.go
│   ├── ierrors/
│   │   ├── codes.go                   # 统一错误码
│   │   ├── mapping.go                 # 错误映射
│   │   └── response.go                # HTTP 响应
│   ├── repository/                    # 数据持久化实现
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
│   │       └── generation_service.go
│   └── transport/
│       └── http/
│           ├── server.go              # HTTP 服务器
│           ├── router.go              # 路由定义
│           ├── middleware/
│           │   ├── auth.go            # JWT 认证
│           │   ├── logging.go         # 日志
│           │   ├── recovery.go        # Panic 恢复
│           │   └── metrics.go         # Prometheus 指标
│           └── handler/
│               ├── identity.go
│               ├── project.go
│               ├── knowledge.go
│               ├── testcase.go
│               ├── testplan.go
│               └── generation.go
├── pkg/                               # 公共工具（可被外部引用）
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

---

## 3. 接口与抽象设计

### 3.1 核心原则

```go
// ✅ 正确：接受接口，返回结构体
func NewUserService(repo identity.UserRepository) *UserService {
    return &UserService{repo: repo}
}

// ❌ 错误：返回接口
func NewUserService(repo identity.UserRepository) identity.UserService { ... }
```

### 3.2 Domain 接口定义

```go
// internal/domain/identity/repository.go
package identity

import (
    "context"
    "github.com/google/uuid"
)

// UserRepository 定义在 domain 层，由 repository 层实现
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
}

// internal/domain/testcase/repository.go
package testcase

type TestCaseRepository interface {
    Save(ctx context.Context, tc *TestCase) error
    FindByID(ctx context.Context, id uuid.UUID) (*TestCase, error)
    FindByModuleID(ctx context.Context, moduleID uuid.UUID, opts QueryOptions) ([]*TestCase, error)
}

// QueryOptions 查询选项值对象
type QueryOptions struct {
    Offset  int
    Limit   int
    OrderBy string
}
```

### 3.3 Service 接口（由 transport 层消费）

```go
// internal/service/identity/auth_service.go
package identity

import (
    "context"
    "github.com/google/uuid"
    "github.com/liang21/aitestos/internal/domain/identity"
)

// AuthService 应用服务接口
type AuthService interface {
    Register(ctx context.Context, req *RegisterRequest) (*identity.User, error)
    Login(ctx context.Context, email, password string) (*Token, error)
    ValidateToken(ctx context.Context, token string) (*Claims, error)
}

// RegisterRequest 请求结构体
type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3,max=32"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// Token 返回结构体
type Token struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}

// Claims 令牌声明
type Claims struct {
    UserID   uuid.UUID
    Username string
    Role     identity.UserRole
}
```

---

## 4. 数据模型定义

### 4.1 聚合根示例

```go
// internal/domain/identity/user.go
package identity

import (
    "time"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

// User 聚合根
type User struct {
    id        uuid.UUID
    username  string
    email     string
    password  string // 加密后的密码
    role      UserRole
    createdAt time.Time
    updatedAt time.Time
}

// NewUser 创建用户（工厂函数）
func NewUser(username, email, rawPassword string, role UserRole) (*User, error) {
    if err := validateEmail(email); err != nil {
        return nil, err
    }
    if err := validateUsername(username); err != nil {
        return nil, err
    }

    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(rawPassword),
        bcrypt.DefaultCost,
    )
    if err != nil {
        return nil, err
    }

    now := time.Now()
    return &User{
        id:        uuid.New(),
        username:  username,
        email:     email,
        password:  string(hashedPassword),
        role:      role,
        createdAt: now,
        updatedAt: now,
    }, nil
}

// ID 返回用户 ID（只读）
func (u *User) ID() uuid.UUID { return u.id }

// Username 返回用户名（只读）
func (u *User) Username() string { return u.username }

// Email 返回邮箱（只读）
func (u *User) Email() string { return u.email }

// Role 返回角色（只读）
func (u *User) Role() UserRole { return u.role }

// VerifyPassword 验证密码
func (u *User) VerifyPassword(rawPassword string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(u.password),
        []byte(rawPassword),
    )
    return err == nil
}

// ChangePassword 修改密码（领域行为）
func (u *User) ChangePassword(oldPassword, newPassword string) error {
    if !u.VerifyPassword(oldPassword) {
        return ErrPasswordMismatch
    }

    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(newPassword),
        bcrypt.DefaultCost,
    )
    if err != nil {
        return err
    }

    u.password = string(hashedPassword)
    u.updatedAt = time.Now()
    return nil
}

// ========== JSON 序列化（用于 API 响应）==========

type UserJSON struct {
    ID        uuid.UUID `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ToJSON 转换为 JSON 结构
func (u *User) ToJSON() *UserJSON {
    return &UserJSON{
        ID:        u.id,
        Username:  u.username,
        Email:     u.email,
        Role:      string(u.role),
        CreatedAt: u.createdAt,
        UpdatedAt: u.updatedAt,
    }
}
```

### 4.2 值对象示例

```go
// internal/domain/identity/user_role.go
package identity

import "errors"

// UserRole 值对象
type UserRole string

const (
    RoleSuperAdmin UserRole = "super_admin"
    RoleAdmin      UserRole = "admin"
    RoleNormal     UserRole = "normal"
)

// ParseUserRole 解析角色
func ParseUserRole(s string) (UserRole, error) {
    switch UserRole(s) {
    case RoleSuperAdmin, RoleAdmin, RoleNormal:
        return UserRole(s), nil
    default:
        return "", errors.New("invalid user role")
    }
}

// IsAdmin 是否为管理员
func (r UserRole) IsAdmin() bool {
    return r == RoleSuperAdmin || r == RoleAdmin
}
```

```go
// internal/domain/testcase/case_number.go
package testcase

import (
    "errors"
    "fmt"
    "regexp"
    "time"
)

// CaseNumber 用例编号值对象
// 格式: PRJ-MOD-YYYYMMDD-NNN
type CaseNumber string

var caseNumberRegex = regexp.MustCompile(`^[A-Z]{3}-[A-Z]{3}-\d{8}-\d{3}$`)

// ParseCaseNumber 解析用例编号
func ParseCaseNumber(s string) (CaseNumber, error) {
    if !caseNumberRegex.MatchString(s) {
        return "", ErrInvalidCaseNumber
    }
    return CaseNumber(s), nil
}

// GenerateCaseNumber 生成用例编号
func GenerateCaseNumber(projectPrefix, modulePrefix string, seq int) CaseNumber {
    date := time.Now().Format("20060102")
    return CaseNumber(fmt.Sprintf("%s-%s-%s-%03d",
        projectPrefix, modulePrefix, date, seq))
}

// String 转换为字符串
func (n CaseNumber) String() string { return string(n) }
```

### 4.3 复杂聚合根示例

```go
// internal/domain/testcase/test_case.go
package testcase

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

// TestCase 聚合根
type TestCase struct {
    id            uuid.UUID
    moduleID      uuid.UUID
    userID        uuid.UUID
    number        CaseNumber
    title         string
    preconditions Preconditions
    steps         Steps
    expected      ExpectedResult
    aiMetadata    AiMetadata
    caseType      CaseType
    priority      Priority
    status        CaseStatus
    createdAt     time.Time
    updatedAt     time.Time
}

// Preconditions 前置条件值对象
type Preconditions []string

// Steps 步骤值对象
type Steps []string

// ExpectedResult 预期结果值对象
type ExpectedResult map[string]any

// AiMetadata AI 元数据值对象
type AiMetadata struct {
    ReferencedChunks []uuid.UUID `json:"referenced_chunks"`
    TaskID           uuid.UUID   `json:"task_id"`
    Confidence       float64     `json:"confidence"`
}

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

// UpdateStatus 更新状态（领域行为）
func (tc *TestCase) UpdateStatus(status CaseStatus) {
    tc.status = status
    tc.updatedAt = time.Now()
}

// ========== 数据库映射 ==========

// TestCaseRow 数据库行结构
type TestCaseRow struct {
    ID            uuid.UUID       `db:"id"`
    ModuleID      uuid.UUID       `db:"module_id"`
    UserID        uuid.UUID       `db:"user_id"`
    Number        string          `db:"number"`
    Title         string          `db:"title"`
    Preconditions json.RawMessage `db:"preconditions"`
    Steps         json.RawMessage `db:"steps"`
    Expected      json.RawMessage `db:"expected"`
    AiMetadata    json.RawMessage `db:"ai_metadata"`
    CaseType      string          `db:"case_type"`
    Priority      string          `db:"priority"`
    Status        string          `db:"status"`
    CreatedAt     time.Time       `db:"created_at"`
    UpdatedAt     time.Time       `db:"updated_at"`
}

// ToRow 转换为数据库行
func (tc *TestCase) ToRow() (*TestCaseRow, error) {
    preconditionsJSON, err := json.Marshal(tc.preconditions)
    if err != nil {
        return nil, err
    }
    stepsJSON, err := json.Marshal(tc.steps)
    if err != nil {
        return nil, err
    }
    expectedJSON, err := json.Marshal(tc.expected)
    if err != nil {
        return nil, err
    }
    aiMetadataJSON, err := json.Marshal(tc.aiMetadata)
    if err != nil {
        return nil, err
    }

    return &TestCaseRow{
        ID:            tc.id,
        ModuleID:      tc.moduleID,
        UserID:        tc.userID,
        Number:        tc.number.String(),
        Title:         tc.title,
        Preconditions: preconditionsJSON,
        Steps:         stepsJSON,
        Expected:      expectedJSON,
        AiMetadata:    aiMetadataJSON,
        CaseType:      string(tc.caseType),
        Priority:      string(tc.priority),
        Status:        string(tc.status),
        CreatedAt:     tc.createdAt,
        UpdatedAt:     tc.updatedAt,
    }, nil
}

// FromRow 从数据库行恢复
func FromRow(row *TestCaseRow) (*TestCase, error) {
    var preconditions Preconditions
    if err := json.Unmarshal(row.Preconditions, &preconditions); err != nil {
        return nil, err
    }

    var steps Steps
    if err := json.Unmarshal(row.Steps, &steps); err != nil {
        return nil, err
    }

    var expected ExpectedResult
    if err := json.Unmarshal(row.Expected, &expected); err != nil {
        return nil, err
    }

    var aiMetadata AiMetadata
    if err := json.Unmarshal(row.AiMetadata, &aiMetadata); err != nil {
        return nil, err
    }

    number, err := ParseCaseNumber(row.Number)
    if err != nil {
        return nil, err
    }

    return &TestCase{
        id:            row.ID,
        moduleID:      row.ModuleID,
        userID:        row.UserID,
        number:        number,
        title:         row.Title,
        preconditions: preconditions,
        steps:         steps,
        expected:      expected,
        aiMetadata:    aiMetadata,
        caseType:      CaseType(row.CaseType),
        priority:      Priority(row.Priority),
        status:        CaseStatus(row.Status),
        createdAt:     row.CreatedAt,
        updatedAt:     row.UpdatedAt,
    }, nil
}
```

---

## 5. 并发模型设计

### 5.1 Context 传递规范

```go
// internal/service/testplan/plan_service.go
package testplan

import (
    "context"
    "time"

    "github.com/liang21/aitestos/internal/domain/testplan"
)

type PlanService struct {
    planRepo   testplan.TestPlanRepository
    resultRepo testplan.TestResultRepository
}

func NewPlanService(
    planRepo testplan.TestPlanRepository,
    resultRepo testplan.TestResultRepository,
) *PlanService {
    return &PlanService{
        planRepo:   planRepo,
        resultRepo: resultRepo,
    }
}

// ExecutePlan 执行测试计划（带超时控制）
func (s *PlanService) ExecutePlan(ctx context.Context, planID uuid.UUID) error {
    // 设置操作超时
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // 查询计划
    plan, err := s.planRepo.FindByID(ctx, planID)
    if err != nil {
        return fmt.Errorf("find plan: %w", err)
    }

    // 执行逻辑...
    return nil
}
```

### 5.2 Goroutine 生命周期管理

```go
// internal/app/app.go
package app

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/rs/zerolog/log"
)

type App struct {
    httpServer *http.Server
    // 其他组件...
}

// Run 启动应用（优雅停机）
func (a *App) Run(ctx context.Context) error {
    // 启动 HTTP 服务器
    errCh := make(chan error, 1)
    go func() {
        log.Info().Str("addr", a.httpServer.Addr).Msg("starting http server")
        if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            errCh <- err
        }
    }()

    // 监听系统信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    select {
    case err := <-errCh:
        return fmt.Errorf("http server error: %w", err)
    case sig := <-quit:
        log.Info().Str("signal", sig.String()).Msg("shutting down server")
    case <-ctx.Done():
        log.Info().Msg("context cancelled, shutting down")
    }

    // 优雅停机（最多等待 30 秒）
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
        log.Error().Err(err).Msg("http server shutdown error")
    }

    log.Info().Msg("server stopped")
    return nil
}
```

### 5.3 并发安全：sync.Mutex vs Channel

```go
// 使用 sync.Mutex 的场景：保护共享状态
// internal/repository/generation/task_cache.go
package generation

import (
    "sync"
    "time"

    "github.com/google/uuid"
)

// TaskCache 任务状态缓存（写多读少）
type TaskCache struct {
    mu    sync.RWMutex
    tasks map[uuid.UUID]*TaskStatus
}

type TaskStatus struct {
    Status    string
    UpdatedAt time.Time
}

func NewTaskCache() *TaskCache {
    return &TaskCache{
        tasks: make(map[uuid.UUID]*TaskStatus),
    }
}

func (c *TaskCache) Set(id uuid.UUID, status *TaskStatus) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.tasks[id] = status
}

func (c *TaskCache) Get(id uuid.UUID) (*TaskStatus, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    status, ok := c.tasks[id]
    return status, ok
}

// 使用 Channel 的场景：任务队列
// internal/service/generation/worker.go
package generation

import (
    "context"

    "github.com/google/uuid"
    "github.com/rs/zerolog/log"
)

type TaskWorker struct {
    taskCh  chan uuid.UUID
    handler func(ctx context.Context, taskID uuid.UUID) error
    workers int
}

func NewTaskWorker(workers int, handler func(ctx context.Context, taskID uuid.UUID) error) *TaskWorker {
    return &TaskWorker{
        taskCh:  make(chan uuid.UUID, 100),
        handler: handler,
        workers: workers,
    }
}

func (w *TaskWorker) Start(ctx context.Context) {
    for i := 0; i < w.workers; i++ {
        go w.runWorker(ctx)
    }
}

func (w *TaskWorker) runWorker(ctx context.Context) {
    for {
        select {
        case taskID := <-w.taskCh:
            if err := w.handler(ctx, taskID); err != nil {
                log.Error().Err(err).Str("task_id", taskID.String()).Msg("task failed")
            }
        case <-ctx.Done():
            return
        }
    }
}

func (w *TaskWorker) Submit(taskID uuid.UUID) bool {
    select {
    case w.taskCh <- taskID:
        return true
    default:
        return false // 队列满
    }
}
```

---

## 6. 存储层设计

### 6.1 Repository 实现

```go
// internal/repository/testcase/case_repo.go
package testcase

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    domain "github.com/liang21/aitestos/internal/domain/testcase"
)

type testCaseRepository struct {
    db *sqlx.DB
}

func NewTestCaseRepository(db *sqlx.DB) domain.TestCaseRepository {
    return &testCaseRepository{db: db}
}

func (r *testCaseRepository) Save(ctx context.Context, tc *domain.TestCase) error {
    row, err := tc.ToRow()
    if err != nil {
        return fmt.Errorf("to row: %w", err)
    }

    query := `
        INSERT INTO test_case (
            id, module_id, user_id, number, title,
            preconditions, steps, expected, ai_metadata,
            case_type, priority, status, created_at, updated_at
        ) VALUES (
            :id, :module_id, :user_id, :number, :title,
            :preconditions, :steps, :expected, :ai_metadata,
            :case_type, :priority, :status, :created_at, :updated_at
        )
    `

    _, err = r.db.NamedExecContext(ctx, query, row)
    if err != nil {
        return fmt.Errorf("insert test case: %w", err)
    }

    return nil
}

func (r *testCaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.TestCase, error) {
    query := `SELECT * FROM test_case WHERE id = $1`

    var row domain.TestCaseRow
    err := r.db.GetContext(ctx, &row, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrCaseNotFound
        }
        return nil, fmt.Errorf("find by id: %w", err)
    }

    return domain.FromRow(&row)
}

func (r *testCaseRepository) FindByModuleID(
    ctx context.Context,
    moduleID uuid.UUID,
    opts domain.QueryOptions,
) ([]*domain.TestCase, error) {
    query := `
        SELECT * FROM test_case
        WHERE module_id = $1
        ORDER BY ` + opts.OrderBy + `
        LIMIT $2 OFFSET $3
    `

    var rows []domain.TestCaseRow
    err := r.db.SelectContext(ctx, &rows, query, moduleID, opts.Limit, opts.Offset)
    if err != nil {
        return nil, fmt.Errorf("find by module id: %w", err)
    }

    cases := make([]*domain.TestCase, 0, len(rows))
    for _, row := range rows {
        tc, err := domain.FromRow(&row)
        if err != nil {
            return nil, fmt.Errorf("from row: %w", err)
        }
        cases = append(cases, tc)
    }

    return cases, nil
}
```

### 6.2 事务处理（Context 传递 Tx）

使用 context 传递事务，支持跨 Repository 调用共享同一事务。

```go
// internal/repository/transaction.go
package repository

import (
    "context"
    "fmt"

    "github.com/jmoiron/sqlx"
)

// txKey 事务 context key
type txKey struct{}

// TxManager 事务管理器
type TxManager struct {
    db *sqlx.DB
}

func NewTxManager(db *sqlx.DB) *TxManager {
    return &TxManager{db: db}
}

// WithTransaction 执行事务（通过 context 传递）
func (m *TxManager) WithTransaction(
    ctx context.Context,
    fn func(ctx context.Context) error,
) error {
    tx, err := m.db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }

    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback()
            panic(p)
        }
    }()

    // 将 tx 注入到 context
    txCtx := context.WithValue(ctx, txKey{}, tx)

    if err := fn(txCtx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("rollback error: %v, original error: %w", rbErr, err)
        }
        return err
    }

    return tx.Commit()
}

// TxFromContext 从 context 获取事务
func TxFromContext(ctx context.Context) *sqlx.Tx {
    if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
        return tx
    }
    return nil
}
```

### 6.3 Repository 支持事务传播

```go
// internal/repository/testcase/case_repo.go
package testcase

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    domain "github.com/liang21/aitestos/internal/domain/testcase"
    "github.com/liang21/aitestos/internal/repository"
)

type testCaseRepository struct {
    db *sqlx.DB
}

func NewTestCaseRepository(db *sqlx.DB) domain.TestCaseRepository {
    return &testCaseRepository{db: db}
}

// getExecutor 获取执行器（优先使用事务）
func (r *testCaseRepository) getExecutor(ctx context.Context) executor {
    if tx := repository.TxFromContext(ctx); tx != nil {
        return tx
    }
    return r.db
}

type executor interface {
    GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
    SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

func (r *testCaseRepository) Save(ctx context.Context, tc *domain.TestCase) error {
    row, err := tc.ToRow()
    if err != nil {
        return fmt.Errorf("to row: %w", err)
    }

    query := `
        INSERT INTO test_case (
            id, module_id, user_id, number, title,
            preconditions, steps, expected, ai_metadata,
            case_type, priority, status, created_at, updated_at
        ) VALUES (
            :id, :module_id, :user_id, :number, :title,
            :preconditions, :steps, :expected, :ai_metadata,
            :case_type, :priority, :status, :created_at, :updated_at
        )
    `

    _, err = r.getExecutor(ctx).NamedExecContext(ctx, query, row)
    if err != nil {
        return fmt.Errorf("insert test case: %w", err)
    }

    return nil
}

func (r *testCaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.TestCase, error) {
    query := `SELECT * FROM test_case WHERE id = $1`

    var row domain.TestCaseRow
    err := r.getExecutor(ctx).GetContext(ctx, &row, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrCaseNotFound
        }
        return nil, fmt.Errorf("find by id: %w", err)
    }

    return domain.FromRow(&row)
}
```

### 6.4 跨 Repository 事务示例

```go
// internal/service/testplan/plan_service.go
package testplan

import (
    "context"
    "fmt"

    "github.com/google/uuid"
)

type PlanService struct {
    txManager  *repository.TxManager
    planRepo   domain.TestPlanRepository
    caseRepo   domain.TestCaseRepository
    resultRepo domain.TestResultRepository
}

// CreatePlanWithCases 创建计划并关联用例（跨 Repository 事务）
func (s *PlanService) CreatePlanWithCases(
    ctx context.Context,
    plan *domain.TestPlan,
    caseIDs []uuid.UUID,
) error {
    return s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
        // 1. 创建计划（planRepo 内部会自动使用 ctx 中的事务）
        if err := s.planRepo.Save(ctx, plan); err != nil {
            return fmt.Errorf("save plan: %w", err)
        }

        // 2. 验证用例存在（caseRepo 内部会自动使用 ctx 中的事务）
        for _, caseID := range caseIDs {
            _, err := s.caseRepo.FindByID(ctx, caseID)
            if err != nil {
                return fmt.Errorf("find case %s: %w", caseID, err)
            }
        }

        // 3. 关联用例
        for _, caseID := range caseIDs {
            if err := s.planRepo.AddCase(ctx, plan.ID(), caseID); err != nil {
                return fmt.Errorf("add case %s: %w", caseID, err)
            }
        }

        return nil
    })
}

// RecordResultAndPlanStatus 记录结果并更新计划状态（跨 Repository 事务）
func (s *PlanService) RecordResultAndPlanStatus(
    ctx context.Context,
    result *domain.TestResult,
) error {
    return s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
        // 1. 保存执行结果
        if err := s.resultRepo.Save(ctx, result); err != nil {
            return fmt.Errorf("save result: %w", err)
        }

        // 2. 更新计划状态（如果所有用例都已执行）
        if err := s.planRepo.UpdateStatusIfComplete(ctx, result.PlanID()); err != nil {
            return fmt.Errorf("update plan status: %w", err)
        }

        return nil
    })
}
```

---

## 7. 性能与内存考量

### 7.1 逃逸分析

```go
// ❌ 逃逸：返回指向栈变量的指针
func getUser() *User {
    u := User{Name: "test"} // u 逃逸到堆
    return &u
}

// ✅ 无逃逸：使用值返回
func getUser() User {
    return User{Name: "test"}
}

// ✅ 无逃逸：预分配切片
func processItems(items []Item) []Result {
    results := make([]Result, 0, len(items)) // 预分配
    for _, item := range items {
        results = append(results, Result{...})
    }
    return results
}
```

### 7.2 sync.Pool 对象复用

```go
// internal/pkg/pool/buffer_pool.go
package pool

import (
    "bytes"
    "sync"
)

var BufferPool = sync.Pool{
    New: func() any {
        return bytes.NewBuffer(make([]byte, 0, 1024))
    },
}

// GetBuffer 获取缓冲区
func GetBuffer() *bytes.Buffer {
    return BufferPool.Get().(*bytes.Buffer)
}

// PutBuffer 归还缓冲区
func PutBuffer(buf *bytes.Buffer) {
    buf.Reset()
    BufferPool.Put(buf)
}

// 使用示例
func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    buf := pool.GetBuffer()
    defer pool.PutBuffer(buf)

    // 使用 buf 处理请求...
    buf.Write(data)
}
```

---

## 8. 错误处理

### 8.1 错误包装规范

```go
// ✅ 正确：使用 %w 包装错误
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("find user by id %s: %w", id, err)
    }
    return user, nil
}

// ✅ 正确：领域错误直接返回
func (s *Service) Login(email, password string) (*Token, error) {
    user, err := s.repo.FindByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, identity.ErrUserNotFound) {
            return nil, identity.ErrUserNotFound // 领域错误
        }
        return nil, fmt.Errorf("find by email: %w", err)
    }

    if !user.VerifyPassword(password) {
        return nil, identity.ErrPasswordMismatch
    }

    return s.generateToken(user)
}
```

### 8.2 Panic 恢复机制

```go
// internal/transport/http/middleware/recovery.go
package middleware

import (
    "fmt"
    "net/http"
    "runtime/debug"

    "github.com/rs/zerolog/log"
    "github.com/liang21/aitestos/internal/ierrors"
)

func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rvr := recover(); rvr != nil {
                log.Error().
                    Str("path", r.URL.Path).
                    Str("method", r.Method).
                    Interface("panic", rvr).
                    Str("stack", string(debug.Stack())).
                    Msg("panic recovered")

                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusInternalServerError)

                resp := ierrors.NewErrorResponse(
                    ierrors.CodeInternalError,
                    r.Header.Get("X-Trace-ID"),
                )
                _, _ = w.Write(resp.ToJSON())
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

---

## 9. 可观测性设计

### 9.1 Prometheus 指标

```go
// internal/transport/http/middleware/metrics.go
package middleware

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "aitestos_http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "aitestos_http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{.01, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"method", "path"},
    )

    testCaseGenerated = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "aitestos_test_cases_generated_total",
            Help: "Total number of test cases generated by AI",
        },
        []string{"project_id", "status"},
    )
)

func Metrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // 包装 ResponseWriter 以获取状态码
        wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

        next.ServeHTTP(wrapped, r)

        duration := time.Since(start).Seconds()

        httpRequestsTotal.WithLabelValues(
            r.Method,
            r.URL.Path,
            fmt.Sprintf("%d", wrapped.status),
        ).Inc()

        httpRequestDuration.WithLabelValues(
            r.Method,
            r.URL.Path,
        ).Observe(duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    status int
}

func (w *responseWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}
```

### 9.2 结构化日志

```go
// internal/pkg/logger/logger.go
package logger

import (
    "context"

    "github.com/rs/zerolog"
)

type ctxKey struct{}

var loggerKey = ctxKey{}

// WithContext 向 context 注入 logger
func WithContext(ctx context.Context, log zerolog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, log)
}

// FromContext 从 context 获取 logger
func FromContext(ctx context.Context) zerolog.Logger {
    if log, ok := ctx.Value(loggerKey).(zerolog.Logger); ok {
        return log
    }
    return zerolog.Nop()
}

// 使用示例
func (s *Service) CreateUser(ctx context.Context, req *RegisterRequest) (*User, error) {
    log := logger.FromContext(ctx)

    log.Info().
        Str("email", req.Email).
        Str("username", req.Username).
        Msg("creating user")

    user, err := s.repo.Save(ctx, newUser)
    if err != nil {
        log.Error().Err(err).
            Str("email", req.Email).
            Msg("failed to create user")
        return nil, err
    }

    log.Info().
        Str("user_id", user.ID().String()).
        Msg("user created successfully")

    return user, nil
}
```

---

## 10. 单元测试方案

### 10.1 表格驱动测试

```go
// internal/domain/identity/user_test.go
package identity_test

import (
    "strings"
    "testing"

    "github.com/liang21/aitestos/internal/domain/identity"
)

func TestNewUser(t *testing.T) {
    tests := []struct {
        name        string
        username    string
        email       string
        password    string
        role        identity.UserRole
        wantErr     bool
        errContains string
    }{
        {
            name:     "valid user",
            username: "testuser",
            email:    "test@example.com",
            password: "password123",
            role:     identity.RoleNormal,
            wantErr:  false,
        },
        {
            name:        "invalid email",
            username:    "testuser",
            email:       "invalid-email",
            password:    "password123",
            role:        identity.RoleNormal,
            wantErr:     true,
            errContains: "invalid email",
        },
        {
            name:        "short username",
            username:    "ab",
            email:       "test@example.com",
            password:    "password123",
            role:        identity.RoleNormal,
            wantErr:     true,
            errContains: "username too short",
        },
        {
            name:        "short password",
            username:    "testuser",
            email:       "test@example.com",
            password:    "short",
            role:        identity.RoleNormal,
            wantErr:     true,
            errContains: "password too short",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := identity.NewUser(tt.username, tt.email, tt.password, tt.role)

            if tt.wantErr {
                if err == nil {
                    t.Errorf("NewUser() expected error, got nil")
                    return
                }
                if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
                    t.Errorf("NewUser() error = %v, want containing %v", err, tt.errContains)
                }
                return
            }

            if err != nil {
                t.Errorf("NewUser() unexpected error = %v", err)
                return
            }

            if user.Username() != tt.username {
                t.Errorf("Username() = %v, want %v", user.Username(), tt.username)
            }
            if user.Email() != tt.email {
                t.Errorf("Email() = %v, want %v", user.Email(), tt.email)
            }
            if user.Role() != tt.role {
                t.Errorf("Role() = %v, want %v", user.Role(), tt.role)
            }
        })
    }
}

func TestUser_VerifyPassword(t *testing.T) {
    tests := []struct {
        name     string
        password string
        input    string
        want     bool
    }{
        {"correct password", "password123", "password123", true},
        {"wrong password", "password123", "wrongpassword", false},
        {"empty password", "password123", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := identity.NewUser("test", "test@example.com", tt.password, identity.RoleNormal)
            if err != nil {
                t.Fatalf("NewUser() error = %v", err)
            }

            got := user.VerifyPassword(tt.input)
            if got != tt.want {
                t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 10.2 集成测试（使用 testcontainers）

```go
// tests/integration/repository/user_repo_test.go
package repository_test

import (
    "context"
    "testing"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"

    "github.com/liang21/aitestos/internal/domain/identity"
    "github.com/liang21/aitestos/internal/repository/identity"
)

func setupTestDB(t *testing.T) *sqlx.DB {
    ctx := context.Background()

    pgContainer, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("aitestos_test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(5*time.Second),
        ),
    )
    if err != nil {
        t.Fatalf("failed to start container: %v", err)
    }

    t.Cleanup(func() {
        if err := pgContainer.Terminate(ctx); err != nil {
            t.Logf("failed to terminate container: %v", err)
        }
    })

    connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    if err != nil {
        t.Fatalf("failed to get connection string: %v", err)
    }

    db, err := sqlx.Connect("postgres", connStr)
    if err != nil {
        t.Fatalf("failed to connect to database: %v", err)
    }

    // 运行迁移
    // migrateDatabase(t, db)

    return db
}

func TestUserRepository_Save(t *testing.T) {
    db := setupTestDB(t)
    repo := repository.NewUserRepository(db)

    ctx := context.Background()
    user, err := identity.NewUser("testuser", "test@example.com", "password123", identity.RoleNormal)
    if err != nil {
        t.Fatalf("NewUser() error = %v", err)
    }

    err = repo.Save(ctx, user)
    if err != nil {
        t.Errorf("Save() error = %v", err)
    }

    // 验证保存成功
    found, err := repo.FindByID(ctx, user.ID())
    if err != nil {
        t.Errorf("FindByID() error = %v", err)
    }

    if found.ID() != user.ID() {
        t.Errorf("FindByID() ID = %v, want %v", found.ID(), user.ID())
    }
}
```

---

## 11. 配置管理

### 11.1 配置结构

```go
// internal/config/config.go
package config

import (
    "fmt"
    "time"

    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
    Host            string        `mapstructure:"host"`
    Port            int           `mapstructure:"port"`
    ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
    Host            string        `mapstructure:"host"`
    Port            int           `mapstructure:"port"`
    User            string        `mapstructure:"user"`
    Password        string        `mapstructure:"password"`
    Database        string        `mapstructure:"database"`
    SSLMode         string        `mapstructure:"ssl_mode"`
    MaxOpenConns    int           `mapstructure:"max_open_conns"`
    MaxIdleConns    int           `mapstructure:"max_idle_conns"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
    Secret     string        `mapstructure:"secret"`
    ExpireTime time.Duration `mapstructure:"expire_time"`
}

type LogConfig struct {
    Level string `mapstructure:"level"`
    Json  bool   `mapstructure:"json"`
}

// Load 加载配置
func Load(path string) (*Config, error) {
    v := viper.New()
    v.SetConfigFile(path)
    v.AutomaticEnv()

    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unmarshal config: %w", err)
    }

    return &cfg, nil
}
```

### 11.2 配置示例

```yaml
# configs/config.example.yaml
server:
  host: "0.0.0.0"
  port: 8080
  shutdown_timeout: 30s

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  database: "aitestos"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key-change-in-production"
  expire_time: 2h

log:
  level: "info"
  json: true
```

---

## 12. Makefile

```makefile
# Makefile
.PHONY: all build test run clean lint fmt vet dev

APP_NAME := aitestos
BUILD_DIR := bin
MAIN_PATH := ./cmd/server

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOVET := $(GOCMD) vet
GOFMT := gofmt
GOLINT := golangci-lint

all: fmt vet lint test build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./tests/integration/...

run: build
	@echo "Running server..."
	./$(BUILD_DIR)/$(APP_NAME)

dev:
	@echo "Starting development environment..."
	docker-compose up -d
	@sleep 5
	$(MAKE) run

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

lint:
	@echo "Running linter..."
	$(GOLINT) run ./...

fmt:
	@echo "Formatting..."
	$(GOFMT) -s -w .

vet:
	@echo "Running vet..."
	$(GOVET) ./...

tidy:
	@echo "Tidying modules..."
	$(GOCMD) mod tidy

migrate-up:
	@echo "Running migrations..."
	@./scripts/migrate.sh up

migrate-down:
	@echo "Rolling back migrations..."
	@./scripts/migrate.sh down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f
```

---

## 13. 实施优先级

| 优先级 | 模块 | 说明 |
|---|---|---|
| P0 | `internal/domain/identity` | 用户认证是所有功能的基础 |
| P0 | `internal/domain/project` | 项目是数据隔离的基本单元 |
| P1 | `internal/domain/testcase` | 核心业务实体 |
| P1 | `internal/domain/testplan` | 测试执行流程 |
| P2 | `internal/domain/knowledge` | AI 能力支撑 |
| P2 | `internal/domain/generation` | 核心竞争力 |

---

## 14. Wire 依赖注入

### 14.1 Wire Provider Set 定义

```go
// internal/app/wire.go
//go:build wireinject

package app

import (
    "github.com/google/wire"
    "github.com/jmoiron/sqlx"

    "github.com/liang21/aitestos/internal/config"
    "github.com/liang21/aitestos/internal/domain/event"
    "github.com/liang21/aitestos/internal/repository/identity"
    "github.com/liang21/aitestos/internal/repository/project"
    "github.com/liang21/aitestos/internal/repository/testcase"
    "github.com/liang21/aitestos/internal/service/auth"
    projectService "github.com/liang21/aitestos/internal/service/project"
    testCaseService "github.com/liang21/aitestos/internal/service/testcase"
    "github.com/liang21/aitestos/internal/transport/http"
)

// InfrastructureSet 基础设施 Provider Set
var InfrastructureSet = wire.NewSet(
    NewDB,
    NewRedisClient,
    NewEventBus,
    NewCacheRepository,
    wire.Bind(new(event.EventBus), new(*RabbitMQEventBus)),
)

// RepositorySet Repository Provider Set
var RepositorySet = wire.NewSet(
    identity.NewUserRepository,
    project.NewProjectRepository,
    project.NewModuleRepository,
    project.NewProjectConfigRepository,
    testcase.NewTestCaseRepository,
    // ... 其他 Repository
)

// ServiceSet Service Provider Set
var ServiceSet = wire.NewSet(
    auth.NewAuthService,
    projectService.NewProjectService,
    testCaseService.NewTestCaseService,
    // ... 其他 Service
)

// HTTPSet HTTP Provider Set
var HTTPSet = wire.NewSet(
    http.NewServer,
    http.NewRouter,
    http.NewIdentityHandler,
    http.NewProjectHandler,
    http.NewTestCaseHandler,
    // ... 其他 Handler
)

// InitializeApp 初始化应用（Wire 生成）
func InitializeApp(cfg *config.Config) (*App, func(), error) {
    wire.Build(
        InfrastructureSet,
        RepositorySet,
        ServiceSet,
        HTTPSet,
        NewApp,
    )
    return nil, nil, nil
}
```

### 14.2 DB 和 Redis 初始化

```go
// internal/app/database.go
package app

import (
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/redis/go-redis/v9"

    "github.com/liang21/aitestos/internal/config"
)

func NewDB(cfg *config.Database) (*sqlx.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
    )

    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("connect database: %w", err)
    }

    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

    return db, nil
}

func NewRedisClient(cfg *config.Redis) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
    })
}
```

---

## 15. 缓存层设计

### 15.1 CacheRepository 接口

```go
// internal/domain/cache/repository.go
package cache

import (
    "context"
    "time"
)

// CacheRepository 缓存仓库接口
type CacheRepository interface {
    // Get 获取缓存
    Get(ctx context.Context, key string) ([]byte, error)
    // Set 设置缓存
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    // Delete 删除缓存
    Delete(ctx context.Context, key string) error
    // GetMulti 批量获取
    GetMulti(ctx context.Context, keys []string) (map[string][]byte, error)
    // SetMulti 批量设置
    SetMulti(ctx context.Context, items map[string][]byte, ttl time.Duration) error
}

// KeyBuilder 缓存键构建器
type KeyBuilder struct {
    prefix string
}

func NewKeyBuilder(prefix string) *KeyBuilder {
    return &KeyBuilder{prefix: prefix}
}

func (b *KeyBuilder) ProjectConfig(projectID string, key string) string {
    return fmt.Sprintf("%s:project:%s:config:%s", b.prefix, projectID, key)
}

func (b *KeyBuilder) User(userID string) string {
    return fmt.Sprintf("%s:user:%s", b.prefix, userID)
}

func (b *KeyBuilder) TestCase(caseID string) string {
    return fmt.Sprintf("%s:testcase:%s", b.prefix, caseID)
}
```

### 15.2 Redis 实现

```go
// internal/infrastructure/cache/redis_repo.go
package cache

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisCacheRepository struct {
    client *redis.Client
}

func NewRedisCacheRepository(client *redis.Client) *RedisCacheRepository {
    return &RedisCacheRepository{client: client}
}

func (r *RedisCacheRepository) Get(ctx context.Context, key string) ([]byte, error) {
    val, err := r.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return nil, nil // 缓存未命中
    }
    if err != nil {
        return nil, fmt.Errorf("redis get: %w", err)
    }
    return val, nil
}

func (r *RedisCacheRepository) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
        return fmt.Errorf("redis set: %w", err)
    }
    return nil
}

func (r *RedisCacheRepository) Delete(ctx context.Context, key string) error {
    if err := r.client.Del(ctx, key).Err(); err != nil {
        return fmt.Errorf("redis del: %w", err)
    }
    return nil
}

func (r *RedisCacheRepository) GetMulti(ctx context.Context, keys []string) (map[string][]byte, error) {
    pipe := r.client.Pipeline()
    cmds := make(map[string]*redis.StringCmd)

    for _, key := range keys {
        cmds[key] = pipe.Get(ctx, key)
    }

    _, err := pipe.Exec(ctx)
    if err != nil && err != redis.Nil {
        return nil, fmt.Errorf("redis pipeline: %w", err)
    }

    result := make(map[string][]byte)
    for key, cmd := range cmds {
        val, err := cmd.Bytes()
        if err == nil {
            result[key] = val
        }
    }

    return result, nil
}

func (r *RedisCacheRepository) SetMulti(ctx context.Context, items map[string][]byte, ttl time.Duration) error {
    pipe := r.client.Pipeline()

    for key, value := range items {
        pipe.Set(ctx, key, value, ttl)
    }

    _, err := pipe.Exec(ctx)
    if err != nil {
        return fmt.Errorf("redis pipeline: %w", err)
    }

    return nil
}
```

### 15.3 带缓存的 Repository 装饰器

```go
// internal/repository/project/cached_config_repo.go
package project

import (
    "context"
    "encoding/json"
    "time"

    "github.com/google/uuid"
    "github.com/liang21/aitestos/internal/domain/cache"
    domain "github.com/liang21/aitestos/internal/domain/project"
)

type CachedProjectConfigRepository struct {
    repo   domain.ProjectConfigRepository
    cache  cache.CacheRepository
    keys   *cache.KeyBuilder
    ttl    time.Duration
}

func NewCachedProjectConfigRepository(
    repo domain.ProjectConfigRepository,
    cache cache.CacheRepository,
    prefix string,
) *CachedProjectConfigRepository {
    return &CachedProjectConfigRepository{
        repo:  repo,
        cache: cache,
        keys:  cache.NewKeyBuilder(prefix),
        ttl:   5 * time.Minute,
    }
}

func (r *CachedProjectConfigRepository) FindByKey(
    ctx context.Context,
    projectID uuid.UUID,
    key string,
) (*domain.ProjectConfig, error) {
    cacheKey := r.keys.ProjectConfig(projectID.String(), key)

    // 尝试从缓存获取
    data, err := r.cache.Get(ctx, cacheKey)
    if err == nil && data != nil {
        var cfg domain.ProjectConfig
        if err := json.Unmarshal(data, &cfg); err == nil {
            return &cfg, nil
        }
    }

    // 缓存未命中，从数据库获取
    cfg, err := r.repo.FindByKey(ctx, projectID, key)
    if err != nil {
        return nil, err
    }

    // 写入缓存
    if data, err := json.Marshal(cfg); err == nil {
        _ = r.cache.Set(ctx, cacheKey, data, r.ttl)
    }

    return cfg, nil
}

func (r *CachedProjectConfigRepository) Save(
    ctx context.Context,
    cfg *domain.ProjectConfig,
) error {
    if err := r.repo.Save(ctx, cfg); err != nil {
        return err
    }

    // 删除缓存
    cacheKey := r.keys.ProjectConfig(cfg.ProjectID().String(), cfg.Key())
    _ = r.cache.Delete(ctx, cacheKey)

    return nil
}
```

---

## 16. 优雅停机设计

### 16.1 组件关闭顺序

```go
// internal/app/shutdown.go
package app

import (
    "context"
    "fmt"
    "time"

    "github.com/rs/zerolog/log"
)

type Closer interface {
    Name() string
    Close(ctx context.Context) error
}

type ShutdownManager struct {
    closers []Closer
    timeout time.Duration
}

func NewShutdownManager(timeout time.Duration) *ShutdownManager {
    return &ShutdownManager{
        closers: make([]Closer, 0),
        timeout: timeout,
    }
}

// Register 注册关闭器（按注册顺序的逆序关闭）
func (m *ShutdownManager) Register(closer Closer) {
    m.closers = append(m.closers, closer)
}

// Shutdown 执行优雅停机
// 关闭顺序: HTTP → gRPC → MQ Consumer → DB Pool
func (m *ShutdownManager) Shutdown(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, m.timeout)
    defer cancel()

    // 逆序关闭
    for i := len(m.closers) - 1; i >= 0; i-- {
        closer := m.closers[i]
        log.Info().Str("component", closer.Name()).Msg("shutting down")

        if err := closer.Close(ctx); err != nil {
            log.Error().Err(err).Str("component", closer.Name()).Msg("shutdown error")
            // 继续关闭其他组件，不中断
        } else {
            log.Info().Str("component", closer.Name()).Msg("shutdown complete")
        }
    }

    return nil
}

// HTTPServerCloser HTTP 服务器关闭器
type HTTPServerCloser struct {
    name string
    srv  interface {
        Shutdown(ctx context.Context) error
    }
}

func (c *HTTPServerCloser) Name() string { return c.name }
func (c *HTTPServerCloser) Close(ctx context.Context) error {
    return c.srv.Shutdown(ctx)
}

// MQConsumerCloser 消息队列消费者关闭器
type MQConsumerCloser struct {
    name string
    close func() error
}

func (c *MQConsumerCloser) Name() string { return c.name }
func (c *MQConsumerCloser) Close(ctx context.Context) error {
    return c.close()
}

// DBCloser 数据库连接池关闭器
type DBCloser struct {
    name string
    close func() error
}

func (c *DBCloser) Name() string { return c.name }
func (c *DBCloser) Close(ctx context.Context) error {
    return c.close()
}
```

### 16.2 完整的 App 生命周期管理

```go
// internal/app/app.go
package app

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/rs/zerolog/log"
)

type App struct {
    httpServer *http.Server
    shutdown   *ShutdownManager
}

func NewApp(
    httpServer *http.Server,
    shutdown *ShutdownManager,
) *App {
    return &App{
        httpServer: httpServer,
        shutdown: shutdown,
    }
}

func (a *App) Run(ctx context.Context) error {
    // 启动 HTTP 服务器
    errCh := make(chan error, 1)
    go func() {
        log.Info().Str("addr", a.httpServer.Addr).Msg("starting http server")
        if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            errCh <- err
        }
    }()

    // 监听系统信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    select {
    case err := <-errCh:
        return fmt.Errorf("http server error: %w", err)
    case sig := <-quit:
        log.Info().Str("signal", sig.String()).Msg("shutting down server")
    case <-ctx.Done():
        log.Info().Msg("context cancelled, shutting down")
    }

    // 执行优雅停机（按注册顺序的逆序关闭）
    return a.shutdown.Shutdown(context.Background())
}
```

---

## 17. Mock 测试方案

### 17.1 使用 mockgen 生成 Mock

```go
// internal/domain/identity/repository.go
package identity

//go:generate mockgen -source=repository.go -destination=../../test/mock/identity/mock_repository.go -package=mock_identity

import (
    "context"
    "github.com/google/uuid"
)

type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### 17.2 生成的 Mock 使用示例

```go
// internal/service/identity/auth_service_test.go
package identity_test

import (
    "context"
    "errors"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"

    mock_identity "github.com/liang21/aitestos/internal/test/mock/identity"
    "github.com/liang21/aitestos/internal/service/identity"
    domain "github.com/liang21/aitestos/internal/domain/identity"
)

func TestAuthService_Register(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // 创建 Mock Repository
    mockRepo := mock_identity.NewMockUserRepository(ctrl)

    // 创建 Service
    svc := identity.NewAuthService(mockRepo)

    tests := []struct {
        name      string
        req       *identity.RegisterRequest
        setupMock func()
        wantErr   bool
        errCode   int
    }{
        {
            name: "success",
            req: &identity.RegisterRequest{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
            },
            setupMock: func() {
                mockRepo.EXPECT().
                    FindByEmail(gomock.Any(), "test@example.com").
                    Return(nil, domain.ErrUserNotFound)
                mockRepo.EXPECT().
                    Save(gomock.Any(), gomock.Any()).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name: "email already exists",
            req: &identity.RegisterRequest{
                Username: "testuser",
                Email:    "existing@example.com",
                Password: "password123",
            },
            setupMock: func() {
                mockRepo.EXPECT().
                    FindByEmail(gomock.Any(), "existing@example.com").
                    Return(&domain.User{}, nil)
            },
            wantErr: true,
            errCode: 10002, // CodeEmailDuplicate
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setupMock()

            user, err := svc.Register(context.Background(), tt.req)

            if tt.wantErr {
                require.Error(t, err)
                var bizErr *ierrors.BizError
                if errors.As(err, &bizErr) {
                    assert.Equal(t, tt.errCode, bizErr.Code)
                }
            } else {
                require.NoError(t, err)
                assert.NotNil(t, user)
                assert.Equal(t, tt.req.Username, user.Username())
            }
        })
    }
}
```

### 17.3 Makefile 添加 Mock 生成

```makefile
# Makefile 追加

.PHONY: mock generate

# 生成所有 Mock
mock:
	@echo "Generating mocks..."
	@find ./internal/domain -name "*.go" -exec grep -l "go:generate mockgen" {} \; | while read f; do \
		dir=$$(dirname $$f); \
		go generate ./...; \
	done

# 或者直接运行
generate:
	@echo "Running go generate..."
	$(GOCMD) generate ./...
```

---

## 18. 验证方式

1. **编译检查**：`make build` 无错误
2. **静态检查**：`make lint` 无警告
3. **单元测试**：`make test` 覆盖率 > 80%
4. **集成测试**：`make test-integration` 通过
5. **API 测试**：使用 Postman/curl 验证 HTTP 接口
6. **性能测试**：使用 k6 进行负载测试
7. **Mock 测试**：`make mock` 生成 Mock，单元测试使用 Mock 隔离依赖
