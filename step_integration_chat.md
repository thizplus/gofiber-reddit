# Chat System Integration Plan

> แผนการพัฒนาและ integrate Chat System เข้ากับระบบ Social Platform ที่มีอยู่

**สถานะ**: Ready for Implementation
**เวลาประมาณการ**: 8-10 สัปดาห์
**วันที่สร้าง**: 2025-01-07

---

## 📋 Table of Contents

1. [สรุปความต้องการจาก Frontend](#1-สรุปความต้องการจาก-frontend)
2. [วิเคราะห์ระบบปัจจุบัน](#2-วิเคราะห์ระบบปัจจุบัน)
3. [ผลกระทบกับ Code เดิม](#3-ผลกระทบกับ-code-เดิม)
4. [แผนการ Implementation (Step-by-Step)](#4-แผนการ-implementation-step-by-step)
5. [Timeline & Resources](#5-timeline--resources)
6. [Testing Strategy](#6-testing-strategy)
7. [Deployment Plan](#7-deployment-plan)
8. [Risks & Mitigation](#8-risks--mitigation)

---

## 1. สรุปความต้องการจาก Frontend

### 1.1 REST API Endpoints (14 endpoints)

#### **Conversations (3 endpoints)**
```
GET    /api/v1/chat/conversations                    - รายการสนทนา (with pagination)
GET    /api/v1/chat/conversations/with/:username     - Get/Create conversation
GET    /api/v1/chat/conversations/unread-count       - จำนวนข้อความยังไม่อ่าน
```

#### **Messages (8 endpoints)**
```
# Core Messages (Phase 1)
GET    /api/v1/chat/conversations/:id/messages       - ดึงข้อความ (with pagination)
POST   /api/v1/chat/conversations/:id/messages       - ส่งข้อความ
POST   /api/v1/chat/conversations/:id/read           - Mark as read
GET    /api/v1/chat/messages/:id                     - ดึงข้อความเดียว

# Jump to Message (Phase 1)
GET    /api/v1/chat/messages/:id/context             - 🆕 Jump to message พร้อม context

# Telegram-style Features (Phase 2 - Optional)
GET    /api/v1/chat/conversations/:id/media          - 🆕 รายการ media ทั้งหมด
GET    /api/v1/chat/conversations/:id/links          - 🆕 รายการ links ทั้งหมด
GET    /api/v1/chat/conversations/:id/files          - 🆕 รายการ files ทั้งหมด
```

#### **Blocks (3 endpoints)**
```
POST   /api/v1/chat/blocks                           - บล็อกผู้ใช้
DELETE /api/v1/chat/blocks/:username                 - ปลดบล็อก
GET    /api/v1/chat/blocks                           - รายการผู้ใช้ที่ถูกบล็อก
GET    /api/v1/chat/blocks/status/:username          - เช็คสถานะการบล็อก (optional)
```

### 1.2 WebSocket Events (8 events)

#### **Client → Server**
- `message.send` - ส่งข้อความ
- `message.read` - Mark as read
- `ping` - Heartbeat
- `auth` - Authentication (optional, ถ้าไม่ส่ง token ใน query string)

#### **Server → Client**
- `message.new` - ข้อความใหม่
- `message.sent` - ส่งสำเร็จ (confirmation)
- `user.online` / `user.offline` - Online status
- `conversation.updated` - Conversation update
- `notification.unread` - Unread count update

### 1.3 Database Schema (3 tables)

#### **conversations**
```sql
- id (UUID, PK)
- user1_id (UUID, FK to users, CHECK: user1_id < user2_id)
- user2_id (UUID, FK to users)
- last_message_id (UUID, FK to messages, nullable)
- last_message_content (TEXT, denormalized)
- last_message_sender_id (UUID)
- last_message_at (TIMESTAMP)
- created_at, updated_at
- UNIQUE(user1_id, user2_id)
```

#### **messages**
```sql
- id (UUID, PK)
- conversation_id (UUID, FK, indexed)
- sender_id (UUID, FK to users, indexed)
- content (TEXT, required)
- is_read (BOOLEAN, default: false)
- read_at (TIMESTAMP, nullable)
- created_at, updated_at
- deleted_at (TIMESTAMP, nullable, for soft delete)
```

#### **blocks**
```sql
- id (UUID, PK)
- blocker_id (UUID, FK to users)
- blocked_id (UUID, FK to users)
- created_at
- UNIQUE(blocker_id, blocked_id)
```

### 1.4 Redis Schema

```
# Online Status
online:{user_id} → timestamp (TTL: 60s)

# Unread Count
unread:{user_id} → total_count
unread:{user_id}:{conversation_id} → count

# Last Message Cache
last_msg:{conversation_id} → hash (id, sender_id, content, created_at)

# WebSocket Connections
ws:{user_id} → set of connection IDs
```

### 1.5 Pagination

**Cursor-based pagination** (ไม่ใช่ offset-based):
```json
{
  "cursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0wMVQxMDowMDowMFoiLCJpZCI6Im1zZy0wNTAifQ==",
  "limit": 20
}
```

**Cursor Structure**:
```json
{
  "created_at": "2024-01-01T10:00:00Z",
  "id": "msg-050"
}
```

### 1.6 Performance Targets

- API Response: < 100ms (p95)
- WebSocket Latency: < 50ms
- Concurrent Connections: 1000+
- Cache Hit Rate: > 80%

---

## 2. วิเคราะห์ระบบปัจจุบัน

### 2.1 ✅ สิ่งที่มีอยู่แล้วและใช้ได้

#### **1. Architecture Pattern**
- ✅ Clean Architecture (4 layers: Domain, Application, Infrastructure, Interface)
- ✅ Dependency Injection (custom DI container)
- ✅ Repository Pattern
- ✅ Service Layer Pattern
- **การใช้งาน**: ทำตาม pattern เดียวกันสำหรับ Chat System

#### **2. Tech Stack**
- ✅ **Go Fiber** v2.52.0 (แทน Gin ที่ spec บอก)
- ✅ **PostgreSQL** 15 + GORM
- ✅ **Redis** 7 (go-redis)
- ✅ **JWT Authentication** (มีระบบ auth ครบ)
- ✅ **WebSocket** (gofiber/websocket)
- **หมายเหตุ**: Spec ใช้ Gin แต่เราใช้ Fiber ซึ่ง syntax คล้ายกัน แค่ปรับเล็กน้อย

#### **3. WebSocket Infrastructure**
- ✅ มี WebSocket Manager แล้ว (`infrastructure/websocket/manager.go`)
- ✅ มี Client, Hub structure
- ✅ มี connection management
- ✅ มี room-based messaging
- **การใช้งาน**: ขยายและเพิ่ม chat-specific events

#### **4. Database**
- ✅ PostgreSQL connection + GORM
- ✅ Auto-migration system
- ✅ UUID support
- ✅ Index support
- ✅ Soft delete support (`gorm.DeletedAt`)
- **การใช้งาน**: เพิ่ม 3 tables ใหม่

#### **5. Redis**
- ✅ Redis client configured (`infrastructure/redis/client.go`)
- ✅ Connection pooling
- ✅ Error handling
- **การใช้งาน**: เพิ่ม functions สำหรับ online status และ cache

#### **6. Authentication**
- ✅ JWT middleware (`interfaces/api/middleware/auth.go`)
- ✅ Protected routes
- ✅ Optional auth
- ✅ User context extraction
- **การใช้งาน**: ใช้ middleware เดิมได้เลย

#### **7. Notification System**
- ✅ In-app notifications
- ✅ Web Push notifications
- ✅ Notification settings
- **การใช้งาน**: Integrate กับระบบแจ้งเตือนที่มี

### 2.2 ❌ สิ่งที่ยังไม่มีและต้องสร้างใหม่

#### **1. Chat-specific Models & DTOs**
- ❌ Conversation model
- ❌ Message model
- ❌ Block model
- ❌ Chat DTOs

#### **2. Chat Repositories**
- ❌ ConversationRepository
- ❌ MessageRepository
- ❌ BlockRepository

#### **3. Chat Services**
- ❌ ChatService (conversations + messages)
- ❌ BlockService
- ❌ OnlineStatusService

#### **4. Chat Handlers**
- ❌ ConversationHandler
- ❌ MessageHandler
- ❌ BlockHandler
- ❌ ChatWebSocketHandler

#### **5. Cursor Pagination**
- ❌ Cursor encoding/decoding utilities
- ❌ Cursor-based queries
- **หมายเหตุ**: ระบบเดิมใช้ offset-based pagination

#### **6. Redis Functions**
- ❌ Online status tracking
- ❌ Unread count tracking
- ❌ Last message caching

---

## 3. ผลกระทบกับ Code เดิม

### 3.1 🟢 ไม่มีผลกระทบ (Safe)

#### **Database Tables**
- ✅ **ไม่กระทบ tables เดิม** - เพิ่ม 3 tables ใหม่เท่านั้น
- ✅ ไม่มีการแก้ไข schema เดิม
- ✅ ไม่มี foreign key conflict

#### **Existing Features**
- ✅ **Posts, Comments, Votes** - ไม่กระทบเลย
- ✅ **User Management** - ใช้ร่วมกันได้ (FK to users table)
- ✅ **Authentication** - ใช้ระบบเดิม
- ✅ **Notifications** - ใช้ร่วมกันได้

#### **WebSocket**
- ✅ **Existing WebSocket** - แยก endpoint (/ws vs /chat/ws)
- ✅ สามารถ coexist ได้ (different routes)
- ✅ ใช้ Manager ตัวเดียวกัน แต่ต่าง namespace

### 3.2 🟡 ต้องปรับเล็กน้อย (Minor Changes)

#### **1. WebSocket Manager**
- **ปัจจุบัน**: มี generic message handling
- **ต้องเพิ่ม**: Chat-specific event handlers
- **วิธีการ**: Extend Manager ด้วย chat methods

```go
// ปัจจุบัน
type Manager struct {
    clients map[string]*Client
    // ...
}

// เพิ่ม
func (m *Manager) BroadcastChatMessage(userID string, msg *ChatMessage)
func (m *Manager) SendOnlineStatus(userIDs []string, status bool)
```

#### **2. Notification Integration**
- **ปัจจุบัน**: มี NotificationService
- **ต้องเพิ่ม**: Chat notification types
- **วิธีการ**: เพิ่ม type ใหม่ในระบบเดิม

```go
// เพิ่ม notification types
const (
    NotificationTypeReply   = "reply"   // มีอยู่แล้ว
    NotificationTypeVote    = "vote"    // มีอยู่แล้ว
    NotificationTypeMessage = "message" // ใหม่
)
```

#### **3. User Model**
- **อาจต้องเพิ่ม** (optional): Last seen online
- **วิธีการ**: เพิ่ม field ใหม่ไม่บังคับ

```go
type User struct {
    // ... existing fields
    LastSeenAt *time.Time `gorm:"index"` // ใหม่ (optional)
}
```

### 3.3 🔴 ต้องระวัง (Caution)

#### **1. Database Connections**
- **ปัจจุบัน**: Connection pool มีอยู่
- **ระวัง**: Chat อาจมี queries บ่อย → ต้อง monitor connection usage
- **แก้ไข**: อาจต้องเพิ่ม pool size

```go
sqlDB.SetMaxOpenConns(100)  // เดิม
sqlDB.SetMaxOpenConns(200)  // อาจต้องเพิ่ม
```

#### **2. WebSocket Connections**
- **ปัจจุบัน**: ไม่มี limit
- **ระวัง**: Chat users อาจ online นานกว่า → เพิ่ม memory usage
- **แก้ไข**: เพิ่ม connection limit per user

#### **3. Redis Memory**
- **ปัจจุบัน**: ใช้สำหรับ session + cache
- **ระวัง**: Online status + unread counts จะเพิ่ม memory usage
- **แก้ไข**: ตั้ง TTL และ monitor memory

#### **4. Notification Spam**
- **ระวัง**: Chat messages อาจสร้าง notification เยอะมาก
- **แก้ไข**: เพิ่ม rate limiting หรือ batch notifications

---

## 4. แผนการ Implementation (Step-by-Step)

### 📌 หมายเหตุสำคัญ
- ทำตาม **Clean Architecture** ของระบบเดิม
- แต่ละ step ต้อง **test ก่อนไป step ถัดไป**
- **ไม่แตะโค้ดเดิม** เว้นแต่จำเป็น

---

### **Phase 1: Foundation (Week 1-2)**

#### **Step 1.1: Database Models & Migrations**
**เวลา**: 2 days

**ทำอะไร**:
1. สร้าง domain models ใหม่

```go
// domain/models/conversation.go
type Conversation struct {
    ID                  uuid.UUID  `gorm:"type:uuid;primaryKey"`
    User1ID             uuid.UUID  `gorm:"type:uuid;not null;index"`
    User2ID             uuid.UUID  `gorm:"type:uuid;not null;index"`

    // Denormalized last message
    LastMessageID       *uuid.UUID `gorm:"type:uuid"`
    LastMessageContent  string
    LastMessageSenderID *uuid.UUID `gorm:"type:uuid"`
    LastMessageAt       *time.Time

    // Relationships
    User1               User      `gorm:"foreignKey:User1ID"`
    User2               User      `gorm:"foreignKey:User2ID"`
    Messages            []Message `gorm:"foreignKey:ConversationID"`

    CreatedAt           time.Time
    UpdatedAt           time.Time
}

// domain/models/message.go
type Message struct {
    ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
    ConversationID uuid.UUID `gorm:"type:uuid;not null;index"`
    SenderID       uuid.UUID `gorm:"type:uuid;not null;index"`
    Content        string    `gorm:"type:text;not null"`

    IsRead         bool       `gorm:"default:false;index"`
    ReadAt         *time.Time

    // Relationships
    Conversation   Conversation `gorm:"foreignKey:ConversationID"`
    Sender         User         `gorm:"foreignKey:SenderID"`

    CreatedAt      time.Time `gorm:"index"`
    UpdatedAt      time.Time
    DeletedAt      *time.Time `gorm:"index"` // Soft delete
}

// domain/models/block.go
type Block struct {
    ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
    BlockerID  uuid.UUID `gorm:"type:uuid;not null;index"`
    BlockedID  uuid.UUID `gorm:"type:uuid;not null;index"`

    Blocker    User `gorm:"foreignKey:BlockerID"`
    Blocked    User `gorm:"foreignKey:BlockedID"`

    CreatedAt  time.Time
}
```

2. สร้าง migration file

```sql
-- infrastructure/postgres/migrations/008_create_chat_tables.sql

-- Conversations table
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    last_message_id UUID,
    last_message_content TEXT,
    last_message_sender_id UUID,
    last_message_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT unique_conversation UNIQUE(user1_id, user2_id),
    CONSTRAINT different_users CHECK (user1_id != user2_id),
    CONSTRAINT ordered_users CHECK (user1_id < user2_id)
);

CREATE INDEX idx_conversations_user1 ON conversations(user1_id, updated_at DESC);
CREATE INDEX idx_conversations_user2 ON conversations(user2_id, updated_at DESC);
CREATE INDEX idx_conversations_updated_at ON conversations(updated_at DESC);

-- Messages table
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,

    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender ON messages(sender_id, created_at DESC);
CREATE INDEX idx_messages_unread ON messages(conversation_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_messages_active ON messages(conversation_id, created_at DESC) WHERE deleted_at IS NULL;

-- Blocks table
CREATE TABLE blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blocker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT unique_block UNIQUE(blocker_id, blocked_id),
    CONSTRAINT different_users_block CHECK (blocker_id != blocked_id)
);

CREATE INDEX idx_blocks_blocker ON blocks(blocker_id, blocked_id);
CREATE INDEX idx_blocks_blocked ON blocks(blocked_id);

-- Trigger: Update conversation timestamp on new message
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

3. เพิ่มใน auto-migration

```go
// infrastructure/postgres/database.go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        // ... existing models
        &models.Conversation{},
        &models.Message{},
        &models.Block{},
    )
}
```

**Test**:
- [ ] Run migration สำเร็จ
- [ ] Tables ถูกสร้างครบ
- [ ] Indexes ถูกต้อง
- [ ] Trigger ทำงานได้

---

#### **Step 1.2: DTOs (Data Transfer Objects)**
**เวลา**: 1 day

**ทำอะไร**:
```go
// domain/dto/chat_dto.go

// Conversation DTOs
type ConversationDTO struct {
    ID            string         `json:"id"`
    OtherUser     *UserShortDTO  `json:"otherUser"`
    LastMessage   *MessageShortDTO `json:"lastMessage"`
    UnreadCount   int            `json:"unreadCount"`
    UpdatedAt     string         `json:"updatedAt"`
    IsBlocked     bool           `json:"isBlocked"`
}

type UserShortDTO struct {
    ID          string  `json:"id"`
    Username    string  `json:"username"`
    DisplayName string  `json:"displayName"`
    Avatar      string  `json:"avatar"`
    IsOnline    bool    `json:"isOnline"`
    LastSeen    *string `json:"lastSeen,omitempty"`
}

type MessageShortDTO struct {
    ID        string `json:"id"`
    SenderID  string `json:"senderId"`
    Content   string `json:"content"`
    CreatedAt string `json:"createdAt"`
    IsRead    bool   `json:"isRead"`
}

// Message DTOs
type MessageDTO struct {
    ID             string  `json:"id"`
    ConversationID string  `json:"conversationId"`
    SenderID       string  `json:"senderId"`
    Sender         *UserShortDTO `json:"sender,omitempty"`
    Content        string  `json:"content"`
    IsRead         bool    `json:"isRead"`
    ReadAt         *string `json:"readAt,omitempty"`
    CreatedAt      string  `json:"createdAt"`
    UpdatedAt      string  `json:"updatedAt"`
}

type SendMessageDTO struct {
    Content string `json:"content" validate:"required,min=1,max=10000"`
}

type MarkAsReadDTO struct {
    MessageID *string `json:"messageId,omitempty"`
}

// Block DTOs
type BlockUserDTO struct {
    Username string `json:"username" validate:"required"`
}

type BlockDTO struct {
    ID          string       `json:"id"`
    BlockedUser *UserShortDTO `json:"blockedUser"`
    CreatedAt   string       `json:"createdAt"`
}

type BlockStatusDTO struct {
    IsBlocked     bool `json:"isBlocked"`
    IsBlockedBy   bool `json:"isBlockedBy"`
    CanSendMessage bool `json:"canSendMessage"`
}

// Pagination
type CursorMeta struct {
    HasMore    bool    `json:"hasMore"`
    NextCursor *string `json:"nextCursor,omitempty"`
}
```

---

#### **Step 1.3: Cursor Pagination Utility**
**เวลา**: 1 day

**ทำอะไร**:
```go
// pkg/utils/cursor.go

type Cursor struct {
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at,omitempty"` // for conversations
    ID        string    `json:"id"`
}

func EncodeCursor(timestamp time.Time, id string, useUpdatedAt bool) (string, error) {
    cursor := Cursor{
        ID: id,
    }

    if useUpdatedAt {
        cursor.UpdatedAt = timestamp
    } else {
        cursor.CreatedAt = timestamp
    }

    jsonBytes, err := json.Marshal(cursor)
    if err != nil {
        return "", err
    }

    encoded := base64.StdEncoding.EncodeToString(jsonBytes)
    return encoded, nil
}

func DecodeCursor(encodedCursor string) (*Cursor, error) {
    if encodedCursor == "" {
        return nil, nil
    }

    jsonBytes, err := base64.StdEncoding.DecodeString(encodedCursor)
    if err != nil {
        return nil, err
    }

    var cursor Cursor
    err = json.Unmarshal(jsonBytes, &cursor)
    if err != nil {
        return nil, err
    }

    return &cursor, nil
}
```

**Test**:
- [ ] Encode/Decode ถูกต้อง
- [ ] Handle empty cursor
- [ ] Handle invalid cursor

---

### **Phase 2: Repository Layer (Week 2-3)**

#### **Step 2.1: Repository Interfaces**
**เวลา**: 1 day

```go
// domain/repositories/conversation_repository.go
type ConversationRepository interface {
    GetByUserID(ctx context.Context, userID uuid.UUID, cursor *utils.Cursor, limit int) ([]*models.Conversation, error)
    GetOrCreate(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Conversation, bool, error)
    GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error)
    Update(ctx context.Context, conv *models.Conversation) error
    IsParticipant(ctx context.Context, convID, userID uuid.UUID) (bool, error)
}

// domain/repositories/message_repository.go
type MessageRepository interface {
    GetByConversationID(ctx context.Context, convID uuid.UUID, cursor *utils.Cursor, limit int) ([]*models.Message, error)
    GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
    Create(ctx context.Context, msg *models.Message) error
    MarkAsRead(ctx context.Context, convID, userID uuid.UUID, messageID *uuid.UUID) (int64, error)
    GetUnreadCount(ctx context.Context, convID, userID uuid.UUID) (int, error)
    GetTotalUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}

// domain/repositories/block_repository.go
type BlockRepository interface {
    Create(ctx context.Context, block *models.Block) error
    Delete(ctx context.Context, blockerID, blockedID uuid.UUID) error
    GetByBlockerID(ctx context.Context, blockerID uuid.UUID, cursor *utils.Cursor, limit int) ([]*models.Block, error)
    IsBlocked(ctx context.Context, user1ID, user2ID uuid.UUID) (bool, bool, error) // (isBlocked, isBlockedBy, error)
}
```

#### **Step 2.2: Repository Implementations**
**เวลา**: 4 days

```go
// infrastructure/postgres/conversation_repository_impl.go
type ConversationRepositoryImpl struct {
    db *gorm.DB
}

func (r *ConversationRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, cursor *utils.Cursor, limit int) ([]*models.Conversation, error) {
    query := r.db.WithContext(ctx).
        Where("user1_id = ? OR user2_id = ?", userID, userID)

    if cursor != nil {
        query = query.Where(
            "(updated_at < ? OR (updated_at = ? AND id < ?))",
            cursor.UpdatedAt, cursor.UpdatedAt, cursor.ID,
        )
    }

    var conversations []*models.Conversation
    err := query.
        Preload("User1").
        Preload("User2").
        Order("updated_at DESC, id DESC").
        Limit(limit + 1). // Fetch one extra to check hasMore
        Find(&conversations).Error

    return conversations, err
}

func (r *ConversationRepositoryImpl) GetOrCreate(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Conversation, bool, error) {
    // Ensure user1 < user2 (for UNIQUE constraint)
    if user1ID.String() > user2ID.String() {
        user1ID, user2ID = user2ID, user1ID
    }

    // Try to find existing
    var conv models.Conversation
    err := r.db.WithContext(ctx).
        Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).
        Preload("User1").
        Preload("User2").
        First(&conv).Error

    if err == nil {
        return &conv, false, nil // Found
    }

    if err != gorm.ErrRecordNotFound {
        return nil, false, err // Error
    }

    // Create new
    conv = models.Conversation{
        ID:      uuid.New(),
        User1ID: user1ID,
        User2ID: user2ID,
    }

    err = r.db.WithContext(ctx).Create(&conv).Error
    if err != nil {
        return nil, false, err
    }

    // Load relationships
    err = r.db.WithContext(ctx).
        Preload("User1").
        Preload("User2").
        First(&conv, "id = ?", conv.ID).Error

    return &conv, true, err // Created
}

// ทำแบบเดียวกันสำหรับ MessageRepository และ BlockRepository
```

**Test**:
- [ ] Unit tests สำหรับแต่ละ repository method
- [ ] Test cursor pagination
- [ ] Test edge cases (empty results, single page, etc.)

---

### **Phase 3: Service Layer (Week 3-4)**

#### **Step 3.1: Service Interfaces**
**เวลา**: 1 day

```go
// domain/services/chat_service.go
type ChatService interface {
    // Conversations
    GetConversations(ctx context.Context, userID uuid.UUID, cursor *utils.Cursor, limit int) ([]*dto.ConversationDTO, *dto.CursorMeta, error)
    GetOrCreateConversation(ctx context.Context, currentUserID uuid.UUID, otherUsername string) (*dto.ConversationDTO, error)
    GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)

    // Messages
    GetMessages(ctx context.Context, convID uuid.UUID, userID uuid.UUID, cursor *utils.Cursor, limit int) ([]*dto.MessageDTO, *dto.CursorMeta, error)
    SendMessage(ctx context.Context, convID uuid.UUID, senderID uuid.UUID, content string) (*dto.MessageDTO, error)
    MarkAsRead(ctx context.Context, convID uuid.UUID, userID uuid.UUID, messageID *uuid.UUID) (int64, error)
    GetMessageByID(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (*dto.MessageDTO, error)
}

// domain/services/block_service.go
type BlockService interface {
    BlockUser(ctx context.Context, blockerID uuid.UUID, blockedUsername string) (*dto.BlockDTO, error)
    UnblockUser(ctx context.Context, blockerID uuid.UUID, blockedUsername string) error
    GetBlockedUsers(ctx context.Context, blockerID uuid.UUID, cursor *utils.Cursor, limit int) ([]*dto.BlockDTO, *dto.CursorMeta, error)
    CheckBlockStatus(ctx context.Context, user1ID uuid.UUID, user2Username string) (*dto.BlockStatusDTO, error)
}

// domain/services/online_status_service.go
type OnlineStatusService interface {
    SetOnline(ctx context.Context, userID uuid.UUID) error
    SetOffline(ctx context.Context, userID uuid.UUID) error
    GetOnlineStatus(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]bool, error)
    GetLastSeen(ctx context.Context, userID uuid.UUID) (*time.Time, error)
}
```

#### **Step 3.2: Service Implementations**
**เวลา**: 5 days

```go
// application/serviceimpl/chat_service_impl.go
type ChatServiceImpl struct {
    convRepo     repositories.ConversationRepository
    msgRepo      repositories.MessageRepository
    blockRepo    repositories.BlockRepository
    userRepo     repositories.UserRepository
    onlineStatusSvc services.OnlineStatusService
    notificationSvc services.NotificationService
    websocketMgr *websocket.Manager
}

func (s *ChatServiceImpl) SendMessage(ctx context.Context, convID uuid.UUID, senderID uuid.UUID, content string) (*dto.MessageDTO, error) {
    // 1. Check conversation exists and user is participant
    conv, err := s.convRepo.GetByID(ctx, convID)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusNotFound, "Conversation not found")
    }

    isParticipant, _ := s.convRepo.IsParticipant(ctx, convID, senderID)
    if !isParticipant {
        return nil, fiber.NewError(fiber.StatusForbidden, "Not a participant")
    }

    // 2. Get receiver ID
    var receiverID uuid.UUID
    if conv.User1ID == senderID {
        receiverID = conv.User2ID
    } else {
        receiverID = conv.User1ID
    }

    // 3. Check if blocked
    isBlocked, isBlockedBy, _ := s.blockRepo.IsBlocked(ctx, senderID, receiverID)
    if isBlocked || isBlockedBy {
        return nil, fiber.NewError(fiber.StatusForbidden, "Cannot send message to this user")
    }

    // 4. Create message
    msg := &models.Message{
        ID:             uuid.New(),
        ConversationID: convID,
        SenderID:       senderID,
        Content:        content,
        IsRead:         false,
    }

    err = s.msgRepo.Create(ctx, msg)
    if err != nil {
        return nil, err
    }

    // 5. Convert to DTO
    msgDTO := s.toMessageDTO(msg)

    // 6. Broadcast via WebSocket
    s.broadcastNewMessage(receiverID, msgDTO)

    // 7. Send notification (if receiver is offline)
    isOnline, _ := s.onlineStatusSvc.GetOnlineStatus(ctx, []uuid.UUID{receiverID})
    if !isOnline[receiverID] {
        s.sendChatNotification(ctx, receiverID, senderID, content)
    }

    return msgDTO, nil
}
```

**Test**:
- [ ] Unit tests สำหรับแต่ละ service method
- [ ] Test business logic (block checking, permissions, etc.)
- [ ] Mock dependencies

---

### **Phase 4: Handler Layer (Week 4-5)**

#### **Step 4.1: HTTP Handlers**
**เวลา**: 4 days

```go
// interfaces/api/handlers/chat_handler.go
type ChatHandler struct {
    chatService  services.ChatService
    blockService services.BlockService
}

func (h *ChatHandler) GetConversations(c *fiber.Ctx) error {
    // Get user from context
    userID := c.Locals("userID").(uuid.UUID)

    // Parse query params
    encodedCursor := c.Query("cursor", "")
    limit := c.QueryInt("limit", 20)

    if limit > 50 {
        limit = 50
    }

    // Decode cursor
    cursor, err := utils.DecodeCursor(encodedCursor)
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid cursor")
    }

    // Get conversations
    convs, meta, err := h.chatService.GetConversations(c.Context(), userID, cursor, limit)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Conversations retrieved successfully",
        "data": fiber.Map{
            "conversations": convs,
            "meta":          meta,
        },
    })
}

