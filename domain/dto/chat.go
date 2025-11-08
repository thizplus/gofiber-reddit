package dto

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Conversation DTOs
// ============================================================================

// ConversationResponse - Single conversation
type ConversationResponse struct {
	ID               uuid.UUID        `json:"id"`
	OtherUser        UserResponse     `json:"otherUser"` // The other participant
	LastMessage      *MessageResponse `json:"lastMessage,omitempty"`
	LastMessageAt    time.Time        `json:"lastMessageAt"`
	UnreadCount      int              `json:"unreadCount"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}

// ConversationListResponse - List of conversations with cursor pagination
type ConversationListResponse struct {
	Conversations []ConversationResponse `json:"conversations"`
	NextCursor    *string                `json:"nextCursor,omitempty"` // Base64 encoded timestamp
	HasMore       bool                   `json:"hasMore"`
}

// UnreadCountResponse - Total unread message count
type UnreadCountResponse struct {
	TotalUnread int `json:"totalUnread"`
}

// CreateConversationRequest - Request to create/get a conversation
type CreateConversationRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
}

// ============================================================================
// Message DTOs
// ============================================================================

// MessageMedia - Media attachment metadata
type MessageMedia struct {
	URL       string  `json:"url"`
	Thumbnail *string `json:"thumbnail,omitempty"`
	Type      string  `json:"type"` // "image", "video", "file"
	Filename  *string `json:"filename,omitempty"`
	MimeType  *string `json:"mimeType,omitempty"`
	Size      *int64  `json:"size,omitempty"`
	Width     *int    `json:"width,omitempty"`
	Height    *int    `json:"height,omitempty"`
	Duration  *int    `json:"duration,omitempty"` // seconds, for videos
	// Video streaming fields (for Bunny Stream HLS)
	VideoID          *string `json:"videoId,omitempty"`
	HLSURL           *string `json:"hlsUrl,omitempty"`
	EncodingStatus   *string `json:"encodingStatus,omitempty"`   // "pending", "processing", "completed", "failed"
	EncodingProgress *int    `json:"encodingProgress,omitempty"` // 0-100
}

// SendMessageRequest - Request to send a message
type SendMessageRequest struct {
	ConversationID uuid.UUID      `json:"conversationId" validate:"required,uuid"`
	Type           string         `json:"type" validate:"required,oneof=text image video file"`
	Content        *string        `json:"content,omitempty" validate:"omitempty,min=1,max=5000"`
	Media          []MessageMedia `json:"media,omitempty"`
	TempID         *string        `json:"tempId,omitempty"` // Client-generated ID for optimistic updates
}

// MessageResponse - Single message
type MessageResponse struct {
	ID             uuid.UUID      `json:"id"`
	ConversationID uuid.UUID      `json:"conversationId"`
	Sender         UserResponse   `json:"sender"`
	Receiver       UserResponse   `json:"receiver"`
	Type           string         `json:"type"` // "text", "image", "video", "file"
	Content        *string        `json:"content,omitempty"`
	Media          []MessageMedia `json:"media,omitempty"`
	IsRead         bool           `json:"isRead"`
	ReadAt         *time.Time     `json:"readAt,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	TempID         *string        `json:"tempId,omitempty"` // Echo back client's tempId if provided

	// Helper fields for frontend (denormalized for convenience)
	SenderId   uuid.UUID `json:"senderId"`   // Same as Sender.ID, for easier access
}

// MessageListResponse - List of messages with cursor pagination
type MessageListResponse struct {
	Messages   []MessageResponse `json:"messages"`
	NextCursor *string           `json:"nextCursor,omitempty"` // Base64 encoded timestamp
	HasMore    bool              `json:"hasMore"`
}

// MessageContextResponse - Message with surrounding context (for jump to message)
type MessageContextResponse struct {
	TargetMessage  MessageResponse   `json:"targetMessage"`
	Before         []MessageResponse `json:"before"` // 20 messages before
	After          []MessageResponse `json:"after"`  // 20 messages after
	BeforeCursor   *string           `json:"beforeCursor,omitempty"`
	AfterCursor    *string           `json:"afterCursor,omitempty"`
	HasMoreBefore  bool              `json:"hasMoreBefore"`
	HasMoreAfter   bool              `json:"hasMoreAfter"`
}

// MarkMessagesAsReadRequest - Request to mark messages as read
type MarkMessagesAsReadRequest struct {
	ConversationID uuid.UUID `json:"conversationId" validate:"required,uuid"`
}

// ============================================================================
// Block DTOs
// ============================================================================

// BlockUserRequest - Request to block a user
type BlockUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
}

// BlockStatusResponse - Block status between two users
type BlockStatusResponse struct {
	IsBlocked       bool `json:"isBlocked"`       // Current user blocked the other user
	IsBlockedBy     bool `json:"isBlockedBy"`     // Current user is blocked by the other user
	CanMessage      bool `json:"canMessage"`      // Whether messaging is allowed
}

// BlockedUserResponse - Blocked user info
type BlockedUserResponse struct {
	User      UserResponse `json:"user"`
	BlockedAt time.Time    `json:"blockedAt"`
}

// BlockedUsersResponse - List of blocked users
type BlockedUsersResponse struct {
	BlockedUsers []BlockedUserResponse `json:"blockedUsers"`
	Meta         PaginationMeta        `json:"meta"`
}

// ============================================================================
// Search Users for Chat DTOs
// ============================================================================

// ChatUserSearchResult - User info for chat search
type ChatUserSearchResult struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	DisplayName string     `json:"displayName"`
	Avatar      string     `json:"avatar,omitempty"`
	Bio         string     `json:"bio,omitempty"`
	IsFollowing bool       `json:"isFollowing"`
	IsOnline    bool       `json:"isOnline"`
	LastActive  *time.Time `json:"lastActive,omitempty"`
}

// ChatUserSearchResponse - Response for chat user search
type ChatUserSearchResponse struct {
	Users     []ChatUserSearchResult `json:"users"`
	Suggested []ChatUserSearchResult `json:"suggested,omitempty"` // Suggested users (followers, following, etc.)
}
