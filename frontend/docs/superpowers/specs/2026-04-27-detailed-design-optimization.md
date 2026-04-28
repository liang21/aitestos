# frontend-detailed-design.md 优化设计方案

**日期**: 2026-04-27
**策略**: 保留原 10 章框架，逐章修补架构冲突、宪法违规和缺失模式
**约束**: 不阅读代码库，仅依据 PRD v2.1 / OpenAPI 3.0.3 / UX 规范 v2.0 / CLAUDE.md / constitution.md

---

## 问题全景

| #   | 问题类别                                             | 涉及章节         | 严重度 |
| --- | ---------------------------------------------------- | ---------------- | ------ |
| 1   | 目录结构与 CLAUDE.md feature-based 架构冲突          | §1.1, §4, §5, §6 | 致命   |
| 2   | 宪法 §3.2 违规：全项目使用 `any` 类型                | §4.2, §4.3, §8   | 致命   |
| 3   | 宪法 §4.1 违规：SearchTable 手动管理服务端状态       | §5.3             | 致命   |
| 4   | 宪法 §4.2 违规：Zustand 存储服务端数据               | §3               | 严重   |
| 5   | 缺失 React Query 模式（query key 工厂、useMutation） | §3, §4, §5, §6   | 严重   |
| 6   | 缺失 React Hook Form + Zod 表单模式                  | §5, §6           | 中等   |
| 7   | API 端点与 OpenAPI 不一致                            | §4.3, §6, §10    | 中等   |

---

## 第一章：项目架构设计

### §1.1 目录结构 — 完全替换

**问题**: 当前使用 `src/api/` + `src/pages/` + `src/hooks/`，与 CLAUDE.md 规定的 feature-based 架构冲突。

**修改方案**: 替换为 feature-based 目录结构，同时保留 `src/components/` (跨 feature 共享组件)、`src/store/` (仅 UI 状态)、`src/types/` (全局类型)。

```
src/
├── app/                          # 应用入口
│   ├── App.tsx                   # 根组件（ConfigProvider + RouterProvider）
│   ├── main.tsx                  # 渲染入口
│   └── providers.tsx             # QueryClientProvider + ConfigProvider
├── router/                       # 路由配置
│   ├── index.tsx                 # 路由定义（lazy loading）
│   └── RouteGuard.tsx            # 认证/权限守卫
├── lib/                          # 基础设施
│   ├── request.ts                # Axios 实例（Token 刷新、错误拦截、typed wrappers）
│   ├── query-client.ts           # React Query 全局配置
│   └── utils.ts                  # cn() 等工具函数
├── types/                        # 全局类型
│   ├── enums.ts                  # 枚举/字面量联合类型
│   └── api.ts                    # API 请求/响应类型（集中定义）
├── features/                     # 业务功能模块 (Feature-Based)
│   ├── auth/                     # 认证
│   │   ├── components/           # LoginPage, RegisterPage
│   │   ├── hooks/                # useAuth, useAuthStore
│   │   ├── schema/               # loginSchema, registerSchema (Zod)
│   │   └── services/             # auth.ts (API function)
│   ├── projects/                 # 项目管理
│   │   ├── components/           # ProjectListPage, ProjectDashboard, CreateProjectModal
│   │   ├── hooks/                # useProjects (query keys + useQuery/useMutation)
│   │   └── services/             # projects.ts
│   ├── modules/                  # 模块管理
│   │   ├── components/           # ModuleManagePage
│   │   ├── hooks/                # useModules
│   │   └── services/             # modules.ts
│   ├── testcases/                # 测试用例
│   │   ├── components/           # CaseListPage, CaseDetailPage, CreateCaseDrawer
│   │   ├── hooks/                # useTestCases
│   │   └── services/             # testcases.ts
│   ├── plans/                    # 测试计划
│   │   ├── components/           # PlanListPage, NewPlanPage, PlanDetailPage, ResultRecordModal
│   │   ├── hooks/                # usePlans
│   │   └── services/             # plans.ts
│   ├── generation/               # AI 生成
│   │   ├── components/           # GenerationTaskListPage, NewGenerationTaskPage, TaskDetailPage
│   │   ├── hooks/                # useGeneration, usePollingTask
│   │   └── services/             # generation.ts
│   ├── drafts/                   # 草稿箱
│   │   ├── components/           # DraftListPage, DraftConfirmPage
│   │   ├── hooks/                # useDrafts
│   │   └── services/             # drafts.ts
│   ├── documents/                # 知识库
│   │   ├── components/           # KnowledgeListPage, DocumentDetailPage, UploadDocumentModal
│   │   ├── hooks/                # useDocuments
│   │   └── services/             # documents.ts
│   └── configs/                  # 项目配置
│       ├── components/           # ConfigManagePage
│       ├── hooks/                # useConfigs
│       └── services/             # configs.ts
├── components/                   # 跨 Feature 共享组件
│   ├── layout/                   # AppLayout, Sidebar, Header, AuthLayout
│   ├── business/                 # StatusTag, SearchTable, ArrayEditor, StatsCard 等
│   ├── ErrorBoundary.tsx         # 全局错误边界
│   └── NotFoundPage.tsx          # 404 页面
├── store/                        # 全局 Zustand store（仅 UI 状态）
│   └── useAppStore.ts            # sidebarCollapsed
├── hooks/                        # 跨 Feature 共享 hooks
│   ├── useDebounce.ts
│   └── useMutationErrorHandler.ts
└── styles/                       # 全局样式
    └── theme.css                 # Tailwind @theme + Arco 主题变量
```

