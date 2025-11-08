# API Endpoints

**Base URL**: `http://localhost:8080/api/v1`

**Authentication**: JWT Bearer Token (Header: `Authorization: Bearer {token}`)

**Response Format**: JSON

## Legend

- 🌐 **Public**: ไม่ต้อง authentication
- 🔓 **Optional Auth**: Login แล้วจะได้ข้อมูลเพิ่มเติม
- 🔒 **Protected**: ต้อง login
- 👑 **Admin Only**: ต้องเป็น admin

---

## 1. Authentication (`/api/v1/auth`)

### 1.1 Register (สมัครสมาชิก)
```http
POST /api/v1/auth/register
```

**Access**: 🌐 Public

**Request Body**:
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "SecurePass123"
}
```

**Validation**:
- `username`: 3-20 characters
- `email`: valid email format
- `password`: min 8 characters

**Response** (201):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "displayName": "",
    "avatar": "",
    "role": "user",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

---

### 1.2 Login (เข้าสู่ระบบ)
```http
POST /api/v1/auth/login
```

**Access**: 🌐 Public

**Request Body**:
```json
{
  "login": "johndoe",  // username or email
  "password": "SecurePass123"
}
```

**Response** (200):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { ... }
}
```

---

### 1.3 Google OAuth URL
```http
GET /api/v1/auth/google
```

**Access**: 🌐 Public

**Response** (200):
```json
{
  "url": "https://accounts.google.com/o/oauth2/auth?..."
}
```

---

### 1.4 Google OAuth Callback
```http
GET /api/v1/auth/google/callback?code={code}&state={state}
```

**Access**: 🌐 Public (called by Google)

**Response**: Redirect to frontend with `code` parameter

---

### 1.5 Exchange OAuth Code
```http
POST /api/v1/auth/exchange
```

**Access**: 🌐 Public

**Request Body**:
```json
{
  "code": "temporary_auth_code"
}
```

**Response** (200):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { ... }
}
```

---

## 2. Users (`/api/v1/users`)

### 2.1 Get Own Profile
```http
GET /api/v1/users/profile
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "username": "johndoe",
  "displayName": "John Doe",
  "avatar": "https://cdn.example.com/avatar.jpg",
  "bio": "Software Developer",
  "location": "Bangkok, Thailand",
  "website": "https://johndoe.com",
  "karma": 150,
  "followersCount": 42,
  "followingCount": 30,
  "createdAt": "2024-01-01T00:00:00Z"
}
```

---

### 2.2 Update Profile
```http
PUT /api/v1/users/profile
```

**Access**: 🔒 Protected

**Request Body**:
```json
{
  "displayName": "John Doe",
  "avatar": "https://cdn.example.com/new-avatar.jpg",
  "bio": "Full-stack Developer",
  "location": "Bangkok",
  "website": "https://johndoe.dev"
}
```

**Response** (200): Updated user object

---

### 2.3 Delete Account
```http
DELETE /api/v1/users/profile
```

**Access**: 🔒 Protected

**Response** (204): No Content

---

### 2.4 List All Users
```http
GET /api/v1/users?page=1&limit=20
```

**Access**: 👑 Admin Only

**Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20)

**Response** (200):
```json
{
  "users": [ ... ],
  "total": 150,
  "page": 1,
  "limit": 20
}
```

---

## 3. Profiles (`/api/v1/profiles`)

### 3.1 Get Public Profile
```http
GET /api/v1/profiles/:username
```

**Access**: 🔓 Optional Auth

**Response** (200):
```json
{
  "id": "uuid",
  "username": "johndoe",
  "displayName": "John Doe",
  "avatar": "https://cdn.example.com/avatar.jpg",
  "bio": "Software Developer",
  "karma": 150,
  "followersCount": 42,
  "followingCount": 30,
  "isFollowing": false,  // if authenticated
  "createdAt": "2024-01-01T00:00:00Z"
}
```

---

## 4. Posts (`/api/v1/posts`)

### 4.1 List Posts
```http
GET /api/v1/posts?sort=hot&page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Query Parameters**:
- `sort`: `hot`, `new`, `top` (default: `hot`)
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20)
- `time`: For `top` sort - `hour`, `day`, `week`, `month`, `year`, `all`

