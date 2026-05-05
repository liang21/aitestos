# 认证与数据库架构简化实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 简化认证机制（延长 access token 至 7 天）和数据库架构（删除软删除、前缀/缩写、用例编号，规范化表名）

**Architecture:** 
- 认证层：调整 access token 过期时间，保留 refresh token 机制
- 领域层：删除 ProjectPrefix、ModuleAbbreviation、CaseNumber 值对象
- 持久层：更新表名和查询，删除相关列
- API 层：删除相关字段和错误码

**Tech Stack:** Go 1.24+, PostgreSQL 16, golang-jwt/jwt/v5, chi router

---

## File Structure

### 新增文件
- `scripts/migrations/004_remove_prefix_deleted_at_rename_tables.sql` - 数据库迁移脚本

### 修改文件（领域层）
- `internal/domain/project/project.go` - 删除 prefix 字段
- `internal/domain/project/module.go` - 删除 abbreviation 字段
- `internal/domain/project/errors.go` - 删除相关错误
- `internal/domain/testcase/test_case.go` - 删除 number 字段
- `internal/domain/testcase/errors.go` - 删除相关错误
- `internal/domain/identity/user.go` - 确认无 deleted_at 依赖

### 删除文件（领域层）
- `internal/domain/project/prefix.go`
- `internal/domain/project/prefix_test.go`
- `internal/domain/project/abbreviation.go`（如果存在）
- `internal/domain/testcase/case_number.go`
- `internal/domain/testcase/case_number_test.go`

### 修改文件（服务层）
- `internal/service/identity/auth_service.go` - 更新 token 过期时间
- `internal/service/project/project_service.go` - 删除 prefix/abbreviation 处理
- `internal/service/testcase/case_service.go` - 删除编号生成

### 修改文件（持久层）
- `internal/repository/project/project_repo.go` - 表名和查询更新
- `internal/repository/project/module_repo.go` - 表名和查询更新
- `internal/repository/identity/user_repo.go` - 删除 deleted_at 查询
- `internal/repository/testcase/case_repo.go` - 删除 number 处理

### 修改文件（传输层）
- `internal/transport/http/handler/project.go` - 删除 prefix/abbreviation 字段
- `internal/transport/http/handler/testcase.go` - 删除 number 字段

### 修改文件（文档）
- `specs/aitestos_optimized.sql` - 数据库 schema
- `specs/openapi.yaml` - API 规范
- `specs/001-core-functionality/spec.md` - 功能规范
- `specs/001-core-functionality/plan.md` - 实现方案
- `specs/001-core-functionality/tasks.md` - 任务列表

---

## Phase 1: 文档更新

### Task 1: 更新数据库 Schema 文件

**Files:**
- Modify: `specs/aitestos_optimized.sql`

- [ ] **Step 1: 更新 users 表定义**

删除 `deleted_at` 列：

```sql
-- 在 specs/aitestos_optimized.sql 中找到 users 表定义（约第74-83行）
-- 删除 deleted_at 行

CREATE TABLE users (
    id         uuid                        DEFAULT gen_random_uuid()       NOT NULL PRIMARY KEY,
    username   varchar(32)                                                  NOT NULL,
    email      varchar(255)                                                 NOT NULL UNIQUE,
    password   varchar(255)                                                 NOT NULL,
    role       user_role_enum              DEFAULT 'normal'::user_role_enum NOT NULL,
    created_at timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
    -- deleted_at 已删除
);
```

- [ ] **Step 2: 更新 project 表名和定义**

重命名为 `projects`，删除 `prefix` 列：

```sql
-- 在 specs/aitestos_optimized.sql 中找到 project 表定义（约第30-38行）
-- 重命名表为 projects，删除 prefix 列和相关注释

CREATE TABLE projects (
    id          uuid                        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name        varchar(255)                                            NOT NULL UNIQUE,
    -- prefix 列已删除
    description text,
    config      jsonb                       DEFAULT '{}'::jsonb,
    created_at  timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);

-- 删除 prefix 相关的注释
-- COMMENT ON COLUMN projects.prefix ...  <-- 删除这行
```

- [ ] **Step 3: 更新 module 表名和定义**

重命名为 `modules`，删除 `abbreviation` 列：

```sql
-- 在 specs/aitestos_optimized.sql 中找到 module 表定义（约第59-67行）
-- 重命名表为 modules，删除 abbreviation 列和唯一约束

CREATE TABLE modules (
    id            uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    project_id    uuid                            NOT NULL REFERENCES projects ON DELETE CASCADE,
    name          varchar(255)                    NOT NULL,
    -- abbreviation 列已删除
    description   text,
    UNIQUE(project_id, name)
    -- UNIQUE(project_id, abbreviation) 已删除
);

-- 删除 abbreviation 相关的注释
-- COMMENT ON COLUMN modules.abbreviation ...  <-- 删除这行
```

- [ ] **Step 4: 更新外键引用**

更新其他表中引用 project/module 的外键：

```sql
-- 更新 project_config 表外键
CREATE TABLE project_config (
    -- ...
    project_id  uuid NOT NULL REFERENCES projects ON DELETE CASCADE,  -- project → projects
    -- ...
);

-- 更新 test_plan 表外键
CREATE TABLE test_plan (
    -- ...
    project_id   uuid NOT NULL REFERENCES projects ON DELETE CASCADE,  -- project → projects
    -- ...
);

-- 更新 module 表外键（在 test_case 中）
CREATE TABLE test_case (
    -- ...
    module_id     uuid NOT NULL REFERENCES modules ON DELETE CASCADE,  -- module → modules
    -- ...
);

-- 更新 document 表外键
CREATE TABLE document (
    -- ...
    project_id   uuid NOT NULL REFERENCES projects ON DELETE CASCADE,  -- project → projects
    -- ...
);

-- 更新 generation_task 表外键
CREATE TABLE generation_task (
    -- ...
    project_id     uuid NOT NULL REFERENCES projects ON DELETE CASCADE,  -- project → projects
    -- ...
);
```

- [ ] **Step 5: 更新 test_case 表定义**

删除 `number` 列：

```sql
-- 在 specs/aitestos_optimized.sql 中找到 test_case 表定义（约第88-103行）
-- 删除 number 列

CREATE TABLE test_case (
    id            uuid                        DEFAULT gen_random_uuid()              NOT NULL PRIMARY KEY,
    module_id     uuid                                                                NOT NULL REFERENCES modules ON DELETE CASCADE,
    user_id       uuid                                                                NOT NULL REFERENCES users ON DELETE NO ACTION,
    -- number 列已删除
    title         varchar(255)                                                        NOT NULL,
    preconditions jsonb                       DEFAULT '[]'::jsonb,
    steps         jsonb                       DEFAULT '[]'::jsonb                     NOT NULL,
    expected      jsonb                       DEFAULT '{}'::jsonb                     NOT NULL,
    ai_metadata   jsonb                       DEFAULT '{}'::jsonb,
    case_type     case_type_enum              DEFAULT 'functionality'::case_type_enum NOT NULL,
    priority      priority_enum               DEFAULT 'P2'::priority_enum             NOT NULL,
    status        case_status_enum            DEFAULT 'unexecuted'::case_status_enum  NOT NULL,
    created_at    timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);
```

- [ ] **Step 6: 删除相关索引**

删除与删除列相关的索引：

```sql
-- 在 specs/aitestos_optimized.sql 中找到索引定义部分
-- 删除以下索引：

-- DROP INDEX IF EXISTS idx_users_deleted_at;  -- 删除这行（约在第212行）

-- 注释：number 列的索引也会自动删除
-- DROP INDEX IF EXISTS idx_test_case_number;  -- 确认此索引被删除
```

- [ ] **Step 7: 删除相关注释**

删除与删除列相关的 COMMENT：

```sql
-- 在 specs/aitestos_optimized.sql 中找到注释部分（第306-328行）
-- 删除以下注释：

-- COMMENT ON COLUMN projects.prefix IS '项目前缀，用于用例编号生成，2-4位大写字母';  -- 删除
-- COMMENT ON COLUMN modules.abbreviation IS '模块缩写，用于用例编号生成，2-4位大写字母，项目内唯一';  -- 删除
```

- [ ] **Step 8: 验证文件语法**

运行：`psql --set=ON_ERROR_STOP=1 -f specs/aitestos_optimized.sql`（可选，如果有 PostgreSQL）
或仅检查 SQL 语法是否正确

- [ ] **Step 9: 提交文档更新**

```bash
git add specs/aitestos_optimized.sql
git commit -m "docs(sql): remove deleted_at, prefix, abbreviation, number; rename tables

- Remove deleted_at from users table
- Remove prefix from project table, rename to projects
- Remove abbreviation from module table, rename to modules  
- Remove number from test_case table
- Update all foreign key references

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 2: 更新 OpenAPI 规范

**Files:**
- Modify: `specs/openapi.yaml`

- [ ] **Step 1: 删除 Project schema 中的 prefix 字段**

在 `components.schemas.Project` 中删除 prefix：

```yaml
# 在 specs/openapi.yaml 中找到 Project schema（约第1586行）
Project:
  type: object
  properties:
    id:
      type: string
      format: uuid
    name:
      type: string
    # prefix:    <-- 删除这整个属性块
    #   type: string
    #   description: 项目前缀，用于用例编号生成
    description:
      type: string
    createdAt:
      type: string
      format: date-time
    updatedAt:
      type: string
      format: date-time
```

- [ ] **Step 2: 删除 ProjectDetail 中的 prefix 字段**

```yaml
# 在 specs/openapi.yaml 中找到 ProjectDetail schema（约第1606行）
ProjectDetail:
  type: object
  properties:
    id:
      type: string
      format: uuid
    name:
      type: string
    # prefix:    <-- 删除
    description:
      type: string
    # ... 其他字段
```

- [ ] **Step 3: 删除 CreateProjectRequest 中的 prefix 字段**

```yaml
# 在 specs/openapi.yaml 中找到 CreateProjectRequest schema（约第1648行）
CreateProjectRequest:
  type: object
  required: [name]  # 从 [name, prefix] 改为 [name]
  properties:
    name:
      type: string
      minLength: 2
      maxLength: 255
    # prefix:    <-- 删除整个属性块
    #   type: string
    #   minLength: 2
    #   maxLength: 4
    #   pattern: '^[A-Z]+$'
    #   description: 项目前缀，2-4位大写字母，用于用例编号生成
    description:
      type: string
```

- [ ] **Step 4: 删除 Module schema 中的 abbreviation 字段**

```yaml
# 在 specs/openapi.yaml 中找到 Module schema（约第1676行）
Module:
  type: object
  properties:
    id:
      type: string
      format: uuid
    projectId:
      type: string
      format: uuid
    name:
      type: string
    # abbreviation:    <-- 删除整个属性块
    #   type: string
    #   description: 模块缩写，用于用例编号生成
    description:
      type: string
    # ... 其他字段
