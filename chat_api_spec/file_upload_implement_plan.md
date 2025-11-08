# File Upload Implementation Plan

> **Status**: 🔴 Ready for Implementation
> **Priority**: HIGH - Phase 1 MVP (Images only)
> **Estimated Time**: 2 days
> **Last Updated**: 2025-01-07

---

## 📋 Overview

Implement file upload functionality for chat messages supporting:
- **Phase 1.0**: Images (JPEG, PNG, GIF, WebP)
- **Phase 1.1**: Videos (MP4, MOV) + Files (PDF, DOCX, etc.)

---

## 1. Current State

### ✅ What We Have
- ✅ Message model with `type` enum and `media` JSONB
- ✅ SendMessageRequest DTO with Media array
- ✅ Bunny Storage service (already integrated)
- ✅ REST endpoint: `POST /chat/conversations/:id/messages`

### ❌ What We Need
- ❌ multipart/form-data handler
- ❌ File validation (size, type, count)
- ❌ Thumbnail generation for images/videos
- ❌ Metadata extraction (dimensions, duration)
- ❌ Integration with Bunny Storage upload
- ❌ JSONB media array population

---

## 2. Architecture

### 2.1 Upload Flow

```
Client (Frontend)
    │
    │ 1. User selects file(s)
    │ 2. Validate client-side (size, type)
    │ 3. Show preview
    │
    ▼
POST /chat/conversations/:id/messages
FormData:
  - type: "image"
  - content: "Check this out!" (optional)
  - media[0]: File blob
  - media[1]: File blob (optional)
    │
    ▼
MessageHandler.SendMessage()
    │
    ├─ Check Content-Type
    │  ├─ application/json → SendTextMessage()
    │  └─ multipart/form-data → SendMediaMessage()
    │
    ▼
SendMediaMessage()
    │
    ├─ 1. Parse FormData (type, content, files)
    ├─ 2. Validate files (size, type, count)
    │      ├─ Image: max 10MB, max 10 files
    │      ├─ Video: max 100MB, max 1 file
    │      └─ File: max 50MB, max 5 files
    │
    ├─ 3. For each file:
    │      ├─ Upload to Bunny Storage
    │      ├─ Generate thumbnail (if image/video)
    │      ├─ Extract metadata
    │      └─ Create MessageMedia object
    │
    ├─ 4. Create Message
    │      ├─ type = "image" | "video" | "file"
    │      ├─ content = caption (nullable)
    │      └─ media = MessageMedia[] (JSONB)
    │
    ├─ 5. Save to database
    │
    └─ 6. Return response with media URLs
```

---

## 3. Implementation Plan

### Day 1: Basic Image Upload (8 hours)

#### Morning (4h): Core Upload Logic

**Step 1.1: Update MessageHandler**

Location: `interfaces/api/handlers/message_handler.go`

```go
// Update SendMessage to detect Content-Type
func (h *MessageHandler) SendMessage(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uuid.UUID)
    conversationID, err := uuid.Parse(c.Params("conversationId"))
    if err != nil {
        return utils.ValidationErrorResponse(c, "Invalid conversation ID")
    }

    // Check Content-Type
    contentType := string(c.Request().Header.ContentType())

    if strings.Contains(contentType, "multipart/form-data") {
        // Handle media upload
        return h.sendMediaMessage(c, userID, conversationID)
    } else {
        // Handle text message (existing logic)
        return h.sendTextMessage(c, userID, conversationID)
    }
}
```

---

**Step 1.2: Implement sendMediaMessage()**

