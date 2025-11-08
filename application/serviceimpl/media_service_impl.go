package serviceimpl

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/redis"
	"gofiber-template/infrastructure/storage"
)

type MediaServiceImpl struct {
	mediaRepo     repositories.MediaRepository
	bunnyStorage  storage.BunnyStorage
	bunnyStream   *storage.BunnyStreamService
	redisService  *redis.RedisService
	allowedImages []string
	allowedVideos []string
	maxImageSize  int64 // bytes
	maxVideoSize  int64 // bytes
}

func NewMediaService(
	mediaRepo repositories.MediaRepository,
	bunnyStorage storage.BunnyStorage,
	bunnyStream *storage.BunnyStreamService,
	redisService *redis.RedisService,
) services.MediaService {
	return &MediaServiceImpl{
		mediaRepo:     mediaRepo,
		bunnyStorage:  bunnyStorage,
		bunnyStream:   bunnyStream,
		redisService:  redisService,
		allowedImages: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		allowedVideos: []string{".mp4", ".mov", ".avi", ".webm"},
		maxImageSize:  10 * 1024 * 1024,  // 10MB
		maxVideoSize:  300 * 1024 * 1024, // 300MB
	}
}

func (s *MediaServiceImpl) UploadImage(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.MediaUploadResponse, error) {
	// Validate file size
	if file.Size > s.maxImageSize {
		return nil, errors.New("image size exceeds maximum allowed (10MB)")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.contains(s.allowedImages, ext) {
		return nil, errors.New("invalid image format. Allowed: jpg, jpeg, png, gif, webp")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Generate unique filename
	filename := fmt.Sprintf("images/%s/%s%s", userID.String(), uuid.New().String(), ext)

	// Upload to Bunny Storage
	cdnURL, err := s.bunnyStorage.UploadFile(src, filename, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("failed to upload to Bunny Storage: %w", err)
	}

	// Create media record
	media := &models.Media{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      "image",
		FileName:  file.Filename,
		Extension: strings.TrimPrefix(ext, "."),
		MimeType:  file.Header.Get("Content-Type"),
		Size:      file.Size,
		URL:       cdnURL,
		CreatedAt: time.Now(),
	}

	// TODO: Extract image dimensions using image library
	// For now, leave Width/Height as 0

	err = s.mediaRepo.Create(ctx, media)
	if err != nil {
		// Cleanup uploaded file
		_ = s.bunnyStorage.DeleteFile(filename)
		return nil, err
	}

	return dto.MediaToMediaUploadResponse(media), nil
}

func (s *MediaServiceImpl) UploadVideo(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.MediaUploadResponse, error) {
	// Validate file size
	if file.Size > s.maxVideoSize {
		return nil, errors.New("video size exceeds maximum allowed (300MB)")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.contains(s.allowedVideos, ext) {
		return nil, errors.New("invalid video format. Allowed: mp4, mov, avi, webm")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Upload to Bunny Stream
	createResp, err := s.bunnyStream.CreateVideo(src, file.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to Bunny Stream: %w", err)
	}

	// Generate media ID
	mediaID := uuid.New()

	// Create media record with pending encoding status
	media := &models.Media{
		ID:               mediaID,
		UserID:           userID,
		Type:             "video",
		FileName:         file.Filename,
		Extension:        strings.TrimPrefix(ext, "."),
		MimeType:         file.Header.Get("Content-Type"),
		Size:             file.Size,
		URL:              s.bunnyStream.GetHLSURL(createResp.VideoID), // HLS URL (will work after encoding)
		Thumbnail:        s.bunnyStream.GetThumbnailURL(createResp.VideoID),
		VideoID:          createResp.VideoID,
		EncodingStatus:   "pending",
		EncodingProgress: 0,
		CreatedAt:        time.Now(),
	}

	// Save to database
	err = s.mediaRepo.Create(ctx, media)
	if err != nil {
		// TODO: Consider deleting from Bunny Stream on DB failure
		return nil, err
	}

	// Enqueue for encoding (async processing)
	if s.redisService != nil {
		err = s.redisService.EnqueueVideoEncoding(ctx, mediaID, createResp.VideoID)
		if err != nil {
			// Log error but don't fail the upload
			// Worker will still process based on database status
			fmt.Printf("Warning: Failed to enqueue video encoding: %v\n", err)
		}
	}

	return dto.MediaToMediaUploadResponse(media), nil
}

func (s *MediaServiceImpl) GetMedia(ctx context.Context, mediaID uuid.UUID) (*dto.MediaResponse, error) {
	media, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	return dto.MediaToMediaResponse(media), nil
}

func (s *MediaServiceImpl) GetUserMedia(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.MediaListResponse, error) {
	mediaList, err := s.mediaRepo.ListByUser(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	count, err := s.mediaRepo.CountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MediaResponse, len(mediaList))
	for i, media := range mediaList {
		responses[i] = *dto.MediaToMediaResponse(media)
	}

	return &dto.MediaListResponse{
		Media: responses,
		Meta: dto.PaginationMeta{
			Total:  count,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (s *MediaServiceImpl) GetUserMediaByType(ctx context.Context, userID uuid.UUID, mediaType string, offset, limit int) (*dto.MediaListResponse, error) {
	mediaList, err := s.mediaRepo.ListByType(ctx, userID, mediaType, offset, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MediaResponse, len(mediaList))
	for i, media := range mediaList {
		responses[i] = *dto.MediaToMediaResponse(media)
	}

	return &dto.MediaListResponse{
		Media: responses,
		Meta: dto.PaginationMeta{
			Total:  int64(len(responses)),
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (s *MediaServiceImpl) DeleteMedia(ctx context.Context, mediaID uuid.UUID, userID uuid.UUID) error {
	// Get media
	media, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return err
	}

	// Check ownership
	if media.UserID != userID {
		return errors.New("unauthorized: not media owner")
	}

	// Check if media is in use
	if media.UsageCount > 0 {
		return errors.New("cannot delete media: still in use")
	}

	// Extract path from URL
	path := strings.TrimPrefix(media.URL, s.bunnyStorage.GetFileURL(""))

	// Delete from Bunny Storage
	err = s.bunnyStorage.DeleteFile(path)
	if err != nil {
		return fmt.Errorf("failed to delete from Bunny Storage: %w", err)
	}

	// Delete from database
	return s.mediaRepo.Delete(ctx, mediaID)
}

func (s *MediaServiceImpl) CleanupUnusedMedia(ctx context.Context, olderThanDays int) (int, error) {
	// Get unused media older than specified days
	mediaList, err := s.mediaRepo.GetUnusedMedia(ctx, olderThanDays)
	if err != nil {
		return 0, err
	}

	deletedCount := 0
	for _, media := range mediaList {
		// Extract path from URL
		path := strings.TrimPrefix(media.URL, s.bunnyStorage.GetFileURL(""))

		// Delete from Bunny Storage
		err := s.bunnyStorage.DeleteFile(path)
		if err != nil {
			// Log error but continue
			continue
		}

		// Delete from database
		err = s.mediaRepo.Delete(ctx, media.ID)
		if err == nil {
			deletedCount++
		}
	}

	return deletedCount, nil
}

// Helper function
func (s *MediaServiceImpl) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

var _ services.MediaService = (*MediaServiceImpl)(nil)