```

- [ ] **Step 5: 删除 CreateModuleRequest 中的 abbreviation 字段**

```yaml
# 在 specs/openapi.yaml 中找到 CreateModuleRequest schema（约第1702行）
CreateModuleRequest:
  type: object
  required: [name]  # 从 [name, abbreviation] 改为 [name]
  properties:
    name:
      type: string
      minLength: 2
      maxLength: 255
    # abbreviation:    <-- 删除整个属性块
    #   type: string
    #   minLength: 2
    #   maxLength: 4
    #   pattern: '^[A-Z]+$'
    #   description: 模块缩写，2-4位大写字母，项目内唯一
    description:
      type: string
```

- [ ] **Step 6: 删除 UpdateModuleRequest 中的 abbreviation 字段**

```yaml
# 在 specs/openapi.yaml 中找到 UpdateModuleRequest schema（约第1719行）
UpdateModuleRequest:
  type: object
  properties:
    name:
      type: string
      minLength: 2
      maxLength: 255
    # abbreviation:    <-- 删除整个属性块
    #   type: string
    #   minLength: 2
    #   maxLength: 4
    #   pattern: '^[A-Z]+$'
    description:
      type: string
```

- [ ] **Step 7: 删除 TestCase schema 中的 number 字段**

```yaml
# 在 specs/openapi.yaml 中找到 TestCase schema（约第1757行）
TestCase:
  type: object
  properties:
    id:
      type: string
      format: uuid
    moduleId:
      type: string
      format: uuid
    userId:
      type: string
      format: uuid
    # number:    <-- 删除整个属性块
    #   type: string
    #   description: 业务编号，如 PRJ-MOD-20260401-001
    title:
      type: string
    # ... 其他字段
```

- [ ] **Step 8: 删除 CaseDetail 中的 number 字段**

```yaml
# 在 specs/openapi.yaml 中找到 CaseDetail schema（约第1828行）
CaseDetail:
  type: object
  description: 测试用例详情（含关联信息）
  properties:
    id:
      type: string
      format: uuid
    moduleId:
      type: string
      format: uuid
    # ... 其他字段
    # number:    <-- 删除
    # project_prefix:    <-- 删除（如果存在）
    title:
      type: string
    # ... 其他字段
```

- [ ] **Step 9: 删除相关错误码**

在错误响应部分删除以下错误码（如果有）：

```yaml
# 检查文档中是否有以下错误码定义，如果有则删除：
# 20007: 项目前缀已存在
# 20008: 项目前缀格式无效  
# 20009: 模块缩写已存在
# 40002: 用例编号已存在
```

- [ ] **Step 10: 验证 YAML 语法**

运行：`yamllint specs/openapi.yaml` 或 `python -c "import yaml; yaml.safe_load(open('specs/openapi.yaml'))"`

- [ ] **Step 11: 提交变更**

```bash
git add specs/openapi.yaml
git commit -m "docs(openapi): remove prefix, abbreviation, number fields

- Remove prefix from Project schemas
- Remove abbreviation from Module schemas
- Remove number from TestCase schemas
- Update required fields accordingly

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 3: 更新功能规范文档

**Files:**
- Modify: `specs/001-core-functionality/spec.md`

- [ ] **Step 1: 删除用户故事中的前缀/缩写需求**

在 US-1.1 和 US-1.2 中删除相关验收标准：

```markdown
# 在 specs/001-core-functionality/spec.md 中找到 US-1.1 创建项目（约第14行）

#### US-1.1 创建项目

```
作为 超级管理员
我想要 创建新项目并设置项目前缀  <-- 修改描述
以便于 隔离不同产品的测试资产

验收标准：
- 输入项目名称（必填，全局唯一）
- 输入项目前缀（必填，2-4位大写字母，全局唯一）  <-- 删除这行
- 输入项目描述（可选）
- 创建成功后跳转到项目详情页
- 项目名称或前缀重复时返回错误码 20002/20007  <-- 修改为：项目名称重复时返回错误码 20002
```

# 修改为：

#### US-1.1 创建项目

```
作为 超级管理员
我想要 创建新项目  <-- 删除"并设置项目前缀"
以便于 隔离不同产品的测试资产

验收标准：
- 输入项目名称（必填，全局唯一）
- 输入项目描述（可选）
- 创建成功后跳转到项目详情页
- 项目名称重复时返回错误码 20002
```
```

- [ ] **Step 2: 删除模块缩写相关验收标准**

```markdown
# 在 specs/001-core-functionality/spec.md 中找到 US-1.2 管理模块（约第29行）

#### US-1.2 管理模块

```
作为 管理员
我想要 在项目下创建功能模块并设置模块缩写  <-- 修改描述
以便于 按功能组织测试用例

验收标准：
- 模块名称在项目内唯一
- 模块缩写在项目内唯一（2-4位大写字母）  <-- 删除这行
- 删除模块时提示将删除关联用例
- 返回错误码 20004/20009 表示重复  <-- 修改为：返回错误码 20004 表示重复
```

# 修改为：

#### US-1.2 管理模块

```
作为 管理员
我想要 在项目下创建功能模块  <-- 删除"并设置模块缩写"
以便于 按功能组织测试用例

验收标准：
- 模块名称在项目内唯一
- 删除模块时提示将删除关联用例
- 返回错误码 20004 表示名称重复
```
```

- [ ] **Step 3: 删除用例编号相关内容**

在 US-3.3 中删除编号生成相关内容：

```markdown
# 在 specs/001-core-functionality/spec.md 中找到 US-3.3 确认草稿（约第125行）

#### US-3.3 确认草稿

```
作为 测试工程师
我想要 确认采用某个草稿
以便于 将其转为正式测试用例

验收标准：
- 确认前可编辑草稿内容
- 确认时自动生成用例编号  <-- 删除这行
- 用例编号格式：{项目前缀}-{模块缩写}-{日期}-{序号}  <-- 删除这行
- 草稿状态变为 confirmed
```

# 修改为：

#### US-3.3 确认草稿

```
作为 测试工程师
我想要 确认采用某个草稿
以便于 将其转为正式测试用例

验收标准：
- 确认前可编辑草稿内容
- 草稿状态变为 confirmed
```
```

- [ ] **Step 4: 删除用例编号格式定义**

删除第413-422行的用例编号规则：

```markdown
# 在 specs/001-core-functionality/spec.md 中找到"用例编号规则"部分（约第412行）

#### 2.4.2 用例编号规则  <-- 删除整个章节

**格式**：`{项目前缀}-{模块缩写}-{日期}-{序号}`

**示例**：`ECO-USR-20260402-001`

**组成**：
- 项目前缀：创建项目时输入，2-4 位大写字母，全局唯一
- 模块缩写：创建模块时输入，2-4 位大写字母，项目内唯一
- 日期：YYYYMMDD 格式
- 序号：当日该模块下 001-999 递增

# 完全删除此章节
```

- [ ] **Step 5: 删除错误码表中的相关错误码**

在错误码速查表中删除：

```markdown
# 在 specs/001-core-functionality/spec.md 中找到错误码速查表（约第762行）

| 错误码 | 说明 | HTTP 状态码 |
|--------|------|-------------|
# 删除以下行：
# 20007 | 项目前缀已存在 | 409
# 20008 | 项目前缀格式无效 | 400
# 20009 | 模块缩写已存在 | 409
# 40002 | 用例编号已存在 | 409
```

- [ ] **Step 6: 更新输出格式示例**

删除相关字段：

```markdown
# 在 specs/001-core-functionality/spec.md 中找到"输出格式示例"部分

# 删除或更新以下示例中的 prefix/abbreviation/number 字段：

### 5.1 项目创建响应
{
  "id": "...",
  "name": "电商平台",
  # "prefix": "ECO",  <-- 删除
  "description": "..."
}

### 5.2 模块创建响应
{
  "id": "...",
  "projectId": "...",
  "name": "用户中心",
  # "abbreviation": "USR",  <-- 删除
  "description": "..."
}

### 5.6 测试用例响应
{
  "id": "...",
  "moduleId": "...",
  # "number": "ECO-USR-20260402-001",  <-- 删除
  "title": "密码错误5次后账号锁定",
  # ...
}
```

- [ ] **Step 7: 更新安全要求中的 token 过期时间**

```markdown
# 在 specs/001-core-functionality/spec.md 中找到安全要求部分（约第504行）

### 3.4 安全要求

| 要求 | 实现方式 |
|------|----------|
| 认证 | JWT Token，有效期 7 天 |  <-- 从"2 小时"改为"7 天"
# ... 其他要求保持不变
```

- [ ] **Step 8: 提交变更**

```bash
git add specs/001-core-functionality/spec.md
git commit -m "docs(spec): remove prefix/abbreviation/number requirements

- Remove project prefix from user stories
- Remove module abbreviation from user stories
- Remove case number format and generation
- Update token expiry to 7 days
- Remove related error codes

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Phase 2: 数据库迁移脚本

### Task 4: 创建数据库迁移脚本

**Files:**
- Create: `scripts/migrations/004_remove_prefix_deleted_at_rename_tables.sql`

- [ ] **Step 1: 创建迁移文件头部**

```sql
-- 创建文件：scripts/migrations/004_remove_prefix_deleted_at_rename_tables.sql

-- ============================================================
-- Migration 004: Remove prefix/deleted_at, rename tables
-- ============================================================
-- Description:
--   - Remove deleted_at column from users table
--   - Remove prefix column from project table
--   - Remove abbreviation column from module table
--   - Remove number column from test_case table
--   - Rename project → projects
--   - Rename module → modules
--   - Drop related indexes and constraints
-- ============================================================

BEGIN;
```

- [ ] **Step 2: 添加 users 表变更**

```sql
-- 1. Remove deleted_at from users table
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
```

- [ ] **Step 3: 添加 project 表变更**

```sql
-- 2. Remove prefix from project table
ALTER TABLE project DROP COLUMN IF EXISTS prefix;
```

- [ ] **Step 4: 添加 module 表变更**

```sql
-- 3. Remove abbreviation from module table
-- First drop the unique constraint on abbreviation
ALTER TABLE module DROP CONSTRAINT IF EXISTS module_abbreviation_key;
-- Then drop the column
ALTER TABLE module DROP COLUMN IF EXISTS abbreviation;
```

- [ ] **Step 5: 添加 test_case 表变更**

```sql
-- 4. Remove number from test_case table
-- First drop the unique constraint on number
ALTER TABLE test_case DROP CONSTRAINT IF EXISTS test_case_number_key;
-- Then drop the column
ALTER TABLE test_case DROP COLUMN IF EXISTS number;
```

- [ ] **Step 6: 添加表重命名**

```sql
-- 5. Rename tables
ALTER TABLE project RENAME TO projects;
ALTER TABLE module RENAME TO modules;
```

- [ ] **Step 7: 添加索引清理**

```sql
-- 6. Drop obsolete indexes
DROP INDEX IF EXISTS idx_users_deleted_at;
```

- [ ] **Step 8: 添加事务提交**

```sql
COMMIT;
```

- [ ] **Step 9: 验证脚本语法**

运行：`psql --set=ON_ERROR_STOP=1 -f scripts/migrations/004_remove_prefix_deleted_at_rename_tables.sql`（可选，如有数据库）

- [ ] **Step 10: 提交迁移脚本**

```bash
git add scripts/migrations/004_remove_prefix_deleted_at_rename_tables.sql
git commit -m "feat(migration): remove prefix/deleted_at, rename tables

