-- Migration: Add file extension field to media table
-- Purpose: Support file uploads with extension tracking for filtering
-- Date: 2025-01-07
-- Phase: File Upload Feature - Foundation

-- =============================================================================
-- Add extension field to media table
-- =============================================================================

ALTER TABLE media
ADD COLUMN IF NOT EXISTS extension VARCHAR(10);

-- =============================================================================
-- Create index for extension lookups (for filtering by file type)
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_media_extension ON media(extension);

-- =============================================================================
-- Comments
-- =============================================================================

COMMENT ON COLUMN media.extension IS 'File extension without dot (e.g., pdf, doc, zip)';

-- =============================================================================
-- Verification Queries
-- =============================================================================

-- Verify column was added
-- SELECT column_name, data_type, character_maximum_length
-- FROM information_schema.columns
-- WHERE table_name = 'media' AND column_name = 'extension';

-- Verify index was created
-- SELECT indexname, indexdef
-- FROM pg_indexes
-- WHERE tablename = 'media' AND indexname = 'idx_media_extension';

-- =============================================================================
-- Rollback (if needed)
-- =============================================================================

-- DROP INDEX IF EXISTS idx_media_extension;
-- ALTER TABLE media DROP COLUMN IF EXISTS extension;