**依赖规则（与 CLAUDE.md §6 对齐）**:

- Page → Feature hooks → Feature services → `@/lib/request`
- `components/` 禁止引用 `features/` 的任何内容
- `services/` 禁止引用 `store/` 或 `hooks/`
- Feature A 禁止引用 Feature B 的内部文件
- `@/types` 可被所有层引用，但自身禁止引用任何业务代码

### §1.2 构建配置 — 微调

**保留**: vite.config.ts、环境变量（无需修改）

**新增**: React Query DevTools 配置建议

```typescript
// vite.config.ts 无需修改，但需确认依赖已安装：
// @tanstack/react-query, react-hook-form, @hookform/resolvers, zod
```

---

## 第二章：路由设计

### §2.1 路由表 — 保留

路由表无需修改，与 PRD 和 UX 规范对齐。

### §2.2 路由配置代码 — 替换

**问题**:

1. 使用 `lazy()` + `Suspense` 组件模式，应改为 React Router v7 的 `lazy` 路由属性
2. 导入路径从 `@/pages/` 改为 `@/features/`
3. `LazyPage` wrapper 不必要

**修改方案**:

```tsx
// src/router/index.tsx
import { createBrowserRouter, Navigate } from 'react-router-dom'
import { RouteGuard } from '@/router/RouteGuard'
import { ErrorBoundary } from '@/components/ErrorBoundary'
import { AppLayout } from '@/components/layout/AppLayout'

export const router = createBrowserRouter([
  // 认证页面（无布局壳）
  {
    path: '/login',
    lazy: () =>
      import('@/features/auth/components/LoginPage').then((m) => ({
        Component: () => (
          <ErrorBoundary>
            <m.LoginPage />
          </ErrorBoundary>
        ),
      })),
  },
  {
    path: '/register',
    lazy: () =>
      import('@/features/auth/components/RegisterPage').then((m) => ({
        Component: () => (
          <ErrorBoundary>
            <m.RegisterPage />
          </ErrorBoundary>
        ),
      })),
  },
  // 需认证页面
  {
    element: (
      <ErrorBoundary>
        <RouteGuard>
          <AppLayout />
        </RouteGuard>
      </ErrorBoundary>
    ),
    children: [
      { index: true, element: <Navigate to="/projects" replace /> },
      {
        path: 'projects',
        lazy: () =>
          import('@/features/projects/components/ProjectListPage').then(
            (m) => ({
              Component: m.ProjectListPage,
            })
          ),
      },
      {
        path: 'projects/:id',
        lazy: () =>
          import('@/features/projects/components/ProjectDashboard').then(
            (m) => ({
              Component: m.ProjectDashboard,
            })
          ),
      },
      // ... 其余路由路径不变，导入改为 features/ 路径
      {
        path: '*',
        element: <Navigate to="/projects" replace />,
      },
    ],
  },
])
```

---

## 第三章：状态管理设计

### §3.1 Store 架构 — 缩减

**问题**: 当前 4 个 Zustand store，其中 `useProjectStore` 和 `useDraftStore` 在 store 中调 API 存服务端数据，违反宪法 §4.2。

**修改方案**: 缩减为 1 个 Zustand store + 1 个 feature 级 auth store

