---
title: 核心功能原子化任务列表
version: 1.0.0
author: 技术组长
status: draft
created: 2026-04-16
based_on: [spec.md, plan.md]
---

# 核心功能原子化任务列表

> **标记说明**：
> - `[P]` = 可并行执行（无前置依赖）
> - `依赖: Txx` = 必须在 Txx 完成后执行
> - **测试任务始终优先于对应实现任务**

---

## Phase 0: 基础设施

### T01 `[P]` 安装核心依赖
- **文件**: `package.json`
- **操作**: 安装 `@tanstack/react-query` `@tanstack/react-query-devtools` `react-hook-form` `@hookform/resolvers` `zod`
- **AC**: `npm install` 无报错，`package.json` 中包含上述依赖

### T02 `[P]` 安装测试依赖
- **文件**: `package.json`
- **操作**: 安装 `vitest` `@testing-library/react` `@testing-library/user-event` `@testing-library/jest-dom` `msw` `jsdom`
- **AC**: `npm install` 无报错

### T03 `[P]` 配置 Vitest
- **文件**: `vitest.config.ts`（新建）, `tsconfig.app.json`
- **操作**: 配置 Vitest 使用 jsdom 环境，路径别名 `@/` → `src/`，添加 `setupFiles`
- **AC**: `npx vitest run` 可执行（空测试通过）

### T04 `[P]` 创建 Vitest setup 文件
- **文件**: `tests/setup.ts`（新建）
- **操作**: 导入 `@testing-library/jest-dom`，`afterEach` 清理 cleanup
- **AC**: setup 文件被 vitest.config.ts 引用

### T05 依赖: T01, T02, T03 | 配置 MSW
- **文件**: `tests/msw/server.ts`（新建）
- **操作**: 使用 `setupServer()` 创建 MSW 实例，在 `tests/setup.ts` 中启动/重置
- **AC**: MSW 可拦截 fetch 请求

### T06 `[P]` 创建类型定义 — 枚举
- **文件**: `src/types/enums.ts`（新建）
- **操作**: 定义所有字面量联合类型（CaseStatus, CaseType, PlanStatus, Priority, ResultStatus, TaskStatus, DraftStatus, DocumentType, DocumentStatus, UserRole, Confidence, SceneType）
- **AC**: TypeScript 编译通过，无 any

### T07 依赖: T06 | 创建类型定义 — API
- **文件**: `src/types/api.ts`（新建）
- **操作**: 定义所有 API 请求/响应类型（PaginatedResponse, UserJSON, Project, ProjectDetail, TestCase, CaseDetail, CaseDraft, TestPlan, PlanDetail, GenerationTask, Document 等）
- **AC**: TypeScript 编译通过，与 OpenAPI schema 完全对齐

### T08 `[P]` 修正 Axios — 消除 any
- **文件**: `src/lib/request.ts`
- **操作**: 替换所有 `<any, T>` 为 `<never, T>`；补全 Token 刷新逻辑（isRefreshing + pendingRequests 队列 + WeakSet 追踪重试）；注册 onAuthExpired 回调
- **AC**: TypeScript 编译通过，`grep -r "any" src/lib/request.ts` 无结果

### T09 依赖: T01 | 创建 React Query Client
- **文件**: `src/lib/query-client.ts`（新建）
- **操作**: 创建 `QueryClient` 实例，默认 staleTime 5min，retry 1 次
- **AC**: 导出 `queryClient` 实例

### T10 依赖: T08, T09 | 创建 Provider 嵌套
- **文件**: `src/app/providers.tsx`（新建）
- **操作**: 创建 `Providers` 组件，嵌套 QueryClientProvider + Arco ConfigProvider（品牌色 #7B61FF）
- **AC**: 组件渲染无报错，React Query Devtools 可用

### T11 依赖: T10 | 重写 App.tsx 入口
- **文件**: `src/app/App.tsx`（新建）, `src/main.tsx`（修改）
- **操作**: App.tsx 使用 Providers 包裹 RouterOutlet；main.tsx 指向新入口；注册 authExpiredHandler
- **AC**: `npm run dev` 启动无报错

### T12 `[P]` 创建工具函数
- **文件**: `src/lib/utils.ts`（修改）
- **操作**: 保留 `cn()` 函数（clsx + tailwind-merge）
- **AC**: `cn('px-4', 'py-2')` 输出正确

---

## Phase 1: 认证模块

### T13 依赖: T05, T07 | 测试: Auth API service
- **文件**: `src/features/auth/services/auth.test.ts`（新建）
- **操作**: 测试 `authApi.login()` 调用 `POST /auth/login` 并返回 token；测试 `authApi.register()` 调用 `POST /auth/register`；测试 `authApi.refresh()` 调用 `POST /auth/refresh`
- **AC**: 3 个测试用例通过（MSW mock）

### T14 依赖: T13 | 实现: Auth API service
- **文件**: `src/features/auth/services/auth.ts`（新建）
- **操作**: 实现 `authApi.login()`, `authApi.register()`, `authApi.refresh()`
- **AC**: T13 测试全部通过

### T15 依赖: T07 | 测试: Auth Zustand store
- **文件**: `src/features/auth/hooks/useAuthStore.test.ts`（新建）
- **操作**: 测试 `login()` 成功写入 localStorage + 更新 state；测试 `logout()` 清理 token + state；测试 `refresh()` 成功更新 token；测试 `refresh()` 失败调用 logout
- **AC**: 4 个测试用例通过

