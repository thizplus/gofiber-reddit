# Search API

## Overview
Endpoints for searching posts and users.

**Base URL:** `/api/search`

---

## Endpoints

### 1. Search Everything (Public)

Searches across posts and users simultaneously.

**Endpoint:** `GET /api/search`

**Access:** Public

**Query Parameters:**
```
q: string (required, search query, min 2 characters)
type: 'all' | 'posts' | 'users' (default: 'all')
page: number (default: 1)
limit: number (default: 20, max: 100)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "posts": {
      "results": [
        {
          "id": "4",
          "title": "แชร์ประสบการณ์เรียน Full Stack Development แบบ Self-taught",
          "content": "สวัสดีครับทุกคน วันนี้อยากมาแชร์ประสบการณ์...",
          "author": {
            "id": "u4",
            "username": "devjourney",
            "displayName": "นักพัฒนามือใหม่",
            "avatar": null
          },
          "votes": 892,
          "userVote": "up",
          "commentCount": 156,
          "createdAt": "2025-01-09T09:00:00Z",
          "tags": ["programming", "career", "self-taught"],
          "relevance": 0.95
        }
      ],
      "total": 5
    },
    "users": {
      "results": [
        {
          "id": "u4",
          "username": "devjourney",
          "displayName": "นักพัฒนามือใหม่",
          "avatar": null,
          "karma": 1456,
          "bio": "Full-stack Developer | Self-taught | แชร์เส้นทางการเรียนรู้ของผม",
          "followersCount": 678,
          "relevance": 0.88
        }
      ],
      "total": 2
    },
    "query": "programming"
  }
}
```

---

### 2. Search Posts (Public)

Searches only posts.

**Endpoint:** `GET /api/search/posts`

**Access:** Public

**Query Parameters:**
```
q: string (required, search query, min 2 characters)
page: number (default: 1)
limit: number (default: 20, max: 100)
sortBy: 'relevance' | 'new' | 'top' | 'comments' (default: 'relevance')
tag: string (optional, filter by tag)
author: string (optional, filter by username)
timeRange: 'today' | 'week' | 'month' | 'year' | 'all' (default: 'all')
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "posts": [
      {
        "id": "2",
        "title": "วิธีทำ Pad Thai ที่บ้านง่ายๆ แบบไทยแท้",
        "content": "วันนี้จะมาแชร์วิธีทำผัดไทยที่บ้านกันครับ...",
        "author": {
          "id": "u2",
          "username": "cookingmaster",
          "displayName": "เชฟมือทอง",
          "avatar": null
        },
        "votes": 567,
        "userVote": "up",
        "commentCount": 89,
        "createdAt": "2025-01-10T10:15:00Z",
        "media": [
          {
            "id": "m1",
            "type": "image",
            "url": "https://images.unsplash.com/photo-1559314809-0d155014e29e?w=800",
            "thumbnail": "https://images.unsplash.com/photo-1559314809-0d155014e29e?w=200"
          }
        ],
        "tags": ["อาหาร", "สูตรอาหาร", "ผัดไทย"],
        "highlight": {
          "title": "วิธีทำ <mark>Pad Thai</mark> ที่บ้านง่ายๆ แบบไทยแท้",
          "content": "วันนี้จะมาแชร์วิธีทำ<mark>ผัดไทย</mark>ที่บ้านกันครับ..."
        }
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 2,
      "totalResults": 25,
      "hasNextPage": true,
      "hasPrevPage": false
    },
    "query": "pad thai"
  }
}
```

---

### 3. Search Users (Public)

Searches only users.

**Endpoint:** `GET /api/search/users`

**Access:** Public

**Query Parameters:**
```
q: string (required, search query, min 2 characters)
page: number (default: 1)
limit: number (default: 20, max: 100)
sortBy: 'relevance' | 'karma' | 'followers' (default: 'relevance')
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "u2",
        "username": "cookingmaster",
        "displayName": "เชฟมือทอง",
        "avatar": null,
        "karma": 2891,
        "bio": "รักการทำอาหาร มาแชร์สูตรอาหารไทยและนานาชาติกัน 🍳👨‍🍳",
        "followersCount": 1205,
        "followingCount": 89,
        "isFollowing": false,
        "highlight": {
          "username": "<mark>cookingmaster</mark>",
          "displayName": "เชฟมือทอง",
          "bio": "รักการทำ<mark>อาหาร</mark> มาแชร์สูตร<mark>อาหาร</mark>ไทยและนานาชาติกัน"
        }
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 1,
      "totalResults": 3,
      "hasNextPage": false,
      "hasPrevPage": false
    },
    "query": "cooking"
  }
}
```

