# Video Streaming Implementation Plan
## Bunny Stream + HLS Adaptive Bitrate Streaming

**Created**: 2025-01-07
**Purpose**: แก้ปัญหา video loading performance ด้วย HLS streaming
**Target**: Post feature & Chat feature

---

## 🎯 ปัญหาที่ต้องแก้

### ปัญหาปัจจุบัน:
```typescript
// ❌ วิธีเก่า: โหลด video source เต็มๆ
<video src="https://cdn.bunny.net/videos/my-video.mp4" controls />
```

**ผลกระทบ:**
- ✗ Video 100MB → User ต้องโหลดหมดก่อนดูได้
- ✗ Mobile ใช้ bandwidth เยอะ (ค่าเน็ตแพง)
- ✗ ไม่สามารถปรับ quality ตาม network ได้
- ✗ Loading นาน UX แย่
- ✗ Seek/Skip ต้องโหลดใหม่
- ✗ Server bandwidth cost สูง

---

## ✅ วิธีแก้: HLS Adaptive Bitrate Streaming

### เทคโนโลยีที่ใช้:
- **Bunny Stream** - Video encoding & streaming service
- **HLS (HTTP Live Streaming)** - Apple's streaming protocol
- **hls.js** - JavaScript library for HLS playback

### การทำงาน:
```
1. Upload video → Bunny Stream API
   ↓
2. Bunny auto-transcode → Multiple qualities:
   - 1080p (Full HD)
   - 720p (HD)
   - 480p (SD)
   - 360p (Mobile)
   ↓
3. Bunny generates HLS playlist (.m3u8)
   ↓
4. Frontend plays via HLS player:
   → โหลดทีละส่วน (chunks ~10 วินาที)
   → ปรับ quality อัตโนมัติตาม network speed
   → Seek/Skip ทันทีไม่ต้องโหลดใหม่
```

---

## 📊 Performance Comparison

| Metric | ก่อน (Direct MP4) | หลัง (HLS Streaming) | Improvement |
|--------|-------------------|---------------------|-------------|
| **Initial Load** | 100MB (ทั้งหมด) | 2-5MB (chunks แรก) | 🔽 95% |
| **Time to Play** | 10-30 วินาที | 1-3 วินาที | 🔽 90% |
| **Bandwidth Usage** | 100MB เสมอ | 20-100MB (adaptive) | 🔽 30-80% |
| **Quality Switching** | ❌ Fixed | ✅ Auto-adaptive | ✅ |
| **Mobile Experience** | ❌ หนักมาก | ✅ เบามาก | ✅ |
| **Seek Performance** | ❌ ต้องโหลดใหม่ | ✅ Instant | ✅ |
| **Cost (Bandwidth)** | 100% | 30-50% | 🔽 50-70% |

---

## 🏗️ Architecture Overview

```
┌─────────────────┐
│  Frontend       │
│  (React/Next)   │
└────────┬────────┘
         │ 1. Upload video file (multipart/form-data)
         ↓
┌─────────────────┐
│  Backend API    │
│  (Golang/Fiber) │
└────────┬────────┘
         │ 2. Upload to Bunny Stream API
         ↓
┌─────────────────────────────────────┐
│  Bunny Stream Service               │
│  ┌──────────────────────────────┐   │
│  │ 1. Store original video      │   │
│  │ 2. Transcode to multiple     │   │
│  │    qualities (1080p-360p)    │   │
│  │ 3. Generate HLS playlist     │   │
│  │ 4. Generate thumbnails       │   │
│  │ 5. Return URLs               │   │
│  └──────────────────────────────┘   │
└────────┬────────────────────────────┘
         │ 3. Return video metadata
         ↓
┌─────────────────┐
│  Database       │
│  (PostgreSQL)   │
│  ┌────────────┐ │
│  │ posts      │ │  ← Store hlsUrl, thumbnail
│  │ messages   │ │  ← Store hlsUrl, thumbnail
│  └────────────┘ │
└─────────────────┘
         │
         │ 4. Fetch video URL
         ↓
┌─────────────────┐
│  Frontend       │
│  (Video Player) │
│  ┌────────────┐ │
│  │ hls.js     │ │  ← Play HLS stream
│  └────────────┘ │
└─────────────────┘
```

