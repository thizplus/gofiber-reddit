package services

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type FileUploadService interface {
	UploadFile(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*models.Media, error)
}