- Remove deleted_at from users (hard delete instead of soft delete)
- Remove prefix from projects
- Remove abbreviation from modules
- Remove number from test_case (UUID only)
- Rename: project → projects, module → modules

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Phase 3: 领域层代码变更

### Task 5: 删除 ProjectPrefix 值对象

**Files:**
- Delete: `internal/domain/project/prefix.go`
- Delete: `internal/domain/project/prefix_test.go`

- [ ] **Step 1: 确认文件存在**

运行：`ls -la internal/domain/project/prefix.go internal/domain/project/prefix_test.go`

- [ ] **Step 2: 删除 prefix.go 文件**

运行：`rm internal/domain/project/prefix.go`

- [ ] **Step 3: 删除 prefix_test.go 文件**

运行：`rm internal/domain/project/prefix_test.go`

- [ ] **Step 4: 提交删除**

```bash
git add internal/domain/project/prefix.go internal/domain/project/prefix_test.go
git commit -m "refactor(domain): remove ProjectPrefix value object

Project prefix concept is removed as part of simplification.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 6: 更新 Project 聚合根

**Files:**
- Modify: `internal/domain/project/project.go`
- Modify: `internal/domain/project/project_test.go`

- [ ] **Step 1: 读取当前 Project 定义**

运行：`cat internal/domain/project/project.go`

- [ ] **Step 2: 更新 Project 结构体**

```go
// 在 internal/domain/project/project.go 中

// 变更前：
type Project struct {
    id          uuid.UUID
    name        string
    prefix      ProjectPrefix    // 删除这行
    description string
    config      map[string]any
    createdAt   time.Time
    updatedAt   time.Time
}

// 变更后：
type Project struct {
    id          uuid.UUID
    name        string
    // prefix 字段已删除
    description string
    config      map[string]any
    createdAt   time.Time
    updatedAt   time.Time
}
```

- [ ] **Step 3: 更新 NewProject 工厂函数**

```go
// 在 internal/domain/project/project.go 中

// 变更前：
func NewProject(name, prefixStr, description string) (*Project, error) {
    prefix, err := ParseProjectPrefix(prefixStr)
    if err != nil {
        return nil, err
    }
    
    if name == "" {
        return nil, ErrProjectNameRequired
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

// 变更后：
func NewProject(name, description string) (*Project, error) {
    if name == "" {
        return nil, ErrProjectNameRequired
    }
    
    now := time.Now()
    return &Project{
        id:          uuid.New(),
        name:        name,
        description: description,
        config:      make(map[string]any),
        createdAt:   now,
        updatedAt:   now,
    }, nil
}
```

- [ ] **Step 4: 更新访问器方法**

```go
// 在 internal/domain/project/project.go 中

// 删除以下方法：
// func (p *Project) Prefix() ProjectPrefix {
//     return p.prefix
// }
```

- [ ] **Step 5: 更新序列化方法（如果存在）**

```go
// 在 internal/domain/project/project.go 中检查是否有 MarshalJSON/UnmarshalJSON
// 如果有，更新以移除 prefix 字段处理

// 变更前：
func (p *Project) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        ID          uuid.UUID        `json:"id"`
        Name        string           `json:"name"`
        Prefix      ProjectPrefix    `json:"prefix"`
        Description string           `json:"description"`
        Config      map[string]any   `json:"config"`
        CreatedAt   time.Time        `json:"createdAt"`
        UpdatedAt   time.Time        `json:"updatedAt"`
    }{
        ID:          p.id,
        Name:        p.name,
        Prefix:      p.prefix,
        Description: p.description,
        Config:      p.config,
        CreatedAt:   p.createdAt,
        UpdatedAt:   p.updatedAt,
    })
}

// 变更后：
func (p *Project) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        ID          uuid.UUID        `json:"id"`
        Name        string           `json:"name"`
        Description string           `json:"description"`
        Config      map[string]any   `json:"config"`
        CreatedAt   time.Time        `json:"createdAt"`
        UpdatedAt   time.Time        `json:"updatedAt"`
    }{
        ID:          p.id,
        Name:        p.name,
        Description: p.description,
        Config:      p.config,
        CreatedAt:   p.createdAt,
        UpdatedAt:   p.updatedAt,
    })
}
```

- [ ] **Step 6: 更新测试文件**

```go
// 在 internal/domain/project/project_test.go 中

// 删除所有与 prefix 相关的测试用例

// 删除类似这样的测试：
// func TestNewProject_InvalidPrefix(t *testing.T) {
//     _, err := NewProject("Test", "ABCD123", "desc")
//     assert.Error(t, err)
// }
```

- [ ] **Step 7: 运行项目领域测试**

运行：`go test ./internal/domain/project/... -v`

预期：所有测试通过（失败的测试应该已删除）

- [ ] **Step 8: 提交变更**

```bash
git add internal/domain/project/project.go internal/domain/project/project_test.go
git commit -m "refactor(domain): remove prefix from Project aggregate

- Remove prefix field from Project struct
- Update NewProject to remove prefix parameter
- Remove Prefix() accessor method
- Update JSON marshaling to exclude prefix
- Remove prefix-related tests

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 7: 更新 Module 实体

**Files:**
- Modify: `internal/domain/project/module.go`
- Modify: `internal/domain/project/module_test.go`

- [ ] **Step 1: 读取当前 Module 定义**

运行：`cat internal/domain/project/module.go`

- [ ] **Step 2: 更新 Module 结构体**

```go
// 在 internal/domain/project/module.go 中

// 变更前：
type Module struct {
    id            uuid.UUID
    projectID     uuid.UUID
    name          string
    abbreviation  ModuleAbbreviation  // 删除这行
    description   string
    createdAt     time.Time
    updatedAt     time.Time
}

// 变更后：
type Module struct {
    id            uuid.UUID
    projectID     uuid.UUID
    name          string
    // abbreviation 字段已删除
    description   string
    createdAt     time.Time
    updatedAt     time.Time
}
```

- [ ] **Step 3: 更新 NewModule 工厂函数**

```go
// 在 internal/domain/project/module.go 中

// 变更前：
func NewModule(projectID uuid.UUID, name, abbrevStr, description string) (*Module, error) {
    abbrev, err := ParseModuleAbbreviation(abbrevStr)
    if err != nil {
        return nil, err
    }
    
    if name == "" {
        return nil, ErrModuleNameRequired
    }
    
    now := time.Now()
    return &Module{
        id:           uuid.New(),
        projectID:    projectID,
        name:         name,
        abbreviation: abbrev,
        description:  description,
        createdAt:    now,
        updatedAt:    now,
    }, nil
}

// 变更后：
func NewModule(projectID uuid.UUID, name, description string) (*Module, error) {
    if name == "" {
        return nil, ErrModuleNameRequired
    }
    
    now := time.Now()
    return &Module{
        id:          uuid.New(),
        projectID:   projectID,
        name:        name,
        description: description,
        createdAt:   now,
        updatedAt:   now,
    }, nil
}
```

- [ ] **Step 4: 删除访问器方法**

```go
// 在 internal/domain/project/module.go 中

// 删除以下方法：
// func (m *Module) Abbreviation() ModuleAbbreviation {
//     return m.abbreviation
// }
```

- [ ] **Step 5: 更新测试文件**

```go
// 在 internal/domain/project/module_test.go 中

// 删除所有与 abbreviation 相关的测试用例

// 删除类似这样的测试：
// func TestNewModule_InvalidAbbreviation(t *testing.T) {
//     _, err := NewModule(uuid.New(), "Test", "ABCD123", "desc")
//     assert.Error(t, err)
// }
```

- [ ] **Step 6: 运行模块领域测试**

运行：`go test ./internal/domain/project/... -v`

预期：所有测试通过

- [ ] **Step 7: 提交变更**

```bash
git add internal/domain/project/module.go internal/domain/project/module_test.go
git commit -m "refactor(domain): remove abbreviation from Module entity

- Remove abbreviation field from Module struct
- Update NewModule to remove abbreviation parameter
- Remove Abbreviation() accessor method
- Remove abbreviation-related tests

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 8: 更新项目领域错误定义

**Files:**
- Modify: `internal/domain/project/errors.go`

- [ ] **Step 1: 读取当前错误定义**

运行：`cat internal/domain/project/errors.go`

- [ ] **Step 2: 删除前缀相关错误**

```go
// 在 internal/domain/project/errors.go 中

// 删除以下错误定义：
var (
    // ErrInvalidProjectPrefix   = errors.New("invalid project prefix: must be 2-4 uppercase letters")
    // ErrProjectPrefixDuplicate = errors.New("project prefix already exists")
)
```

- [ ] **Step 3: 删除缩写相关错误**

```go
// 在 internal/domain/project/errors.go 中

// 删除以下错误定义：
var (
    // ErrInvalidModuleAbbrev   = errors.New("invalid module abbreviation: must be 2-4 uppercase letters")
    // ErrModuleAbbrevDuplicate = errors.New("module abbreviation already exists in project")
)
```

- [ ] **Step 4: 运行项目领域测试**

运行：`go test ./internal/domain/project/... -v`

- [ ] **Step 5: 提交变更**

```bash
git add internal/domain/project/errors.go
git commit -m "refactor(domain): remove prefix/abbreviation errors

- Remove ErrInvalidProjectPrefix
- Remove ErrProjectPrefixDuplicate
- Remove ErrInvalidModuleAbbrev
- Remove ErrModuleAbbrevDuplicate

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 9: 删除 CaseNumber 值对象

**Files:**
- Delete: `internal/domain/testcase/case_number.go`
- Delete: `internal/domain/testcase/case_number_test.go`

- [ ] **Step 1: 确认文件存在**

运行：`ls -la internal/domain/testcase/case_number.go internal/domain/testcase/case_number_test.go`

- [ ] **Step 2: 删除 case_number.go**

运行：`rm internal/domain/testcase/case_number.go`

- [ ] **Step 3: 删除 case_number_test.go**

运行：`rm internal/domain/testcase/case_number_test.go`

- [ ] **Step 4: 提交删除**

```bash
git add internal/domain/testcase/case_number.go internal/domain/testcase/case_number_test.go
git commit -m "refactor(domain): remove CaseNumber value object

Case numbering is removed; UUIDs are used for identification.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 10: 更新 TestCase 聚合根

**Files:**
- Modify: `internal/domain/testcase/test_case.go`
- Modify: `internal/domain/testcase/test_case_test.go`

- [ ] **Step 1: 读取当前 TestCase 定义**

运行：`cat internal/domain/testcase/test_case.go`

- [ ] **Step 2: 更新 TestCase 结构体**

```go
// 在 internal/domain/testcase/test_case.go 中

