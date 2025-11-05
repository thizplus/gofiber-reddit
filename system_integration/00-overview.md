# System Integration Overview

## 📋 สรุปโปรเจค

นี่คือโปรเจค **Social Media Platform** แบบ Reddit-like ที่พัฒนาด้วย Go Fiber + PostgreSQL + Bunny Storage

### ฟีเจอร์หลัก
- ✅ โพสต์ + Crosspost (แชร์โพสต์)
- ✅ Comment System (Nested replies, max depth 10)
- ✅ Vote System (Upvote/Downvote)
- ✅ Follow/Followers System
- ✅ Karma Score
- ✅ Notification System
- ✅ Saved Posts (Bookmark)
- ✅ Full-text Search + Trending
- ✅ Tag System
- ✅ Media Upload (Bunny Storage CDN)

---

## 📊 Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                     Client                          │
│                   (Frontend)                        │
└───────────────────┬─────────────────────────────────┘
                    │ HTTP/REST API
┌───────────────────┴─────────────────────────────────┐
│                Go Fiber Backend                     │
│  ┌──────────────────────────────────────────────┐  │
│  │  Interfaces Layer (HTTP Handlers)            │  │
│  │  - Routes, Middleware, WebSocket             │  │
│  └─────────────────────┬────────────────────────┘  │
│  ┌─────────────────────┴────────────────────────┐  │
│  │  Application Layer (Services)                │  │
│  │  - Business Logic, Validation                │  │
│  └─────────────────────┬────────────────────────┘  │
│  ┌─────────────────────┴────────────────────────┐  │
│  │  Domain Layer (Models, Interfaces)           │  │
│  │  - DTOs, Service Contracts, Repositories     │  │
│  └─────────────────────┬────────────────────────┘  │
│  ┌─────────────────────┴────────────────────────┐  │
│  │  Infrastructure Layer                        │  │
│  │  - PostgreSQL, Redis, Bunny CDN, WebSocket   │  │
│  └──────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                    │
        ┌───────────┴───────────┐
        │                       │
   ┌────┴─────┐          ┌─────┴──────┐
   │PostgreSQL│          │Bunny Storage│
   │ Database │          │    CDN      │
   └──────────┘          └────────────┘
```

---

## 📁 Project Structure

```
gofiber-backend/
├── cmd/api/                     # Entry point
├── domain/                      # Core business logic
│   ├── models/                  # Domain entities
│   ├── dto/                     # Data Transfer Objects
│   ├── services/                # Service interfaces
│   └── repositories/            # Repository interfaces
├── application/                 # Application services
│   └── serviceimpl/            # Service implementations
├── infrastructure/              # External services
│   ├── postgres/               # PostgreSQL + Repositories
│   ├── redis/                  # Redis client
│   ├── storage/                # Bunny CDN integration
│   └── websocket/              # WebSocket server
├── interfaces/                  # API layer
│   └── api/
│       ├── handlers/           # HTTP handlers
│       ├── middleware/         # Middleware
│       ├── routes/             # Route definitions
│       └── websocket/          # WebSocket handlers
├── pkg/                        # Utilities
│   ├── config/                 # Configuration
│   ├── di/                     # Dependency Injection
│   ├── scheduler/              # Background jobs
│   └── utils/                  # Helper functions
├── backend_spec/               # API specifications
├── system_integration/         # Implementation guides
└── docker-compose.yml          # Local development
```

---

## 🎯 Implementation Roadmap

### Phase 1: Foundation (Week 1)
- Database schema migration
- User model enhancement
- Bunny Storage setup

### Phase 2: Core Features (Week 2)
- Posts API (8 endpoints)
- Comments API (6 endpoints)
- Vote System

### Phase 3: Social Features (Week 3)
- Users API (Follow system)
- Saved Posts API
- Notifications API

### Phase 4: Advanced Features (Week 4)
- Search API + Full-text search
- Media API + Optimization
- Tag System

### Phase 5: Testing & Polish (Week 5)
- Integration testing
- Performance optimization
- Documentation

---

## 📈 Statistics

### API Endpoints
- **Total:** 61 endpoints
- **Public:** 16 endpoints
- **Private:** 45 endpoints

### Database Tables
- **Existing:** 4 tables (User, Task, File, Job)
- **New:** 12+ tables (Posts, Comments, Votes, Follows, etc.)
- **Total:** 16+ tables

### Features
- **Authentication:** 5 endpoints
- **Posts:** 8 endpoints
- **Comments:** 6 endpoints
- **Users:** 10 endpoints
- **Notifications:** 8 endpoints
- **Saved Posts:** 6 endpoints
- **Search:** 8 endpoints
- **Media:** 6 endpoints

---

## 🔑 Key Technologies

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Backend** | Go Fiber | Web framework |
| **Database** | PostgreSQL | Primary data store |
| **Cache** | Redis | Session & caching |
| **Storage** | Bunny CDN | Media files |
| **ORM** | GORM | Database operations |
| **Auth** | JWT | Token-based authentication |
| **Search** | PostgreSQL FTS | Full-text search |
| **Jobs** | gocron | Background scheduler |
| **WebSocket** | Fiber WebSocket | Real-time updates |

---

## 📚 Documentation Files

| File | Description |
|------|-------------|
| `00-overview.md` | This file - Project overview |
| `01-database-schema.md` | Complete database schema |
| `02-implementation-phases.md` | Step-by-step implementation plan |
| `03-bunny-storage-setup.md` | Bunny Storage integration guide |
| `04-api-endpoints-checklist.md` | All 61 endpoints with examples |
| `05-testing-checklist.md` | Testing guidelines |
| `06-deployment.md` | Deployment instructions |

---

## ⚡ Quick Start

### 1. Read Documentation Order
1. 📖 Start here: `00-overview.md`
2. 🗄️ Database: `01-database-schema.md`
3. 🚀 Implementation: `02-implementation-phases.md`
4. ☁️ Storage: `03-bunny-storage-setup.md`
5. 📡 APIs: `04-api-endpoints-checklist.md`
6. ✅ Testing: `05-testing-checklist.md`

### 2. Setup Development Environment
```bash
# 1. Start PostgreSQL & Redis
docker-compose up -d

