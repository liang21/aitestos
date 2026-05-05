# OpenAPI 规范与前端实现差异审查

生成时间: 2026-05-05

## 一、枚举类型对比

### ✅ 一致的枚举
| 枚举名 | OpenAPI | 前端 | 状态 |
|--------|---------|------|------|
| CaseStatus | unexecuted, pass, block, fail | unexecuted, pass, block, fail | ✅ |
| CaseType | functionality, performance, api, ui, security | functionality, performance, api, ui, security | ✅ |
| PlanStatus | draft, active, completed, archived | draft, active, completed, archived | ✅ |
| Priority | P0, P1, P2, P3 | P0, P1, P2, P3 | ✅ |
| ResultStatus | pass, fail, block, skip | pass, fail, block, skip | ✅ |
| TaskStatus | pending, processing, completed, failed | pending, processing, completed, failed | ✅ |
| DraftStatus | pending, confirmed, rejected | pending, confirmed, rejected | ✅ |
| DocumentType | prd, figma, api_spec, swagger, markdown | prd, figma, api_spec, swagger, markdown | ✅ |
| DocumentStatus | pending, processing, completed, failed | pending, processing, completed, failed | ✅ |
| UserRole | super_admin, admin, normal | super_admin, admin, normal | ✅ |

### ✅ 前端独有的枚举
- `Confidence` (high, medium, low) - 前端扩展，合理
- `SceneType` (positive, negative, boundary) - 前端扩展，合理

---

## 二、类型定义差异

### 1. Project (项目)

#### OpenAPI Schema
```yaml
Project:
  properties:
    id: string (uuid)
    name: string
    description: string
    createdAt: string (date-time)
    updatedAt: string (date-time)
```

#### 前端类型
```typescript
interface Project {
  id: string
  name: string
  prefix: string              // ❌ 前端独有
  description: string
  createdAt: string
  updatedAt: string
}
```

**差异**: `prefix` 字段存在于前端但不在 OpenAPI 中

---

### 2. ProjectDetail (项目详情)

#### OpenAPI Schema
```yaml
ProjectDetail:
  properties:
    id: string
    name: string
    description: string
    createdAt: string
    updatedAt: string
    module_count: integer
    case_count: integer
    document_count: integer
```

#### 前端类型
```typescript
interface ProjectDetail extends Project {
  moduleCount: number
  caseCount: number
  draftCount: number         // ❌ 前端独有
}
```

**差异**:
1. 命名风格: OpenAPI 用 snake_case (module_count)，前端用 camelCase (moduleCount) - 这是合理的，前端遵循 JS 约定
2. `document_count` vs `draftCount` - 字段语义不同
3. 前端继承了 Project 的 `prefix` 字段

---

### 3. Module (模块)

#### OpenAPI Schema
```yaml
Module:
  properties:
    id: string (uuid)
    projectId: string (uuid)
    name: string
    description: string
    createdBy: string (uuid)
    createdAt: string (date-time)
    updatedAt: string (date-time)
```

#### 前端类型
```typescript
interface Module {
  id: string
  projectId: string
  name: string
  abbreviation: string        // ❌ 前端独有
  createdAt: string
  updatedAt: string
  caseCount?: number          // ❌ 前端独有
}
```

**差异**:
1. `abbreviation` vs `description` - 完全不同的字段
2. 前端缺少 `createdBy`
3. 前端多了 `caseCount` (关联数据)

---

### 4. TestCase (测试用例)

#### OpenAPI Schema
```yaml
TestCase:
  properties:
    id: string (uuid)
    moduleId: string (uuid)
    userId: string (uuid)
    title: string
    preconditions: array[string]
    steps: array[string]
    expected: object
    caseType: CaseType
    priority: Priority
    status: CaseStatus
    aiMetadata: object
    createdAt: string (date-time)
    updatedAt: string (date-time)
```

#### 前端类型
```typescript
interface TestCase {
  id: string
  moduleId: string
  userId: string
  number: string              // ❌ 前端独有，格式: {prefix}-{abbreviation}-{YYYYMMDD}-{001}
  title: string
  preconditions: string[]
  steps: string[]
  expected: Record<string, unknown>
  caseType: CaseType
  priority: Priority
  status: CaseStatus
  aiMetadata?: AiMetadata
  createdAt: string
  updatedAt: string
  // Joined fields
  moduleName?: string
  projectName?: string
  projectPrefix?: string
  createdByName?: string
}
```

