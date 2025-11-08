package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type BlockRepository interface {
	// Basic CRUD
	Create(ctx context.Context, block *models.Block) error
	Delete(ctx context.Context, blockerID, blockedID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Block, error)

	// Check block status
	IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error)
	GetBlockStatus(ctx context.Context, user1ID, user2ID uuid.UUID) (user1BlockedUser2 bool, user2BlockedUser1 bool, err error)

	// List blocks
	ListByBlocker(ctx context.Context, blockerID uuid.UUID, offset, limit int) ([]*models.Block, error)
	ListBlockedUsers(ctx context.Context, blockerID uuid.UUID, offset, limit int) ([]*models.User, error) // Returns blocked users with details

	// Stats
	Count(ctx context.Context) (int64, error)
	CountByBlocker(ctx context.Context, blockerID uuid.UUID) (int64, error)
}