// Implement ทุก endpoint ที่ระบุไว้ใน spec
```

**Test**:
- [ ] Integration tests สำหรับแต่ละ endpoint
- [ ] Test auth middleware
- [ ] Test validation
- [ ] Test error handling

---

### **Phase 5: WebSocket Integration (Week 5-6)**

#### **Step 5.1: Chat WebSocket Handler**
**เวลา**: 4 days

```go
// interfaces/api/websocket/chat_handler.go
type ChatWebSocketHandler struct {
    chatService services.ChatService
    onlineStatusSvc services.OnlineStatusService
    manager *websocket.Manager
}

func (h *ChatWebSocketHandler) HandleConnection(c *websocket.Conn) {
    // Handle chat-specific WebSocket events
}

func (h *ChatWebSocketHandler) handleMessage(client *websocket.Client, msgType string, payload map[string]interface{}) {
    switch msgType {
    case "message.send":
        h.handleMessageSend(client, payload)
    case "message.read":
        h.handleMessageRead(client, payload)
    case "ping":
        h.handlePing(client, payload)
    default:
        // Unknown message type
    }
}

func (h *ChatWebSocketHandler) handleMessageSend(client *websocket.Client, payload map[string]interface{}) {
    // 1. Parse payload
    convID := payload["conversationId"].(string)
    content := payload["content"].(string)
    tempID := payload["tempId"].(string) // optional

    // 2. Send message via service
    msgDTO, err := h.chatService.SendMessage(
        context.Background(),
        uuid.MustParse(convID),
        client.UserID,
        content,
    )

    if err != nil {
        // Send error back to client
        client.Send <- []byte(json.Marshal(map[string]interface{}{
            "type": "message.error",
            "payload": map[string]interface{}{
                "tempId": tempID,
                "error": err.Error(),
            },
        }))
        return
    }

    // 3. Send confirmation to sender
    client.Send <- []byte(json.Marshal(map[string]interface{}{
        "type": "message.sent",
        "payload": map[string]interface{}{
            "tempId": tempID,
            "message": msgDTO,
        },
    }))
}
```

#### **Step 5.2: Online Status Tracking**
**เวลา**: 2 days

```go
// application/serviceimpl/online_status_service_impl.go
type OnlineStatusServiceImpl struct {
    redis *redis.Client
}

