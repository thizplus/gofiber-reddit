# Chat API Specification - REST API Endpoints

## Base URL
```
Production: https://api.voobize.com/v1
Development: http://localhost:8080/v1
```

## Authentication
ทุก endpoint ต้องใช้ JWT authentication ผ่าน header:
```
Authorization: Bearer <token>
```

---

## 1. Conversations APIs

### 1.1 Get Conversations List

**Endpoint**: `GET /chat/conversations`

**Description**: ดึงรายการสนทนาทั้งหมดของผู้ใช้ พร้อม last message และ unread count

**Query Parameters**:
```typescript
{
  cursor?: string;      // Cursor สำหรับ pagination (base64 encoded)
  limit?: number;       // จำนวนรายการต่อหน้า (default: 20, max: 50)
}
```

**Request Example**:
```bash
GET /chat/conversations?limit=20
GET /chat/conversations?cursor=eyJ1cGRhdGVkX2F0IjoiMjAyNC0wMS0wMVQxMDowMDowMFoiLCJpZCI6ImNvbnYtMDAxIn0&limit=20
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Conversations retrieved successfully",
  "data": {
    "conversations": [
      {
        "id": "conv-001",
        "otherUser": {
          "id": "user-123",
          "username": "somchai",
          "displayName": "สมชาย มีสุข",
          "avatar": "https://cdn.voobize.com/avatars/user-123.jpg",
          "isOnline": true,
          "lastSeen": "2024-01-01T10:30:00Z"
        },
        "lastMessage": {
          "id": "msg-001",
          "senderId": "user-123",
          "content": "สวัสดีครับ",
          "createdAt": "2024-01-01T10:00:00Z",
          "isRead": false
        },
        "unreadCount": 2,
        "updatedAt": "2024-01-01T10:00:00Z",
        "isBlocked": false
      }
    ],
    "meta": {
      "hasMore": true,
      "nextCursor": "eyJ1cGRhdGVkX2F0IjoiMjAyNC0wMS0wMVQwOTowMDowMFoiLCJpZCI6ImNvbnYtMDAyIn0"
    }
  }
}
```

