# Aitestos Frontend Spec v2.0 Alignment Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Align the existing frontend implementation with spec v2.0 — fix routing, close page-level gaps, add missing pages, and ensure all acceptance tests pass.

**Architecture:** Feature-Based React 19 + TypeScript app using Arco Design, React Query for server state, Zustand for UI state, React Router 7 with lazy loading. The codebase already has all feature modules scaffolded with services/hooks/components — this plan closes the quality and spec-compliance gaps.

**Tech Stack:** React 19, TypeScript 5.9, Arco Design 2.66, TanStack React Query 5.x, Zustand 5.x, React Router 7.x, Vitest 4.x, MSW 2.x, Tailwind CSS 4.x, Zod 4.x

---

## Scope Check

This spec covers the full platform. The codebase already has ~80 files with substantial implementations. This plan focuses on **closing the gap** between current state and spec v2.0, organized into 6 phases by dependency order:

1. **Foundation fixes** — routing, layout, sidebar (blocks everything else)
2. **Shared component polish** — StatusTag color maps, missing features
3. **Auth & projects** — login/register polish, project list as cards, dashboard
4. **Knowledge & generation** — documents, generation tasks, drafts (core flow)
5. **Test cases & plans** — case management, plan execution
6. **Settings & configs** — module management, config page

---

## File Structure

### Files to Create
- `src/features/documents/components/FigmaIntegrationPage.tsx` — Figma import page
- `src/features/auth/schema/registerSchema.ts` — Zod schema for registration
- `src/features/testcases/schema/createCaseSchema.ts` — Zod schema for case creation
- `src/features/generation/schema/createTaskSchema.ts` — Zod schema for generation task
- `src/features/plans/schema/createPlanSchema.ts` — Zod schema for plan creation
- `src/styles/animations.css` — AI-specific animations (glow-pulse, ai-reveal, shimmer)

### Files to Modify (Major)
- `src/router/index.tsx` — Restructure to nested project-scoped routes
- `src/components/layout/Sidebar.tsx` — Add project submenu, draft badge, AI icon styling
- `src/components/layout/Header.tsx` — Add breadcrumb with project context
- `src/components/layout/AppLayout.tsx` — Pass project context to sidebar/header
- `src/components/business/StatusTag.tsx` — Add all color maps from UX spec §2.1
- `src/features/projects/components/ProjectListPage.tsx` — Table → Card grid
- `src/features/projects/components/ProjectDashboard.tsx` — Add trend chart, recent tasks, quick-start
- `src/features/drafts/components/DraftConfirmPage.tsx` — Add draft navigation, reference panel, AI styling
- `src/features/drafts/components/DraftListPage.tsx` — Add cross-project filters
- `src/features/generation/components/GenerationTaskListPage.tsx` — AI visual styling
- `src/features/generation/components/NewGenerationTaskPage.tsx` — Knowledge readiness indicator
- `src/features/generation/components/TaskDetailPage.tsx` — AI status card, batch operations
- `src/features/plans/components/PlanDetailPage.tsx` — Stats cards, inline result recording, state actions
- `src/features/plans/components/NewPlanPage.tsx` — Case selector panel
- `src/features/plans/components/ResultRecordModal.tsx` — Polish to spec
- `src/features/testcases/components/CaseListPage.tsx` — Split panel with module tree
- `src/features/testcases/components/CaseDetailPage.tsx` — AI metadata collapse, execution history
- `src/features/testcases/components/CreateCaseDrawer.tsx` — ArrayEditor integration
- `src/features/modules/components/ModuleManagePage.tsx` — Split panel with tree
- `src/features/configs/components/ConfigManagePage.tsx` — Import/export JSON
- `src/features/documents/components/KnowledgeListPage.tsx` — Status filters, icons
- `src/features/documents/components/DocumentDetailPage.tsx` — Split layout, chunks
- `src/features/documents/components/UploadDocumentModal.tsx` — File type restrictions

### Files to Modify (Minor)
- `src/features/auth/components/LoginPage.tsx` — Brand area styling
- `src/features/auth/components/RegisterPage.tsx` — Role select (no super_admin)
- `src/store/useAppStore.ts` — Remove currentProject from store (use React Query)

---

## Phase 1: Routing & Layout Foundation

### Task 1: Restructure Router to Nested Project-Scoped Routes

**Files:**
- Modify: `src/router/index.tsx`

The current router uses flat paths (`/testcases`, `/plans`, `/generation`). Spec requires project-scoped nesting (`/projects/:id/cases`, `/projects/:id/plans`, etc.). This is the single highest-priority change because every page component expects route params.

- [ ] **Step 1: Rewrite router/index.tsx with nested project routes**

