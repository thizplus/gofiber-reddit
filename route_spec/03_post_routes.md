# Post Routes

Base URL: `/api/v1/posts`

## Create Post

**POST** `/api/v1/posts`

Create a new post.

**Authentication:** Required (JWT Token)

**Request Body:**
```json
{
  "title": "My First Post",
  "content": "This is the content of my post. It can be quite long...",
  "mediaIds": ["uuid1", "uuid2"],
  "tags": ["golang", "backend", "api"],
  "sourcePostId": null
}
```

**Validation:**
- `title`: required, min 1 char, max 300 chars
- `content`: required, min 1 char, max 40000 chars
- `mediaIds`: optional, array of UUIDs (max 10)
- `tags`: optional, array of strings (max 5 tags, each max 50 chars)
- `sourcePostId`: optional, UUID (for crossposts)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Post created successfully",
  "data": {
    "id": "uuid",
    "title": "My First Post",
    "content": "This is the content of my post...",
    "author": {
      "id": "uuid",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": "https://cdn.example.com/avatars/user.jpg"
    },
    "votes": 0,
    "commentCount": 0,
    "media": [
      {
        "id": "uuid1",
        "type": "image",
        "url": "https://cdn.example.com/images/photo.jpg",
        "thumbnail": "https://cdn.example.com/images/photo_thumb.jpg"
      }
    ],
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
    "userVote": null,
    "isSaved": false,
    "hotScore": 0.0
  }
}
```

**Error Responses:**
- **400 Bad Request:** Validation failed
- **401 Unauthorized:** Missing or invalid token

---

## Get Post

**GET** `/api/v1/posts/:id`

Get a single post by ID.

**Authentication:** Optional (Public, but includes user-specific data if authenticated)

**Path Parameters:**
- `id`: Post UUID

**Example Request:**
```
GET /api/v1/posts/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Post retrieved successfully",
  "data": {
    "id": "uuid",
    "title": "My First Post",
    "content": "This is the content of my post...",
    "author": {
      "id": "uuid",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": "https://cdn.example.com/avatars/user.jpg"
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
    "hotScore": 15.5
  }
}
```

**Error Responses:**
- **404 Not Found:** Post not found

---

## Update Post

**PUT** `/api/v1/posts/:id`

Update an existing post (author only).

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: Post UUID

**Request Body:**
```json
{
  "title": "Updated Title",
  "content": "Updated content..."
}
```

**Validation:**
- `title`: required, min 1 char, max 300 chars
- `content`: required, min 1 char, max 40000 chars

**Success Response (200):**
```json
{
  "success": true,
  "message": "Post updated successfully",
  "data": {
    "id": "uuid",
    "title": "Updated Title",
    "content": "Updated content...",
    "author": {
      "id": "uuid",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": "https://cdn.example.com/avatars/user.jpg"
    },
    "votes": 42,
    "commentCount": 15,
    "media": [],
    "tags": [],
    "sourcePost": null,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T11:00:00Z",
    "userVote": "upvote",
    "isSaved": true,
    "hotScore": 15.5
  }
}
```

**Error Responses:**
- **400 Bad Request:** Validation failed or unauthorized (not author)
- **401 Unauthorized:** Missing or invalid token
- **404 Not Found:** Post not found

---

## Delete Post

**DELETE** `/api/v1/posts/:id`

Delete a post (soft delete, author only).

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: Post UUID

**Success Response (200):**
```json
{
  "success": true,
  "message": "Post deleted successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Unauthorized (not author)
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- This is a soft delete - post is marked as deleted but retained in database
- Deleted posts are not visible in listings

---

## List Posts

**GET** `/api/v1/posts`

List posts with pagination and sorting.

**Authentication:** Optional

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20, max: 100)
- `sortBy`: Sort order - `hot`, `new`, `top`, `controversial` (default: hot)

