# Bunny Storage Integration Guide

## 📦 Overview

คู่มือการ setup และใช้งาน Bunny Storage CDN สำหรับ media upload แทน AWS S3

**Bunny.net** เป็น CDN service ที่:
- ⚡ เร็วกว่า AWS S3
- 💰 ถูกกว่า AWS S3 มาก
- 🌏 มี CDN ทั่วโลก
- 📦 มี Storage + CDN ในที่เดียว
- 🔧 API ง่ายต่อการใช้งาน

---

## 🎯 Architecture

```
Client Upload Request
        ↓
Go Fiber Backend
        ↓
Process Image/Video
        ↓
Upload to Bunny Storage
        ↓
Get CDN URL
        ↓
Save to PostgreSQL
        ↓
Return CDN URL to Client
        ↓
Client loads media from Bunny CDN
```

---

## 🔑 Step 1: Get Bunny Storage Credentials

### 1.1 Create Bunny.net Account
1. Go to https://bunny.net
2. Sign up for free account
3. Verify email

### 1.2 Create Storage Zone
1. Go to **Storage** → **Add Storage Zone**
2. Enter Storage Zone Name: `social-media-storage`
3. Select Region: Choose closest to your users
   - Falkenstein (Europe)
   - New York (North America)
   - Singapore (Asia)
   - Sydney (Oceania)
4. Click **Add Storage Zone**

### 1.3 Get Access Credentials
1. Click on your storage zone
2. Go to **FTP & API Access**
3. Copy:
   - **Storage Zone Name:** `social-media-storage`
   - **Password (API Key):** `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`
   - **Hostname:** `storage.bunnycdn.com`

### 1.4 Connect CDN (Pull Zone)
1. Go to **CDN** → **Add Pull Zone**
2. Link to your storage zone
3. Get CDN URL: `https://your-storage.b-cdn.net`

---

## ⚙️ Step 2: Configure in Go Fiber

### 2.1 Update Environment Variables

**File:** `.env`

```env
# Bunny Storage Configuration
BUNNY_STORAGE_ZONE=social-media-storage
BUNNY_ACCESS_KEY=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
BUNNY_BASE_URL=https://storage.bunnycdn.com
BUNNY_CDN_URL=https://your-storage.b-cdn.net
```

### 2.2 Update Config Struct

**File:** `pkg/config/config.go`

```go
type BunnyConfig struct {
    StorageZone string
    AccessKey   string
    BaseURL     string
    CDNUrl      string
}

func LoadConfig() (*Config, error) {
    // ... existing code ...

    config := &Config{
        // ... other configs ...

        Bunny: BunnyConfig{
            StorageZone: getEnv("BUNNY_STORAGE_ZONE", ""),
            AccessKey:   getEnv("BUNNY_ACCESS_KEY", ""),
            BaseURL:     getEnv("BUNNY_BASE_URL", "https://storage.bunnycdn.com"),
            CDNUrl:      getEnv("BUNNY_CDN_URL", ""),
        },
    }

    return config, nil
}
```

---

## 🔧 Step 3: Implement Bunny Storage Service

### 3.1 Create Storage Service Interface

**File:** `domain/services/storage_service.go`

```go
package services

import (
    "context"
    "mime/multipart"
)

type StorageService interface {
    Upload(ctx context.Context, file multipart.File, filename string) (string, error)
    UploadBytes(ctx context.Context, data []byte, filename string) (string, error)
    Delete(ctx context.Context, filename string) error
    GetURL(filename string) string
}
```

### 3.2 Implement Bunny Storage Service

**File:** `infrastructure/storage/bunny_storage.go`