```
┌──────────────┐  ┌──────────────────────────────────────────────┐
│ useAppStore  │  │ features/auth/hooks/useAuthStore             │
│ (Zustand)    │  │ (Zustand — feature 级)                       │
│              │  │                                              │
│ sidebarColl- │  │ user / token / refreshToken / isAuthenticated │
│  apsed       │  │ login() / logout() / refresh() / setTokens() │
│ toggleSidebar│  │                                              │
└──────────────┘  └──────────────────────────────────────────────┘
  纯 UI 状态               认证状态（唯一允许存 token 的例外）

其余所有服务端数据 → React Query (useQuery / useMutation)
```

**删除**:

- `useProjectStore` — 当前项目信息通过 React Query `useProjectDetail(id)` + URL params 管理
- `useDraftStore` — 草稿计数通过 React Query `useDraftCount()` + `refetchInterval` 实现

### §3.2 useAuthStore — 保留但重写

**问题**: 当前版本在 store 内 `import` API，耦合过重。

**修改方案**: 保留在 `features/auth/hooks/useAuthStore.ts`，但 token 存储和 API 调用分离

```typescript
// src/features/auth/hooks/useAuthStore.ts
import { create } from 'zustand'
import type { UserJSON } from '@/types/api'

interface AuthState {
  user: UserJSON | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isInitialized: boolean

  initialize: () => void
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  setTokens: (accessToken: string, refreshToken: string) => void
  setUser: (user: UserJSON) => void
}

// Token 持久化接口
const tokenStorage = {
  getItem: (key: string): string | null => {
    try {
      return localStorage.getItem(key)
    } catch {
      return null
    }
  },
  setItem: (key: string, value: string): void => {
    try {
      localStorage.setItem(key, value)
    } catch {
      /* 静默 */
    }
  },
  removeItem: (key: string): void => {
    try {
      localStorage.removeItem(key)
    } catch {
      /* 静默 */
    }
  },
}

// JWT 过期检查
function isTokenExpired(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.exp ? Date.now() >= payload.exp * 1000 : true
  } catch {
    return true
  }
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  token: null,
  refreshToken: null,
  isAuthenticated: false,
  isInitialized: false,

  initialize: () => {
    const accessToken = tokenStorage.getItem('access_token')
    const refreshToken = tokenStorage.getItem('refresh_token')

    if (!accessToken || !refreshToken || isTokenExpired(accessToken)) {
      tokenStorage.removeItem('access_token')
      tokenStorage.removeItem('refresh_token')
      set({ isInitialized: true })
      return
    }

    set({
      token: accessToken,
      refreshToken,
      isAuthenticated: true,
      isInitialized: true,
    })
  },

  login: async (email: string, password: string) => {
    // 动态导入避免循环依赖
    const { authApi } = await import('../services/auth')
    const response = await authApi.login({ email, password })

    tokenStorage.setItem('access_token', response.access_token)
    tokenStorage.setItem('refresh_token', response.refresh_token)

    set({
      user: response.user,
      token: response.access_token,
      refreshToken: response.refresh_token,
      isAuthenticated: true,
      isInitialized: true,
    })
  },

  logout: () => {
    tokenStorage.removeItem('access_token')
    tokenStorage.removeItem('refresh_token')
    set({ user: null, token: null, refreshToken: null, isAuthenticated: false })
  },

  setTokens: (accessToken: string, newRefreshToken: string) => {
    tokenStorage.setItem('access_token', accessToken)
    tokenStorage.setItem('refresh_token', newRefreshToken)
    set({ token: accessToken, refreshToken: newRefreshToken })
  },

  setUser: (user: UserJSON) => set({ user }),
}))
```

### §3.3 useAppStore — 精简

**问题**: 当前包含 `notifications`、`pendingDraftCount`、`fetchPendingDraftCount` — 这些是服务端数据。

**修改方案**: 仅保留纯 UI 状态

```typescript
// src/store/useAppStore.ts
import { create } from 'zustand'

interface AppState {
  sidebarCollapsed: boolean
  toggleSidebar: () => void
  setSidebarCollapsed: (collapsed: boolean) => void
}

export const useAppStore = create<AppState>((set) => ({
  sidebarCollapsed:
    typeof window !== 'undefined' ? window.innerWidth < 1280 : false,

  toggleSidebar: () => set((s) => ({ sidebarCollapsed: !s.sidebarCollapsed })),
  setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
}))
```

