# Media Routes

Base URL: `/api/v1/media`

## Upload Image

**POST** `/api/v1/media/upload/image`

Upload an image file to Bunny CDN.

**Authentication:** Required (JWT Token)

**Content-Type:** `multipart/form-data`

**Form Data:**
- `image`: Image file (required)

**Supported Formats:**
- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- WebP (.webp)

**File Size Limit:** 10 MB

**Example Request (curl):**
```bash
curl -X POST http://localhost:8080/api/v1/media/upload/image \
  -H "Authorization: Bearer <token>" \
  -F "image=@/path/to/photo.jpg"
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "data": {
    "id": "uuid",
    "type": "image",
    "fileName": "photo.jpg",
    "mimeType": "image/jpeg",
    "size": 2048576,
    "url": "https://cdn.bunny.net/storage/images/uuid-photo.jpg",
    "thumbnail": "https://cdn.bunny.net/storage/images/uuid-photo_thumb.jpg",
    "width": 1920,
    "height": 1080,
    "createdAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request:** Missing file, invalid format, or file too large
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "Image file is required"
}
```
- **400 Bad Request:** Invalid image type
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "Invalid image type. Supported: jpeg, png, gif, webp"
}
```
- **400 Bad Request:** File too large
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "Image size must be less than 10MB"
}
```
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Image is automatically compressed and optimized
- Thumbnail is generated automatically
- Image dimensions are extracted if possible
- File is stored on Bunny CDN
- Returns direct CDN URL for immediate use

---

## Upload Video

**POST** `/api/v1/media/upload/video`

Upload a video file to Bunny CDN.

**Authentication:** Required (JWT Token)

**Content-Type:** `multipart/form-data`

**Form Data:**
- `video`: Video file (required)

**Supported Formats:**
- MP4 (.mp4)
- WebM (.webm)
- OGG (.ogg)

**File Size Limit:** 100 MB

**Example Request (curl):**
```bash
curl -X POST http://localhost:8080/api/v1/media/upload/video \
  -H "Authorization: Bearer <token>" \
  -F "video=@/path/to/video.mp4"
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Video uploaded successfully",
  "data": {
    "id": "uuid",
    "type": "video",
    "fileName": "video.mp4",
    "mimeType": "video/mp4",
    "size": 52428800,
    "url": "https://cdn.bunny.net/storage/videos/uuid-video.mp4",
    "thumbnail": "https://cdn.bunny.net/storage/videos/uuid-video_thumb.jpg",
    "width": 1920,
    "height": 1080,
    "duration": 125.5,
    "createdAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request:** Missing file, invalid format, or file too large
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "Video file is required"
}
```
- **400 Bad Request:** Invalid video type
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "Invalid video type. Supported: mp4, webm, ogg"
}
```
- **400 Bad Request:** File too large
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "Video size must be less than 100MB"
}
```
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Video is stored as-is (no transcoding)
- Thumbnail is extracted from first frame
- Video duration is extracted if possible
- Returns direct CDN URL
- Consider using HLS/DASH for large videos

---

## Get Media

**GET** `/api/v1/media/:id`

Get media metadata by ID.

**Authentication:** Public

**Path Parameters:**
- `id`: Media UUID

**Example Request:**
```
GET /api/v1/media/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Media retrieved successfully",
  "data": {
    "id": "uuid",
    "type": "image",
    "fileName": "photo.jpg",
    "mimeType": "image/jpeg",
    "size": 2048576,
    "url": "https://cdn.bunny.net/storage/images/uuid-photo.jpg",
    "thumbnail": "https://cdn.bunny.net/storage/images/uuid-photo_thumb.jpg",
    "width": 1920,
    "height": 1080,
    "createdAt": "2025-11-04T10:00:00Z"
  }
}
```

**Error Responses:**
- **404 Not Found:** Media not found

**Notes:**
- Returns metadata only, not the actual file
- File is served directly from CDN URL

---

## Get User Media

**GET** `/api/v1/media/user/:userId`

Get all media uploaded by a specific user.

**Authentication:** Public

**Path Parameters:**
- `userId`: User UUID

**Query Parameters:**
- `offset`: Page offset (default: 0)
- `limit`: Items per page (default: 20)

**Example Request:**
```
GET /api/v1/media/user/123e4567-e89b-12d3-a456-426614174000?limit=50
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "User media retrieved successfully",
  "data": {
    "media": [
      {
        "id": "uuid",
        "type": "image",
        "fileName": "photo.jpg",
        "mimeType": "image/jpeg",
        "size": 2048576,
        "url": "https://cdn.bunny.net/storage/images/uuid-photo.jpg",
        "thumbnail": "https://cdn.bunny.net/storage/images/uuid-photo_thumb.jpg",
        "width": 1920,
        "height": 1080,
        "createdAt": "2025-11-04T10:00:00Z"
      },
      {
        "id": "uuid",
        "type": "video",
        "fileName": "video.mp4",
        "mimeType": "video/mp4",
        "size": 52428800,
        "url": "https://cdn.bunny.net/storage/videos/uuid-video.mp4",
        "thumbnail": "https://cdn.bunny.net/storage/videos/uuid-video_thumb.jpg",
        "width": 1920,
        "height": 1080,
        "duration": 125.5,
        "createdAt": "2025-11-03T15:30:00Z"
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
- **400 Bad Request:** Invalid user ID
- **500 Internal Server Error:** Failed to retrieve media

**Notes:**
- Ordered by upload date (most recent first)
- Includes both images and videos
- Only returns media uploaded by user

---

## Delete Media

**DELETE** `/api/v1/media/:id`

Delete a media file (owner only).

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id`: Media UUID

**Example Request:**
```
DELETE /api/v1/media/123e4567-e89b-12d3-a456-426614174000
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Media deleted successfully",
  "data": null
}
```

**Error Responses:**
- **400 Bad Request:** Not owner, media not found, or deletion failed
- **401 Unauthorized:** Missing or invalid token

**Notes:**
- Only media owner can delete
- File is removed from Bunny CDN
- Database record is deleted
- Cannot delete media that's currently in use by posts
- Orphaned media (not used in posts) can be cleaned up via background job