```go
package storage

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "path/filepath"
    "time"

    "your-project/pkg/config"
)

type BunnyStorage struct {
    storageZone string
    accessKey   string
    baseURL     string
    cdnURL      string
    httpClient  *http.Client
}

func NewBunnyStorage(config config.BunnyConfig) *BunnyStorage {
    return &BunnyStorage{
        storageZone: config.StorageZone,
        accessKey:   config.AccessKey,
        baseURL:     config.BaseURL,
        cdnURL:      config.CDNUrl,
        httpClient: &http.Client{
            Timeout: 60 * time.Second,
        },
    }
}

// Upload file to Bunny Storage
func (b *BunnyStorage) Upload(ctx context.Context, file multipart.File, filename string) (string, error) {
    // Read file content
    fileBytes, err := io.ReadAll(file)
    if err != nil {
        return "", fmt.Errorf("failed to read file: %w", err)
    }

    return b.UploadBytes(ctx, fileBytes, filename)
}

// Upload bytes to Bunny Storage
func (b *BunnyStorage) UploadBytes(ctx context.Context, data []byte, filename string) (string, error) {
    // Build upload URL
    // Format: https://storage.bunnycdn.com/{storage-zone}/{path}/{filename}
    uploadURL := fmt.Sprintf("%s/%s/%s", b.baseURL, b.storageZone, filename)

    // Create HTTP request
    req, err := http.NewRequestWithContext(ctx, "PUT", uploadURL, bytes.NewReader(data))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("AccessKey", b.accessKey)
    req.Header.Set("Content-Type", "application/octet-stream")

    // Send request
    resp, err := b.httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to upload: %w", err)
    }
    defer resp.Body.Close()

    // Check response
    if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
    }

    // Return CDN URL
    cdnURL := fmt.Sprintf("%s/%s", b.cdnURL, filename)
    return cdnURL, nil
}

// Delete file from Bunny Storage
func (b *BunnyStorage) Delete(ctx context.Context, filename string) error {
    // Build delete URL
    deleteURL := fmt.Sprintf("%s/%s/%s", b.baseURL, b.storageZone, filename)

    // Create HTTP request
    req, err := http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("AccessKey", b.accessKey)

    // Send request
    resp, err := b.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to delete: %w", err)
    }
    defer resp.Body.Close()

    // Check response
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
    }

    return nil
}

// Get CDN URL for filename
func (b *BunnyStorage) GetURL(filename string) string {
    return fmt.Sprintf("%s/%s", b.cdnURL, filename)
}
```

---

## 🖼️ Step 4: Image Processing

### 4.1 Install Image Processing Library

```bash
go get -u github.com/disintegration/imaging
```

### 4.2 Create Image Processor

**File:** `pkg/utils/image_processor.go`

```go
package utils

import (
    "bytes"
    "fmt"
    "image"
    "mime/multipart"

    "github.com/disintegration/imaging"
)

// ProcessImage resizes and compresses image
func ProcessImage(file multipart.File, maxWidth, maxHeight, quality int) ([]byte, error) {
    // Decode image
    img, err := imaging.Decode(file)
    if err != nil {
        return nil, fmt.Errorf("failed to decode image: %w", err)
    }

    // Get dimensions
    bounds := img.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()

    // Resize if needed
    if width > maxWidth || height > maxHeight {
        img = imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
    }

    // Encode as JPEG
    var buf bytes.Buffer
    if err := imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
        return nil, fmt.Errorf("failed to encode image: %w", err)
    }

    return buf.Bytes(), nil
}

// GenerateThumbnail creates square thumbnail
func GenerateThumbnail(file multipart.File, size int) ([]byte, error) {
    // Decode image
    img, err := imaging.Decode(file)
    if err != nil {
        return nil, fmt.Errorf("failed to decode image: %w", err)
    }

    // Create square thumbnail (crop from center)
    thumb := imaging.Fill(img, size, size, imaging.Center, imaging.Lanczos)

    // Encode as JPEG
    var buf bytes.Buffer
    if err := imaging.Encode(&buf, thumb, imaging.JPEG, imaging.JPEGQuality(80)); err != nil {
        return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
    }

    return buf.Bytes(), nil
}

// GetImageDimensions returns width and height
func GetImageDimensions(file multipart.File) (int, int, error) {
    img, _, err := image.DecodeConfig(file)
    if err != nil {
        return 0, 0, err
    }
    return img.Width, img.Height, nil
}
```

---

## 📤 Step 5: Implement Media Upload Handler

### 5.1 Media Service

**File:** `application/serviceimpl/media_service_impl.go`

