# Deployment Guide

## 1. Overview

ระบบนี้สามารถ deploy ได้หลายวิธี:
- **Docker Compose** (แนะนำสำหรับ development & staging)
- **Docker** (single container)
- **Native** (build & run binary)
- **Kubernetes** (สำหรับ production scaling)

---

## 2. Prerequisites

### 2.1 Required Software

**สำหรับ Docker Deployment**:
- Docker Engine 20.10+
- Docker Compose v2.0+

**สำหรับ Native Deployment**:
- Go 1.24+
- PostgreSQL 15+
- Redis 7+

### 2.2 Required Services

**External Services**:
- **Bunny CDN**: สำหรับเก็บ media files
  - Storage Zone
  - Access Key
- **Google OAuth** (ถ้าใช้ OAuth):
  - Client ID
  - Client Secret
  - Redirect URL

---

## 3. Environment Configuration

### 3.1 Create Environment File

สร้างไฟล์ `.env` จาก `.env.example`:

```bash
cp .env.example .env
```

### 3.2 Environment Variables

**Required Variables**:

```bash
# ========================================
# Application Configuration
# ========================================
APP_NAME=Social Media Platform
APP_PORT=3000
APP_ENV=production  # development, staging, production
FRONTEND_URL=https://yourfrontend.com

# ========================================
# Database Configuration (PostgreSQL)
# ========================================
DB_HOST=postgres  # or localhost for native
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password_here
DB_NAME=social_platform
DB_SSL_MODE=disable  # require for production

# ========================================
# Redis Configuration
# ========================================
REDIS_HOST=redis  # or localhost for native
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password  # optional
REDIS_DB=0

# ========================================
# JWT Configuration
# ========================================
JWT_SECRET=your-super-secret-jwt-key-change-this
JWT_EXPIRY=24h  # 24 hours

# ========================================
# Bunny CDN Configuration
# ========================================
BUNNY_STORAGE_ZONE=your-storage-zone
BUNNY_ACCESS_KEY=your-access-key-here
BUNNY_BASE_URL=https://storage.bunnycdn.com
BUNNY_CDN_URL=https://your-cdn-zone.b-cdn.net

# ========================================
# Google OAuth Configuration (Optional)
# ========================================
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=https://api.yourdomain.com/api/v1/auth/google/callback

# ========================================
# Web Push (VAPID) Configuration
# ========================================
VAPID_PUBLIC_KEY=your-vapid-public-key
VAPID_PRIVATE_KEY=your-vapid-private-key
VAPID_SUBJECT=mailto:admin@yourdomain.com

# ========================================
# CORS Configuration
# ========================================
ALLOWED_ORIGINS=https://yourfrontend.com,https://app.yourdomain.com

# ========================================
# Optional: Email Configuration (Planned)
# ========================================
# SMTP_HOST=smtp.gmail.com
# SMTP_PORT=587
# SMTP_USER=your-email@gmail.com
# SMTP_PASSWORD=your-app-password
```

### 3.3 Generate VAPID Keys

สำหรับ Web Push notifications:

```bash
# ใช้ web-push library (Node.js)
npx web-push generate-vapid-keys

# หรือใช้ online tool
# https://vapidkeys.com/
```

### 3.4 Security Checklist

**Before Production**:
- [ ] เปลี่ยน `JWT_SECRET` เป็นค่าที่ปลอดภัย (min 32 characters)
- [ ] ใช้ strong passwords สำหรับ database และ Redis
- [ ] เปิด SSL mode สำหรับ database (`DB_SSL_MODE=require`)
- [ ] ตั้งค่า `ALLOWED_ORIGINS` ให้ถูกต้อง
- [ ] ใช้ HTTPS สำหรับ production
- [ ] ปิด debug mode (`APP_ENV=production`)
- [ ] ตั้งค่า firewall (เปิดเฉพาะ port ที่จำเป็น)
- [ ] เปิด rate limiting (planned)

---

## 4. Docker Compose Deployment (แนะนำ)

### 4.1 Quick Start

```bash
# 1. Clone repository
git clone https://github.com/yourusername/gofiber-backend.git
cd gofiber-backend

# 2. Create .env file
cp .env.example .env
nano .env  # แก้ไข environment variables

# 3. Start services
docker-compose up -d

# 4. Check logs
docker-compose logs -f app

# 5. Check health
curl http://localhost:8080/health
```

### 4.2 Docker Compose File

**docker-compose.yml**:

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: social-app
    ports:
      - "8080:3000"
    environment:
      # Load from .env file
      - APP_NAME=${APP_NAME}
      - APP_PORT=${APP_PORT}
      - APP_ENV=${APP_ENV}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      # ... (all other env vars)
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    restart: unless-stopped
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:15-alpine
    container_name: social-postgres
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: social-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped
    networks:
      - app-network
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

networks:
  app-network:
    driver: bridge

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local
```

### 4.3 Docker Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Restart services
docker-compose restart

# View logs
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f redis

# Execute commands in container
docker-compose exec app sh

# Check service status
docker-compose ps

# Rebuild after code changes
docker-compose up -d --build

# Clean up (WARNING: deletes volumes)
docker-compose down -v
```

