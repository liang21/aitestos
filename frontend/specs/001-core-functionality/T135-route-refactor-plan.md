# T135 路由重构详细计划

## 1. 概述

### 目标
将当前路由结构对齐到 plan.md §3.1 规定的完整路由清单。

### 当前问题
1. **结构不一致** - 业务路由未嵌套在 `/projects/:projectId` 下
2. **命名不统一** - `/testcases` vs `/cases`，`/documents` vs `/knowledge`
3. **权限守卫缺失** - admin 路由未配置 `requireAdmin`
4. **组件缺失** - FigmaIntegrationPage 未实现

---

## 2. 路由对比表

### 规格要求 vs 当前实现

| 功能 | plan.md 规格 | 当前实际 | 偏差类型 |
|------|-------------|---------|---------|
| **认证** | | | |
| 登录 | `/login` | `/login` | ✅ 匹配 |
| 注册 | `/register` | `/register` | ✅ 匹配 |
| **项目管理** | | | |
| 项目列表 | `/projects` | `/projects` | ✅ 匹配 |
| 项目仪表盘 | `/projects/:projectId` | `/projects/:projectId` | ✅ 匹配 |
| **知识库** | | | |
| 文档列表 | `/projects/:projectId/knowledge` | `/documents` | 🔴 结构错误 |
| 文档详情 | `/projects/:projectId/knowledge/:docId` | `/documents/:documentId` | 🔴 结构错误 |
| Figma 集成 | `/projects/:projectId/knowledge/figima` | **缺失** | 🔴 未实现 |
| **AI 生成** | | | |
| 任务列表 | `/projects/:projectId/generation` | `/generation` | 🔴 缺少 projectId |
| 新建任务 | `/projects/:projectId/generation/new` | `/generation/tasks/new` | 🔴 路径错误 |
| 任务详情 | `/projects/:projectId/generation/:taskId` | `/generation/tasks/:taskId` | 🔴 路径错误 |
| **用例管理** | | | |
| 用例列表 | `/projects/:projectId/cases` | `/testcases` | 🔴 结构错误 |
| 用例详情 | `/projects/:projectId/cases/:caseId` | `/testcases/:caseId` | 🔴 结构错误 |
| **测试计划** | | | |
| 计划列表 | `/projects/:projectId/plans` | `/plans` | 🔴 缺少 projectId |
| 新建计划 | `/projects/:projectId/plans/new` | `/plans/new` | 🔴 缺少 projectId |
| 计划详情 | `/projects/:projectId/plans/:planId` | `/plans/:planId` | 🔴 缺少 projectId |
| **草稿箱** | | | |
| 草稿列表 | `/drafts` | `/drafts` | ✅ 匹配（全局） |
| 草稿确认 | `/drafts/:draftId` | `/drafts/:draftId` | ✅ 匹配（全局） |
| **配置管理** | | | |
| 配置管理 | `/projects/:projectId/configs` | `/projects/:projectId/configs` | ✅ 匹配 |
| **模块管理** | | | |
| 模块管理 | `/projects/:projectId/settings/modules` | `/projects/:projectId/modules` | 🟡 缺少 /settings |

---

## 3. 影响范围分析

### 3.1 受影响的导航链接（13 处）

| 文件 | 当前链接 | 目标链接 | 优先级 |
|------|---------|---------|--------|
| `Sidebar.tsx` | `/testcases` | 动态生成 | P0 |
| `Sidebar.tsx` | `/plans` | 动态生成 | P0 |
| `Sidebar.tsx` | `/generation` | 动态生成 | P0 |
| `Sidebar.tsx` | `/documents` | 动态生成 | P0 |
| `Sidebar.tsx` | `/drafts` | `/drafts` | P0 |
| `ProjectDashboard.tsx` | 操作按钮 | 动态生成 | P0 |
| `CreateCaseDrawer.tsx` | 提交后跳转 | 动态生成 | P1 |
| `DraftConfirmPage.tsx` | 确认后跳转 | 动态生成 | P1 |

### 3.2 受影响的组件（35 个文件使用路由 hooks）

高影响文件（直接使用 useParams）：
- `CaseDetailPage.tsx` - 使用 `caseId` 参数
- `DocumentDetailPage.tsx` - 使用 `documentId` 参数
- `TaskDetailPage.tsx` - 使用 `taskId` 参数
- `DraftConfirmPage.tsx` - 使用 `draftId` 参数
- `PlanDetailPage.tsx` - 使用 `planId` 参数
- `ProjectDashboard.tsx` - 使用 `projectId` 参数
- `ModuleManagePage.tsx` - 使用 `projectId` 参数
- `ConfigManagePage.tsx` - 使用 `projectId` 参数
- `KnowledgeListPage.tsx` - 需要添加 `projectId` 参数

