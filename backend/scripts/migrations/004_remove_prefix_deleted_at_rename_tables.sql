-- ============================================================
-- Migration 004: Remove prefix/deleted_at, rename tables
-- ============================================================
-- Description:
--   - Remove deleted_at column from users table
--   - Remove prefix column from project table
--   - Remove abbreviation column from module table
--   - Remove number column from test_case table
--   - Rename project → projects
--   - Rename module → modules
--   - Drop related indexes and constraints
-- ============================================================

BEGIN;

-- 1. Remove deleted_at from users table
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;

-- 2. Remove prefix from project table
ALTER TABLE project DROP COLUMN IF EXISTS prefix;

-- 3. Remove abbreviation from module table
-- First drop the unique constraint on abbreviation
ALTER TABLE module DROP CONSTRAINT IF EXISTS module_abbreviation_key;
-- Then drop the column
ALTER TABLE module DROP COLUMN IF EXISTS abbreviation;

-- 4. Remove number from test_case table
-- First drop the unique constraint on number
ALTER TABLE test_case DROP CONSTRAINT IF EXISTS test_case_number_key;
-- Then drop the column
ALTER TABLE test_case DROP COLUMN IF EXISTS number;

-- 5. Rename tables
ALTER TABLE project RENAME TO projects;
ALTER TABLE module RENAME TO modules;

-- 6. Drop obsolete indexes
DROP INDEX IF EXISTS idx_users_deleted_at;

COMMIT;
