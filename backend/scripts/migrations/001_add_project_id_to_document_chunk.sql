-- Migration: Add project_id to document_chunk table
-- Description: Add project_id column to support vector filtering by project
-- Version: 001

-- Add project_id column
ALTER TABLE document_chunk
ADD COLUMN project_id uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES project(id) ON DELETE CASCADE;

-- Create index for filtering
CREATE INDEX idx_document_chunk_project_id ON document_chunk(project_id);

-- Update existing records with their associated project_id
UPDATE document_chunk dc
SET project_id = (
    SELECT d.project_id
    FROM document d
    WHERE d.id = dc.document_id
)
WHERE dc.project_id = '00000000-0000-0000-0000-000000000000';

-- Add comment
COMMENT ON COLUMN document_chunk.project_id IS 'Associated project ID for vector filtering';
