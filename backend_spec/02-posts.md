# Posts API

## Overview
Endpoints for creating, reading, updating, deleting, and interacting with posts.

**Base URL:** `/api/posts`

---

## Endpoints

### 1. Get All Posts (Public)

Retrieves a paginated list of posts.

**Endpoint:** `GET /api/posts`

**Access:** Public (userVote only available if authenticated)

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20, max: 100)
sortBy: 'new' | 'hot' | 'top' (default: 'hot')
timeRange: 'today' | 'week' | 'month' | 'year' | 'all' (default: 'all', used with sortBy='top')
tag: string (optional, filter by tag)
author: string (optional, filter by username)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "posts": [
      {
        "id": "1",
        "title": "สวัสดีครับ! ขอต้อนรับสู่ชุมชนของเรา",
        "content": "นี่คือโพสต์แรกของผม หวังว่าทุกคนจะชอบนะครับ...",
        "author": {
          "id": "u1",
          "username": "thepthai",
          "displayName": "เทพไท ใจน้อม",
          "avatar": "/logo.png"
        },
        "votes": 142,
        "userVote": null,
        "commentCount": 23,
        "createdAt": "2025-01-10T08:30:00Z",
        "updatedAt": null,
        "media": null,
        "tags": ["แนะนำตัว", "สวัสดี"],
        "sourcePostId": null,
        "sourcePost": null
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 5,
      "totalPosts": 95,
      "hasNextPage": true,
      "hasPrevPage": false
    }
  }
}
```

---

### 2. Get Post by ID (Public)

Retrieves a single post by its ID.

**Endpoint:** `GET /api/posts/:id`

**Access:** Public (userVote only available if authenticated)

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "post": {
      "id": "7",
      "title": "สูตรนี้ดีจริงครับ ทำตามแล้วอร่อยมาก!",
      "content": "เห็นสูตรผัดไทยของคุณเชฟมือทอง ลองทำตามแล้วครับ...",
      "author": {
        "id": "u1",
        "username": "thepthai",
        "displayName": "เทพไท ใจน้อม",
        "avatar": "/logo.png"
      },
      "votes": 89,
      "userVote": "up",
      "commentCount": 12,
      "createdAt": "2025-01-11T12:00:00Z",
      "updatedAt": null,
      "media": null,
      "tags": ["อาหาร", "ผัดไทย"],
      "sourcePostId": "2",
      "sourcePost": {
        "id": "2",
        "title": "วิธีทำ Pad Thai ที่บ้านง่ายๆ แบบไทยแท้",
        "content": "วันนี้จะมาแชร์วิธีทำผัดไทยที่บ้านกันครับ...",
        "author": {
          "id": "u2",
          "username": "cookingmaster",
          "displayName": "เชฟมือทอง",
          "avatar": null
        },
        "votes": 567,
        "commentCount": 89,
        "createdAt": "2025-01-10T10:15:00Z",
        "media": [
          {
            "id": "m1",
            "type": "image",
            "url": "https://images.unsplash.com/photo-1559314809-0d155014e29e?w=800",
            "thumbnail": "https://images.unsplash.com/photo-1559314809-0d155014e29e?w=200"
          }
        ],
        "tags": ["อาหาร", "สูตรอาหาร", "ผัดไทย"]
      }
    }
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

### 3. Create Post (Private)

Creates a new post.

**Endpoint:** `POST /api/posts`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
Content-Type: multipart/form-data (if uploading media)
```

**Request Body (FormData):**
```json
{
  "title": "string (required, max 300 characters)",
  "content": "string (required, max 10000 characters)",
  "tags": "string[] (optional, max 5 tags, each max 30 characters)",
  "sourcePostId": "string (optional, for crosspost)",
  "media": "File[] (optional, max 10 files, each max 50MB)"
}
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "สร้างโพสต์สำเร็จ",
  "data": {
    "post": {
      "id": "8",
      "title": "My New Post",
      "content": "This is the content...",
      "author": {
        "id": "u1",
        "username": "thepthai",
        "displayName": "เทพไท ใจน้อม",
        "avatar": "/logo.png"
      },
      "votes": 0,
      "userVote": null,
      "commentCount": 0,
      "createdAt": "2025-01-12T10:00:00Z",
      "updatedAt": null,
      "media": [
        {
          "id": "m3",
          "type": "image",
          "url": "https://storage.example.com/uploads/abc123.jpg",
          "thumbnail": "https://storage.example.com/uploads/abc123_thumb.jpg"
        }
      ],
      "tags": ["technology", "programming"],
      "sourcePostId": null,
      "sourcePost": null
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
    "title": "กรุณากรอกหัวข้อ",
    "content": "กรุณากรอกเนื้อหา",
    "tags": "แท็กต้องไม่เกิน 5 แท็ก"
  }
}
```

