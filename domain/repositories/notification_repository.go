package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	// Create notification
	Create(ctx context.Context, notification *models.Notification) error

	// Get notifications
	GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Notification, error)
	ListUnreadByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Notification, error)

	// Mark as read
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error

	// Delete
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteAllByUser(ctx context.Context, userID uuid.UUID) error

	// Count
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	CountUnreadByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}
