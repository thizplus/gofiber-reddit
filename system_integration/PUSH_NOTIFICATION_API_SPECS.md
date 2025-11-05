# Push Notification API Specifications

## 📋 Overview

ระบบ Push Notification สำหรับ VOOBIZE ใช้ Web Push API (VAPID) เพื่อส่งการแจ้งเตือนไปยัง users แม้ว่าจะปิดเว็บไซต์แล้ว

**Strategy:** Hybrid Notifications
- **WebSocket** → Real-time สำหรับ users ที่ออนไลน์อยู่
- **Push Notifications** → สำหรับ users ที่ปิดเว็บ/app

---

## 🔑 VAPID Keys (เก็บใน Backend ENV)

```env
VAPID_PUBLIC_KEY=BIC9GBiayeWgHZXvxam9S1G_xCR5OYKA0NcfhXGhZ2KA3sNA4Wi5n38QXCUQV_jlN7yTd5bSyBNQe0NispxkKYk
VAPID_PRIVATE_KEY=imjcz4yam3RM1WqztZfrNrXpexrtwudhsS7rs9_Xniw
VAPID_SUBJECT=mailto:admin@voobize.com
```

---

## 💾 Database Schema

### Table: `push_subscriptions`

```sql
CREATE TABLE push_subscriptions (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  endpoint TEXT NOT NULL,
  p256dh TEXT NOT NULL,
  auth TEXT NOT NULL,
  expiration_time BIGINT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_id, endpoint)
);

CREATE INDEX idx_push_subscriptions_user_id ON push_subscriptions(user_id);
CREATE INDEX idx_push_subscriptions_endpoint ON push_subscriptions(endpoint);
```

**หมายเหตุ:**
- `endpoint` - URL สำหรับส่ง push notification
- `p256dh`, `auth` - Encryption keys จาก browser
- User 1 คน สามารถมีหลาย subscriptions ได้ (หลาย devices)

---

## 🛠️ API Endpoints

### 1. POST `/api/v1/push/subscribe`

**บันทึก push subscription จาก user**

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "endpoint": "https://fcm.googleapis.com/fcm/send/xxxx",
  "expirationTime": null,
  "keys": {
    "p256dh": "BDjASz8xxxxxxxxxx...",
    "auth": "xcvbnmxxxxxxxxxx..."
  }
}
```

**Response - Success (200):**
```json
{
  "success": true,
  "message": "Subscription saved successfully",
  "data": {
    "id": 123,
    "userId": 1,
    "endpoint": "https://fcm.googleapis.com/fcm/send/xxxx"
  }
}
```

**Response - Error (400):**
```json
{
  "success": false,
  "error": "Invalid subscription data"
}
```

**Business Logic:**
1. ✅ ตรวจสอบ auth token (ดึง user ID)
2. ✅ Validate subscription format
3. ✅ เช็คว่า subscription นี้มีอยู่แล้วหรือยัง (endpoint + user_id)
4. ✅ ถ้ามี → UPDATE updated_at
5. ✅ ถ้าไม่มี → INSERT ใหม่
6. ✅ Return success response

---

### 2. POST `/api/v1/push/unsubscribe`

**ลบ push subscription**

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "endpoint": "https://fcm.googleapis.com/fcm/send/xxxx",
  "expirationTime": null,
  "keys": {
    "p256dh": "BDjASz8xxxxxxxxxx...",
    "auth": "xcvbnmxxxxxxxxxx..."
  }
}
```

**Response - Success (200):**
```json
{
  "success": true,
  "message": "Subscription removed successfully"
}
```

**Business Logic:**
1. ✅ ตรวจสอบ auth token
2. ✅ DELETE subscription WHERE user_id = ? AND endpoint = ?
3. ✅ Return success (แม้ไม่เจอก็ return success)

---

## 📤 Sending Push Notifications

### Go Example (using webpush-go)

**Install:**
```bash
go get github.com/SherClockHolmes/webpush-go
```

