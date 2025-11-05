package serviceimpl

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/websocket"
)

type NotificationServiceImpl struct {
	notifRepo         repositories.NotificationRepository
	notifSettingsRepo repositories.NotificationSettingsRepository
	userRepo          repositories.UserRepository
	pushService       services.PushService
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
		pushService:       nil, // Will be set later via SetPushService
	}
}

// SetPushService sets the push service (to avoid circular dependency)
func (s *NotificationServiceImpl) SetPushService(pushService services.PushService) {
	s.pushService = pushService
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

	err = s.notifRepo.MarkAsRead(ctx, notificationID)
	if err != nil {
		return err
	}

	// Send real-time update via WebSocket
	websocket.Manager.BroadcastToUser(userID, "notification_read", map[string]interface{}{
		"notificationId": notificationID,
		"unreadCount":    s.getUnreadCount(ctx, userID),
	})

	return nil
}

func (s *NotificationServiceImpl) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	err := s.notifRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return err
	}

	// Send real-time update via WebSocket
	websocket.Manager.BroadcastToUser(userID, "notification_read_all", map[string]interface{}{
		"unreadCount": 0,
	})

	return nil
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

	// Create notification in database
	err := s.notifRepo.Create(ctx, notification)
	if err != nil {
		return err
	}

	// Fetch notification with relations for real-time broadcast
	createdNotification, err := s.notifRepo.GetByID(ctx, notification.ID)
	if err != nil {
		log.Printf("Warning: Failed to fetch notification for WebSocket broadcast: %v", err)
		return nil // Don't fail the whole operation
	}

	// Convert to DTO
	notificationDTO := dto.NotificationToNotificationResponse(createdNotification)

	// Send real-time notification via WebSocket
	websocket.Manager.BroadcastToUser(userID, "notification", map[string]interface{}{
		"notification": notificationDTO,
		"unreadCount":  s.getUnreadCount(ctx, userID),
	})

	log.Printf("📬 Real-time notification sent to user %s: %s", userID.String(), message)

	// Send push notification (if user is offline and pushService is available)
	if s.pushService != nil {
		// Prepare push payload
		pushPayload := &dto.PushNotificationPayload{
			Title: "VOOBIZE",
			Body:  message,
			Icon:  "/logo.png",
			Badge: "/logo.png",
			Tag:   notifType,
			Data: map[string]interface{}{
				"notificationId": notification.ID.String(),
				"url":            s.buildNotificationURL(postID, commentID),
			},
		}

		// Send push notification (non-blocking)
		go func() {
			if err := s.pushService.SendToUser(context.Background(), userID, pushPayload); err != nil {
				log.Printf("⚠️  Failed to send push notification: %v", err)
			}
		}()
	}

	return nil
}

// Helper function to get unread count
func (s *NotificationServiceImpl) getUnreadCount(ctx context.Context, userID uuid.UUID) int64 {
	count, err := s.notifRepo.CountUnreadByUser(ctx, userID)
	if err != nil {
		return 0
	}
	return count
}

// Helper function to build notification URL
func (s *NotificationServiceImpl) buildNotificationURL(postID, commentID *uuid.UUID) string {
	if commentID != nil {
		return "/post/" + postID.String() + "#comment-" + commentID.String()
	}
	if postID != nil {
		return "/post/" + postID.String()
	}
	return "/notifications"
}

var _ services.NotificationService = (*NotificationServiceImpl)(nil)