### T16 依赖: T14, T15 | 实现: Auth Zustand store
- **文件**: `src/features/auth/hooks/useAuthStore.ts`（新建）
- **操作**: 实现 AuthState（user, token, refreshToken, isAuthenticated, login, logout, refresh, setUser）
- **AC**: T15 测试全部通过

### T17 依赖: T14 | 测试: useAuth hooks
- **文件**: `src/features/auth/hooks/useAuth.test.ts`（新建）
- **操作**: 测试 `useLogin` mutation 成功调用 store.login；测试 `useRegister` mutation 成功调用 authApi.register；测试失败时抛出错误
- **AC**: 3 个测试用例通过

### T18 依赖: T16, T17 | 实现: useAuth hooks
- **文件**: `src/features/auth/hooks/useAuth.ts`（新建）
- **操作**: 实现 `useLogin`（useMutation + store.login），`useRegister`（useMutation），`useRefresh`
- **AC**: T17 测试全部通过

### T19 依赖: T18 | 测试: LoginPage
- **文件**: `src/features/auth/components/LoginPage.test.tsx`（新建）
- **操作**: 测试渲染邮箱/密码输入框和登录按钮；测试无效邮箱显示验证错误；测试登录成功跳转 /projects；测试登录失败显示错误提示；测试密码少于 8 位显示验证错误
- **AC**: 5 个测试用例通过

### T20 依赖: T19 | 实现: LoginPage
- **文件**: `src/features/auth/components/LoginPage.tsx`（新建）
- **操作**: 使用 Arco Form + React Hook Form 实现；表单校验用 zod schema；提交调用 useLogin mutation；成功跳转 useNavigate('/projects')；失败显示 Message.error
- **AC**: T19 测试全部通过

### T21 `[P]` 依赖: T18 | 测试: RegisterPage
- **文件**: `src/features/auth/components/RegisterPage.test.tsx`（新建）
- **操作**: 测试渲染用户名/邮箱/密码/角色表单；测试字段校验（用户名 ≥3 字符，邮箱格式，密码 ≥8 位）；测试注册成功跳转 /login
- **AC**: 3 个测试用例通过

### T22 依赖: T21 | 实现: RegisterPage
- **文件**: `src/features/auth/components/RegisterPage.tsx`（新建）
- **操作**: Arco Form + React Hook Form + zod；角色选择 Radio（admin/normal）；注册成功跳转登录页
- **AC**: T21 测试全部通过

### T23 依赖: T16 | 测试: RouteGuard
- **文件**: `src/router/RouteGuard.test.tsx`（新建）
- **操作**: 测试未认证时跳转 /login；测试 token 过期时跳转 /login；测试正常 token 渲染 children；测试 requireAdmin=true 且 normal 角色跳转 /projects
- **AC**: 4 个测试用例通过

### T24 依赖: T23 | 实现: RouteGuard
- **文件**: `src/router/RouteGuard.tsx`（新建）
- **操作**: 读取 useAuthStore 判断 isAuthenticated + token 是否过期；未认证 Navigate to /login（保留 location.state.from）；requireAdmin 检查 role
- **AC**: T23 测试全部通过

### T25 `[P]` 实现: AuthLayout
- **文件**: `src/components/layout/AuthLayout.tsx`（新建）
- **操作**: 居中布局容器，灰色背景，`<Outlet />`
- **AC**: 页面居中展示，无侧边栏

---

## Phase 2: 共享业务组件

### T26 `[P]` 测试: StatusTag
- **文件**: `src/components/business/StatusTag.test.tsx`（新建）
- **操作**: 测试 case_status 的 4 种状态颜色正确；测试 priority 的 4 级颜色正确；测试不存在的 status 返回 null；测试自定义 label 覆盖默认文本
- **AC**: 4 个测试用例通过

### T27 依赖: T26 | 实现: StatusTag
- **文件**: `src/components/business/StatusTag.tsx`（新建）
- **操作**: 完整 COLOR_MAP 映射（9 个 category），Arco Tag + Tailwind 样式
- **AC**: T26 测试全部通过

### T28 `[P]` 测试: SearchTable
- **文件**: `src/components/business/SearchTable.test.tsx`（新建）
- **操作**: 测试渲染表格列和数据行；测试 loading 时显示 Spin；测试分页控件渲染；测试传入 columns/data 正确展示
- **AC**: 3 个测试用例通过

### T29 依赖: T28 | 实现: SearchTable
- **文件**: `src/components/business/SearchTable.tsx`（新建）
- **操作**: 封装 Arco Table + Pagination，统一 loading/error/empty 状态，Props 泛型化
- **AC**: T28 测试全部通过

### T30 `[P]` 测试: ArrayEditor
- **文件**: `src/components/business/ArrayEditor.test.tsx`（新建）
- **操作**: 测试初始渲染一行；测试点击添加按钮新增行；测试点击删除按钮移除行；测试 onChange 回调传递最新值
- **AC**: 4 个测试用例通过

### T31 依赖: T30 | 实现: ArrayEditor
- **文件**: `src/components/business/ArrayEditor.tsx`（新建）
- **操作**: 动态行列表，每行 Input + 删除按钮，底部添加按钮，onChange(value: string[]) 回调
- **AC**: T30 测试全部通过

