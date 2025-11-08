# Upload Endpoints Specification
## File Upload API for Posts & Chat Features

**Created**: 2025-01-07
**Purpose**: Add missing file upload endpoints for better UX
**Target**: Posts feature & Chat feature
**Priority**: 🟡 MEDIUM (Recommended before Chat launch)
**Effort**: 4-6 hours

---

## 🎯 Executive Summary

### Current Status:
```
✅ POST /media/upload/image - Exists (for Posts)
✅ POST /media/upload/video - Exists (for Posts)
❌ POST /upload/file        - Missing! (Need for Chat & Posts)
```

### Problem:
Currently, **Chat feature** requires uploading files (PDF, DOC, XLS, ZIP) as part of the message send request. This creates poor UX:
- ❌ No upload progress indicator
- ❌ Cannot preview before sending
- ❌ Cannot retry upload separately if fails
- ❌ Duplicate implementation (Chat vs Posts)

### Solution:
Add **dedicated upload endpoints** that can be shared across features:
- ✅ Separate upload from message sending
- ✅ Show upload progress to user
- ✅ Allow retry without resending message
- ✅ Consistent API across Posts & Chat

---

## 🏗️ Storage Architecture Overview

### ⚠️ IMPORTANT: Two Different Upload Systems

**1. Bunny Storage** (Static Files)
- **Used for**: Images, Files (PDF, DOC, ZIP, etc.)
- **Purpose**: Static file storage & CDN delivery
- **Endpoints**:
  - `POST /media/upload/image` → Bunny Storage
  - `POST /upload/file` → Bunny Storage
- **Implementation**: Direct file upload to storage zone
- **Result**: Direct URL to file

**2. Bunny Stream** (Video Platform)
- **Used for**: Videos only
- **Purpose**: Video encoding, transcoding, HLS streaming
- **Endpoint**:
  - `POST /media/upload/video` → Bunny Stream API (NOT Storage!)
- **Implementation**: Upload to Bunny Stream Library → Auto transcode → HLS playlist
- **Result**: Video ID + Playlist URL (not direct MP4!)
- **Documentation**: See `vdo_stream_implement.md` for complete implementation

### Why Different Systems?

| Feature | Bunny Storage | Bunny Stream |
|---------|--------------|--------------|
| File Type | Images, Documents | Videos |
| Delivery | Direct CDN URL | HLS Adaptive Streaming |
| Processing | None | Auto transcode (360p, 720p, 1080p) |
| Performance | Standard | Optimized (95% faster initial load) |
| Cost | Pay per storage/bandwidth | Pay per encoding + streaming |

---

## 📋 Current Implementation (Posts)

### 1. POST /media/upload/image ✅

**Endpoint**: `POST /api/v1/media/upload/image`
**Status**: ✅ Already implemented
**Storage**: **Bunny Storage** (CDN)

**Request:**
```bash
POST /api/v1/media/upload/image
Content-Type: multipart/form-data
Authorization: Bearer <token>

FormData:
  image: File
```

**Response:**
```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "data": {
    "id": "media-001",
    "url": "https://cdn.bunny.net/images/xxx.jpg",
    "thumbnail": "https://cdn.bunny.net/images/thumb_xxx.jpg",
    "width": 1920,
    "height": 1080,
    "size": 1024000,
    "mimeType": "image/jpeg",
    "createdAt": "2024-01-01T10:00:00Z"
  }
}
```

**Validation:**
- File types: `image/jpeg`, `image/png`, `image/gif`, `image/webp`
- Max size: 10 MB
- Required: Yes

**Backend Implementation:**
- Upload to Bunny Storage zone
- Generate thumbnail automatically
- Return direct CDN URL

---

### 2. POST /media/upload/video ✅

**Endpoint**: `POST /api/v1/media/upload/video`
**Status**: ✅ Already implemented
**Storage**: **Bunny Stream** (Video Platform) ⚠️ **NOT Bunny Storage!**

**Request:**
```bash
POST /api/v1/media/upload/video
Content-Type: multipart/form-data
Authorization: Bearer <token>

FormData:
  video: File
```

