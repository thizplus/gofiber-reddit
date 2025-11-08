# สถาปัตยกรรมระบบ (Architecture)

## Clean Architecture Overview

ระบบนี้ใช้ **Clean Architecture** (Uncle Bob) ซึ่งแบ่งออกเป็น 4 layers หลัก:

```
┌─────────────────────────────────────────────────────┐
│           Interface Layer (Presentation)             │
│   - HTTP Handlers                                    │
│   - Middleware                                       │
│   - Routes                                           │
│   - WebSocket Handler                                │
└────────────────────┬────────────────────────────────┘
                     │ Depends on ↓
┌────────────────────▼────────────────────────────────┐
│          Application Layer (Use Cases)               │
│   - Service Implementations                          │
│   - Business Logic                                   │
│   - Orchestration                                    │
└────────────────────┬────────────────────────────────┘
                     │ Depends on ↓
┌────────────────────▼────────────────────────────────┐
│              Domain Layer (Entities)                 │
│   - Models (Database Entities)                       │
│   - DTOs (Data Transfer Objects)                     │
│   - Repository Interfaces                            │
│   - Service Interfaces                               │
└────────────────────┬────────────────────────────────┘
                     │ Implemented by ↑
┌────────────────────▼────────────────────────────────┐
│        Infrastructure Layer (External)               │
│   - PostgreSQL Repositories                          │
│   - Redis Client                                     │
│   - Bunny CDN Storage                                │
│   - WebSocket Manager                                │
└─────────────────────────────────────────────────────┘
```

## Layer Details

### 1. Domain Layer (`domain/`)

**หน้าที่**: กำหนด business entities และ contracts (interfaces)

**โครงสร้าง**:
```
domain/
├── models/              # Database models (15 models)
│   ├── user.go
│   ├── post.go
│   ├── comment.go
│   ├── vote.go
│   ├── follow.go
│   ├── saved_post.go
│   ├── notification.go
│   ├── notification_settings.go
│   ├── push_subscription.go
│   ├── media.go
│   ├── tag.go
│   ├── search_history.go
│   ├── task.go (legacy)
│   ├── file.go (legacy)
│   └── job.go (legacy)
├── dto/                 # Data Transfer Objects
│   ├── auth_dto.go
│   ├── user_dto.go
│   ├── post_dto.go
│   ├── comment_dto.go
│   ├── vote_dto.go
│   ├── notification_dto.go
│   └── ...
├── repositories/        # Repository interfaces (contracts)
│   ├── user_repository.go
│   ├── post_repository.go
│   ├── comment_repository.go
│   └── ...
└── services/           # Service interfaces (contracts)
    ├── user_service.go
    ├── post_service.go
    ├── comment_service.go
    └── ...
```

**ความสำคัญ**:
- Layer นี้ **ไม่ depend** กับ layer อื่นใด
- เป็นศูนย์กลางของ business logic
- Interfaces ในนี้จะถูก implement โดย layer อื่น

### 2. Application Layer (`application/`)

**หน้าที่**: ประมวลผล business logic และ use cases

**โครงสร้าง**:
```
application/
└── serviceimpl/        # Service implementations
    ├── user_service_impl.go
    ├── post_service_impl.go
    ├── comment_service_impl.go
    ├── vote_service_impl.go
    ├── follow_service_impl.go
    ├── saved_post_service_impl.go
    ├── notification_service_impl.go
    ├── notification_settings_service_impl.go
    ├── push_service_impl.go
    ├── media_service_impl.go
    ├── tag_service_impl.go
    ├── search_service_impl.go
    ├── task_service_impl.go (legacy)
    ├── file_service_impl.go (legacy)
    └── job_service_impl.go (legacy)
```

**ตัวอย่าง Use Case**:
```go
// CreatePost use case
func (s *PostServiceImpl) CreatePost(dto *dto.CreatePostDTO) (*models.Post, error) {
    // 1. Validate input
    // 2. Create post entity
    // 3. Handle tags (create if not exist)
    // 4. Handle media attachments
    // 5. Save to database via repository
    // 6. Create notification for followers
    // 7. Return created post
}
```

**หน้าที่หลัก**:
- รับ DTO จาก interface layer
- ประมวลผล business logic
- เรียกใช้ repositories
- จัดการ transactions
- Error handling
- คืนผลลัพธ์กลับไปยัง interface layer

### 3. Infrastructure Layer (`infrastructure/`)