### T32 `[P]` 测试: StatsCard
- **文件**: `src/components/business/StatsCard.test.tsx`（新建）
- **操作**: 测试渲染标题和数值；测试趋势箭头（up/down）；测试自定义 icon
- **AC**: 3 个测试用例通过

### T33 依赖: T32 | 实现: StatsCard
- **文件**: `src/components/business/StatsCard.tsx`（新建）
- **操作**: Arco Card + Statistic，支持 title/value/trend/icon props
- **AC**: T32 测试全部通过

### T34 `[P]` 测试: SplitPanel
- **文件**: `src/components/business/SplitPanel.test.tsx`（新建）
- **操作**: 测试左右面板渲染 children；测试默认分割比例
- **AC**: 2 个测试用例通过

### T35 依赖: T34 | 实现: SplitPanel
- **文件**: `src/components/business/SplitPanel.tsx`（新建）
- **操作**: CSS resize 或拖拽分割线实现左右分栏，最小宽度限制
- **AC**: T34 测试全部通过

### T36 `[P]` 测试: ReferencePanel
- **文件**: `src/components/business/ReferencePanel.test.tsx`（新建）
- **操作**: 测试渲染引用块列表；测试空列表显示"无引用来源"
- **AC**: 2 个测试用例通过

### T37 依赖: T36 | 实现: ReferencePanel
- **文件**: `src/components/business/ReferencePanel.tsx`（新建）
- **操作**: 接收 ReferencedChunk[] props，展示文档标题 + 相似度分数 + "查看原文"链接
- **AC**: T36 测试全部通过

---

## Phase 3: 项目管理模块

### T38 依赖: T05, T07 | 测试: Projects API service
- **文件**: `src/features/projects/services/projects.test.ts`（新建）
- **操作**: 测试 list() 调用 GET /projects；测试 get() 调用 GET /projects/:id；测试 create() 调用 POST /projects；测试 update() 调用 PUT /projects/:id；测试 delete() 调用 DELETE /projects/:id；测试 getStats() 调用 GET /projects/:id/stats
- **AC**: 6 个测试用例通过

### T39 依赖: T38 | 实现: Projects API service
- **文件**: `src/features/projects/services/projects.ts`（新建）
- **操作**: 实现 projectsApi.list/get/create/update/delete/getStats
- **AC**: T38 测试全部通过

### T40 依赖: T39 | 测试: useProjects hooks
- **文件**: `src/features/projects/hooks/useProjects.test.ts`（新建）
- **操作**: 测试 useProjectList 返回数据；测试 useProjectDetail 按 ID 查询；测试 useCreateProject mutation 成功后失效 ['projects'] 缓存；测试 useUpdateProject mutation 成功后失效详情缓存；测试 useDeleteProject mutation 成功后失效列表缓存；测试 useProjectStats 查询
- **AC**: 6 个测试用例通过

### T41 依赖: T40 | 实现: useProjects hooks
- **文件**: `src/features/projects/hooks/useProjects.ts`（新建）
- **操作**: 实现 useProjectList, useProjectDetail, useProjectStats, useCreateProject, useUpdateProject, useDeleteProject
- **AC**: T40 测试全部通过

### T42 依赖: T41, T29 | 测试: ProjectListPage
- **文件**: `src/features/projects/components/ProjectListPage.test.tsx`（新建）
- **操作**: 测试渲染项目表格（名称、前缀、描述、创建时间）；测试搜索框输入触发查询；测试分页切换；测试"新建项目"按钮打开弹窗
- **AC**: 4 个测试用例通过

### T43 依赖: T42 | 实现: ProjectListPage
- **文件**: `src/features/projects/components/ProjectListPage.tsx`（新建）
- **操作**: PageHeader + 搜索 Input + SearchTable + CreateProjectModal 触发按钮
- **AC**: T42 测试全部通过

### T44 依赖: T41 | 测试: CreateProjectModal
- **文件**: `src/features/projects/components/CreateProjectModal.test.tsx`（新建）
- **操作**: 测试表单渲染（名称/前缀/描述）；测试名称为空时校验失败；测试前缀非 2-4 位大写字母校验失败；测试提交成功关闭弹窗
- **AC**: 4 个测试用例通过

### T45 依赖: T44 | 实现: CreateProjectModal
- **文件**: `src/features/projects/components/CreateProjectModal.tsx`（新建）
- **操作**: Arco Modal + Form + React Hook Form + zod（name 必填 2-255 字符，prefix 必填 2-4 位大写 `/^[A-Z]{2,4}$/`）
- **AC**: T44 测试全部通过

### T46 依赖: T41, T33 | 测试: ProjectDashboard
- **文件**: `src/features/projects/components/ProjectDashboard.test.tsx`（新建）
- **操作**: 测试渲染 4 个统计卡片；测试趋势图区域存在；测试最近任务列表渲染
- **AC**: 3 个测试用例通过

### T47 依赖: T46 | 实现: ProjectDashboard
- **文件**: `src/features/projects/components/ProjectDashboard.tsx`（新建）
- **操作**: 4 个 StatsCard + 通过率趋势折线图（Arco Chart 或简单 CSS）+ 最近任务列表
- **AC**: T46 测试全部通过

### T48 依赖: T05, T07 | 测试: Modules API service
- **文件**: `src/features/modules/services/modules.test.ts`（新建）
- **操作**: 测试 list(projectId)；测试 create(projectId, data)；测试 delete(moduleId)
- **AC**: 3 个测试用例通过