**Response:**
```json
{
  "success": true,
  "message": "Video uploaded successfully",
  "data": {
    "id": "media-002",
    "videoId": "b1631ae0-4c8a-47b0-430d-c9ad4914",
    "playlistUrl": "https://vz-b1631ae0-4c8.b-cdn.net/b1631ae0-4c8a-47b0-430d-c9ad4914/playlist.m3u8",
    "thumbnail": "https://vz-b1631ae0-4c8.b-cdn.net/b1631ae0-4c8a-47b0-430d-c9ad4914/thumbnail.jpg",
    "duration": 120.5,
    "width": 1920,
    "height": 1080,
    "size": 50000000,
    "mimeType": "video/mp4",
    "availableResolutions": ["360p", "720p", "1080p"],
    "createdAt": "2024-01-01T10:00:00Z"
  }
}
```

**Validation:**
- File types: `video/mp4`, `video/quicktime`, `video/webm`
- Max size: 300 MB (or 100MB for chat as per chat spec)
- Required: Yes

**Backend Implementation:**
- ⚠️ **MUST use Bunny Stream API** (NOT Storage!)
- Upload to Bunny Stream Library ID: `533535`
- CDN: `vz-b1631ae0-4c8.b-cdn.net`
- Auto-transcode to multiple resolutions (360p, 720p, 1080p)
- Generate HLS playlist (.m3u8) for adaptive streaming
- Auto-generate thumbnail
- Return playlist URL (not direct MP4!)

**📖 Complete Implementation Guide:**
See `vdo_stream_implement.md` for:
- Bunny Stream API integration
- Go code implementation
- HLS player frontend code
- Performance optimization (95% faster)

---

## 🆕 New Requirements

### 3. POST /upload/file ❌ NEW!

**Endpoint**: `POST /api/v1/upload/file` (or `/media/upload/file`)
**Status**: ❌ **NEEDS IMPLEMENTATION**
**Storage**: **Bunny Storage** (CDN)
**Priority**: 🟡 MEDIUM

**Purpose:**
- Upload document files for Chat attachments (PDF, DOC, XLS, ZIP, etc.)
- Upload files for Posts (future feature)
- Consistent with existing image upload pattern
- ⚠️ **NOT for videos** - Videos use Bunny Stream instead

---

## 📝 API Specification: POST /upload/file

### Endpoint Details

**URL**: `POST /api/v1/upload/file`
**Authentication**: Required (JWT Bearer token)
**Content-Type**: `multipart/form-data`
**Rate Limit**: 20 uploads/minute (recommended)

---

### Request Format

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
```

**Body (FormData):**
```
file: File (required)
```

**Example (cURL):**
```bash
curl -X POST https://api.voobize.com/v1/upload/file \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "file=@document.pdf"
```

**Example (JavaScript):**
```typescript
const formData = new FormData();
formData.append('file', fileObject);

const response = await fetch('https://api.voobize.com/v1/upload/file', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: formData
});

const data = await response.json();
```

---

### Response Format

#### Success Response (201 Created)

```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "id": "file-001",
    "url": "https://cdn.bunny.net/files/abc123-document.pdf",
    "filename": "document.pdf",
    "size": 5000000,
    "mimeType": "application/pdf",
    "extension": "pdf",
    "createdAt": "2024-01-01T10:00:00Z"
  }
}
```

**Response Fields:**
| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique file ID (for reference) |
| `url` | string | **Direct CDN URL from Bunny Storage** |
| `filename` | string | Original filename (sanitized) |
| `size` | number | File size in bytes |
| `mimeType` | string | MIME type |
| `extension` | string | File extension (pdf, docx, etc.) |
| `createdAt` | string | ISO 8601 timestamp |

**Note:** This endpoint uses **Bunny Storage** (not Stream) for direct file delivery.

---

### Validation Rules

#### Allowed File Types:
```typescript
const ALLOWED_MIME_TYPES = [
  'application/pdf',                    // PDF
  'application/msword',                 // DOC
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document', // DOCX
  'application/vnd.ms-excel',           // XLS
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',       // XLSX
  'application/vnd.ms-powerpoint',      // PPT
  'application/vnd.openxmlformats-officedocument.presentationml.presentation', // PPTX
  'application/zip',                    // ZIP
  'application/x-rar-compressed',       // RAR
  'application/x-7z-compressed',        // 7Z
  'text/plain',                         // TXT
  'text/csv',                           // CSV
];
```

#### File Size Limits:
```
Max size: 50 MB (52,428,800 bytes)
Min size: 1 byte (prevent empty files)
```

#### Filename Validation:
- Preserve original filename (sanitized)
- Remove special characters: `/:*?"<>|`
- Max filename length: 255 characters
- Generate unique filename for storage: `{uuid}-{sanitized-filename}`

