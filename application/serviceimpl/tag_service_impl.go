package serviceimpl

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type TagServiceImpl struct {
	tagRepo repositories.TagRepository
}

func NewTagService(tagRepo repositories.TagRepository) services.TagService {
	return &TagServiceImpl{
		tagRepo: tagRepo,
	}
}

func (s *TagServiceImpl) GetTag(ctx context.Context, tagID uuid.UUID) (*dto.TagResponse, error) {
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, err
	}

	return dto.TagToTagResponse(tag), nil
}

func (s *TagServiceImpl) GetTagByName(ctx context.Context, name string) (*dto.TagResponse, error) {
	tag, err := s.tagRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return dto.TagToTagResponse(tag), nil
}

func (s *TagServiceImpl) ListTags(ctx context.Context, offset, limit int) (*dto.TagListResponse, error) {
	tags, err := s.tagRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	count, err := s.tagRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = *dto.TagToTagResponse(tag)
	}

	return &dto.TagListResponse{
		Tags: responses,
		Meta: dto.PaginationMeta{
			Total:  count,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (s *TagServiceImpl) GetPopularTags(ctx context.Context, limit int) (*dto.PopularTagsResponse, error) {
	tags, err := s.tagRepo.ListPopular(ctx, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = *dto.TagToTagResponse(tag)
	}

	return &dto.PopularTagsResponse{
		Tags: responses,
	}, nil
}

func (s *TagServiceImpl) SearchTags(ctx context.Context, query string, limit int) (*dto.TagListResponse, error) {
	tags, err := s.tagRepo.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = *dto.TagToTagResponse(tag)
	}

	return &dto.TagListResponse{
		Tags: responses,
		Meta: dto.PaginationMeta{
			Total:  int64(len(responses)),
			Offset: 0,
			Limit:  limit,
		},
	}, nil
}

func (s *TagServiceImpl) GetOrCreateTags(ctx context.Context, tagNames []string) ([]uuid.UUID, error) {
	tagIDs := make([]uuid.UUID, 0, len(tagNames))

	for _, name := range tagNames {
		tag, err := s.tagRepo.GetOrCreate(ctx, name)
		if err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, tag.ID)

		// Increment post count
		_ = s.tagRepo.IncrementPostCount(ctx, tag.ID)
	}

	return tagIDs, nil
}

var _ services.TagService = (*TagServiceImpl)(nil)
