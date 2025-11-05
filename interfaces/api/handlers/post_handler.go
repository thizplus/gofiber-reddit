package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type PostHandler struct {
	postService services.PostService
}

func NewPostHandler(postService services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost creates a new post
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req dto.CreatePostRequest
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

	post, err := h.postService.CreatePost(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create post", err)
	}

	return utils.SuccessResponse(c, "Post created successfully", post)
}

// GetPost retrieves a single post by ID
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	post, err := h.postService.GetPost(c.Context(), postID, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Post not found", err)
	}

	return utils.SuccessResponse(c, "Post retrieved successfully", post)
}

// UpdatePost updates an existing post
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	var req dto.UpdatePostRequest
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

	post, err := h.postService.UpdatePost(c.Context(), postID, userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update post", err)
	}

	return utils.SuccessResponse(c, "Post updated successfully", post)
}

// DeletePost deletes a post
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	err = h.postService.DeletePost(c.Context(), postID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete post", err)
	}

	return utils.SuccessResponse(c, "Post deleted successfully", nil)
}

// ListPosts retrieves a list of posts with pagination and sorting
func (h *PostHandler) ListPosts(c *fiber.Ctx) error {
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	sortBy := c.Query("sort", "hot") // hot, new, top, controversial

	// Check if filtering by tag (query param) - รองรับภาษาไทย
	tagQuery := c.Query("tag")
	if tagQuery != "" {
		log.Printf("🏷️  Searching posts by tag (query param): '%s'", tagQuery)
		return h.listPostsByTagQuery(c, tagQuery, offset, limit, sortBy)
	}

	// Validate and convert sortBy
	var sortByEnum repositories.PostSortBy
	switch sortBy {
	case "hot":
		sortByEnum = repositories.SortByHot
	case "new":
		sortByEnum = repositories.SortByNew
	case "top":
		sortByEnum = repositories.SortByTop
	case "controversial":
		sortByEnum = repositories.SortByControversial
	default:
		sortByEnum = repositories.SortByHot
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	posts, err := h.postService.ListPosts(c.Context(), offset, limit, sortByEnum, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve posts", err)
	}

	return utils.SuccessResponse(c, "Posts retrieved successfully", posts)
}

// Helper function to handle tag filtering via query param
func (h *PostHandler) listPostsByTagQuery(c *fiber.Ctx, tagName string, offset, limit int, sortBy string) error {
	var sortByEnum repositories.PostSortBy
	switch sortBy {
	case "hot":
		sortByEnum = repositories.SortByHot
	case "new":
		sortByEnum = repositories.SortByNew
	case "top":
		sortByEnum = repositories.SortByTop
	default:
		sortByEnum = repositories.SortByHot
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	posts, err := h.postService.ListPostsByTag(c.Context(), tagName, offset, limit, sortByEnum, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve posts", err)
	}

	return utils.SuccessResponse(c, "Posts retrieved successfully", posts)
}

// ListPostsByAuthor retrieves posts by a specific author
func (h *PostHandler) ListPostsByAuthor(c *fiber.Ctx) error {
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

	posts, err := h.postService.ListPostsByAuthor(c.Context(), authorID, offset, limit, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve posts", err)
	}

	return utils.SuccessResponse(c, "Posts retrieved successfully", posts)
}

// ListPostsByTag retrieves posts by tag
func (h *PostHandler) ListPostsByTag(c *fiber.Ctx) error {
	tagName := c.Params("tagName")
	if tagName == "" {
		return utils.ValidationErrorResponse(c, "Tag name is required")
	}

	// URL decode (Fiber should auto-decode, but just in case)
	// Log for debugging
	log.Printf("🏷️  Searching for tag: '%s' (len: %d)", tagName, len(tagName))

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	sortBy := c.Query("sort", "hot")

	var sortByEnum repositories.PostSortBy
	switch sortBy {
	case "hot":
		sortByEnum = repositories.SortByHot
	case "new":
		sortByEnum = repositories.SortByNew
	case "top":
		sortByEnum = repositories.SortByTop
	default:
		sortByEnum = repositories.SortByHot
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	posts, err := h.postService.ListPostsByTag(c.Context(), tagName, offset, limit, sortByEnum, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve posts", err)
	}

	return utils.SuccessResponse(c, "Posts retrieved successfully", posts)
}

// ListPostsByTagID retrieves posts by tag ID
func (h *PostHandler) ListPostsByTagID(c *fiber.Ctx) error {
	tagID, err := uuid.Parse(c.Params("tagId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tag ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	sortBy := c.Query("sort", "hot")

	var sortByEnum repositories.PostSortBy
	switch sortBy {
	case "hot":
		sortByEnum = repositories.SortByHot
	case "new":
		sortByEnum = repositories.SortByNew
	case "top":
		sortByEnum = repositories.SortByTop
	default:
		sortByEnum = repositories.SortByHot
	}

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	posts, err := h.postService.ListPostsByTagID(c.Context(), tagID, offset, limit, sortByEnum, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve posts", err)
	}

	return utils.SuccessResponse(c, "Posts retrieved successfully", posts)
}

// SearchPosts searches for posts
func (h *PostHandler) SearchPosts(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return utils.ValidationErrorResponse(c, "Search query is required")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	posts, err := h.postService.SearchPosts(c.Context(), query, offset, limit, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search posts", err)
	}

	return utils.SuccessResponse(c, "Posts retrieved successfully", posts)
}

// CreateCrosspost creates a crosspost
func (h *PostHandler) CreateCrosspost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	sourcePostID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	var req dto.CreatePostRequest
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

	post, err := h.postService.CreateCrosspost(c.Context(), userID, sourcePostID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create crosspost", err)
	}

	return utils.SuccessResponse(c, "Crosspost created successfully", post)
}

// GetCrossposts retrieves crossposts of a post
func (h *PostHandler) GetCrossposts(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid post ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Get userID if authenticated (optional)
	var userIDPtr *uuid.UUID
	if userID, ok := c.Locals("userID").(uuid.UUID); ok {
		userIDPtr = &userID
	}

	posts, err := h.postService.GetCrossposts(c.Context(), postID, offset, limit, userIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve crossposts", err)
	}

	return utils.SuccessResponse(c, "Crossposts retrieved successfully", posts)
}

// GetFeed retrieves personalized feed for authenticated user
func (h *PostHandler) GetFeed(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	sortBy := c.Query("sort", "hot")

	var sortByEnum repositories.PostSortBy
	switch sortBy {
	case "hot":
		sortByEnum = repositories.SortByHot
	case "new":
		sortByEnum = repositories.SortByNew
	case "top":
		sortByEnum = repositories.SortByTop
	default:
		sortByEnum = repositories.SortByHot
	}

	feed, err := h.postService.GetFeed(c.Context(), userID, offset, limit, sortByEnum)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve feed", err)
	}

	return utils.SuccessResponse(c, "Feed retrieved successfully", feed)
}
