# 前端项目开发宪法 (Ultimate Edition)

# Version: 4.0 | Scope: React 18+ & Enterprise Architecture

本文件定义了项目不可动摇的核心开发原则。所有 AI Agent 与开发者必须无条件遵循。

---

# 第一条：简单性与并发安全 (Simplicity & Concurrent Safety)

* **1.1 YAGNI:** 严禁实现 `spec.md` 之外的任何功能。
* **1.2 渲染纯净性:** 严禁在渲染过程中产生副作用（如修改外部变量）。
* **1.3 拒绝过度抽象:** ≥3 个复用场景才允许抽象 Hook/组件。
* **1.4 状态本地化:** * URL 状态 (Search Params) 优先。
    * Server State (React Query) 为主。
    * UI Global State (Zustand) 仅用于跨 Feature 共享。

---

# 第二条：测试先行铁律 (TDD First)

* **2.1 红绿循环:** 修改前必须提供失败测试。
* **2.2 用户视角:** 必须使用 Testing Library (getByRole / userEvent) 模拟真实交互。
* **2.3 真实依赖:** API 必须使用 MSW 模拟，严禁直接 mock fetch/axios。
* **2.4 异步闭环:** useEffect 必须返回清理函数，并处理竞态条件（使用 AbortController）。

---

# 第三条：全链路类型安全 (Type Safety)

* **3.1 Strict 模式:** tsconfig 必须开启 `strict: true`。
* **3.2 运行期校验:** 所有 API Response 必须使用 Zod 定义 Schema 并进行运行时校验。
* **3.3 禁止滥用 as:** 除非在处理 legacy 代码或无法避免的外部库类型。

---

# 第四条：副作用隔离 (Explicit Side Effects)

* **4.1 容器/逻辑解耦:** 逻辑必须封装在 hooks 中，UI 组件内严禁编写超过 15 行的逻辑运算。
* **4.2 依赖诚实:** 严禁使用 `eslint-disable` 规避 hook 依赖检查。

---

# 第五条：组件设计与可访问性 (Design & A11y)

* **5.1 组合优于继承:** 优先使用 Render Props 或 Composition。
* **5.2 A11y 强制性:** 必须使用语义化 HTML（<main>, <section>），严禁 <div> 堆叠；必须包含 aria-label。
* **5.3 稳定性:** key 严禁使用 index。

---

# 第六条：工程一致性 (Mandatory Stack)

* **6.1 状态管理:** Zustand (唯一全局库)。
* **6.2 数据流:** React Query (必须启用 Suspense)。
* **6.3 样式:** Tailwind CSS (严禁 inline-style)。
* **6.4 校验:** Zod / Valibot。

---

# 第七条：架构组织

* **7.1 Feature-Based:** 严格遵守 `features/{name}/{components,hooks,services,types,schema.ts}`。
* **7.2 命名规范:** 组件 PascalCase，Hook useXXX，类型 .types.ts。

---

# 治理 (Governance)

本宪法效力高于任何单次指令。若 AI 输出违反宪法，**必须拒绝执行并指出违规点。**