package serviceimpl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type BlockServiceImpl struct {
	blockRepo repositories.BlockRepository
	userRepo  repositories.UserRepository
}

func NewBlockService(
	blockRepo repositories.BlockRepository,
	userRepo repositories.UserRepository,
) services.BlockService {
	return &BlockServiceImpl{
		blockRepo: blockRepo,
		userRepo:  userRepo,
	}
}

func (s *BlockServiceImpl) BlockUser(ctx context.Context, blockerID uuid.UUID, blockedUsername string) error {
	// Get blocked user by username
	blockedUser, err := s.userRepo.GetByUsername(ctx, blockedUsername)
	if err != nil {
		return errors.New("user not found")
	}

	// Can't block yourself
	if blockerID == blockedUser.ID {
		return errors.New("cannot block yourself")
	}

	// Check if already blocked
	isBlocked, err := s.blockRepo.IsBlocked(ctx, blockerID, blockedUser.ID)
	if err != nil {
		return err
	}

	if isBlocked {
		return errors.New("user is already blocked")
	}

	// Create block
	block := &models.Block{
		BlockerID: blockerID,
		BlockedID: blockedUser.ID,
		CreatedAt: time.Now(),
	}

	return s.blockRepo.Create(ctx, block)
}

func (s *BlockServiceImpl) UnblockUser(ctx context.Context, blockerID uuid.UUID, blockedUsername string) error {
	// Get blocked user by username
	blockedUser, err := s.userRepo.GetByUsername(ctx, blockedUsername)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if currently blocked
	isBlocked, err := s.blockRepo.IsBlocked(ctx, blockerID, blockedUser.ID)
	if err != nil {
		return err
	}

	if !isBlocked {
		return errors.New("user is not blocked")
	}

	// Remove block
	return s.blockRepo.Delete(ctx, blockerID, blockedUser.ID)
}

func (s *BlockServiceImpl) GetBlockStatus(ctx context.Context, userID uuid.UUID, otherUsername string) (*dto.BlockStatusResponse, error) {
	// Get other user by username
	otherUser, err := s.userRepo.GetByUsername(ctx, otherUsername)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get block status
	isBlocked, isBlockedBy, err := s.blockRepo.GetBlockStatus(ctx, userID, otherUser.ID)
	if err != nil {
		return nil, err
	}

	// Can message only if neither user is blocking the other
	canMessage := !isBlocked && !isBlockedBy

	return &dto.BlockStatusResponse{
		IsBlocked:   isBlocked,
		IsBlockedBy: isBlockedBy,
		CanMessage:  canMessage,
	}, nil
}

func (s *BlockServiceImpl) ListBlockedUsers(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.BlockedUsersResponse, error) {
	// Set default pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Get blocks
	blocks, err := s.blockRepo.ListByBlocker(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.blockRepo.CountByBlocker(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	blockedUsers := make([]dto.BlockedUserResponse, len(blocks))
	for i, block := range blocks {
		blockedUsers[i] = *dto.BlockToBlockedUserResponse(block)
	}

	return &dto.BlockedUsersResponse{
		BlockedUsers: blockedUsers,
		Meta: dto.PaginationMeta{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

// Ensure interface compliance
var _ services.BlockService = (*BlockServiceImpl)(nil)
