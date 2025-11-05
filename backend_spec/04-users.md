# Users API

## Overview
Endpoints for user profiles, following, and user management.

**Base URL:** `/api/users`

---

## Endpoints

### 1. Get User Profile (Public)

Retrieves a user's public profile information.

**Endpoint:** `GET /api/users/:username`

**Access:** Public

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "u1",
      "username": "thepthai",
      "displayName": "เทพไท ใจน้อม",
      "avatar": "/logo.png",
      "karma": 1247,
      "bio": "นักพัฒนาซอฟต์แวร์ที่หลงรักการเขียนโค้ดและแบ่งปันความรู้ 💻✨",
      "coverImage": null,
      "joinedAt": "2024-06-15T10:00:00Z",
      "location": "กรุงเทพมหานคร, ประเทศไทย",
      "website": "https://thepthai.dev",
      "followersCount": 432,
      "followingCount": 156,
      "isFollowing": false
    }
  }
}
```

**Note:** `isFollowing` is only available if the requester is authenticated.

**Error Responses:**

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบผู้ใช้"
}
```

---

### 2. Update Profile (Private)

Updates the authenticated user's profile.

**Endpoint:** `PUT /api/users/me`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
Content-Type: multipart/form-data (if uploading avatar)
```

**Request Body (FormData):**
```json
{
  "displayName": "string (optional, max 50 characters)",
  "bio": "string (optional, max 200 characters)",
  "location": "string (optional, max 100 characters)",
  "website": "string (optional, valid URL)",
  "avatar": "File (optional, image file, max 5MB)"
}
```

**Note:** Cannot change username or email through this endpoint.

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "อัปเดตโปรไฟล์สำเร็จ",
  "data": {
    "user": {
      "id": "u1",
      "username": "thepthai",
      "email": "thepthai@example.com",
      "displayName": "เทพไท ใจน้อม (Updated)",
      "avatar": "https://storage.example.com/avatars/u1_abc123.jpg",
      "karma": 1247,
      "bio": "Updated bio",
      "coverImage": null,
      "joinedAt": "2024-06-15T10:00:00Z",
      "location": "เชียงใหม่, ประเทศไทย",
      "website": "https://newwebsite.dev",
      "followersCount": 432,
      "followingCount": 156
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
    "displayName": "ชื่อที่แสดงต้องไม่เกิน 50 ตัวอักษร",
    "website": "รูปแบบ URL ไม่ถูกต้อง",
    "avatar": "ไฟล์ต้องเป็นรูปภาพและขนาดไม่เกิน 5MB"
  }
}
```

---

### 3. Follow User (Private)

Follows a user.

**Endpoint:** `POST /api/users/:username/follow`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ติดตามผู้ใช้สำเร็จ",
  "data": {
    "isFollowing": true,
    "followersCount": 433
  }
}
```

**Error Responses:**

400 Bad Request (Already Following):
```json
{
  "success": false,
  "message": "คุณติดตามผู้ใช้นี้อยู่แล้ว"
}
```

400 Bad Request (Self Follow):
```json
{
  "success": false,
  "message": "ไม่สามารถติดตามตัวเองได้"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบผู้ใช้"
}
```

---

### 4. Unfollow User (Private)

Unfollows a user.

**Endpoint:** `DELETE /api/users/:username/follow`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "เลิกติดตามผู้ใช้สำเร็จ",
  "data": {
    "isFollowing": false,
    "followersCount": 432
  }
}
```

**Error Responses:**

400 Bad Request (Not Following):
```json
{
  "success": false,
  "message": "คุณไม่ได้ติดตามผู้ใช้นี้อยู่"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบผู้ใช้"
}
```

---

### 5. Get User's Followers (Public)

Gets a list of users following the specified user.

**Endpoint:** `GET /api/users/:username/followers`

**Access:** Public

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20, max: 100)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "followers": [
      {
        "id": "u2",
        "username": "cookingmaster",
        "displayName": "เชฟมือทอง",
        "avatar": null,
        "karma": 2891,
        "bio": "รักการทำอาหาร มาแชร์สูตรอาหารไทยและนานาชาติกัน 🍳👨‍🍳",
        "isFollowing": true,
        "followedAt": "2024-12-01T10:00:00Z"
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 22,
      "totalFollowers": 432,
      "hasNextPage": true,
      "hasPrevPage": false
    }
  }
}
```

---

### 6. Get User's Following (Public)

Gets a list of users that the specified user is following.

**Endpoint:** `GET /api/users/:username/following`

**Access:** Public

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20, max: 100)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "following": [
      {
        "id": "u3",
        "username": "traveler2024",
        "displayName": "นักเดินทาง",
        "avatar": null,
        "karma": 567,
        "bio": "ชอบเที่ยว ชอบถ่ายรูป แชร์ประสบการณ์ท่องเที่ยวทั่วไทย 🌏📸",
        "isFollowing": true,
        "followedAt": "2024-11-15T08:00:00Z"
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 8,
      "totalFollowing": 156,
      "hasNextPage": true,
      "hasPrevPage": false
    }
  }
}
```

