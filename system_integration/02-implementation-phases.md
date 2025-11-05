# Implementation Phases - Step by Step Guide

## 🎯 Overview

แผนการพัฒนาแบบ step-by-step สำหรับ Social Media Platform Backend
แบ่งเป็น 5 Phases ใช้เวลาประมาณ 5 สัปดาห์

---

## 📅 Timeline Overview

```
Week 1: Foundation (Database + Auth + Media)
Week 2: Core Features (Posts + Comments + Votes)
Week 3: Social Features (Follow + Saved + Notifications)
Week 4: Advanced Features (Search + Tags + Media Processing)
Week 5: Testing + Optimization + Documentation
```

---

# Phase 1: Foundation Setup 🔧

**Duration:** Week 1 (5-7 days)
**Goal:** Setup database, authentication, and Bunny Storage integration

## Step 1.1: Database Migration

### Tasks:
- [ ] Backup existing database
- [ ] Create all new models in `domain/models/`
- [ ] Update User model with new fields
- [ ] Run GORM AutoMigrate
- [ ] Verify tables created correctly

### Files to Create/Modify:
```
domain/models/
├── user.go (update)
├── post.go (new)
├── comment.go (new)
├── media.go (update)
├── vote.go (new)
├── follow.go (new)
├── saved_post.go (new)
├── notification.go (new)
├── notification_settings.go (new)
├── tag.go (new)
└── search_history.go (new)
```

### Code Example:
```go
// infrastructure/postgres/database.go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Post{},
        &models.Comment{},
        &models.Media{},
        &models.Vote{},
        &models.Follow{},
        &models.SavedPost{},
        &models.Notification{},
        &models.NotificationSettings{},
        &models.Tag{},
        &models.SearchHistory{},
    )
}
```

### Testing:
```bash
# Run migration
go run cmd/api/main.go

# Check tables
psql -U postgres -d gofiber_template
\dt

# Should see all new tables
```

---

## Step 1.2: Update User Service

### Tasks:
- [ ] Update User model DTOs
- [ ] Add new fields to UserResponse DTO
- [ ] Update Register endpoint (add displayName)
- [ ] Update Profile endpoints
- [ ] Create notification settings on user registration

### Files to Update:
```
domain/dto/
├── user_dto.go (update)
└── mappers.go (update)

application/serviceimpl/
└── user_service_impl.go (update)
```

### Updated User DTO:
```go
// domain/dto/user_dto.go
type CreateUserRequest struct {
    Username    string `json:"username" validate:"required,min=3,max=20"`
    Email       string `json:"email" validate:"required,email"`
    Password    string `json:"password" validate:"required,min=8"`
    DisplayName string `json:"displayName" validate:"required,max=50"`
}

type UserResponse struct {
    ID              string  `json:"id"`
    Username        string  `json:"username"`
    Email           string  `json:"email,omitempty"` // Only for owner
    DisplayName     string  `json:"displayName"`
    Avatar          *string `json:"avatar"`
    Karma           int     `json:"karma"`
    Bio             *string `json:"bio"`
    CoverImage      *string `json:"coverImage"`
    JoinedAt        string  `json:"joinedAt"`
    Location        *string `json:"location"`
    Website         *string `json:"website"`
    FollowersCount  int     `json:"followersCount"`
    FollowingCount  int     `json:"followingCount"`
    IsFollowing     *bool   `json:"isFollowing,omitempty"` // Only when authenticated
}
```

---

## Step 1.3: Bunny Storage Integration

### Tasks:
- [ ] Update BunnyStorage service
- [ ] Implement Upload method
- [ ] Implement Delete method
- [ ] Test file upload to Bunny CDN
- [ ] Verify CDN URLs work

### Files to Create/Update:
```
infrastructure/storage/
└── bunny_storage.go (update)
```