### 3.3 受影响的 hooks（6 个）

| Hook | 当前用途 | 需要修改 |
|------|---------|---------|
| `useProjects.ts` | 项目列表查询 | ✅ 无需修改 |
| `useTestCases.ts` | 用例查询 | 🔴 需要传递 projectId |
| `useDocuments.ts` | 文档查询 | 🔴 需要传递 projectId |
| `useGeneration.ts` | 生成任务查询 | 🔴 需要传递 projectId |
| `usePlans.ts` | 计划查询 | 🔴 需要传递 projectId |
| `useDrafts.ts` | 草稿查询 | ✅ 全局查询 |

---

## 4. 重构策略

### 4.1 阶段划分（4 个阶段）

#### 阶段 1：基础设施（P0，无破坏性变更）
- [ ] 创建 `src/lib/routes.ts` - 路径常量工厂
- [ ] 创建 `src/hooks/useProjectRoutes.ts` - 动态路由生成 hook
- [ ] 更新 `src/router/index.tsx` - 添加新路由结构
- [ ] 添加旧路径重定向（向后兼容）

#### 阶段 2：核心组件迁移（P0）
- [ ] 更新 `Sidebar.tsx` - 使用动态路由
- [ ] 更新 `Header.tsx` - 使用动态路由
- [ ] 更新 `ProjectDashboard.tsx` - 使用动态路由
- [ ] 更新 `AppLayout.tsx` - 传递 projectId 上下文

#### 阶段 3：业务组件迁移（P1）
- [ ] 更新所有 useParams 调用（10 个组件）
- [ ] 更新所有 hooks 查询（5 个 hooks）
- [ ] 更新所有 navigate 调用（13 处）

#### 阶段 4：缺失功能实现（P1）
- [ ] 创建 `FigmaIntegrationPage.tsx`
- [ ] 添加权限守卫（requireAdmin）
- [ ] 移除旧路径重定向（可选，保留向后兼容）

---

## 5. 详细实施步骤

### 5.1 阶段 1：基础设施

#### 5.1.1 创建路径常量工厂 (`src/lib/routes.ts`)

```typescript
// 动态生成项目相关路由
export const buildProjectRoutes = (projectId: string) => ({
  dashboard: `/projects/${projectId}`,
  knowledge: {
    list: `/projects/${projectId}/knowledge`,
    detail: (docId: string) => `/projects/${projectId}/knowledge/${docId}`,
    figma: `/projects/${projectId}/knowledge/figma`,
  },
  generation: {
    list: `/projects/${projectId}/generation`,
    new: `/projects/${projectId}/generation/new`,
    detail: (taskId: string) => `/projects/${projectId}/generation/${taskId}`,
  },
  cases: {
    list: `/projects/${projectId}/cases`,
    detail: (caseId: string) => `/projects/${projectId}/cases/${caseId}`,
  },
  plans: {
    list: `/projects/${projectId}/plans`,
    new: `/projects/${projectId}/plans/new`,
    detail: (planId: string) => `/projects/${projectId}/plans/${planId}`,
  },
  settings: {
    modules: `/projects/${projectId}/settings/modules`,
    configs: `/projects/${projectId}/configs`,
  },
})

// 全局路由
export const GLOBAL_ROUTES = {
  projects: '/projects',
  drafts: '/drafts',
  draftDetail: (draftId: string) => `/drafts/${draftId}`,
  login: '/login',
  register: '/register',
} as const
```

#### 5.1.2 创建动态路由 Hook (`src/hooks/useProjectRoutes.ts`)

```typescript
import { useParams } from 'react-router-dom'
import { buildProjectRoutes } from '@/lib/routes'

export function useProjectRoutes() {
  const { projectId } = useParams<{ projectId: string }>()
  if (!projectId) throw new Error('useProjectRoutes must be used within a project route')
  return buildProjectRoutes(projectId)
}
```

#### 5.1.3 更新路由配置 (`src/router/index.tsx`)

添加新路由结构（保留旧路由用于向后兼容）：

