---
title: 核心功能规格说明书
version: 1.0.0
author: 需求编译器
status: draft
created: 2026-04-16
---

# 核心功能规格说明书

> 本文档基于 PRD v2.1、OpenAPI 3.0.3、UX 设计规范 v1.0、前端详细设计 v1.1 编译生成，
> 覆盖 Aitestos 智能测试管理平台的核心业务流程。

---

## 1. 用户故事

### US-C01 核心故事：AI 用例生成与确认全流程

```
作为 测试团队负责人
我想要 描述测试需求并由 AI 自动生成测试用例草稿，审核后一键转为正式用例
以便于 将用例编写效率提升 60% 以上，同时保持对需求变更的快速响应能力

验收标准：
- 用户可在指定项目/模块下发起 AI 生成任务
- 生成完成后可在草稿箱查看结果，包含标题、步骤、预期结果
- 草稿支持编辑、确认（转正式用例）、拒绝（带原因反馈）
- 确认后自动生成业务编号（格式：项目前缀-模块缩写-日期-序号）
- 草稿详情页右侧展示引用来源（知识库文档块），支持跳转原文
```

### US-C02 边缘场景：知识库内容不足时的降级生成

```
作为 测试工程师
我想要 在知识库文档不足的情况下仍能发起 AI 生成任务
以便于 在项目初期文档尚未完备时即可开始用例设计工作

验收标准：
- 当 RAG 检索结果 < 1 个文档块时，弹窗提示"知识库内容不足，生成质量可能较低"
- 用户确认后允许继续生成，置信度强制标记为"低"
- 当项目完全无知识库文档时，禁止发起生成任务，提示"请先上传需求文档"
- 降级生成的草稿仍可在草稿箱正常编辑和确认
```

### US-C03 边缘场景：JWT Token 过期无感刷新

```
作为 已登录用户
我想要 在 Token 过期时系统自动刷新而无需重新登录
以便于 在长时间操作（如等待 AI 生成结果）期间不会因会话过期而丢失工作状态

验收标准：
- 前端检测到 401 响应时自动使用 refresh_token 尝试刷新
- 刷新成功后自动重放原始请求
- 刷新失败时跳转登录页，保留当前路由 state 以便登录后回跳
- 并发请求时仅触发一次刷新，其余请求排队等待刷新完成后重放
```

---

## 2. 功能性需求

### 2.1 认证与会话管理

- 用户注册：输入用户名、邮箱、密码、角色，密码 bcrypt 加密存储
- 用户登录：邮箱 + 密码认证，返回 JWT access_token（2h 有效期）和 refresh_token
- Token 刷新：使用 refresh_token 换取新的 token 对，实现无感续期
- 路由守卫：前端通过 `RouteGuard` 组件拦截未认证访问，解析 JWT `exp` 字段判断过期

### 2.2 项目与模块管理

- 项目 CRUD：创建（名称 + 前缀 + 描述）、查看列表（分页/搜索）、编辑、软删除
- 模块管理：在项目下创建模块（名称 + 缩写），删除模块时级联删除关联用例
- 项目仪表盘：展示统计数据（用例总数、通过率、覆盖率、AI 生成数）和趋势图表
- 项目配置：支持 KV 配置管理（LLM 模型、温度参数等），支持 JSON 批量导入/导出

### 2.3 知识库文档管理

- 文档上传：支持 PRD（Word/PDF/Markdown）、API 规范（OpenAPI YAML/JSON）等类型
- 异步解析：上传后进入消息队列，后台解析文档并分块（DocumentChunk），生成向量存入 Milvus
- 文档状态：pending → processing → completed / failed
- Figma 集成：通过个人令牌或 OAuth 连接，导入指定页面节点树内容
- 文档分块查看：在文档详情页展示分块列表，支持查看块内容

### 2.4 AI 用例生成（核心）