### T49 依赖: T48 | 实现: Modules API service
- **文件**: `src/features/modules/services/modules.ts`（新建）
- **操作**: 实现 modulesApi.list/create/delete
- **AC**: T48 测试全部通过

### T50 依赖: T49 | 测试: useModules hooks
- **文件**: `src/features/modules/hooks/useModules.test.ts`（新建）
- **操作**: 测试 useModuleList 按 projectId 查询；测试 useCreateModule mutation；测试 useDeleteModule mutation
- **AC**: 3 个测试用例通过

### T51 依赖: T50 | 实现: useModules hooks
- **文件**: `src/features/modules/hooks/useModules.ts`（新建）
- **操作**: 实现 useModuleList, useCreateModule, useDeleteModule
- **AC**: T50 测试全部通过

### T52 依赖: T51, T29 | 测试: ModuleManagePage
- **文件**: `src/features/modules/components/ModuleManagePage.test.tsx`（新建）
- **操作**: 测试渲染模块表格；测试创建模块弹窗；测试删除模块二次确认
- **AC**: 3 个测试用例通过

### T53 依赖: T52 | 实现: ModuleManagePage
- **文件**: `src/features/modules/components/ModuleManagePage.tsx`（新建）
- **操作**: 模块列表 Table + 创建 Modal（名称 + 缩写 2-4 位大写）+ 删除 Popconfirm
- **AC**: T52 测试全部通过

---

## Phase 4: 知识库模块

### T54 依赖: T05, T07 | 测试: Documents API service
- **文件**: `src/features/documents/services/documents.test.ts`（新建）
- **操作**: 测试 list(projectId)；测试 get(id) 返回 DocumentDetail 含 chunk_count；测试 create() 上传文档；测试 delete(id)；测试 getChunks(docId)
- **AC**: 5 个测试用例通过

### T55 依赖: T54 | 实现: Documents API service
- **文件**: `src/features/documents/services/documents.ts`（新建）
- **操作**: 实现 documentsApi.list/get/create/delete/getChunks
- **AC**: T54 测试全部通过

### T56 依赖: T55 | 测试: useDocuments hooks
- **文件**: `src/features/documents/hooks/useDocuments.test.ts`（新建）
- **操作**: 测试 useDocumentList；测试 useDocumentDetail；测试 useUploadDocument mutation；测试 useDeleteDocument mutation；测试 useDocumentChunks
- **AC**: 5 个测试用例通过

### T57 依赖: T56 | 实现: useDocuments hooks
- **文件**: `src/features/documents/hooks/useDocuments.ts`（新建）
- **操作**: 实现 useDocumentList, useDocumentDetail, useUploadDocument, useDeleteDocument, useDocumentChunks
- **AC**: T56 测试全部通过

### T58 依赖: T57, T27, T29 | 测试: KnowledgeListPage
- **文件**: `src/features/documents/components/KnowledgeListPage.test.tsx`（新建）
- **操作**: 测试渲染文档列表（名称、类型 Tag、状态 Tag、上传时间）；测试按类型筛选；测试"上传文档"按钮打开弹窗
- **AC**: 3 个测试用例通过

### T59 依赖: T58 | 实现: KnowledgeListPage
- **文件**: `src/features/documents/components/KnowledgeListPage.tsx`（新建）
- **操作**: SearchTable + 文档类型/状态筛选 + UploadDocumentModal 触发
- **AC**: T58 测试全部通过

### T60 依赖: T57 | 测试: UploadDocumentModal
- **文件**: `src/features/documents/components/UploadDocumentModal.test.tsx`（新建）
- **操作**: 测试表单渲染（项目 ID、名称、类型选择）；测试名称为空校验失败；测试提交成功关闭弹窗
- **AC**: 3 个测试用例通过

### T61 依赖: T60 | 实现: UploadDocumentModal
- **文件**: `src/features/documents/components/UploadDocumentModal.tsx`（新建）
- **操作**: Arco Modal + Form（名称 2-255 字符 + 类型 Select）+ React Hook Form + zod
- **AC**: T60 测试全部通过

### T62 `[P]` 依赖: T57 | 测试: DocumentDetailPage
- **文件**: `src/features/documents/components/DocumentDetailPage.test.tsx`（新建）
- **操作**: 测试渲染文档基本信息；测试渲染分块列表；测试处理中状态显示 Spin
- **AC**: 3 个测试用例通过

### T63 依赖: T62 | 实现: DocumentDetailPage
- **文件**: `src/features/documents/components/DocumentDetailPage.tsx`（新建）
- **操作**: 文档信息卡片 + 分块列表（Arco List）
- **AC**: T62 测试全部通过

---

## Phase 5: AI 生成模块

### T64 依赖: T05, T07 | 测试: Generation API service
- **文件**: `src/features/generation/services/generation.test.ts`（新建）
- **操作**: 测试 createTask(); 测试 getTask(id); 测试 listTasks(params); 测试 getTaskDrafts(taskId)
- **AC**: 4 个测试用例通过

### T65 依赖: T64 | 实现: Generation API service
- **文件**: `src/features/generation/services/generation.ts`（新建）
- **操作**: 实现 generationApi.createTask/getTask/listTasks/getTaskDrafts
- **AC**: T64 测试全部通过

### T66 依赖: T65 | 测试: useGeneration hooks
- **文件**: `src/features/generation/hooks/useGeneration.test.ts`（新建）
- **操作**: 测试 useGenerationTasks 列表查询；测试 useGenerationTask 详情查询；测试 useCreateGenerationTask mutation
- **AC**: 3 个测试用例通过

