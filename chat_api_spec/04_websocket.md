# Chat API Specification - WebSocket Protocol

## WebSocket Endpoint

### Connection URL
```
Production: wss://api.voobize.com/v1/chat/ws
Development: ws://localhost:8080/v1/chat/ws
```

### Authentication
ส่ง JWT token ตอน handshake:
```
ws://localhost:8080/v1/chat/ws?token=<jwt_token>
```

หรือส่งผ่าน header (preferred):
```javascript
const ws = new WebSocket('wss://api.voobize.com/v1/chat/ws', {
  headers: {
    'Authorization': 'Bearer <jwt_token>'
  }
});
```

---

## Connection Lifecycle

### 1. Connection Establishment

**Client → Server** (First message after connect):
```json
{
  "type": "auth",
  "payload": {
    "token": "<jwt_token>"
  }
}
```

**Server → Client** (Authentication success):
```json
{
  "type": "auth_success",
  "payload": {
    "userId": "user-123",
    "connectedAt": "2024-01-01T10:00:00Z"
  }
}
```

**Server → Client** (Authentication failed):
```json
{
  "type": "auth_failed",
  "payload": {
    "message": "Invalid token"
  }
}
```
*Note: Connection จะถูกปิดทันทีหลัง auth_failed*

---

### 2. Heartbeat / Keep-Alive

**Client → Server** (ทุก 30 วินาที):
```json
{
  "type": "ping",
  "payload": {
    "timestamp": "2024-01-01T10:00:00Z"
  }
}
```

**Server → Client**:
```json
{
  "type": "pong",
  "payload": {
    "timestamp": "2024-01-01T10:00:00Z",
    "serverTime": "2024-01-01T10:00:00Z"
  }
}
```

**Timeout**: ถ้าไม่ได้รับ ping ภายใน 60 วินาที → disconnect และ mark offline

---

### 3. Disconnection

**Client → Server** (Graceful disconnect):
```json
{
  "type": "disconnect",
  "payload": {
    "reason": "user_logout"
  }
}
```

**Server → Client** (Before server disconnect):
```json
{
  "type": "disconnect",
  "payload": {
    "reason": "server_restart",
    "message": "Server is restarting. Please reconnect.",
    "reconnectIn": 5000
  }
}
```

---

## Message Events

### 1. Send Message

**Client → Server**:
```json
{
  "type": "message.send",
  "payload": {
    "conversationId": "conv-001",
    "content": "สวัสดีครับ",
    "tempId": "temp-msg-123"
  }
}
```

**Field Details**:
- `conversationId`: required, UUID
- `content`: required, string, 1-10000 chars
- `tempId`: optional, temporary ID สำหรับ client tracking

---

**Server → Client** (Success - to sender):
```json
{
  "type": "message.sent",
  "payload": {
    "tempId": "temp-msg-123",
    "message": {
      "id": "msg-new",
      "conversationId": "conv-001",
      "senderId": "current-user-id",
      "content": "สวัสดีครับ",
      "isRead": false,
      "readAt": null,
      "createdAt": "2024-01-01T11:00:00Z",
      "updatedAt": "2024-01-01T11:00:00Z"
    }
  }
}
```

---

**Server → Client** (Delivery - to receiver):
```json
{
  "type": "message.new",
  "payload": {
    "message": {
      "id": "msg-new",
      "conversationId": "conv-001",
      "senderId": "user-456",
      "sender": {
        "id": "user-456",
        "username": "somchai",
        "displayName": "สมชาย มีสุข",
        "avatar": "https://cdn.voobize.com/avatars/user-456.jpg"
      },
      "content": "สวัสดีครับ",
      "isRead": false,
      "readAt": null,
      "createdAt": "2024-01-01T11:00:00Z",
      "updatedAt": "2024-01-01T11:00:00Z"
    }
  }
}
```

---

**Server → Client** (Error):
```json
{
  "type": "message.error",
  "payload": {
    "tempId": "temp-msg-123",
    "error": "BLOCKED",
    "message": "You cannot send messages to this user"
  }
}
```

**Error Codes**:
- `BLOCKED`: ถูกบล็อกหรือบล็อกผู้รับ
- `CONVERSATION_NOT_FOUND`: ไม่พบการสนทนา
- `RATE_LIMIT_EXCEEDED`: ส่งข้อความมากเกินไป
- `VALIDATION_ERROR`: ข้อมูลไม่ถูกต้อง
- `INTERNAL_ERROR`: Server error

---

### 2. Mark as Read

**Client → Server**:
```json
{
  "type": "message.read",
  "payload": {
    "conversationId": "conv-001",
    "messageId": "msg-123"
  }
}
```

**Field Details**:
- `conversationId`: required, UUID
- `messageId`: optional, ถ้าไม่ระบุ = mark all unread messages

---

**Server → Client** (Acknowledgement - to reader):
```json
{
  "type": "message.read_ack",
  "payload": {
    "conversationId": "conv-001",
    "markedCount": 3,
    "readAt": "2024-01-01T11:00:00Z"
  }
}
```

