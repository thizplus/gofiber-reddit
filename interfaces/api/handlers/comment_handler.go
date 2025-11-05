package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type CommentHandler struct {
	commentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment creates a new comment or reply
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req dto.CreateCommentRequest
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

	comment, err := h.commentService.CreateComment(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create comment", err)
	}

	return utils.SuccessResponse(c, "Comment created successfully", comment)
}

// GetComment retrieves a single comment by ID
func (h *CommentHandler) GetComment(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	comment, err := h.commentService.GetComment(c.Context(), commentID, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Comment not found", err)
	}

	return utils.SuccessResponse(c, "Comment retrieved successfully", comment)
}

// UpdateComment updates an existing comment
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	var req dto.UpdateCommentRequest
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

	comment, err := h.commentService.UpdateComment(c.Context(), commentID, userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update comment", err)
	}

	return utils.SuccessResponse(c, "Comment updated successfully", comment)
}

// DeleteComment deletes a comment
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	err = h.commentService.DeleteComment(c.Context(), commentID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete comment", err)
	}

	return utils.SuccessResponse(c, "Comment deleted successfully", nil)
}

// ListCommentsByPost retrieves comments for a specific post
func (h *CommentHandler) ListCommentsByPost(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("postId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	sortBy := c.Query("sort", "new") // hot, new, top, old

	var sortByEnum repositories.CommentSortBy
	switch sortBy {
	case "hot":
		sortByEnum = repositories.CommentSortByHot
	case "new":
		sortByEnum = repositories.CommentSortByNew
	case "top":
		sortByEnum = repositories.CommentSortByTop
	case "old":
		sortByEnum = repositories.CommentSortByOld
	default:
		sortByEnum = repositories.CommentSortByNew
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	comments, err := h.commentService.ListCommentsByPost(c.Context(), postID, offset, limit, sortByEnum, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve comments", err)
	}

	return utils.SuccessResponse(c, "Comments retrieved successfully", comments)
}

// ListCommentsByAuthor retrieves comments by a specific author
func (h *CommentHandler) ListCommentsByAuthor(c *fiber.Ctx) error {
	authorID, err := uuid.Parse(c.Params("authorId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid author ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	comments, err := h.commentService.ListCommentsByAuthor(c.Context(), authorID, offset, limit, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve comments", err)
	}

	return utils.SuccessResponse(c, "Comments retrieved successfully", comments)
}

// ListReplies retrieves replies to a specific comment
func (h *CommentHandler) ListReplies(c *fiber.Ctx) error {
	parentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	sortBy := c.Query("sort", "new")

	var sortByEnum repositories.CommentSortBy
	switch sortBy {
	case "hot":
		sortByEnum = repositories.CommentSortByHot
	case "new":
		sortByEnum = repositories.CommentSortByNew
	case "top":
		sortByEnum = repositories.CommentSortByTop
	case "old":
		sortByEnum = repositories.CommentSortByOld
	default:
		sortByEnum = repositories.CommentSortByNew
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	comments, err := h.commentService.ListReplies(c.Context(), parentID, offset, limit, sortByEnum, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve replies", err)
	}

	return utils.SuccessResponse(c, "Replies retrieved successfully", comments)
}

// GetCommentTree retrieves nested comment tree for a post
func (h *CommentHandler) GetCommentTree(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("postId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	maxDepth, _ := strconv.Atoi(c.Query("maxDepth", "10"))
	if maxDepth > 10 {
		maxDepth = 10
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	tree, err := h.commentService.GetCommentTree(c.Context(), postID, maxDepth, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve comment tree", err)
	}

	return utils.SuccessResponse(c, "Comment tree retrieved successfully", tree)
}

// GetParentChain retrieves parent chain (breadcrumb) for a comment
func (h *CommentHandler) GetParentChain(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	chain, err := h.commentService.GetParentChain(c.Context(), commentID, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve parent chain", err)
	}

	return utils.SuccessResponse(c, "Parent chain retrieved successfully", chain)
}
