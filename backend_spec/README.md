# Backend API Specification

## Overview
This directory contains comprehensive API specifications for the Social Media Platform backend. Each file documents specific API modules with detailed request/response formats, error handling, and business logic.

---

## Quick Start

1. **Start Here:** Read `00-overview.md` for API architecture and general concepts
2. **Authentication:** Review `01-authentication.md` for auth flow
3. **Core Features:** Read modules 02-08 based on feature you're implementing
4. **Error Handling:** Reference `09-error-codes.md` for error codes

---

## File Structure

```
backend_spec/
├── README.md                    # This file
├── 00-overview.md              # API overview and general concepts
├── 01-authentication.md        # Auth: login, register, logout
├── 02-posts.md                 # Posts: CRUD, voting, crossposting
├── 03-comments.md              # Comments: CRUD, replies, voting
├── 04-users.md                 # Users: profiles, following, management
├── 05-notifications.md         # Notifications: delivery, settings
├── 06-saved-posts.md           # Saved: bookmark system
├── 07-search.md                # Search: full-text search for posts/users
├── 08-media.md                 # Media: file upload and management
└── 09-error-codes.md           # Complete error code reference
```

---

## API Modules

### 00. Overview (`00-overview.md`)
**Purpose:** High-level API architecture, authentication flow, response formats

**Key Topics:**
- Base URL and versioning
- Authentication mechanism (JWT)
- Public vs Private routes
- Common response format
- Pagination format
- Rate limiting
- CORS policy
- Security considerations

**When to read:** Before starting any implementation

---

### 01. Authentication (`01-authentication.md`)
**Endpoints:** 5 endpoints

**Key Features:**
- User registration with validation
- Login with username/email
- JWT token-based auth
- Token refresh mechanism
- Logout and session invalidation
- Get current user info

**Dependencies:** None (entry point)

**When to implement:** First (required for all private endpoints)

---

### 02. Posts (`02-posts.md`)
**Endpoints:** 8 endpoints

**Key Features:**
- Get posts (paginated, filtered, sorted)
- Get single post with sourcePost
- Create post with media and tags
- Update post (title, content, tags)
- Delete post (soft delete)
- Vote system (upvote/downvote)
- Get user's posts
- Get own posts

**Dependencies:**
- Authentication (for private endpoints)
- Media API (for file uploads)

**When to implement:** After authentication

---

### 03. Comments (`03-comments.md`)
**Endpoints:** 6 endpoints

**Key Features:**
- Get comments with nested replies
- Create comment/reply (unlimited nesting)
- Update comment
- Delete comment (soft delete)
- Vote on comments
- Get single comment with thread

**Dependencies:**
- Authentication (for private endpoints)
- Posts API (comments belong to posts)

**When to implement:** After posts

---

### 04. Users (`04-users.md`)
**Endpoints:** 10 endpoints

**Key Features:**
- Get user profile (public)
- Update own profile (with avatar upload)
- Follow/unfollow users
- Get followers/following lists
- Get user's comments
- Get karma history
- Search users
- Delete account

**Dependencies:**
- Authentication
- Media API (for avatar upload)

**When to implement:** Can implement in parallel with posts

---

### 05. Notifications (`05-notifications.md`)
**Endpoints:** 8 endpoints

**Key Features:**
- Get notifications (filtered, paginated)
- Get unread count
- Mark as read (single/all)
- Delete notifications
- Notification settings/preferences

**Notification Types:**
- Reply to comment/post
- Vote on post/comment
- Mention (@username)
- New follower

**Dependencies:**
- Authentication
- Posts, Comments, Users (create notifications)

**When to implement:** After core features (posts, comments, users)

---

### 06. Saved Posts (`06-saved-posts.md`)
**Endpoints:** 6 endpoints

**Key Features:**
- Get saved posts (paginated)
- Check if post is saved
- Save post
- Unsave post
- Clear all saved
- Get saved count

**Dependencies:**
- Authentication
- Posts API

**When to implement:** Can implement anytime after posts

---

### 07. Search (`07-search.md`)
**Endpoints:** 8 endpoints

**Key Features:**
- Search everything (posts + users)
- Search posts only (with filters)
- Search users only
- Autocomplete suggestions
- Popular tags
- Trending searches
- Search history (private)