### Implementation:
```go
// infrastructure/storage/bunny_storage.go
type BunnyStorage struct {
    storageZone string
    accessKey   string
    baseURL     string
    cdnURL      string
    httpClient  *http.Client
}

func NewBunnyStorage(config config.BunnyConfig) *BunnyStorage {
    return &BunnyStorage{
        storageZone: config.StorageZone,
        accessKey:   config.AccessKey,
        baseURL:     config.BaseURL,
        cdnURL:      config.CDNUrl,
        httpClient:  &http.Client{Timeout: 30 * time.Second},
    }
}

func (b *BunnyStorage) Upload(ctx context.Context, file multipart.File, filename string) (string, error) {
    // Read file content
    fileBytes, err := io.ReadAll(file)
    if err != nil {
        return "", err
    }

    // Build upload URL
    uploadURL := fmt.Sprintf("%s/%s/%s", b.baseURL, b.storageZone, filename)

    // Create request
    req, err := http.NewRequestWithContext(ctx, "PUT", uploadURL, bytes.NewReader(fileBytes))
    if err != nil {
        return "", err
    }

    // Set headers
    req.Header.Set("AccessKey", b.accessKey)
    req.Header.Set("Content-Type", "application/octet-stream")

    // Send request
    resp, err := b.httpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("bunny storage upload failed: %d", resp.StatusCode)
    }

    // Return CDN URL
    cdnURL := fmt.Sprintf("%s/%s", b.cdnURL, filename)
    return cdnURL, nil
}

func (b *BunnyStorage) Delete(ctx context.Context, filename string) error {
    deleteURL := fmt.Sprintf("%s/%s/%s", b.baseURL, b.storageZone, filename)

    req, err := http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
    if err != nil {
        return err
    }

    req.Header.Set("AccessKey", b.accessKey)

    resp, err := b.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("bunny storage delete failed: %d", resp.StatusCode)
    }

    return nil
}
```

### Testing:
```bash
# Test upload
curl -X POST http://localhost:3000/api/v1/media/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "files=@test-image.jpg"

# Check response contains Bunny CDN URL
```

---

## Phase 1 Completion Checklist

- [ ] All database tables created
- [ ] Indexes added
- [ ] User model updated with new fields
- [ ] Register/Login works with displayName
- [ ] Bunny Storage upload works
- [ ] Bunny Storage delete works
- [ ] CDN URLs accessible
- [ ] All migrations pass
- [ ] No breaking changes to existing API

---

# Phase 2: Core Features 📝

**Duration:** Week 2 (5-7 days)
**Goal:** Implement Posts, Comments, and Vote system

## Step 2.1: Posts Module

### Tasks:
- [ ] Create Post repository interface
- [ ] Implement Post repository
- [ ] Create Post service interface
- [ ] Implement Post service
- [ ] Create Post handlers
- [ ] Setup Post routes
- [ ] Implement pagination
- [ ] Implement sorting (hot/new/top)

### Files to Create:
```
domain/
├── repositories/
│   └── post_repository.go
└── services/
    └── post_service.go

infrastructure/postgres/
└── post_repository_impl.go

application/serviceimpl/
└── post_service_impl.go

interfaces/api/
├── handlers/
│   └── post_handler.go
└── routes/
    └── post_routes.go

domain/dto/
└── post_dto.go
```

### Post Repository Interface:
```go
// domain/repositories/post_repository.go
type PostRepository interface {
    Create(ctx context.Context, post *models.Post) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error)
    List(ctx context.Context, params ListPostsParams) ([]*models.Post, int64, error)
    Update(ctx context.Context, id uuid.UUID, post *models.Post) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetByAuthor(ctx context.Context, authorID uuid.UUID, offset, limit int) ([]*models.Post, int64, error)
    IncrementCommentCount(ctx context.Context, postID uuid.UUID) error
    DecrementCommentCount(ctx context.Context, postID uuid.UUID) error
}

type ListPostsParams struct {
    Page      int
    Limit     int
    SortBy    string // hot, new, top
    TimeRange string // today, week, month, year, all
    Tag       string
    Author    string
}
```

