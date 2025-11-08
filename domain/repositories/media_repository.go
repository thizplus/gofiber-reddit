package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type MediaRepository interface {
	// Create media record
	Create(ctx context.Context, media *models.Media) error

	// Get media
	GetByID(ctx context.Context, id uuid.UUID) (*models.Media, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Media, error)

	// List media
	ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Media, error)
	ListByType(ctx context.Context, userID uuid.UUID, mediaType string, offset, limit int) ([]*models.Media, error)

	// Count
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	// Update
	Update(ctx context.Context, media *models.Media) error

	// Update usage count
	IncrementUsageCount(ctx context.Context, mediaID uuid.UUID) error
	DecrementUsageCount(ctx context.Context, mediaID uuid.UUID) error

	// Delete
	Delete(ctx context.Context, id uuid.UUID) error

	// Get unused media (for cleanup jobs)
	GetUnusedMedia(ctx context.Context, olderThan int) ([]*models.Media, error) // olderThan in days
}