# 2. Copy environment variables
cp .env.example .env

# 3. Update .env with your credentials
# - Database credentials
# - JWT secret
# - Bunny Storage credentials

# 4. Run migrations
go run cmd/api/main.go

# 5. Server starts at http://localhost:3000
```

### 3. Test API
```bash
# Health check
curl http://localhost:3000/health

# Register user
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test1234",
    "displayName": "Test User"
  }'
```

---

## 🎓 Learning Resources

### Go Fiber Documentation
- Official Docs: https://docs.gofiber.io/
- GitHub: https://github.com/gofiber/fiber

### GORM Documentation
- Official Docs: https://gorm.io/docs/
- Guides: https://gorm.io/docs/guides.html

### Bunny CDN Documentation
- Storage API: https://docs.bunny.net/reference/storage-api
- CDN Guide: https://docs.bunny.net/docs/stream

### Backend Specification
- Complete API Spec: `../backend_spec/README.md`
- Error Codes: `../backend_spec/09-error-codes.md`

---

## 💡 Best Practices

### Code Organization
- ✅ Follow Clean Architecture principles
- ✅ Use dependency injection
- ✅ Keep business logic in services
- ✅ Use DTOs for API requests/responses
- ✅ Write tests for critical paths

### API Design
- ✅ RESTful endpoints
- ✅ Consistent response format
- ✅ Proper HTTP status codes
- ✅ Thai language error messages
- ✅ Pagination for list endpoints

### Security
- ✅ JWT token authentication
- ✅ Password hashing (bcrypt)
- ✅ Input validation
- ✅ Rate limiting
- ✅ CORS configuration

### Performance
- ✅ Database indexing
- ✅ Redis caching
- ✅ Pagination
- ✅ Eager loading for relationships
- ✅ CDN for media files

---

## 🚨 Important Notes

### Differences from Original System
- ❌ Remove: Task, Job, File models (replaced)
- ✅ Add: Post, Comment, Vote, Follow models
- ✅ Enhance: User model (karma, bio, followers)
- ✅ Change: Media upload (AWS → Bunny Storage)

### Critical Requirements
- 🔴 **Must use Bunny Storage** (not AWS S3)
- 🔴 **Response format** must match existing pattern
- 🔴 **Error messages in Thai** language
- 🔴 **JWT expiry** changed from 24h to 7 days
- 🔴 **Soft delete** for Posts and Comments

---

## 📞 Support

- 📖 Backend Spec: `../backend_spec/`
- 🐛 Issues: Create issues in GitHub
- 💬 Questions: Check documentation first

---

**Ready to start? Proceed to `01-database-schema.md`** →
