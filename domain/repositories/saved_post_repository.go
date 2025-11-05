package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type SavedPostRepository interface {
	// Save/Unsave post
	SavePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error
	UnsavePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error

	// Check if saved
	IsSaved(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (bool, error)

	// Get saved posts
	GetSavedPosts(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Post, error)
	CountSavedPosts(ctx context.Context, userID uuid.UUID) (int64, error)

	// Batch check (for checking multiple posts at once)
	GetSavedStatus(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}
