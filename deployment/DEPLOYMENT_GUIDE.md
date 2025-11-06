# VOOBIZE Backend Deployment Guide

## 📋 Pre-Deployment Checklist

### 1. Domain & DNS Setup
- [ ] Domain: `backend.voobize.com` ชี้ไปยัง server IP
- [ ] SSL Certificate: ติดตั้ง Let's Encrypt (Certbot)
- [ ] DNS A Record: `backend.voobize.com → YOUR_SERVER_IP`

### 2. Environment Variables
- [ ] คัดลอก `.env.production` → `.env` บน server
- [ ] เปลี่ยน `DB_PASSWORD` เป็น production password
- [ ] เปลี่ยน `JWT_SECRET` เป็น strong secret
- [ ] เปลี่ยน `ALLOWED_ORIGINS` เป็น `https://voobize.com`
- [ ] เปลี่ยน `FRONTEND_URL` เป็น `https://voobize.com`
- [ ] เปลี่ยน `GOOGLE_REDIRECT_URL` เป็น `https://backend.voobize.com/...`

### 3. Database
- [ ] PostgreSQL ติดตั้งและ running
- [ ] Database `gofiber_social` สร้างแล้ว
- [ ] Run migrations (รวมถึง `add_push_subscriptions_unique_constraint.sql`)
- [ ] Database accessible จาก backend app

### 4. Server Dependencies
- [ ] Go 1.21+ ติดตั้งแล้ว
- [ ] PostgreSQL client (psql) ติดตั้งแล้ว
- [ ] Nginx ติดตั้งแล้ว
- [ ] Systemd service สำหรับ Go app

---

## 🚀 Deployment Steps

### Step 1: Clone & Build Backend

```bash
# SSH to server
ssh user@YOUR_SERVER_IP

# Clone repository
cd /opt
git clone YOUR_REPO_URL voobize-backend
cd voobize-backend

# Install dependencies
go mod download

# Build production binary
go build -o bin/voobize-api cmd/api/main.go

# Make executable
chmod +x bin/voobize-api
```

### Step 2: Setup Environment

```bash
# Copy production env
cp .env.production .env

# Edit environment variables
nano .env

# Update these values:
# - DB_PASSWORD
# - JWT_SECRET
# - REDIS_PASSWORD
# - ALLOWED_ORIGINS
# - FRONTEND_URL
# - GOOGLE_REDIRECT_URL
```

### Step 3: Run Database Migration

```bash
# Run push subscription migration
PGPASSWORD=YOUR_DB_PASSWORD psql -h localhost -p 5432 -U postgres -d gofiber_social -f migrations/add_push_subscriptions_unique_constraint.sql

# Verify constraint was added
PGPASSWORD=YOUR_DB_PASSWORD psql -h localhost -p 5432 -U postgres -d gofiber_social -c "SELECT conname FROM pg_constraint WHERE conrelid = 'push_subscriptions'::regclass;"
```

### Step 4: Setup Systemd Service

```bash
# Copy systemd service file
sudo cp deployment/voobize-backend.service /etc/systemd/system/

# Or create manually:
sudo nano /etc/systemd/system/voobize-backend.service
```

**Systemd Service File:**
```ini
[Unit]
Description=VOOBIZE Backend API
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/voobize-backend

# Load environment variables from .env file
EnvironmentFile=/opt/voobize-backend/.env

# Start the application
ExecStart=/opt/voobize-backend/bin/voobize-api

# Restart configuration
Restart=always
RestartSec=5

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=voobize-backend

# Security hardening
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

**⚠️ IMPORTANT**: The service file must include `EnvironmentFile=/opt/voobize-backend/.env` to load all environment variables including `ALLOWED_ORIGINS` for CORS to work properly!

**Enable and start service:**
```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service (auto-start on boot)
sudo systemctl enable voobize-backend

# Start service
sudo systemctl start voobize-backend

# Check status
sudo systemctl status voobize-backend

# View logs
sudo journalctl -u voobize-backend -f
```

### Step 5: Setup Nginx Reverse Proxy

```bash
# Copy nginx configuration
sudo cp deployment/nginx-backend.conf /etc/nginx/sites-available/backend.voobize.com