```typescript
// 新路由结构
{
  path: 'projects/:projectId',
  element: <RouteGuard><AppLayout /></RouteGuard>,
  children: [
    { index: true, element: <Navigate to="dashboard" replace /> },
    { path: 'dashboard', lazy: () => import(...) },
    {
      path: 'knowledge',
      children: [
        { index: true, lazy: () => import(...) },
        { path: ':docId', lazy: () => import(...) },
        { path: 'figma', element: <RouteGuard requireAdmin><FigmaIntegrationPage /></RouteGuard> }
      ]
    },
    { path: 'generation', ... },
    { path: 'generation/new', ... },
    { path: 'generation/:taskId', ... },
    { path: 'cases', ... },
    { path: 'cases/:caseId', ... },
    { path: 'plans', ... },
    { path: 'plans/new', ... },
    { path: 'plans/:planId', ... },
    {
      path: 'settings/modules',
      element: <RouteGuard requireAdmin><ModuleManagePage /></RouteGuard>
    },
    {
      path: 'configs',
      element: <RouteGuard requireAdmin><ConfigManagePage /></RouteGuard>
    },
  ]
}

// 旧路径重定向（向后兼容）
{
  path: 'testcases',
  element: <Navigate to="/projects/:projectId/cases" replace />
},
{
  path: 'documents',
  element: <Navigate to="/projects/:projectId/knowledge" replace />
},
// ... 其他重定向
```

---

### 5.2 阶段 2：核心组件迁移

#### 5.2.1 更新 Sidebar.tsx

```typescript
// 修改前
<NavLink to="/testcases">

// 修改后
const { projectId } = useParams() // 从上下文获取
<NavLink to={`/projects/${projectId}/cases`}>
```

**问题**：Sidebar 需要知道当前 projectId。

**解决方案**：
1. 选项 A：从 URL 获取（当前在 `/projects/:projectId` 下有效）
2. 选项 B：从全局状态获取（useCurrentProject store）
3. 选项 C：仅显示项目列表，子页面隐藏 Sidebar

**推荐**：选项 A + 选项 C 组合

#### 5.2.2 更新 AppLayout.tsx

```typescript
// 添加 projectId 上下文
interface AppLayoutProps {
  projectId?: string
  // ...
}

export function AppLayout({ projectId, ... }: AppLayoutProps) {
  // 通过 Context 或 state 传递 projectId
}
```

---

### 5.3 阶段 3：业务组件迁移

#### 5.3.1 更新 useParams 调用

```typescript
// CaseDetailPage.tsx
// 修改前
const { caseId } = useParams<{ caseId: string }>()

// 修改后
const { projectId, caseId } = useParams<{ projectId: string; caseId: string }>()
```

#### 5.3.2 更新 hooks 查询

```typescript
// useTestCases.ts
// 修改前
export function useCaseList(params: { status?: string }) {
  return useQuery({
    queryKey: ['testcases', params],
    queryFn: () => testcasesApi.list({ project_id: '???, ...params }) // 缺少 projectId
  })
}

// 修改后
export function useCaseList(projectId: string, params: { status?: string }) {
  return useQuery({
    queryKey: ['testcases', projectId, params],
    queryFn: () => testcasesApi.list({ project_id: projectId, ...params })
  })
}
```

#### 5.3.3 更新所有调用点

```typescript
// CaseListPage.tsx
// 修改前
const { data } = useCaseList({ status: 'pass' })

// 修改后
const { projectId } = useParams()
const { data } = useCaseList(projectId, { status: 'pass' })
```

---

### 5.4 阶段 4：缺失功能实现

#### 5.4.1 创建 FigmaIntegrationPage.tsx

规格：plan.md §5.4
- 路由：`/projects/:projectId/knowledge/figma`
- 权限：admin+
- 功能：连接配置、导入文件、节点选择

#### 5.4.2 添加权限守卫

```typescript
// 模块管理
<RouteGuard requireAdmin>
  <ModuleManagePage />
</RouteGuard>

// Figma 集成
<RouteGuard requireAdmin>
  <FigmaIntegrationPage />
</RouteGuard>
```

---

## 6. 风险评估

### 6.1 高风险区域

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 现有功能破坏 | 🔴 高 | 保留旧路径重定向（向后兼容） |
| hooks 签名变更 | 🔴 高 | 使用 TypeScript 严格检查，编译时发现所有调用点 |
| projectId 丢失 | 🟡 中 | 添加运行时检查，抛出友好错误 |
| 测试失败 | 🟡 中 | 分阶段更新测试，每阶段验证 |

### 6.2 回滚策略

1. **Git 分支保护** - 在新分支执行重构
2. **渐进式合并** - 每个阶段独立 commit
3. **功能开关** - 可通过环境变量切换新旧路由
4. **回滚命令** - `git revert <commit-range>`

---

## 7. 验收标准

### 7.1 功能验收

