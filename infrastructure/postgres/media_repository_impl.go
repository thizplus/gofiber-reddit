package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type MediaRepositoryImpl struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) repositories.MediaRepository {
	return &MediaRepositoryImpl{db: db}
}

func (r *MediaRepositoryImpl) Create(ctx context.Context, media *models.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

func (r *MediaRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Media, error) {
	var media models.Media
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&media).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepositoryImpl) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Media, error) {
	var mediaList []*models.Media
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id IN ?", ids).
		Find(&mediaList).Error
	return mediaList, err
}

func (r *MediaRepositoryImpl) ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Media, error) {
	var mediaList []*models.Media
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&mediaList).Error
	return mediaList, err
}

func (r *MediaRepositoryImpl) ListByType(ctx context.Context, userID uuid.UUID, mediaType string, offset, limit int) ([]*models.Media, error) {
	var mediaList []*models.Media
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, mediaType).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&mediaList).Error
	return mediaList, err
}

func (r *MediaRepositoryImpl) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Media{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *MediaRepositoryImpl) Update(ctx context.Context, media *models.Media) error {
	return r.db.WithContext(ctx).Save(media).Error
}

func (r *MediaRepositoryImpl) IncrementUsageCount(ctx context.Context, mediaID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Media{}).
		Where("id = ?", mediaID).
		UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func (r *MediaRepositoryImpl) DecrementUsageCount(ctx context.Context, mediaID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Media{}).
		Where("id = ?", mediaID).
		UpdateColumn("usage_count", gorm.Expr("usage_count - ?", 1)).Error
}

func (r *MediaRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Media{}).Error
}

func (r *MediaRepositoryImpl) GetUnusedMedia(ctx context.Context, olderThan int) ([]*models.Media, error) {
	var mediaList []*models.Media
	cutoffDate := time.Now().AddDate(0, 0, -olderThan)
	err := r.db.WithContext(ctx).
		Where("usage_count = ? AND created_at < ?", 0, cutoffDate).
		Find(&mediaList).Error
	return mediaList, err
}

var _ repositories.MediaRepository = (*MediaRepositoryImpl)(nil)
