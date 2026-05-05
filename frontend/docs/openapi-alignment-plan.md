# OpenAPI 与前端实现对齐方案

生成时间: 2026-05-05

## 一、背景

基于 `docs/openapi-review-diff.md` 的审查结果，本文档提供具体的对齐行动方案。

## 二、优先级分类

### 🔴 P0 - 阻塞性问题（必须立即解决）

#### 1. 移除已废弃的字段

根据 git log 显示，后端正在进行重构，移除 `prefix`、`abbreviation`、`number` 字段。前端需要同步移除这些字段。

**受影响的文件**:
```
src/types/api.ts
  - Project.prefix
  - Module.abbreviation
  - TestCase.number
  - CreateTestCaseRequest (无影响，已不包含)
  - ProjectDetail (继承自 Project)
```

**具体变更**:
```typescript
// 修改前
interface Project {
  id: string
  name: string
  prefix: string              // ❌ 删除
  description: string
  createdAt: string
  updatedAt: string
}

// 修改后
interface Project {
  id: string
  name: string
  description: string
  createdAt: string
  updatedAt: string
}
```

```typescript
// 修改前
interface Module {
  id: string
  projectId: string
  name: string
  abbreviation: string        // ❌ 删除
  createdAt: string
  updatedAt: string
  caseCount?: number
}

// 修改后
interface Module {
  id: string
  projectId: string
  name: string
  description?: string        // ✅ 添加（与 OpenAPI 对齐）
  createdAt: string
  updatedAt: string
  caseCount?: number
}
```

```typescript
// 修改前
interface TestCase {
  id: string
  moduleId: string
  userId: string
  number: string              // ❌ 删除
  title: string
  // ...
}

// 修改后
interface TestCase {
  id: string
  moduleId: string
  userId: string
  title: string
  // ...
}
```

**影响范围评估**:
- 需要检查所有使用这些字段的组件和 hooks
- 需要更新表单验证 schema
- 需要更新测试文件

---

#### 2. 修正 Module 字段结构

OpenAPI 定义 Module 包含 `createdBy` 字段，前端缺失。

**具体变更**:
```typescript
// 修改前
interface Module {
  id: string
  projectId: string
  name: string
  description?: string
  createdAt: string
  updatedAt: string
  caseCount?: number
}

// 修改后
interface Module {
  id: string
  projectId: string
  name: string
  description?: string
  createdBy?: string          // ✅ 添加
  createdAt: string
  updatedAt: string
  caseCount?: number
}
```

---

#### 3. 修正 ProjectStatistics 结构

OpenAPI 定义的 ProjectStatistics 与前端的 ProjectStats 结构完全不同。

**OpenAPI 定义**:
```yaml
ProjectStatistics:
  module_count: integer
  case_count: integer
  document_count: integer
  pass_rate: number (double)
  coverage_rate: number (double)
  ai_generated_count: integer
  recent_tasks: array[TaskSummary]
  pass_rate_trend: array[TrendData]
  updated_at: string (date-time)
```

**前端当前定义**:
```typescript
interface ProjectStats {
  totalCases: number
  passRate: number
  coverage: number
  aiGeneratedCount: number
  trend: Array<{
    date: string
    passRate: number
  }>
}
```

**对齐方案**:
```typescript
// 添加 TaskSummary 和 TrendData 类型
interface TaskSummary {
  id: string
  status: TaskStatus
  resultSummary?: TaskResultSummary
  createdAt: string
}

interface TaskResultSummary {
  totalDrafts: number
  confirmedCount: number
  rejectedCount: number
}

interface TrendData {
  date: string  // YYYY-MM-DD
  rate: number  // 通过率百分比
}

// 重新定义 ProjectStats
interface ProjectStats {
  moduleCount: number
  caseCount: number
  documentCount: number
  passRate: number
  coverageRate: number
  aiGeneratedCount: number
  recentTasks: TaskSummary[]
  passRateTrend: TrendData[]
  updatedAt: string
}
```

---

#### 4. 修正 CreateTaskRequest.sceneTypes 类型

OpenAPI 定义 `scene_types` 为数组，前端定义为单值。

**具体变更**:
```typescript
// 修改前
interface CreateTaskRequest {
  projectId: string
  moduleId: string
  prompt: string
  count?: number
  caseType?: CaseType
  priority?: Priority
  sceneType?: SceneType       // ❌ 单值
}

// 修改后
interface CreateTaskRequest {
  projectId: string
  moduleId: string
  prompt: string
  caseCount?: number          // ✅ 重命名为 caseCount
  sceneTypes?: SceneType[]    // ✅ 改为数组
  priority?: Priority
  caseType?: CaseType
}
```

**同时需要更新枚举名**:
```typescript
// src/types/enums.ts
/**
 * Test scene type for AI generation
 */
export type SceneType = 'positive' | 'negative' | 'boundary'
```

---

### 🟡 P1 - 重要问题（应尽快解决）

#### 1. ProjectDetail 字段对齐

**问题**:
- OpenAPI: `document_count` (文档数量)
- 前端: `draftCount` (草稿数量)

**决策**: 需要与后端确认，是前端显示草稿数还是文档数

**临时方案**: 保留前端字段名（draftCount），但确保后端 API 返回相应数据

---

#### 2. CaseDraft.projectId 字段

**问题**: 前端定义了 `projectId` 字段，但 OpenAPI 的 CaseDraft 中未定义