### Post Service Interface:
```go
// domain/services/post_service.go
type PostService interface {
    CreatePost(ctx context.Context, req *dto.CreatePostRequest, userID uuid.UUID) (*models.Post, error)
    GetPost(ctx context.Context, postID uuid.UUID, userID *uuid.UUID) (*dto.PostResponse, error)
    ListPosts(ctx context.Context, params ListPostsParams, userID *uuid.UUID) ([]*dto.PostResponse, *dto.PaginationMeta, error)
    UpdatePost(ctx context.Context, postID uuid.UUID, req *dto.UpdatePostRequest, userID uuid.UUID) (*models.Post, error)
    DeletePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
    GetUserPosts(ctx context.Context, username string, page, limit int) ([]*dto.PostResponse, *dto.PaginationMeta, error)
}
```

### Post Routes:
```go
// interfaces/api/routes/post_routes.go
func SetupPostRoutes(api fiber.Router, h *handlers.Handlers) {
    posts := api.Group("/posts")

    // Public routes
    posts.Get("/", h.PostHandler.ListPosts)
    posts.Get("/:id", h.PostHandler.GetPost)
    posts.Get("/user/:username", h.PostHandler.GetUserPosts)

    // Protected routes
    posts.Use(middleware.Protected())
    posts.Post("/", h.PostHandler.CreatePost)
    posts.Put("/:id", middleware.OwnerOnly(), h.PostHandler.UpdatePost)
    posts.Delete("/:id", middleware.OwnerOnly(), h.PostHandler.DeletePost)
}
```

### Testing:
```bash
# Create post
curl -X POST http://localhost:3000/api/v1/posts \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Post",
    "content": "This is a test post",
    "tags": ["test", "golang"]
  }'

# Get posts
curl http://localhost:3000/api/v1/posts?page=1&limit=20&sortBy=hot
```

---

## Step 2.2: Comments Module

### Tasks:
- [ ] Create Comment repository
- [ ] Implement Comment service
- [ ] Create Comment handlers
- [ ] Support nested replies (max depth 10)
- [ ] Auto-update post comment_count
- [ ] Implement comment sorting

### Key Implementation:
```go
// When creating a comment
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *dto.CreateCommentRequest) (*models.Comment, error) {
    comment := &models.Comment{
        ID:       uuid.New(),
        PostID:   req.PostID,
        AuthorID: req.AuthorID,
        Content:  req.Content,
        Depth:    0,
    }

    // If reply, set parent and calculate depth
    if req.ParentID != nil {
        parent, err := s.commentRepo.GetByID(ctx, *req.ParentID)
        if err != nil {
            return nil, err
        }

        comment.ParentID = req.ParentID
        comment.Depth = parent.Depth + 1

        // Check max depth
        if comment.Depth > 10 {
            return nil, errors.New("maximum nesting depth reached")
        }
    }

    // Create comment
    if err := s.commentRepo.Create(ctx, comment); err != nil {
        return nil, err
    }

    // Increment post comment count
    if err := s.postRepo.IncrementCommentCount(ctx, comment.PostID); err != nil {
        // Log error but don't fail
    }

    return comment, nil
}
```

---

## Step 2.3: Vote System

### Tasks:
- [ ] Create Vote repository
- [ ] Implement Vote service
- [ ] Create Vote handlers
- [ ] Update karma on vote changes
- [ ] Support polymorphic voting (posts + comments)
- [ ] Handle vote changes (up→down, remove vote)

