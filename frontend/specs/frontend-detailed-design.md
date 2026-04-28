# Aitestos 智能测试管理平台 — 前端详细设计文档

| 项目       | 内容                                                                           |
| ---------- | ------------------------------------------------------------------------------ |
| 产品名称   | Aitestos 智能测试管理平台                                                      |
| 文档版本   | v2.0                                                                           |
| 文档状态   | 正式发布                                                                       |
| 最后更新   | 2026-04-27                                                                     |
| 文档所有者 | 前端开发团队                                                                   |
| 关联规范   | UX 设计规范 v2.0 · 产品 PRD v2.1 · OpenAPI 3.0.3 · CLAUDE.md · constitution.md |

---

# 第一章：项目架构设计

## 1.1 目录结构

> **架构决策**: 采用 Feature-Based 架构，每个业务模块内聚 components/hooks/services，与 CLAUDE.md §6 对齐。

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
├── types/                        # 全局类型（可被所有层引用，禁止引用业务代码）
│   ├── enums.ts                  # 枚举/字面量联合类型
│   └── api.ts                    # API 请求/响应类型
├── features/                     # 业务功能模块 (Feature-Based)
│   ├── auth/                     # 认证
│   │   ├── components/           # LoginPage, RegisterPage
│   │   ├── hooks/                # useAuthStore (Zustand)
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
│   ├── business/                 # StatusTag, SearchTable, ArrayEditor, StatsCard, SplitPanel 等
│   ├── ErrorBoundary.tsx         # 全局错误边界
│   └── NotFoundPage.tsx          # 404 页面
├── store/                        # 全局 Zustand store（仅 UI 状态）
│   └── useAppStore.ts            # sidebarCollapsed + toggleSidebar
├── hooks/                        # 跨 Feature 共享 hooks
│   ├── useDebounce.ts
│   └── useMutationErrorHandler.ts
└── styles/                       # 全局样式
    └── theme.css                 # Tailwind @theme + Arco 主题变量
```

### 依赖规则

- **Page → Feature hooks → Feature services → `@/lib/request`**
- `components/` 禁止引用 `features/` 的任何内容
- `services/` 禁止引用 `store/` 或 `hooks/`
- Feature A 禁止引用 Feature B 的内部文件
- `@/types` 可被所有层引用，但自身禁止引用任何业务代码
- `@/lib` 可被 services/hooks 引用，禁止反向依赖

## 1.2 构建配置

### vite.config.ts

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import arcoReactPlugin from '@arco-plugins/vite-react'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [
    react(),
    arcoReactPlugin({
      style: 'css',
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

### 环境变量

```bash
# .env.development
VITE_API_BASE_URL=/api/v1

# .env.production
VITE_API_BASE_URL=https://api.aitestos.com/api/v1
```

---

# 第二章：路由设计

## 2.1 完整路由表

| 路由                               | 页面组件               | 布局       | 认证 | 权限   |
| ---------------------------------- | ---------------------- | ---------- | ---- | ------ |
| `/login`                           | LoginPage              | AuthLayout | 否   | —      |
| `/register`                        | RegisterPage           | AuthLayout | 否   | —      |
| `/projects`                        | ProjectListPage        | AppLayout  | 是   | —      |
| `/projects/:id`                    | ProjectDashboard       | AppLayout  | 是   | —      |
| `/projects/:id/knowledge`          | KnowledgeListPage      | AppLayout  | 是   | —      |
| `/projects/:id/knowledge/figma`    | FigmaIntegrationPage   | AppLayout  | 是   | —      |
| `/projects/:id/knowledge/:docId`   | DocumentDetailPage     | AppLayout  | 是   | —      |
| `/projects/:id/generation`         | GenerationTaskListPage | AppLayout  | 是   | —      |
| `/projects/:id/generation/new`     | NewGenerationTaskPage  | AppLayout  | 是   | —      |
| `/projects/:id/generation/:taskId` | TaskDetailPage         | AppLayout  | 是   | —      |
| `/drafts`                          | DraftListPage          | AppLayout  | 是   | —      |
| `/drafts/:draftId`                 | DraftConfirmPage       | AppLayout  | 是   | —      |
| `/projects/:id/cases`              | CaseListPage           | AppLayout  | 是   | —      |
| `/projects/:id/cases/:caseId`      | CaseDetailPage         | AppLayout  | 是   | —      |
| `/projects/:id/plans`              | PlanListPage           | AppLayout  | 是   | —      |
| `/projects/:id/plans/new`          | NewPlanPage            | AppLayout  | 是   | —      |
| `/projects/:id/plans/:planId`      | PlanDetailPage         | AppLayout  | 是   | —      |
| `/projects/:id/settings`           | ProjectSettingsPage    | AppLayout  | 是   | admin+ |
| `/projects/:id/settings/modules`   | ModuleManagePage       | AppLayout  | 是   | admin+ |
| `/projects/:id/settings/configs`   | ConfigManagePage       | AppLayout  | 是   | admin+ |

## 2.2 路由配置代码

### src/router/RouteGuard.tsx

```tsx
import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'

interface RouteGuardProps {
  children: React.ReactNode
  requireAdmin?: boolean
}

/** 解析 JWT payload 中的 exp 字段，判断 token 是否过期 */
function isTokenExpired(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.exp ? Date.now() >= payload.exp * 1000 : false
  } catch {
    return true
  }
}

export function RouteGuard({ children, requireAdmin }: RouteGuardProps) {
  const { isAuthenticated, token, user, logout } = useAuthStore()
  const location = useLocation()

  // 未认证或 token 已过期 → 跳转登录
  if (!isAuthenticated || !token || isTokenExpired(token)) {
    logout() // 清理无效 token
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  if (requireAdmin && user?.role === 'normal') {
    return <Navigate to="/projects" replace />
  }

  return <>{children}</>
}
```

### src/pages/NotFoundPage.tsx

```tsx
import { Result, Button } from '@arco-design/web-react'
import { useNavigate } from 'react-router-dom'

export function NotFoundPage() {
  const navigate = useNavigate()
  return (
    <div className="flex items-center justify-center h-screen bg-gray-1">
      <Result
        status="404"
        title="404"
        subTitle="页面不存在或资源已删除"
        extra={
          <Button type="primary" onClick={() => navigate('/projects')}>
            返回首页
          </Button>
        }
      />
    </div>
  )
}
```

### src/router/index.tsx

```tsx
import { createBrowserRouter, Navigate } from 'react-router-dom'
import { RouteGuard } from './RouteGuard'
import { ErrorBoundary } from '@/components/ErrorBoundary'
import { AppLayout } from '@/components/layout/AppLayout'

export const router = createBrowserRouter([
  // 认证页面（无布局壳，包裹 ErrorBoundary）
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
      {
        path: 'projects/:id/knowledge',
        lazy: () =>
          import('@/features/documents/components/KnowledgeListPage').then(
            (m) => ({
              Component: m.KnowledgeListPage,
            })
          ),
      },
      {
        path: 'projects/:id/knowledge/figma',
        lazy: () =>
          import('@/features/documents/components/FigmaIntegrationPage').then(
            (m) => ({
              Component: m.FigmaIntegrationPage,
            })
          ),
      },
      {
        path: 'projects/:id/knowledge/:docId',
        lazy: () =>
          import('@/features/documents/components/DocumentDetailPage').then(
            (m) => ({
              Component: m.DocumentDetailPage,
            })
          ),
      },
      {
        path: 'projects/:id/generation',
        lazy: () =>
          import('@/features/generation/components/GenerationTaskListPage').then(
            (m) => ({ Component: m.GenerationTaskListPage })
          ),
      },
      {
        path: 'projects/:id/generation/new',
        lazy: () =>
          import('@/features/generation/components/NewGenerationTaskPage').then(
            (m) => ({ Component: m.NewGenerationTaskPage })
          ),
      },
      {
        path: 'projects/:id/generation/:taskId',
        lazy: () =>
          import('@/features/generation/components/TaskDetailPage').then(
            (m) => ({
              Component: m.TaskDetailPage,
            })
          ),
      },
      {
        path: 'drafts',
        lazy: () =>
          import('@/features/drafts/components/DraftListPage').then((m) => ({
            Component: m.DraftListPage,
          })),
      },
      {
        path: 'drafts/:draftId',
        lazy: () =>
          import('@/features/drafts/components/DraftConfirmPage').then((m) => ({
            Component: m.DraftConfirmPage,
          })),
      },
      {
        path: 'projects/:id/cases',
        lazy: () =>
          import('@/features/testcases/components/CaseListPage').then((m) => ({
            Component: m.CaseListPage,
          })),
      },
      {
        path: 'projects/:id/cases/:caseId',
        lazy: () =>
          import('@/features/testcases/components/CaseDetailPage').then(
            (m) => ({
              Component: m.CaseDetailPage,
            })
          ),
      },
      {
        path: 'projects/:id/plans',
        lazy: () =>
          import('@/features/plans/components/PlanListPage').then((m) => ({
            Component: m.PlanListPage,
          })),
      },
      {
        path: 'projects/:id/plans/new',
        lazy: () =>
          import('@/features/plans/components/NewPlanPage').then((m) => ({
            Component: m.NewPlanPage,
          })),
      },
      {
        path: 'projects/:id/plans/:planId',
        lazy: () =>
          import('@/features/plans/components/PlanDetailPage').then((m) => ({
            Component: m.PlanDetailPage,
          })),
      },
      {
        path: 'projects/:id/settings',
        element: (
          <RouteGuard requireAdmin>
            {/*
              设置页使用 Element 而非 lazy，因为需要 RouteGuard 包裹。
              也可以在 lazy 回调中包裹 RouteGuard。
            */}
          </RouteGuard>
        ),
        children: [
          {
            index: true,
            lazy: () =>
              import('@/features/configs/components/ConfigManagePage').then(
                (m) => ({
                  Component: m.ConfigManagePage,
                })
              ),
          },
        ],
      },
      {
        path: 'projects/:id/settings/modules',
        element: <RouteGuard requireAdmin>{null}</RouteGuard>,
        children: [
          {
            index: true,
            lazy: () =>
              import('@/features/modules/components/ModuleManagePage').then(
                (m) => ({
                  Component: m.ModuleManagePage,
                })
              ),
          },
        ],
      },
      {
        path: '*',
        element: <Navigate to="/projects" replace />,
      },
    ],
  },
  // 404 兜底
  {
    path: '*',
    lazy: () =>
      import('@/components/NotFoundPage').then((m) => ({
        Component: m.NotFoundPage,
      })),
  },
])
```

---

# 第三章：状态管理设计

> **架构决策**: Zustand 仅管理纯 UI 状态（sidebar 折叠）和认证 token。所有服务端数据通过 React Query 管理。宪法 §4.1/§4.2 强制要求。

## 3.1 状态管理职责划分

```
┌──────────────┐  ┌────────────────────────────────────────────────┐
│ useAppStore  │  │ features/auth/hooks/useAuthStore               │
│ (全局 Zustand)│  │ (Feature 级 Zustand — 唯一存 token 的例外)     │
│              │  │                                                │
│ sidebarColl- │  │ user / token / refreshToken / isAuthenticated  │
│  apsed       │  │ initialize() / login() / logout() / setTokens()│
│ toggleSidebar│  │                                                │
│ setSidebar-  │  │ Token 持久化: localStorage                     │
│  Collapsed() │  │                                                │
└──────────────┘  └────────────────────────────────────────────────┘
  纯 UI 状态                 认证状态（允许跨 Feature 访问）

其余所有服务端数据 → React Query (useQuery / useMutation)
```

## 3.2 useAppStore — 全局 UI 状态（仅此一个全局 store）

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

> **禁止存入 useAppStore 的数据**: notifications、pendingDraftCount、currentProject — 这些都是服务端数据，必须通过 React Query 管理。

## 3.3 useAuthStore — 认证状态（Feature 级）

> UX 规范：§6.1 认证流程

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

// Token 持久化接口（便于测试和 SSR 兼容）
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

  // 从 localStorage 恢复认证状态
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

  // 登录：动态导入 authApi 避免循环依赖
  login: async (email: string, password: string) => {
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

  // 由 request.ts 拦截器调用，保持 store 与 localStorage 同步
  setTokens: (accessToken: string, newRefreshToken: string) => {
    tokenStorage.setItem('access_token', accessToken)
    tokenStorage.setItem('refresh_token', newRefreshToken)
    set({ token: accessToken, refreshToken: newRefreshToken })
  },

  setUser: (user: UserJSON) => set({ user }),
}))
```

## 3.4 React Query 模式指南

> **宪法 §4.1**: 所有服务端数据必须通过 React Query 获取和变更。

### Query Key 工厂模式

每个 feature 的 hooks 文件中定义 query key 工厂，确保缓存粒度精确：

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

### 查询 Hook 模式

```typescript
// 列表查询
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

// 详情查询（enabled 守卫防止无效请求）
export function useProjectDetail(id: string) {
  return useQuery({
    queryKey: projectKeys.detail(id),
    queryFn: () => projectsApi.get(id),
    enabled: !!id,
  })
}
```

### 变更 Hook 模式

```typescript
// 创建 + 自动刷新列表缓存
export function useCreateProject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateProjectRequest) => projectsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() })
    },
  })
}

// 更新 + 刷新详情和列表缓存
export function useUpdateProject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateProjectRequest }) =>
      projectsApi.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: projectKeys.detail(variables.id),
      })
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() })
    },
  })
}
```

### 轮询模式（替代自定义 usePolling hook）

```typescript
// 使用 refetchInterval 实现任务轮询，任务完成后自动停止
export function useTaskDetail(id: string) {
  return useQuery({
    queryKey: generationKeys.task(id),
    queryFn: () => generationApi.getTask(id),
    enabled: !!id,
    refetchInterval: (query) => {
      const status = query.state.data?.status
      return status === 'pending' || status === 'processing' ? 5000 : false
    },
  })
}
```

### 各 Feature Query Key 清单

| Feature    | Key 工厂         | 包含                                                           |
| ---------- | ---------------- | -------------------------------------------------------------- |
| projects   | `projectKeys`    | all, lists, list(params), details, detail(id), stats(id)       |
| modules    | `moduleKeys`     | all, list(projectId), detail(id)                               |
| testcases  | `caseKeys`       | all, lists, list(params), details, detail(id)                  |
| plans      | `planKeys`       | all, lists, list(params), details, detail(id), results(planId) |
| generation | `generationKeys` | all, tasks, list(params), task(id), drafts(taskId)             |
| drafts     | `draftKeys`      | all, lists, list(params), detail(id), pendingCount             |
| documents  | `documentKeys`   | all, lists, list(params), details, detail(id), chunks(docId)   |
| configs    | `configKeys`     | all, list(projectId)                                           |

