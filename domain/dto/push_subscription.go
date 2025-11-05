package dto

import (
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

// PushSubscriptionRequest represents the subscription data from the browser
type PushSubscriptionRequest struct {
	Endpoint       string                 `json:"endpoint" validate:"required,url"`
	ExpirationTime *int64                 `json:"expirationTime"`
	Keys           PushSubscriptionKeys   `json:"keys" validate:"required"`
}

type PushSubscriptionKeys struct {
	P256dh string `json:"p256dh" validate:"required"`
	Auth   string `json:"auth" validate:"required"`
}

// PushSubscriptionResponse represents the response after saving subscription
type PushSubscriptionResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"userId"`
	Endpoint string    `json:"endpoint"`
}

// PushNotificationPayload represents the notification sent to the user
type PushNotificationPayload struct {
	Title string                 `json:"title"`
	Body  string                 `json:"body"`
	Icon  string                 `json:"icon"`
	Badge string                 `json:"badge,omitempty"`
	Tag   string                 `json:"tag,omitempty"`
	Data  map[string]interface{} `json:"data,omitempty"`
}

// RequestToModel converts DTO to model
func PushSubscriptionRequestToModel(req *PushSubscriptionRequest, userID uuid.UUID) *models.PushSubscription {
	return &models.PushSubscription{
		ID:             uuid.New(),
		UserID:         userID,
		Endpoint:       req.Endpoint,
		P256dh:         req.Keys.P256dh,
		Auth:           req.Keys.Auth,
		ExpirationTime: req.ExpirationTime,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// ModelToResponse converts model to response DTO
func PushSubscriptionToResponse(subscription *models.PushSubscription) *PushSubscriptionResponse {
	return &PushSubscriptionResponse{
		ID:       subscription.ID,
		UserID:   subscription.UserID,
		Endpoint: subscription.Endpoint,
	}
}
