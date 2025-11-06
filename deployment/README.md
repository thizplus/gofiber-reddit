# VOOBIZE Deployment Files

## 📁 Files in this directory

### 1. `DEPLOYMENT_GUIDE.md` (Main Guide)
คำแนะนำทีละขั้นตอนสำหรับการ deploy backend ไปยัง production server

**Topics:**
- Pre-deployment checklist
- Step-by-step deployment
- Systemd service setup
- Nginx reverse proxy
- SSL certificate (Let's Encrypt)
- Testing & troubleshooting
- Monitoring & logging

### 2. `nginx-backend.conf`
Nginx configuration สำหรับ `backend.voobize.com`

**Features:**
- HTTP → HTTPS redirect
- WebSocket (WSS) support
- SSL termination
- Reverse proxy to Go backend (localhost:8080)
- Security headers
- CORS support

**Usage:**
```bash
sudo cp nginx-backend.conf /etc/nginx/sites-available/backend.voobize.com
sudo ln -s /etc/nginx/sites-available/backend.voobize.com /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 3. `FRONTEND_ENV_REFERENCE.md`
คำแนะนำสำหรับ Frontend Developer เรื่อง environment variables

**Topics:**
- Production environment variables
- API & WebSocket URLs
- Testing procedures
- Common issues & solutions
- Performance tips

---

## 🚀 Quick Start

### For Backend Developer:

1. **อ่าน:** `DEPLOYMENT_GUIDE.md`
2. **Setup:** Environment variables (`.env.production`)
3. **Deploy:** Follow step-by-step guide
4. **Configure:** Nginx (`nginx-backend.conf`)
5. **Setup SSL:** Let's Encrypt
6. **Test:** Health check, WebSocket, Push API

### For Frontend Developer:

1. **อ่าน:** `FRONTEND_ENV_REFERENCE.md`
2. **Setup:** Environment variables
   ```env
   NEXT_PUBLIC_API_URL=https://backend.voobize.com/api/v1
   NEXT_PUBLIC_WS_URL=wss://backend.voobize.com/ws
   NEXT_PUBLIC_VAPID_PUBLIC_KEY=BIC9GBi...
   ```
3. **Deploy:** Vercel / Own server
4. **Test:** API, WebSocket, Push Notifications

---

## 🔑 Important URLs

| Service | Development | Production |
|---------|-------------|-----------|
| Frontend | http://localhost:3000 | https://voobize.com |
| Backend API | http://localhost:8080/api/v1 | https://backend.voobize.com/api/v1 |
| WebSocket | ws://localhost:8080/ws | wss://backend.voobize.com/ws |
| Health Check | http://localhost:8080/health | https://backend.voobize.com/health |

---

## 🔐 Security Checklist

- [ ] SSL Certificate ติดตั้งแล้ว (Let's Encrypt)
- [ ] HTTPS Redirect เปิดใช้งาน
- [ ] CORS ตั้งค่าถูกต้อง (เฉพาะ production domain)
- [ ] JWT Secret เปลี่ยนจาก default
- [ ] Database Password strong enough
- [ ] `.env` file permissions: 600 (owner read/write only)
- [ ] Firewall (UFW) เปิดแค่ port 22, 80, 443
- [ ] Database ไม่ expose ออก internet (localhost only)
- [ ] Redis password ตั้งไว้

---

## 📊 Monitoring Commands

```bash
# Backend logs (real-time)
sudo journalctl -u voobize-backend -f

# Backend status
sudo systemctl status voobize-backend

# Nginx error logs
sudo tail -f /var/log/nginx/backend.voobize.com.error.log

# WebSocket connections
sudo grep "GET /ws" /var/log/nginx/backend.voobize.com.access.log

# System resources
htop
```

---

## 🔄 Update Workflow

```bash
# 1. SSH to server
ssh user@server

# 2. Pull latest code
cd /opt/voobize-backend
git pull origin main

# 3. Run migrations (if any)
PGPASSWORD=xxx psql -d gofiber_social -f migrations/xxx.sql

# 4. Rebuild
go build -o bin/voobize-api cmd/api/main.go

# 5. Restart
sudo systemctl restart voobize-backend

# 6. Check status
sudo systemctl status voobize-backend
sudo journalctl -u voobize-backend -n 50
```

---

## 🆘 Emergency Procedures

### Backend Down

```bash
# Check status
sudo systemctl status voobize-backend

# View recent logs
sudo journalctl -u voobize-backend -n 100 --no-pager

# Restart service
sudo systemctl restart voobize-backend

# If still failing, check process
ps aux | grep voobize

# Check port availability
sudo netstat -tulpn | grep 8080
```

### Database Issues

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Restart PostgreSQL
sudo systemctl restart postgresql

# Check connections
sudo -u postgres psql -c "SELECT * FROM pg_stat_activity WHERE datname = 'gofiber_social';"

# Kill idle connections
sudo -u postgres psql -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'gofiber_social' AND state = 'idle';"
```

### SSL Certificate Expired

```bash
# Check expiry
sudo certbot certificates

# Renew
sudo certbot renew --force-renewal

# Reload nginx
sudo systemctl reload nginx
```

### Nginx Issues

```bash
# Test config
sudo nginx -t

# Reload
sudo systemctl reload nginx

# Restart
sudo systemctl restart nginx

# Check error logs
sudo tail -n 100 /var/log/nginx/error.log
```

---

## 📞 Contact

- **Backend Developer:** [Your Name]
- **DevOps:** [DevOps Contact]
- **Server IP:** [YOUR_SERVER_IP]
- **SSH Access:** `ssh user@YOUR_SERVER_IP`

---

## 📚 Related Documentation

- [Backend API Specs](../backend_spec/README.md)
- [Push Notification API](../system_integration/PUSH_NOTIFICATION_API_SPECS.md)
- [Database Migrations](../migrations/README.md)

---

**Last Updated:** 2025-01-06
**Version:** 1.0
