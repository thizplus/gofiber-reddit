-- Migration: Create Chat System Tables
-- Purpose: Add support for 1-on-1 messaging with conversations, messages, and blocks
-- Date: 2025-01-07
-- Phase: Chat System Phase 1 - Foundation

-- =============================================================================
-- Table: conversations
-- Purpose: Store 1-on-1 conversations between two users
-- =============================================================================

CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Participants (ordered by UUID for consistency)
    user1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Last Message (denormalized for performance)
    last_message_id UUID,
    last_message_at TIMESTAMP WITH TIME ZONE NOT NULL,

    -- Unread Counts (denormalized for performance)
    user1_unread_count INTEGER DEFAULT 0 NOT NULL,
    user2_unread_count INTEGER DEFAULT 0 NOT NULL,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Constraints
    CONSTRAINT conversations_user1_user2_unique UNIQUE (user1_id, user2_id),
    CONSTRAINT conversations_users_different CHECK (user1_id != user2_id)
);

-- Indexes for conversations
CREATE INDEX IF NOT EXISTS idx_conversations_user1_id ON conversations(user1_id);
CREATE INDEX IF NOT EXISTS idx_conversations_user2_id ON conversations(user2_id);
CREATE INDEX IF NOT EXISTS idx_conversations_user1_user2 ON conversations(user1_id, user2_id);
CREATE INDEX IF NOT EXISTS idx_conversations_last_message_at ON conversations(last_message_at DESC);
CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at);

-- =============================================================================
-- Table: messages
-- Purpose: Store individual chat messages
-- =============================================================================

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Conversation reference
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,

    -- Sender & Receiver
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Message Type & Content
    type VARCHAR(20) NOT NULL DEFAULT 'text',
    content TEXT,
    media JSONB,

    -- Read Status
    is_read BOOLEAN DEFAULT FALSE NOT NULL,
    read_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps (for cursor pagination)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Constraints
    CONSTRAINT messages_sender_receiver_different CHECK (sender_id != receiver_id),
    CONSTRAINT messages_content_or_media_required CHECK (content IS NOT NULL OR media IS NOT NULL)
);

-- Indexes for messages (optimized for cursor pagination)
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_created_at ON messages(conversation_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_is_read ON messages(is_read);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_type ON messages(type);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_type ON messages(conversation_id, type);
-- GIN index for JSONB media queries (for Media Gallery feature)
CREATE INDEX IF NOT EXISTS idx_messages_media ON messages USING GIN(media);

-- =============================================================================
-- Table: blocks
-- Purpose: Store user blocking relationships
-- =============================================================================

CREATE TABLE IF NOT EXISTS blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Blocker (user who blocks)
    blocker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Blocked (user being blocked)
    blocked_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Constraints
    CONSTRAINT blocks_blocker_blocked_unique UNIQUE (blocker_id, blocked_id),
    CONSTRAINT blocks_users_different CHECK (blocker_id != blocked_id)
);

-- Indexes for blocks
CREATE INDEX IF NOT EXISTS idx_blocks_blocker_id ON blocks(blocker_id);
CREATE INDEX IF NOT EXISTS idx_blocks_blocked_id ON blocks(blocked_id);
CREATE INDEX IF NOT EXISTS idx_blocks_blocker_blocked ON blocks(blocker_id, blocked_id);
CREATE INDEX IF NOT EXISTS idx_blocks_created_at ON blocks(created_at);

-- =============================================================================
-- Foreign Key: Add last_message_id FK to conversations after messages table exists
-- =============================================================================

ALTER TABLE conversations
ADD CONSTRAINT fk_conversations_last_message
FOREIGN KEY (last_message_id) REFERENCES messages(id) ON DELETE SET NULL;

-- =============================================================================
-- Triggers: Auto-update conversation on message insert
-- =============================================================================

-- Trigger function to update conversation when new message is inserted
CREATE OR REPLACE FUNCTION update_conversation_on_message()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE conversations
    SET
        last_message_id = NEW.id,
        last_message_at = NEW.created_at,
        updated_at = NEW.created_at,
        -- Increment unread count for receiver
        user1_unread_count = CASE
            WHEN user1_id = NEW.receiver_id THEN user1_unread_count + 1
            ELSE user1_unread_count
        END,
        user2_unread_count = CASE
            WHEN user2_id = NEW.receiver_id THEN user2_unread_count + 1
            ELSE user2_unread_count
        END
    WHERE id = NEW.conversation_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
DROP TRIGGER IF EXISTS trigger_update_conversation_on_message ON messages;
CREATE TRIGGER trigger_update_conversation_on_message
    AFTER INSERT ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_conversation_on_message();

-- =============================================================================
-- Verification Queries
-- =============================================================================

-- Verify tables were created
-- SELECT table_name FROM information_schema.tables
-- WHERE table_schema = 'public'
-- AND table_name IN ('conversations', 'messages', 'blocks')
-- ORDER BY table_name;

-- Verify indexes
-- SELECT tablename, indexname, indexdef
-- FROM pg_indexes
-- WHERE tablename IN ('conversations', 'messages', 'blocks')
-- ORDER BY tablename, indexname;

-- Verify constraints
-- SELECT conname, contype, pg_get_constraintdef(oid)
-- FROM pg_constraint
-- WHERE conrelid IN ('conversations'::regclass, 'messages'::regclass, 'blocks'::regclass)
-- ORDER BY conrelid::regclass::text, conname;

-- =============================================================================
-- Rollback (if needed)
-- =============================================================================

-- DROP TRIGGER IF EXISTS trigger_update_conversation_on_message ON messages;
-- DROP FUNCTION IF EXISTS update_conversation_on_message();
-- DROP TABLE IF EXISTS messages CASCADE;
-- DROP TABLE IF EXISTS conversations CASCADE;
-- DROP TABLE IF EXISTS blocks CASCADE;
