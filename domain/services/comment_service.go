package services

import (
	"context"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/repositories"
	"github.com/google/uuid"
)

type CommentService interface {
	// Create and manage comments
	CreateComment(ctx context.Context, userID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetComment(ctx context.Context, commentID uuid.UUID, userID *uuid.UUID) (*dto.CommentResponse, error)
	UpdateComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	DeleteComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error

	// List comments
	ListCommentsByPost(ctx context.Context, postID uuid.UUID, offset, limit int, sortBy repositories.CommentSortBy, userID *uuid.UUID) (*dto.CommentListResponse, error)
	ListCommentsByAuthor(ctx context.Context, authorID uuid.UUID, offset, limit int, userID *uuid.UUID) (*dto.CommentListResponse, error)
	ListReplies(ctx context.Context, parentID uuid.UUID, offset, limit int, sortBy repositories.CommentSortBy, userID *uuid.UUID) (*dto.CommentListResponse, error)

	// Tree structure
	GetCommentTree(ctx context.Context, postID uuid.UUID, maxDepth int, userID *uuid.UUID) (*dto.CommentTreeResponse, error)
	GetParentChain(ctx context.Context, commentID uuid.UUID, userID *uuid.UUID) ([]*dto.CommentResponse, error)
}