```go
func (h *MessageHandler) sendMediaMessage(c *fiber.Ctx, userID uuid.UUID, conversationID uuid.UUID) error {
    // 1. Parse form data
    messageType := c.FormValue("type") // "image", "video", "file"
    content := c.FormValue("content")  // Optional caption

    // 2. Get uploaded files
    form, err := c.MultipartForm()
    if err != nil {
        return utils.ValidationErrorResponse(c, "Failed to parse form data")
    }

    files := form.File["media"]
    if len(files) == 0 {
        return utils.ValidationErrorResponse(c, "No files uploaded")
    }

    // 3. Validate message type
    validTypes := map[string]bool{"image": true, "video": true, "file": true}
    if !validTypes[messageType] {
        return utils.ValidationErrorResponse(c, "Invalid message type")
    }

    // 4. Validate file count
    maxFiles := map[string]int{"image": 10, "video": 1, "file": 5}
    if len(files) > maxFiles[messageType] {
        return utils.ValidationErrorResponse(c, fmt.Sprintf("Maximum %d files allowed for type %s", maxFiles[messageType], messageType))
    }

    // 5. Process each file
    mediaItems := make([]dto.MessageMediaItem, 0, len(files))

    for _, fileHeader := range files {
        // Open file
        file, err := fileHeader.Open()
        if err != nil {
            return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to open file", err)
        }
        defer file.Close()

        // Validate file
        if err := h.validateFile(fileHeader, messageType); err != nil {
            return utils.ValidationErrorResponse(c, err.Error())
        }

        // Upload to Bunny Storage
        mediaItem, err := h.uploadFile(file, fileHeader, messageType)
        if err != nil {
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload file", err)
        }

        mediaItems = append(mediaItems, mediaItem)
    }

    // 6. Create message via service
    var contentPtr *string
    if content != "" {
        contentPtr = &content
    }

    req := &dto.SendMessageRequest{
        ConversationID: conversationID,
        Type:           messageType,
        Content:        contentPtr,
        Media:          mediaItems,
    }

    message, err := h.messageService.SendMessage(c.Context(), userID, req)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send message", err)
    }

    return utils.SuccessResponse(c, "Message sent successfully", message)
}
```

---

**Step 1.3: File Validation**

```go
func (h *MessageHandler) validateFile(fileHeader *multipart.FileHeader, messageType string) error {
    // Check file size
    maxSizes := map[string]int64{
        "image": 10 * 1024 * 1024,  // 10MB
        "video": 100 * 1024 * 1024, // 100MB
        "file":  50 * 1024 * 1024,  // 50MB
    }

    if fileHeader.Size > maxSizes[messageType] {
        return fmt.Errorf("file size exceeds maximum %d MB", maxSizes[messageType]/(1024*1024))
    }

    // Check MIME type
    allowedMimeTypes := map[string][]string{
        "image": {"image/jpeg", "image/png", "image/gif", "image/webp"},
        "video": {"video/mp4", "video/quicktime", "video/x-matroska"},
        "file":  {
            "application/pdf",
            "application/msword",
            "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
            "application/vnd.ms-excel",
            "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
            "application/zip",
            "text/plain",
        },
    }

    // Detect MIME type from file content
    file, err := fileHeader.Open()
    if err != nil {
        return err
    }
    defer file.Close()

    buffer := make([]byte, 512)
    _, err = file.Read(buffer)
    if err != nil {
        return err
    }

    mimeType := http.DetectContentType(buffer)

    // Check if MIME type is allowed
    allowed := false
    for _, allowedType := range allowedMimeTypes[messageType] {
        if strings.HasPrefix(mimeType, allowedType) {
            allowed = true
            break
        }
    }

    if !allowed {
        return fmt.Errorf("file type %s not allowed for %s messages", mimeType, messageType)
    }

    return nil
}
```

---

#### Afternoon (4h): Bunny Storage Integration

**Step 1.4: Upload to Bunny Storage**

Create new service: `infrastructure/storage/media_upload_service.go`

