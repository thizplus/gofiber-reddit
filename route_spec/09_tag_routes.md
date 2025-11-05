# Tag Routes

Base URL: `/api/v1/tags`

All tag routes are **public** (no authentication required).

## Get Tag

**GET** `/api/v1/tags/:id`

Get a tag by ID.

**Authentication:** Public

**Path Parameters:**
- `id`: Tag UUID

**Example Request:**
```
GET /api/v1/tags/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Tag retrieved successfully",
  "data": {
    "id": "uuid",
    "name": "golang",
    "postCount": 150,
    "createdAt": "2025-11-01T10:00:00Z"
  }
}
```

**Error Responses:**
- **404 Not Found:** Tag not found

---

## Get Tag By Name

**GET** `/api/v1/tags/name/:name`

Get a tag by its name.

**Authentication:** Public

**Path Parameters:**
- `name`: Tag name (e.g., "golang")

**Example Request:**
```
GET /api/v1/tags/name/golang
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Tag retrieved successfully",
  "data": {
    "id": "uuid",
    "name": "golang",
    "postCount": 150,
    "createdAt": "2025-11-01T10:00:00Z"
  }
}
```

**Error Responses:**
- **404 Not Found:** Tag not found

**Notes:**
- Tag name is case-insensitive
- Useful for tag-based navigation

---

## List Tags

**GET** `/api/v1/tags`

Get all tags with pagination.

**Authentication:** Public

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 50, max: 100)

**Example Request:**
```
GET /api/v1/tags?offset=0&limit=100
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Tags retrieved successfully",
  "data": {
    "tags": [
      {
        "id": "uuid",
        "name": "golang",
        "postCount": 150,
        "createdAt": "2025-11-01T10:00:00Z"
      },
      {
        "id": "uuid",
        "name": "javascript",
        "postCount": 245,
        "createdAt": "2025-11-01T10:00:00Z"
      },
      {
        "id": "uuid",
        "name": "python",
        "postCount": 320,
        "createdAt": "2025-11-01T10:00:00Z"
      }
    ],
    "meta": {
      "total": 500,
      "offset": 0,
      "limit": 50
    }
  }
}
```

**Error Responses:**
- **500 Internal Server Error:** Failed to retrieve tags

**Notes:**
- Tags are ordered alphabetically by name
- Useful for tag browsing/discovery

---

## Get Popular Tags

**GET** `/api/v1/tags/popular`

Get most popular tags by post count.

**Authentication:** Public

**Query Parameters:**
- `limit`: Number of tags (default: 20, max: 100)

**Example Request:**
```
GET /api/v1/tags/popular?limit=50
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Popular tags retrieved successfully",
  "data": {
    "tags": [
      {
        "id": "uuid",
        "name": "python",
        "postCount": 320,
        "createdAt": "2025-11-01T10:00:00Z"
      },
      {
        "id": "uuid",
        "name": "javascript",
        "postCount": 245,
        "createdAt": "2025-11-01T10:00:00Z"
      },
      {
        "id": "uuid",
        "name": "golang",
        "postCount": 150,
        "createdAt": "2025-11-01T10:00:00Z"
      }
    ]
  }
}
```

**Error Responses:**
- **500 Internal Server Error:** Failed to retrieve popular tags

**Notes:**
- Ordered by post count (descending)
- Useful for trending tags widget
- No pagination meta (fixed limit)

---

## Search Tags

**GET** `/api/v1/tags/search`

Search for tags by name prefix.

**Authentication:** Public

**Query Parameters:**
- `q`: Search query (required)
- `limit`: Number of results (default: 20, max: 50)

**Example Request:**
```
GET /api/v1/tags/search?q=go&limit=10
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Tags search results retrieved successfully",
  "data": {
    "tags": [
      {
        "id": "uuid",
        "name": "golang",
        "postCount": 150,
        "createdAt": "2025-11-01T10:00:00Z"
      },
      {
        "id": "uuid",
        "name": "google",
        "postCount": 45,
        "createdAt": "2025-11-01T10:00:00Z"
      },
      {
        "id": "uuid",
        "name": "goroutine",
        "postCount": 32,
        "createdAt": "2025-11-01T10:00:00Z"
      }
    ]
  }
}
```

**Error Responses:**
- **400 Bad Request:** Missing search query

**Notes:**
- Case-insensitive search
- Matches tags starting with query string
- Ordered by post count (most popular first)
- Useful for tag autocomplete/suggestions
