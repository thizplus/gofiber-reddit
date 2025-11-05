package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type SavedPostHandler struct {
	savedPostService services.SavedPostService
}

func NewSavedPostHandler(savedPostService services.SavedPostService) *SavedPostHandler {
	return &SavedPostHandler{
		savedPostService: savedPostService,
	}
}

// SavePost saves a post
func (h *SavedPostHandler) SavePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	postID, err := uuid.Parse(c.Params("postId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	savedPost, err := h.savedPostService.SavePost(c.Context(), userID, postID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to save post", err)
	}

	return utils.SuccessResponse(c, "Post saved successfully", savedPost)
}

// UnsavePost unsaves a post
func (h *SavedPostHandler) UnsavePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	postID, err := uuid.Parse(c.Params("postId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	err = h.savedPostService.UnsavePost(c.Context(), userID, postID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unsave post", err)
	}

	return utils.SuccessResponse(c, "Post unsaved successfully", nil)
}

// IsSaved checks if a post is saved
func (h *SavedPostHandler) IsSaved(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	postID, err := uuid.Parse(c.Params("postId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	status, err := h.savedPostService.IsSaved(c.Context(), userID, postID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check saved status", err)
	}

	return utils.SuccessResponse(c, "Saved status retrieved successfully", status)
}

// GetSavedPosts retrieves all saved posts
func (h *SavedPostHandler) GetSavedPosts(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	savedPosts, err := h.savedPostService.GetSavedPosts(c.Context(), userID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve saved posts", err)
	}

	return utils.SuccessResponse(c, "Saved posts retrieved successfully", savedPosts)
}
