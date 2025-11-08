# WebSocket Chat Implementation Plan

> **Status**: 🔴 Ready for Implementation
> **Priority**: CRITICAL - Phase 1 MVP
> **Estimated Time**: 3-4 days
> **Last Updated**: 2025-01-07

---

## 📋 Table of Contents

1. [Current State Analysis](#1-current-state-analysis)
2. [Architecture Design](#2-architecture-design)
3. [Implementation Plan](#3-implementation-plan)
4. [Code Structure](#4-code-structure)
5. [Integration Points](#5-integration-points)
6. [Testing Strategy](#6-testing-strategy)
7. [Timeline & Milestones](#7-timeline--milestones)

---

## 1. Current State Analysis

### ✅ What We Have (Generic WebSocket)

**Files:**
- `infrastructure/websocket/websocket.go` - WebSocketManager (Hub pattern)
- `interfaces/api/websocket/websocket_handler.go` - HTTP upgrade handler
- `interfaces/api/routes/websocket_routes.go` - Route setup

**Features:**
- ✅ Hub pattern (clients, register, unregister, broadcast)
- ✅ Basic ping/pong
- ✅ Room-based messaging (join_room, leave_room)
- ✅ JWT authentication via query param
- ✅ Connection lifecycle management

**Endpoint:**
- Current: `/ws`
- Need: `/chat/ws` (specific for chat)

### ❌ What We're Missing (Chat-Specific)

**Events:**
- ❌ Chat message events (send, sent, new)
- ❌ Read receipt events (read, read_ack, read_update)
- ❌ Online status events (online, offline, status.bulk)
- ❌ Notification events (unread, conversation.updated)
- ❌ Block events (block.add, block.remove)

**Integration:**
- ❌ No integration with MessageService
- ❌ No integration with ConversationService
- ❌ No integration with BlockService
- ❌ No Redis online status tracking
- ❌ No Redis Pub/Sub for multi-server

**Authentication:**
- ❌ No `auth` message after connect
- ❌ No `auth_success` / `auth_failed` response

---

## 2. Architecture Design

### 2.1 Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Layer                         │
│  (Browser WebSocket / React useWebSocket hook)             │
└────────────────────────┬────────────────────────────────────┘
                         │ WSS Connection
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    WebSocket Handler                        │
│  - Upgrade HTTP → WebSocket                                 │
│  - JWT Authentication                                       │
│  - Register client to Hub                                   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Chat Hub                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Clients Map: userID → []*Client                      │  │
│  │ Register/Unregister Channels                         │  │
│  │ Message Router (by event type)                       │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  Event Handlers:                                            │
│  ├─ HandleMessageSend()     → MessageService               │
│  ├─ HandleMessageRead()     → ConversationService          │
│  ├─ HandlePing()            → Update online status         │
│  ├─ HandleBlockAdd()        → BlockService                 │
│  └─ HandleBlockRemove()     → BlockService                 │
└────────────────────────┬────────────────────────────────────┘
                         │
            ┌────────────┼────────────┐
            ▼            ▼            ▼
    ┌──────────┐  ┌──────────┐  ┌──────────┐
    │ Message  │  │Conversa- │  │  Block   │
    │ Service  │  │  tion    │  │ Service  │
    │          │  │ Service  │  │          │
    └─────┬────┘  └─────┬────┘  └─────┬────┘
          │             │             │
          └─────────────┼─────────────┘
                        ▼
              ┌──────────────────┐
              │   PostgreSQL     │
              │   (Messages,     │
              │  Conversations,  │
              │    Blocks)       │
              └──────────────────┘

    ┌───────────────────────────────────────┐
    │         Redis Integration             │
    │                                       │
    │  ┌─────────────────────────────────┐ │
    │  │ Online Status                   │ │
    │  │ - online:{userId} (TTL 60s)     │ │
    │  │ - Updated on ping               │ │
    │  └─────────────────────────────────┘ │
    │                                       │
    │  ┌─────────────────────────────────┐ │
    │  │ Redis Pub/Sub                   │ │
    │  │ - Channel: chat:user:{userId}   │ │
    │  │ - For multi-server broadcast    │ │
    │  └─────────────────────────────────┘ │
    └───────────────────────────────────────┘
```

---

### 2.2 Message Flow Examples

#### Example 1: User A sends message to User B

```
User A (Client)
    │
    │ 1. Send via WebSocket
    │ { type: "message.send", payload: { conversationId, content } }
    │
    ▼
Chat Hub (Server)
    │
    │ 2. Route to HandleMessageSend()
    │
    ▼
MessageService
    │
    │ 3. Validate (check block, permissions)
    │ 4. Save to PostgreSQL
    │ 5. Update conversation (last_message, updated_at)
    │ 6. Increment unread count
    │
    ▼
Chat Hub
    │
    ├─ 7a. Send "message.sent" to User A (sender)
    │      { type: "message.sent", payload: { tempId, message } }
    │
    └─ 7b. Send "message.new" to User B (receiver)
           { type: "message.new", payload: { message } }
           │
           └─ 8. If User B offline → Send push notification
```

---

#### Example 2: User B marks messages as read

```
User B (Client)
    │
    │ 1. Send via WebSocket
    │ { type: "message.read", payload: { conversationId } }
    │
    ▼
Chat Hub
    │
    │ 2. Route to HandleMessageRead()
    │
    ▼
ConversationService
    │
    │ 3. Mark all unread messages as read
    │ 4. Reset unread count in Redis
    │ 5. Update read_at timestamp
    │
    ▼
Chat Hub
    │
    ├─ 6a. Send "message.read_ack" to User B (reader)
    │      { type: "message.read_ack", payload: { markedCount } }
    │
    └─ 6b. Send "message.read_update" to User A (sender)
           { type: "message.read_update", payload: { messageIds, readBy, readAt } }
```

---

## 3. Implementation Plan

### Phase 1: Refactor Infrastructure (Day 1 Morning)

#### Step 1.1: Rename Generic WebSocket
```bash
# Rename to avoid confusion
mv infrastructure/websocket/websocket.go infrastructure/websocket/generic_websocket.go
```

#### Step 1.2: Create Chat Hub Structure
```go
// infrastructure/websocket/chat_hub.go

type ChatHub struct {
    // Client management
    clients         map[uuid.UUID]*ChatClient  // userID → Client
    clientsMutex    sync.RWMutex

    // Channels
    register        chan *ChatClient
    unregister      chan *ChatClient

    // Services (Dependency Injection)
    messageService      services.MessageService
    conversationService services.ConversationService
    blockService        services.BlockService

    // Redis for online status & pub/sub
    redisClient     *redis.Client
    redisPubSub     *redis.PubSub

    // Context for graceful shutdown
    ctx             context.Context
    cancel          context.CancelFunc
}

type ChatClient struct {
    ID              string          // Connection ID (unique per connection)
    UserID          uuid.UUID       // User ID
    Conn            *websocket.Conn
    Send            chan []byte
    Hub             *ChatHub
    LastPing        time.Time
}
```

**Location:** `infrastructure/websocket/chat_hub.go`

---

### Phase 2: Implement Core Hub Logic (Day 1 Afternoon)

#### Step 2.1: Hub Run Loop
```go
func (h *ChatHub) Run() {
    for {
        select {
        case client := <-h.register:
            h.registerClient(client)

        case client := <-h.unregister:
            h.unregisterClient(client)

        case <-h.ctx.Done():
            log.Println("Chat Hub shutting down...")
            return
        }
    }
}

func (h *ChatHub) registerClient(client *ChatClient) {
    h.clientsMutex.Lock()
    defer h.clientsMutex.Unlock()

    h.clients[client.UserID] = client

    // Set user online in Redis
    h.setUserOnline(client.UserID)

    // Broadcast online status to friends
    h.broadcastUserOnline(client.UserID)

    // Send initial bulk status
    h.sendBulkOnlineStatus(client)

    log.Printf("✅ User %s connected", client.UserID)
}

func (h *ChatHub) unregisterClient(client *ChatClient) {
    h.clientsMutex.Lock()
    defer h.clientsMutex.Unlock()

    delete(h.clients, client.UserID)
    close(client.Send)

    // Set user offline in Redis
    h.setUserOffline(client.UserID)

    // Broadcast offline status
    h.broadcastUserOffline(client.UserID)

    log.Printf("❌ User %s disconnected", client.UserID)
}
```

**Location:** `infrastructure/websocket/chat_hub.go`

---

#### Step 2.2: Client Read/Write Pumps
```go
// ReadPump reads messages from WebSocket and routes to handlers
func (c *ChatClient) ReadPump() {
    defer func() {
        c.Hub.unregister <- c
        c.Conn.Close()
    }()

    // Set read deadline
    c.Conn.SetReadDeadline(time.Now().Add(90 * time.Second))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(90 * time.Second))
        return nil
    })

    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }

        // Parse message
        var wsMessage WebSocketMessage
        if err := json.Unmarshal(message, &wsMessage); err != nil {
            c.sendError("INVALID_MESSAGE", "Invalid message format")
            continue
        }

        // Route to handler
        c.Hub.RouteMessage(c, &wsMessage)
    }
}

// WritePump sends messages from Send channel to WebSocket
func (c *ChatClient) WritePump() {
    ticker := time.NewTicker(30 * time.Second)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }

        case <-ticker.C:
            // Send ping every 30 seconds
            c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

**Location:** `infrastructure/websocket/chat_client.go`

---

### Phase 3: Implement Event Handlers (Day 2 Morning)

#### Step 3.1: Message Router
```go
// infrastructure/websocket/message_router.go

func (h *ChatHub) RouteMessage(client *ChatClient, msg *WebSocketMessage) {
    switch msg.Type {
    case "auth":
        h.handleAuth(client, msg.Payload)
    case "ping":
        h.handlePing(client, msg.Payload)
    case "message.send":
        h.handleMessageSend(client, msg.Payload)
    case "message.read":
        h.handleMessageRead(client, msg.Payload)
    case "block.add":
        h.handleBlockAdd(client, msg.Payload)
    case "block.remove":
        h.handleBlockRemove(client, msg.Payload)
    default:
        client.sendError("UNKNOWN_EVENT", fmt.Sprintf("Unknown event type: %s", msg.Type))
    }
}
```

---

#### Step 3.2: Message Send Handler
```go
// infrastructure/websocket/handlers/message_send.go

func (h *ChatHub) handleMessageSend(client *ChatClient, payload map[string]interface{}) {
    // Parse payload
    conversationID, ok := payload["conversationId"].(string)
    if !ok {
        client.sendError("VALIDATION_ERROR", "conversationId is required")
        return
    }

    content, _ := payload["content"].(string)
    tempID, _ := payload["tempId"].(string)

    // Validate content
    if content == "" {
        client.sendError("VALIDATION_ERROR", "content cannot be empty")
        return
    }

    // Parse UUID
    convUUID, err := uuid.Parse(conversationID)
    if err != nil {
        client.sendError("VALIDATION_ERROR", "Invalid conversationId")
        return
    }

    // Create SendMessageRequest
    req := &dto.SendMessageRequest{
        ConversationID: convUUID,
        Type:           "text",
        Content:        &content,
        TempID:         &tempID,
    }

    // Call MessageService
    message, err := h.messageService.SendMessage(context.Background(), client.UserID, req)
    if err != nil {
        if strings.Contains(err.Error(), "blocked") {
            client.sendError("BLOCKED", "Cannot send message: user is blocked")
        } else {
            client.sendError("INTERNAL_ERROR", "Failed to send message")
        }
        return
    }

    // Send confirmation to sender
    client.sendMessage("message.sent", map[string]interface{}{
        "tempId":  tempID,
        "message": message,
    })

    // Get receiver ID
    receiverID, err := h.getOtherUserID(convUUID, client.UserID)
    if err != nil {
        log.Printf("Error getting receiver ID: %v", err)
        return
    }

    // Send to receiver (if online)
    h.sendToUser(receiverID, "message.new", map[string]interface{}{
        "message": message,
    })

    // Update conversation for both users
    h.sendToUser(client.UserID, "conversation.updated", map[string]interface{}{
        "conversationId": conversationID,
        "lastMessage":    message,
    })
    h.sendToUser(receiverID, "conversation.updated", map[string]interface{}{
        "conversationId": conversationID,
        "lastMessage":    message,
    })

    // Send unread count update to receiver
    unreadCount, _ := h.conversationService.GetUnreadCount(context.Background(), receiverID)
    h.sendToUser(receiverID, "notification.unread", map[string]interface{}{
        "conversationId": conversationID,
        "totalUnread":    unreadCount,
    })
}
```

**Location:** `infrastructure/websocket/handlers/message_send.go`

---

#### Step 3.3: Message Read Handler
```go
// infrastructure/websocket/handlers/message_read.go

func (h *ChatHub) handleMessageRead(client *ChatClient, payload map[string]interface{}) {
    conversationID, ok := payload["conversationId"].(string)
    if !ok {
        client.sendError("VALIDATION_ERROR", "conversationId is required")
        return
    }

    convUUID, err := uuid.Parse(conversationID)
    if err != nil {
        client.sendError("VALIDATION_ERROR", "Invalid conversationId")
        return
    }

    // Mark as read via ConversationService
    err = h.conversationService.MarkAsRead(context.Background(), convUUID, client.UserID)
    if err != nil {
        client.sendError("INTERNAL_ERROR", "Failed to mark as read")
        return
    }

    // Send acknowledgement to reader
    client.sendMessage("message.read_ack", map[string]interface{}{
        "conversationId": conversationID,
        "readAt":         time.Now(),
    })

    // Notify sender about read status
    senderID, err := h.getOtherUserID(convUUID, client.UserID)
    if err == nil {
        h.sendToUser(senderID, "message.read_update", map[string]interface{}{
            "conversationId": conversationID,
            "readBy":         client.UserID.String(),
            "readAt":         time.Now(),
        })
    }
}
```

**Location:** `infrastructure/websocket/handlers/message_read.go`

---

#### Step 3.4: Ping Handler
```go
// infrastructure/websocket/handlers/ping.go

func (h *ChatHub) handlePing(client *ChatClient, payload map[string]interface{}) {
    // Update last ping time
    client.LastPing = time.Now()

    // Update Redis online status (extend TTL)
    h.setUserOnline(client.UserID)

    // Send pong response
    client.sendMessage("pong", map[string]interface{}{
        "timestamp":  time.Now(),
        "serverTime": time.Now(),
    })
}
```

**Location:** `infrastructure/websocket/handlers/ping.go`

---

### Phase 4: Redis Integration (Day 2 Afternoon)

#### Step 4.1: Online Status Tracking
```go
// infrastructure/websocket/redis_online.go

func (h *ChatHub) setUserOnline(userID uuid.UUID) error {
    key := fmt.Sprintf("online:%s", userID.String())
    return h.redisClient.Set(context.Background(), key, time.Now().Unix(), 60*time.Second).Err()
}

func (h *ChatHub) setUserOffline(userID uuid.UUID) error {
    key := fmt.Sprintf("online:%s", userID.String())
    return h.redisClient.Set(context.Background(), key, time.Now().Unix(), 0).Err()
}

func (h *ChatHub) isUserOnline(userID uuid.UUID) (bool, time.Time) {
    key := fmt.Sprintf("online:%s", userID.String())
    val, err := h.redisClient.Get(context.Background(), key).Int64()
    if err != nil {
        return false, time.Time{}
    }
    return true, time.Unix(val, 0)
}

func (h *ChatHub) getBulkOnlineStatus(userIDs []uuid.UUID) map[string]bool {
    result := make(map[string]bool)
    for _, userID := range userIDs {
        online, _ := h.isUserOnline(userID)
        result[userID.String()] = online
    }
    return result
}
```

---

#### Step 4.2: Redis Pub/Sub (Multi-Server Support)
```go
// infrastructure/websocket/redis_pubsub.go

func (h *ChatHub) subscribeToPubSub() {
    pubsub := h.redisClient.Subscribe(context.Background(), "chat:broadcast")
    h.redisPubSub = pubsub

    go func() {
        for msg := range pubsub.Channel() {
            var wsMessage WebSocketMessage
            if err := json.Unmarshal([]byte(msg.Payload), &wsMessage); err != nil {
                continue
            }

            // Route message to local clients
            if userID, ok := wsMessage.Payload["userId"].(string); ok {
                uid, _ := uuid.Parse(userID)
                h.sendToUser(uid, wsMessage.Type, wsMessage.Payload)
            }
        }
    }()
}

func (h *ChatHub) publishToRedis(userID uuid.UUID, eventType string, payload map[string]interface{}) {
    msg := WebSocketMessage{
        Type:    eventType,
        Payload: payload,
    }
    msg.Payload["userId"] = userID.String()

    data, _ := json.Marshal(msg)
    h.redisClient.Publish(context.Background(), "chat:broadcast", data)
}
```

---

### Phase 5: Update HTTP Handler (Day 3 Morning)

#### Step 5.1: New Chat WebSocket Handler
```go
// interfaces/api/websocket/chat_websocket_handler.go

type ChatWebSocketHandler struct {
    chatHub *websocket.ChatHub
}

func NewChatWebSocketHandler(chatHub *websocket.ChatHub) *ChatWebSocketHandler {
    return &ChatWebSocketHandler{
        chatHub: chatHub,
    }
}

func (h *ChatWebSocketHandler) HandleChatWebSocket(c *fiber.Ctx) error {
    // Upgrade to WebSocket
    return websocket.New(func(conn *websocket.Conn) {
        // Get token from query or header
        token := conn.Query("token")
        if token == "" {
            token = conn.Headers("Authorization")
            token = strings.TrimPrefix(token, "Bearer ")
        }

        // Validate JWT
        jwtSecret := os.Getenv("JWT_SECRET")
        userCtx, err := utils.ValidateTokenStringToUUID(token, jwtSecret)
        if err != nil {
            log.Printf("❌ WebSocket auth failed: %v", err)
            conn.WriteJSON(map[string]interface{}{
                "type": "auth_failed",
                "payload": map[string]string{
                    "message": "Invalid token",
                },
            })
            conn.Close()
            return
        }

        // Create ChatClient
        client := &websocket.ChatClient{
            ID:       uuid.New().String(),
            UserID:   userCtx.ID,
            Conn:     conn,
            Send:     make(chan []byte, 256),
            Hub:      h.chatHub,
            LastPing: time.Now(),
        }

        // Register client
        h.chatHub.Register <- client

        // Send auth success
        conn.WriteJSON(map[string]interface{}{
            "type": "auth_success",
            "payload": map[string]interface{}{
                "userId":      userCtx.ID.String(),
                "connectedAt": time.Now(),
            },
        })

        // Start pumps
        go client.WritePump()
        client.ReadPump()
    })(c)
}
```

**Location:** `interfaces/api/websocket/chat_websocket_handler.go`

---

#### Step 5.2: Update Routes
```go
// interfaces/api/routes/chat_routes.go

func SetupChatRoutes(api fiber.Router, h *handlers.Handlers) {
    // ... existing REST routes ...

    // WebSocket endpoint
    api.Get("/chat/ws", h.ChatWebSocketHandler.HandleChatWebSocket)
}
```

---

### Phase 6: Container Integration (Day 3 Afternoon)

#### Step 6.1: Update Container
```go
// infrastructure/container/container.go

type Container struct {
    // ... existing fields ...

    // WebSocket
    ChatHub *websocket.ChatHub
}

func NewContainer() (*Container, error) {
    // ... existing initialization ...

    // Initialize ChatHub
    chatHub := websocket.NewChatHub(
        messageService,
        conversationService,
        blockService,
        redisClient,
    )

    // Start hub
    go chatHub.Run()

    return &Container{
        // ... existing fields ...
        ChatHub: chatHub,
    }, nil
}
```

---

#### Step 6.2: Update Handlers Struct
```go
// interfaces/api/handlers/handlers.go

type Handlers struct {
    // ... existing handlers ...
    ChatWebSocketHandler *websocketHandlers.ChatWebSocketHandler
}

func NewHandlers(container *container.Container) *Handlers {
    return &Handlers{
        // ... existing handlers ...
        ChatWebSocketHandler: websocketHandlers.NewChatWebSocketHandler(container.ChatHub),
    }
}
```

---

## 4. Code Structure

### Final File Structure

```
gofiber-backend/
├── infrastructure/
│   └── websocket/
│       ├── chat_hub.go              # Main hub logic
│       ├── chat_client.go           # Client read/write pumps
│       ├── message_router.go        # Event routing
│       ├── redis_online.go          # Online status tracking
│       ├── redis_pubsub.go          # Multi-server support
│       └── handlers/
│           ├── message_send.go      # Handle message.send
│           ├── message_read.go      # Handle message.read
│           ├── ping.go              # Handle ping/pong
│           ├── block_add.go         # Handle block.add
│           └── block_remove.go      # Handle block.remove
│
├── interfaces/api/
│   ├── websocket/
│   │   └── chat_websocket_handler.go  # HTTP → WebSocket upgrade
│   └── routes/
│       └── chat_routes.go           # Updated with /chat/ws
│
└── domain/
    └── dto/
        └── websocket_dto.go         # WebSocket message DTOs
```

---

## 5. Integration Points

### 5.1 Service Dependencies

**MessageService:**
```go
// Already exists - use as-is
messageService.SendMessage(ctx, userID, req)
```

**ConversationService:**
```go
// Already exists - use as-is
conversationService.MarkAsRead(ctx, conversationID, userID)
conversationService.GetUnreadCount(ctx, userID)
```

**BlockService:**
```go
// Already exists - use as-is
blockService.BlockUser(ctx, blockerID, username)
blockService.UnblockUser(ctx, blockerID, username)
```

**Redis:**
```go
// Already configured in container
redisClient.Set(ctx, key, value, ttl)
redisClient.Get(ctx, key)
redisClient.Publish(ctx, channel, message)
```

---

### 5.2 Helper Functions Needed

```go
// Get other user ID in 1-on-1 conversation
func (h *ChatHub) getOtherUserID(conversationID uuid.UUID, currentUserID uuid.UUID) (uuid.UUID, error) {
    conv, err := h.conversationRepo.GetByID(context.Background(), conversationID)
    if err != nil {
        return uuid.Nil, err
    }

    if conv.User1ID == currentUserID {
        return conv.User2ID, nil
    }
    return conv.User1ID, nil
}

// Send message to specific user
func (h *ChatHub) sendToUser(userID uuid.UUID, eventType string, payload map[string]interface{}) {
    h.clientsMutex.RLock()
    client, exists := h.clients[userID]
    h.clientsMutex.RUnlock()

    if !exists {
        // User offline - publish to Redis for other servers
        h.publishToRedis(userID, eventType, payload)
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
        // Buffer full - close connection
        close(client.Send)
        delete(h.clients, userID)
    }
}

// Broadcast online status to friends
func (h *ChatHub) broadcastUserOnline(userID uuid.UUID) {
    // Get user's conversations
    conversations, _ := h.conversationRepo.GetByUserID(context.Background(), userID)

    // Get unique friend IDs
    friendIDs := make(map[uuid.UUID]bool)
    for _, conv := range conversations {
        if conv.User1ID == userID {
            friendIDs[conv.User2ID] = true
        } else {
            friendIDs[conv.User1ID] = true
        }
    }

    // Broadcast to each friend
    for friendID := range friendIDs {
        h.sendToUser(friendID, "user.online", map[string]interface{}{
            "userId":   userID.String(),
            "onlineAt": time.Now(),
        })
    }
}
```

---

## 6. Testing Strategy

### 6.1 Unit Tests

```go
// infrastructure/websocket/chat_hub_test.go

func TestChatHub_RegisterClient(t *testing.T) {
    hub := NewChatHub(mockMessageService, mockConvService, mockBlockService, mockRedis)
    client := &ChatClient{UserID: testUserID}

    hub.registerClient(client)

    assert.Contains(t, hub.clients, testUserID)
}

func TestChatHub_HandleMessageSend(t *testing.T) {
    // Mock services
    // Test message sending
    // Verify service calls
    // Check WebSocket responses
}
```

---

### 6.2 Integration Tests

```go
// tests/integration/websocket_test.go

func TestWebSocket_SendAndReceiveMessage(t *testing.T) {
    // Start test server
    server := setupTestServer()
    defer server.Close()

    // Connect User A
    wsA, _, err := websocket.DefaultDialer.Dial(
        "ws://localhost:8080/chat/ws?token="+tokenA,
        nil,
    )
    require.NoError(t, err)
    defer wsA.Close()

    // Connect User B
    wsB, _, err := websocket.DefaultDialer.Dial(
        "ws://localhost:8080/chat/ws?token="+tokenB,
        nil,
    )
    require.NoError(t, err)
    defer wsB.Close()

    // User A sends message
    wsA.WriteJSON(map[string]interface{}{
        "type": "message.send",
        "payload": map[string]interface{}{
            "conversationId": testConvID,
            "content":        "Hello!",
            "tempId":         "temp-123",
        },
    })

    // User A receives "message.sent"
    var sentMsg WebSocketMessage
    wsA.ReadJSON(&sentMsg)
    assert.Equal(t, "message.sent", sentMsg.Type)

    // User B receives "message.new"
    var newMsg WebSocketMessage
    wsB.ReadJSON(&newMsg)
    assert.Equal(t, "message.new", newMsg.Type)
    assert.Equal(t, "Hello!", newMsg.Payload["message"].(map[string]interface{})["content"])
}
```

---

### 6.3 Load Testing

```javascript
// tests/load/websocket_load_test.js (k6)

import ws from 'k6/ws';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '1m', target: 100 },
    { duration: '2m', target: 500 },
    { duration: '1m', target: 0 },
  ],
};

export default function () {
  const url = 'ws://localhost:8080/chat/ws?token=' + __ENV.TOKEN;

  ws.connect(url, function (socket) {
    socket.on('open', () => {
      // Send auth
      socket.send(JSON.stringify({
        type: 'auth',
        payload: { token: __ENV.TOKEN }
      }));

      // Send messages every 5 seconds
      socket.setInterval(() => {
        socket.send(JSON.stringify({
          type: 'message.send',
          payload: {
            conversationId: 'test-conv-001',
            content: 'Load test message',
            tempId: 'temp-' + Date.now(),
          }
        }));
      }, 5000);
    });

    socket.on('message', (data) => {
      const msg = JSON.parse(data);
      check(msg, {
        'is auth_success': (m) => m.type === 'auth_success',
        'is message.sent': (m) => m.type === 'message.sent',
      });
    });

    socket.setTimeout(() => {
      socket.close();
    }, 30000);
  });
}
```

---

## 7. Timeline & Milestones

### Day 1: Infrastructure Setup (8 hours)

**Morning (4h):**
- ✅ Create ChatHub structure
- ✅ Implement register/unregister logic
- ✅ Setup client management

**Afternoon (4h):**
- ✅ Implement ReadPump/WritePump
- ✅ Create message router
- ✅ Setup basic ping/pong

**Deliverable:** Working WebSocket connection with auth

---

### Day 2: Event Handlers (8 hours)

**Morning (4h):**
- ✅ Implement message.send handler
- ✅ Implement message.read handler
- ✅ Test message flow

**Afternoon (4h):**
- ✅ Implement Redis online status
- ✅ Implement online/offline broadcasting
- ✅ Test online status

**Deliverable:** Working message send/receive and online status

---

### Day 3: Integration & Testing (8 hours)

**Morning (4h):**
- ✅ Update HTTP handler
- ✅ Update routes to /chat/ws
- ✅ Integrate with container
- ✅ Test end-to-end flow

**Afternoon (4h):**
- ✅ Implement block events
- ✅ Add error handling
- ✅ Write unit tests
- ✅ Write integration tests

**Deliverable:** Complete WebSocket implementation with tests

---

### Day 4: Polish & Load Testing (4-6 hours)

**Morning (2-3h):**
- ✅ Redis Pub/Sub for multi-server
- ✅ Graceful shutdown
- ✅ Connection cleanup

**Afternoon (2-3h):**
- ✅ Load testing with k6
- ✅ Performance tuning
- ✅ Documentation

**Deliverable:** Production-ready WebSocket

---

## 8. Rollout Checklist

### Before Starting
- [ ] Review this plan with team
- [ ] Backup current websocket code
- [ ] Create feature branch `feature/chat-websocket`
- [ ] Setup test environment

### Implementation
- [ ] Day 1: Infrastructure ✅
- [ ] Day 2: Event Handlers ✅
- [ ] Day 3: Integration ✅
- [ ] Day 4: Testing & Polish ✅

### Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Load test: 500 concurrent connections
- [ ] Manual testing with 2+ users
- [ ] Test reconnection scenarios
- [ ] Test error handling

### Deployment
- [ ] Merge to develop
- [ ] Deploy to staging
- [ ] Smoke tests on staging
- [ ] Monitor logs and metrics
- [ ] Deploy to production
- [ ] Monitor for 24 hours

---

## 9. Risk Mitigation

### Risk 1: Race Conditions
**Mitigation:** Use proper mutex locking, test with concurrent connections

### Risk 2: Memory Leaks
**Mitigation:** Properly close channels, implement connection cleanup

### Risk 3: Redis Connection Issues
**Mitigation:** Implement reconnection logic, fallback to direct delivery

### Risk 4: High Load
**Mitigation:** Load test early, implement rate limiting, use Redis Pub/Sub

---

## 10. Success Criteria

### Functional Requirements
- ✅ Real-time message delivery (< 100ms latency)
- ✅ Online status tracking
- ✅ Reconnection handling
- ✅ Multi-user support
- ✅ Block functionality via WebSocket

### Performance Requirements
- ✅ Support 500+ concurrent connections
- ✅ Message delivery < 50ms
- ✅ No message loss
- ✅ Graceful degradation

### Quality Requirements
- ✅ 80% code coverage
- ✅ All integration tests pass
- ✅ No memory leaks
- ✅ Proper error handling

---

## 11. Next Steps After Completion

1. **File Upload via REST** (Day 5-6)
   - multipart/form-data handler
   - Integration with Bunny Storage

2. **Frontend Integration** (Week 2)
   - Update useWebSocket hook
   - Connect to real WebSocket
   - Test E2E flows

3. **Phase 2 Features** (Later)
   - Typing indicators
   - Message reactions
   - Group chat

---

**Document Status:** ✅ Ready for Implementation
**Next Action:** Review with team → Start Day 1
**Questions?** Contact Backend Team Lead