**หน้าที่**: Implement รายละเอียดทางเทคนิคและ external dependencies

**โครงสร้าง**:
```
infrastructure/
├── postgres/           # PostgreSQL implementation
│   ├── database.go                    # DB connection & migration
│   ├── migrations/                    # SQL migration files
│   ├── user_repository_impl.go
│   ├── post_repository_impl.go
│   ├── comment_repository_impl.go
│   ├── vote_repository_impl.go
│   ├── follow_repository_impl.go
│   ├── saved_post_repository_impl.go
│   ├── notification_repository_impl.go
│   ├── notification_settings_repository_impl.go
│   ├── push_subscription_repository_impl.go
│   ├── media_repository_impl.go
│   ├── tag_repository_impl.go
│   ├── search_history_repository_impl.go
│   ├── task_repository_impl.go (legacy)
│   ├── file_repository_impl.go (legacy)
│   └── job_repository_impl.go (legacy)
├── redis/              # Redis client
│   └── client.go
├── storage/            # Bunny CDN client
│   └── bunny_storage.go
└── websocket/          # WebSocket manager
    └── manager.go
```

**ตัวอย่าง Repository Implementation**:
```go
type UserRepositoryImpl struct {
    db *gorm.DB
}

func (r *UserRepositoryImpl) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}
```

**External Dependencies**:
- PostgreSQL (GORM)
- Redis (go-redis)
- Bunny CDN (HTTP client)
- WebSocket (gorilla/websocket)

### 4. Interface Layer (`interfaces/`)

**หน้าที่**: รับ HTTP requests, WebSocket connections และแปลงเป็น DTOs

**โครงสร้าง**:
```
interfaces/
└── api/
    ├── handlers/           # HTTP request handlers
    │   ├── auth_handler.go
    │   ├── user_handler.go
    │   ├── profile_handler.go
    │   ├── post_handler.go
    │   ├── comment_handler.go
    │   ├── vote_handler.go
    │   ├── follow_handler.go
    │   ├── saved_post_handler.go
    │   ├── notification_handler.go
    │   ├── tag_handler.go
    │   ├── search_handler.go
    │   ├── media_handler.go
    │   ├── push_handler.go
    │   ├── task_handler.go (legacy)
    │   ├── file_handler.go (legacy)
    │   ├── job_handler.go (legacy)
    │   └── health_handler.go
    ├── middleware/         # HTTP middleware
    │   ├── auth.go
    │   ├── cors.go
    │   ├── error_handler.go
    │   └── logger.go
    ├── routes/            # Route definitions
    │   ├── auth_routes.go
    │   ├── user_routes.go
    │   ├── profile_routes.go
    │   ├── post_routes.go
    │   ├── comment_routes.go
    │   ├── vote_routes.go
    │   ├── follow_routes.go
    │   ├── saved_post_routes.go
    │   ├── notification_routes.go
    │   ├── tag_routes.go
    │   ├── search_routes.go
    │   ├── media_routes.go
    │   ├── push_routes.go
    │   ├── task_routes.go (legacy)
    │   ├── file_routes.go (legacy)
    │   ├── job_routes.go (legacy)
    │   ├── health_routes.go
    │   └── routes.go (main router)
    └── websocket/         # WebSocket handler
        └── handler.go
```

**ตัวอย่าง Handler**:
```go
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
    // 1. Parse request body to DTO
    var req dto.CreatePostDTO
    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
    }

    // 2. Validate DTO
    if err := validate.Struct(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }

    // 3. Call service
    post, err := h.postService.CreatePost(&req)
    if err != nil {
        return err
    }

    // 4. Return response
    return c.Status(fiber.StatusCreated).JSON(post)
}
```

## Shared Packages (`pkg/`)

**หน้าที่**: Utilities และ shared logic ที่ใช้ร่วมกันทุก layer

```
pkg/
├── config/             # Configuration management
│   └── config.go       # Load .env และ validate config
├── di/                 # Dependency Injection container
│   └── container.go    # Initialize & wire all dependencies
├── utils/              # Utility functions
│   ├── jwt.go          # JWT generation/validation
│   ├── password.go     # Password hashing/verification
│   ├── response.go     # Standard API responses
│   └── validator.go    # Custom validators
├── scheduler/          # Job scheduler
│   └── scheduler.go    # Cron job management
└── auth_code_store/    # OAuth code storage (Redis)
    └── store.go
```

## Dependency Injection Flow