````

### useDraftCount — 草稿未处理计数（Sidebar Badge）

> 用于 Sidebar 草稿箱 Badge 显示，数据通过 React Query 定时刷新。

```typescript
// src/features/drafts/hooks/useDrafts.ts
import { useQuery } from '@tanstack/react-query'
import { generationApi } from '../services/generation'

export const draftKeys = {
  all: ['drafts'] as const,
  lists: () => [...draftKeys.all, 'list'] as const,
  list: (params: Record<string, unknown>) => [...draftKeys.lists(), params] as const,
  detail: (id: string) => [...draftKeys.all, 'detail', id] as const,
  pendingCount: () => [...draftKeys.all, 'pendingCount'] as const,
}

/** 获取未处理草稿数量，用于 Sidebar Badge */
export function useDraftCount() {
  return useQuery({
    queryKey: draftKeys.pendingCount(),
    queryFn: async () => {
      const res = await generationApi.getDrafts({ status: 'pending', limit: 1 })
      return res.total
    },
    refetchInterval: 30_000, // 每 30 秒刷新
    staleTime: 15_000,
  })
}
````

---

# 第四章：API 层设计

## 4.1 Axios 配置

> **宪法 §3.2**: 禁止 `any`。提供类型安全的 wrapper 函数，用 `<never, T>` 替代 `<any, T>`。

### src/lib/request.ts

```typescript
import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
})

// ==================== Token 刷新机制 ====================

interface QueuedRequest {
  request: InternalAxiosRequestConfig
  resolve: (value: unknown) => void
  reject: (reason: unknown) => void
}

let isRefreshing = false
const pendingRequests: QueuedRequest[] = []

let authExpiredHandler: (() => void) | null = null
let tokenUpdatedHandler:
  | ((accessToken: string, refreshToken: string) => void)
  | null = null

export function setAuthExpiredHandler(handler: () => void) {
  authExpiredHandler = handler
}
export function setTokenUpdatedHandler(
  handler: (a: string, r: string) => void
) {
  tokenUpdatedHandler = handler
}

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
    } catch {}
  },
  removeItem: (key: string): void => {
    try {
      localStorage.removeItem(key)
    } catch {}
  },
}

// ==================== 请求拦截器 ====================

