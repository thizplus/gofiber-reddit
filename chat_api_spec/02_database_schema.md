# Chat API Specification - Database Schema

## PostgreSQL Schema Design

### 1. Conversations Table

```sql
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Last message info (denormalized for performance)
    last_message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    last_message_content TEXT,
    last_message_sender_id UUID,
    last_message_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT unique_conversation UNIQUE(user1_id, user2_id),
    CONSTRAINT different_users CHECK (user1_id != user2_id),
    CONSTRAINT ordered_users CHECK (user1_id < user2_id)
);

-- Indexes
CREATE INDEX idx_conversations_user1 ON conversations(user1_id, updated_at DESC);
CREATE INDEX idx_conversations_user2 ON conversations(user2_id, updated_at DESC);
CREATE INDEX idx_conversations_updated_at ON conversations(updated_at DESC);

-- Comments
COMMENT ON TABLE conversations IS 'Stores 1-on-1 conversations between two users';
COMMENT ON COLUMN conversations.user1_id IS 'First user ID (always smaller UUID)';
COMMENT ON COLUMN conversations.user2_id IS 'Second user ID (always larger UUID)';
COMMENT ON CONSTRAINT ordered_users ON conversations IS 'Ensures user1_id < user2_id to avoid duplicate conversations';
```

**Design Notes**:
- `user1_id` และ `user2_id` เรียงลำดับตาม UUID (user1 < user2) เพื่อป้องกัน duplicate conversations
- Denormalize `last_message_*` fields เพื่อเพิ่ม performance การดึงรายการสนทนา
- Index ทั้ง user1 และ user2 เพื่อ query ได้รวดเร็ว

---

### 2. Messages Table

```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Message content
    content TEXT NOT NULL,

    -- Read status
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Soft delete (for future)
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender ON messages(sender_id, created_at DESC);
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);
CREATE INDEX idx_messages_unread ON messages(conversation_id, is_read) WHERE is_read = FALSE;

-- Partial index for active messages
CREATE INDEX idx_messages_active ON messages(conversation_id, created_at DESC)
    WHERE deleted_at IS NULL;

-- Comments
COMMENT ON TABLE messages IS 'Stores all chat messages';
COMMENT ON COLUMN messages.content IS 'Message text content (Phase 1: text only)';
COMMENT ON COLUMN messages.is_read IS 'Whether the message has been read by the receiver';
```

**Design Notes**:
- Index `conversation_id + created_at DESC` สำหรับ pagination ย้อนหลัง
- Partial index สำหรับ unread messages เพิ่ม performance
- `deleted_at` เตรียมไว้สำหรับ soft delete (Phase 2)

---

### 3. Conversation Participants (Alternative Design)

```sql
-- Alternative: ถ้าต้องการรองรับ group chat ในอนาคต
CREATE TABLE conversation_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Participant-specific data
    unread_count INTEGER DEFAULT 0,
    last_read_message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    last_read_at TIMESTAMP WITH TIME ZONE,

    -- Settings
    is_muted BOOLEAN DEFAULT FALSE,
    is_archived BOOLEAN DEFAULT FALSE,

    -- Timestamps
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    left_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_participant UNIQUE(conversation_id, user_id)
);

-- Indexes
CREATE INDEX idx_participants_user ON conversation_participants(user_id, conversation_id);
CREATE INDEX idx_participants_conversation ON conversation_participants(conversation_id);
CREATE INDEX idx_participants_unread ON conversation_participants(user_id, unread_count)
    WHERE unread_count > 0;
```

**Note**: Phase 1 ใช้ `conversations` table แบบง่าย, แต่ถ้าต้องการ track unread per user หรือเตรียมพร้อม group chat ควรใช้ `conversation_participants`

---

### 4. Blocks Table

```sql
CREATE TABLE blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blocker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT unique_block UNIQUE(blocker_id, blocked_id),
    CONSTRAINT different_users_block CHECK (blocker_id != blocked_id)
);

-- Indexes
CREATE INDEX idx_blocks_blocker ON blocks(blocker_id, blocked_id);
CREATE INDEX idx_blocks_blocked ON blocks(blocked_id);

-- Comments
COMMENT ON TABLE blocks IS 'Stores user blocking relationships';
COMMENT ON COLUMN blocks.blocker_id IS 'User who initiated the block';
COMMENT ON COLUMN blocks.blocked_id IS 'User who got blocked';
```

**Design Notes**:
- เก็บเฉพาะ active blocks (ไม่ต้อง soft delete)
- Query ได้ทั้งทิศทาง: "ใครที่ฉันบล็อก" และ "ใครบล็อกฉัน"

