# Chat API Specification - Overview

## ภาพรวมระบบ Chat Phase 1 MVP

### วัตถุประสงค์
พัฒนาระบบ Chat แบบ Real-time สำหรับ VOOBIZE Social Platform โดยมีความสามารถหลักดังนี้:
- แชทแบบ 1-on-1 (Direct Message)
- แสดงสถานะออนไลน์/ออฟไลน์
- การแจ้งเตือนข้อความใหม่
- บล็อกผู้ใช้
- Infinite Scroll สำหรับข้อความและรายการสนทนา

### เทคโนโลยีที่ใช้

#### Backend Stack
- **Framework**: Golang (Gin Framework)
- **Database**: PostgreSQL
- **Real-time**: WebSocket
- **Cache**: Redis (สำหรับ online status และ unread counts)
- **Message Queue**: (Optional) RabbitMQ สำหรับ async notifications

#### Communication Patterns
1. **REST API**: สำหรับ CRUD operations และดึงข้อมูลเริ่มต้น
2. **WebSocket**: สำหรับ real-time messaging และ online status
3. **Cursor-based Pagination**: สำหรับ infinite scroll

### Core Features

#### 1. Conversations (รายการสนทนา)
- ดึงรายการสนทนาทั้งหมดของผู้ใช้
- แสดง last message, unread count, updated time
- Infinite scroll แบบ cursor-based pagination
- Real-time update เมื่อมีข้อความใหม่

#### 2. Messages (ข้อความ)
- ส่งข้อความแบบ text (Phase 1)
- ดึงประวัติข้อความย้อนหลัง (infinite scroll)
- แสดงสถานะอ่าน/ไม่อ่าน
- Real-time delivery ผ่าน WebSocket
- Mark as read

#### 3. Online Status (สถานะออนไลน์)
- เก็บใน Redis สำหรับ performance
- Update ทุก 30 วินาที (heartbeat)
- Broadcast สถานะเมื่อเปลี่ยนแปลง
- แสดง "last seen" เมื่อออฟไลน์

#### 4. Notifications (การแจ้งเตือน)
- แจ้งเตือนข้อความใหม่
- Unread count per conversation
- Total unread count
- ใช้ระบบ notification ที่มีอยู่แล้วใน VOOBIZE

#### 5. Block User (บล็อกผู้ใช้)
- บล็อก/ปลดบล็อกผู้ใช้
- ผู้ที่ถูกบล็อกจะไม่สามารถส่งข้อความมาได้
- ซ่อนการสนทนากับผู้ที่ถูกบล็อก

### Performance Considerations

#### Database Optimization
- **Indexes**: conversation_id, sender_id, receiver_id, created_at
- **Partitioning**: แบ่ง messages table ตาม created_at (หากมีข้อมูลมาก)
- **Archiving**: ย้ายข้อความเก่า (>1 ปี) ไป archive table

#### Caching Strategy
- **Redis Keys**:
  - `online:{user_id}` - Online status (TTL 60s)
  - `unread:{user_id}` - Total unread count
  - `unread:{user_id}:{conversation_id}` - Unread per conversation
  - `last_message:{conversation_id}` - Cache last message

#### WebSocket Optimization
- Connection pooling
- Message batching (รวมข้อความหลายๆ อันส่งพร้อมกัน)
- Automatic reconnection logic
- Heartbeat every 30s

### Security Considerations
- Authentication ผ่าน JWT token (ใช้ระบบ auth ที่มีอยู่)
- Authorization: ตรวจสอบสิทธิ์ก่อนดึง/ส่งข้อความ
- Rate limiting: จำกัดจำนวนข้อความต่อนาที
- Input validation: ป้องกัน XSS, SQL Injection
- WebSocket authentication: ส่ง token เมื่อ connect

### Scalability Plan (Future)
- **Horizontal Scaling**:
  - WebSocket servers หลายตัว
  - ใช้ Redis Pub/Sub สำหรับ broadcast ข้าม servers
- **Database Sharding**: แบ่ง database ตาม user_id
- **CDN**: สำหรับ media files (Phase 2)
- **Load Balancer**: Sticky sessions สำหรับ WebSocket

### Data Flow

#### การส่งข้อความ
```
Client A ──(WebSocket)──> Server ──(Save)──> Database
                            │
                            ├──(Cache)──> Redis
                            │
                            └──(WebSocket)──> Client B
```

#### การดึงข้อมูล
```
Client ──(REST)──> Server ──(Check Cache)──> Redis
                      │                         ↓
                      │                      (Miss)
                      │                         ↓
                      └────────────────> Database
                                              ↓
                                         (Update Cache)
```

### API Surface Summary
- **REST Endpoints**: 10 endpoints
  - 3 Conversations APIs
  - 4 Messages APIs
  - 2 Block APIs
  - 1 Online Status API
- **WebSocket Events**: 8 events
  - 4 Client → Server
  - 4 Server → Client

### Phase 1 Limitations
- Text messages only (ไม่มี images, files, voice)
- 1-on-1 chat only (ไม่มี group chat)
- No message edit/delete
- No typing indicator
- No read receipts broadcast (แค่ mark as read)

### Next Documents
1. `02_database_schema.md` - โครงสร้างฐานข้อมูล
2. `03_rest_api.md` - REST API endpoints
3. `04_websocket.md` - WebSocket events และ protocol
4. `05_pagination.md` - Infinite scroll strategy
5. `06_implementation_plan.md` - แผนการพัฒนาทีละขั้นตอน