request.interceptors.request.use((config) => {
  const token = tokenStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// ==================== 响应拦截器 ====================

request.interceptors.response.use(
  (response) => response.data,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean
      _queued?: boolean
    }
    const isAuthEndpoint =
      originalRequest.url?.includes('/auth/login') ||
      originalRequest.url?.includes('/auth/register')

    if (
      error.response?.status === 401 &&
      originalRequest &&
      !originalRequest._retry &&
      !isAuthEndpoint
    ) {
      if (isRefreshing) {
        if (originalRequest._queued) return Promise.reject(error)
        originalRequest._queued = true
        return new Promise((resolve, reject) => {
          pendingRequests.push({ request: originalRequest, resolve, reject })
        })
      }
      originalRequest._retry = true
      isRefreshing = true
      const refreshToken = tokenStorage.getItem('refresh_token')
      if (!refreshToken) {
        authExpiredHandler?.()
        pendingRequests.forEach(({ reject }) => reject(error))
        pendingRequests.length = 0
        return Promise.reject(error)
      }
      try {
        const response = await axios.post(
          `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/auth/refresh`,
          { refresh_token: refreshToken },
          { headers: { 'Content-Type': 'application/json' } }
        )
        const { access_token, refresh_token: newRefreshToken } = response.data
        tokenStorage.setItem('access_token', access_token)
        tokenStorage.setItem('refresh_token', newRefreshToken)
        tokenUpdatedHandler?.(access_token, newRefreshToken)
        await Promise.allSettled(
          pendingRequests.map(async ({ request: q, resolve }) => {
            if (q.headers) q.headers.Authorization = `Bearer ${access_token}`
            resolve(await request(q))
          })
        )
        pendingRequests.length = 0
        if (originalRequest.headers)
          originalRequest.headers.Authorization = `Bearer ${access_token}`
        return request(originalRequest)
      } catch (refreshError) {
        tokenStorage.removeItem('access_token')
        tokenStorage.removeItem('refresh_token')
        authExpiredHandler?.()
        pendingRequests.forEach(({ reject }) => reject(refreshError))
        pendingRequests.length = 0
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }
    return Promise.reject(error)
  }
)

// ==================== 类型安全的 API Wrappers ====================

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

### App 层注册认证回调

```typescript
// src/app/App.tsx
import { setAuthExpiredHandler, setTokenUpdatedHandler } from '@/lib/request'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'

setAuthExpiredHandler(() => {
  useAuthStore.getState().logout()
  window.location.href = '/login'
})
setTokenUpdatedHandler((accessToken, refreshToken) => {
  useAuthStore.getState().setTokens(accessToken, refreshToken)
})
```

## 4.2 TypeScript 类型定义

### src/types/enums.ts

```typescript
// 用例状态
export type CaseStatus = 'unexecuted' | 'pass' | 'block' | 'fail'

// 用例类型
export type CaseType =
  | 'functionality'
  | 'performance'
  | 'api'
  | 'ui'
  | 'security'

// 计划状态
export type PlanStatus = 'draft' | 'active' | 'completed' | 'archived'

// 优先级
export type Priority = 'P0' | 'P1' | 'P2' | 'P3'

// 执行结果状态
export type ResultStatus = 'pass' | 'fail' | 'block' | 'skip'

// 任务状态
export type TaskStatus = 'pending' | 'processing' | 'completed' | 'failed'

// 草稿状态
export type DraftStatus = 'pending' | 'confirmed' | 'rejected'

// 文档类型
export type DocumentType = 'prd' | 'figma' | 'api_spec' | 'swagger' | 'markdown'

// 文档状态
export type DocumentStatus = 'pending' | 'processing' | 'completed' | 'failed'

// 用户角色
export type UserRole = 'super_admin' | 'admin' | 'normal'

// AI 置信度
export type Confidence = 'high' | 'medium' | 'low'

// 场景类型
export type SceneType = 'positive' | 'negative' | 'boundary'
```

### src/types/api.ts

```typescript
import type {
  CaseStatus,
  CaseType,
  PlanStatus,
  Priority,
  ResultStatus,
  TaskStatus,
  DraftStatus,
  DocumentType,
  DocumentStatus,
  UserRole,
  Confidence,
  SceneType,
} from './enums'

// ==================== 通用 ====================

export interface ErrorResponse {
  error: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  offset: number
  limit: number
}

// ==================== 认证 ====================

export interface RegisterRequest {
  username: string
  email: string
  password: string
  role: 'admin' | 'normal'
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  user: UserJSON
}

export interface UserJSON {
  id: string
  username: string
  email: string
  role: UserRole
  createdAt: string
  updatedAt: string
}

// ==================== 项目 ====================

export interface Project {
  id: string
  name: string
  prefix: string
  description: string
  createdAt: string
  updatedAt: string
}

export interface ProjectDetail extends Project {
  module_count: number
  case_count: number
  document_count: number
}

export type ProjectListResponse = PaginatedResponse<Project>

export interface CreateProjectRequest {
  name: string
  prefix: string
  description?: string
}

export interface UpdateProjectRequest {
  name?: string
  description?: string
}

// ==================== 模块 ====================

export interface Module {
  id: string
  projectId: string
  name: string
  abbreviation: string
  description: string
  createdBy: string
  createdAt: string
  updatedAt: string
}

export interface CreateModuleRequest {
  name: string
  abbreviation: string
  description?: string
}

export interface UpdateModuleRequest {
  name?: string
  abbreviation?: string
  description?: string
}

// ==================== 项目配置 ====================

export interface ProjectConfig {
  id: string
  projectId: string
  key: string
  value: Record<string, unknown>
  createdAt: string
  updatedAt: string
}

// ==================== 测试用例 ====================

export interface ReferencedChunk {
  chunkId: string
  documentId: string
  documentTitle: string
  similarityScore: number
}

/**
 * 引用来源详情（含内容预览）。
 * ReferencedChunk 来自 AI 元数据，不含 content 字段。
 * 展示引用内容时需额外调用 GET /knowledge/documents/{id}/chunks 获取。
 */
export interface ReferencedChunkDetail extends ReferencedChunk {
  content: string
}

export interface AiMetadata {
  generationTaskId: string
  confidence: Confidence
  referencedChunks: ReferencedChunk[]
  modelVersion: string
  generatedAt: string
}

export interface TestCase {
  id: string
  moduleId: string
  userId: string
  number: string
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

export interface CaseDetail extends TestCase {
  module_name: string
  project_name: string
  project_prefix: string
  created_by_name: string
}

export type TestCaseListResponse = PaginatedResponse<TestCase>

export interface CreateTestCaseRequest {
  module_id: string
  title: string
  preconditions?: string[]
  steps: string[]
  expected: Record<string, unknown>
  case_type: CaseType
  priority: Priority
}

export interface UpdateTestCaseRequest {
  title?: string
  preconditions?: string[]
  steps?: string[]
  expected?: Record<string, unknown>
  case_type?: CaseType
  priority?: Priority
  status?: CaseStatus
}

// ==================== 测试计划 ====================

export interface TestPlan {
  id: string
  projectId: string
  name: string
  description: string
  status: PlanStatus
  createdBy: string
  createdAt: string
  updatedAt: string
}

export interface PlanStatistics {
  total_cases: number
  passed_cases: number
  failed_cases: number
  blocked_cases: number
  skipped_cases: number
  unexecuted: number
}

export interface PlanDetail extends TestPlan {
  cases: TestCase[]
  results: TestResult[]
  stats: PlanStatistics
}

export type TestPlanListResponse = PaginatedResponse<TestPlan>

export interface CreateTestPlanRequest {
  project_id: string
  name: string
  description?: string
  case_ids?: string[]
}

export interface UpdateTestPlanRequest {
  name?: string
  description?: string
}

// ==================== 测试结果 ====================

export interface TestResult {
  id: string
  planId: string
  caseId: string
  executedBy: string
  status: ResultStatus
  note: string
  executedAt: string
  updatedAt: string
}

export interface RecordResultRequest {
  plan_id: string
  case_id: string
  status: ResultStatus
  note?: string
}

// ==================== AI 生成 ====================

export interface CreateGenerationTaskRequest {
  project_id: string
  module_id: string
  prompt: string
  case_count?: number
  scene_types?: SceneType[]
  priority?: Priority
  case_type?: CaseType
}

export interface GenerationTask {
  id: string
  projectId: string
  moduleId: string
  status: TaskStatus
  prompt: string
  result: Record<string, unknown>
  createdAt: string
  updatedAt: string
}

export interface CaseDraft {
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

export interface BatchConfirmResult {
  success_count: number
  failed_count: number
  success_ids: string[]
  failed_ids: string[]
  errors: string[]
}

// ==================== 文档 ====================

export interface Document {
  id: string
  projectId: string
  name: string
  type: DocumentType
  url: string
  contentText: string
  status: DocumentStatus
  createdBy: string
  createdAt: string
  updatedAt: string
}

export interface DocumentDetail extends Document {
  chunk_count: number
}

export type DocumentListResponse = PaginatedResponse<Document>

export interface UploadDocumentRequest {
  project_id: string
  name: string
  type: DocumentType
}

export interface ChunkInfo {
  id: string
  document_id: string
  chunk_index: number
  content: string
  created_at: number
}

// ==================== 项目统计 ====================

export interface ProjectStatistics {
  module_count: number
  case_count: number
  document_count: number
  pass_rate: number
  coverage_rate: number
  ai_generated_count: number
  recent_tasks: TaskSummary[]
  pass_rate_trend: TrendData[]
  updated_at: string
}

export interface TaskSummary {
  id: string
  status: TaskStatus
  result_summary: TaskResultSummary
  created_at: string
}

export interface TaskResultSummary {
  total_drafts: number
  confirmed_count: number
  rejected_count: number
}

export interface TrendData {
  date: string
  rate: number
}
```

## 4.3 API 模块

### src/features/auth/services/auth.ts

```typescript
import { post } from '@/lib/request'
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RefreshTokenRequest,
  UserJSON,
} from '@/types/api'

export const authApi = {
  login: (data: LoginRequest) =>
    post<LoginRequest, LoginResponse>('/auth/login', data),

  register: (data: RegisterRequest) =>
    post<RegisterRequest, UserJSON>('/auth/register', data),

  refresh: (data: RefreshTokenRequest) =>
    post<RefreshTokenRequest, LoginResponse>('/auth/refresh', data),
}
```

### src/features/projects/services/projects.ts

```typescript
import { get, post, put, del } from '@/lib/request'
import type {
  Project,
  ProjectDetail,
  ProjectListResponse,
  CreateProjectRequest,
  UpdateProjectRequest,
  ProjectStatistics,
} from '@/types/api'

export const projectsApi = {
  list: (params?: { offset?: number; limit?: number; keywords?: string }) =>
    get<ProjectListResponse>('/projects', { params }),

  get: (id: string) => get<ProjectDetail>(`/projects/${id}`),

  create: (data: CreateProjectRequest) =>
    post<CreateProjectRequest, Project>('/projects', data),

  update: (id: string, data: UpdateProjectRequest) =>
    put<UpdateProjectRequest, Project>(`/projects/${id}`, data),

  delete: (id: string) => del<void>(`/projects/${id}`),

  getStats: (id: string) => get<ProjectStatistics>(`/projects/${id}/stats`),
}
```

### src/features/modules/services/modules.ts

```typescript
import { get, post, put, del } from '@/lib/request'
import type {
  Module,
  CreateModuleRequest,
  UpdateModuleRequest,
} from '@/types/api'

export const modulesApi = {
  list: (projectId: string) => get<Module[]>(`/projects/${projectId}/modules`),

  create: (projectId: string, data: CreateModuleRequest) =>
    post<CreateModuleRequest, Module>(`/projects/${projectId}/modules`, data),

  update: (
    projectId: string,
    moduleId: string,
    data: Partial<CreateModuleRequest>
  ) =>
    put<UpdateModuleRequest, Module>(
      `/projects/${projectId}/modules/${moduleId}`,
      data
    ),

  delete: (moduleId: string) => del<void>(`/modules/${moduleId}`),
}
```

### src/features/configs/services/configs.ts

```typescript
import { get, post, put, del } from '@/lib/request'
import type { ProjectConfig } from '@/types/api'

interface ConfigExportItem {
  key: string
  value: Record<string, unknown>
  description?: string
}

interface ConfigImportResult {
  imported: number
  failed: number
  errors: string[]
}

interface ConfigImportRequest {
  configs: ConfigExportItem[]
}

export const configsApi = {
  list: (projectId: string) =>
    get<ProjectConfig[]>(`/projects/${projectId}/configs`),

  set: (projectId: string, key: string, value: Record<string, unknown>) =>
    put<{ value: Record<string, unknown> }, void>(
      `/projects/${projectId}/configs/${key}`,
      { value }
    ),

  delete: (projectId: string, key: string) =>
    del<void>(`/projects/${projectId}/configs/${key}`),

  import: (projectId: string, configs: ConfigExportItem[]) =>
    post<ConfigImportRequest, ConfigImportResult>(
      `/projects/${projectId}/configs/import`,
      { configs }
    ),

  export: (projectId: string) =>
    get<ConfigExportItem[]>(`/projects/${projectId}/configs/export`),
}
```

### src/features/testcases/services/testcases.ts

```typescript
import { get, post, put, del } from '@/lib/request'
import type {
  TestCase,
  CaseDetail,
  TestCaseListResponse,
  CreateTestCaseRequest,
  UpdateTestCaseRequest,
  CaseStatus,
  CaseType,
  Priority,
} from '@/types/api'

interface TestCaseListParams {
  project_id: string
  module_id?: string
  status?: CaseStatus
  case_type?: CaseType
  priority?: Priority
  keywords?: string
  offset?: number
  limit?: number
}

export const testcasesApi = {
  list: (params: TestCaseListParams) =>
    get<TestCaseListResponse>('/testcases', { params }),

  get: (id: string) => get<CaseDetail>(`/testcases/${id}`),

  create: (data: CreateTestCaseRequest) =>
    post<CreateTestCaseRequest, TestCase>('/testcases', data),

  update: (id: string, data: UpdateTestCaseRequest) =>
    put<UpdateTestCaseRequest, TestCase>(`/testcases/${id}`, data),

  delete: (id: string) => del<void>(`/testcases/${id}`),
}
```

### src/features/plans/services/plans.ts

```typescript
import { get, post, put, patch, del } from '@/lib/request'
import type {
  TestPlan,
  PlanDetail,
  TestPlanListResponse,
  CreateTestPlanRequest,
  UpdateTestPlanRequest,
  TestResult,
  RecordResultRequest,
  PlanStatus,
} from '@/types/api'

interface PlanListParams {
  project_id: string
  status?: PlanStatus
  keywords?: string
  offset?: number
  limit?: number
}

interface AddCasesRequest {
  case_ids: string[]
}

interface UpdateStatusRequest {
  status: PlanStatus
}

export const plansApi = {
  list: (params: PlanListParams) =>
    get<TestPlanListResponse>('/plans', { params }),

  get: (id: string) => get<PlanDetail>(`/plans/${id}`),

  create: (data: CreateTestPlanRequest) =>
    post<CreateTestPlanRequest, TestPlan>('/plans', data),

  update: (id: string, data: UpdateTestPlanRequest) =>
    put<UpdateTestPlanRequest, TestPlan>(`/plans/${id}`, data),

  delete: (id: string) => del<void>(`/plans/${id}`),

  addCases: (planId: string, caseIds: string[]) =>
    post<AddCasesRequest, void>(`/plans/${planId}/cases`, {
      case_ids: caseIds,
    }),

  removeCase: (planId: string, caseId: string) =>
    del<void>(`/plans/${planId}/cases/${caseId}`),

  getResults: (planId: string) => get<TestResult[]>(`/plans/${planId}/results`),

  recordResult: (planId: string, data: RecordResultRequest) =>
    post<RecordResultRequest, TestResult>(`/plans/${planId}/results`, data),

  updateStatus: (planId: string, status: PlanStatus) =>
    patch<UpdateStatusRequest, void>(`/plans/${planId}/status`, { status }),
}
```

### src/features/generation/services/generation.ts

```typescript
import { get, post } from '@/lib/request'
import type {
  GenerationTask,
  CaseDraft,
  TestCase,
  CreateGenerationTaskRequest,
  BatchConfirmResult,
  TaskStatus,
  DraftStatus,
} from '@/types/api'

interface GenerationTaskListParams {
  project_id: string
  status?: TaskStatus
  offset?: number
  limit?: number
}

interface DraftListParams {
  project_id?: string
  module_id?: string
  status?: DraftStatus
  keywords?: string
  offset?: number
  limit?: number
}

interface DraftListResponse {
  data: CaseDraft[]
  total: number
  offset: number
  limit: number
}

interface ConfirmDraftRequest {
  module_id: string
}

interface RejectDraftRequest {
  reason?: string
  feedback?: string
}

interface BatchConfirmRequest {
  draft_ids: string[]
  module_id: string
}

export const generationApi = {
  // 生成任务
  createTask: (data: CreateGenerationTaskRequest) =>
    post<CreateGenerationTaskRequest, GenerationTask>(
      '/generation/tasks',
      data
    ),

  getTask: (id: string) => get<GenerationTask>(`/generation/tasks/${id}`),

  listTasks: (params: GenerationTaskListParams) =>
    get<{ data: GenerationTask[]; total: number }>('/generation/tasks', {
      params,
    }),

  // 草稿
  getTaskDrafts: (taskId: string) =>
    get<CaseDraft[]>(`/generation/tasks/${taskId}/drafts`),

  getDrafts: (params: DraftListParams) =>
    get<DraftListResponse>('/drafts', { params }),

  confirmDraft: (draftId: string, moduleId: string) =>
    post<ConfirmDraftRequest, TestCase>(
      `/generation/drafts/${draftId}/confirm`,
      { module_id: moduleId }
    ),

  rejectDraft: (draftId: string, data: RejectDraftRequest) =>
    post<RejectDraftRequest, void>(
      `/generation/drafts/${draftId}/reject`,
      data
    ),

  batchConfirm: (draftIds: string[], moduleId: string) =>
    post<BatchConfirmRequest, BatchConfirmResult>(
      '/generation/drafts/batch-confirm',
      {
        draft_ids: draftIds,
        module_id: moduleId,
      }
    ),
}
```

### src/features/documents/services/documents.ts

```typescript
import { get, post, del } from '@/lib/request'
import type {
  Document,
  DocumentDetail,
  DocumentListResponse,
  UploadDocumentRequest,
  ChunkInfo,
  DocumentType,
  DocumentStatus,
} from '@/types/api'

interface DocumentListParams {
  project_id: string
  type?: DocumentType
  status?: DocumentStatus
  offset?: number
  limit?: number
}

export const documentsApi = {
  list: (params: DocumentListParams) =>
    get<DocumentListResponse>('/knowledge/documents', { params }),

  get: (id: string) => get<DocumentDetail>(`/knowledge/documents/${id}`),

  upload: (data: UploadDocumentRequest) =>
    post<UploadDocumentRequest, Document>('/knowledge/documents', data),

  delete: (id: string) => del<void>(`/knowledge/documents/${id}`),

  getChunks: (id: string) =>
    get<ChunkInfo[]>(`/knowledge/documents/${id}/chunks`),
}
```

---

# 第五章：组件架构设计

## 5.1 Arco 主题配置

### src/styles/arco-theme.ts

```typescript
import type { Token } from '@arco-design/web-react/es/theme/interface'

export const arcoTheme: Partial<Token> = {
  // 品牌色
  colorPrimary: '#7B61FF',
  colorPrimaryHover: '#6B4FE0',
  colorPrimaryActive: '#5A3DC0',

  // 功能色
  colorSuccess: '#00B42A',
  colorWarning: '#FF7D00',
  colorDanger: '#F53F3F',
  colorInfo: '#86909C',

  // 文字色
  colorText: '#4E5969',
  colorTextSecondary: '#86909C',
  colorTextDisabled: '#C9CDD4',
  colorTextPlaceholder: '#A9AEB8',

  // 边框
  colorBorder: '#E5E6EB',
  colorBorderSecondary: '#F2F3F5',

  // 填充
  colorFill: '#F7F8FA',
  colorFillSecondary: '#F2F3F5',

  // 背景
  colorBgContainer: '#FFFFFF',
  colorBgElevated: '#FFFFFF',
  colorBgLayout: '#F7F8FA',

  // 圆角
  borderRadius: 4,
  borderRadiusSmall: 2,
  borderRadiusLarge: 8,

  // 字号
  fontSize: 14,
  fontSizeBody: 14,
  fontSizeTitleH1: 24,
  fontSizeTitleH2: 20,
  fontSizeTitleH3: 16,

  // 高度
  controlHeight: 32,
  controlHeightSmall: 28,
  controlHeightMini: 24,
}
```

## 5.2 布局组件

### AppLayout 结构

```
┌──────────────────────────────────────────────────┐
│  Header (48px, fixed)                             │
├─────────┬────────────────────────────────────────┤
│         │  面包屑栏 (40px)                         │
│ Sidebar ├────────────────────────────────────────┤
│ (220px) │                                        │
│         │  <Outlet /> (主内容区, padding: 24px)    │
│         │                                        │
└─────────┴────────────────────────────────────────┘
```

### src/components/layout/AppLayout.tsx

```tsx
import { Layout } from '@arco-design/web-react'
import { Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { Header } from './Header'
import { useAppStore } from '@/store/useAppStore'

const { Sider, Content } = Layout

export function AppLayout() {
  const sidebarCollapsed = useAppStore((s) => s.sidebarCollapsed)

  return (
    <Layout className="h-screen">
      <Header />
      <Layout>
        <Sider
          width={sidebarCollapsed ? 64 : 220}
          collapsible
          collapsed={sidebarCollapsed}
          trigger={null}
          className="border-r border-gray-3 bg-white"
          style={{ transition: 'width 300ms cubic-bezier(0.4, 0, 0.2, 1)' }}
        >
          <Sidebar />
        </Sider>
        <Layout>
          <Content className="bg-gray-1 p-6 overflow-auto">
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  )
}
```

### src/components/layout/Sidebar.tsx

```tsx
import { Menu } from '@arco-design/web-react'
import { useNavigate, useLocation, useParams } from 'react-router-dom'
import {
  FolderOpen,
  LayoutDashboard,
  BookOpen,
  Sparkles,
  FileCheck,
  ClipboardList,
  Settings,
  FileEdit,
} from 'lucide-react'
import { useAppStore } from '@/store/useAppStore'
import { Badge } from '@arco-design/web-react'

const MenuItem = Menu.Item
const SubMenu = Menu.SubMenu

export function Sidebar() {
  const navigate = useNavigate()
  const location = useLocation()
  const params = useParams()
  const projectId = params.id
  // 项目信息通过 React Query 获取（替代 useProjectStore）
  const { data: currentProject } = useProjectDetail(projectId ?? '')
  const sidebarCollapsed = useAppStore((s) => s.sidebarCollapsed)
  // 草稿计数通过 React Query 轮询获取（替代 useAppStore.pendingDraftCount）
  // useDraftCount 定义于 src/features/drafts/hooks/useDrafts.ts
  const { data: draftCount } = useDraftCount()
  const pendingDraftCount = draftCount ?? 0
  const projId = params.id

  // 使用路径前缀匹配，支持子路由高亮
  const selectedKeys = (() => {
    const path = location.pathname
    // 精确匹配优先
    if (path === '/projects') return ['/projects']
    if (path === '/drafts') return ['/drafts']
    // 前缀匹配项目子路由
    if (projectId) {
      if (path.includes('/settings')) return [`/projects/${projectId}/settings`]
      if (path.includes('/plans')) return [`/projects/${projectId}/plans`]
      if (path.includes('/cases')) return [`/projects/${projectId}/cases`]
      if (path.includes('/generation'))
        return [`/projects/${projectId}/generation`]
      if (path.includes('/knowledge'))
        return [`/projects/${projectId}/knowledge`]
      return [`/projects/${projectId}`]
    }
    return [path]
  })()
  const openKeys = currentProject ? ['project'] : []

  return (
    <Menu
      selectedKeys={selectedKeys}
      defaultOpenKeys={openKeys}
      style={{ width: '100%' }}
      collapse={sidebarCollapsed}
      onClickMenuItem={(key) => navigate(key)}
    >
      <MenuItem key="/projects">
        <FolderOpen size={16} />
        项目列表
      </MenuItem>

      {currentProject && projectId && (
        <SubMenu
          key="project"
          title={<span className="truncate">{currentProject.name}</span>}
        >
          <MenuItem key={`/projects/${projectId}`}>
            <LayoutDashboard size={16} />
            仪表盘
          </MenuItem>
          <MenuItem key={`/projects/${projectId}/knowledge`}>
            <BookOpen size={16} />
            知识库
          </MenuItem>
          <MenuItem key={`/projects/${projectId}/generation`}>
            <Sparkles size={16} />
            AI 生成
          </MenuItem>
          <MenuItem key={`/projects/${projectId}/cases`}>
            <FileCheck size={16} />
            测试用例
          </MenuItem>
          <MenuItem key={`/projects/${projectId}/plans`}>
            <ClipboardList size={16} />
            测试计划
          </MenuItem>
          <MenuItem key={`/projects/${projectId}/settings`}>
            <Settings size={16} />
            项目设置
          </MenuItem>
        </SubMenu>
      )}

      <MenuItem key="/drafts">
        <Badge count={pendingDraftCount} offset={[4, -2]}>
          <div className="flex items-center gap-2">
            <FileEdit size={16} />
            {!sidebarCollapsed && '草稿箱'}
          </div>
        </Badge>
      </MenuItem>
    </Menu>
  )
}
```

### src/components/layout/Header.tsx

```tsx
import {
  Breadcrumb,
  Button,
  Badge,
  Dropdown,
  Avatar,
  Popconfirm,
} from '@arco-design/web-react'
import {
  IconLeft,
  IconRight,
  IconNotification,
} from '@arco-design/web-react/icon'
import { useNavigate, useLocation } from 'react-router-dom'
import { Bell, PanelLeftClose, PanelLeftOpen } from 'lucide-react'
import { useAppStore } from '@/store/useAppStore'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'

export function Header() {
  const navigate = useNavigate()
  const location = useLocation()
  const { sidebarCollapsed, toggleSidebar } = useAppStore()
  const { user, logout } = useAuthStore()
  // 项目信息从 URL params + React Query 获取（替代 useProjectStore）
  const params = useParams()
  const { data: currentProject } = useProjectDetail(params.id ?? '')

  // 生成面包屑
  const breadcrumbs = generateBreadcrumbs(location.pathname, currentProject)

  return (
    <div className="h-12 border-b border-gray-3 bg-white flex items-center justify-between px-4 sticky top-0 z-50">
      <div className="flex items-center gap-3">
        <Button
          type="text"
          size="small"
          icon={
            sidebarCollapsed ? (
              <PanelLeftOpen size={16} />
            ) : (
              <PanelLeftClose size={16} />
            )
          }
          onClick={toggleSidebar}
        />
        <Breadcrumb>
          {breadcrumbs.map((item) => (
            <Breadcrumb.Item key={item.path ?? item.label}>
              {item.path ? (
                <span
                  className="cursor-pointer text-gray-7"
                  onClick={() => navigate(item.path!)}
                >
                  {item.label}
                </span>
              ) : (
                item.label
              )}
            </Breadcrumb.Item>
          ))}
        </Breadcrumb>
      </div>

      <div className="flex items-center gap-4">
        <Badge count={0} dot>
          <Bell size={18} className="text-gray-7 cursor-pointer" />
        </Badge>

        <Dropdown
          trigger="click"
          droplist={
            <div className="py-1">
              <Popconfirm
                title="确认退出登录？"
                onOk={() => {
                  logout()
                  navigate('/login')
                }}
              >
                <div className="px-3 py-1.5 text-sm cursor-pointer hover:bg-gray-2">
                  退出登录
                </div>
              </Popconfirm>
            </div>
          }
        >
          <div className="flex items-center gap-2 cursor-pointer">
            <Avatar size={28} className="bg-primary text-white text-xs">
              {user?.username?.charAt(0).toUpperCase() || 'U'}
            </Avatar>
            <span className="text-sm text-gray-8">{user?.username}</span>
          </div>
        </Dropdown>
      </div>
    </div>
  )
}

function generateBreadcrumbs(
  pathname: string,
  project: { name: string } | null
) {
  const crumbs: Array<{ label: string; path?: string }> = [
    { label: '首页', path: '/projects' },
  ]
  // 按路径分段匹配，避免字符串 includes 误判
  const segments = pathname.split('/').filter(Boolean)

  if (segments[0] === 'drafts') {
    crumbs.push({ label: '草稿箱' })
    return crumbs
  }

  if (segments[0] === 'projects' && segments[1] && project) {
    crumbs.push({ label: project.name, path: `/projects/${segments[1]}` })

    const sub = segments[2]
    if (sub === 'knowledge') crumbs.push({ label: '知识库' })
    else if (sub === 'generation') crumbs.push({ label: 'AI 生成' })
    else if (sub === 'cases') crumbs.push({ label: '测试用例' })
    else if (sub === 'plans') crumbs.push({ label: '测试计划' })
    else if (sub === 'settings') crumbs.push({ label: '项目设置' })
  }

  return crumbs
}
```

## 5.3 核心业务组件

### StatusTag — 状态标签（统一色彩映射）

> UX 规范：§5.6, §2.1 语义色映射表

```tsx
// src/components/business/StatusTag.tsx
import { Tag } from '@arco-design/web-react'
import type {
  CaseStatus,
  CaseType,
  Confidence,
  DocumentStatus,
  DocumentType,
  DraftStatus,
  PlanStatus,
  Priority,
  TaskStatus,
} from '@/types/enums'

type StatusCategory =
  | 'caseStatus'
  | 'planStatus'
  | 'taskStatus'
  | 'draftStatus'
  | 'priority'
  | 'confidence'
  | 'caseType'
  | 'documentType'
  | 'documentStatus'

interface StatusTagProps {
  type: StatusCategory
  value: string
  size?: 'default' | 'small'
}

// 色彩映射配置
const COLOR_MAP = {
  caseStatus: {
    unexecuted: { label: '未执行', color: '#86909C', textColor: '#4E5969' },
    pass: { label: '通过', color: '#00B42A', textColor: '#008A27' },
    block: { label: '阻塞', color: '#FF7D00', textColor: '#CC6200' },
    fail: { label: '失败', color: '#F53F3F', textColor: '#CB2634' },
  },
  planStatus: {
    draft: { label: '草稿', color: '#86909C', textColor: '#4E5969' },
    active: { label: '进行中', color: '#7B61FF', textColor: '#5A3DC0' },
    completed: { label: '已完成', color: '#00B42A', textColor: '#008A27' },
    archived: { label: '已归档', color: '#C9CDD4', textColor: '#86909C' },
  },
  taskStatus: {
    pending: { label: '待处理', color: '#86909C', textColor: '#4E5969' },
    processing: { label: '处理中', color: '#7B61FF', textColor: '#5A3DC0' },
    completed: { label: '已完成', color: '#00B42A', textColor: '#008A27' },
    failed: { label: '失败', color: '#F53F3F', textColor: '#CB2634' },
  },
  draftStatus: {
    pending: { label: '待处理', color: '#FF7D00', textColor: '#CC6200' },
    confirmed: { label: '已确认', color: '#00B42A', textColor: '#008A27' },
    rejected: { label: '已拒绝', color: '#F53F3F', textColor: '#CB2634' },
  },
  priority: {
    P0: { label: 'P0 紧急', color: '#F53F3F', textColor: '#CB2634' },
    P1: { label: 'P1 高', color: '#FF7D00', textColor: '#CC6200' },
    P2: { label: 'P2 中', color: '#7B61FF', textColor: '#5A3DC0' },
    P3: { label: 'P3 低', color: '#86909C', textColor: '#4E5969' },
  },
  confidence: {
    high: { label: '高置信度', color: '#00B42A', textColor: '#008A27' },
    medium: { label: '中置信度', color: '#FF7D00', textColor: '#CC6200' },
    low: { label: '低置信度', color: '#F53F3F', textColor: '#CB2634' },
  },
  caseType: {
    functionality: {
      label: '功能测试',
      color: '#7B61FF',
      textColor: '#5A3DC0',
    },
    performance: { label: '性能测试', color: '#3491FA', textColor: '#1677FF' },
    api: { label: 'API 测试', color: '#0FC6C2', textColor: '#0A8A87' },
    ui: { label: 'UI 测试', color: '#F77234', textColor: '#CC5E2A' },
    security: { label: '安全测试', color: '#722ED1', textColor: '#531DAB' },
  },
  documentType: {
    prd: { label: 'PRD 文档', color: '#7B61FF', textColor: '#5A3DC0' },
    figma: { label: 'Figma 设计', color: '#A259FF', textColor: '#8B3DD9' },
    api_spec: { label: 'API 规范', color: '#0FC6C2', textColor: '#0A8A87' },
    swagger: { label: 'Swagger', color: '#00B42A', textColor: '#008A27' },
    markdown: { label: 'Markdown', color: '#86909C', textColor: '#4E5969' },
  },
  documentStatus: {
    pending: { label: '待处理', color: '#86909C', textColor: '#4E5969' },
    processing: { label: '处理中', color: '#7B61FF', textColor: '#5A3DC0' },
    completed: { label: '已完成', color: '#00B42A', textColor: '#008A27' },
    failed: { label: '失败', color: '#F53F3F', textColor: '#CB2634' },
  },
} as const

export function StatusTag({ type, value, size = 'default' }: StatusTagProps) {
  const category = COLOR_MAP[type as keyof typeof COLOR_MAP]
  const config = category?.[value as keyof typeof category]

  if (!config) {
    return <Tag size={size}>{value}</Tag>
  }

  return (
    <Tag
      size={size}
      style={{
        color: config.textColor,
        backgroundColor: `${config.color}1A`, // 10% opacity
        borderColor: 'transparent',
        fontSize: size === 'small' ? 11 : 12,
      }}
    >
      {config.label}
    </Tag>
  )
}
```

### ArrayEditor — 数组编辑器

> UX 规范：§5.2 数组编辑器

```tsx
// src/components/business/ArrayEditor.tsx
import { Input, Button } from '@arco-design/web-react'
import { Plus, Trash2, ChevronUp, ChevronDown } from 'lucide-react'
import { useRef } from 'react'

interface ArrayEditorProps {
  value: string[]
  onChange: (value: string[]) => void
  minItems?: number
  addButtonText?: string
  placeholder?: string
}

/** 生成稳定的 item ID，避免 index 作为 key 导致排序时焦点错位 */
function useStableKeys(length: number): string[] {
  const keysRef = useRef<string[]>([])
  // 当长度增加时，为新项追加新 key
  while (keysRef.current.length < length) {
    keysRef.current.push(
      `item-${keysRef.current.length}-${Math.random().toString(36).slice(2, 9)}`
    )
  }
  // 当长度减少时，裁剪多余 key
  keysRef.current = keysRef.current.slice(0, length)
  return keysRef.current
}

export function ArrayEditor({
  value,
  onChange,
  minItems = 0,
  addButtonText = '添加项目',
  placeholder = '请输入内容',
}: ArrayEditorProps) {
  const stableKeys = useStableKeys(value.length)

  const addItem = () => {
    onChange([...value, ''])
  }

  const removeItem = (index: number) => {
    onChange(value.filter((_, i) => i !== index))
  }

  const updateItem = (index: number, newValue: string) => {
    const updated = [...value]
    updated[index] = newValue
    onChange(updated)
  }

  const moveUp = (index: number) => {
    if (index === 0) return
    const updated = [...value]
    ;[updated[index - 1], updated[index]] = [updated[index], updated[index - 1]]
    onChange(updated)
  }

  const moveDown = (index: number) => {
    if (index === value.length - 1) return
    const updated = [...value]
    ;[updated[index], updated[index + 1]] = [updated[index + 1], updated[index]]
    onChange(updated)
  }

  return (
    <div className="space-y-2">
      {value.map((item, index) => (
        <div key={stableKeys[index]} className="flex items-center gap-2">
          <span className="text-gray-6 text-sm w-5 text-right">
            {index + 1}.
          </span>
          <Input
            value={item}
            onChange={(val) => updateItem(index, val)}
            placeholder={placeholder}
            className="flex-1"
          />
          <div className="flex items-center gap-1">
            <Button
              type="text"
              size="mini"
              icon={<ChevronUp size={14} />}
              disabled={index === 0}
              onClick={() => moveUp(index)}
            />
            <Button
              type="text"
              size="mini"
              icon={<ChevronDown size={14} />}
              disabled={index === value.length - 1}
              onClick={() => moveDown(index)}
            />
            <Button
              type="text"
              size="mini"
              status="danger"
              icon={<Trash2 size={14} />}
              disabled={value.length <= minItems}
              onClick={() => removeItem(index)}
            />
          </div>
        </div>
      ))}
      <Button
        type="dashed"
        long
        icon={<Plus size={14} />}
        onClick={addItem}
        className="border-dashed"
      >
        {addButtonText}
      </Button>
    </div>
  )
}
```

### SearchTable — 搜索筛选表格

> UX 规范：§5.1 表格组件

> **宪法 §4.1 合规**: SearchTable 为纯展示组件，**不持有服务端数据**。所有数据通过 React Query hook 获取后以 props 传入。组件仅负责筛选栏 UI 和表格展示。筛选条件和分页状态通过回调上抛，由父组件透传给 React Query hook 的 queryKey 参数。

```tsx
// src/components/business/SearchTable.tsx
import {
  Table,
  Input,
  Select,
  Button,
  Empty,
  Space,
} from '@arco-design/web-react'
import { IconSearch } from '@arco-design/web-react/icon'
import type { TableColumnProps } from '@arco-design/web-react'

interface FilterOption {
  key: string
  placeholder: string
  options: Array<{ label: string; value: string }>
}

interface SearchTableProps<T> {
  columns: TableColumnProps[]
  /** 服务端数据，由父组件通过 React Query hook 获取 */
  data: T[]
  /** 加载状态，来自 React Query 的 isFetching */
  loading: boolean
  /** 总记录数，来自 API 响应的 total */
  total: number
  /** 当前分页，受控模式 */
  pagination: { offset: number; limit: number }
  /** 分页变更回调 */
  onPaginationChange: (pagination: { offset: number; limit: number }) => void
  /** 筛选配置 */
  filters?: FilterOption[]
  /** 当前筛选值，受控模式 */
  filterValues?: Record<string, string>
  /** 筛选变更回调 */
  onFilterChange?: (filters: Record<string, string>) => void
  /** 搜索关键词，受控模式 */
  keywords?: string
  /** 搜索变更回调（父组件配合 useDebounce 防抖后传入 React Query） */
  onKeywordsChange?: (keywords: string) => void
  searchPlaceholder?: string
  rowKey?: string
  toolbar?: React.ReactNode
  onRowClick?: (record: T) => void
}

export function SearchTable<T extends Record<string, unknown>>({
  columns,
  data,
  loading,
  total,
  pagination,
  onPaginationChange,
  filters = [],
  filterValues = {},
  onFilterChange,
  keywords = '',
  onKeywordsChange,
  searchPlaceholder = '搜索关键词',
  rowKey = 'id',
  toolbar,
  onRowClick,
}: SearchTableProps<T>) {
  const handlePageChange = (page: number, pageSize: number) => {
    onPaginationChange({ offset: (page - 1) * pageSize, limit: pageSize })
  }

  return (
    <div>
      {/* 筛选栏 */}
      <div className="flex items-center justify-between mb-4">
        <Space>
          {filters.map((f) => (
            <Select
              key={f.key}
              placeholder={f.placeholder}
              style={{ width: 140 }}
              allowClear
              value={filterValues[f.key] || undefined}
              onChange={(val) =>
                onFilterChange?.({ ...filterValues, [f.key]: val || '' })
              }
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

      {/* 表格 */}
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
          current: Math.floor(pagination.offset / pagination.limit) + 1,
          pageSize: pagination.limit,
          onChange: handlePageChange,
          showTotal: true,
          pageSizeChangeResetCurrent: true,
          sizeCanChange: true,
        }}
        empty={<Empty />}
      />
    </div>
  )
}
```

### SplitPanel — 分栏面板

> UX 规范：§5.8

```tsx
// src/components/business/SplitPanel.tsx
import { ResizeBox } from '@arco-design/web-react'

interface SplitPanelProps {
  leftContent: React.ReactNode
  rightContent: React.ReactNode
  defaultLeftWidth?: number
  minLeftWidth?: number
}

export function SplitPanel({
  leftContent,
  rightContent,
  defaultLeftWidth = 240,
  minLeftWidth = 180,
}: SplitPanelProps) {
  return (
    <ResizeBox.Split
      direction="horizontal"
      defaultSize={defaultLeftWidth}
      min={minLeftWidth}
      max={-400}
      style={{ height: '100%' }}
      panes={[
        <div key="left" className="h-full overflow-auto border-r border-gray-3">
          {leftContent}
        </div>,
        <div key="right" className="h-full overflow-auto">
          {rightContent}
        </div>,
      ]}
    />
  )
}
```

### StatsCard — 统计卡片

> UX 规范：§5.3 统计卡片

```tsx
// src/components/business/StatsCard.tsx
import { Card, Statistic } from '@arco-design/web-react'
import { TrendingUp, TrendingDown } from 'lucide-react'

interface StatsCardProps {
  title: string
  value: number | string
  trend?: number
  suffix?: string
  prefix?: React.ReactNode
  icon?: React.ReactNode
  valueStyle?: React.CSSProperties
}

export function StatsCard({
  title,
  value,
  trend,
  suffix,
  prefix,
  icon,
  valueStyle,
}: StatsCardProps) {
  return (
    <Card className="hover:shadow-sm transition-shadow duration-100" bordered>
      <div className="flex items-start justify-between">
        <div>
          <div className="text-xs text-gray-6 mb-1">{title}</div>
          <Statistic
            value={value}
            suffix={suffix}
            prefix={prefix}
            valueStyle={
              valueStyle || { fontSize: 24, fontWeight: 600, color: '#272E3B' }
            }
          />
          {trend !== undefined && (
            <div
              className={`text-xs mt-1 flex items-center gap-1 ${trend >= 0 ? 'text-green-6' : 'text-red-6'}`}
            >
              {trend >= 0 ? (
                <TrendingUp size={12} />
              ) : (
                <TrendingDown size={12} />
              )}
              <span>{Math.abs(trend)}% 较上月</span>
            </div>
          )}
        </div>
        {icon && <div className="text-gray-4">{icon}</div>}
      </div>
    </Card>
  )
}
```

### ReferencePanel — 引用来源面板

> UX 规范：§9.4 右侧引用来源区

```tsx
// src/components/business/ReferencePanel.tsx
import { Card, Badge } from '@arco-design/web-react'
import { FileText } from 'lucide-react'
import type { ReferencedChunkDetail } from '@/types/api'

interface ReferencePanelProps {
  /** 引用来源列表（需包含 content 字段，由页面组件从 chunks API 获取后传入） */
  references: ReferencedChunkDetail[]
  onJumpToSource?: (chunkId: string, documentId: string) => void
}

export function ReferencePanel({
  references,
  onJumpToSource,
}: ReferencePanelProps) {
  if (references.length === 0) {
    return (
      <div className="text-center py-8 text-gray-6 text-sm">
        此草稿未引用知识库内容
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center gap-2 mb-4">
        <h3 className="text-base font-medium text-gray-9">引用来源</h3>
        <Badge
          count={references.length}
          style={{ backgroundColor: '#7B61FF' }}
        />
      </div>
      <div className="space-y-3">
        {references.map((ref) => (
          <Card key={ref.chunkId} className="border border-gray-3" size="small">
            <div className="flex items-center gap-2 mb-2">
              <FileText size={14} className="text-primary" />
              <span
                className="text-sm text-primary cursor-pointer hover:underline"
                onClick={() => onJumpToSource?.(ref.chunkId, ref.documentId)}
              >
                {ref.documentTitle}
              </span>
            </div>
            <div className="text-xs text-gray-6 mb-2">
              相似度: {Math.round(ref.similarityScore * 100)}%
            </div>
            {ref.content ? (
              <div className="text-sm text-gray-7 bg-gray-1 p-3 rounded leading-relaxed line-clamp-5">
                {ref.content}
              </div>
            ) : (
              <div className="text-sm text-gray-5 bg-gray-1 p-3 rounded italic">
                内容加载中...
              </div>
            )}
            {onJumpToSource && (
              <div className="mt-2">
                <span
                  className="text-xs text-primary cursor-pointer hover:underline"
                  onClick={() => onJumpToSource(ref.chunkId, ref.documentId)}
                >
                  查看原文
                </span>
              </div>
            )}
          </Card>
        ))}
      </div>
    </div>
  )
}
```

### CaseSelector — 用例选择器

> UX 规范：§11.2 新建计划页的用例选择面板

> **宪法 §4.1 合规**: CaseSelector 为纯展示组件，可选用例数据由父组件通过 React Query hook 获取后以 props 传入。筛选条件通过回调上抛。

```tsx
// src/components/business/CaseSelector.tsx
import {
  Table,
  Tabs,
  Select,
  Checkbox,
  Button,
  Badge,
} from '@arco-design/web-react'
import { StatusTag } from './StatusTag'
import type { TestCase, CaseType, Priority } from '@/types/api'

interface CaseSelectorProps {
  /** 可选用例列表，由父组件通过 React Query hook 获取 */
  availableCases: TestCase[]
  /** 加载状态 */
  loading: boolean
  /** 已选用例变更回调 */
  onChange: (selected: TestCase[]) => void
  /** 当前已选中的用例列表 */
  value: TestCase[]
  /** 筛选变更回调，父组件用于更新 React Query 参数 */
  onFilterChange?: (filters: {
    module_id?: string
    case_type?: CaseType
    priority?: Priority
  }) => void
}

export function CaseSelector({
  availableCases,
  loading,
  onChange,
  value,
  onFilterChange,
}: CaseSelectorProps) {
  const selectedIds = new Set(value.map((c) => c.id))

  const toggleSelect = (testCase: TestCase) => {
    if (selectedIds.has(testCase.id)) {
      onChange(value.filter((c) => c.id !== testCase.id))
    } else {
      onChange([...value, testCase])
    }
  }

  const removeAll = () => onChange([])

  const columns = [
    {
      title: '编号',
      dataIndex: 'number',
      width: 180,
      render: (val: string) => <span className="font-mono text-sm">{val}</span>,
    },
    { title: '标题', dataIndex: 'title', ellipsis: true },
    {
      title: '类型',
      dataIndex: 'caseType',
      width: 100,
      render: (val: string) => <StatusTag type="caseType" value={val} />,
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      width: 80,
      render: (val: string) => <StatusTag type="priority" value={val} />,
    },
  ]

  return (
    <Tabs defaultActiveTab="available">
      <Tabs.TabPane key="available" title="可选用例">
        <div className="flex gap-2 mb-3">
          <Select
            placeholder="模块"
            style={{ width: 120 }}
            allowClear
            onChange={(val) => onFilterChange?.({ module_id: val })}
          />
          <Select
            placeholder="类型"
            style={{ width: 120 }}
            allowClear
            onChange={(val) => onFilterChange?.({ case_type: val })}
          />
          <Select
            placeholder="优先级"
            style={{ width: 120 }}
            allowClear
            onChange={(val) => onFilterChange?.({ priority: val })}
          />
          <span className="text-sm text-gray-6 self-center ml-2">
            已选择 {value.length} 条
          </span>
          {value.length > 0 && (
            <Button type="text" size="small" onClick={removeAll}>
              清空选择
            </Button>
          )}
        </div>
        <Table
          columns={[
            {
              title: '',
              width: 48,
              render: (_: unknown, record: TestCase) => (
                <Checkbox
                  checked={selectedIds.has(record.id)}
                  onChange={() => toggleSelect(record)}
                />
              ),
            },
            ...columns,
          ]}
          data={availableCases}
          loading={loading}
          rowKey="id"
          size="small"
          pagination={false}
        />
      </Tabs.TabPane>
      <Tabs.TabPane
        key="selected"
        title={
          <span>
            已选用例{' '}
            <Badge
              count={value.length}
              style={{ backgroundColor: '#7B61FF' }}
            />
          </span>
        }
      >
        {value.length === 0 ? (
          <div className="text-center py-8 text-gray-5 text-sm">
            未选择任何用例
          </div>
        ) : (
          <Table
            columns={[
              ...columns,
              {
                title: '操作',
                width: 80,
                render: (_: unknown, record: TestCase) => (
                  <Button
                    type="text"
                    size="small"
                    status="danger"
                    onClick={() => toggleSelect(record)}
                  >
                    移除
                  </Button>
                ),
              },
            ]}
            data={value}
            rowKey="id"
            size="small"
            pagination={false}
          />
        )}
        {value.length > 0 && (
          <div className="mt-2">
            <Button type="text" size="small" onClick={removeAll}>
              全部移除
            </Button>
          </div>
        )}
      </Tabs.TabPane>
    </Tabs>
  )
}
```

### JsonEditor — JSON 编辑器

> UX 规范：§7.5 配置管理页

```tsx
// src/components/business/JsonEditor.tsx
import { useState, useEffect } from 'react'
import { Input, Message } from '@arco-design/web-react'

const { TextArea } = Input

interface JsonEditorProps {
  value: unknown
  onChange: (value: unknown, isValid: boolean) => void
  height?: number
}

export function JsonEditor({ value, onChange, height = 200 }: JsonEditorProps) {
  const [text, setText] = useState(() => JSON.stringify(value, null, 2))
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setText(JSON.stringify(value, null, 2))
    setError(null)
  }, [value])

  const handleChange = (newText: string) => {
    setText(newText)
    try {
      const parsed = JSON.parse(newText)
      setError(null)
      onChange(parsed, true)
    } catch (e) {
      setError((e as Error).message)
      onChange(null, false)
    }
  }

  return (
    <div>
      <TextArea
        value={text}
        onChange={handleChange}
        style={{
          height,
          fontFamily: 'var(--font-mono, monospace)',
          fontSize: 13,
        }}
        status={error ? 'error' : undefined}
      />
      {error && (
        <div className="text-red-500 text-xs mt-1">JSON 格式错误: {error}</div>
      )}
    </div>
  )
}
```

---

# 第六章：页面详细设计

> **数据获取模式**: 所有页面数据通过 React Query hooks 获取（见 §3.4），禁止 `useEffect` + API 调用。表单使用 React Hook Form + Zod schema 校验（见 §3.4 表单模式）。

## 6.0 Zod Schema 定义（表单校验）

> **宪法 §3.3**: 表单输入使用 Zod schema 校验。以下为各页面的 schema 定义。

### src/features/auth/schema/loginSchema.ts

```typescript
import { z } from 'zod'

export const loginSchema = z.object({
  email: z.string().min(1, '请输入邮箱').email('请输入有效的邮箱地址'),
  password: z.string().min(1, '请输入密码'),
})

export type LoginInput = z.infer<typeof loginSchema>
```

### src/features/auth/schema/registerSchema.ts

```typescript
import { z } from 'zod'

export const registerSchema = z.object({
  username: z
    .string()
    .min(3, '用户名至少 3 个字符')
    .max(32, '用户名最多 32 个字符'),
  email: z.string().min(1, '请输入邮箱').email('请输入有效的邮箱地址'),
  password: z
    .string()
    .min(8, '密码至少 8 位字符')
    .max(100, '密码最多 100 个字符'),
  // 注意：super_admin 由系统分配，注册页面不开放此角色
  role: z.enum(['admin', 'normal'], { errorMap: () => '请选择用户角色' }),
})

export type RegisterInput = z.infer<typeof registerSchema>
```

### src/features/projects/schema/projectSchema.ts

```typescript
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

### 表单使用模式

```tsx
// 在组件中使用 React Hook Form + Zod
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  createProjectSchema,
  type CreateProjectInput,
} from '../schema/projectSchema'

