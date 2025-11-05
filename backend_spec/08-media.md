# Media Upload API

## Overview
Endpoints for uploading and managing media files (images and videos).

**Base URL:** `/api/media`

---

## Endpoints

### 1. Upload Media (Private)

Uploads one or more media files.

**Endpoint:** `POST /api/media/upload`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Request Body (FormData):**
```
files: File[] (required, max 10 files per request)
type: 'post' | 'avatar' | 'comment' (optional, default: 'post')
```

**File Constraints:**
- **Images**: JPG, PNG, GIF, WEBP
  - Max size: 10MB per file
  - Max dimensions: 4096x4096px
  - Auto-generate thumbnail (200x200px)

- **Videos**: MP4, WEBM, MOV
  - Max size: 100MB per file
  - Max duration: 10 minutes
  - Max dimensions: 1920x1080px
  - Auto-generate thumbnail from first frame

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "อัปโหลดไฟล์สำเร็จ",
  "data": {
    "media": [
      {
        "id": "m123",
        "type": "image",
        "url": "https://storage.example.com/uploads/2025/01/abc123.jpg",
        "thumbnail": "https://storage.example.com/uploads/2025/01/abc123_thumb.jpg",
        "width": 1920,
        "height": 1080,
        "size": 2458624,
        "filename": "my-image.jpg",
        "uploadedAt": "2025-01-12T10:00:00Z"
      },
      {
        "id": "m124",
        "type": "video",
        "url": "https://storage.example.com/uploads/2025/01/def456.mp4",
        "thumbnail": "https://storage.example.com/uploads/2025/01/def456_thumb.jpg",
        "width": 1920,
        "height": 1080,
        "duration": 45.5,
        "size": 15728640,
        "filename": "my-video.mp4",
        "uploadedAt": "2025-01-12T10:00:30Z"
      }
    ]
  }
}
```

**Error Responses:**

400 Bad Request (Invalid File Type):
```json
{
  "success": false,
  "message": "ประเภทไฟล์ไม่ถูกต้อง",
  "errors": {
    "file_0": "รองรับเฉพาะไฟล์รูปภาพ (JPG, PNG, GIF, WEBP) และวิดีโอ (MP4, WEBM, MOV)"
  }
}
```

400 Bad Request (File Too Large):
```json
{
  "success": false,
  "message": "ไฟล์มีขนาดใหญ่เกินไป",
  "errors": {
    "file_1": "ขนาดไฟล์ต้องไม่เกิน 10MB สำหรับรูปภาพ และ 100MB สำหรับวิดีโอ"
  }
}
```

413 Payload Too Large:
```json
{
  "success": false,
  "message": "ไฟล์มากเกินไป สามารถอัปโหลดได้สูงสุด 10 ไฟล์ต่อครั้ง"
}
```

---

### 2. Get Media Details (Public)

Gets details of a specific media file.

**Endpoint:** `GET /api/media/:id`

**Access:** Public

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "media": {
      "id": "m123",
      "type": "image",
      "url": "https://storage.example.com/uploads/2025/01/abc123.jpg",
      "thumbnail": "https://storage.example.com/uploads/2025/01/abc123_thumb.jpg",
      "width": 1920,
      "height": 1080,
      "size": 2458624,
      "filename": "my-image.jpg",
      "uploadedAt": "2025-01-12T10:00:00Z",
      "uploadedBy": {
        "id": "u1",
        "username": "thepthai",
        "displayName": "เทพไท ใจน้อม"
      }
    }
  }
}
```

**Error Responses:**

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบไฟล์"
}
```

---

### 3. Delete Media (Private)

Deletes a media file (only uploader can delete).

**Endpoint:** `DELETE /api/media/:id`

**Access:** Private (Requires Authentication + Ownership)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ลบไฟล์สำเร็จ"
}
```

**Error Responses:**

403 Forbidden:
```json
{
  "success": false,
  "message": "คุณไม่มีสิทธิ์ลบไฟล์นี้"
}
```

404 Not Found:
```json
{
  "success": false,
  "message": "ไม่พบไฟล์"
}
```

409 Conflict (Media In Use):
```json
{
  "success": false,
  "message": "ไม่สามารถลบไฟล์ที่กำลังใช้งานอยู่ในโพสต์หรือคอมเมนต์"
}
```