```go
package serviceimpl

import (
    "context"
    "fmt"
    "mime/multipart"
    "path/filepath"
    "time"

    "github.com/google/uuid"
    "your-project/domain/models"
    "your-project/domain/repositories"
    "your-project/domain/services"
    "your-project/pkg/utils"
)

type MediaServiceImpl struct {
    mediaRepo      repositories.MediaRepository
    storageService services.StorageService
}

func NewMediaService(mediaRepo repositories.MediaRepository, storage services.StorageService) services.MediaService {
    return &MediaServiceImpl{
        mediaRepo:      mediaRepo,
        storageService: storage,
    }
}

func (s *MediaServiceImpl) UploadMedia(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID uuid.UUID) (*models.Media, error) {
    // Get file extension
    ext := filepath.Ext(header.Filename)

    // Generate unique filename
    timestamp := time.Now().Unix()
    filename := fmt.Sprintf("%s_%d_%s%s", userID.String(), timestamp, uuid.New().String()[:8], ext)

    // Detect mime type
    mimeType := header.Header.Get("Content-Type")

    var processedData []byte
    var thumbnailData []byte
    var width, height int
    var err error

    if isImage(mimeType) {
        // Process image
        processedData, err = utils.ProcessImage(file, 1920, 1920, 85)
        if err != nil {
            return nil, err
        }

        // Generate thumbnail
        file.Seek(0, 0) // Reset file pointer
        thumbnailData, err = utils.GenerateThumbnail(file, 200)
        if err != nil {
            return nil, err
        }

        // Get dimensions
        file.Seek(0, 0)
        width, height, _ = utils.GetImageDimensions(file)

        // Upload main image
        mainPath := fmt.Sprintf("images/%s", filename)
        url, err := s.storageService.UploadBytes(ctx, processedData, mainPath)
        if err != nil {
            return nil, err
        }

        // Upload thumbnail
        thumbPath := fmt.Sprintf("thumbnails/%s", filename)
        thumbnailURL, err := s.storageService.UploadBytes(ctx, thumbnailData, thumbPath)
        if err != nil {
            return nil, err
        }

        // Create media record
        media := &models.Media{
            ID:        uuid.New(),
            UserID:    userID,
            Type:      "image",
            FileName:  header.Filename,
            MimeType:  mimeType,
            Size:      int64(len(processedData)),
            URL:       url,
            Thumbnail: thumbnailURL,
            Width:     width,
            Height:    height,
            CreatedAt: time.Now(),
        }

        if err := s.mediaRepo.Create(ctx, media); err != nil {
            return nil, err
        }

        return media, nil
    }

    // Handle videos (simplified - just upload)
    if isVideo(mimeType) {
        file.Seek(0, 0)
        videoPath := fmt.Sprintf("videos/%s", filename)
        url, err := s.storageService.Upload(ctx, file, videoPath)
        if err != nil {
            return nil, err
        }

        media := &models.Media{
            ID:       uuid.New(),
            UserID:   userID,
            Type:     "video",
            FileName: header.Filename,
            MimeType: mimeType,
            Size:     header.Size,
            URL:      url,
            CreatedAt: time.Now(),
        }

        if err := s.mediaRepo.Create(ctx, media); err != nil {
            return nil, err
        }

        return media, nil
    }

    return nil, fmt.Errorf("unsupported file type: %s", mimeType)
}

func isImage(mimeType string) bool {
    return mimeType == "image/jpeg" || mimeType == "image/png" || mimeType == "image/gif" || mimeType == "image/webp"
}

func isVideo(mimeType string) bool {
    return mimeType == "video/mp4" || mimeType == "video/webm" || mimeType == "video/quicktime"
}
```

---

## 🧪 Step 6: Testing

### 6.1 Test Upload

```bash
# Upload image
curl -X POST http://localhost:3000/api/v1/media/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "files=@test-image.jpg"

# Expected response:
{
  "success": true,
  "message": "อัปโหลดไฟล์สำเร็จ",
  "data": {
    "media": [
      {
        "id": "uuid",
        "type": "image",
        "url": "https://your-storage.b-cdn.net/images/filename.jpg",
        "thumbnail": "https://your-storage.b-cdn.net/thumbnails/filename.jpg",
        "width": 1920,
        "height": 1080,
        "size": 245678
      }
    ]
  }
}
```

### 6.2 Test CDN Access

```bash
# Test main image URL
curl -I https://your-storage.b-cdn.net/images/filename.jpg

# Should return 200 OK
```

