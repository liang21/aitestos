-- ============================================================
-- Migration 004 Rollback: Restore prefix/deleted_at, old table names
-- ============================================================
-- Description:
--   - Restore table names: projects → project, modules → module
--   - Restore deleted_at column to users table
--   - Restore prefix column to project table
--   - Restore abbreviation column to module table
--   - Restore number column to test_case table
--   - Restore related indexes and constraints
-- WARNING: Only use if migration 004 was applied and no data was
--          added that depends on the new schema. Data loss may occur.
-- ============================================================

BEGIN;

-- 1. Restore table names
ALTER TABLE projects RENAME TO project;
ALTER TABLE modules RENAME TO module;

-- 2. Restore deleted_at to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at timestamp(3) with time zone;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- 3. Restore prefix to project table
ALTER TABLE project ADD COLUMN IF NOT EXISTS prefix varchar(4);
ALTER TABLE project ADD CONSTRAINT project_prefix_key UNIQUE (prefix);

-- 4. Restore abbreviation to module table
ALTER TABLE module ADD COLUMN IF NOT EXISTS abbreviation varchar(4);
ALTER TABLE module ADD CONSTRAINT module_abbreviation_key UNIQUE (abbreviation);

-- 5. Restore number to test_case table
ALTER TABLE test_case ADD COLUMN IF NOT EXISTS number varchar(50);
ALTER TABLE test_case ADD CONSTRAINT test_case_number_key UNIQUE (number);

COMMIT;
