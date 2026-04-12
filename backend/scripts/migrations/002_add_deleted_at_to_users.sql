-- Migration: Add deleted_at column to users table for soft-delete support
-- Description: Add deleted_at column and index to support soft-delete pattern
-- Version: 002
-- Date: 2026-04-12

-- Add deleted_at column (nullable for backward compatibility)
ALTER TABLE users
ADD COLUMN deleted_at TIMESTAMP(3) WITH TIME ZONE;

-- Add index for performance on soft-delete queries
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NOT NULL;

-- Add comment for documentation
COMMENT ON COLUMN users.deleted_at IS 'Soft-delete timestamp. NULL indicates active record, non-NULL indicates deleted record.';