**Response** (200):
```json
{
  "posts": [
    {
      "id": "uuid",
      "title": "My First Post",
      "content": "This is the content...",
      "authorId": "uuid",
      "author": {
        "username": "johndoe",
        "avatar": "..."
      },
      "votes": 42,
      "commentCount": 15,
      "media": [],
      "tags": ["golang", "fiber"],
      "createdAt": "2024-01-01T00:00:00Z",
      "userVote": null  // if authenticated: "up", "down", or null
    }
  ],
  "total": 150,
  "page": 1,
  "limit": 20
}
```

---

### 4.2 Get Single Post
```http
GET /api/v1/posts/:id
```

**Access**: 🔓 Optional Auth

**Response** (200): Single post object (same structure as above)

---

### 4.3 Create Post
```http
POST /api/v1/posts
```

**Access**: 🔒 Protected

**Request Body**:
```json
{
  "title": "My Post Title",
  "content": "This is the content...",
  "tags": ["golang", "fiber"],
  "mediaIds": ["uuid1", "uuid2"]
}
```

**Response** (201): Created post object

---

### 4.4 Update Post
```http
PUT /api/v1/posts/:id
```

**Access**: 🔒 Protected (owner only)

**Request Body**:
```json
{
  "title": "Updated Title",
  "content": "Updated content...",
  "tags": ["golang", "updated"]
}
```

**Response** (200): Updated post object

---

### 4.5 Delete Post
```http
DELETE /api/v1/posts/:id
```

**Access**: 🔒 Protected (owner only)

**Response** (204): No Content

---

### 4.6 Create Crosspost
```http
POST /api/v1/posts/:id/crosspost
```

**Access**: 🔒 Protected

**Request Body**:
```json
{
  "title": "Check this out!",
  "tags": ["shared"]
}
```

**Response** (201): Created crosspost (references original via `sourcePostId`)

---

### 4.7 Get Crossposts
```http
GET /api/v1/posts/:id/crossposts
```

**Access**: 🔓 Optional Auth

**Response** (200): Array of crossposts

---

### 4.8 Get Posts by Author
```http
GET /api/v1/posts/author/:authorId?page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Response** (200): Paginated posts by author

---

### 4.9 Get Posts by Tag Name
```http
GET /api/v1/posts/tag/:tagName?page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Response** (200): Paginated posts with tag

---

### 4.10 Get Posts by Tag ID
```http
GET /api/v1/posts/tag-id/:tagId?page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Response** (200): Paginated posts with tag

---

### 4.11 Search Posts
```http
GET /api/v1/posts/search?q=golang&page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Query Parameters**:
- `q`: Search query (required)
- `page`, `limit`: Pagination

**Response** (200): Paginated search results

---

### 4.12 Get Personalized Feed
```http
GET /api/v1/posts/feed?page=1&limit=20
```

**Access**: 🔒 Protected

**Description**: Posts from followed users

**Response** (200): Paginated posts

---

## 5. Comments (`/api/v1/comments`)

### 5.1 Get Comment
```http
GET /api/v1/comments/:id
```

**Access**: 🔓 Optional Auth

**Response** (200):
```json
{
  "id": "uuid",
  "postId": "uuid",
  "authorId": "uuid",
  "author": {
    "username": "johndoe",
    "avatar": "..."
  },
  "content": "This is a comment",
  "votes": 10,
  "parentId": null,  // or parent comment UUID
  "depth": 0,  // 0-10
  "createdAt": "2024-01-01T00:00:00Z",
  "userVote": null  // if authenticated
}
```

---

### 5.2 Get Comments by Post
```http
GET /api/v1/comments/post/:postId?page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Response** (200): Paginated comments (flat list)

---

### 5.3 Get Comment Tree
```http
GET /api/v1/comments/post/:postId/tree
```

**Access**: 🔓 Optional Auth

**Description**: Nested comment structure

**Response** (200):
```json
{
  "comments": [
    {
      "id": "uuid",
      "content": "Top-level comment",
      "replies": [
        {
          "id": "uuid",
          "content": "Nested reply",
          "replies": [ ... ]
        }
      ]
    }
  ]
}
```

---

### 5.4 Create Comment
```http
POST /api/v1/comments
```

**Access**: 🔒 Protected

**Request Body**:
```json
{
  "postId": "uuid",
  "content": "This is my comment",
  "parentId": null  // or parent comment UUID for reply
}
```

**Response** (201): Created comment object

---

### 5.5 Update Comment
```http
PUT /api/v1/comments/:id
```

**Access**: 🔒 Protected (owner only)

**Request Body**:
```json
{
  "content": "Updated comment text"
}
```

**Response** (200): Updated comment object

---

### 5.6 Delete Comment
```http
DELETE /api/v1/comments/:id
```

**Access**: 🔒 Protected (owner only)

**Response** (204): No Content (soft delete)

---

### 5.7 Get Comments by Author
```http
GET /api/v1/comments/author/:authorId?page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Response** (200): Paginated comments