- [ ] 所有 19 条 plan.md 路由可访问
- [ ] 权限守卫正确（admin 路由）
- [ ] 懒加载正常工作
- [ ] 旧路径重定向有效
- [ ] 404 页面兜底

### 7.2 质量验收

- [ ] TypeScript 编译无错误
- [ ] 所有测试通过
- [ ] 无 `any` 类型
- [ ] ESLint 无警告

### 7.3 用户体验验收

- [ ] 所有导航链接工作正常
- [ ] 浏览器后退/前进按钮正常
- [ ] 书签/直接访问 URL 正常
- [ ] 页面刷新保持状态

---

## 8. 时间估算

| 阶段 | 预计时间 | 依赖 |
|------|---------|------|
| 阶段 1：基础设施 | 1 小时 | 无 |
| 阶段 2：核心组件 | 1.5 小时 | 阶段 1 |
| 阶段 3：业务组件 | 3 小时 | 阶段 2 |
| 阶段 4：缺失功能 | 1.5 小时 | 阶段 1 |
| 测试与修复 | 1 小时 | 全部 |
| **总计** | **8 小时** | |

---

## 9. 待审核问题

### 9.1 需要决策

1. **向后兼容策略**
   - A) 永久保留旧路径重定向
   - B) 临时保留（1-2 个版本后移除）
   - C) 不保留，一次性迁移

   **推荐**：A（用户友好，技术成本低）

2. **projectId 传递方式**
   - A) 仅从 URL 获取
   - B) 全局状态管理
   - C) 混合方式

   **推荐**：A（简单、符合 RESTful 规范）

3. **重构时机**
   - A) 立即执行
   - B) 等待下一个版本
   - C) 分阶段发布

   **推荐**：A（越晚成本越高）

### 9.2 技术债务确认

当前是否接受：
- ✅ 路径命名不一致（testcases vs cases）
- ✅ 结构不统一（有/无 projectId）
- ✅ 缺失功能（Figma 集成）

---

## 10. 附录

### 10.1 完整路由清单（19 条）

```
/login                                  → LoginPage
/register                               → RegisterPage
/projects                               → ProjectListPage
/projects/:projectId                    → ProjectDashboard
/projects/:projectId/knowledge          → KnowledgeListPage
/projects/:projectId/knowledge/:docId   → DocumentDetailPage
/projects/:projectId/knowledge/figma    → FigmaIntegrationPage [admin]
/projects/:projectId/generation         → GenerationTaskListPage
/projects/:projectId/generation/new     → NewGenerationTaskPage
/projects/:projectId/generation/:taskId → TaskDetailPage
/projects/:projectId/cases              → CaseListPage
/projects/:projectId/cases/:caseId      → CaseDetailPage
/projects/:projectId/plans              → PlanListPage
/projects/:projectId/plans/new          → NewPlanPage
/projects/:projectId/plans/:planId      → PlanDetailPage
/projects/:projectId/settings/modules   → ModuleManagePage [admin]
/projects/:projectId/configs            → ConfigManagePage [admin]
/drafts                                 → DraftListPage
/drafts/:draftId                        → DraftConfirmPage
```

### 10.2 受影响文件清单（38 个）

```
src/router/index.tsx
src/lib/routes.ts (新建)
src/hooks/useProjectRoutes.ts (新建)
src/components/layout/Sidebar.tsx
src/components/layout/Header.tsx
src/components/layout/AppLayout.tsx
src/features/auth/components/LoginPage.tsx
src/features/configs/components/ConfigManagePage.tsx
src/features/documents/components/DocumentDetailPage.tsx
src/features/documents/components/KnowledgeListPage.tsx
src/features/documents/components/FigmaIntegrationPage.tsx (新建)
src/features/drafts/components/DraftConfirmPage.tsx
src/features/drafts/components/DraftListPage.tsx
src/features/generation/components/GenerationTaskListPage.tsx
src/features/generation/components/NewGenerationTaskPage.tsx
src/features/generation/components/TaskDetailPage.tsx
src/features/generation/hooks/useGeneration.ts
src/features/modules/components/ModuleManagePage.tsx
src/features/plans/components/NewPlanPage.tsx
src/features/plans/components/PlanDetailPage.tsx
src/features/plans/components/PlanListPage.tsx
src/features/plans/hooks/usePlans.ts
src/features/projects/components/ProjectDashboard.tsx
src/features/projects/hooks/useProjects.ts
src/features/testcases/components/CaseDetailPage.tsx
src/features/testcases/components/CaseListPage.tsx
src/features/testcases/components/CreateCaseDrawer.tsx
src/features/testcases/hooks/useTestCases.ts
src/components/NotFoundPage.tsx
```