- 生成任务发起：选择目标模块，输入需求描述（≥10 字），可选数量（1-20）、类型、优先级、场景类型
- RAG 检索流程：需求描述 → Embedding 向量化 → Milvus Top-K 检索 → 组装 Prompt → LLM 生成
- 置信度判定：高（≥2 个块，相似度 > 0.8）、中（1 个块，0.5-0.8）、低（无匹配）
- 草稿管理：生成结果存入草稿箱，支持按项目/模块/时间分组查看
- 草稿确认：编辑后确认转正，自动生成编号；支持批量确认（最多 50 条）
- 草稿拒绝：选择原因（重复/无关/低质量/其他）+ 文字反馈

### 2.5 测试用例管理

- 用例属性：编号（系统生成）、标题、前置条件、步骤、预期结果、类型、优先级、状态、AI 元数据
- 用例 CRUD：创建、编辑、软删除、查看详情（含 AI 来源追溯）
- 用例列表：分页展示，支持按项目/模块/状态/类型/优先级/关键词筛选
- 用例编号规则：`{项目前缀}-{模块缩写}-{YYYYMMDD}-{001-999}`

### 2.6 测试计划与执行

- 计划 CRUD：创建（关联项目 + 用例列表）、编辑、删除、查看详情
- 计划状态：draft → active → completed → archived
- 用例关联：向计划添加/移除用例
- 结果录入：对计划中的用例录入执行结果（pass/fail/block/skip + 备注）
- 统计汇总：计划详情页展示执行统计（总数/通过/失败/阻塞/跳过/未执行）

### 2.7 业务逻辑流转

```
注册/登录 → 获取 Token → 创建项目 → 创建模块
                                    ↓
                            上传知识库文档 → 异步解析 → 分块 + 向量化
                                    ↓
                      发起 AI 生成任务 → RAG 检索 → LLM 生成 → 草稿入库
                                    ↓
                      草稿审核 → 编辑 → 确认转正 → 正式用例（自动编号）
                                    ↓
                      创建测试计划 → 关联用例 → 执行测试 → 录入结果
                                    ↓
                      项目仪表盘 ← 统计数据（通过率/覆盖率/趋势）
```

---

## 3. 非功能性需求

### 3.1 架构解耦要求

| 解耦维度           | 要求                                                                   | 实现方式                                                                                 |
| ------------------ | ---------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| **API 抽象层**     | 前端不直接依赖 Axios 实例，通过 `src/api/` 模块封装所有接口调用        | 每个业务域独立文件（auth.ts、projects.ts、testcases.ts 等），统一导出类型安全的函数      |
| **状态管理独立**   | Store 按职责拆分，单一职责，不跨 Store 直接引用                        | Zustand 分为 useAuthStore、useAppStore、useProjectStore、useDraftStore，各自管理独立状态 |
| **路由与页面解耦** | 路由配置（router/index.tsx）仅定义映射，页面组件通过 `lazy()` 懒加载   | 路由守卫（RouteGuard）独立组件，权限逻辑与页面渲染分离                                   |
| **组件分层**       | 布局组件（layout/）、业务组件（business/）、页面组件（pages/）三层分离 | 通用组件封装为 SearchTable、StatusTag、ArrayEditor 等，页面组件组合使用                  |
| **类型定义集中**   | 所有 API 请求/响应类型、枚举、模型定义集中在 `src/types/`              | api.ts（请求响应）、enums.ts（枚举映射）、models.ts（数据模型）                          |

### 3.2 错误处理策略

