# Authentication API

## Overview
Authentication endpoints for user registration, login, logout, and session management.

**Base URL:** `/api/auth`

---

## Endpoints

### 1. Register (Public)

Creates a new user account.

**Endpoint:** `POST /api/auth/register`

**Access:** Public

**Request Body:**
```json
{
  "username": "string (required, 3-20 characters, alphanumeric + underscore)",
  "email": "string (required, valid email format)",
  "password": "string (required, min 8 characters)",
  "displayName": "string (required, max 50 characters)"
}
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "ลงทะเบียนสำเร็จ",
  "data": {
    "user": {
      "id": "u1",
      "username": "thepthai",
      "email": "thepthai@example.com",
      "displayName": "เทพไท ใจน้อม",
      "avatar": null,
      "karma": 0,
      "bio": null,
      "coverImage": null,
      "joinedAt": "2025-01-12T10:00:00Z",
      "location": null,
      "website": null,
      "followersCount": 0,
      "followingCount": 0
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**

400 Bad Request:
```json
{
  "success": false,
  "message": "ข้อมูลไม่ถูกต้อง",
  "errors": {
    "username": "ชื่อผู้ใช้ต้องมี 3-20 ตัวอักษร",
    "email": "รูปแบบอีเมลไม่ถูกต้อง",
    "password": "รหัสผ่านต้องมีอย่างน้อย 8 ตัวอักษร"
  }
}
```

409 Conflict:
```json
{
  "success": false,
  "message": "ชื่อผู้ใช้หรืออีเมลนี้ถูกใช้งานแล้ว",
  "field": "username" // or "email"
}
```

---

### 2. Login (Public)

Authenticates a user and returns a JWT token.

**Endpoint:** `POST /api/auth/login`

**Access:** Public

**Request Body:**
```json
{
  "username": "string (required, username or email)",
  "password": "string (required)"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "เข้าสู่ระบบสำเร็จ",
  "data": {
    "user": {
      "id": "u1",
      "username": "thepthai",
      "email": "thepthai@example.com",
      "displayName": "เทพไท ใจน้อม",
      "avatar": "/logo.png",
      "karma": 1247,
      "bio": "นักพัฒนาซอฟต์แวร์ที่หลงรักการเขียนโค้ดและแบ่งปันความรู้ 💻✨",
      "coverImage": null,
      "joinedAt": "2024-06-15T10:00:00Z",
      "location": "กรุงเทพมหานคร, ประเทศไทย",
      "website": "https://thepthai.dev",
      "followersCount": 432,
      "followingCount": 156
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**

400 Bad Request:
```json
{
  "success": false,
  "message": "กรุณากรอกชื่อผู้ใช้และรหัสผ่าน"
}
```

401 Unauthorized:
```json
{
  "success": false,
  "message": "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง"
}
```

---

### 3. Logout (Private)

Logs out the current user and invalidates the token.

**Endpoint:** `POST /api/auth/logout`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Request Body:** None

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ออกจากระบบสำเร็จ"
}
```

**Error Responses:**

401 Unauthorized:
```json
{
  "success": false,
  "message": "ไม่ได้รับอนุญาต กรุณาเข้าสู่ระบบ"
}
```

---

### 4. Get Current User (Private)

Retrieves the currently authenticated user's information.

**Endpoint:** `GET /api/auth/me`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "u1",
      "username": "thepthai",
      "email": "thepthai@example.com",
      "displayName": "เทพไท ใจน้อม",
      "avatar": "/logo.png",
      "karma": 1247,
      "bio": "นักพัฒนาซอฟต์แวร์ที่หลงรักการเขียนโค้ดและแบ่งปันความรู้ 💻✨",
      "coverImage": null,
      "joinedAt": "2024-06-15T10:00:00Z",
      "location": "กรุงเทพมหานคร, ประเทศไทย",
      "website": "https://thepthai.dev",
      "followersCount": 432,
      "followingCount": 156
    }
  }
}
```

**Error Responses:**

401 Unauthorized:
```json
{
  "success": false,
  "message": "ไม่ได้รับอนุญาต กรุณาเข้าสู่ระบบ"
}
```

---

### 5. Refresh Token (Private)

Refreshes the authentication token.

**Endpoint:** `POST /api/auth/refresh`

**Access:** Private (Requires valid token)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 86400
  }
}
```

**Error Responses:**

401 Unauthorized:
```json
{
  "success": false,
  "message": "Token ไม่ถูกต้องหรือหมดอายุ"
}
```

---

## Authentication Flow

1. User registers via `/api/auth/register` → receives token
2. User logs in via `/api/auth/login` → receives token
3. Include token in `Authorization: Bearer {token}` header for all private endpoints
4. Token expires after 24 hours (configurable)
5. Use `/api/auth/refresh` to get a new token
6. User logs out via `/api/auth/logout` → token is invalidated

---

## Security Notes

- Passwords must be hashed using bcrypt (min 10 rounds)
- JWT tokens should include: `userId`, `username`, `iat`, `exp`
- Tokens expire after 24 hours by default
- Implement rate limiting on login/register endpoints (max 5 attempts per 15 minutes)
- Username/email uniqueness must be enforced at database level
- Validate all input fields for XSS prevention
