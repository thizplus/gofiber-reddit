# Backend API Overview

## Introduction
This document provides a comprehensive overview of the backend API for the Social Media Platform. The API follows RESTful principles and uses JSON for request/response payloads.

---

## Base URL
```
Production: https://api.example.com
Development: http://localhost:3000
```

---

## Authentication

### Token-Based Authentication (JWT)
- All private endpoints require JWT token in the Authorization header
- Token format: `Authorization: Bearer {token}`
- Token expires after 24 hours
- Refresh token endpoint available to get new token

### Public vs Private Routes

**Public Routes** (No authentication required):
- Authentication: Login, Register
- Posts: Get all posts, Get post by ID, Get user's posts
- Comments: Get comments for post, Get comment by ID
- Users: Get user profile, Search users, Get followers/following
- Search: All search endpoints
- Media: Get media details

**Private Routes** (Authentication required):
- Authentication: Logout, Get current user, Refresh token
- Posts: Create, Update, Delete, Vote
- Comments: Create, Update, Delete, Vote
- Users: Update profile, Follow/Unfollow, Delete account, Get karma history
- Notifications: All notification endpoints
- Saved Posts: All saved posts endpoints
- Media: Upload, Delete, Optimize, Get user's media

---

## API Modules

### 1. Authentication (`01-authentication.md`)
User registration, login, logout, and session management.

**Key Endpoints:**
- `POST /api/auth/register` - Create new account
- `POST /api/auth/login` - Sign in
- `POST /api/auth/logout` - Sign out
- `GET /api/auth/me` - Get current user

---

### 2. Posts (`02-posts.md`)
CRUD operations for posts, voting, and crossposting.

**Key Endpoints:**
- `GET /api/posts` - Get all posts (paginated)
- `GET /api/posts/:id` - Get single post
- `POST /api/posts` - Create new post
- `PUT /api/posts/:id` - Update post
- `DELETE /api/posts/:id` - Delete post
- `POST /api/posts/:id/vote` - Vote on post

---

### 3. Comments (`03-comments.md`)
CRUD operations for comments and nested replies.

**Key Endpoints:**
- `GET /api/posts/:postId/comments` - Get all comments for post
- `POST /api/posts/:postId/comments` - Create comment/reply
- `PUT /api/comments/:id` - Update comment
- `DELETE /api/comments/:id` - Delete comment
- `POST /api/comments/:id/vote` - Vote on comment

---

### 4. Users (`04-users.md`)
User profiles, following system, and user management.

**Key Endpoints:**
- `GET /api/users/:username` - Get user profile
- `PUT /api/users/me` - Update own profile
- `POST /api/users/:username/follow` - Follow user
- `DELETE /api/users/:username/follow` - Unfollow user
- `GET /api/users/:username/followers` - Get followers list
- `GET /api/users/:username/following` - Get following list
- `GET /api/users/:username/comments` - Get user's comments

---

### 5. Notifications (`05-notifications.md`)
User notification system for interactions.

**Key Endpoints:**
- `GET /api/notifications` - Get all notifications
- `GET /api/notifications/unread-count` - Get unread count
- `PUT /api/notifications/:id/read` - Mark as read
- `PUT /api/notifications/read-all` - Mark all as read
- `DELETE /api/notifications/:id` - Delete notification

---

### 6. Saved Posts (`06-saved-posts.md`)
Bookmark system for posts.

**Key Endpoints:**
- `GET /api/saved` - Get saved posts
- `POST /api/saved/:postId` - Save post
- `DELETE /api/saved/:postId` - Unsave post
- `GET /api/saved/count` - Get saved count

---

### 7. Search (`07-search.md`)
Full-text search for posts and users.

**Key Endpoints:**
- `GET /api/search` - Search everything
- `GET /api/search/posts` - Search posts only
- `GET /api/search/users` - Search users only
- `GET /api/search/suggestions` - Get autocomplete suggestions
- `GET /api/search/tags` - Get popular tags

---

### 8. Media (`08-media.md`)
File upload and media management.

**Key Endpoints:**
- `POST /api/media/upload` - Upload media files
- `GET /api/media/:id` - Get media details
- `DELETE /api/media/:id` - Delete media
- `GET /api/media/me` - Get user's media
- `GET /api/media/storage` - Get storage usage

---

