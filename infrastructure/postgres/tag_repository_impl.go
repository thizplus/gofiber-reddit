package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type TagRepositoryImpl struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) repositories.TagRepository {
	return &TagRepositoryImpl{db: db}
}

func (r *TagRepositoryImpl) Create(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *TagRepositoryImpl) GetOrCreate(ctx context.Context, name string) (*models.Tag, error) {
	var tag models.Tag

	// Try to get existing tag
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&tag).Error

	if err == nil {
		return &tag, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new tag
	tag = models.Tag{
		ID:        uuid.New(),
		Name:      name,
		PostCount: 0,
		CreatedAt: time.Now(),
	}

	err = r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		Create(&tag).Error

	if err != nil {
		// If conflict, try to get it again
		err = r.db.WithContext(ctx).
			Where("name = ?", name).
			First(&tag).Error
	}

	return &tag, err
}

func (r *TagRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepositoryImpl) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepositoryImpl) GetByNames(ctx context.Context, names []string) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.WithContext(ctx).
		Where("name IN ?", names).
		Find(&tags).Error
	return tags, err
}

func (r *TagRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Offset(offset).Limit(limit).
		Find(&tags).Error
	return tags, err
}

func (r *TagRepositoryImpl) ListPopular(ctx context.Context, limit int) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.WithContext(ctx).
		Order("post_count DESC").
		Limit(limit).
		Find(&tags).Error
	return tags, err
}

func (r *TagRepositoryImpl) Search(ctx context.Context, query string, limit int) ([]*models.Tag, error) {
	var tags []*models.Tag
	searchQuery := "%" + query + "%"
	err := r.db.WithContext(ctx).
		Where("name ILIKE ?", searchQuery).
		Order("post_count DESC").
		Limit(limit).
		Find(&tags).Error
	return tags, err
}

func (r *TagRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Tag{}).Count(&count).Error
	return count, err
}

func (r *TagRepositoryImpl) IncrementPostCount(ctx context.Context, tagID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Tag{}).
		Where("id = ?", tagID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
}

func (r *TagRepositoryImpl) DecrementPostCount(ctx context.Context, tagID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Tag{}).
		Where("id = ?", tagID).
		UpdateColumn("post_count", gorm.Expr("post_count - ?", 1)).Error
}

func (r *TagRepositoryImpl) UpdatePostCount(ctx context.Context, tagID uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.Tag{}).
		Where("id = ?", tagID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", delta)).Error
}

func (r *TagRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Tag{}).Error
}

var _ repositories.TagRepository = (*TagRepositoryImpl)(nil)
