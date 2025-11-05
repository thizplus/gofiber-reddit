package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type SearchHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewSearchHistoryRepository(db *gorm.DB) repositories.SearchHistoryRepository {
	return &SearchHistoryRepositoryImpl{db: db}
}

func (r *SearchHistoryRepositoryImpl) Create(ctx context.Context, history *models.SearchHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *SearchHistoryRepositoryImpl) ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.SearchHistory, error) {
	var history []*models.SearchHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("searched_at DESC").
		Offset(offset).Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *SearchHistoryRepositoryImpl) GetPopularSearches(ctx context.Context, limit int) ([]string, error) {
	var results []struct {
		Query string
		Count int64
	}

	err := r.db.WithContext(ctx).
		Model(&models.SearchHistory{}).
		Select("query, COUNT(*) as count").
		Group("query").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	queries := make([]string, len(results))
	for i, result := range results {
		queries[i] = result.Query
	}

	return queries, nil
}

func (r *SearchHistoryRepositoryImpl) ListByUserAndType(ctx context.Context, userID uuid.UUID, searchType string, limit int) ([]*models.SearchHistory, error) {
	var history []*models.SearchHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, searchType).
		Order("searched_at DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *SearchHistoryRepositoryImpl) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.SearchHistory{}).Error
}

func (r *SearchHistoryRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.SearchHistory{}).Error
}

var _ repositories.SearchHistoryRepository = (*SearchHistoryRepositoryImpl)(nil)