## Common Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "errors": {
    "field": "Field-specific error message"
  }
}
```

---

## HTTP Status Codes

- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required or failed
- `403 Forbidden` - Authenticated but not authorized
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., duplicate)
- `413 Payload Too Large` - Request body too large
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

---

## Pagination

All paginated endpoints follow this format:

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20, max: 100)
```

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "currentPage": 1,
      "totalPages": 10,
      "totalItems": 195,
      "hasNextPage": true,
      "hasPrevPage": false
    }
  }
}
```

---

## Rate Limiting

Rate limits are applied per IP address or authenticated user:

- **Authentication endpoints**: 5 requests per 15 minutes
- **Search endpoints**: 30 requests per minute
- **Media upload**: 10 uploads per minute
- **Post/Comment creation**: 10 per minute
- **Vote actions**: 60 per minute
- **General API**: 100 requests per minute

**Rate Limit Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704876543
```

**Rate Limit Exceeded Response (429):**
```json
{
  "success": false,
  "message": "คำขอมากเกินไป กรุณาลองใหม่อีกครั้งใน 30 วินาที",
  "retryAfter": 30
}
```

---

## CORS

Cross-Origin Resource Sharing is enabled for:
- Development: `http://localhost:3000`, `http://localhost:3001`
- Production: `https://yourdomain.com`

**Allowed Methods:** GET, POST, PUT, DELETE, OPTIONS
**Allowed Headers:** Content-Type, Authorization

---

## Security

### Request Validation
- All inputs validated and sanitized
- SQL injection prevention (parameterized queries)
- XSS prevention (output encoding)
- CSRF tokens for state-changing operations

### Password Security
- Minimum 8 characters
- Hashed with bcrypt (10+ rounds)
- Never returned in API responses

### File Upload Security
- File type validation by content
- Malware scanning
- File size limits enforced
- EXIF data stripped from images

---

## Error Handling

All errors follow consistent format with:
- Human-readable message in Thai
- Error code for programmatic handling
- Field-specific errors when applicable

Example:
```json
{
  "success": false,
  "message": "ข้อมูลไม่ถูกต้อง",
  "code": "VALIDATION_ERROR",
  "errors": {
    "email": "รูปแบบอีเมลไม่ถูกต้อง",
    "password": "รหัสผ่านต้องมีอย่างน้อย 8 ตัวอักษร"
  }
}
```

See `09-error-codes.md` for complete error code reference.

---

## Data Types

### Date/Time Format
All dates use ISO 8601 format in UTC:
```
2025-01-12T10:00:00Z
```

### User Object
```typescript
{
  id: string
  username: string
  displayName: string
  avatar: string | null
  karma: number
  bio: string | null
  coverImage: string | null
  joinedAt: string
  location: string | null
  website: string | null
  followersCount: number
  followingCount: number
}
```

### Post Object
```typescript
{
  id: string
  title: string
  content: string
  author: User
  votes: number
  userVote: 'up' | 'down' | null
  commentCount: number
  createdAt: string
  updatedAt: string | null
  media: Media[] | null
  tags: string[]
  sourcePostId: string | null
  sourcePost: Post | null
}
```

### Comment Object
```typescript
{
  id: string
  postId: string
  author: User
  content: string
  votes: number
  userVote: 'up' | 'down' | null
  parentId: string | null
  replies: Comment[]
  depth: number
  createdAt: string
  updatedAt: string | null
}
```

---

## WebSocket/Real-time (Future)

Real-time features planned:
- Live notification delivery
- Real-time vote updates
- Live comment updates
- Online user presence

**WebSocket Endpoint:** `wss://api.example.com/ws`

---

## API Versioning

Current version: `v1`

Version can be specified in:
1. URL: `/api/v1/posts`
2. Header: `X-API-Version: 1`

Breaking changes will introduce new version.

---

## Development Notes

### Testing
- Use Postman collection (available in `/docs/postman`)
- Mock data available in frontend `/lib/data`
- Test users available (username: `testuser`, password: `Test1234`)

### Database Schema
- See `/docs/database-schema.md` for complete schema
- Migrations in `/migrations` directory

### Environment Variables
Required environment variables:
```
DATABASE_URL=postgresql://...
JWT_SECRET=your-secret-key
AWS_S3_BUCKET=your-bucket
AWS_ACCESS_KEY=...
AWS_SECRET_KEY=...
```

---

## Support

For API questions or issues:
- Documentation: https://docs.example.com
- GitHub Issues: https://github.com/yourrepo/issues
- Email: support@example.com