**Example Request:**
```
GET /api/v1/posts?offset=0&limit=20&sortBy=hot
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Posts retrieved successfully",
  "data": {
    "posts": [
      {
        "id": "uuid",
        "title": "Post Title",
        "content": "Post content...",
        "author": {
          "id": "uuid",
          "username": "johndoe",
          "displayName": "John Doe",
          "avatar": "https://cdn.example.com/avatars/user.jpg"
        },
        "votes": 42,
        "commentCount": 15,
        "media": [],
        "tags": [],
        "sourcePost": null,
        "createdAt": "2025-11-04T10:00:00Z",
        "updatedAt": "2025-11-04T10:00:00Z",
        "userVote": null,
        "isSaved": false,
        "hotScore": 15.5
      }
    ],
    "meta": {
      "total": 1000,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Sorting Options:**
- `hot`: Reddit-style hot ranking (votes / (hours + 2)^1.5)
- `new`: Most recent first
- `top`: Highest votes first
- `controversial`: Most controversial (high vote activity, close to 0 score)

---

## List Posts By Author

**GET** `/api/v1/posts/author/:userId`

Get all posts by a specific author.

**Authentication:** Optional

**Path Parameters:**
- `userId`: User UUID

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/posts/author/123e4567-e89b-12d3-a456-426614174000?offset=0&limit=20
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Posts retrieved successfully",
  "data": {
    "posts": [...],
    "meta": {
      "total": 50,
      "offset": 0,
      "limit": 20
    }
  }
}
```

---

## List Posts By Tag

**GET** `/api/v1/posts/tag/:tagName`

Get all posts with a specific tag.

**Authentication:** Optional

**Path Parameters:**
- `tagName`: Tag name (e.g., "golang")

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)
- `sortBy`: Sort order - `hot`, `new`, `top`, `controversial` (default: hot)

**Example Request:**
```
GET /api/v1/posts/tag/golang?offset=0&limit=20&sortBy=hot
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Posts retrieved successfully",
  "data": {
    "posts": [...],
    "meta": {
      "total": 150,
      "offset": 0,
      "limit": 20
    }
  }
}
```

---

## Search Posts

**GET** `/api/v1/posts/search`

Search posts by title or content.

**Authentication:** Optional

**Query Parameters:**
- `q`: Search query (required)
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/posts/search?q=golang&offset=0&limit=20
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Posts retrieved successfully",
  "data": {
    "posts": [...],
    "meta": {
      "total": 75,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Error Responses:**
- **400 Bad Request:** Missing search query

---

## Create Crosspost

**POST** `/api/v1/posts/:id/crosspost`

Create a crosspost (share) of an existing post.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: Source post UUID

**Request Body:**
```json
{
  "title": "Check out this interesting post!",
  "content": "I found this really helpful..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Crosspost created successfully",
  "data": {
    "id": "new-uuid",
    "title": "Check out this interesting post!",
    "content": "I found this really helpful...",
    "author": {
      "id": "uuid",
      "username": "janedoe",
      "displayName": "Jane Doe",
      "avatar": null
    },
    "votes": 0,
    "commentCount": 0,
    "media": [],
    "tags": [],
    "sourcePost": {
      "id": "original-uuid",
      "title": "Original Post Title",
      "content": "Original content...",
      "author": {
        "id": "uuid",
        "username": "johndoe",
        "displayName": "John Doe",
        "avatar": null
      },
      "votes": 42,
      "commentCount": 15
    },
    "createdAt": "2025-11-04T11:00:00Z",
    "updatedAt": "2025-11-04T11:00:00Z",
    "userVote": null,
    "isSaved": false,
    "hotScore": 0.0
  }
}
```

**Error Responses:**
- **400 Bad Request:** Source post not found or validation failed
- **401 Unauthorized:** Missing or invalid token

---

## Get Crossposts

**GET** `/api/v1/posts/:id/crossposts`

Get all crossposts of a post.

**Authentication:** Optional

**Path Parameters:**
- `id`: Source post UUID

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/posts/123e4567-e89b-12d3-a456-426614174000/crossposts
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Crossposts retrieved successfully",
  "data": {
    "posts": [...],
    "meta": {
      "total": 5,
      "offset": 0,
      "limit": 20
    }
  }
}
```

---

## Get Feed

**GET** `/api/v1/posts/feed`

Get personalized feed of posts from followed users.

**Authentication:** Required (JWT Token)

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)
- `sortBy`: Sort order - `hot`, `new` (default: hot)

**Example Request:**
```
GET /api/v1/posts/feed?offset=0&limit=20&sortBy=hot
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Feed retrieved successfully",
  "data": {
    "posts": [...],
    "meta": {
      "total": 200,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Failed to retrieve feed

**Notes:**
- Feed only includes posts from users that the authenticated user follows
- If user follows nobody, feed will be empty