```tsx
// src/router/index.tsx
import { createBrowserRouter, Navigate } from 'react-router-dom'
import { RouteGuard } from '@/router/RouteGuard'
import { AuthErrorBoundary } from '@/components/ErrorBoundary'
import { App } from '@/app/App'
import { AppLayout } from '@/components/layout/AppLayout'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <Navigate to="/projects" replace /> },
      // Public routes
      {
        path: 'login',
        lazy: () =>
          import('../features/auth/components/LoginPage').then((m) => ({
            Component: () => (
              <AuthErrorBoundary>
                <m.LoginPage />
              </AuthErrorBoundary>
            ),
          })),
      },
      {
        path: 'register',
        lazy: () =>
          import('../features/auth/components/RegisterPage').then((m) => ({
            Component: () => (
              <AuthErrorBoundary>
                <m.RegisterPage />
              </AuthErrorBoundary>
            ),
          })),
      },
      // Protected routes
      {
        path: '/',
        element: (
          <AuthErrorBoundary>
            <RouteGuard>
              <AppLayout />
            </RouteGuard>
          </AuthErrorBoundary>
        ),
        children: [
          {
            path: 'projects',
            lazy: () =>
              import('../features/projects/components/ProjectListPage').then(
                (m) => ({ Component: m.ProjectListPage })
              ),
          },
          // Project-scoped routes — all nested under /projects/:id
          {
            path: 'projects/:id',
            children: [
              {
                index: true,
                lazy: () =>
                  import('../features/projects/components/ProjectDashboard').then(
                    (m) => ({ Component: m.ProjectDashboard })
                  ),
              },
              // Knowledge Base
              {
                path: 'knowledge',
                lazy: () =>
                  import('../features/documents/components/KnowledgeListPage').then(
                    (m) => ({ Component: m.KnowledgeListPage })
                  ),
              },
              {
                path: 'knowledge/figma',
                lazy: () =>
                  import('../features/documents/components/FigmaIntegrationPage').then(
                    (m) => ({ Component: m.FigmaIntegrationPage })
                  ),
              },
              {
                path: 'knowledge/:docId',
                lazy: () =>
                  import('../features/documents/components/DocumentDetailPage').then(
                    (m) => ({ Component: m.DocumentDetailPage })
                  ),
              },
              // AI Generation
              {
                path: 'generation',
                lazy: () =>
                  import('../features/generation/components/GenerationTaskListPage').then(
                    (m) => ({ Component: m.GenerationTaskListPage })
                  ),
              },
              {
                path: 'generation/new',
                lazy: () =>
                  import('../features/generation/components/NewGenerationTaskPage').then(
                    (m) => ({ Component: m.NewGenerationTaskPage })
                  ),
              },
              {
                path: 'generation/:taskId',
                lazy: () =>
                  import('../features/generation/components/TaskDetailPage').then(
                    (m) => ({ Component: m.TaskDetailPage })
                  ),
              },
              // Test Cases
              {
                path: 'cases',
                lazy: () =>
                  import('../features/testcases/components/CaseListPage').then(
                    (m) => ({ Component: m.CaseListPage })
                  ),
              },
              {
                path: 'cases/:caseId',
                lazy: () =>
                  import('../features/testcases/components/CaseDetailPage').then(
                    (m) => ({ Component: m.CaseDetailPage })
                  ),
              },
              // Test Plans
              {
                path: 'plans',
                lazy: () =>
                  import('../features/plans/components/PlanListPage').then((m) => ({
                    Component: m.PlanListPage,
                  })),
              },
              {
                path: 'plans/new',
                lazy: () =>
                  import('../features/plans/components/NewPlanPage').then((m) => ({
                    Component: m.NewPlanPage,
                  })),
              },
              {
                path: 'plans/:planId',
                lazy: () =>
                  import('../features/plans/components/PlanDetailPage').then(
                    (m) => ({ Component: m.PlanDetailPage })
                  ),
              },
              // Settings (admin+)
              {
                path: 'settings',
                element: <RouteGuard requireAdmin>{null}</RouteGuard>,
                children: [
                  {
                    index: true,
                    lazy: () =>
                      import('../features/configs/components/ConfigManagePage').then(
                        (m) => ({ Component: m.ConfigManagePage })
                      ),
                  },
                ],
              },
              {
                path: 'settings/modules',
                element: <RouteGuard requireAdmin>{null}</RouteGuard>,
                children: [
                  {
                    index: true,
                    lazy: () =>
                      import('../features/modules/components/ModuleManagePage').then(
                        (m) => ({ Component: m.ModuleManagePage })
                      ),
                  },
                ],
              },
              {
                path: 'settings/configs',
                element: <RouteGuard requireAdmin>{null}</RouteGuard>,
                children: [
                  {
                    index: true,
                    lazy: () =>
                      import('../features/configs/components/ConfigManagePage').then(
                        (m) => ({ Component: m.ConfigManagePage })
                      ),
                  },
                ],
              },
            ],
          },
          // Global Drafts
          {
            path: 'drafts',
            lazy: () =>
              import('../features/drafts/components/DraftListPage').then(
                (m) => ({ Component: m.DraftListPage })
              ),
          },
          {
            path: 'drafts/:draftId',
            lazy: () =>
              import('../features/drafts/components/DraftConfirmPage').then(
                (m) => ({ Component: m.DraftConfirmPage })
              ),
          },
          { path: '*', element: <Navigate to="/projects" replace /> },
        ],
      },
      // 404 fallback
      {
        path: '*',
        lazy: () =>
          import('../components/NotFoundPage').then((m) => ({
            Component: m.NotFoundPage,
          })),
      },
    ],
  },
])
```

