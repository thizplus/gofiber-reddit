package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type MessageService interface {
	// Send message
	SendMessage(ctx context.Context, userID uuid.UUID, req *dto.SendMessageRequest) (*dto.MessageResponse, error)

	// Get message
	GetMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (*dto.MessageResponse, error)

	// List messages with cursor pagination
	ListMessages(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, cursor *string, limit int) (*dto.MessageListResponse, error)

	// Jump to message with context (for search, media tabs, etc.)
	GetMessageContext(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (*dto.MessageContextResponse, error)

	// Mark as read
	MarkAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error

	// Phase 2: Media/Links/Files Queries
	ListMediaMessages(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, mediaType *string, cursor *string, limit int) (*dto.MessageListResponse, error)
	ListMessagesWithLinks(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, cursor *string, limit int) (*dto.MessageListResponse, error)
	ListFileMessages(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, cursor *string, limit int) (*dto.MessageListResponse, error)
}
