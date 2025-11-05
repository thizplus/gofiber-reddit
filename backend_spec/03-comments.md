# Comments API

## Overview
Endpoints for creating, reading, updating, deleting, and voting on comments and replies.

**Base URL:** `/api/comments`

---

## Endpoints

### 1. Get Comments for Post (Public)

Retrieves all comments for a specific post with nested replies.

**Endpoint:** `GET /api/posts/:postId/comments`

**Access:** Public (userVote only available if authenticated)

**Query Parameters:**
```
sortBy: 'new' | 'top' | 'old' (default: 'top')
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "comments": [
      {
        "id": "c1",
        "postId": "1",
        "author": {
          "id": "u2",
          "username": "cookingmaster",
          "displayName": "เชฟมือทอง",
          "avatar": null
        },
        "content": "ยินดีต้อนรับครับ! หวังว่าจะได้แลกเปลี่ยนเรื่องราวกันเยอะๆ",
        "votes": 12,
        "userVote": null,
        "parentId": null,
        "depth": 0,
        "createdAt": "2025-01-10T09:00:00Z",
        "updatedAt": null,
        "replies": []
      },
      {
        "id": "c2",
        "postId": "1",
        "author": {
          "id": "u3",
          "username": "traveler2024",
          "displayName": "นักเดินทาง",
          "avatar": null
        },
        "content": "สวัสดีครับผม! ชอบบรรยากาศที่นี่มากเลย",
        "votes": 8,
        "userVote": "up",
        "parentId": null,
        "depth": 0,
        "createdAt": "2025-01-10T10:30:00Z",
        "updatedAt": null,
        "replies": [
          {
            "id": "c3",
            "postId": "1",
            "author": {
              "id": "u1",
              "username": "thepthai",
              "displayName": "เทพไท ใจน้อม",
              "avatar": "/logo.png"
            },
            "content": "ขอบคุณครับ! ยินดีที่ได้รู้จัก 😊",
            "votes": 5,
            "userVote": null,
            "parentId": "c2",
            "depth": 1,
            "createdAt": "2025-01-10T11:00:00Z",
            "updatedAt": null,
            "replies": [
              {
                "id": "c3-1",
                "postId": "1",
                "author": {
                  "id": "u4",
                  "username": "devjourney",
                  "displayName": "นักพัฒนามือใหม่",
                  "avatar": null
                },
                "content": "ผมก็เพิ่งมาเหมือนกันครับ ยินดีที่ได้รู้จักทุกคน",
                "votes": 3,
                "userVote": null,
                "parentId": "c3",
                "depth": 2,
                "createdAt": "2025-01-10T11:30:00Z",
                "updatedAt": null,
                "replies": []
              }
            ]
          }
        ]
      }
    ],
    "totalComments": 23
  }
}
```

**Error Responses:**

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบโพสต์"
}
```

---

### 2. Create Comment (Private)

Creates a new comment or reply.

**Endpoint:** `POST /api/posts/:postId/comments`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "content": "string (required, max 5000 characters)",
  "parentId": "string (optional, for replies)"
}
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "สร้างความคิดเห็นสำเร็จ",
  "data": {
    "comment": {
      "id": "c100",
      "postId": "1",
      "author": {
        "id": "u1",
        "username": "thepthai",
        "displayName": "เทพไท ใจน้อม",
        "avatar": "/logo.png"
      },
      "content": "This is my comment!",
      "votes": 0,
      "userVote": null,
      "parentId": null,
      "depth": 0,
      "createdAt": "2025-01-12T10:00:00Z",
      "updatedAt": null,
      "replies": []
    }
  }
}
```

**Error Responses:**

400 Bad Request:
```json
{
  "success": false,
  "message": "ข้อมูลไม่ถูกต้อง",
  "errors": {
    "content": "กรุณากรอกเนื้อหาความคิดเห็น"
  }
}
```

404 Not Found (Post):
```json
{
  "success": false,
  "message": "ไม่พบโพสต์"
}
```

