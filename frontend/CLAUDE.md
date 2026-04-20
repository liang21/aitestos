# CLAUDE.md - Aitestos 前端工程操作手册

## 1. 指令优先级

1. **执行宪法:** 任何行动前必须阅读并对齐 `@constitution.md`。
2. **环境感知:** 修改前执行 `ls src/` 了解当前代码状态。
3. **TDD First:** 先生成 MSW handler 与测试用例，再生成实现。

---

## 2. 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| React | 19.2 | UI 框架 |
| TypeScript | 5.9 | 类型系统 (strict: true) |
| React Router | 7.x | 路由 |
| TanStack React Query | 5.x | 服务端状态管理 |
| Zustand | 5.x | 客户端全局状态 |
| React Hook Form | 7.x | 表单管理 |
| Zod | 4.x | 运行时校验 |
| Arco Design | 2.66 | UI 组件库 |
| Tailwind CSS | 4.x | 样式 |
| Axios | 1.14 | HTTP 客户端（仅 React Query queryFn 底层） |
| Vitest | 4.x | 测试运行器 |
| MSW | 2.x | API Mock |

---

## 3. 代码风格

* **Hook 优先:** 禁止 class component。
* **解耦:** UI 组件仅负责展示，业务逻辑封装在 `useFeature.ts` 中。
* **导包顺序:** React → 第三方库 → @/ 内部模块（分块排列，空行分隔）。

### 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 页面组件 | PascalCase + Page 后缀 | `CaseListPage.tsx` |
| 业务组件 | PascalCase + 类型后缀 | `CreateProjectModal.tsx`、`ResultRecordModal.tsx` |
| 共享组件 | PascalCase | `StatusTag.tsx`、`SearchTable.tsx` |
| 自定义 Hook | camelCase + use 前缀 | `useProjects.ts`、`usePollingTask.ts` |
| API service | camelCase | `projects.ts`、`generation.ts` |
| 类型文件 | enums.ts / api.ts | 集中在 `src/types/` |
| 测试文件 | 同目录 `.test.tsx` | `useProjects.test.ts` 放在 `hooks/` 旁 |

### 测试文件约定

* 位置：与源文件同目录（如 `src/features/auth/hooks/useAuth.test.ts`）
* 匹配：`src/**/*.test.{ts,tsx}`（由 vitest.config.ts 配置）
* MSW handlers：集中管理在 `tests/msw/handlers/`

---

## 4. Agent 思考链 (CoT) 要求

执行任何代码编写前，必须按顺序输出：

1. **Spec 对齐:** 确认功能点与 spec.md 对应。
2. **ADR:** 状态放在哪里？（React Query / Zustand / useState）关键依赖是什么？
3. **测试方案:** 列出关键路径的测试点和 MSW mock 策略。

---

## 5. 数据流规范

* **查询:** `useQuery` — 列表、详情、统计数据。
* **变更:** `useMutation` — 创建、更新、删除。成功后 `invalidateQueries` 刷新缓存。
* **轮询:** `useQuery` + `refetchInterval` — AI 任务状态、草稿 Badge。
* **表单:** React Hook Form + zod resolver。
* **API 边界:** `services/` 层封装 Axios 调用，组件禁止直接 import axios。

---

## 6. 目录架构

```text
src/
├── app/                    # 应用入口
│   ├── App.tsx             # 根组件（Provider 嵌套）
│   ├── main.tsx            # 渲染入口
│   └── providers.tsx       # QueryClientProvider + ConfigProvider
├── router/                 # 路由
│   ├── index.tsx           # 路由定义（lazy loading）
│   └── RouteGuard.tsx      # 认证/权限守卫
├── lib/                    # 基础设施
│   ├── request.ts          # Axios 实例（Token 刷新、错误拦截）
│   ├── query-client.ts     # React Query 全局配置
│   └── utils.ts            # cn() 等工具函数
├── types/                  # 全局类型
│   ├── enums.ts            # 枚举/字面量联合类型
│   └── api.ts              # API 请求/响应类型
├── features/               # 业务功能模块 (Feature-Based)
│   ├── auth/               # 认证
│   │   ├── components/     # LoginPage, RegisterPage
│   │   ├── hooks/          # useAuth, useAuthStore
│   │   └── services/       # auth.ts (API function)
│   ├── projects/           # 项目管理
│   ├── modules/            # 模块管理
│   ├── testcases/          # 测试用例
│   ├── plans/              # 测试计划
│   ├── generation/         # AI 生成
│   ├── drafts/             # 草稿箱
│   ├── documents/          # 知识库
│   └── configs/            # 项目配置
├── components/             # 跨 Feature 共享组件
│   ├── layout/             # AppLayout, Sidebar, Header, AuthLayout
│   └── business/           # StatusTag, SearchTable, ArrayEditor, StatsCard 等
├── store/                  # 全局 Zustand store
│   └── useAppStore.ts      # sidebarCollapsed + notifications
└── styles/                 # 全局样式
    └── theme.css           # Arco 主题变量 + Tailwind @theme
```

### 依赖规则

* Page → Feature hooks → Feature services → `@/lib/request`
* `components/` 禁止引用 `features/` 的任何内容
* `services/` 禁止引用 `store/` 或 `hooks/`
* Feature A 禁止引用 Feature B 的内部文件
* `@/types` 可被所有层引用，但自身禁止引用任何业务代码
* `@/lib` 可被 services/hooks 引用，禁止反向依赖

---

## 7. 开发工作流

通过 `Makefile` 执行所有操作：

| 命令 | 说明 |
|------|------|
| `make dev` | 启动开发服务器 |
| `make test` | 运行单元测试 |
| `make test-watch` | 测试 watch 模式 |
| `make check` | 一键全检查（lint + format + type-check） |
| `make build` | 生产构建 |
| `make ci` | CI 完整流水线 |

---

## 8. Git 规范

* **格式:** `<type>(<scope>): <subject>` — feat, fix, refactor, test, docs, chore
* **提交前:** 必须通过 `make check`（Husky pre-commit 自动执行 format）
* **Scope:** 取自 feature 目录名，如 `feat(auth): add login page`

---

## 9. 禁止行为 (Zero Tolerance)

* ❌ 组件内直接调用 axios
* ❌ 使用 `any` 类型
* ❌ `useEffect` 中获取数据
* ❌ `eslint-disable` 规避 hook 依赖检查
* ❌ `console.log`（使用结构化日志工具）
* ❌ 未达 3 次复用就抽象公共组件
* ❌ 服务端数据存入 Zustand

---

## 10. 成功标准

* ✅ **合宪性:** 是否符合 constitution.md
* ✅ **类型闭环:** 运行时数据是否与 TS 类型严格匹配（Zod 校验）
* ✅ **测试覆盖:** 核心逻辑 Hook 必须有单元测试
* ✅ **无副作用:** 组件渲染函数内无数据获取和全局状态修改

---

🚨 **若无法满足上述任一标准，请立即停止并向用户报告架构冲突。**
