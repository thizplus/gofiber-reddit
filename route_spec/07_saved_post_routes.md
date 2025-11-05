# Saved Post Routes

Base URL: `/api/v1/saved`

## Save Post

**POST** `/api/v1/saved/posts/:postId`

Bookmark/save a post for later.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `postId`: Post UUID

**Example Request:**
```
POST /api/v1/saved/posts/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Post saved successfully",
  "data": {
    "userId": "current-user-uuid",
    "postId": "123e4567-e89b-12d3-a456-426614174000",
    "createdAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request:** Post not found or already saved
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Cannot save deleted posts
- Error if post is already saved by user

---

## Unsave Post

**DELETE** `/api/v1/saved/posts/:postId`

Remove a post from saved/bookmarked posts.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `postId`: Post UUID

**Example Request:**
```
DELETE /api/v1/saved/posts/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Post unsaved successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Post not saved
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- No error if post wasn't saved

---

## Check Saved Status

**GET** `/api/v1/saved/posts/:postId/status`

Check if current user has saved a post.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `postId`: Post UUID

**Example Request:**
```
GET /api/v1/saved/posts/123e4567-e89b-12d3-a456-426614174000/status
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Saved status retrieved successfully",
  "data": {
    "isSaved": true,
    "savedAt": "2025-11-04T10:00:00Z"
  }
}
```

**Response when not saved:**
```json
{
  "success": true,
  "message": "Saved status retrieved successfully",
  "data": {
    "isSaved": false,
    "savedAt": null
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Failed to check saved status

---

## Get Saved Posts

**GET** `/api/v1/saved/posts`

Get all saved posts for current user.

**Authentication:** Required (JWT Token)

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/saved/posts?offset=0&limit=50
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Saved posts retrieved successfully",
  "data": {
    "posts": [
      {
        "id": "uuid",
        "title": "Saved Post Title",
        "content": "Post content...",
        "author": {
          "id": "uuid",
          "username": "johndoe",
          "displayName": "John Doe",
          "avatar": null
        },
        "votes": 42,
        "commentCount": 15,
        "media": [],
        "tags": [
          {
            "id": "uuid",
            "name": "golang",
            "postCount": 150
          }
        ],
        "sourcePost": null,
        "createdAt": "2025-11-04T10:00:00Z",
        "updatedAt": "2025-11-04T10:00:00Z",
        "userVote": "upvote",
        "isSaved": true,
        "hotScore": 15.5,
        "savedAt": "2025-11-04T12:00:00Z"
      }
    ],
    "meta": {
      "total": 45,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Failed to retrieve saved posts

**Notes:**
- Posts are ordered by save date (most recently saved first)
- Includes `savedAt` field for each post
- Deleted posts are excluded from results
- Includes user-specific data (userVote, isSaved)
