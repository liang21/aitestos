-- Rollback: Remove deleted_at column from users table
-- Description: Remove soft-delete support
-- Version: 002

-- Drop index
DROP INDEX IF EXISTS idx_users_deleted_at;

-- Drop column (CAUTION: This will lose all soft-delete information)
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