---

## 🎬 Bunny Stream API Integration

### 1. ✅ Bunny Stream Account (Ready!)

**Your Bunny Stream Configuration:**
```
Video Library ID: 533535
CDN Hostname: vz-b1631ae0-4c8.b-cdn.net
Pull Zone: vz-b1631ae0-4c8
API Key: 4c1ec80d-8809-4852-89f647b0430d-c9ad-4914
```

**Environment Variables:**
```bash
# backend/.env
# ⚠️ IMPORTANT: ห้าม commit .env file ลง git!
BUNNY_STREAM_API_KEY=4c1ec80d-8809-4852-89f647b0430d-c9ad-4914
BUNNY_STREAM_LIBRARY_ID=533535
BUNNY_STREAM_CDN_HOSTNAME=vz-b1631ae0-4c8.b-cdn.net
BUNNY_STREAM_BASE_URL=https://video.bunnycdn.com
```

**Security Note:**
- ⚠️ API Key เป็นข้อมูลสำคัญ ห้าม commit ลง git
- ⚠️ เก็บไว้ใน `.env` file เท่านั้น
- ⚠️ เพิ่ม `.env` ใน `.gitignore`

---

### 2. Backend Implementation

#### 2.1 Install Dependencies

```bash
go get github.com/gofiber/fiber/v2
go get github.com/google/uuid
```

#### 2.2 Create Bunny Stream Service

**File**: `backend/services/bunny_stream.go`

