package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type NotificationService interface {
	// Get notifications
	GetNotifications(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.NotificationListResponse, error)
	GetUnreadNotifications(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.NotificationListResponse, error)
	GetNotification(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) (*dto.NotificationResponse, error)

	// Mark as read
	MarkAsRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error

	// Delete notifications
	DeleteNotification(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error
	DeleteAllNotifications(ctx context.Context, userID uuid.UUID) error

	// Count
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)

	// Notification settings
	GetSettings(ctx context.Context, userID uuid.UUID) (*dto.NotificationSettingsResponse, error)
	UpdateSettings(ctx context.Context, userID uuid.UUID, req *dto.NotificationSettingsRequest) (*dto.NotificationSettingsResponse, error)

	// Internal methods for creating notifications (used by other services)
	CreateNotification(ctx context.Context, userID uuid.UUID, senderID uuid.UUID, notifType string, message string, postID *uuid.UUID, commentID *uuid.UUID) error
}
