package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PushSubscriptionRepositoryImpl struct {
	db *gorm.DB
}

func NewPushSubscriptionRepository(db *gorm.DB) repositories.PushSubscriptionRepository {
	return &PushSubscriptionRepositoryImpl{db: db}
}

func (r *PushSubscriptionRepositoryImpl) Upsert(ctx context.Context, subscription *models.PushSubscription) error {
	// Use GORM's Clauses for UPSERT (INSERT ... ON CONFLICT ... DO UPDATE)
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "endpoint"}},
			DoUpdates: clause.AssignmentColumns([]string{"p256dh", "auth", "expiration_time", "updated_at"}),
		}).
		Create(subscription).Error
}

func (r *PushSubscriptionRepositoryImpl) Delete(ctx context.Context, userID uuid.UUID, endpoint string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND endpoint = ?", userID, endpoint).
		Delete(&models.PushSubscription{}).Error
}

func (r *PushSubscriptionRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PushSubscription, error) {
	var subscriptions []*models.PushSubscription
	now := time.Now().Unix() * 1000 // Convert to milliseconds

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("expiration_time IS NULL OR expiration_time > ?", now).
		Find(&subscriptions).Error

	return subscriptions, err
}

func (r *PushSubscriptionRepositoryImpl) DeleteByEndpoint(ctx context.Context, endpoint string) error {
	return r.db.WithContext(ctx).
		Where("endpoint = ?", endpoint).
		Delete(&models.PushSubscription{}).Error
}

func (r *PushSubscriptionRepositoryImpl) DeleteExpired(ctx context.Context) error {
	now := time.Now().Unix() * 1000 // Convert to milliseconds

	return r.db.WithContext(ctx).
		Where("expiration_time IS NOT NULL AND expiration_time < ?", now).
		Delete(&models.PushSubscription{}).Error
}

// Compiler check to ensure implementation satisfies interface
var _ repositories.PushSubscriptionRepository = (*PushSubscriptionRepositoryImpl)(nil)