| 错误类型               | 处理策略                                                                         | 用户感知                                                  |
| ---------------------- | -------------------------------------------------------------------------------- | --------------------------------------------------------- |
| **全局 HTTP 错误拦截** | Axios 拦截器统一处理：401 → 自动刷新 Token；403 → 提示无权限；5xx → 提示服务异常 | 全局 Message 提示，不中断当前页面                         |
| **Token 过期**         | 检测 JWT `exp` 字段或 401 响应，自动使用 refresh_token 刷新，刷新失败跳转登录    | 用户无感知（成功时）；跳转登录并保留来源路由（失败时）    |
| **表单验证错误**       | Arco Form 组件内置验证，字段级错误提示                                           | 字段下方红色提示文字，提交按钮禁用                        |
| **AI 生成超时**        | LLM 单次调用 30s 超时，总任务 120s 超时                                          | 任务状态置为 failed，提示"生成超时，请重试"，支持手动重试 |
| **文档解析失败**       | 后台异步处理，失败后文档状态置为 failed                                          | 列表页显示"失败"状态 Tag，支持重新解析                    |
| **网络断开**           | Axios 超时 + 网络错误拦截                                                        | 全局 Message "网络连接失败，请检查网络设置"               |
| **并发请求排队**       | Token 刷新期间并发请求放入等待队列，刷新完成后统一重放                           | 用户无感知                                                |

### 3.3 性能要求

| 指标                    | 目标值                        |
| ----------------------- | ----------------------------- |
| API 响应 P99            | < 500ms                       |
| 用例列表加载（1000 条） | < 2s                          |
| 文档解析（50 页 PDF）   | < 30s                         |
| AI 生成（5 个用例）     | < 60s                         |
| 页面首屏渲染            | < 1.5s（路由懒加载 + 骨架屏） |

### 3.4 安全要求

| 要求         | 实现                                                    |
| ------------ | ------------------------------------------------------- |
| JWT 认证     | access_token 2h 有效期，存储于 localStorage             |
| 密码存储     | bcrypt 加密                                             |
| 前端权限控制 | RouteGuard 拦截，admin+ 路由需 requireAdmin             |
| XSS 防护     | React 默认转义 + 输入内容不使用 dangerouslySetInnerHTML |
| CSRF 防护    | JWT Bearer 方式，不依赖 Cookie                          |

---

## 4. 验收标准

### TC-01 用户登录

| 项目         | 内容                                                            |
| ------------ | --------------------------------------------------------------- |
| **用例标题** | 使用有效凭据成功登录                                            |
| **前置条件** | 已注册用户（email: `admin@aitestos.com`，password: `Test1234`） |
| **操作步骤** | 1. 访问 `/login` 页面                                           |
|              | 2. 输入邮箱 `admin@aitestos.com`                                |
|              | 3. 输入密码 `Test1234`                                          |
|              | 4. 点击"登录"按钮                                               |
| **预期结果** | 1. 请求 `POST /api/v1/auth/login` 返回 200                      |
|              | 2. token 和 refresh_token 写入 localStorage                     |
|              | 3. 跳转至 `/projects` 页面                                      |
|              | 4. 侧边栏显示当前用户头像和用户名                               |

### TC-02 登录失败

| 项目         | 内容                                       |
| ------------ | ------------------------------------------ |
| **用例标题** | 使用错误密码登录失败                       |
| **前置条件** | 已注册用户（email: `admin@aitestos.com`）  |
| **操作步骤** | 1. 访问 `/login` 页面                      |
|              | 2. 输入邮箱 `admin@aitestos.com`           |
|              | 3. 输入错误密码 `WrongPassword`            |
|              | 4. 点击"登录"按钮                          |
| **预期结果** | 1. 请求 `POST /api/v1/auth/login` 返回 401 |
|              | 2. 页面显示错误提示"邮箱或密码错误"        |
|              | 3. 页面停留在 `/login`，不跳转             |

### TC-03 创建项目

| 项目         | 内容                                     |
| ------------ | ---------------------------------------- |
| **用例标题** | 成功创建新项目                           |
| **前置条件** | 已登录为 admin 或 super_admin 角色       |
| **操作步骤** | 1. 在项目列表页点击"新建项目"按钮        |
|              | 2. 在弹窗中输入项目名称"ECommerce"       |
|              | 3. 输入项目前缀"ECO"                     |
|              | 4. 输入描述"电商平台测试项目"            |
|              | 5. 点击"确定"                            |
| **预期结果** | 1. 请求 `POST /api/v1/projects` 返回 201 |
|              | 2. 项目列表刷新，新项目出现在列表中      |
|              | 3. 弹窗关闭，显示成功 Message            |