func (s *OnlineStatusServiceImpl) SetOnline(ctx context.Context, userID uuid.UUID) error {
    key := fmt.Sprintf("online:%s", userID.String())
    return s.redis.Set(ctx, key, time.Now().Unix(), 60*time.Second).Err()
}

func (s *OnlineStatusServiceImpl) SetOffline(ctx context.Context, userID uuid.UUID) error {
    key := fmt.Sprintf("online:%s", userID.String())

    // Set last seen
    lastSeenKey := fmt.Sprintf("last_seen:%s", userID.String())
    s.redis.Set(ctx, lastSeenKey, time.Now().Unix(), 0) // no expiry

    return s.redis.Del(ctx, key).Err()
}

func (s *OnlineStatusServiceImpl) GetOnlineStatus(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
    result := make(map[uuid.UUID]bool)

    for _, userID := range userIDs {
        key := fmt.Sprintf("online:%s", userID.String())
        exists, _ := s.redis.Exists(ctx, key).Result()
        result[userID] = exists > 0
    }

    return result, nil
}
```

**Test**:
- [ ] WebSocket connection test
- [ ] Message delivery test
- [ ] Online status test
- [ ] Heartbeat test
- [ ] Reconnection test

---

### **Phase 6: Routes & DI (Week 6)**

#### **Step 6.1: Route Definitions**
**เวลา**: 1 day

```go
// interfaces/api/routes/chat_routes.go
func SetupChatRoutes(router fiber.Router, container *di.Container) {
    chatHandler := container.ChatHandler
    blockHandler := container.BlockHandler

    chat := router.Group("/chat")
    chat.Use(middleware.Protected()) // All chat routes require auth

    // Conversations
    chat.Get("/conversations", chatHandler.GetConversations)
    chat.Get("/conversations/with/:username", chatHandler.GetOrCreateConversation)
    chat.Get("/conversations/unread-count", chatHandler.GetUnreadCount)

    // Messages
    chat.Get("/conversations/:id/messages", chatHandler.GetMessages)
    chat.Post("/conversations/:id/messages", chatHandler.SendMessage)
    chat.Post("/conversations/:id/read", chatHandler.MarkAsRead)
    chat.Get("/messages/:id", chatHandler.GetMessageByID)

    // Blocks
    chat.Post("/blocks", blockHandler.BlockUser)
    chat.Delete("/blocks/:username", blockHandler.UnblockUser)
    chat.Get("/blocks", blockHandler.GetBlockedUsers)
    chat.Get("/blocks/status/:username", blockHandler.CheckBlockStatus)

    // WebSocket
    chat.Get("/ws", websocket.New(chatWebSocketHandler.HandleConnection))
}
```

#### **Step 6.2: DI Container Updates**
**เวลา**: 1 day

```go
// pkg/di/container.go
type Container struct {
    // ... existing fields

    // Chat repositories
    ConversationRepo repositories.ConversationRepository
    MessageRepo      repositories.MessageRepository
    BlockRepo        repositories.BlockRepository

    // Chat services
    ChatService         services.ChatService
    BlockService        services.BlockService
    OnlineStatusService services.OnlineStatusService

    // Chat handlers
    ChatHandler         *handlers.ChatHandler
    BlockHandler        *handlers.BlockHandler
    ChatWebSocketHandler *websocket.ChatWebSocketHandler
}

