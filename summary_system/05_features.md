# Features & Capabilities

## 1. Authentication System

### 1.1 ระบบ Authentication แบบ Dual-mode

**Standard Authentication (Email/Password)**:
- ✅ สมัครสมาชิกด้วย email, username, password
- ✅ Login ด้วย username หรือ email
- ✅ Password hashing ด้วย bcrypt (secure)
- ✅ JWT token-based authentication
- ✅ Token expiry (24 hours, configurable)

**OAuth 2.0 (Google)**:
- ✅ Google Sign-in integration
- ✅ Authorization Code Flow (secure)
- ✅ Temporary code exchange (5-minute TTL)
- ✅ Auto-create user from Google profile
- ✅ Link Google account to existing user

**Security Features**:
- ✅ Password strength validation (min 8 chars)
- ✅ Email format validation
- ✅ Username length validation (3-20 chars)
- ✅ Protected routes with middleware
- ✅ Role-based access control (user/admin)
- ✅ Optional authentication (for public endpoints)

---

## 2. User Management

### 2.1 User Profiles

**Profile Information**:
- ✅ Basic: username, email, display name
- ✅ Avatar: profile picture URL
- ✅ Bio: personal description
- ✅ Location: user's location
- ✅ Website: personal website/portfolio
- ✅ Social stats: karma, followers, following

**Profile Features**:
- ✅ View own profile (private info)
- ✅ View public profile (by username)
- ✅ Update profile information
- ✅ Upload avatar (via media system)
- ✅ Delete account (soft delete)

### 2.2 Karma System

**Karma Calculation**:
- +1 karma: Receive upvote on post
- +1 karma: Receive upvote on comment
- -1 karma: Receive downvote on post
- -1 karma: Receive downvote on comment

**Features**:
- ✅ Cumulative karma score
- ✅ Indexed for leaderboard queries
- ✅ Display on profile
- ✅ Reputation indicator

---

## 3. Content Management

### 3.1 Posts System

**Create & Manage Posts**:
- ✅ Create text posts (title + content)
- ✅ Attach multiple images (via media IDs)
- ✅ Attach multiple videos (via media IDs)
- ✅ Add tags (auto-create if not exist)
- ✅ Update own posts
- ✅ Delete own posts (soft delete)
- ✅ View single post
- ✅ List posts (paginated)

**Post Features**:
- ✅ Rich text content (Markdown support on frontend)
- ✅ Media gallery (images + videos)
- ✅ Tag system (multi-tag support)
- ✅ Vote count display
- ✅ Comment count tracking
- ✅ Author information
- ✅ Timestamps (created, updated)

**Content Discovery**:
- ✅ Sort by: Hot, New, Top
- ✅ Filter by tag
- ✅ Filter by author
- ✅ Search posts
- ✅ Personalized feed (followed users)

### 3.2 Comments System

**Nested Comments**:
- ✅ Threaded comment structure
- ✅ Max depth: 10 levels
- ✅ Parent-child relationships
- ✅ Reply to comments
- ✅ Comment tree view
- ✅ Parent chain (breadcrumb)

**Comment Features**:
- ✅ Create comments
- ✅ Update own comments
- ✅ Delete own comments (soft delete)
- ✅ Vote on comments
- ✅ View comments by post
- ✅ View comments by author
- ✅ Pagination support

**Comment Display**:
- ✅ Show author info
- ✅ Show vote count
- ✅ Show reply count
- ✅ Show depth level
- ✅ Deleted state ([deleted])

### 3.3 Crosspost System

**Features**:
- ✅ Share existing post to own feed
- ✅ Add custom title to crosspost
- ✅ Add custom tags
- ✅ Track original post (sourcePostId)
- ✅ View all crossposts of a post
- ✅ Prevent recursive crossposts

**Use Cases**:
- Share interesting content
- Introduce posts to different audiences
- Boost visibility of good content

---

## 4. Voting System

### 4.1 Voting Mechanism

**Vote Types**:
- ✅ Upvote (+1)
- ✅ Downvote (-1)
- ✅ Remove vote (0)
- ✅ Change vote (up → down or down → up)

**Vote Targets**:
- ✅ Posts
- ✅ Comments