function CreateProjectModal() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateProjectInput>({
    resolver: zodResolver(createProjectSchema),
    mode: 'onChange',
  })

  const createProject = useCreateProject() // React Query mutation hook

  const onSubmit = (data: CreateProjectInput) => {
    createProject.mutate(data, {
      onSuccess: () => Message.success('项目创建成功'),
    })
  }

  // ... JSX
}
```

### Schema 清单

| 页面         | Schema 文件                                          | 字段                                                 |
| ------------ | ---------------------------------------------------- | ---------------------------------------------------- |
| 登录         | `features/auth/schema/loginSchema.ts`                | email, password                                      |
| 注册         | `features/auth/schema/registerSchema.ts`             | username, email, password, role                      |
| 创建项目     | `features/projects/schema/projectSchema.ts`          | name, prefix, description                            |
| 新建用例     | 内联或 `features/testcases/schema/caseSchema.ts`     | moduleId, title, steps, expected, caseType, priority |
| 新建计划     | 内联或 `features/plans/schema/planSchema.ts`         | name, description                                    |
| 新建生成任务 | 内联或 `features/generation/schema/taskSchema.ts`    | moduleId, prompt, count, caseType, priority          |
| 上传文档     | 内联或 `features/documents/schema/documentSchema.ts` | name, type                                           |

## 6.1 登录页 (/login)

> UX 规范：§6.1

**路由**：`/login`
**布局**：AuthLayout（左右分栏 55%:45%）
**认证**：公开

### API

| 时机     | 端点          | 方法 | 说明              |
| -------- | ------------- | ---- | ----------------- |
| 点击登录 | `/auth/login` | POST | 返回 token + user |

### 组件树

```
AuthLayout
├── LeftPanel（品牌展示区）
│   ├── 品牌插图
│   └── 标语
└── RightPanel（登录表单区）
    ├── Logo + 产品名
    ├── Form
    │   ├── Input（邮箱，必填，email 格式校验）
    │   ├── Input.Password（密码，必填）
    │   ├── Checkbox（记住我）
    │   └── Button（登录，primary，loading 态）
    └── Link → /register