### 6.3 Test Delete

```bash
curl -X DELETE http://localhost:3000/api/v1/media/:id \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📁 File Organization in Bunny Storage

```
social-media-storage/
├── images/
│   ├── {userId}_{timestamp}_{random}.jpg
│   ├── {userId}_{timestamp}_{random}.png
│   └── ...
├── thumbnails/
│   ├── {userId}_{timestamp}_{random}.jpg
│   └── ...
├── videos/
│   ├── {userId}_{timestamp}_{random}.mp4
│   └── ...
└── avatars/
    ├── {userId}_avatar.jpg
    └── ...
```

---

## 💰 Pricing (as of 2025)

### Bunny Storage
- **$0.01/GB/month** - Storage
- **$0.01/GB** - Bandwidth (first 500GB free)
- No minimum fees
- No request charges

### Comparison with AWS S3
- AWS S3: ~$0.023/GB/month storage + $0.09/GB bandwidth
- **Bunny is ~70% cheaper**

---

## 🚀 Performance Tips

### 1. Enable Edge Rules (Bunny CDN)
- Auto-convert images to WebP
- Lazy loading
- Image optimization

### 2. Set Cache Headers
```go
// When uploading
req.Header.Set("Cache-Control", "public, max-age=31536000")
```

### 3. Use CDN Pull Zones
- Bunny automatically caches globally
- No additional configuration needed

### 4. Implement Cleanup Job
```go
// Delete orphaned media after 30 days
func CleanupOrphanedMedia(ctx context.Context) {
    // Find media not attached to any post/comment
    // Older than 30 days
    // Delete from Bunny and database
}
```

---

## 🔒 Security Best Practices

### 1. Validate File Types
```go
func ValidateFileType(mimeType string) error {
    allowed := []string{
        "image/jpeg", "image/png", "image/gif", "image/webp",
        "video/mp4", "video/webm", "video/quicktime",
    }

    for _, a := range allowed {
        if mimeType == a {
            return nil
        }
    }

    return errors.New("file type not allowed")
}
```

### 2. Limit File Size
```go
// Max 10MB for images, 100MB for videos
const MaxImageSize = 10 * 1024 * 1024  // 10MB
const MaxVideoSize = 100 * 1024 * 1024 // 100MB

if header.Size > MaxImageSize && isImage(mimeType) {
    return errors.New("image too large")
}
```

### 3. Rate Limiting
```go
// Limit uploads per user
// 10 uploads per minute
```

### 4. Secure Access Keys
- Never commit `.env` to git
- Use different keys for dev/staging/prod
- Rotate keys periodically

---

## 🐛 Troubleshooting

### Issue: Upload returns 401 Unauthorized
**Solution:** Check AccessKey is correct

```bash
# Test with curl
curl -X PUT "https://storage.bunnycdn.com/YOUR-ZONE/test.txt" \
  -H "AccessKey: YOUR-ACCESS-KEY" \
  -d "test content"
```

### Issue: CDN URL returns 404
**Solution:**
1. Check Pull Zone is connected to Storage Zone
2. Wait 1-2 minutes for CDN propagation
3. Verify file exists in storage zone

### Issue: Images not loading
**Solution:** Check CORS settings in Bunny dashboard

---

## ✅ Bunny Storage Integration Checklist

- [ ] Bunny.net account created
- [ ] Storage Zone created
- [ ] Pull Zone (CDN) connected
- [ ] Access credentials obtained
- [ ] Environment variables configured
- [ ] BunnyStorage service implemented
- [ ] Image processing implemented
- [ ] Upload endpoint working
- [ ] Delete endpoint working
- [ ] CDN URLs accessible
- [ ] Thumbnails generating correctly
- [ ] File validation working
- [ ] Rate limiting implemented
- [ ] Cleanup job scheduled

---

## 📚 Resources

- Bunny.net Docs: https://docs.bunny.net/
- Storage API Reference: https://docs.bunny.net/reference/storage-api
- CDN Documentation: https://docs.bunny.net/docs/cdn
- Pricing: https://bunny.net/pricing/

---

**Bunny Storage Setup Complete? → Proceed to `04-api-endpoints-checklist.md`**
