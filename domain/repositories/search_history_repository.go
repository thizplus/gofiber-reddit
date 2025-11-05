package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type SearchHistoryRepository interface {
	// Create search history entry
	Create(ctx context.Context, history *models.SearchHistory) error

	// Get user's search history
	ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.SearchHistory, error)

	// Get popular searches (global)
	GetPopularSearches(ctx context.Context, limit int) ([]string, error)

	// Get recent searches by type
	ListByUserAndType(ctx context.Context, userID uuid.UUID, searchType string, limit int) ([]*models.SearchHistory, error)

	// Clear history
	DeleteByUser(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}
