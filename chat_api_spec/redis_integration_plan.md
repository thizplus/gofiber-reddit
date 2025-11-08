# Redis Integration Plan

> **Status**: 🔴 Ready for Implementation
> **Priority**: HIGH - Phase 1 MVP (Performance Critical)
> **Estimated Time**: 2 days
> **Last Updated**: 2025-01-07

---

## 📋 Overview

Integrate Redis for:
1. **Online Status Tracking** (TTL-based)
2. **Unread Count Caching** (Fast reads)
3. **WebSocket Pub/Sub** (Multi-server support)
4. **Last Message Caching** (Performance boost)

---

## 1. Current State

### ✅ What We Have
- ✅ Redis client configured in container
- ✅ Redis connection warning (but not used)
- ✅ Database queries for all data

### ❌ What We Need
- ❌ Online status tracking with TTL
- ❌ Unread count caching per user
- ❌ Unread count per conversation
- ❌ Last message caching
- ❌ Redis Pub/Sub for WebSocket
- ❌ Cache invalidation strategies

### 🔴 Current Problems (Without Redis)
- Online status requires database query
- Unread count requires COUNT(*) on messages
- Last message requires ORDER BY + LIMIT
- WebSocket broadcast doesn't work across servers
- Performance degrades with scale

---

## 2. Redis Architecture

### 2.1 Key Design

```
┌─────────────────────────────────────────────────────────────┐
│                       Redis Keys                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Online Status (TTL 60s)                                 │
│     Key:   "online:{userId}"                                │
│     Type:  String                                           │
│     Value: Unix timestamp (last seen)                       │
│     TTL:   60 seconds                                       │
│                                                              │
│     Example:                                                │
│     online:550e8400-e29b-41d4-a716-446655440000 = 1704067200│
│                                                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  2. Total Unread Count                                      │
│     Key:   "unread:total:{userId}"                          │
│     Type:  Integer                                          │
│     Value: Total unread messages count                      │
│     TTL:   None (persistent, manual invalidation)           │
│                                                              │
│     Example:                                                │
│     unread:total:550e8400-e29b-41d4-a716-446655440000 = 12  │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  3. Unread Count Per Conversation                           │
│     Key:   "unread:conv:{userId}:{conversationId}"          │
│     Type:  Integer                                          │
│     Value: Unread messages in specific conversation         │
│     TTL:   None (persistent)                                │
│                                                              │
│     Example:                                                │
│     unread:conv:550e...:abc123... = 3                       │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  4. Last Message Cache                                      │
│     Key:   "last_msg:{conversationId}"                      │
│     Type:  Hash                                             │
│     Fields: id, sender_id, content, created_at, type        │
│     TTL:   1 hour                                           │
│                                                              │
│     Example:                                                │
│     HGETALL last_msg:abc123...                              │
│     {                                                        │
│       "id": "msg-001",                                      │
│       "sender_id": "user-456",                              │
│       "content": "Hello!",                                  │
│       "created_at": "2024-01-01T10:00:00Z",                 │
│       "type": "text"                                        │
│     }                                                        │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  5. WebSocket Pub/Sub                                       │
│     Channel: "chat:user:{userId}"                           │
│     Payload: WebSocket message JSON                         │
│                                                              │
│     Example:                                                │
│     PUBLISH chat:user:550e8400... '{"type":"message.new"...}'│
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

### 2.2 Data Flow

#### Scenario 1: User goes online

```
1. User connects to WebSocket
   │
   ├─> ChatHub.registerClient()
   │   └─> SET online:{userId} {timestamp} EX 60
   │
   ├─> Broadcast to friends
   │   └─> For each friend:
   │       └─> Send "user.online" event
   │
   └─> Send bulk status to user
       └─> MGET online:{friendId1} online:{friendId2} ...
```

#### Scenario 2: New message arrives

```
1. User A sends message via WebSocket
   │
   ├─> MessageService.SendMessage()
   │   ├─> Save to PostgreSQL
   │   └─> Update conversation
   │
   ├─> Redis Updates:
   │   ├─> INCR unread:total:{userB}
   │   ├─> INCR unread:conv:{userB}:{convId}
   │   └─> HSET last_msg:{convId} ...
   │
   └─> Broadcast via Redis Pub/Sub:
       └─> PUBLISH chat:user:{userB} '{"type":"message.new"...}'
