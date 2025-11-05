# Saved Posts API

## Overview
Endpoints for saving and managing saved posts.

**Base URL:** `/api/saved`

---

## Endpoints

### 1. Get Saved Posts (Private)

Retrieves all posts saved by the authenticated user.

**Endpoint:** `GET /api/saved`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20, max: 100)
sortBy: 'new' | 'old' (default: 'new', refers to save date)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "posts": [
      {
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
        "userVote": "up",
        "commentCount": 89,
        "createdAt": "2025-01-10T10:15:00Z",
        "updatedAt": null,
        "media": [
          {
            "id": "m1",
            "type": "image",
            "url": "https://images.unsplash.com/photo-1559314809-0d155014e29e?w=800",
            "thumbnail": "https://images.unsplash.com/photo-1559314809-0d155014e29e?w=200"
          }
        ],
        "tags": ["อาหาร", "สูตรอาหาร", "ผัดไทย"],
        "savedAt": "2025-01-11T15:30:00Z"
      },
      {
        "id": "4",
        "title": "แชร์ประสบการณ์เรียน Full Stack Development แบบ Self-taught",
        "content": "สวัสดีครับทุกคน วันนี้อยากมาแชร์ประสบการณ์...",
        "author": {
          "id": "u4",
          "username": "devjourney",
          "displayName": "นักพัฒนามือใหม่",
          "avatar": null
        },
        "votes": 892,
        "userVote": "up",
        "commentCount": 156,
        "createdAt": "2025-01-09T09:00:00Z",
        "updatedAt": null,
        "media": null,
        "tags": ["programming", "career", "self-taught"],
        "savedAt": "2025-01-10T08:20:00Z"
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 2,
      "totalSaved": 38,
      "hasNextPage": true,
      "hasPrevPage": false
    }
  }
}
```

---

### 2. Check if Post is Saved (Private)

Checks if a specific post is saved by the authenticated user.

**Endpoint:** `GET /api/saved/:postId`

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
    "isSaved": true,
    "savedAt": "2025-01-11T15:30:00Z"
  }
}
```

**If Not Saved:**
```json
{
  "success": true,
  "data": {
    "isSaved": false,
    "savedAt": null
  }
}
```

---

### 3. Save Post (Private)

Saves a post for the authenticated user.

**Endpoint:** `POST /api/saved/:postId`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "บันทึกโพสต์สำเร็จ",
  "data": {
    "postId": "2",
    "savedAt": "2025-01-12T10:00:00Z"
  }
}
```

**Error Responses:**

400 Bad Request (Already Saved):
```json
{
  "success": false,
  "message": "คุณได้บันทึกโพสต์นี้แล้ว"
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

### 4. Unsave Post (Private)

Removes a post from saved posts.

**Endpoint:** `DELETE /api/saved/:postId`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ยกเลิกการบันทึกโพสต์สำเร็จ"
}
```

**Error Responses:**

400 Bad Request (Not Saved):
```json
{
  "success": false,
  "message": "คุณไม่ได้บันทึกโพสต์นี้"
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

### 5. Clear All Saved Posts (Private)

Removes all saved posts for the authenticated user.

**Endpoint:** `DELETE /api/saved`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ลบโพสต์ที่บันทึกทั้งหมดสำเร็จ",
  "data": {
    "deletedCount": 38
  }
}
```

---

### 6. Get Saved Posts Count (Private)

Gets the total count of saved posts for the authenticated user.

**Endpoint:** `GET /api/saved/count`

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
    "count": 38
  }
}
```

---

## Business Logic

### Saving Posts
- Users can save unlimited posts
- Can't save the same post twice
- Saving a post doesn't notify the post author
- Saved posts are private (not visible to others)

### Saved Post Data
- Store: `userId`, `postId`, `savedAt` timestamp
- Don't duplicate post data, just reference
- If post is deleted, keep saved reference but mark as deleted

### Sorting
- Default sort: Most recently saved first
- Can sort by: save date (new/old)
- Post's own creation date doesn't affect saved list order

### Saved Post Display
- Show deleted posts in saved list with `[โพสต์ถูกลบ]`
- Allow unsaving deleted posts
- Don't count deleted posts in statistics

### Use Cases
- Bookmark interesting posts to read later
- Create personal collection of helpful resources
- Reference posts for future discussion