```go
package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type BunnyStreamService struct {
    APIKey    string
    LibraryID string
    BaseURL   string
}

// Video upload response
type VideoUploadResponse struct {
    VideoID   string  `json:"videoId"`
    HLSUrl    string  `json:"hlsUrl"`
    Thumbnail string  `json:"thumbnail"`
    Duration  float64 `json:"duration"`
    Width     int     `json:"width"`
    Height    int     `json:"height"`
    Size      int64   `json:"size"`
    Status    string  `json:"status"`
}

// Bunny API video info
type BunnyVideoInfo struct {
    GUID              string    `json:"guid"`
    Title             string    `json:"title"`
    Status            int       `json:"status"` // 0=created, 1=uploading, 2=processing, 3=encoding, 4=ready
    PlayURL           string    `json:"playUrl"` // HLS playlist URL
    ThumbnailURL      string    `json:"thumbnailUrl"`
    Length            float64   `json:"length"` // duration in seconds
    Width             int       `json:"width"`
    Height            int       `json:"height"`
    AvailableResolutions string `json:"availableResolutions"` // "240,360,480,720,1080"
}

// NewBunnyStreamService creates new instance
func NewBunnyStreamService() *BunnyStreamService {
    return &BunnyStreamService{
        APIKey:    os.Getenv("BUNNY_STREAM_API_KEY"),
        LibraryID: os.Getenv("BUNNY_STREAM_LIBRARY_ID"),
        BaseURL:   "https://video.bunnycdn.com",
    }
}

// UploadVideo uploads video to Bunny Stream
func (s *BunnyStreamService) UploadVideo(file io.Reader, filename string, fileSize int64) (*VideoUploadResponse, error) {
    // Step 1: Create video object
    videoID, err := s.createVideoObject(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to create video object: %w", err)
    }

    // Step 2: Upload video file
    if err := s.uploadVideoFile(videoID, file, fileSize); err != nil {
        return nil, fmt.Errorf("failed to upload video file: %w", err)
    }

    // Step 3: Wait for encoding to complete
    videoInfo, err := s.waitForEncoding(videoID, 10*time.Minute)
    if err != nil {
        return nil, fmt.Errorf("encoding failed or timed out: %w", err)
    }

    // Step 4: Return video metadata
    return &VideoUploadResponse{
        VideoID:   videoInfo.GUID,
        HLSUrl:    videoInfo.PlayURL,
        Thumbnail: videoInfo.ThumbnailURL,
        Duration:  videoInfo.Length,
        Width:     videoInfo.Width,
        Height:    videoInfo.Height,
        Size:      fileSize,
        Status:    "ready",
    }, nil
}

// Step 1: Create video object in Bunny Stream
func (s *BunnyStreamService) createVideoObject(title string) (string, error) {
    url := fmt.Sprintf("%s/library/%s/videos", s.BaseURL, s.LibraryID)

    payload := map[string]interface{}{
        "title": title,
    }
    jsonData, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }

    req.Header.Set("AccessKey", s.APIKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("bunny API error: %s - %s", resp.Status, string(body))
    }

    var result BunnyVideoInfo
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    return result.GUID, nil
}

// Step 2: Upload video file to Bunny Stream
func (s *BunnyStreamService) uploadVideoFile(videoID string, file io.Reader, fileSize int64) error {
    url := fmt.Sprintf("%s/library/%s/videos/%s", s.BaseURL, s.LibraryID, videoID)

    req, err := http.NewRequest("PUT", url, file)
    if err != nil {
        return err
    }

    req.Header.Set("AccessKey", s.APIKey)
    req.Header.Set("Content-Type", "application/octet-stream")
    req.ContentLength = fileSize

    client := &http.Client{
        Timeout: 30 * time.Minute, // Long timeout for large files
    }

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("upload failed: %s - %s", resp.Status, string(body))
    }

    return nil
}

// Step 3: Wait for video encoding to complete
func (s *BunnyStreamService) waitForEncoding(videoID string, maxWait time.Duration) (*BunnyVideoInfo, error) {
    checkInterval := 5 * time.Second
    timeout := time.After(maxWait)
    ticker := time.NewTicker(checkInterval)
    defer ticker.Stop()

    for {
        select {
        case <-timeout:
            return nil, fmt.Errorf("encoding timeout after %v", maxWait)

        case <-ticker.C:
            videoInfo, err := s.getVideoInfo(videoID)
            if err != nil {
                // Continue checking on error
                continue
            }

            // Check encoding status
            // 0 = Created, 1 = Uploading, 2 = Processing, 3 = Encoding, 4 = Ready
            if videoInfo.Status == 4 {
                return videoInfo, nil
            }

            // Log progress
            statusText := map[int]string{
                0: "Created",
                1: "Uploading",
                2: "Processing",
                3: "Encoding",
                4: "Ready",
            }
            fmt.Printf("Video %s status: %s\n", videoID, statusText[videoInfo.Status])
        }
    }
}

// Get video info from Bunny Stream
func (s *BunnyStreamService) getVideoInfo(videoID string) (*BunnyVideoInfo, error) {
    url := fmt.Sprintf("%s/library/%s/videos/%s", s.BaseURL, s.LibraryID, videoID)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("AccessKey", s.APIKey)

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("get video info failed: %s - %s", resp.Status, string(body))
    }

    var videoInfo BunnyVideoInfo
    if err := json.NewDecoder(resp.Body).Decode(&videoInfo); err != nil {
        return nil, err
    }

    return &videoInfo, nil
}

// DeleteVideo deletes video from Bunny Stream
func (s *BunnyStreamService) DeleteVideo(videoID string) error {
    url := fmt.Sprintf("%s/library/%s/videos/%s", s.BaseURL, s.LibraryID, videoID)

    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return err
    }

    req.Header.Set("AccessKey", s.APIKey)

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("delete failed: %s - %s", resp.Status, string(body))
    }

    return nil
}
```

---

#### 2.3 Update Upload Video Endpoint

**File**: `backend/handlers/media_handler.go`