# Create symlink
sudo ln -s /etc/nginx/sites-available/backend.voobize.com /etc/nginx/sites-enabled/

# Test nginx configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

### Step 6: Setup SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Get SSL certificate
sudo certbot --nginx -d backend.voobize.com

# Certbot will automatically update nginx config

# Test auto-renewal
sudo certbot renew --dry-run

# Reload nginx
sudo systemctl reload nginx
```

---

## 🧪 Testing After Deployment

### 1. Test Backend Health

```bash
# HTTP (should redirect to HTTPS)
curl http://backend.voobize.com/health

# HTTPS
curl https://backend.voobize.com/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-01-06T12:00:00Z"
}
```

### 2. Test WebSocket (WSS)

**From browser console (on https://voobize.com):**

```javascript
// Test WebSocket connection
const ws = new WebSocket('wss://backend.voobize.com/ws?token=YOUR_JWT_TOKEN');

ws.onopen = () => {
  console.log('✅ WebSocket connected');
};

ws.onmessage = (event) => {
  console.log('📨 Message:', event.data);
};

ws.onerror = (error) => {
  console.error('❌ WebSocket error:', error);
};

ws.onclose = () => {
  console.log('🔌 WebSocket closed');
};
```

### 3. Test Push Notification API

```bash
# Get JWT token first (login)
TOKEN="YOUR_JWT_TOKEN_HERE"

# Test subscribe
curl -X POST https://backend.voobize.com/api/v1/push/subscribe \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "endpoint": "https://fcm.googleapis.com/fcm/send/test123",
    "keys": {
      "p256dh": "test-key",
      "auth": "test-auth"
    }
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Subscription saved successfully"
}
```

### 4. Test CORS

```bash
# Test CORS preflight
curl -X OPTIONS https://backend.voobize.com/api/v1/posts \
  -H "Origin: https://voobize.com" \
  -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: Authorization" \
  -v
```

**Expected Headers:**
```
Access-Control-Allow-Origin: https://voobize.com
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: GET,POST,PUT,DELETE,PATCH,OPTIONS
```

---

## 🔧 Troubleshooting

### Issue 1: WebSocket 502 Bad Gateway

**Symptoms:**
- WebSocket connection fails
- Nginx error: "upstream prematurely closed connection"

**Solutions:**
1. Check backend is running:
   ```bash
   sudo systemctl status voobize-backend
   ```

2. Check backend logs:
   ```bash
   sudo journalctl -u voobize-backend -n 100
   ```

3. Verify nginx WebSocket config:
   ```bash
   sudo nginx -T | grep -A 10 "location /ws"
   ```

4. Test backend directly (bypass nginx):
   ```bash
   # From server
   curl http://localhost:8080/health
   ```

### Issue 2: CORS Error

**Symptoms:**
- Browser console: "CORS policy: No 'Access-Control-Allow-Origin' header"

**Solutions:**
1. Check `ALLOWED_ORIGINS` in `.env`:
   ```bash
   cat .env | grep ALLOWED_ORIGINS
   ```

2. Restart backend:
   ```bash
   sudo systemctl restart voobize-backend
   ```

3. Check backend CORS middleware logs:
   ```bash
   sudo journalctl -u voobize-backend | grep CORS
   ```

### Issue 3: Push Notification ON CONFLICT Error

**Symptoms:**
- POST /api/v1/push/subscribe returns 500
- Error: "there is no unique or exclusion constraint"

**Solutions:**
1. Check constraint exists:
   ```sql
   SELECT conname FROM pg_constraint
   WHERE conrelid = 'push_subscriptions'::regclass
     AND conname = 'idx_user_endpoint';
   ```

2. If not exists, run migration:
   ```bash
   PGPASSWORD=YOUR_PASSWORD psql -d gofiber_social -f migrations/add_push_subscriptions_unique_constraint.sql
   ```

3. Restart backend:
   ```bash
   sudo systemctl restart voobize-backend
   ```

### Issue 4: SSL Certificate Issues

**Symptoms:**
- Browser: "Your connection is not private"
- curl: "SSL certificate problem"

**Solutions:**
1. Check certificate status:
   ```bash
   sudo certbot certificates
   ```

2. Renew certificate:
   ```bash
   sudo certbot renew
   sudo systemctl reload nginx
   ```

3. Check nginx SSL config:
   ```bash
   sudo nginx -T | grep ssl_certificate
   ```

---

## 📊 Monitoring

### View Backend Logs

```bash
# Real-time logs
sudo journalctl -u voobize-backend -f

