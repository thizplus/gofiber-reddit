package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreatePostRequest - Request for creating a new post
type CreatePostRequest struct {
	Title        string      `json:"title" validate:"required,min=1,max=300"`
	Content      string      `json:"content" validate:"required,min=1,max=40000"`
	MediaIDs     []uuid.UUID `json:"mediaIds" validate:"omitempty,dive,uuid"`
	Tags         []string    `json:"tags" validate:"omitempty,max=5,dive,min=1,max=50"`
	SourcePostID *uuid.UUID  `json:"sourcePostId" validate:"omitempty,uuid"` // For crossposting
}

// UpdatePostRequest - Request for updating a post
type UpdatePostRequest struct {
	Title   string   `json:"title" validate:"omitempty,min=1,max=300"`
	Content string   `json:"content" validate:"omitempty,min=1,max=40000"`
	Tags    []string `json:"tags" validate:"omitempty,max=5,dive,min=1,max=50"`
}

// PostResponse - Response for a single post
type PostResponse struct {
	ID           uuid.UUID      `json:"id"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	Author       UserResponse   `json:"author"`
	Votes        int            `json:"votes"`
	CommentCount int            `json:"commentCount"`
	Media        []MediaResponse `json:"media,omitempty"`
	Tags         []TagResponse  `json:"tags,omitempty"`
	SourcePost   *PostResponse  `json:"sourcePost,omitempty"` // For crossposts
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`

	// User-specific fields (when authenticated)
	UserVote  *string `json:"userVote,omitempty"`  // "up", "down", or null
	IsSaved   *bool   `json:"isSaved,omitempty"`   // true/false
	HotScore  *float64 `json:"hotScore,omitempty"` // For debugging/sorting
}

// PostListResponse - Response for listing posts
type PostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Meta  PaginationMeta `json:"meta"`
}

// PostFeedResponse - Response for feed with mixed content types
type PostFeedResponse struct {
	Posts []PostResponse `json:"posts"`
	Meta  PaginationMeta `json:"meta"`
}

// PostSummaryResponse - Lightweight post info for nested responses (comments, etc.)
type PostSummaryResponse struct {
	ID        uuid.UUID    `json:"id"`
	Title     string       `json:"title"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"createdAt"`
}
