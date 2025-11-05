package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type SavedPostRepositoryImpl struct {
	db *gorm.DB
}

func NewSavedPostRepository(db *gorm.DB) repositories.SavedPostRepository {
	return &SavedPostRepositoryImpl{db: db}
}

func (r *SavedPostRepositoryImpl) SavePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error {
	savedPost := &models.SavedPost{
		UserID:  userID,
		PostID:  postID,
		SavedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(savedPost).Error
}

func (r *SavedPostRepositoryImpl) UnsavePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&models.SavedPost{}).Error
}

func (r *SavedPostRepositoryImpl) IsSaved(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.SavedPost{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}

func (r *SavedPostRepositoryImpl) GetSavedPosts(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Post, error) {
	var posts []*models.Post
	err := r.db.WithContext(ctx).
		Joins("JOIN saved_posts ON saved_posts.post_id = posts.id").
		Preload("Author").
		Preload("Media").
		Preload("Tags").
		Preload("SourcePost").
		Preload("SourcePost.Author").
		Preload("SourcePost.Media").
		Preload("SourcePost.Tags").
		Where("saved_posts.user_id = ? AND posts.is_deleted = ?", userID, false).
		Order("saved_posts.saved_at DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *SavedPostRepositoryImpl) CountSavedPosts(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.SavedPost{}).
		Joins("JOIN posts ON posts.id = saved_posts.post_id").
		Where("saved_posts.user_id = ? AND posts.is_deleted = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *SavedPostRepositoryImpl) GetSavedStatus(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	var savedPosts []models.SavedPost
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Find(&savedPosts).Error
	if err != nil {
		return nil, err
	}

	// Convert to map
	statusMap := make(map[uuid.UUID]bool)
	for _, postID := range postIDs {
		statusMap[postID] = false
	}
	for _, savedPost := range savedPosts {
		statusMap[savedPost.PostID] = true
	}

	return statusMap, nil
}

var _ repositories.SavedPostRepository = (*SavedPostRepositoryImpl)(nil)
