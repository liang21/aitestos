---
title: 核心功能技术实现方案
version: 1.0.0
author: 首席前端架构师
status: draft
created: 2026-04-16
reviewers: [前端团队]
---

# 核心功能技术实现方案

> 基于 `spec.md`、OpenAPI 3.0.3、UX 设计规范 v1.0、前端详细设计 v1.1 编制。
> 本方案可直接驱动 TDD 开发。

---

## 1. 架构总览

```
┌─────────────────────────────────────────────────────────────┐
│                       React 19 + TypeScript                  │
├──────────┬──────────┬──────────┬──────────┬─────────────────┤
│  Pages   │ Features │  Shared  │  Stores  │    Queries      │
│ (Route)  │ (Domain) │   (UI)   │ (Zustand)│ (React Query)   │
├──────────┴──────────┴──────────┴──────────┴─────────────────┤
│                    Axios (request instance)                   │
│                    Token 刷新 / 错误拦截                      │
├─────────────────────────────────────────────────────────────┤
│                    React Router v7                            │
│                    Arco Design v2.66                          │
│                    Tailwind CSS v4                            │
└─────────────────────────────────────────────────────────────┘
```

**数据流方向**：`Page → Query/Hook → API function → Axios → Backend`

**状态分层原则**：

| 数据类型       | 管理方式                   | 示例                         |
| -------------- | -------------------------- | ---------------------------- |
| 服务端数据     | React Query                | 项目列表、用例详情、草稿列表 |
| 客户端全局状态 | Zustand                    | auth token、侧边栏折叠、通知 |
| 客户端局部状态 | useState/useReducer        | 表单输入、弹窗开关、筛选条件 |
| URL 状态       | React Router params/search | 当前项目 ID、分页页码        |

---

## 2. 技术上下文

### 2.1 技术栈选型

| 技术                 | 版本   | 职责        | 选型理由                                     |
| -------------------- | ------ | ----------- | -------------------------------------------- |
| React                | 19.2   | UI 框架     | 项目已安装，生态成熟                         |
| TypeScript           | 5.9    | 类型系统    | strict: true，编译期捕获错误                 |
| React Router         | 7.x    | 路由        | 项目已安装，支持 lazy loading                |
| Zustand              | 5.x    | 全局状态    | 项目已安装，极简 API                         |
| TanStack React Query | 5.x    | 服务端状态  | 缓存/重试/轮询/失效一体化                    |
| React Hook Form      | 7.x    | 表单管理    | 非受控性能优势，与 Arco 集成简洁             |
| Arco Design          | 2.66   | UI 组件库   | 项目已安装，企业级组件完备                   |
| Tailwind CSS         | 4.x    | 样式        | 项目已安装，原子化样式                       |
| Axios                | 1.14   | HTTP 客户端 | 项目已安装，作为 React Query 的 queryFn 底层 |
| dayjs                | 1.11   | 日期处理    | 项目已安装，轻量                             |
| Lucide React         | 1.7    | 图标        | 项目已安装，一致性好                         |
| Testing Library      | latest | 组件测试    | React 社区标准                               |
| MSW                  | latest | API Mock    | 浏览器层拦截，与 React Query 配合好          |
| Vitest               | latest | 测试运行器  | 与 Vite 原生集成                             |
| Playwright           | latest | E2E 测试    | 关键路径端到端验证                           |

### 2.2 简单性原则验证

| 检查项               | 结果 | 说明                                                  |
| -------------------- | ---- | ----------------------------------------------------- |
| 无多余抽象层         | ✅   | API function → React Query hook → Component，三层即可 |
| 状态管理最小化       | ✅   | 仅 auth + UI 用 Zustand，所有服务端数据走 React Query |
| 无重复请求           | ✅   | React Query 自动去重和缓存                            |
| 无手写 loading/error | ✅   | React Query 提供 isLoading/isError                    |
| 表单不手动管理       | ✅   | React Hook Form 管理 register/validate/submit         |

### 2.3 可测试性验证

| 测试层级 | 工具                     | 覆盖目标                              |
| -------- | ------------------------ | ------------------------------------- |
| 单元测试 | Vitest + Testing Library | 组件渲染、用户交互、Hook 逻辑         |
| API Mock | MSW                      | 拦截 HTTP 请求，模拟成功/失败/超时    |
| 集成测试 | Testing Library + MSW    | 完整的用户操作流程                    |
| E2E      | Playwright               | 登录 → 创建项目 → 发起生成 → 确认草稿 |

---

## 3. 合宪性审查

> `constitution.md` 当前为空文件，以下基于 Prompt 中声明的约束逐条审查。

### Constitutional Audit

| 规则                         | 状态    | 说明                                                                                                                       |
| ---------------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------- |
| **禁止直接 fetch/axios**     | ✅ 通过 | 组件中禁止直接调用 `request.get()`。所有数据获取通过 React Query hook 包装。`src/lib/request.ts` 仅作为 queryFn 底层传输层 |
| **禁止 any**                 | ✅ 通过 | API function 使用泛型 `request.post<never, ResponseType>()` 替代 `any`。类型定义在 `src/types/api.ts` 中集中管理           |
| **避免过度抽象**             | ✅ 通过 | 不创建 HOC、render props、context wrapper 等中间层。三层结构：Component → Hook → API function                              |
| **状态本地化**               | ✅ 通过 | 表单状态归 React Hook Form，服务端状态归 React Query，仅 token/sidebar 归 Zustand                                          |
| **Hook 提取规则**            | ✅ 通过 | 每个自定义 Hook 单一职责：`useProjects()` 管理项目列表查询，`useCreateProject()` 管理创建 mutation                         |
| **TDD First**                | ✅ 通过 | 每个实现任务先写测试，验收标准包含测试通过条件                                                                             |
| **禁止 useEffect 数据获取**  | ✅ 通过 | 所有数据获取使用 `useQuery()`，轮询使用 `refetchInterval`，不手写 `useEffect + fetch`                                      |
| **React Query 统一数据获取** | ✅ 通过 | 列表、详情、统计全部使用 useQuery；创建、更新、删除全部使用 useMutation                                                    |