---

### 7. Get User's Karma History (Private)

Gets the authenticated user's karma history.

**Endpoint:** `GET /api/users/me/karma`

**Access:** Private (Requires Authentication, own profile only)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 50)
timeRange: 'day' | 'week' | 'month' | 'all' (default: 'all')
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "currentKarma": 1247,
    "history": [
      {
        "id": "k1",
        "type": "post_vote",
        "change": 1,
        "postId": "7",
        "postTitle": "สูตรนี้ดีจริงครับ ทำตามแล้วอร่อยมาก!",
        "createdAt": "2025-01-12T09:30:00Z"
      },
      {
        "id": "k2",
        "type": "comment_vote",
        "change": 1,
        "commentId": "c10",
        "postId": "6",
        "createdAt": "2025-01-12T08:15:00Z"
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 25,
      "totalRecords": 1247,
      "hasNextPage": true,
      "hasPrevPage": false
    }
  }
}
```

---

### 8. Get User's Comments (Public)

Gets all comments made by a specific user.

**Endpoint:** `GET /api/users/:username/comments`

**Access:** Public

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20, max: 100)
sortBy: 'new' | 'top' | 'old' (default: 'new')
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "comments": [
      {
        "id": "c10",
        "postId": "6",
        "post": {
          "id": "6",
          "title": "พาชมวิวทะเลกระบี่สวยๆ จากมุมสูง 🌊",
          "author": {
            "id": "u3",
            "username": "traveler2024",
            "displayName": "นักเดินทาง"
          }
        },
        "author": {
          "id": "u1",
          "username": "thepthai",
          "displayName": "เทพไท ใจน้อม",
          "avatar": "/logo.png"
        },
        "content": "สวยมากเลยครับ! ถ่ายด้วย drone รุ่นอะไรครับ?",
        "votes": 18,
        "userVote": "up",
        "parentId": null,
        "depth": 0,
        "createdAt": "2025-01-11T09:00:00Z",
        "updatedAt": null
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 5,
      "totalComments": 92,
      "hasNextPage": true,
      "hasPrevPage": false
    }
  }
}
```

**Note:** Each comment includes basic post info (id, title, author) so user can see context.

**Error Responses:**

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบผู้ใช้"
}
```

---

### 9. Search Users (Public)

Searches for users by username, display name, or bio.

**Endpoint:** `GET /api/users/search`

**Access:** Public

**Query Parameters:**
```
q: string (required, search query)
page: number (default: 1)
limit: number (default: 20, max: 100)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "u1",
        "username": "thepthai",
        "displayName": "เทพไท ใจน้อม",
        "avatar": "/logo.png",
        "karma": 1247,
        "bio": "นักพัฒนาซอฟต์แวร์ที่หลงรักการเขียนโค้ดและแบ่งปันความรู้ 💻✨",
        "followersCount": 432
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 1,
      "totalUsers": 5,
      "hasNextPage": false,
      "hasPrevPage": false
    }
  }
}
```

---

### 10. Delete Account (Private)

Deletes the authenticated user's account.

**Endpoint:** `DELETE /api/users/me`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "password": "string (required, for confirmation)",
  "confirmation": "DELETE" (required, must type exactly "DELETE")
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ลบบัญชีสำเร็จ"
}
```

**Error Responses:**

400 Bad Request:
```json
{
  "success": false,
  "message": "กรุณายืนยันการลบบัญชีโดยพิมพ์ DELETE"
}
```

401 Unauthorized:
```json
{
  "success": false,
  "message": "รหัสผ่านไม่ถูกต้อง"
}
```

---

## Business Logic

### Karma Calculation
Karma is the sum of all votes received on posts and comments:
- Post upvote: +1 karma
- Post downvote: -1 karma
- Comment upvote: +1 karma
- Comment downvote: -1 karma

### Follow System
- Can't follow yourself
- Follow creates notification to target user
- Unfollow doesn't create notification
- Following list affects feed algorithm (future feature)

### Profile Visibility
- All profiles are public
- Email address is private (never shown in API responses except to owner)
- Karma history is private (only visible to owner)

### Account Deletion
- Soft delete: Keep user data but mark as deleted
- Posts/comments remain but author shown as "[deleted]"
- Username becomes available after 30 days
- Can't recover account after deletion