```go
package handlers

import (
    "fmt"
    "path/filepath"
    "strings"

    "github.com/gofiber/fiber/v2"
    "your-project/services"
)

type MediaHandler struct {
    bunnyStream *services.BunnyStreamService
}

func NewMediaHandler() *MediaHandler {
    return &MediaHandler{
        bunnyStream: services.NewBunnyStreamService(),
    }
}

// UploadVideo handles video upload
// POST /upload/video
func (h *MediaHandler) UploadVideo(c *fiber.Ctx) error {
    // Get uploaded file
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "No file uploaded",
            "error":   "FILE_REQUIRED",
        })
    }

    // Validate file type
    ext := strings.ToLower(filepath.Ext(file.Filename))
    allowedExts := []string{".mp4", ".mov", ".avi", ".webm", ".mkv"}
    if !contains(allowedExts, ext) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid file type. Allowed: mp4, mov, avi, webm, mkv",
            "error":   "INVALID_FILE_TYPE",
        })
    }

    // Validate file size (max 500MB)
    maxSize := int64(500 * 1024 * 1024) // 500MB
    if file.Size > maxSize {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "File too large. Maximum size is 500MB",
            "error":   "FILE_TOO_LARGE",
        })
    }

    // Open file
    fileReader, err := file.Open()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to read file",
            "error":   "FILE_READ_ERROR",
        })
    }
    defer fileReader.Close()

    // Upload to Bunny Stream
    result, err := h.bunnyStream.UploadVideo(fileReader, file.Filename, file.Size)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to upload video to streaming service",
            "error":   "UPLOAD_FAILED",
            "details": err.Error(),
        })
    }

    // Return success response
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Video uploaded and encoded successfully",
        "data": fiber.Map{
            "videoId":   result.VideoID,
            "hlsUrl":    result.HLSUrl,      // ← HLS playlist URL
            "thumbnail": result.Thumbnail,   // ← Auto-generated thumbnail
            "duration":  result.Duration,    // ← Duration in seconds
            "width":     result.Width,
            "height":    result.Height,
            "size":      result.Size,
            "mimeType":  "application/vnd.apple.mpegurl", // HLS MIME type
        },
    })
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

---

#### 2.4 Update Database Schema

**สำหรับ Posts:**

```sql
-- Add video streaming fields to posts table
ALTER TABLE posts
ADD COLUMN video_id VARCHAR(255),        -- Bunny Stream video ID
ADD COLUMN hls_url TEXT,                 -- HLS playlist URL (replaces video_url)
ADD COLUMN video_thumbnail_url TEXT,     -- Auto-generated thumbnail
ADD COLUMN video_duration FLOAT,         -- Duration in seconds
ADD COLUMN video_width INTEGER,
ADD COLUMN video_height INTEGER;

-- Index for video_id lookups
CREATE INDEX idx_posts_video_id ON posts(video_id);

-- Migrate existing data (if needed)
-- UPDATE posts SET hls_url = video_url WHERE video_url IS NOT NULL;
```

**สำหรับ Chat Messages:**

```sql
-- Update messages table for video support
-- The media column is already JSONB, just update the structure

-- Example media JSONB structure:
{
  "type": "video",
  "videoId": "abc123",
  "hlsUrl": "https://vz-b1631ae0-4c8.b-cdn.net/abc123/playlist.m3u8",
  "thumbnail": "https://vz-b1631ae0-4c8.b-cdn.net/abc123/thumbnail.jpg",
  "duration": 120.5,
  "width": 1920,
  "height": 1080,
  "size": 50000000
}
```

---

#### 2.5 Update API Response Models

**File**: `backend/models/media.go`

```go
package models

// VideoUploadResponse - Response for video upload
type VideoUploadResponse struct {
    VideoID   string  `json:"videoId"`
    HLSUrl    string  `json:"hlsUrl"`       // ← Changed from videoUrl
    Thumbnail string  `json:"thumbnail"`    // ← New field
    Duration  float64 `json:"duration"`     // ← New field
    Width     int     `json:"width"`
    Height    int     `json:"height"`
    Size      int64   `json:"size"`
    MimeType  string  `json:"mimeType"`     // application/vnd.apple.mpegurl
}

// MessageMedia - Media in chat messages
type MessageMedia struct {
    Type      string  `json:"type"`      // "video"
    VideoID   string  `json:"videoId"`   // Bunny Stream video ID
    HLSUrl    string  `json:"hlsUrl"`    // HLS playlist URL
    Thumbnail string  `json:"thumbnail"` // Thumbnail URL
    Duration  float64 `json:"duration"`  // seconds
    Width     int     `json:"width"`
    Height    int     `json:"height"`
    Size      int64   `json:"size"`
}
```

---

## 📱 Frontend Implementation

### 1. Install Dependencies

```bash
npm install hls.js
npm install --save-dev @types/hls.js
```

---

### 2. Create Video Player Component

**File**: `components/video/VideoPlayer.tsx`

```typescript
"use client";

import { useEffect, useRef, useState } from "react";
import Hls from "hls.js";
import { Play, Pause, Volume2, VolumeX, Maximize } from "lucide-react";
import { cn } from "@/lib/utils";

