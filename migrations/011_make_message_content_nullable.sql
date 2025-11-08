-- Migration: Make message content nullable for media-only messages
-- Purpose: Allow sending images/videos/files without text caption
-- Date: 2025-11-07

-- Remove NOT NULL constraint from content column (if exists)
ALTER TABLE messages
ALTER COLUMN content DROP NOT NULL;

-- Verify the constraint was removed
-- You can check with:
-- SELECT column_name, is_nullable, data_type
-- FROM information_schema.columns
-- WHERE table_name = 'messages' AND column_name = 'content';

-- Note: The CHECK constraint already ensures either content OR media must be provided:
-- CONSTRAINT messages_content_or_media_required CHECK (content IS NOT NULL OR media IS NOT NULL)
