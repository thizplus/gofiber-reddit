# Notification Routes

Base URL: `/api/v1/notifications`

## Get Notifications

**GET** `/api/v1/notifications`

Get all notifications for current user.

**Authentication:** Required (JWT Token)

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/notifications?offset=0&limit=50
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Notifications retrieved successfully",
  "data": {
    "notifications": [
      {
        "id": "uuid",
        "type": "reply",
        "title": "New reply to your comment",
        "message": "johndoe replied to your comment on 'Post Title'",
        "isRead": false,
        "actorId": "uuid",
        "actor": {
          "id": "uuid",
          "username": "johndoe",
          "displayName": "John Doe",
          "avatar": "https://cdn.example.com/avatars/user.jpg"
        },
        "targetId": "comment-uuid",
        "targetType": "comment",
        "link": "/posts/post-uuid/comments/comment-uuid",
        "createdAt": "2025-11-04T10:00:00Z"
      },
      {
        "id": "uuid",
        "type": "vote",
        "title": "Someone upvoted your post",
        "message": "Your post received an upvote",
        "isRead": true,
        "actorId": null,
        "actor": null,
        "targetId": "post-uuid",
        "targetType": "post",
        "link": "/posts/post-uuid",
        "createdAt": "2025-11-03T15:30:00Z"
      }
    ],
    "meta": {
      "total": 125,
      "offset": 0,
      "limit": 20,
      "unreadCount": 15
    }
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Failed to retrieve notifications

**Notes:**
- Ordered by creation date (most recent first)
- Includes unread count in meta
- Notification types: `reply`, `vote`, `mention`, `follow`

---

## Get Unread Notifications

**GET** `/api/v1/notifications/unread`

Get only unread notifications.

**Authentication:** Required (JWT Token)

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/notifications/unread
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Unread notifications retrieved successfully",
  "data": {
    "notifications": [...],
    "meta": {
      "total": 15,
      "offset": 0,
      "limit": 20,
      "unreadCount": 15
    }
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Failed to retrieve notifications

---

## Get Unread Count

**GET** `/api/v1/notifications/unread/count`

Get count of unread notifications.

**Authentication:** Required (JWT Token)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Unread count retrieved successfully",
  "data": {
    "count": 15
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Failed to retrieve unread count

**Notes:**
- Useful for displaying badge count in UI
- Lightweight endpoint for polling

---

## Get Single Notification

**GET** `/api/v1/notifications/:id`

Get a specific notification.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: Notification UUID

**Success Response (200):**
```json
{
  "success": true,
  "message": "Notification retrieved successfully",
  "data": {
    "id": "uuid",
    "type": "reply",
    "title": "New reply to your comment",
    "message": "johndoe replied to your comment on 'Post Title'",
    "isRead": false,
    "actorId": "uuid",
    "actor": {
      "id": "uuid",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": null
    },
    "targetId": "comment-uuid",
    "targetType": "comment",
    "link": "/posts/post-uuid/comments/comment-uuid",
    "createdAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **404 Not Found:** Notification not found or doesn't belong to user

---

## Mark As Read

**PUT** `/api/v1/notifications/:id/read`

Mark a notification as read.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: Notification UUID

**Success Response (200):**
```json
{
  "success": true,
  "message": "Notification marked as read",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Notification not found or doesn't belong to user
- **401 Unauthorized:** Missing or invalid token

---

## Mark All As Read

**PUT** `/api/v1/notifications/read-all`

Mark all notifications as read.

**Authentication:** Required (JWT Token)

**Success Response (200):**
```json
{
  "success": true,
  "message": "All notifications marked as read",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Failed to mark notifications as read
- **401 Unauthorized:** Missing or invalid token

---

## Delete Notification

**DELETE** `/api/v1/notifications/:id`

Delete a specific notification.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: Notification UUID

**Success Response (200):**
```json
{
  "success": true,
  "message": "Notification deleted successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Notification not found or doesn't belong to user
- **401 Unauthorized:** Missing or invalid token

---

## Delete All Notifications

**DELETE** `/api/v1/notifications`

Delete all notifications for current user.

**Authentication:** Required (JWT Token)

**Success Response (200):**
```json
{
  "success": true,
  "message": "All notifications deleted successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Failed to delete notifications
- **401 Unauthorized:** Missing or invalid token

---

## Get Notification Settings

**GET** `/api/v1/notifications/settings`

Get current user's notification preferences.

**Authentication:** Required (JWT Token)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Notification settings retrieved successfully",
  "data": {
    "userId": "uuid",
    "replyNotifications": true,
    "voteNotifications": false,
    "mentionNotifications": true,
    "followNotifications": true,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Failed to retrieve settings

**Notes:**
- Default settings are created if not exist
- All notification types are enabled by default

---

## Update Notification Settings

**PUT** `/api/v1/notifications/settings`

Update notification preferences.

**Authentication:** Required (JWT Token)

**Request Body:**
```json
{
  "replyNotifications": true,
  "voteNotifications": false,
  "mentionNotifications": true,
  "followNotifications": false
}
```

**Validation:**
- All fields are optional booleans

**Success Response (200):**
```json
{
  "success": true,
  "message": "Notification settings updated successfully",
  "data": {
    "userId": "uuid",
    "replyNotifications": true,
    "voteNotifications": false,
    "mentionNotifications": true,
    "followNotifications": false,
    "createdAt": "2025-11-04T10:00:00Z",
    "updatedAt": "2025-11-04T11:00:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request:** Validation failed
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Only sends specified notification types
- Changes apply to future notifications only
- Does not affect existing notifications
