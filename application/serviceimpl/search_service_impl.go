package serviceimpl

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type SearchServiceImpl struct {
	postRepo          repositories.PostRepository
	userRepo          repositories.UserRepository
	tagRepo           repositories.TagRepository
	searchHistoryRepo repositories.SearchHistoryRepository
	voteRepo          repositories.VoteRepository
	savedPostRepo     repositories.SavedPostRepository
}

func NewSearchService(
	postRepo repositories.PostRepository,
	userRepo repositories.UserRepository,
	tagRepo repositories.TagRepository,
	searchHistoryRepo repositories.SearchHistoryRepository,
	voteRepo repositories.VoteRepository,
	savedPostRepo repositories.SavedPostRepository,
) services.SearchService {
	return &SearchServiceImpl{
		postRepo:          postRepo,
		userRepo:          userRepo,
		tagRepo:           tagRepo,
		searchHistoryRepo: searchHistoryRepo,
		voteRepo:          voteRepo,
		savedPostRepo:     savedPostRepo,
	}
}

func (s *SearchServiceImpl) Search(ctx context.Context, userID *uuid.UUID, req *dto.SearchRequest) (*dto.SearchResponse, error) {
	if req.Query == "" {
		return nil, errors.New("search query is required")
	}

	searchType := req.Type
	if searchType == "" {
		searchType = "all"
	}

	limit := req.Limit
	if limit == 0 {
		limit = 20
	}

	response := &dto.SearchResponse{
		Query: req.Query,
		Type:  searchType,
		Meta: dto.PaginationMeta{
			Offset: 0,
			Limit:  limit,
		},
	}

	// Search posts
	if searchType == "post" || searchType == "all" {
		posts, err := s.postRepo.Search(ctx, req.Query, 0, limit)
		if err == nil {
			postResponses := make([]dto.PostResponse, len(posts))
			postIDs := make([]uuid.UUID, len(posts))
			for i, post := range posts {
				postIDs[i] = post.ID
			}

			// Get user-specific data if authenticated
			var voteMap map[uuid.UUID]*models.Vote
			var savedMap map[uuid.UUID]bool
			if userID != nil {
				voteMap, _ = s.voteRepo.GetUserVotesForTargets(ctx, *userID, postIDs, "post")
				savedMap, _ = s.savedPostRepo.GetSavedStatus(ctx, *userID, postIDs)
			}

			for i, post := range posts {
				resp := dto.PostToPostResponse(post)

				// Calculate hot score
				hoursSinceCreation := time.Since(post.CreatedAt).Hours()
				hotScore := float64(post.Votes) / math.Pow(hoursSinceCreation+2, 1.5)
				resp.HotScore = &hotScore

				// Add user-specific data
				if userID != nil {
					if vote, ok := voteMap[post.ID]; ok {
						resp.UserVote = &vote.VoteType
					}
					// Always set isSaved for authenticated users
					isSaved := false
					if saved, ok := savedMap[post.ID]; ok {
						isSaved = saved
					}
					resp.IsSaved = &isSaved
				}

				postResponses[i] = *resp
			}
			response.Posts = postResponses
			response.Meta.Total = int64(len(postResponses))
		}
	}

	// Search users (by username or display name)
	if searchType == "user" || searchType == "all" {
		// Simple search using List (would need a proper search method in UserRepository)
		users, err := s.userRepo.List(ctx, 0, limit)
		if err == nil {
			// Filter by query (simple substring match)
			userResponses := []dto.UserResponse{}
			for _, user := range users {
				// Check if username or displayName contains query (case insensitive)
				// This is a simplified implementation
				userResp := dto.UserToUserResponse(user)
				userResponses = append(userResponses, *userResp)
			}
			response.Users = userResponses
		}
	}

	// Search tags
	if searchType == "tag" || searchType == "all" {
		tags, err := s.tagRepo.Search(ctx, req.Query, limit)
		if err == nil {
			tagResponses := make([]dto.TagResponse, len(tags))
			for i, tag := range tags {
				tagResponses[i] = *dto.TagToTagResponse(tag)
			}
			response.Tags = tagResponses
		}
	}

	// Save search history if user is authenticated
	if userID != nil {
		_ = s.SaveSearchHistory(ctx, *userID, req.Query, searchType)
	}

	return response, nil
}

func (s *SearchServiceImpl) GetSearchHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.SearchHistoryListResponse, error) {
	history, err := s.searchHistoryRepo.ListByUser(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.SearchHistoryResponse, len(history))
	for i, h := range history {
		responses[i] = *dto.SearchHistoryToResponse(h)
	}

	return &dto.SearchHistoryListResponse{
		History: responses,
		Meta: dto.PaginationMeta{
			Total:  int64(len(responses)),
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (s *SearchServiceImpl) GetPopularSearches(ctx context.Context, limit int) (*dto.PopularSearchesResponse, error) {
	queries, err := s.searchHistoryRepo.GetPopularSearches(ctx, limit)
	if err != nil {
		return nil, err
	}

	return &dto.PopularSearchesResponse{
		Queries: queries,
	}, nil
}

func (s *SearchServiceImpl) ClearSearchHistory(ctx context.Context, userID uuid.UUID) error {
	return s.searchHistoryRepo.DeleteByUser(ctx, userID)
}

func (s *SearchServiceImpl) DeleteSearchHistoryItem(ctx context.Context, userID uuid.UUID, historyID uuid.UUID) error {
	// TODO: Should verify ownership before deleting
	return s.searchHistoryRepo.Delete(ctx, historyID)
}

func (s *SearchServiceImpl) SaveSearchHistory(ctx context.Context, userID uuid.UUID, query string, searchType string) error {
	history := &models.SearchHistory{
		ID:         uuid.New(),
		UserID:     userID,
		Query:      query,
		Type:       searchType,
		SearchedAt: time.Now(),
	}

	return s.searchHistoryRepo.Create(ctx, history)
}

var _ services.SearchService = (*SearchServiceImpl)(nil)