404 Not Found (Parent Comment):
```json
{
  "success": false,
  "message": "ไม่พบความคิดเห็นหลัก"
}
```

---

### 3. Update Comment (Private)

Updates an existing comment (only author can update).

**Endpoint:** `PUT /api/comments/:id`

**Access:** Private (Requires Authentication + Ownership)

**Headers:**
```
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "content": "string (required, max 5000 characters)"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "แก้ไขความคิดเห็นสำเร็จ",
  "data": {
    "comment": {
      "id": "c100",
      "postId": "1",
      "author": {
        "id": "u1",
        "username": "thepthai",
        "displayName": "เทพไท ใจน้อม",
        "avatar": "/logo.png"
      },
      "content": "Updated comment content",
      "votes": 5,
      "userVote": "up",
      "parentId": null,
      "depth": 0,
      "createdAt": "2025-01-12T10:00:00Z",
      "updatedAt": "2025-01-12T10:30:00Z",
      "replies": []
    }
  }
}
```

**Error Responses:**

403 Forbidden:
```json
{
  "success": false,
  "message": "คุณไม่มีสิทธิ์แก้ไขความคิดเห็นนี้"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบความคิดเห็น"
}
```

---

### 4. Delete Comment (Private)

Deletes a comment (only author can delete).

**Endpoint:** `DELETE /api/comments/:id`

**Access:** Private (Requires Authentication + Ownership)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ลบความคิดเห็นสำเร็จ"
}
```

**Error Responses:**

403 Forbidden:
```json
{
  "success": false,
  "message": "คุณไม่มีสิทธิ์ลบความคิดเห็นนี้"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบความคิดเห็น"
}
```

---

### 5. Vote Comment (Private)

Upvote or downvote a comment.

**Endpoint:** `POST /api/comments/:id/vote`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "voteType": "up" | "down" | null
}
```

**Note:**
- `"up"` = upvote
- `"down"` = downvote
- `null` = remove vote

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "โหวตสำเร็จ",
  "data": {
    "votes": 13,
    "userVote": "up"
  }
}
```

**Error Responses:**

400 Bad Request:
```json
{
  "success": false,
  "message": "ประเภทการโหวตไม่ถูกต้อง"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบความคิดเห็น"
}
```

---

### 6. Get Comment by ID (Public)

Retrieves a single comment with its replies.

**Endpoint:** `GET /api/comments/:id`

**Access:** Public

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "comment": {
      "id": "c2",
      "postId": "1",
      "author": {
        "id": "u3",
        "username": "traveler2024",
        "displayName": "นักเดินทาง",
        "avatar": null
      },
      "content": "สวัสดีครับผม! ชอบบรรยากาศที่นี่มากเลย",
      "votes": 8,
      "userVote": "up",
      "parentId": null,
      "depth": 0,
      "createdAt": "2025-01-10T10:30:00Z",
      "updatedAt": null,
      "replies": [...]
    }
  }
}
```

**Error Responses:**

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบความคิดเห็น"
}
```

---

## Business Logic

### Comment Nesting
- Maximum depth: 10 levels (configurable)
- `depth` field indicates nesting level (0 = top-level)
- Replies are sorted by votes (top) or time (new)

### Vote Calculation
- Each upvote adds +1 to votes
- Each downvote adds -1 to votes
- User can change vote
- Comment author's karma affected by votes received

### Comment Deletion
- Soft delete: Set `isDeleted: true`, keep in database
- Content shown as "[คอมเมนต์ถูกลบ]"
- Author shown as "[deleted]"
- Nested replies remain visible
- Can't reply to deleted comments

### Notifications
- Reply to comment → notify parent comment author
- Vote on comment → notify comment author (optional, can be disabled)
- Mention user (@username) → notify mentioned user

### Post Comment Count
- Automatically update post's `commentCount` when:
  - Comment created: +1
  - Comment deleted: -1
  - Nested replies also counted
