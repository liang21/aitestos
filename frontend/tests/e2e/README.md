# E2E 测试使用指南

## 概述

本目录包含 Aitestos 前端的端到端 (E2E) 测试，使用 Playwright 框架编写。

## 目录结构

```
tests/e2e/
├── fixtures/           # 测试 fixtures 和扩展
├── helpers/           # 辅助工具（配置、工厂函数、清理）
├── pages/             # Page Object Model
├── *.spec.ts          # 测试文件
├── global-setup.ts    # 全局测试设置
└── README.md          # 本文档
```

## 运行测试

### 运行所有 E2E 测试

```bash
npm run test:e2e
```

### 运行特定测试文件

```bash
npx playwright test auth.spec.ts
```

### 运行特定测试

```bash
npx playwright test -g "should login successfully"
```

### 带界面的交互模式

```bash
npm run test:e2e:ui
```

### 调试模式

```bash
npm run test:e2e:debug
```

### 指定浏览器

```bash
npx playwright test --project chromium    # Chrome
npx playwright test --project firefox     # Firefox
npx playwright test --project webkit      # Safari
```

### 查看测试报告

测试运行后，HTML 报告会生成在 `playwright-report/index.html`。

```bash
npx playwright show-report
```

## Page Object Model

我们使用 Page Object Model (POM) 模式来组织测试代码，提高可维护性。

### 基础用法

```typescript
import { LoginPage } from './pages/LoginPage'

test('should login', async ({ page }) => {
  const loginPage = new LoginPage(page)
  await loginPage.goto()
  await loginPage.login('user@example.com', 'password')
})
```

### 可用的 Page Objects

- `LoginPage` - 登录页面
- `ProjectListPage` - 项目列表页
- `CreateProjectModal` - 创建项目弹窗

## 测试数据工厂

使用工厂函数生成一致的测试数据：

```typescript
import { createProject, createModule } from './helpers/factories'

// 使用默认值
const project = createProject()

// 覆盖特定字段
const project = createProject({ name: 'Custom Name' })
```

## 测试数据清理

测试数据会自动标记并在测试结束后清理：

```typescript
import { generateCleanupId, cleanupTestData } from './helpers/cleanup'

// 生成带清理 ID 的数据
const cleanupId = generateCleanupId('project')
const project = createProject({ name: `Project ${cleanupId}` })

// 在测试后清理
test.afterAll(async ({ page }) => {
  await cleanupTestData(page)
})
```

## 配置

### 环境变量

- `BASE_URL` - 应用基础 URL (默认: `http://localhost:5174`)

### 测试配置

配置文件：`playwright.config.ts`

```typescript
// 修改默认超时
use: {
  actionTimeout: 10000,
  navigationTimeout: 30000,
}

// 修改重试次数
retries: 2
```

## 编写新测试

### 1. 创建测试文件

在 `tests/e2e/` 目录下创建 `.spec.ts` 文件：

```typescript
import { test, expect } from '@playwright/test'
import { LoginPage } from './pages/LoginPage'

test.describe('Feature Name', () => {
  test('should do something', async ({ page }) => {
    // 测试代码
  })
})
```

### 2. 使用 Page Objects

```typescript
import { ProjectListPage } from './pages/ProjectListPage'

const projectListPage = new ProjectListPage(page)
await projectListPage.goto()
```

### 3. 使用测试工厂

```typescript
import { createProject } from './helpers/factories'

const project = createProject()
```

### 4. 添加 data-testid 选择器

在组件中添加 `data-testid` 属性：

```typescript
<Button data-testid="login-submit-button">登录</Button>
```

然后在测试中使用：

```typescript
await page.click('[data-testid="login-submit-button"]')
```

## 最佳实践

### ✅ DO

- 使用 `data-testid` 选择器而不是文本选择器
- 使用 Page Object Model 封装页面逻辑
- 使用工厂函数生成测试数据
- 在 `test.afterAll` 中清理测试数据
- 使用 `waitForSelector` 等待元素可见
- 使用描述性的测试名称

### ❌ DON'T

- 使用硬编码的等待时间（`page.waitForTimeout`）
- 使用可能因国际化而改变的文本选择器
- 在测试中直接使用魔法数字/字符串
- 忽略测试数据清理
- 编写过长的测试（超过 100 行考虑拆分）

## 故障排查

### 测试超时

如果测试超时，检查：

1. 应用是否正在运行
2. BASE_URL 是否正确
3. 选择器是否正确

### 元素找不到

如果出现 "Element not found" 错误：

1. 使用 `--debug` 模式运行查看实际页面
2. 检查元素是否有 `data-testid`
3. 检查元素是否在 iframe 中

### 数据冲突

如果测试间有数据冲突：

1. 确保使用 `generateCleanupId` 生成唯一 ID
2. 检查 `test.afterAll` 清理是否正确执行

## CI/CD 集成

在 CI 环境中运行：

```bash
# 启动应用
npm run dev &

# 等待应用就绪
sleep 10

# 运行 E2E 测试
npm run test:e2e

# 上传测试报告
# (根据 CI 平台配置)
```

## 参考资源

- [Playwright 官方文档](https://playwright.dev)
- [Page Object Model 最佳实践](https://playwright.dev/docs/pom)
- [测试最佳实践](https://playwright.dev/docs/best-practices)
