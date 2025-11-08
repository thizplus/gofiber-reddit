# Chat API Specification - Implementation Plan

## Development Phases

### Phase 1: Core Infrastructure (Week 1-2)

#### 1.1 Database Setup
**Tasks**:
- [ ] Create PostgreSQL migrations
  - `conversations` table
  - `messages` table
  - `blocks` table
  - Indexes and constraints
- [ ] Setup Redis instance
  - Configure persistence (RDB + AOF)
  - Setup Redis Cluster (if needed)
- [ ] Create database triggers
  - `update_conversation_timestamp()`
- [ ] Seed test data
  - 10 users
  - 50 conversations
  - 1000 messages

**Deliverables**:
- Migration files in `migrations/`
- Database schema documentation
- Test data seeder script

**Estimated Time**: 3 days

---

#### 1.2 Project Structure
**Tasks**:
- [ ] Setup Go project structure
  ```
  chat-service/
  ├── cmd/
  │   └── server/
  │       └── main.go
  ├── internal/
  │   ├── handlers/
  │   │   ├── conversation.go
  │   │   ├── message.go
  │   │   ├── block.go
  │   │   └── websocket.go
  │   ├── models/
  │   │   ├── conversation.go
  │   │   ├── message.go
  │   │   └── block.go
  │   ├── repositories/
  │   │   ├── conversation_repo.go
  │   │   ├── message_repo.go
  │   │   └── block_repo.go
  │   ├── services/
  │   │   ├── chat_service.go
  │   │   ├── block_service.go
  │   │   └── notification_service.go
  │   ├── websocket/
  │   │   ├── hub.go
  │   │   ├── client.go
  │   │   └── message.go
  │   ├── cache/
  │   │   └── redis.go
  │   ├── middleware/
  │   │   ├── auth.go
  │   │   ├── rate_limit.go
  │   │   └── cors.go
  │   └── utils/
  │       ├── pagination.go
  │       └── error.go
  ├── pkg/
  │   └── common/
  ├── config/
  │   └── config.yaml
  ├── migrations/
  └── tests/
  ```

- [ ] Install dependencies
  ```bash
  go get github.com/gin-gonic/gin
  go get github.com/gorilla/websocket
  go get gorm.io/gorm
  go get github.com/go-redis/redis/v8
  go get github.com/golang-jwt/jwt
  ```

- [ ] Setup configuration management
  - Environment variables
  - Config file (YAML)
  - Secrets management

**Deliverables**:
- Project scaffold with all folders
- Makefile for common tasks
- `.env.example` file

**Estimated Time**: 2 days

---

#### 1.3 Authentication Middleware
**Tasks**:
- [ ] Integrate existing JWT auth system
- [ ] Create auth middleware
  ```go
  func AuthMiddleware() gin.HandlerFunc {
      return func(c *gin.Context) {
          token := c.GetHeader("Authorization")
          // Validate JWT
          // Set user context
          c.Next()
      }
  }
  ```
- [ ] Test auth flow

**Deliverables**:
- `middleware/auth.go`
- Unit tests for auth

**Estimated Time**: 1 day

---

### Phase 2: REST API Development (Week 3-4)

#### 2.1 Conversations API
**Tasks**:
- [ ] Implement repository layer
  ```go
  type ConversationRepository interface {
      GetByUserID(userID string, cursor *Cursor, limit int) ([]*Conversation, error)
      GetOrCreate(user1ID, user2ID string) (*Conversation, bool, error)
      Update(conv *Conversation) error
  }
  ```

- [ ] Implement service layer
  ```go
  type ChatService interface {
      GetConversations(userID string, cursor *Cursor, limit int) ([]*ConversationDTO, *Meta, error)
      GetOrCreateConversation(currentUserID, otherUsername string) (*ConversationDTO, error)
      GetUnreadCount(userID string) (int, error)
  }
  ```

- [ ] Implement handlers
  - `GET /chat/conversations`
  - `GET /chat/conversations/with/:username`
  - `GET /chat/conversations/unread-count`