### ⚠️ 现有代码修正项

| 现有文件                   | 问题                                                          | 修正方案                                                                          |
| -------------------------- | ------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| `src/lib/request.ts`       | 组件直接调用 axios，绕过 React Query 缓存                     | 保留为 React Query 的 queryFn 底层，组件层不再直接 import                         |
| `src/store/useAppStore.ts` | 包含样板 counter 代码，且 `fetchPendingDraftCount` 直接调 API | 清理为仅管理 sidebarCollapsed 和 notifications，待处理草稿数通过 React Query 查询 |
| `src/App.tsx`              | Vite 脚手架模板                                               | 完全重写为路由入口                                                                |

---

## 4. 项目结构

```
src/
├── app/                          # 应用入口
│   ├── App.tsx                   # 根组件（Provider 嵌套）
│   ├── main.tsx                  # 渲染入口
│   └── providers.tsx             # QueryClientProvider + ConfigProvider
│
├── router/                       # 路由
│   ├── index.tsx                 # 路由定义（lazy loading）
│   └── RouteGuard.tsx            # 认证/权限守卫
│
├── lib/                          # 基础设施
│   ├── request.ts                # Axios 实例（Token 刷新、错误拦截）
│   ├── query-client.ts           # React Query 全局配置
│   └── utils.ts                  # cn() 等工具函数
│
├── types/                        # 全局类型
│   ├── enums.ts                  # 枚举/字面量联合类型
│   └── api.ts                    # API 请求/响应类型
│
├── features/                     # 业务功能模块（Feature-Based）
│   ├── auth/                     # 认证
│   │   ├── components/
│   │   │   ├── LoginPage.tsx
│   │   │   └── RegisterPage.tsx
│   │   ├── hooks/
│   │   │   ├── useAuth.ts        # useLogin, useRegister, useRefresh
│   │   │   └── useAuthStore.ts   # Zustand store（token/user）
│   │   ├── services/
│   │   │   └── auth.ts           # API function
│   │   └── types.ts              # 认证相关类型（如果与全局不同）
│   │
│   ├── projects/                 # 项目管理
│   │   ├── components/
│   │   │   ├── ProjectListPage.tsx
│   │   │   ├── ProjectDashboard.tsx
│   │   │   └── CreateProjectModal.tsx
│   │   ├── hooks/
│   │   │   ├── useProjects.ts    # useProjectList, useProjectDetail, useProjectStats
│   │   │   └── useProjectMutations.ts  # useCreateProject, useUpdateProject, useDeleteProject
│   │   └── services/
│   │       └── projects.ts
│   │
│   ├── modules/                  # 模块管理
│   │   ├── components/
│   │   │   ├── ModuleManagePage.tsx
│   │   │   └── CreateModuleModal.tsx
│   │   ├── hooks/
│   │   │   └── useModules.ts
│   │   └── services/
│   │       └── modules.ts
│   │
│   ├── testcases/                # 测试用例
│   │   ├── components/
│   │   │   ├── CaseListPage.tsx
│   │   │   ├── CaseDetailPage.tsx
│   │   │   └── CreateCaseDrawer.tsx
│   │   ├── hooks/
│   │   │   └── useTestCases.ts
│   │   └── services/
│   │       └── testcases.ts
│   │
│   ├── plans/                    # 测试计划
│   │   ├── components/
│   │   │   ├── PlanListPage.tsx
│   │   │   ├── NewPlanPage.tsx
│   │   │   ├── PlanDetailPage.tsx
│   │   │   └── ResultRecordModal.tsx
│   │   ├── hooks/
│   │   │   └── usePlans.ts
│   │   └── services/
│   │       └── plans.ts
│   │
│   ├── generation/               # AI 生成
│   │   ├── components/
│   │   │   ├── GenerationTaskListPage.tsx
│   │   │   ├── NewGenerationTaskPage.tsx
│   │   │   └── TaskDetailPage.tsx
│   │   ├── hooks/
│   │   │   ├── useGeneration.ts  # useCreateTask, useTask, useTaskDrafts
│   │   │   └── usePollingTask.ts # 轮询任务状态
│   │   └── services/
│   │       └── generation.ts
│   │
│   ├── drafts/                   # 草稿箱
│   │   ├── components/
│   │   │   ├── DraftListPage.tsx
│   │   │   └── DraftConfirmPage.tsx
│   │   ├── hooks/
│   │   │   └── useDrafts.ts      # useDraftList, useConfirmDraft, useRejectDraft, useBatchConfirm
│   │   └── services/
│   │       └── drafts.ts
│   │
│   ├── documents/                # 知识库
│   │   ├── components/
│   │   │   ├── KnowledgeListPage.tsx
│   │   │   ├── DocumentDetailPage.tsx
│   │   │   └── UploadDocumentModal.tsx
│   │   ├── hooks/
│   │   │   └── useDocuments.ts
│   │   └── services/
│   │       └── documents.ts
│   │
│   └── configs/                  # 项目配置
│       ├── components/
│       │   ├── ConfigManagePage.tsx
│       │   └── ConfigEditModal.tsx
│       ├── hooks/
│       │   └── useConfigs.ts
│       └── services/
│           └── configs.ts
│
├── components/                   # 跨 Feature 共享组件
│   ├── layout/
│   │   ├── AppLayout.tsx         # 主布局（Sidebar + Header + Content）
│   │   ├── Sidebar.tsx           # 侧边栏
│   │   ├── Header.tsx            # 顶部栏
│   │   └── AuthLayout.tsx        # 登录/注册布局
│   ├── business/
│   │   ├── StatusTag.tsx         # 统一状态标签（色彩映射）
│   │   ├── SearchTable.tsx       # 搜索筛选表格
│   │   ├── ArrayEditor.tsx       # 数组编辑器（前置条件/步骤）
│   │   ├── StatsCard.tsx         # 统计卡片
│   │   ├── SplitPanel.tsx        # 分栏面板（草稿确认页）
│   │   └── ReferencePanel.tsx    # 引用来源面板
│   └── NotFoundPage.tsx          # 404 页面
│
├── hooks/                        # 跨 Feature 共享 Hook
│   └── usePagination.ts          # 分页逻辑（与 React Query 集成）
│
├── store/                        # 全局 Zustand store
│   └── useAppStore.ts            # sidebarCollapsed + notifications
│
├── styles/
│   ├── theme.css                 # Arco 主题变量 + Tailwind @theme
│   └── global.css                # 全局样式
│
├── index.css
└── vite-env.d.ts
```

