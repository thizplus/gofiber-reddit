package dto

import (
	"time"

	"github.com/google/uuid"
)

// TagResponse - Response for a single tag
type TagResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	PostCount int       `json:"postCount"`
	CreatedAt time.Time `json:"createdAt"`
}

// TagListResponse - Response for listing tags
type TagListResponse struct {
	Tags []TagResponse  `json:"tags"`
	Meta PaginationMeta `json:"meta"`
}

// PopularTagsResponse - Response for popular tags
type PopularTagsResponse struct {
	Tags []TagResponse `json:"tags"`
}

// SearchTagsRequest - Request for searching tags
type SearchTagsRequest struct {
	Query string `json:"query" validate:"required,min=1,max=50"`
	Limit int    `json:"limit" validate:"omitempty,min=1,max=50"`
}
