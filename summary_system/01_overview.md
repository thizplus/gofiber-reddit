# ภาพรวมระบบ Social Media Platform

## ข้อมูลพื้นฐาน

**ชื่อโปรเจกต์**: Social Media Platform (Go Fiber Backend)
**ภาษาที่ใช้**: Go 1.24.0
**Framework**: Fiber v2.52.0 (Express-style web framework)
**สถาปัตยกรรม**: Clean Architecture (4 layers)
**ฐานข้อมูล**: PostgreSQL 15
**Cache**: Redis 7
**Storage**: Bunny CDN

## วัตถุประสงค์ของระบบ

ระบบ Social Media Platform นี้เป็นแพลตฟอร์มโซเชียลมีเดียแบบครบวงจร ที่ผสมผสานฟีเจอร์จากหลายแพลตฟอร์มยอดนิยม:

- **Reddit-style**: ระบบโหวต (upvote/downvote), karma, nested comments
- **Twitter-style**: ระบบ follow/followers, personalized feed
- **Instagram-style**: การแชร์รูปภาพและวิดีโอ, media gallery
- **Modern features**: Real-time notifications, WebSocket, Push notifications

## Tech Stack

### Backend
- **Web Framework**: Fiber v2.52.0
- **ORM**: GORM v1.25.6
- **Authentication**: JWT (golang-jwt/jwt v5.2.1)
- **OAuth**: golang.org/x/oauth2 (Google OAuth2)
- **Validation**: go-playground/validator v10.16.0
- **Password**: golang.org/x/crypto (bcrypt)
- **WebSocket**: gofiber/websocket v2.2.1
- **Push Notifications**: webpush-go v1.4.0
- **Job Scheduler**: gocron v1.37.0

### Infrastructure
- **Database**: PostgreSQL 15 (Alpine)
- **Cache**: Redis 7 (Alpine)
- **Storage**: Bunny CDN
- **Container**: Docker + Docker Compose

### Development Tools
- **Package Manager**: Go Modules
- **Environment**: godotenv v1.5.1
- **Version Control**: Git

## คุณลักษณะเด่น

### 1. Clean Architecture
- แยก layer อย่างชัดเจน: Domain, Application, Infrastructure, Interface
- Dependency Injection แบบ custom
- Testable และ maintainable
- Domain-driven design

### 2. ครบครันและ Production-Ready
- Docker containerization พร้อม health checks
- Graceful shutdown
- Comprehensive error handling
- Security best practices
- CORS configuration
- Environment-based configuration

### 3. Scalable Design
- Redis caching layer
- Database indexing ที่เหมาะสม
- Pagination ทุก list endpoints
- WebSocket สำหรับ real-time (ไม่ใช่ polling)
- CDN สำหรับ media delivery
- พร้อมสำหรับ horizontal scaling

### 4. Feature-Rich
- มากกว่า 60+ API endpoints
- 15 database models
- 15 services และ repositories
- Real-time WebSocket
- Push notifications
- OAuth integration
- Media upload system

## โครงสร้างโปรเจกต์หลัก

```
gofiber-backend/
├── cmd/api/                 # Entry point (main.go)
├── domain/                  # Domain Layer
│   ├── models/             # Database Models
│   ├── dto/                # Data Transfer Objects
│   ├── repositories/       # Repository Interfaces
│   └── services/           # Service Interfaces
├── application/             # Application Layer
│   └── serviceimpl/        # Service Implementations
├── infrastructure/          # Infrastructure Layer
│   ├── postgres/           # PostgreSQL + Repositories
│   ├── redis/              # Redis Client
│   ├── storage/            # Bunny CDN
│   └── websocket/          # WebSocket Manager
├── interfaces/              # Presentation Layer
│   └── api/
│       ├── handlers/       # HTTP Handlers
│       ├── middleware/     # Middleware
│       ├── routes/         # Route Definitions
│       └── websocket/      # WebSocket Handler
└── pkg/                     # Shared Packages
    ├── config/             # Configuration
    ├── di/                 # Dependency Injection
    ├── utils/              # Utilities
    ├── scheduler/          # Job Scheduler
    └── auth_code_store/    # OAuth Code Storage
```

## Feature List

### Core Features (Implemented)
- ✅ User Authentication (Email/Password + Google OAuth)
- ✅ User Profiles (Bio, Avatar, Stats)
- ✅ Posts (Create, Read, Update, Delete)
- ✅ Nested Comments (max 10 levels)
- ✅ Voting System (Posts & Comments)
- ✅ Follow System
- ✅ Tag System
- ✅ Media Upload (Images & Videos)
- ✅ Notifications (Real-time)
- ✅ Web Push Notifications
- ✅ Search (Posts, Users, Tags)
- ✅ Saved Posts
- ✅ Crossposting
- ✅ Personalized Feed
- ✅ WebSocket (Real-time)
- ✅ Karma System

### Planned Features
- 🔄 Direct Messaging (1-on-1 Chat)
- 🔄 Group Chat
- 🔄 Rate Limiting
- 🔄 Email Notifications
- 🔄 Admin Dashboard
- 🔄 Reporting System
- 🔄 Moderation Tools

## ตัวเลขสถิติ

- **API Endpoints**: 60+ endpoints
- **Database Models**: 15 models
- **Services**: 15 service implementations
- **Repositories**: 15 repository implementations
- **Handlers**: 18 HTTP handlers
- **Middleware**: 4 middleware functions
- **Route Files**: 20 route definition files
- **Lines of Code**: ~10,000+ lines (estimated)

## Use Cases หลัก

1. **สร้างและแชร์เนื้อหา**: ผู้ใช้สามารถโพสต์ข้อความ, รูปภาพ, วิดีโอ
2. **โต้ตอบกับเนื้อหา**: Comment, Vote, Save, Crosspost
3. **สร้างชุมชน**: Follow ผู้ใช้อื่น, Tag-based discovery
4. **ค้นหาเนื้อหา**: Search posts, users, tags
5. **รับการแจ้งเตือน**: Real-time notifications, Push notifications
6. **Chat (กำลังพัฒนา)**: Direct messaging แบบ real-time

## Performance Targets

- **API Response Time**: < 100ms (average)
- **Database Queries**: Optimized with indexes
- **Real-time Latency**: < 50ms (WebSocket)
- **Media Upload**: Up to 300MB per file
- **Concurrent Users**: Designed for 1000+ concurrent connections

## Security Features

- ✅ Password hashing (bcrypt)
- ✅ JWT authentication
- ✅ OAuth 2.0 (Google)
- ✅ CORS configuration
- ✅ Input validation
- ✅ SQL injection prevention (GORM)
- ✅ XSS prevention
- ⏳ Rate limiting (planned)
- ⏳ CSRF protection (planned)

## Next Steps

1. ทำเอกสารโดยละเอียดเพิ่มเติมใน `summary_system/`
2. พัฒนาระบบ Chat (ตาม chat_api_spec)
3. เพิ่ม Rate Limiting
4. พัฒนา Admin Dashboard
5. เพิ่ม Testing Coverage
6. Performance Optimization
7. Deploy to Production