### 目录职责与依赖规则

```
Page Component
  ├── 可引用 → 同 Feature 内的 hooks/
  ├── 可引用 → 同 Feature 内的 components/
  ├── 可引用 → @/components/ (共享组件)
  ├── 可引用 → @/hooks/ (共享 Hook)
  │
  hooks/ (React Query hooks)
  └── 可引用 → 同 Feature 内的 services/
      │
      services/ (API functions)
      └── 可引用 → @/lib/request.ts
      └── 可引用 → @/types/api.ts
```

**禁止**：

- Feature A 引用 Feature B 的内部文件
- components/ 引用 features/ 的任何内容
- services/ 引用 store/ 或 hooks/

---

## 5. 数据模型

### 5.1 枚举类型

```typescript
// src/types/enums.ts
export type CaseStatus = 'unexecuted' | 'pass' | 'block' | 'fail'
export type CaseType =
  | 'functionality'
  | 'performance'
  | 'api'
  | 'ui'
  | 'security'
export type PlanStatus = 'draft' | 'active' | 'completed' | 'archived'
export type Priority = 'P0' | 'P1' | 'P2' | 'P3'
export type ResultStatus = 'pass' | 'fail' | 'block' | 'skip'
export type TaskStatus = 'pending' | 'processing' | 'completed' | 'failed'
export type DraftStatus = 'pending' | 'confirmed' | 'rejected'
export type DocumentType = 'prd' | 'figma' | 'api_spec' | 'swagger' | 'markdown'
export type DocumentStatus = 'pending' | 'processing' | 'completed' | 'failed'
export type UserRole = 'super_admin' | 'admin' | 'normal'
export type Confidence = 'high' | 'medium' | 'low'
export type SceneType = 'positive' | 'negative' | 'boundary'
```

### 5.2 核心 API 类型

> 完整类型定义见 `specs/frontend-detailed-design.md` 第 4.2 节，此处列出核心模型。

```typescript
// src/types/api.ts — 关键结构摘要

interface PaginatedResponse<T> {
  data: T[]
  total: number
  offset: number
  limit: number
}

interface UserJSON {
  id: string
  username: string
  email: string
  role: UserRole
  createdAt: string
  updatedAt: string
}

interface Project {
  id: string
  name: string
  prefix: string
  description: string
  createdAt: string
  updatedAt: string
}

interface TestCase {
  id: string
  moduleId: string
  userId: string
  number: string // 格式: ECO-USR-20260416-001
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
}

interface AiMetadata {
  generationTaskId: string
  confidence: Confidence
  referencedChunks: ReferencedChunk[]
  modelVersion: string
  generatedAt: string
}

interface CaseDraft {
  id: string
  taskId: string
  title: string
  preconditions: string[]
  steps: string[]
  expected: Record<string, unknown>
  caseType: CaseType
  priority: Priority
  status: DraftStatus
  feedback?: string
  createdAt: string
  updatedAt: string
}

interface TestPlan {
  id: string
  projectId: string
  name: string
  description: string
  status: PlanStatus
  createdBy: string
  createdAt: string
  updatedAt: string
}

interface GenerationTask {
  id: string
  projectId: string
  moduleId: string
  status: TaskStatus
  prompt: string
  result: Record<string, unknown>
  createdAt: string
  updatedAt: string
}
```

### 5.3 禁止 any 的 Axios 调用规范

```typescript
// ❌ 禁止
request.post<any, LoginResponse>('/auth/login', data)

// ✅ 正确
request.post<never, LoginResponse>('/auth/login', data)
```

全局替换：所有 API function 中 `<any, T>` 改为 `<never, T>`。

---

