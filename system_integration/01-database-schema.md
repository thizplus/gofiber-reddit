# Database Schema Implementation Guide

## 📊 Overview

ไฟล์นี้อธิบาย database schema ทั้งหมดที่ต้องสร้างและแก้ไขสำหรับ Social Media Platform

---

## 🔄 Migration Strategy

### Option 1: Fresh Start (Recommended for Development)
```bash
# Drop existing database
DROP DATABASE gofiber_template;
CREATE DATABASE gofiber_template;

# Run new migrations
go run cmd/api/main.go
```

### Option 2: Incremental Migration (Production)
```bash
# Create migration files
# Apply migrations one by one
# Keep existing data intact
```

---

## 🗄️ Complete Schema

### 1. Users Table (Enhanced)

**Purpose:** Store user accounts and profile information

```go
// File: domain/models/user.go
type User struct {
    // Core Fields
    ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Username  string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`

    // Profile Fields (NEW)
    DisplayName string `gorm:"not null"`
    Avatar      string
    Bio         string `gorm:"type:text"`
    CoverImage  string
    Location    string
    Website     string

    // Social Stats (NEW)
    Karma           int  `gorm:"default:0"`
    FollowersCount  int  `gorm:"default:0"`
    FollowingCount  int  `gorm:"default:0"`

    // Status
    Role     string `gorm:"default:'user'"` // user, admin
    IsActive bool   `gorm:"default:true"`

    // Timestamps
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (User) TableName() string { return "users" }
```

**SQL:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    avatar VARCHAR(500),
    bio TEXT,
    cover_image VARCHAR(500),
    location VARCHAR(100),
    website VARCHAR(255),
    karma INTEGER DEFAULT 0,
    followers_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0,
    role VARCHAR(20) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_karma ON users(karma DESC);
```

---

### 2. Posts Table (NEW)

**Purpose:** Store user posts and crossposts

```go
// File: domain/models/post.go
type Post struct {
    ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Title    string    `gorm:"not null;type:varchar(300)"`
    Content  string    `gorm:"not null;type:text"`

    // Author
    AuthorID uuid.UUID `gorm:"not null;index"`
    Author   User      `gorm:"foreignKey:AuthorID"`

    // Stats
    Votes        int `gorm:"default:0"`
    CommentCount int `gorm:"default:0"`

    // Crosspost (optional)
    SourcePostID *uuid.UUID `gorm:"index"`
    SourcePost   *Post      `gorm:"foreignKey:SourcePostID"`

    // Media & Tags (relationships)
    Media []Media `gorm:"many2many:post_media;"`
    Tags  []Tag   `gorm:"many2many:post_tags;"`

    // Status
    IsDeleted bool `gorm:"default:false;index"`

    // Timestamps
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `gorm:"index"`
}

func (Post) TableName() string { return "posts" }
```

**SQL:**
```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(300) NOT NULL,
    content TEXT NOT NULL,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    votes INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    source_post_id UUID REFERENCES posts(id) ON DELETE SET NULL,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_source_post_id ON posts(source_post_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_votes ON posts(votes DESC);
CREATE INDEX idx_posts_is_deleted ON posts(is_deleted);

-- Full-text search index
CREATE INDEX idx_posts_title_search ON posts USING gin(to_tsvector('english', title));
CREATE INDEX idx_posts_content_search ON posts USING gin(to_tsvector('english', content));
```

---

### 3. Comments Table (NEW)

**Purpose:** Store comments and nested replies

