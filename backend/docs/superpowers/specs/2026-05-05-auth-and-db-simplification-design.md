# 认证与数据库架构简化设计

**版本**: 1.0
**日期**: 2026-05-05
**状态**: 已批准

---

## 1. 概述

### 1.1 变更目标

简化系统架构，移除不必要的复杂度：
- **认证调整**：延长 access token 有效期至 7 天，保留 refresh token 机制
- **数据库简化**：移除软删除和前缀/缩写概念
- **表名规范化**：使用复数形式命名表

### 1.2 认证机制变更

| 项目 | 变更前 | 变更后 |
|------|--------|--------|
| Access Token 过期时间 | 15 分钟 | 7 天 |
| Refresh Token 过期时间 | 7 天 | 7 天（保持不变） |
| Refresh 机制 | 保留 | 保留 |
| /auth/refresh 端点 | 保留 | 保留 |

---

## 2. 数据库架构变更

### 2.1 表结构变更

#### 2.1.1 users 表

删除 `deleted_at` 列：

```sql
-- 变更前
CREATE TABLE users (
    -- ...
    deleted_at timestamp(3) with time zone
);

-- 变更后
CREATE TABLE users (
    -- ...
    -- deleted_at 列已删除
);
```

#### 2.1.2 project → projects 表

删除 `prefix` 列，重命名表：

```sql
-- 变更前
CREATE TABLE project (
    -- ...
    prefix varchar(4) NOT NULL UNIQUE,
    -- ...
);

-- 变更后
CREATE TABLE projects (
    -- ...
    -- prefix 列已删除
    -- ...
);
```

#### 2.1.3 module → modules 表

删除 `abbreviation` 列，重命名表：

```sql
-- 变更前
CREATE TABLE module (
    -- ...
    abbreviation varchar(4) NOT NULL,
    UNIQUE(project_id, abbreviation)
);

-- 变更后
CREATE TABLE modules (
    -- ...
    -- abbreviation 列已删除
    UNIQUE(project_id, name)
);
```

#### 2.1.4 test_case 表

删除 `number` 列：

```sql
-- 变更前
CREATE TABLE test_case (
    -- ...
    number varchar(32) NOT NULL UNIQUE,
    -- ...
);

-- 变更后
CREATE TABLE test_case (
    -- ...
    -- number 列已删除
    -- ...
);
```

### 2.2 迁移脚本

创建 `scripts/migrations/004_remove_prefix_deleted_at_rename_tables.sql`：

```sql
-- ============================================================
-- Migration 004: Remove prefix/deleted_at, rename tables
-- ============================================================

BEGIN;

-- 1. 删除 users 表的 deleted_at 列
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;

-- 2. 删除 project 表的 prefix 列
ALTER TABLE project DROP COLUMN IF EXISTS prefix;

-- 3. 删除 module 表的 abbreviation 列及相关约束
ALTER TABLE module DROP CONSTRAINT IF EXISTS module_abbreviation_key;
ALTER TABLE module DROP COLUMN IF EXISTS abbreviation;

-- 4. 删除 test_case 表的 number 列
ALTER TABLE test_case DROP COLUMN IF EXISTS test_case_number_key;
ALTER TABLE test_case DROP COLUMN IF EXISTS number;

-- 5. 重命名表
ALTER TABLE project RENAME TO projects;
ALTER TABLE module RENAME TO modules;

-- 6. 删除相关索引
DROP INDEX IF EXISTS idx_users_deleted_at;

COMMIT;
```

---

## 3. API 变更

### 3.1 项目 API

删除 `prefix` 字段：

```json
// 创建项目请求
{
  "name": "电商平台",
  "description": "电商平台测试项目"
}

// 项目响应
{
  "id": "...",
  "name": "电商平台",
  "description": "...",
  "createdAt": "2026-05-05T10:00:00Z"
}
```

### 3.2 模块 API

删除 `abbreviation` 字段：

```json
// 创建模块请求
{
  "name": "用户中心",
  "description": "用户相关功能模块"
}

// 模块响应
{
  "id": "...",
  "projectId": "...",
  "name": "用户中心",
  "description": "..."
}
```

### 3.3 测试用例 API

删除 `number` 字段：

```json
// 用例响应
{
  "id": "...",
  "moduleId": "...",
  "title": "密码错误5次后账号锁定",
  "steps": [...],
  "expected": {...}
}
```

### 3.4 删除错误码

- 20007: 项目前缀已存在
- 20008: 项目前缀格式无效
- 20009: 模块缩写已存在
- 40002: 用例编号已存在

---

## 4. 代码实现变更

### 4.1 领域层

**删除文件**：
- `internal/domain/project/prefix.go`
- `internal/domain/project/prefix_test.go`
- `internal/domain/project/abbreviation_test.go`
- `internal/domain/testcase/case_number.go`
- `internal/domain/testcase/case_number_test.go`

**更新聚合根/实体**：
- `Project`: 删除 `prefix` 字段和 `Prefix()` 访问器
- `Module`: 删除 `abbreviation` 字段和 `Abbreviation()` 访问器
- `TestCase`: 删除 `number` 字段和 `Number()` 访问器

### 4.2 服务层

**认证服务**：
```go
const accessTokenExpiry = 7 * 24 * time.Hour  // 15 分钟 → 7 天
```

**项目服务**：
- 删除 `CreateProjectRequest.Prefix`
- 删除前缀唯一性检查
- 更新 `NewProject` 调用

**模块服务**：
- 删除 `CreateModuleRequest.Abbreviation`
- 删除缩写唯一性检查
- 更新 `NewModule` 调用

**用例服务**：
- 删除编号生成逻辑
- 删除 `CountByDate` 调用

### 4.3 持久层

**Repository 查询更新**：
- 删除 `prefix`, `abbreviation`, `number`, `deleted_at` 列
- 表名更新：`project` → `projects`, `module` → `modules`
- 删除方法：`FindByPrefix`, `FindByAbbreviation`, `FindByNumber`, `CountByDate`

### 4.4 传输层

**Handler 响应更新**：
- 删除 `prefix`, `abbreviation`, `number` 字段
- 更新 `toProjectResponse`, `toModuleResponse`, `toTestCaseResponse`

---

## 5. 文档变更

| 文件 | 主要变更 |
|------|----------|
| `specs/aitestos_optimized.sql` | 删除列和重命名表 |
| `specs/openapi.yaml` | 删除 schema 字段和错误码 |
| `specs/001-core-functionality/spec.md` | 删除编号规则和前缀/缩写验证 |
| `specs/001-core-functionality/plan.md` | 删除值对象定义 |
| `specs/001-core-functionality/tasks.md` | 删除相关任务 |

---

## 6. 实施顺序

1. 更新文档
2. 创建数据库迁移脚本
3. 更新领域层代码
4. 更新服务层代码
5. 更新持久层代码
6. 更新传输层代码
7. 运行测试验证

---

## 7. 完成标准

- [ ] 所有文档更新完成
- [ ] 数据库迁移脚本创建并测试
- [ ] 代码变更实现并测试通过
- [ ] `go test ./...` 通过
- [ ] `go fmt ./...` 和 `go vet ./...` 通过
- [ ] `golangci-lint run` 通过