## 6. 状态与数据流

### 6.1 React Query 使用规范

```typescript
// ✅ 正确的数据获取模式
export function useProjectList(params?: { keywords?: string }) {
  return useQuery({
    queryKey: ['projects', params],
    queryFn: () => projectsApi.list(params),
  })
}

export function useProjectDetail(id: string) {
  return useQuery({
    queryKey: ['projects', id],
    queryFn: () => projectsApi.get(id),
    enabled: !!id,
  })
}

// ✅ 正确的 Mutation 模式
export function useCreateProject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: projectsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}
```

### 6.2 轮询策略

```typescript
// AI 生成任务轮询
export function usePollingTask(taskId: string) {
  return useQuery({
    queryKey: ['generation-task', taskId],
    queryFn: () => generationApi.getTask(taskId),
    refetchInterval: (query) => {
      const status = query.state.data?.status
      return status === 'pending' || status === 'processing' ? 3000 : false
    },
    enabled: !!taskId,
  })
}

// 待处理草稿数量轮询（侧边栏 Badge）
export function usePendingDraftCount() {
  return useQuery({
    queryKey: ['drafts', 'pending-count'],
    queryFn: async () => {
      const res = await generationApi.getDrafts({ status: 'pending', limit: 1 })
      return res.total
    },
    refetchInterval: 60_000, // 每 60s 刷新
  })
}
```

### 6.3 Zustand 使用范围

仅以下状态使用 Zustand：

```typescript
// src/features/auth/hooks/useAuthStore.ts
interface AuthState {
  user: UserJSON | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  refresh: () => Promise<void>
}

// src/store/useAppStore.ts
interface AppState {
  sidebarCollapsed: boolean
  toggleSidebar: () => void
}
```

**不使用 Zustand 管理的数据**：项目列表、用例列表、计划详情、草稿列表等所有服务端数据。

---

## 7. 组件设计

### 7.1 组件分层

| 层级                  | 职责                         | 数据来源                  | 示例                       |
| --------------------- | ---------------------------- | ------------------------- | -------------------------- |
| **Page**              | 路由对应、布局组合、数据编排 | React Query hooks         | `CaseListPage`             |
| **Feature Component** | 业务逻辑、用户交互           | Props + React Query hooks | `CreateCaseDrawer`         |
| **Shared Component**  | 可复用 UI，无业务逻辑        | Props only                | `StatusTag`、`SearchTable` |

### 7.2 核心 Page 组件

```typescript
// Page 组件职责：组合 hooks + 布局，不包含复杂逻辑
function CaseListPage() {
  const { projectId } = useParams()
  const [filters, setFilters] = useState<CaseFilters>({})
  const { data, isLoading } = useCaseList(projectId, filters)

  return (
    <div className="p-6">
      <PageHeader title="测试用例" action={<CreateCaseButton />} />
      <CaseFilterBar filters={filters} onChange={setFilters} />
      <SearchTable
        loading={isLoading}
        data={data?.data ?? []}
        total={data?.total ?? 0}
        columns={caseColumns}
      />
    </div>
  )
}
```

### 7.3 核心 Feature Component

```typescript
// Feature Component：拥有自己的 React Query mutation
function CreateCaseDrawer({ visible, onClose }: CreateCaseDrawerProps) {
  const { projectId } = useParams()
  const form = useForm<CreateTestCaseFormData>({
    resolver: zodResolver(createCaseSchema),
    defaultValues: { steps: [''] },
  })
  const createCase = useCreateTestCase()

  const onSubmit = (data: CreateTestCaseFormData) => {
    createCase.mutate(data, { onSuccess: onClose })
  }

  return (
    <Drawer visible={visible} onClose={onClose} title="新建用例">
      <form onSubmit={form.handleSubmit(onSubmit)}>
        {/* Arco Form.Item + React Hook Form register */}
      </form>
    </Drawer>
  )
}
```

### 7.4 共享 Business Component

```typescript
// StatusTag：纯 UI，通过 props 驱动
interface StatusTagProps {
  status: string
  category: StatusCategory
  label?: string
  size?: 'small' | 'default'
}

function StatusTag({ status, category, label, size }: StatusTagProps) {
  const mapping = COLOR_MAP[category]?.[status]
  if (!mapping) return null

  return (
    <Tag
      size={size ?? 'small'}
      className={cn(
        'border-transparent',
        category === 'confidence' && 'font-medium'
      )}
      style={{ color: mapping.color, backgroundColor: mapping.bg }}
    >
      {label ?? mapping.text}
    </Tag>
  )
}
```

---

## 8. 副作用设计

### 8.1 数据获取（全部使用 React Query）

```typescript
// ✅ 列表查询
const { data, isLoading, error } = useProjectList({ keywords })

// ✅ 详情查询（依赖 ID）
const { data } = useProjectDetail(projectId)

// ✅ 条件查询（ID 存在时才发起）
const { data } = useProjectStats(projectId)

// ✅ 轮询（AI 任务状态）
const { data } = usePollingTask(taskId)
```

### 8.2 数据变更（全部使用 useMutation）

```typescript
// ✅ 创建
const createProject = useCreateProject()
createProject.mutate({ name: 'ECommerce', prefix: 'ECO' })

// ✅ 更新（自动失效缓存）
const updateCase = useUpdateTestCase()
updateCase.mutate({ id: caseId, data: { title: '新标题' } })

// ✅ 删除（乐观更新）
const deleteProject = useDeleteProject()
deleteProject.mutate(projectId)

// ✅ 确认草稿
const confirmDraft = useConfirmDraft()
confirmDraft.mutate({ draftId: id, moduleId: 'usr-module-id' })
```