### TC-04 发起 AI 生成任务

| 项目         | 内容                                                          |
| ------------ | ------------------------------------------------------------- |
| **用例标题** | 成功发起用例生成任务并查看草稿                                |
| **前置条件** | 1. 已创建项目 ECommerce 及模块"用户中心"（缩写 USR）          |
|              | 2. 已上传 PRD 文档且状态为 completed                          |
|              | 3. 已登录                                                     |
| **操作步骤** | 1. 导航至项目 ECommerce → AI 生成 → 点击"新建任务"            |
|              | 2. 选择模块"用户中心"                                         |
|              | 3. 输入需求描述"测试用户注册功能，包括邮箱验证和密码强度校验" |
|              | 4. 设置用例数量为 5                                           |
|              | 5. 点击"立即生成"                                             |
|              | 6. 等待任务完成（列表页轮询状态）                             |
|              | 7. 点击任务进入详情，查看草稿列表                             |
| **预期结果** | 1. `POST /api/v1/generation/tasks` 返回 201                   |
|              | 2. 任务列表显示新任务，状态为 processing                      |
|              | 3. 任务完成后状态变为 completed                               |
|              | 4. 草稿列表展示 5 条草稿，每条包含标题、类型、优先级、置信度  |
|              | 5. 草稿箱 Badge 数字更新                                      |

### TC-05 确认草稿转正式用例

| 项目         | 内容                                                      |
| ------------ | --------------------------------------------------------- |
| **用例标题** | 确认草稿并验证自动编号                                    |
| **前置条件** | 1. 存在状态为 pending 的草稿                              |
|              | 2. 草稿属于项目 ECommerce（前缀 ECO）、模块 USR           |
|              | 3. 当前日期 2026-04-16                                    |
| **操作步骤** | 1. 进入草稿箱，点击目标草稿进入确认页                     |
|              | 2. 左侧编辑区查看/编辑用例内容                            |
|              | 3. 右侧查看引用来源                                       |
|              | 4. 点击"确认并转为正式用例"                               |
|              | 5. 选择目标模块"用户中心"                                 |
|              | 6. 确认操作                                               |
| **预期结果** | 1. `POST /api/v1/generation/drafts/{id}/confirm` 返回 201 |
|              | 2. 返回的 TestCase 包含编号 `ECO-USR-20260416-001`        |
|              | 3. 草稿状态变为 confirmed                                 |
|              | 4. 用例库中出现新的正式用例                               |
|              | 5. 草稿箱 Badge 数字减 1                                  |

### TC-06 批量确认草稿

| 项目         | 内容                                                       |
| ------------ | ---------------------------------------------------------- |
| **用例标题** | 批量确认多个草稿                                           |
| **前置条件** | 草稿箱中至少有 3 条 pending 状态草稿                       |
| **操作步骤** | 1. 进入草稿箱列表页                                        |
|              | 2. 勾选 3 条草稿                                           |
|              | 3. 点击"批量确认"按钮                                      |
|              | 4. 选择目标模块                                            |
|              | 5. 确认操作                                                |
| **预期结果** | 1. `POST /api/v1/generation/drafts/batch-confirm` 返回 200 |
|              | 2. 响应体中 `success_count` = 3，`failed_count` = 0        |
|              | 3. 列表中 3 条草稿状态变为 confirmed                       |
|              | 4. 显示成功 Message "成功确认 3 条草稿"                    |

### TC-07 创建测试计划并录入结果