---

### Error Responses

#### 1. File Too Large (413)
```json
{
  "success": false,
  "message": "File size exceeds maximum allowed size",
  "error": "FILE_TOO_LARGE",
  "details": {
    "maxSize": 52428800,
    "receivedSize": 60000000,
    "maxSizeMB": 50
  }
}
```

#### 2. Invalid File Type (400)
```json
{
  "success": false,
  "message": "File type not allowed",
  "error": "INVALID_FILE_TYPE",
  "details": {
    "receivedType": "application/exe",
    "allowedTypes": [
      "application/pdf",
      "application/msword",
      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      "application/zip",
      "text/plain"
    ]
  }
}
```

#### 3. No File Provided (400)
```json
{
  "success": false,
  "message": "No file provided",
  "error": "FILE_REQUIRED"
}
```

#### 4. Upload Failed (500)
```json
{
  "success": false,
  "message": "Failed to upload file to storage",
  "error": "UPLOAD_FAILED",
  "details": {
    "reason": "CDN connection timeout"
  }
}
```

#### 5. Unauthorized (401)
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "UNAUTHORIZED"
}
```

#### 6. Rate Limit Exceeded (429)
```json
{
  "success": false,
  "message": "Too many upload requests. Please try again later.",
  "error": "RATE_LIMIT_EXCEEDED",
  "details": {
    "retryAfter": 60
  }
}
```

---

## 🔧 Implementation Guide (Backend)

### ⚠️ Important: Storage vs Stream

This implementation guide covers **file uploads to Bunny Storage only**.

**For video uploads:**
- Videos use **Bunny Stream API** (different implementation)
- See `vdo_stream_implement.md` for complete video upload guide
- Do NOT mix Storage and Stream implementations!

### File Structure

**Recommended location:**
```
backend/
├── interfaces/
│   └── api/
│       └── handlers/
│           └── media_handler.go (add UploadFile method)
├── application/
│   └── services/
│       └── media_service.go (add UploadFile logic)
├── infrastructure/
│   └── storage/
│       ├── bunny_storage_client.go (existing - for images & files)
│       ├── bunny_stream_client.go (existing - for videos only!)
│       └── file_upload_service.go (NEW - uses BunnyStorageClient)
└── domain/
    └── models/
        └── media.go (update Media model if needed)
```

---

### Step 1: Create File Upload Service

**File**: `infrastructure/storage/file_upload_service.go`

**⚠️ IMPORTANT:** This service uses `BunnyStorageClient` (NOT BunnyStreamClient!)

```go
package storage

import (
    "fmt"
    "io"
    "mime/multipart"
    "path/filepath"
    "strings"

    "github.com/google/uuid"
)

// FileUploadService handles file uploads to Bunny Storage (not Stream!)
// For video uploads, use VideoUploadService with Bunny Stream instead
type FileUploadService struct {
    bunnyStorage *BunnyStorageClient // ⚠️ Use Storage client (not Stream!)
}

// Allowed MIME types for file upload
var AllowedFileMIMETypes = []string{
    "application/pdf",
    "application/msword",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "application/vnd.ms-excel",
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    "application/vnd.ms-powerpoint",
    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    "application/zip",
    "application/x-rar-compressed",
    "application/x-7z-compressed",
    "text/plain",
    "text/csv",
}