```
main.go
  ↓
Initialize Config
  ↓
Initialize Database (PostgreSQL)
  ↓
Initialize Redis
  ↓
Initialize Storage (Bunny CDN)
  ↓
Initialize Repositories (15 repos)
  ↓
Initialize Services (15 services)
  ↓
Initialize Handlers (18 handlers)
  ↓
Setup Routes
  ↓
Start Fiber Server
```

**DI Container (`pkg/di/container.go`)**:
```go
type Container struct {
    // Infrastructure
    DB              *gorm.DB
    RedisClient     *redis.Client
    BunnyStorage    *storage.BunnyStorage
    WebSocketMgr    *websocket.Manager

    // Repositories
    UserRepo        repositories.UserRepository
    PostRepo        repositories.PostRepository
    // ... 13 more repos

    // Services
    UserService     services.UserService
    PostService     services.PostService
    // ... 13 more services

    // Handlers
    AuthHandler     *handlers.AuthHandler
    UserHandler     *handlers.UserHandler
    // ... 16 more handlers
}
```

## Request Flow Example

**ตัวอย่าง: สร้าง Post**

```
Client (POST /api/v1/posts)
  ↓
Fiber Router
  ↓
Logger Middleware (log request)
  ↓
CORS Middleware (check origin)
  ↓
Auth Middleware (validate JWT)
  ↓
PostHandler.CreatePost()
  ↓ (parse DTO, validate)
PostService.CreatePost()
  ↓ (business logic)
├─ TagRepository.FindOrCreate() (handle tags)
├─ MediaRepository.AttachToPost() (handle media)
├─ PostRepository.Create() (save post)
└─ NotificationService.NotifyFollowers() (async)
  ↓
Return Response (201 Created)
```

## Database Transaction Management

**ระดับ Service Layer**:
```go
func (s *PostServiceImpl) CreatePost(dto *dto.CreatePostDTO) (*models.Post, error) {
    // เริ่ม transaction
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Business logic with transaction
    post := &models.Post{...}
    if err := tx.Create(post).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // Handle tags within transaction
    for _, tagName := range dto.Tags {
        // ...
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return post, nil
}
```

## Error Handling Strategy

**3 ระดับของ Error Handling**:

1. **Repository Level**: Database errors
```go
if err := r.db.Create(user).Error; err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

2. **Service Level**: Business logic errors
```go
if existingUser != nil {
    return nil, fiber.NewError(fiber.StatusConflict, "Username already exists")
}
```

3. **Handler Level**: HTTP errors
```go
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }
    return fiber.NewError(fiber.StatusInternalServerError, err.Error())
}
```

## Middleware Chain

```
Request
  ↓
LoggerMiddleware (log all requests)
  ↓
CorsMiddleware (handle CORS)
  ↓
AuthMiddleware (validate JWT) [if protected]
  ↓
RoleMiddleware (check permissions) [if admin-only]
  ↓
Handler
  ↓
ErrorHandler (global error handling)
  ↓
Response
```

## WebSocket Architecture

```
Client
  ↓ (ws://host/ws?token=xxx&room=yyy)
WebSocket Handler
  ↓ (authenticate)
WebSocket Manager
  ↓ (register client)
├─ Store in clients map (userID → connections)
├─ Join room (roomID → clients)
└─ Start heartbeat (30s ping/pong)
  ↓
Message Loop
├─ Receive message
├─ Process message type
└─ Broadcast to room/user/all
```

**Manager Methods**:
- `RegisterClient(client *Client)`
- `UnregisterClient(client *Client)`
- `BroadcastToAll(message []byte)`
- `BroadcastToRoom(roomID string, message []byte)`
- `SendToUser(userID string, message []byte)`

## Configuration Management

**Environment Variables** → **Config Struct** → **DI Container**

```go
// Load from .env
config := config.LoadConfig()

// Validate required fields
if config.DBHost == "" {
    log.Fatal("DB_HOST is required")
}

// Pass to services
container.InitializeAll(config)
```

## Advantages of This Architecture

1. **Testability**: แต่ละ layer สามารถ test แยกได้
2. **Maintainability**: แยก concerns ชัดเจน
3. **Scalability**: เพิ่ม feature ง่ายโดยไม่กระทบส่วนอื่น
4. **Flexibility**: เปลี่ยน database/framework ได้ง่าย
5. **Reusability**: Share business logic ได้
6. **Team Collaboration**: แต่ละ team ทำงานใน layer ของตัวเองได้
