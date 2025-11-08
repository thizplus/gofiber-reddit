-- Migration: Add video streaming fields to media table
-- Purpose: Support Bunny Stream HLS video streaming with encoding tracking
-- Date: 2025-01-07
-- Phase: Video Streaming Implementation

-- =============================================================================
-- Add video streaming fields to media table
-- =============================================================================

-- Bunny Stream Video ID
ALTER TABLE media
ADD COLUMN IF NOT EXISTS video_id VARCHAR(255);

-- HLS Playlist URL (for streaming)
ALTER TABLE media
ADD COLUMN IF NOT EXISTS hls_url TEXT;

-- Encoding status: pending, processing, completed, failed
ALTER TABLE media
ADD COLUMN IF NOT EXISTS encoding_status VARCHAR(20) DEFAULT 'pending';

-- Encoding progress: 0-100 (percentage)
ALTER TABLE media
ADD COLUMN IF NOT EXISTS encoding_progress INTEGER DEFAULT 0;

-- =============================================================================
-- Create indexes for video streaming queries
-- =============================================================================

-- Index for video_id lookups (Bunny Stream API queries)
CREATE INDEX IF NOT EXISTS idx_media_video_id ON media(video_id);

-- Index for encoding status queries (for worker polling)
CREATE INDEX IF NOT EXISTS idx_media_encoding_status ON media(encoding_status);

-- Composite index for type + encoding_status (efficient filtering)
CREATE INDEX IF NOT EXISTS idx_media_type_encoding_status ON media(type, encoding_status);

-- =============================================================================
-- Comments
-- =============================================================================

COMMENT ON COLUMN media.video_id IS 'Bunny Stream video ID (UUID from Bunny API)';
COMMENT ON COLUMN media.hls_url IS 'HLS playlist URL for adaptive streaming (m3u8)';
COMMENT ON COLUMN media.encoding_status IS 'Video encoding status: pending, processing, completed, failed';
COMMENT ON COLUMN media.encoding_progress IS 'Encoding progress percentage (0-100)';

-- =============================================================================
-- Verification Queries
-- =============================================================================

-- Verify columns were added
-- SELECT column_name, data_type, character_maximum_length, column_default
-- FROM information_schema.columns
-- WHERE table_name = 'media'
-- AND column_name IN ('video_id', 'hls_url', 'encoding_status', 'encoding_progress');

-- Verify indexes were created
-- SELECT indexname, indexdef
-- FROM pg_indexes
-- WHERE tablename = 'media'
-- AND indexname IN ('idx_media_video_id', 'idx_media_encoding_status', 'idx_media_type_encoding_status');

-- =============================================================================
-- Rollback (if needed)
-- =============================================================================

-- DROP INDEX IF EXISTS idx_media_type_encoding_status;
-- DROP INDEX IF EXISTS idx_media_encoding_status;
-- DROP INDEX IF EXISTS idx_media_video_id;
-- ALTER TABLE media DROP COLUMN IF EXISTS encoding_progress;
-- ALTER TABLE media DROP COLUMN IF EXISTS encoding_status;
-- ALTER TABLE media DROP COLUMN IF EXISTS hls_url;
-- ALTER TABLE media DROP COLUMN IF EXISTS video_id;