---

### 5.8 Get Direct Replies
```http
GET /api/v1/comments/:id/replies?page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Response** (200): Direct child comments

---

### 5.9 Get Parent Chain
```http
GET /api/v1/comments/:id/parent-chain
```

**Access**: 🔓 Optional Auth

**Description**: Get breadcrumb to root comment

**Response** (200): Array of parent comments

---

## 6. Votes (`/api/v1/votes`)

### 6.1 Vote on Target
```http
POST /api/v1/votes
```

**Access**: 🔒 Protected

**Request Body**:
```json
{
  "targetId": "uuid",
  "targetType": "post",  // or "comment"
  "voteType": "up"  // or "down"
}
```

**Response** (200):
```json
{
  "message": "Vote recorded",
  "newVoteCount": 43
}
```

---

### 6.2 Remove Vote
```http
DELETE /api/v1/votes/:targetType/:targetId
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "message": "Vote removed",
  "newVoteCount": 42
}
```

---

### 6.3 Get Vote Status
```http
GET /api/v1/votes/:targetType/:targetId
```

**Access**: 🌐 Public

**Response** (200):
```json
{
  "voteType": "up",  // or "down" or null
  "voteCount": 42
}
```

---

### 6.4 Get Vote Count
```http
GET /api/v1/votes/:targetType/:targetId/count
```

**Access**: 🌐 Public

**Response** (200):
```json
{
  "count": 42
}
```

---

### 6.5 Get User's Votes
```http
GET /api/v1/votes/user?page=1&limit=20
```

**Access**: 🔒 Protected

**Response** (200): Paginated user votes

---

## 7. Follows (`/api/v1/follows`)

### 7.1 Follow User
```http
POST /api/v1/follows/user/:userId
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "message": "Followed successfully"
}
```

---

### 7.2 Unfollow User
```http
DELETE /api/v1/follows/user/:userId
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "message": "Unfollowed successfully"
}
```

---

### 7.3 Get Followers
```http
GET /api/v1/follows/user/:userId/followers?page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Response** (200):
```json
{
  "followers": [
    {
      "id": "uuid",
      "username": "follower1",
      "avatar": "...",
      "isFollowing": false  // if authenticated
    }
  ],
  "total": 42,
  "page": 1
}
```

---

