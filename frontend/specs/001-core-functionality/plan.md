# Aitestos 前端实现计划

## Context

基于 `specs/001-core-functionality/spec.md` v2.0，制定前端实现计划。该规格书覆盖认证、项目管理、知识库、AI 生成、草稿确认、用例管理、测试执行 7 大模块，包含 20+ 页面、30+ API 端点、12 种枚举类型。

---

## Phase 依赖关系

```
Phase 1 基础设施 ──→ Phase 2 通用组件 ──→ Phase 3 项目管理 ──→ Phase 4 模块管理
                                                                │
                        ┌───────────────────────────────────────┤
                        ▼                                       ▼
                  Phase 5 知识库                           Phase 7 用例管理
                        │                                       │
                        ▼                                       ▼
                  Phase 6 AI 生成 ──→ Phase 8 草稿箱 ──→ Phase 9 测试执行
                                                                │
                                                                ▼
                                                         Phase 10 配置管理
```

---

## Phase 1: 基础设施与认证 (P0)

> **验收用例**: TC-A01 登录成功 / TC-A02 登录失败 / TC-A03 Token 刷新 / TC-A04 注册

### 1.1 类型定义
- **文件**: `src/types/enums.ts` — 定义全部 12 个枚举（CaseStatus, CaseType, PlanStatus, Priority, ResultStatus, TaskStatus, DraftStatus, DocumentType, DocumentStatus, UserRole, Confidence, SceneType）
- **文件**: `src/types/api.ts` — 定义所有 API 请求/响应类型（Project, Module, Document, DocumentChunk, GenerationTask, CaseDraft, TestCase, TestPlan, TestResult, ProjectStatistics, User 等）

### 1.2 基础设施层
- **文件**: `src/lib/request.ts` — Axios 实例，Token 刷新（401 → refresh → 重放，并发排队），错误拦截
- **文件**: `src/lib/query-client.ts` — React Query 全局配置（staleTime, retry, Query Key 工厂）
- **文件**: `src/lib/utils.ts` — cn() 等工具函数

### 1.3 错误处理策略

统一 HTTP 错误处理，在 `src/lib/request.ts` 拦截器中实现：

| HTTP 状态码 | 处理方式 | UI 反馈 |
|---|---|---|
| 400 | 表单验证提示 | 字段下方红色错误 |
| 401 | 跳转登录页 | 清除 token + redirect |
| 403 | 无权限提示 | Notification.warning |
| 404 | 404 页面 | "返回首页"按钮 |
| 409 | 冲突提示 | 字段级错误 |
| 500 | 服务器错误 | Notification.error + 重试按钮 |

### 1.4 全局状态
- **文件**: `src/store/useAppStore.ts` — Zustand: sidebarCollapsed, notifications
- **文件**: `src/features/auth/hooks/useAuthStore.ts` — Zustand: tokens, user, login/logout

### 1.5 认证模块
- **路由**: `/login`, `/register`
- **页面**: LoginPage, RegisterPage（AuthLayout 左右分栏 55:45）
- **组件**: LoginForm, LoginBanner, AuthProvider
- **Hooks**: useAuth（login/register mutations）
- **Services**: auth.ts（POST /auth/login, /auth/register, /auth/refresh）
- **表单验证**: Zod schema（邮箱格式、密码 ≥8、用户名 3-32、确认密码一致）
- **关键逻辑**: Token 刷新（401 自动重放）、来源路由保留

---

## Phase 2: 通用组件 (P0)

> 被后续所有 Phase 依赖，必须在业务模块之前完成。

### 2.1 StatusTag
- **支持 type**: caseStatus | planStatus | taskStatus | draftStatus | priority | confidence | caseType | documentType | documentStatus
- **色彩映射**: 根据 ux-design-spec.md §2.1 语义色表
- **尺寸**: default(24px) / small(20px)

### 2.2 SearchTable
- 基于 Arco Table，行高 48px/36px(紧凑)，斑马纹，默认 20 条/页

### 2.3 ArrayEditor
- 序号 + 输入框 + 上移/下移/删除，底部虚线"添加"
- 前置条件 min 0，测试步骤 min 1

