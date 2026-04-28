---
title: 核心功能原子化任务列表
version: 2.1.0
author: 技术组长
status: draft
created: 2026-04-28
updated: 2026-04-28
based_on: [spec.md v2.0, plan.md v2.0]
review_notes: 修订 [P] 标记矛盾、补充遗漏任务、拆分过粗任务、对齐验收用例
---

# 核心功能原子化任务列表

> **标记说明**：
>
> - `[P]` = 可并行执行（无前置依赖）
> - `依赖: Txx` = 必须在 Txx 完成后执行
> - **测试任务始终优先于对应实现任务**
> - 每个 Task 对应 plan.md 中的一个具体交付物
> - **验收用例映射**：每个 Phase 标注对应 spec §7 的 TC 编号

---

## Phase 0: MSW 基础设施

### T01 `[P]` 配置 MSW 集中管理

- **文件**: `tests/msw/server.ts`（新建）, `tests/msw/handlers/`（新建目录）
- **操作**: 使用 `setupServer()` 创建 MSW 实例；创建 `tests/msw/handlers/index.ts` 汇总所有 handler；在 `tests/setup.ts` 中启动/重置/关闭 server
- **AC**: MSW 可拦截 fetch 请求；各 feature service 测试可 import 集中管理的 handler

---

## Phase 1: 基础设施与认证 (plan §1)

> **验收用例**: TC-A01 登录成功 / TC-A02 登录失败 / TC-A03 Token 刷新 / TC-A04 注册

### T02 `[P]` 类型定义 — 枚举

- **对应 plan**: §1.1
- **文件**: `src/types/enums.ts`
- **操作**: 定义全部 12 个字面量联合类型：CaseStatus, CaseType, PlanStatus, Priority, ResultStatus, TaskStatus, DraftStatus, DocumentType, DocumentStatus, UserRole, Confidence, SceneType
- **AC**: TypeScript strict 编译通过，无 `any`

### T03 依赖: T02 | 类型定义 — API

- **对应 plan**: §1.1
- **文件**: `src/types/api.ts`
- **操作**: 定义所有 API 请求/响应类型：PaginatedResponse\<T\>, UserJSON, Project, ProjectDetail, ProjectStatistics, Module, Document, DocumentDetail, DocumentChunk, GenerationTask, GenerationTaskCreate, CaseDraft, DraftConfirmRequest, DraftRejectRequest, TestCase, CaseDetail, TestCaseCreate, TestPlan, PlanDetail, PlanCreate, TestResult, ResultRecordRequest, ProjectConfig 等
- **AC**: 与 openapi.yaml schema 对齐，TypeScript 编译通过

### T04 `[P]` Axios 实例与 Token 刷新

- **对应 plan**: §1.2 + §1.3
- **文件**: `src/lib/request.ts`
- **操作**: 封装类型安全的 `get<T>`, `post<T,R>`, `put<T,R>`, `delete<T>` 方法；实现 401 → refresh → 重放逻辑（isRefreshing + pendingQueue + WeakSet 防循环）；注册 onAuthExpired 回调；实现统一错误拦截（400/401/403/404/409/500 → 对应 UI 反馈）
- **AC**: `grep -r "any" src/lib/request.ts` 无结果；Token 刷新并发安全

### T05 依赖: T04 | 测试: Axios 实例

- **文件**: `src/lib/request.test.ts`
- **操作**: 测试 401 触发 refresh 并重放原始请求；测试并发请求仅触发一次 refresh；测试 refresh 失败清除 token 并跳转 /login；测试 409 返回字段级错误
- **AC**: 4 个测试通过

### T06 `[P]` React Query Client 配置

- **对应 plan**: §1.2
- **文件**: `src/lib/query-client.ts`
- **操作**: 创建 QueryClient 实例（staleTime 5min, retry 1, refetchOnWindowFocus false）；导出 queryClient 和 Query Key 工厂函数
- **AC**: 导出 queryClient 实例

### T07 `[P]` 工具函数

- **对应 plan**: §1.2
- **文件**: `src/lib/utils.ts`
- **操作**: 保留 `cn()` (clsx + tailwind-merge)，添加 `formatDate()`, `truncateText()` 等工具
- **AC**: 工具函数可正常调用

### T08 `[P]` Zustand: useAppStore

- **对应 plan**: §1.4
- **文件**: `src/store/useAppStore.ts`
- **操作**: 实现 sidebarCollapsed + toggleSidebar + notifications
- **AC**: store 创建无报错

### T09 依赖: T08 | 测试: useAppStore

- **文件**: `src/store/useAppStore.test.ts`
- **操作**: 测试 sidebarCollapsed 默认值和切换；测试 notifications 增删
- **AC**: 2 个测试通过

### T10 `[P]` 实现: ErrorBoundary

- **文件**: `src/components/ErrorBoundary.tsx`
- **操作**: React ErrorBoundary 组件，捕获渲染错误显示 fallback UI（用于路由包裹）
- **AC**: 组件渲染无报错，子组件 throw 时显示 fallback

### T11 依赖: T04, T01 | 测试: Auth API service

- **对应 plan**: §1.5
- **文件**: `src/features/auth/services/auth.test.ts`
- **操作**: MSW mock 测试 `POST /auth/login` 返回 token；测试 `POST /auth/register` 返回 UserJSON；测试 `POST /auth/refresh` 返回新 token
- **AC**: 3 个测试通过

### T12 依赖: T11 | 实现: Auth API service

- **文件**: `src/features/auth/services/auth.ts`
- **操作**: 实现 authApi.login / register / refresh
- **AC**: T11 测试全部通过

### T13 依赖: T08 | 测试: Auth Zustand store

- **对应 plan**: §1.4
- **文件**: `src/features/auth/hooks/useAuthStore.test.ts`
- **操作**: 测试 login() 写入 localStorage + 更新 state；测试 logout() 清理 token + state；测试 refresh() 成功更新 token；测试 refresh() 失败调用 logout
- **AC**: 4 个测试通过

### T14 依赖: T13 | 实现: Auth Zustand store

- **文件**: `src/features/auth/hooks/useAuthStore.ts`
- **操作**: 实现 AuthState（user, token, refreshToken, isAuthenticated, login, logout, refresh, setUser）
- **AC**: T13 测试全部通过

### T15 依赖: T12 | 测试: useAuth hooks

- **对应 plan**: §1.5
- **文件**: `src/features/auth/hooks/useAuth.test.ts`
- **操作**: 测试 useLogin mutation 成功调用 store.login + invalidateQueries；测试 useRegister mutation 成功后跳转；测试失败时抛出错误
- **AC**: 3 个测试通过

### T16 依赖: T14, T15 | 实现: useAuth hooks

- **文件**: `src/features/auth/hooks/useAuth.ts`
- **操作**: 实现 useLogin(useMutation + store.login + navigate), useRegister, useRefresh
- **AC**: T15 测试全部通过