- [ ] **Step 2: Update all page components that read route params**

Every page currently reads `projectId` or `caseId` from route params. The param names changed (e.g. `projectId` → `id`, `caseId` → `caseId` stays). Audit and fix:

```bash
grep -rn "useParams" src/features/ --include="*.tsx"
```

Key changes needed in each page component:
- `ProjectDashboard`: `useParams<{ id: string }>()` — param is now `id` not `projectId`
- `CaseListPage`: `useParams<{ id: string }>()` — use `id` as `projectId`
- `CaseDetailPage`: `useParams<{ id: string; caseId: string }>()`
- `PlanListPage`: `useParams<{ id: string }>()`
- `NewPlanPage`: `useParams<{ id: string }>()`
- `PlanDetailPage`: `useParams<{ id: string; planId: string }>()`
- `GenerationTaskListPage`: `useParams<{ id: string }>()`
- `NewGenerationTaskPage`: `useParams<{ id: string }>()`
- `TaskDetailPage`: `useParams<{ id: string; taskId: string }>()`
- `KnowledgeListPage`: `useParams<{ id: string }>()`
- `DocumentDetailPage`: `useParams<{ id: string; docId: string }>()`
- `ModuleManagePage`: `useParams<{ id: string }>()`
- `ConfigManagePage`: `useParams<{ id: string }>()`

For each file, change `projectId` to `id` in the `useParams` type and destructuring.

- [ ] **Step 3: Update Sidebar links to use new route paths**

In `src/components/layout/Sidebar.tsx`, update all `<Link>` paths from flat (`/testcases`) to nested (`/projects/${currentProjectId}/cases`), etc.

- [ ] **Step 4: Update navigation calls across all components**

Search for all `navigate('/testcases')`, `navigate('/plans')`, etc. and update to project-scoped paths. Key patterns to fix:

```bash
grep -rn "navigate.*testcases\|navigate.*plans\|navigate.*generation\|navigate.*documents\|navigate.*modules\|navigate.*configs" src/features/ --include="*.tsx"
```

Each `navigate('/foo')` becomes `navigate(`/projects/${projectId}/foo`)`.

- [ ] **Step 5: Run tests and fix param name changes**

Run: `npx vitest run --reporter=verbose 2>&1 | tail -60`
Fix any failing tests due to route param name changes.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor(router): restructure to nested project-scoped routes per spec v2.0"
```

---

### Task 2: Update Sidebar with Project Submenu & Draft Badge

**Files:**
- Modify: `src/components/layout/Sidebar.tsx`
- Modify: `src/store/useAppStore.ts`

- [ ] **Step 1: Write failing test for Sidebar project submenu**

Add test case in `src/components/layout/Sidebar.test.tsx` verifying:
- When `currentProjectId` is set, project submenu items render
- Draft Badge shows pending count
- AI generation menu item has purple icon
- Menu item click navigates to correct project-scoped path

- [ ] **Step 2: Run test to verify it fails**

Run: `npx vitest run src/components/layout/Sidebar.test.tsx`

- [ ] **Step 3: Implement Sidebar project submenu**

Update `Sidebar.tsx`:
- Add project submenu section with items: 仪表盘, 知识库, AI 生成 (purple icon), 测试用例, 测试计划, 项目设置
- Use `useDraftCount()` hook (from `useDrafts`) for Badge on drafts menu item
- Selected state uses purple background `rgba(123,97,255,0.10)` + purple text + left 3px border
- Collapsed state shows only icons with tooltip

- [ ] **Step 4: Run test to verify it passes**

Run: `npx vitest run src/components/layout/Sidebar.test.tsx`

- [ ] **Step 5: Commit**

```bash
git add src/components/layout/Sidebar.tsx src/components/layout/Sidebar.test.tsx
git commit -m "feat(sidebar): add project submenu, draft badge, and AI purple styling"
```

---

### Task 3: Add Breadcrumb to Header

**Files:**
- Modify: `src/components/layout/Header.tsx`

- [ ] **Step 1: Write failing test for breadcrumb**

Test that breadcrumb renders: `首页 > 项目列表 > {ProjectName} > {PageName}` based on current route.

- [ ] **Step 2: Run test to verify it fails**

- [ ] **Step 3: Implement breadcrumb**

Use `useLocation` + `useParams` to build breadcrumb path. Map route segments to Chinese labels:
- `/projects` → "首页"
- `/projects/:id` → "项目列表 > {projectName}"
- `/projects/:id/cases` → "项目列表 > {projectName} > 测试用例"
- `/projects/:id/cases/:caseId` → "... > 测试用例 > {caseNumber}"
- etc.

Use `Arco Breadcrumb` component. Skip link for last item (current page).

- [ ] **Step 4: Run test to verify it passes**

- [ ] **Step 5: Commit**

```bash
git commit -m "feat(header): add breadcrumb navigation with project context"
```

---

### Task 4: Add AI Animation Styles

**Files:**
- Create: `src/styles/animations.css`
- Modify: `src/app/providers.tsx` (import animations.css)

- [ ] **Step 1: Create animations.css with keyframes**

```css
/* src/styles/animations.css */

