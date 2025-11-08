package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type BlockRepositoryImpl struct {
	db *gorm.DB
}

func NewBlockRepository(db *gorm.DB) repositories.BlockRepository {
	return &BlockRepositoryImpl{db: db}
}

func (r *BlockRepositoryImpl) Create(ctx context.Context, block *models.Block) error {
	return r.db.WithContext(ctx).Create(block).Error
}

func (r *BlockRepositoryImpl) Delete(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).
		Delete(&models.Block{}).Error
}

func (r *BlockRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Block, error) {
	var block models.Block
	err := r.db.WithContext(ctx).
		Preload("Blocker").
		Preload("Blocked").
		First(&block, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *BlockRepositoryImpl) IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Block{}).
		Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).
		Count(&count).Error
	return count > 0, err
}

func (r *BlockRepositoryImpl) GetBlockStatus(ctx context.Context, user1ID, user2ID uuid.UUID) (user1BlockedUser2 bool, user2BlockedUser1 bool, err error) {
	// Check if user1 blocked user2
	user1BlockedUser2, err = r.IsBlocked(ctx, user1ID, user2ID)
	if err != nil {
		return false, false, err
	}

	// Check if user2 blocked user1
	user2BlockedUser1, err = r.IsBlocked(ctx, user2ID, user1ID)
	if err != nil {
		return false, false, err
	}

	return user1BlockedUser2, user2BlockedUser1, nil
}

func (r *BlockRepositoryImpl) ListByBlocker(ctx context.Context, blockerID uuid.UUID, offset, limit int) ([]*models.Block, error) {
	var blocks []*models.Block
	err := r.db.WithContext(ctx).
		Preload("Blocker").
		Preload("Blocked").
		Where("blocker_id = ?", blockerID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&blocks).Error
	return blocks, err
}

func (r *BlockRepositoryImpl) ListBlockedUsers(ctx context.Context, blockerID uuid.UUID, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN blocks ON blocks.blocked_id = users.id").
		Where("blocks.blocker_id = ?", blockerID).
		Order("blocks.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	return users, err
}

func (r *BlockRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Block{}).Count(&count).Error
	return count, err
}

func (r *BlockRepositoryImpl) CountByBlocker(ctx context.Context, blockerID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Block{}).
		Where("blocker_id = ?", blockerID).
		Count(&count).Error
	return count, err
}

// Ensure interface compliance
var _ repositories.BlockRepository = (*BlockRepositoryImpl)(nil)