---

**Server → Client** (Notification - to sender):
```json
{
  "type": "message.read_update",
  "payload": {
    "conversationId": "conv-001",
    "messageIds": ["msg-121", "msg-122", "msg-123"],
    "readBy": "user-456",
    "readAt": "2024-01-01T11:00:00Z"
  }
}
```

---

### 3. Typing Indicator (Future - Phase 2)

**Client → Server**:
```json
{
  "type": "typing.start",
  "payload": {
    "conversationId": "conv-001"
  }
}
```

**Server → Client** (to other user):
```json
{
  "type": "typing.user",
  "payload": {
    "conversationId": "conv-001",
    "userId": "user-456",
    "username": "somchai",
    "isTyping": true
  }
}
```

---

## Online Status Events

### 1. User Online

**Server → Client** (Broadcast to friends):
```json
{
  "type": "user.online",
  "payload": {
    "userId": "user-456",
    "username": "somchai",
    "onlineAt": "2024-01-01T11:00:00Z"
  }
}
```

---

### 2. User Offline

**Server → Client** (Broadcast to friends):
```json
{
  "type": "user.offline",
  "payload": {
    "userId": "user-456",
    "username": "somchai",
    "lastSeen": "2024-01-01T11:30:00Z"
  }
}
```

**Trigger Conditions**:
- WebSocket disconnected
- ไม่ได้รับ ping ภายใน 60 วินาที
- User logout

---

### 3. Bulk Status Update

**Server → Client** (เมื่อ connect ครั้งแรก):
```json
{
  "type": "status.bulk",
  "payload": {
    "users": [
      {
        "userId": "user-456",
        "isOnline": true,
        "lastSeen": "2024-01-01T11:00:00Z"
      },
      {
        "userId": "user-789",
        "isOnline": false,
        "lastSeen": "2024-01-01T10:30:00Z"
      }
    ]
  }
}
```

**Note**: ส่ง status ของผู้ใช้ทั้งหมดที่มี active conversation

---

## Notification Events

### 1. Unread Count Update

**Server → Client**:
```json
{
  "type": "notification.unread",
  "payload": {
    "conversationId": "conv-001",
    "unreadCount": 2,
    "totalUnread": 5
  }
}
```

**Trigger**: เมื่อมีข้อความใหม่หรือ mark as read

---

### 2. Conversation Updated

**Server → Client**:
```json
{
  "type": "conversation.updated",
  "payload": {
    "conversationId": "conv-001",
    "lastMessage": {
      "id": "msg-new",
      "senderId": "user-456",
      "content": "สวัสดีครับ",
      "createdAt": "2024-01-01T11:00:00Z"
    },
    "unreadCount": 1,
    "updatedAt": "2024-01-01T11:00:00Z"
  }
}
```

**Trigger**: เมื่อมีข้อความใหม่ในการสนทนา

---

## Block Events

### 1. User Blocked

**Client → Server**:
```json
{
  "type": "block.add",
  "payload": {
    "username": "somchai"
  }
}
```

**Server → Client** (Success):
```json
{
  "type": "block.added",
  "payload": {
    "blockId": "block-001",
    "blockedUser": {
      "id": "user-456",
      "username": "somchai"
    },
    "createdAt": "2024-01-01T11:00:00Z"
  }
}
```

---

### 2. User Unblocked

**Client → Server**:
```json
{
  "type": "block.remove",
  "payload": {
    "username": "somchai"
  }
}
```

**Server → Client** (Success):
```json
{
  "type": "block.removed",
  "payload": {
    "username": "somchai",
    "removedAt": "2024-01-01T11:00:00Z"
  }
}
```

---

## Error Handling

### General Error Response

**Server → Client**:
```json
{
  "type": "error",
  "payload": {
    "code": "INTERNAL_ERROR",
    "message": "An unexpected error occurred",
    "details": {
      "requestType": "message.send",
      "timestamp": "2024-01-01T11:00:00Z"
    }
  }
}
```

**Error Codes**:
- `AUTH_FAILED`: Authentication ไม่สำเร็จ
- `INVALID_MESSAGE`: รูปแบบข้อความไม่ถูกต้อง
- `RATE_LIMIT_EXCEEDED`: ส่งข้อความมากเกินไป
- `BLOCKED`: ถูกบล็อกหรือบล็อกผู้ใช้
- `CONVERSATION_NOT_FOUND`: ไม่พบการสนทนา
- `INTERNAL_ERROR`: Server error

---

## Client Implementation

### React/Next.js Example

