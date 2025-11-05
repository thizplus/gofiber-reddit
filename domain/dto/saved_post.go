package dto

import (
	"time"

	"github.com/google/uuid"
)

// SavePostRequest - Request for saving a post
type SavePostRequest struct {
	PostID uuid.UUID `json:"postId" validate:"required,uuid"`
}

// UnsavePostRequest - Request for unsaving a post
type UnsavePostRequest struct {
	PostID uuid.UUID `json:"postId" validate:"required,uuid"`
}

// SavedPostResponse - Response for saved post
type SavedPostResponse struct {
	PostID  uuid.UUID `json:"postId"`
	SavedAt time.Time `json:"savedAt"`
}

// SavedPostListResponse - Response for saved posts list
type SavedPostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Meta  PaginationMeta `json:"meta"`
}

// SaveStatusResponse - Response for checking if post is saved
type SaveStatusResponse struct {
	IsSaved bool `json:"isSaved"`
}