```

#### Scenario 3: User marks as read

```
1. User B marks conversation as read
   │
   ├─> ConversationService.MarkAsRead()
   │   ├─> Update PostgreSQL (is_read = true)
   │   │
   │   └─> Redis Updates:
   │       ├─> GET unread:conv:{userB}:{convId}  (get count)
   │       ├─> DECRBY unread:total:{userB} {count}
   │       └─> DEL unread:conv:{userB}:{convId}
   │
   └─> Broadcast to sender (User A):
       └─> Send "message.read_update" event
```

---

## 3. Implementation Plan

### Day 1: Online Status & Unread Counts (8 hours)

#### Morning (4h): Online Status Service

**Step 1.1: Create Redis Service**

Location: `infrastructure/redis/redis_service.go`

```go
package redis

import (
    "context"
    "fmt"
    "strconv"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"
)

type RedisService struct {
    client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
    return &RedisService{
        client: client,
    }
}

// ========== Online Status ==========

func (r *RedisService) SetUserOnline(ctx context.Context, userID uuid.UUID) error {
    key := fmt.Sprintf("online:%s", userID.String())
    timestamp := time.Now().Unix()

    return r.client.Set(ctx, key, timestamp, 60*time.Second).Err()
}

func (r *RedisService) SetUserOffline(ctx context.Context, userID uuid.UUID) error {
    key := fmt.Sprintf("online:%s", userID.String())
    timestamp := time.Now().Unix()

    // Set with no TTL (will remain as last seen)
    return r.client.Set(ctx, key, timestamp, 0).Err()
}

func (r *RedisService) IsUserOnline(ctx context.Context, userID uuid.UUID) (bool, time.Time, error) {
    key := fmt.Sprintf("online:%s", userID.String())

    // Check if key exists (with TTL = online)
    ttl, err := r.client.TTL(ctx, key).Result()
    if err != nil {
        return false, time.Time{}, err
    }

    // If TTL > 0, user is online
    isOnline := ttl > 0

    // Get last seen timestamp
    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return false, time.Time{}, nil // Never seen
        }
        return false, time.Time{}, err
    }

    timestamp, _ := strconv.ParseInt(val, 10, 64)
    lastSeen := time.Unix(timestamp, 0)

    return isOnline, lastSeen, nil
}

func (r *RedisService) GetBulkOnlineStatus(ctx context.Context, userIDs []uuid.UUID) (map[string]bool, error) {
    if len(userIDs) == 0 {
        return map[string]bool{}, nil
    }

    // Build keys
    keys := make([]string, len(userIDs))
    for i, id := range userIDs {
        keys[i] = fmt.Sprintf("online:%s", id.String())
    }

    // MGET all keys
    vals, err := r.client.MGet(ctx, keys...).Result()
    if err != nil {
        return nil, err
    }

    // Check TTL for each key to determine online status
    result := make(map[string]bool)
    for i, userID := range userIDs {
        if vals[i] != nil {
            // Check TTL
            ttl, _ := r.client.TTL(ctx, keys[i]).Result()
            result[userID.String()] = ttl > 0
        } else {
            result[userID.String()] = false
        }
    }

    return result, nil
}
```

---

#### Afternoon (4h): Unread Count Service

**Step 1.2: Add Unread Count Methods**

Continue in `infrastructure/redis/redis_service.go`:

```go
// ========== Unread Counts ==========

func (r *RedisService) GetTotalUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
    key := fmt.Sprintf("unread:total:%s", userID.String())

    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return 0, nil // No unread
        }
        return 0, err
    }

    count, err := strconv.Atoi(val)
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (r *RedisService) IncrementTotalUnread(ctx context.Context, userID uuid.UUID) error {
    key := fmt.Sprintf("unread:total:%s", userID.String())
    return r.client.Incr(ctx, key).Err()
}

func (r *RedisService) DecrementTotalUnread(ctx context.Context, userID uuid.UUID, count int) error {
    if count <= 0 {
        return nil
    }

    key := fmt.Sprintf("unread:total:%s", userID.String())

    // DECRBY
    err := r.client.DecrBy(ctx, key, int64(count)).Err()
    if err != nil {
        return err
    }

    // Ensure it doesn't go negative
    val, err := r.client.Get(ctx, key).Int()
    if err == nil && val < 0 {
        r.client.Set(ctx, key, 0, 0)
    }

    return nil
}

func (r *RedisService) GetConversationUnreadCount(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) (int, error) {
    key := fmt.Sprintf("unread:conv:%s:%s", userID.String(), conversationID.String())

    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return 0, nil
        }
        return 0, err
    }

    return strconv.Atoi(val)
}

