# Search Routes

Base URL: `/api/v1/search`

## Search

**GET** `/api/v1/search`

Search across posts, users, and tags.

**Authentication:** Optional (includes user-specific data if authenticated)

**Query Parameters:**
- `q`: Search query (required)
- `type`: Search type - `all`, `post`, `user`, `tag` (default: all)
- `limit`: Results per type (default: 20, max: 50)

**Example Request:**
```
GET /api/v1/search?q=golang&type=all&limit=10
```

**Success Response (200) - Type: "all":**
```json
{
  "success": true,
  "message": "Search completed successfully",
  "data": {
    "query": "golang",
    "posts": [
      {
        "id": "uuid",
        "title": "Getting Started with Golang",
        "content": "This is a tutorial...",
        "author": {
          "id": "uuid",
          "username": "johndoe",
          "displayName": "John Doe",
          "avatar": null
        },
        "votes": 42,
        "commentCount": 15,
        "createdAt": "2025-11-04T10:00:00Z",
        "userVote": null,
        "isSaved": false
      }
    ],
    "users": [
      {
        "id": "uuid",
        "username": "golangdev",
        "displayName": "Golang Developer",
        "avatar": null,
        "bio": "Passionate about Go",
        "karma": 150,
        "followersCount": 42,
        "isFollowing": false
      }
    ],
    "tags": [
      {
        "id": "uuid",
        "name": "golang",
        "postCount": 150
      }
    ]
  }
}
```

**Success Response (200) - Type: "post":**
```json
{
  "success": true,
  "message": "Search completed successfully",
  "data": {
    "query": "golang",
    "posts": [...]
  }
}
```

**Success Response (200) - Type: "user":**
```json
{
  "success": true,
  "message": "Search completed successfully",
  "data": {
    "query": "golang",
    "users": [...]
  }
}
```

**Success Response (200) - Type: "tag":**
```json
{
  "success": true,
  "message": "Search completed successfully",
  "data": {
    "query": "golang",
    "tags": [...]
  }
}
```

**Error Responses:**
- **400 Bad Request:** Missing search query or invalid type
- **500 Internal Server Error:** Search failed

**Notes:**
- If authenticated, search history is automatically saved
- Posts searched by title and content
- Users searched by username and displayName
- Tags searched by name
- Results are ranked by relevance and popularity

---

## Get Popular Searches

**GET** `/api/v1/search/popular`

Get most popular search queries.

**Authentication:** Public

**Query Parameters:**
- `limit`: Number of results (default: 10, max: 50)

**Example Request:**
```
GET /api/v1/search/popular?limit=20
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Popular searches retrieved successfully",
  "data": {
    "searches": [
      {
        "query": "golang",
        "count": 1250
      },
      {
        "query": "react",
        "count": 980
      },
      {
        "query": "python",
        "count": 875
      }
    ]
  }
}
```

**Error Responses:**
- **500 Internal Server Error:** Failed to retrieve popular searches

**Notes:**
- Ordered by search count (descending)
- Useful for trending searches widget
- Shows aggregate search counts across all users

---

## Get Search History

**GET** `/api/v1/search/history`

Get current user's search history.

**Authentication:** Required (JWT Token)

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/search/history?limit=50
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Search history retrieved successfully",
  "data": {
    "history": [
      {
        "id": "uuid",
        "query": "golang tutorial",
        "searchedAt": "2025-11-04T10:00:00Z"
      },
      {
        "id": "uuid",
        "query": "react hooks",
        "searchedAt": "2025-11-03T15:30:00Z"
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
- **500 Internal Server Error:** Failed to retrieve search history

**Notes:**
- Ordered by search date (most recent first)
- Duplicate searches create new history entries
- Useful for recent searches feature

---

## Clear Search History

**DELETE** `/api/v1/search/history`

Clear all search history for current user.

**Authentication:** Required (JWT Token)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Search history cleared successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Failed to clear search history
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Permanently deletes all search history
- Cannot be undone

---

## Delete Search History Item

**DELETE** `/api/v1/search/history/:id`

Delete a specific search history item.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: History item UUID

**Example Request:**
```
DELETE /api/v1/search/history/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Search history item deleted successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** History item not found or doesn't belong to user
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Only deletes single history entry
- Useful for removing specific unwanted searches