**迁移说明**:

- `notifications` → 改用 React Query 轮询 + Arco Notification API
- `pendingDraftCount` → 改用 `useDraftCount()` hook (React Query + refetchInterval)
- `currentProject` → 已删除，通过 `useProjectDetail(id)` + URL params 管理

### §3.4 新增：React Query 模式指南

这是原设计完全缺失的核心章节。

#### Query Key 工厂模式

每个 feature 的 hooks 文件中定义 query key 工厂：

```typescript
// src/features/projects/hooks/useProjects.ts
export const projectKeys = {
  all: ['projects'] as const,
  lists: () => [...projectKeys.all, 'list'] as const,
  list: (params: Record<string, unknown>) =>
    [...projectKeys.lists(), params] as const,
  details: () => [...projectKeys.all, 'detail'] as const,
  detail: (id: string) => [...projectKeys.details(), id] as const,
  stats: (id: string) => [...projectKeys.all, 'stats', id] as const,
}
```

#### useQuery 查询模式

```typescript
export function useProjectList(params?: {
  keywords?: string
  offset?: number
  limit?: number
}) {
  return useQuery({
    queryKey: projectKeys.list(params ?? {}),
    queryFn: () => projectsApi.list(params),
  })
}

export function useProjectDetail(id: string) {
  return useQuery({
    queryKey: projectKeys.detail(id),
    queryFn: () => projectsApi.get(id),
    enabled: !!id,
  })
}
```

#### useMutation + invalidateQueries 模式

```typescript
export function useCreateProject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateProjectRequest) => projectsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() })
    },
  })
}
```

#### 轮询模式（替代 usePolling hook）

```typescript
// 用 refetchInterval 替代自定义 usePolling
export function useTaskDetail(id: string) {
  const { data: task } = useQuery({
    queryKey: generationKeys.task(id),
    queryFn: () => generationApi.getTask(id),
    enabled: !!id,
    refetchInterval: (query) => {
      const status = query.state.data?.status
      return status === 'pending' || status === 'processing' ? 5000 : false
    },
  })
  return task
}
```

---

## 第四章：API 层设计

### §4.1 Axios 配置 — 重写

**问题**: 使用 `any` 类型泛型，违反宪法 §3.2。

**修改方案**: 提供类型安全的 wrapper 函数，用 `never` 替代 `any`

```typescript
// src/lib/request.ts
import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'

// ... token 刷新逻辑保留不变 ...

// 类型安全的 API wrappers（替代直接使用 request.get/post/put/delete）
export function get<TResponse>(
  url: string,
  config?: InternalAxiosRequestConfig
) {
  return request.get<never, TResponse>(url, config)
}

export function post<TRequest, TResponse>(
  url: string,
  data?: TRequest,
  config?: InternalAxiosRequestConfig
) {
  return request.post<never, TResponse>(url, data, config)
}

export function put<TRequest, TResponse>(
  url: string,
  data?: TRequest,
  config?: InternalAxiosRequestConfig
) {
  return request.put<never, TResponse>(url, data, config)
}

export function patch<TRequest, TResponse>(
  url: string,
  data?: TRequest,
  config?: InternalAxiosRequestConfig
) {
  return request.patch<never, TResponse>(url, data, config)
}

export function del<TResponse>(
  url: string,
  config?: InternalAxiosRequestConfig
) {
  return request.delete<never, TResponse>(url, config)
}

export default request
```

### §4.2 TypeScript 类型定义 — 保留

`src/types/enums.ts` 和 `src/types/api.ts` 保留不变，类型定义与 OpenAPI 对齐。

### §4.3 API 模块 — 迁移 + 修 any

**问题**:

1. 文件位置从 `src/api/` 迁移到 `src/features/*/services/`
2. 所有 `request.post<any, T>` 改为使用 typed wrappers
3. API 端点与 OpenAPI 不一致

**修改方案**（以 projects 为例，其他 feature 同理）:

```typescript
// src/features/projects/services/projects.ts
import { get, post, put, del } from '@/lib/request'
import type {
  Project,
  ProjectDetail,
  ProjectStats,
  CreateProjectRequest,
  UpdateProjectRequest,
} from '@/types/api'

export const projectsApi = {
  list: (params?: { keywords?: string; offset?: number; limit?: number }) =>
    get<PaginatedResponse<Project>>('/projects', { params }),

  get: (id: string) => get<ProjectDetail>(`/projects/${id}`),

  getStats: (id: string) => get<ProjectStats>(`/projects/${id}/stats`),

  create: (data: CreateProjectRequest) =>
    post<CreateProjectRequest, Project>('/projects', data),

  update: (id: string, data: UpdateProjectRequest) =>
    put<UpdateProjectRequest, Project>(`/projects/${id}`, data),

  delete: (id: string) => del<void>(`/projects/${id}`),
}
```