| 项目         | 内容                                                    |
| ------------ | ------------------------------------------------------- |
| **用例标题** | 创建计划、关联用例、录入执行结果                        |
| **前置条件** | 1. 项目 ECommerce 下有至少 2 条正式测试用例             |
|              | 2. 已登录                                               |
| **操作步骤** | 1. 导航至 ECommerce → 测试计划 → 点击"新建计划"         |
|              | 2. 输入计划名称"Sprint 12 回归测试"                     |
|              | 3. 从用例库选择 2 条用例                                |
|              | 4. 点击"创建"                                           |
|              | 5. 进入计划详情页                                       |
|              | 6. 对第一条用例录入结果：状态 = pass，备注 = "功能正常" |
|              | 7. 对第二条用例录入结果：状态 = fail，备注 = "登录超时" |
| **预期结果** | 1. `POST /api/v1/plans` 返回 201，计划状态为 draft      |
|              | 2. `POST /api/v1/plans/{id}/results` 返回 201           |
|              | 3. 计划详情统计更新：total=2, passed=1, failed=1        |
|              | 4. 对应用例状态更新为 pass / fail                       |

### TC-08 项目仪表盘统计

| 项目         | 内容                                                                 |
| ------------ | -------------------------------------------------------------------- |
| **用例标题** | 验证仪表盘统计数据准确性                                             |
| **前置条件** | 项目下有完整的用例和执行记录                                         |
| **操作步骤** | 1. 导航至项目 ECommerce 仪表盘                                       |
|              | 2. 查看 4 个统计卡片                                                 |
|              | 3. 查看通过率趋势图                                                  |
|              | 4. 查看最近 AI 生成任务列表                                          |
| **预期结果** | 1. `GET /api/v1/projects/{id}/stats` 返回统计数据                    |
|              | 2. 4 个卡片数值与实际数据一致（用例总数、通过率、覆盖率、AI 生成数） |
|              | 3. 趋势折线图正确渲染，日期轴按天展示                                |
|              | 4. 最近任务列表按创建时间倒序排列                                    |

---

## 5. 输出格式示例

### 5.1 API 响应体示例 — 测试用例详情

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "moduleId": "f7e6d5c4-b3a2-1098-7654-321fedcba098",
  "userId": "12345678-90ab-cdef-1234-567890abcdef",
  "number": "ECO-USR-20260416-001",
  "title": "验证邮箱注册流程 - 有效邮箱格式",
  "preconditions": ["用户未注册过该邮箱", "邮件服务正常运行"],
  "steps": [
    "打开注册页面",
    "输入有效邮箱地址 user@example.com",
    "输入符合要求的密码 Test@1234",
    "点击注册按钮",
    "检查邮箱收到验证邮件"
  ],
  "expected": {
    "step_1": "注册页面正常加载，表单元素完整",
    "step_2": "邮箱格式校验通过，无错误提示",
    "step_3": "密码强度校验通过，无错误提示",
    "step_4": "显示注册成功提示，跳转至邮箱验证页",
    "step_5": "验证邮件在 60s 内送达，内容包含验证链接"
  },
  "caseType": "functionality",
  "priority": "P1",
  "status": "pass",
  "aiMetadata": {
    "generationTaskId": "b2c3d4e5-f6a7-8901-bcde-f23456789012",
    "confidence": "high",
    "referencedChunks": [
      {
        "chunkId": "c3d4e5f6-a7b8-9012-cdef-345678901234",
        "documentId": "d4e5f6a7-b8c9-0123-defa-456789012345",
        "documentTitle": "用户注册模块 PRD v2.0",
        "similarityScore": 0.92
      }
    ],
    "modelVersion": "deepseek-chat-v3",
    "generatedAt": "2026-04-16T08:30:00Z"
  },
  "createdAt": "2026-04-16T08:35:00Z",
  "updatedAt": "2026-04-16T10:20:00Z",
  "module_name": "用户中心",
  "project_name": "ECommerce",
  "project_prefix": "ECO",
  "created_by_name": "张三"
}
```

### 5.2 前端组件 Props 接口示例 — StatusTag

```typescript
// src/components/business/StatusTag.tsx
import { Tag } from '@arco-design/web-react'