```typescript
import { useEffect, useRef, useState } from 'react';

interface WebSocketMessage {
  type: string;
  payload: any;
}

export function useWebSocket(token: string) {
  const ws = useRef<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const reconnectTimeout = useRef<NodeJS.Timeout>();

  useEffect(() => {
    const connect = () => {
      const socket = new WebSocket(`ws://localhost:8080/v1/chat/ws?token=${token}`);

      socket.onopen = () => {
        console.log('✅ WebSocket connected');
        setIsConnected(true);

        // Send auth message
        socket.send(JSON.stringify({
          type: 'auth',
          payload: { token }
        }));

        // Start heartbeat
        const heartbeat = setInterval(() => {
          if (socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({
              type: 'ping',
              payload: { timestamp: new Date().toISOString() }
            }));
          }
        }, 30000);

        return () => clearInterval(heartbeat);
      };

      socket.onmessage = (event) => {
        const message: WebSocketMessage = JSON.parse(event.data);
        handleMessage(message);
      };

      socket.onerror = (error) => {
        console.error('❌ WebSocket error:', error);
      };

      socket.onclose = () => {
        console.log('🔌 WebSocket disconnected');
        setIsConnected(false);

        // Auto reconnect after 3 seconds
        reconnectTimeout.current = setTimeout(() => {
          console.log('🔄 Reconnecting...');
          connect();
        }, 3000);
      };

      ws.current = socket;
    };

    connect();

    return () => {
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
      }
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [token]);

  const sendMessage = (type: string, payload: any) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ type, payload }));
    } else {
      console.warn('⚠️ WebSocket not connected');
    }
  };

  const handleMessage = (message: WebSocketMessage) => {
    switch (message.type) {
      case 'auth_success':
        console.log('✅ Authenticated:', message.payload);
        break;
      case 'message.new':
        console.log('📨 New message:', message.payload);
        // Handle new message
        break;
      case 'user.online':
        console.log('🟢 User online:', message.payload);
        // Update online status
        break;
      case 'user.offline':
        console.log('⚪ User offline:', message.payload);
        // Update offline status
        break;
      default:
        console.log('📩 Received:', message);
    }
  };

  return {
    isConnected,
    sendMessage,
    ws: ws.current
  };
}
```

### Usage Example

```typescript
function ChatComponent() {
  const { isConnected, sendMessage } = useWebSocket(authToken);

  const handleSendMessage = (content: string) => {
    sendMessage('message.send', {
      conversationId: 'conv-001',
      content,
      tempId: `temp-${Date.now()}`
    });
  };

  const handleMarkAsRead = (conversationId: string) => {
    sendMessage('message.read', {
      conversationId
    });
  };

  return (
    <div>
      {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
      {/* Chat UI */}
    </div>
  );
}
```

---

## Server Implementation Notes

### Connection Management

```go
type Client struct {
    ID       string
    UserID   string
    Conn     *websocket.Conn
    Send     chan []byte
    Hub      *Hub
}

type Hub struct {
    clients    map[string]*Client           // userID -> Client
    register   chan *Client
    unregister chan *Client
    broadcast  chan *Message
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client.UserID] = client
            h.broadcastOnlineStatus(client.UserID, true)

        case client := <-h.unregister:
            if _, ok := h.clients[client.UserID]; ok {
                delete(h.clients, client.UserID)
                close(client.Send)
                h.broadcastOnlineStatus(client.UserID, false)
            }

        case message := <-h.broadcast:
            h.handleBroadcast(message)
        }
    }
}
```

### Message Broadcasting

```go
func (h *Hub) BroadcastToUser(userID string, message []byte) {
    if client, ok := h.clients[userID]; ok {
        select {
        case client.Send <- message:
        default:
            // Client buffer full, disconnect
            close(client.Send)
            delete(h.clients, userID)
        }
    }
}
```

---

## Performance Optimization

### 1. Connection Pooling
- **Max connections per user**: 3 (web, mobile app, desktop)
- **Idle timeout**: 60 seconds
- **Max message size**: 64 KB

### 2. Message Batching
- รวมข้อความหลายๆ อันส่งพร้อมกันถ้าส่งภายใน 100ms
- ลดจำนวน WebSocket frames

### 3. Redis Pub/Sub
- ใช้สำหรับ broadcast ข้าม WebSocket servers
- Channel pattern: `chat:user:{userId}`

```redis
PUBLISH chat:user:123e4567 "{\"type\":\"message.new\",\"payload\":{...}}"
```

---

## Monitoring & Debugging

### Metrics to Track
- Active connections count
- Messages per second
- Average message latency
- Connection errors
- Reconnection rate

### Debug Events

**Server → Client** (Debug mode only):
```json
{
  "type": "debug",
  "payload": {
    "event": "message_processed",
    "duration": "5ms",
    "timestamp": "2024-01-01T11:00:00Z"
  }
}
```

---

## Security Considerations

### 1. Authentication
- Validate JWT on connection
- Re-validate every 1 hour (send new token via `auth` message)

### 2. Rate Limiting
- Max 30 messages per minute per connection
- Max 1000 messages per hour per user

### 3. Input Validation
- Sanitize all message content
- Validate message type and payload structure
- Limit payload size

### 4. Connection Limits
- Max 3 concurrent connections per user
- Reject if exceeded (close oldest connection)