**API 端点对齐修复**:

| #   | 原设计端点                              | OpenAPI 定义               | 修复                                     |
| --- | --------------------------------------- | -------------------------- | ---------------------------------------- |
| 1   | `GET /generation/tasks`                 | 不存在                     | 需后端补充（已在 §10 标注）              |
| 2   | `GET /drafts`                           | 不存在                     | 需后端补充（已在 §10 标注）              |
| 3   | `PUT /projects/{id}/modules/{moduleId}` | 不存在                     | 需后端补充                               |
| 4   | `PATCH /plans/{id}/status`              | 不存在                     | 用 `PUT /plans/{id}` + `status` 字段替代 |
| 5   | `POST /plans/{planId}/results`          | `POST /plans/{id}/results` | 路径参数名统一                           |
| 6   | `POST /testcases/import`                | 不存在                     | 需后端补充（P1）                         |
| 7   | `GET /testcases/export`                 | 不存在                     | 需后端补充（P1）                         |

---

## 第五章：组件架构设计

### §5.1 Arco 主题配置 — 保留

无需修改。

### §5.2 布局组件 — 导入路径修改

所有 `@/store/useProjectStore` 引用删除，改为从 URL params 获取项目 ID 后通过 React Query 获取项目信息。

**Sidebar 修改要点**:

- 删除 `useProjectStore` 引用
- 项目上下文通过 `useProjectDetail(projectId)` 获取
- `pendingDraftCount` 改为 React Query hook

### §5.3 SearchTable — 核心重写

**问题**: 当前 SearchTable 用 `useState` + `useEffect` 管理服务端数据，违反宪法 §4.1。

**修改方案**: SearchTable 改为纯展示组件，数据获取由 React Query 管理。

```tsx
// src/components/business/SearchTable.tsx
// 设计原则：SearchTable 不管理数据，仅负责筛选栏 + 表格展示
// 数据获取由各页面的 React Query hook 负责

interface SearchTableProps<T> {
  columns: TableColumnProps[]
  data: T[]
  total: number
  loading: boolean
  currentPage: number
  pageSize: number
  onPageChange: (page: number, pageSize: number) => void
  filters?: FilterOption[]
  filterValues?: Record<string, string>
  onFilterChange?: (key: string, value: string) => void
  searchPlaceholder?: string
  keywords?: string
  onKeywordsChange?: (value: string) => void
  rowKey?: string
  toolbar?: React.ReactNode
  onRowClick?: (record: T) => void
}

export function SearchTable<T extends Record<string, unknown>>({
  columns,
  data,
  total,
  loading,
  currentPage,
  pageSize,
  onPageChange,
  filters,
  filterValues,
  onFilterChange,
  searchPlaceholder,
  keywords,
  onKeywordsChange,
  rowKey = 'id',
  toolbar,
  onRowClick,
}: SearchTableProps<T>) {
  // 纯展示组件：所有状态由父组件（React Query）传入
  return (
    <div>
      {/* 筛选栏 */}
      <div className="flex items-center justify-between mb-4">
        <Space>
          {filters?.map((f) => (
            <Select
              key={f.key}
              placeholder={f.placeholder}
              style={{ width: 140 }}
              allowClear
              value={filterValues?.[f.key]}
              onChange={(val) => onFilterChange?.(f.key, val ?? '')}
            >
              {f.options.map((opt) => (
                <Select.Option key={opt.value} value={opt.value}>
                  {opt.label}
                </Select.Option>
              ))}
            </Select>
          ))}
          <Input
            prefix={<IconSearch />}
            placeholder={searchPlaceholder}
            style={{ width: 220 }}
            allowClear
            value={keywords}
            onChange={onKeywordsChange}
          />
        </Space>
        {toolbar && <Space>{toolbar}</Space>}
      </div>

      <Table
        columns={columns}
        data={data}
        loading={loading}
        rowKey={rowKey}
        border
        stripe
        hoverable
        onRowClick={onRowClick}
        pagination={{
          total,
          current: currentPage,
          pageSize,
          onChange: onPageChange,
          showTotal: true,
          pageSizeChangeResetCurrent: true,
          sizeCanChange: true,
        }}
      />
    </div>
  )
}
```