/* AI content reveal animation */
@keyframes ai-reveal {
  from {
    opacity: 0;
    transform: translateY(8px) scale(0.97);
    filter: blur(2px);
  }
  to {
    opacity:1;
    transform: translateY(0) scale(1);
    filter: blur(0);
  }
}

/* Glow pulse for processing state */
@keyframes glow-pulse {
  0%, 100% { box-shadow: 0 0 15px rgba(123,97,255,0.15); }
  50% { box-shadow: 0 0 25px rgba(123,97,255,0.30); }
}

/* Progress shimmer */
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

/* Breathing glow for high confidence */
@keyframes breath-glow {
  0%, 100% { box-shadow: 0 0 8px rgba(0,180,42,0.2); }
  50% { box-shadow: 0 0 12px rgba(0,180,42,0.4); }
}

/* Spring celebration */
@keyframes spring-celebrate {
  0% { transform: scale(1) rotate(0deg); }
  30% { transform: scale(1.2) rotate(-5deg); }
  60% { transform: scale(0.95) rotate(3deg); }
  100% { transform: scale(1) rotate(0deg); }
}

.animate-ai-reveal {
  animation: ai-reveal 400ms ease-out both;
}

.animate-glow-pulse {
  animation: glow-pulse 2s ease-in-out infinite;
}

.animate-shimmer {
  background: linear-gradient(90deg, transparent 0%, rgba(123,97,255,0.08) 50%, transparent 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s linear infinite;
}

.animate-breath-glow {
  animation: breath-glow 2s ease-in-out infinite;
}

.animate-spring {
  animation: spring-celebrate 500ms cubic-bezier(0.34, 1.56, 0.64, 1);
}

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  .animate-ai-reveal,
  .animate-glow-pulse,
  .animate-shimmer,
  .animate-breath-glow,
  .animate-spring {
    animation: none !important;
  }
}
```

- [ ] **Step 2: Import in providers.tsx**

Add `import '@/styles/animations.css'` to `src/app/providers.tsx`.

- [ ] **Step 3: Commit**

```bash
git add src/styles/animations.css src/app/providers.tsx
git commit -m "feat(styles): add AI-specific animation keyframes"
```

---

## Phase 2: Shared Component Polish

### Task 5: Complete StatusTag Color Maps

**Files:**
- Modify: `src/components/business/StatusTag.tsx`

- [ ] **Step 1: Write failing tests for missing color maps**

In `src/components/business/StatusTag.test.tsx`, add test cases for all `type` values:
- `planStatus`: draft(灰), active(紫+glow), completed(绿), archived(灰淡)
- `taskStatus`: pending(灰), processing(紫+glow), completed(绿), failed(红)
- `draftStatus`: pending(橙), confirmed(绿), rejected(红)
- `confidence`: high(绿), medium(黄), low(红)
- `caseType`: functionality(紫), performance(蓝), api(青), ui(橙), security(紫深)
- `documentType`: prd(紫), figma(紫), api_spec(青), swagger(绿), markdown(灰)
- `documentStatus`: pending(灰), processing(紫), completed(绿), failed(红)

Verify each renders correct `color`, `backgroundColor`, and `children` text.

- [ ] **Step 2: Run tests to verify they fail**

- [ ] **Step 3: Update StatusTag with complete color maps**

Replace the existing COLOR_MAP in `StatusTag.tsx` with the full set from UX spec §2.1. Use the exact hex values and rgba backgrounds specified. Add text labels for all Chinese translations.

- [ ] **Step 4: Run tests to verify they pass**

- [ ] **Step 5: Commit**

```bash
git add src/components/business/StatusTag.tsx src/components/business/StatusTag.test.tsx
git commit -m "feat(statustag): add all color maps per UX spec §2.1"
```

---

### Task 6: Enhance SearchTable Component

**Files:**
- Modify: `src/components/business/SearchTable.tsx`

- [ ] **Step 1: Update SearchTable with spec styling defaults**

Set default props matching UX spec §5.1:
- row height 48px (default), 36px (compact via `size="small"`)
- header background `#F7F8FA`, font 13px Medium
- zebra stripe on even rows `#FAFBFC`
- border `1px solid #E5E6EB`
- hover row `#F2F3F5`
- default pagination: 20/page, options [10, 20, 50, 100], showTotal
- empty state: Arco Empty with icon + message

