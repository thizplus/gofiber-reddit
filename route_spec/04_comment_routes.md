# Comment Routes

Base URL: `/api/v1/comments`

## Create Comment

**POST** `/api/v1/comments`

Create a new comment or reply.

**Authentication:** Required (JWT Token)

**Request Body:**
```json
{
  "postId": "uuid",
  "parentId": null,
  "content": "This is my comment on the post."
}
```

**Validation:**
- `postId`: required, UUID
- `parentId`: optional, UUID (for replies to comments)
- `content`: required, min 1 char, max 10000 chars

**Success Response (200):**
```json
{
  "success": true,
  "message": "Comment created successfully",
  "data": {
    "id": "uuid",
    "content": "This is my comment on the post.",
    "postId": "uuid",
    "parentId": null,
    "author": {
      "id": "uuid",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": "https://cdn.example.com/avatars/user.jpg"
    },
    "votes": 0,
    "replyCount": 0,
    "depth": 0,
    "parentChain": [],
    "isDeleted": false,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T10:00:00Z",
    "userVote": null,
    "hotScore": 0.0
  }
}
```

**Error Responses:**
- **400 Bad Request:** Validation failed, post not found, or max depth exceeded (10 levels)
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Maximum nesting depth is 10 levels
- Replying to a deleted comment is not allowed
- Post's comment count is automatically incremented

---

## Get Comment

**GET** `/api/v1/comments/:id`

Get a single comment by ID.

**Authentication:** Optional

**Path Parameters:**
- `id`: Comment UUID

**Success Response (200):**
```json
{
  "success": true,
  "message": "Comment retrieved successfully",
  "data": {
    "id": "uuid",
    "content": "This is my comment.",
    "postId": "uuid",
    "parentId": null,
    "author": {
      "id": "uuid",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": null
    },
    "votes": 15,
    "replyCount": 3,
    "depth": 0,
    "parentChain": [],
    "isDeleted": false,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T10:00:00Z",
    "userVote": "upvote",
    "hotScore": 5.2
  }
}
```

**Error Responses:**
- **404 Not Found:** Comment not found

---

## Update Comment

**PUT** `/api/v1/comments/:id`

Update an existing comment (author only).

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: Comment UUID

**Request Body:**
```json
{
  "content": "Updated comment content..."
}
```

**Validation:**
- `content`: required, min 1 char, max 10000 chars

**Success Response (200):**
```json
{
  "success": true,
  "message": "Comment updated successfully",
  "data": {
    "id": "uuid",
    "content": "Updated comment content...",
    "postId": "uuid",
    "parentId": null,
    "author": {
      "id": "uuid",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": null
    },
    "votes": 15,
    "replyCount": 3,
    "depth": 0,
    "parentChain": [],
    "isDeleted": false,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T11:00:00Z",
    "userVote": "upvote",
    "hotScore": 5.2
  }
}
```

**Error Responses:**
- **400 Bad Request:** Validation failed or unauthorized (not author)
- **401 Unauthorized:** Missing or invalid token
- **404 Not Found:** Comment not found

---

## Delete Comment

**DELETE** `/api/v1/comments/:id`

Delete a comment (soft delete, author only).

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: Comment UUID