```go
// File: domain/models/comment.go
type Comment struct {
    ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    PostID  uuid.UUID `gorm:"not null;index"`
    Post    Post      `gorm:"foreignKey:PostID"`

    AuthorID uuid.UUID `gorm:"not null;index"`
    Author   User      `gorm:"foreignKey:AuthorID"`

    Content string `gorm:"not null;type:text"`
    Votes   int    `gorm:"default:0"`

    // Nested replies
    ParentID *uuid.UUID `gorm:"index"`
    Parent   *Comment   `gorm:"foreignKey:ParentID"`
    Replies  []Comment  `gorm:"foreignKey:ParentID"`
    Depth    int        `gorm:"default:0;index"` // 0 = top-level, max 10

    // Status
    IsDeleted bool `gorm:"default:false"`

    // Timestamps
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}

func (Comment) TableName() string { return "comments" }
```

**SQL:**
```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    votes INTEGER DEFAULT 0,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    depth INTEGER DEFAULT 0,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_depth ON comments(depth);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);
CREATE INDEX idx_comments_votes ON comments(votes DESC);

-- Constraint: max depth 10
ALTER TABLE comments ADD CONSTRAINT chk_comments_depth CHECK (depth >= 0 AND depth <= 10);
```

---

### 4. Media Table (Enhanced for Bunny Storage)

**Purpose:** Store media files metadata (Bunny CDN URLs)

```go
// File: domain/models/media.go
type Media struct {
    ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID   uuid.UUID `gorm:"not null;index"`
    User     User      `gorm:"foreignKey:UserID"`

    // File info
    Type     string `gorm:"not null"` // image, video
    FileName string `gorm:"not null"`
    MimeType string `gorm:"not null"`
    Size     int64  `gorm:"not null"` // bytes

    // URLs (Bunny CDN)
    URL       string `gorm:"not null"` // Full CDN URL
    Thumbnail string                   // Thumbnail URL

    // Dimensions
    Width  int
    Height int

    // Video specific
    Duration float64 // seconds (for videos)

    // Usage tracking
    Posts    []Post    `gorm:"many2many:post_media;"`
    UsageCount int     `gorm:"default:0"`

    // Timestamps
    CreatedAt time.Time
}

func (Media) TableName() string { return "media" }
```

**SQL:**
```sql
CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- 'image' or 'video'
    file_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    url VARCHAR(500) NOT NULL,
    thumbnail VARCHAR(500),
    width INTEGER,
    height INTEGER,
    duration DECIMAL(10,2), -- for videos
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_media_user_id ON media(user_id);
CREATE INDEX idx_media_type ON media(type);
CREATE INDEX idx_media_created_at ON media(created_at DESC);
```

---

### 5. Post_Media Junction Table (NEW)

**Purpose:** Many-to-many relationship between posts and media

```go
// File: domain/models/post_media.go
type PostMedia struct {
    PostID  uuid.UUID `gorm:"primaryKey"`
    MediaID uuid.UUID `gorm:"primaryKey"`
    Order   int       `gorm:"default:0"` // Display order
}

func (PostMedia) TableName() string { return "post_media" }
```

**SQL:**
```sql
CREATE TABLE post_media (
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    "order" INTEGER DEFAULT 0,
    PRIMARY KEY (post_id, media_id)
);

CREATE INDEX idx_post_media_post_id ON post_media(post_id);
CREATE INDEX idx_post_media_media_id ON post_media(media_id);
```

---

### 6. Votes Table (NEW)

**Purpose:** Store upvotes/downvotes for posts and comments (polymorphic)

```go
// File: domain/models/vote.go
type Vote struct {
    UserID     uuid.UUID `gorm:"primaryKey"`
    User       User      `gorm:"foreignKey:UserID"`

    TargetID   uuid.UUID `gorm:"primaryKey;index"` // post_id or comment_id
    TargetType string    `gorm:"primaryKey"`        // 'post' or 'comment'

    VoteType   string    `gorm:"not null"` // 'up' or 'down'

    CreatedAt  time.Time
}

func (Vote) TableName() string { return "votes" }
```