// 变更前：
type TestCase struct {
    id            uuid.UUID
    moduleID      uuid.UUID
    userID        uuid.UUID
    number        CaseNumber          // 删除这行
    title         string
    preconditions Preconditions
    steps         Steps
    expected      ExpectedResult
    aiMetadata    *AiMetadata
    caseType      CaseType
    priority      Priority
    status        CaseStatus
    createdAt     time.Time
    updatedAt     time.Time
}

// 变更后：
type TestCase struct {
    id            uuid.UUID
    moduleID      uuid.UUID
    userID        uuid.UUID
    // number 字段已删除
    title         string
    preconditions Preconditions
    steps         Steps
    expected      ExpectedResult
    aiMetadata    *AiMetadata
    caseType      CaseType
    priority      Priority
    status        CaseStatus
    createdAt     time.Time
    updatedAt     time.Time
}
```

- [ ] **Step 3: 更新 NewTestCase 工厂函数**

```go
// 在 internal/domain/testcase/test_case.go 中

// 变更前：
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
    // ...
}

// 变更后：
func NewTestCase(
    moduleID, userID uuid.UUID,
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
```

- [ ] **Step 4: 删除访问器方法**

```go
// 在 internal/domain/testcase/test_case.go 中

// 删除以下方法：
// func (tc *TestCase) Number() CaseNumber {
//     return tc.number
// }
```

- [ ] **Step 5: 更新测试文件**

```go
// 在 internal/domain/testcase/test_case_test.go 中

// 更新所有调用 NewTestCase 的测试

// 变更前：
func TestNewTestCase_Success(t *testing.T) {
    tc, err := NewTestCase(
        moduleID, userID,
        testcase.CaseNumber("TEST-001"),
        "Test Title",
        testcase.Preconditions{"pre1"},
        testcase.Steps{"step1"},
        testcase.ExpectedResult{"result": "expected"},
        testcase.CaseTypeFunctionality,
        testcase.PriorityP1,
    )
    // ...
}

// 变更后：
func TestNewTestCase_Success(t *testing.T) {
    tc, err := NewTestCase(
        moduleID, userID,
        "Test Title",
        testcase.Preconditions{"pre1"},
        testcase.Steps{"step1"},
        testcase.ExpectedResult{"result": "expected"},
        testcase.CaseTypeFunctionality,
        testcase.PriorityP1,
    )
    // ...
}

// 删除所有与 number 相关的测试用例
```

- [ ] **Step 6: 运行测试用例领域测试**

运行：`go test ./internal/domain/testcase/... -v`

预期：所有测试通过

- [ ] **Step 7: 提交变更**

```bash
git add internal/domain/testcase/test_case.go internal/domain/testcase/test_case_test.go
git commit -m "refactor(domain): remove number from TestCase aggregate

- Remove number field from TestCase struct
- Update NewTestCase to remove number parameter
- Remove Number() accessor method
- Update all tests to remove number parameter

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 11: 更新测试用例领域错误定义

**Files:**
- Modify: `internal/domain/testcase/errors.go`

- [ ] **Step 1: 读取当前错误定义**

运行：`cat internal/domain/testcase/errors.go`

- [ ] **Step 2: 删除编号相关错误**

```go
// 在 internal/domain/testcase/errors.go 中

// 删除以下错误定义：
var (
    // ErrInvalidCaseNumber  = errors.New("invalid case number format")
    // ErrCaseNumberDuplicate = errors.New("case number already exists")
)
```

- [ ] **Step 3: 运行测试用例领域测试**

运行：`go test ./internal/domain/testcase/... -v`

- [ ] **Step 4: 提交变更**

```bash
git add internal/domain/testcase/errors.go
git commit -m "refactor(domain): remove case number errors

- Remove ErrInvalidCaseNumber
- Remove ErrCaseNumberDuplicate

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Phase 4: 服务层代码变更

### Task 12: 更新认证服务 Token 过期时间

**Files:**
- Modify: `internal/service/identity/auth_service.go`

- [ ] **Step 1: 读取当前认证服务**

运行：`cat internal/service/identity/auth_service.go`

- [ ] **Step 2: 更新 access token 过期时间**

```go
// 在 internal/service/identity/auth_service.go 中

// 变更前：
const accessTokenExpiry = 15 * time.Minute

// 变更后：
const accessTokenExpiry = 7 * 24 * time.Hour  // 7 天
```

- [ ] **Step 3: 运行认证服务测试**

运行：`go test ./internal/service/identity/... -v`

预期：所有测试通过

- [ ] **Step 4: 提交变更**

```bash
git add internal/service/identity/auth_service.go
git commit -m "feat(auth): extend access token expiry to 7 days

- Change access token expiry from 15 minutes to 7 days
- Refresh token mechanism remains unchanged

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 13: 更新项目服务

**Files:**
- Modify: `internal/service/project/project_service.go`
- Modify: `internal/service/project/project_service_test.go`

- [ ] **Step 1: 读取当前项目服务**

运行：`cat internal/service/project/project_service.go`

- [ ] **Step 2: 更新 CreateProjectRequest 结构**

```go
// 在 internal/service/project/project_service.go 中

// 变更前：
type CreateProjectRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Prefix      string `json:"prefix" validate:"required,min=2,max=4"`
    Description string `json:"description"`
}

// 变更后：
type CreateProjectRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    // Prefix 字段已删除
    Description string `json:"description"`
}
```

- [ ] **Step 3: 更新 UpdateProjectRequest 结构（如果有 prefix 字段）**

```go
// 检查 UpdateProjectRequest 是否有 Prefix 字段，如果有则删除
type UpdateProjectRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Description string `json:"description"`
}
```

- [ ] **Step 4: 更新 CreateProject 方法**

```go
// 在 internal/service/project/project_service.go 中

// 变更前：
func (s *ProjectServiceImpl) CreateProject(ctx context.Context, req *CreateProjectRequest) (*ProjectDetail, error) {
    // 前缀验证
    if err := s.validatePrefix(ctx, req.Prefix); err != nil {
        return nil, err
    }
    
    project, err := project.NewProject(req.Name, req.Prefix, req.Description)
    // ...
}

// 变更后：
func (s *ProjectServiceImpl) CreateProject(ctx context.Context, req *CreateProjectRequest) (*ProjectDetail, error) {
    // 前缀验证已删除
    
    project, err := project.NewProject(req.Name, req.Description)
    if err != nil {
        return nil, fmt.Errorf("create project: %w", err)
    }
    
    // 保存项目（前缀唯一性检查已删除）
    if err := s.projectRepo.Save(ctx, project); err != nil {
        if errors.Is(err, project.ErrProjectNameDuplicate) {
            return nil, fmt.Errorf("project name already exists: %w", err)
        }
        return nil, fmt.Errorf("save project: %w", err)
    }
    
    // ... 其余代码保持不变
}
```

- [ ] **Step 5: 删除前缀验证方法（如果存在）**

```go
// 在 internal/service/project/project_service.go 中

// 删除类似这样的方法：
// func (s *ProjectServiceImpl) validatePrefix(ctx context.Context, prefix string) error {
//     _, err := s.projectRepo.FindByPrefix(ctx, project.ProjectPrefix(prefix))
//     if err == nil {
//         return project.ErrProjectPrefixDuplicate
//     }
//     if !errors.Is(err, project.ErrProjectNotFound) {
//         return err
//     }
//     return nil
// }
```

- [ ] **Step 6: 更新响应转换方法**

```go
// 检查 toProjectResponse 等方法，删除 prefix 字段处理

// 变更前：
func toProjectResponse(p *project.Project) ProjectResponse {
    return ProjectResponse{
        ID:          p.ID(),
        Name:        p.Name(),
        Prefix:      p.Prefix().String(),  // 删除
        Description: p.Description(),
        CreatedAt:   p.CreatedAt(),
        UpdatedAt:   p.UpdatedAt(),
    }
}

// 变更后：
func toProjectResponse(p *project.Project) ProjectResponse {
    return ProjectResponse{
        ID:          p.ID(),
        Name:        p.Name(),
        Description: p.Description(),
        CreatedAt:   p.CreatedAt(),
        UpdatedAt:   p.UpdatedAt(),
    }
}
```

- [ ] **Step 7: 更新 ProjectDetail 结构**

```go
// 在 internal/service/project/project_service.go 中

// 变更前：
type ProjectDetail struct {
    *project.Project
    ModuleCount   int64 `json:"module_count"`
    CaseCount     int64 `json:"case_count"`
    DocumentCount int64 `json:"document_count"`
}

// 如果需要序列化，确保 JSON 输出不包含 prefix
```

- [ ] **Step 8: 更新测试文件**

```go
// 在 internal/service/project/project_service_test.go 中

// 更新所有调用 CreateProject 的测试

// 变更前：
req := &CreateProjectRequest{
    Name:        "Test Project",
    Prefix:      "TST",
    Description: "Test Description",
}

// 变更后：
req := &CreateProjectRequest{
    Name:        "Test Project",
    Description: "Test Description",
}
```

- [ ] **Step 9: 运行项目服务测试**

运行：`go test ./internal/service/project/... -v`

预期：所有测试通过

- [ ] **Step 10: 提交变更**

```bash
git add internal/service/project/project_service.go internal/service/project/project_service_test.go
git commit -m "refactor(service): remove prefix from project service

- Remove Prefix field from CreateProjectRequest
- Remove prefix validation logic
- Update CreateProject to call NewProject without prefix
- Update response conversion methods

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 14: 更新模块服务

**Files:**
- Modify: `internal/service/project/project_service.go`（模块相关方法）
- Modify: `internal/service/project/project_service_test.go`

- [ ] **Step 1: 更新 CreateModuleRequest 结构**

```go
// 在 internal/service/project/project_service.go 中

// 变更前：
type CreateModuleRequest struct {
    Name         string `json:"name" validate:"required,min=2,max=255"`
    Abbreviation string `json:"abbreviation" validate:"required,min=2,max=4"`
    Description  string `json:"description"`
}

// 变更后：
type CreateModuleRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    // Abbreviation 字段已删除
    Description string `json:"description"`
}
```

- [ ] **Step 2: 更新 UpdateModuleRequest 结构**

```go
// 在 internal/service/project/project_service.go 中

// 变更前：
type UpdateModuleRequest struct {
    Name         string `json:"name" validate:"required,min=2,max=255"`
    Abbreviation string `json:"abbreviation" validate:"min=2,max=4"`
    Description  string `json:"description"`
}

// 变更后：
type UpdateModuleRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Description string `json:"description"`
}
```

- [ ] **Step 3: 更新 CreateModule 方法**

```go
// 在 internal/service/project/project_service.go 中

// 变更前：
func (s *ProjectServiceImpl) CreateModule(ctx context.Context, projectID uuid.UUID, req *CreateModuleRequest) (*project.Module, error) {
    // 缩写验证
    if err := s.validateAbbreviation(ctx, projectID, req.Abbreviation); err != nil {
        return nil, err
    }
    
    module, err := project.NewModule(projectID, req.Name, req.Abbreviation, req.Description)
    // ...
}

