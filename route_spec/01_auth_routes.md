# Authentication Routes

Base URL: `/api/v1/auth`

## Register

**POST** `/api/v1/auth/register`

Create a new user account.

**Authentication:** Public

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securePassword123",
  "displayName": "John Doe"
}
```

**Validation:**
- `email`: required, valid email format
- `username`: required, min 3 chars, max 50 chars, alphanumeric + underscore
- `password`: required, min 8 chars
- `displayName`: optional, max 100 chars

**Success Response (200):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "displayName": "John Doe",
    "avatar": null,
    "bio": null,
    "location": null,
    "website": null,
    "karma": 0,
    "followersCount": 0,
    "followingCount": 0,
    "role": "user",
    "isActive": true,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request:** Validation failed or user already exists
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "email": "Email already exists",
    "username": "Username must be at least 3 characters"
  }
}
```

---

## Login

**POST** `/api/v1/auth/login`

Authenticate user and receive JWT token.

**Authentication:** Public

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Validation:**
- `email`: required, valid email format
- `password`: required

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": null,
      "bio": null,
      "location": null,
      "website": null,
      "karma": 0,
      "followersCount": 0,
      "followingCount": 0,
      "role": "user",
      "isActive": true,
      "createdAt": "2025-11-04T10:00:00Z",
      "updatedAt": "2025-11-04T10:00:00Z"
    }
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Invalid credentials
```json
{
  "success": false,
  "message": "Login failed",
  "error": "Invalid email or password"
}
```

**Notes:**
- JWT token expires in 7 days
- Include token in `Authorization` header as `Bearer <token>` for protected routes

---

## Google OAuth - Get Authorization URL

**GET** `/api/v1/auth/google`

Get Google OAuth authorization URL to start OAuth flow.

**Authentication:** Public

**Success Response (200):**
```json
{
  "success": true,
  "message": "Google OAuth URL generated",
  "data": {
    "url": "https://accounts.google.com/o/oauth2/v2/auth?client_id=...&redirect_uri=...&response_type=code&scope=email%20profile&state=..."
  }
}
```

**Notes:**
- The returned URL should be used to redirect the user to Google's OAuth consent screen
- A state parameter is stored in a secure cookie for CSRF protection
- After user authorization, Google redirects to the callback URL

---

## Google OAuth - Callback

**GET** `/api/v1/auth/google/callback`

Handle Google OAuth callback and login/register user.

**Authentication:** Public

**Query Parameters:**
- `code` (required): Authorization code from Google
- `state` (required): State parameter for CSRF protection

**Success Response (200):**
```json
{
  "success": true,
  "message": "OAuth authentication successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "email": "user@gmail.com",
      "username": "johndoe123",
      "displayName": "John Doe",
      "avatar": "https://lh3.googleusercontent.com/...",
      "bio": null,
      "location": null,
      "website": null,
      "karma": 0,
      "followersCount": 0,
      "followingCount": 0,
      "role": "user",
      "isActive": true,
      "createdAt": "2025-11-04T10:00:00Z",
      "updatedAt": "2025-11-04T10:00:00Z"
    },
    "isNewUser": true,
    "needsProfile": false
  }
}
```

**Error Responses:**
- **400 Bad Request:** Missing code or invalid state
```json
{
  "success": false,
  "message": "OAuth authentication failed",
  "error": "Authorization code is required"
}
```

- **400 Bad Request:** OAuth authentication failed
```json
{
  "success": false,
  "message": "OAuth authentication failed",
  "error": "failed to exchange code: invalid_grant"
}
```

**OAuth Flow:**
1. Existing user with OAuth: Logs in existing user
2. Existing user by email (no OAuth): Links Google account to existing user
3. New user: Creates new user account with Google profile info

**Notes:**
- JWT token expires in 7 days
- `isNewUser` indicates if a new account was created
- `needsProfile` indicates if user needs to complete profile (always false for Google OAuth)
- Username is automatically generated from email or name if not available
- If email exists but not linked to Google, the accounts are automatically linked
- OAuth users may not have a password (can only login via OAuth)