**SQL:**
```sql
CREATE TABLE votes (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_id UUID NOT NULL,
    target_type VARCHAR(20) NOT NULL, -- 'post' or 'comment'
    vote_type VARCHAR(10) NOT NULL, -- 'up' or 'down'
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, target_id, target_type)
);

CREATE INDEX idx_votes_target ON votes(target_id, target_type);
CREATE INDEX idx_votes_user_id ON votes(user_id);

-- Constraint: vote_type must be 'up' or 'down'
ALTER TABLE votes ADD CONSTRAINT chk_votes_type CHECK (vote_type IN ('up', 'down'));

-- Constraint: target_type must be 'post' or 'comment'
ALTER TABLE votes ADD CONSTRAINT chk_votes_target_type CHECK (target_type IN ('post', 'comment'));
```

---

### 7. Follows Table (NEW)

**Purpose:** Store follow relationships between users

```go
// File: domain/models/follow.go
type Follow struct {
    FollowerID  uuid.UUID `gorm:"primaryKey"`
    Follower    User      `gorm:"foreignKey:FollowerID"`

    FollowingID uuid.UUID `gorm:"primaryKey"`
    Following   User      `gorm:"foreignKey:FollowingID"`

    CreatedAt   time.Time
}

func (Follow) TableName() string { return "follows" }
```

**SQL:**
```sql
CREATE TABLE follows (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id) -- Can't follow yourself
);

CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);
```

---

### 8. Saved_Posts Table (NEW)

**Purpose:** Store saved/bookmarked posts per user

```go
// File: domain/models/saved_post.go
type SavedPost struct {
    UserID  uuid.UUID `gorm:"primaryKey"`
    User    User      `gorm:"foreignKey:UserID"`

    PostID  uuid.UUID `gorm:"primaryKey"`
    Post    Post      `gorm:"foreignKey:PostID"`

    SavedAt time.Time `gorm:"index"`
}

func (SavedPost) TableName() string { return "saved_posts" }
```

**SQL:**
```sql
CREATE TABLE saved_posts (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    saved_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

CREATE INDEX idx_saved_posts_user_id ON saved_posts(user_id);
CREATE INDEX idx_saved_posts_saved_at ON saved_posts(saved_at DESC);
```

---

### 9. Notifications Table (NEW)

**Purpose:** Store user notifications

```go
// File: domain/models/notification.go
type Notification struct {
    ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

    UserID   uuid.UUID `gorm:"not null;index"` // Recipient
    User     User      `gorm:"foreignKey:UserID"`

    SenderID uuid.UUID `gorm:"not null"` // Who triggered notification
    Sender   User      `gorm:"foreignKey:SenderID"`

    Type     string `gorm:"not null;index"` // reply, vote, mention, follow
    Message  string `gorm:"not null"`

    // Optional references
    PostID    *uuid.UUID
    Post      *Post `gorm:"foreignKey:PostID"`
    CommentID *uuid.UUID
    Comment   *Comment `gorm:"foreignKey:CommentID"`

    IsRead    bool `gorm:"default:false;index"`
    CreatedAt time.Time `gorm:"index"`
}

func (Notification) TableName() string { return "notifications" }
```

**SQL:**
```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- reply, vote, mention, follow
    message TEXT NOT NULL,
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);

-- Constraint: type must be valid
ALTER TABLE notifications ADD CONSTRAINT chk_notifications_type
    CHECK (type IN ('reply', 'vote', 'mention', 'follow'));
```

---

### 10. Notification_Settings Table (NEW)

**Purpose:** Store user's notification preferences

```go
// File: domain/models/notification_settings.go
type NotificationSettings struct {
    UserID             uuid.UUID `gorm:"primaryKey"`
    User               User      `gorm:"foreignKey:UserID"`

    Replies            bool `gorm:"default:true"`
    Mentions           bool `gorm:"default:true"`
    Votes              bool `gorm:"default:false"`
    Follows            bool `gorm:"default:true"`
    EmailNotifications bool `gorm:"default:false"`

    UpdatedAt time.Time
}

func (NotificationSettings) TableName() string { return "notification_settings" }
```