- [ ] Add input validation
- [ ] Add error handling
- [ ] Write unit tests
- [ ] Write integration tests

**Deliverables**:
- `repositories/conversation_repo.go`
- `services/chat_service.go`
- `handlers/conversation.go`
- Test files

**Estimated Time**: 4 days

---

#### 2.2 Messages API
**Tasks**:
- [ ] Implement repository layer
  ```go
  type MessageRepository interface {
      GetByConversationID(convID string, cursor *Cursor, limit int) ([]*Message, error)
      Create(msg *Message) error
      MarkAsRead(convID, userID string, messageID *string) (int, error)
      GetByID(messageID string) (*Message, error)
  }
  ```

- [ ] Implement service layer
  ```go
  type MessageService interface {
      GetMessages(convID, userID string, cursor *Cursor, limit int) ([]*MessageDTO, *Meta, error)
      SendMessage(convID, senderID, content string) (*MessageDTO, error)
      MarkAsRead(convID, userID string, messageID *string) (int, error)
  }
  ```

- [ ] Implement handlers
  - `GET /chat/conversations/:conversationId/messages`
  - `POST /chat/conversations/:conversationId/messages`
  - `POST /chat/conversations/:conversationId/read`
  - `GET /chat/messages/:messageId`

- [ ] Add authorization checks
  - User must be participant of conversation
  - Cannot send to blocked users

- [ ] Write tests

**Deliverables**:
- `repositories/message_repo.go`
- `services/message_service.go`
- `handlers/message.go`
- Test files

**Estimated Time**: 5 days

---

#### 2.3 Block API
**Tasks**:
- [ ] Implement repository layer
  ```go
  type BlockRepository interface {
      Create(blockerID, blockedID string) error
      Delete(blockerID, blockedID string) error
      GetByBlockerID(blockerID string, cursor *Cursor, limit int) ([]*Block, error)
      IsBlocked(user1ID, user2ID string) (bool, bool, error)
  }
  ```

- [ ] Implement service layer
  ```go
  type BlockService interface {
      BlockUser(blockerID, blockedUsername string) (*BlockDTO, error)
      UnblockUser(blockerID, blockedUsername string) error
      GetBlockedUsers(blockerID string, cursor *Cursor, limit int) ([]*BlockDTO, *Meta, error)
      CheckBlockStatus(user1ID, user2Username string) (*BlockStatusDTO, error)
  }
  ```

- [ ] Implement handlers
  - `POST /chat/blocks`
  - `DELETE /chat/blocks/:username`
  - `GET /chat/blocks`
  - `GET /chat/blocks/status/:username`

- [ ] Write tests

**Deliverables**:
- `repositories/block_repo.go`
- `services/block_service.go`
- `handlers/block.go`
- Test files

**Estimated Time**: 3 days

---

#### 2.4 Pagination Implementation
**Tasks**:
- [ ] Create cursor utilities
  ```go
  func EncodeCursor(timestamp time.Time, id string) (string, error)
  func DecodeCursor(encoded string) (*Cursor, error)
  ```

- [ ] Apply to all list endpoints
- [ ] Test edge cases
  - Empty results
  - Single page
  - Multiple pages
  - Real-time updates during pagination

**Deliverables**:
- `utils/pagination.go`
- Pagination tests

**Estimated Time**: 2 days

---

### Phase 3: WebSocket Implementation (Week 5-6)

#### 3.1 WebSocket Hub
**Tasks**:
- [ ] Implement Hub struct
  ```go
  type Hub struct {
      clients    map[string]*Client
      register   chan *Client
      unregister chan *Client
      broadcast  chan *Message
      redis      *redis.Client
  }

  func (h *Hub) Run()
  func (h *Hub) BroadcastToUser(userID string, msg []byte)
  ```

- [ ] Implement Client struct
  ```go
  type Client struct {
      ID     string
      UserID string
      Conn   *websocket.Conn
      Send   chan []byte
      Hub    *Hub
  }

  func (c *Client) ReadPump()
  func (c *Client) WritePump()
  ```

