package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type TagHandler struct {
	tagService services.TagService
}

func NewTagHandler(tagService services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// GetTag retrieves a tag by ID
func (h *TagHandler) GetTag(c *fiber.Ctx) error {
	tagID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tag ID")
	}

	tag, err := h.tagService.GetTag(c.Context(), tagID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Tag not found", err)
	}

	return utils.SuccessResponse(c, "Tag retrieved successfully", tag)
}

// GetTagByName retrieves a tag by name
func (h *TagHandler) GetTagByName(c *fiber.Ctx) error {
	tagName := c.Params("name")
	if tagName == "" {
		return utils.ValidationErrorResponse(c, "Tag name is required")
	}

	tag, err := h.tagService.GetTagByName(c.Context(), tagName)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Tag not found", err)
	}

	return utils.SuccessResponse(c, "Tag retrieved successfully", tag)
}

// ListTags retrieves all tags with pagination
func (h *TagHandler) ListTags(c *fiber.Ctx) error {
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	tags, err := h.tagService.ListTags(c.Context(), offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve tags", err)
	}

	return utils.SuccessResponse(c, "Tags retrieved successfully", tags)
}

// GetPopularTags retrieves popular tags
func (h *TagHandler) GetPopularTags(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	tags, err := h.tagService.GetPopularTags(c.Context(), limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve popular tags", err)
	}

	return utils.SuccessResponse(c, "Popular tags retrieved successfully", tags)
}

// SearchTags searches for tags
func (h *TagHandler) SearchTags(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return utils.ValidationErrorResponse(c, "Search query is required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	tags, err := h.tagService.SearchTags(c.Context(), query, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search tags", err)
	}

	return utils.SuccessResponse(c, "Tags search results retrieved successfully", tags)
}
