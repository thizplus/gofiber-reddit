package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type TagService interface {
	// Get tags
	GetTag(ctx context.Context, tagID uuid.UUID) (*dto.TagResponse, error)
	GetTagByName(ctx context.Context, name string) (*dto.TagResponse, error)

	// List tags
	ListTags(ctx context.Context, offset, limit int) (*dto.TagListResponse, error)
	GetPopularTags(ctx context.Context, limit int) (*dto.PopularTagsResponse, error)

	// Search tags
	SearchTags(ctx context.Context, query string, limit int) (*dto.TagListResponse, error)

	// Internal methods (used by PostService)
	GetOrCreateTags(ctx context.Context, tagNames []string) ([]uuid.UUID, error)
}