- [ ] Handle client connections
  - Register/unregister
  - Authentication
  - Heartbeat

**Deliverables**:
- `websocket/hub.go`
- `websocket/client.go`
- Connection management tests

**Estimated Time**: 4 days

---

#### 3.2 Message Events
**Tasks**:
- [ ] Implement message handlers
  ```go
  func HandleMessageSend(client *Client, payload map[string]interface{})
  func HandleMessageRead(client *Client, payload map[string]interface{})
  func HandlePing(client *Client, payload map[string]interface{})
  ```

- [ ] Implement message broadcasting
  - Send to specific user
  - Send to conversation participants

- [ ] Add message validation
- [ ] Add rate limiting
- [ ] Write tests

**Deliverables**:
- `websocket/message.go`
- Message handler tests

**Estimated Time**: 3 days

---

#### 3.3 Online Status
**Tasks**:
- [ ] Implement Redis-based online tracking
  ```go
  func SetUserOnline(userID string) error
  func SetUserOffline(userID string) error
  func GetOnlineStatus(userIDs []string) (map[string]bool, error)
  func GetLastSeen(userID string) (time.Time, error)
  ```

- [ ] Broadcast online/offline events
- [ ] Implement heartbeat mechanism
- [ ] Add cleanup for stale connections

**Deliverables**:
- `cache/online_status.go`
- Online status tests

**Estimated Time**: 2 days

---

#### 3.4 WebSocket Testing
**Tasks**:
- [ ] Write WebSocket integration tests
  ```go
  func TestWebSocketConnection(t *testing.T)
  func TestMessageDelivery(t *testing.T)
  func TestOnlineStatus(t *testing.T)
  func TestReconnection(t *testing.T)
  ```

- [ ] Test concurrent connections
- [ ] Test error scenarios
- [ ] Load testing (100+ concurrent connections)

**Deliverables**:
- WebSocket test suite
- Load test results

**Estimated Time**: 3 days

---

### Phase 4: Redis Integration & Caching (Week 7)

#### 4.1 Caching Layer
**Tasks**:
- [ ] Implement cache service
  ```go
  type CacheService interface {
      GetConversations(userID string) ([]*Conversation, error)
      SetConversations(userID string, convs []*Conversation) error
      InvalidateConversations(userID string) error

      GetUnreadCount(userID string) (int, error)
      SetUnreadCount(userID string, count int) error
      IncrementUnreadCount(userID, convID string) error
      DecrementUnreadCount(userID, convID string, count int) error

      GetLastMessage(convID string) (*Message, error)
      SetLastMessage(convID string, msg *Message) error
  }
  ```

- [ ] Add cache middleware
- [ ] Implement cache invalidation strategy
- [ ] Add cache warming on server start

**Deliverables**:
- `cache/cache_service.go`
- Cache tests

**Estimated Time**: 3 days

---

#### 4.2 Redis Pub/Sub
**Tasks**:
- [ ] Setup Redis Pub/Sub for multi-server WebSocket
  ```go
  func SubscribeUserChannel(userID string) (*redis.PubSub, error)
  func PublishToUser(userID string, msg []byte) error
  ```

- [ ] Handle cross-server message delivery
- [ ] Test with multiple WebSocket servers

**Deliverables**:
- `cache/pubsub.go`
- Multi-server tests

**Estimated Time**: 2 days

---

### Phase 5: Notifications Integration (Week 8)

#### 5.1 Notification Service
**Tasks**:
- [ ] Integrate with existing notification system
  ```go
  func SendChatNotification(receiverID, senderID, messagePreview string) error
  ```

- [ ] Create notification on new message
- [ ] Only notify if receiver is offline
- [ ] Add notification preferences

**Deliverables**:
- `services/notification_service.go`
- Notification tests

**Estimated Time**: 2 days

---

#### 5.2 Unread Count API
**Tasks**:
- [ ] Implement efficient unread count tracking
- [ ] Update unread count on:
  - New message received
  - Messages marked as read