```go
package storage

import (
    "bytes"
    "context"
    "fmt"
    "image"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
    "io"
    "mime/multipart"
    "path/filepath"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/nfnt/resize"
)

type MediaUploadService struct {
    bunnyClient *BunnyStorageClient // Assuming already exists
}

type UploadResult struct {
    URL       string
    Thumbnail string
    MimeType  string
    Size      int64
    Width     int
    Height    int
    Duration  int // For videos
}

func NewMediaUploadService(bunnyClient *BunnyStorageClient) *MediaUploadService {
    return &MediaUploadService{
        bunnyClient: bunnyClient,
    }
}

func (s *MediaUploadService) UploadImage(ctx context.Context, file multipart.File, filename string) (*UploadResult, error) {
    // Read file into buffer
    var buf bytes.Buffer
    fileSize, err := io.Copy(&buf, file)
    if err != nil {
        return nil, err
    }

    // Decode image to get dimensions
    img, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
    if err != nil {
        return nil, err
    }

    width := img.Bounds().Dx()
    height := img.Bounds().Dy()

    // Generate unique filename
    ext := filepath.Ext(filename)
    uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
    uploadPath := fmt.Sprintf("chat/images/%s", uniqueName)

    // Upload original to Bunny
    originalURL, err := s.bunnyClient.Upload(ctx, uploadPath, buf.Bytes())
    if err != nil {
        return nil, err
    }

    // Generate thumbnail (max 300x300)
    thumbnail := resize.Thumbnail(300, 300, img, resize.Lanczos3)

    // Encode thumbnail
    var thumbBuf bytes.Buffer
    if err := s.encodeImage(&thumbBuf, thumbnail, format); err != nil {
        return nil, err
    }

    // Upload thumbnail
    thumbPath := fmt.Sprintf("chat/thumbnails/%s", uniqueName)
    thumbnailURL, err := s.bunnyClient.Upload(ctx, thumbPath, thumbBuf.Bytes())
    if err != nil {
        return nil, err
    }

    return &UploadResult{
        URL:       originalURL,
        Thumbnail: thumbnailURL,
        MimeType:  fmt.Sprintf("image/%s", format),
        Size:      fileSize,
        Width:     width,
        Height:    height,
    }, nil
}

func (s *MediaUploadService) encodeImage(w io.Writer, img image.Image, format string) error {
    switch format {
    case "jpeg":
        return jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
    case "png":
        return png.Encode(w, img)
    case "gif":
        return gif.Encode(w, img, nil)
    default:
        return fmt.Errorf("unsupported format: %s", format)
    }
}
```

---

**Step 1.5: Update uploadFile() in Handler**

```go
func (h *MessageHandler) uploadFile(file multipart.File, fileHeader *multipart.FileHeader, messageType string) (dto.MessageMediaItem, error) {
    switch messageType {
    case "image":
        result, err := h.mediaUploadService.UploadImage(context.Background(), file, fileHeader.Filename)
        if err != nil {
            return dto.MessageMediaItem{}, err
        }

        return dto.MessageMediaItem{
            URL:       result.URL,
            Thumbnail: &result.Thumbnail,
            Type:      "image",
            MimeType:  &result.MimeType,
            Size:      &result.Size,
            Width:     &result.Width,
            Height:    &result.Height,
        }, nil

    case "video":
        // TODO: Implement video upload
        return dto.MessageMediaItem{}, fmt.Errorf("video upload not implemented yet")

    case "file":
        // TODO: Implement file upload
        return dto.MessageMediaItem{}, fmt.Errorf("file upload not implemented yet")

    default:
        return dto.MessageMediaItem{}, fmt.Errorf("unsupported message type: %s", messageType)
    }
}
```

---

**Step 1.6: Update Container**

```go
// infrastructure/container/container.go

type Container struct {
    // ... existing fields ...
    MediaUploadService *storage.MediaUploadService
}

func NewContainer() (*Container, error) {
    // ... existing initialization ...

    // Initialize MediaUploadService
    mediaUploadService := storage.NewMediaUploadService(bunnyClient)

    return &Container{
        // ... existing fields ...
        MediaUploadService: mediaUploadService,
    }, nil
}
```

---

**Step 1.7: Update MessageHandler Constructor**