**SQL:**
```sql
CREATE TABLE notification_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    replies BOOLEAN DEFAULT true,
    mentions BOOLEAN DEFAULT true,
    votes BOOLEAN DEFAULT false,
    follows BOOLEAN DEFAULT true,
    email_notifications BOOLEAN DEFAULT false,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

### 11. Tags Table (NEW)

**Purpose:** Store post tags

```go
// File: domain/models/tag.go
type Tag struct {
    ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name      string    `gorm:"uniqueIndex;not null"`
    PostCount int       `gorm:"default:0"`

    Posts     []Post    `gorm:"many2many:post_tags;"`

    CreatedAt time.Time
}

func (Tag) TableName() string { return "tags" }
```

**SQL:**
```sql
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    post_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_tags_post_count ON tags(post_count DESC);
```

---

### 12. Post_Tags Junction Table (NEW)

**Purpose:** Many-to-many relationship between posts and tags

```sql
CREATE TABLE post_tags (
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);
```

---

### 13. Search_History Table (NEW - Optional)

**Purpose:** Store user's search history

```go
// File: domain/models/search_history.go
type SearchHistory struct {
    ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID     uuid.UUID `gorm:"not null;index"`
    User       User      `gorm:"foreignKey:UserID"`

    Query      string `gorm:"not null"`
    Type       string // 'posts', 'users', 'all'

    SearchedAt time.Time `gorm:"index"`
}

func (SearchHistory) TableName() string { return "search_history" }
```

**SQL:**
```sql
CREATE TABLE search_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    query VARCHAR(200) NOT NULL,
    type VARCHAR(20),
    searched_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_search_history_user_id ON search_history(user_id);
CREATE INDEX idx_search_history_searched_at ON search_history(searched_at DESC);

-- Auto-delete old history (30 days)
-- Can be handled by scheduled job
```

---

## 🔧 GORM Auto-Migrate Code

**File:** `infrastructure/postgres/database.go`

```go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        // Enhanced existing
        &models.User{},

        // New core models
        &models.Post{},
        &models.Comment{},
        &models.Media{},

        // Voting & Social
        &models.Vote{},
        &models.Follow{},
        &models.SavedPost{},

        // Notifications
        &models.Notification{},
        &models.NotificationSettings{},

        // Tags
        &models.Tag{},

        // Search
        &models.SearchHistory{},
    )
}
```

---

## 📈 Database Indexes Summary

### Critical Indexes (Must Have)
```sql
-- User indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_karma ON users(karma DESC);

-- Post indexes
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_votes ON posts(votes DESC);
CREATE INDEX idx_posts_is_deleted ON posts(is_deleted);

-- Comment indexes
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);

-- Vote indexes
CREATE INDEX idx_votes_target ON votes(target_id, target_type);
CREATE INDEX idx_votes_user_id ON votes(user_id);

-- Full-text search (PostgreSQL)
CREATE INDEX idx_posts_title_search ON posts USING gin(to_tsvector('english', title));
CREATE INDEX idx_posts_content_search ON posts USING gin(to_tsvector('english', content));
```

---

## 🔄 Migration Checklist

- [ ] Backup existing database
- [ ] Update User model with new fields
- [ ] Create all new tables (Posts, Comments, Votes, etc.)
- [ ] Create junction tables (PostMedia, PostTags)
- [ ] Add foreign key constraints
- [ ] Create indexes for performance
- [ ] Add check constraints for data integrity
- [ ] Create full-text search indexes
- [ ] Test migrations on development database
- [ ] Run GORM AutoMigrate
- [ ] Verify all relationships work
- [ ] Seed test data

---

## ✅ Next Steps

After completing database schema:
1. ✅ Proceed to `02-implementation-phases.md`
2. ✅ Implement Repository interfaces
3. ✅ Create Service layer
4. ✅ Build API handlers

---

**Schema Complete? → Proceed to `02-implementation-phases.md`**
