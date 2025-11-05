package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type SavedPostService interface {
	// Save/Unsave posts
	SavePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (*dto.SavedPostResponse, error)
	UnsavePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error

	// Check if saved
	IsSaved(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (*dto.SaveStatusResponse, error)

	// Get saved posts
	GetSavedPosts(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.SavedPostListResponse, error)
}