```

### 交互逻辑

1. 填写邮箱密码 → 点击登录 → Button loading
2. 成功 → `useAuthStore.login()` 存储 token → 跳转 `/projects`
3. 失败 401 → 顶部 Alert（error）"邮箱或密码错误"
4. 已登录访问 `/login` → 重定向 `/projects`

### 本地状态

| 状态    | 类型         | 说明           |
| ------- | ------------ | -------------- |
| loading | boolean      | 按钮加载态     |
| form    | FormInstance | Arco Form 实例 |

## 6.2 注册页 (/register)

> UX 规范：§6.2

**路由**：`/register`
**布局**：AuthLayout
**认证**：公开

### API

| 时机     | 端点             | 方法 |
| -------- | ---------------- | ---- |
| 点击注册 | `/auth/register` | POST |

### 表单字段

| 字段     | 组件           | 必填 | 验证           |
| -------- | -------------- | ---- | -------------- |
| 用户名   | Input          | 是   | 3-32 字符      |
| 邮箱     | Input          | 是   | email 格式     |
| 密码     | Input.Password | 是   | ≥8 字符        |
| 确认密码 | Input.Password | 是   | 与密码一致     |
| 角色     | Select         | 是   | admin / normal |

### 交互逻辑

1. 角色选项不含 super_admin（由系统分配）
2. 409 错误 → 字段级错误（"邮箱已存在"或"用户名已存在"）
3. 成功 → Message.success("注册成功") → 跳转 `/login`

## 6.3 项目列表页 (/projects)

> UX 规范：§7.1

**路由**：`/projects`
**布局**：AppLayout
**认证**：需认证

### API

| 时机     | 端点                     | 方法 |
| -------- | ------------------------ | ---- |
| 页面加载 | `/projects`              | GET  |
| 搜索     | `/projects?keywords=xxx` | GET  |
| 创建项目 | `/projects`              | POST |

### 页面布局

```
┌──────────────────────────────────────────────────┐
│  🔍 搜索项目名称                    [+ 新建项目]   │
├──────────────────────────────────────────────────┤
│  ┌──────┐  ┌──────┐  ┌──────┐                    │
│  │ Card │  │ Card │  │ Card │  ... (3列网格)      │
│  └──────┘  └──────┘  └──────┘                    │
└──────────────────────────────────────────────────┘
```

### 组件树

```
ProjectListPage
├── 操作栏（Input 搜索 + Button 新建项目）
├── Card 网格 (grid grid-cols-3 gap-4)
│   └── Card × N
│       ├── h3 项目名称 + Tag 前缀(monospace)
│       ├── caption 描述(2行截断)
│       ├── caption 统计行（📄 数量 📋 数量）
│       └── Button 组（进入项目 / 设置）
└── CreateProjectModal
    ├── Input（项目名称，必填，2-255 字符）
    ├── Input（项目前缀，必填，2-4 位大写字母）
    ├── TextArea（项目描述，可选）
    └── Button（取消 + 确认创建）
