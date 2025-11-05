package serviceimpl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type FollowServiceImpl struct {
	followRepo   repositories.FollowRepository
	userRepo     repositories.UserRepository
	notifService services.NotificationService
}

func NewFollowService(
	followRepo repositories.FollowRepository,
	userRepo repositories.UserRepository,
	notifService services.NotificationService,
) services.FollowService {
	return &FollowServiceImpl{
		followRepo:   followRepo,
		userRepo:     userRepo,
		notifService: notifService,
	}
}

func (s *FollowServiceImpl) Follow(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) (*dto.FollowResponse, error) {
	// Can't follow yourself
	if followerID == followingID {
		return nil, errors.New("cannot follow yourself")
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if already following
	isFollowing, _ := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if isFollowing {
		return nil, errors.New("already following")
	}

	// Create follow relationship
	err = s.followRepo.Follow(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	// Update follower/following counts
	_ = s.followRepo.UpdateFollowerCount(ctx, followingID, 1)
	_ = s.followRepo.UpdateFollowingCount(ctx, followerID, 1)

	// Send notification
	_ = s.notifService.CreateNotification(
		ctx,
		followingID,
		followerID,
		"follow",
		"เริ่มติดตามคุณ",
		nil,
		nil,
	)

	return &dto.FollowResponse{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *FollowServiceImpl) Unfollow(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) error {
	// Check if following
	isFollowing, _ := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if !isFollowing {
		return errors.New("not following")
	}

	// Remove follow relationship
	err := s.followRepo.Unfollow(ctx, followerID, followingID)
	if err != nil {
		return err
	}

	// Update follower/following counts
	_ = s.followRepo.UpdateFollowerCount(ctx, followingID, -1)
	_ = s.followRepo.UpdateFollowingCount(ctx, followerID, -1)

	return nil
}

func (s *FollowServiceImpl) IsFollowing(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) (*dto.FollowStatusResponse, error) {
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	// Check if mutual (both follow each other)
	isMutual := false
	if isFollowing {
		isMutual, _ = s.followRepo.IsFollowing(ctx, followingID, followerID)
	}

	return &dto.FollowStatusResponse{
		IsFollowing: isFollowing,
		IsMutual:    isMutual,
	}, nil
}

func (s *FollowServiceImpl) GetFollowers(ctx context.Context, userID uuid.UUID, offset, limit int, currentUserID *uuid.UUID) (*dto.FollowerListResponse, error) {
	users, err := s.followRepo.GetFollowers(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	count, err := s.followRepo.CountFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		resp := dto.UserToUserResponse(user)

		// Add isFollowing status if current user is authenticated
		if currentUserID != nil {
			isFollowing, _ := s.followRepo.IsFollowing(ctx, *currentUserID, user.ID)
			resp.IsFollowing = &isFollowing
		}

		userResponses[i] = *resp
	}

	return &dto.FollowerListResponse{
		Users: userResponses,
		Meta: dto.PaginationMeta{
			Total:  count,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (s *FollowServiceImpl) GetFollowing(ctx context.Context, userID uuid.UUID, offset, limit int, currentUserID *uuid.UUID) (*dto.FollowingListResponse, error) {
	users, err := s.followRepo.GetFollowing(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	count, err := s.followRepo.CountFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		resp := dto.UserToUserResponse(user)

		// Add isFollowing status if current user is authenticated
		if currentUserID != nil {
			isFollowing, _ := s.followRepo.IsFollowing(ctx, *currentUserID, user.ID)
			resp.IsFollowing = &isFollowing
		}

		userResponses[i] = *resp
	}

	return &dto.FollowingListResponse{
		Users: userResponses,
		Meta: dto.PaginationMeta{
			Total:  count,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (s *FollowServiceImpl) GetMutualFollows(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.FollowerListResponse, error) {
	users, err := s.followRepo.GetMutualFollows(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		resp := dto.UserToUserResponse(user)
		isMutual := true
		resp.IsFollowing = &isMutual // Always true for mutual follows
		userResponses[i] = *resp
	}

	return &dto.FollowerListResponse{
		Users: userResponses,
		Meta: dto.PaginationMeta{
			Total:  int64(len(users)),
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

var _ services.FollowService = (*FollowServiceImpl)(nil)
