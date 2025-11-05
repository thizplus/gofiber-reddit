# Vote Routes

Base URL: `/api/v1/votes`

## Vote

**POST** `/api/v1/votes`

Upvote or downvote a post or comment.

**Authentication:** Required (JWT Token)

**Request Body:**
```json
{
  "targetId": "uuid",
  "targetType": "post",
  "voteType": "upvote"
}
```

**Validation:**
- `targetId`: required, UUID
- `targetType`: required, enum (`post`, `comment`)
- `voteType`: required, enum (`upvote`, `downvote`)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Vote recorded successfully",
  "data": {
    "targetId": "uuid",
    "targetType": "post",
    "voteType": "upvote",
    "createdAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request:** Validation failed or target not found
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- If user already voted, their vote is updated (upvote ↔ downvote)
- Voting on own content is allowed
- Target's vote count is automatically updated
- Notifications are sent to target author (upvotes only)

---

## Unvote

**DELETE** `/api/v1/votes/:targetType/:targetId`

Remove your vote from a post or comment.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `targetType`: `post` or `comment`
- `targetId`: Target UUID

**Example Request:**
```
DELETE /api/v1/votes/post/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Vote removed successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Invalid target type or target not found
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Target's vote count is automatically decremented
- No error if user hasn't voted on the target

---

## Get Vote

**GET** `/api/v1/votes/:targetType/:targetId`

Get current user's vote on a post or comment.

**Authentication:** Optional (returns null if not authenticated)

**Path Parameters:**
- `targetType`: `post` or `comment`
- `targetId`: Target UUID

**Example Request:**
```
GET /api/v1/votes/post/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Vote status retrieved successfully",
  "data": {
    "voteType": "upvote",
    "createdAt": "2025-11-04T10:00:00Z"
  }
}
```

**Response when not voted:**
```json
{
  "success": true,
  "message": "Vote status retrieved successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Invalid target type
- **404 Not Found:** Target not found

---

## Get Vote Count

**GET** `/api/v1/votes/:targetType/:targetId/count`

Get total vote count for a post or comment.

**Authentication:** Public

**Path Parameters:**
- `targetType`: `post` or `comment`
- `targetId`: Target UUID

**Example Request:**
```
GET /api/v1/votes/post/123e4567-e89b-12d3-a456-426614174000/count
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Vote count retrieved successfully",
  "data": {
    "targetId": "uuid",
    "targetType": "post",
    "voteCount": 42
  }
}
```

**Error Responses:**
- **400 Bad Request:** Invalid target type
- **500 Internal Server Error:** Failed to retrieve vote count

**Notes:**
- Vote count = upvotes - downvotes
- Can be negative if more downvotes than upvotes

---

## Get User Votes

**GET** `/api/v1/votes/user`

Get all votes by current user with pagination.

**Authentication:** Required (JWT Token)

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/votes/user?offset=0&limit=50
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "User votes retrieved successfully",
  "data": {
    "votes": [
      {
        "targetId": "uuid",
        "targetType": "post",
        "voteType": "upvote",
        "createdAt": "2025-11-04T10:00:00Z"
      },
      {
        "targetId": "uuid",
        "targetType": "comment",
        "voteType": "downvote",
        "createdAt": "2025-11-03T15:30:00Z"
      }
    ],
    "meta": {
      "total": 125,
      "offset": 0,
      "limit": 20
    }
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Failed to retrieve votes

**Notes:**
- Votes are ordered by creation date (most recent first)
- Useful for tracking user's voting history
