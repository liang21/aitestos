# Aitestos Frontend

AI 驱动的测试用例生成与管理平台前端。支持项目管理、测试用例 CRUD、AI 生成测试用例、草稿箱审批、测试计划执行、知识库管理及项目配置等功能。

## 技术栈

| 技术                 | 版本 | 用途                    |
| -------------------- | ---- | ----------------------- |
| React                | 19.2 | UI 框架                 |
| TypeScript           | 5.9  | 类型系统 (strict: true) |
| React Router         | 7.x  | 路由 (lazy loading)     |
| TanStack React Query | 5.x  | 服务端状态管理          |
| Zustand              | 5.x  | 客户端全局状态          |
| React Hook Form      | 7.x  | 表单管理                |
| Zod                  | 4.x  | 运行时校验              |
| Arco Design          | 2.66 | UI 组件库               |
| Tailwind CSS         | 4.x  | 样式                    |
| Axios                | 1.14 | HTTP 客户端             |
| Vitest               | 4.x  | 单元测试                |
| Playwright           | 1.59 | E2E 测试                |
| MSW                  | 2.x  | API Mock                |

## 目录架构

```
src/
├── app/                    # 应用入口
│   ├── App.tsx             # 根组件（Provider 嵌套）
│   └── main.tsx            # 渲染入口
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
│   ├── auth/               # 认证（登录、注册）
│   ├── projects/           # 项目管理
│   ├── testcases/          # 测试用例
│   ├── plans/              # 测试计划
│   ├── generation/         # AI 生成
│   ├── drafts/             # 草稿箱
│   ├── documents/          # 知识库
│   ├── modules/            # 模块管理
│   └── configs/            # 项目配置
├── components/             # 跨 Feature 共享组件
│   ├── layout/             # AppLayout, Sidebar, Header
│   └── business/           # StatusTag, SearchTable, ArrayEditor, StatsCard 等
├── store/                  # 全局 Zustand store
│   └── useAppStore.ts      # sidebarCollapsed + notifications
├── hooks/                  # 全局共享 hooks
└── styles/                 # 全局样式
```

每个 feature 模块遵循统一结构：

```
src/features/<name>/
├── components/     # 页面与业务组件
├── hooks/          # useFeature hooks（业务逻辑）
└── services/       # API 调用层
```

## 快速开始

### 前置要求

- Node.js >= 20
- Yarn 1.x（Classic）

### 安装依赖

```bash
yarn install
```

### 开发模式

```bash
make dev
# 或
yarn dev
```

访问 http://localhost:5173

### 构建部署

```bash
make build
# 或
yarn build
```

构建产物注入版本元数据（`VITE_APP_VERSION`、`VITE_APP_COMMIT`、`VITE_APP_BUILD_TIME`）。

## 开发指南

### 核心原则

- **Hook 优先**: 禁止 class component，业务逻辑封装在 `useFeature.ts` 中
- **Feature-Based 架构**: 模块间禁止互相引用内部文件
- **API 边界**: 组件禁止直接 import axios，必须通过 `services/` 层
- **服务端数据走 React Query**: 禁止存入 Zustand，禁止 `useEffect` 中获取数据
- **严格类型**: 禁止 `any` 类型，运行时数据使用 Zod 校验

完整规范参见 [CLAUDE.md](./CLAUDE.md)。

### 依赖规则

```
Page → Feature hooks → Feature services → @/lib/request
components/ 禁止引用 features/ 的任何内容
services/ 禁止引用 store/ 或 hooks/
Feature A 禁止引用 Feature B 的内部文件
```

### Git 提交规范

```
<type>(<scope>): <subject>
```

- **type**: feat, fix, refactor, test, docs, chore, style, build
- **scope**: feature 目录名（如 auth, projects, configs）
- **subject**: 祈使句，首字母小写，不加句号

示例：`feat(configs): add config import with validation`

提交前自动执行 `make check`（Husky pre-commit hook）。

### 测试

```bash
make test              # 运行单元测试 (Vitest)
make test-watch        # watch 模式
make test-coverage     # 生成覆盖率报告
make test-e2e          # 运行 E2E 测试 (Playwright)
```

测试文件与源文件同目录，匹配 `src/**/*.test.{ts,tsx}`。API Mock 使用 MSW，handlers 集中管理在 `tests/msw/handlers/`。

### 代码检查

```bash
make check             # lint + format check + type-check
make lint              # ESLint
make format            # Prettier 格式化
make type-check        # TypeScript 类型检查
```

## Makefile 命令参考

| 命令                 | 说明                                    |
| -------------------- | --------------------------------------- |
| `make dev`           | 启动 Vite 开发服务器                    |
| `make test`          | 运行单元测试 (Vitest)                   |
| `make test-watch`    | 测试 watch 模式                         |
| `make test-coverage` | 生成覆盖率报告                          |
| `make test-e2e`      | 运行 Playwright E2E 测试                |
| `make build`         | 生产构建（含 type-check）               |
| `make preview`       | 预览生产构建                            |
| `make check`         | 一键全检查 (lint + format + type-check) |
| `make lint`          | ESLint 检查                             |
| `make format`        | Prettier 格式化                         |
| `make type-check`    | TypeScript 类型检查                     |
| `make setup`         | 安装依赖（lockfile 一致性校验）         |
| `make ci`            | CI 完整流水线 (install → check → test)  |
| `make clean`         | 清理构建产物和缓存                      |
| `make clean-all`     | 深度清理（含 node_modules）             |
| `make help`          | 显示帮助信息                            |

## 环境变量

复制 `.env.example` 为 `.env` 并按需修改：

```bash
cp .env.example .env
```

| 变量名              | 说明         | 默认值    |
| ------------------- | ------------ | --------- |
| `VITE_API_BASE_URL` | API 基础路径 | `/api/v1` |

本地开发后端运行在 8080 端口时，设置为 `http://localhost:8080/api/v1`。

## 许可证

Copyright © 2026 Aitestos