```go
type MessageHandler struct {
    messageService      services.MessageService
    mediaUploadService  *storage.MediaUploadService
}

func NewMessageHandler(messageService services.MessageService, mediaUploadService *storage.MediaUploadService) *MessageHandler {
    return &MessageHandler{
        messageService:     messageService,
        mediaUploadService: mediaUploadService,
    }
}
```

---

### Day 2: Testing & Video/File Support (8 hours)

#### Morning (4h): Testing

**Step 2.1: Unit Tests**

```go
// interfaces/api/handlers/message_handler_test.go

func TestSendMediaMessage_ValidImage(t *testing.T) {
    // Create test image file
    img := createTestImage(t, 800, 600)

    // Create multipart form
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    // Add type field
    writer.WriteField("type", "image")
    writer.WriteField("content", "Test caption")

    // Add image file
    part, err := writer.CreateFormFile("media", "test.jpg")
    require.NoError(t, err)

    jpeg.Encode(part, img, nil)
    writer.Close()

    // Create request
    req := httptest.NewRequest("POST", "/chat/conversations/conv-001/messages", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // Execute
    resp := executeRequest(req)

    // Assert
    assert.Equal(t, 201, resp.StatusCode)

    var result map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &result)

    assert.True(t, result["success"].(bool))
    assert.NotNil(t, result["data"])

    message := result["data"].(map[string]interface{})
    assert.Equal(t, "image", message["type"])
    assert.NotEmpty(t, message["media"])
}

func TestSendMediaMessage_FileTooLarge(t *testing.T) {
    // Create 11MB file (exceeds 10MB limit)
    largeFile := make([]byte, 11*1024*1024)

    // Create multipart form
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    writer.WriteField("type", "image")

    part, _ := writer.CreateFormFile("media", "large.jpg")
    part.Write(largeFile)
    writer.Close()

    // Create request
    req := httptest.NewRequest("POST", "/chat/conversations/conv-001/messages", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // Execute
    resp := executeRequest(req)

    // Assert
    assert.Equal(t, 400, resp.StatusCode)

    var result map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &result)

    assert.False(t, result["success"].(bool))
    assert.Contains(t, result["message"].(string), "exceeds maximum")
}
```

---

**Step 2.2: Integration Tests**

```go
// tests/integration/file_upload_test.go

func TestFileUpload_EndToEnd(t *testing.T) {
    // Setup test server
    server := setupTestServer()
    defer server.Close()

    // Login and get token
    token := loginTestUser(t, "user1")

    // Create conversation
    convID := createTestConversation(t, token)

    // Upload image
    imageFile := loadTestImage(t, "testdata/sample.jpg")

    resp := uploadImageMessage(t, server.URL, token, convID, imageFile, "Nice photo!")

    assert.Equal(t, 201, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    message := result["data"].(map[string]interface{})
    media := message["media"].([]interface{})[0].(map[string]interface{})

    // Verify CDN URLs are accessible
    assert.True(t, strings.HasPrefix(media["url"].(string), "https://cdn.bunny.net"))
    assert.NotEmpty(t, media["thumbnail"])
    assert.NotZero(t, media["width"])
    assert.NotZero(t, media["height"])

    // Download and verify
    downloadedImage := downloadFile(t, media["url"].(string))
    assert.NotEmpty(t, downloadedImage)
}
```

---

#### Afternoon (4h): Video & File Support

**Step 2.3: Implement Video Upload**

