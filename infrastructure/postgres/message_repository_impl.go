package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type MessageRepositoryImpl struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repositories.MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

func (r *MessageRepositoryImpl) Create(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *MessageRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	var message models.Message
	err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepositoryImpl) Update(ctx context.Context, id uuid.UUID, message *models.Message) error {
	return r.db.WithContext(ctx).
		Model(&models.Message{}).
		Where("id = ?", id).
		Updates(message).Error
}

func (r *MessageRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Message{}, "id = ?", id).Error
}

func (r *MessageRepositoryImpl) ListByConversation(ctx context.Context, conversationID uuid.UUID, beforeCursor *time.Time, limit int) ([]*models.Message, error) {
	query := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC") // Most recent first

	// Apply cursor pagination (messages created before cursor)
	if beforeCursor != nil {
		query = query.Where("created_at < ?", *beforeCursor)
	}

	var messages []*models.Message
	err := query.Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *MessageRepositoryImpl) GetMessagesBeforeTimestamp(ctx context.Context, conversationID uuid.UUID, timestamp time.Time, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		Where("conversation_id = ? AND created_at < ?", conversationID, timestamp).
		Order("created_at DESC"). // Most recent first
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *MessageRepositoryImpl) GetMessagesAfterTimestamp(ctx context.Context, conversationID uuid.UUID, timestamp time.Time, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		Where("conversation_id = ? AND created_at > ?", conversationID, timestamp).
		Order("created_at ASC"). // Oldest first (to get next messages)
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *MessageRepositoryImpl) MarkAsRead(ctx context.Context, messageID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *MessageRepositoryImpl) MarkAllAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Message{}).
		Where("conversation_id = ? AND receiver_id = ? AND is_read = ?", conversationID, userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *MessageRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Message{}).Count(&count).Error
	return count, err
}

func (r *MessageRepositoryImpl) CountByConversation(ctx context.Context, conversationID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Message{}).
		Where("conversation_id = ?", conversationID).
		Count(&count).Error
	return count, err
}

func (r *MessageRepositoryImpl) CountUnread(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Message{}).
		Where("conversation_id = ? AND receiver_id = ? AND is_read = ?", conversationID, userID, false).
		Count(&count).Error
	return count, err
}

// Phase 2: Media/Links/Files Queries

func (r *MessageRepositoryImpl) ListMediaMessages(ctx context.Context, conversationID uuid.UUID, mediaType *string, beforeCursor *time.Time, limit int) ([]*models.Message, error) {
	query := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		Where("conversation_id = ?", conversationID).
		Where("type IN (?)", []string{"image", "video"}). // Media messages only
		Where("media IS NOT NULL AND media != '[]'").     // Has media content
		Order("created_at DESC")

	// Filter by specific media type if provided
	if mediaType != nil && *mediaType != "" {
		query = query.Where("type = ?", *mediaType)
	}

	// Apply cursor pagination
	if beforeCursor != nil {
		query = query.Where("created_at < ?", *beforeCursor)
	}

	var messages []*models.Message
	err := query.Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *MessageRepositoryImpl) ListMessagesWithLinks(ctx context.Context, conversationID uuid.UUID, beforeCursor *time.Time, limit int) ([]*models.Message, error) {
	query := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		Where("conversation_id = ?", conversationID).
		Where("type = ?", "text").
		Where("content IS NOT NULL").
		// PostgreSQL pattern matching for URLs (http://, https://, www.)
		Where("content ~ ?", `(https?://|www\.)[^\s]+`).
		Order("created_at DESC")

	// Apply cursor pagination
	if beforeCursor != nil {
		query = query.Where("created_at < ?", *beforeCursor)
	}

	var messages []*models.Message
	err := query.Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *MessageRepositoryImpl) ListFileMessages(ctx context.Context, conversationID uuid.UUID, beforeCursor *time.Time, limit int) ([]*models.Message, error) {
	query := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		Where("conversation_id = ?", conversationID).
		Where("type = ?", "file"). // File type messages
		Where("media IS NOT NULL AND media != '[]'").
		Order("created_at DESC")

	// Apply cursor pagination
	if beforeCursor != nil {
		query = query.Where("created_at < ?", *beforeCursor)
	}

	var messages []*models.Message
	err := query.Limit(limit).Find(&messages).Error
	return messages, err
}

// Ensure interface compliance
var _ repositories.MessageRepository = (*MessageRepositoryImpl)(nil)
