package services

import (
	"context"
	"gofiber-template/domain/dto"
	"mime/multipart"
	"github.com/google/uuid"
)

type MediaService interface {
	// Upload media
	UploadImage(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.MediaUploadResponse, error)
	UploadVideo(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.MediaUploadResponse, error)

	// Get media
	GetMedia(ctx context.Context, mediaID uuid.UUID) (*dto.MediaResponse, error)
	GetUserMedia(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.MediaListResponse, error)
	GetUserMediaByType(ctx context.Context, userID uuid.UUID, mediaType string, offset, limit int) (*dto.MediaListResponse, error)

	// Delete media
	DeleteMedia(ctx context.Context, mediaID uuid.UUID, userID uuid.UUID) error

	// Cleanup unused media (for background jobs)
	CleanupUnusedMedia(ctx context.Context, olderThanDays int) (int, error)
}