- [ ] **Step 2: Add batch operation bar support**

When rows are selected, render a sticky bar above the table: "已选择 N 项 + 清空" on left, action buttons on right. Animate with slideDown 200ms.

- [ ] **Step 3: Write tests for batch bar**

Test that selecting rows shows the bar, clicking "清空" deselects all, bar disappears when no rows selected.

- [ ] **Step 4: Commit**

```bash
git add src/components/business/SearchTable.tsx src/components/business/SearchTable.test.tsx
git commit -m "feat(searchtable): add spec styling, pagination, batch operation bar"
```

---

### Task 7: Enhance ArrayEditor with Validation

**Files:**
- Modify: `src/components/business/ArrayEditor.tsx`

- [ ] **Step 1: Add `minItems` prop support**

ArrayEditor should accept `minItems` prop. When `minItems > 0` and items are fewer, show validation error. Default delete button disabled when at `minItems`.

- [ ] **Step 2: Add reorder buttons (up/down)**

Each row gets ChevronUp/ChevronDown buttons. Clicking swaps with adjacent item. Top item's up button and bottom item's down button are disabled.

- [ ] **Step 3: Write tests**

Test minItems validation, reorder behavior, add/remove item.

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(arrayeditor): add minItems validation and reorder buttons"
```

---

## Phase 3: Auth & Projects

### Task 8: Polish Login Page with Brand Area

**Files:**
- Modify: `src/features/auth/components/LoginPage.tsx`
- Modify: `src/features/auth/components/LoginBanner.tsx`

- [ ] **Step 1: Update LoginPage to left-right split layout (55%:45%)**

Left side: LoginBanner with mesh gradient background (#7B61FF → #5A3DC0 → #3B1FA0), floating geometric shapes (CSS animation), product slogan.

Right side: Clean white background, Aitestos logo with ✦ icon, email/password form with large rounded inputs (border-radius: 8px, padding: 12px 16px). Add "记住我" checkbox.

- [ ] **Step 2: Add proper error display**

Replace any inline error with `Arco Alert` (type=error) at top of form area on login failure. Show "邮箱或密码错误".

- [ ] **Step 3: Update existing tests**

Fix any broken tests from layout change. Verify login flow still works.

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(login): implement split layout with brand area per UX spec"
```

---

### Task 9: Polish Register Page

**Files:**
- Modify: `src/features/auth/components/RegisterPage.tsx`

- [ ] **Step 1: Add role select (admin/normal only, no super_admin)**

Ensure the role Arco Select only shows "管理员 (admin)" and "测试工程师 (normal)". No super_admin option.

- [ ] **Step 2: Add confirm password field with Zod validation**

Add register Zod schema at `src/features/auth/schema/registerSchema.ts`:
```typescript
import { z } from 'zod/v4'
export const registerSchema = z
  .object({
    username: z.string().min(3).max(32),
    email: z.email(),
    password: z.string().min(8),
    confirmPassword: z.string(),
    role: z.enum(['admin', 'normal']),
  })
  .refine((d) => d.password === d.confirmPassword, {
    message: '两次密码输入不一致',
    path: ['confirmPassword'],
  })
```

- [ ] **Step 2: Handle 409 conflict errors at field level**

On 409 response, check error message for "邮箱" or "用户名" and set the corresponding field error.

- [ ] **Step 3: Commit**

```bash
git commit -m "feat(register): add confirm password, role select, field-level 409 handling"
```

---

### Task 10: Convert Project List to Card Grid

**Files:**
- Modify: `src/features/projects/components/ProjectListPage.tsx`

- [ ] **Step 1: Replace table with Card grid (3 columns)**

Use Arco `Grid` with `Col span={8}` for 3-column layout. Each card:
- Top 3px color bar (hash project prefix to pick from palette)
- Project name as h3
- Prefix as Tag with monospace font
- Description (2-line truncation with line-clamp)
- Stats row: document count + case count
- Hover: shadow-md elevation + action buttons fadeIn
- Action buttons: "进入项目" (navigates to `/projects/${id}`) + "设置" (navigates to `/projects/${id}/settings`)

- [ ] **Step 2: Add empty state**

When no projects: `FolderOpen` icon + "创建第一个项目开始测试管理之旅" + "新建项目" button.

