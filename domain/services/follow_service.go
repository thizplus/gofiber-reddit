package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type FollowService interface {
	// Follow/Unfollow
	Follow(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) (*dto.FollowResponse, error)
	Unfollow(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) error

	// Check relationship
	IsFollowing(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) (*dto.FollowStatusResponse, error)

	// Get followers/following
	GetFollowers(ctx context.Context, userID uuid.UUID, offset, limit int, currentUserID *uuid.UUID) (*dto.FollowerListResponse, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, offset, limit int, currentUserID *uuid.UUID) (*dto.FollowingListResponse, error)

	// Mutual follows
	GetMutualFollows(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.FollowerListResponse, error)
}
