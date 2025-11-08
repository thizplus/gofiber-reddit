package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type ConversationService interface {
	// Get or create conversation
	GetOrCreateConversation(ctx context.Context, userID uuid.UUID, otherUsername string) (*dto.ConversationResponse, error)
	GetConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*dto.ConversationResponse, error)

	// List conversations with cursor pagination
	ListConversations(ctx context.Context, userID uuid.UUID, cursor *string, limit int) (*dto.ConversationListResponse, error)

	// Unread counts
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (*dto.UnreadCountResponse, error)

	// Mark as read
	MarkAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error

	// Search users for chat
	SearchUsersForChat(ctx context.Context, userID uuid.UUID, query string, limit int) (*dto.ChatUserSearchResponse, error)
}
