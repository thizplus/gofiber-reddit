# 🔧 CORS Error Fix - Production Deployment

## ปัญหา

เมื่อ deploy production แล้วเจอ CORS error:
```
Access to XMLHttpRequest at 'https://backend.voobize.com/api/v1/...'
from origin 'https://voobize.com' has been blocked by CORS policy:
Response to preflight request doesn't pass access control check:
The 'Access-Control-Allow-Origin' header contains the invalid value ''.
```

## สาเหตุ

Systemd service **ไม่ได้โหลด** environment variables จากไฟล์ `.env` ทำให้:
- `ALLOWED_ORIGINS` เป็น empty string
- Backend ส่ง CORS header ที่ไม่ถูกต้อง
- Browser block ทุก request

## วิธีแก้ไข

### Step 1: อัพเดท Systemd Service File

SSH เข้า production server:

```bash
ssh user@YOUR_SERVER_IP
```

แก้ไข systemd service file:

```bash
sudo nano /etc/systemd/system/voobize-backend.service
```

**เพิ่ม** บรรทัดนี้ในส่วน `[Service]`:

```ini
# Load environment variables from .env file
EnvironmentFile=/opt/voobize-backend/.env
```

**ไฟล์ที่ถูกต้องควรมีหน้าตาแบบนี้:**

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

### Step 2: ตรวจสอบ .env File

ตรวจสอบว่าไฟล์ `.env` มี `ALLOWED_ORIGINS`:

```bash
cd /opt/voobize-backend
cat .env | grep ALLOWED_ORIGINS
```

**ควรแสดง:**
```
ALLOWED_ORIGINS=https://voobize.com,https://www.voobize.com
```

ถ้า**ไม่มี** ให้เพิ่ม:

```bash
nano .env
```

เพิ่มบรรทัด:
```env
ALLOWED_ORIGINS=https://voobize.com,https://www.voobize.com
FRONTEND_URL=https://voobize.com
```

### Step 3: Restart Service

```bash
# Reload systemd configuration
sudo systemctl daemon-reload

# Restart backend service
sudo systemctl restart voobize-backend

# Check status
sudo systemctl status voobize-backend
```

### Step 4: ตรวจสอบ Logs

```bash
# ดู logs แบบ real-time
sudo journalctl -u voobize-backend -f

# ตรวจสอบว่าโหลด environment variables ได้
sudo journalctl -u voobize-backend | grep "ALLOWED_ORIGINS\|CORS\|Starting"
```

### Step 5: Test CORS

จาก browser console บน `https://voobize.com`:

```javascript
// Test CORS
fetch('https://backend.voobize.com/api/v1/health', {
  credentials: 'include'
})
  .then(res => res.json())
  .then(data => console.log('✅ CORS OK:', data))
  .catch(err => console.error('❌ CORS Error:', err));
```

**Expected Output:**
```
✅ CORS OK: {status: "ok", timestamp: "..."}
```

---

## ตรวจสอบเพิ่มเติม

### ดู Environment Variables ที่ Service ใช้

```bash
# ดูว่า service ใช้ env อะไรบ้าง
sudo systemctl show voobize-backend | grep Environment
```

### ทดสอบ CORS Headers

```bash
# ทดสอบ preflight request
curl -X OPTIONS https://backend.voobize.com/api/v1/health \
  -H "Origin: https://voobize.com" \
  -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v
```

**ควรเห็น headers:**
```
< Access-Control-Allow-Origin: https://voobize.com
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Methods: GET,POST,PUT,DELETE,PATCH,OPTIONS
```

### ถ้ายังไม่ได้

1. **ตรวจสอบ file permissions:**
   ```bash
   ls -la /opt/voobize-backend/.env
   # ควรเป็น readable โดย www-data
   sudo chmod 640 /opt/voobize-backend/.env
   sudo chown www-data:www-data /opt/voobize-backend/.env
   ```

2. **Rebuild binary:**
   ```bash
   cd /opt/voobize-backend
   go build -o bin/voobize-api cmd/api/main.go
   sudo systemctl restart voobize-backend
   ```

3. **ตรวจสอบ Nginx:**
   ```bash
   sudo nginx -t
   sudo systemctl reload nginx
   ```

---

## สรุป

การแก้ปัญหา CORS ใน production:

✅ เพิ่ม `EnvironmentFile=/opt/voobize-backend/.env` ใน systemd service
✅ ตรวจสอบว่า `.env` มี `ALLOWED_ORIGINS=https://voobize.com`
✅ Reload systemd และ restart service
✅ ทดสอบ CORS จาก browser

หลังจากแก้แล้ว frontend ควรเชื่อมต่อ backend ได้ปกติโดยไม่มี CORS error!