```

### 交互逻辑

- 搜索：Input onChange 防抖 300ms → 重新请求
- 新建项目：点击 → Modal → 前缀实时校验格式 + 失焦校验唯一
- 409 → 字段错误"项目名称已存在"或"项目前缀已存在"
- 空状态：`FolderOpen` 图标 + "创建第一个项目" 按钮

## 6.4 项目仪表盘 (/projects/:id)

> UX 规范：§7.3

**路由**：`/projects/:id`
**布局**：AppLayout

### API

| 时机     | 端点                   | 方法 |
| -------- | ---------------------- | ---- |
| 页面加载 | `/projects/{id}`       | GET  |
| 页面加载 | `/projects/{id}/stats` | GET  |

### 组件树

```
ProjectDashboard
├── 标题行（h1 项目名称 + caption 描述 + 操作按钮组）
│   └── Button × 3（上传文档 / 发起生成 / 新建用例，outlined）
├── StatsCard 网格 (grid grid-cols-4 gap-4)
│   ├── StatsCard（用例总数）
│   ├── StatsCard（通过率，后缀 %）
│   ├── StatsCard（覆盖率，后缀 %）
│   └── StatsCard（AI 生成数）
├── 底部双栏 (grid grid-cols-5 gap-4)
│   ├── 趋势图 (col-span-3) — 折线图
│   └── 最近任务列表 (col-span-2)
│       └── List × N（状态 Tag + 摘要 + 时间）
└── 新项目引导（当 case_count=0 时显示）
    └── 快速开始卡片（3 步引导）
```

### 交互逻辑

- 统计卡片数据来自 `ProjectStatistics`
- 新项目（0 用例/0 文档）显示引导卡片替代趋势图
- 趋势图使用 ECharts 或 Recharts，品牌色 `#7B61FF`

## 6.5 知识库文档列表 (/projects/:id/knowledge)

> UX 规范：§8.1

**路由**：`/projects/:id/knowledge`

### API

| 时机     | 端点                                  | 方法   |
| -------- | ------------------------------------- | ------ |
| 页面加载 | `/knowledge/documents?project_id=xxx` | GET    |
| 筛选     | 同上 + type/status 参数               | GET    |
| 上传文档 | `/knowledge/documents`                | POST   |
| 删除文档 | `/knowledge/documents/{id}`           | DELETE |

### 组件树

```
KnowledgeListPage
├── SearchTable
│   ├── filters: [文档类型, 状态]
│   ├── searchPlaceholder: "搜索文档名称"
│   └── columns: [文档名称, 类型Tag, 状态Tag, 分块数, 上传时间, 操作]
│       └── 操作: 查看 / 删除(Popconfirm)
└── UploadDocumentModal
    ├── Input（文档名称，必填）
    ├── Select（文档类型，必填）
    ├── Upload（拖拽上传区）
    └── Button（取消 + 确认上传）
```

### 表格列

| 列名     | 宽度  | 说明                                            |
| -------- | ----- | ----------------------------------------------- |
| 文档名称 | 弹性  | 可点击跳转详情                                  |
| 类型     | 100px | StatusTag(documentType)                         |
| 状态     | 100px | StatusTag(documentStatus)，processing 显示 Spin |
| 分块数   | 80px  | 整数                                            |
| 上传时间 | 160px | YYYY-MM-DD HH:mm                                |
| 操作     | 120px | 查看 / 删除                                     |

## 6.6 文档详情页 (/projects/:id/knowledge/:docId)

> UX 规范：§8.3

### API

| 时机     | 端点                               | 方法 |
| -------- | ---------------------------------- | ---- |
| 页面加载 | `/knowledge/documents/{id}`        | GET  |
| 页面加载 | `/knowledge/documents/{id}/chunks` | GET  |

### 布局：SplitPanel（左 300px 信息面板 + 右侧分块列表）

**左侧面板**：

- 文档名称(h2) + 类型 Tag + 状态 Tag
- 元信息：上传人、上传时间、分块数
- Arco Steps：pending → processing → completed
- 失败时：Alert(error) + "重新解析"按钮

**右侧面板**：

- 标题"文档分块" + 数量 Badge
- Arco List：序号 + 内容预览（3行截断）+ 展开按钮

## 6.7 Figma 集成页 (/projects/:id/knowledge/figma)

> UX 规范：§8.4

**全页面表单，分 3 个区域**：

| 区域     | 内容                                                          |
| -------- | ------------------------------------------------------------- |
| 连接配置 | Radio（认证方式）+ Input.Password（令牌）+ Button（测试连接） |
| 导入文件 | Input（Figma URL）+ Button（解析）                            |
| 节点选择 | Tree（带 Checkbox），全选/取消，确认导入                      |

## 6.8 AI 生成任务列表 (/projects/:id/generation)

> UX 规范：§9.1

### API

| 时机     | 端点                               | 方法 |
| -------- | ---------------------------------- | ---- |
| 页面加载 | `/generation/tasks?project_id=xxx` | GET  |
| 筛选     | 同上 + status 参数                 | GET  |

### 组件树

```
GenerationTaskListPage
├── SearchTable
│   ├── filters: [状态]
│   ├── searchPlaceholder: "搜索需求描述"
│   ├── toolbar: [✨ 新建生成任务] Button(primary)
│   └── columns: [任务ID, 需求描述, 用例数量, 状态Tag, 创建时间, 操作]
│       └── 操作: 详情 / 重试(仅 failed)
└── 空状态: Sparkles 图标 + "暂无生成任务"
```

### 表格列

| 列名     | 宽度  | 说明                  |
| -------- | ----- | --------------------- |
| 任务 ID  | 120px | 截断 + 复制按钮       |
| 需求描述 | 弹性  | 截断 50 字 + Tooltip  |
| 用例数量 | 80px  | 整数                  |
| 状态     | 100px | StatusTag(taskStatus) |
| 创建时间 | 160px | YYYY-MM-DD HH:mm      |
| 操作     | 120px | 详情 / 重试           |

## 6.9 新建生成任务 (/projects/:id/generation/new)

> UX 规范：§9.2

### API

| 时机     | 端点                     | 方法 |
| -------- | ------------------------ | ---- |
| 页面加载 | `/projects/{id}/modules` | GET  |
| 提交     | `/generation/tasks`      | POST |

### 表单字段

| 区域     | 字段     | 组件             | 必填 | 说明                     |
| -------- | -------- | ---------------- | ---- | ------------------------ |
| 基本信息 | 目标模块 | Select           | 是   | 数据来自 modules API     |
| 基本信息 | 需求描述 | TextArea(rows=5) | 是   | ≥10 字符，显示字数       |
| 基本信息 | 用例数量 | InputNumber      | 否   | 1-20，默认 5             |
| 文档范围 | 范围     | Radio.Group      | 否   | 全部文档/仅PRD/仅Figma   |
| 高级选项 | 场景类型 | Checkbox.Group   | 否   | 正向/异常/边界，默认全选 |
| 高级选项 | 优先级   | Select           | 否   | P0-P3，默认 P2           |
| 高级选项 | 用例类型 | Select           | 否   | 枚举值                   |
| 高级选项 | 生成模式 | Radio.Group      | 否   | 常规/深度覆盖            |

### 知识库检查逻辑

| 条件       | 显示                              | 操作按钮状态   |
| ---------- | --------------------------------- | -------------- |
| 项目无文档 | Alert(warning) "请先上传需求文档" | 禁用"立即生成" |
| 文档不足   | Alert(warning) "生成质量可能较低" | 正常可用       |

## 6.10 生成任务详情 (/projects/:id/generation/:taskId)

> UX 规范：§9.3

### API

| 时机     | 端点                            | 方法 |
| -------- | ------------------------------- | ---- | ------------------- |
| 页面加载 | `/generation/tasks/{id}`        | GET  |
| 页面加载 | `/generation/tasks/{id}/drafts` | GET  |
| 轮询(5s) | `/generation/tasks/{id}`        | GET  | status≠completed 时 |

### 状态映射

| 任务状态   | 显示                                  |
| ---------- | ------------------------------------- |
| pending    | Spin + "任务排队中..."                |
| processing | Progress + "正在生成..." + 每 5s 轮询 |
| completed  | "生成完成" Badge + 草稿列表           |
| failed     | Alert(error) + 错误信息 + "重试"按钮  |

### 草稿列表

- 批量工具栏：全选 Checkbox + "批量确认"(primary) + "批量拒绝"(default)
- 表格列：Checkbox / 序号 / 标题 / 类型 Tag / 优先级 Tag / 置信度 Tag / 引用来源 / 操作
- 点击行 → 草稿确认页
- 批量确认流程：选中 → Modal 选择模块 → POST batch-confirm

## 6.11 草稿确认页 (/drafts/:draftId) ⭐核心页面

> UX 规范：§9.4（核心交互页，需重点打磨）

**路由**：`/drafts/:draftId`
**布局**：AppLayout（无侧边栏模块树）

### API

| 时机     | 端点                              | 方法                            |
| -------- | --------------------------------- | ------------------------------- |
| 页面加载 | `/generation/drafts/{id}`         | GET（草稿详情，**需后端补充**） |
| 保存修改 | `PUT /generation/drafts/{id}`     | PUT（**需后端补充**，见 §10.6） |
| 确认     | `/generation/drafts/{id}/confirm` | POST                            |
| 拒绝     | `/generation/drafts/{id}/reject`  | POST                            |

### 草稿编辑保存机制

**保存策略**：使用 `PUT /generation/drafts/{id}` 将编辑内容持久化到后端。

**保存时机**：

1. 用户点击"保存修改"按钮 → 显式保存
2. 用户切换到另一条草稿（圆点导航 / 键盘 ← →）→ 自动保存当前编辑
3. 用户离开草稿确认页 → `useEffect` cleanup 中自动保存

**降级方案**：若后端 `PUT /generation/drafts/{id}` 暂未实现，前端使用 `sessionStorage` 作为临时存储：

- key: `draft_edit_{draftId}`
- value: 编辑中的 JSON 序列化数据
- 页面加载时先检查 sessionStorage，有未保存数据则恢复并提示"有未保存的编辑内容"

### 页面布局

```
┌──────────────────────────────────────────────────────────┐
│  ← 返回草稿列表      第 2 / 5 条                [●●◉○○] │  ← 草稿导航栏
├───────────────────────────┬──────────────────────────────┤
│       草稿编辑区 (60%)     │      引用来源区 (40%)         │
│                           │                              │
│  标题: [Input          ] │  📄 登录模块 PRD              │
│                           │  相似度: 92%                  │
│  前置条件:                 │  "引用内容预览..."            │
│  [ArrayEditor           ] │  [查看原文]                   │
│                           │                              │
│  测试步骤:                 │  📄 安全规范文档              │
│  [ArrayEditor(min=1)    ] │  相似度: 78%                  │
│                           │  "引用内容预览..."            │
│  预期结果:                 │                              │
│  [TextArea              ] │                              │
│                           │                              │
│  类型:[Select] 优先级:[Sel] │                              │
├───────────────────────────┴──────────────────────────────┤
│  [拒绝] [保存修改] [✅ 确认并转为正式用例]                   │  ← 固定底部
└──────────────────────────────────────────────────────────┘
```

### 组件树

```
DraftConfirmPage
├── DraftNavBar（返回 + 进度 + 圆点导航）
├── SplitPanel（60:40）
│   ├── 左侧编辑区
│   │   ├── Form
│   │   │   ├── Input（标题）
│   │   │   ├── ArrayEditor（前置条件，minItems=0）
│   │   │   ├── ArrayEditor（测试步骤，minItems=1）
│   │   │   ├── TextArea（预期结果，rows=4）
│   │   │   ├── Select（用例类型）
│   │   │   └── Select（优先级）
│   │   └——
│   └── 右侧引用区
│       └── ReferencePanel
├── 底部操作栏（sticky bottom）
│   ├── Button（拒绝，danger）
│   ├── Button（保存修改，default）
│   └── Button（确认并转为正式用例，primary）
├── RejectModal
│   ├── Radio.Group（拒绝原因：重复/无关/质量低/其他）
│   ├── TextArea（详细反馈，选填）
│   └── Button（取消 + 确认拒绝）
└── ConfirmModal
    ├── Select（目标模块，必填）
    └── Button（取消 + 确认）
```

### 草稿间导航交互

| 元素               | 行为                            |
| ------------------ | ------------------------------- |
| "← 返回草稿列表"   | 导航回 `/drafts`                |
| "第 N / M 条"      | 文字进度指示                    |
| 圆点导航 `[●●◉○○]` | 点击跳转对应草稿                |
| 键盘 ← →           | 切换上/下一条（未聚焦输入框时） |
| 切换前             | 自动保存当前编辑内容            |

### 交互流程

1. 确认：点击"确认" → Modal 选择模块 → POST confirm → 成功 → Message + 跳转用例详情
2. 拒绝：点击"拒绝" → Modal 填原因 → POST reject → 回到草稿列表
3. 保存：暂存编辑内容到前端 state，不调用转正 API

## 6.12 草稿箱全局视图 (/drafts)

> UX 规范：§9.5

**路由**：`/drafts`（全局，跨项目）

### API

| 时机     | 端点                                               | 方法 |
| -------- | -------------------------------------------------- | ---- |
| 页面加载 | `/drafts?status=pending`                           | GET  |
| 筛选     | `/drafts?project_id=&module_id=&status=&keywords=` | GET  |

### SearchTable 筛选

| 筛选项 | 组件   | 选项                            |
| ------ | ------ | ------------------------------- |
| 项目   | Select | 全部 / 按项目                   |
| 模块   | Select | 随项目联动                      |
| 状态   | Select | 全部 / 待处理 / 已确认 / 已拒绝 |

## 6.13 测试用例库 (/projects/:id/cases)

