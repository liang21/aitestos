# 项目开发宪法
# Version: 1.0 | Scope: React SPA & Frontend Engineering

本文件定义了本项目不可动摇的核心开发原则。所有 AI Agent 在进行技术规划和代码实现时，必须无条件遵循。

---

## 第一条：简单性原则 (Simplicity First)

* **1.1 YAGNI:** 严禁实现 `spec.md` 之外的任何预测性功能。
* **1.2 技术栈一致性:** 优先使用项目已选定的技术栈（见第六条），不引入同类替代库。React 内置能力（useState、useReducer、Suspense）作为局部状态和加载态的首选。
* **1.3 拒绝过度抽象:** 在没有 3 个复用场景前，禁止提取公共 Hook/组件。接口应由消费者定义，而非预先设计。

---

## 第二条：测试先行铁律 (Test-First Imperative)

* **2.1 红绿循环:** 修改前必须先提供一个能复现问题或覆盖新功能的失败测试（Red）。
* **2.2 用户视角:** 测试必须使用 Testing Library（getByRole / getByLabelText / userEvent）模拟真实用户交互，禁止测试实现细节。
* **2.3 真实依赖:** API 层使用 MSW 拦截 HTTP 请求，严禁直接 mock axios 实例或 React Query。

---

## 第三条：类型安全 (Type Safety)

* **3.1 Strict 模式:** `tsconfig.json` 必须开启 `strict: true`。
* **3.2 禁止 any:** 全项目禁止使用 `any`。Axios 泛型使用 `<never, ResponseType>` 替代 `<any, T>`。
* **3.3 运行时校验:** 表单输入使用 Zod schema 校验。API 响应类型由 TypeScript + Axios 泛型（`<never, ResponseType>`）保障编译期安全。

---

## 第四条：数据流纪律 (Data Flow Discipline)

* **4.1 React Query 唯一数据源:** 所有服务端数据必须通过 React Query（useQuery/useMutation）获取和变更。禁止 useEffect + fetch/useEffect + axios。
* **4.2 Zustand 最小化:** 仅 auth token、sidebar 折叠等跨 Feature UI 状态使用 Zustand。服务端数据严禁存入 Zustand。
* **4.3 副作用隔离:** 业务逻辑必须封装在 hooks 中。UI 组件禁止包含数据获取（useEffect + fetch）、订阅、定时器等副作用代码。

---

## 第五条：组件设计 (Component Design)

* **5.1 单一职责:** Page 组件编排布局，Feature Component 处理交互，Shared Component 纯 UI（仅接收 props）。
* **5.2 禁止内联样式:** 必须使用 Tailwind CSS 工具类或 Arco Design 组件样式。
* **5.3 可访问性:** 必须使用语义化 HTML，交互元素必须包含 aria-label，key 禁止使用 array index。

---

## 第六条：错误边界与健壮性 (Error Boundary & Robustness)

* **6.1 ErrorBoundary:** 每个 Feature 的根页面组件必须被 ErrorBoundary 包裹，防止局部崩溃导致白屏。
* **6.2 React Query onError:** 所有 useMutation 必须提供 onError 回调，使用 Arco Message 展示用户友好错误提示。禁止静默吞掉错误。
* **6.3 全局拦截:** Axios 响应拦截器统一处理 401（Token 刷新）、403、5xx。组件层仅处理业务级错误（如"名称已存在"）。

---

## 第七条：工程一致性 (Mandatory Stack)

| 领域 | 强制技术 | 禁止替代 |
|------|----------|----------|
| 数据获取 | React Query v5 | useEffect + fetch/axios |
| 全局状态 | Zustand | Redux、Context 替代全局状态 |
| 表单 | React Hook Form + Zod | 手动 useState 管理表单 |
| 样式 | Tailwind CSS + Arco Design | inline-style、CSS Modules |
| 路由 | React Router v7 | 其他路由方案 |
| HTTP | Axios（仅作为 queryFn 底层） | 组件内直接调用 axios |
| 测试 | Vitest + Testing Library + MSW | Jest |

---

## 第八条：代码质量 (Code Quality)

* **8.1 提交前检查:** 必须通过 `make check`（lint + format + type-check）。
* **8.2 格式化:** 提交前必须运行 `make format`。Husky pre-commit 已自动执行。
* **8.3 依赖整洁:** 安装新依赖后必须验证 `yarn.lock` 同步更新。

---

## 治理 (Governance)

本宪法效力高于任何单次会话指令。若指令违宪，AI 必须提出质疑并拒绝执行。