**决策**: 保留此字段（用于列表筛选），作为前端扩展

---

#### 3. 统一命名风格

OpenAPI 使用 snake_case，前端使用 camelCase。这是合理的，但需要确保类型转换一致。

**行动**: 在 `src/lib/request.ts` 中确保使用 axios 的 `transformResponse` 统一处理

---

### 🟢 P2 - 优化问题（可延后处理）

#### 1. 补充 OpenAPI 中缺失的端点定义

前端实现了 3 个扩展端点，建议补充到 OpenAPI 规范中：

```yaml
# 建议添加到 openapi.yaml

DELETE /projects/{id}/configs/{key}:
  tags: [ProjectConfig]
  summary: 删除配置项
  operationId: deleteProjectConfig
  security:
    - bearerAuth: []
  parameters:
    - name: id
      in: path
      required: true
      schema:
        type: string
        format: uuid
    - name: key
      in: path
      required: true
      schema:
        type: string
  responses:
    '204':
      description: 删除成功

DELETE /plans/{planId}/results/{caseId}:
  tags: [TestPlans]
  summary: 删除测试结果
  operationId: deleteTestResult
  security:
    - bearerAuth: []
  parameters:
    - name: planId
      in: path
      required: true
      schema:
        type: string
        format: uuid
    - name: caseId
      in: path
      required: true
      schema:
        type: string
        format: uuid
  responses:
    '204':
      description: 删除成功

GET /generation/drafts/count:
  tags: [Generation]
  summary: 获取待处理草稿数量
  operationId: getPendingDraftsCount
  security:
    - bearerAuth: []
  responses:
    '200':
      description: 成功
      content:
        application/json:
          schema:
            type: object
            properties:
              count:
                type: integer
```

---

## 三、实施计划

### 阶段 1: 类型定义修正（1-2 天）

1. **修改 `src/types/api.ts`**
   - [ ] 移除 Project.prefix
   - [ ] 移除 Module.abbreviation，添加 description 和 createdBy
   - [ ] 移除 TestCase.number
   - [ ] 重构 ProjectStats 结构
   - [ ] 修正 CreateTaskRequest

2. **修改 `src/types/enums.ts`**
   - [ ] 确保 SceneType 导出正确

3. **运行类型检查**
   ```bash
   make type-check
   ```

### 阶段 2: 服务层更新（1-2 天）

1. **更新 services 层**
   - [ ] 检查并更新所有使用已移除字段的 API 调用
   - [ ] 更新请求/响应类型转换

2. **更新测试**
   - [ ] 更新所有单元测试中的 mock 数据
   - [ ] 确保 MSW handlers 与新类型一致

### 阶段 3: 组件层修复（2-3 天）

1. **检查所有使用受影响字段的组件**
   ```bash
   # 搜索使用这些字段的地方
   grep -r "\.prefix" src/
   grep -r "\.abbreviation" src/
   grep -r "\.number" src/features/testcases/
   ```

2. **更新组件和 hooks**
   - [ ] 移除对已废弃字段的引用
   - [ ] 更新表单验证 schema
   - [ ] 更新显示逻辑

3. **更新集成测试**
   - [ ] 确保 E2E 测试使用新的数据结构

### 阶段 4: OpenAPI 规范更新（1 天）

1. **补充缺失的端点定义**
   - [ ] 添加 DELETE /configs/{key}
   - [ ] 添加 DELETE /plans/{planId}/results/{caseId}
   - [ ] 添加 GET /drafts/count

2. **验证 OpenAPI 规范**
   ```bash
   # 使用 spectral 或其他工具验证
   npm install -g @stoplight/spectral-cli
   spectral lint specs/openapi.yaml
   ```

---

## 四、风险评估

### 高风险区域

1. **ProjectDetail 相关功能**
   - 风险: 如果后端仍在返回 `prefix`，移除后可能导致显示异常
   - 缓解: 先检查后端实际返回数据，再决定是否移除

2. **Module.abbreviation → description**
   - 风险: 业务逻辑可能依赖 abbreviation
   - 缓解: 全面搜索使用该字段的地方

3. **TestCase.number**
   - 风险: 如果业务流程中使用编号作为唯一标识
   - 缓解: 确认后端 API 不再返回此字段

### 回滚计划

如果对齐后发现严重问题，可以：
1. 暂时保留旧字段，标记为 `@deprecated`
2. 在类型中添加可选的旧字段，逐步迁移

---

## 五、验收标准

### 类型层面
- [ ] `make type-check` 通过，无 any 错误
- [ ] 所有枚举类型与 OpenAPI 一致

### 运行时层面
- [ ] 所有 API 调用正常响应
- [ ] 表单提交验证通过
- [ ] 数据显示正确

### 测试层面
- [ ] 所有单元测试通过
- [ ] E2E 测试覆盖关键路径
- [ ] MSW handlers 正确 mock 新数据结构

---

## 六、后续建议

1. **建立自动化类型生成流程**
   - 考虑使用 openapi-typescript 从 OpenAPI 规范自动生成类型
   - 配置 CI 检查类型一致性

2. **添加契约测试**
   - 使用 MSW 验证前端请求与 OpenAPI 规范的一致性
   - 添加后端 API 变更检测机制

3. **文档更新**
   - 更新 CLAUDE.md 中的类型定义说明
   - 更新组件文档，反映新的数据结构