/**
 * 状态枚举与色彩映射的统一 Props
 * 所有 StatusTag、表格 Tag、筛选 Tag 统一使用此接口
 */
interface StatusTagProps {
  /** 状态枚举值，如 'pass' | 'fail' | 'pending' */
  status: string
  /** 状态类别，决定色彩映射表 */
  category:
    | 'case_status'
    | 'plan_status'
    | 'task_status'
    | 'draft_status'
    | 'priority'
    | 'confidence'
    | 'case_type'
    | 'document_type'
    | 'document_status'
  /** 自定义中文标签（不传则使用默认映射） */
  label?: string
  /** 尺寸 */
  size?: 'small' | 'default'
}

// 色彩映射表示例（完整映射见 UX 设计规范 2.1 节）
const COLOR_MAP: Record<string, Record<string, { color: string; text: string; bg: string }>> = {
  case_status: {
    unexecuted: { color: '#86909C', text: '未执行', bg: 'rgba(134,144,156,0.10)' },
    pass:       { color: '#00B42A', text: '通过',   bg: 'rgba(0,180,42,0.10)' },
    block:      { color: '#FF7D00', text: '阻塞',   bg: 'rgba(255,125,0,0.10)' },
    fail:       { color: '#F53F3F', text: '失败',   bg: 'rgba(245,63,63,0.10)' },
  },
  priority: {
    P0: { color: '#F53F3F', text: 'P0 紧急', bg: 'rgba(245,63,63,0.10)' },
    P1: { color: '#FF7D00', text: 'P1 高',   bg: 'rgba(255,125,0,0.10)' },
    P2: { color: '#7B61FF', text: 'P2 中',   bg: 'rgba(123,97,255,0.10)' },
    P3: { color: '#86909C', text: 'P3 低',   bg: 'rgba(134,144,156,0.10)' },
  },
}

export function StatusTag({ status, category, label, size }: StatusTagProps) {
  const mapping = COLOR_MAP[category]?.[status]
  if (!mapping) return null

  return (
    <Tag
      size={size ?? 'small'}
      style={{
        color: mapping.color,
        backgroundColor: mapping.bg,
        borderColor: 'transparent',
      }}
    >
      {label ?? mapping.text}
    </Tag>
  )
}
```

### 5.3 API 分页响应体 JSON Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "TestCaseListResponse",
  "type": "object",
  "required": ["data", "total", "offset", "limit"],
  "properties": {
    "data": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/TestCase"
      },
      "description": "测试用例列表"
    },
    "total": {
      "type": "integer",
      "description": "符合筛选条件的总记录数",
      "example": 128
    },
    "offset": {
      "type": "integer",
      "description": "当前偏移量",
      "example": 0
    },
    "limit": {
      "type": "integer",
      "description": "每页数量",
      "example": 10
    }
  },
  "definitions": {
    "TestCase": {
      "type": "object",
      "properties": {
        "id": { "type": "string", "format": "uuid" },
        "moduleId": { "type": "string", "format": "uuid" },
        "number": {
          "type": "string",
          "pattern": "^[A-Z]{2,4}-[A-Z]{2,4}-\\d{8}-\\d{3}$"
        },
        "title": { "type": "string" },
        "caseType": { "$ref": "#/definitions/CaseType" },
        "priority": { "$ref": "#/definitions/Priority" },
        "status": { "$ref": "#/definitions/CaseStatus" },
        "createdAt": { "type": "string", "format": "date-time" },
        "updatedAt": { "type": "string", "format": "date-time" }
      }
    },
    "CaseType": {
      "type": "string",
      "enum": ["functionality", "performance", "api", "ui", "security"]
    },
    "Priority": { "type": "string", "enum": ["P0", "P1", "P2", "P3"] },
    "CaseStatus": {
      "type": "string",
      "enum": ["unexecuted", "pass", "block", "fail"]
    }
  }
}
```