### T67 依赖: T66 | 实现: useGeneration hooks
- **文件**: `src/features/generation/hooks/useGeneration.ts`（新建）
- **操作**: 实现 useGenerationTasks, useGenerationTask, useCreateGenerationTask
- **AC**: T66 测试全部通过

### T68 依赖: T65 | 测试: usePollingTask
- **文件**: `src/features/generation/hooks/usePollingTask.test.ts`（新建）
- **操作**: 测试 pending 状态触发轮询（refetchInterval=3000）；测试 completed 状态停止轮询；测试 taskId 为空时不发起查询
- **AC**: 3 个测试用例通过

### T69 依赖: T68 | 实现: usePollingTask
- **文件**: `src/features/generation/hooks/usePollingTask.ts`（新建）
- **操作**: useQuery + refetchInterval 动态判断（pending/processing → 3000ms，其他 → false）
- **AC**: T68 测试全部通过

### T70 依赖: T67 | 测试: NewGenerationTaskPage
- **文件**: `src/features/generation/components/NewGenerationTaskPage.test.tsx`（新建）
- **操作**: 测试渲染模块选择（必填）；测试需求描述输入（≥10 字校验）；测试用例数量输入（1-20）；测试高级选项折叠/展开；测试提交成功跳转任务详情
- **AC**: 5 个测试用例通过

### T71 依赖: T70 | 实现: NewGenerationTaskPage
- **文件**: `src/features/generation/components/NewGenerationTaskPage.tsx`（新建）
- **操作**: Arco Form + 模块 Select + 需求描述 Textarea + 用例数量 InputNumber + Collapse 高级选项（场景类型/优先级/用例类型）+ React Hook Form + zod
- **AC**: T70 测试全部通过

### T72 依赖: T67, T69, T27 | 测试: GenerationTaskListPage
- **文件**: `src/features/generation/components/GenerationTaskListPage.test.tsx`（新建）
- **操作**: 测试渲染任务列表（prompt 摘要、状态 Tag、创建时间）；测试状态筛选（pending/processing/completed/failed）；测试点击任务跳转详情
- **AC**: 3 个测试用例通过

### T73 依赖: T72 | 实现: GenerationTaskListPage
- **文件**: `src/features/generation/components/GenerationTaskListPage.tsx`（新建）
- **操作**: SearchTable + 状态筛选 Select + StatusTag 渲染 + "新建任务"按钮
- **AC**: T72 测试全部通过

### T74 依赖: T69, T65 | 测试: TaskDetailPage
- **文件**: `src/features/generation/components/TaskDetailPage.test.tsx`（新建）
- **操作**: 测试渲染任务详情（prompt、状态、时间）；测试 processing 状态显示轮询进度；测试 completed 显示草稿列表；测试草稿列表含标题/类型/优先级/置信度
- **AC**: 4 个测试用例通过

### T75 依赖: T74 | 实现: TaskDetailPage
- **文件**: `src/features/generation/components/TaskDetailPage.tsx`（新建）
- **操作**: 任务信息卡片 + usePollingTask + 草稿列表（StatusTag 展示置信度）
- **AC**: T74 测试全部通过

---

## Phase 6: 草稿箱模块

### T76 依赖: T05, T07 | 测试: Drafts API service
- **文件**: `src/features/drafts/services/drafts.test.ts`（新建）
- **操作**: 测试 getDrafts(params)；测试 confirmDraft(draftId, moduleId)；测试 rejectDraft(draftId, data)；测试 batchConfirm(draftIds, moduleId)
- **AC**: 4 个测试用例通过

### T77 依赖: T76 | 实现: Drafts API service
- **文件**: `src/features/drafts/services/drafts.ts`（新建）
- **操作**: 实现 draftsApi.getDrafts/confirmDraft/rejectDraft/batchConfirm
- **AC**: T76 测试全部通过

### T78 依赖: T77 | 测试: useDrafts hooks
- **文件**: `src/features/drafts/hooks/useDrafts.test.ts`（新建）
- **操作**: 测试 useDraftList 查询；测试 useConfirmDraft mutation 成功后失效缓存；测试 useRejectDraft mutation；测试 useBatchConfirm mutation；测试 usePendingDraftCount 轮询
- **AC**: 5 个测试用例通过

### T79 依赖: T78 | 实现: useDrafts hooks
- **文件**: `src/features/drafts/hooks/useDrafts.ts`（新建）
- **操作**: 实现 useDraftList, useConfirmDraft, useRejectDraft, useBatchConfirm, usePendingDraftCount
- **AC**: T78 测试全部通过

### T80 依赖: T79, T27, T29 | 测试: DraftListPage
- **文件**: `src/features/drafts/components/DraftListPage.test.tsx`（新建）
- **操作**: 测试渲染草稿列表（标题、来源、置信度 Tag、时间）；测试批量勾选 checkbox；测试"批量确认"按钮；测试按项目/模块/状态筛选
- **AC**: 4 个测试用例通过

### T81 依赖: T80 | 实现: DraftListPage
- **文件**: `src/features/drafts/components/DraftListPage.tsx`（新建）
- **操作**: SearchTable + checkbox 批量选择 + 筛选条件 + 批量确认/拒绝按钮 + 点击跳转确认页
- **AC**: T80 测试全部通过