- [ ] **Step 3: Update tests**

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(projects): convert list to card grid with hover effects per UX spec"
```

---

### Task 11: Enhance Project Dashboard

**Files:**
- Modify: `src/features/projects/components/ProjectDashboard.tsx`

- [ ] **Step 1: Add quick-start guide for empty projects**

When all stats are 0, show 3-step guide: ❶ Upload documents → ❷ AI generate → ❸ Create plan. Each step is a clickable card linking to the respective page.

- [ ] **Step 2: Add stats cards with decoration lines**

4 StatsCards in a row: 用例总数(case_count), 通过率(pass_rate%), 需求覆盖率(coverage_rate%), AI生成数(ai_generated_count). Each with colored left decoration line. AI card gets purple line + AI Glow shadow.

- [ ] **Step 3: Add recent generation tasks section**

Show `recent_tasks` from stats response. Each item: status Tag + result summary (X/Y confirmed) + created time. Link to task detail.

- [ ] **Step 4: Add pass rate trend chart**

Render `pass_rate_trend` data as a line chart using ECharts or Recharts. Brand purple line (#7B61FF), grid lines #E5E6EB dashed, axis labels 12px #86909C. Show empty state "暂无数据" when no trend data.

- [ ] **Step 5: Commit**

```bash
git commit -m "feat(dashboard): add stats cards, trend chart, recent tasks, quick-start guide"
```

---

### Task 12: Polish CreateProjectModal

**Files:**
- Modify: `src/features/projects/components/CreateProjectModal.tsx`

- [ ] **Step 1: Add prefix format validation**

Real-time validation on prefix input: must be 2-4 uppercase letters matching `^[A-Z]+$`. Show inline error if format wrong. Convert lowercase input to uppercase automatically.

- [ ] **Step 2: Add Zod schema**

Create `src/features/projects/schema/projectSchema.ts`:
```typescript
import { z } from 'zod/v4'
export const createProjectSchema = z.object({
  name: z.string().min(2).max(255),
  prefix: z.string().min(2).max(4).regex(/^[A-Z]+$/, '前缀必须为2-4位大写字母'),
  description: z.string().optional(),
})
```

- [ ] **Step 3: Handle 409 conflict at field level**

Parse 409 error to determine if name or prefix conflicts, set field error accordingly.

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(create-project): add prefix validation, Zod schema, 409 field errors"
```

---

## Phase 4: Knowledge & Generation (Core Flow)

### Task 13: Polish Knowledge List Page

**Files:**
- Modify: `src/features/documents/components/KnowledgeListPage.tsx`

- [ ] **Step 1: Add filter bar with type, status, search**

Add Arco Select for document type and status filtering, plus Input.Search for name search. Pass filters to `useDocuments` hook query params.

- [ ] **Step 2: Add document type icons**

Prepend each document name with the appropriate Lucide icon from the document type table (FileText for PRD, Figma icon, FileCode2 for api_spec, etc.) with the matching type color.

- [ ] **Step 3: Add empty state**

FileText icon + "暂无文档" + "上传第一份需求文档，开启智能测试之旅" + "上传文档" button.

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(documents): add filters, type icons, empty state to knowledge list"
```

---

### Task 14: Polish Document Detail Page

**Files:**
- Modify: `src/features/documents/components/DocumentDetailPage.tsx`

- [ ] **Step 1: Implement split layout**

Use SplitPanel: left 300px info panel, right chunks list. Left panel shows: doc name(h2) + type Tag + status Tag + uploader info + chunk count + Arco Steps status flow. Failed status shows Alert + retry button.

- [ ] **Step 2: Implement chunks list**

Right side: title "文档分块" + count Badge. List of items with chunk_index + content preview (3-line truncation) + "展开" toggle. Show reference count icon for cited chunks.

- [ ] **Step 3: Commit**

```bash
git commit -m "feat(document-detail): implement split layout with info panel and chunks"
```

---

### Task 15: Create Figma Integration Page (Stub)

**Files:**
- Create: `src/features/documents/components/FigmaIntegrationPage.tsx`

- [ ] **Step 1: Create FigmaIntegrationPage component**

Page with 3 sections:
1. Connection config: Radio.Group (personal token / OAuth) + Input.Password + "测试连接" button
2. Import file: URL Input + "解析" button
3. Node selection: Arco Tree with Checkbox mode, default all selected

Since backend Figma API doesn't exist yet, stub the tree with placeholder data and disable the import action. Show a "即将推出" notice.

- [ ] **Step 2: Commit**

```bash
git commit -m "feat(figma): create Figma integration page stub"
```

---

### Task 16: Enhance Generation Task List Page

**Files:**
- Modify: `src/features/generation/components/GenerationTaskListPage.tsx`

- [ ] **Step 1: Add AI visual styling**

- "新建生成任务" button: AI Gradient background + Sparkles icon
- Processing rows: `animate-glow-pulse` class + row background `rgba(123,97,255,0.04)`
- Task ID column: monospace font, purple color (#7B61FF), hover shows full ID

- [ ] **Step 2: Add status filter and empty state**

Arco Select for status filtering. Empty state: Sparkles icon + "暂无生成任务" + "新建生成任务" button (AI Gradient).

- [ ] **Step 3: Add retry action for failed tasks**

Failed rows show "重试" button → Popconfirm → resubmit task with same params.

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(generation-list): add AI styling, status filter, retry action"
```