// 变更后：
func (s *ProjectServiceImpl) CreateModule(ctx context.Context, projectID uuid.UUID, req *CreateModuleRequest) (*project.Module, error) {
    // 缩写验证已删除
    
    module, err := project.NewModule(projectID, req.Name, req.Description)
    if err != nil {
        return nil, fmt.Errorf("create module: %w", err)
    }
    
    // 保存模块（缩写唯一性检查已删除）
    if err := s.moduleRepo.Save(ctx, module); err != nil {
        if errors.Is(err, project.ErrModuleNameDuplicate) {
            return nil, fmt.Errorf("module name already exists: %w", err)
        }
        return nil, fmt.Errorf("save module: %w", err)
    }
    
    return module, nil
}
```

- [ ] **Step 4: 删除缩写验证方法（如果存在）**

```go
// 在 internal/service/project/project_service.go 中

// 删除类似这样的方法：
// func (s *ProjectServiceImpl) validateAbbreviation(ctx context.Context, projectID uuid.UUID, abbrev string) error {
//     _, err := s.moduleRepo.FindByAbbreviation(ctx, projectID, project.ModuleAbbreviation(abbrev))
//     if err == nil {
//         return project.ErrModuleAbbrevDuplicate
//     }
//     if !errors.Is(err, project.ErrModuleNotFound) {
//         return err
//     }
//     return nil
// }
```

- [ ] **Step 5: 更新响应转换方法**

```go
// 检查 toModuleResponse 等方法，删除 abbreviation 字段处理

// 变更前：
func toModuleResponse(m *project.Module) ModuleResponse {
    return ModuleResponse{
        ID:           m.ID(),
        ProjectID:    m.ProjectID(),
        Name:         m.Name(),
        Abbreviation: m.Abbreviation().String(),  // 删除
        Description:  m.Description(),
        CreatedAt:    m.CreatedAt(),
        UpdatedAt:    m.UpdatedAt(),
    }
}

// 变更后：
func toModuleResponse(m *project.Module) ModuleResponse {
    return ModuleResponse{
        ID:          m.ID(),
        ProjectID:   m.ProjectID(),
        Name:        m.Name(),
        Description: m.Description(),
        CreatedAt:   m.CreatedAt(),
        UpdatedAt:   m.UpdatedAt(),
    }
}
```

- [ ] **Step 6: 更新测试文件**

```go
// 在 internal/service/project/project_service_test.go 中

// 更新所有调用 CreateModule 的测试

// 变更前：
req := &CreateModuleRequest{
    Name:         "Test Module",
    Abbreviation: "TM",
    Description:  "Test Description",
}

// 变更后：
req := &CreateModuleRequest{
    Name:        "Test Module",
    Description: "Test Description",
}
```

- [ ] **Step 7: 运行项目服务测试**

运行：`go test ./internal/service/project/... -v`

预期：所有测试通过

- [ ] **Step 8: 提交变更**

```bash
git add internal/service/project/project_service.go internal/service/project/project_service_test.go
git commit -m "refactor(service): remove abbreviation from module service

- Remove Abbreviation field from CreateModuleRequest and UpdateModuleRequest
- Remove abbreviation validation logic
- Update CreateModule to call NewModule without abbreviation
- Update response conversion methods

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 15: 更新测试用例服务

**Files:**
- Modify: `internal/service/testcase/case_service.go`
- Modify: `internal/service/testcase/case_service_test.go`

- [ ] **Step 1: 读取当前测试用例服务**

运行：`cat internal/service/testcase/case_service.go`

- [ ] **Step 2: 检查 CreateTestCaseRequest 结构**

```go
// 确认 CreateTestCaseRequest 是否有 Number 字段
// 通常不需要，因为编号是自动生成的

type CreateTestCaseRequest struct {
    ModuleID      uuid.UUID              `json:"module_id" validate:"required"`
    Title         string                 `json:"title" validate:"required,min=2,max=500"`
    Preconditions []string               `json:"preconditions"`
    Steps         []string               `json:"steps" validate:"required,min=1"`
    Expected      map[string]interface{} `json:"expected" validate:"required"`
    CaseType      CaseType               `json:"case_type" validate:"required"`
    Priority      Priority               `json:"priority" validate:"required"`
}
```

- [ ] **Step 3: 更新 CreateTestCase 方法**

```go
// 在 internal/service/testcase/case_service.go 中

// 变更前：
func (s *TestCaseServiceImpl) CreateTestCase(ctx context.Context, req *CreateTestCaseRequest) (*testcase.TestCase, error) {
    // 获取用户 ID
    userID := s.getUserID(ctx)
    
    // 获取项目和模块以生成编号
    proj, err := s.projectRepo.FindByID(ctx, req.ModuleID)
    // ...
    
    mod, err := s.moduleRepo.FindByID(ctx, req.ModuleID)
    // ...
    
    // 生成用例编号
    today := time.Now().Truncate(24 * time.Hour)
    count, err := s.caseRepo.CountByDate(ctx, req.ModuleID, today)
    // ...
    
    number := testcase.GenerateCaseNumber(
        proj.Prefix().String(),
        mod.Abbreviation().String(),
        int(count)+1,
    )
    
    tc, err := testcase.NewTestCase(
        req.ModuleID,
        userID,
        number,  // <-- 删除这个参数
        req.Title,
        testcase.Preconditions(req.Preconditions),
        testcase.Steps(req.Steps),
        testcase.ExpectedResult(req.Expected),
        testcase.CaseType(req.CaseType),
        testcase.Priority(req.Priority),
    )
    // ...
}

// 变更后：
func (s *TestCaseServiceImpl) CreateTestCase(ctx context.Context, req *CreateTestCaseRequest) (*testcase.TestCase, error) {
    // 获取用户 ID
    userID := s.getUserID(ctx)
    
    // 编号生成逻辑已删除
    
    tc, err := testcase.NewTestCase(
        req.ModuleID,
        userID,
        req.Title,
        testcase.Preconditions(req.Preconditions),
        testcase.Steps(req.Steps),
        testcase.ExpectedResult(req.Expected),
        testcase.CaseType(req.CaseType),
        testcase.Priority(req.Priority),
    )
    if err != nil {
        return nil, fmt.Errorf("create test case: %w", err)
    }
    
    if err := s.caseRepo.Save(ctx, tc); err != nil {
        return nil, fmt.Errorf("save test case: %w", err)
    }
    
    return tc, nil
}
```

- [ ] **Step 4: 更新响应转换方法**

```go
// 检查 toTestCaseResponse 方法，删除 number 字段处理

// 变更前：
func toTestCaseResponse(tc *testcase.TestCase) TestCaseResponse {
    return TestCaseResponse{
        ID:            tc.ID(),
        ModuleID:      tc.ModuleID(),
        UserID:        tc.UserID(),
        Number:        tc.Number().String(),  // 删除
        Title:         tc.Title(),
        Preconditions: tc.Preconditions(),
        Steps:         tc.Steps(),
        Expected:      tc.Expected(),
        CaseType:      string(tc.CaseType()),
        Priority:      string(tc.Priority()),
        Status:        string(tc.Status()),
        AiMetadata:    toAiMetadataJSON(tc.AiMetadata()),
        CreatedAt:     tc.CreatedAt(),
        UpdatedAt:     tc.UpdatedAt(),
    }
}

// 变更后：
func toTestCaseResponse(tc *testcase.TestCase) TestCaseResponse {
    return TestCaseResponse{
        ID:            tc.ID(),
        ModuleID:      tc.ModuleID(),
        UserID:        tc.UserID(),
        Title:         tc.Title(),
        Preconditions: tc.Preconditions(),
        Steps:         tc.Steps(),
        Expected:      tc.Expected(),
        CaseType:      string(tc.CaseType()),
        Priority:      string(tc.Priority()),
        Status:        string(tc.Status()),
        AiMetadata:    toAiMetadataJSON(tc.AiMetadata()),
        CreatedAt:     tc.CreatedAt(),
        UpdatedAt:     tc.UpdatedAt(),
    }
}
```

- [ ] **Step 5: 更新 TestCaseResponse 结构**

```go
// 在 internal/service/testcase/case_service.go 中

// 变更前：
type TestCaseResponse struct {
    ID            uuid.UUID              `json:"id"`
    ModuleID      uuid.UUID              `json:"moduleId"`
    UserID        uuid.UUID              `json:"userId"`
    Number        string                 `json:"number"`  // 删除
    Title         string                 `json:"title"`
    Preconditions []string               `json:"preconditions"`
    Steps         []string               `json:"steps"`
    Expected      map[string]interface{} `json:"expected"`
    CaseType      string                 `json:"caseType"`
    Priority      string                 `json:"priority"`
    Status        string                 `json:"status"`
    AiMetadata    *AiMetadataJSON        `json:"aiMetadata,omitempty"`
    CreatedAt     time.Time              `json:"createdAt"`
    UpdatedAt     time.Time              `json:"updatedAt"`
}

// 变更后：
type TestCaseResponse struct {
    ID            uuid.UUID              `json:"id"`
    ModuleID      uuid.UUID              `json:"moduleId"`
    UserID        uuid.UUID              `json:"userId"`
    Title         string                 `json:"title"`
    Preconditions []string               `json:"preconditions"`
    Steps         []string               `json:"steps"`
    Expected      map[string]interface{} `json:"expected"`
    CaseType      string                 `json:"caseType"`
    Priority      string                 `json:"priority"`
    Status        string                 `json:"status"`
    AiMetadata    *AiMetadataJSON        `json:"aiMetadata,omitempty"`
    CreatedAt     time.Time              `json:"createdAt"`
    UpdatedAt     time.Time              `json:"updatedAt"`
}
```

- [ ] **Step 6: 删除不必要的依赖注入**

```go
// 检查 TestCaseServiceImpl 结构，如果不再需要 projectRepo 用于编号生成，可以移除

// 变更前：
type TestCaseServiceImpl struct {
    caseRepo    testcase.TestCaseRepository
    moduleRepo  project.ModuleRepository  // 如果仅用于编号生成，可能删除
    projectRepo project.ProjectRepository  // 如果仅用于编号生成，可能删除
}

// 变更后（如果这些 repo 仅用于编号生成）：
type TestCaseServiceImpl struct {
    caseRepo   testcase.TestCaseRepository
}
```

注意：只有在 projectRepo 和 moduleRepo 不再被其他方法使用时才能删除。

- [ ] **Step 7: 更新测试文件**

```go
// 在 internal/service/testcase/case_service_test.go 中

// 更新所有调用 CreateTestCase 的测试
// 删除所有与编号相关的断言

// 变更前：
assert.Equal(t, "ECO-USR-20260505-001", result.Number())

// 删除这行或验证 ID 不为空
assert.NotEqual(t, uuid.Nil, result.ID())
```

- [ ] **Step 8: 运行测试用例服务测试**

运行：`go test ./internal/service/testcase/... -v`

预期：所有测试通过

- [ ] **Step 9: 提交变更**