### T82 依赖: T79, T31, T35, T37 | 测试: DraftConfirmPage
- **文件**: `src/features/drafts/components/DraftConfirmPage.test.tsx`（新建）
- **操作**: 测试左右分栏布局；测试左侧编辑区（标题、前置条件 ArrayEditor、步骤 ArrayEditor）；测试右侧引用来源面板；测试"确认"按钮触发 mutation；测试"拒绝"按钮弹出原因选择
- **AC**: 5 个测试用例通过

### T83 依赖: T82 | 实现: DraftConfirmPage
- **文件**: `src/features/drafts/components/DraftConfirmPage.tsx`（新建）
- **操作**: SplitPanel + 左侧表单（标题 Input + 前置条件 ArrayEditor + 步骤 ArrayEditor + 预期结果 Textarea + 类型/优先级 Select）+ 右侧 ReferencePanel + 底部按钮组（确认/保存/拒绝）
- **AC**: T82 测试全部通过

---

## Phase 7: 测试用例管理

### T84 依赖: T05, T07 | 测试: TestCases API service
- **文件**: `src/features/testcases/services/testcases.test.ts`（新建）
- **操作**: 测试 list(params) 含筛选参数；测试 get(id) 返回 CaseDetail；测试 create(data)；测试 update(id, data)；测试 delete(id)
- **AC**: 5 个测试用例通过

### T85 依赖: T84 | 实现: TestCases API service
- **文件**: `src/features/testcases/services/testcases.ts`（新建）
- **操作**: 实现 testcasesApi.list/get/create/update/delete
- **AC**: T84 测试全部通过

### T86 依赖: T85 | 测试: useTestCases hooks
- **文件**: `src/features/testcases/hooks/useTestCases.test.ts`（新建）
- **操作**: 测试 useCaseList 按 project_id 查询；测试 useCaseDetail 查询；测试 useCreateTestCase mutation；测试 useUpdateTestCase mutation；测试 useDeleteTestCase mutation
- **AC**: 5 个测试用例通过

### T87 依赖: T86 | 实现: useTestCases hooks
- **文件**: `src/features/testcases/hooks/useTestCases.ts`（新建）
- **操作**: 实现 useCaseList, useCaseDetail, useCreateTestCase, useUpdateTestCase, useDeleteTestCase
- **AC**: T86 测试全部通过

### T88 依赖: T87, T27, T29 | 测试: CaseListPage
- **文件**: `src/features/testcases/components/CaseListPage.test.tsx`（新建）
- **操作**: 测试渲染用例表格（编号、标题、类型 Tag、优先级 Tag、状态 Tag）；测试筛选栏（类型/优先级/状态）；测试分页；测试点击行跳转详情
- **AC**: 4 个测试用例通过

### T89 依赖: T88 | 实现: CaseListPage
- **文件**: `src/features/testcases/components/CaseListPage.tsx`（新建）
- **操作**: 筛选栏（Arco Select × 3）+ SearchTable + "新建用例"按钮 + 行点击导航
- **AC**: T88 测试全部通过

### T90 依赖: T87, T27, T37 | 测试: CaseDetailPage
- **文件**: `src/features/testcases/components/CaseDetailPage.test.tsx`（新建）
- **操作**: 测试渲染用例完整信息（编号、标题、前置条件、步骤、预期结果）；测试 AI 元数据区域展示（置信度 Tag、引用块列表）；测试编辑/删除按钮存在
- **AC**: 3 个测试用例通过

### T91 依赖: T90 | 实现: CaseDetailPage
- **文件**: `src/features/testcases/components/CaseDetailPage.tsx`（新建）
- **操作**: 用例信息卡片 + 步骤列表 + AI 来源区域（StatusTag confidence + ReferencePanel）+ 操作按钮
- **AC**: T90 测试全部通过

### T92 依赖: T87, T31 | 测试: CreateCaseDrawer
- **文件**: `src/features/testcases/components/CreateCaseDrawer.test.tsx`（新建）
- **操作**: 测试 Drawer 打开时渲染表单；测试模块选择必填校验；测试标题必填校验；测试步骤至少 1 条校验；测试提交成功关闭 Drawer
- **AC**: 5 个测试用例通过

### T93 依赖: T92 | 实现: CreateCaseDrawer
- **文件**: `src/features/testcases/components/CreateCaseDrawer.tsx`（新建）
- **操作**: Arco Drawer + React Hook Form + zod（module_id 必填, title 2-500 字符, steps ≥1 条, expected 对象, case_type 必填, priority 必填）+ ArrayEditor（步骤）+ ArrayEditor（前置条件）
- **AC**: T92 测试全部通过

---

## Phase 8: 测试计划与执行

### T94 依赖: T05, T07 | 测试: Plans API service
- **文件**: `src/features/plans/services/plans.test.ts`（新建）
- **操作**: 测试 list(params); 测试 get(id) 返回 PlanDetail 含 stats; 测试 create(data); 测试 addCases(planId, caseIds); 测试 removeCase(planId, caseId); 测试 recordResult(planId, data)
- **AC**: 6 个测试用例通过

### T95 依赖: T94 | 实现: Plans API service
- **文件**: `src/features/plans/services/plans.ts`（新建）
- **操作**: 实现 plansApi.list/get/create/update/delete/addCases/removeCase/getResults/recordResult
- **AC**: T94 测试全部通过