### Vote Service Implementation:
```go
// application/serviceimpl/vote_service_impl.go
func (s *VoteServiceImpl) VotePost(ctx context.Context, postID, userID uuid.UUID, voteType *string) error {
    // Get existing vote
    existingVote, _ := s.voteRepo.GetVote(ctx, userID, postID, "post")

    var deltaVotes int
    var deltaKarma int

    if voteType == nil {
        // Remove vote
        if existingVote != nil {
            s.voteRepo.DeleteVote(ctx, userID, postID, "post")
            if existingVote.VoteType == "up" {
                deltaVotes = -1
                deltaKarma = -1
            } else {
                deltaVotes = 1
                deltaKarma = 1
            }
        }
    } else if existingVote == nil {
        // Create new vote
        vote := &models.Vote{
            UserID:     userID,
            TargetID:   postID,
            TargetType: "post",
            VoteType:   *voteType,
        }
        s.voteRepo.Create(ctx, vote)

        if *voteType == "up" {
            deltaVotes = 1
            deltaKarma = 1
        } else {
            deltaVotes = -1
            deltaKarma = -1
        }
    } else {
        // Update existing vote
        if existingVote.VoteType != *voteType {
            existingVote.VoteType = *voteType
            s.voteRepo.Update(ctx, existingVote)

            if *voteType == "up" {
                deltaVotes = 2
                deltaKarma = 2
            } else {
                deltaVotes = -2
                deltaKarma = -2
            }
        }
    }

    // Update post votes
    if deltaVotes != 0 {
        s.postRepo.UpdateVotes(ctx, postID, deltaVotes)
    }

    // Update author karma
    if deltaKarma != 0 {
        post, _ := s.postRepo.GetByID(ctx, postID)
        if post != nil {
            s.userRepo.UpdateKarma(ctx, post.AuthorID, deltaKarma)
        }
    }

    return nil
}
```

---

## Phase 2 Completion Checklist

- [ ] Posts CRUD working
- [ ] Comments CRUD working
- [ ] Nested comments working (max depth 10)
- [ ] Vote system working (posts)
- [ ] Vote system working (comments)
- [ ] Karma calculation working
- [ ] Post comment_count auto-updates
- [ ] Pagination working
- [ ] Sorting working (hot/new/top)
- [ ] All endpoints tested

---

# Phase 3: Social Features 👥

**Duration:** Week 3 (5-7 days)
**Goal:** Follow system, Saved Posts, Notifications

## Step 3.1: Follow System

### Tasks:
- [ ] Create Follow repository
- [ ] Implement Follow service
- [ ] Create Follow/Unfollow handlers
- [ ] Auto-update followers_count
- [ ] Auto-update following_count
- [ ] Prevent self-follow
- [ ] Get followers/following lists

### Implementation:
```go
func (s *FollowServiceImpl) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
    // Check not following self
    if followerID == followingID {
        return errors.New("cannot follow yourself")
    }

    // Check not already following
    exists, _ := s.followRepo.IsFollowing(ctx, followerID, followingID)
    if exists {
        return errors.New("already following")
    }

    // Create follow
    follow := &models.Follow{
        FollowerID:  followerID,
        FollowingID: followingID,
        CreatedAt:   time.Now(),
    }

    if err := s.followRepo.Create(ctx, follow); err != nil {
        return err
    }

    // Update counts
    s.userRepo.IncrementFollowersCount(ctx, followingID)
    s.userRepo.IncrementFollowingCount(ctx, followerID)

    // Create notification
    s.notificationService.CreateFollowNotification(ctx, followingID, followerID)

    return nil
}
```

---

## Step 3.2: Saved Posts

### Tasks:
- [ ] Create SavedPost repository
- [ ] Implement SavedPost service
- [ ] Create Save/Unsave handlers
- [ ] Get saved posts list
- [ ] Check if post is saved
- [ ] Clear all saved posts

### Routes:
```go
func SetupSavedRoutes(api fiber.Router, h *handlers.Handlers) {
    saved := api.Group("/saved")
    saved.Use(middleware.Protected())

    saved.Get("/", h.SavedHandler.GetSavedPosts)
    saved.Get("/count", h.SavedHandler.GetSavedCount)
    saved.Get("/:postId", h.SavedHandler.CheckSaved)
    saved.Post("/:postId", h.SavedHandler.SavePost)
    saved.Delete("/:postId", h.SavedHandler.UnsavePost)
    saved.Delete("/", h.SavedHandler.ClearAll)
}
```

---

## Step 3.3: Notifications

### Tasks:
- [ ] Create Notification repository
- [ ] Create NotificationSettings repository
- [ ] Implement Notification service
- [ ] Create notification handlers
- [ ] Auto-create notifications on events
- [ ] Filter by type/read status
- [ ] Mark as read (single/all)
- [ ] Settings CRUD

