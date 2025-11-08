package repositories

import (
	"context"
	"time"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type ConversationRepository interface {
	// Basic CRUD
	Create(ctx context.Context, conversation *models.Conversation) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error)
	Update(ctx context.Context, id uuid.UUID, conversation *models.Conversation) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Get or Create conversation between two users
	GetOrCreateByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Conversation, bool, error) // returns conversation, isNew, error

	// Get conversation by participants
	GetByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Conversation, error)

	// List conversations for a user (cursor-based pagination)
	ListByUser(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int) ([]*models.Conversation, error)

	// Unread count
	GetTotalUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)

	// Mark as read
	ResetUnreadCount(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error

	// Update conversation metadata
	UpdateLastMessage(ctx context.Context, conversationID uuid.UUID, messageID uuid.UUID, timestamp time.Time) error
	IncrementUnreadCount(ctx context.Context, conversationID uuid.UUID, receiverID uuid.UUID) error

	// Stats
	Count(ctx context.Context) (int64, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}