### 8.3 Token 刷新流程

```
Axios 拦截器检测 401
  → 并发请求入队
  → 用 refresh_token 调用 /auth/refresh
  → 成功：更新 localStorage，重放队列中的请求
  → 失败：清理 token，触发 onAuthExpired → 跳转 /login
```

**注意**：Token 刷新逻辑保留在 Axios 拦截器中（`src/lib/request.ts`），这是底层传输层职责，不通过 React Query 管理。

### 8.4 禁止模式

```typescript
// ❌ 禁止：useEffect 中获取数据
useEffect(() => {
  fetchProjects()
}, [])

// ❌ 禁止：手动管理 loading 状态
const [loading, setLoading] = useState(false)

// ❌ 禁止：组件中直接调用 axios
request.get('/projects')

// ❌ 禁止：在 useEffect 中轮询
useEffect(() => {
  const timer = setInterval(() => refetch(), 3000)
  return () => clearInterval(timer)
}, [])
```

---

## 9. TDD 计划

### 9.1 测试架构

```
tests/
├── setup.ts                      # 测试入口（beforeEach cleanup）
├── msw/
│   ├── handlers/
│   │   ├── auth.ts               # 认证接口 mock
│   │   ├── projects.ts           # 项目接口 mock
│   │   ├── testcases.ts          # 用例接口 mock
│   │   ├── generation.ts         # 生成接口 mock
│   │   └── drafts.ts             # 草稿接口 mock
│   ├── server.ts                 # MSW setupServer
│   └── browser.ts                # MSW setupWorker（开发调试用）
│
└── e2e/
    └── core-flow.spec.ts         # Playwright 关键路径
```

### 9.2 单元/UI 测试矩阵

| 测试目标               | 测试内容                                 | Mock 策略                                        |
| ---------------------- | ---------------------------------------- | ------------------------------------------------ |
| **LoginPage**          | 渲染邮箱/密码输入框和登录按钮            | 无需 Mock                                        |
| **LoginPage**          | 输入无效邮箱时显示验证错误               | 无需 Mock                                        |
| **LoginPage**          | 提交后调用 login API，成功跳转 /projects | MSW mock `POST /auth/login` 200                  |
| **LoginPage**          | 登录失败显示错误提示                     | MSW mock `POST /auth/login` 401                  |
| **useProjectList**     | 返回项目列表数据                         | MSW mock `GET /projects`                         |
| **useProjectList**     | 支持关键词搜索                           | MSW mock `GET /projects?keywords=xxx`            |
| **useCreateProject**   | 创建成功后失效项目列表缓存               | MSW mock `POST /projects` 201                    |
| **ProjectListPage**    | 渲染项目表格，含搜索和分页               | MSW mock `GET /projects`                         |
| **CreateProjectModal** | 提交表单后关闭弹窗并刷新列表             | MSW mock `POST /projects` 201                    |
| **useCaseList**        | 返回筛选后的用例列表                     | MSW mock `GET /testcases?project_id=xxx`         |
| **CaseDetailPage**     | 显示用例详情和 AI 来源追溯               | MSW mock `GET /testcases/:id`                    |
| **usePollingTask**     | pending 时每 3s 轮询                     | MSW mock 序列：pending → processing → completed  |
| **DraftConfirmPage**   | 左右分栏布局，编辑+引用来源              | MSW mock `GET /generation/drafts/:id`            |
| **useConfirmDraft**    | 确认后草稿状态变为 confirmed             | MSW mock `POST /generation/drafts/:id/confirm`   |
| **useBatchConfirm**    | 批量确认返回成功/失败计数                | MSW mock `POST /generation/drafts/batch-confirm` |
| **StatusTag**          | 根据 category+status 渲染正确颜色        | 无需 Mock（纯 UI）                               |
| **SearchTable**        | 渲染分页表格，支持排序                   | 无需 Mock（纯 UI，传 props）                     |
| **ArrayEditor**        | 添加/删除/拖拽排序行                     | 无需 Mock（纯 UI）                               |

### 9.3 测试示例

```typescript
// tests/features/auth/LoginPage.test.tsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { server } from '../../msw/server'
import { LoginPage } from '@/features/auth/components/LoginPage'

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        {ui}
      </BrowserRouter>
    </QueryClientProvider>
  )
}

describe('LoginPage', () => {
  it('should show validation error for invalid email', async () => {
    renderWithProviders(<LoginPage />)
    const user = userEvent.setup()

    await user.type(screen.getByLabelText('邮箱'), 'invalid-email')
    await user.click(screen.getByRole('button', { name: /登录/ }))

    expect(screen.getByText(/请输入有效的邮箱地址/)).toBeInTheDocument()
  })

  it('should redirect to /projects on successful login', async () => {
    server.use(
      http.post('/api/v1/auth/login', () =>
        HttpResponse.json({
          access_token: 'mock-token',
          refresh_token: 'mock-refresh',
          user: { id: '1', username: 'admin', email: 'admin@test.com', role: 'admin' },
        })
      )
    )

    renderWithProviders(<LoginPage />)
    const user = userEvent.setup()

    await user.type(screen.getByLabelText('邮箱'), 'admin@test.com')
    await user.type(screen.getByLabelText('密码'), 'Test1234')
    await user.click(screen.getByRole('button', { name: /登录/ }))

    expect(window.location.pathname).toBe('/projects')
  })

  it('should show error message on login failure', async () => {
    server.use(
      http.post('/api/v1/auth/login', () =>
        HttpResponse.json({ error: '邮箱或密码错误' }, { status: 401 })
      )
    )

    renderWithProviders(<LoginPage />)
    const user = userEvent.setup()

    await user.type(screen.getByLabelText('邮箱'), 'admin@test.com')
    await user.type(screen.getByLabelText('密码'), 'wrong')
    await user.click(screen.getByRole('button', { name: /登录/ }))

    expect(screen.getByText('邮箱或密码错误')).toBeInTheDocument()
  })
})
```

