# aitestos

AI-powered intelligent test case management platform. Built with Go 1.26, DDD architecture, and integrated LLM capabilities for automated test case generation from PRDs, API specs, and design documents.

## Features

- **AI Test Generation** - Leverage DeepSeek LLM + RAG to auto-generate test cases from requirement documents
- **Knowledge Base** - Upload PRDs, Figma exports, OpenAPI specs; automatic chunking & vectorization via Milvus
- **Test Lifecycle** - Full CRUD for test cases, test plans, execution tracking with pass/fail/block/skip
- **Project Isolation** - Multi-project, multi-module architecture with per-project configuration
- **Observability** - Structured logging (zerolog), Prometheus metrics, OpenTelemetry tracing
- **Graceful Degradation** - Redis cache with in-memory fallback; LLM/RAG/Milvus optional
- **Production Ready** - Distroless Docker image (~20MB), graceful shutdown, configurable timeouts

## Architecture

```
cmd/server/          Entry point (graceful shutdown)
internal/
  app/               DI & lifecycle orchestration
  config/            YAML + env config (Viper)
  domain/            Core models & interfaces (no external deps)
    identity/        User, auth, roles
    project/         Project, module, config
    testcase/        Test case, versioning
    testplan/        Test plan, execution results
    knowledge/       Document, chunks, vector search
    generation/      AI generation tasks & drafts
  service/           Business logic
  repository/        Data persistence (PostgreSQL + Redis cache)
  transport/http/    HTTP handlers, middleware, routing
  infrastructure/    External integrations (LLM, Milvus, Redis, MinIO)
  ierrors/           Unified error codes
api/                 OpenAPI 3.0 contracts
scripts/             Database migrations
tests/integration/   Integration tests (testcontainers)
```

## Quick Start

### Prerequisites

- Go 1.26+
- Docker & Docker Compose (for infrastructure)
- PostgreSQL 15+, Redis 7+ (or use Docker)

### 1. Clone & Install Tools

```bash
git clone https://github.com/liang21/aitestos.git
cd aitestos/backend
make install-tools
```

### 2. Start Infrastructure

```bash
make docker-up
```

This starts PostgreSQL, Redis, Milvus, MinIO, and RabbitMQ via Docker Compose.

### 3. Configure

```bash
cp configs/config.yaml configs/config.local.yaml
# Edit config.local.yaml with your settings
# Or set environment variables:
export DEEPSEEK_API_KEY=your_api_key
export JWT_SECRET=your_jwt_secret_min_8_chars
```

### 4. Run Database Migrations

```bash
make migrate-up
```

### 5. Start the Server

```bash
make dev    # Docker infra + app
# or
make run    # Build & run only
```

Server starts at `http://localhost:8080`.

## API Endpoints

### System

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |

### Authentication

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/register` | Register user |
| POST | `/auth/login` | Login (returns JWT) |
| POST | `/auth/refresh` | Refresh access token |

### Projects

| Method | Path | Description |
|--------|------|-------------|
| POST | `/projects` | Create project |
| GET | `/projects` | List projects |
| GET | `/projects/{id}` | Get project with stats |
| PUT | `/projects/{id}` | Update project |
| DELETE | `/projects/{id}` | Delete project |
| GET | `/projects/{id}/statistics` | Project statistics |
| POST | `/projects/{id}/modules` | Create module |
| GET | `/projects/{id}/modules` | List modules |
| POST | `/projects/{id}/config` | Set config |
| GET | `/projects/{id}/config` | List configs |
| POST | `/projects/{id}/config/import` | Import configs |
| POST | `/projects/{id}/config/export` | Export configs |

### Test Cases

| Method | Path | Description |
|--------|------|-------------|
| POST | `/testcases` | Create test case |
| GET | `/testcases` | List (filter by module/project/status/type) |
| GET | `/testcases/{id}` | Get case detail |
| PUT | `/testcases/{id}` | Update test case |
| DELETE | `/testcases/{id}` | Delete test case |

### Test Plans

| Method | Path | Description |
|--------|------|-------------|
| POST | `/testplans` | Create test plan |
| GET | `/testplans` | List test plans |
| GET | `/testplans/{id}` | Get plan with results |
| PUT | `/testplans/{id}` | Update plan |
| POST | `/testplans/{id}/results` | Record execution result |

### AI Generation

| Method | Path | Description |
|--------|------|-------------|
| POST | `/generation/tasks` | Create generation task |
| GET | `/generation/tasks` | List tasks |
| GET | `/generation/drafts` | Get generated drafts |
| POST | `/generation/drafts/{id}/confirm` | Approve draft -> test case |
| POST | `/generation/drafts/{id}/reject` | Reject draft |
| POST | `/generation/batch-confirm` | Batch approve drafts |

### Knowledge Management

| Method | Path | Description |
|--------|------|-------------|
| POST | `/documents` | Upload document |
| GET | `/documents` | List documents |
| GET | `/documents/{id}` | Get document |
| DELETE | `/documents/{id}` | Delete document |
| POST | `/documents/{id}/process` | Chunk & vectorize |

## Makefile Reference

| Command | Description |
|---------|-------------|
| `make all` | Format, vet, lint, test and build |
| `make build` | Build binary to `bin/aitestos` |
| `make run` | Build and run the server |
| `make dev` | Start Docker infra + run app |
| `make test` | Run unit tests with coverage |
| `make test-integration` | Run integration tests |
| `make fmt` | Format code with `gofmt -s` |
| `make vet` | Run `go vet` |
| `make lint` | Run `golangci-lint` |
| `make check` | Run all quality checks (fmt + vet + lint + test) |
| `make tidy` | Run `go mod tidy` |
| `make generate` | Run `go generate` (Wire, mocks) |
| `make migrate-up` | Apply database migrations |
| `make migrate-down` | Rollback migrations |
| `make docker-up` | Start infrastructure containers |
| `make docker-full` | Start full stack (infra + app) |
| `make docker-down` | Stop all containers |
| `make docker-logs` | Tail container logs |
| `make install-tools` | Install golangci-lint, wire, mockery |
| `make clean` | Remove build artifacts |

## Configuration

Configuration is loaded from `configs/config.yaml` with environment variable substitution:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  shutdown_timeout: 30s

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  database: "aitestos"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

llm:
  provider: "deepseek"
  api_key: "${DEEPSEEK_API_KEY}"
  base_url: "https://api.deepseek.com"
  model: "deepseek-chat"
  embedding_model: "deepseek-embedding"
  timeout: 60s
  max_retries: 3

milvus:
  host: "localhost"
  port: 19530
  database: "aitestos"
  collection: "document_chunks"

jwt:
  secret: "${JWT_SECRET}"
  expire_time: 2h

log:
  level: "info"
  json: true
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DEEPSEEK_API_KEY` | Yes* | DeepSeek LLM API key |
| `JWT_SECRET` | Yes | JWT signing secret (min 8 chars) |
| `MINIO_ACCESS_KEY` | No | MinIO access key |
| `MINIO_SECRET_KEY` | No | MinIO secret key |
| `DB_PASSWORD` | No | PostgreSQL password |
| `REDIS_PASSWORD` | No | Redis password |

*Required for AI features; the server starts without it but generation endpoints return errors.

## Docker

Build and run with Docker:

```bash
docker build \
  --build-arg VERSION=$(git describe --tags --always) \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  -t aitestos:latest .

docker compose --profile app up -d
```

Multi-stage build produces a ~20MB distroless image with static binary.

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.26 |
| Router | Chi v5 |
| Database | PostgreSQL (sqlx) |
| Cache | Redis with in-memory fallback |
| Vector DB | Milvus |
| LLM | DeepSeek API |
| Auth | JWT (golang-jwt/v5) |
| Logging | Zerolog (structured, JSON) |
| Metrics | Prometheus |
| Config | Viper (YAML + env vars) |
| DI | Google Wire |
| Container | Distroless multi-stage |

## License

MIT