const (
    MaxFileSize = 50 * 1024 * 1024 // 50 MB
    FilesFolder = "files"           // Folder in Bunny Storage
)

// UploadFile uploads a file to Bunny Storage (CDN)
// ⚠️ This is for files only (PDF, DOC, ZIP, etc.)
// For videos, use Bunny Stream API instead (see vdo_stream_implement.md)
func (s *FileUploadService) UploadFile(file *multipart.FileHeader) (*FileUploadResult, error) {
    // Validate file size
    if file.Size > MaxFileSize {
        return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
    }

    if file.Size == 0 {
        return nil, fmt.Errorf("empty file not allowed")
    }

    // Validate MIME type (should NOT include video types!)
    if !isAllowedFileMIMEType(file.Header.Get("Content-Type")) {
        return nil, fmt.Errorf("file type not allowed: %s", file.Header.Get("Content-Type"))
    }

    // Open file
    src, err := file.Open()
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer src.Close()

    // Generate unique filename
    filename := sanitizeFilename(file.Filename)
    uniqueFilename := fmt.Sprintf("%s-%s", uuid.New().String(), filename)
    storagePath := fmt.Sprintf("%s/%s", FilesFolder, uniqueFilename)

    // Upload to Bunny Storage (returns direct CDN URL)
    url, err := s.bunnyStorage.Upload(storagePath, src, file.Size)
    if err != nil {
        return nil, fmt.Errorf("failed to upload to storage: %w", err)
    }

    // Get file extension
    extension := strings.TrimPrefix(filepath.Ext(filename), ".")

    return &FileUploadResult{
        URL:       url,
        Filename:  filename,
        Size:      file.Size,
        MimeType:  file.Header.Get("Content-Type"),
        Extension: extension,
    }, nil
}

// isAllowedFileMIMEType checks if MIME type is allowed
func isAllowedFileMIMEType(mimeType string) bool {
    for _, allowed := range AllowedFileMIMETypes {
        if mimeType == allowed {
            return true
        }
    }
    return false
}

// sanitizeFilename removes special characters from filename
func sanitizeFilename(filename string) string {
    // Remove path separators and special characters
    filename = filepath.Base(filename)

    // Replace special characters
    replacer := strings.NewReplacer(
        "/", "_",
        "\\", "_",
        ":", "_",
        "*", "_",
        "?", "_",
        "\"", "_",
        "<", "_",
        ">", "_",
        "|", "_",
    )

    filename = replacer.Replace(filename)

    // Limit length
    if len(filename) > 255 {
        ext := filepath.Ext(filename)
        name := filename[:255-len(ext)]
        filename = name + ext
    }

    return filename
}