**差异**:
1. `number` 字段存在于前端但不在 OpenAPI 中
2. 前端扩展了关联字段 (moduleName, projectName 等) - 这是合理的

---

### 5. GenerationTask (生成任务)

#### OpenAPI Schema
```yaml
GenerationTask:
  properties:
    id: string (uuid)
    projectId: string (uuid)
    moduleId: string (uuid)
    status: TaskStatus
    prompt: string
    result: object
    createdAt: string (date-time)
    updatedAt: string (date-time)
```

#### 前端类型
```typescript
interface GenerationTask {
  id: string
  projectId: string
  moduleId: string
  status: TaskStatus
  prompt: string
  result: GenerationTaskResult | null
  createdAt: string
  updatedAt: string
  // Joined fields
  projectName?: string
  moduleName?: string
  createdByName?: string
}
```

**差异**: 结构一致，前端扩展了合理的关联字段

---

### 6. CreateGenerationTaskRequest

#### OpenAPI Schema
```yaml
CreateGenerationTaskRequest:
  required: [project_id, module_id, prompt]
  properties:
    project_id: string (uuid)
    module_id: string (uuid)
    prompt: string
      minLength: 10
    case_count: integer
      minimum: 1
      maximum: 20
      default: 5
    scene_types: array[string]
      enum: [positive, negative, boundary]
    priority: Priority
    case_type: CaseType
```

#### 前端类型
```typescript
interface CreateTaskRequest {
  projectId: string
  moduleId: string
  prompt: string
  count?: number
  caseType?: CaseType
  priority?: Priority
  sceneType?: SceneType      // ❌ 单数 vs 复数
}
```

**差异**:
1. `case_count` vs `count`
2. `scene_types` (数组) vs `sceneType` (单值) - 类型不一致

---

### 7. ProjectStatistics (项目统计)

#### OpenAPI Schema
```yaml
ProjectStatistics:
  properties:
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

#### 前端类型
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

**差异**: 完全不同的结构，需要重新对齐

---

### 8. CaseDraft (草稿)

#### OpenAPI Schema
```yaml
CaseDraft:
  properties:
    id: string (uuid)
    taskId: string (uuid)
    title: string
    preconditions: array[string]
    steps: array[string]
    expected: object
    caseType: CaseType
    priority: Priority
    status: DraftStatus
    feedback: string
    createdAt: string (date-time)
    updatedAt: string (date-time)
```

#### 前端类型
```typescript
interface CaseDraft {
  id: string
  taskId: string
  projectId: string          // ❌ 前端独有
  title: string
  preconditions: string[]
  steps: string[]
  expected: Record<string, unknown>
  caseType: CaseType
  priority: Priority
  status: DraftStatus
  feedback?: string
  aiMetadata?: { ... }
  createdAt: string
  updatedAt: string
  // Joined fields
  projectName?: string
  moduleName?: string
}
```

**差异**:
1. 前端多了 `projectId`
2. 前端扩展了 `aiMetadata`

---

### 9. CreateTestCaseRequest

#### OpenAPI Schema
```yaml
CreateTestCaseRequest:
  required: [module_id, title, steps, expected, case_type, priority]
  properties:
    module_id: string (uuid)
    title: string
      minLength: 2
      maxLength: 500
    preconditions: array[string]
    steps: array[string]
      minItems: 1
    expected: object
    case_type: CaseType
    priority: Priority