**Search Features:**
- Full-text search with relevance ranking
- Highlighting matched terms
- Advanced filters (tag, author, time range)
- Special queries (@username, #tag)

**Dependencies:**
- Posts API
- Users API
- Full-text search engine (Elasticsearch recommended)

**When to implement:** After core features are stable

---

### 08. Media (`08-media.md`)
**Endpoints:** 6 endpoints

**Key Features:**
- Upload media (images/videos)
- Get media details
- Delete media
- Get user's media library
- Get storage usage
- Optimize images

**Media Processing:**
- Image: compression, thumbnails, WebP conversion
- Video: transcoding, thumbnails, quality versions
- Security: malware scan, EXIF stripping

**Dependencies:**
- Authentication
- Cloud storage (S3, Cloudinary, etc.)
- Image processing library

**When to implement:** Can implement early (needed for posts/users)

---

### 09. Error Codes (`09-error-codes.md`)
**Purpose:** Complete error code reference

**Categories:**
- Authentication errors (AUTH_*)
- Validation errors (VAL_*)
- Resource errors (RES_*)
- Permission errors (PERM_*)
- Conflict errors (CONF_*)
- File upload errors (FILE_*)
- Rate limit errors (RATE_*)
- Server errors (SRV_*)
- Business logic errors (BIZ_*)
- Search errors (SRCH_*)

**When to use:** Reference during error handling implementation

---

## Implementation Order

### Phase 1: Foundation
1. Setup project (Express/Fastify + TypeScript)
2. Database schema and migrations
3. Authentication system (`01-authentication.md`)
4. Error handling middleware (`09-error-codes.md`)

### Phase 2: Core Features
5. Media upload (`08-media.md`)
6. Posts API (`02-posts.md`)
7. Comments API (`03-comments.md`)
8. Users API (`04-users.md`)

### Phase 3: Social Features
9. Saved posts (`06-saved-posts.md`)
10. Notifications (`05-notifications.md`)
11. Search (`07-search.md`)

### Phase 4: Optimization
12. Caching (Redis)
13. Rate limiting
14. Performance optimization
15. Real-time features (WebSocket)

---

## Technology Stack Recommendations

### Backend Framework
- **Node.js + Express** (popular, flexible)
- **Node.js + Fastify** (faster, modern)
- **Nest.js** (enterprise, TypeScript-first)

### Database
- **PostgreSQL** (recommended, supports full-text search)
- **MySQL** (alternative)
- **MongoDB** (document-based alternative)

### ORM
- **Prisma** (modern, type-safe)
- **TypeORM** (feature-rich)
- **Drizzle** (lightweight)

### Authentication
- **JWT** (jsonwebtoken package)
- **bcrypt** (password hashing)

### File Storage
- **AWS S3** (reliable, scalable)
- **Cloudinary** (image optimization built-in)
- **DigitalOcean Spaces** (S3-compatible, cheaper)

### Search
- **PostgreSQL Full-Text Search** (built-in, good for start)
- **Elasticsearch** (powerful, recommended for scale)
- **MeiliSearch** (fast, easier than ES)

### Caching
- **Redis** (in-memory cache)
- **Node-cache** (in-process, simple)

### Real-time
- **Socket.io** (WebSocket)
- **Server-Sent Events** (simpler alternative)

---

## Database Schema Overview

### Core Tables
```
users
  - id, username, email, password_hash
  - display_name, avatar, bio, karma
  - created_at, updated_at

posts
  - id, title, content, author_id
  - votes, comment_count
  - source_post_id (for crosspost)
  - created_at, updated_at, deleted_at

comments
  - id, post_id, author_id, parent_id
  - content, votes, depth
  - created_at, updated_at, deleted_at

media
  - id, user_id, type (image/video)
  - url, thumbnail, size, width, height
  - created_at

post_media (junction)
  - post_id, media_id

votes
  - user_id, target_id, target_type (post/comment)
  - vote_type (up/down)
  - created_at

follows
  - follower_id, following_id
  - created_at

saved_posts
  - user_id, post_id
  - saved_at

notifications
  - id, user_id, sender_id, type
  - post_id, comment_id, message
  - is_read, created_at

tags
  - id, name, post_count

post_tags (junction)
  - post_id, tag_id
```

---

## API Testing

### Recommended Tools
- **Postman** - API testing, collection sharing
- **Insomnia** - Alternative to Postman
- **Thunder Client** - VS Code extension
- **cURL** - Command line testing

### Test Data
Use mock data from frontend:
```
/lib/data/mock-users.json
/lib/data/mock-posts.json
/lib/data/mock-comments.json
/lib/data/mock-notifications.json
```

### Test Users
Create these test accounts:
```
Username: thepthai
Password: Test1234
Email: thepthai@example.com

Username: cookingmaster
Password: Test1234
Email: cooking@example.com
```

---

## Security Checklist

### Authentication
- [ ] Passwords hashed with bcrypt (10+ rounds)
- [ ] JWT tokens with expiration
- [ ] Token refresh mechanism
- [ ] Logout invalidates token
- [ ] Rate limiting on login/register

### Input Validation
- [ ] All inputs validated
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (sanitize HTML)
- [ ] CSRF protection
- [ ] File upload validation (type, size, content)

### Authorization
- [ ] Verify ownership before update/delete
- [ ] Check permissions on all private routes
- [ ] Validate user can access resource

### Data Protection
- [ ] Never expose passwords in responses
- [ ] Sanitize sensitive data in logs
- [ ] HTTPS only in production
- [ ] Secure cookie flags

### Rate Limiting
- [ ] Login attempts limited
- [ ] API rate limits per user/IP
- [ ] File upload limits
- [ ] Search query limits

---

## Environment Variables

Create `.env` file:
```env
# Server
NODE_ENV=development
PORT=3000
API_URL=http://localhost:3000

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/social_db

# Authentication
JWT_SECRET=your-super-secret-key-change-this
JWT_EXPIRE=24h

# File Storage (AWS S3)
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=us-east-1
AWS_S3_BUCKET=your-bucket-name

# Or Cloudinary
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret

# Redis (optional)
REDIS_URL=redis://localhost:6379

# Email (optional, for notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-password

# Rate Limiting
RATE_LIMIT_WINDOW=60000
RATE_LIMIT_MAX=100
```

---

## Common Patterns

### Controller Pattern
```typescript
// posts.controller.ts
export const getPosts = async (req, res) => {
  try {
    const { page = 1, limit = 20 } = req.query;
    const posts = await postService.getPosts(page, limit);

    res.json({
      success: true,
      data: posts
    });
  } catch (error) {
    next(error);
  }
};
```

### Service Pattern
```typescript
// posts.service.ts
export const getPosts = async (page: number, limit: number) => {
  const offset = (page - 1) * limit;

  const [posts, total] = await Promise.all([
    db.posts.findMany({ skip: offset, take: limit }),
    db.posts.count()
  ]);

  return {
    posts,
    pagination: {
      currentPage: page,
      totalPages: Math.ceil(total / limit),
      totalPosts: total,
      hasNextPage: offset + limit < total,
      hasPrevPage: page > 1
    }
  };
};
```

### Middleware Pattern
```typescript
// auth.middleware.ts
export const requireAuth = async (req, res, next) => {
  const token = req.headers.authorization?.replace('Bearer ', '');

  if (!token) {
    return res.status(401).json({
      success: false,
      code: 'AUTH_002',
      message: 'ไม่พบ Token กรุณาเข้าสู่ระบบ'
    });
  }

  try {
    const decoded = jwt.verify(token, JWT_SECRET);
    req.user = await db.users.findUnique({ where: { id: decoded.userId } });
    next();
  } catch (error) {
    return res.status(401).json({
      success: false,
      code: 'AUTH_003',
      message: 'Token ไม่ถูกต้อง'
    });
  }
};
```

---

## Need Help?

- Check `00-overview.md` for general concepts
- Check specific module files for detailed endpoints
- Check `09-error-codes.md` for error handling
- Create GitHub issue for questions
- Contact: support@example.com

---

## Contributing

When updating this specification:
1. Keep response formats consistent
2. Update error codes in `09-error-codes.md`
3. Update this README if adding new modules
4. Use Thai language for user-facing messages
5. Include examples for all endpoints

---

## Version History

- **v1.0** (2025-01-12) - Initial specification based on frontend requirements
  - All core features documented
  - Complete error code reference
  - Implementation guidelines included