type FileUploadResult struct {
    URL       string
    Filename  string
    Size      int64
    MimeType  string
    Extension string
}
```

---

### Step 2: Add Handler

**File**: `interfaces/api/handlers/media_handler.go`

```go
// UploadFile handles file upload
// POST /upload/file
func (h *MediaHandler) UploadFile(c *fiber.Ctx) error {
    // Get file from request
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "No file provided",
            "error":   "FILE_REQUIRED",
        })
    }

    // Upload file
    result, err := h.fileUploadService.UploadFile(file)
    if err != nil {
        // Handle specific errors
        if strings.Contains(err.Error(), "exceeds maximum") {
            return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
                "success": false,
                "message": "File size exceeds maximum allowed size",
                "error":   "FILE_TOO_LARGE",
                "details": fiber.Map{
                    "maxSize":   50 * 1024 * 1024,
                    "maxSizeMB": 50,
                },
            })
        }

        if strings.Contains(err.Error(), "not allowed") {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "File type not allowed",
                "error":   "INVALID_FILE_TYPE",
            })
        }

        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to upload file",
            "error":   "UPLOAD_FAILED",
        })
    }

    // Save to database (optional - for tracking)
    media := &models.Media{
        UserID:   c.Locals("userId").(string),
        Type:     "file",
        URL:      result.URL,
        Filename: result.Filename,
        Size:     result.Size,
        MimeType: result.MimeType,
    }

    if err := h.mediaService.Create(media); err != nil {
        // Log error but don't fail request (file already uploaded)
        log.Printf("Failed to save media record: %v", err)
    }

    // Return success response
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "File uploaded successfully",
        "data": fiber.Map{
            "id":        media.ID,
            "url":       result.URL,
            "filename":  result.Filename,
            "size":      result.Size,
            "mimeType":  result.MimeType,
            "extension": result.Extension,
            "createdAt": media.CreatedAt,
        },
    })
}
```

---

### Step 3: Add Route

**File**: `interfaces/api/routes/media_routes.go`

```go
func SetupMediaRoutes(app *fiber.App, handler *handlers.MediaHandler) {
    media := app.Group("/api/v1/upload")

    // Apply authentication middleware
    media.Use(middleware.Protected())

    // Existing routes
    media.Post("/image", handler.UploadImage)
    media.Post("/video", handler.UploadVideo)

    // 🆕 New route
    media.Post("/file", handler.UploadFile)
}
```

---

### Step 4: Update Media Model (Optional)

**File**: `domain/models/media.go`

```go
type Media struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    UserID    string    `json:"userId" gorm:"not null;index"`
    Type      string    `json:"type" gorm:"not null"` // "image", "video", "file"
    URL       string    `json:"url" gorm:"not null"`
    Thumbnail string    `json:"thumbnail,omitempty"`

    // Image/Video specific
    Width     int       `json:"width,omitempty"`
    Height    int       `json:"height,omitempty"`
    Duration  float64   `json:"duration,omitempty"`

    // File specific
    Filename  string    `json:"filename,omitempty"`  // 🆕 Add this
    Extension string    `json:"extension,omitempty"` // 🆕 Add this

    // Common
    Size      int64     `json:"size" gorm:"not null"`
    MimeType  string    `json:"mimeType" gorm:"not null"`
    CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
```

---

## 🧪 Testing Examples

### 1. Test with cURL

```bash
# Test PDF upload
curl -X POST http://localhost:8080/api/v1/upload/file \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test-document.pdf"

# Test DOCX upload
curl -X POST http://localhost:8080/api/v1/upload/file \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test-document.docx"

# Test ZIP upload
curl -X POST http://localhost:8080/api/v1/upload/file \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test-archive.zip"

# Test error: file too large (51MB)
curl -X POST http://localhost:8080/api/v1/upload/file \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@large-file.pdf"

# Test error: invalid file type (.exe)
curl -X POST http://localhost:8080/api/v1/upload/file \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test.exe"
```

---

### 2. Test with Postman

**Request:**
```
POST http://localhost:8080/api/v1/upload/file
Headers:
  Authorization: Bearer YOUR_TOKEN
Body:
  form-data
    file: [Select File]
```

**Expected Response (Success):**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "id": "file-abc123",
    "url": "https://cdn.bunny.net/files/uuid-document.pdf",
    "filename": "document.pdf",
    "size": 5000000,
    "mimeType": "application/pdf",
    "extension": "pdf",
    "createdAt": "2024-01-01T10:00:00Z"
  }
}
```

---

### 3. Frontend Integration Test

**File**: `lib/services/api/media.service.ts`

```typescript
// Add to existing mediaService
export const mediaService = {
  // ... existing uploadImage, uploadVideo ...

  /**
   * Upload file (NEW)
   */
  uploadFile: async (
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<UploadFileResponse> => {
    try {
      // Client-side validation
      const maxSize = 50 * 1024 * 1024; // 50MB
      if (file.size > maxSize) {
        throw new Error('File size exceeds 50MB');
      }

      const allowedTypes = [
        'application/pdf',
        'application/msword',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
        'application/zip',
        'text/plain',
      ];

      if (!allowedTypes.includes(file.type)) {
        throw new Error('File type not allowed');
      }

      // Upload
      const formData = new FormData();
      formData.append('file', file);

      return await apiService.upload<UploadFileResponse>(
        '/upload/file',
        formData,
        onProgress
      );
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`File upload failed: ${error.message}`);
      }
      throw new Error('File upload failed');
    }
  },
};

// Add response type
export interface UploadFileResponse {
  success: boolean;
  message: string;
  data: {
    id: string;
    url: string;
    filename: string;
    size: number;
    mimeType: string;
    extension: string;
    createdAt: string;
  };
}
```

