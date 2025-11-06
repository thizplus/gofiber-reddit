# Frontend Environment Variables (Production)

## 📝 คำแนะนำสำหรับ Frontend Developer

เมื่อ deploy frontend ไป production ต้องตั้งค่า environment variables ดังนี้:

---

## Production Environment Variables

สร้างไฟล์ `.env.production` ใน Next.js project:

```env
# API Backend URL (HTTPS)
NEXT_PUBLIC_API_URL=https://backend.voobize.com/api/v1

# WebSocket URL (WSS - WebSocket Secure)
NEXT_PUBLIC_WS_URL=wss://backend.voobize.com/ws

# VAPID Public Key (for Push Notifications)
NEXT_PUBLIC_VAPID_PUBLIC_KEY=BIC9GBiayeWgHZXvxam9S1G_xCR5OYKA0NcfhXGhZ2KA3sNA4Wi5n38QXCUQV_jlN7yTd5bSyBNQe0NispxkKYk

# Google OAuth (if used on frontend)
# NEXT_PUBLIC_GOOGLE_CLIENT_ID=274539164677-j3lpqtctkr1kmbkfprb43fatuiq5og80.apps.googleusercontent.com

# Site URL (for SEO)
NEXT_PUBLIC_SITE_URL=https://voobize.com

# Analytics (if any)
# NEXT_PUBLIC_GA_ID=G-XXXXXXXXXX
```

---

## ตรวจสอบว่าถูกต้อง

### 1. API URL
```bash
# ต้องเป็น HTTPS (ไม่ใช่ HTTP)
✅ NEXT_PUBLIC_API_URL=https://backend.voobize.com/api/v1
❌ NEXT_PUBLIC_API_URL=http://backend.voobize.com/api/v1
```

### 2. WebSocket URL
```bash
# ต้องเป็น WSS (ไม่ใช่ WS)
✅ NEXT_PUBLIC_WS_URL=wss://backend.voobize.com/ws
❌ NEXT_PUBLIC_WS_URL=ws://backend.voobize.com/ws

# Path ต้องเป็น /ws (ตาม backend configuration)
✅ wss://backend.voobize.com/ws
❌ wss://backend.voobize.com/websocket
❌ wss://backend.voobize.com/api/v1/ws
```

### 3. VAPID Public Key
```bash
# ใช้ key เดียวกับที่ backend ตั้งไว้
✅ NEXT_PUBLIC_VAPID_PUBLIC_KEY=BIC9GBiayeWgHZXvxam9S1G_xCR5OYKA0NcfhXGhZ2KA3sNA4Wi5n38QXCUQV_jlN7yTd5bSyBNQe0NispxkKYk
```

---

## Build & Deploy Frontend

### Vercel Deployment

```bash
# Install Vercel CLI (if not installed)
npm i -g vercel

# Login
vercel login

# Deploy to production
vercel --prod

# Environment variables จะถูก set ผ่าน Vercel Dashboard:
# Project Settings → Environment Variables
```

**ใน Vercel Dashboard:**
1. เข้า Project Settings
2. คลิก Environment Variables
3. เพิ่ม variables ทั้งหมดจาก `.env.production`
4. เลือก Environment: **Production**
5. Save

### Manual Build (if deploying to own server)

```bash
# Build Next.js for production
npm run build

# Start production server
npm run start

# หรือใช้ PM2
pm2 start npm --name "voobize-frontend" -- start
```

---

## ทดสอบหลัง Deploy

### 1. ทดสอบ API Connection

เปิด browser console ที่ `https://voobize.com`:

```javascript
// ทดสอบ API
fetch('https://backend.voobize.com/api/v1/health')
  .then(res => res.json())
  .then(data => console.log('API Health:', data));

// Expected: { status: "ok", timestamp: "..." }
```

### 2. ทดสอบ WebSocket

```javascript
// ทดสอบ WebSocket
const ws = new WebSocket('wss://backend.voobize.com/ws');

ws.onopen = () => console.log('✅ WebSocket connected');
ws.onerror = (err) => console.error('❌ WebSocket error:', err);
ws.onclose = () => console.log('🔌 WebSocket closed');

// Expected: ✅ WebSocket connected
```

### 3. ทดสอบ CORS

```javascript
// ทดสอบ CORS (with credentials)
fetch('https://backend.voobize.com/api/v1/posts', {
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json'
  }
})
  .then(res => {
    console.log('CORS OK:', res.status);
    return res.json();
  })
  .then(data => console.log('Posts:', data));

// Expected: CORS OK: 200
```

### 4. ทดสอบ Push Notifications