---

### Task 17: Enhance New Generation Task Page

**Files:**
- Modify: `src/features/generation/components/NewGenerationTaskPage.tsx`

- [ ] **Step 1: Add knowledge readiness indicator**

At page top, fetch document count for current project. Show:
- Green "就绪" if documents exist and are completed
- Yellow "内容有限" if only a few
- Red "请先上传需求文档" + disable submit button if zero documents

- [ ] **Step 2: Add advanced options collapse**

Arco Collapse wrapping: scene type (Checkbox.Group), priority preference (Select), case type (Select), generation mode (Radio.Group). Default collapsed.

- [ ] **Step 3: Add Zod schema and form validation**

Create `src/features/generation/schema/createTaskSchema.ts`:
```typescript
import { z } from 'zod/v4'
export const createTaskSchema = z.object({
  module_id: z.string().min(1, '请选择目标模块'),
  prompt: z.string().min(10, '需求描述至少10个字符'),
  case_count: z.number().min(1).max(20).optional(),
  scene_types: z.array(z.enum(['positive', 'negative', 'boundary'])).optional(),
  priority: z.enum(['P0', 'P1', 'P2', 'P3']).optional(),
  case_type: z.enum(['functionality', 'performance', 'api', 'ui', 'security']).optional(),
})
```

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(new-task): add readiness indicator, advanced options, Zod schema"
```

---

### Task 18: Enhance Task Detail Page with AI Status Card

**Files:**
- Modify: `src/features/generation/components/TaskDetailPage.tsx`

- [ ] **Step 1: Add AI status card**

Top card showing task status with AI visual effects:
- pending: Spin + "任务排队中" + faint AI Glow
- processing: Progress bar + shimmer + estimated time + glow-pulse
- completed: Success Badge + spring celebration animation on transition
- failed: Error Alert + retry button

- [ ] **Step 2: Add draft list with batch operations**

When completed, show drafts table with: Checkbox, title, type Tag, priority Tag, confidence Tag, actions. Batch toolbar: "批量确认" (AI Gradient) + "批量拒绝". Click row → navigate to draft confirm page.

- [ ] **Step 3: Add ai-reveal animation on draft list items**

Apply `animate-ai-reveal` with staggered delay (index * 80ms) to each draft row.

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(task-detail): add AI status card, draft list with batch operations"
```

---

### Task 19: Enhance Draft Confirm Page

**Files:**
- Modify: `src/features/drafts/components/DraftConfirmPage.tsx`

- [ ] **Step 1: Add draft navigation bar**

Top bar: "← 返回草稿列表" link | "第 N / M 条" progress | dot navigation. Keyboard ←/→ to switch drafts. Auto-save edits before switching.

- [ ] **Step 2: Implement split layout with reference panel**

Use SplitPanel (60/40). Left: edit form with ArrayEditor for preconditions/steps, Input for title, TextArea for expected, Selects for type/priority. Right: ReferencePanel with AI glow background.

- [ ] **Step 3: Implement bottom action bar**

Fixed bottom bar: "拒绝" (danger) → reject Modal with reason Radio + feedback TextArea. "保存修改" (default) → save edits. "确认并转为正式用例" (primary, AI Gradient) → confirm Modal with module Select.

- [ ] **Step 4: Add unsaved changes protection**

Use `form.isDirty` check before navigation/closing. Show Popconfirm: "有未保存的变更，确认离开？"

- [ ] **Step 5: Commit**

```bash
git commit -m "feat(draft-confirm): add navigation, split layout, reference panel, action bar"
```

---

## Phase 5: Test Cases & Plans

### Task 20: Enhance Case List Page with Split Panel

**Files:**
- Modify: `src/features/testcases/components/CaseListPage.tsx`

- [ ] **Step 1: Add module tree on left panel**

Use SplitPanel. Left: Arco Tree with "全部" option + module nodes (name + case count Badge). Click node → filter table by module_id. Fetch modules via `useModules(projectId)`.

- [ ] **Step 2: Add toolbar with filters and actions**

Right side toolbar: search Input + status/type/priority Select filters + "新建用例"(primary) + "导入"(default) + "导出"(default). Use URL query params for filter persistence.

- [ ] **Step 3: Add batch operation bar**

Selected rows → sticky bar: "修改优先级" / "修改状态" / "加入计划" / "删除".

- [ ] **Step 4: Add empty state**

FileCheck icon + two buttons: "手动创建" (primary) + "使用 AI 生成" (AI Gradient, links to `/projects/${id}/generation/new`).

- [ ] **Step 5: Commit**

