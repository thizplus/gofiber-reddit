package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type ConversationRepositoryImpl struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) repositories.ConversationRepository {
	return &ConversationRepositoryImpl{db: db}
}

func (r *ConversationRepositoryImpl) Create(ctx context.Context, conversation *models.Conversation) error {
	return r.db.WithContext(ctx).Create(conversation).Error
}

func (r *ConversationRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		First(&conversation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepositoryImpl) Update(ctx context.Context, id uuid.UUID, conversation *models.Conversation) error {
	return r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Where("id = ?", id).
		Updates(conversation).Error
}

func (r *ConversationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Conversation{}, "id = ?", id).Error
}

func (r *ConversationRepositoryImpl) GetOrCreateByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Conversation, bool, error) {
	// Ensure consistent ordering (smaller UUID first)
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	// Try to get existing conversation
	var conversation models.Conversation
	err := r.db.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).
		First(&conversation).Error

	if err == nil {
		// Conversation exists
		return &conversation, false, nil
	}

	if err != gorm.ErrRecordNotFound {
		// Unexpected error
		return nil, false, err
	}

	// Create new conversation
	now := time.Now()
	conversation = models.Conversation{
		User1ID:          user1ID,
		User2ID:          user2ID,
		LastMessageAt:    now,
		User1UnreadCount: 0,
		User2UnreadCount: 0,
	}

	if err := r.db.WithContext(ctx).Create(&conversation).Error; err != nil {
		return nil, false, err
	}

	// Reload with preloaded users
	err = r.db.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		First(&conversation, "id = ?", conversation.ID).Error
	if err != nil {
		return nil, false, err
	}

	return &conversation, true, nil
}

func (r *ConversationRepositoryImpl) GetByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Conversation, error) {
	// Ensure consistent ordering
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	var conversation models.Conversation
	err := r.db.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).
		First(&conversation).Error

	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepositoryImpl) ListByUser(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int) ([]*models.Conversation, error) {
	query := r.db.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("last_message_at DESC")

	// Apply cursor pagination
	if cursor != nil {
		query = query.Where("last_message_at < ?", *cursor)
	}

	var conversations []*models.Conversation
	err := query.Limit(limit).Find(&conversations).Error
	return conversations, err
}

func (r *ConversationRepositoryImpl) GetTotalUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var totalUnread int64

	// Sum unread counts where user is user1
	var unreadAsUser1 int64
	err := r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Where("user1_id = ?", userID).
		Select("COALESCE(SUM(user1_unread_count), 0)").
		Scan(&unreadAsUser1).Error
	if err != nil {
		return 0, err
	}

	// Sum unread counts where user is user2
	var unreadAsUser2 int64
	err = r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Where("user2_id = ?", userID).
		Select("COALESCE(SUM(user2_unread_count), 0)").
		Scan(&unreadAsUser2).Error
	if err != nil {
		return 0, err
	}

	totalUnread = unreadAsUser1 + unreadAsUser2
	return int(totalUnread), nil
}

func (r *ConversationRepositoryImpl) ResetUnreadCount(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	// Get conversation to determine which field to reset
	var conversation models.Conversation
	err := r.db.WithContext(ctx).First(&conversation, "id = ?", conversationID).Error
	if err != nil {
		return err
	}

	// Reset appropriate unread count based on which user this is
	if conversation.User1ID == userID {
		return r.db.WithContext(ctx).
			Model(&models.Conversation{}).
			Where("id = ?", conversationID).
			Update("user1_unread_count", 0).Error
	} else if conversation.User2ID == userID {
		return r.db.WithContext(ctx).
			Model(&models.Conversation{}).
			Where("id = ?", conversationID).
			Update("user2_unread_count", 0).Error
	}

	return nil
}

func (r *ConversationRepositoryImpl) UpdateLastMessage(ctx context.Context, conversationID uuid.UUID, messageID uuid.UUID, timestamp time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Updates(map[string]interface{}{
			"last_message_id": messageID,
			"last_message_at": timestamp,
		}).Error
}

func (r *ConversationRepositoryImpl) IncrementUnreadCount(ctx context.Context, conversationID uuid.UUID, receiverID uuid.UUID) error {
	// Get conversation to determine which field to increment
	var conversation models.Conversation
	err := r.db.WithContext(ctx).First(&conversation, "id = ?", conversationID).Error
	if err != nil {
		return err
	}

	// Increment appropriate unread count
	if conversation.User1ID == receiverID {
		return r.db.WithContext(ctx).
			Model(&models.Conversation{}).
			Where("id = ?", conversationID).
			UpdateColumn("user1_unread_count", gorm.Expr("user1_unread_count + ?", 1)).Error
	} else if conversation.User2ID == receiverID {
		return r.db.WithContext(ctx).
			Model(&models.Conversation{}).
			Where("id = ?", conversationID).
			UpdateColumn("user2_unread_count", gorm.Expr("user2_unread_count + ?", 1)).Error
	}

	return nil
}

func (r *ConversationRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Conversation{}).Count(&count).Error
	return count, err
}

func (r *ConversationRepositoryImpl) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Count(&count).Error
	return count, err
}

// Ensure interface compliance
var _ repositories.ConversationRepository = (*ConversationRepositoryImpl)(nil)
