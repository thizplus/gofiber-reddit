package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type TagRepository interface {
	// Create or get existing tag
	Create(ctx context.Context, tag *models.Tag) error
	GetOrCreate(ctx context.Context, name string) (*models.Tag, error)

	// Get tags
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error)
	GetByName(ctx context.Context, name string) (*models.Tag, error)
	GetByNames(ctx context.Context, names []string) ([]*models.Tag, error)

	// List tags
	List(ctx context.Context, offset, limit int) ([]*models.Tag, error)
	ListPopular(ctx context.Context, limit int) ([]*models.Tag, error)
	Search(ctx context.Context, query string, limit int) ([]*models.Tag, error)

	// Count
	Count(ctx context.Context) (int64, error)

	// Update post count
	IncrementPostCount(ctx context.Context, tagID uuid.UUID) error
	DecrementPostCount(ctx context.Context, tagID uuid.UUID) error
	UpdatePostCount(ctx context.Context, tagID uuid.UUID, delta int) error

	// Delete
	Delete(ctx context.Context, id uuid.UUID) error
}
