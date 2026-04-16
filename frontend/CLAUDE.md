# CLAUDE.md - 前端 AI 工程深度操作手册

---

# 1️⃣ 指令优先级

1. **宪法优先:** 读取 `@constitution.md`
2. **决策驱动:** 执行前必须输出架构决策记录 (ADR)
3. **TDD First:** 先生成 MSW 定义与测试用例，再生成实现

---

# 2️⃣ React 专家规范 (React 18+)

## 技术栈 (强制)
* React 18 (Suspense/Transition)
* TypeScript (Strict)
* React Query (v5+)
* Zustand
* Tailwind CSS
* Zod (Runtime Validation)

## 代码风格
* **Hook 优先:** 禁止 class components。
* **解耦:** UI 组件仅负责展示，业务逻辑、数据转换必须位于 `useFeatureName.ts`。

---

# 3️⃣ Agent 思考链 (CoT) 要求

在执行任何代码编写前，Agent 必须按顺序输出：
1. **Spec 对齐:** 确认功能点。
2. **ADR (Architecture Decision Record):** - 状态放在哪里？为什么？
   - 关键 Hook 依赖是什么？
3. **测试方案:** 列出关键路径的测试点。

---

# 4️⃣ 数据流与安全性

* **数据获取:** 必须封装为 `useQuery` / `useMutation` 自定义 Hook。
* **输入校验:** 表单必须使用 `react-hook-form` + `zod`。
* **API 边界:** `services/` 下的方法必须返回经 Zod 校验后的对象。

---

# 5️⃣ 目录结构约束

```text
src/
├── api/                        # API 请求模块
│   ├── auth.ts                 # 认证接口
│   ├── projects.ts             # 项目管理接口
│   ├── modules.ts              # 模块管理接口
│   ├── configs.ts              # 项目配置接口
│   ├── testcases.ts            # 测试用例接口
│   ├── plans.ts                # 测试计划接口
│   ├── generation.ts           # AI 生成接口
│   └── documents.ts            # 知识库文档接口
├── assets/                     # 静态资源
│   └── images/                 # 图片资源
├── components/                 # 通用组件
│   ├── business/               # 业务组件
│   │   ├── StatusTag.tsx        # 状态标签（统一色彩映射）
│   │   ├── ArrayEditor.tsx      # 数组编辑器（前置条件/步骤）
│   │   ├── SearchTable.tsx      # 搜索筛选表格
│   │   ├── StatsCard.tsx        # 统计卡片
│   │   ├── SplitPanel.tsx       # 分栏面板
│   │   ├── ReferencePanel.tsx   # 引用来源面板
│   │   ├── CaseSelector.tsx     # 用例选择器
│   │   └── JsonEditor.tsx       # JSON 编辑器
│   └── layout/                 # 布局组件
│       ├── AppLayout.tsx        # 应用主布局
│       ├── Sidebar.tsx          # 侧边栏导航
│       ├── Header.tsx           # 顶部导航栏
│       └── AuthLayout.tsx       # 认证页布局
├── hooks/                      # 自定义 Hooks
│   ├── useAuth.ts              # 认证 Hook
│   ├── useProject.ts           # 当前项目 Hook
│   └── usePolling.ts           # 轮询 Hook
├── pages/                      # 页面组件
│   ├── auth/                   # 认证页面
│   │   ├── LoginPage.tsx
│   │   └── RegisterPage.tsx
│   ├── projects/               # 项目页面
│   │   ├── ProjectListPage.tsx
│   │   ├── ProjectDashboard.tsx
│   │   └── CreateProjectModal.tsx
│   ├── knowledge/              # 知识库页面
│   │   ├── KnowledgeListPage.tsx
│   │   ├── DocumentDetailPage.tsx
│   │   ├── FigmaIntegrationPage.tsx
│   │   └── UploadDocumentModal.tsx
│   ├── generation/             # AI 生成页面
│   │   ├── GenerationTaskListPage.tsx
│   │   ├── NewGenerationTaskPage.tsx
│   │   └── TaskDetailPage.tsx
│   ├── drafts/                 # 草稿箱页面
│   │   ├── DraftListPage.tsx
│   │   └── DraftConfirmPage.tsx
│   ├── cases/                  # 测试用例页面
│   │   ├── CaseListPage.tsx
│   │   ├── CaseDetailPage.tsx
│   │   └── CreateCaseDrawer.tsx
│   ├── plans/                  # 测试计划页面
│   │   ├── PlanListPage.tsx
│   │   ├── NewPlanPage.tsx
│   │   ├── PlanDetailPage.tsx
│   │   └── ResultRecordModal.tsx
│   └── settings/               # 项目设置页面
│       ├── ProjectSettingsPage.tsx
│       ├── ModuleManagePage.tsx
│       └── ConfigManagePage.tsx
├── router/                     # 路由配置
│   ├── index.tsx               # 路由定义
│   └── RouteGuard.tsx          # 路由守卫
├── store/                      # Zustand 状态管理
│   ├── useAuthStore.ts         # 认证状态
│   ├── useAppStore.ts          # 全局 UI 状态
│   ├── useProjectStore.ts      # 当前项目状态
│   └── useDraftStore.ts        # 草稿箱状态
├── styles/                     # 全局样式
│   ├── theme.css               # 主题变量
│   └── global.css               # 全局样式
├── types/                      # TypeScript 类型定义
│   ├── api.ts                  # API 请求/响应类型
│   ├── enums.ts                # 枚举类型
│   └── models.ts               # 数据模型类型
├── lib/                        # 工具库
│   ├── request.ts              # Axios 实例
│   └── utils.ts                # 通用工具函数
├── App.tsx                     # 应用入口
├── main.tsx                    # 渲染入口
└── index.css                   # 全局 CSS
```

