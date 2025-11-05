# User Routes

Base URL: `/api/v1/users`

## Get Profile

**GET** `/api/v1/users/profile`

Get current authenticated user's profile.

**Authentication:** Required (JWT Token)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "displayName": "John Doe",
    "avatar": "https://cdn.example.com/avatars/user.jpg",
    "bio": "Software developer passionate about Go",
    "location": "San Francisco, CA",
    "website": "https://johndoe.com",
    "karma": 150,
    "followersCount": 42,
    "followingCount": 28,
    "role": "user",
    "isActive": true,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **404 Not Found:** User not found

---

## Update Profile

**PUT** `/api/v1/users/profile`

Update current user's profile information.

**Authentication:** Required (JWT Token)

**Request Body:**
```json
{
  "displayName": "John Doe Updated",
  "avatar": "https://cdn.example.com/avatars/new.jpg",
  "bio": "Updated bio text",
  "location": "New York, NY",
  "website": "https://newwebsite.com"
}
```

**Validation:**
- `displayName`: optional, max 100 chars
- `avatar`: optional, valid URL
- `bio`: optional, max 500 chars
- `location`: optional, max 100 chars
- `website`: optional, valid URL

**Success Response (200):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "displayName": "John Doe Updated",
    "avatar": "https://cdn.example.com/avatars/new.jpg",
    "bio": "Updated bio text",
    "location": "New York, NY",
    "website": "https://newwebsite.com",
    "karma": 150,
    "followersCount": 42,
    "followingCount": 28,
    "role": "user",
    "isActive": true,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T10:30:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request:** Validation failed
- **401 Unauthorized:** Missing or invalid token

---

## Delete User

**DELETE** `/api/v1/users/profile`

Delete current user's account (soft delete).

**Authentication:** Required (JWT Token)

**Success Response (200):**
```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** User deletion failed
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- This is a soft delete - user data is retained but marked as inactive
- User will need to contact support to restore account

---

## List Users (Admin Only)

**GET** `/api/v1/users`

List all users with pagination.

**Authentication:** Required (JWT Token) + Admin Role

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 10, max: 100)

**Example Request:**
```
GET /api/v1/users?offset=0&limit=20
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": "uuid",
        "email": "user1@example.com",
        "username": "user1",
        "displayName": "User One",
        "avatar": null,
        "bio": null,
        "location": null,
        "website": null,
        "karma": 50,
        "followersCount": 10,
        "followingCount": 15,
        "role": "user",
        "isActive": true,
        "createdAt": "2025-11-01T10:00:00Z",
        "updatedAt": "2025-11-01T10:00:00Z"
      },
      {
        "id": "uuid",
        "email": "user2@example.com",
        "username": "user2",
        "displayName": "User Two",
        "avatar": null,
        "bio": null,
        "location": null,
        "website": null,
        "karma": 75,
        "followersCount": 20,
        "followingCount": 18,
        "role": "user",
        "isActive": true,
        "createdAt": "2025-11-02T10:00:00Z",
        "updatedAt": "2025-11-02T10:00:00Z"
      }
    ],
    "meta": {
      "total": 100,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **403 Forbidden:** User is not an admin
- **500 Internal Server Error:** Failed to retrieve users