### 7.4 Get Following
```http
GET /api/v1/follows/user/:userId/following?page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Response** (200): Same structure as followers

---

### 7.5 Check Follow Status
```http
GET /api/v1/follows/user/:userId/status
```

**Access**: 🔓 Optional Auth

**Response** (200):
```json
{
  "isFollowing": true
}
```

---

### 7.6 Get Mutual Follows
```http
GET /api/v1/follows/mutuals?page=1&limit=20
```

**Access**: 🔒 Protected

**Description**: Users who follow you back

**Response** (200): Paginated mutual follows

---

## 8. Saved Posts (`/api/v1/saved`)

### 8.1 Save Post
```http
POST /api/v1/saved/posts/:postId
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "message": "Post saved"
}
```

---

### 8.2 Unsave Post
```http
DELETE /api/v1/saved/posts/:postId
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "message": "Post unsaved"
}
```

---

### 8.3 Check if Saved
```http
GET /api/v1/saved/posts/:postId/status
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "isSaved": true
}
```

---

### 8.4 Get Saved Posts
```http
GET /api/v1/saved/posts?page=1&limit=20
```

**Access**: 🔒 Protected

**Response** (200): Paginated saved posts

---

## 9. Notifications (`/api/v1/notifications`)

### 9.1 Get Notifications
```http
GET /api/v1/notifications?page=1&limit=20
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "notifications": [
    {
      "id": "uuid",
      "type": "reply",
      "message": "johndoe replied to your post",
      "sender": {
        "username": "johndoe",
        "avatar": "..."
      },
      "postId": "uuid",
      "commentId": "uuid",
      "isRead": false,
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 50,
  "page": 1
}
```

---

### 9.2 Get Unread Notifications
```http
GET /api/v1/notifications/unread?limit=20
```

**Access**: 🔒 Protected

**Response** (200): Array of unread notifications

---

### 9.3 Get Unread Count
```http
GET /api/v1/notifications/unread/count
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "count": 5
}
```

---

### 9.4 Mark as Read
```http
PUT /api/v1/notifications/:id/read
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "message": "Marked as read"
}
```

---

### 9.5 Mark All as Read
```http
PUT /api/v1/notifications/read-all
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "message": "All marked as read"
}
```

---

### 9.6 Delete Notification
```http
DELETE /api/v1/notifications/:id
```

**Access**: 🔒 Protected

**Response** (204): No Content

---

### 9.7 Delete All Notifications
```http
DELETE /api/v1/notifications
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "message": "All notifications deleted"
}
```

---

### 9.8 Get Notification Settings
```http
GET /api/v1/notifications/settings
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "replies": true,
  "mentions": true,
  "votes": true,
  "follows": true,
  "emailNotifications": false
}
```

---

### 9.9 Update Notification Settings
```http
PUT /api/v1/notifications/settings
```

**Access**: 🔒 Protected

**Request Body**:
```json
{
  "replies": true,
  "mentions": false,
  "votes": true,
  "follows": true
}
```

**Response** (200): Updated settings

---

## 10. Tags (`/api/v1/tags`)

### 10.1 List Tags
```http
GET /api/v1/tags?page=1&limit=50
```

**Access**: 🌐 Public

**Response** (200):
```json
{
  "tags": [
    {
      "id": "uuid",
      "name": "golang",
      "postCount": 150
    }
  ],
  "total": 100
}
```

---

### 10.2 Get Popular Tags
```http
GET /api/v1/tags/popular?limit=20
```

**Access**: 🌐 Public

**Response** (200): Array of tags sorted by popularity

---

### 10.3 Search Tags
```http
GET /api/v1/tags/search?q=go&limit=20
```

**Access**: 🌐 Public

**Response** (200): Array of matching tags

---

### 10.4 Get Tag by ID
```http
GET /api/v1/tags/:id
```

**Access**: 🌐 Public

**Response** (200): Single tag object

---

### 10.5 Get Tag by Name
```http
GET /api/v1/tags/name/:name
```

**Access**: 🌐 Public

**Response** (200): Single tag object

---

## 11. Search (`/api/v1/search`)

### 11.1 Universal Search
```http
GET /api/v1/search?q=golang&type=all&page=1&limit=20
```

**Access**: 🔓 Optional Auth

**Query Parameters**:
- `q`: Search query (required)
- `type`: `posts`, `users`, `tags`, `all` (default: `all`)
- `page`, `limit`: Pagination

**Response** (200):
```json
{
  "posts": [ ... ],
  "users": [ ... ],
  "tags": [ ... ],
  "total": {
    "posts": 10,
    "users": 5,
    "tags": 3
  }
}
```

---

### 11.2 Get Popular Searches
```http
GET /api/v1/search/popular?limit=10
```

**Access**: 🌐 Public

**Response** (200):
```json
{
  "searches": [
    {
      "query": "golang",
      "count": 150
    }
  ]
}
```

---

### 11.3 Get Search History
```http
GET /api/v1/search/history?page=1&limit=20
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "history": [
    {
      "id": "uuid",
      "query": "golang fiber",
      "type": "posts",
      "searchedAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### 11.4 Clear Search History
```http
DELETE /api/v1/search/history
```

**Access**: 🔒 Protected

**Response** (200):
```json
{
  "message": "Search history cleared"
}
```

---

### 11.5 Delete Search Item
```http
DELETE /api/v1/search/history/:id
```

**Access**: 🔒 Protected

**Response** (204): No Content

---

## 12. Media (`/api/v1/media`)

### 12.1 Upload Image
```http
POST /api/v1/media/upload/image
```

**Access**: 🔒 Protected

**Content-Type**: `multipart/form-data`

**Form Data**:
- `image`: File (max 300MB)

**Response** (201):
```json
{
  "id": "uuid",
  "url": "https://cdn.bunny.net/storage/image.jpg",
  "type": "image",
  "width": 1920,
  "height": 1080,
  "size": 2048576
}
```

---

### 12.2 Upload Video
```http
POST /api/v1/media/upload/video
```

**Access**: 🔒 Protected

**Content-Type**: `multipart/form-data`

**Form Data**:
- `video`: File (max 300MB)

**Response** (201): Same as image upload + `duration` field

---

### 12.3 Get Media Details
```http
GET /api/v1/media/:id
```

**Access**: 🌐 Public

**Response** (200):
```json
{
  "id": "uuid",
  "userId": "uuid",
  "type": "image",
  "url": "https://cdn.bunny.net/storage/image.jpg",
  "thumbnail": "https://cdn.bunny.net/storage/thumb.jpg",
  "width": 1920,
  "height": 1080,
  "size": 2048576,
  "createdAt": "2024-01-01T00:00:00Z"
}
```

---

### 12.4 Get User's Media
```http
GET /api/v1/media/user/:userId?page=1&limit=20
```

**Access**: 🌐 Public

**Response** (200): Paginated media files

---

### 12.5 Delete Media
```http
DELETE /api/v1/media/:id
```

**Access**: 🔒 Protected (owner only)

**Response** (204): No Content

---

## 13. Push Notifications (`/api/v1/push`)

### 13.1 Get VAPID Public Key
```http
GET /api/v1/push/public-key
```

**Access**: 🌐 Public

**Response** (200):
```json
{
  "publicKey": "BNxX..."
}
```

---

### 13.2 Subscribe to Push
```http
POST /api/v1/push/subscribe
```

**Access**: 🔒 Protected

**Request Body**:
```json
{
  "endpoint": "https://fcm.googleapis.com/fcm/send/...",
  "keys": {
    "p256dh": "BNxX...",
    "auth": "xYz..."
  }
}
```

**Response** (200):
```json
{
  "message": "Subscribed successfully"
}
```

---

### 13.3 Unsubscribe from Push
```http
POST /api/v1/push/unsubscribe
```

**Access**: 🔒 Protected

**Request Body**:
```json
{
  "endpoint": "https://fcm.googleapis.com/fcm/send/..."
}
```

**Response** (200):
```json
{
  "message": "Unsubscribed successfully"
}
```

---

## 14. WebSocket (`/ws`)

### 14.1 WebSocket Connection
```
ws://localhost:8080/ws?token={jwt}&room={roomId}
```

**Access**: 🔓 Optional Auth (via query param or header)

**Query Parameters**:
- `token`: JWT token (optional)
- `room`: Room ID (optional)

**Message Format**:
```json
{
  "type": "message.send",
  "data": {
    "content": "Hello World"
  }
}
```

**Event Types**:
- `message.send` - Send message
- `message.new` - Receive new message
- `notification.new` - Receive notification
- `user.online` - User went online
- `user.offline` - User went offline
- `heartbeat` - Keep connection alive

---

## 15. Health & Info

### 15.1 Health Check
```http
GET /health
```

**Access**: 🌐 Public

**Response** (200):
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

### 15.2 API Info
```http
GET /api/v1
```

**Access**: 🌐 Public

**Response** (200):
```json
{
  "name": "Social Media API",
  "version": "1.0.0",
  "description": "Go Fiber Social Media Platform"
}
```

---

## Error Responses

**Standard Error Format**:
```json
{
  "error": "Error message here"
}
```

**HTTP Status Codes**:
- `200 OK`: Success
- `201 Created`: Resource created
- `204 No Content`: Success with no body
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Duplicate resource
- `500 Internal Server Error`: Server error

---

## Rate Limiting (Planned)

**Not yet implemented**, but recommended:
- Anonymous: 100 requests/hour
- Authenticated: 1000 requests/hour
- Admin: Unlimited

---

## Pagination

**Standard Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

**Standard Response**:
```json
{
  "data": [ ... ],
  "total": 150,
  "page": 1,
  "limit": 20,
  "totalPages": 8
}
```

---

## Sorting

**Common Sort Options**:
- `hot`: Trending (votes + time decay)
- `new`: Latest first
- `top`: Most votes (with time filter)

**Time Filters** (for `top`):
- `hour`, `day`, `week`, `month`, `year`, `all`