func (c *Container) InitializeChatSystem() {
    // Repositories
    c.ConversationRepo = postgres.NewConversationRepository(c.DB)
    c.MessageRepo = postgres.NewMessageRepository(c.DB)
    c.BlockRepo = postgres.NewBlockRepository(c.DB)

    // Services
    c.OnlineStatusService = serviceimpl.NewOnlineStatusService(c.RedisClient)
    c.ChatService = serviceimpl.NewChatService(
        c.ConversationRepo,
        c.MessageRepo,
        c.BlockRepo,
        c.UserRepo,
        c.OnlineStatusService,
        c.NotificationService,
        c.WebSocketMgr,
    )
    c.BlockService = serviceimpl.NewBlockService(c.BlockRepo, c.UserRepo)

    // Handlers
    c.ChatHandler = handlers.NewChatHandler(c.ChatService, c.BlockService)
    c.BlockHandler = handlers.NewBlockHandler(c.BlockService)
    c.ChatWebSocketHandler = websocket.NewChatWebSocketHandler(
        c.ChatService,
        c.OnlineStatusService,
        c.WebSocketMgr,
    )
}
```

---

### **Phase 7: Testing & Polish (Week 7-8)**

#### **Step 7.1: Integration Tests**
**เวลา**: 3 days

#### **Step 7.2: E2E Tests**
**เวลา**: 2 days

#### **Step 7.3: Load Testing**
**เวลา**: 2 days

#### **Step 7.4: Bug Fixes**
**เวลา**: 3 days

---

## 5. Timeline & Resources

### 5.1 ประมาณการเวลา

| Phase | Tasks | Duration | Cumulative |
|-------|-------|----------|------------|
| 1 | Foundation (DB, Models, DTOs) | 2 weeks | 2 weeks |
| 2 | Repository Layer | 1 week | 3 weeks |
| 3 | Service Layer | 1.5 weeks | 4.5 weeks |
| 4 | Handler Layer | 1 week | 5.5 weeks |
| 5 | WebSocket Integration | 1.5 weeks | 7 weeks |
| 6 | Routes & DI | 0.5 week | 7.5 weeks |
| 7 | Testing & Polish | 2 weeks | 9.5 weeks |

**Total: 9-10 สัปดาห์** (ประมาณ 2-2.5 เดือน)

### 5.2 Resources Required

**Backend Developer**: 1 person (full-time)
- Must know: Go, Fiber, PostgreSQL, Redis, WebSocket
- Nice to have: GORM, Clean Architecture experience

**Optional**:
- QA Tester: 0.5 person (for testing phase)
- DevOps: 0.2 person (for deployment support)

---

## 6. Testing Strategy

### 6.1 Unit Tests
- [ ] Repository layer (mock database)
- [ ] Service layer (mock repositories)
- [ ] Cursor encoding/decoding
- [ ] Business logic (block checks, permissions)

### 6.2 Integration Tests
- [ ] API endpoints (with test database)
- [ ] WebSocket connections
- [ ] Database transactions
- [ ] Redis operations

### 6.3 E2E Tests
- [ ] Send/receive messages flow
- [ ] Block/unblock flow
- [ ] Online status updates
- [ ] Pagination (infinite scroll)

### 6.4 Load Tests
- [ ] 100 concurrent WebSocket connections
- [ ] 50 messages/second
- [ ] 1000 API requests/minute

### 6.5 Test Coverage Target
- **Minimum**: 70%
- **Target**: 80%
- **Critical paths**: 100% (send message, block check)

---

## 7. Deployment Plan

### 7.1 Database Migration
```bash
# Run migrations
psql -U postgres -d social_platform -f migrations/008_create_chat_tables.sql

