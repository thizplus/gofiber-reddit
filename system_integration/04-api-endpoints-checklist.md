# API Endpoints Checklist

## 📋 Complete List of 61 Endpoints

สรุปทุก endpoint พร้อมตัวอย่าง request/response และ status

---

## 📊 Endpoints Summary

| Module | Public | Private | Total |
|--------|--------|---------|-------|
| **Authentication** | 2 | 3 | 5 |
| **Posts** | 3 | 5 | 8 |
| **Comments** | 2 | 4 | 6 |
| **Users** | 4 | 6 | 10 |
| **Notifications** | 0 | 8 | 8 |
| **Saved Posts** | 0 | 6 | 6 |
| **Search** | 7 | 1 | 8 |
| **Media** | 1 | 5 | 6 |
| **WebSocket** | 1 | 1 | 2 |
| **Health** | 2 | 0 | 2 |
| **TOTAL** | **22** | **39** | **61** |

---

# 1. Authentication Module (5 endpoints)

## ✅ POST /api/v1/auth/register
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test1234",
    "displayName": "Test User"
  }'
```

**Response (201):**
```json
{
  "success": true,
  "message": "ลงทะเบียนสำเร็จ",
  "data": {
    "user": {...},
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

---

## ✅ POST /api/v1/auth/login
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test1234"
  }'
```

---

## ✅ POST /api/v1/auth/logout
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X POST http://localhost:3000/api/v1/auth/logout \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ GET /api/v1/auth/me
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/auth/me \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ POST /api/v1/auth/refresh
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X POST http://localhost:3000/api/v1/auth/refresh \
  -H "Authorization: Bearer TOKEN"
```

---

# 2. Posts Module (8 endpoints)

## ✅ GET /api/v1/posts
**Access:** Public (userVote only if authenticated)
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
# Basic
curl "http://localhost:3000/api/v1/posts?page=1&limit=20"

# With filters
curl "http://localhost:3000/api/v1/posts?page=1&limit=20&sortBy=hot&timeRange=week&tag=golang"
```

**Query Parameters:**
- `page`: number (default: 1)
- `limit`: number (default: 20, max: 100)
- `sortBy`: hot | new | top (default: hot)
- `timeRange`: today | week | month | year | all (for sortBy=top)
- `tag`: string (filter by tag)
- `author`: string (filter by username)

---

## ✅ GET /api/v1/posts/:id
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/posts/POST_ID
```

---

## ✅ POST /api/v1/posts
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X POST http://localhost:3000/api/v1/posts \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my post...",
    "tags": ["golang", "backend"],
    "sourcePostId": null
  }'
```

---

## ✅ PUT /api/v1/posts/:id
**Access:** Private (Owner only)
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X PUT http://localhost:3000/api/v1/posts/POST_ID \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content...",
    "tags": ["golang", "fiber"]
  }'
```

---

## ✅ DELETE /api/v1/posts/:id
**Access:** Private (Owner only)
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/posts/POST_ID \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ POST /api/v1/posts/:id/vote
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
# Upvote
curl -X POST http://localhost:3000/api/v1/posts/POST_ID/vote \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"voteType": "up"}'

# Downvote
curl -X POST http://localhost:3000/api/v1/posts/POST_ID/vote \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"voteType": "down"}'

# Remove vote
curl -X POST http://localhost:3000/api/v1/posts/POST_ID/vote \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"voteType": null}'
```

---

## ✅ GET /api/v1/posts/user/:username
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/posts/user/testuser?page=1&limit=20"
```

---

## ✅ GET /api/v1/posts/me
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/posts/me?page=1&limit=20" \
  -H "Authorization: Bearer TOKEN"
```

---

# 3. Comments Module (6 endpoints)

## ✅ GET /api/v1/posts/:postId/comments
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/posts/POST_ID/comments?sortBy=top"
```

**Query Parameters:**
- `sortBy`: top | new | old (default: top)

---

## ✅ POST /api/v1/posts/:postId/comments
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
# Top-level comment
curl -X POST http://localhost:3000/api/v1/posts/POST_ID/comments \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This is my comment"
  }'

# Reply to comment
curl -X POST http://localhost:3000/api/v1/posts/POST_ID/comments \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This is a reply",
    "parentId": "PARENT_COMMENT_ID"
  }'
```

---

## ✅ GET /api/v1/comments/:id
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/comments/COMMENT_ID
```

---

## ✅ PUT /api/v1/comments/:id
**Access:** Private (Owner only)
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X PUT http://localhost:3000/api/v1/comments/COMMENT_ID \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated comment content"
  }'