```bash
git add internal/service/testcase/case_service.go internal/service/testcase/case_service_test.go
git commit -m "refactor(service): remove case number generation from test case service

- Remove number generation logic from CreateTestCase
- Remove Number field from TestCaseResponse
- Update CreateTestCase to call NewTestCase without number
- Update response conversion methods
- Remove unnecessary repo dependencies if applicable

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Phase 5: 持久层代码变更

### Task 16: 更新项目 Repository

**Files:**
- Modify: `internal/repository/project/project_repo.go`
- Modify: `internal/repository/project/project_repo_test.go`

- [ ] **Step 1: 读取当前项目 repository**

运行：`cat internal/repository/project/project_repo.go`

- [ ] **Step 2: 更新表名常量和查询**

```go
// 在 internal/repository/project/project_repo.go 中

// 变更前：
const projectTableName = "project"

const projectSelectQuery = `
    SELECT id, name, prefix, description, config, created_at, updated_at
    FROM project
    WHERE id = $1
`

// 变更后：
const projectTableName = "projects"

const projectSelectQuery = `
    SELECT id, name, description, config, created_at, updated_at
    FROM projects
    WHERE id = $1
`
```

- [ ] **Step 3: 更新 FindByID 方法**

```go
// 在 internal/repository/project/project_repo.go 中

// 变更前：
func (r *ProjectRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*project.Project, error) {
    var p project.Project
    var prefixStr string
    
    err := r.db.QueryRowContext(ctx, projectSelectQuery, id).Scan(
        &p.id, &p.name, &prefixStr, &p.description,
        &p.config, &p.createdAt, &p.updatedAt,
    )
    // ...
}

// 变更后：
func (r *ProjectRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*project.Project, error) {
    var p project.Project
    
    err := r.db.QueryRowContext(ctx, projectSelectQuery, id).Scan(
        &p.id, &p.name, &p.description,
        &p.config, &p.createdAt, &p.updatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, project.ErrProjectNotFound
        }
        return nil, fmt.Errorf("find project by id: %w", err)
    }
    
    // prefix 处理已删除
    return &p, nil
}
```

- [ ] **Step 4: 更新 Save 方法**

```go
// 在 internal/repository/project/project_repo.go 中

// 变更前：
const projectInsertQuery = `
    INSERT INTO project (id, name, prefix, description, config, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT (name) DO UPDATE SET
        name = EXCLUDED.name,
        prefix = EXCLUDED.prefix,
        description = EXCLUDED.description,
        config = EXCLUDED.config,
        updated_at = CURRENT_TIMESTAMP
`

func (r *ProjectRepoImpl) Save(ctx context.Context, p *project.Project) error {
    _, err := r.db.ExecContext(ctx, projectInsertQuery,
        p.ID(), p.Name(), p.Prefix(), p.Description(),
        p.Config(), p.CreatedAt(), p.UpdatedAt(),
    )
    // ...
}

// 变更后：
const projectInsertQuery = `
    INSERT INTO projects (id, name, description, config, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (name) DO UPDATE SET
        name = EXCLUDED.name,
        description = EXCLUDED.description,
        config = EXCLUDED.config,
        updated_at = CURRENT_TIMESTAMP
`

func (r *ProjectRepoImpl) Save(ctx context.Context, p *project.Project) error {
    _, err := r.db.ExecContext(ctx, projectInsertQuery,
        p.ID(), p.Name(), p.Description(),
        p.Config(), p.CreatedAt(), p.UpdatedAt(),
    )
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
            if strings.Contains(pgErr.ConstraintName, "name") {
                return project.ErrProjectNameDuplicate
            }
        }
        return fmt.Errorf("save project: %w", err)
    }
    return nil
}
```

- [ ] **Step 5: 删除 FindByPrefix 方法**

```go
// 在 internal/repository/project/project_repo.go 中

// 删除整个方法：
// func (r *ProjectRepoImpl) FindByPrefix(ctx context.Context, prefix project.ProjectPrefix) (*project.Project, error) {
//     // ...
// }
```

- [ ] **Step 6: 更新其他查询方法**

更新所有涉及 `project` 表的查询字符串，将 `project` 改为 `projects`。

- [ ] **Step 7: 更新测试文件**

```go
// 在 internal/repository/project/project_repo_test.go 中

// 更新所有测试中的 SQL 断言（如果有的话）
// 更新 mock 数据（删除 prefix 字段）
```

- [ ] **Step 8: 运行项目 repository 测试**

运行：`go test ./internal/repository/project/... -v`

预期：所有测试通过

- [ ] **Step 9: 提交变更**

```bash
git add internal/repository/project/project_repo.go internal/repository/project/project_repo_test.go
git commit -m "refactor(repository): update project repo for table rename and prefix removal

- Rename table: project → projects
- Remove prefix from SELECT and INSERT queries
- Remove FindByPrefix method
- Update constraint handling

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 17: 更新模块 Repository

**Files:**
- Modify: `internal/repository/project/module_repo.go`
- Modify: `internal/repository/project/module_repo_test.go`

- [ ] **Step 1: 读取当前模块 repository**

运行：`cat internal/repository/project/module_repo.go`

- [ ] **Step 2: 更新表名常量和查询**

```go
// 在 internal/repository/project/module_repo.go 中

// 变更前：
const moduleTableName = "module"

const moduleSelectQuery = `
    SELECT id, project_id, name, abbreviation, description, created_at, updated_at
    FROM module
    WHERE id = $1
`

// 变更后：
const moduleTableName = "modules"

const moduleSelectQuery = `
    SELECT id, project_id, name, description, created_at, updated_at
    FROM modules
    WHERE id = $1
`
```

- [ ] **Step 3: 更新 FindByID 方法**

```go
// 在 internal/repository/project/module_repo.go 中

// 变更前：
func (r *ModuleRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*project.Module, error) {
    var m project.Module
    var abbrevStr string
    
    err := r.db.QueryRowContext(ctx, moduleSelectQuery, id).Scan(
        &m.id, &m.projectID, &m.name, &abbrevStr, &m.description,
        &m.createdAt, &m.updatedAt,
    )
    // ...
}

// 变更后：
func (r *ModuleRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*project.Module, error) {
    var m project.Module
    
    err := r.db.QueryRowContext(ctx, moduleSelectQuery, id).Scan(
        &m.id, &m.projectID, &m.name, &m.description,
        &m.createdAt, &m.updatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, project.ErrModuleNotFound
        }
        return nil, fmt.Errorf("find module by id: %w", err)
    }
    
    return &m, nil
}
```

- [ ] **Step 4: 更新 Save 方法**

```go
// 在 internal/repository/project/module_repo.go 中

// 变更前：
const moduleInsertQuery = `
    INSERT INTO module (id, project_id, name, abbreviation, description, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT (project_id, name) DO UPDATE SET
        name = EXCLUDED.name,
        abbreviation = EXCLUDED.abbreviation,
        description = EXCLUDED.description,
        updated_at = CURRENT_TIMESTAMP
`

func (r *ModuleRepoImpl) Save(ctx context.Context, m *project.Module) error {
    _, err := r.db.ExecContext(ctx, moduleInsertQuery,
        m.ID(), m.ProjectID(), m.Name(), m.Abbreviation(), m.Description(),
        m.CreatedAt(), m.UpdatedAt(),
    )
    // ...
}

// 变更后：
const moduleInsertQuery = `
    INSERT INTO modules (id, project_id, name, description, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (project_id, name) DO UPDATE SET
        name = EXCLUDED.name,
        description = EXCLUDED.description,
        updated_at = CURRENT_TIMESTAMP
`

func (r *ModuleRepoImpl) Save(ctx context.Context, m *project.Module) error {
    _, err := r.db.ExecContext(ctx, moduleInsertQuery,
        m.ID(), m.ProjectID(), m.Name(), m.Description(),
        m.CreatedAt(), m.UpdatedAt(),
    )
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            if strings.Contains(pgErr.ConstraintName, "name") {
                return project.ErrModuleNameDuplicate
            }
        }
        return fmt.Errorf("save module: %w", err)
    }
    return nil
}
```

- [ ] **Step 5: 删除 FindByAbbreviation 方法**

```go
// 在 internal/repository/project/module_repo.go 中

// 删除整个方法：
// func (r *ModuleRepoImpl) FindByAbbreviation(ctx context.Context, projectID uuid.UUID, abbrev project.ModuleAbbreviation) (*project.Module, error) {
//     // ...
// }
```

- [ ] **Step 6: 更新 FindByProjectID 方法**

```go
// 在 internal/repository/project/module_repo.go 中

// 变更前：
const findByProjectIDQuery = `
    SELECT id, project_id, name, abbreviation, description, created_at, updated_at
    FROM module
    WHERE project_id = $1
    ORDER BY created_at DESC
`

// 变更后：
const findByProjectIDQuery = `
    SELECT id, project_id, name, description, created_at, updated_at
    FROM modules
    WHERE project_id = $1
    ORDER BY created_at DESC
`
```

- [ ] **Step 7: 更新测试文件**

```go
// 在 internal/repository/project/module_repo_test.go 中

// 更新所有测试中的 SQL 断言
// 更新 mock 数据（删除 abbreviation 字段）
```

- [ ] **Step 8: 运行模块 repository 测试**

运行：`go test ./internal/repository/project/... -v`

预期：所有测试通过

- [ ] **Step 9: 提交变更**

```bash
git add internal/repository/project/module_repo.go internal/repository/project/module_repo_test.go
git commit -m "refactor(repository): update module repo for table rename and abbreviation removal

- Rename table: module → modules
- Remove abbreviation from SELECT and INSERT queries
- Remove FindByAbbreviation method
- Update all query methods

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 18: 更新用户 Repository

**Files:**
- Modify: `internal/repository/identity/user_repo.go`
- Modify: `internal/repository/identity/user_repo_test.go`

- [ ] **Step 1: 读取当前用户 repository**

运行：`cat internal/repository/identity/user_repo.go`

- [ ] **Step 2: 更新查询删除 deleted_at 检查**

```go
// 在 internal/repository/identity/user_repo.go 中

// 变更前：
const userSelectQuery = `
    SELECT id, username, email, password, role, created_at, updated_at, deleted_at
    FROM users
    WHERE id = $1 AND deleted_at IS NULL
`

// 变更后：
const userSelectQuery = `
    SELECT id, username, email, password, role, created_at, updated_at
    FROM users
    WHERE id = $1
`
```

- [ ] **Step 3: 更新 FindByID 方法**

```go
// 在 internal/repository/identity/user_repo.go 中

// 变更前：
func (r *UserRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*identity.User, error) {
    var u identity.User
    var deletedAt sql.NullTime
    
    err := r.db.QueryRowContext(ctx, userSelectQuery, id).Scan(
        &u.id, &u.username, &u.email, &u.password, &u.role,
        &u.createdAt, &u.updatedAt, &deletedAt,
    )
    // ...
}

// 变更后：
func (r *UserRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*identity.User, error) {
    var u identity.User
    
    err := r.db.QueryRowContext(ctx, userSelectQuery, id).Scan(
        &u.id, &u.username, &u.email, &u.password, &u.role,
        &u.createdAt, &u.updatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, identity.ErrUserNotFound
        }
        return nil, fmt.Errorf("find user by id: %w", err)
    }
    
    return &u, nil
}
```

- [ ] **Step 4: 更新 Save 方法**

```go
// 在 internal/repository/identity/user_repo.go 中

// 变更前：
const userInsertQuery = `
    INSERT INTO users (id, username, email, password, role, created_at, updated_at, deleted_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)
    ON CONFLICT (email) DO UPDATE SET
        username = EXCLUDED.username,
        password = EXCLUDED.password,
        role = EXCLUDED.role,
        updated_at = CURRENT_TIMESTAMP
`

// 变更后：
const userInsertQuery = `
    INSERT INTO users (id, username, email, password, role, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT (email) DO UPDATE SET
        username = EXCLUDED.username,
        password = EXCLUDED.password,
        role = EXCLUDED.role,
        updated_at = CURRENT_TIMESTAMP
`
```

- [ ] **Step 5: 更新其他查询方法**

更新 FindByEmail、FindByUsername 等方法，删除 `deleted_at IS NULL` 条件。

- [ ] **Step 6: 删除 SoftDelete 方法（如果存在）**

```go
// 在 internal/repository/identity/user_repo.go 中

// 删除整个方法：
// func (r *UserRepoImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
//     query := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
//     // ...
// }
```

- [ ] **Step 7: 更新 Delete 方法（如果存在）**

如果 Delete 方法之前是软删除，改为硬删除：

```go
// 在 internal/repository/identity/user_repo.go 中

// 变更前（软删除）：
func (r *UserRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
    query := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
    // ...
}

// 变更后（硬删除）：
func (r *UserRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("delete user: %w", err)
    }
    return nil
}
```

- [ ] **Step 8: 更新测试文件**

```go
// 在 internal/repository/identity/user_repo_test.go 中

