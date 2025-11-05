package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type NotificationSettingsRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationSettingsRepository(db *gorm.DB) repositories.NotificationSettingsRepository {
	return &NotificationSettingsRepositoryImpl{db: db}
}

func (r *NotificationSettingsRepositoryImpl) Create(ctx context.Context, settings *models.NotificationSettings) error {
	return r.db.WithContext(ctx).Create(settings).Error
}

func (r *NotificationSettingsRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.NotificationSettings, error) {
	var settings models.NotificationSettings
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *NotificationSettingsRepositoryImpl) Update(ctx context.Context, userID uuid.UUID, settings *models.NotificationSettings) error {
	// Use map to ensure boolean false values are updated
	return r.db.WithContext(ctx).
		Model(&models.NotificationSettings{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"replies":             settings.Replies,
			"mentions":            settings.Mentions,
			"votes":               settings.Votes,
			"follows":             settings.Follows,
			"email_notifications": settings.EmailNotifications,
			"updated_at":          settings.UpdatedAt,
		}).Error
}

func (r *NotificationSettingsRepositoryImpl) ShouldNotify(ctx context.Context, userID uuid.UUID, notificationType string) (bool, error) {
	settings, err := r.GetByUserID(ctx, userID)
	if err != nil {
		// If no settings found, use default (true for most types)
		if err == gorm.ErrRecordNotFound {
			return notificationType != "votes", nil // Default: all except votes
		}
		return false, err
	}

	switch notificationType {
	case "reply", "replies":
		return settings.Replies, nil
	case "mention", "mentions":
		return settings.Mentions, nil
	case "vote", "votes":
		return settings.Votes, nil
	case "follow", "follows":
		return settings.Follows, nil
	default:
		return true, nil
	}
}

var _ repositories.NotificationSettingsRepository = (*NotificationSettingsRepositoryImpl)(nil)