### 9.4 E2E 关键路径

```typescript
// tests/e2e/core-flow.spec.ts (Playwright)
test.describe('AI 生成核心流程', () => {
  test('登录 → 创建项目 → 上传文档 → 发起生成 → 确认草稿', async ({ page }) => {
    // 1. 登录
    await page.goto('/login')
    await page.fill('[aria-label="邮箱"]', 'admin@aitestos.com')
    await page.fill('[aria-label="密码"]', 'Test1234')
    await page.click('button:has-text("登录")')
    await expect(page).toHaveURL('/projects')

    // 2. 创建项目
    await page.click('button:has-text("新建项目")')
    await page.fill('[placeholder="项目名称"]', 'ECommerce')
    await page.fill('[placeholder="项目前缀"]', 'ECO')
    await page.click('button:has-text("确定")')
    await expect(page.locator('text=ECommerce')).toBeVisible()

    // 3. 进入项目 → 创建模块
    await page.click('text=ECommerce')
    await page.click('text=项目设置')
    await page.click('text=模块管理')
    // ... 后续步骤
  })
})
```

---

## 10. 性能与可访问性

### 10.1 性能

| 优化项                    | 策略                                               | 适用场景                               |
| ------------------------- | -------------------------------------------------- | -------------------------------------- |
| **路由懒加载**            | `React.lazy()` + `Suspense`                        | 所有 Page 组件                         |
| **列表 key**              | 使用 `item.id`（UUID）                             | 所有 Table/List 渲染                   |
| **虚拟滚动**              | 暂不引入，1000 条用例 < 2s 满足要求                | 若后续超阈值再引入 `@tanstack/virtual` |
| **React.memo**            | 仅在 `StatusTag`、`SearchTable` 等高频渲染组件使用 | 避免不必要的重渲染                     |
| **React Query staleTime** | 项目列表 5min，详情 2min，统计数据 1min            | 减少重复请求                           |
| **图片**                  | 使用 `loading="lazy"`                              | 知识库文档缩略图                       |
| **Bundle 分析**           | `vite-plugin-visualizer`                           | 构建后检查包体积                       |

### 10.2 可访问性

| 要求              | 实现                                           | 验证方式                   |
| ----------------- | ---------------------------------------------- | -------------------------- |
| **语义化 HTML**   | 使用 `<nav>`, `<main>`, `<article>`, `<aside>` | axe audit                  |
| **aria-label**    | 所有按钮、链接、输入框添加描述性 label         | 人工 Review                |
| **键盘导航**      | 侧边栏 `Tab` + `Enter`，表格 `Arrow` 键        | 手动测试                   |
| **色彩对比度**    | 所有文字/背景组合 ≥ 4.5:1 (WCAG AA)            | UX 规范已验证              |
| **Focus 管理**    | Modal 打开时 trap focus，关闭时恢复            | Testing Library `tab` 测试 |
| **Screen Reader** | StatusTag 添加 `aria-label="状态：通过"`       | axe audit                  |

---

## 11. 任务清单

### Phase 0: 基础设施（TDD 基座）

| 优先级 | 任务                                                                                    | 路径                                     | AC                                                 |
| ------ | --------------------------------------------------------------------------------------- | ---------------------------------------- | -------------------------------------------------- |
| P0     | 安装依赖：@tanstack/react-query, react-hook-form, @testing-library/\*, msw, vitest, zod | `package.json`                           | `npm install` 无错误；`npm test` 可运行            |
| P0     | 配置 Vitest                                                                             | `vitest.config.ts`                       | `npm test` 通过空测试                              |
| P0     | 配置 MSW                                                                                | `tests/msw/server.ts`                    | handler 可拦截请求                                 |
| P0     | 创建 React Query Client                                                                 | `src/lib/query-client.ts`                | 默认 staleTime 5min，retry 1 次                    |
| P0     | 创建 Provider 嵌套                                                                      | `src/app/providers.tsx`                  | QueryClientProvider + Arco ConfigProvider 包裹     |
| P0     | 类型定义                                                                                | `src/types/enums.ts`, `src/types/api.ts` | TypeScript 编译通过，无 any                        |
| P0     | 修正 Axios 调用                                                                         | `src/lib/request.ts`                     | 所有 `<any, T>` → `<never, T>`；Token 刷新逻辑完整 |

### Phase 1: 认证模块

