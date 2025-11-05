package dto

import (
	"time"

	"github.com/google/uuid"
)

// NotificationResponse - Response for a single notification
type NotificationResponse struct {
	ID        uuid.UUID    `json:"id"`
	User      UserResponse `json:"user"`
	Sender    UserResponse `json:"sender"`
	Type      string       `json:"type"` // "reply", "vote", "mention", "follow"
	Message   string       `json:"message"`
	PostID    *uuid.UUID   `json:"postId,omitempty"`
	CommentID *uuid.UUID   `json:"commentId,omitempty"`
	IsRead    bool         `json:"isRead"`
	CreatedAt time.Time    `json:"createdAt"`
}

// NotificationListResponse - Response for listing notifications
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int64                  `json:"unreadCount"`
	Meta          PaginationMeta         `json:"meta"`
}

// MarkAsReadRequest - Request for marking notification as read
type MarkAsReadRequest struct {
	NotificationID uuid.UUID `json:"notificationId" validate:"required,uuid"`
}

// NotificationSettingsRequest - Request for updating notification settings
type NotificationSettingsRequest struct {
	Replies            *bool `json:"replies" validate:"omitempty"`
	Mentions           *bool `json:"mentions" validate:"omitempty"`
	Votes              *bool `json:"votes" validate:"omitempty"`
	Follows            *bool `json:"follows" validate:"omitempty"`
	EmailNotifications *bool `json:"emailNotifications" validate:"omitempty"`
}

// NotificationSettingsResponse - Response for notification settings
type NotificationSettingsResponse struct {
	UserID             uuid.UUID `json:"userId"`
	Replies            bool      `json:"replies"`
	Mentions           bool      `json:"mentions"`
	Votes              bool      `json:"votes"`
	Follows            bool      `json:"follows"`
	EmailNotifications bool      `json:"emailNotifications"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
