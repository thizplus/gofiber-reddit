package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateCommentRequest - Request for creating a new comment
type CreateCommentRequest struct {
	PostID   uuid.UUID  `json:"postId" validate:"required,uuid"`
	ParentID *uuid.UUID `json:"parentId" validate:"omitempty,uuid"` // For replies
	Content  string     `json:"content" validate:"required,min=1,max=10000"`
}

// UpdateCommentRequest - Request for updating a comment
type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=10000"`
}

// CommentResponse - Response for a single comment
type CommentResponse struct {
	ID        uuid.UUID            `json:"id"`
	PostID    uuid.UUID            `json:"postId"`
	Post      *PostSummaryResponse `json:"post,omitempty"` // Post info (title, author)
	ParentID  *uuid.UUID           `json:"parentId,omitempty"`
	Author    UserResponse         `json:"author"`
	Content   string               `json:"content"`
	Votes     int                  `json:"votes"`
	Depth     int                  `json:"depth"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`

	// User-specific fields (when authenticated)
	UserVote   *string `json:"userVote,omitempty"`   // "up", "down", or null
	ReplyCount *int    `json:"replyCount,omitempty"` // Number of direct replies
	IsDeleted  bool    `json:"isDeleted"`
}

// CommentWithRepliesResponse - Comment with nested replies
type CommentWithRepliesResponse struct {
	CommentResponse
	Replies []CommentWithRepliesResponse `json:"replies"`
}

// CommentListResponse - Response for listing comments
type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
	Meta     PaginationMeta    `json:"meta"`
}

// CommentTreeResponse - Response for comment tree (nested structure)
type CommentTreeResponse struct {
	Comments []CommentWithRepliesResponse `json:"comments"`
	Total    int64                        `json:"total"`
}
