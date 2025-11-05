package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type FollowHandler struct {
	followService services.FollowService
}

func NewFollowHandler(followService services.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

// Follow follows a user
func (h *FollowHandler) Follow(c *fiber.Ctx) error {
	followerID := c.Locals("userID").(uuid.UUID)

	followingID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	follow, err := h.followService.Follow(c.Context(), followerID, followingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to follow user", err)
	}

	return utils.SuccessResponse(c, "Followed successfully", follow)
}

// Unfollow unfollows a user
func (h *FollowHandler) Unfollow(c *fiber.Ctx) error {
	followerID := c.Locals("userID").(uuid.UUID)

	followingID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	err = h.followService.Unfollow(c.Context(), followerID, followingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unfollow user", err)
	}

	return utils.SuccessResponse(c, "Unfollowed successfully", nil)
}

// IsFollowing checks if current user is following another user
func (h *FollowHandler) IsFollowing(c *fiber.Ctx) error {
	followerID := c.Locals("userID").(uuid.UUID)

	followingID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	status, err := h.followService.IsFollowing(c.Context(), followerID, followingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check follow status", err)
	}

	return utils.SuccessResponse(c, "Follow status retrieved successfully", status)
}

// GetFollowers retrieves followers of a user
func (h *FollowHandler) GetFollowers(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Get current user ID if authenticated (optional)
	var currentUserIDPtr *uuid.UUID
	if currentUserID, ok := c.Locals("userID").(uuid.UUID); ok {
		currentUserIDPtr = &currentUserID
	}

	followers, err := h.followService.GetFollowers(c.Context(), userID, offset, limit, currentUserIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve followers", err)
	}

	return utils.SuccessResponse(c, "Followers retrieved successfully", followers)
}

// GetFollowing retrieves users that a user is following
func (h *FollowHandler) GetFollowing(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Get current user ID if authenticated (optional)
	var currentUserIDPtr *uuid.UUID
	if currentUserID, ok := c.Locals("userID").(uuid.UUID); ok {
		currentUserIDPtr = &currentUserID
	}

	following, err := h.followService.GetFollowing(c.Context(), userID, offset, limit, currentUserIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve following", err)
	}

	return utils.SuccessResponse(c, "Following retrieved successfully", following)
}

// GetMutualFollows retrieves mutual follows (friends)
func (h *FollowHandler) GetMutualFollows(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	mutuals, err := h.followService.GetMutualFollows(c.Context(), userID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve mutual follows", err)
	}

	return utils.SuccessResponse(c, "Mutual follows retrieved successfully", mutuals)
}