---

### 4. Get User's Media (Private)

Gets all media uploaded by the authenticated user.

**Endpoint:** `GET /api/media/me`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
```
page: number (default: 1)
limit: number (default: 20, max: 100)
type: 'image' | 'video' | 'all' (default: 'all')
sortBy: 'new' | 'old' | 'size' (default: 'new')
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "media": [
      {
        "id": "m123",
        "type": "image",
        "url": "https://storage.example.com/uploads/2025/01/abc123.jpg",
        "thumbnail": "https://storage.example.com/uploads/2025/01/abc123_thumb.jpg",
        "width": 1920,
        "height": 1080,
        "size": 2458624,
        "filename": "my-image.jpg",
        "uploadedAt": "2025-01-12T10:00:00Z",
        "usedInPosts": ["1", "7"],
        "usedInComments": []
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 5,
      "totalMedia": 87,
      "hasNextPage": true,
      "hasPrevPage": false
    },
    "storage": {
      "used": 245760000,
      "limit": 5368709120,
      "usedFormatted": "234 MB",
      "limitFormatted": "5 GB",
      "percentUsed": 4.58
    }
  }
}
```

---

### 5. Get Storage Info (Private)

Gets storage usage information for the authenticated user.

**Endpoint:** `GET /api/media/storage`

**Access:** Private (Requires Authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "used": 245760000,
    "limit": 5368709120,
    "usedFormatted": "234 MB",
    "limitFormatted": "5 GB",
    "percentUsed": 4.58,
    "breakdown": {
      "images": {
        "count": 45,
        "size": 123904000,
        "sizeFormatted": "118 MB"
      },
      "videos": {
        "count": 8,
        "size": 121856000,
        "sizeFormatted": "116 MB"
      }
    }
  }
}
```

---

### 6. Optimize Image (Private)

Requests image optimization/compression.

**Endpoint:** `POST /api/media/:id/optimize`

**Access:** Private (Requires Authentication + Ownership)

**Headers:**
```
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "quality": 80 (optional, 1-100, default: 80),
  "maxWidth": 1920 (optional, max width in pixels),
  "maxHeight": 1080 (optional, max height in pixels)
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "ปรับแต่งรูปภาพสำเร็จ",
  "data": {
    "media": {
      "id": "m123",
      "url": "https://storage.example.com/uploads/2025/01/abc123_optimized.jpg",
      "originalSize": 2458624,
      "optimizedSize": 856320,
      "savings": 65.17,
      "savingsFormatted": "1.53 MB saved"
    }
  }
}
```

---

## Media Processing

### Image Processing
1. **Validation**: Check file type, size, dimensions
2. **Security**: Scan for malware, strip EXIF data (except orientation)
3. **Optimization**:
   - Compress to WebP format (fallback to original format)
   - Generate thumbnail (200x200px, cropped)
   - Generate medium size (800px max dimension)
4. **Storage**: Upload to cloud storage (S3/CDN)
5. **Database**: Store metadata in database

### Video Processing
1. **Validation**: Check file type, size, duration, codec
2. **Security**: Scan for malware
3. **Processing**:
   - Extract first frame as thumbnail
   - Transcode to H.264 MP4 (if needed)
   - Generate multiple quality versions (360p, 720p, 1080p)
4. **Storage**: Upload to video CDN
5. **Database**: Store metadata

### Storage Strategy
- **Free tier**: 5GB per user
- **Premium tier**: 50GB per user (future feature)
- Auto-cleanup unused media after 30 days
- CDN caching for better performance

### File Naming
- Original filename is sanitized and stored in database
- Storage filename: `{userId}_{timestamp}_{random}.{ext}`
- Example: `u1_1704876543_a8f3c2.jpg`

### Security
- Validate file type by content (not just extension)
- Scan all uploads for malware
- Strip EXIF data to remove location/metadata
- Rate limit: 10 uploads per minute per user
- Max total size: 100MB per upload request

### URLs
- All media served through CDN
- URLs are permanent (don't change)
- Signed URLs for private content (future feature)
- Auto-generate responsive image srcset

### Cleanup
- Orphaned media (not used in posts/comments) deleted after 30 days
- Deleted post/comment media marked for cleanup
- User can manually delete unused media
- Deleted user's media removed after 7 days