# Last 100 lines
sudo journalctl -u voobize-backend -n 100

# Errors only
sudo journalctl -u voobize-backend -p err

# Today's logs
sudo journalctl -u voobize-backend --since today
```

### View Nginx Logs

```bash
# Access logs
sudo tail -f /var/log/nginx/backend.voobize.com.access.log

# Error logs
sudo tail -f /var/log/nginx/backend.voobize.com.error.log

# WebSocket connections
sudo grep "GET /ws" /var/log/nginx/backend.voobize.com.access.log
```

### Check Service Status

```bash
# Backend
sudo systemctl status voobize-backend

# Nginx
sudo systemctl status nginx

# PostgreSQL
sudo systemctl status postgresql

# All at once
sudo systemctl status voobize-backend nginx postgresql
```

---

## 🔐 Security Hardening

### 1. Firewall (UFW)

```bash
# Enable UFW
sudo ufw enable

# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# PostgreSQL (only from localhost - already restricted by default)
# sudo ufw deny 5432/tcp

# Check status
sudo ufw status
```

### 2. Fail2Ban (Optional)

```bash
# Install
sudo apt install fail2ban

# Create jail for nginx
sudo nano /etc/fail2ban/jail.d/nginx.conf
```

**nginx.conf:**
```ini
[nginx-http-auth]
enabled = true
port = http,https
logpath = /var/log/nginx/*error.log

[nginx-limit-req]
enabled = true
port = http,https
logpath = /var/log/nginx/*error.log
```

### 3. Environment Variables Security

```bash
# Secure .env file
sudo chown www-data:www-data /opt/voobize-backend/.env
sudo chmod 600 /opt/voobize-backend/.env

# Only www-data can read
ls -la /opt/voobize-backend/.env
# Should show: -rw------- 1 www-data www-data
```

---

## 🔄 Update Deployment

```bash
# SSH to server
ssh user@YOUR_SERVER_IP

# Go to project directory
cd /opt/voobize-backend

# Pull latest changes
git pull origin main

# Rebuild
go build -o bin/voobize-api cmd/api/main.go

# Restart service
sudo systemctl restart voobize-backend

# Check status
sudo systemctl status voobize-backend

# View logs
sudo journalctl -u voobize-backend -f
```

---

## 📝 Quick Reference

### Important Files
- Backend binary: `/opt/voobize-backend/bin/voobize-api`
- Environment: `/opt/voobize-backend/.env`
- Systemd service: `/etc/systemd/system/voobize-backend.service`
- Nginx config: `/etc/nginx/sites-available/backend.voobize.com`
- SSL certificate: `/etc/letsencrypt/live/backend.voobize.com/`

### Important Commands
```bash
# Restart backend
sudo systemctl restart voobize-backend

# View backend logs
sudo journalctl -u voobize-backend -f

# Restart nginx
sudo systemctl restart nginx

# Test nginx config
sudo nginx -t

# Renew SSL certificate
sudo certbot renew
```

### Important URLs
- API Base: `https://backend.voobize.com/api/v1`
- WebSocket: `wss://backend.voobize.com/ws`
- Health Check: `https://backend.voobize.com/health`
- Frontend: `https://voobize.com`

---

## 🆘 Support

หากพบปัญหา:
1. ตรวจสอบ logs ก่อน: `sudo journalctl -u voobize-backend -f`
2. ตรวจสอบ service status: `sudo systemctl status voobize-backend`
3. ตรวจสอบ nginx error logs: `sudo tail -f /var/log/nginx/backend.voobize.com.error.log`

---

**Last Updated:** 2025-01-06
**Version:** 1.0