interface VideoPlayerProps {
  hlsUrl: string;
  thumbnail?: string;
  width?: number;
  height?: number;
  autoplay?: boolean;
  className?: string;
}

export function VideoPlayer({
  hlsUrl,
  thumbnail,
  width,
  height,
  autoplay = false,
  className,
}: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<Hls | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [isMuted, setIsMuted] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    // Safari ใช้ Native HLS
    if (video.canPlayType("application/vnd.apple.mpegurl")) {
      video.src = hlsUrl;

      if (autoplay) {
        video.play().catch((err) => {
          console.warn("Autoplay prevented:", err);
        });
      }
    }
    // Chrome, Firefox, Edge ใช้ hls.js
    else if (Hls.isSupported()) {
      const hls = new Hls({
        enableWorker: true,
        lowLatencyMode: true,
        backBufferLength: 90,
      });

      hls.loadSource(hlsUrl);
      hls.attachMedia(video);

      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        console.log("✅ HLS manifest loaded");
        if (autoplay) {
          video.play().catch((err) => {
            console.warn("Autoplay prevented:", err);
          });
        }
      });

      hls.on(Hls.Events.ERROR, (event, data) => {
        console.error("❌ HLS error:", data);

        if (data.fatal) {
          switch (data.type) {
            case Hls.ErrorTypes.NETWORK_ERROR:
              console.error("Network error, trying to recover...");
              hls.startLoad();
              break;
            case Hls.ErrorTypes.MEDIA_ERROR:
              console.error("Media error, trying to recover...");
              hls.recoverMediaError();
              break;
            default:
              setError("Failed to load video");
              hls.destroy();
              break;
          }
        }
      });

      hlsRef.current = hls;

      return () => {
        hls.destroy();
      };
    } else {
      setError("Your browser doesn't support video streaming");
    }
  }, [hlsUrl, autoplay]);

  const togglePlay = () => {
    const video = videoRef.current;
    if (!video) return;

    if (video.paused) {
      video.play();
      setIsPlaying(true);
    } else {
      video.pause();
      setIsPlaying(false);
    }
  };

  const toggleMute = () => {
    const video = videoRef.current;
    if (!video) return;

    video.muted = !video.muted;
    setIsMuted(video.muted);
  };

  const toggleFullscreen = () => {
    const video = videoRef.current;
    if (!video) return;

    if (document.fullscreenElement) {
      document.exitFullscreen();
    } else {
      video.requestFullscreen();
    }
  };

  if (error) {
    return (
      <div className="flex items-center justify-center bg-muted rounded-lg p-8">
        <p className="text-destructive">{error}</p>
      </div>
    );
  }

  return (
    <div className={cn("relative rounded-lg overflow-hidden bg-black group", className)}>
      <video
        ref={videoRef}
        className="w-full h-auto"
        poster={thumbnail}
        preload="metadata"
        playsInline
        onClick={togglePlay}
        onPlay={() => setIsPlaying(true)}
        onPause={() => setIsPlaying(false)}
      />

      {/* Custom Controls (Optional) */}
      <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-4 opacity-0 group-hover:opacity-100 transition-opacity">
        <div className="flex items-center gap-2">
          <button
            onClick={togglePlay}
            className="text-white hover:text-primary transition-colors"
          >
            {isPlaying ? <Pause className="h-6 w-6" /> : <Play className="h-6 w-6" />}
          </button>

          <button
            onClick={toggleMute}
            className="text-white hover:text-primary transition-colors"
          >
            {isMuted ? <VolumeX className="h-5 w-5" /> : <Volume2 className="h-5 w-5" />}
          </button>

          <div className="flex-1" />

          <button
            onClick={toggleFullscreen}
            className="text-white hover:text-primary transition-colors"
          >
            <Maximize className="h-5 w-5" />
          </button>
        </div>
      </div>
    </div>
  );
}
```

---

### 3. Update Chat Video Component

**File**: `components/chat/ChatMessageVideo.tsx`

```typescript
"use client";

import { useState } from "react";
import { Play } from "lucide-react";
import { VideoPlayer } from "@/components/video/VideoPlayer";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import type { MessageMedia } from "@/lib/data/mockChats";

interface ChatMessageVideoProps {
  media: MessageMedia[];
  isOwnMessage: boolean;
}