- [ ] Sync with notification badge

**Deliverables**:
- Unread count tracking
- Sync tests

**Estimated Time**: 2 days

---

### Phase 6: Rate Limiting & Security (Week 9)

#### 6.1 Rate Limiting
**Tasks**:
- [ ] Implement rate limiter middleware
  ```go
  func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc
  ```

- [ ] Apply rate limits:
  - Send message: 30/min
  - Create conversation: 10/min
  - WebSocket messages: 30/min

- [ ] Use Redis for distributed rate limiting

**Deliverables**:
- `middleware/rate_limit.go`
- Rate limit tests

**Estimated Time**: 2 days

---

#### 6.2 Security Hardening
**Tasks**:
- [ ] Input sanitization
  - XSS prevention
  - SQL injection prevention (via ORM)
- [ ] Add CORS configuration
- [ ] Setup HTTPS/WSS
- [ ] Add request validation
- [ ] Implement IP-based blocking
- [ ] Add security headers

**Deliverables**:
- Security middleware
- Security tests
- Security documentation

**Estimated Time**: 3 days

---

### Phase 7: Performance Optimization (Week 10)

#### 7.1 Database Optimization
**Tasks**:
- [ ] Analyze slow queries
- [ ] Add missing indexes
- [ ] Optimize N+1 queries
- [ ] Add database connection pooling
- [ ] Setup read replicas (if needed)

**Deliverables**:
- Performance audit report
- Optimized queries
- Index recommendations

**Estimated Time**: 2 days

---

#### 7.2 Load Testing
**Tasks**:
- [ ] Setup load testing tools (k6, Artillery)
- [ ] Create load test scenarios
  - 1000 concurrent users
  - 100 messages/second
  - 500 WebSocket connections
- [ ] Run tests and analyze results
- [ ] Optimize bottlenecks

**Test Scenarios**:
```javascript
// k6 load test
import http from 'k6/http';
import ws from 'k6/ws';

export let options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 500 },
    { duration: '2m', target: 1000 },
    { duration: '5m', target: 1000 },
    { duration: '2m', target: 0 },
  ],
};

export default function() {
  // Test REST API
  http.get('http://localhost:8080/api/chat/conversations');

  // Test WebSocket
  ws.connect('ws://localhost:8080/api/chat/ws', function(socket) {
    socket.on('open', () => socket.send(JSON.stringify({
      type: 'message.send',
      payload: { conversationId: 'conv-1', content: 'Hello' }
    })));
  });
}
```

**Deliverables**:
- Load test scripts
- Performance reports
- Optimization recommendations

**Estimated Time**: 3 days

---

### Phase 8: Monitoring & Logging (Week 11)

#### 8.1 Structured Logging
**Tasks**:
- [ ] Setup logging library (zap, logrus)
  ```go
  import "go.uber.org/zap"

  logger.Info("Message sent",
      zap.String("user_id", userID),
      zap.String("conversation_id", convID),
      zap.String("message_id", msgID),
  )
  ```

- [ ] Add logging to all handlers
- [ ] Log errors with context
- [ ] Setup log rotation

**Deliverables**:
- Logging configuration
- Log format documentation

**Estimated Time**: 2 days

---

#### 8.2 Monitoring & Metrics
**Tasks**:
- [ ] Setup Prometheus metrics
  ```go
  var (
      messagesTotal = prometheus.NewCounterVec(
          prometheus.CounterOpts{
              Name: "chat_messages_total",
              Help: "Total number of messages sent",
          },
          []string{"conversation_id"},
      )

      wsConnections = prometheus.NewGauge(
          prometheus.GaugeOpts{
              Name: "chat_ws_connections",
              Help: "Number of active WebSocket connections",
          },
      )
  )
  ```

- [ ] Setup Grafana dashboards
- [ ] Add health check endpoints
- [ ] Setup alerts

**Metrics to Track**:
- Active WebSocket connections
- Messages sent per second
- API response times
- Error rates
- Cache hit/miss rates
- Database query times