# Verify tables
psql -U postgres -d social_platform -c "\dt"
```

### 7.2 Environment Variables
```bash
# เพิ่มใน .env (ไม่มีตัวแปรใหม่ที่ต้องเพิ่ม - ใช้ของเดิมได้หมด)
# DB, Redis, JWT มีครบแล้ว
```

### 7.3 Deployment Steps
1. [ ] Deploy to staging
2. [ ] Run database migrations
3. [ ] Smoke test (manual)
4. [ ] Load test
5. [ ] Deploy to production (off-peak hours)
6. [ ] Monitor for 24 hours

---

## 8. Risks & Mitigation

### 8.1 Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| WebSocket scaling issues | Medium | High | Load test early, use Redis Pub/Sub for multi-server |
| Database performance | Medium | High | Proper indexes, query optimization, monitor slow queries |
| Race conditions | Medium | Medium | Use database transactions, atomic Redis operations |
| Memory leaks (WebSocket) | Low | High | Connection limits, heartbeat timeouts, monitoring |

### 8.2 Schedule Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Underestimated complexity | High | High | Add 20% buffer time, prioritize core features first |
| Breaking changes | Low | High | Don't touch existing code, add new code only |
| Integration issues | Medium | Medium | Test incrementally, integrate early |

### 8.3 Safety Measures

**ไม่กระทบระบบเดิม**:
- ✅ แยก tables ใหม่ทั้งหมด
- ✅ แยก routes (/api/v1/chat/*)
- ✅ ไม่แก้ไข models เดิม
- ✅ ใช้ services เดิมผ่าน interface (Notification)

**Rollback Plan**:
```sql
-- ถ้ามีปัญหา สามารถ drop tables ได้เลย
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS conversations CASCADE;
DROP TABLE IF EXISTS blocks CASCADE;
DROP FUNCTION IF EXISTS update_conversation_timestamp CASCADE;
```

---

## ✅ Checklist ก่อนเริ่มพัฒนา

### Prerequisites
- [ ] อ่าน spec ทั้งหมดใน `chat_api_spec/` แล้ว
- [ ] เข้าใจ Clean Architecture ของระบบเดิม
- [ ] Setup development environment (Go, PostgreSQL, Redis)
- [ ] Backup database ก่อนทำ migration

### Phase 1 Ready
- [ ] เข้าใจ database schema
- [ ] เข้าใจ cursor pagination
- [ ] เตรียม migration script
- [ ] เตรียม test data

---

## 📝 Notes สำคัญ

1. **ไม่แตะโค้ดเดิม**: ทำได้ 95% โดยไม่แก้ไขโค้ดเดิมเลย เพียงเพิ่มโค้ดใหม่
2. **ทำทีละ Phase**: อย่าข้าม phase ต้อง test แต่ละ phase ให้เรียบร้อยก่อน
3. **Follow Pattern**: ทำตาม Clean Architecture เหมือนโค้ดเดิมทุกประการ
4. **Fiber vs Gin**: Spec ใช้ Gin แต่เราใช้ Fiber - syntax คล้ายกันมาก แค่เปลี่ยน `gin` → `fiber`
5. **Test First**: เขียน test พร้อมกับ implementation ไม่ใช่ทีหลัง

---

## 🎯 Success Criteria

### Must Have (Launch Blockers)
- [ ] 10 REST endpoints ทำงานถูกต้อง
- [ ] WebSocket real-time messaging ทำงาน
- [ ] Cursor pagination ทำงานถูกต้อง (no duplicates)
- [ ] Block system ป้องกัน spam ได้
- [ ] Online status accurate
- [ ] Test coverage > 70%

### Should Have
- [ ] API response < 100ms (p95)
- [ ] WebSocket latency < 50ms
- [ ] Support 100+ concurrent connections
- [ ] Unread count real-time update

### Nice to Have
- [ ] Support 1000+ concurrent connections
- [ ] Cache hit rate > 80%
- [ ] Monitoring dashboard
- [ ] Admin tools

---

## 📚 References

- Chat API Spec: `chat_api_spec/`
- Existing Architecture: `summary_system/02_architecture.md`
- Existing Database: `summary_system/03_database.md`
- Fiber Docs: https://docs.gofiber.io/
- GORM Docs: https://gorm.io/docs/

---

**ถัดไป**: เริ่ม Phase 1 - สร้าง database models และ migrations

**คำถาม?** อ่าน spec ใน `chat_api_spec/` หรือถาม team lead