### Notification Creation Logic:
```go
// Create notification on comment reply
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *dto.CreateCommentRequest) (*models.Comment, error) {
    // ... create comment logic ...

    // If reply to another comment, notify parent author
    if comment.ParentID != nil {
        parent, _ := s.commentRepo.GetByID(ctx, *comment.ParentID)
        if parent != nil && parent.AuthorID != comment.AuthorID {
            // Check user notification settings
            settings, _ := s.notificationService.GetSettings(ctx, parent.AuthorID)
            if settings.Replies {
                notification := &models.Notification{
                    UserID:    parent.AuthorID,
                    SenderID:  comment.AuthorID,
                    Type:      "reply",
                    Message:   "ตอบกลับคอมเมนต์ของคุณ",
                    PostID:    &comment.PostID,
                    CommentID: &comment.ID,
                }
                s.notificationRepo.Create(ctx, notification)
            }
        }
    }

    return comment, nil
}
```

---

## Phase 3 Completion Checklist

- [ ] Follow/Unfollow working
- [ ] Follower counts auto-update
- [ ] Get followers/following lists
- [ ] Save/Unsave posts working
- [ ] Get saved posts list
- [ ] Notifications created on events
- [ ] Get notifications working
- [ ] Filter notifications (type/read)
- [ ] Mark as read working
- [ ] Notification settings CRUD

---

# Phase 4: Advanced Features 🔍

**Duration:** Week 4 (5-7 days)
**Goal:** Search, Tags, Media Processing

## Step 4.1: Search System

### Tasks:
- [ ] Setup PostgreSQL full-text search
- [ ] Implement search repository
- [ ] Search posts by title/content
- [ ] Search users by username/displayName/bio
- [ ] Implement relevance ranking
- [ ] Add search filters
- [ ] Popular tags endpoint
- [ ] Search history (optional)

### PostgreSQL Full-Text Search:
```go
// infrastructure/postgres/search_repository_impl.go
func (r *SearchRepositoryImpl) SearchPosts(ctx context.Context, query string, params SearchParams) ([]*models.Post, int64, error) {
    var posts []*models.Post
    var total int64

    db := r.db.WithContext(ctx)

    // Build full-text search query
    tsQuery := "to_tsquery('english', ?)"
    searchQuery := strings.ReplaceAll(query, " ", " & ")

    // Search in title and content
    db = db.Where(
        db.Where("to_tsvector('english', title) @@ "+tsQuery, searchQuery).
        Or("to_tsvector('english', content) @@ "+tsQuery, searchQuery),
    )

    // Count total
    db.Model(&models.Post{}).Count(&total)

    // Add ordering by relevance
    db = db.Order(gorm.Expr(`
        ts_rank(to_tsvector('english', title), to_tsquery('english', ?)) * 2 +
        ts_rank(to_tsvector('english', content), to_tsquery('english', ?)) DESC
    `, searchQuery, searchQuery))

    // Pagination
    db = db.Offset((params.Page - 1) * params.Limit).Limit(params.Limit)

    // Execute
    if err := db.Preload("Author").Preload("Media").Find(&posts).Error; err != nil {
        return nil, 0, err
    }

    return posts, total, nil
}
```

---

## Step 4.2: Tag System

### Tasks:
- [ ] Create Tag repository
- [ ] Implement tag auto-creation on post
- [ ] Update post_count on tag operations
- [ ] Get posts by tag
- [ ] Get popular tags
- [ ] Trending tags calculation

### Tag Handler:
```go
func (h *TagHandler) GetPopularTags(c *fiber.Ctx) error {
    limit, _ := strconv.Atoi(c.Query("limit", "20"))
    timeRange := c.Query("timeRange", "all")

    tags, err := h.tagService.GetPopularTags(c.Context(), limit, timeRange)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get tags", err)
    }

    return utils.SuccessResponse(c, "Tags retrieved successfully", fiber.Map{
        "tags": tags,
    })
}
```

---