func (r *RedisService) IncrementConversationUnread(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) error {
    key := fmt.Sprintf("unread:conv:%s:%s", userID.String(), conversationID.String())
    return r.client.Incr(ctx, key).Err()
}

func (r *RedisService) ResetConversationUnread(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) (int, error) {
    key := fmt.Sprintf("unread:conv:%s:%s", userID.String(), conversationID.String())

    // Get current count before deleting
    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return 0, nil // Already 0
        }
        return 0, err
    }

    count, _ := strconv.Atoi(val)

    // Delete key
    r.client.Del(ctx, key)

    return count, nil
}

// Invalidate all unread caches for a user (use when rebuilding)
func (r *RedisService) InvalidateUnreadCache(ctx context.Context, userID uuid.UUID) error {
    // Delete total unread
    totalKey := fmt.Sprintf("unread:total:%s", userID.String())
    r.client.Del(ctx, totalKey)

    // Delete all conversation unread keys
    pattern := fmt.Sprintf("unread:conv:%s:*", userID.String())
    keys, err := r.scanKeys(ctx, pattern)
    if err != nil {
        return err
    }

    if len(keys) > 0 {
        r.client.Del(ctx, keys...)
    }

    return nil
}

func (r *RedisService) scanKeys(ctx context.Context, pattern string) ([]string, error) {
    var keys []string
    iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()

    for iter.Next(ctx) {
        keys = append(keys, iter.Val())
    }

    return keys, iter.Err()
}
```

---

**Step 1.3: Update MessageService to use Redis**

Location: `application/serviceimpl/message_service_impl.go`

```go
type MessageServiceImpl struct {
    messageRepo      repositories.MessageRepository
    conversationRepo repositories.ConversationRepository
    blockRepo        repositories.BlockRepository
    userRepo         repositories.UserRepository
    redisService     *redis.RedisService  // Add this
}

func NewMessageService(
    messageRepo repositories.MessageRepository,
    conversationRepo repositories.ConversationRepository,
    blockRepo repositories.BlockRepository,
    userRepo repositories.UserRepository,
    redisService *redis.RedisService,  // Add this
) services.MessageService {
    return &MessageServiceImpl{
        messageRepo:      messageRepo,
        conversationRepo: conversationRepo,
        blockRepo:        blockRepo,
        userRepo:         userRepo,
        redisService:     redisService,  // Add this
    }
}

// Update SendMessage to increment Redis unread counts
func (s *MessageServiceImpl) SendMessage(ctx context.Context, userID uuid.UUID, req *dto.SendMessageRequest) (*dto.MessageResponse, error) {
    // ... existing validation and message creation ...

    if err := s.messageRepo.Create(ctx, message); err != nil {
        return nil, err
    }

    // Update conversation last message
    _ = s.conversationRepo.UpdateLastMessage(ctx, req.ConversationID, message.ID, now)

    // 🆕 Increment unread counts in Redis
    _ = s.redisService.IncrementTotalUnread(ctx, receiverID)
    _ = s.redisService.IncrementConversationUnread(ctx, receiverID, req.ConversationID)

    // ... rest of the method ...
}
```

---

**Step 1.4: Update ConversationService to use Redis**

Location: `application/serviceimpl/conversation_service_impl.go`

```go
type ConversationServiceImpl struct {
    conversationRepo repositories.ConversationRepository
    messageRepo      repositories.MessageRepository
    userRepo         repositories.UserRepository
    blockRepo        repositories.BlockRepository
    redisService     *redis.RedisService  // Add this
}

// Update GetUnreadCount to use Redis
func (s *ConversationServiceImpl) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
    // Try Redis first
    count, err := s.redisService.GetTotalUnreadCount(ctx, userID)
    if err == nil && count > 0 {
        return count, nil
    }

    // Fallback to database if Redis empty/error
    dbCount, err := s.messageRepo.GetUnreadCount(ctx, userID)
    if err != nil {
        return 0, err
    }

    // Rebuild Redis cache
    if dbCount > 0 {
        // Set total
        key := fmt.Sprintf("unread:total:%s", userID.String())
        s.redisService.client.Set(ctx, key, dbCount, 0)

        // Rebuild per-conversation counts
        conversations, _ := s.conversationRepo.GetByUserID(ctx, userID)
        for _, conv := range conversations {
            convCount, _ := s.messageRepo.GetConversationUnreadCount(ctx, conv.ID, userID)
            if convCount > 0 {
                s.redisService.IncrementConversationUnread(ctx, userID, conv.ID)
            }
        }
    }

    return dbCount, nil
}

