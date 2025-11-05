package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type SearchHandler struct {
	searchService services.SearchService
}

func NewSearchHandler(searchService services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// Search performs a search across posts, users, and tags
func (h *SearchHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return utils.ValidationErrorResponse(c, "Search query is required")
	}

	searchType := c.Query("type", "all") // all, post, user, tag
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	req := &dto.SearchRequest{
		Query: query,
		Type:  searchType,
		Limit: limit,
	}

	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	results, err := h.searchService.Search(c.Context(), userIDPtr, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Search failed", err)
	}

	return utils.SuccessResponse(c, "Search completed successfully", results)
}

// GetSearchHistory retrieves user's search history
func (h *SearchHandler) GetSearchHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	history, err := h.searchService.GetSearchHistory(c.Context(), userID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve search history", err)
	}

	return utils.SuccessResponse(c, "Search history retrieved successfully", history)
}

// GetPopularSearches retrieves popular searches
func (h *SearchHandler) GetPopularSearches(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	searches, err := h.searchService.GetPopularSearches(c.Context(), limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve popular searches", err)
	}

	return utils.SuccessResponse(c, "Popular searches retrieved successfully", searches)
}

// ClearSearchHistory clears user's search history
func (h *SearchHandler) ClearSearchHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	err := h.searchService.ClearSearchHistory(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to clear search history", err)
	}

	return utils.SuccessResponse(c, "Search history cleared successfully", nil)
}

// DeleteSearchHistoryItem deletes a search history item
func (h *SearchHandler) DeleteSearchHistoryItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	historyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid history ID")
	}

	err = h.searchService.DeleteSearchHistoryItem(c.Context(), userID, historyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete search history item", err)
	}

	return utils.SuccessResponse(c, "Search history item deleted successfully", nil)
}