---

### 5. Redis Schema (Cache)

#### Online Status
```
Key: "online:{user_id}"
Type: String
Value: timestamp (last seen)
TTL: 60 seconds

Example:
online:123e4567-e89b-12d3-a456-426614174000 = "1704067200"
```

#### Unread Count (Per User)
```
Key: "unread:{user_id}"
Type: Integer
Value: total unread count
TTL: None (persistent, updated on new message)

Example:
unread:123e4567-e89b-12d3-a456-426614174000 = "5"
```

#### Unread Count (Per Conversation)
```
Key: "unread:{user_id}:{conversation_id}"
Type: Integer
Value: unread count for specific conversation
TTL: None (persistent)

Example:
unread:123e4567-e89b-12d3-a456-426614174000:conv-001 = "2"
```

#### Last Message Cache
```
Key: "last_msg:{conversation_id}"
Type: Hash
Fields:
  - id
  - sender_id
  - content
  - created_at
TTL: 1 hour

Example:
HGETALL last_msg:conv-001
{
  "id": "msg-001",
  "sender_id": "user-123",
  "content": "Hello!",
  "created_at": "2024-01-01T10:00:00Z"
}
```

#### WebSocket Connection Tracking
```
Key: "ws:{user_id}"
Type: Set
Value: set of connection IDs
TTL: None (removed on disconnect)

Example:
ws:123e4567-e89b-12d3-a456-426614174000 = {"conn-1", "conn-2"}
```

---

### Database Triggers

#### 1. Update conversation.updated_at on new message
```sql
CREATE OR REPLACE FUNCTION update_conversation_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE conversations
    SET updated_at = NEW.created_at,
        last_message_id = NEW.id,
        last_message_content = NEW.content,
        last_message_sender_id = NEW.sender_id,
        last_message_at = NEW.created_at
    WHERE id = NEW.conversation_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_conversation
AFTER INSERT ON messages
FOR EACH ROW
EXECUTE FUNCTION update_conversation_timestamp();
```

#### 2. Auto-create conversation if not exists (Optional)
```sql
-- ถ้าต้องการให้สร้าง conversation อัตโนมัติเมื่อส่งข้อความครั้งแรก
-- ควรทำใน application layer แทน
```

---

### Migrations Plan

#### Migration 1: Initial Schema
```sql
-- Create conversations table
-- Create messages table
-- Create blocks table
-- Create indexes
-- Create triggers
```

#### Migration 2: Add conversation_participants (if needed)
```sql
-- Create conversation_participants table
-- Migrate data from conversations
-- Add indexes
```

#### Migration 3: Add soft delete support
```sql
-- Add deleted_at to messages
-- Add deleted_by to messages
-- Create partial indexes
```

---

### Data Retention Policy

#### Active Data (Hot Storage)
- Messages: เก็บทั้งหมดใน main table
- Conversations: เก็บทั้งหมด

#### Archive Strategy (Future)
- Archive messages > 1 ปี ไป `messages_archive` table
- Keep index ใน main table สำหรับ search
- ใช้ partitioning by created_at

#### Deletion Policy
- User deleted: CASCADE ลบ conversations และ messages
- Conversation deleted: ลบ messages ที่เกี่ยวข้อง
- Block deleted: เก็บ conversations ไว้

---

### Estimated Storage

#### Assumptions
- Active users: 100,000
- Messages per user per day: 20
- Average message size: 100 bytes
- Retention: 2 years

#### Calculations
```
Messages per day: 100,000 × 20 = 2,000,000
Messages per year: 2,000,000 × 365 = 730,000,000
Total messages (2 years): 1,460,000,000

Storage per message: ~300 bytes (including indexes)
Total storage: 1.46B × 300 bytes = 438 GB

Conversations: 100,000 users × 50 conversations = 5,000,000
Conversation storage: 5M × 500 bytes = 2.5 GB

Total: ~440 GB (within reasonable limits for PostgreSQL)
```

---

### Performance Targets

| Operation | Target | Notes |
|-----------|--------|-------|
| Get conversations list | < 100ms | With pagination |
| Get messages | < 100ms | With pagination |
| Send message | < 50ms | Write to DB + cache |
| Mark as read | < 50ms | Update + cache invalidation |
| Check block status | < 10ms | From cache/index |

---

### Backup Strategy

1. **PostgreSQL**
   - Full backup: Daily
   - Incremental backup: Every 6 hours
   - WAL archiving: Continuous
   - Retention: 30 days

2. **Redis**
   - RDB snapshots: Every hour
   - AOF: Enabled
   - Retention: 7 days
   - Note: Redis data is cache, can be rebuilt