**Test in Component:**
```typescript
// Test upload file
const handleFileUpload = async (file: File) => {
  try {
    const response = await mediaService.uploadFile(file, (progress) => {
      console.log(`Upload progress: ${progress}%`);
    });

    console.log('File uploaded:', response.data.url);
    // Use the URL to send message or create post
  } catch (error) {
    console.error('Upload failed:', error);
  }
};
```

---

## 🔗 Integration with Chat Feature

### Updated Chat Flow

**Before (Current):**
```typescript
// ❌ Old way: Upload + send in one request
const handleSend = async (content: string, files: File[]) => {
  const formData = new FormData();
  formData.append('type', 'file');
  formData.append('content', content);
  files.forEach(f => formData.append('media[]', f));

  await chatService.sendMessage(conversationId, formData);
};
```

**After (New - Recommended):**
```typescript
// ✅ New way: Upload first, then send with URLs
const handleSend = async (content: string, files: File[]) => {
  try {
    // Step 1: Upload files with progress
    setUploading(true);
    const uploadedFiles = [];

    for (const file of files) {
      const response = await mediaService.uploadFile(file, (progress) => {
        setUploadProgress(progress); // Show progress bar
      });
      uploadedFiles.push(response.data);
    }

    // Step 2: Send message with file URLs
    await chatService.sendMessage({
      conversationId,
      type: 'file',
      content,
      mediaUrls: uploadedFiles.map(f => f.url),
      mediaMetadata: uploadedFiles.map(f => ({
        filename: f.filename,
        size: f.size,
        mimeType: f.mimeType,
      })),
    });

    setUploading(false);
  } catch (error) {
    setUploading(false);
    // Show error, allow retry upload
  }
};
```

**Benefits:**
- ✅ Show upload progress (better UX)
- ✅ Retry upload if fails (without resending message)
- ✅ Preview files before sending
- ✅ Upload multiple files in parallel

---

## 📊 Summary & Checklist

### What to Implement:

#### Backend (4-6 hours):
- [ ] Create `FileUploadService` in `infrastructure/storage/`
- [ ] Add `UploadFile()` handler in `media_handler.go`
- [ ] Add route `POST /upload/file` in `media_routes.go`
- [ ] Update `Media` model (add filename, extension fields)
- [ ] Add validation (MIME type, size, filename sanitization)
- [ ] Add error handling (413, 400, 500)
- [ ] Test with cURL + Postman
- [ ] Update API documentation

#### Frontend (2-3 hours):
- [ ] Add `uploadFile()` to `media.service.ts`
- [ ] Update `API.MEDIA.UPLOAD_FILE` constant
- [ ] Add `UploadFileResponse` type
- [ ] Update `ChatInput` to use new flow
- [ ] Add upload progress UI
- [ ] Test upload + retry flow

---

## 🎯 Priority & Timeline

| Task | Priority | Effort | When |
|------|----------|--------|------|
| **Implement POST /upload/file** | 🟡 MEDIUM | 4-6h | Before Chat launch |
| **Update Chat to use uploads** | 🟡 MEDIUM | 2-3h | After upload endpoint ready |
| **Add upload progress UI** | 🟢 LOW | 1-2h | Nice to have |

**Total Effort**: 6-11 hours (Backend + Frontend)
**Recommended Timeline**: Implement before Chat feature launch

---

## 📞 Questions & Support

### For Backend Team:

**Common Questions:**
1. **Q**: Should we create a new folder structure or reuse `/media/upload`?
   **A**: Recommend `/api/v1/upload/file` to match existing pattern

2. **Q**: Should we track uploads in database?
   **A**: Yes, recommended for cleanup and user media management