// Update MarkAsRead to decrement Redis
func (s *ConversationServiceImpl) MarkAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
    // Get current unread count before marking
    unreadCount, _ := s.redisService.GetConversationUnreadCount(ctx, userID, conversationID)

    // Mark as read in database
    if err := s.messageRepo.MarkAllAsRead(ctx, conversationID, userID); err != nil {
        return err
    }

    // 🆕 Update Redis
    if unreadCount > 0 {
        // Decrement total unread
        _ = s.redisService.DecrementTotalUnread(ctx, userID, unreadCount)

        // Reset conversation unread
        _ = s.redisService.ResetConversationUnread(ctx, userID, conversationID)
    }

    return nil
}
```

---

### Day 2: Caching & Pub/Sub (8 hours)

#### Morning (4h): Last Message Caching

**Step 2.1: Add Last Message Cache**

Continue in `infrastructure/redis/redis_service.go`:

```go
// ========== Last Message Cache ==========

func (r *RedisService) CacheLastMessage(ctx context.Context, conversationID uuid.UUID, message *models.Message) error {
    key := fmt.Sprintf("last_msg:%s", conversationID.String())

    // Store as hash
    data := map[string]interface{}{
        "id":         message.ID.String(),
        "sender_id":  message.SenderID.String(),
        "content":    "",
        "type":       message.Type,
        "created_at": message.CreatedAt.Format(time.RFC3339),
    }

    if message.Content != nil {
        data["content"] = *message.Content
    }

    // HSET with 1 hour TTL
    pipe := r.client.Pipeline()
    pipe.HSet(ctx, key, data)
    pipe.Expire(ctx, key, 1*time.Hour)
    _, err := pipe.Exec(ctx)

    return err
}

func (r *RedisService) GetCachedLastMessage(ctx context.Context, conversationID uuid.UUID) (map[string]string, error) {
    key := fmt.Sprintf("last_msg:%s", conversationID.String())

    result, err := r.client.HGetAll(ctx, key).Result()
    if err != nil {
        return nil, err
    }

    if len(result) == 0 {
        return nil, redis.Nil // Cache miss
    }

    return result, nil
}

func (r *RedisService) InvalidateLastMessage(ctx context.Context, conversationID uuid.UUID) error {
    key := fmt.Sprintf("last_msg:%s", conversationID.String())
    return r.client.Del(ctx, key).Err()
}
```

---

**Step 2.2: Update ConversationService to use cache**

```go
func (s *ConversationServiceImpl) ListConversations(ctx context.Context, userID uuid.UUID, cursorStr *string, limit int) (*dto.ConversationListResponse, error) {
    // ... existing cursor decoding and query ...

    // Fetch conversations
    conversations, err := s.conversationRepo.GetByUserID(ctx, userID, cursor, limit+1)
    if err != nil {
        return nil, err
    }

    // ... hasMore logic ...

    // Convert to DTOs
    convResponses := make([]dto.ConversationResponse, len(conversations))
    for i, conv := range conversations {
        // 🆕 Try to get last message from cache
        cachedLastMsg, err := s.redisService.GetCachedLastMessage(ctx, conv.ID)

        var lastMsgResp *dto.MessageResponse
        if err == nil && len(cachedLastMsg) > 0 {
            // Cache hit - use cached data
            lastMsgResp = s.parseCachedLastMessage(cachedLastMsg)
        } else if conv.LastMessageID != nil {
            // Cache miss - fetch from database
            lastMsg, err := s.messageRepo.GetByID(ctx, *conv.LastMessageID)
            if err == nil {
                lastMsgResp = dto.MessageToMessageResponse(lastMsg)

                // Update cache for next time
                _ = s.redisService.CacheLastMessage(ctx, conv.ID, lastMsg)
            }
        }

        // Get unread count from Redis
        unreadCount, _ := s.redisService.GetConversationUnreadCount(ctx, userID, conv.ID)

        convResponses[i] = *dto.ConversationToConversationResponse(conv, lastMsgResp, unreadCount)
    }

    // ... generate nextCursor and return ...
}