```

---

## ✅ DELETE /api/v1/comments/:id
**Access:** Private (Owner only)
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/comments/COMMENT_ID \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ POST /api/v1/comments/:id/vote
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X POST http://localhost:3000/api/v1/comments/COMMENT_ID/vote \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"voteType": "up"}'
```

---

# 4. Users Module (10 endpoints)

## ✅ GET /api/v1/users/:username
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/users/testuser
```

---

## ✅ PUT /api/v1/users/me
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
# Without avatar
curl -X PUT http://localhost:3000/api/v1/users/me \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "displayName": "New Display Name",
    "bio": "My new bio",
    "location": "Bangkok, Thailand",
    "website": "https://example.com"
  }'

# With avatar (multipart)
curl -X PUT http://localhost:3000/api/v1/users/me \
  -H "Authorization: Bearer TOKEN" \
  -F "displayName=New Name" \
  -F "bio=My bio" \
  -F "avatar=@avatar.jpg"
```

---

## ✅ POST /api/v1/users/:username/follow
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X POST http://localhost:3000/api/v1/users/testuser/follow \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ DELETE /api/v1/users/:username/follow
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/users/testuser/follow \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ GET /api/v1/users/:username/followers
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/users/testuser/followers?page=1&limit=20"
```

---

## ✅ GET /api/v1/users/:username/following
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/users/testuser/following?page=1&limit=20"
```

---

## ✅ GET /api/v1/users/:username/comments
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/users/testuser/comments?page=1&limit=20&sortBy=new"
```

---

## ✅ GET /api/v1/users/me/karma
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/users/me/karma?page=1&limit=50&timeRange=week" \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ GET /api/v1/users/search
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/users/search?q=test&page=1&limit=20"
```

---

## ✅ DELETE /api/v1/users/me
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/users/me \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "Test1234",
    "confirmation": "DELETE"
  }'
```

---

# 5. Notifications Module (8 endpoints)

## ✅ GET /api/v1/notifications
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
# All notifications
curl "http://localhost:3000/api/v1/notifications?page=1&limit=20" \
  -H "Authorization: Bearer TOKEN"

# Unread only
curl "http://localhost:3000/api/v1/notifications?filter=unread&page=1&limit=20" \
  -H "Authorization: Bearer TOKEN"

# By type
curl "http://localhost:3000/api/v1/notifications?type=reply&page=1&limit=20" \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ GET /api/v1/notifications/unread-count
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/notifications/unread-count \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ PUT /api/v1/notifications/:id/read
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X PUT http://localhost:3000/api/v1/notifications/NOTIFICATION_ID/read \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ PUT /api/v1/notifications/read-all
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X PUT http://localhost:3000/api/v1/notifications/read-all \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ DELETE /api/v1/notifications/:id
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/notifications/NOTIFICATION_ID \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ DELETE /api/v1/notifications/read
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/notifications/read \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ GET /api/v1/notifications/settings
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/notifications/settings \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ PUT /api/v1/notifications/settings
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X PUT http://localhost:3000/api/v1/notifications/settings \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "replies": true,
    "mentions": true,
    "votes": false,
    "follows": true,
    "emailNotifications": false
  }'
```

---

# 6. Saved Posts Module (6 endpoints)

## ✅ GET /api/v1/saved
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/saved?page=1&limit=20&sortBy=new" \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ GET /api/v1/saved/:postId
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/saved/POST_ID \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ POST /api/v1/saved/:postId
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X POST http://localhost:3000/api/v1/saved/POST_ID \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ DELETE /api/v1/saved/:postId
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/saved/POST_ID \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ DELETE /api/v1/saved
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/saved \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ GET /api/v1/saved/count
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/saved/count \
  -H "Authorization: Bearer TOKEN"
```

---

# 7. Search Module (8 endpoints)

## ✅ GET /api/v1/search
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/search?q=golang&type=all&page=1&limit=20"
```

---

## ✅ GET /api/v1/search/posts
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/search/posts?q=golang&sortBy=relevance&page=1&limit=20"

# With filters
curl "http://localhost:3000/api/v1/search/posts?q=golang&tag=backend&author=testuser&timeRange=week"
```

---

## ✅ GET /api/v1/search/users
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/search/users?q=test&sortBy=karma&page=1&limit=20"
```

---

## ✅ GET /api/v1/search/suggestions
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/search/suggestions?q=gol&limit=10"
```

---

## ✅ GET /api/v1/search/tags
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/search/tags?limit=20&timeRange=week"
```

---

## ✅ GET /api/v1/search/trending
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/search/trending?limit=10&timeRange=today"
```

---

## ✅ GET /api/v1/search/history
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/search/history?limit=20" \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ DELETE /api/v1/search/history
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/search/history \
  -H "Authorization: Bearer TOKEN"
```