> UX 规范：§10.1

**路由**：`/projects/:id/cases`

### API

| 时机     | 端点                                                             | 方法   |
| -------- | ---------------------------------------------------------------- | ------ |
| 页面加载 | `/projects/{id}/modules`                                         | GET    |
| 页面加载 | `/testcases?project_id=xxx`                                      | GET    |
| 筛选     | `/testcases?project_id=&module_id=&status=&case_type=&priority=` | GET    |
| 新建用例 | `/testcases`                                                     | POST   |
| 删除用例 | `/testcases/{id}`                                                | DELETE |

### 布局：SplitPanel（左 240px 模块树 + 右侧表格）

**左侧模块树**：

- "全部"选项 + Arco Tree（模块名 + 用例数 Badge）
- 点击节点 → 右侧表格按 module_id 筛选

**右侧表格列**：

| 列名     | 宽度  | 说明                      |
| -------- | ----- | ------------------------- |
| Checkbox | 48px  | 批量选择                  |
| 编号     | 200px | monospace                 |
| 标题     | 弹性  | ellipsis + Tooltip        |
| 用例类型 | 100px | StatusTag                 |
| 优先级   | 80px  | StatusTag                 |
| 状态     | 80px  | StatusTag                 |
| 创建人   | 80px  | 用户名                    |
| 更新时间 | 160px | YYYY-MM-DD HH:mm          |
| 操作     | 100px | 查看 / 编辑 / 复制 / 删除 |

**批量操作栏**（选中行后浮现）：

- 修改优先级 / 修改状态 / 加入计划 / 删除(danger)

### 空状态

`FileCheck` 图标 + "暂无测试用例" + 两个按钮："手动创建"(primary) + "使用 AI 生成"(default)

## 6.14 新建用例 (Drawer)

> UX 规范：§10.2

**触发**：点击"新建用例" → Drawer（宽 640px）

### 表单字段

| 字段     | 组件               | 必填 | 验证       |
| -------- | ------------------ | ---- | ---------- |
| 目标模块 | Select             | 是   | 非空       |
| 标题     | Input              | 是   | 2-500 字符 |
| 前置条件 | ArrayEditor        | 否   | —          |
| 测试步骤 | ArrayEditor(min=1) | 是   | ≥1 项      |
| 预期结果 | TextArea           | 是   | 非空       |
| 用例类型 | Select             | 是   | 枚举值     |
| 优先级   | Select             | 是   | P0-P3      |

### API

`POST /testcases` → `TestCase`

## 6.15 用例详情页 (/projects/:id/cases/:caseId)

> UX 规范：§10.3

### API

| 时机     | 端点              | 方法 |
| -------- | ----------------- | ---- |
| 页面加载 | `/testcases/{id}` | GET  |

### 页面布局

```
┌──────────────────────────────────────────────────────┐
│  ← 返回用例库                                         │
│  ECO-USR-20260402-001  密码错误超过5次锁定  [✅ 通过]   │
│                                    [编辑][复制][删除]  │
├──────────────────────────────────────────────────────┤
│  ┌─ 基本信息 ─────────────────────────────────────┐  │
│  │  类型: [功能测试]  优先级: [P1]  模块: 登录模块  │  │
│  │  创建人: 张三  创建时间: 04-02  更新: 04-03     │  │
│  └────────────────────────────────────────────────┘  │
│  ┌─ 用例内容 ─────────────────────────────────────┐  │
│  │  前置条件: 1. xxx  2. xxx                      │  │
│  │  测试步骤: ① xxx  ② xxx  ③ xxx               │  │
│  │  预期结果: 灰色背景文本块                        │  │
│  └────────────────────────────────────────────────┘  │
│  ▶ AI 生成信息（Collapse，仅 AI 生成用例显示）         │
│  ▶ 执行历史（表格：计划名称/结果/执行人/时间/备注）     │
└──────────────────────────────────────────────────────┘
```

### AI 元数据区（Collapse）

展开内容：生成任务（可点击跳转）+ 置信度 StatusTag + 引用文档块列表 + 模型版本 + 生成时间

**特殊状态**：

- 源文档已变更 → Alert(warning) "源文档已变更，请核实"
- 源文档已删除 → Alert(warning) "源文档已删除" + "移除引用"按钮

## 6.16 测试计划列表 (/projects/:id/plans)

> UX 规范：§11.1

### API

| 时机     | 端点                    | 方法 |
| -------- | ----------------------- | ---- |
| 页面加载 | `/plans?project_id=xxx` | GET  |
| 筛选     | 同上 + status/keywords  | GET  |

### 表格列

| 列名     | 说明                  |
| -------- | --------------------- |
| 计划名称 | 可点击跳转详情        |
| 状态     | StatusTag(planStatus) |
| 用例数   | 关联用例总数          |
| 通过率   | 百分比 + 迷你进度条   |
| 创建人   | 用户名                |
| 创建时间 | YYYY-MM-DD HH:mm      |
| 操作     | 查看 / 编辑 / 删除    |

## 6.17 新建计划 (/projects/:id/plans/new)

> UX 规范：§11.2

**布局**：SplitPanel（左侧表单 + 右侧用例选择面板）

### API

| 时机     | 端点                        | 方法 |
| -------- | --------------------------- | ---- |
| 页面加载 | `/testcases?project_id=xxx` | GET  |
| 提交     | `/plans`                    | POST |

### 用例选择面板

使用 Tab 切换两个视图：

- **可选用例 Tab**：筛选栏（模块/类型/优先级）+ 表格（Checkbox + 编号 + 标题 + 类型 + 优先级）
- **已选用例 Tab**：已选列表 + 移除按钮 + "全部移除"

## 6.18 计划详情页 (/projects/:id/plans/:planId)

> UX 规范：§11.3

### API

| 时机     | 端点                  | 方法  |
| -------- | --------------------- | ----- |
| 页面加载 | `/plans/{id}`         | GET   |
| 状态变更 | `/plans/{id}/status`  | PATCH |
| 录入结果 | `/plans/{id}/results` | POST  |

### 操作按钮（按状态变化）

| 状态      | 可用操作                                |
| --------- | --------------------------------------- |
| draft     | 编辑 + 开始执行(primary) + 删除(danger) |
| active    | 编辑 + 标记完成(primary)                |
| completed | 重新执行 + 归档(primary)                |
| archived  | 取消归档                                |

### 统计卡片（5 列）

总用例 / 通过(绿) / 失败(红) / 阻塞(橙) / 跳过(灰) + 未执行数

### 执行结果录入

**快捷录入**：点击执行结果列 → 内联 Select → 自动提交 → Message("已录入：通过") + 3s 撤销链接

**详细录入**：点击"录入"按钮 → Modal(500px) → Radio.Group(通过/失败/阻塞/跳过) + TextArea(备注)

## 6.19 模块管理 (/projects/:id/settings/modules)

> UX 规范：§7.4

**布局**：SplitPanel（左 280px 模块树 + 右侧编辑区）

### API

| 时机     | 端点                                | 方法   |
| -------- | ----------------------------------- | ------ |
| 页面加载 | `/projects/{id}/modules`            | GET    |
| 创建     | `/projects/{id}/modules`            | POST   |
| 编辑     | `/projects/{id}/modules/{moduleId}` | PUT    |
| 删除     | `/modules/{id}`                     | DELETE |

**左侧模块树**：模块名 + 缩写 Tag + 悬浮操作（编辑/删除/拖拽排序）
**右侧编辑区**：名称/缩写/描述/用例数量 + 编辑表单 + 删除按钮(Popconfirm)

### 级联删除

删除模块 → Modal 确认："删除模块将级联删除其下所有用例，确认删除？"

## 6.20 配置管理 (/projects/:id/settings/configs)

> UX 规范：§7.5

### API

| 时机      | 端点                            | 方法   |
| --------- | ------------------------------- | ------ |
| 页面加载  | `/projects/{id}/configs`        | GET    |
| 新增/编辑 | `/projects/{id}/configs/{key}`  | PUT    |
| 删除      | `/projects/{id}/configs/{key}`  | DELETE |
| 导入      | `/projects/{id}/configs/import` | POST   |
| 导出      | `/projects/{id}/configs/export` | GET    |

**表格**：键名 / 值(JSON预览) / 描述 / 更新时间 / 操作(编辑/删除)

---

# 第七章：设计令牌实现

## 7.1 Tailwind CSS v4 配置

### src/styles/theme.css

```css
@import 'tailwindcss';

@theme {
  /* ===== 色彩系统 ===== */

  /* 品牌色 */
  --color-primary: #7b61ff;
  --color-primary-hover: #6b4fe0;
  --color-primary-active: #5a3dc0;
  --color-primary-light: rgba(123, 97, 255, 0.1);

  /* 功能色 */
  --color-success: #00b42a;
  --color-success-hover: #009a29;
  --color-success-light: rgba(0, 180, 42, 0.1);
  --color-warning: #ff7d00;
  --color-warning-hover: #e66f00;
  --color-warning-light: rgba(255, 125, 0, 0.1);
  --color-error: #f53f3f;
  --color-error-hover: #cb2634;
  --color-error-light: rgba(245, 63, 63, 0.1);
  --color-info: #86909c;
  --color-info-light: rgba(134, 144, 156, 0.1);

  /* 中性色阶 */
  --color-gray-1: #f7f8fa;
  --color-gray-2: #f2f3f5;
  --color-gray-3: #e5e6eb;
  --color-gray-4: #c9cdd4;
  --color-gray-5: #a9aeb8;
  --color-gray-6: #86909c;
  --color-gray-7: #6b7785;
  --color-gray-8: #4e5969;
  --color-gray-9: #272e3b;
  --color-gray-10: #1d2129;

  /* ===== 字体 ===== */
  --font-sans:
    system-ui, -apple-system, 'PingFang SC', 'Microsoft YaHei', sans-serif;
  --font-mono: 'JetBrains Mono', 'Fira Code', Consolas, ui-monospace, monospace;

  /* 字号 */
  --text-h1: 24px;
  --text-h2: 20px;
  --text-h3: 16px;
  --text-body: 14px;
  --text-body-sm: 13px;
  --text-caption: 12px;
  --text-micro: 11px;

  /* ===== 间距 ===== */
  /* 使用 Tailwind v4 标准 spacing scale 扩展，自动映射为 p-xs/m-xs 等工具类 */
  --spacing-0-5: 2px; /* 0.125rem → p-0.5 */
  --spacing-1: 4px; /* xs */
  --spacing-1-5: 6px;
  --spacing-2: 8px; /* sm */
  --spacing-3: 12px; /* md */
  --spacing-4: 16px; /* base */
  --spacing-5: 20px; /* lg */
  --spacing-6: 24px; /* xl */
  --spacing-8: 32px; /* 2xl */
  --spacing-10: 40px; /* 3xl */
  --spacing-12: 48px; /* 4xl */
  --spacing-16: 64px; /* 5xl */

  /* ===== 圆角 ===== */
  --radius-small: 2px;
  --radius-medium: 4px;
  --radius-large: 8px;
  --radius-xl: 12px;
  --radius-2xl: 16px;

  /* ===== 阴影 ===== */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.07);
  --shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1);
  --shadow-xl: 0 20px 25px rgba(0, 0, 0, 0.1);

  /* ===== 动效时长 ===== */
  --duration-micro: 100ms;
  --duration-fast: 200ms;
  --duration-normal: 300ms;
  --duration-slow: 500ms;

  /* ===== 缓动 ===== */
  --ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);
  --ease-out: cubic-bezier(0, 0, 0.2, 1);
  --ease-in: cubic-bezier(0.4, 0, 1, 1);
}
```

## 7.2 Arco ConfigProvider 主题

### src/App.tsx

```tsx
import { ConfigProvider } from '@arco-design/web-react'
import { RouterProvider } from 'react-router-dom'
import { arcoTheme } from '@/styles/arco-theme'
import { router } from '@/router'
import zhCN from '@arco-design/web-react/es/locale/zh-CN'

export default function App() {
  return (
    <ConfigProvider locale={zhCN} theme={arcoTheme}>
      <RouterProvider router={router} />
    </ConfigProvider>
  )
}
```

## 7.3 CSS 变量（语义色）

以上 Tailwind @theme 已定义完整色值，组件中可直接使用 `text-primary`、`bg-gray-1`、`border-gray-3` 等工具类。需要额外定义的语义变量：

```css
/* src/styles/global.css */
:root {
  /* 侧边栏 */
  --sidebar-width: 220px;
  --sidebar-collapsed-width: 64px;
  --header-height: 48px;

  /* 内容区 */
  --content-padding: 24px;
  --content-min-width: 1024px;

  /* 面板 */
  --panel-split-width: 2px;
  --panel-split-hover-width: 4px;
}
```

---

# 第八章：错误处理与反馈

## 8.1 全局错误边界

### src/components/ErrorBoundary.tsx

```tsx
import { Component, type ReactNode } from 'react'
import { Button, Result } from '@arco-design/web-react'

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error }
  }

  render() {
    if (this.state.hasError) {
      return (
        <Result
          status="error"
          title="页面出现错误"
          subTitle={this.state.error?.message || '请刷新页面重试'}
          extra={
            <Button type="primary" onClick={() => window.location.reload()}>
              刷新页面
            </Button>
          }
        />
      )
    }
    return this.props.children
  }
}
```

### 使用方式

在 `App.tsx` 中包裹 `RouterProvider`：

```tsx
<ErrorBoundary>
  <RouterProvider router={router} />
</ErrorBoundary>
```

## 8.2 API 错误码映射

> UX 规范：§12.3

| HTTP 状态码 | 处理方式     | UI 反馈                        | 处理位置       |
| ----------- | ------------ | ------------------------------ | -------------- |
| 400         | 表单验证提示 | 对应字段下方红色错误提示       | 页面组件 catch |
| 401         | 跳转登录页   | 清除 token + redirect `/login` | Axios 拦截器   |
| 403         | 无权限提示   | Message.warning("无操作权限")  | Axios 拦截器   |
| 404         | 404 页面     | 专属 404 页 + "返回首页"按钮   | 路由兜底       |
| 409         | 冲突提示     | 表单字段错误（"名称已存在"）   | 页面组件 catch |
| 500         | 服务器错误   | Message.error("服务器异常")    | Axios 拦截器   |

### 表单级错误处理模式

```tsx
// 通用表单提交错误处理
// 需要导入: import { isAxiosError } from 'axios'
const handleSubmit = async (values: FormData) => {
  try {
    setLoading(true)
    await someApi.create(values)
    Message.success('创建成功')
    onClose()
  } catch (err: unknown) {
    if (isAxiosError(err) && err.response?.status === 400) {
      // 后端返回字段级错误
      const msg = (err.response.data as { error?: string }).error
      form.setFields({
        name: { value: values.name, errors: [msg ?? '请求参数错误'] },
      })
    } else if (isAxiosError(err) && err.response?.status === 409) {
      Message.warning('资源已存在')
    }
  } finally {
    setLoading(false)
  }
}
```

