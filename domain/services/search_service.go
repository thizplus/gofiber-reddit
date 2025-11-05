package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type SearchService interface {
	// Search
	Search(ctx context.Context, userID *uuid.UUID, req *dto.SearchRequest) (*dto.SearchResponse, error)

	// Search history
	GetSearchHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.SearchHistoryListResponse, error)
	GetPopularSearches(ctx context.Context, limit int) (*dto.PopularSearchesResponse, error)
	ClearSearchHistory(ctx context.Context, userID uuid.UUID) error
	DeleteSearchHistoryItem(ctx context.Context, userID uuid.UUID, historyID uuid.UUID) error

	// Internal method to save search history
	SaveSearchHistory(ctx context.Context, userID uuid.UUID, query string, searchType string) error
}