```javascript
// ทดสอบ Push Notification subscription
if ('serviceWorker' in navigator && 'PushManager' in window) {
  navigator.serviceWorker.ready.then(async (registration) => {
    // Subscribe
    const subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: 'BIC9GBiayeWgHZXvxam9S1G_xCR5OYKA0NcfhXGhZ2KA3sNA4Wi5n38QXCUQV_jlN7yTd5bSyBNQe0NispxkKYk'
    });

    console.log('✅ Push subscription:', subscription);

    // Send to backend
    const token = 'YOUR_JWT_TOKEN'; // Get from auth
    const response = await fetch('https://backend.voobize.com/api/v1/push/subscribe', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(subscription)
    });

    const result = await response.json();
    console.log('Backend response:', result);
    // Expected: { success: true, message: "Subscription saved successfully" }
  });
}
```

---

## Common Issues

### Issue 1: Mixed Content Error

**Symptoms:**
```
Mixed Content: The page at 'https://voobize.com' was loaded over HTTPS,
but requested an insecure resource 'http://backend.voobize.com/...'.
This request has been blocked.
```

**Solution:**
- ตรวจสอบว่า `NEXT_PUBLIC_API_URL` เป็น `https://` (ไม่ใช่ `http://`)
- ตรวจสอบว่า `NEXT_PUBLIC_WS_URL` เป็น `wss://` (ไม่ใช่ `ws://`)

### Issue 2: WebSocket Connection Failed

**Symptoms:**
```
WebSocket connection to 'wss://backend.voobize.com/ws' failed
```

**Solutions:**
1. ตรวจสอบ backend ทำงานอยู่:
   ```bash
   curl https://backend.voobize.com/health
   ```

2. ตรวจสอบ Nginx WebSocket config
3. ตรวจสอบ SSL certificate ใช้งานได้
4. ลองเปิด browser console ดู error detail

### Issue 3: CORS Error

**Symptoms:**
```
Access to fetch at 'https://backend.voobize.com/api/v1/posts' from origin
'https://voobize.com' has been blocked by CORS policy
```

**Solutions:**
1. ตรวจสอบ backend `.env`:
   ```
   ALLOWED_ORIGINS=https://voobize.com
   ```

2. Restart backend service:
   ```bash
   sudo systemctl restart voobize-backend
   ```

3. ตรวจสอบ CORS headers:
   ```bash
   curl -I -X OPTIONS https://backend.voobize.com/api/v1/posts \
     -H "Origin: https://voobize.com"
   ```

### Issue 4: Push Notification Permission Denied

**Symptoms:**
- User คลิก "Allow" แต่ได้ error
- Browser console: "Registration failed"

**Solutions:**
1. ตรวจสอบ Service Worker registered:
   ```javascript
   navigator.serviceWorker.getRegistrations()
     .then(regs => console.log('SW:', regs));
   ```

2. ตรวจสอบ VAPID key ถูกต้อง
3. ลอง unregister SW แล้ว register ใหม่:
   ```javascript
   navigator.serviceWorker.getRegistrations()
     .then(regs => regs.forEach(reg => reg.unregister()));
   ```

---

## Performance Tips

### 1. CDN (Bunny CDN)

หากใช้ Bunny CDN สำหรับ images/media:

```env
# Use CDN URL instead of storage URL
NEXT_PUBLIC_CDN_URL=https://voobizethailand.b-cdn.net
```

### 2. Image Optimization

```javascript
// ใช้ Next.js Image component
import Image from 'next/image';

<Image
  src={`${process.env.NEXT_PUBLIC_CDN_URL}/uploads/image.jpg`}
  alt="..."
  width={800}
  height={600}
  loading="lazy"
/>
```

### 3. WebSocket Reconnection

```javascript
// Implement auto-reconnect
let ws;
let reconnectInterval = 1000;

function connectWebSocket() {
  ws = new WebSocket(process.env.NEXT_PUBLIC_WS_URL);

  ws.onopen = () => {
    console.log('✅ WebSocket connected');
    reconnectInterval = 1000; // Reset interval
  };

  ws.onerror = (err) => {
    console.error('❌ WebSocket error:', err);
  };

  ws.onclose = () => {
    console.log('🔌 WebSocket closed, reconnecting...');
    setTimeout(() => {
      reconnectInterval *= 2; // Exponential backoff
      if (reconnectInterval > 30000) reconnectInterval = 30000; // Max 30s
      connectWebSocket();
    }, reconnectInterval);
  };
}

connectWebSocket();
```

---

## Checklist ก่อน Deploy

- [ ] `NEXT_PUBLIC_API_URL` เป็น HTTPS
- [ ] `NEXT_PUBLIC_WS_URL` เป็น WSS
- [ ] `NEXT_PUBLIC_VAPID_PUBLIC_KEY` ตรงกับ backend
- [ ] Backend CORS ตั้งค่า `ALLOWED_ORIGINS` ถูกต้อง
- [ ] Service Worker ทำงานบน HTTPS
- [ ] ทดสอบ API connection
- [ ] ทดสอบ WebSocket connection
- [ ] ทดสอบ Push Notifications
- [ ] ทดสอบ image loading (CDN)

---

**Last Updated:** 2025-01-06
**Version:** 1.0
