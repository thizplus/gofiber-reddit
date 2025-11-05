package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type VoteHandler struct {
	voteService services.VoteService
}

func NewVoteHandler(voteService services.VoteService) *VoteHandler {
	return &VoteHandler{
		voteService: voteService,
	}
}

// Vote creates or updates a vote (upvote/downvote)
func (h *VoteHandler) Vote(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req dto.VoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	vote, err := h.voteService.Vote(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to vote", err)
	}

	return utils.SuccessResponse(c, "Voted successfully", vote)
}

// Unvote removes a vote
func (h *VoteHandler) Unvote(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get parameters from path
	targetType := c.Params("targetType")
	if targetType != "post" && targetType != "comment" {
		return utils.ValidationErrorResponse(c, "Invalid target type. Must be 'post' or 'comment'")
	}

	targetID, err := uuid.Parse(c.Params("targetId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid target ID")
	}

	// Create request from path parameters
	req := &dto.UnvoteRequest{
		TargetID:   targetID,
		TargetType: targetType,
	}

	err = h.voteService.Unvote(c.Context(), userID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unvote", err)
	}

	return utils.SuccessResponse(c, "Unvoted successfully", nil)
}

// GetVote retrieves user's vote on a target
func (h *VoteHandler) GetVote(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	targetID, err := uuid.Parse(c.Params("targetId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid target ID")
	}

	targetType := c.Query("targetType", "post") // post or comment
	if targetType != "post" && targetType != "comment" {
		return utils.ValidationErrorResponse(c, "Invalid target type. Must be 'post' or 'comment'")
	}

	vote, err := h.voteService.GetVote(c.Context(), userID, targetID, targetType)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Vote not found", err)
	}

	if vote == nil {
		return utils.SuccessResponse(c, "No vote found", nil)
	}

	return utils.SuccessResponse(c, "Vote retrieved successfully", vote)
}

// GetVoteCount retrieves vote counts for a target
func (h *VoteHandler) GetVoteCount(c *fiber.Ctx) error {
	targetID, err := uuid.Parse(c.Params("targetId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid target ID")
	}

	targetType := c.Query("targetType", "post")
	if targetType != "post" && targetType != "comment" {
		return utils.ValidationErrorResponse(c, "Invalid target type. Must be 'post' or 'comment'")
	}

	voteCount, err := h.voteService.GetVoteCount(c.Context(), targetID, targetType)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve vote count", err)
	}

	return utils.SuccessResponse(c, "Vote count retrieved successfully", voteCount)
}

// GetUserVotes retrieves user's votes
func (h *VoteHandler) GetUserVotes(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	targetType := c.Query("targetType", "") // optional: post, comment, or empty for all
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	votes, err := h.voteService.GetUserVotes(c.Context(), userID, targetType, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve votes", err)
	}

	return utils.SuccessResponse(c, "Votes retrieved successfully", votes)
}