**Code:**
```go
package notification

import (
    "encoding/json"
    "log"
    "os"

    webpush "github.com/SherClockHolmes/webpush-go"
)

// PushSubscription represents a user's push subscription
type PushSubscription struct {
    Endpoint       string `json:"endpoint"`
    P256dh         string `json:"p256dh"`
    Auth           string `json:"auth"`
    ExpirationTime *int64 `json:"expirationTime"`
}

// NotificationPayload is the data sent to the user
type NotificationPayload struct {
    Title string                 `json:"title"`
    Body  string                 `json:"body"`
    Icon  string                 `json:"icon"`
    Badge string                 `json:"badge"`
    Tag   string                 `json:"tag"`
    Data  map[string]interface{} `json:"data"`
}

// SendPushNotification ส่ง push notification ไปยัง user
func SendPushNotification(subscription PushSubscription, payload NotificationPayload) error {
    // Convert payload to JSON
    payloadJSON, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    // Create webpush subscription
    s := &webpush.Subscription{
        Endpoint: subscription.Endpoint,
        Keys: webpush.Keys{
            P256dh: subscription.P256dh,
            Auth:   subscription.Auth,
        },
    }

    // Send notification
    resp, err := webpush.SendNotification(payloadJSON, s, &webpush.Options{
        Subscriber:      os.Getenv("VAPID_SUBJECT"), // mailto:admin@voobize.com
        VAPIDPublicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
        VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE_KEY"),
        TTL:             30,
    })

    if err != nil {
        log.Printf("❌ Push notification error: %v", err)
        return err
    }
    defer resp.Body.Close()

    // Check response status
    if resp.StatusCode != 201 {
        log.Printf("⚠️ Push notification failed with status: %d", resp.StatusCode)

        // ถ้าเป็น 410 Gone หรือ 404 Not Found = subscription หมดอายุ
        if resp.StatusCode == 410 || resp.StatusCode == 404 {
            // ลบ subscription ออกจาก database
            log.Printf("🗑️ Removing expired subscription: %s", subscription.Endpoint)
            // TODO: DeleteSubscription(subscription.Endpoint)
        }
        return nil
    }

    log.Printf("✅ Push notification sent successfully to: %s", subscription.Endpoint)
    return nil
}

// NotifyUser ส่ง notification ทั้ง WebSocket และ Push
func NotifyUser(userID int, notificationType string, message string, url string) error {
    // 1. ส่งผ่าน WebSocket (ถ้าออนไลน์)
    if websocket.IsUserOnline(userID) {
        websocket.SendToUser(userID, map[string]interface{}{
            "type":    notificationType,
            "message": message,
            "url":     url,
        })
    }

    // 2. ดึง push subscriptions ของ user
    subscriptions, err := GetUserSubscriptions(userID)
    if err != nil {
        log.Printf("Error getting subscriptions: %v", err)
        return err
    }

    // 3. ส่ง push notification ไปทุก subscriptions
    payload := NotificationPayload{
        Title: "VOOBIZE",
        Body:  message,
        Icon:  "/logo.png",
        Badge: "/logo.png",
        Tag:   notificationType,
        Data: map[string]interface{}{
            "url": url,
        },
    }

    for _, sub := range subscriptions {
        go func(s PushSubscription) {
            if err := SendPushNotification(s, payload); err != nil {
                log.Printf("Failed to send push to subscription: %v", err)
            }
        }(sub)
    }

    return nil
}
```

---

## 🔔 Notification Triggers

**เมื่อควรส่ง Notification:**

### 1. Comment (ความคิดเห็นใหม่)
```go
// เมื่อมีคนแสดงความคิดเห็นในโพสต์ของ user
func OnNewComment(postAuthorID int, commenterName string, postID int) {
    NotifyUser(
        postAuthorID,
        "comment",
        fmt.Sprintf("%s แสดงความคิดเห็นในโพสต์ของคุณ", commenterName),
        fmt.Sprintf("/post/%d", postID),
    )
}
```

### 2. Like (กดถูกใจ)
```go
// เมื่อมีคนกดถูกใจโพสต์
func OnNewLike(postAuthorID int, likerName string, postID int) {
    NotifyUser(
        postAuthorID,
        "like",
        fmt.Sprintf("%s กดถูกใจโพสต์ของคุณ", likerName),
        fmt.Sprintf("/post/%d", postID),
    )
}
```

### 3. Follow (ติดตามใหม่)
```go
// เมื่อมีคนติดตาม
func OnNewFollower(userID int, followerName string) {
    NotifyUser(
        userID,
        "follow",
        fmt.Sprintf("%s เริ่มติดตามคุณแล้ว", followerName),
        fmt.Sprintf("/profile/%s", followerName),
    )
}
```

### 4. Mention (@mention)
```go
// เมื่อมีคน mention ในโพสต์หรือความคิดเห็น
func OnMention(mentionedUserID int, mentionerName string, postID int) {
    NotifyUser(
        mentionedUserID,
        "mention",
        fmt.Sprintf("%s กล่าวถึงคุณในโพสต์", mentionerName),
        fmt.Sprintf("/post/%d", postID),
    )
}
```

---

## 🧪 Testing