### T96 依赖: T95 | 测试: usePlans hooks
- **文件**: `src/features/plans/hooks/usePlans.test.ts`（新建）
- **操作**: 测试 usePlanList; 测试 usePlanDetail; 测试 useCreatePlan mutation; 测试 useAddCases mutation; 测试 useRecordResult mutation
- **AC**: 5 个测试用例通过

### T97 依赖: T96 | 实现: usePlans hooks
- **文件**: `src/features/plans/hooks/usePlans.ts`（新建）
- **操作**: 实现 usePlanList, usePlanDetail, useCreatePlan, useUpdatePlan, useDeletePlan, useAddCases, useRemoveCase, useRecordResult
- **AC**: T96 测试全部通过

### T98 依赖: T97, T27, T29 | 测试: PlanListPage
- **文件**: `src/features/plans/components/PlanListPage.test.tsx`（新建）
- **操作**: 测试渲染计划列表（名称、状态 Tag、创建时间）；测试状态筛选；测试"新建计划"按钮
- **AC**: 3 个测试用例通过

### T99 依赖: T98 | 实现: PlanListPage
- **文件**: `src/features/plans/components/PlanListPage.tsx`（新建）
- **操作**: SearchTable + 状态筛选 + StatusTag + "新建计划"按钮
- **AC**: T98 测试全部通过

### T100 依赖: T97 | 测试: NewPlanPage
- **文件**: `src/features/plans/components/NewPlanPage.test.tsx`（新建）
- **操作**: 测试渲染计划名称输入；测试用例选择面板；测试名称为空校验失败；测试提交成功跳转计划详情
- **AC**: 4 个测试用例通过

### T101 依赖: T100 | 实现: NewPlanPage
- **文件**: `src/features/plans/components/NewPlanPage.tsx`（新建）
- **操作**: 计划名称 Input + 描述 Textarea + 用例选择（checkbox 列表）+ React Hook Form + zod
- **AC**: T100 测试全部通过

### T102 依赖: T97, T27 | 测试: PlanDetailPage
- **文件**: `src/features/plans/components/PlanDetailPage.test.tsx`（新建）
- **操作**: 测试渲染计划信息 + 执行统计（总数/通过/失败/阻塞/跳过）；测试用例列表渲染；测试"录入结果"按钮打开弹窗
- **AC**: 3 个测试用例通过

### T103 依赖: T102 | 实现: PlanDetailPage
- **文件**: `src/features/plans/components/PlanDetailPage.tsx`（新建）
- **操作**: 计划信息卡片 + 统计进度条 + 用例列表 Table + ResultRecordModal 触发
- **AC**: T102 测试全部通过

### T104 依赖: T97 | 测试: ResultRecordModal
- **文件**: `src/features/plans/components/ResultRecordModal.test.tsx`（新建）
- **操作**: 测试渲染状态选择（pass/fail/block/skip）；测试备注输入；测试提交成功关闭弹窗
- **AC**: 3 个测试用例通过

### T105 依赖: T104 | 实现: ResultRecordModal
- **文件**: `src/features/plans/components/ResultRecordModal.tsx`（新建）
- **操作**: Arco Modal + Radio 组（执行状态）+ Textarea（备注）+ React Hook Form + zod
- **AC**: T104 测试全部通过

---

## Phase 9: 全局布局与路由集成

### T106 依赖: T16 | 测试: useAppStore（精简版）
- **文件**: `src/store/useAppStore.test.ts`（新建）
- **操作**: 测试 sidebarCollapsed 默认值；测试 toggleSidebar 切换；测试初始值响应窗口宽度
- **AC**: 3 个测试用例通过

### T107 依赖: T106 | 实现: useAppStore（精简版）
- **文件**: `src/store/useAppStore.ts`（修改）
- **操作**: 清理 counter 样板代码，仅保留 sidebarCollapsed + toggleSidebar
- **AC**: T106 测试全部通过

### T108 `[P]` 依赖: T79 | 测试: Sidebar
- **文件**: `src/components/layout/Sidebar.test.tsx`（新建）
- **操作**: 测试渲染菜单项（项目列表、当前项目子菜单、草稿箱）；测试选中态样式；测试折叠/展开切换；测试草稿箱 Badge 显示 pendingDraftCount
- **AC**: 4 个测试用例通过

### T109 依赖: T108 | 实现: Sidebar
- **文件**: `src/components/layout/Sidebar.tsx`（新建）
- **操作**: Arco Menu + Lucide icons + useAppStore(sidebarCollapsed) + usePendingDraftCount(Badge) + 用户信息区（Avatar + Dropdown）
- **AC**: T108 测试全部通过

### T110 `[P]` 测试: Header
- **文件**: `src/components/layout/Header.test.tsx`（新建）
- **操作**: 测试渲染折叠按钮；测试面包屑导航；测试通知图标；测试用户下拉菜单（退出登录）
- **AC**: 4 个测试用例通过

### T111 依赖: T110 | 实现: Header
- **文件**: `src/components/layout/Header.tsx`（新建）
- **操作**: 折叠 Button + Arco Breadcrumb + Bell Badge + Avatar Dropdown
- **AC**: T110 测试全部通过

### T112 依赖: T109, T111 | 测试: AppLayout
- **文件**: `src/components/layout/AppLayout.test.tsx`（新建）
- **操作**: 测试渲染 Sidebar + Header + Content 三区布局；测试侧边栏折叠时宽度变为 64px
- **AC**: 2 个测试用例通过

