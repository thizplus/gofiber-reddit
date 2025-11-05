package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type FollowRepositoryImpl struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) repositories.FollowRepository {
	return &FollowRepositoryImpl{db: db}
}

func (r *FollowRepositoryImpl) Follow(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) error {
	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return r.db.WithContext(ctx).Create(follow).Error
}

func (r *FollowRepositoryImpl) Unfollow(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&models.Follow{}).Error
}

func (r *FollowRepositoryImpl) IsFollowing(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *FollowRepositoryImpl) GetFollowers(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN follows ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userID).
		Offset(offset).Limit(limit).
		Find(&users).Error
	return users, err
}

func (r *FollowRepositoryImpl) CountFollowers(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Follow{}).
		Where("following_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *FollowRepositoryImpl) GetFollowing(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN follows ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userID).
		Offset(offset).Limit(limit).
		Find(&users).Error
	return users, err
}

func (r *FollowRepositoryImpl) CountFollowing(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *FollowRepositoryImpl) GetFollowStatus(ctx context.Context, followerID uuid.UUID, userIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	var follows []models.Follow
	err := r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id IN ?", followerID, userIDs).
		Find(&follows).Error
	if err != nil {
		return nil, err
	}

	// Convert to map
	statusMap := make(map[uuid.UUID]bool)
	for _, userID := range userIDs {
		statusMap[userID] = false
	}
	for _, follow := range follows {
		statusMap[follow.FollowingID] = true
	}

	return statusMap, nil
}

func (r *FollowRepositoryImpl) GetMutualFollows(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	// Find users who follow userID and are followed by userID
	err := r.db.WithContext(ctx).
		Joins("JOIN follows f1 ON f1.following_id = users.id").
		Joins("JOIN follows f2 ON f2.follower_id = users.id").
		Where("f1.follower_id = ? AND f2.following_id = ?", userID, userID).
		Offset(offset).Limit(limit).
		Find(&users).Error
	return users, err
}

func (r *FollowRepositoryImpl) UpdateFollowerCount(ctx context.Context, userID uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("followers_count", gorm.Expr("followers_count + ?", delta)).Error
}

func (r *FollowRepositoryImpl) UpdateFollowingCount(ctx context.Context, userID uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("following_count", gorm.Expr("following_count + ?", delta)).Error
}

var _ repositories.FollowRepository = (*FollowRepositoryImpl)(nil)
