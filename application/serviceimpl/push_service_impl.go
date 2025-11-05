package serviceimpl

import (
	"context"
	"encoding/json"
	"log"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/config"
)

type PushServiceImpl struct {
	pushRepo repositories.PushSubscriptionRepository
	config   *config.Config
}

func NewPushService(
	pushRepo repositories.PushSubscriptionRepository,
	config *config.Config,
) services.PushService {
	return &PushServiceImpl{
		pushRepo: pushRepo,
		config:   config,
	}
}

func (s *PushServiceImpl) Subscribe(ctx context.Context, userID uuid.UUID, req *dto.PushSubscriptionRequest) (*dto.PushSubscriptionResponse, error) {
	// Convert DTO to model
	subscription := dto.PushSubscriptionRequestToModel(req, userID)

	// Upsert subscription (INSERT or UPDATE if exists)
	err := s.pushRepo.Upsert(ctx, subscription)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Push subscription saved for user %s: %s", userID.String(), subscription.Endpoint)

	return dto.PushSubscriptionToResponse(subscription), nil
}

func (s *PushServiceImpl) Unsubscribe(ctx context.Context, userID uuid.UUID, req *dto.PushSubscriptionRequest) error {
	err := s.pushRepo.Delete(ctx, userID, req.Endpoint)
	if err != nil {
		return err
	}

	log.Printf("🗑️  Push subscription removed for user %s: %s", userID.String(), req.Endpoint)

	return nil
}

func (s *PushServiceImpl) SendToUser(ctx context.Context, userID uuid.UUID, payload *dto.PushNotificationPayload) error {
	// Get all active subscriptions for the user
	subscriptions, err := s.pushRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error getting subscriptions for user %s: %v", userID.String(), err)
		return err
	}

	if len(subscriptions) == 0 {
		log.Printf("📭 No push subscriptions found for user %s", userID.String())
		return nil
	}

	// Convert payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Send to all subscriptions concurrently
	for _, sub := range subscriptions {
		go func(subscription *dto.PushSubscriptionRequest, endpoint string) {
			// Create webpush subscription
			webpushSub := &webpush.Subscription{
				Endpoint: subscription.Endpoint,
				Keys: webpush.Keys{
					P256dh: subscription.Keys.P256dh,
					Auth:   subscription.Keys.Auth,
				},
			}

			// Send notification
			resp, err := webpush.SendNotification(payloadJSON, webpushSub, &webpush.Options{
				Subscriber:      s.config.VAPID.Subject,
				VAPIDPublicKey:  s.config.VAPID.PublicKey,
				VAPIDPrivateKey: s.config.VAPID.PrivateKey,
				TTL:             30,
			})

			if err != nil {
				log.Printf("❌ Push notification error for %s: %v", endpoint, err)
				return
			}
			defer resp.Body.Close()

			// Check response status
			if resp.StatusCode == 201 {
				log.Printf("✅ Push notification sent successfully to: %s", endpoint)
			} else {
				log.Printf("⚠️  Push notification failed with status %d for: %s", resp.StatusCode, endpoint)

				// If subscription expired (410 Gone or 404 Not Found), delete it
				if resp.StatusCode == 410 || resp.StatusCode == 404 {
					log.Printf("🗑️  Removing expired subscription: %s", endpoint)
					s.pushRepo.DeleteByEndpoint(context.Background(), endpoint)
				}
			}
		}(&dto.PushSubscriptionRequest{
			Endpoint: sub.Endpoint,
			Keys: dto.PushSubscriptionKeys{
				P256dh: sub.P256dh,
				Auth:   sub.Auth,
			},
		}, sub.Endpoint)
	}

	log.Printf("📤 Push notifications sent to %d subscription(s) for user %s", len(subscriptions), userID.String())

	return nil
}

func (s *PushServiceImpl) GetPublicKey() string {
	return s.config.VAPID.PublicKey
}

// Compiler check to ensure implementation satisfies interface
var _ services.PushService = (*PushServiceImpl)(nil)
