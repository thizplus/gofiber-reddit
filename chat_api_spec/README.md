# VOOBIZE Chat API Specification

> Complete specification for Phase 1 MVP Chat System - Real-time 1-on-1 messaging with online status, notifications, and block features.

## 📚 Documentation Index

### 1. [Overview](./01_overview.md)
ภาพรวมระบบ Chat ทั้งหมด รวมถึง:
- วัตถุประสงค์และ features
- เทคโนโลยีที่ใช้ (Go, PostgreSQL, Redis, WebSocket)
- Performance considerations
- Security & scalability plan
- Data flow และ architecture

**อ่านนี้ก่อน** เพื่อเข้าใจภาพรวมทั้งหมด

---

### 2. [Database Schema](./02_database_schema.md)
โครงสร้างฐานข้อมูล PostgreSQL และ Redis:
- **Tables**: conversations, messages, blocks
- **Indexes**: สำหรับ performance optimization
- **Redis schema**: online status, unread counts, cache
- **Triggers**: auto-update conversation timestamps
- **Migrations plan**: ทีละขั้นตอน
- **Storage estimation**: คำนวณพื้นที่ใช้งาน

**Key Points**:
- Cursor-based pagination support
- Denormalized last_message for performance
- Efficient indexes for all queries

---

### 3. [REST API Endpoints](./03_rest_api.md)
API endpoints ทั้งหมด (14 endpoints):

#### Conversations (3 endpoints)
- `GET /chat/conversations` - รายการสนทนา (with pagination)
- `GET /chat/conversations/with/:username` - Get/Create conversation
- `GET /chat/conversations/unread-count` - จำนวนข้อความยังไม่อ่าน

#### Messages (8 endpoints)
- `GET /chat/conversations/:id/messages` - ดึงข้อความ (with pagination)
- `POST /chat/conversations/:id/messages` - ส่งข้อความ
- `POST /chat/conversations/:id/read` - Mark as read
- `GET /chat/messages/:id` - ดึงข้อความเดียว
- `GET /chat/messages/:id/context` - 🆕 Jump to message (with context)
- `GET /chat/conversations/:id/media` - 🆕 รายการ media ทั้งหมด (Phase 2)
- `GET /chat/conversations/:id/links` - 🆕 รายการ links ทั้งหมด (Phase 2)
- `GET /chat/conversations/:id/files` - 🆕 รายการ files ทั้งหมด (Phase 2)

#### Blocks (3 endpoints)
- `POST /chat/blocks` - บล็อกผู้ใช้
- `DELETE /chat/blocks/:username` - ปลดบล็อก
- `GET /chat/blocks` - รายการผู้ใช้ที่ถูกบล็อก
- `GET /chat/blocks/status/:username` - เช็คสถานะการบล็อก

**Key Features**:
- Complete request/response examples
- Error codes และ handling
- Rate limiting rules
- Input validation

---

### 4. [WebSocket Protocol](./04_websocket.md)
WebSocket events และ real-time communication:

#### Connection Lifecycle
- Authentication flow
- Heartbeat/ping-pong (every 30s)
- Graceful disconnect

#### Events (8 events)
**Client → Server**:
- `message.send` - ส่งข้อความ
- `message.read` - Mark as read
- `ping` - Heartbeat

**Server → Client**:
- `message.new` - ข้อความใหม่
- `message.sent` - ส่งสำเร็จ
- `user.online` / `user.offline` - Online status
- `conversation.updated` - Conversation update
- `notification.unread` - Unread count update

**Key Features**:
- Complete message format examples
- Error handling
- React/TypeScript implementation guide
- Go server implementation notes
- Redis Pub/Sub for multi-server support

---

### 5. [Pagination & Infinite Scroll](./05_pagination.md)
Cursor-based pagination strategy:

#### Cursor Design
```json
{
  "created_at": "2024-01-01T10:00:00Z",
  "id": "msg-050"
}
```
Encoded เป็น base64: `eyJjcmVhdGVkX2F0Ij...`