### 1. Test Subscribe

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/push/subscribe \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "endpoint": "https://fcm.googleapis.com/fcm/send/test123",
    "keys": {
      "p256dh": "test-p256dh-key",
      "auth": "test-auth-key"
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

### 2. Test Sending Push Notification

**Go Test:**
```go
func TestSendPush() {
    subscription := PushSubscription{
        Endpoint: "https://fcm.googleapis.com/fcm/send/...",
        P256dh:   "...",
        Auth:     "...",
    }

    payload := NotificationPayload{
        Title: "Test Notification",
        Body:  "This is a test",
        Icon:  "/logo.png",
        Data: map[string]interface{}{
            "url": "/notifications",
        },
    }

    err := SendPushNotification(subscription, payload)
    if err != nil {
        t.Errorf("Failed to send push: %v", err)
    }
}
```

### 3. Test from Frontend

**Browser Console:**
```javascript
// ดู subscription ที่ส่งไป
navigator.serviceWorker.ready.then(reg => {
  reg.pushManager.getSubscription().then(sub => {
    console.log('Current subscription:', JSON.stringify(sub));
  });
});
```

---

## 📊 Database Queries

### Get User Subscriptions
```sql
SELECT endpoint, p256dh, auth, expiration_time
FROM push_subscriptions
WHERE user_id = $1
  AND (expiration_time IS NULL OR expiration_time > EXTRACT(EPOCH FROM NOW()) * 1000);
```

### Clean Expired Subscriptions
```sql
DELETE FROM push_subscriptions
WHERE expiration_time IS NOT NULL
  AND expiration_time < EXTRACT(EPOCH FROM NOW()) * 1000;
```

### Count Active Subscriptions
```sql
SELECT user_id, COUNT(*) as subscription_count
FROM push_subscriptions
GROUP BY user_id;
```

---

## ⚠️ Error Handling

### 410 Gone / 404 Not Found
- Subscription หมดอายุหรือถูกลบแล้ว
- **Action:** ลบ subscription จาก database

### 429 Too Many Requests
- ส่ง notification บ่อยเกินไป
- **Action:** Implement rate limiting

### 500 Internal Server Error
- Push service มีปัญหา
- **Action:** Retry later, log error

---

## 🔐 Security

1. ✅ **ตรวจสอบ Authentication** - ใช้ JWT token
2. ✅ **Validate Input** - ตรวจสอบ subscription format
3. ✅ **Rate Limiting** - จำกัดจำนวน subscriptions ต่อ user (เช่น max 5 devices)
4. ✅ **HTTPS Only** - Push API ต้องใช้ HTTPS
5. ✅ **Private Key Security** - เก็บ VAPID private key ใน env อย่างปลอดภัย

---

## 📝 Implementation Checklist

### Phase 1: Basic Setup
- [ ] ติดตั้ง `webpush-go` library
- [ ] เพิ่ม VAPID keys ใน env
- [ ] สร้าง database table `push_subscriptions`
- [ ] Implement `/api/v1/push/subscribe` endpoint
- [ ] Implement `/api/v1/push/unsubscribe` endpoint
- [ ] ทดสอบ subscribe/unsubscribe ผ่าน Postman

### Phase 2: Send Notifications
- [ ] สร้าง `SendPushNotification` function
- [ ] สร้าง `NotifyUser` function (hybrid WebSocket + Push)
- [ ] Integrate กับ comment system
- [ ] Integrate กับ like system
- [ ] Integrate กับ follow system
- [ ] ทดสอบการส่ง notification

### Phase 3: Optimization
- [ ] Handle expired subscriptions (410/404 responses)
- [ ] Implement rate limiting
- [ ] Add logging และ monitoring
- [ ] Optimize database queries
- [ ] Add admin dashboard สำหรับดู push stats

---

## 🎯 Expected Behavior

### User Journey:

1. **User เปิดเว็บ** → เห็นปุ่ม "เปิดการแจ้งเตือน" ในหน้า Notifications
2. **คลิก "เปิดการแจ้งเตือน"** → Browser ขอ permission
3. **Allow** → Frontend ส่ง subscription ไป backend
4. **Backend บันทึก subscription** → Response success
5. **User ปิดเว็บ**
6. **มี notification ใหม่** (comment/like/follow)
   - ถ้าออนไลน์ → WebSocket ส่งทันที
   - ถ้าปิดเว็บ → Push Notification ส่งไป
7. **User เห็น notification** → คลิก → เปิดเว็บที่ URL ที่กำหนด

---

## 📞 Contact

หากมีคำถามหรือต้องการความช่วยเหลือ:
- Frontend Developer: [Your Name]
- API Endpoint: http://localhost:8080/api/v1/push/*

---

**Last Updated:** 2025-01-06
**Version:** 1.0