## 8.3 加载状态与反馈

> UX 规范：§12.1, §12.4

### 加载状态层级

| 级别     | 组件                         | 使用场景          |
| -------- | ---------------------------- | ----------------- |
| 页面级   | `Spin`（全页居中）           | 首次加载页面      |
| 组件级   | `Table.loading` / `Skeleton` | 表格/列表刷新     |
| 按钮级   | `Button.loading`             | 提交/保存操作     |
| 长时任务 | `Progress` + 轮询            | AI 生成、文档解析 |

### Message 使用规范

| 操作   | 类型    | 文案模板               |
| ------ | ------- | ---------------------- |
| 创建   | success | "{资源}创建成功"       |
| 更新   | success | "{资源}更新成功"       |
| 删除   | success | "{资源}已删除"         |
| 确认   | success | "草稿已确认为正式用例" |
| 导入   | success | "导入成功，共 N 条"    |
| 导出   | success | "导出成功"             |
| 无权限 | warning | "无操作权限"           |
| 失败   | error   | 具体错误信息           |

### Notification 使用规范

| 场景         | 标题              | 内容              | 操作       |
| ------------ | ----------------- | ----------------- | ---------- |
| AI 生成完成  | "AI 用例生成完成" | 产出 N 条草稿     | "查看"链接 |
| 文档解析完成 | "文档解析完成"    | 文档名 + N 个分块 | "查看"链接 |
| 文档解析失败 | "文档解析失败"    | 文档名 + 错误原因 | "重试"链接 |

## 8.4 安全编码规范

### XSS 防护

| 规则                               | 说明                                                                                  |
| ---------------------------------- | ------------------------------------------------------------------------------------- |
| **禁止 `dangerouslySetInnerHTML`** | 全项目禁止使用，任何富文本渲染须通过 DOMPurify 消毒                                   |
| **用户输入渲染**                   | 所有用户输入（用例标题、步骤、文档名称等）通过 React JSX 自动转义，不做额外 HTML 解析 |
| **URL 跳转白名单**                 | 外部链接跳转须校验协议头（仅允许 `http://` / `https://`），禁止 `javascript:` 协议    |
| **JSON 配置展示**                  | `JsonEditor` 组件必须使用 `<pre><code>` 文本渲染，禁止将用户 JSON 作为 HTML 注入      |

### 敏感数据处理

| 数据类型      | 存储方式     | 说明                                                       |
| ------------- | ------------ | ---------------------------------------------------------- |
| JWT Token     | localStorage | v1.0 使用 localStorage，后续版本建议迁移至 httpOnly Cookie |
| Refresh Token | localStorage | 同上                                                       |
| Figma Token   | 不在前端存储 | 通过后端代理加密存储，前端仅触发后端 API                   |
| 用户密码      | 不在前端存储 | 登录后立即丢弃，不缓存在任何前端状态                       |

### CSRF 防护

- 所有修改类请求（POST/PUT/PATCH/DELETE）携带 JWT Token 在 `Authorization` Header 中
- JWT 放在 Header 而非 Cookie 中，天然防御 CSRF 攻击
- 后续若迁移至 Cookie 方案，需配合 `SameSite=Strict` + CSRF Token

---

# 第九章：性能与优化

## 9.1 路由懒加载

所有页面组件使用 React Router `lazy` 路由属性实现按需加载，已在第二章路由配置中实现。

```tsx
// 路由级 lazy loading（无需 React.lazy + Suspense）
{
  path: 'projects/:id/cases',
  lazy: () => import('@/features/testcases/components/CaseListPage').then((m) => ({
    Component: m.CaseListPage,
  })),
}
```

## 9.2 组件按需加载

Arco Design 通过 `@arco-plugins/vite-react` 插件实现组件级按需加载（已在 vite.config.ts 中配置 `style: 'css'`）。

## 9.3 大数据量表格优化

当用例数超过 500 条时，考虑使用虚拟滚动：

```tsx
import { VirtualList } from '@arco-design/web-react'
;<Table virtual scroll={{ y: 600 }} columns={columns} data={data} />
```

## 9.4 API 请求缓存

> **宪法 §4.1**: 服务端数据必须通过 React Query 管理，禁止 Zustand 缓存。

对于不频繁变化的数据（如模块列表、项目配置），使用 React Query 的 `staleTime` 控制缓存有效期：

```tsx
// 在 hooks 中设置 staleTime，5 分钟内不重新请求
export function useModuleList(projectId: string) {
  return useQuery({
    queryKey: moduleKeys.list(projectId),
    queryFn: () => modulesApi.list(projectId),
    enabled: !!projectId,
    staleTime: 5 * 60 * 1000, // 5 分钟
  })
}
```

## 9.5 轮询优化

> **宪法 §4.1**: 使用 React Query `refetchInterval` 替代自定义轮询 hook（见 §3.4 轮询模式）。**禁止自行实现 `usePolling` hook。**

AI 生成任务、文档解析等长时任务均使用 `refetchInterval` + 条件停止模式：

```tsx
// 示例：AI 生成任务轮询
export function useTaskDetail(id: string) {
  return useQuery({
    queryKey: generationKeys.task(id),
    queryFn: () => generationApi.getTask(id),
    enabled: !!id,
    refetchInterval: (query) => {
      const status = query.state.data?.status
      return status === 'pending' || status === 'processing' ? 5000 : false
    },
  })
}
```

## 9.6 防抖搜索

列表页搜索输入使用防抖，避免频繁 API 调用。统一使用 `@/hooks/useDebounce.ts`：

```tsx
// src/hooks/useDebounce.ts
import { useState, useEffect } from 'react'

export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState(value)
  useEffect(() => {
    const timer = setTimeout(() => setDebouncedValue(value), delay)
    return () => clearTimeout(timer)
  }, [value, delay])
  return debouncedValue
}
```

页面使用示例：

```tsx
// 配合 SearchTable 受控模式 + React Query
const [keywords, setKeywords] = useState('')
const debouncedKeywords = useDebounce(keywords, 300)

const { data, isFetching } = useProjectList({ keywords: debouncedKeywords })

<SearchTable
  data={data?.data ?? []}
  loading={isFetching}
  total={data?.total ?? 0}
  keywords={keywords}
  onKeywordsChange={setKeywords}
  // ...
/>
```

---

# 第十章：缺失 API 与开发建议

## 10.1 缺失 API 清单

> UX 规范：§14.7

以下接口在 openapi.yaml 中未定义，需后端补充：

| #   | 端点                                | 方法  | 优先级 | 说明                                                |
| --- | ----------------------------------- | ----- | ------ | --------------------------------------------------- |
| 1   | `/generation/tasks`                 | GET   | P0     | 生成任务列表（分页、按项目+状态筛选）               |
| 2   | `/drafts`                           | GET   | P0     | 全局草稿列表（跨项目，分页+筛选）                   |
| 3   | `/projects/{id}/modules/{moduleId}` | PUT   | P0     | 模块编辑（名称、缩写、描述）                        |
| 4   | `/plans/{id}/status`                | PATCH | P1     | 计划状态变更（draft→active→completed→archived）     |
| 5   | `/testcases/import`                 | POST  | P1     | 用例批量导入（上传 xlsx/csv）                       |
| 6   | `/testcases/export`                 | GET   | P1     | 用例导出（返回文件流）                              |
| 7   | `/generation/drafts/{id}`           | GET   | P0     | 草稿详情（草稿确认页依赖）                          |
| 8   | `/generation/drafts/{id}`           | PUT   | P1     | 草稿编辑保存（草稿确认页"保存修改"依赖）            |
| 9   | `/auth/me`                          | GET   | P1     | 当前用户信息验证（RouteGuard token 有效性校验依赖） |

## 10.2 前端 Mock 策略

开发阶段使用 Mock Service Worker (MSW) 模拟后端 API：

```bash
npm install msw --save-dev
npx msw init public/ --save
```

### Mock 配置示例

```typescript
// src/mocks/handlers.ts
import { http, HttpResponse } from 'msw'

export const handlers = [
  // 登录
  http.post('/api/v1/auth/login', async ({ request }) => {
    const body = (await request.json()) as { email: string; password: string }
    return HttpResponse.json({
      access_token: 'mock-token-123',
      refresh_token: 'mock-refresh-123',
      user: {
        id: '1',
        username: body.email.split('@')[0],
        email: body.email,
        role: 'admin',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    })
  }),

  // 项目列表
  http.get('/api/v1/projects', () => {
    return HttpResponse.json({
      data: [
        {
          id: '1',
          name: 'ECommerce',
          prefix: 'ECO',
          description: '电商平台测试项目',
          createdAt: '2026-04-01T00:00:00Z',
          updatedAt: '2026-04-02T00:00:00Z',
        },
      ],
      total: 1,
      offset: 0,
      limit: 20,
    })
  }),
]
```

## 10.3 开发优先级建议

### P0（MVP 核心，优先开发）

| 序号 | 页面/功能             | 前端工作量 | 依赖 API          |
| ---- | --------------------- | ---------- | ----------------- |
| 1    | 认证流程（登录/注册） | 1 天       | auth              |
| 2    | 全局布局 + 路由       | 1 天       | —                 |
| 3    | 项目列表 + 创建       | 1 天       | projects          |
| 4    | 项目仪表盘            | 2 天       | projects/stats    |
| 5    | 模块管理              | 1 天       | modules           |
| 6    | 测试用例库 + CRUD     | 3 天       | testcases         |
| 7    | 知识库文档列表 + 上传 | 2 天       | documents         |
| 8    | AI 生成任务 + 新建    | 2 天       | generation/tasks  |
| 9    | 草稿确认页（核心）    | 3 天       | generation/drafts |

**P0 预计总工时：16 天**

### P1（完善功能）

| 序号 | 页面/功能                  | 前端工作量 |
| ---- | -------------------------- | ---------- |
| 1    | 测试计划列表 + 新建 + 详情 | 3 天       |
| 2    | 执行结果录入（快捷+详细）  | 2 天       |
| 3    | 项目配置管理               | 1 天       |
| 4    | 用例导入/导出              | 1 天       |
| 5    | 文档详情页 + Figma 集成    | 2 天       |
| 6    | 草稿箱全局视图             | 1 天       |

**P1 预计总工时：10 天**

### P2（增强体验）

| 序号 | 页面/功能          |
| ---- | ------------------ |
| 1    | 全局搜索（Ctrl+K） |
| 2    | 表格列自定义       |
| 3    | 批量编辑统一规范   |
| 4    | 快捷键系统         |
| 5    | 通知面板           |
| 6    | 趋势图表优化       |

## 10.4 测试策略

### 单元测试

| 工具       | Vitest + React Testing Library                               |
| ---------- | ------------------------------------------------------------ |
| 覆盖范围   | Store 逻辑、API 模块、工具函数、StatusTag/ArrayEditor 等组件 |
| 运行时机   | git pre-commit hook + CI pipeline                            |
| 覆盖率目标 | 核心业务逻辑 > 80%                                           |

**优先编写单测的模块**：

1. `request.ts` — Token 刷新逻辑（含并发刷新场景）
2. `StatusTag` — 色彩映射完整性
3. `ArrayEditor` — 增删排序行为
4. `useAuthStore` — 登录/登出/token 刷新
5. `useAppStore` — 通知状态管理

### E2E 测试

| 工具     | Playwright                   |
| -------- | ---------------------------- |
| 覆盖范围 | 核心业务流程（3 条关键路径） |
| 运行时机 | PR 合并前 CI                 |

**核心 E2E 场景**：

1. **登录 → 创建项目 → 创建模块 → 上传文档 → AI 生成 → 确认草稿 → 查看用例**
2. **创建计划 → 选择用例 → 开始执行 → 录入结果 → 查看统计**
3. **搜索用例 → 编辑 → 复制 → 删除 → 验证列表更新**

## 10.5 PRD 功能点追溯矩阵

| PRD 模块 | PRD 功能     | 优先级 | UX 章节  | 本文档章节 |
| -------- | ------------ | ------ | -------- | ---------- |
| 项目管理 | 创建项目     | P0     | §7.1-7.2 | §6.3       |
| 项目管理 | 编辑项目     | P0     | §7.2     | §6.3       |
| 项目管理 | 删除项目     | P1     | §12.7    | §8.2       |
| 项目管理 | 项目列表     | P0     | §7.1     | §6.3       |
| 项目管理 | 创建模块     | P0     | §7.4     | §6.19      |
| 项目管理 | 编辑模块     | P0     | §7.4     | §6.19      |
| 项目管理 | 删除模块     | P1     | §12.7    | §6.19      |
| 项目管理 | 项目配置     | P2     | §7.5     | §6.20      |
| 知识库   | PRD 上传解析 | P0     | §8.1-8.3 | §6.5-6.6   |
| 知识库   | Figma 集成   | P0     | §8.4     | §6.7       |
| 知识库   | API 规范导入 | P1     | §8.2     | §6.5       |
| 知识库   | 文档更新策略 | P1     | §4.3     | §6.6       |
| AI 生成  | 发起生成任务 | P0     | §9.2     | §6.9       |
| AI 生成  | RAG 检索流程 | P0     | §9.2     | §6.9       |
| AI 生成  | 草稿确认     | P0     | §9.4     | §6.11      |
| AI 生成  | 草稿拒绝     | P0     | §9.4     | §6.11      |
| AI 生成  | 批量确认     | P0     | §9.3     | §6.10      |
| 用例管理 | 创建用例     | P0     | §10.2    | §6.14      |
| 用例管理 | 编辑用例     | P0     | §10.3    | §6.15      |
| 用例管理 | 删除用例     | P1     | §12.7    | §6.13      |
| 用例管理 | 用例详情     | P0     | §10.3    | §6.15      |
| 用例管理 | 批量导入     | P1     | §10.1    | §6.13      |
| 用例管理 | 导出用例     | P1     | §10.1    | §6.13      |
| 测试执行 | 创建计划     | P1     | §11.2    | §6.17      |
| 测试执行 | 关联用例     | P1     | §11.2    | §6.17      |
| 测试执行 | 结果录入     | P1     | §11.4    | §6.18      |
| 测试执行 | 执行历史     | P1     | §10.3    | §6.15      |

---

> **文档结束** — 本文档覆盖 Aitestos 智能测试管理平台前端完整详细设计，包括项目架构、路由、状态管理、API 层、组件实现、页面设计、设计令牌、错误处理、性能优化和开发建议。
