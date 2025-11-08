package repositories

import (
	"context"
	"time"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type MessageRepository interface {
	// Basic CRUD
	Create(ctx context.Context, message *models.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	Update(ctx context.Context, id uuid.UUID, message *models.Message) error
	Delete(ctx context.Context, id uuid.UUID) error

	// List messages (cursor-based pagination)
	// Cursor is based on created_at timestamp
	ListByConversation(ctx context.Context, conversationID uuid.UUID, beforeCursor *time.Time, limit int) ([]*models.Message, error)

	// Jump to message with context (for search results, media tabs, etc.)
	GetMessagesBeforeTimestamp(ctx context.Context, conversationID uuid.UUID, timestamp time.Time, limit int) ([]*models.Message, error)
	GetMessagesAfterTimestamp(ctx context.Context, conversationID uuid.UUID, timestamp time.Time, limit int) ([]*models.Message, error)

	// Mark messages as read
	MarkAsRead(ctx context.Context, messageID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error

	// Stats
	Count(ctx context.Context) (int64, error)
	CountByConversation(ctx context.Context, conversationID uuid.UUID) (int64, error)
	CountUnread(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (int64, error)

	// Phase 2: Media/Links/Files Queries
	ListMediaMessages(ctx context.Context, conversationID uuid.UUID, mediaType *string, beforeCursor *time.Time, limit int) ([]*models.Message, error)
	ListMessagesWithLinks(ctx context.Context, conversationID uuid.UUID, beforeCursor *time.Time, limit int) ([]*models.Message, error)
	ListFileMessages(ctx context.Context, conversationID uuid.UUID, beforeCursor *time.Time, limit int) ([]*models.Message, error)
}
