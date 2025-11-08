package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type BlockService interface {
	// Block user
	BlockUser(ctx context.Context, blockerID uuid.UUID, blockedUsername string) error

	// Unblock user
	UnblockUser(ctx context.Context, blockerID uuid.UUID, blockedUsername string) error

	// Check block status
	GetBlockStatus(ctx context.Context, userID uuid.UUID, otherUsername string) (*dto.BlockStatusResponse, error)

	// List blocked users
	ListBlockedUsers(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.BlockedUsersResponse, error)
}