---

## 5. Native Deployment

### 5.1 Install Dependencies

**PostgreSQL**:
```bash
# Ubuntu/Debian
sudo apt-get install postgresql-15

# macOS
brew install postgresql@15

# Start PostgreSQL
sudo systemctl start postgresql  # Linux
brew services start postgresql@15  # macOS
```

**Redis**:
```bash
# Ubuntu/Debian
sudo apt-get install redis-server

# macOS
brew install redis

# Start Redis
sudo systemctl start redis  # Linux
brew services start redis  # macOS
```

### 5.2 Setup Database

```bash
# Create database
psql -U postgres -c "CREATE DATABASE social_platform;"

# Create user (optional)
psql -U postgres -c "CREATE USER social_user WITH PASSWORD 'your_password';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE social_platform TO social_user;"
```

### 5.3 Build & Run

```bash
# 1. Install Go dependencies
go mod download

# 2. Build binary
go build -o bin/app cmd/api/main.go

# 3. Run application
./bin/app

# Or run directly
go run cmd/api/main.go
```

### 5.4 Run as Service (Linux)

**Create systemd service** (`/etc/systemd/system/social-app.service`):

```ini
[Unit]
Description=Social Media Platform API
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/social-app
EnvironmentFile=/var/www/social-app/.env
ExecStart=/var/www/social-app/bin/app
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

**Enable & start service**:
```bash
sudo systemctl daemon-reload
sudo systemctl enable social-app
sudo systemctl start social-app
sudo systemctl status social-app
```

---

## 6. Production Deployment

### 6.1 Architecture (Single Server)

```
Internet
    │
    ↓
┌───────────┐
│  Nginx    │  (Reverse Proxy, SSL, Load Balancer)
│  Port 80  │
│  Port 443 │
└─────┬─────┘
      │
      ↓
┌───────────┐
│  Go App   │  (Fiber Application)
│  Port 3000│
└─────┬─────┘
      │
      ├──────────┐
      ↓          ↓
┌──────────┐  ┌──────────┐
│PostgreSQL│  │  Redis   │
│ Port 5432│  │ Port 6379│
└──────────┘  └──────────┘
      │
      ↓
┌──────────┐
│ Bunny CDN│  (Media Storage)
└──────────┘
```

### 6.2 Nginx Configuration

**Install Nginx**:
```bash
sudo apt-get install nginx
```

**Nginx config** (`/etc/nginx/sites-available/social-app`):

```nginx
# Upstream (Go Fiber app)
upstream go_backend {
    server localhost:3000;
    # Add more servers for load balancing
    # server localhost:3001;
    # server localhost:3002;
}

# HTTP (redirect to HTTPS)
server {
    listen 80;
    listen [::]:80;
    server_name api.yourdomain.com;

    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name api.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Logging
    access_log /var/log/nginx/social-app-access.log;
    error_log /var/log/nginx/social-app-error.log;

    # Max upload size (match Fiber config)
    client_max_body_size 300M;

    # Proxy settings
    location / {
        proxy_pass http://go_backend;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # WebSocket support
    location /ws {
        proxy_pass http://go_backend;
        proxy_http_version 1.1;

        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;

        # WebSocket timeouts
        proxy_read_timeout 86400s;
        proxy_send_timeout 86400s;
    }

    # Health check (bypass proxy for faster response)
    location /health {
        proxy_pass http://go_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        access_log off;
    }
}
```

**Enable site**:
```bash
sudo ln -s /etc/nginx/sites-available/social-app /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 6.3 SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d api.yourdomain.com

# Auto-renewal (already set up by certbot)
sudo certbot renew --dry-run
```

### 6.4 Firewall Setup

```bash
# UFW (Ubuntu)
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# Optional: allow direct access to app (for debugging)
# sudo ufw allow 3000/tcp
```

---

## 7. Database Management

### 7.1 Backup Strategy

**Automated Daily Backup**:

```bash
#!/bin/bash
# /usr/local/bin/backup-db.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/postgres"
DB_NAME="social_platform"
DB_USER="postgres"

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup database
pg_dump -U $DB_USER $DB_NAME | gzip > $BACKUP_DIR/backup_$DATE.sql.gz

# Keep only last 30 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete

echo "Backup completed: backup_$DATE.sql.gz"
```

**Cron job** (run daily at 2 AM):
```bash
crontab -e

# Add this line:
0 2 * * * /usr/local/bin/backup-db.sh >> /var/log/backup-db.log 2>&1
```

### 7.2 Restore Database

```bash
# Restore from backup
gunzip < backup_20240101_020000.sql.gz | psql -U postgres social_platform

# Or direct restore
psql -U postgres social_platform < backup.sql
```

### 7.3 Database Migrations

**Auto-migration** (on app startup):
- GORM auto-migrates models
- Safe for development
- Use manual migrations for production

**Manual migrations**:
```bash
# Location: infrastructure/postgres/migrations/

# Create migration
# 20240101000000_add_new_feature.sql

# Apply manually
psql -U postgres social_platform < migrations/20240101000000_add_new_feature.sql
```

---

## 8. Monitoring & Logging

### 8.1 Application Logs

**Docker Compose**:
```bash
# View logs
docker-compose logs -f app

# Save logs to file
docker-compose logs app > app.log
```

**Native**:
```bash
# Redirect to file
./bin/app >> /var/log/social-app/app.log 2>&1

# Use log rotation
sudo apt-get install logrotate

# Config: /etc/logrotate.d/social-app
/var/log/social-app/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data www-data
    sharedscripts
    postrotate
        systemctl reload social-app
    endscript
}
```

### 8.2 Monitoring Tools (Recommended)

**Prometheus + Grafana**:
- Metrics collection
- Real-time dashboards
- Alerting

**Sentry**:
- Error tracking
- Performance monitoring

**Uptime Robot**:
- Uptime monitoring
- Alert on downtime

### 8.3 Health Checks

```bash
# HTTP health check
curl http://localhost:8080/health

