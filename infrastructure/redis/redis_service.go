package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(redisClient *RedisClient) *RedisService {
	return &RedisService{
		client: redisClient.client,
	}
}

// ========== Online Status ==========

// SetUserOnline marks a user as online with 60s TTL
func (r *RedisService) SetUserOnline(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("online:%s", userID.String())
	timestamp := time.Now().Unix()

	return r.client.Set(ctx, key, timestamp, 60*time.Second).Err()
}

// SetUserOffline marks a user as offline (no TTL - persists as last seen)
func (r *RedisService) SetUserOffline(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("online:%s", userID.String())
	timestamp := time.Now().Unix()

	// Set with no TTL (will remain as last seen)
	return r.client.Set(ctx, key, timestamp, 0).Err()
}

// IsUserOnline checks if a user is online and returns last seen timestamp
func (r *RedisService) IsUserOnline(ctx context.Context, userID uuid.UUID) (bool, time.Time, error) {
	key := fmt.Sprintf("online:%s", userID.String())

	// Check if key exists (with TTL = online)
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return false, time.Time{}, err
	}

	// If TTL > 0, user is online
	isOnline := ttl > 0

	// Get last seen timestamp
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, time.Time{}, nil // Never seen
		}
		return false, time.Time{}, err
	}

	timestamp, _ := strconv.ParseInt(val, 10, 64)
	lastSeen := time.Unix(timestamp, 0)

	return isOnline, lastSeen, nil
}

// GetBulkOnlineStatus retrieves online status for multiple users (efficient MGET)
func (r *RedisService) GetBulkOnlineStatus(ctx context.Context, userIDs []uuid.UUID) (map[string]bool, error) {
	if len(userIDs) == 0 {
		return map[string]bool{}, nil
	}

	// Build keys
	keys := make([]string, len(userIDs))
	for i, id := range userIDs {
		keys[i] = fmt.Sprintf("online:%s", id.String())
	}

	// MGET all keys
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	// Check TTL for each key to determine online status
	result := make(map[string]bool)
	for i, userID := range userIDs {
		if vals[i] != nil {
			// Check TTL
			ttl, _ := r.client.TTL(ctx, keys[i]).Result()
			result[userID.String()] = ttl > 0
		} else {
			result[userID.String()] = false
		}
	}

	return result, nil
}

// ========== Unread Counts ==========

// GetTotalUnreadCount retrieves total unread message count for a user
func (r *RedisService) GetTotalUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	key := fmt.Sprintf("unread:total:%s", userID.String())

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // No unread
		}
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// IncrementTotalUnread increments total unread count for a user
func (r *RedisService) IncrementTotalUnread(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("unread:total:%s", userID.String())
	return r.client.Incr(ctx, key).Err()
}

// DecrementTotalUnread decrements total unread count (prevents negative values)
func (r *RedisService) DecrementTotalUnread(ctx context.Context, userID uuid.UUID, count int) error {
	if count <= 0 {
		return nil
	}

	key := fmt.Sprintf("unread:total:%s", userID.String())

	// DECRBY
	err := r.client.DecrBy(ctx, key, int64(count)).Err()
	if err != nil {
		return err
	}

	// Ensure it doesn't go negative
	val, err := r.client.Get(ctx, key).Int()
	if err == nil && val < 0 {
		r.client.Set(ctx, key, 0, 0)
	}

	return nil
}

// GetConversationUnreadCount retrieves unread count for a specific conversation
func (r *RedisService) GetConversationUnreadCount(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) (int, error) {
	key := fmt.Sprintf("unread:conv:%s:%s", userID.String(), conversationID.String())

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return strconv.Atoi(val)
}

// IncrementConversationUnread increments unread count for a specific conversation
func (r *RedisService) IncrementConversationUnread(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) error {
	key := fmt.Sprintf("unread:conv:%s:%s", userID.String(), conversationID.String())
	return r.client.Incr(ctx, key).Err()
}

// ResetConversationUnread resets unread count for a conversation and returns the previous count
func (r *RedisService) ResetConversationUnread(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) (int, error) {
	key := fmt.Sprintf("unread:conv:%s:%s", userID.String(), conversationID.String())

	// Get current count before deleting
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Already 0
		}
		return 0, err
	}

	count, _ := strconv.Atoi(val)

	// Delete key
	r.client.Del(ctx, key)

	return count, nil
}

// InvalidateUnreadCache invalidates all unread caches for a user (use when rebuilding)
func (r *RedisService) InvalidateUnreadCache(ctx context.Context, userID uuid.UUID) error {
	// Delete total unread
	totalKey := fmt.Sprintf("unread:total:%s", userID.String())
	r.client.Del(ctx, totalKey)

	// Delete all conversation unread keys
	pattern := fmt.Sprintf("unread:conv:%s:*", userID.String())
	keys, err := r.scanKeys(ctx, pattern)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		r.client.Del(ctx, keys...)
	}

	return nil
}

// scanKeys scans Redis for keys matching a pattern (internal helper)
func (r *RedisService) scanKeys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	return keys, iter.Err()
}

// ========== Last Message Cache ==========

// CacheLastMessage caches the last message of a conversation with 1 hour TTL
func (r *RedisService) CacheLastMessage(ctx context.Context, conversationID uuid.UUID, messageID uuid.UUID, senderID uuid.UUID, content *string, messageType string, createdAt time.Time) error {
	key := fmt.Sprintf("last_msg:%s", conversationID.String())

	// Store as hash
	data := map[string]interface{}{
		"id":         messageID.String(),
		"sender_id":  senderID.String(),
		"content":    "",
		"type":       messageType,
		"created_at": createdAt.Format(time.RFC3339),
	}

	if content != nil {
		data["content"] = *content
	}

	// HSET with 1 hour TTL
	pipe := r.client.Pipeline()
	pipe.HSet(ctx, key, data)
	pipe.Expire(ctx, key, 1*time.Hour)
	_, err := pipe.Exec(ctx)

	return err
}

// GetCachedLastMessage retrieves cached last message for a conversation
func (r *RedisService) GetCachedLastMessage(ctx context.Context, conversationID uuid.UUID) (map[string]string, error) {
	key := fmt.Sprintf("last_msg:%s", conversationID.String())

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, redis.Nil // Cache miss
	}

	return result, nil
}

// InvalidateLastMessage invalidates cached last message for a conversation
func (r *RedisService) InvalidateLastMessage(ctx context.Context, conversationID uuid.UUID) error {
	key := fmt.Sprintf("last_msg:%s", conversationID.String())
	return r.client.Del(ctx, key).Err()
}

// ========== Pub/Sub ==========

// PublishToUser publishes a message to a user's channel (for multi-server WebSocket)
func (r *RedisService) PublishToUser(ctx context.Context, userID uuid.UUID, message interface{}) error {
	channel := fmt.Sprintf("chat:user:%s", userID.String())

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.client.Publish(ctx, channel, data).Err()
}

// SubscribeToUser subscribes to a user's channel
func (r *RedisService) SubscribeToUser(ctx context.Context, userID uuid.UUID) *redis.PubSub {
	channel := fmt.Sprintf("chat:user:%s", userID.String())
	return r.client.Subscribe(ctx, channel)
}

// UnsubscribeUser closes a Pub/Sub subscription
func (r *RedisService) UnsubscribeUser(ctx context.Context, pubsub *redis.PubSub) error {
	return pubsub.Close()
}
