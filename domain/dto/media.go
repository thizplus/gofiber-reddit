package dto

import (
	"time"

	"github.com/google/uuid"
)

// MediaUploadResponse - Response for media upload
type MediaUploadResponse struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	FileName  string    `json:"fileName"`
	MimeType  string    `json:"mimeType"`
	Size      int64     `json:"size"`
	URL       string    `json:"url"`
	Thumbnail string    `json:"thumbnail,omitempty"`
	Width     int       `json:"width,omitempty"`
	Height    int       `json:"height,omitempty"`
	Duration  float64   `json:"duration,omitempty"` // For videos
	CreatedAt time.Time `json:"createdAt"`
}

// MediaResponse - Response for a single media
type MediaResponse struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"` // "image", "video"
	FileName  string    `json:"fileName"`
	MimeType  string    `json:"mimeType"`
	Size      int64     `json:"size"`
	URL       string    `json:"url"`
	Thumbnail string    `json:"thumbnail,omitempty"`
	Width     int       `json:"width,omitempty"`
	Height    int       `json:"height,omitempty"`
	Duration  float64   `json:"duration,omitempty"` // For videos
	CreatedAt time.Time `json:"createdAt"`
}

// MediaListResponse - Response for listing media
type MediaListResponse struct {
	Media []MediaResponse `json:"media"`
	Meta  PaginationMeta  `json:"meta"`
}

// MediaUploadRequest - Request for uploading media (multipart form)
type MediaUploadRequest struct {
	Type string `form:"type" validate:"required,oneof=image video"` // Validated from form
	// File will be handled separately via c.FormFile("file")
}