---

# 8. Media Module (6 endpoints)

## ✅ POST /api/v1/media/upload
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
# Single file
curl -X POST http://localhost:3000/api/v1/media/upload \
  -H "Authorization: Bearer TOKEN" \
  -F "files=@image.jpg"

# Multiple files
curl -X POST http://localhost:3000/api/v1/media/upload \
  -H "Authorization: Bearer TOKEN" \
  -F "files=@image1.jpg" \
  -F "files=@image2.jpg" \
  -F "type=post"
```

---

## ✅ GET /api/v1/media/:id
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/media/MEDIA_ID
```

---

## ✅ DELETE /api/v1/media/:id
**Access:** Private (Owner only)
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X DELETE http://localhost:3000/api/v1/media/MEDIA_ID \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ GET /api/v1/media/me
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl "http://localhost:3000/api/v1/media/me?page=1&limit=20&type=image&sortBy=new" \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ GET /api/v1/media/storage
**Access:** Private
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/media/storage \
  -H "Authorization: Bearer TOKEN"
```

---

## ✅ POST /api/v1/media/:id/optimize
**Access:** Private (Owner only)
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl -X POST http://localhost:3000/api/v1/media/MEDIA_ID/optimize \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quality": 80,
    "maxWidth": 1920,
    "maxHeight": 1080
  }'
```

---

# 9. Additional Endpoints

## ✅ GET /health
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/health
```

---

## ✅ GET /api/v1/health
**Access:** Public
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```bash
curl http://localhost:3000/api/v1/health
```

---

## ✅ GET /ws (WebSocket)
**Access:** Public (upgrade to WebSocket)
**Status:** `[ ]` Not Started | `[ ]` In Progress | `[ ]` Done

```javascript
// Frontend example
const ws = new WebSocket('ws://localhost:3000/ws');
ws.onmessage = (event) => {
  const notification = JSON.parse(event.data);
  console.log('New notification:', notification);
};
```

---

# 📝 Testing Checklist

## Module Testing Status

### Authentication
- [ ] Register new user
- [ ] Login existing user
- [ ] Get current user
- [ ] Refresh token
- [ ] Logout

### Posts
- [ ] Create post
- [ ] Get all posts (sorted by hot/new/top)
- [ ] Get single post
- [ ] Update own post
- [ ] Delete own post
- [ ] Vote on post (up/down/remove)
- [ ] Get posts by user
- [ ] Get own posts

### Comments
- [ ] Create top-level comment
- [ ] Create reply comment
- [ ] Get comments for post
- [ ] Get single comment
- [ ] Update own comment
- [ ] Delete own comment
- [ ] Vote on comment
- [ ] Test nested replies (max depth 10)

### Users
- [ ] Get user profile
- [ ] Update own profile (with avatar)
- [ ] Follow user
- [ ] Unfollow user
- [ ] Get followers list
- [ ] Get following list
- [ ] Get user's comments
- [ ] Get karma history
- [ ] Search users
- [ ] Delete account

### Notifications
- [ ] Get all notifications
- [ ] Get unread count
- [ ] Mark as read (single)
- [ ] Mark all as read
- [ ] Delete notification
- [ ] Delete all read
- [ ] Get settings
- [ ] Update settings
- [ ] Test auto-creation on events

### Saved Posts
- [ ] Save post
- [ ] Unsave post
- [ ] Get saved posts
- [ ] Check if post is saved
- [ ] Get saved count
- [ ] Clear all saved

### Search
- [ ] Search everything
- [ ] Search posts only
- [ ] Search users only
- [ ] Get suggestions
- [ ] Get popular tags
- [ ] Get trending searches
- [ ] Search history (save/get/delete)

### Media
- [ ] Upload single image
- [ ] Upload multiple images
- [ ] Upload video
- [ ] Get media details
- [ ] Delete media
- [ ] Get user's media
- [ ] Get storage usage
- [ ] Optimize image

---

## 🎯 Final Checklist

- [ ] All 61 endpoints implemented
- [ ] All endpoints tested manually
- [ ] Integration tests written
- [ ] Error handling tested
- [ ] Rate limiting working
- [ ] Authentication working
- [ ] Authorization working (owner/admin)
- [ ] Pagination working everywhere
- [ ] Sorting/filtering working
- [ ] Media upload to Bunny CDN working
- [ ] Notifications auto-creating
- [ ] Karma updating correctly
- [ ] Vote system working
- [ ] Search relevance ranking working
- [ ] WebSocket real-time updates
- [ ] Documentation complete

---

**All Endpoints Complete? → Ready for Production!** 🚀
