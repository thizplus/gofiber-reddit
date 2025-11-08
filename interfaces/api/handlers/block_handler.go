package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type BlockHandler struct {
	blockService services.BlockService
}

func NewBlockHandler(blockService services.BlockService) *BlockHandler {
	return &BlockHandler{
		blockService: blockService,
	}
}

// BlockUser blocks a user
// POST /blocks
func (h *BlockHandler) BlockUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req dto.BlockUserRequest
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

	if err := h.blockService.BlockUser(c.Context(), userID, req.Username); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to block user", err)
	}

	return utils.SuccessResponse(c, "User blocked successfully", nil)
}

// UnblockUser unblocks a user
// DELETE /blocks/:username
func (h *BlockHandler) UnblockUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	username := c.Params("username")
	if username == "" {
		return utils.ValidationErrorResponse(c, "Username is required")
	}

	if err := h.blockService.UnblockUser(c.Context(), userID, username); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unblock user", err)
	}

	return utils.SuccessResponse(c, "User unblocked successfully", nil)
}

// GetBlockStatus checks block status with another user
// GET /blocks/status/:username
func (h *BlockHandler) GetBlockStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	username := c.Params("username")
	if username == "" {
		return utils.ValidationErrorResponse(c, "Username is required")
	}

	status, err := h.blockService.GetBlockStatus(c.Context(), userID, username)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve block status", err)
	}

	return utils.SuccessResponse(c, "Block status retrieved successfully", status)
}

// ListBlockedUsers retrieves all blocked users
// GET /blocks?offset=0&limit=20
func (h *BlockHandler) ListBlockedUsers(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get offset from query params (default: 0)
	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get limit from query params (default: 20, max: 100)
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}
	}

	blockedUsers, err := h.blockService.ListBlockedUsers(c.Context(), userID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve blocked users", err)
	}

	return utils.SuccessResponse(c, "Blocked users retrieved successfully", blockedUsers)
}