| 优先级 | 任务                   | 路径                                            | AC                                             |
| ------ | ---------------------- | ----------------------------------------------- | ---------------------------------------------- |
| P0     | Auth API service       | `src/features/auth/services/auth.ts`            | 类型安全，编译通过                             |
| P0     | Auth Zustand store     | `src/features/auth/hooks/useAuthStore.ts`       | login/logout/refresh 正确操作 localStorage     |
| P0     | `useLogin` hook + 测试 | `src/features/auth/hooks/useAuth.ts`            | 测试：成功登录更新 store，失败抛错             |
| P0     | LoginPage + 测试       | `src/features/auth/components/LoginPage.tsx`    | 测试：渲染表单、校验邮箱、成功跳转、失败提示   |
| P0     | RegisterPage + 测试    | `src/features/auth/components/RegisterPage.tsx` | 测试：渲染表单、字段校验、提交成功             |
| P0     | RouteGuard + 测试      | `src/router/RouteGuard.tsx`                     | 测试：未认证跳转 /login，token 过期跳转 /login |
| P0     | AuthLayout             | `src/components/layout/AuthLayout.tsx`          | 居中布局，无侧边栏                             |

### Phase 2: 项目管理模块

| 优先级 | 任务                       | 路径                                                      | AC                                                         |
| ------ | -------------------------- | --------------------------------------------------------- | ---------------------------------------------------------- |
| P0     | Projects API service       | `src/features/projects/services/projects.ts`              | 类型安全，与 OpenAPI 对齐                                  |
| P0     | `useProjects` hooks + 测试 | `src/features/projects/hooks/useProjects.ts`              | 测试：列表查询、详情查询、创建/更新/删除 mutation          |
| P0     | ProjectListPage + 测试     | `src/features/projects/components/ProjectListPage.tsx`    | 测试：渲染表格、搜索、分页、新建弹窗                       |
| P0     | CreateProjectModal + 测试  | `src/features/projects/components/CreateProjectModal.tsx` | 测试：表单校验（名称/前缀必填，前缀 2-4 位大写）、提交成功 |
| P0     | ProjectDashboard + 测试    | `src/features/projects/components/ProjectDashboard.tsx`   | 测试：渲染统计卡片、趋势图、最近任务                       |
| P0     | StatsCard 组件 + 测试      | `src/components/business/StatsCard.tsx`                   | 测试：渲染标题、数值、趋势箭头                             |
| P0     | Module API + hooks + 测试  | `src/features/modules/`                                   | 测试：CRUD 操作、列表按项目查询                            |
| P0     | ModuleManagePage + 测试    | `src/features/modules/components/ModuleManagePage.tsx`    | 测试：模块列表、创建/删除                                  |

### Phase 3: 知识库模块

| 优先级 | 任务                        | 路径                                                        | AC                                    |
| ------ | --------------------------- | ----------------------------------------------------------- | ------------------------------------- |
| P0     | Documents API service       | `src/features/documents/services/documents.ts`              | 类型安全                              |
| P0     | `useDocuments` hooks + 测试 | `src/features/documents/hooks/useDocuments.ts`              | 测试：列表查询、上传、删除、详情+分块 |
| P0     | KnowledgeListPage + 测试    | `src/features/documents/components/KnowledgeListPage.tsx`   | 测试：文档列表、状态 Tag、上传弹窗    |
| P0     | UploadDocumentModal + 测试  | `src/features/documents/components/UploadDocumentModal.tsx` | 测试：文件类型选择、名称输入、提交    |
| P1     | DocumentDetailPage + 测试   | `src/features/documents/components/DocumentDetailPage.tsx`  | 测试：分块列表展示                    |

### Phase 4: AI 生成模块（核心）

| 优先级 | 任务                          | 路径                                                            | AC                                                    |
| ------ | ----------------------------- | --------------------------------------------------------------- | ----------------------------------------------------- |
| P0     | Generation API service        | `src/features/generation/services/generation.ts`                | 类型安全                                              |
| P0     | `useGeneration` hooks + 测试  | `src/features/generation/hooks/useGeneration.ts`                | 测试：创建任务、查询任务、查询草稿                    |
| P0     | `usePollingTask` hook + 测试  | `src/features/generation/hooks/usePollingTask.ts`               | 测试：pending 时轮询、completed 时停止                |
| P0     | NewGenerationTaskPage + 测试  | `src/features/generation/components/NewGenerationTaskPage.tsx`  | 测试：模块选择、需求描述（≥10字）、高级选项折叠、提交 |
| P0     | GenerationTaskListPage + 测试 | `src/features/generation/components/GenerationTaskListPage.tsx` | 测试：任务列表、状态 Tag、进度展示                    |
| P1     | TaskDetailPage + 测试         | `src/features/generation/components/TaskDetailPage.tsx`         | 测试：草稿列表、置信度标签                            |

### Phase 5: 草稿箱模块（核心）

| 优先级 | 任务                       | 路径                                                  | AC                                          |
| ------ | -------------------------- | ----------------------------------------------------- | ------------------------------------------- |
| P0     | Drafts API service         | `src/features/drafts/services/drafts.ts`              | 类型安全                                    |
| P0     | `useDrafts` hooks + 测试   | `src/features/drafts/hooks/useDrafts.ts`              | 测试：列表查询、确认、拒绝、批量确认        |
| P0     | DraftListPage + 测试       | `src/features/drafts/components/DraftListPage.tsx`    | 测试：草稿列表、批量勾选、确认/拒绝操作     |
| P0     | DraftConfirmPage + 测试    | `src/features/drafts/components/DraftConfirmPage.tsx` | 测试：左右分栏、编辑区、引用来源、确认/拒绝 |
| P0     | SplitPanel 组件 + 测试     | `src/components/business/SplitPanel.tsx`              | 测试：拖拽分割、最小宽度                    |
| P0     | ReferencePanel 组件 + 测试 | `src/components/business/ReferencePanel.tsx`          | 测试：引用块列表、查看原文                  |
| P0     | ArrayEditor 组件 + 测试    | `src/components/business/ArrayEditor.tsx`             | 测试：添加/删除行、拖拽排序                 |

