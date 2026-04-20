# ==============================================================================
# Aitestos Frontend — 工程管理控制台
# React 19 · Vite 8 · Yarn 1.x · TypeScript 5.9
# ==============================================================================

.DEFAULT_GOAL := help

# --- 元数据 ---
APP_NAME    := aitest-frontend
VERSION     := $(shell node -p "require('./package.json').version")
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
IMAGE_TAG   ?= latest

# --- 依赖管理 ---

.PHONY: setup
setup: ## 安装依赖（lockfile 一致性校验）
	yarn install --frozen-lockfile

.PHONY: tidy
tidy: ## 清理缓存并重校验依赖
	yarn install --check-files
	yarn cache clean

# --- 开发 ---

.PHONY: dev
dev: ## 启动 Vite 开发服务器 (http://localhost:5173)
	yarn dev

# --- 代码质量 ---

.PHONY: type-check
type-check: ## TypeScript 类型检查 (strict mode)
	yarn tsc -b

.PHONY: lint
lint: ## ESLint 检查
	yarn lint

.PHONY: format
format: ## Prettier 格式化
	yarn format

.PHONY: format-check
format-check: ## Prettier 格式检查（不写入）
	yarn format:check

.PHONY: check
check: lint format-check type-check ## 一键全检查（lint + format + type-check）

# --- 测试 ---

.PHONY: test
test: ## 运行单元测试 (Vitest)
	yarn vitest run

.PHONY: test-watch
test-watch: ## 运行测试 (watch 模式)
	yarn vitest

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	yarn vitest run --coverage

.PHONY: test-e2e
test-e2e: ## 运行 Playwright E2E 测试
	npx playwright test

# --- 构建 ---

.PHONY: build
build: type-check ## 生产构建（注入版本元数据）
	VITE_APP_VERSION=$(VERSION) VITE_APP_COMMIT=$(COMMIT_HASH) VITE_APP_BUILD_TIME=$(BUILD_TIME) yarn build

.PHONY: preview
preview: build ## 预览生产构建
	yarn preview

# --- Docker ---

.PHONY: docker-build
docker-build: ## 使用 BuildKit 构建镜像
	DOCKER_BUILDKIT=1 docker build \
		--build-arg BUILD_METADATA='{"version":"$(VERSION)","commit":"$(COMMIT_HASH)","buildTime":"$(BUILD_TIME)"}' \
		-t $(APP_NAME):$(IMAGE_TAG) .

# --- CI ---

.PHONY: ci
ci: setup check test ## CI 完整流水线（install → check → test）

# --- 清理 ---

.PHONY: clean
clean: ## 清理构建产物、缓存
	rm -rf dist coverage node_modules/.vite

.PHONY: clean-all
clean-all: clean ## 深度清理（含 node_modules）
	rm -rf node_modules

# --- 帮助 ---

.PHONY: help
help: ## 显示此帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'
