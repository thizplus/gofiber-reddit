# Notifications API

## Overview
Endpoints for managing user notifications.

**Base URL:** `/api/notifications`

---

## Endpoints

### 1. Get All Notifications (Private)

Retrieves all notifications for the authenticated user.

**Endpoint:** `GET /api/notifications`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20, max: 100)
filter: 'all' | 'unread' (default: 'all')
type: 'reply' | 'vote' | 'mention' | 'follow' (optional, filter by type)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "notifications": [
      {
        "id": "n1",
        "type": "reply",
        "sender": {
          "id": "u2",
          "username": "cookingmaster",
          "displayName": "เชฟมือทอง",
          "avatar": null
        },
        "postId": "1",
        "commentId": "c123",
        "message": "ตอบกลับคอมเมนต์ของคุณในโพสต์ \"สวัสดีครับ! ขอต้อนรับสู่ชุมชนของเรา\"",
        "isRead": false,
        "createdAt": "2025-01-11T14:30:00Z"
      },
      {
        "id": "n2",
        "type": "vote",
        "sender": {
          "id": "u3",
          "username": "traveler2024",
          "displayName": "นักเดินทาง",
          "avatar": null
        },
        "postId": "7",
        "commentId": null,
        "message": "โหวตโพสต์ของคุณ \"สูตรนี้ดีจริงครับ ทำตามแล้วอร่อยมาก!\"",
        "isRead": false,
        "createdAt": "2025-01-11T13:15:00Z"
      },
      {
        "id": "n3",
        "type": "mention",
        "sender": {
          "id": "u4",
          "username": "devjourney",
          "displayName": "นักพัฒนามือใหม่",
          "avatar": null
        },
        "postId": "4",
        "commentId": "c456",
        "message": "กล่าวถึงคุณในคอมเมนต์",
        "isRead": true,
        "createdAt": "2025-01-11T10:00:00Z"
      },
      {
        "id": "n4",
        "type": "follow",
        "sender": {
          "id": "u5",
          "username": "techseeker",
          "displayName": "คนหาของ",
          "avatar": null
        },
        "postId": null,
        "commentId": null,
        "message": "เริ่มติดตามคุณ",
        "isRead": true,
        "createdAt": "2025-01-10T18:45:00Z"
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 1,
      "totalNotifications": 6,
      "unreadCount": 2,
      "hasNextPage": false,
      "hasPrevPage": false
    }
  }
}
```

---

### 2. Get Unread Count (Private)

Gets the count of unread notifications.

**Endpoint:** `GET /api/notifications/unread-count`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "unreadCount": 2
  }
}
```

---

### 3. Mark Notification as Read (Private)

Marks a specific notification as read.

**Endpoint:** `PUT /api/notifications/:id/read`

**Access:** Private (Requires Authentication + Ownership)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ทำเครื่องหมายว่าอ่านแล้ว",
  "data": {
    "notification": {
      "id": "n1",
      "type": "reply",
      "sender": {
        "id": "u2",
        "username": "cookingmaster",
        "displayName": "เชฟมือทอง",
        "avatar": null
      },
      "postId": "1",
      "commentId": "c123",
      "message": "ตอบกลับคอมเมนต์ของคุณในโพสต์ \"สวัสดีครับ! ขอต้อนรับสู่ชุมชนของเรา\"",
      "isRead": true,
      "createdAt": "2025-01-11T14:30:00Z"
    }
  }
}
```

**Error Responses:**

403 Forbidden:
```json
{
  "success": false,
  "message": "ไม่มีสิทธิ์เข้าถึงการแจ้งเตือนนี้"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบการแจ้งเตือน"
}
```

---

### 4. Mark All as Read (Private)

Marks all notifications as read for the authenticated user.

**Endpoint:** `PUT /api/notifications/read-all`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ทำเครื่องหมายทั้งหมดว่าอ่านแล้ว",
  "data": {
    "markedCount": 5
  }
}
```

---

### 5. Delete Notification (Private)

Deletes a specific notification.

**Endpoint:** `DELETE /api/notifications/:id`

**Access:** Private (Requires Authentication + Ownership)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ลบการแจ้งเตือนสำเร็จ"
}
```

**Error Responses:**

403 Forbidden:
```json
{
  "success": false,
  "message": "ไม่มีสิทธิ์ลบการแจ้งเตือนนี้"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบการแจ้งเตือน"
}
```

---

### 6. Delete All Read Notifications (Private)

Deletes all read notifications for the authenticated user.

**Endpoint:** `DELETE /api/notifications/read`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ลบการแจ้งเตือนที่อ่านแล้วทั้งหมดสำเร็จ",
  "data": {
    "deletedCount": 12
  }
}
```

---

### 7. Get Notification Settings (Private)

Gets the user's notification preferences.

**Endpoint:** `GET /api/notifications/settings`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "settings": {
      "replies": true,
      "mentions": true,
      "votes": false,
      "follows": true,
      "emailNotifications": false
    }
  }
}
```

---

### 8. Update Notification Settings (Private)

Updates the user's notification preferences.

**Endpoint:** `PUT /api/notifications/settings`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "replies": "boolean (optional)",
  "mentions": "boolean (optional)",
  "votes": "boolean (optional)",
  "follows": "boolean (optional)",
  "emailNotifications": "boolean (optional)"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "อัปเดตการตั้งค่าการแจ้งเตือนสำเร็จ",
  "data": {
    "settings": {
      "replies": true,
      "mentions": true,
      "votes": false,
      "follows": true,
      "emailNotifications": false
    }
  }
}
```

---

## Notification Types

### 1. Reply (reply)
Triggered when someone replies to your comment or post.

**Format:**
```
"ตอบกลับคอมเมนต์ของคุณในโพสต์ \"{postTitle}\""
```

**Fields:**
- `postId`: ID of the post
- `commentId`: ID of the reply comment

---

### 2. Vote (vote)
Triggered when someone votes on your post or comment.

**Format:**
```
"โหวตโพสต์ของคุณ \"{postTitle}\""
"โหวตคอมเมนต์ของคุณ"
```

**Fields:**
- `postId`: ID of the post
- `commentId`: ID of the comment (if comment vote)

**Note:** Vote notifications can be disabled by user preference.

---

### 3. Mention (mention)
Triggered when someone mentions you (@username) in a post or comment.

**Format:**
```
"กล่าวถึงคุณในโพสต์ \"{postTitle}\""
"กล่าวถึงคุณในคอมเมนต์"
```

**Fields:**
- `postId`: ID of the post
- `commentId`: ID of the comment (if mentioned in comment)

---

### 4. Follow (follow)
Triggered when someone follows you.

**Format:**
```
"เริ่มติดตามคุณ"
```

**Fields:** None (only sender info)

---

## Business Logic

### Notification Creation Rules
1. Don't create notification if:
   - User performs action on their own content
   - User has disabled that notification type
   - Duplicate notification within 5 minutes

2. Group similar notifications:
   - Multiple votes on same post within 1 hour → single notification
   - Show "X people voted on your post"

3. Auto-delete old notifications:
   - Delete read notifications older than 30 days
   - Delete unread notifications older than 90 days

### Real-time Updates
- Use WebSocket/SSE for real-time notification delivery
- Push new notifications immediately when created
- Update unread count in real-time

### Notification Priority
- High: mentions, replies
- Medium: follows
- Low: votes (if enabled)
