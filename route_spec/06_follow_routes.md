# Follow Routes

Base URL: `/api/v1/follows`

## Follow User

**POST** `/api/v1/follows/user/:userId`

Follow another user.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `userId`: User UUID to follow

**Example Request:**
```
POST /api/v1/follows/user/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Followed successfully",
  "data": {
    "followerId": "current-user-uuid",
    "followingId": "123e4567-e89b-12d3-a456-426614174000",
    "createdAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request:** Already following, cannot follow yourself, or user not found
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Cannot follow yourself
- Sends notification to followed user
- Both users' follower/following counts are updated

---

## Unfollow User

**DELETE** `/api/v1/follows/user/:userId`

Unfollow a user.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `userId`: User UUID to unfollow

**Example Request:**
```
DELETE /api/v1/follows/user/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Unfollowed successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Not following this user
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Both users' follower/following counts are updated
- No error if you weren't following the user

---

## Check Follow Status

**GET** `/api/v1/follows/user/:userId/status`

Check if current user is following another user.

**Authentication:** Optional (returns false if not authenticated)

**Path Parameters:**
- `userId`: User UUID to check

**Example Request:**
```
GET /api/v1/follows/user/123e4567-e89b-12d3-a456-426614174000/status
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Follow status retrieved successfully",
  "data": {
    "isFollowing": true,
    "followedAt": "2025-11-04T10:00:00Z"
  }
}
```

**Response when not following:**
```json
{
  "success": true,
  "message": "Follow status retrieved successfully",
  "data": {
    "isFollowing": false,
    "followedAt": null
  }
}
```

**Error Responses:**
- **500 Internal Server Error:** Failed to check follow status

---

## Get Followers

**GET** `/api/v1/follows/user/:userId/followers`

Get list of users following a specific user.

**Authentication:** Optional

**Path Parameters:**
- `userId`: User UUID

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/follows/user/123e4567-e89b-12d3-a456-426614174000/followers?limit=50
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Followers retrieved successfully",
  "data": {
    "users": [
      {
        "id": "uuid",
        "username": "follower1",
        "displayName": "Follower One",
        "avatar": "https://cdn.example.com/avatars/user.jpg",
        "bio": "User bio",
        "karma": 150,
        "followersCount": 42,
        "followingCount": 28,
        "isFollowing": true,
        "followedAt": "2025-11-03T10:00:00Z"
      }
    ],
    "meta": {
      "total": 42,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Error Responses:**
- **400 Bad Request:** Invalid user ID
- **500 Internal Server Error:** Failed to retrieve followers

**Notes:**
- If authenticated, includes `isFollowing` to show if you follow each follower back
- Ordered by follow date (most recent first)

---

## Get Following

**GET** `/api/v1/follows/user/:userId/following`

Get list of users that a specific user is following.

**Authentication:** Optional

**Path Parameters:**
- `userId`: User UUID

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/follows/user/123e4567-e89b-12d3-a456-426614174000/following
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Following retrieved successfully",
  "data": {
    "users": [
      {
        "id": "uuid",
        "username": "following1",
        "displayName": "Following One",
        "avatar": null,
        "bio": null,
        "karma": 200,
        "followersCount": 100,
        "followingCount": 50,
        "isFollowing": false,
        "followedAt": "2025-11-02T10:00:00Z"
      }
    ],
    "meta": {
      "total": 28,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Error Responses:**
- **400 Bad Request:** Invalid user ID
- **500 Internal Server Error:** Failed to retrieve following

**Notes:**
- If authenticated, includes `isFollowing` to show if you also follow each user
- Ordered by follow date (most recent first)

---

## Get Mutual Follows

**GET** `/api/v1/follows/mutuals`

Get list of mutual follows (friends - users you follow who follow you back).

**Authentication:** Required (JWT Token)

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/follows/mutuals?limit=50
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Mutual follows retrieved successfully",
  "data": {
    "users": [
      {
        "id": "uuid",
        "username": "friend1",
        "displayName": "Friend One",
        "avatar": "https://cdn.example.com/avatars/user.jpg",
        "bio": "We follow each other!",
        "karma": 250,
        "followersCount": 150,
        "followingCount": 120,
        "isFollowing": true,
        "followedAt": "2025-11-01T10:00:00Z"
      }
    ],
    "meta": {
      "total": 15,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Failed to retrieve mutual follows

**Notes:**
- Only returns users where both users follow each other
- Useful for finding close connections/friends
- Ordered by follow date (most recent first)