### T113 依赖: T112 | 实现: AppLayout
- **文件**: `src/components/layout/AppLayout.tsx`（新建）
- **操作**: Arco Layout + Layout.Sider（220px/64px）+ Layout.Header（48px）+ Layout.Content（独立滚动，padding 24px）
- **AC**: T112 测试全部通过

### T114 `[P]` 实现: NotFoundPage
- **文件**: `src/components/NotFoundPage.tsx`（新建）
- **操作**: Arco Result 404 + 返回首页按钮
- **AC**: 渲染 404 页面

### T115 依赖: T24, T113, T114, T20, T22, T43, T47, T59, T63, T73, T75, T81, T83, T89, T91, T93, T99, T101, T103 | 实现: 路由配置
- **文件**: `src/router/index.tsx`（新建）
- **操作**: createBrowserRouter 定义所有路由，AuthLayout 包裹登录/注册，AppLayout + RouteGuard 包裹所有业务页面，lazy() 加载，NotFoundPage 兜底
- **AC**: 所有路由可访问，懒加载正常

---

## Phase 10: 配置管理模块

### T116 依赖: T05, T07 | 测试: Configs API service
- **文件**: `src/features/configs/services/configs.test.ts`（新建）
- **操作**: 测试 list(projectId); 测试 set(projectId, key, value); 测试 delete(projectId, key); 测试 import(projectId, configs); 测试 export(projectId)
- **AC**: 5 个测试用例通过

### T117 依赖: T116 | 实现: Configs API service
- **文件**: `src/features/configs/services/configs.ts`（新建）
- **操作**: 实现 configsApi.list/set/delete/import/export
- **AC**: T116 测试全部通过

### T118 依赖: T117 | 测试: useConfigs hooks
- **文件**: `src/features/configs/hooks/useConfigs.test.ts`（新建）
- **操作**: 测试 useConfigList; 测试 useSetConfig mutation; 测试 useDeleteConfig mutation; 测试 useImportConfigs mutation; 测试 useExportConfigs
- **AC**: 5 个测试用例通过

### T119 依赖: T118 | 实现: useConfigs hooks
- **文件**: `src/features/configs/hooks/useConfigs.ts`（新建）
- **操作**: 实现 useConfigList, useSetConfig, useDeleteConfig, useImportConfigs, useExportConfigs
- **AC**: T118 测试全部通过

### T120 依赖: T119, T29 | 测试: ConfigManagePage
- **文件**: `src/features/configs/components/ConfigManagePage.test.tsx`（新建）
- **操作**: 测试渲染配置表格（键名、值、描述、操作）；测试"新增配置"按钮；测试编辑弹窗；测试"导入/导出"按钮
- **AC**: 4 个测试用例通过

### T121 依赖: T120 | 实现: ConfigManagePage
- **文件**: `src/features/configs/components/ConfigManagePage.tsx`（新建）
- **操作**: 配置列表 Table + 新增/编辑 Modal（键名 Input + 值 JsonEditor + 描述 Input）+ 导入/导出按钮
- **AC**: T120 测试全部通过

---

## Phase 11: E2E 测试

### T122 依赖: T115 | E2E: 登录流程
- **文件**: `tests/e2e/auth.spec.ts`（新建）
- **操作**: 测试登录成功跳转 /projects；测试登录失败显示错误；测试注册成功跳转 /login；测试未认证访问自动跳转 /login
- **AC**: 4 个 E2E 测试通过

### T123 依赖: T115 | E2E: 项目管理流程
- **文件**: `tests/e2e/project.spec.ts`（新建）
- **操作**: 测试创建项目 → 查看列表 → 进入仪表盘 → 创建模块 → 删除模块
- **AC**: E2E 测试通过

### T124 依赖: T115 | E2E: AI 生成核心流程
- **文件**: `tests/e2e/core-flow.spec.ts`（新建）
- **操作**: 测试登录 → 创建项目 → 上传文档 → 发起 AI 生成 → 查看草稿 → 确认草稿 → 验证用例编号格式 → 创建测试计划 → 录入执行结果
- **AC**: E2E 测试通过（覆盖 spec.md TC-01~TC-08 关键路径）

---

## 依赖关系图

```
Phase 0 (T01-T12)
  │
  ├─→ Phase 1: Auth (T13-T25)
  │     │
  │     └─→ Phase 9: Layout (T106-T115)
  │              │
  │              └─→ Phase 11: E2E (T122-T124)
  │
  ├─→ Phase 2: Shared Components (T26-T37) [P]
  │     │
  │     ├─→ Phase 3: Projects (T38-T53)
  │     │
  │     ├─→ Phase 4: Documents (T54-T63)
  │     │
  │     ├─→ Phase 5: Generation (T64-T75)
  │     │     │
  │     │     └─→ Phase 6: Drafts (T76-T83)
  │     │
  │     ├─→ Phase 7: TestCases (T84-T93)
  │     │
  │     └─→ Phase 8: Plans (T94-T105)
  │
  └─→ Phase 10: Configs (T116-T121) [独立，可与 Phase 3-8 并行]
```

## 统计

| 指标 | 数值 |
|------|------|
| 总任务数 | 124 |
| 测试任务 | 62（50%） |
| 实现任务 | 56 |
| 基础设施任务 | 6 |
| 可并行任务（`[P]`） | 18 |
| Phase 数 | 12 |
