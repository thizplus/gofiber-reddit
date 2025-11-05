package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type CommentSortBy string

const (
	CommentSortByHot CommentSortBy = "hot" // votes / (hours + 2)^1.5
	CommentSortByNew CommentSortBy = "new" // created_at DESC
	CommentSortByTop CommentSortBy = "top" // votes DESC
	CommentSortByOld CommentSortBy = "old" // created_at ASC
)

type CommentRepository interface {
	// Basic CRUD
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error)
	Update(ctx context.Context, id uuid.UUID, comment *models.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error // Soft delete

	// List & Filter
	ListByPost(ctx context.Context, postID uuid.UUID, offset, limit int, sortBy CommentSortBy) ([]*models.Comment, error)
	ListByAuthor(ctx context.Context, authorID uuid.UUID, offset, limit int) ([]*models.Comment, error)
	ListReplies(ctx context.Context, parentID uuid.UUID, offset, limit int, sortBy CommentSortBy) ([]*models.Comment, error)

	// Tree structure
	GetCommentTree(ctx context.Context, postID uuid.UUID, maxDepth int) ([]*models.Comment, error)
	GetParentChain(ctx context.Context, commentID uuid.UUID) ([]*models.Comment, error)

	// Stats
	Count(ctx context.Context) (int64, error)
	CountByPost(ctx context.Context, postID uuid.UUID) (int64, error)
	CountByAuthor(ctx context.Context, authorID uuid.UUID) (int64, error)
	CountReplies(ctx context.Context, parentID uuid.UUID) (int64, error)

	// Vote management
	UpdateVoteCount(ctx context.Context, commentID uuid.UUID, voteChange int) error
}