// 删除所有与软删除相关的测试
// 更新 mock 数据（删除 deleted_at 字段）
```

- [ ] **Step 9: 运行用户 repository 测试**

运行：`go test ./internal/repository/identity/... -v`

预期：所有测试通过

- [ ] **Step 10: 提交变更**

```bash
git add internal/repository/identity/user_repo.go internal/repository/identity/user_repo_test.go
git commit -m "refactor(repository): remove soft delete from user repo

- Remove deleted_at from all user queries
- Change Delete from soft delete to hard delete
- Remove SoftDelete method
- Update all SELECT queries to remove deleted_at IS NULL check

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 19: 更新测试用例 Repository

**Files:**
- Modify: `internal/repository/testcase/case_repo.go`
- Modify: `internal/repository/testcase/case_repo_test.go`

- [ ] **Step 1: 读取当前测试用例 repository**

运行：`cat internal/repository/testcase/case_repo.go`

- [ ] **Step 2: 更新表名常量和查询**

```go
// 在 internal/repository/testcase/case_repo.go 中

// 更新所有涉及 module 的外键引用（module → modules）

const caseSelectQuery = `
    SELECT id, module_id, user_id, title, preconditions, steps, expected,
           ai_metadata, case_type, priority, status, created_at, updated_at
    FROM test_case
    WHERE id = $1
`
```

- [ ] **Step 3: 更新 FindByID 方法**

```go
// 在 internal/repository/testcase/case_repo.go 中

// 变更前：
func (r *TestCaseRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*testcase.TestCase, error) {
    var tc testcase.TestCase
    var numberStr string
    
    err := r.db.QueryRowContext(ctx, caseSelectQuery, id).Scan(
        &tc.id, &tc.moduleID, &tc.userID, &numberStr, &tc.title,
        // ...
    )
    // ...
}

// 变更后：
func (r *TestCaseRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*testcase.TestCase, error) {
    var tc testcase.TestCase
    
    err := r.db.QueryRowContext(ctx, caseSelectQuery, id).Scan(
        &tc.id, &tc.moduleID, &tc.userID, &tc.title,
        &tc.preconditions, &tc.steps, &tc.expected,
        &tc.aiMetadata, &tc.caseType, &tc.priority,
        &tc.status, &tc.createdAt, &tc.updatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, testcase.ErrCaseNotFound
        }
        return nil, fmt.Errorf("find test case by id: %w", err)
    }
    
    return &tc, nil
}
```

- [ ] **Step 4: 更新 Save 方法**

```go
// 在 internal/repository/testcase/case_repo.go 中

// 变更前：
const caseInsertQuery = `
    INSERT INTO test_case (id, module_id, user_id, number, title, preconditions, 
                          steps, expected, ai_metadata, case_type, priority, status, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    ON CONFLICT (id) DO UPDATE SET
        title = EXCLUDED.title,
        -- ...
`

// 变更后：
const caseInsertQuery = `
    INSERT INTO test_case (id, module_id, user_id, title, preconditions, 
                          steps, expected, ai_metadata, case_type, priority, status, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    ON CONFLICT (id) DO UPDATE SET
        title = EXCLUDED.title,
        preconditions = EXCLUDED.preconditions,
        steps = EXCLUDED.steps,
        expected = EXCLUDED.expected,
        case_type = EXCLUDED.case_type,
        priority = EXCLUDED.priority,
        status = EXCLUDED.status,
        ai_metadata = EXCLUDED.ai_metadata,
        updated_at = CURRENT_TIMESTAMP
`

func (r *TestCaseRepoImpl) Save(ctx context.Context, tc *testcase.TestCase) error {
    _, err := r.db.ExecContext(ctx, caseInsertQuery,
        tc.ID(), tc.ModuleID(), tc.UserID(), tc.Title(),
        tc.Preconditions(), tc.Steps(), tc.Expected(),
        tc.AiMetadata(), tc.CaseType(), tc.Priority(),
        tc.Status(), tc.CreatedAt(), tc.UpdatedAt(),
    )
    if err != nil {
        return fmt.Errorf("save test case: %w", err)
    }
    return nil
}
```

- [ ] **Step 5: 删除 FindByNumber 方法**

```go
// 在 internal/repository/testcase/case_repo.go 中

// 删除整个方法：
// func (r *TestCaseRepoImpl) FindByNumber(ctx context.Context, number testcase.CaseNumber) (*testcase.TestCase, error) {
//     // ...
// }
```

- [ ] **Step 6: 删除 CountByDate 方法**

```go
// 在 internal/repository/testcase/case_repo.go 中

// 删除整个方法（不再需要编号生成）：
// func (r *TestCaseRepoImpl) CountByDate(ctx context.Context, moduleID uuid.UUID, date time.Time) (int64, error) {
//     query := `
//         SELECT COUNT(*)
//         FROM test_case
//         WHERE module_id = $1 
//           AND created_at >= $2 
//           AND created_at < $2 + INTERVAL '1 day'
//     `
//     // ...
// }
```

- [ ] **Step 7: 更新接口定义（如果有）**

检查 TestCaseRepository 接口，删除不需要的方法：

```go
// 在 internal/domain/testcase/repository.go 中

// 删除接口中的方法定义：
// type TestCaseRepository interface {
//     // ...
//     FindByNumber(ctx context.Context, number CaseNumber) (*TestCase, error)  // 删除
//     CountByDate(ctx context.Context, moduleID uuid.UUID, date time.Time) (int64, error)  // 删除
// }
```

- [ ] **Step 8: 更新测试文件**

```go
// 在 internal/repository/testcase/case_repo_test.go 中

// 删除所有与编号相关的测试
// 更新 mock 数据（删除 number 字段）
```

- [ ] **Step 9: 运行测试用例 repository 测试**

运行：`go test ./internal/repository/testcase/... -v`

预期：所有测试通过

- [ ] **Step 10: 提交变更**

```bash
git add internal/repository/testcase/case_repo.go internal/repository/testcase/case_repo_test.go
git add internal/domain/testcase/repository.go
git commit -m "refactor(repository): remove case number from test case repo

- Remove number from SELECT and INSERT queries
- Remove FindByNumber method
- Remove CountByDate method (no longer needed for numbering)
- Update TestCaseRepository interface
- Update foreign key reference: module → modules

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Phase 6: 传输层代码变更

### Task 20: 更新项目 Handler

**Files:**
- Modify: `internal/transport/http/handler/project.go`
- Modify: `internal/transport/http/handler/project_test.go`

- [ ] **Step 1: 读取当前项目 handler**

运行：`cat internal/transport/http/handler/project.go`

- [ ] **Step 2: 更新 CreateProjectRequest 结构**

```go
// 在 internal/transport/http/handler/project.go 中

// 变更前：
type CreateProjectRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Prefix      string `json:"prefix" validate:"required,min=2,max=4"`
    Description string `json:"description"`
}

// 变更后：
type CreateProjectRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Description string `json:"description"`
}
```

- [ ] **Step 3: 更新 ProjectResponse 结构**

```go
// 在 internal/transport/http/handler/project.go 中

// 变更前：
type ProjectResponse struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Prefix      string    `json:"prefix"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

// 变更后：
type ProjectResponse struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

- [ ] **Step 4: 更新 CreateProject handler**

```go
// 在 internal/transport/http/handler/project.go 中

// 变更前：
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
    var req CreateProjectRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httperror.RespondWithError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    project, err := h.projectService.CreateProject(r.Context(), &project_service.CreateProjectRequest{
        Name:        req.Name,
        Prefix:      req.Prefix,
        Description: req.Description,
    })
    // ...
}

// 变更后：
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
    var req CreateProjectRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httperror.RespondWithError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    project, err := h.projectService.CreateProject(r.Context(), &project_service.CreateProjectRequest{
        Name:        req.Name,
        Description: req.Description,
    })
    if err != nil {
        httperror.RespondWithError(w, httpcode.GetHTTPStatus(err), err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusCreated, toProjectResponse(project))
}
```

- [ ] **Step 5: 更新 UpdateProjectRequest 结构**

```go
// 检查 UpdateProjectRequest 是否有 Prefix 字段，如果有则删除

type UpdateProjectRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Description string `json:"description"`
}
```

- [ ] **Step 6: 更新响应转换函数**

```go
// 在 internal/transport/http/handler/project.go 中

func toProjectResponse(p *project.Project) ProjectResponse {
    return ProjectResponse{
        ID:          p.ID(),
        Name:        p.Name(),
        Description: p.Description(),
        CreatedAt:   p.CreatedAt(),
        UpdatedAt:   p.UpdatedAt(),
    }
}
```

- [ ] **Step 7: 更新 ProjectDetailResponse 结构**

```go
// 如果有 ProjectDetailResponse，删除相关字段

type ProjectDetailResponse struct {
    ID            uuid.UUID `json:"id"`
    Name          string    `json:"name"`
    Description   string    `json:"description"`
    ModuleCount   int64     `json:"module_count"`
    CaseCount     int64     `json:"case_count"`
    DocumentCount int64     `json:"document_count"`
    CreatedAt     time.Time `json:"createdAt"`
    UpdatedAt     time.Time `json:"updatedAt"`
}
```

- [ ] **Step 8: 运行项目 handler 测试**

运行：`go test ./internal/transport/http/handler/... -v -run TestProject`

预期：所有测试通过

- [ ] **Step 9: 提交变更**

```bash
git add internal/transport/http/handler/project.go internal/transport/http/handler/project_test.go
git commit -m "refactor(handler): remove prefix from project handler

- Remove Prefix from CreateProjectRequest
- Remove Prefix from ProjectResponse
- Update CreateProject handler
- Update response conversion functions

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 21: 更新模块 Handler

**Files:**
- Modify: `internal/transport/http/handler/project.go`（模块相关 handler）
- Modify: `internal/transport/http/handler/project_test.go`

- [ ] **Step 1: 更新 CreateModuleRequest 结构**

```go
// 在 internal/transport/http/handler/project.go 中

// 变更前：
type CreateModuleRequest struct {
    Name         string `json:"name" validate:"required,min=2,max=255"`
    Abbreviation string `json:"abbreviation" validate:"required,min=2,max=4"`
    Description  string `json:"description"`
}

// 变更后：
type CreateModuleRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Description string `json:"description"`
}
```

- [ ] **Step 2: 更新 UpdateModuleRequest 结构**

```go
// 在 internal/transport/http/handler/project.go 中

// 变更前：
type UpdateModuleRequest struct {
    Name         string `json:"name" validate:"required,min=2,max=255"`
    Abbreviation string `json:"abbreviation" validate:"min=2,max=4"`
    Description  string `json:"description"`
}

// 变更后：
type UpdateModuleRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Description string `json:"description"`
}
```

- [ ] **Step 3: 更新 ModuleResponse 结构**

```go
// 在 internal/transport/http/handler/project.go 中

