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
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/storage"
)

type FileUploadServiceImpl struct {
	mediaRepo        repositories.MediaRepository
	mediaUploadSvc   *storage.MediaUploadService
}

func NewFileUploadService(
	mediaRepo repositories.MediaRepository,
	mediaUploadSvc *storage.MediaUploadService,
) services.FileUploadService {
	return &FileUploadServiceImpl{
		mediaRepo:      mediaRepo,
		mediaUploadSvc: mediaUploadSvc,
	}
}

// File upload constraints
const (
	MaxFileSize = 50 * 1024 * 1024 // 50 MB
)

var AllowedMimeTypes = map[string]bool{
	// Documents
	"application/pdf":                                                        true,
	"application/msword":                                                     true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,

	// Spreadsheets
	"application/vnd.ms-excel":                                              true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":     true,

	// Presentations
	"application/vnd.ms-powerpoint":                                              true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation":  true,

	// Archives
	"application/zip":         true,
	"application/x-zip-compressed": true,
	"application/x-rar-compressed": true,

	// Text
	"text/plain": true,
}

func (s *FileUploadServiceImpl) UploadFile(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*models.Media, error) {
	// Validate file size
	if header.Size > MaxFileSize {
		return nil, errors.New("file size exceeds 50MB limit")
	}

	// Validate MIME type
	mimeType := header.Header.Get("Content-Type")
	if !AllowedMimeTypes[mimeType] {
		return nil, fmt.Errorf("file type not allowed: %s", mimeType)
	}

	// Extract filename and extension
	filename := header.Filename
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")

	// Upload to Bunny Storage
	result, err := s.mediaUploadSvc.UploadFile(ctx, file, filename, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Create media record
	now := time.Now()
	media := &models.Media{
		UserID:    userID,
		Type:      "file",
		FileName:  filename,
		Extension: ext,
		MimeType:  mimeType,
		Size:      header.Size,
		URL:       result.URL,
		CreatedAt: now,
	}

	// Save to database
	if err := s.mediaRepo.Create(ctx, media); err != nil {
		// TODO: Consider deleting from Bunny Storage on DB failure
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return media, nil
}

// Ensure interface compliance
var _ services.FileUploadService = (*FileUploadServiceImpl)(nil)