**Vote Features**:
- ✅ One vote per user per target
- ✅ Real-time vote count update
- ✅ User vote status (up/down/null)
- ✅ Vote history (user's votes)

**Karma Impact**:
- ✅ Upvote increases author's karma
- ✅ Downvote decreases author's karma
- ✅ Remove vote reverts karma change

---

## 5. Social Features

### 5.1 Follow System

**Follow Features**:
- ✅ Follow users
- ✅ Unfollow users
- ✅ View followers list
- ✅ View following list
- ✅ Check follow status
- ✅ Get mutual follows
- ✅ Follower/following count

**Follow Impact**:
- ✅ Personalized feed (see posts from followed users)
- ✅ Notifications on follow
- ✅ Social graph building

### 5.2 Saved Posts

**Features**:
- ✅ Save posts for later
- ✅ Unsave posts
- ✅ View saved posts
- ✅ Check save status
- ✅ Personal bookmark collection

**Use Cases**:
- Read later functionality
- Personal content curation
- Reference collection

---

## 6. Tag System

### 6.1 Tag Management

**Tag Features**:
- ✅ Auto-create tags on post creation
- ✅ Lowercase normalization
- ✅ Max 50 characters per tag
- ✅ Multi-tag support per post
- ✅ Tag popularity tracking (post count)

**Tag Discovery**:
- ✅ List all tags
- ✅ Get popular tags
- ✅ Search tags
- ✅ Get tag by name/ID
- ✅ View posts by tag

**Tag Display**:
- ✅ Tag name
- ✅ Post count
- ✅ Creation date

---

## 7. Search System

### 7.1 Universal Search

**Search Capabilities**:
- ✅ Search posts (title + content)
- ✅ Search users (username + display name)
- ✅ Search tags (tag name)
- ✅ Universal search (all types)

**Search Features**:
- ✅ Full-text search
- ✅ Pagination
- ✅ Result count per type
- ✅ Relevance-based results

### 7.2 Search History

**Features**:
- ✅ Track user searches
- ✅ View search history
- ✅ Clear all history
- ✅ Delete individual search
- ✅ Popular searches (global)

**Privacy**:
- ✅ Per-user history (not shared)
- ✅ User-controlled deletion

---

## 8. Media Management

### 8.1 Media Upload

**Upload Features**:
- ✅ Upload images (JPG, PNG, GIF, WebP)
- ✅ Upload videos (MP4, WebM, MOV)
- ✅ Max file size: 300 MB
- ✅ Generate unique filename (UUID)
- ✅ Store on Bunny CDN
- ✅ Track metadata (dimensions, size, MIME type)

**Media Processing**:
- ✅ Extract dimensions (width, height)
- ✅ Generate thumbnail (for videos)
- ✅ Calculate file size
- ✅ Store duration (for videos)

### 8.2 Media Management

**Features**:
- ✅ View media details
- ✅ View user's media gallery
- ✅ Attach media to posts
- ✅ Delete media
- ✅ Track media usage count
- ✅ Prevent deletion if in use

**Storage**:
- ✅ Bunny CDN integration
- ✅ Direct CDN URLs
- ✅ Fast global delivery
- ✅ Cost-effective storage

---

## 9. Notification System

### 9.1 In-App Notifications

**Notification Types**:
- ✅ Reply: Someone replied to your post/comment
- ✅ Vote: Someone voted on your content
- ✅ Mention: Someone mentioned you (planned)
- ✅ Follow: Someone followed you

**Notification Features**:
- ✅ Real-time notification delivery
- ✅ Unread count badge
- ✅ Mark as read (single/all)
- ✅ Delete notifications (single/all)
- ✅ Notification list (paginated)
- ✅ Unread filter

**Notification Data**:
- ✅ Sender information
- ✅ Message text
- ✅ Related post/comment
- ✅ Timestamp
- ✅ Read status

### 9.2 Notification Settings

**User Preferences**:
- ✅ Toggle replies notifications
- ✅ Toggle mentions notifications
- ✅ Toggle votes notifications
- ✅ Toggle follows notifications
- ✅ Email notifications (planned)

**Privacy Control**:
- ✅ Per-notification-type control
- ✅ User-controlled preferences

### 9.3 Web Push Notifications

**Push Features**:
- ✅ VAPID-based Web Push
- ✅ Browser notifications
- ✅ Multiple device support
- ✅ Subscribe/unsubscribe
- ✅ Background notifications

**Push Events**:
- ✅ New reply notification
- ✅ New vote notification
- ✅ New follower notification
- ✅ Custom notification messages

---

## 10. Real-time Features

### 10.1 WebSocket System

**Connection Features**:
- ✅ Authenticated WebSocket connections
- ✅ Anonymous connections (with UUID)
- ✅ Room-based messaging
- ✅ Heartbeat mechanism (30s ping/pong)
- ✅ Auto-reconnection support

**Real-time Events**:
- ✅ Message delivery
- ✅ Online/offline status
- ✅ Notification delivery (planned)
- ✅ Live vote updates (planned)
- ✅ Live comment updates (planned)

**WebSocket Manager**:
- ✅ Client registration/unregistration
- ✅ Broadcast to all clients
- ✅ Broadcast to room
- ✅ Send to specific user
- ✅ Connection tracking

---

## 11. Sorting & Ranking

### 11.1 Hot Algorithm

**Hot Score Calculation**:
```
score = (votes + 1) / (age_in_hours + 2)^1.5
```

**Features**:
- ✅ Time decay factor
- ✅ Vote weight
- ✅ Trending content discovery

### 11.2 Top Sorting

**Time Filters**:
- ✅ Hour: Top posts in last hour
- ✅ Day: Top posts in last 24 hours
- ✅ Week: Top posts in last 7 days
- ✅ Month: Top posts in last 30 days
- ✅ Year: Top posts in last 365 days
- ✅ All Time: Top posts ever

### 11.3 New Sorting

**Features**:
- ✅ Latest first (created_at DESC)
- ✅ Real-time updates

---

## 12. Personalization

### 12.1 Personalized Feed

**Feed Algorithm**:
- ✅ Posts from followed users
- ✅ Sorted by recency
- ✅ Optional relevance scoring (planned)

**Features**:
- ✅ Authenticated users only
- ✅ Pagination support
- ✅ Real-time updates (via WebSocket)

### 12.2 Recommendations (Planned)

**Planned Features**:
- 🔄 Recommended posts (based on interests)
- 🔄 Recommended users (based on follows)
- 🔄 Recommended tags (based on views)
- 🔄 Similar posts

---

## 13. Admin Features

### 13.1 User Management

**Admin Capabilities**:
- ✅ List all users
- ✅ View user details
- 🔄 Ban/unban users (planned)
- 🔄 Delete users (planned)
- 🔄 Change user roles (planned)

### 13.2 Content Moderation (Planned)

**Planned Features**:
- 🔄 Remove posts/comments
- 🔄 Hide content
- 🔄 Ban content
- 🔄 Moderation queue
- 🔄 Report system
- 🔄 Automated moderation (AI)

---

## 14. Performance Features

### 14.1 Caching Strategy

**Redis Caching**:
- ✅ OAuth code storage (5-min TTL)
- ✅ Session management
- ✅ Unread counts caching (planned)
- ✅ Online status tracking (planned)

**Database Optimizations**:
- ✅ Indexed queries
- ✅ Eager loading (prevent N+1)
- ✅ Select specific fields
- ✅ Denormalized counts

### 14.2 Pagination

**Features**:
- ✅ Cursor-based pagination (planned for chat)
- ✅ Offset-based pagination (current)
- ✅ Configurable page size
- ✅ Total count tracking

---

## 15. Security Features

### 15.1 Authentication Security

**Implemented**:
- ✅ Password hashing (bcrypt)
- ✅ JWT tokens
- ✅ Token expiry
- ✅ Secure OAuth flow
- ✅ CORS configuration

**Planned**:
- 🔄 Rate limiting
- 🔄 CSRF protection
- 🔄 IP blocking
- 🔄 2FA support

### 15.2 Data Security

**Implemented**:
- ✅ Input validation (go-playground/validator)
- ✅ SQL injection prevention (GORM)
- ✅ XSS prevention (input sanitization)
- ✅ Soft delete (data recovery)

**Planned**:
- 🔄 GDPR compliance
- 🔄 Data export
- 🔄 Right to be forgotten

---

## 16. Developer Features

### 16.1 API Documentation

**Documentation**:
- ✅ Backend API spec (backend_spec/)
- ✅ Chat API spec (chat_api_spec/)
- ✅ Postman collection (planned)
- ✅ OpenAPI/Swagger (planned)

### 16.2 Development Tools

**Features**:
- ✅ Comprehensive error handling
- ✅ Logging middleware
- ✅ Health check endpoint
- ✅ Hot reload (via Air, planned)
- ✅ Database migrations

---

## 17. Chat System (In Development)

### 17.1 Direct Messaging (Phase 1 MVP)

**Planned Features**:
- 🔄 1-on-1 direct messaging
- 🔄 Real-time message delivery (WebSocket)
- 🔄 Message history
- 🔄 Unread count tracking
- 🔄 Read receipts
- 🔄 Online/offline status
- 🔄 Last seen tracking
- 🔄 Typing indicators

**Database Schema**:
- 🔄 conversations table
- 🔄 messages table
- 🔄 blocks table

**Performance Targets**:
- Get conversations: < 100ms
- Get messages: < 100ms
- Send message: < 50ms

### 17.2 Group Chat (Future)

**Planned Features**:
- 🔄 Create group chats
- 🔄 Add/remove members
- 🔄 Group admin roles
- 🔄 Group settings

---

## 18. Future Enhancements

### 18.1 Content Features

**Planned**:
- 🔄 Polls & surveys
- 🔄 Live streaming
- 🔄 Stories (24-hour posts)
- 🔄 Scheduled posts
- 🔄 Draft posts

### 18.2 Social Features

**Planned**:
- 🔄 User blocking
- 🔄 Mute users
- 🔄 Private profiles
- 🔄 Verified accounts
- 🔄 User badges

### 18.3 Discovery Features

**Planned**:
- 🔄 Trending topics
- 🔄 Featured posts
- 🔄 Explore page
- 🔄 Categories/Communities
- 🔄 Hashtag following

### 18.4 Monetization (Optional)

**Planned**:
- 🔄 Premium membership
- 🔄 Tipping system
- 🔄 Creator monetization
- 🔄 Ad platform

---

## Feature Summary

### Implemented Features (✅)

**Core**:
- Authentication (Email/Password + Google OAuth)
- User profiles & management
- Posts (create, read, update, delete)
- Nested comments (10 levels)
- Voting system (posts & comments)
- Karma system

**Social**:
- Follow/unfollow users
- Saved posts
- Personalized feed
- Tag system
- Universal search

**Media**:
- Image upload (300MB max)
- Video upload (300MB max)
- Bunny CDN integration
- Media gallery

**Notifications**:
- In-app notifications
- Web Push notifications
- Notification settings
- Unread count

**Real-time**:
- WebSocket connections
- Room-based messaging
- Online status tracking
- Heartbeat mechanism

### In Development (🔄)

- Direct messaging (1-on-1 chat)
- Group chat
- Rate limiting
- Advanced moderation tools
- Admin dashboard

### Planned for Future

- 2FA authentication
- User blocking
- Polls & surveys
- Live streaming
- Trending algorithm improvements
- Full-text search optimization
- Email notifications
- Mobile app (React Native)

---

## Technical Highlights

1. **Clean Architecture**: 4-layer separation for maintainability
2. **Dependency Injection**: Custom DI container
3. **Repository Pattern**: Abstracted data access
4. **Service Layer**: Centralized business logic
5. **Middleware Chain**: Composable request processing
6. **Error Handling**: Comprehensive error management
7. **Testing Ready**: Interface-based design for easy mocking
8. **Scalable**: Redis caching, database indexing
9. **Secure**: JWT, bcrypt, CORS, input validation
10. **Real-time**: WebSocket for live features

---

## Performance Metrics (Target)

- **API Response Time**: < 100ms (average)
- **WebSocket Latency**: < 50ms
- **Media Upload**: 300MB max file size
- **Database Queries**: Optimized with indexes
- **Concurrent Users**: 1000+ simultaneous connections
- **Uptime**: 99.9% (planned)