**Deliverables**:
- Prometheus configuration
- Grafana dashboards
- Alert rules

**Estimated Time**: 3 days

---

### Phase 9: Frontend Integration (Week 12)

#### 9.1 Update Frontend API Client
**Tasks**:
- [ ] Create chat service
  ```typescript
  // lib/services/api/chat.service.ts
  const chatService = {
    getConversations: async (cursor?: string) => { /* ... */ },
    getMessages: async (convId: string, cursor?: string) => { /* ... */ },
    sendMessage: async (convId: string, content: string) => { /* ... */ },
    markAsRead: async (convId: string) => { /* ... */ },
    blockUser: async (username: string) => { /* ... */ },
  };
  ```

- [ ] Replace mock data with real API calls
- [ ] Add error handling
- [ ] Add loading states

**Deliverables**:
- `lib/services/api/chat.service.ts`
- Updated chat components

**Estimated Time**: 2 days

---

#### 9.2 WebSocket Integration
**Tasks**:
- [ ] Update `useWebSocket` hook to use real WebSocket
  ```typescript
  export function useWebSocket() {
    const [ws, setWs] = useState<WebSocket | null>(null);

    useEffect(() => {
      const socket = new WebSocket(
        `${WS_URL}/chat/ws?token=${getToken()}`
      );

      socket.onopen = () => { /* ... */ };
      socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        handleMessage(message);
      };

      setWs(socket);

      return () => socket.close();
    }, []);

    return { ws, isConnected: ws?.readyState === WebSocket.OPEN };
  }
  ```

- [ ] Handle all WebSocket events
- [ ] Add reconnection logic
- [ ] Test real-time features

**Deliverables**:
- Updated `lib/hooks/useWebSocket.ts`
- WebSocket event handlers

**Estimated Time**: 3 days

---

#### 9.3 Infinite Scroll Implementation
**Tasks**:
- [ ] Update `useConversations` hook
  ```typescript
  export function useConversations() {
    return useInfiniteQuery({
      queryKey: ['conversations'],
      queryFn: ({ pageParam }) => chatService.getConversations(pageParam),
      getNextPageParam: (lastPage) => lastPage.data.meta.nextCursor,
    });
  }
  ```

- [ ] Update `useMessages` hook with cursor pagination
- [ ] Implement scroll-to-load logic
- [ ] Add loading indicators
- [ ] Test scroll behavior

**Deliverables**:
- Infinite scroll hooks
- Updated chat UI components

**Estimated Time**: 3 days

---

### Phase 10: Testing & Bug Fixes (Week 13)

#### 10.1 End-to-End Testing
**Tasks**:
- [ ] Write E2E tests (Playwright/Cypress)
  ```typescript
  test('should send and receive messages', async ({ page }) => {
    // Login as User A
    await page.goto('/chat/user-b');

    // Send message
    await page.fill('input[name="message"]', 'Hello!');
    await page.click('button[type="submit"]');

    // Verify message appears
    await expect(page.locator('text=Hello!')).toBeVisible();

    // Login as User B (new context)
    // Verify message received
  });
  ```

- [ ] Test all user flows
  - Send/receive messages
  - Block/unblock users
  - Mark as read
  - Infinite scroll
  - WebSocket reconnection

**Deliverables**:
- E2E test suite
- Test coverage report

**Estimated Time**: 4 days

---

#### 10.2 Bug Fixes & Polish
**Tasks**:
- [ ] Fix bugs found during testing
- [ ] Improve error messages
- [ ] Add loading states
- [ ] Optimize UI performance
- [ ] Cross-browser testing
- [ ] Mobile responsiveness testing

**Deliverables**:
- Bug fix commits
- Polished UI

**Estimated Time**: 3 days

---

### Phase 11: Documentation & Deployment (Week 14)

#### 11.1 Documentation
**Tasks**:
- [ ] API documentation (Swagger/OpenAPI)
- [ ] WebSocket protocol documentation
- [ ] Deployment guide
- [ ] Troubleshooting guide
- [ ] Developer setup guide