func (s *ConversationServiceImpl) parseCachedLastMessage(cached map[string]string) *dto.MessageResponse {
    id, _ := uuid.Parse(cached["id"])
    senderID, _ := uuid.Parse(cached["sender_id"])
    createdAt, _ := time.Parse(time.RFC3339, cached["created_at"])

    content := cached["content"]
    var contentPtr *string
    if content != "" {
        contentPtr = &content
    }

    return &dto.MessageResponse{
        ID:        id,
        SenderID:  senderID,
        Type:      cached["type"],
        Content:   contentPtr,
        CreatedAt: createdAt,
    }
}
```

---

#### Afternoon (4h): Redis Pub/Sub for WebSocket

**Step 2.3: Add Pub/Sub Methods**

Continue in `infrastructure/redis/redis_service.go`:

```go
// ========== Pub/Sub ==========

func (r *RedisService) PublishToUser(ctx context.Context, userID uuid.UUID, message interface{}) error {
    channel := fmt.Sprintf("chat:user:%s", userID.String())

    data, err := json.Marshal(message)
    if err != nil {
        return err
    }

    return r.client.Publish(ctx, channel, data).Err()
}

func (r *RedisService) SubscribeToUser(ctx context.Context, userID uuid.UUID) *redis.PubSub {
    channel := fmt.Sprintf("chat:user:%s", userID.String())
    return r.client.Subscribe(ctx, channel)
}

func (r *RedisService) UnsubscribeUser(ctx context.Context, pubsub *redis.PubSub) error {
    return pubsub.Close()
}
```

---

**Step 2.4: Integrate with ChatHub**

Location: `infrastructure/websocket/chat_hub.go`

```go
type ChatHub struct {
    // ... existing fields ...
    redisService *redis.RedisService
}

func (h *ChatHub) Run() {
    // Subscribe to global broadcast channel
    pubsub := h.redisService.client.Subscribe(context.Background(), "chat:broadcast")
    defer pubsub.Close()

    // Start goroutine to listen to Redis Pub/Sub
    go func() {
        for msg := range pubsub.Channel() {
            var wsMessage WebSocketMessage
            if err := json.Unmarshal([]byte(msg.Payload), &wsMessage); err != nil {
                continue
            }

            // Extract userID from payload
            if userIDStr, ok := wsMessage.Payload["userId"].(string); ok {
                userID, _ := uuid.Parse(userIDStr)
                h.sendToLocalClient(userID, wsMessage.Type, wsMessage.Payload)
            }
        }
    }()

    // Main hub loop
    for {
        select {
        case client := <-h.register:
            h.registerClient(client)

        case client := <-h.unregister:
            h.unregisterClient(client)

        case <-h.ctx.Done():
            return
        }
    }
}

// Send to local client or publish to Redis if not local
func (h *ChatHub) sendToUser(userID uuid.UUID, eventType string, payload map[string]interface{}) {
    h.clientsMutex.RLock()
    client, exists := h.clients[userID]
    h.clientsMutex.RUnlock()

    if exists {
        // User connected to this server - send directly
        h.sendToLocalClient(userID, eventType, payload)
    } else {
        // User might be on another server - publish to Redis
        payload["userId"] = userID.String()
        h.redisService.PublishToUser(context.Background(), userID, map[string]interface{}{
            "type":    eventType,
            "payload": payload,
        })
    }
}

func (h *ChatHub) sendToLocalClient(userID uuid.UUID, eventType string, payload map[string]interface{}) {
    h.clientsMutex.RLock()
    client, exists := h.clients[userID]
    h.clientsMutex.RUnlock()

    if !exists {
        return
    }

    msg := WebSocketMessage{
        Type:    eventType,
        Payload: payload,
    }

    data, _ := json.Marshal(msg)

    select {
    case client.Send <- data:
    default:
        close(client.Send)
        delete(h.clients, userID)
    }
}
```

---

## 4. Testing Strategy

### Unit Tests

```go
// infrastructure/redis/redis_service_test.go

func TestRedisService_OnlineStatus(t *testing.T) {
    redis := setupTestRedis(t)
    defer redis.Close()

    service := NewRedisService(redis)
    userID := uuid.New()

    // Set online
    err := service.SetUserOnline(context.Background(), userID)
    assert.NoError(t, err)

    // Check online
    isOnline, lastSeen, err := service.IsUserOnline(context.Background(), userID)
    assert.NoError(t, err)
    assert.True(t, isOnline)
    assert.WithinDuration(t, time.Now(), lastSeen, 2*time.Second)

    // Wait for TTL expire
    time.Sleep(61 * time.Second)

    // Check offline
    isOnline, _, _ = service.IsUserOnline(context.Background(), userID)
    assert.False(t, isOnline)
}

