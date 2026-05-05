-- ============================================================
-- Migration 005: Rename document tables to plural form
-- ============================================================
-- Description:
--   - Rename document → documents to match plural naming convention
--   - Rename document_chunk → document_chunks to match plural naming convention
--   - This aligns with the migration 004 which renamed project→projects
--     and module→modules
-- ============================================================

BEGIN;

-- Rename document_chunk table to document_chunks (must be done first due to foreign key)
ALTER TABLE document_chunk RENAME TO document_chunks;

-- Rename document table to documents
ALTER TABLE document RENAME TO documents;

COMMIT;