### T17 依赖: T16 | 测试: LoginForm 组件

- **对应 plan**: §1.5
- **文件**: `src/features/auth/components/LoginForm.test.tsx`
- **操作**: 测试渲染邮箱/密码输入框和登录按钮；测试无效邮箱显示验证错误；测试密码为空校验失败；测试登录成功跳转 /projects；测试登录失败显示错误提示
- **AC**: 5 个测试通过

### T18 依赖: T17 | 实现: LoginForm 组件

- **文件**: `src/features/auth/components/LoginForm.tsx`
- **操作**: Arco Form + React Hook Form + zod（email 格式, password 非空）；提交调用 useLogin mutation；loading 态禁用按钮
- **AC**: T17 测试全部通过

### T19 依赖: T16 | 测试: RegisterPage

- **对应 plan**: §1.5
- **文件**: `src/features/auth/components/RegisterPage.test.tsx`
- **操作**: 测试渲染用户名/邮箱/密码/确认密码/角色表单；测试字段校验（用户名 3-32, 邮箱格式, 密码 ≥8, 确认密码一致）；测试角色选项仅有 admin/normal（无 super_admin）；测试 409 字段级错误；测试注册成功跳转 /login
- **AC**: 5 个测试通过

### T20 依赖: T19 | 实现: RegisterPage

- **文件**: `src/features/auth/components/RegisterPage.tsx`
- **操作**: AuthLayout + Arco Form + React Hook Form + zod；角色 Select(admin/normal)；409 错误 → 字段级提示
- **AC**: T19 测试全部通过

### T21 `[P]` 依赖: T02 | Zod Schema: 认证表单

- **文件**: `src/features/auth/schema/loginSchema.ts`, `src/features/auth/schema/registerSchema.ts`（新建）
- **操作**: 独立导出 loginSchema(email格式+password非空) 和 registerSchema(username 3-32+email+password≥8+confirmPassword+role)，供 LoginForm/RegisterPage 引用
- **AC**: schema 可独立 import，单元测试通过

### T22 依赖: T14 | 测试: RouteGuard

- **对应 plan**: §1.5
- **文件**: `src/router/RouteGuard.test.tsx`
- **操作**: 测试未认证跳转 /login（保留来源路由）；测试正常 token 渲染 children；测试 requireAdmin=true 且 normal 角色跳转 /projects
- **AC**: 3 个测试通过

### T23 依赖: T22 | 实现: RouteGuard

- **文件**: `src/router/RouteGuard.tsx`
- **操作**: 读取 useAuthStore 判断 isAuthenticated + token exp；未认证 Navigate to /login（state.from）；requireAdmin 检查 role
- **AC**: T22 测试全部通过

### T24 `[P]` 实现: AuthLayout