---

# 6️⃣ 性能与异常治理

* **性能:** 关键路径必须使用 `React.Suspense` 处理 Loading。
* **健壮性:** 每个 Feature 根组件必须包裹 `ErrorBoundary`。
* **日志:** 禁止 `console.log`，必须使用项目统一的 `logger` 工具。

---

# 7️⃣ Git & 提交规范

* 格式：`<type>: <description>` (feat, fix, refactor, test, docs)。
* 提交前必须通过：`pnpm type-check` & `pnpm test`。

---

# 8️⃣ 禁止行为 (Zero Tolerance)

* ❌ **使用 any**: 哪怕是 unknown + guard 也不要用 any。
* ❌ **隐式副作用**: 在渲染周期内修改全局变量或 Ref。
* ❌ **过度封装**: 在没有 3 个复用案例前，不要创建“通用”组件。
* ❌ **裸 Fetch**: 禁止直接调用 fetch，必须经过 `services/` 层的 Zod 校验。

---

# 🎯 成功标准

* ✅ **合宪性:** 是否符合 constitution.md。
* ✅ **可回溯:** 是否有 ADR 记录。
* ✅ **类型闭环:** 运行时数据是否与 TS 类型严格匹配。
* ✅ **100% 测试覆盖:** 核心逻辑 Hook 必须有 Unit Test。

---
🚨 **若无法满足上述任一标准，请立即停止并向用户报告架构冲突。**
```

---

### 👨‍💻 架构师执行建议：

1.  **自动化治理：** 建议在项目中配置 `husky`，在 `pre-commit` 阶段强制运行 `pnpm type-check` 和 `pnpm test`，配合这套文档实现“流程闭环”。
2.  **AI 交互提示：** 之后你每次开启新对话，可以直接把这两份文件喂给 AI，并告诉它：“你是这个项目的架构师，请按照 `CLAUDE.md` 的工作流，并在不违反 `constitution.md` 的前提下开始任务。”