# WebSocket health check
wscat -c ws://localhost:8080/ws

# Database check
psql -U postgres -c "SELECT 1;"

# Redis check
redis-cli ping
```

---

## 9. Scaling

### 9.1 Horizontal Scaling

**Load Balancer** (Nginx):
```nginx
upstream go_backend {
    least_conn;  # Load balancing method
    server app1:3000;
    server app2:3000;
    server app3:3000;
}
```

**Session Management**:
- Use Redis for session storage (already implemented)
- Stateless JWT tokens (already implemented)

### 9.2 Database Scaling

**Read Replicas**:
```go
// GORM supports read/write splitting
db, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{})
db.Use(dbresolver.Register(dbresolver.Config{
    Replicas: []gorm.Dialector{
        postgres.Open(replica1DSN),
        postgres.Open(replica2DSN),
    },
}))
```

**Connection Pooling**:
```go
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(100)
sqlDB.SetMaxIdleConns(10)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 9.3 Redis Scaling

**Redis Cluster** (multi-node):
- Horizontal partitioning
- High availability

**Redis Sentinel** (master-slave):
- Automatic failover
- Read scaling

---

## 10. Troubleshooting

### 10.1 Common Issues

**Cannot connect to database**:
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Check connection
psql -U postgres -h localhost -p 5432

# Check logs
sudo tail -f /var/log/postgresql/postgresql-15-main.log
```

**Cannot connect to Redis**:
```bash
# Check Redis is running
sudo systemctl status redis

# Check connection
redis-cli ping

# Check logs
sudo tail -f /var/log/redis/redis-server.log
```

**Port already in use**:
```bash
# Find process using port
sudo lsof -i :3000

# Kill process
sudo kill -9 <PID>
```

**High memory usage**:
```bash
# Check memory
free -h

# Check Go app memory
docker stats social-app

# Adjust Go GC
GOGC=50 ./bin/app  # More aggressive GC
```

### 10.2 Performance Tuning

**PostgreSQL**:
```sql
-- Check slow queries
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;

-- Add indexes
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
```

**Redis**:
```bash
# Monitor Redis
redis-cli --stat

# Check memory usage
redis-cli INFO memory
```

**Go App**:
```bash
# Profile CPU
go tool pprof http://localhost:3000/debug/pprof/profile

# Profile memory
go tool pprof http://localhost:3000/debug/pprof/heap
```

---

## 11. Deployment Checklist

### Pre-Deployment
- [ ] Update `.env` with production values
- [ ] Generate secure JWT secret
- [ ] Setup Bunny CDN
- [ ] Setup Google OAuth (if using)
- [ ] Generate VAPID keys
- [ ] Create database backup
- [ ] Test all API endpoints
- [ ] Run security audit

### Deployment
- [ ] Build Docker image
- [ ] Push to registry (if using)
- [ ] Deploy to server
- [ ] Run database migrations
- [ ] Configure Nginx
- [ ] Setup SSL certificate
- [ ] Configure firewall
- [ ] Setup monitoring

### Post-Deployment
- [ ] Verify health check
- [ ] Test critical endpoints
- [ ] Check logs for errors
- [ ] Monitor performance
- [ ] Setup automated backups
- [ ] Document deployment process
- [ ] Create rollback plan

---

## 12. Rollback Strategy

**Docker Compose**:
```bash
# Tag images before deployment
docker tag social-app:latest social-app:v1.0.0

# Rollback
docker-compose down
docker tag social-app:v1.0.0 social-app:latest
docker-compose up -d
```

**Database Rollback**:
```bash
# Restore from backup
psql -U postgres social_platform < backup_before_deploy.sql
```

**Quick Rollback Script**:
```bash
#!/bin/bash
# rollback.sh

echo "Rolling back to previous version..."
docker-compose down
docker tag social-app:previous social-app:latest
docker-compose up -d
echo "Rollback completed!"
```

---

## 13. Additional Resources

- [Go Fiber Documentation](https://docs.gofiber.io/)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Let's Encrypt](https://letsencrypt.org/)
