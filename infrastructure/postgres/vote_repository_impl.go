package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type VoteRepositoryImpl struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) repositories.VoteRepository {
	return &VoteRepositoryImpl{db: db}
}

func (r *VoteRepositoryImpl) Vote(ctx context.Context, vote *models.Vote) error {
	// Upsert: Insert or update if exists
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "target_id"}, {Name: "target_type"}},
			DoUpdates: clause.AssignmentColumns([]string{"vote_type"}),
		}).
		Create(vote).Error
}

func (r *VoteRepositoryImpl) Unvote(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, targetType string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Delete(&models.Vote{}).Error
}

func (r *VoteRepositoryImpl) GetVote(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, targetType string) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		First(&vote).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepositoryImpl) HasVoted(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, targetType string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Vote{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Count(&count).Error
	return count > 0, err
}

func (r *VoteRepositoryImpl) GetVoteCount(ctx context.Context, targetID uuid.UUID, targetType string) (upvotes int64, downvotes int64, err error) {
	// Count upvotes
	err = r.db.WithContext(ctx).
		Model(&models.Vote{}).
		Where("target_id = ? AND target_type = ? AND vote_type = ?", targetID, targetType, "up").
		Count(&upvotes).Error
	if err != nil {
		return 0, 0, err
	}

	// Count downvotes
	err = r.db.WithContext(ctx).
		Model(&models.Vote{}).
		Where("target_id = ? AND target_type = ? AND vote_type = ?", targetID, targetType, "down").
		Count(&downvotes).Error
	if err != nil {
		return 0, 0, err
	}

	return upvotes, downvotes, nil
}

func (r *VoteRepositoryImpl) ListByUser(ctx context.Context, userID uuid.UUID, targetType string, offset, limit int) ([]*models.Vote, error) {
	var votes []*models.Vote
	query := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID)

	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&votes).Error
	return votes, err
}

func (r *VoteRepositoryImpl) ListByTarget(ctx context.Context, targetID uuid.UUID, targetType string, offset, limit int) ([]*models.Vote, error) {
	var votes []*models.Vote
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("target_id = ? AND target_type = ?", targetID, targetType).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&votes).Error
	return votes, err
}

func (r *VoteRepositoryImpl) GetUserVotesForTargets(ctx context.Context, userID uuid.UUID, targetIDs []uuid.UUID, targetType string) (map[uuid.UUID]*models.Vote, error) {
	var votes []*models.Vote
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND target_id IN ? AND target_type = ?", userID, targetIDs, targetType).
		Find(&votes).Error
	if err != nil {
		return nil, err
	}

	// Convert to map
	voteMap := make(map[uuid.UUID]*models.Vote)
	for _, vote := range votes {
		voteMap[vote.TargetID] = vote
	}

	return voteMap, nil
}

var _ repositories.VoteRepository = (*VoteRepositoryImpl)(nil)
