# API Route Specifications

Complete documentation for all API endpoints in the social media backend.

**Base URL:** `http://localhost:8080/api/v1`

**Server Status:** ✅ Running on port 8080

## Quick Links

1. [Authentication Routes](01_auth_routes.md) - Registration & Login
2. [User Routes](02_user_routes.md) - Profile management
3. [Post Routes](03_post_routes.md) - Create, read, update posts
4. [Comment Routes](04_comment_routes.md) - Nested comments & replies
5. [Vote Routes](05_vote_routes.md) - Upvote/downvote system
6. [Follow Routes](06_follow_routes.md) - User following
7. [Saved Post Routes](07_saved_post_routes.md) - Bookmark posts
8. [Notification Routes](08_notification_routes.md) - User notifications
9. [Tag Routes](09_tag_routes.md) - Content tagging
10. [Search Routes](10_search_routes.md) - Global search
11. [Media Routes](11_media_routes.md) - Image/video upload

---

## Overview

### Total Endpoints: 71

#### By Category:
- **Authentication & Users:** 7 endpoints
- **Posts:** 11 endpoints
- **Comments:** 9 endpoints
- **Votes:** 5 endpoints
- **Follows:** 6 endpoints
- **Saved Posts:** 4 endpoints
- **Notifications:** 10 endpoints
- **Tags:** 5 endpoints
- **Search:** 5 endpoints
- **Media:** 5 endpoints
- **Legacy (Tasks, Files, Jobs):** 4 endpoints

---

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

**Token Expiry:** 7 days

**Get Token:**
```bash
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}
```

---

## Response Format

All responses follow a consistent format:

**Success Response:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Operation failed",
  "error": "Error details"
}
```

**Validation Error Response:**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "field1": "Error message",
    "field2": "Error message"
  }
}
```

---

## Common Parameters

### Pagination

Most list endpoints support pagination:

**Query Parameters:**
- `offset`: Starting position (default: 0)
- `limit`: Number of items (default: 20)

**Example:**
```
GET /api/v1/posts?offset=20&limit=10
```

**Response includes meta:**
```json
{
  "data": {
    "posts": [...],
    "meta": {
      "total": 100,
      "offset": 20,
      "limit": 10
    }
  }
}
```

### Sorting

Post and comment endpoints support sorting:

**Sort Options:**
- `hot` - Reddit-style hot ranking (default)
- `new` - Most recent first
- `top` - Highest votes first
- `controversial` - Most controversial

**Example:**
```
GET /api/v1/posts?sortBy=top
```

---

## Feature Highlights

### 1. Nested Comments (Reddit-style)
- Up to 10 levels of nesting
- Tree structure with parent chains
- Reply counting
- Hot score algorithm

### 2. Voting System
- Polymorphic (posts & comments)
- Vote change tracking
- Automatic score updates
- Karma calculation

### 3. Smart Notifications
- Type-based (reply, vote, mention, follow)
- User preferences/settings
- Read/unread tracking
- Batch operations

### 4. Content Discovery
- Tag-based organization
- Full-text search
- Popular/trending content
- Personalized feed

### 5. Media Handling
- Bunny CDN integration
- Automatic thumbnails
- Image optimization
- Video support (up to 100MB)

### 6. Social Features
- Follow/unfollow users
- Mutual follows (friends)
- Saved posts (bookmarks)
- User karma tracking

---

## Quick Start Examples

### 1. Register & Login

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "password123",
    "displayName": "John Doe"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 2. Create a Post

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "Hello world!",
    "tags": ["golang", "tutorial"]
  }'
```

### 3. Add a Comment

```bash
curl -X POST http://localhost:8080/api/v1/comments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "postId": "<post-uuid>",
    "content": "Great post!"
  }'
```

### 4. Upvote a Post

```bash
curl -X POST http://localhost:8080/api/v1/votes \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "targetId": "<post-uuid>",
    "targetType": "post",
    "voteType": "upvote"
  }'
```

### 5. Upload an Image

```bash
curl -X POST http://localhost:8080/api/v1/media/upload/image \
  -H "Authorization: Bearer <token>" \
  -F "image=@/path/to/photo.jpg"
```

---

## Technical Details

### Database
- **PostgreSQL** with GORM
- 11 new tables for social features
- Proper indexes and foreign keys
- Soft delete support

### Storage
- **Bunny CDN** for media files
- Direct URL access
- Automatic thumbnail generation
- Image optimization

### Architecture
- **Clean Architecture** (Domain, Application, Infrastructure, Interfaces)
- Repository pattern
- Service layer
- Dependency injection

### Security
- **JWT Authentication** (7-day expiry)
- Password hashing (bcrypt)
- Input validation
- Authorization checks

### Performance
- Batch operations for efficiency
- Database indexes
- Hot score caching
- Pagination

---

## Testing

### Health Check

```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-11-04T10:00:00Z"
}
```

### WebSocket Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  console.log('Received:', event.data);
};
```

---

## Error Handling

### Common HTTP Status Codes

- **200 OK** - Success
- **400 Bad Request** - Validation error or invalid input
- **401 Unauthorized** - Missing or invalid token
- **403 Forbidden** - No permission
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server error

### Error Response Examples

**Validation Error:**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "email": "Email is required",
    "password": "Password must be at least 8 characters"
  }
}
```

**Authentication Error:**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Invalid or expired token"
}
```

**Not Found Error:**
```json
{
  "success": false,
  "message": "Post not found",
  "error": "Post with ID xyz does not exist"
}
```

---

## Rate Limiting

Currently **no rate limiting** is implemented. Consider adding rate limiting for production:

- Authentication endpoints: 5 requests/minute
- API endpoints: 100 requests/minute
- Upload endpoints: 10 requests/minute

---

## Changelog

### Version 1.0.0 (2025-11-04)

**Initial Release**
- Complete social media API
- 71 endpoints across 11 feature areas
- JWT authentication
- Nested comments
- Voting system
- File uploads
- Full-text search
- Notifications
- User following
- Post bookmarking

---

## Support & Contact

For issues or questions:
- Check individual route documentation
- Review error messages
- Test with provided curl examples
- Verify authentication token

---

## Next Steps for Frontend Development

1. **Start with Authentication**
   - Implement login/register pages
   - Store JWT token securely
   - Add token refresh logic

2. **Build Core Features**
   - Post listing and detail pages
   - Comment threads (nested)
   - Voting UI
   - User profiles

3. **Add Social Features**
   - Follow/unfollow buttons
   - Notification dropdown
   - Saved posts page
   - User search

4. **Implement Media Upload**
   - Image upload with preview
   - Video upload with progress
   - Media gallery

5. **Polish & Optimize**
   - Pagination/infinite scroll
   - Real-time updates (WebSocket)
   - Caching strategies
   - Error handling

---

**Last Updated:** 2025-11-04
**API Version:** 1.0.0
**Server:** Running on http://localhost:8080