export function ChatMessageVideo({ media, isOwnMessage }: ChatMessageVideoProps) {
  const [lightboxOpen, setLightboxOpen] = useState(false);
  const video = media[0];

  const formatDuration = (seconds?: number) => {
    if (!seconds) return "";
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, "0")}`;
  };

  return (
    <>
      {/* Thumbnail Preview */}
      <div
        className="relative rounded-lg overflow-hidden cursor-pointer group max-w-xs"
        onClick={() => setLightboxOpen(true)}
      >
        <img
          src={video.thumbnail}
          alt="Video thumbnail"
          className="w-full h-auto object-cover"
          style={{ maxHeight: "300px" }}
        />

        {/* Play Button Overlay */}
        <div className="absolute inset-0 flex items-center justify-center bg-black/30 group-hover:bg-black/40 transition-colors">
          <div className="bg-white/90 rounded-full p-4">
            <Play className="h-8 w-8 text-black fill-black" />
          </div>
        </div>

        {/* Duration Badge */}
        {video.duration && (
          <div className="absolute bottom-2 right-2 bg-black/80 text-white text-xs px-2 py-1 rounded">
            {formatDuration(video.duration)}
          </div>
        )}
      </div>

      {/* Fullscreen Video Player */}
      <Dialog open={lightboxOpen} onOpenChange={setLightboxOpen}>
        <DialogContent className="max-w-4xl p-0">
          <VideoPlayer
            hlsUrl={video.url}        // HLS URL from Bunny Stream
            thumbnail={video.thumbnail}
            width={video.width}
            height={video.height}
            autoplay={true}
          />
        </DialogContent>
      </Dialog>
    </>
  );
}
```

---

## 🔧 Testing & Validation

### 1. Backend Testing

```bash
# Test video upload endpoint
curl -X POST http://localhost:8080/v1/upload/video \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test-video.mp4"

# Expected response:
{
  "success": true,
  "message": "Video uploaded and encoded successfully",
  "data": {
    "videoId": "abc123-def456-ghi789",
    "hlsUrl": "https://vz-b1631ae0-4c8.b-cdn.net/abc123/playlist.m3u8",
    "thumbnail": "https://vz-b1631ae0-4c8.b-cdn.net/abc123/thumbnail.jpg",
    "duration": 120.5,
    "width": 1920,
    "height": 1080,
    "size": 50000000,
    "mimeType": "application/vnd.apple.mpegurl"
  }
}
```

### 2. HLS Playback Testing

```bash
# Test HLS playlist
curl https://vz-b1631ae0-4c8.b-cdn.net/abc123/playlist.m3u8

# Expected output (HLS manifest):
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360
360p/index.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=1400000,RESOLUTION=842x480
480p/index.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2800000,RESOLUTION=1280x720
720p/index.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080
1080p/index.m3u8
```

### 3. Frontend Testing

**Test Cases:**
1. ✅ Upload video → ได้ HLS URL
2. ✅ Play video → โหลดเร็ว (2-5 วินาที)
3. ✅ Thumbnail แสดงก่อนเล่น
4. ✅ Quality switching automatic
5. ✅ Seek/Skip ทำงานได้
6. ✅ Fullscreen mode
7. ✅ Mobile responsive
8. ✅ Safari + Chrome compatibility

---

## 💰 Cost Estimation (Bunny Stream)

### Pricing:
- **Storage**: $0.01/GB/month
- **Streaming**: $0.01/GB delivered
- **Encoding**: FREE
- **Thumbnails**: FREE

### Example Calculation:

**Scenario**: Social media platform
- 1,000 video uploads/month
- Average size: 50MB/video
- Average views: 100 views/video
- Average watch time: 50% (25MB streamed per view)

**Monthly Cost:**
```
Storage: 1,000 videos × 50MB = 50GB × $0.01 = $0.50
Streaming: 1,000 videos × 100 views × 25MB = 2,500GB × $0.01 = $25.00
Total: $25.50/month
```

**Comparison with Direct MP4 delivery:**
```
Direct MP4: 1,000 × 100 × 50MB = 5,000GB × $0.01 = $50.00/month
Savings: $24.50/month (49% cheaper!)
```

---

## ✅ Implementation Checklist

### Backend Tasks:

- [ ] **Environment Setup**
  - [ ] สมัคร Bunny Stream account
  - [ ] สร้าง Video Library (Asia region)
  - [ ] Add BUNNY_STREAM_API_KEY, BUNNY_STREAM_LIBRARY_ID to .env

- [ ] **Code Implementation**
  - [ ] สร้าง `services/bunny_stream.go`
  - [ ] แก้ไข `handlers/media_handler.go` → UploadVideo method
  - [ ] เพิ่ม video validation (type, size)
  - [ ] Implement encoding wait logic with timeout

- [ ] **Database Changes**
  - [ ] เพิ่ม columns: video_id, hls_url, video_thumbnail_url, video_duration
  - [ ] Create indexes for video_id
  - [ ] Migrate existing video_url data (if any)

- [ ] **API Response Updates**
  - [ ] Update POST /upload/video response structure
  - [ ] Update GET /posts response (include hlsUrl, thumbnail)
  - [ ] Update POST /chat/messages response

- [ ] **Testing**
  - [ ] Test video upload endpoint
  - [ ] Test encoding wait (small + large files)
  - [ ] Test error handling (invalid file, timeout)
  - [ ] Verify HLS playlist generation
  - [ ] Test thumbnail generation

### Frontend Tasks:

- [ ] **Dependencies**
  - [ ] Install hls.js package
  - [ ] Install @types/hls.js

- [ ] **Components**
  - [ ] สร้าง VideoPlayer component
  - [ ] แก้ไข ChatMessageVideo → use VideoPlayer
  - [ ] แก้ไข PostCard → use VideoPlayer (for posts feature)
  - [ ] Add custom controls (optional)

- [ ] **API Integration**
  - [ ] Update media.service.ts → uploadVideo method
  - [ ] Update response types (UploadVideoResponse)
  - [ ] Update MessageMedia interface

- [ ] **Testing**
  - [ ] Test video upload with progress
  - [ ] Test HLS playback (Chrome, Safari, Mobile)
  - [ ] Test quality switching
  - [ ] Test fullscreen mode
  - [ ] Test thumbnail display

---

## 📚 Resources

### Documentation:
- [Bunny Stream API Docs](https://docs.bunny.net/docs/stream)
- [HLS.js Documentation](https://github.com/video-dev/hls.js/)
- [HLS Protocol Spec](https://datatracker.ietf.org/doc/html/rfc8216)

### Example Code:
- [Bunny Stream Go Examples](https://github.com/BunnyWay/bunnycdn-stream-go-sdk)
- [HLS.js React Example](https://github.com/video-dev/hls.js/blob/master/docs/API.md#react-integration)

---

## 🚨 Important Notes

1. **Encoding Time**:
   - Small video (10MB): 1-2 minutes
   - Medium video (100MB): 5-10 minutes
   - Large video (500MB): 15-30 minutes
   - Implement proper timeout handling!

2. **Mobile Compatibility**:
   - iOS Safari: Native HLS support (no hls.js needed)
   - Android Chrome: Requires hls.js
   - Test on real devices!

3. **Thumbnail Generation**:
   - Bunny auto-generates at 0s, 25%, 50%, 75%, 100%
   - Can request specific timestamp: `thumbnail.jpg?time=10`
   - Always use thumbnail for preview (faster load)

4. **Quality Selection**:
   - HLS automatically chooses best quality
   - User can manually select in video player
   - Mobile typically starts at 360p or 480p

5. **Storage Cleanup**:
   - Implement video deletion when post/message deleted
   - Call BunnyStreamService.DeleteVideo()
   - Prevent orphaned videos

---

## 📝 Next Steps

1. **Immediate (Phase 1)**:
   - Backend: Implement Bunny Stream integration
   - Frontend: Create VideoPlayer component
   - Testing: Basic upload + playback

2. **Short-term (Phase 2)**:
   - Add video upload progress bar
   - Add quality selector UI
   - Implement video deletion on post/message delete

3. **Long-term (Phase 3)**:
   - Video analytics (views, watch time)
   - Video trimming/editing in upload
   - Live streaming support

---

**Document Version**: 1.0.0
**Last Updated**: 2025-01-07
**Status**: ✅ Ready for Implementation
**Estimated Time**: 2-3 days (Backend + Frontend)

---

**End of Document**