3. **Q**: What about file scanning (virus check)?
   **A**: Phase 2 - can add ClamAV or similar later

4. **Q**: Rate limiting?
   **A**: Yes, recommend 20 uploads/minute per user

### For Frontend Team:

**Common Questions:**
1. **Q**: Should we always use separate upload + send?
   **A**: Yes, recommended for better UX

2. **Q**: What about backward compatibility with current chat?
   **A**: Keep both methods during transition, migrate gradually

3. **Q**: Upload progress implementation?
   **A**: Use `XMLHttpRequest.upload.onprogress` or Axios progress callback

---

## ✅ Conclusion

### Summary:

**What's Missing:**
- ❌ `POST /upload/file` endpoint for file uploads

**What to Add:**
- ✅ New endpoint: `POST /api/v1/upload/file`
- ✅ File validation (MIME type, size, filename)
- ✅ Integration with **Bunny Storage** (not Stream!)
- ✅ Error handling
- ✅ Frontend service update

**Benefits:**
- ✅ Better UX (upload progress, preview, retry)
- ✅ Consistent API across Posts & Chat
- ✅ Shared infrastructure (DRY principle)
- ✅ Future-proof architecture

**Timeline:**
- Backend: 4-6 hours
- Frontend: 2-3 hours
- **Total: 6-11 hours**

---

## 🎯 Critical Reminders for Backend Team

### ⚠️ Two Different Upload Systems

1. **Images & Files → Bunny Storage**
   - `POST /media/upload/image` → BunnyStorageClient
   - `POST /upload/file` → BunnyStorageClient
   - Returns: Direct CDN URL
   - Example: `https://cdn.bunny.net/files/abc-document.pdf`

2. **Videos → Bunny Stream**
   - `POST /media/upload/video` → BunnyStreamClient (NOT Storage!)
   - Returns: Playlist URL (.m3u8) for HLS streaming
   - Example: `https://vz-b1631ae0-4c8.b-cdn.net/{videoId}/playlist.m3u8`
   - See: `vdo_stream_implement.md` for complete implementation

### Environment Variables

```bash
# Bunny Storage (for images & files)
BUNNY_STORAGE_API_KEY=your-storage-api-key
BUNNY_STORAGE_ZONE=your-storage-zone

# Bunny Stream (for videos only!)
BUNNY_STREAM_API_KEY=4c1ec80d-8809-4852-89f647b0430d-c9ad-4914
BUNNY_STREAM_LIBRARY_ID=533535
BUNNY_STREAM_CDN=vz-b1631ae0-4c8.b-cdn.net
```

### Do NOT:
- ❌ Upload videos to Bunny Storage (use Stream instead!)
- ❌ Upload files to Bunny Stream (use Storage instead!)
- ❌ Mix Storage and Stream clients
- ❌ Return direct MP4 URLs for videos (use HLS playlist!)

---

**Document Version**: 2.0.0
**Created**: 2025-01-07
**Updated**: 2025-01-07 (Added Storage vs Stream clarification)
**Status**: ✅ Ready for Implementation
**Priority**: 🟡 MEDIUM (Recommended)

---

## 📖 Quick Reference Guide

### Which Upload System to Use?

| File Type | Upload Endpoint | Storage System | Client to Use | Returns |
|-----------|----------------|----------------|---------------|---------|
| **Images** (JPG, PNG, GIF, WEBP) | `POST /media/upload/image` | Bunny Storage | BunnyStorageClient | Direct CDN URL |
| **Videos** (MP4, MOV, WEBM) | `POST /media/upload/video` | **Bunny Stream** | **BunnyStreamClient** | **Playlist URL (.m3u8)** |
| **Files** (PDF, DOC, XLS, ZIP) | `POST /upload/file` | Bunny Storage | BunnyStorageClient | Direct CDN URL |

### Related Documentation

- **Video Implementation**: `vdo_stream_implement.md`
- **Chat Backend Status**: `final_result_chat_from_backend/`
- **Chat Integration**: `chat_integration_checklist.md`

---

**End of Specification** 📄