**Response Error (401)**:
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "UNAUTHORIZED"
}
```

---

### 1.2 Get or Create Conversation

**Endpoint**: `GET /chat/conversations/with/:username`

**Description**: ดึงการสนทนากับผู้ใช้ที่ระบุ หรือสร้างใหม่ถ้ายังไม่มี

**Path Parameters**:
```typescript
{
  username: string;     // Username ของผู้ใช้ที่ต้องการสนทนา
}
```

**Request Example**:
```bash
GET /chat/conversations/with/somchai
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Conversation retrieved successfully",
  "data": {
    "id": "conv-001",
    "otherUser": {
      "id": "user-123",
      "username": "somchai",
      "displayName": "สมชาย มีสุข",
      "avatar": "https://cdn.voobize.com/avatars/user-123.jpg",
      "isOnline": true,
      "lastSeen": "2024-01-01T10:30:00Z"
    },
    "lastMessage": {
      "id": "msg-001",
      "senderId": "user-123",
      "content": "สวัสดีครับ",
      "createdAt": "2024-01-01T10:00:00Z",
      "isRead": false
    },
    "unreadCount": 2,
    "updatedAt": "2024-01-01T10:00:00Z",
    "isBlocked": false,
    "createdAt": "2023-12-01T08:00:00Z"
  }
}
```

**Response Success (201)** - Created:
```json
{
  "success": true,
  "message": "Conversation created successfully",
  "data": {
    "id": "conv-new",
    "otherUser": { /* ... */ },
    "lastMessage": null,
    "unreadCount": 0,
    "updatedAt": "2024-01-01T11:00:00Z",
    "isBlocked": false,
    "createdAt": "2024-01-01T11:00:00Z"
  }
}
```

**Response Error (403)** - Blocked:
```json
{
  "success": false,
  "message": "You have been blocked by this user",
  "error": "BLOCKED"
}
```

**Response Error (404)**:
```json
{
  "success": false,
  "message": "User not found",
  "error": "USER_NOT_FOUND"
}
```

---

### 1.3 Get Unread Count

**Endpoint**: `GET /chat/conversations/unread-count`

**Description**: ดึงจำนวนข้อความที่ยังไม่ได้อ่านทั้งหมด

**Request Example**:
```bash
GET /chat/conversations/unread-count
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Unread count retrieved successfully",
  "data": {
    "count": 5
  }
}
```

---

## 2. Messages APIs

### 2.1 Get Messages

**Endpoint**: `GET /chat/conversations/:conversationId/messages`

**Description**: ดึงข้อความในการสนทนา พร้อม cursor-based pagination

**Path Parameters**:
```typescript
{
  conversationId: string;   // ID ของการสนทนา
}
```

**Query Parameters**:
```typescript
{
  cursor?: string;          // Cursor สำหรับ pagination
  limit?: number;           // จำนวนข้อความต่อหน้า (default: 50, max: 100)
}
```

**Request Example**:
```bash
GET /chat/conversations/conv-001/messages?limit=50
GET /chat/conversations/conv-001/messages?cursor=eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0wMVQwOTowMDowMFoiLCJpZCI6Im1zZy0wNTAifQ&limit=50
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Messages retrieved successfully",
  "data": {
    "messages": [
      {
        "id": "msg-001",
        "conversationId": "conv-001",
        "senderId": "user-123",
        "content": "สวัสดีครับ",
        "isRead": true,
        "readAt": "2024-01-01T10:05:00Z",
        "createdAt": "2024-01-01T10:00:00Z",
        "updatedAt": "2024-01-01T10:05:00Z"
      },
      {
        "id": "msg-002",
        "conversationId": "conv-001",
        "senderId": "current-user-id",
        "content": "สวัสดีครับ",
        "isRead": true,
        "readAt": "2024-01-01T10:03:00Z",
        "createdAt": "2024-01-01T10:02:00Z",
        "updatedAt": "2024-01-01T10:03:00Z"
      }
    ],
    "meta": {
      "hasMore": true,
      "nextCursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0wMVQwODowMDowMFoiLCJpZCI6Im1zZy0wNTEifQ"
    }
  }
}
```

**Response Error (403)**:
```json
{
  "success": false,
  "message": "You don't have permission to access this conversation",
  "error": "FORBIDDEN"
}
```

**Response Error (404)**:
```json
{
  "success": false,
  "message": "Conversation not found",
  "error": "CONVERSATION_NOT_FOUND"
}
```

---

### 2.2 Send Message (via REST)

**Endpoint**: `POST /chat/conversations/:conversationId/messages`

**Description**: ส่งข้อความใหม่ (สามารถใช้ WebSocket แทนได้)

**Path Parameters**:
```typescript
{
  conversationId: string;   // ID ของการสนทนา
}
```

**Request Body**:
```json
{
  "content": "สวัสดีครับ วันนี้เป็นยังไงบ้าง?"
}
```

**Validation**:
- `content`: required, string, min: 1, max: 10000 characters
- ต้องไม่ถูกบล็อกโดยผู้รับ

**Response Success (201)**:
```json
{
  "success": true,
  "message": "Message sent successfully",
  "data": {
    "id": "msg-new",
    "conversationId": "conv-001",
    "senderId": "current-user-id",
    "content": "สวัสดีครับ วันนี้เป็นยังไงบ้าง?",
    "isRead": false,
    "readAt": null,
    "createdAt": "2024-01-01T11:00:00Z",
    "updatedAt": "2024-01-01T11:00:00Z"
  }
}
```

**Response Error (403)** - Blocked:
```json
{
  "success": false,
  "message": "You cannot send messages to this user",
  "error": "BLOCKED"
}
```

**Response Error (429)** - Rate Limit:
```json
{
  "success": false,
  "message": "Too many messages. Please try again later.",
  "error": "RATE_LIMIT_EXCEEDED"
}
```

---

### 2.3 Mark Messages as Read

**Endpoint**: `POST /chat/conversations/:conversationId/read`

**Description**: ทำเครื่องหมายข้อความทั้งหมดในการสนทนาว่าอ่านแล้ว

**Path Parameters**:
```typescript
{
  conversationId: string;   // ID ของการสนทนา
}
```

**Request Body** (Optional):
```json
{
  "messageId": "msg-123"    // Mark ถึง message นี้ (ถ้าไม่ระบุ = mark all)
}
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Messages marked as read",
  "data": {
    "conversationId": "conv-001",
    "markedCount": 3,
    "readAt": "2024-01-01T11:00:00Z"
  }
}
```

**Response Error (403)**:
```json
{
  "success": false,
  "message": "You don't have permission to access this conversation",
  "error": "FORBIDDEN"
}
```

---

### 2.4 Get Message by ID

**Endpoint**: `GET /chat/messages/:messageId`

**Description**: ดึงข้อความเดียวตาม ID (ใช้สำหรับ deep link)

**Path Parameters**:
```typescript
{
  messageId: string;        // ID ของข้อความ
}
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Message retrieved successfully",
  "data": {
    "id": "msg-001",
    "conversationId": "conv-001",
    "senderId": "user-123",
    "content": "สวัสดีครับ",
    "isRead": true,
    "readAt": "2024-01-01T10:05:00Z",
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T10:05:00Z"
  }
}
```

---

### 2.5 Get Message Context (Jump to Message)

**Endpoint**: `GET /chat/messages/:messageId/context`

**Description**: ดึงข้อความเป้าหมายพร้อมกับข้อความก่อนหน้าและถัดไป (สำหรับ jump to message feature)

**Path Parameters**:
```typescript
{
  messageId: string;        // ID ของข้อความเป้าหมาย
}
```

**Query Parameters**:
```typescript
{
  before?: number;          // จำนวนข้อความก่อนหน้า (default: 20, max: 50)
  after?: number;           // จำนวนข้อความถัดไป (default: 20, max: 50)
}
```

**Request Example**:
```bash
GET /chat/messages/msg-123/context?before=20&after=20
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Message context retrieved successfully",
  "data": {
    "targetMessage": {
      "id": "msg-123",
      "conversationId": "conv-001",
      "senderId": "user-456",
      "content": "This is the target message",
      "isRead": true,
      "readAt": "2024-01-15T10:35:00Z",
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T10:30:00Z"
    },
    "before": [
      {
        "id": "msg-122",
        "conversationId": "conv-001",
        "senderId": "current-user-id",
        "content": "Message before target",
        "isRead": true,
        "createdAt": "2024-01-15T10:29:00Z"
      }
      // ... 19 more messages
    ],
    "after": [
      {
        "id": "msg-124",
        "conversationId": "conv-001",
        "senderId": "user-456",
        "content": "Message after target",
        "isRead": false,
        "createdAt": "2024-01-15T10:31:00Z"
      }
      // ... 19 more messages
    ],
    "cursors": {
      "beforeCursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0xNVQxMDoyMDowMFoiLCJpZCI6Im1zZy0xMDMifQ",
      "afterCursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0xNVQxMDo0MDowMFoiLCJpZCI6Im1zZy0xNDMifQ"
    }
  }
}
```

**Response Error (404)**:
```json
{
  "success": false,
  "message": "Message not found",
  "error": "MESSAGE_NOT_FOUND"
}
```

**Response Error (403)**:
```json
{
  "success": false,
  "message": "You don't have permission to access this conversation",
  "error": "FORBIDDEN"
}
```

**Use Cases**:
- Jump to message จาก search results
- Jump to message จาก media/links/files tab
- Jump to quoted/replied message
- Jump to pinned message

---

### 2.6 Get Conversation Media

**Endpoint**: `GET /chat/conversations/:conversationId/media`

**Description**: ดึงรายการ media (images, videos) ทั้งหมดใน conversation (Telegram-style)

**Path Parameters**:
```typescript
{
  conversationId: string;   // ID ของการสนทนา
}
```

**Query Parameters**:
```typescript
{
  cursor?: string;          // Cursor สำหรับ pagination
  limit?: number;           // จำนวนรายการต่อหน้า (default: 50, max: 100)
  type?: string;            // 'image' | 'video' | 'all' (default: 'all')
}
```

**Request Example**:
```bash
GET /chat/conversations/conv-001/media?limit=50&type=image
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Media retrieved successfully",
  "data": {
    "media": [
      {
        "id": "media-001",
        "messageId": "msg-123",
        "type": "image",
        "url": "https://cdn.bunny.net/images/photo.jpg",
        "thumbnail": "https://cdn.bunny.net/images/photo_thumb.jpg",
        "width": 1920,
        "height": 1080,
        "size": 2048576,
        "fileName": "photo.jpg",
        "mimeType": "image/jpeg",
        "sender": {
          "id": "user-456",
          "username": "somchai",
          "displayName": "สมชาย มีสุข",
          "avatar": "https://cdn.bunny.net/avatars/user-456.jpg"
        },
        "createdAt": "2024-01-15T10:30:00Z"
      },
      {
        "id": "media-002",
        "messageId": "msg-145",
        "type": "video",
        "url": "https://cdn.bunny.net/videos/video.mp4",
        "thumbnail": "https://cdn.bunny.net/videos/video_thumb.jpg",
        "width": 1280,
        "height": 720,
        "duration": 120,
        "size": 10485760,
        "fileName": "video.mp4",
        "mimeType": "video/mp4",
        "sender": {
          "id": "current-user-id",
          "username": "me",
          "displayName": "Me"
        },
        "createdAt": "2024-01-14T15:20:00Z"
      }
    ],
    "meta": {
      "hasMore": true,
      "nextCursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0xNFQxMDowMDowMFoiLCJpZCI6Im1lZGlhLTA1MCJ9",
      "totalCount": 150
    }
  }
}
```

**Response Error (403)**:
```json
{
  "success": false,
  "message": "You don't have permission to access this conversation",
  "error": "FORBIDDEN"
}
```

**Note**: Phase 2 feature (requires file upload implementation)

---

### 2.7 Get Conversation Links

**Endpoint**: `GET /chat/conversations/:conversationId/links`

**Description**: ดึงรายการ links ทั้งหมดใน conversation (Telegram-style)

**Path Parameters**:
```typescript
{
  conversationId: string;   // ID ของการสนทนา
}
```

**Query Parameters**:
```typescript
{
  cursor?: string;          // Cursor สำหรับ pagination
  limit?: number;           // จำนวนรายการต่อหน้า (default: 50, max: 100)
}
```

**Request Example**:
```bash
GET /chat/conversations/conv-001/links?limit=50
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Links retrieved successfully",
  "data": {
    "links": [
      {
        "id": "link-001",
        "messageId": "msg-234",
        "url": "https://github.com/user/repo",
        "domain": "github.com",
        "title": "user/repo: Awesome Project",
        "description": "An awesome open source project",
        "imageUrl": "https://opengraph.githubassets.com/...",
        "sender": {
          "id": "user-456",
          "username": "somchai",
          "displayName": "สมชาย มีสุข"
        },
        "createdAt": "2024-01-15T10:30:00Z"
      },
      {
        "id": "link-002",
        "messageId": "msg-245",
        "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
        "domain": "youtube.com",
        "title": "Rick Astley - Never Gonna Give You Up",
        "description": "Official music video",
        "imageUrl": "https://i.ytimg.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
        "sender": {
          "id": "current-user-id",
          "username": "me"
        },
        "createdAt": "2024-01-14T18:45:00Z"
      }
    ],
    "meta": {
      "hasMore": true,
      "nextCursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0xNFQxMDowMDowMFoiLCJpZCI6ImxpbmstMDUwIn0",
      "totalCount": 87
    }
  }
}
```

**Response Error (403)**:
```json
{
  "success": false,
  "message": "You don't have permission to access this conversation",
  "error": "FORBIDDEN"
}
```

**Implementation Notes**:
- Extract URLs from message content using regex
- Optionally fetch Open Graph metadata (title, description, image)
- Store in `message_links` table for faster queries
- Cache metadata to avoid re-fetching

**Note**: Phase 2 feature

---

### 2.8 Get Conversation Files

**Endpoint**: `GET /chat/conversations/:conversationId/files`

**Description**: ดึงรายการ files ทั้งหมดใน conversation (Telegram-style)

**Path Parameters**:
```typescript
{
  conversationId: string;   // ID ของการสนทนา
}
```

**Query Parameters**:
```typescript
{
  cursor?: string;          // Cursor สำหรับ pagination
  limit?: number;           // จำนวนรายการต่อหน้า (default: 50, max: 100)
  fileType?: string;        // MIME type filter (e.g., 'application/pdf')
}
```

**Request Example**:
```bash
GET /chat/conversations/conv-001/files?limit=50
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Files retrieved successfully",
  "data": {
    "files": [
      {
        "id": "file-001",
        "messageId": "msg-345",
        "fileName": "document.pdf",
        "fileSize": 1024000,
        "fileType": "application/pdf",
        "mimeType": "application/pdf",
        "url": "https://cdn.bunny.net/files/document.pdf",
        "thumbnailUrl": "https://cdn.bunny.net/files/document_thumb.jpg",
        "sender": {
          "id": "user-456",
          "username": "somchai",
          "displayName": "สมชาย มีสุข"
        },
        "createdAt": "2024-01-15T10:30:00Z"
      },
      {
        "id": "file-002",
        "messageId": "msg-356",
        "fileName": "presentation.pptx",
        "fileSize": 5242880,
        "fileType": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
        "mimeType": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
        "url": "https://cdn.bunny.net/files/presentation.pptx",
        "sender": {
          "id": "current-user-id",
          "username": "me"
        },
        "createdAt": "2024-01-14T16:15:00Z"
      }
    ],
    "meta": {
      "hasMore": true,
      "nextCursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0xNFQxMDowMDowMFoiLCJpZCI6ImZpbGUtMDUwIn0",
      "totalCount": 23
    }
  }
}
```

**Response Error (403)**:
```json
{
  "success": false,
  "message": "You don't have permission to access this conversation",
  "error": "FORBIDDEN"
}
```

**Supported File Types**:
- Documents: PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX
- Archives: ZIP, RAR, 7Z
- Text: TXT, CSV, JSON, XML
- Code: JS, TS, PY, GO, etc.

**Note**: Phase 2 feature (requires file upload implementation)

---

## 3. Block APIs

### 3.1 Block User

**Endpoint**: `POST /chat/blocks`

**Description**: บล็อกผู้ใช้

**Request Body**:
```json
{
  "username": "somchai"
}
```

**Response Success (201)**:
```json
{
  "success": true,
  "message": "User blocked successfully",
  "data": {
    "id": "block-001",
    "blockerId": "current-user-id",
    "blockedId": "user-123",
    "blockedUsername": "somchai",
    "createdAt": "2024-01-01T11:00:00Z"
  }
}
```

**Response Error (400)**:
```json
{
  "success": false,
  "message": "User is already blocked",
  "error": "ALREADY_BLOCKED"
}
```

**Response Error (404)**:
```json
{
  "success": false,
  "message": "User not found",
  "error": "USER_NOT_FOUND"
}
```

---

### 3.2 Unblock User

**Endpoint**: `DELETE /chat/blocks/:username`

**Description**: ปลดบล็อกผู้ใช้

**Path Parameters**:
```typescript
{
  username: string;         // Username ของผู้ใช้ที่ต้องการปลดบล็อก
}
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "User unblocked successfully",
  "data": null
}
```

**Response Error (404)**:
```json
{
  "success": false,
  "message": "Block not found",
  "error": "BLOCK_NOT_FOUND"
}
```

---

### 3.3 Get Blocked Users

**Endpoint**: `GET /chat/blocks`

**Description**: ดึงรายการผู้ใช้ที่ถูกบล็อก

**Query Parameters**:
```typescript
{
  cursor?: string;
  limit?: number;           // default: 20, max: 50
}
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Blocked users retrieved successfully",
  "data": {
    "blocks": [
      {
        "id": "block-001",
        "blockedUser": {
          "id": "user-123",
          "username": "somchai",
          "displayName": "สมชาย มีสุข",
          "avatar": "https://cdn.voobize.com/avatars/user-123.jpg"
        },
        "createdAt": "2024-01-01T11:00:00Z"
      }
    ],
    "meta": {
      "hasMore": false,
      "nextCursor": null
    }
  }
}
```

---

### 3.4 Check Block Status

**Endpoint**: `GET /chat/blocks/status/:username`

**Description**: ตรวจสอบว่าบล็อกผู้ใช้นี้อยู่หรือไม่ (หรือถูกบล็อกโดยผู้ใช้นี้)

**Path Parameters**:
```typescript
{
  username: string;
}
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Block status retrieved successfully",
  "data": {
    "isBlocked": true,          // คุณบล็อกผู้ใช้นี้
    "isBlockedBy": false,       // ผู้ใช้นี้บล็อกคุณ
    "canSendMessage": false     // สามารถส่งข้อความได้หรือไม่
  }
}
```

---

## 4. Online Status API

### 4.1 Get Online Status

**Endpoint**: `GET /chat/online-status`

**Description**: ดึงสถานะออนไลน์ของผู้ใช้หลายคน

**Query Parameters**:
```typescript
{
  userIds: string;          // Comma-separated user IDs
}
```

**Request Example**:
```bash
GET /chat/online-status?userIds=user-123,user-456,user-789
```

**Response Success (200)**:
```json
{
  "success": true,
  "message": "Online status retrieved successfully",
  "data": {
    "statuses": [
      {
        "userId": "user-123",
        "isOnline": true,
        "lastSeen": "2024-01-01T11:00:00Z"
      },
      {
        "userId": "user-456",
        "isOnline": false,
        "lastSeen": "2024-01-01T10:30:00Z"
      },
      {
        "userId": "user-789",
        "isOnline": true,
        "lastSeen": "2024-01-01T11:01:00Z"
      }
    ]
  }
}
```

---

## Error Codes Summary

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | ไม่มี token หรือ token ไม่ถูกต้อง |
| `FORBIDDEN` | 403 | ไม่มีสิทธิ์เข้าถึง resource |
| `USER_NOT_FOUND` | 404 | ไม่พบผู้ใช้ |
| `CONVERSATION_NOT_FOUND` | 404 | ไม่พบการสนทนา |
| `MESSAGE_NOT_FOUND` | 404 | ไม่พบข้อความ |
| `BLOCKED` | 403 | ถูกบล็อกหรือบล็อกผู้ใช้ |
| `ALREADY_BLOCKED` | 400 | บล็อกอยู่แล้ว |
| `BLOCK_NOT_FOUND` | 404 | ไม่พบการบล็อก |
| `RATE_LIMIT_EXCEEDED` | 429 | ส่งข้อความมากเกินไป |
| `VALIDATION_ERROR` | 400 | ข้อมูล input ไม่ถูกต้อง |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Rate Limiting

### Per User Limits
- **Send Message**: 30 ข้อความต่อนาที
- **Create Conversation**: 10 ครั้งต่อนาที
- **Mark as Read**: 60 ครั้งต่อนาที
- **Get Conversations**: 60 ครั้งต่อนาที
- **Get Messages**: 120 ครั้งต่อนาที

### Response Headers
```
X-RateLimit-Limit: 30
X-RateLimit-Remaining: 25
X-RateLimit-Reset: 1704067800
```

---

## CORS Configuration

```
Access-Control-Allow-Origin: https://voobize.com
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Authorization, Content-Type
Access-Control-Max-Age: 86400
```

---

## API Versioning Strategy

- **Current**: v1 (stable)
- **Breaking changes**: สร้าง v2
- **Deprecation**: แจ้งเตือน 3 เดือนก่อน sunset
- **Version header**: `Accept: application/vnd.voobize.v1+json` (optional)
