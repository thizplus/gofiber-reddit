package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type FollowRepository interface {
	// Follow/Unfollow
	Follow(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) error

	// Check relationship
	IsFollowing(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) (bool, error)

	// Get followers
	GetFollowers(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.User, error)
	CountFollowers(ctx context.Context, userID uuid.UUID) (int64, error)

	// Get following
	GetFollowing(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.User, error)
	CountFollowing(ctx context.Context, userID uuid.UUID) (int64, error)

	// Batch check (for checking multiple users at once)
	GetFollowStatus(ctx context.Context, followerID uuid.UUID, userIDs []uuid.UUID) (map[uuid.UUID]bool, error)

	// Mutual follows
	GetMutualFollows(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.User, error)

	// Update user counts (followers_count, following_count in users table)
	UpdateFollowerCount(ctx context.Context, userID uuid.UUID, delta int) error
	UpdateFollowingCount(ctx context.Context, userID uuid.UUID, delta int) error
}
