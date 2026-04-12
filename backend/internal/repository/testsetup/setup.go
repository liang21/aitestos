package testsetup

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "github.com/lib/pq"
)

// PostgresContainer 封装 PostgreSQL 容器配置
type PostgresContainer struct {
	container testcontainers.Container
	ConnStr   string
}

// SetupPostgres 创建并启动 PostgreSQL 测试容器
func SetupPostgres(ctx context.Context) (*PostgresContainer, error) {
	const (
		dbname   = "aitestos_test"
		user     = "testuser"
		password = "testpass"
		port     = "5432"
	)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{port + "/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       dbname,
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("create postgres container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, port)
	if err != nil {
		return nil, fmt.Errorf("get mapped port: %w", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, mappedPort.Port(), user, password, dbname,
	)

	return &PostgresContainer{
		container: container,
		ConnStr:   connStr,
	}, nil
}

// CreateDB 从容器连接创建数据库连接池
func (pc *PostgresContainer) CreateDB(ctx context.Context) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", pc.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// Teardown 清理容器
func (pc *PostgresContainer) Teardown(ctx context.Context) error {
	if pc.container != nil {
		return pc.container.Terminate(ctx)
	}
	return nil
}

// RunMigrations 执行数据库迁移
func RunMigrations(ctx context.Context, db *sqlx.DB) error {
	// 启用 uuid-ossp 扩展（必须在创建表之前，因为表使用了 uuid_generate_v4()）
	if _, err := db.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		return fmt.Errorf("create uuid extension: %w", err)
	}

	// 创建 ENUM 类型
	enumSQL := []string{
		`CREATE TYPE IF NOT EXISTS case_status_enum AS ENUM ('unexecuted', 'pass', 'block', 'fail');`,
		`CREATE TYPE IF NOT EXISTS case_type_enum AS ENUM ('functionality', 'performance', 'api', 'ui', 'security');`,
		`CREATE TYPE IF NOT EXISTS plan_status_enum AS ENUM ('draft', 'active', 'completed', 'archived');`,
		`CREATE TYPE IF NOT EXISTS priority_enum AS ENUM ('P0', 'P1', 'P2', 'P3');`,
		`CREATE TYPE IF NOT EXISTS result_status_enum AS ENUM ('pass', 'fail', 'block', 'skip');`,
		`CREATE TYPE IF NOT EXISTS user_role_enum AS ENUM ('super_admin', 'admin', 'normal');`,
		`CREATE TYPE IF NOT EXISTS document_type_enum AS ENUM ('prd', 'figma', 'api_spec');`,
		`CREATE TYPE IF NOT EXISTS document_status_enum AS ENUM ('pending', 'processing', 'completed', 'failed');`,
	}

	for _, sql := range enumSQL {
		if _, err := db.ExecContext(ctx, sql); err != nil {
			return fmt.Errorf("create enum types: %w", err)
		}
	}

	// 创建表
	tableSQL := []string{
		// users 表
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(32) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role user_role_enum NOT NULL DEFAULT 'normal',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);`,

		// project 表
		`CREATE TABLE IF NOT EXISTS project (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL UNIQUE,
			prefix VARCHAR(4) NOT NULL UNIQUE,
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);`,

		// module 表
		`CREATE TABLE IF NOT EXISTS module (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			project_id UUID NOT NULL REFERENCES project(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			abbreviation VARCHAR(4) NOT NULL,
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(project_id, name),
			UNIQUE(project_id, abbreviation)
		);`,

		// project_config 表
		`CREATE TABLE IF NOT EXISTS project_config (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			project_id UUID NOT NULL REFERENCES project(id) ON DELETE CASCADE,
			key VARCHAR(255) NOT NULL,
			value JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(project_id, key)
		);`,

		// test_case 表
		`CREATE TABLE IF NOT EXISTS test_case (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			module_id UUID NOT NULL REFERENCES module(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id),
			number VARCHAR(32) NOT NULL UNIQUE,
			title VARCHAR(255) NOT NULL,
			preconditions JSONB DEFAULT '[]'::jsonb,
			steps JSONB DEFAULT '[]'::jsonb NOT NULL,
			expected JSONB DEFAULT '{}'::jsonb NOT NULL,
			ai_metadata JSONB DEFAULT '{}'::jsonb,
			case_type case_type_enum NOT NULL DEFAULT 'functionality',
			priority priority_enum NOT NULL DEFAULT 'P2',
			status case_status_enum NOT NULL DEFAULT 'unexecuted',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX IF NOT EXISTS idx_test_case_ai ON test_case USING gin (ai_metadata);
		CREATE INDEX IF NOT EXISTS idx_test_case_steps ON test_case USING gin (steps);`,

		// test_plan 表
		`CREATE TABLE IF NOT EXISTS test_plan (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			project_id UUID NOT NULL REFERENCES project(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id),
			name VARCHAR(255) NOT NULL,
			description TEXT DEFAULT '',
			status plan_status_enum NOT NULL DEFAULT 'draft',
			extra_config JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);`,

		// plan_cases 关联表
		`CREATE TABLE IF NOT EXISTS plan_cases (
			plan_id UUID NOT NULL REFERENCES test_plan(id) ON DELETE CASCADE,
			case_id UUID NOT NULL REFERENCES test_case(id) ON DELETE CASCADE,
			PRIMARY KEY (plan_id, case_id)
		);`,

		// test_result 表
		`CREATE TABLE IF NOT EXISTS test_result (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			case_id UUID NOT NULL REFERENCES test_case(id) ON DELETE CASCADE,
			plan_id UUID NOT NULL REFERENCES test_plan(id) ON DELETE CASCADE,
			executor_id UUID NOT NULL REFERENCES users(id),
			status result_status_enum NOT NULL,
			note TEXT DEFAULT '',
			executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// documents 表
		`CREATE TABLE IF NOT EXISTS documents (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			project_id UUID NOT NULL REFERENCES project(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			type document_type_enum NOT NULL,
			url TEXT,
			content_text TEXT,
			metadata JSONB DEFAULT '{}'::jsonb,
			status document_status_enum NOT NULL DEFAULT 'pending',
			created_by UUID NOT NULL REFERENCES users(id),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);`,

		// document_chunks 表
		`CREATE TABLE IF NOT EXISTS document_chunks (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
			chunk_index INTEGER NOT NULL,
			content TEXT NOT NULL,
			embedding BYTEA,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(document_id, chunk_index)
		);`,

		// generation_tasks 表
		`CREATE TABLE IF NOT EXISTS generation_tasks (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			project_id UUID NOT NULL REFERENCES project(id) ON DELETE CASCADE,
			module_id UUID REFERENCES module(id) ON DELETE SET NULL,
			user_id UUID NOT NULL REFERENCES users(id),
			status VARCHAR(32) NOT NULL DEFAULT 'pending',
			prompt TEXT,
			result_summary JSONB DEFAULT '{}'::jsonb,
			error_message TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// case_drafts 表
		`CREATE TABLE IF NOT EXISTS case_drafts (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			task_id UUID NOT NULL REFERENCES generation_tasks(id) ON DELETE CASCADE,
			module_id UUID REFERENCES module(id) ON DELETE SET NULL,
			title VARCHAR(255) NOT NULL,
			preconditions JSONB DEFAULT '[]'::jsonb,
			steps JSONB DEFAULT '[]'::jsonb NOT NULL,
			expected_result JSONB DEFAULT '{}'::jsonb NOT NULL,
			case_type case_type_enum NOT NULL DEFAULT 'functionality',
			priority priority_enum NOT NULL DEFAULT 'P2',
			ai_metadata JSONB DEFAULT '{}'::jsonb,
			status VARCHAR(32) NOT NULL DEFAULT 'pending',
			feedback TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, sql := range tableSQL {
		if _, err := db.ExecContext(ctx, sql); err != nil {
			return fmt.Errorf("create tables: %w", err)
		}
	}

	return nil
}

// TruncateAllTables 清空所有表数据（测试间清理）
func TruncateAllTables(ctx context.Context, db *sqlx.DB) error {
	tables := []string{
		"test_result", "plan_cases", "test_plan", "test_case", "module",
		"project_config", "project", "users", "document_chunks", "documents",
		"case_drafts", "generation_tasks",
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		if _, err := tx.ExecContext(ctx, query); err != nil {
			// 表可能不存在，继续清理其他表
			continue
		}
	}

	return tx.Commit()
}

// TestContext 封装测试上下文
type TestContext struct {
	Container *PostgresContainer
	DB        *sqlx.DB
	Cancel    context.CancelFunc
	T         *testing.T
}

// SetupTest 初始化测试环境
func SetupTest(t *testing.T) *TestContext {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	container, err := SetupPostgres(ctx)
	if err != nil {
		cancel()
		t.Fatalf("setup postgres container: %v", err)
	}

	db, err := container.CreateDB(ctx)
	if err != nil {
		_ = container.Teardown(ctx)
		cancel()
		t.Fatalf("create database connection: %v", err)
	}

	// 执行迁移
	if err := RunMigrations(ctx, db); err != nil {
		_ = container.Teardown(ctx)
		cancel()
		t.Fatalf("run migrations: %v", err)
	}

	tc := &TestContext{
		Container: container,
		DB:        db,
		Cancel:    cancel,
		T:         t,
	}

	// 注册清理函数
	t.Cleanup(func() {
		ctx := context.Background()
		if err := container.Teardown(ctx); err != nil {
			t.Logf("warning: failed to teardown container: %v", err)
		}
		cancel()
	})

	return tc
}

// CleanupTest 清理测试数据
func (tc *TestContext) CleanupTest() {
	tc.T.Helper()

	ctx := context.Background()
	if err := TruncateAllTables(ctx, tc.DB); err != nil {
		tc.T.Logf("warning: failed to cleanup tables: %v", err)
	}
}
