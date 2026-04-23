-- Migration: Add prefix and deleted_at columns to project table
-- This migration aligns the project table schema with the domain model expectations

-- Add prefix column (nullable first, then set unique constraint after populating existing rows)
ALTER TABLE project ADD COLUMN IF NOT EXISTS prefix varchar(4);

-- For existing projects without a prefix, generate one from the name (first 4 chars, uppercase)
-- This is a safe default that can be manually updated later
UPDATE project SET prefix = UPPER(SUBSTRING(name, 1, 4)) WHERE prefix IS NULL;

-- Now make it NOT NULL and UNIQUE
ALTER TABLE project ALTER COLUMN prefix SET NOT NULL;
ALTER TABLE project ADD CONSTRAINT project_prefix_key UNIQUE (prefix);

-- Add deleted_at column for soft delete support
ALTER TABLE project ADD COLUMN IF NOT EXISTS deleted_at timestamp(3) with time zone;

-- Remove unused config column if it exists
ALTER TABLE project DROP COLUMN IF EXISTS config;

-- Create index on deleted_at for better query performance on soft deletes
CREATE INDEX IF NOT EXISTS idx_project_deleted_at ON project(deleted_at) WHERE deleted_at IS NOT NULL;

-- Verify the schema
SELECT
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_name = 'project'
ORDER BY ordinal_position;