```go
// infrastructure/storage/media_upload_service.go

func (s *MediaUploadService) UploadVideo(ctx context.Context, file multipart.File, filename string) (*UploadResult, error) {
    // Read file
    var buf bytes.Buffer
    fileSize, err := io.Copy(&buf, file)
    if err != nil {
        return nil, err
    }

    // Generate unique filename
    ext := filepath.Ext(filename)
    uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
    uploadPath := fmt.Sprintf("chat/videos/%s", uniqueName)

    // Upload to Bunny
    videoURL, err := s.bunnyClient.Upload(ctx, uploadPath, buf.Bytes())
    if err != nil {
        return nil, err
    }

    // Extract video metadata using ffmpeg
    metadata, err := s.extractVideoMetadata(buf.Bytes())
    if err != nil {
        log.Printf("Failed to extract video metadata: %v", err)
        metadata = &VideoMetadata{} // Use empty metadata
    }

    // Generate thumbnail (first frame)
    thumbnailURL, err := s.generateVideoThumbnail(ctx, buf.Bytes(), uniqueName)
    if err != nil {
        log.Printf("Failed to generate video thumbnail: %v", err)
        thumbnailURL = "" // No thumbnail
    }

    return &UploadResult{
        URL:       videoURL,
        Thumbnail: thumbnailURL,
        MimeType:  "video/mp4",
        Size:      fileSize,
        Width:     metadata.Width,
        Height:    metadata.Height,
        Duration:  metadata.Duration,
    }, nil
}

// Simple video metadata extraction (without ffmpeg for MVP)
func (s *MediaUploadService) extractVideoMetadata(data []byte) (*VideoMetadata, error) {
    // For MVP: return default values
    // For production: use ffmpeg or similar library
    return &VideoMetadata{
        Width:    1920,
        Height:   1080,
        Duration: 0, // Unknown
    }, nil
}

func (s *MediaUploadService) generateVideoThumbnail(ctx context.Context, videoData []byte, videoName string) (string, error) {
    // For MVP: skip thumbnail generation
    // For production: use ffmpeg to extract first frame
    return "", nil
}
```

---

**Step 2.4: Implement File Upload**

```go
func (s *MediaUploadService) UploadFile(ctx context.Context, file multipart.File, filename string, mimeType string) (*UploadResult, error) {
    // Read file
    var buf bytes.Buffer
    fileSize, err := io.Copy(&buf, file)
    if err != nil {
        return nil, err
    }

    // Generate unique filename
    ext := filepath.Ext(filename)
    uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
    uploadPath := fmt.Sprintf("chat/files/%s", uniqueName)

    // Upload to Bunny
    fileURL, err := s.bunnyClient.Upload(ctx, uploadPath, buf.Bytes())
    if err != nil {
        return nil, err
    }

    return &UploadResult{
        URL:      fileURL,
        MimeType: mimeType,
        Size:     fileSize,
    }, nil
}
```

---

**Step 2.5: Update uploadFile() with all types**

```go
func (h *MessageHandler) uploadFile(file multipart.File, fileHeader *multipart.FileHeader, messageType string) (dto.MessageMediaItem, error) {
    ctx := context.Background()

    switch messageType {
    case "image":
        result, err := h.mediaUploadService.UploadImage(ctx, file, fileHeader.Filename)
        if err != nil {
            return dto.MessageMediaItem{}, err
        }

        return dto.MessageMediaItem{
            URL:       result.URL,
            Thumbnail: &result.Thumbnail,
            Type:      "image",
            Filename:  &fileHeader.Filename,
            MimeType:  &result.MimeType,
            Size:      &result.Size,
            Width:     &result.Width,
            Height:    &result.Height,
        }, nil

    case "video":
        result, err := h.mediaUploadService.UploadVideo(ctx, file, fileHeader.Filename)
        if err != nil {
            return dto.MessageMediaItem{}, err
        }

        item := dto.MessageMediaItem{
            URL:      result.URL,
            Type:     "video",
            Filename: &fileHeader.Filename,
            MimeType: &result.MimeType,
            Size:     &result.Size,
        }

        if result.Thumbnail != "" {
            item.Thumbnail = &result.Thumbnail
        }
        if result.Width > 0 {
            item.Width = &result.Width
            item.Height = &result.Height
        }
        if result.Duration > 0 {
            item.Duration = &result.Duration
        }

        return item, nil

    case "file":
        // Detect MIME type
        buffer := make([]byte, 512)
        file.Read(buffer)
        file.Seek(0, 0) // Reset
        mimeType := http.DetectContentType(buffer)

        result, err := h.mediaUploadService.UploadFile(ctx, file, fileHeader.Filename, mimeType)
        if err != nil {
            return dto.MessageMediaItem{}, err
        }

        return dto.MessageMediaItem{
            URL:      result.URL,
            Type:     "file",
            Filename: &fileHeader.Filename,
            MimeType: &result.MimeType,
            Size:     &result.Size,
        }, nil

    default:
        return dto.MessageMediaItem{}, fmt.Errorf("unsupported message type: %s", messageType)
    }
}
```