#### Implementation
- **SQL queries**: สำหรับทั้ง conversations และ messages
- **Frontend**: React Query + useInfiniteQuery
- **Backend**: Go encoding/decoding
- **Optimizations**: LIMIT+1 pattern, indexes

#### Features
- Consistent results (ไม่มี duplicates)
- Better performance than offset-based
- Support real-time updates
- Reverse infinite scroll (messages)

---

### 6. [Implementation Plan](./06_implementation_plan.md)
แผนการพัฒนาทีละขั้นตอน (14 weeks):

| Phase | Duration | Focus |
|-------|----------|-------|
| 1 | Week 1-2 | Core Infrastructure (DB, Redis, Auth) |
| 2 | Week 3-4 | REST API Development |
| 3 | Week 5-6 | WebSocket Implementation |
| 4 | Week 7 | Redis Integration & Caching |
| 5 | Week 8 | Notifications Integration |
| 6 | Week 9 | Rate Limiting & Security |
| 7 | Week 10 | Performance Optimization |
| 8 | Week 11 | Monitoring & Logging |
| 9 | Week 12 | Frontend Integration |
| 10 | Week 13 | Testing & Bug Fixes |
| 11 | Week 14 | Documentation & Deployment |

**Includes**:
- Detailed task breakdown
- Code examples
- Testing strategies
- Risk mitigation
- Success criteria

---

## 🚀 Quick Start Guide

### For Backend Developers

1. **Read in this order**:
   - Overview → Database Schema → REST API → WebSocket → Pagination

2. **Start Development**:
   ```bash
   # 1. Setup database
   psql -U postgres -c "CREATE DATABASE voobize_chat;"
   psql -U postgres -d voobize_chat -f migrations/001_initial.sql

   # 2. Setup Redis
   redis-server

   # 3. Clone project structure from Implementation Plan
   # 4. Follow Phase 1 tasks
   ```

3. **Testing**:
   ```bash
   # Run tests
   go test ./...

   # Load testing
   k6 run load-test.js
   ```

### For Frontend Developers

1. **Read**:
   - Overview → REST API → WebSocket → Pagination

2. **Integration**:
   - See section 9 in Implementation Plan
   - Replace mock data with real API calls
   - Implement WebSocket hooks
   - Add infinite scroll

---

## 📊 Architecture Diagram

```
┌─────────────┐         ┌──────────────┐
│  Next.js    │◄──REST──►│   Go API     │
│  Frontend   │         │   (Gin)      │
│             │         │              │
│  WebSocket  │◄────────┤  WebSocket   │
│  Client     │   WSS   │  Hub         │
└─────────────┘         └───────┬──────┘
                                │
                    ┌───────────┼───────────┐
                    ▼           ▼           ▼
              ┌──────────┐ ┌────────┐ ┌────────┐
              │PostgreSQL│ │ Redis  │ │ Redis  │
              │  (Main)  │ │ Cache  │ │Pub/Sub │
              └──────────┘ └────────┘ └────────┘
```

---

## 🎯 Core Features Summary

### Phase 1 MVP (Current)
✅ 1-on-1 Chat (text only)
✅ Real-time messaging (WebSocket)
✅ Online/Offline status
✅ Last seen timestamp
✅ Unread count tracking
✅ Mark messages as read
✅ Block/Unblock users
✅ Infinite scroll (conversations & messages)
✅ Push notifications
✅ Responsive UI (mobile + desktop)

### Phase 1 Limitations
❌ No images/files
❌ No group chat
❌ No message edit/delete
❌ No typing indicators
❌ No read receipts broadcast
❌ No voice/video calls

### Phase 2 Ideas (Future)
- Image/File sharing
- Voice messages
- Group chat
- Typing indicators
- Message reactions
- Message search
- Video calls

---

## 🔧 Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **WebSocket**: gorilla/websocket
- **ORM**: GORM