## Step 4.3: Media Processing

### Tasks:
- [ ] Image compression
- [ ] Thumbnail generation (200x200)
- [ ] WebP conversion
- [ ] Video thumbnail extraction
- [ ] Storage usage calculation
- [ ] Media optimization endpoint

### Image Processing:
```go
// Use github.com/disintegration/imaging
func (s *MediaServiceImpl) ProcessImage(file multipart.File) ([]byte, error) {
    // Decode image
    img, err := imaging.Decode(file)
    if err != nil {
        return nil, err
    }

    // Resize if too large (max 1920px)
    if img.Bounds().Dx() > 1920 || img.Bounds().Dy() > 1920 {
        img = imaging.Fit(img, 1920, 1920, imaging.Lanczos)
    }

    // Convert to JPEG with quality 85
    var buf bytes.Buffer
    if err := imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(85)); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

func (s *MediaServiceImpl) GenerateThumbnail(file multipart.File) ([]byte, error) {
    img, err := imaging.Decode(file)
    if err != nil {
        return nil, err
    }

    // Create 200x200 thumbnail (crop center)
    thumb := imaging.Fill(img, 200, 200, imaging.Center, imaging.Lanczos)

    var buf bytes.Buffer
    if err := imaging.Encode(&buf, thumb, imaging.JPEG, imaging.JPEGQuality(80)); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}
```

---

## Phase 4 Completion Checklist

- [ ] Search posts working
- [ ] Search users working
- [ ] Search relevance ranking
- [ ] Tag system working
- [ ] Popular tags endpoint
- [ ] Image compression working
- [ ] Thumbnail generation working
- [ ] Video thumbnail extraction
- [ ] Storage usage calculation
- [ ] All media endpoints tested

---

# Phase 5: Testing & Optimization ✅

**Duration:** Week 5 (5-7 days)
**Goal:** Testing, optimization, documentation

## Step 5.1: Integration Testing

### Tasks:
- [ ] Test all public endpoints
- [ ] Test all private endpoints
- [ ] Test authentication flow
- [ ] Test vote system
- [ ] Test notification creation
- [ ] Test follow system
- [ ] Test search functionality
- [ ] Test media upload
- [ ] Load testing with k6 or similar

---

## Step 5.2: Performance Optimization

### Tasks:
- [ ] Add Redis caching for hot posts
- [ ] Cache user profiles
- [ ] Cache popular tags
- [ ] Add database query logging
- [ ] Optimize N+1 queries
- [ ] Add pagination everywhere
- [ ] Implement rate limiting

### Redis Caching Example:
```go
func (s *PostServiceImpl) GetHotPosts(ctx context.Context) ([]*models.Post, error) {
    cacheKey := "hot_posts"

    // Try cache first
    cached, err := s.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var posts []*models.Post
        json.Unmarshal([]byte(cached), &posts)
        return posts, nil
    }

    // Cache miss, get from DB
    posts, err := s.postRepo.GetHotPosts(ctx, 20)
    if err != nil {
        return nil, err
    }

    // Cache for 5 minutes
    data, _ := json.Marshal(posts)
    s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

    return posts, nil
}
```

---

## Step 5.3: Documentation

### Tasks:
- [ ] API documentation (Swagger/Postman)
- [ ] README updates
- [ ] Environment variables documentation
- [ ] Deployment guide
- [ ] Code comments
- [ ] Architecture diagram

---

## Phase 5 Completion Checklist

- [ ] All endpoints tested
- [ ] Integration tests passing
- [ ] Load testing completed
- [ ] Redis caching implemented
- [ ] Rate limiting working
- [ ] Documentation complete
- [ ] Deployment guide ready
- [ ] Code reviewed

---

## 🎉 Project Complete!

After completing all 5 phases, you should have:
- ✅ 61 working API endpoints
- ✅ Complete social media features
- ✅ Bunny Storage integration
- ✅ Full-text search
- ✅ Notification system
- ✅ Vote & Karma system
- ✅ Production-ready code

**Next:** Deploy to production! See `06-deployment.md`