---

## 4. Testing Checklist

### Unit Tests
- [ ] File validation (size, type)
- [ ] Image upload success
- [ ] Video upload success
- [ ] File upload success
- [ ] Multiple files upload
- [ ] File too large error
- [ ] Invalid file type error
- [ ] Thumbnail generation

### Integration Tests
- [ ] End-to-end image upload
- [ ] End-to-end video upload
- [ ] End-to-end file upload
- [ ] Bunny Storage integration
- [ ] Database persistence
- [ ] Media URLs accessibility

### Manual Tests
- [ ] Upload from Postman
- [ ] Upload 10 images at once
- [ ] Upload large video (90MB)
- [ ] Upload PDF file
- [ ] Check CDN URLs work
- [ ] Check thumbnails display correctly

---

## 5. Timeline Summary

| Day | Task | Hours | Status |
|-----|------|-------|--------|
| **Day 1 Morning** | Core upload logic + validation | 4h | ⏳ Pending |
| **Day 1 Afternoon** | Bunny Storage integration | 4h | ⏳ Pending |
| **Day 2 Morning** | Testing (unit + integration) | 4h | ⏳ Pending |
| **Day 2 Afternoon** | Video & File support | 4h | ⏳ Pending |

**Total**: 16 hours (2 days)

---

## 6. Dependencies

### Required Libraries

```bash
go get github.com/nfnt/resize  # Image resizing
```

### Optional (for production):
```bash
# FFmpeg for video processing (install separately)
apt-get install ffmpeg

go get github.com/u2takey/ffmpeg-go  # Go bindings
```

---

## 7. Success Criteria

### Functional
- ✅ Upload images (JPEG, PNG, GIF, WebP)
- ✅ Upload videos (MP4, MOV)
- ✅ Upload files (PDF, DOCX, etc.)
- ✅ Generate thumbnails for images
- ✅ Multiple files support (up to 10 images)
- ✅ File validation (size, type)

### Performance
- ✅ Upload time < 3 seconds (10MB image)
- ✅ Thumbnail generation < 500ms
- ✅ CDN URLs accessible < 100ms

### Quality
- ✅ All tests pass
- ✅ No memory leaks
- ✅ Proper error handling
- ✅ Clean code structure

---

## 8. Rollout Checklist

### Before Starting
- [ ] Review plan with team
- [ ] Setup test images/videos
- [ ] Verify Bunny Storage credentials
- [ ] Create feature branch

### Implementation
- [ ] Day 1: Core + Bunny integration ✅
- [ ] Day 2: Testing + Video/File support ✅

### Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual upload tests
- [ ] Check CDN accessibility

### Deployment
- [ ] Merge to develop
- [ ] Deploy to staging
- [ ] Test on staging
- [ ] Deploy to production

---

## 9. Phase 1.1 Extensions (Optional - Week 2)

### Video Enhancements
- [ ] Extract video duration with ffmpeg
- [ ] Generate video thumbnail (first frame)
- [ ] Support more video formats (MKV, AVI)

### File Enhancements
- [ ] Preview images for PDFs
- [ ] File type icons
- [ ] Virus scanning

### Performance
- [ ] Parallel upload for multiple files
- [ ] Client-side image compression
- [ ] Progress bar support

---

**Document Status:** ✅ Ready for Implementation
**Next Action:** Review → Start Day 1
**Questions?** Contact Backend Team