### Frontend
- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript
- **State**: Zustand + React Query
- **WebSocket**: Native WebSocket API
- **UI**: shadcn/ui + Tailwind CSS

### DevOps
- **Containerization**: Docker
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus + Grafana
- **Logging**: Zap (structured logs)

---

## 📈 Performance Targets

| Metric | Target | Notes |
|--------|--------|-------|
| API Response Time | < 100ms | p95, cached |
| WebSocket Latency | < 50ms | Message delivery |
| Concurrent WS Connections | 1000+ | Per server |
| Database Query Time | < 50ms | With indexes |
| Cache Hit Rate | > 80% | Redis |
| Uptime | 99.9% | ~43min downtime/month |

---

## 🔐 Security Features

- JWT authentication
- Rate limiting (30 msg/min)
- Input sanitization (XSS prevention)
- SQL injection prevention (via ORM)
- HTTPS/WSS only in production
- IP-based blocking
- Block user functionality

---

## 📝 API Examples

### Send Message (REST)
```bash
curl -X POST https://api.voobize.com/v1/chat/conversations/conv-001/messages \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"content": "สวัสดีครับ"}'
```

### Send Message (WebSocket)
```json
{
  "type": "message.send",
  "payload": {
    "conversationId": "conv-001",
    "content": "สวัสดีครับ",
    "tempId": "temp-123"
  }
}
```

---

## 🧪 Testing Strategy

### Backend
- **Unit Tests**: 80% coverage target
- **Integration Tests**: All API endpoints
- **Load Tests**: 1000 concurrent users
- **WebSocket Tests**: Connection, delivery, reconnection

### Frontend
- **Unit Tests**: React Testing Library
- **E2E Tests**: Playwright
- **Visual Tests**: Storybook

---

## 📖 Additional Resources

### Internal Links
- [API Constants](../lib/constants/api.ts)
- [Mock Data](../lib/data/mockChats.ts)
- [Chat Components](../components/chat/)
- [WebSocket Hook](../lib/hooks/useWebSocket.ts)

### External References
- [WebSocket Protocol](https://datatracker.ietf.org/doc/html/rfc6455)
- [Cursor Pagination Best Practices](https://www.postgresql.org/docs/current/queries-limit.html)
- [JWT Authentication](https://jwt.io/introduction)

---

## 🤝 Contributing

### Development Workflow
1. Create feature branch from `develop`
2. Implement following this spec
3. Write tests (unit + integration)
4. Create PR with spec reference
5. Code review
6. Merge to `develop`

### Code Standards
- **Go**: Follow [Effective Go](https://go.dev/doc/effective_go)
- **TypeScript**: ESLint + Prettier
- **Commits**: Conventional Commits

---

## 📞 Support

### Questions?
- Backend API: See [REST API](./03_rest_api.md)
- WebSocket: See [WebSocket Protocol](./04_websocket.md)
- Database: See [Database Schema](./02_database_schema.md)
- Implementation: See [Implementation Plan](./06_implementation_plan.md)

### Issues
- Create GitHub issue with spec reference
- Tag with `chat-api` label

---

## ✅ Checklist Before Starting

Backend Developer:
- [ ] Read all 6 specification documents
- [ ] Setup PostgreSQL and Redis
- [ ] Clone project structure from Implementation Plan
- [ ] Understand cursor pagination logic
- [ ] Review WebSocket protocol

Frontend Developer:
- [ ] Read Overview, REST API, WebSocket docs
- [ ] Review current mock implementation
- [ ] Understand infinite scroll pattern
- [ ] Plan migration from mock to real API

DevOps:
- [ ] Review deployment requirements
- [ ] Setup monitoring infrastructure
- [ ] Prepare CI/CD pipeline

---

## 📅 Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2024-01-XX | Initial specification for Phase 1 MVP |

---

**Last Updated**: January 2025
**Status**: Ready for Development
**Estimated Completion**: 14 weeks from start

---

> 💡 **Tip**: Bookmark this README and refer back frequently during development. Each linked document contains detailed implementation guidelines.