**使用示例**（展示 React Query 如何驱动 SearchTable）:

```tsx
// 在页面组件中组合 React Query + SearchTable
function ProjectListPage() {
  const [params, setParams] = useState({ offset: 0, limit: 20, keywords: '' })
  const debouncedKeywords = useDebounce(params.keywords, 300)

  const { data, isLoading } = useProjectList({
    ...params,
    keywords: debouncedKeywords || undefined,
  })

  return (
    <SearchTable
      columns={columns}
      data={data?.data ?? []}
      total={data?.total ?? 0}
      loading={isLoading}
      currentPage={Math.floor(params.offset / params.limit) + 1}
      pageSize={params.limit}
      onPageChange={(page, pageSize) =>
        setParams((p) => ({
          ...p,
          offset: (page - 1) * pageSize,
          limit: pageSize,
        }))
      }
      keywords={params.keywords}
      onKeywordsChange={(val) => setParams((p) => ({ ...p, keywords: val }))}
    />
  )
}
```

### §5.4 其他业务组件 — 保留

StatusTag、ArrayEditor、SplitPanel、StatsCard、ReferencePanel、CaseSelector、JsonEditor — 代码不变，仅修改导入路径。

---

## 第六章：页面详细设计

**整体修改**: 所有页面的 API 数据获取描述从"页面加载时调用 API"改为"通过 React Query hook 自动管理"。

**每个页面的新增内容**:

### 表单页面新增 Zod Schema 模式

为所有含表单的页面补充 Zod schema 定义：

```typescript
// 示例：src/features/projects/schema/projectSchema.ts
import { z } from 'zod'

export const createProjectSchema = z.object({
  name: z
    .string()
    .min(2, '项目名称至少 2 个字符')
    .max(255, '项目名称最多 255 个字符'),
  prefix: z
    .string()
    .min(2, '前缀至少 2 个字符')
    .max(4, '前缀最多 4 个字符')
    .regex(/^[A-Z]+$/, '前缀仅支持大写字母'),
  description: z.string().optional(),
})

export type CreateProjectInput = z.infer<typeof createProjectSchema>
```

**所有需补充 schema 的页面**:

| 页面         | Schema 名称            | 字段                                                 |
| ------------ | ---------------------- | ---------------------------------------------------- |
| 登录         | `loginSchema`          | username, password                                   |
| 注册         | `registerSchema`       | username, email, password, role                      |
| 创建项目     | `createProjectSchema`  | name, prefix, description                            |
| 新建用例     | `createCaseSchema`     | moduleId, title, steps, expected, caseType, priority |
| 新建计划     | `createPlanSchema`     | name, description                                    |
| 新建生成任务 | `createTaskSchema`     | moduleId, prompt, count, caseType, priority          |
| 上传文档     | `uploadDocumentSchema` | name, type                                           |

### §6.x 每页修改要点

**6.3 项目列表页**: 新增 `useProjectList` hook 描述，搜索用 React Query + debounce

**6.4 项目仪表盘**: 新增 `useProjectStats` hook 描述

**6.11 草稿确认页 (核心)**:

- 草稿间导航的自动保存改为 React Query mutation
- `useDraftUpdate` mutation 描述

---

## 第七章：设计令牌实现

**无需修改**。Tailwind CSS @theme 和 Arco ConfigProvider 配置与 UX 规范对齐，保持不变。

---

## 第八章：错误处理与反馈

### §8.2 API 错误码映射 — 修复 `any`

**问题**: 表单级错误处理使用 `catch (err: any)`，违反宪法 §3.2。

**修改方案**:

```typescript
import { isAxiosError } from 'axios'

// 类型安全的错误处理模式
const handleSubmit = async (values: FormData) => {
  try {
    setLoading(true)
    await someApi.create(values)
    Message.success('创建成功')
    onClose()
  } catch (err: unknown) {
    if (isAxiosError(err) && err.response) {
      const status = err.response.status
      const message = (err.response.data as { error?: string }).error

      if (status === 400 && message) {
        form.setFields({ name: { value: values.name, errors: [message] } })
      } else if (status === 409) {
        Message.warning('资源已存在')
      }
    }
  } finally {
    setLoading(false)
  }
}
```

