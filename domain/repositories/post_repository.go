package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type PostSortBy string

const (
	SortByHot       PostSortBy = "hot"        // votes / (hours + 2)^1.5
	SortByNew       PostSortBy = "new"        // created_at DESC
	SortByTop       PostSortBy = "top"        // votes DESC
	SortByControversial PostSortBy = "controversial" // high engagement but mixed votes
)

type PostRepository interface {
	// Basic CRUD
	Create(ctx context.Context, post *models.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error)
	Update(ctx context.Context, id uuid.UUID, post *models.Post) error
	Delete(ctx context.Context, id uuid.UUID) error // Soft delete

	// List & Filter
	List(ctx context.Context, offset, limit int, sortBy PostSortBy) ([]*models.Post, error)
	ListByAuthor(ctx context.Context, authorID uuid.UUID, offset, limit int) ([]*models.Post, error)
	ListByTag(ctx context.Context, tagName string, offset, limit int, sortBy PostSortBy) ([]*models.Post, error)
	ListByTagID(ctx context.Context, tagID uuid.UUID, offset, limit int, sortBy PostSortBy) ([]*models.Post, error)

	// Search
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Post, error)

	// Crosspost
	GetCrossposts(ctx context.Context, postID uuid.UUID, offset, limit int) ([]*models.Post, error)

	// Stats
	Count(ctx context.Context) (int64, error)
	CountByAuthor(ctx context.Context, authorID uuid.UUID) (int64, error)

	// Comment count management
	IncrementCommentCount(ctx context.Context, postID uuid.UUID) error
	DecrementCommentCount(ctx context.Context, postID uuid.UUID) error

	// Vote management
	UpdateVoteCount(ctx context.Context, postID uuid.UUID, voteChange int) error

	// Media association
	AttachMedia(ctx context.Context, postID uuid.UUID, mediaIDs []uuid.UUID) error
	DetachMedia(ctx context.Context, postID uuid.UUID, mediaIDs []uuid.UUID) error

	// Tag association
	AttachTags(ctx context.Context, postID uuid.UUID, tagIDs []uuid.UUID) error
	DetachTags(ctx context.Context, postID uuid.UUID, tagIDs []uuid.UUID) error
	SyncTags(ctx context.Context, postID uuid.UUID, tagIDs []uuid.UUID) error
}
