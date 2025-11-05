package serviceimpl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type NotificationServiceImpl struct {
	notifRepo         repositories.NotificationRepository
	notifSettingsRepo repositories.NotificationSettingsRepository
	userRepo          repositories.UserRepository
}

func NewNotificationService(
	notifRepo repositories.NotificationRepository,
	notifSettingsRepo repositories.NotificationSettingsRepository,
	userRepo repositories.UserRepository,
) services.NotificationService {
	return &NotificationServiceImpl{
		notifRepo:         notifRepo,
		notifSettingsRepo: notifSettingsRepo,
		userRepo:          userRepo,
	}
}

func (s *NotificationServiceImpl) GetNotifications(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.NotificationListResponse, error) {
	notifications, err := s.notifRepo.ListByUser(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.notifRepo.CountUnreadByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	count, err := s.notifRepo.CountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.NotificationResponse, len(notifications))
	for i, notif := range notifications {
		responses[i] = *dto.NotificationToNotificationResponse(notif)
	}

	return &dto.NotificationListResponse{
		Notifications: responses,
		UnreadCount:   unreadCount,
		Meta: dto.PaginationMeta{
			Total:  count,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (s *NotificationServiceImpl) GetUnreadNotifications(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.NotificationListResponse, error) {
	notifications, err := s.notifRepo.ListUnreadByUser(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.notifRepo.CountUnreadByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.NotificationResponse, len(notifications))
	for i, notif := range notifications {
		responses[i] = *dto.NotificationToNotificationResponse(notif)
	}

	return &dto.NotificationListResponse{
		Notifications: responses,
		UnreadCount:   unreadCount,
		Meta: dto.PaginationMeta{
			Total:  unreadCount,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (s *NotificationServiceImpl) GetNotification(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) (*dto.NotificationResponse, error) {
	notification, err := s.notifRepo.GetByID(ctx, notificationID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if notification.UserID != userID {
		return nil, errors.New("unauthorized: not notification owner")
	}

	return dto.NotificationToNotificationResponse(notification), nil
}

func (s *NotificationServiceImpl) MarkAsRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
	// Get notification
	notification, err := s.notifRepo.GetByID(ctx, notificationID)
	if err != nil {
		return err
	}

	// Check ownership
	if notification.UserID != userID {
		return errors.New("unauthorized: not notification owner")
	}

	return s.notifRepo.MarkAsRead(ctx, notificationID)
}

func (s *NotificationServiceImpl) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.notifRepo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationServiceImpl) DeleteNotification(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
	// Get notification
	notification, err := s.notifRepo.GetByID(ctx, notificationID)
	if err != nil {
		return err
	}

	// Check ownership
	if notification.UserID != userID {
		return errors.New("unauthorized: not notification owner")
	}

	return s.notifRepo.Delete(ctx, notificationID)
}

func (s *NotificationServiceImpl) DeleteAllNotifications(ctx context.Context, userID uuid.UUID) error {
	return s.notifRepo.DeleteAllByUser(ctx, userID)
}

func (s *NotificationServiceImpl) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.notifRepo.CountUnreadByUser(ctx, userID)
}

func (s *NotificationServiceImpl) GetSettings(ctx context.Context, userID uuid.UUID) (*dto.NotificationSettingsResponse, error) {
	settings, err := s.notifSettingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		// If not found, create default settings
		defaultSettings := &models.NotificationSettings{
			UserID:             userID,
			Replies:            true,
			Mentions:           true,
			Votes:              false,
			Follows:            true,
			EmailNotifications: false,
			UpdatedAt:          time.Now(),
		}
		err = s.notifSettingsRepo.Create(ctx, defaultSettings)
		if err != nil {
			return nil, err
		}
		return dto.NotificationSettingsToResponse(defaultSettings), nil
	}

	return dto.NotificationSettingsToResponse(settings), nil
}

func (s *NotificationServiceImpl) UpdateSettings(ctx context.Context, userID uuid.UUID, req *dto.NotificationSettingsRequest) (*dto.NotificationSettingsResponse, error) {
	// Get existing settings
	settings, err := s.notifSettingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Create if not exists
		settings = &models.NotificationSettings{
			UserID: userID,
		}
	}

	// Update fields
	if req.Replies != nil {
		settings.Replies = *req.Replies
	}
	if req.Mentions != nil {
		settings.Mentions = *req.Mentions
	}
	if req.Votes != nil {
		settings.Votes = *req.Votes
	}
	if req.Follows != nil {
		settings.Follows = *req.Follows
	}
	if req.EmailNotifications != nil {
		settings.EmailNotifications = *req.EmailNotifications
	}
	settings.UpdatedAt = time.Now()

	err = s.notifSettingsRepo.Update(ctx, userID, settings)
	if err != nil {
		return nil, err
	}

	return dto.NotificationSettingsToResponse(settings), nil
}

func (s *NotificationServiceImpl) CreateNotification(ctx context.Context, userID uuid.UUID, senderID uuid.UUID, notifType string, message string, postID *uuid.UUID, commentID *uuid.UUID) error {
	// Check if user wants to receive this notification type
	shouldNotify, _ := s.notifSettingsRepo.ShouldNotify(ctx, userID, notifType)
	if !shouldNotify {
		return nil // User has disabled this notification type
	}

	notification := &models.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		SenderID:  senderID,
		Type:      notifType,
		Message:   message,
		PostID:    postID,
		CommentID: commentID,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	return s.notifRepo.Create(ctx, notification)
}

var _ services.NotificationService = (*NotificationServiceImpl)(nil)