### 2.4 StatsCard
- Arco Card + Statistic，左侧 4px 装饰线，数值 28px Bold tabular-nums，计数动画 800ms

### 2.5 SplitPanel
- 基于 Arco ResizeBox.Split，左默认 240px，拖拽条 2px→4px(#7B61FF) 悬浮效果

### 2.6 其他
- ReferencePanel — 草稿引用来源展示
- useFormError — 表单错误处理 hook

---

## Phase 3: 项目管理 (P0)

> **验收用例**: TC-P01 创建项目 / TC-P02 前缀重复 / TC-P03 仪表盘统计

### 3.1 布局与路由
- **文件**: `src/components/layout/AppLayout.tsx` — 侧边栏 + 内容区
- **文件**: `src/components/layout/Sidebar.tsx` — 项目导航菜单（仪表盘/用例/计划/知识库/AI生成/草稿箱）
- **文件**: `src/components/layout/Header.tsx` — 项目选择器 + 用户菜单 + 通知
- **文件**: `src/router/RouteGuard.tsx` — 认证守卫 + 权限检查（admin+ 路由 requireAdmin 属性）

### 3.2 项目列表页
- **路由**: `/projects`
- **页面**: ProjectListPage — 搜索 + Card 网格（3列）
- **组件**: 项目卡片（颜色装饰条 + 前缀 Tag + 统计行 + 操作按钮）
- **操作**: 卡片"进入"按钮跳转仪表盘，"设置"按钮进入项目设置
- **空状态**: FolderOpen 图标 + "创建第一个项目"

### 3.3 创建项目模态框
- **组件**: CreateProjectModal（600px）
- **字段**: 名称(2-255)、前缀(2-4大写字母, `^[A-Z]+$`)、描述
- **验证**: 前缀实时格式校验 + 失焦唯一性校验，409 字段级错误
- **API**: `POST /projects`

### 3.4 项目编辑/删除
- **编辑**: 项目卡片"设置"或仪表盘操作按钮 → Modal（预填名称/前缀/描述）
- **删除**: 仪表盘操作 → Popconfirm 确认 → `DELETE /projects/{id}` → 跳转项目列表
- **API**: `PUT /projects/{id}`, `DELETE /projects/{id}`

### 3.5 项目仪表盘
- **路由**: `/projects/:projectId`
- **页面**: ProjectDashboard — 标题行 + 4 统计卡片 + 趋势图(60%) + 最近任务(40%)
- **统计卡片**: 用例总数、通过率、覆盖率、AI 生成数（左侧 4px 装饰线，计数动画 800ms）
- **新项目引导**: ❶上传文档 → ❷AI生成 → ❸创建计划
- **操作按钮**: 上传文档 / 发起生成 / 新建用例
- **API**: `GET /projects/{id}/stats`

---

## Phase 4: 模块管理 (P0)

### 4.1 模块管理页
- **路由**: `/projects/:projectId/settings/modules`（权限 admin+）
- **布局**: SplitPanel — 左 280px 模块树(Arco Tree) + 右编辑区
- **操作**:
  - 新增：左侧底部"新增模块"按钮 → 弹出表单
  - 编辑：点击模块节点 → 右侧编辑表单（名称、缩写、描述）
  - 删除：悬浮删除图标 → Modal 输入模块名称确认（提示级联删除用例）
- **API**: `GET/POST /projects/{id}/modules`, `DELETE /modules/{id}`
- **注意**: 缺少 `PUT /modules/{id}`，需后端补充（已列入缺失 API 清单）

---

## Phase 5: 知识库 (P0)

> **验收用例**: TC-K01 上传文档 / TC-K02 文档详情与分块

### 5.1 文档列表页
- **路由**: `/projects/:projectId/knowledge`
- **布局**: 筛选栏（类型/状态/搜索）+ 表格
- **表格**: 文档名称(前缀类型图标) | 类型 | 状态(StatusTag, processing 辉光动画) | 分块数 | 上传时间 | 操作
- **空状态**: FileText 图标 + "暂无文档" + "上传文档"按钮
- **API**: `GET /knowledge/documents`

### 5.2 文档上传模态框
- **组件**: UploadDocumentModal（600px）
- **字段**: 名称、类型(PRD/Figma/API Spec/Swagger/Markdown)、文件拖拽上传
- **文件限制**: PRD=.docx/.pdf/.md, API=.json/.yaml, ≤50MB
- **API**: `POST /knowledge/documents`

### 5.3 文档详情页
- **路由**: `/projects/:projectId/knowledge/:docId`
- **布局**: SplitPanel — 左 300px 信息面板 + 右分块列表
- **左侧**: 文档名称(h2) + 类型 Tag + 状态 Tag + 上传人/时间 + 分块数 + Steps 状态流
- **右侧分块列表**: 序号 + 内容预览(3行截断) + 展开按钮 + 引用计数(链接图标)
- **API**: `GET /knowledge/documents/{id}`, `GET /knowledge/documents/{id}/chunks`

### 5.4 Figma 集成页
- **路由**: `/projects/:projectId/knowledge/figma`
- **布局**: 全页面表单，分区域
- **区域 1 连接配置**: 认证方式(Radio: 个人令牌/OAuth 2.0) + 令牌(Input.Password) + 测试连接按钮
- **区域 2 导入文件**: Figma URL(Input) + 解析(Button)
- **区域 3 节点选择**: Arco Tree(带 Checkbox) + 底部取消/"确认导入"
- **对应 US**: US-2.2 关联 Figma 设计稿

---

## Phase 6: AI 生成模块 (P0 核心)

> 核心差异化功能。AI 视觉特征：品牌紫色辉光、Sparkles 图标、AI Gradient 按钮、glow-pulse 动画。
>
> **验收用例**: TC-G01 生成任务(充足) / TC-G02 生成任务(为空) / TC-G03 查看草稿

### 6.1 生成任务列表页
- **路由**: `/projects/:projectId/generation`
- **特征**: "新建生成任务"按钮 AI Gradient + Sparkles 图标，Processing 行辉光脉冲 + 紫色背景
- **操作**: 双击行进入详情，failed 状态显示"重试"按钮(Popconfirm 确认)
- **API**: `GET /generation/tasks`（**缺失，需后端补充**）

### 6.2 新建生成任务页
- **路由**: `/projects/:projectId/generation/new`
- **布局**: 居中表单(max-width 720px)
- **知识库就绪度指示器**（页面顶部）:
  - 充足(绿): "N 份文档 · M 个分块 · 就绪"
  - 不足(黄): "N 份文档 · M 个分块 · 内容有限" + 弹窗警告，用户确认后可继续
  - 为空(红): "请先上传需求文档" + 禁用提交按钮
- **RAG 降级策略**:
  - 知识库为空 → 禁止生成，禁用按钮
  - 检索结果 < 1 个文档块 → 弹窗提示"知识库内容不足，生成质量可能较低"，用户确认后继续，置信度强制为"低"
- **字段**: 目标模块(Select)、需求描述(TextArea, ≥10字)、用例数量(1-20, 默认5)
- **高级选项**(Collapse 默认折叠): 场景类型、优先级偏好、用例类型、生成模式
- **API**: `POST /generation/tasks`

### 6.3 任务详情页
- **路由**: `/projects/:projectId/generation/:taskId`
- **状态卡片**: pending(Spin + 边框微弱 Glow) / processing(Progress+shimmer+轮询5s+glow-pulse) / completed(spring 庆祝动画) / failed(Alert+重试按钮)
- **草稿列表**: ai-reveal 揭示动画(交错80ms)，工具栏：全选+批量确认(AI Gradient)+批量拒绝
- **API**: `GET /generation/tasks/{id}`(轮询), `GET /generation/tasks/{id}/drafts`

---

## Phase 7: 用例管理 (P0)

> **验收用例**: TC-C01 创建用例 / TC-C02 AI用例详情 / TC-C03 编号验证

### 7.1 用例库列表页
- **路由**: `/projects/:projectId/cases`
- **布局**: SplitPanel — 左 240px 模块树("全部"+各模块名称+用例数 Badge) + 右表格
- **工具栏**: 搜索 + 状态/类型/优先级筛选(Select) + 新建(primary) + 导入 + 导出(default)
- **表格**: Checkbox(48px) | 编号(200px monospace) | 标题(弹性) | 类型 Tag(100px) | 优先级 Tag(80px) | 状态 Tag(80px) | 创建人(80px) | 更新时间(160px) | 操作(100px)
- **批量操作栏**(选中行后浮现): 修改优先级 / 修改状态 / 加入计划 / 删除
- **导入**: Modal → Upload(drag, .xlsx/.csv) + 模板下载 → `POST /testcases/import`
- **导出**: Dropdown → 全部/已选中/筛选结果 → `POST /testcases/export` → .xlsx 下载
- **空状态**: FileCheck 图标 + "手动创建"(primary) + "使用 AI 生成"(AI Gradient)
- **API**: `GET /testcases`

### 7.2 创建用例抽屉
- **组件**: CreateCaseDrawer（640px 右侧 Drawer）
- **字段**: 目标模块(Select) | 标题(2-500) | 前置条件(ArrayEditor, min 0) | 测试步骤(ArrayEditor, min 1) | 预期结果(TextArea) | 用例类型(Select) | 优先级(Select, P0-P3)
- **API**: `POST /testcases`

### 7.3 用例编辑
- **触发**: 用例详情页点击"编辑" → Drawer（复用 CreateCaseDrawer 组件，预填数据）
- **API**: `PUT /testcases/{id}`

### 7.4 用例复制
- **触发**: 用例详情页点击"复制" → 打开新建 Drawer，预填源用例所有字段（标题前缀"[副本]"）
- **实现**: 复用 CreateCaseDrawer，传入初始值
- **API**: `POST /testcases`

### 7.5 用例详情页
- **路由**: `/projects/:projectId/cases/:caseId`
- **页头**: 返回链接 + 编号(monospace h1) + 标题 + 状态 Tag(大号) + 编辑/复制/删除按钮
- **基本信息卡片**: 3列网格（类型、优先级、模块、创建人、创建/更新时间）
- **AI 元数据区**(Collapse，仅 AI 生成用例显示):
  - 标题栏：`rgba(123,97,255,0.06)` 紫色背景 + Sparkles 图标
  - 内容：生成任务 ID(可点击链接) + 置信度 StatusTag + 引用文档块(相似度着色) + 模型版本 + 生成时间
  - 源文档已变更 → ⚠️ Alert(warning)
  - 源文档已删除 → Alert(error) + "移除引用"按钮
- **执行历史**: 表格 — 计划名称(可点击跳转) | 执行结果 StatusTag | 执行人 | 时间 | 备注
- **API**: `GET /testcases/{id}`, `PUT /testcases/{id}`, `DELETE /testcases/{id}`

---

## Phase 8: 草稿箱 (P0)

> **验收用例**: TC-D01 单条确认 / TC-D02 批量确认 / TC-D03 拒绝草稿

### 8.1 草稿列表页
- **路由**: `/drafts`（全局视图，跨项目）
- **筛选**: 项目(Select, 随模块联动) + 模块(Select) + 状态(Select: 待处理/已确认/已拒绝) + 搜索(Input)
- **侧边栏 Badge**: pending 数量，品牌紫色，每 30s 刷新
- **API**: `GET /drafts`（**缺失，需后端补充**）

### 8.2 草稿确认页 ⭐ 核心页面
- **路由**: `/drafts/:draftId`
- **导航栏**: 返回链接 + "第N/M条"进度 + 圆点导航(当前紫色高亮) + 键盘←/→切换(未聚焦输入框时)
- **布局**: 左 60% 编辑区（精准工作区风格） + 右 40% 引用来源（AI 辉光风格，紫色底色 `rgba(123,97,255,0.03)`）
- **左侧编辑字段**: 标题(Input) | 前置条件(ArrayEditor) | 测试步骤(ArrayEditor) | 预期结果(TextArea) | 类型(Select) | 优先级(Select)
- **右侧引用来源**: 文档标题(链接) + 类型 Tag + 相似度(>0.8绿/0.5-0.8黄/<0.5红) + 引用原文(灰色背景,5行截断) + "查看原文"链接。无引用时："此草稿未引用知识库内容"
- **底部操作栏**(固定底部):

  | 按钮 | 类型 | 行为 |
  |---|---|---|
  | 拒绝 | danger | Modal: 原因(Radio: 重复/无关/低质量/其他) + 反馈(TextArea) → `POST /generation/drafts/{id}/reject` |
  | 保存修改 | default | 前端 React state 暂存编辑内容，不调用 API，切换草稿时自动恢复 |
  | 确认并转为正式用例 | primary(AI Gradient) | Modal: 选择目标模块(Select) → `POST /generation/drafts/{id}/confirm` → Message.success("用例 {number} 已创建") → 跳转用例详情页 |

- **切换保护**: 切换草稿/关闭页面前，如有未保存编辑，弹出确认 Dialog
- **切换前自动保存**: 当前编辑内容暂存于 React state

---

## Phase 9: 测试执行 (P1)

> **验收用例**: TC-PL01 创建计划 / TC-PL02 录入结果 / TC-PL03 状态流转

### 9.1 计划列表页
- **路由**: `/projects/:projectId/plans`
- **表格**: 计划名称(可点击) | 状态 StatusTag | 用例数 | 通过率(百分比+迷你进度条) | 创建人 | 创建时间 | 操作
- **API**: `GET /plans`

### 9.2 新建计划页
- **路由**: `/projects/:projectId/plans/new`
- **布局**: SplitPanel — 左表单(名称 Input 必填 2-255字符 + 描述 TextArea 选填) + 右用例选择面板
- **右侧用例选择**(Tab 切换):
  - "可选用例" Tab：模块/类型/优先级筛选 + Checkbox 表格
  - "已选用例" Tab：已选列表 + 移除按钮 + Badge 显示数量
- **API**: `POST /plans`

### 9.3 计划详情页
- **路由**: `/projects/:projectId/plans/:planId`
- **操作按钮**(按状态变化):

  | 当前状态 | 可用操作 |
  |---|---|
  | draft | 编辑 + 开始执行(primary) + 删除(danger) |
  | active | 编辑 + 标记完成(primary) |
  | completed | 重新执行 + 归档(primary) |
  | archived | 取消归档 |

- **统计卡片**(5列): 总用例(灰线) + 通过(绿线) + 失败(红线) + 阻塞(橙线) + 跳过(灰线)
- **执行进度条**: Arco Progress，品牌紫色
- **用例执行表格**(紧凑变体 36px 行高):
  - Checkbox | 编号 | 标题 | 类型+优先级 | 执行结果(StatusTag/内联Select) | 执行人 | 时间 | "录入"按钮
- **快捷录入**: 点击执行结果列 → 内联 Select → 选择值自动提交 → 整行闪烁对应色 → Toast "已录入：通过" + "撤销"链接(3s 内)
- **详细录入**: 点击"录入" → Modal(500px) — 用例信息(只读) + 执行结果(Radio) + 备注(TextArea)
- **批量录入**: 选中多条 → "批量录入" → Modal — 结果选择(Radio) + 备注 → 批量调用
- **API**: `GET /plans/{id}`, `POST /plans/{id}/results`, `PATCH /plans/{id}/status`

---

## Phase 10: 配置管理 (P1)

### 10.1 配置管理页
- **路由**: `/projects/:projectId/settings/configs`（权限 admin+）
- **布局**: 工具栏 + 表格
- **操作**:
  - 新增/编辑: Modal(500px) — 键名 + 值(JSON, JsonEditor) + 描述
  - 删除: Popconfirm 确认
  - 导入 JSON: Modal — 粘贴 JSON + 预览 + 确认
  - 导出 JSON: 直接下载
- **API**: `GET/PUT /projects/{id}/configs/{key}`, `POST /projects/{id}/configs/import`, `GET /projects/{id}/configs/export`

---

## 前端路由完整清单

| 路由 | 页面组件 | Phase | 权限 |
|---|---|---|---|
| `/login` | LoginPage | 1 | 公开 |
| `/register` | RegisterPage | 1 | 公开 |
| `/projects` | ProjectListPage | 3 | 登录 |
| `/projects/:projectId` | ProjectDashboard | 3 | 登录 |
| `/projects/:projectId/settings/modules` | ModuleManagePage | 4 | admin+ |
| `/projects/:projectId/configs` | ConfigManagePage | 10 | admin+ |
| `/projects/:projectId/knowledge` | KnowledgeListPage | 5 | 登录 |
| `/projects/:projectId/knowledge/:docId` | DocumentDetailPage | 5 | 登录 |
| `/projects/:projectId/knowledge/figma` | FigmaIntegrationPage | 5 | admin+ |
| `/projects/:projectId/generation` | GenerationTaskListPage | 6 | 登录 |
| `/projects/:projectId/generation/new` | NewGenerationTaskPage | 6 | 登录 |
| `/projects/:projectId/generation/:taskId` | TaskDetailPage | 6 | 登录 |
| `/drafts` | DraftListPage | 8 | 登录 |
| `/drafts/:draftId` | DraftConfirmPage | 8 | 登录 |
| `/projects/:projectId/cases` | CaseListPage | 7 | 登录 |
| `/projects/:projectId/cases/:caseId` | CaseDetailPage | 7 | 登录 |
| `/projects/:projectId/plans` | PlanListPage | 9 | 登录 |
| `/projects/:projectId/plans/new` | NewPlanPage | 9 | 登录 |
| `/projects/:projectId/plans/:planId` | PlanDetailPage | 9 | 登录 |

---

## 缺失 API 清单（需后端补充）

| 功能 | 端点 | 方法 | 说明 |
|---|---|---|---|
| 生成任务列表 | `GET /generation/tasks` | GET | project_id, status, offset, limit |
| 全局草稿列表 | `GET /drafts` | GET | project_id, module_id, status, keywords, offset, limit |
| 模块编辑 | `PUT /projects/{id}/modules/{moduleId}` | PUT | name, abbreviation, description |
| 计划状态变更 | `PATCH /plans/{id}/status` | PATCH | status |
| 用例导入 | `POST /testcases/import` | POST | multipart/form-data |
| 用例导出 | `POST /testcases/export` | POST | project_id, filters, format |

---

## 设计令牌速查

| 类别 | 值 |
|---|---|
| 品牌色 | Primary `#7B61FF` |
| AI Glow | `rgba(123,97,255,0.15)` |
| AI Gradient | `linear-gradient(135deg, #7B61FF, #9B7BFF)` |
| 标题字体 | DM Sans |
| 编号字体 | JetBrains Mono (monospace) |
| 间距基础 | 4px，常用 8/12/16/20/24px |

---

## 验证方式

1. `make check` — lint + format + type-check 全通过
2. `make test` — 所有单元测试通过（每个 Hook/Service 必须有测试）
3. `make build` — 生产构建无错误
4. 逐页面验收：对照 spec §7 验收测试用例

| TC 编号 | 验收项 | 对应 Phase |
|---|---|---|
| TC-A01 | 用户登录成功 | Phase 1 |
| TC-A02 | 登录失败 | Phase 1 |
| TC-A03 | Token 自动刷新 | Phase 1 |
| TC-A04 | 用户注册 | Phase 1 |
| TC-P01 | 创建项目 | Phase 3 |
| TC-P02 | 项目前缀重复 | Phase 3 |
| TC-P03 | 项目仪表盘统计 | Phase 3 |
| TC-K01 | 上传文档 | Phase 5 |
| TC-K02 | 文档详情与分块 | Phase 5 |
| TC-G01 | 创建生成任务(充足) | Phase 6 |
| TC-G02 | 创建生成任务(为空) | Phase 6 |
| TC-G03 | 查看草稿 | Phase 6 |
| TC-D01 | 单条确认草稿 | Phase 8 |
| TC-D02 | 批量确认草稿 | Phase 8 |
| TC-D03 | 拒绝草稿 | Phase 8 |
| TC-C01 | 创建用例(手动) | Phase 7 |
| TC-C02 | AI 用例详情 | Phase 7 |
| TC-C03 | 编号验证 | Phase 7 |
| TC-PL01 | 创建测试计划 | Phase 9 |
| TC-PL02 | 录入执行结果 | Phase 9 |
| TC-PL03 | 计划状态流转 | Phase 9 |