func TestRedisService_UnreadCounts(t *testing.T) {
    redis := setupTestRedis(t)
    service := NewRedisService(redis)

    userID := uuid.New()
    convID := uuid.New()

    // Initially 0
    count, err := service.GetTotalUnreadCount(context.Background(), userID)
    assert.NoError(t, err)
    assert.Equal(t, 0, count)

    // Increment 3 times
    service.IncrementTotalUnread(context.Background(), userID)
    service.IncrementTotalUnread(context.Background(), userID)
    service.IncrementTotalUnread(context.Background(), userID)

    count, _ = service.GetTotalUnreadCount(context.Background(), userID)
    assert.Equal(t, 3, count)

    // Decrement 2
    service.DecrementTotalUnread(context.Background(), userID, 2)

    count, _ = service.GetTotalUnreadCount(context.Background(), userID)
    assert.Equal(t, 1, count)
}
```

---

### Integration Tests

```go
// tests/integration/redis_integration_test.go

func TestRedis_MessageUnreadFlow(t *testing.T) {
    // Setup
    server := setupTestServer()
    redis := server.Redis

    user1 := createTestUser(t)
    user2 := createTestUser(t)
    conv := createTestConversation(t, user1.ID, user2.ID)

    // User1 sends message to User2
    sendMessage(t, user1.Token, conv.ID, "Hello!")

    // Check User2 unread count in Redis
    count, err := redis.GetTotalUnreadCount(context.Background(), user2.ID)
    assert.NoError(t, err)
    assert.Equal(t, 1, count)

    convCount, _ := redis.GetConversationUnreadCount(context.Background(), user2.ID, conv.ID)
    assert.Equal(t, 1, convCount)

    // User2 marks as read
    markAsRead(t, user2.Token, conv.ID)

    // Check unread count is 0
    count, _ = redis.GetTotalUnreadCount(context.Background(), user2.ID)
    assert.Equal(t, 0, count)

    convCount, _ = redis.GetConversationUnreadCount(context.Background(), user2.ID, conv.ID)
    assert.Equal(t, 0, convCount)
}
```

---

## 5. Performance Metrics

### Before Redis (Database Only)

| Operation | Time |
|-----------|------|
| Get online status | 50-100ms (DB query) |
| Get unread count | 100-200ms (COUNT query) |
| List conversations | 150-300ms (multiple queries) |

### After Redis

| Operation | Time |
|-----------|------|
| Get online status | < 5ms (Redis GET) |
| Get unread count | < 5ms (Redis GET) |
| List conversations | 50-100ms (cache hits) |

**Performance Improvement**: 5-10x faster

---

## 6. Timeline Summary

| Day | Task | Hours | Status |
|-----|------|-------|--------|
| **Day 1 Morning** | Online Status Service | 4h | ⏳ Pending |
| **Day 1 Afternoon** | Unread Count Service | 4h | ⏳ Pending |
| **Day 2 Morning** | Last Message Caching | 4h | ⏳ Pending |
| **Day 2 Afternoon** | Redis Pub/Sub | 4h | ⏳ Pending |

**Total**: 16 hours (2 days)

---

## 7. Rollout Checklist

### Before Starting
- [ ] Verify Redis is running
- [ ] Test Redis connection
- [ ] Review cache invalidation strategy
- [ ] Create feature branch

### Implementation
- [ ] Day 1: Online Status + Unread ✅
- [ ] Day 2: Caching + Pub/Sub ✅

### Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Load test with Redis
- [ ] Monitor cache hit rates

### Deployment
- [ ] Merge to develop
- [ ] Deploy to staging
- [ ] Monitor Redis memory usage
- [ ] Deploy to production
- [ ] Monitor for 24 hours

---

## 8. Success Criteria

### Functional
- ✅ Online status tracking works
- ✅ Unread counts are accurate
- ✅ Cache invalidation works correctly
- ✅ Multi-server WebSocket works

### Performance
- ✅ Online status query < 5ms
- ✅ Unread count query < 5ms
- ✅ Cache hit rate > 80%
- ✅ Redis memory usage < 100MB (1000 users)

### Quality
- ✅ All tests pass
- ✅ No race conditions
- ✅ Proper error handling
- ✅ Monitoring in place

---

**Document Status:** ✅ Ready for Implementation
**Next Action:** Review → Start Day 1
**Questions?** Contact Backend Team
