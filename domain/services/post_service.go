package services

import (
	"context"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/repositories"
	"github.com/google/uuid"
)

type PostService interface {
	// Create and manage posts
	CreatePost(ctx context.Context, userID uuid.UUID, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	GetPost(ctx context.Context, postID uuid.UUID, userID *uuid.UUID) (*dto.PostResponse, error)
	UpdatePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	DeletePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error

	// List and filter posts
	ListPosts(ctx context.Context, offset, limit int, sortBy repositories.PostSortBy, userID *uuid.UUID) (*dto.PostListResponse, error)
	ListPostsByAuthor(ctx context.Context, authorID uuid.UUID, offset, limit int, userID *uuid.UUID) (*dto.PostListResponse, error)
	ListPostsByTag(ctx context.Context, tagName string, offset, limit int, sortBy repositories.PostSortBy, userID *uuid.UUID) (*dto.PostListResponse, error)
	ListPostsByTagID(ctx context.Context, tagID uuid.UUID, offset, limit int, sortBy repositories.PostSortBy, userID *uuid.UUID) (*dto.PostListResponse, error)

	// Search
	SearchPosts(ctx context.Context, query string, offset, limit int, userID *uuid.UUID) (*dto.PostListResponse, error)

	// Crossposting
	CreateCrosspost(ctx context.Context, userID uuid.UUID, sourcePostID uuid.UUID, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	GetCrossposts(ctx context.Context, postID uuid.UUID, offset, limit int, userID *uuid.UUID) (*dto.PostListResponse, error)

	// Feed
	GetFeed(ctx context.Context, userID uuid.UUID, offset, limit int, sortBy repositories.PostSortBy) (*dto.PostFeedResponse, error)
}