```bash
git commit -m "feat(case-list): add module tree, filters, batch ops, empty state"
```

---

### Task 21: Enhance Case Detail Page with AI Metadata

**Files:**
- Modify: `src/features/testcases/components/CaseDetailPage.tsx`

- [ ] **Step 1: Add page header with back link, number, status**

Back link (ArrowLeft + "返回用例库") + number (monospace, h1) + title + large StatusTag + action buttons (编辑/复制/删除).

- [ ] **Step 2: Add basic info card**

3-column grid: 用例类型, 优先级, 所属模块, 创建人, 创建时间, 更新时间.

- [ ] **Step 3: Add AI metadata collapse panel**

Arco Collapse with purple background header (#Sparkles icon). Content: generation task link, confidence StatusTag, referenced chunks (similarity-colored), model version, generation time. Special alerts for changed/deleted source documents.

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(case-detail): add header, info card, AI metadata panel"
```

---

### Task 22: Enhance Plan Detail Page with Stats & Inline Recording

**Files:**
- Modify: `src/features/plans/components/PlanDetailPage.tsx`

- [ ] **Step 1: Add stats cards (5 columns)**

StatsCard row: 总用例(灰), 通过(绿), 失败(红), 阻塞(橙), 跳过(灰) + unexecuted count. Read from `plan.stats`.

- [ ] **Step 2: Add execution progress bar**

Arco Progress bar, brand purple, showing executed/total percentage.

- [ ] **Step 3: Add state-dependent action buttons**

Render different button sets based on plan status (draft → 编辑+开始执行+删除, active → 编辑+标记完成, completed → 重新执行+归档, archived → 取消归档).

- [ ] **Step 4: Add inline result recording (quick entry)**

In the execution table, make the result column an inline Arco Select. On change → auto-submit result → flash row with result color → Toast "已录入：通过" + "撤销" link (3s timeout).

- [ ] **Step 5: Commit**

```bash
git commit -m "feat(plan-detail): add stats cards, progress bar, state actions, inline recording"
```

---

### Task 23: Enhance New Plan Page with Case Selector

**Files:**
- Modify: `src/features/plans/components/NewPlanPage.tsx`

- [ ] **Step 1: Implement split layout**

SplitPanel: left form (name + description), right case selector.

- [ ] **Step 2: Implement case selector with tabs**

Two tabs: "可选用例" (table with filters + checkbox) and "已选用例" (selected list + remove buttons, Badge count). Case selection stored in local state, submitted as `case_ids` on form submit.

- [ ] **Step 3: Commit**

```bash
git commit -m "feat(new-plan): add case selector with tabs and filters"
```

---

## Phase 6: Settings & Configs

### Task 24: Enhance Module Manage Page

**Files:**
- Modify: `src/features/modules/components/ModuleManagePage.tsx`

- [ ] **Step 1: Add split layout with module tree**

SplitPanel: left module tree (Arco Tree) with nodes showing name + abbreviation Tag. Hover shows edit/delete icons. Bottom "新增模块" button (dashed style).

- [ ] **Step 2: Add type-to-confirm delete Modal**

For module deletion, use Modal (not Popconfirm) with text input: "请输入模块名称 '{name}' 确认删除". Delete button only enabled when input matches exactly.

- [ ] **Step 3: Commit**

```bash
git commit -m "feat(modules): add split layout, tree, type-to-confirm delete"
```

---

### Task 25: Enhance Config Manage Page with Import/Export

**Files:**
- Modify: `src/features/configs/components/ConfigManagePage.tsx`

- [ ] **Step 1: Add import JSON functionality**

"导入 JSON" button → Modal with TextArea for JSON paste + preview table showing parsed configs + "确认导入" button. Call `configsApi.import()` on submit.

- [ ] **Step 2: Add export JSON functionality**

"导出 JSON" button → call `configsApi.export()` → download as .json file using Blob + URL.createObjectURL.

- [ ] **Step 3: Commit**

```bash
git commit -m "feat(configs): add JSON import/export, preview table"
```

---

## Verification Checklist

After all tasks complete, verify against spec §7 acceptance tests:

- [ ] **TC-A01~A04**: Login/register/token refresh work
- [ ] **TC-P01~P03**: Project CRUD, prefix validation, dashboard stats
- [ ] **TC-K01~K02**: Document upload, detail with chunks
- [ ] **TC-G01~G03**: Generation tasks, knowledge check, draft viewing
- [ ] **TC-D01~D03**: Draft confirm (single/batch/reject), numbering
- [ ] **TC-C01~C03**: Case create, detail with AI metadata, numbering
- [ ] **TC-PL01~PL03**: Plan create, result recording, state transitions

Run: `npx vitest run --reporter=verbose` to ensure all unit tests pass.

Run: `make check` to verify lint + format + type-check.

Run: `make build` to verify production build succeeds.