**Success Response (200):**
```json
{
  "success": true,
  "message": "Comment deleted successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Unauthorized (not author)
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Soft delete - comment is marked as deleted but retained
- Content is replaced with "[deleted]"
- Replies to deleted comments are still visible
- Post's comment count is automatically decremented

---

## List Comments By Post

**GET** `/api/v1/comments/post/:postId`

Get all top-level comments for a post.

**Authentication:** Optional

**Path Parameters:**
- `postId`: Post UUID

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)
- `sortBy`: Sort order - `hot`, `new`, `top`, `controversial` (default: hot)

**Example Request:**
```
GET /api/v1/comments/post/123e4567-e89b-12d3-a456-426614174000?sortBy=hot&limit=50
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Comments retrieved successfully",
  "data": {
    "comments": [
      {
        "id": "uuid",
        "content": "Great post!",
        "postId": "uuid",
        "parentId": null,
        "author": {
          "id": "uuid",
          "username": "johndoe",
          "displayName": "John Doe",
          "avatar": null
        },
        "votes": 25,
        "replyCount": 5,
        "depth": 0,
        "parentChain": [],
        "isDeleted": false,
        "createdAt": "2025-11-04T10:00:00Z",
        "updatedAt": "2025-11-04T10:00:00Z",
        "userVote": null,
        "hotScore": 8.5
      }
    ],
    "meta": {
      "total": 150,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Notes:**
- Returns only top-level comments (parentId = null)
- Use "Get Comment Tree" or "List Replies" to get nested replies

---

## List Comments By Author

**GET** `/api/v1/comments/author/:userId`

Get all comments by a specific author.

**Authentication:** Optional

**Path Parameters:**
- `userId`: User UUID

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/comments/author/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Comments retrieved successfully",
  "data": {
    "comments": [...],
    "meta": {
      "total": 75,
      "offset": 0,
      "limit": 20
    }
  }
}
```

---

## List Replies

**GET** `/api/v1/comments/:id/replies`

Get direct replies to a comment.

**Authentication:** Optional

**Path Parameters:**
- `id`: Parent comment UUID

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/comments/123e4567-e89b-12d3-a456-426614174000/replies
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Replies retrieved successfully",
  "data": {
    "comments": [
      {
        "id": "uuid",
        "content": "I agree with this!",
        "postId": "uuid",
        "parentId": "123e4567-e89b-12d3-a456-426614174000",
        "author": {
          "id": "uuid",
          "username": "janedoe",
          "displayName": "Jane Doe",
          "avatar": null
        },
        "votes": 10,
        "replyCount": 2,
        "depth": 1,
        "parentChain": ["123e4567-e89b-12d3-a456-426614174000"],
        "isDeleted": false,
        "createdAt": "2025-11-04T10:30:00Z",
        "updatedAt": "2025-11-04T10:30:00Z",
        "userVote": null,
        "hotScore": 3.2
      }
    ],
    "meta": {
      "total": 8,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Notes:**
- Returns only direct children (depth + 1)
- Does not include nested replies (use Get Comment Tree for that)

---

## Get Comment Tree

**GET** `/api/v1/comments/:id/tree`

Get a comment and all its nested replies in tree structure.

**Authentication:** Optional

**Path Parameters:**
- `id`: Comment UUID

**Query Parameters:**
- `maxDepth`: Maximum depth to fetch (default: 10)

**Example Request:**
```
GET /api/v1/comments/123e4567-e89b-12d3-a456-426614174000/tree?maxDepth=3
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Comment tree retrieved successfully",
  "data": {
    "id": "uuid",
    "content": "Parent comment",
    "postId": "uuid",
    "parentId": null,
    "author": {
      "id": "uuid",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": null
    },
    "votes": 25,
    "replyCount": 5,
    "depth": 0,
    "parentChain": [],
    "isDeleted": false,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T10:00:00Z",
    "userVote": null,
    "hotScore": 8.5,
    "replies": [
      {
        "id": "uuid",
        "content": "First reply",
        "postId": "uuid",
        "parentId": "parent-uuid",
        "author": {
          "id": "uuid",
          "username": "janedoe",
          "displayName": "Jane Doe",
          "avatar": null
        },
        "votes": 10,
        "replyCount": 2,
        "depth": 1,
        "parentChain": ["parent-uuid"],
        "isDeleted": false,
        "createdAt": "2025-11-04T10:30:00Z",
        "updatedAt": "2025-11-04T10:30:00Z",
        "userVote": null,
        "hotScore": 3.2,
        "replies": [
          {
            "id": "uuid",
            "content": "Nested reply",
            "replies": []
          }
        ]
      }
    ]
  }
}
```

**Error Responses:**
- **404 Not Found:** Comment not found

**Notes:**
- Returns full nested structure up to maxDepth
- Useful for rendering threaded discussions
- Performance may be impacted for deeply nested trees

---

## Get Parent Chain

**GET** `/api/v1/comments/:id/parent-chain`

Get all parent comments from root to current comment.

**Authentication:** Optional

**Path Parameters:**
- `id`: Comment UUID

**Example Request:**
```
GET /api/v1/comments/123e4567-e89b-12d3-a456-426614174000/parent-chain
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Parent chain retrieved successfully",
  "data": {
    "comments": [
      {
        "id": "root-uuid",
        "content": "Root comment",
        "depth": 0,
        "author": {...},
        "votes": 100
      },
      {
        "id": "middle-uuid",
        "content": "Reply to root",
        "depth": 1,
        "author": {...},
        "votes": 50
      },
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "content": "Current comment",
        "depth": 2,
        "author": {...},
        "votes": 25
      }
    ]
  }
}
```

**Error Responses:**
- **404 Not Found:** Comment not found

**Notes:**
- Returns comments ordered from root to current
- Useful for showing context in deep threads
- Array length indicates nesting depth