### §8.1 ErrorBoundary — 保留

无需修改。

---

## 第九章：性能与优化

### §9.4 API 请求缓存 — 替换

**问题**: 当前建议用 Zustand 缓存 API 数据，违反宪法 §4.1。

**修改方案**: 使用 React Query 内置缓存机制

```typescript
// 使用 staleTime 控制缓存有效期
export function useModuleList(projectId: string) {
  return useQuery({
    queryKey: moduleKeys.list(projectId),
    queryFn: () => modulesApi.list(projectId),
    enabled: !!projectId,
    staleTime: 5 * 60 * 1000, // 5 分钟内不重新请求
  })
}
```

### §9.5 轮询优化 — 简化

**问题**: 自定义 `usePolling` hook 可用 React Query `refetchInterval` 替代。

**修改方案**: 删除 `usePolling` hook，改用 `refetchInterval`（已在 §3.4 展示）。

---

## 第十章：缺失 API 与开发建议

### §10.1 缺失 API 清单 — 更新

与 OpenAPI 交叉验证后更新：

| #   | 端点                                | 方法  | 优先级 | 说明                                  |
| --- | ----------------------------------- | ----- | ------ | ------------------------------------- |
| 1   | `/generation/tasks`                 | GET   | P0     | 生成任务列表（分页、按项目+状态筛选） |
| 2   | `/drafts`                           | GET   | P0     | 全局草稿列表（跨项目，分页+筛选）     |
| 3   | `/projects/{id}/modules/{moduleId}` | PUT   | P0     | 模块编辑（名称、缩写、描述）          |
| 4   | `/plans/{id}/status`                | PATCH | P1     | 计划状态变更                          |
| 5   | `/testcases/import`                 | POST  | P1     | 用例批量导入                          |
| 6   | `/testcases/export`                 | GET   | P1     | 用例导出                              |
| 7   | `/generation/drafts/{id}`           | GET   | P0     | 草稿详情                              |
| 8   | `/generation/drafts/{id}`           | PUT   | P1     | 草稿编辑保存                          |
| 9   | `/auth/me`                          | GET   | P1     | 当前用户信息验证                      |

### §10.2 Mock 策略 — 修复 `any`

**问题**: Mock 配置使用 `(await request.json()) as any`，违反宪法 §3.2。

**修改方案**: 使用类型断言替代 `any`

```typescript
import { http, HttpResponse } from 'msw'

const mockHandlers = [
  http.post('/api/v1/auth/login', async ({ request }) => {
    const body = (await request.json()) as { email: string; password: string }
    return HttpResponse.json({
      access_token: 'mock-access-token',
      refresh_token: 'mock-refresh-token',
      user: {
        id: 'mock-user-id',
        username: body.email.split('@')[0],
        email: body.email,
        role: 'admin' as const,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    })
  }),
]
```

---

## 修改影响矩阵

| 章节         | 改动级别 | 主要修改                                      |
| ------------ | -------- | --------------------------------------------- |
| §1 项目架构  | 🔴 重写  | 目录结构完全替换为 feature-based              |
| §2 路由设计  | 🟡 中度  | 导入路径改为 features/，去掉 LazyPage wrapper |
| §3 状态管理  | 🔴 重写  | 4 store → 1+1，新增 React Query 模式指南      |
| §4 API 层    | 🟡 中度  | 迁移路径 + `any` → `never` + API 端点对齐     |
| §5 组件架构  | 🟡 中度  | SearchTable 核心重写，其他组件路径调整        |
| §6 页面设计  | 🟡 中度  | 补充 Zod schema、React Query hook 描述        |
| §7 设计令牌  | 🟢 无改  | —                                             |
| §8 错误处理  | 🟢 微调  | `any` → `unknown` + `isAxiosError`            |
| §9 性能优化  | 🟡 中度  | Zustand 缓存 → React Query staleTime          |
| §10 开发建议 | 🟢 微调  | 更新 API 清单、修复 mock any                  |

---

## 验证方案

1. **宪法合规检查**: 逐条对照 constitution.md 8 条规则，确认无违规
2. **CLAUDE.md 对齐检查**: 目录结构、命名规范、依赖规则全部对齐
3. **OpenAPI 一致性**: 所有 API 端点与 openapi.yaml 交叉验证
4. **UX 规范覆盖**: 每个页面的设计描述与 ux-design-spec.md 章节对应
