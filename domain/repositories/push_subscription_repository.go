package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type PushSubscriptionRepository interface {
	// Create or update subscription
	Upsert(ctx context.Context, subscription *models.PushSubscription) error

	// Delete subscription
	Delete(ctx context.Context, userID uuid.UUID, endpoint string) error

	// Get all active subscriptions for a user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PushSubscription, error)

	// Delete subscription by endpoint (used when subscription expires)
	DeleteByEndpoint(ctx context.Context, endpoint string) error

	// Delete all expired subscriptions
	DeleteExpired(ctx context.Context) error
}