```

#### 前端类型
```typescript
interface CreateTestCaseRequest {
  moduleId: string
  title: string
  preconditions: string[]
  steps: string[]
  expected: Record<string, unknown>
  caseType: CaseType
  priority: Priority
}
```

**差异**: 结构一致，字段映射正确

---

## 三、API 端点覆盖度

### 1. System (系统接口)

| OpenAPI 端点 | 方法 | 前端实现 | 状态 |
|--------------|------|----------|------|
| /health | GET | ❌ 未实现 | ⚪ 无需实现（基础设施） |
| /metrics | GET | ❌ 未实现 | ⚪ 无需实现（运维监控） |

**说明**: 健康检查和指标端点由基础设施层处理，前端无需实现

---

### 2. Auth (认证)

| OpenAPI 端点 | 方法 | 前端实现 | 状态 |
|--------------|------|----------|------|
| /auth/register | POST | authApi.register | ✅ 已实现 |
| /auth/login | POST | authApi.login | ✅ 已实现 |
| /auth/refresh | POST | authApi.refresh | ✅ 已实现 |

**覆盖率**: 3/3 (100%)

---

### 3. Projects (项目管理)

| OpenAPI 端点 | 方法 | 前端实现 | 状态 |
|--------------|------|----------|------|
| /projects | GET | projectsApi.list | ✅ 已实现 |
| /projects | POST | projectsApi.create | ✅ 已实现 |
| /projects/{id} | GET | projectsApi.get | ✅ 已实现 |
| /projects/{id} | PUT | projectsApi.update | ✅ 已实现 |
| /projects/{id} | DELETE | projectsApi.delete | ✅ 已实现 |
| /projects/{id}/stats | GET | projectsApi.getStats | ✅ 已实现 |

**覆盖率**: 6/6 (100%)

---

### 4. Modules (模块管理)

| OpenAPI 端点 | 方法 | 前端实现 | 状态 |
|--------------|------|----------|------|
| /projects/{id}/modules | GET | modulesApi.list | ✅ 已实现 |
| /projects/{id}/modules | POST | modulesApi.create | ✅ 已实现 |
| /modules/{id} | PUT | modulesApi.update | ✅ 已实现 |
| /modules/{id} | DELETE | modulesApi.delete | ✅ 已实现 |

**覆盖率**: 4/4 (100%)

---

### 5. ProjectConfig (项目配置)

| OpenAPI 端点 | 方法 | 前端实现 | 状态 |
|--------------|------|----------|------|
| /projects/{id}/configs | GET | configsApi.list | ✅ 已实现 |
| /projects/{id}/configs/{key} | PUT | configsApi.set | ✅ 已实现 |
| /projects/{id}/configs/{key} | DELETE | configsApi.delete | ⚠️ **API 中未定义** | ⚠️ 前端扩展 |
| /projects/{id}/configs/import | POST | configsApi.import | ✅ 已实现 |
| /projects/{id}/configs/export | GET | configsApi.export | ✅ 已实现 |

**覆盖率**: 4/4 (OpenAPI 规范) + 1 个前端扩展

**注意**: 前端实现了 DELETE /configs/{key}，但 OpenAPI 中未定义此端点

---

### 6. TestCases (测试用例)

| OpenAPI 端点 | 方法 | 前端实现 | 状态 |
|--------------|------|----------|------|
| /testcases | GET | testcasesApi.list | ✅ 已实现 |
| /testcases | POST | testcasesApi.create | ✅ 已实现 |
| /testcases/{id} | GET | testcasesApi.get | ✅ 已实现 |
| /testcases/{id} | PUT | testcasesApi.update | ✅ 已实现 |
| /testcases/{id} | DELETE | testcasesApi.delete | ✅ 已实现 |

**覆盖率**: 5/5 (100%)

---

### 7. TestPlans (测试计划)

| OpenAPI 端点 | 方法 | 前端实现 | 状态 |
|--------------|------|----------|------|
| /plans | GET | plansApi.list | ✅ 已实现 |
| /plans | POST | plansApi.create | ✅ 已实现 |
| /plans/{id} | GET | plansApi.get | ✅ 已实现 |
| /plans/{id} | PUT | plansApi.update | ✅ 已实现 |
| /plans/{id} | PATCH | plansApi.updateStatus | ✅ 已实现 |
| /plans/{id} | DELETE | plansApi.delete | ✅ 已实现 |
| /plans/{id}/cases | POST | plansApi.addCases | ✅ 已实现 |
| /plans/{id}/cases/{caseId} | DELETE | plansApi.removeCase | ✅ 已实现 |
| /plans/{id}/results | GET | plansApi.getResults | ✅ 已实现 |
| /plans/{id}/results | POST | plansApi.recordResult | ✅ 已实现 |

**覆盖率**: 10/10 (100%)

**注意**: 前端还实现了 DELETE /plans/{planId}/results/{caseId}，但 OpenAPI 中未定义

---

### 8. Generation (AI 生成)

| OpenAPI 端点 | 方法 | 前端实现 | 状态 |
|--------------|------|----------|------|
| /generation/tasks | GET | generationApi.listTasks | ✅ 已实现 |
| /generation/tasks | POST | generationApi.createTask | ✅ 已实现 |
| /generation/tasks/{id} | GET | generationApi.getTask | ✅ 已实现 |
| /generation/tasks/{id}/drafts | GET | generationApi.getTaskDrafts | ✅ 已实现 |
| /generation/drafts | GET | draftsApi.getDrafts | ✅ 已实现 |
| /generation/drafts/{id}/confirm | POST | draftsApi.confirmDraft | ✅ 已实现 |
| /generation/drafts/{id}/reject | POST | draftsApi.rejectDraft | ✅ 已实现 |
| /generation/drafts/batch-confirm | POST | draftsApi.batchConfirm | ✅ 已实现 |

**覆盖率**: 8/8 (100%)

**注意**: 前端还实现了 GET /generation/drafts/count，但 OpenAPI 中未定义（用于 badge 计数）

---

### 9. Documents (知识库文档)

| OpenAPI 端点 | 方法 | 前端实现 | 状态 |
|--------------|------|----------|------|
| /knowledge/documents | GET | documentsApi.list | ✅ 已实现 |
| /knowledge/documents | POST | documentsApi.create | ✅ 已实现 |
| /knowledge/documents/{id} | GET | documentsApi.get | ✅ 已实现 |
| /knowledge/documents/{id} | DELETE | documentsApi.delete | ✅ 已实现 |
| /knowledge/documents/{id}/chunks | GET | documentsApi.getChunks | ✅ 已实现 |

**覆盖率**: 5/5 (100%)

---

## 四、总体覆盖度统计

| API 分组 | OpenAPI 端点数 | 已实现 | 覆盖率 | 前端扩展 |
|----------|----------------|--------|--------|----------|
| System | 2 | 0 | 0% (无需) | - |
| Auth | 3 | 3 | 100% | - |
| Projects | 6 | 6 | 100% | - |
| Modules | 4 | 4 | 100% | - |
| ProjectConfig | 4 | 5 | 100% | +1 (DELETE) |
| TestCases | 5 | 5 | 100% | - |
| TestPlans | 10 | 11 | 100% | +1 (DELETE result) |
| Generation | 8 | 9 | 100% | +1 (GET count) |
| Documents | 5 | 5 | 100% | - |
| **总计** | **47** | **48** | **100%** | **+3** |

### 前端扩展的端点

1. **DELETE /projects/{id}/configs/{key}** — 删除配置项（configsApi.delete）
2. **DELETE /plans/{planId}/results/{caseId}** — 删除测试结果（plansApi.deleteResult）
3. **GET /generation/drafts/count** — 获取待处理草稿数量（draftsApi.getPendingCount）

**建议**: 这些扩展端点功能合理，但需要在 OpenAPI 规范中补充定义

---

## 四、关键问题总结

### 🔴 高优先级问题
1. **Project.prefix** - 前端有但 OpenAPI 无，需确认后端是否已移除
2. **Module.abbreviation vs description** - 字段名称和语义不同
3. **TestCase.number** - 前端有但 OpenAPI 无
4. **ProjectStatistics** - 结构完全不同，需重新设计
5. **CreateTaskRequest.sceneType vs scene_types** - 单值 vs 数组

### 🟡 中优先级问题
1. **命名风格** - OpenAPI snake_case vs 前端 camelCase (这是合理的，但需保持一致)
2. **关联字段** - 前端扩展的 joined fields 需要文档化
3. **ProjectDetail.document_count vs draftCount** - 字段语义不同

### 🟢 低优先级
1. **可选字段** - `aiMetadata?` 等可选性的差异
2. **类型细化** - `object` vs `Record<string, unknown>`

---

## 五、建议行动

### 立即行动
1. 与后端确认 `prefix`, `abbreviation`, `number` 字段的状态
2. 统一 `ProjectStatistics` 结构
3. 修正 `scene_types` 类型定义

### 短期行动
1. 移除已废弃的字段
2. 添加缺失的字段
3. 统一命名风格转换规则

### 长期行动
1. 建立自动化 OpenAPI 类型生成流程
2. 添加类型一致性测试