- **对应 plan**: §1.5 (spec §4.1.1 AuthLayout 左右分栏 55:45)
- **文件**: `src/features/auth/components/AuthLayout.tsx`（新建）
- **操作**: 左右分栏布局(55%:45%)，左侧 Mesh Gradient 背景(#7B61FF→#5A3DC0→#3B1FA0) + LoginBanner，右侧白色表单区 + `<Outlet />`
- **AC**: 布局渲染正确，左右比例正确

### T25 `[P]` 实现: Providers

- **对应 plan**: §1.2
- **文件**: `src/app/providers.tsx`
- **操作**: QueryClientProvider + Arco ConfigProvider（品牌色 #7B61FF）嵌套
- **AC**: 组件渲染无报错

### T26 依赖: T25, T23, T24, T10 | 实现: App 入口 + 路由骨架

- **文件**: `src/app/App.tsx`, `src/router/index.tsx`
- **操作**: App 使用 Providers 包裹 RouterOutlet；router/index.tsx 定义骨架路由（/login 用 AuthLayout 包裹 + ErrorBoundary, /register, /projects 占位）
- **AC**: `npm run dev` 启动无报错，/login 可访问

---

## Phase 2: 通用组件 (plan §2)

> 被后续所有 Phase 依赖，必须在业务模块之前完成。

### T27 `[P]` 测试: StatusTag

- **对应 plan**: §2.1
- **文件**: `src/components/business/StatusTag.test.tsx`
- **操作**: 测试 9 种 type 的色彩映射正确；测试 default/small 尺寸；测试无效 value 返回 null
- **AC**: 3 个测试通过

### T28 依赖: T27 | 实现: StatusTag

- **文件**: `src/components/business/StatusTag.tsx`
- **操作**: 完整 COLOR_MAP 映射（9 category），Arco Tag + Tailwind，支持 size prop
- **AC**: T27 测试全部通过

### T29 `[P]` 测试: SearchTable

- **对应 plan**: §2.2
- **文件**: `src/components/business/SearchTable.test.tsx`
- **操作**: 测试渲染表格列和数据行；测试 loading 显示 Spin；测试分页控件；测试 empty 状态
- **AC**: 4 个测试通过

### T30 依赖: T29 | 实现: SearchTable

- **文件**: `src/components/business/SearchTable.tsx`
- **操作**: 封装 Arco Table + Pagination，统一 loading/error/empty 状态，Props 泛型化，行高 48px/36px(紧凑)
- **AC**: T29 测试全部通过

### T31 `[P]` 测试: ArrayEditor

- **对应 plan**: §2.3
- **文件**: `src/components/business/ArrayEditor.test.tsx`
- **操作**: 测试初始渲染；测试添加行；测试删除行；测试 min 约束（steps min 1 不可全删）；测试 onChange 回调
- **AC**: 5 个测试通过

### T32 依赖: T31 | 实现: ArrayEditor

- **文件**: `src/components/business/ArrayEditor.tsx`
- **操作**: 动态行列表（序号 + Input + 上移/下移/删除），底部虚线"添加"，min/max 约束
- **AC**: T31 测试全部通过

### T33 `[P]` 测试: StatsCard

- **对应 plan**: §2.4
- **文件**: `src/components/business/StatsCard.test.tsx`
- **操作**: 测试渲染标题和数值；测试趋势箭头（up/down/flat）；测试装饰线颜色
- **AC**: 3 个测试通过

### T34 依赖: T33 | 实现: StatsCard

- **文件**: `src/components/business/StatsCard.tsx`
- **操作**: Arco Card + Statistic，左侧 4px 装饰线，数值 28px Bold tabular-nums，支持 title/value/trend/color props
- **AC**: T33 测试全部通过

### T35 `[P]` 测试: SplitPanel

- **对应 plan**: §2.5
- **文件**: `src/components/business/SplitPanel.test.tsx`
- **操作**: 测试左右面板渲染 children；测试默认分割位置
- **AC**: 2 个测试通过

### T36 依赖: T35 | 实现: SplitPanel

- **文件**: `src/components/business/SplitPanel.tsx`
- **操作**: Arco ResizeBox.Split 封装，左默认 240px(min 180px)，拖拽条 2px→4px(#7B61FF) 悬浮
- **AC**: T35 测试全部通过

### T37 `[P]` 测试: ReferencePanel

- **对应 plan**: §2.6
- **文件**: `src/components/business/ReferencePanel.test.tsx`
- **操作**: 测试渲染引用块列表（文档标题、相似度、引用原文）；测试空列表显示"无引用来源"
- **AC**: 2 个测试通过

### T38 依赖: T37 | 实现: ReferencePanel

- **文件**: `src/components/business/ReferencePanel.tsx`
- **操作**: 接收 ReferencedChunk[] props，展示文档标题 + 类型 Tag + 相似度(>0.8绿/0.5-0.8黄/<0.5红) + 引用原文(5行截断) + "查看原文"链接
- **AC**: T37 测试全部通过

### T39 `[P]` 依赖: T02 | 实现: useFormError hook

- **对应 plan**: §2.6
- **文件**: `src/components/business/useFormError.ts`
- **操作**: 封装表单错误处理（409 字段级错误映射到 React Hook Form setError）
- **AC**: hook 可正常调用

---

## Phase 3: 项目管理 (plan §3)

> **验收用例**: TC-P01 创建项目 / TC-P02 前缀重复 / TC-P03 仪表盘统计

### T40 依赖: T04, T03, T01 | 测试: Projects API service

- **对应 plan**: §3.2 ~ §3.5
- **文件**: `src/features/projects/services/projects.test.ts`
- **操作**: MSW mock 测试 list() GET /projects；测试 get(id)；测试 create() POST /projects；测试 update(id) PUT /projects/{id}；测试 remove(id) DELETE /projects/{id}；测试 getStats(id) GET /projects/{id}/stats
- **AC**: 6 个测试通过

### T41 依赖: T40 | 实现: Projects API service

- **文件**: `src/features/projects/services/projects.ts`
- **操作**: 实现 projectsApi.list / get / create / update / remove / getStats
- **AC**: T40 测试全部通过

### T42 依赖: T41 | 测试: useProjects hooks

- **文件**: `src/features/projects/hooks/useProjects.test.ts`
- **操作**: 测试 useProjectList 查询 + 搜索/分页参数；测试 useProjectDetail 按 ID 查询；测试 useProjectStats 查询；测试 useCreateProject mutation 成功后 invalidateQueries(['projects'])；测试 useUpdateProject mutation；测试 useDeleteProject mutation
- **AC**: 6 个测试通过

### T43 依赖: T42 | 实现: useProjects hooks

- **文件**: `src/features/projects/hooks/useProjects.ts`
- **操作**: 实现 useProjectList, useProjectDetail, useProjectStats, useCreateProject, useUpdateProject, useDeleteProject
- **AC**: T42 测试全部通过

### T44 依赖: T43, T30 | 测试: ProjectListPage

- **对应 plan**: §3.2
- **文件**: `src/features/projects/components/ProjectListPage.test.tsx`
- **操作**: 测试渲染项目卡片网格（3列）；测试搜索框输入触发查询；测试"新建项目"按钮打开 Modal；测试空状态显示 FolderOpen + "创建第一个项目"
- **AC**: 4 个测试通过

### T45 依赖: T44 | 实现: ProjectListPage

- **文件**: `src/features/projects/components/ProjectListPage.tsx`
- **操作**: 搜索 Input + Card 网格（颜色装饰条 + 前缀 Tag + 统计行 + 操作按钮）+ CreateProjectModal 触发 + 空状态
- **AC**: T44 测试全部通过

### T46 依赖: T43 | 测试: CreateProjectModal

- **对应 plan**: §3.3
- **文件**: `src/features/projects/components/CreateProjectModal.test.tsx`
- **操作**: 测试表单渲染（名称/前缀/描述）；测试名称为空校验失败；测试前缀非 2-4 位大写字母校验失败；测试前缀唯一性（409 模拟）；测试提交成功关闭 + 刷新列表
- **AC**: 5 个测试通过

### T47 依赖: T46 | 实现: CreateProjectModal

- **文件**: `src/features/projects/components/CreateProjectModal.tsx`
- **操作**: Arco Modal(600px) + React Hook Form + zod（name 2-255, prefix `/^[A-Z]{2,4}$/`, description optional）；409 → 字段错误
- **AC**: T46 测试全部通过

### T48 依赖: T43 | 测试: 项目编辑/删除

- **对应 plan**: §3.4
- **文件**: `src/features/projects/components/EditProjectModal.test.tsx`（新建）
- **操作**: 测试编辑 Modal 预填数据；测试修改成功关闭 + 刷新；测试删除 Popconfirm 确认后跳转项目列表
- **AC**: 3 个测试通过

### T49 依赖: T48 | 实现: 项目编辑/删除

- **文件**: `src/features/projects/components/EditProjectModal.tsx`（新建）
- **操作**: 复用 CreateProjectModal 表单，预填名称/前缀/描述，调用 useUpdateProject；删除用 Popconfirm + useDeleteProject
- **AC**: T48 测试全部通过

### T50 依赖: T43, T34 | 测试: ProjectDashboard

- **对应 plan**: §3.5
- **文件**: `src/features/projects/components/ProjectDashboard.test.tsx`
- **操作**: 测试渲染 4 个统计卡片（用例总数/通过率/覆盖率/AI生成数）；测试新项目空状态显示快速引导；测试操作按钮存在（上传文档/发起生成/新建用例）
- **AC**: 3 个测试通过

### T51 依赖: T50 | 实现: ProjectDashboard

- **文件**: `src/features/projects/components/ProjectDashboard.tsx`
- **操作**: 标题行 + 操作按钮 + 4 StatsCard + 趋势图区域(60%) + 最近任务列表(40%) + 新项目引导
- **AC**: T50 测试全部通过

---

## Phase 4: 模块管理 (plan §4)

### T52 依赖: T04, T03, T01 | 测试: Modules API service

- **文件**: `src/features/modules/services/modules.test.ts`
- **操作**: 测试 list(projectId) GET /projects/{id}/modules；测试 create(projectId, data) POST；测试 remove(moduleId) DELETE /modules/{id}
- **AC**: 3 个测试通过

### T53 依赖: T52 | 实现: Modules API service

- **文件**: `src/features/modules/services/modules.ts`
- **操作**: 实现 modulesApi.list / create / remove
- **AC**: T52 测试全部通过

### T54 依赖: T53 | 测试: useModules hooks

- **文件**: `src/features/modules/hooks/useModules.test.ts`
- **操作**: 测试 useModuleList 按 projectId 查询；测试 useCreateModule mutation + invalidate；测试 useDeleteModule mutation + invalidate
- **AC**: 3 个测试通过

### T55 依赖: T54 | 实现: useModules hooks

- **文件**: `src/features/modules/hooks/useModules.ts`
- **操作**: 实现 useModuleList, useCreateModule, useDeleteModule
- **AC**: T54 测试全部通过

### T56 依赖: T55, T36 | 测试: ModuleManagePage

- **对应 plan**: §4.1
- **文件**: `src/features/modules/components/ModuleManagePage.test.tsx`
- **操作**: 测试 SplitPanel 布局（左模块树 + 右编辑区）；测试新增模块；测试删除确认（输入名称匹配）；测试 admin+ 权限守卫
- **AC**: 4 个测试通过

### T57 依赖: T56 | 实现: ModuleManagePage

- **文件**: `src/features/modules/components/ModuleManagePage.tsx`
- **操作**: SplitPanel — 左 Arco Tree(280px, 节点: 名称+缩写Tag) + 右编辑表单；新增 Modal(名称+缩写)；删除 Modal(输入名称确认+级联提示)
- **AC**: T56 测试全部通过

---

## Phase 5: 知识库 (plan §5)

> **验收用例**: TC-K01 上传文档 / TC-K02 文档详情与分块

### T58 依赖: T04, T03, T01 | 测试: Documents API service

- **文件**: `src/features/documents/services/documents.test.ts`
- **操作**: 测试 list(projectId, params)；测试 get(id) 返回 DocumentDetail；测试 create() multipart 上传；测试 remove(id)；测试 getChunks(docId)
- **AC**: 5 个测试通过

### T59 依赖: T58 | 实现: Documents API service

- **文件**: `src/features/documents/services/documents.ts`
- **操作**: 实现 documentsApi.list / get / create / remove / getChunks
- **AC**: T58 测试全部通过

### T60 依赖: T59 | 测试: useDocuments hooks

- **文件**: `src/features/documents/hooks/useDocuments.test.ts`
- **操作**: 测试 useDocumentList + 筛选参数；测试 useDocumentDetail；测试 useUploadDocument mutation + invalidate；测试 useDeleteDocument mutation；测试 useDocumentChunks
- **AC**: 5 个测试通过

### T61 依赖: T60 | 实现: useDocuments hooks

- **文件**: `src/features/documents/hooks/useDocuments.ts`
- **操作**: 实现 useDocumentList, useDocumentDetail, useUploadDocument, useDeleteDocument, useDocumentChunks
- **AC**: T60 测试全部通过

### T62 依赖: T61, T28, T30 | 测试: KnowledgeListPage

- **对应 plan**: §5.1
- **文件**: `src/features/documents/components/KnowledgeListPage.test.tsx`
- **操作**: 测试渲染文档表格（名称+类型图标, 类型Tag, 状态Tag, 分块数, 上传时间）；测试按类型/状态筛选；测试空状态；测试"上传文档"按钮打开 Modal
- **AC**: 4 个测试通过

### T63 依赖: T62 | 实现: KnowledgeListPage

- **文件**: `src/features/documents/components/KnowledgeListPage.tsx`
- **操作**: 筛选栏(类型/状态/搜索) + SearchTable + 空状态(FileText) + UploadDocumentModal 触发
- **AC**: T62 测试全部通过

### T64 依赖: T61 | 测试: UploadDocumentModal

- **对应 plan**: §5.2
- **文件**: `src/features/documents/components/UploadDocumentModal.test.tsx`
- **操作**: 测试表单渲染（名称/类型/文件上传）；测试文件类型限制；测试提交成功关闭 + 刷新列表
- **AC**: 3 个测试通过

### T65 依赖: T64 | 实现: UploadDocumentModal

- **文件**: `src/features/documents/components/UploadDocumentModal.tsx`
- **操作**: Arco Modal(600px) + 名称(2-255) + 类型 Select(PRD/Figma/API Spec/Swagger/Markdown) + Upload(drag, 限制 .docx/.pdf/.md/.json/.yaml, ≤50MB)
- **AC**: T64 测试全部通过

### T66 依赖: T61, T36 | 测试: DocumentDetailPage

- **对应 plan**: §5.3
- **文件**: `src/features/documents/components/DocumentDetailPage.test.tsx`
- **操作**: 测试 SplitPanel 布局（左信息面板 + 右分块列表）；测试文档信息区（名称/类型Tag/状态Tag/Steps状态流）；测试分块列表（序号/预览/展开/引用计数）；测试 processing 状态显示 Spin
- **AC**: 4 个测试通过

### T67 依赖: T66 | 实现: DocumentDetailPage

- **文件**: `src/features/documents/components/DocumentDetailPage.tsx`
- **操作**: SplitPanel — 左 300px(名称h2 + 类型Tag + 状态Tag + 上传人/时间 + 分块数 + Steps) + 右分块列表(序号 + 3行截断 + 展开 + 引用图标)
- **AC**: T66 测试全部通过

### T68 依赖: T61 | 测试: FigmaIntegrationPage

- **对应 plan**: §5.4
- **文件**: `src/features/documents/components/FigmaIntegrationPage.test.tsx`（新建）
- **操作**: 测试区域1 连接配置渲染（认证方式Radio + 令牌Input + 测试连接按钮）；测试区域2 URL输入 + 解析按钮；测试区域3 节点树(Arco Tree with Checkbox)
- **AC**: 3 个测试通过

### T69 依赖: T68 | 实现: FigmaIntegrationPage

- **文件**: `src/features/documents/components/FigmaIntegrationPage.tsx`（新建）
- **操作**: 分区域表单 — 区域1: Radio(个人令牌/OAuth) + Input.Password + 测试连接；区域2: Figma URL + 解析；区域3: Arco Tree(Checkbox) + 底部取消/确认导入
- **AC**: T68 测试全部通过

---

## Phase 6: AI 生成模块 (plan §6)

> 核心差异化功能。AI 视觉特征：品牌紫色辉光、Sparkles 图标、AI Gradient 按钮、glow-pulse 动画。
>
> **验收用例**: TC-G01 生成任务(充足) / TC-G02 生成任务(为空) / TC-G03 查看草稿

### T70 依赖: T04, T03, T01 | 测试: Generation API service

- **文件**: `src/features/generation/services/generation.test.ts`
- **操作**: 测试 createTask() POST /generation/tasks；测试 getTask(id)；测试 listTasks(projectId, params)（缺失API，前端仍需实现调用）；测试 getTaskDrafts(taskId)
- **AC**: 4 个测试通过

### T71 依赖: T70 | 实现: Generation API service

- **文件**: `src/features/generation/services/generation.ts`
- **操作**: 实现 generationApi.createTask / getTask / listTasks / getTaskDrafts
- **AC**: T70 测试全部通过

### T72 依赖: T71 | 测试: useGeneration hooks

- **文件**: `src/features/generation/hooks/useGeneration.test.ts`
- **操作**: 测试 useGenerationTasks 列表查询；测试 useGenerationTask 详情查询；测试 useCreateGenerationTask mutation + invalidate + navigate
- **AC**: 3 个测试通过

### T73 依赖: T72 | 实现: useGeneration hooks

- **文件**: `src/features/generation/hooks/useGeneration.ts`
- **操作**: 实现 useGenerationTasks, useGenerationTask, useCreateGenerationTask
- **AC**: T72 测试全部通过

### T74 依赖: T71 | 测试: usePollingTask

- **对应 plan**: §6.3 轮询
- **文件**: `src/features/generation/hooks/usePollingTask.test.ts`
- **操作**: 测试 pending/processing 状态触发轮询(refetchInterval=5000)；测试 completed/failed 停止轮询；测试 taskId 为空不查询
- **AC**: 3 个测试通过

### T75 依赖: T74 | 实现: usePollingTask

- **文件**: `src/features/generation/hooks/usePollingTask.ts`
- **操作**: useQuery + refetchInterval 动态（pending/processing → 5000ms，其他 → false）
- **AC**: T74 测试全部通过

### T76 依赖: T73, T71 | 测试: GenerationTaskListPage

- **对应 plan**: §6.1
- **文件**: `src/features/generation/components/GenerationTaskListPage.test.tsx`
- **操作**: 测试渲染任务列表（任务ID紫色monospace, 需求摘要, 数量, 状态Tag, 时间）；测试"新建生成任务"按钮(AI Gradient)；测试状态筛选；测试双击行跳转详情；测试 failed 行显示"重试"
- **AC**: 5 个测试通过

### T77 依赖: T76 | 实现: GenerationTaskListPage

- **文件**: `src/features/generation/components/GenerationTaskListPage.tsx`
- **操作**: SearchTable + 状态筛选 + "新建任务"按钮(AI Gradient + Sparkles 图标) + 任务ID紫色monospace + failed重试(Popconfirm) + processing行辉光动画
- **AC**: T76 测试全部通过

### T78 依赖: T73 | 测试: NewGenerationTaskPage

- **对应 plan**: §6.2
- **文件**: `src/features/generation/components/NewGenerationTaskPage.test.tsx`
- **操作**: 测试知识库就绪度指示器（充足绿/不足黄/为空红+禁用按钮）；测试模块选择(必填)；测试需求描述(≥10字)；测试用例数量(1-20, 默认5)；测试高级选项折叠/展开；测试提交成功跳转任务详情；测试知识库为空时禁止生成
- **AC**: 7 个测试通过

### T79 依赖: T78 | 实现: NewGenerationTaskPage

- **文件**: `src/features/generation/components/NewGenerationTaskPage.tsx`
- **操作**: 居中表单(max-width 720px) + 知识库就绪度指示器 + 模块Select + 需求描述TextArea(≥10字) + 用例数量InputNumber(1-20, 默认5) + Collapse高级选项(场景类型/优先级/用例类型/生成模式) + RAG降级弹窗
- **AC**: T78 测试全部通过

### T80 依赖: T75, T71, T28 | 测试: TaskDetailPage

- **对应 plan**: §6.3
- **文件**: `src/features/generation/components/TaskDetailPage.test.tsx`
- **操作**: 测试 pending 状态显示 Spin + 微弱 Glow；测试 processing 显示 Progress + shimmer + 轮询；测试 completed 显示草稿列表(ai-reveal动画)；测试 failed 显示 Alert + 重试按钮；测试草稿列表工具栏（全选+批量确认+批量拒绝）
- **AC**: 5 个测试通过

### T81 依赖: T80 | 实现: TaskDetailPage

- **文件**: `src/features/generation/components/TaskDetailPage.tsx`
- **操作**: 状态卡片(按状态渲染不同UI) + usePollingTask + 草稿列表(checkbox + 标题/类型/优先级/置信度Tag + ai-reveal动画) + 工具栏(全选+批量确认AI Gradient+批量拒绝)
- **AC**: T80 测试全部通过

---

## Phase 7: 用例管理 (plan §7)

> **验收用例**: TC-C01 创建用例 / TC-C02 AI用例详情 / TC-C03 编号验证

### T82 依赖: T04, T03, T01 | 测试: TestCases API service

- **文件**: `src/features/testcases/services/testcases.test.ts`
- **操作**: 测试 list(params 含筛选)；测试 get(id) 返回 CaseDetail；测试 create(data)；测试 update(id, data)；测试 remove(id)
- **AC**: 5 个测试通过

### T83 依赖: T82 | 实现: TestCases API service

- **文件**: `src/features/testcases/services/testcases.ts`
- **操作**: 实现 testcasesApi.list / get / create / update / remove
- **AC**: T82 测试全部通过

### T84 依赖: T83 | 测试: useTestCases hooks

- **文件**: `src/features/testcases/hooks/useTestCases.test.ts`
- **操作**: 测试 useCaseList 按 project_id + 筛选查询；测试 useCaseDetail 查询；测试 useCreateTestCase mutation + invalidate；测试 useUpdateTestCase mutation；测试 useDeleteTestCase mutation
- **AC**: 5 个测试通过

### T85 依赖: T84 | 实现: useTestCases hooks

- **文件**: `src/features/testcases/hooks/useTestCases.ts`
- **操作**: 实现 useCaseList, useCaseDetail, useCreateTestCase, useUpdateTestCase, useDeleteTestCase
- **AC**: T84 测试全部通过

### T86 依赖: T85, T28, T30, T36 | 测试: CaseListPage

- **对应 plan**: §7.1
- **文件**: `src/features/testcases/components/CaseListPage.test.tsx`
- **操作**: 测试 SplitPanel 布局（左模块树+右表格）；测试左侧模块树（全部+各模块+用例数Badge）；测试表格列（编号monospace/标题/类型Tag/优先级Tag/状态Tag/创建人/更新时间/操作）；测试筛选栏（搜索/状态/类型/优先级）；测试批量操作栏浮现；测试空状态（FileCheck+手动创建+AI生成按钮）；测试导入/导出按钮存在
- **AC**: 7 个测试通过

### T87 依赖: T86 | 实现: CaseListPage

- **文件**: `src/features/testcases/components/CaseListPage.tsx`
- **操作**: SplitPanel(左240px模块树) + 筛选栏 + SearchTable + 批量操作栏 + 导入Modal(Upload drag .xlsx/.csv) + 导出Dropdown(全部/已选中/筛选结果) + 空状态
- **AC**: T86 测试全部通过

### T88 依赖: T85, T32 | 测试: CreateCaseDrawer

- **对应 plan**: §7.2
- **文件**: `src/features/testcases/components/CreateCaseDrawer.test.tsx`
- **操作**: 测试 Drawer 打开渲染表单；测试模块选择必填校验；测试标题必填(2-500)校验；测试测试步骤≥1条校验；测试提交成功关闭+刷新列表
- **AC**: 5 个测试通过

### T89 依赖: T88 | 实现: CreateCaseDrawer

- **文件**: `src/features/testcases/components/CreateCaseDrawer.tsx`
- **操作**: Arco Drawer(640px) + React Hook Form + zod（module_id必填, title 2-500, steps≥1, expected必填, case_type必填, priority必填）+ ArrayEditor(steps) + ArrayEditor(前置条件, min 0)
- **AC**: T88 测试全部通过

### T90 依赖: T85, T28, T38 | 测试: CaseDetailPage 页面

- **对应 plan**: §7.5
- **文件**: `src/features/testcases/components/CaseDetailPage.test.tsx`
- **操作**: 测试页头（返回链接+编号monospace h1+标题+状态Tag大号+编辑/复制/删除按钮）；测试基本信息卡片(3列网格)；测试 AI 元数据区(Collapse, 紫色背景+Sparkles, 置信度+引用+模型版本, 源文档变更⚠️/删除Alert)；测试执行历史表格；测试编号格式渲染（TC-C03: `ECO-USR-YYYYMMDD-NNN` monospace 样式）
- **AC**: 5 个测试通过

### T91 依赖: T90 | 实现: CaseDetailPage 页面

- **文件**: `src/features/testcases/components/CaseDetailPage.tsx`
- **操作**: 页头(返回+编号monospace h1+标题+状态Tag大号+操作按钮) + 基本信息卡片(3列网格) + AI元数据区(Collapse) + 执行历史表格 + 编号 monospace 渲染
- **AC**: T90 测试全部通过

### T92 依赖: T91, T89 | 测试: 用例编辑/复制

- **对应 plan**: §7.3, §7.4
- **文件**: `src/features/testcases/components/CaseDetailPage.edit.test.tsx`（新建）
- **操作**: 测试编辑按钮打开 Drawer 预填数据；测试编辑成功调用 PUT /testcases/{id}；测试复制按钮打开 Drawer 预填+标题前缀"[副本]"；测试复制成功调用 POST /testcases
- **AC**: 4 个测试通过

### T93 依赖: T92 | 实现: 用例编辑/复制集成

- **文件**: `src/features/testcases/components/CaseDetailPage.tsx`（修改）
- **操作**: 在 CaseDetailPage 中集成编辑/复制功能：编辑按钮→复用 CreateCaseDrawer(预填) + useUpdateTestCase；复制按钮→复用 CreateCaseDrawer(预填+标题"[副本]") + useCreateTestCase
- **AC**: T92 测试全部通过

---

## Phase 8: 草稿箱 (plan §8)

> **验收用例**: TC-D01 单条确认 / TC-D02 批量确认 / TC-D03 拒绝草稿

### T94 依赖: T04, T03, T01 | 测试: Drafts API service

- **文件**: `src/features/drafts/services/drafts.test.ts`
- **操作**: 测试 getDrafts(params)（缺失API，前端实现调用）；测试 confirmDraft(draftId, moduleId) POST /generation/drafts/{id}/confirm；测试 rejectDraft(draftId, data) POST /generation/drafts/{id}/reject；测试 batchConfirm(draftIds, moduleId) POST /generation/drafts/batch-confirm
- **AC**: 4 个测试通过

### T95 依赖: T94 | 实现: Drafts API service

- **文件**: `src/features/drafts/services/drafts.ts`
- **操作**: 实现 draftsApi.getDrafts / confirmDraft / rejectDraft / batchConfirm
- **AC**: T94 测试全部通过

### T96 依赖: T95 | 测试: useDrafts hooks

- **文件**: `src/features/drafts/hooks/useDrafts.test.ts`
- **操作**: 测试 useDraftList 查询+筛选；测试 useConfirmDraft mutation + invalidate + navigate；测试 useRejectDraft mutation；测试 useBatchConfirm mutation（验证调用 batchConfirm API + invalidate + 更新列表状态）；测试 usePendingDraftCount(refetchInterval=30000)
- **AC**: 5 个测试通过

### T97 依赖: T96 | 实现: useDrafts hooks

- **文件**: `src/features/drafts/hooks/useDrafts.ts`
- **操作**: 实现 useDraftList, useConfirmDraft, useRejectDraft, useBatchConfirm, usePendingDraftCount
- **AC**: T96 测试全部通过

### T98 依赖: T97, T28 | 测试: DraftListPage

- **对应 plan**: §8.1
- **文件**: `src/features/drafts/components/DraftListPage.test.tsx`
- **操作**: 测试渲染草稿列表（标题/项目/模块/置信度Tag/时间）；测试筛选（项目联动模块+状态+搜索）；测试批量勾选 checkbox；测试勾选后点击"批量确认"按钮调用 batchConfirm mutation + success_count=3（TC-D02）；测试"批量拒绝"按钮；测试侧边栏Badge显示pending数
- **AC**: 6 个测试通过

### T99 依赖: T98 | 实现: DraftListPage

- **文件**: `src/features/drafts/components/DraftListPage.tsx`
- **操作**: 筛选栏(项目Select联动模块+状态+搜索) + SearchTable + checkbox批量选择 + 批量确认(AI Gradient)/拒绝按钮(选中后显示) + 点击跳转确认页
- **AC**: T98 测试全部通过

### T100 依赖: T97, T32, T36, T38 | 测试: DraftConfirmPage 骨架 ⭐

- **对应 plan**: §8.2 核心页面
- **文件**: `src/features/drafts/components/DraftConfirmPage.test.tsx`
- **操作**: 测试导航栏（返回链接+"第N/M条"+圆点导航）；测试左右分栏布局（60%编辑+40%引用）；测试左侧编辑区（标题Input+前置条件ArrayEditor+步骤ArrayEditor+预期结果TextArea+类型Select+优先级Select）；测试右侧引用来源面板（文档标题+相似度+引用原文）；测试底部操作栏三个按钮存在
- **AC**: 5 个测试通过

### T101 依赖: T100 | 实现: DraftConfirmPage 骨架 ⭐

- **文件**: `src/features/drafts/components/DraftConfirmPage.tsx`
- **操作**: 导航栏(返回+进度+圆点+键盘←/→) + SplitPanel(左60%表单 + 右40%ReferencePanel紫色底色 rgba(123,97,255,0.03)) + 底部固定操作栏(拒绝danger/保存default/确认AI Gradient)
- **AC**: T100 测试全部通过

### T102 依赖: T101 | 测试: DraftConfirmPage 操作与切换保护 ⭐

- **对应 plan**: §8.2 操作栏 + 切换保护
- **文件**: `src/features/drafts/components/DraftConfirmPage.actions.test.tsx`（新建）
- **操作**: 测试"拒绝"弹出Modal（原因Radio: 重复/无关/低质量/其他 + 反馈TextArea）→ 调用 rejectDraft；测试"确认"弹出Modal（选择目标模块Select）→ 调用 confirmDraft → Message.success("用例 {number} 已创建")；测试"保存修改"仅暂存 React state（无API调用）；测试切换草稿前未保存编辑弹出确认Dialog；测试草稿间切换 auto-save 到 React state
- **AC**: 5 个测试通过

### T103 依赖: T102 | 实现: DraftConfirmPage 操作与切换保护 ⭐

- **文件**: `src/features/drafts/components/DraftConfirmPage.tsx`（修改）
- **操作**: 实现拒绝Modal + 确认Modal(选模块) + 保存暂存(React state) + 切换保护Dialog(useBlocker/beforeunload) + 草稿间切换 auto-save
- **AC**: T102 测试全部通过

---

## Phase 9: 测试执行 (plan §9)

> **验收用例**: TC-PL01 创建计划 / TC-PL02 录入结果 / TC-PL03 状态流转

### T104 依赖: T04, T03, T01 | 测试: Plans API service

- **文件**: `src/features/plans/services/plans.test.ts`
- **操作**: 测试 list(params)；测试 get(id) 返回 PlanDetail 含 stats+cases+results；测试 create(data)；测试 addCases(planId, caseIds)；测试 removeCase(planId, caseId)；测试 recordResult(planId, data)；测试 updateStatus(planId, status) PATCH /plans/{id}/status
- **AC**: 7 个测试通过

### T105 依赖: T104 | 实现: Plans API service

- **文件**: `src/features/plans/services/plans.ts`
- **操作**: 实现 plansApi.list / get / create / addCases / removeCase / recordResult / updateStatus
- **AC**: T104 测试全部通过

### T106 依赖: T105 | 测试: usePlans hooks

- **文件**: `src/features/plans/hooks/usePlans.test.ts`
- **操作**: 测试 usePlanList；测试 usePlanDetail；测试 useCreatePlan mutation + navigate；测试 useRecordResult mutation + invalidate；测试 useUpdatePlanStatus mutation
- **AC**: 5 个测试通过

### T107 依赖: T106 | 实现: usePlans hooks

- **文件**: `src/features/plans/hooks/usePlans.ts`
- **操作**: 实现 usePlanList, usePlanDetail, useCreatePlan, useRecordResult, useUpdatePlanStatus
- **AC**: T106 测试全部通过

### T108 依赖: T107, T28 | 测试: PlanListPage

- **对应 plan**: §9.1
- **文件**: `src/features/plans/components/PlanListPage.test.tsx`
- **操作**: 测试渲染计划列表（名称可点击/状态Tag/用例数/通过率+迷你进度条/创建人/时间）；测试状态筛选；测试"新建计划"按钮
- **AC**: 3 个测试通过

### T109 依赖: T108 | 实现: PlanListPage

- **文件**: `src/features/plans/components/PlanListPage.tsx`
- **操作**: SearchTable + 状态筛选 + StatusTag + 通过率(百分比+迷你Progress) + "新建计划"按钮
- **AC**: T108 测试全部通过

### T110 依赖: T107 | 测试: NewPlanPage

- **对应 plan**: §9.2
- **文件**: `src/features/plans/components/NewPlanPage.test.tsx`
- **操作**: 测试 SplitPanel 布局（左表单+右用例选择）；测试计划名称必填(2-255)；测试用例选择面板 Tab 切换（可选用例/已选用例）；测试用例筛选（模块/类型/优先级）；测试已选用例 Badge + 移除；测试提交成功跳转计划详情
- **AC**: 6 个测试通过

### T111 依赖: T110 | 实现: NewPlanPage

- **文件**: `src/features/plans/components/NewPlanPage.tsx`
- **操作**: SplitPanel — 左(名称Input + 描述TextArea) + 右(Tab: 可选用例(筛选+Checkbox表) / 已选用例(列表+移除+Badge)) + React Hook Form + zod
- **AC**: T110 测试全部通过

### T112 依赖: T107, T34 | 测试: PlanDetailPage 骨架

- **对应 plan**: §9.3
- **文件**: `src/features/plans/components/PlanDetailPage.test.tsx`
- **操作**: 测试操作按钮按状态变化（draft/active/completed/archived）；测试 5 列统计卡片（总数/通过/失败/阻塞/跳过）；测试执行进度条(Arco Progress, 品牌紫色)；测试用例执行表格渲染（紧凑36px行高）
- **AC**: 4 个测试通过

### T113 依赖: T112 | 实现: PlanDetailPage 骨架

- **文件**: `src/features/plans/components/PlanDetailPage.tsx`
- **操作**: 操作按钮(按状态, 4种状态组) + 5 StatsCard(绿/红/橙/灰装饰线) + Progress进度条 + 紧凑Table(36px行高, 编号/标题/类型+优先级/执行结果/执行人/时间)
- **AC**: T112 测试全部通过

### T114 依赖: T113 | 测试: 快捷录入

- **对应 plan**: §9.3 快捷录入
- **文件**: `src/features/plans/components/PlanDetailPage.quickEntry.test.tsx`（新建）
- **操作**: 测试点击执行结果列弹出内联 Select；测试选择值自动提交(POST /plans/{id}/results)；测试提交后整行闪烁对应色；测试 Toast "已录入：通过" + "撤销"链接显示；测试 3s 内点击撤销恢复原值
- **AC**: 5 个测试通过

### T115 依赖: T114 | 实现: 快捷录入

- **文件**: `src/features/plans/components/PlanDetailPage.tsx`（修改）
- **操作**: 执行结果列可点击→内联Arco Select→选择值触发 useRecordResult→整行闪烁CSS动画→Toast+撤销链接(3s定时器，撤销调用 delete result)
- **AC**: T114 测试全部通过

### T116 依赖: T113 | 测试: ResultRecordModal (详细录入)

- **对应 plan**: §9.3 详细录入
- **文件**: `src/features/plans/components/ResultRecordModal.test.tsx`
- **操作**: 测试渲染用例信息(只读)；测试执行结果 Radio（pass/fail/block/skip）；测试备注 TextArea；测试提交成功关闭+刷新
- **AC**: 4 个测试通过

### T117 依赖: T116 | 实现: ResultRecordModal

- **文件**: `src/features/plans/components/ResultRecordModal.tsx`
- **操作**: Arco Modal(500px) + 用例信息(只读) + Radio(pass/fail/block/skip) + TextArea(备注) + React Hook Form + zod
- **AC**: T116 测试全部通过

### T118 依赖: T115, T117 | 测试: 批量录入

- **对应 plan**: §9.3 批量录入
- **文件**: `src/features/plans/components/PlanDetailPage.batchEntry.test.tsx`（新建）
- **操作**: 测试选中多条用例后"批量录入"按钮显示；测试点击打开 Modal（结果Radio+备注TextArea）；测试提交批量调用 recordResult；测试批量完成后 invalidate 缓存 + 统计更新
- **AC**: 4 个测试通过

### T119 依赖: T118 | 实现: 批量录入

- **文件**: `src/features/plans/components/PlanDetailPage.tsx`（修改）
- **操作**: 选中行→"批量录入"按钮→Modal(Radio+TextArea)→Promise.all 批量调用→完成刷新统计
- **AC**: T118 测试全部通过

---

## Phase 10: 配置管理 (plan §10)

### T120 依赖: T04, T03, T01 | 测试: Configs API service

- **文件**: `src/features/configs/services/configs.test.ts`
- **操作**: 测试 list(projectId)；测试 set(projectId, key, value) PUT；测试 remove(projectId, key)；测试 importConfigs(projectId, configs) POST；测试 exportConfigs(projectId) GET
- **AC**: 5 个测试通过

### T121 依赖: T120 | 实现: Configs API service

- **文件**: `src/features/configs/services/configs.ts`
- **操作**: 实现 configsApi.list / set / remove / importConfigs / exportConfigs
- **AC**: T120 测试全部通过

### T122 依赖: T121 | 测试: useConfigs hooks

- **文件**: `src/features/configs/hooks/useConfigs.test.ts`
- **操作**: 测试 useConfigList；测试 useSetConfig mutation + invalidate；测试 useDeleteConfig mutation；测试 useImportConfigs mutation；测试 useExportConfigs
- **AC**: 5 个测试通过

### T123 依赖: T122 | 实现: useConfigs hooks

- **文件**: `src/features/configs/hooks/useConfigs.ts`
- **操作**: 实现 useConfigList, useSetConfig, useDeleteConfig, useImportConfigs, useExportConfigs
- **AC**: T122 测试全部通过

### T124 依赖: T123, T30 | 测试: ConfigManagePage

- **对应 plan**: §10.1
- **文件**: `src/features/configs/components/ConfigManagePage.test.tsx`
- **操作**: 测试渲染配置表格（键名/值/描述/操作）；测试"新增配置"按钮打开 Modal；测试编辑 Modal 预填数据；测试删除 Popconfirm；测试"导入 JSON"按钮打开 Modal；测试"导出 JSON"按钮
- **AC**: 6 个测试通过

### T125 依赖: T124 | 实现: ConfigManagePage

- **文件**: `src/features/configs/components/ConfigManagePage.tsx`
- **操作**: 配置表格 + 新增/编辑 Modal(500px, 键名Input+值JsonEditor+描述Input) + 删除 Popconfirm + 导入 Modal(粘贴JSON+预览+确认) + 导出(下载)
- **AC**: T124 测试全部通过

---

## Phase 11: 全局布局与路由集成 (plan §3.1)

### T126 依赖: T08 | 测试: useAppStore 精简

- **文件**: `src/store/useAppStore.test.ts`
- **操作**: 测试 sidebarCollapsed 默认值和切换
- **AC**: 测试通过

### T127 依赖: T126 | 实现: useAppStore 精简

- **文件**: `src/store/useAppStore.ts`
- **操作**: 仅保留 sidebarCollapsed + toggleSidebar
- **AC**: T126 测试通过

### T128 依赖: T97 | 测试: Sidebar

- **文件**: `src/components/layout/Sidebar.test.tsx`
- **操作**: 测试渲染菜单项（仪表盘/用例/计划/知识库/AI生成/草稿箱）；测试选中态样式；测试折叠/展开；测试草稿箱 Badge 显示 pendingDraftCount
- **AC**: 4 个测试通过

### T129 依赖: T128 | 实现: Sidebar

- **文件**: `src/components/layout/Sidebar.tsx`
- **操作**: Arco Menu + 图标 + useAppStore(sidebarCollapsed) + usePendingDraftCount(Badge, 30s刷新) + 用户信息区
- **AC**: T128 测试全部通过

### T130 `[P]` 测试: Header

- **文件**: `src/components/layout/Header.test.tsx`
- **操作**: 测试折叠按钮；测试面包屑；测试用户下拉菜单（退出登录）
- **AC**: 3 个测试通过

### T131 依赖: T130 | 实现: Header

- **文件**: `src/components/layout/Header.tsx`
- **操作**: 折叠 Button + Arco Breadcrumb + Avatar Dropdown(退出)
- **AC**: T130 测试全部通过

### T132 依赖: T129, T131 | 测试: AppLayout

- **文件**: `src/components/layout/AppLayout.test.tsx`
- **操作**: 测试三区布局（Sidebar + Header + Content）；测试侧边栏折叠宽度变化
- **AC**: 2 个测试通过

### T133 依赖: T132 | 实现: AppLayout

- **文件**: `src/components/layout/AppLayout.tsx`
- **操作**: Arco Layout + Layout.Sider(220px/64px) + Layout.Header(48px) + Layout.Content(padding 24px, 独立滚动)
- **AC**: T132 测试全部通过

### T134 `[P]` 实现: NotFoundPage

- **文件**: `src/components/NotFoundPage.tsx`
- **操作**: Arco Result 404 + 返回首页按钮
- **AC**: 渲染 404 页面

### T135 依赖: T23, T133, T134, T10, 全部页面组件 | 实现: 完整路由配置

- **文件**: `src/router/index.tsx`
- **操作**: createBrowserRouter 定义 plan.md 路由清单全部 19 条路由；AuthLayout 包裹 /login, /register + ErrorBoundary；AppLayout + RouteGuard 包裹业务页面；lazy() 加载；NotFoundPage 兜底
- **AC**: 所有路由可访问，懒加载正常

---

## 依赖关系图

```
Phase 0 MSW (T01)
  │
  └─→ Phase 1 基础设施与认证 (T02-T26)
        │
        ├─→ Phase 2 通用组件 (T27-T39) [P]
        │     │
        │     ├─→ Phase 3 项目管理 (T40-T51)
        │     │     │
        │     │     └─→ Phase 4 模块管理 (T52-T57)
        │     │
        │     ├─→ Phase 5 知识库 (T58-T69) [可并行于 Phase 3]
        │     │     │
        │     │     └─→ Phase 6 AI 生成 (T70-T81)
        │     │           │
        │     │           └─→ Phase 8 草稿箱 (T94-T103)
        │     │
        │     ├─→ Phase 7 用例管理 (T82-T93) [可并行于 Phase 5-6]
        │     │
        │     ├─→ Phase 9 测试执行 (T104-T119) [依赖 Phase 7]
        │     │
        │     └─→ Phase 10 配置管理 (T120-T125) [独立，可并行]
        │
        └─→ Phase 11 布局与路由集成 (T126-T135) [依赖全部页面组件]
```

## 统计

| 指标               | 数值            |
| ------------------ | --------------- |
| 总任务数           | 135             |
| 测试任务           | 67 (50%)        |
| 实现任务           | 68              |
| 可并行任务 (`[P]`) | 18              |
| Phase 数           | 12 (含 Phase 0) |