401 Unauthorized:
```json
{
  "success": false,
  "message": "ไม่ได้รับอนุญาต กรุณาเข้าสู่ระบบ"
}
```

---

### 4. Update Post (Private)

Updates an existing post (only author can update).

**Endpoint:** `PUT /api/posts/:id`

**Access:** Private (Requires Authentication + Ownership)

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "string (optional, max 300 characters)",
  "content": "string (optional, max 10000 characters)",
  "tags": "string[] (optional, max 5 tags)"
}
```

**Note:** Cannot update media or sourcePost after creation.

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "แก้ไขโพสต์สำเร็จ",
  "data": {
    "post": {
      "id": "8",
      "title": "Updated Title",
      "content": "Updated content...",
      "author": {
        "id": "u1",
        "username": "thepthai",
        "displayName": "เทพไท ใจน้อม",
        "avatar": "/logo.png"
      },
      "votes": 5,
      "userVote": "up",
      "commentCount": 2,
      "createdAt": "2025-01-12T10:00:00Z",
      "updatedAt": "2025-01-12T11:30:00Z",
      "media": [...],
      "tags": ["updated", "tags"],
      "sourcePostId": null,
      "sourcePost": null
    }
  }
}
```

**Error Responses:**

403 Forbidden:
```json
{
  "success": false,
  "message": "คุณไม่มีสิทธิ์แก้ไขโพสต์นี้"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบโพสต์"
}
```

---

### 5. Delete Post (Private)

Deletes a post (only author can delete).

**Endpoint:** `DELETE /api/posts/:id`

**Access:** Private (Requires Authentication + Ownership)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ลบโพสต์สำเร็จ"
}
```

**Error Responses:**

403 Forbidden:
```json
{
  "success": false,
  "message": "คุณไม่มีสิทธิ์ลบโพสต์นี้"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบโพสต์"
}
```

---

### 6. Vote Post (Private)

Upvote or downvote a post.

**Endpoint:** `POST /api/posts/:id/vote`

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
    "votes": 143,
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
  "message": "ไม่พบโพสต์"
}
```

---

### 7. Get User's Posts (Public)

Gets all posts by a specific user.

**Endpoint:** `GET /api/posts/user/:username`

**Access:** Public

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20)
sortBy: 'new' | 'top' (default: 'new')
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "posts": [...],
    "pagination": {
      "currentPage": 1,
      "totalPages": 3,
      "totalPosts": 45,
      "hasNextPage": true,
      "hasPrevPage": false
    }
  }
}
```

---

### 8. Get My Posts (Private)

Gets all posts by the authenticated user.

**Endpoint:** `GET /api/posts/me`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20)
sortBy: 'new' | 'top' (default: 'new')
includeDeleted: boolean (default: false)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "posts": [...],
    "pagination": {
      "currentPage": 1,
      "totalPages": 2,
      "totalPosts": 28,
      "hasNextPage": true,
      "hasPrevPage": false
    }
  }
}
```

---

## Business Logic

### Vote Calculation
- Each upvote adds +1 to votes
- Each downvote adds -1 to votes
- Removing vote returns to previous state
- User can change vote (up → down or vice versa)
- User's karma is affected by post votes received

### Hot Score Algorithm
```
hot_score = (votes / (hours_since_posted + 2)^1.5)
```

### Post Deletion
- Soft delete: Set `isDeleted: true`, keep in database
- Author name shown as "[deleted]" for deleted posts
- Content shown as "[โพสต์ถูกลบ]"
- Comments remain visible but can't add new comments