---

### 4. Get Search Suggestions (Public)

Gets search suggestions/autocomplete based on partial query.

**Endpoint:** `GET /api/search/suggestions`

**Access:** Public

**Query Parameters:**
```
q: string (required, partial search query, min 2 characters)
limit: number (default: 10, max: 20)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "suggestions": [
      {
        "type": "user",
        "text": "cookingmaster",
        "displayText": "@cookingmaster - เชฟมือทอง",
        "avatar": null
      },
      {
        "type": "tag",
        "text": "อาหาร",
        "displayText": "#อาหาร",
        "count": 15
      },
      {
        "type": "post",
        "text": "วิธีทำ Pad Thai ที่บ้านง่ายๆ",
        "displayText": "วิธีทำ Pad Thai ที่บ้านง่ายๆ แบบไทยแท้",
        "postId": "2"
      }
    ]
  }
}
```

---

### 5. Get Popular Tags (Public)

Gets the most popular tags in the system.

**Endpoint:** `GET /api/search/tags`

**Access:** Public

**Query Parameters:**
```
limit: number (default: 20, max: 100)
timeRange: 'today' | 'week' | 'month' | 'all' (default: 'all')
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "tags": [
      {
        "name": "programming",
        "count": 245,
        "trendingScore": 8.5
      },
      {
        "name": "อาหาร",
        "count": 189,
        "trendingScore": 7.2
      },
      {
        "name": "ท่องเที่ยว",
        "count": 156,
        "trendingScore": 6.8
      }
    ]
  }
}
```

---

### 6. Get Trending Searches (Public)

Gets the most popular search queries.

**Endpoint:** `GET /api/search/trending`

**Access:** Public

**Query Parameters:**
```
limit: number (default: 10, max: 20)
timeRange: 'today' | 'week' | 'month' (default: 'today')
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "trending": [
      {
        "query": "macbook m3",
        "count": 1250,
        "change": "+15%"
      },
      {
        "query": "เชียงใหม่",
        "count": 980,
        "change": "+8%"
      },
      {
        "query": "ผัดไทย",
        "count": 745,
        "change": "new"
      }
    ]
  }
}
```

---

### 7. Search History (Private)

Gets the authenticated user's search history.

**Endpoint:** `GET /api/search/history`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
```
limit: number (default: 20, max: 50)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "history": [
      {
        "id": "h1",
        "query": "macbook m3",
        "type": "posts",
        "searchedAt": "2025-01-12T09:30:00Z"
      },
      {
        "id": "h2",
        "query": "cookingmaster",
        "type": "users",
        "searchedAt": "2025-01-11T15:20:00Z"
      }
    ]
  }
}
```

---

### 8. Delete Search History (Private)

Deletes the authenticated user's search history.

**Endpoint:** `DELETE /api/search/history`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ลบประวัติการค้นหาสำเร็จ",
  "data": {
    "deletedCount": 45
  }
}
```

---

## Search Features

### Full-Text Search
Posts are searched in:
1. Title (highest weight)
2. Content (medium weight)
3. Tags (medium weight)
4. Author username/displayName (low weight)

Users are searched in:
1. Username (highest weight)
2. DisplayName (high weight)
3. Bio (medium weight)

### Search Ranking

**Relevance Score Formula:**
```
score = (text_match_score * 0.6) + (popularity_score * 0.3) + (freshness_score * 0.1)

Where:
- text_match_score: How well the query matches the content (0-1)
- popularity_score: Normalized vote/follower count (0-1)
- freshness_score: Recency factor (0-1)
```

### Highlighting
- Return matched terms wrapped in `<mark>` tags
- Highlight in: title, content (first 200 chars), username, displayName, bio

### Search Filters
- **Posts**: tag, author, timeRange, sortBy
- **Users**: sortBy (relevance, karma, followers)
- Combined searches show top 5 of each type

### Special Queries
- `@username`: Search for specific user
- `#tag`: Search for specific tag
- `"exact phrase"`: Search for exact phrase match

### Performance
- Index: title, content, tags, username, displayName, bio
- Cache popular searches (5 minutes)
- Limit: 2 searches per second per user
- Min query length: 2 characters
- Max query length: 100 characters

### Privacy
- Search history is private
- No personalized results (same results for everyone)
- Deleted content excluded from results
