# Database Schema

## ภาพรวม

**ฐานข้อมูล**: PostgreSQL 15
**ORM**: GORM v1.25.6
**จำนวน Tables**: 15 tables
**Migration**: Auto-migration + SQL files

## Database Models (ทั้งหมด)

### 1. Users (ผู้ใช้)

**ตาราง**: `users`

```go
type User struct {
    ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
    Email             string    `gorm:"uniqueIndex;not null"`
    Username          string    `gorm:"uniqueIndex;not null"`
    Password          string    `gorm:"not null"` // bcrypt hashed

    // OAuth fields
    OAuthProvider     string    `gorm:"index"` // "google", etc.
    OAuthID           string    `gorm:"index"`
    IsOAuthUser       bool      `gorm:"default:false"`

    // Profile fields
    DisplayName       string
    Avatar            string
    Bio               string
    Location          string
    Website           string

    // Social stats
    Karma             int       `gorm:"default:0;index"`
    FollowersCount    int       `gorm:"default:0"`
    FollowingCount    int       `gorm:"default:0"`

    // Account status
    Role              string    `gorm:"default:'user'"` // "user" or "admin"
    IsActive          bool      `gorm:"default:true"`

    // Relationships
    Posts             []Post    `gorm:"foreignKey:AuthorID"`
    Comments          []Comment `gorm:"foreignKey:AuthorID"`
    Followers         []Follow  `gorm:"foreignKey:FollowingID"`
    Following         []Follow  `gorm:"foreignKey:FollowerID"`

    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

**Indexes**:
- `email` (unique)
- `username` (unique)
- `oauth_provider`, `oauth_id`
- `karma` (for leaderboard)

**Use Cases**:
- Authentication (email/password or OAuth)
- Profile display
- User discovery
- Karma leaderboard

---

### 2. Posts (โพสต์)

**ตาราง**: `posts`

```go
type Post struct {
    ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
    Title           string    `gorm:"size:300;not null;index"`
    Content         string    `gorm:"type:text"`
    AuthorID        uuid.UUID `gorm:"type:uuid;not null;index"`

    // Stats
    Votes           int       `gorm:"default:0;index"`
    CommentCount    int       `gorm:"default:0"`

    // Crosspost support
    SourcePostID    *uuid.UUID `gorm:"type:uuid;index"` // null = original post

    // Relationships
    Author          User      `gorm:"foreignKey:AuthorID"`
    SourcePost      *Post     `gorm:"foreignKey:SourcePostID"`
    Comments        []Comment `gorm:"foreignKey:PostID"`
    Media           []Media   `gorm:"many2many:post_media;"`
    Tags            []Tag     `gorm:"many2many:post_tags;"`

    // Soft delete
    IsDeleted       bool      `gorm:"default:false;index"`
    DeletedAt       *time.Time

    CreatedAt       time.Time `gorm:"index"`
    UpdatedAt       time.Time
}
```

**Indexes**:
- `title` (for search)
- `author_id` (for author's posts)
- `votes` (for sorting)
- `is_deleted` (for filtering)
- `created_at` (for sorting)
- `source_post_id` (for crossposts)

**Use Cases**:
- Create/read/update/delete posts
- Vote on posts
- Search posts
- Get posts by author/tag
- Crosspost system

---

### 3. Comments (คอมเมนต์)

**ตาราง**: `comments`

```go
type Comment struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
    PostID      uuid.UUID `gorm:"type:uuid;not null;index"`
    AuthorID    uuid.UUID `gorm:"type:uuid;not null;index"`
    Content     string    `gorm:"type:text;not null"`

    // Stats
    Votes       int       `gorm:"default:0;index"`

    // Nested structure
    ParentID    *uuid.UUID `gorm:"type:uuid;index"` // null = top-level
    Depth       int        `gorm:"default:0"` // max 10

    // Relationships
    Post        Post      `gorm:"foreignKey:PostID"`
    Author      User      `gorm:"foreignKey:AuthorID"`
    Parent      *Comment  `gorm:"foreignKey:ParentID"`
    Replies     []Comment `gorm:"foreignKey:ParentID"`

    // Soft delete
    IsDeleted   bool      `gorm:"default:false;index"`
    DeletedAt   *time.Time

    CreatedAt   time.Time `gorm:"index"`
    UpdatedAt   time.Time
}
```

**Indexes**:
- `post_id` (for post's comments)
- `author_id` (for author's comments)
- `votes` (for sorting)
- `parent_id` (for nested comments)
- `is_deleted` (for filtering)
- `created_at` (for sorting)

**Features**:
- Nested comments (max 10 levels)
- Vote on comments
- Soft delete (แสดงว่า [deleted])

---

### 4. Votes (การโหวต)

**ตาราง**: `votes`

```go
type Vote struct {
    UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
    TargetID    uuid.UUID `gorm:"type:uuid;primaryKey"`
    TargetType  string    `gorm:"primaryKey"` // "post" or "comment"
    VoteType    string    `gorm:"not null"`   // "up" or "down"

    CreatedAt   time.Time `gorm:"index"`
}
```

**Composite Primary Key**: `(user_id, target_id, target_type)`

**Indexes**:
- Composite PK (auto-indexed)
- `created_at` (for analytics)

**Business Rules**:
- 1 user = 1 vote per target
- Can change vote (up → down หรือ down → up)
- Delete vote = remove vote

---

### 5. Follows (การติดตาม)

**ตาราง**: `follows`

```go
type Follow struct {
    FollowerID  uuid.UUID `gorm:"type:uuid;primaryKey"`
    FollowingID uuid.UUID `gorm:"type:uuid;primaryKey"`

    Follower    User      `gorm:"foreignKey:FollowerID"`
    Following   User      `gorm:"foreignKey:FollowingID"`

    CreatedAt   time.Time `gorm:"index"`
}
```

**Composite Primary Key**: `(follower_id, following_id)`

**Indexes**:
- `follower_id` (for "following list")
- `following_id` (for "followers list")
- `created_at` (for sorting)

**Use Cases**:
- Follow/unfollow users
- Get followers/following lists
- Check follow status
- Personalized feed

---

### 6. Saved Posts (โพสต์ที่บันทึก)

**ตาราง**: `saved_posts`

```go
type SavedPost struct {
    UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
    PostID    uuid.UUID `gorm:"type:uuid;primaryKey"`

    User      User      `gorm:"foreignKey:UserID"`
    Post      Post      `gorm:"foreignKey:PostID"`

    SavedAt   time.Time `gorm:"index"`
}
```

**Composite Primary Key**: `(user_id, post_id)`

**Use Cases**:
- Save posts for later
- Bookmark functionality
- Personal collection

---

### 7. Notifications (การแจ้งเตือน)

**ตาราง**: `notifications`

```go
type Notification struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
    UserID      uuid.UUID `gorm:"type:uuid;not null;index"` // recipient
    SenderID    uuid.UUID `gorm:"type:uuid"`

    Type        string    `gorm:"not null;index"` // "reply", "vote", "mention", "follow"
    Message     string    `gorm:"not null"`

    // Optional references
    PostID      *uuid.UUID `gorm:"type:uuid"`
    CommentID   *uuid.UUID `gorm:"type:uuid"`

    // Status
    IsRead      bool      `gorm:"default:false;index"`

    // Relationships
    User        User      `gorm:"foreignKey:UserID"`
    Sender      User      `gorm:"foreignKey:SenderID"`
    Post        *Post     `gorm:"foreignKey:PostID"`
    Comment     *Comment  `gorm:"foreignKey:CommentID"`

    CreatedAt   time.Time `gorm:"index"`
}
```

**Indexes**:
- `user_id` (for user's notifications)
- `type` (for filtering)
- `is_read` (for unread count)
- `created_at` (for sorting)

**Notification Types**:
- `reply`: Someone replied to your post/comment
- `vote`: Someone voted on your post/comment
- `mention`: Someone mentioned you
- `follow`: Someone followed you

---

### 8. Notification Settings (การตั้งค่าการแจ้งเตือน)

**ตาราง**: `notification_settings`

```go
type NotificationSettings struct {
    UserID              uuid.UUID `gorm:"type:uuid;primaryKey"`

    // In-app notifications
    Replies             bool      `gorm:"default:true"`
    Mentions            bool      `gorm:"default:true"`
    Votes               bool      `gorm:"default:true"`
    Follows             bool      `gorm:"default:true"`

    // Email notifications (planned)
    EmailNotifications  bool      `gorm:"default:false"`

    User                User      `gorm:"foreignKey:UserID"`

    UpdatedAt           time.Time
}
```

**Use Cases**:
- User preferences for notifications
- Control what notifications to receive

---

### 9. Push Subscriptions (การสมัครรับ Push)

**ตาราง**: `push_subscriptions`

```go
type PushSubscription struct {
    ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
    UserID          uuid.UUID `gorm:"type:uuid;not null;index"`

    // Web Push fields (VAPID)
    Endpoint        string    `gorm:"not null"`
    P256dh          string    `gorm:"not null"` // Public key
    Auth            string    `gorm:"not null"` // Auth secret

    ExpirationTime  *time.Time

    User            User      `gorm:"foreignKey:UserID"`

    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

**Indexes**:
- `user_id` (ผู้ใช้อาจมีหลาย device)
- `endpoint` (unique per device)

**Use Cases**:
- Subscribe to push notifications
- Send push notifications to user
- Manage device subscriptions

---

### 10. Media (ไฟล์มีเดีย)

**ตาราง**: `media`

```go
type Media struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
    UserID      uuid.UUID `gorm:"type:uuid;not null;index"`

    Type        string    `gorm:"not null;index"` // "image" or "video"
    FileName    string    `gorm:"not null"`
    MimeType    string    `gorm:"not null"`
    Size        int64     `gorm:"not null"`

    // Storage
    URL         string    `gorm:"not null"` // Bunny CDN URL
    Thumbnail   string    // Thumbnail URL (for videos)

    // Metadata
    Width       int
    Height      int
    Duration    int       // For videos (seconds)

    // Relationships
    User        User      `gorm:"foreignKey:UserID"`
    Posts       []Post    `gorm:"many2many:post_media;"`

    // Usage tracking
    UsageCount  int       `gorm:"default:0"`

    CreatedAt   time.Time `gorm:"index"`
}
```

**Indexes**:
- `user_id` (for user's media)
- `type` (for filtering)
- `created_at` (for sorting)

**Supported Types**:
- `image`: JPG, PNG, GIF, WebP
- `video`: MP4, WebM, MOV

**Max Size**: 300 MB per file

---

### 11. Tags (แท็ก)

**ตาราง**: `tags`

```go
type Tag struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
    Name        string    `gorm:"unique;not null;size:50"`
    PostCount   int       `gorm:"default:0;index"`

    Posts       []Post    `gorm:"many2many:post_tags;"`

    CreatedAt   time.Time
}
```

**Indexes**:
- `name` (unique, for lookup)
- `post_count` (for popular tags)

**Features**:
- Auto-create tags when used
- Lowercase normalization
- Track popularity

---

### 12. Search History (ประวัติการค้นหา)

**ตาราง**: `search_history`

```go
type SearchHistory struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
    UserID      uuid.UUID `gorm:"type:uuid;not null;index"`

    Query       string    `gorm:"not null"`
    Type        string    `gorm:"not null"` // "posts", "users", "all"

    User        User      `gorm:"foreignKey:UserID"`

    SearchedAt  time.Time `gorm:"index"`
}
```

**Indexes**:
- `user_id` (for user's search history)
- `searched_at` (for sorting)

**Use Cases**:
- Track user search behavior
- Show recent searches
- Analytics

---

### 13-15. Legacy Models (เก็บไว้เพื่อ backward compatibility)

**Task, File, Job**: เป็น models เดิมจาก template โปรเจกต์ ยังไม่ได้ใช้งานในระบบ social

## Many-to-Many Relationships

### 1. Post-Media (post_media)
```sql
CREATE TABLE post_media (
    post_id UUID REFERENCES posts(id),
    media_id UUID REFERENCES media(id),
    PRIMARY KEY (post_id, media_id)
);
```

### 2. Post-Tags (post_tags)
```sql
CREATE TABLE post_tags (
    post_id UUID REFERENCES posts(id),
    tag_id UUID REFERENCES tags(id),
    PRIMARY KEY (post_id, tag_id)
);
```

## Database Performance Optimizations

### 1. Indexes Strategy
- **Primary Keys**: UUID (กระจายข้อมูลดีกว่า auto-increment)
- **Foreign Keys**: ทุก FK มี index
- **Search Fields**: title, content (full-text search)
- **Sort Fields**: votes, created_at
- **Filter Fields**: is_deleted, is_read, type

### 2. Query Optimizations
- **Pagination**: Limit + Offset
- **Eager Loading**: Preload relationships ที่จำเป็น
- **Select Specific Fields**: หลีกเลี่ยง SELECT *
- **Count Caching**: เก็บ counts ไว้ใน model (FollowersCount, CommentCount)

### 3. Soft Delete Pattern
- ใช้ `is_deleted` flag แทน hard delete
- เก็บ `deleted_at` timestamp
- ฟื้นฟูข้อมูลได้

### 4. Denormalization
- **Karma**: เก็บไว้ใน user แทนนับจาก votes
- **FollowersCount, FollowingCount**: เก็บไว้ใน user
- **CommentCount**: เก็บไว้ใน post
- **PostCount**: เก็บไว้ใน tag

## Migrations

**Location**: `infrastructure/postgres/migrations/`

**Format**: `YYYYMMDDHHMMSS_description.sql`

**Example**:
```sql
-- 20240101000000_create_users_table.sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    -- ...
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

**Migration Strategy**:
- Auto-migration on startup (GORM)
- Manual SQL files for complex migrations
- Rollback support (planned)

## Backup Strategy (Recommended)

1. **Daily Backups**: pg_dump ทุกวัน
2. **Point-in-Time Recovery**: WAL archiving
3. **Replication**: Streaming replication (planned)
4. **Retention**: เก็บ 30 days

## Database Connection

**Configuration**:
```go
config := &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
    NowFunc: func() time.Time {
        return time.Now().UTC()
    },
}

dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
)

db, err := gorm.Open(postgres.Open(dsn), config)
```

**Connection Pool** (Default):
- MaxIdleConns: 10
- MaxOpenConns: 100
- ConnMaxLifetime: 1 hour

## ER Diagram (สรุป)

```
Users ──< Posts ──< Comments
  │       │          │
  │       └──< Votes │
  │                  │
  ├──< Votes ────────┘
  ├──< Follows
  ├──< SavedPosts ──> Posts
  ├──< Notifications
  ├──< NotificationSettings
  ├──< PushSubscriptions
  ├──< Media ──< PostMedia >──> Posts
  └──< SearchHistory

Tags ──< PostTags >──> Posts
```

## Future Enhancements

- [ ] Full-text search (PostgreSQL FTS)
- [ ] Read replicas (scalability)
- [ ] Partitioning (for large tables)
- [ ] Time-series data (for analytics)
- [ ] Graph database (for social graph)
