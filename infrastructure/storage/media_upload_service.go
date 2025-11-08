package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp" // WebP support
)

type MediaUploadService struct {
	bunnyStorage BunnyStorage
	bunnyStream  *BunnyStreamService
}

type UploadResult struct {
	URL              string
	Thumbnail        string
	MimeType         string
	Size             int64
	Width            int
	Height           int
	Duration         int // For videos
	// Video streaming fields (for Bunny Stream)
	VideoID          string
	HLSURL           string
	EncodingStatus   string // "pending", "processing", "completed", "failed"
	EncodingProgress int
}

type VideoMetadata struct {
	Width    int
	Height   int
	Duration int
}

func NewMediaUploadService(bunnyStorage BunnyStorage, bunnyStream *BunnyStreamService) *MediaUploadService {
	return &MediaUploadService{
		bunnyStorage: bunnyStorage,
		bunnyStream:  bunnyStream,
	}
}

// UploadImage uploads an image to Bunny Storage and generates thumbnail
func (s *MediaUploadService) UploadImage(ctx context.Context, file multipart.File, filename string) (*UploadResult, error) {
	// Read file into buffer
	var buf bytes.Buffer
	fileSize, err := io.Copy(&buf, file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Decode image to get dimensions and format
	img, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Generate unique filename
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = "." + format // Use detected format if no extension
	}
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadPath := fmt.Sprintf("chat/images/%s", uniqueName)

	// Detect content type
	contentType := fmt.Sprintf("image/%s", format)
	if format == "jpeg" {
		contentType = "image/jpeg"
	}

	// Upload original to Bunny
	originalURL, err := s.bunnyStorage.UploadFile(bytes.NewReader(buf.Bytes()), uploadPath, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload original image: %w", err)
	}

	// Generate thumbnail (max 300x300)
	thumbnail := resize.Thumbnail(300, 300, img, resize.Lanczos3)

	// Encode thumbnail
	var thumbBuf bytes.Buffer
	if err := s.encodeImage(&thumbBuf, thumbnail, format); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	// Upload thumbnail
	thumbPath := fmt.Sprintf("chat/thumbnails/%s", uniqueName)
	thumbnailURL, err := s.bunnyStorage.UploadFile(bytes.NewReader(thumbBuf.Bytes()), thumbPath, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	return &UploadResult{
		URL:       originalURL,
		Thumbnail: thumbnailURL,
		MimeType:  contentType,
		Size:      fileSize,
		Width:     width,
		Height:    height,
	}, nil
}

// encodeImage encodes image to appropriate format
func (s *MediaUploadService) encodeImage(w io.Writer, img image.Image, format string) error {
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
	case "png":
		return png.Encode(w, img)
	case "gif":
		return gif.Encode(w, img, nil)
	default:
		// Default to JPEG for unknown formats
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
	}
}

// UploadVideo uploads a video to Bunny Stream (HLS streaming)
func (s *MediaUploadService) UploadVideo(ctx context.Context, file multipart.File, filename string) (*UploadResult, error) {
	// Upload to Bunny Stream (pass file directly as it implements io.Reader)
	createResp, err := s.bunnyStream.CreateVideo(file, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to upload video to Bunny Stream: %w", err)
	}

	// Get HLS and thumbnail URLs
	hlsURL := s.bunnyStream.GetHLSURL(createResp.VideoID)
	thumbnailURL := s.bunnyStream.GetThumbnailURL(createResp.VideoID)

	return &UploadResult{
		URL:              hlsURL,              // HLS URL for playback
		Thumbnail:        thumbnailURL,        // Thumbnail URL
		MimeType:         "application/x-mpegURL", // HLS MIME type
		Size:             0,                   // Size not available immediately
		VideoID:          createResp.VideoID,
		HLSURL:           hlsURL,
		EncodingStatus:   "pending", // Initial status
		EncodingProgress: 0,
		// Width, Height, Duration will be available after encoding
	}, nil
}

// UploadFile uploads a generic file to Bunny Storage
func (s *MediaUploadService) UploadFile(ctx context.Context, file multipart.File, filename string, mimeType string) (*UploadResult, error) {
	// Read file
	var buf bytes.Buffer
	fileSize, err := io.Copy(&buf, file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(filename)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadPath := fmt.Sprintf("chat/files/%s", uniqueName)

	// Upload to Bunny
	fileURL, err := s.bunnyStorage.UploadFile(bytes.NewReader(buf.Bytes()), uploadPath, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &UploadResult{
		URL:      fileURL,
		MimeType: mimeType,
		Size:     fileSize,
	}, nil
}
