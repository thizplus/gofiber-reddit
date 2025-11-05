package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type PushService interface {
	// Subscribe a user to push notifications
	Subscribe(ctx context.Context, userID uuid.UUID, req *dto.PushSubscriptionRequest) (*dto.PushSubscriptionResponse, error)

	// Unsubscribe a user from push notifications
	Unsubscribe(ctx context.Context, userID uuid.UUID, req *dto.PushSubscriptionRequest) error

	// Send push notification to a user
	SendToUser(ctx context.Context, userID uuid.UUID, payload *dto.PushNotificationPayload) error

	// Get VAPID public key (for frontend)
	GetPublicKey() string
}