**Deliverables**:
- `docs/` folder with all documentation
- Swagger UI for API

**Estimated Time**: 3 days

---

#### 11.2 Deployment
**Tasks**:
- [ ] Setup Docker containers
  ```dockerfile
  FROM golang:1.21-alpine
  WORKDIR /app
  COPY . .
  RUN go build -o server ./cmd/server
  CMD ["./server"]
  ```

- [ ] Create docker-compose for local dev
  ```yaml
  version: '3.8'
  services:
    chat-api:
      build: .
      ports:
        - "8080:8080"
      depends_on:
        - postgres
        - redis

    postgres:
      image: postgres:15-alpine
      environment:
        POSTGRES_DB: voobize_chat
        POSTGRES_PASSWORD: password

    redis:
      image: redis:7-alpine
  ```

- [ ] Setup CI/CD pipeline
- [ ] Deploy to staging
- [ ] Run smoke tests
- [ ] Deploy to production

**Deliverables**:
- Dockerfile
- docker-compose.yml
- CI/CD configuration
- Deployment scripts

**Estimated Time**: 4 days

---

## Summary Timeline

| Phase | Duration | Key Deliverables |
|-------|----------|------------------|
| 1. Core Infrastructure | 1-2 weeks | Database, Project structure, Auth |
| 2. REST API | 2 weeks | All REST endpoints, Pagination |
| 3. WebSocket | 2 weeks | Real-time messaging, Online status |
| 4. Redis Integration | 1 week | Caching, Pub/Sub |
| 5. Notifications | 1 week | Push notifications, Unread count |
| 6. Security | 1 week | Rate limiting, Security hardening |
| 7. Performance | 1 week | Optimization, Load testing |
| 8. Monitoring | 1 week | Logging, Metrics, Alerts |
| 9. Frontend Integration | 1 week | API integration, WebSocket, Infinite scroll |
| 10. Testing | 1 week | E2E tests, Bug fixes |
| 11. Documentation & Deployment | 1 week | Docs, Docker, CI/CD, Deploy |

**Total Estimated Time**: 14 weeks (3.5 months)

---

## Team Structure

### Backend Team (2 developers)
- **Senior Go Developer**: Core API, WebSocket, Redis
- **Junior Go Developer**: Testing, Documentation, Support

### Frontend Team (1 developer)
- **React/Next.js Developer**: Frontend integration, UI polish

### DevOps (0.5 allocation)
- **DevOps Engineer**: CI/CD, Deployment, Monitoring

---

## Risk Mitigation

### Technical Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| WebSocket scalability issues | High | Load test early, use Redis Pub/Sub |
| Database performance | Medium | Add indexes, optimize queries, test with realistic data |
| Real-time sync bugs | High | Comprehensive testing, error handling |
| Race conditions | Medium | Use proper locking, atomic operations |

### Schedule Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| Underestimated complexity | High | Add 20% buffer, prioritize core features |
| Team member unavailable | Medium | Knowledge sharing, pair programming |
| Integration issues | Medium | Weekly integration testing |

---

## Success Criteria

### Performance
- [ ] API response time < 100ms (p95)
- [ ] WebSocket message latency < 50ms
- [ ] Support 1000+ concurrent WebSocket connections
- [ ] Cache hit rate > 80%

### Reliability
- [ ] 99.9% uptime
- [ ] Zero data loss
- [ ] Automatic failover
- [ ] Graceful degradation

### Quality
- [ ] 80% code coverage
- [ ] Zero critical bugs
- [ ] All E2E tests passing
- [ ] Security audit passed

---

## Post-Launch (Phase 2 Ideas)

### Features to Add Later
- [ ] Image/File sharing
- [ ] Voice messages
- [ ] Group chat
- [ ] Typing indicators
- [ ] Message reactions
- [ ] Message search
- [ ] Message edit/delete
- [ ] Video calls
- [ ] Chat themes
- [ ] Message forwarding
- [ ] @mentions
- [ ] Read receipts broadcast