### Phase 6: 测试用例管理

| 优先级 | 任务                        | 路径                                                     | AC                                               |
| ------ | --------------------------- | -------------------------------------------------------- | ------------------------------------------------ |
| P0     | TestCases API service       | `src/features/testcases/services/testcases.ts`           | 类型安全                                         |
| P0     | `useTestCases` hooks + 测试 | `src/features/testcases/hooks/useTestCases.ts`           | 测试：列表查询（含筛选）、详情、CRUD mutation    |
| P0     | CaseListPage + 测试         | `src/features/testcases/components/CaseListPage.tsx`     | 测试：表格渲染、筛选栏（类型/优先级/状态）、分页 |
| P0     | CaseDetailPage + 测试       | `src/features/testcases/components/CaseDetailPage.tsx`   | 测试：用例信息展示、AI 元数据展示、编号显示      |
| P0     | CreateCaseDrawer + 测试     | `src/features/testcases/components/CreateCaseDrawer.tsx` | 测试：表单校验、步骤编辑器、提交                 |
| P0     | StatusTag 组件 + 测试       | `src/components/business/StatusTag.tsx`                  | 测试：各 category+status 颜色映射正确            |
| P0     | SearchTable 组件 + 测试     | `src/components/business/SearchTable.tsx`                | 测试：渲染表格、分页、loading/error 状态         |

### Phase 7: 测试计划与执行

| 优先级 | 任务                     | 路径                                                  | AC                                              |
| ------ | ------------------------ | ----------------------------------------------------- | ----------------------------------------------- |
| P0     | Plans API service        | `src/features/plans/services/plans.ts`                | 类型安全                                        |
| P0     | `usePlans` hooks + 测试  | `src/features/plans/hooks/usePlans.ts`                | 测试：列表、详情、创建、添加/移除用例、录入结果 |
| P1     | PlanListPage + 测试      | `src/features/plans/components/PlanListPage.tsx`      | 测试：计划列表、状态筛选                        |
| P1     | NewPlanPage + 测试       | `src/features/plans/components/NewPlanPage.tsx`       | 测试：创建表单、用例选择                        |
| P1     | PlanDetailPage + 测试    | `src/features/plans/components/PlanDetailPage.tsx`    | 测试：用例列表、执行统计、结果录入弹窗          |
| P1     | ResultRecordModal + 测试 | `src/features/plans/components/ResultRecordModal.tsx` | 测试：状态选择、备注输入、提交                  |

### Phase 8: 全局组件与集成

| 优先级 | 任务             | 路径                                  | AC                                       |
| ------ | ---------------- | ------------------------------------- | ---------------------------------------- |
| P0     | AppLayout + 测试 | `src/components/layout/AppLayout.tsx` | 测试：侧边栏+顶栏+内容区布局             |
| P0     | Sidebar + 测试   | `src/components/layout/Sidebar.tsx`   | 测试：菜单渲染、选中态、折叠、草稿 Badge |
| P0     | Header + 测试    | `src/components/layout/Header.tsx`    | 测试：面包屑、通知图标、用户下拉         |
| P0     | 路由配置         | `src/router/index.tsx`                | 所有路由可访问，lazy loading 正常        |
| P0     | 404 页面         | `src/components/NotFoundPage.tsx`     | 未知路由显示 404                         |
| P1     | 主题配置         | `src/styles/theme.css`                | Arco 品牌色 #7B61FF 生效                 |

### Phase 9: E2E 与验收

| 优先级 | 任务                 | 路径                               | AC                                       |
| ------ | -------------------- | ---------------------------------- | ---------------------------------------- |
| P0     | E2E: 登录流程        | `tests/e2e/auth.spec.ts`           | 登录成功跳转、失败提示                   |
| P0     | E2E: AI 生成核心流程 | `tests/e2e/core-flow.spec.ts`      | 登录→创建项目→发起生成→确认草稿→验证编号 |
| P1     | E2E: 测试计划执行    | `tests/e2e/plan-execution.spec.ts` | 创建计划→关联用例→录入结果→统计更新      |

---

## ✅ Self Review

| 检查项          | 状态 | 说明                                                           |
| --------------- | ---- | -------------------------------------------------------------- |
| 符合 React 宪法 | ✅   | 无过度抽象、状态本地化、Hook 单一职责                          |
| 避免过度设计    | ✅   | 三层架构（Component → Hook → API），Feature-Based 但不额外嵌套 |
| TDD-ready       | ✅   | 每个任务包含测试要求，MSW 覆盖所有 API mock                    |
| 完全类型安全    | ✅   | 禁止 any，枚举和接口完整定义，Zod 校验运行时类型               |
| 数据获取规范    | ✅   | 全部使用 React Query，禁止 useEffect + fetch                   |
| 全局状态最小化  | ✅   | 仅 auth token 和 sidebar 用 Zustand                            |
| 组件职责单一    | ✅   | Page 编排、Feature 交互、Shared 纯 UI，三层清晰                |
| 依赖方向正确    | ✅   | 单向依赖：Component → Hook → Service → Request，无循环         |