// 变更前：
type ModuleResponse struct {
    ID           uuid.UUID `json:"id"`
    ProjectID    uuid.UUID `json:"projectId"`
    Name         string    `json:"name"`
    Abbreviation string    `json:"abbreviation"`
    Description  string    `json:"description"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}

// 变更后：
type ModuleResponse struct {
    ID          uuid.UUID `json:"id"`
    ProjectID   uuid.UUID `json:"projectId"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

- [ ] **Step 4: 更新 CreateModule handler**

```go
// 在 internal/transport/http/handler/project.go 中

// 变更前：
func (h *ProjectHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
    projectID, err := uuid.Parse(chi.URLParam(r, "id"))
    // ...
    
    var req CreateModuleRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httperror.RespondWithError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    module, err := h.projectService.CreateModule(r.Context(), projectID, &project_service.CreateModuleRequest{
        Name:         req.Name,
        Abbreviation: req.Abbreviation,
        Description:  req.Description,
    })
    // ...
}

// 变更后：
func (h *ProjectHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
    projectID, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        httperror.RespondWithError(w, http.StatusBadRequest, "invalid project id")
        return
    }
    
    var req CreateModuleRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httperror.RespondWithError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    module, err := h.projectService.CreateModule(r.Context(), projectID, &project_service.CreateModuleRequest{
        Name:        req.Name,
        Description: req.Description,
    })
    if err != nil {
        httperror.RespondWithError(w, httpcode.GetHTTPStatus(err), err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusCreated, toModuleResponse(module))
}
```

- [ ] **Step 5: 更新响应转换函数**

```go
// 在 internal/transport/http/handler/project.go 中

func toModuleResponse(m *project.Module) ModuleResponse {
    return ModuleResponse{
        ID:          m.ID(),
        ProjectID:   m.ProjectID(),
        Name:        m.Name(),
        Description: m.Description(),
        CreatedAt:   m.CreatedAt(),
        UpdatedAt:   m.UpdatedAt(),
    }
}
```

- [ ] **Step 6: 运行模块 handler 测试**

运行：`go test ./internal/transport/http/handler/... -v -run TestModule`

预期：所有测试通过

- [ ] **Step 7: 提交变更**

```bash
git add internal/transport/http/handler/project.go internal/transport/http/handler/project_test.go
git commit -m "refactor(handler): remove abbreviation from module handler

- Remove Abbreviation from CreateModuleRequest and UpdateModuleRequest
- Remove Abbreviation from ModuleResponse
- Update CreateModule handler
- Update response conversion functions

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

### Task 22: 更新测试用例 Handler

**Files:**
- Modify: `internal/transport/http/handler/testcase.go`
- Modify: `internal/transport/http/handler/testcase_test.go`

- [ ] **Step 1: 读取当前测试用例 handler**

运行：`cat internal/transport/http/handler/testcase.go`

- [ ] **Step 2: 更新 TestCaseResponse 结构**

```go
// 在 internal/transport/http/handler/testcase.go 中

// 变更前：
type TestCaseResponse struct {
    ID            uuid.UUID              `json:"id"`
    ModuleID      uuid.UUID              `json:"moduleId"`
    UserID        uuid.UUID              `json:"userId"`
    Number        string                 `json:"number"`
    Title         string                 `json:"title"`
    Preconditions []string               `json:"preconditions"`
    Steps         []string               `json:"steps"`
    Expected      map[string]interface{} `json:"expected"`
    CaseType      string                 `json:"caseType"`
    Priority      string                 `json:"priority"`
    Status        string                 `json:"status"`
    AiMetadata    *AiMetadataJSON        `json:"aiMetadata,omitempty"`
    CreatedAt     time.Time              `json:"createdAt"`
    UpdatedAt     time.Time              `json:"updatedAt"`
}

// 变更后：
type TestCaseResponse struct {
    ID            uuid.UUID              `json:"id"`
    ModuleID      uuid.UUID              `json:"moduleId"`
    UserID        uuid.UUID              `json:"userId"`
    Title         string                 `json:"title"`
    Preconditions []string               `json:"preconditions"`
    Steps         []string               `json:"steps"`
    Expected      map[string]interface{} `json:"expected"`
    CaseType      string                 `json:"caseType"`
    Priority      string                 `json:"priority"`
    Status        string                 `json:"status"`
    AiMetadata    *AiMetadataJSON        `json:"aiMetadata,omitempty"`
    CreatedAt     time.Time              `json:"createdAt"`
    UpdatedAt     time.Time              `json:"updatedAt"`
}
```

- [ ] **Step 3: 更新 CaseDetailResponse 结构**

```go
// 如果有 CaseDetailResponse，同样删除 Number 字段

type CaseDetailResponse struct {
    ID            uuid.UUID              `json:"id"`
    ModuleID      uuid.UUID              `json:"moduleId"`
    UserID        uuid.UUID              `json:"userId"`
    ModuleName    string                 `json:"module_name,omitempty"`
    ProjectName   string                 `json:"project_name,omitempty"`
    Title         string                 `json:"title"`
    Preconditions []string               `json:"preconditions"`
    Steps         []string               `json:"steps"`
    Expected      map[string]interface{} `json:"expected"`
    CaseType      string                 `json:"caseType"`
    Priority      string                 `json:"priority"`
    Status        string                 `json:"status"`
    AiMetadata    *AiMetadataJSON        `json:"aiMetadata,omitempty"`
    CreatedAt     time.Time              `json:"createdAt"`
    UpdatedAt     time.Time              `json:"updatedAt"`
}
```

- [ ] **Step 4: 更新响应转换函数**

```go
// 在 internal/transport/http/handler/testcase.go 中

// 变更前：
func toTestCaseResponse(tc *testcase.TestCase) TestCaseResponse {
    return TestCaseResponse{
        ID:            tc.ID(),
        ModuleID:      tc.ModuleID(),
        UserID:        tc.UserID(),
        Number:        tc.Number().String(),
        Title:         tc.Title(),
        Preconditions: tc.Preconditions(),
        Steps:         tc.Steps(),
        Expected:      tc.Expected(),
        CaseType:      string(tc.CaseType()),
        Priority:      string(tc.Priority()),
        Status:        string(tc.Status()),
        AiMetadata:    toAiMetadataJSON(tc.AiMetadata()),
        CreatedAt:     tc.CreatedAt(),
        UpdatedAt:     tc.UpdatedAt(),
    }
}

// 变更后：
func toTestCaseResponse(tc *testcase.TestCase) TestCaseResponse {
    return TestCaseResponse{
        ID:            tc.ID(),
        ModuleID:      tc.ModuleID(),
        UserID:        tc.UserID(),
        Title:         tc.Title(),
        Preconditions: tc.Preconditions(),
        Steps:         tc.Steps(),
        Expected:      tc.Expected(),
        CaseType:      string(tc.CaseType()),
        Priority:      string(tc.Priority()),
        Status:        string(tc.Status()),
        AiMetadata:    toAiMetadataJSON(tc.AiMetadata()),
        CreatedAt:     tc.CreatedAt(),
        UpdatedAt:     tc.UpdatedAt(),
    }
}
```

- [ ] **Step 5: 运行测试用例 handler 测试**

运行：`go test ./internal/transport/http/handler/... -v -run TestCase`

预期：所有测试通过

- [ ] **Step 6: 提交变更**

```bash
git add internal/transport/http/handler/testcase.go internal/transport/http/handler/testcase_test.go
git commit -m "refactor(handler): remove number from test case handler

- Remove Number from TestCaseResponse
- Remove Number from CaseDetailResponse
- Update response conversion functions

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Phase 7: 最终验证

### Task 23: 运行完整测试套件

**Files:**
- None (验证任务)

- [ ] **Step 1: 格式化代码**

运行：`go fmt ./...`

- [ ] **Step 2: 运行代码检查**

运行：`go vet ./...`

- [ ] **Step 3: 运行所有单元测试**

运行：`go test ./... -v`

预期：所有测试通过

- [ ] **Step 4: 运行静态分析**

运行：`golangci-lint run ./...`

预期：无警告或错误

- [ ] **Step 5: 检查构建**

运行：`go build ./cmd/server`

预期：构建成功

- [ ] **Step 6: 检查数据库迁移**

如果有测试数据库，运行迁移：

```bash
# 使用 psql 运行迁移
psql -h localhost -U postgres -d aitestos_test -f scripts/migrations/004_remove_prefix_deleted_at_rename_tables.sql
```

预期：迁移成功执行

- [ ] **Step 7: 提交最终更改**

如果有任何调整：

```bash
git add .
git commit -m "chore: final adjustments after simplification

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## 完成标准检查

在完成所有任务后，验证：

- [ ] 所有文档已更新
- [ ] 数据库迁移脚本已创建
- [ ] 领域层代码已更新
- [ ] 服务层代码已更新
- [ ] 持久层代码已更新
- [ ] 传输层代码已更新
- [ ] 所有测试通过
- [ ] `go fmt ./...` 无输出
- [ ] `go vet ./...` 无警告
- [ ] `golangci-lint run` 无警告
- [ ] 代码已提交到 Git

---

## 总计

- **Phase 1**: 文档更新 (3 个任务)
- **Phase 2**: 数据库迁移 (1 个任务)
- **Phase 3**: 领域层 (7 个任务)
- **Phase 4**: 服务层 (4 个任务)
- **Phase 5**: 持久层 (4 个任务)
- **Phase 6**: 传输层 (3 个任务)
- **Phase 7**: 验证 (1 个任务)

**总计**: 23 个任务

---

**计划创建时间**: 2026-05-05
**预计执行时间**: 4-6 小时
**风险等级**: 中等（涉及数据库表结构变更）
