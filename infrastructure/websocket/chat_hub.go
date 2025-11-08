package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/redis"
)

// ChatHub manages chat-specific WebSocket connections
type ChatHub struct {
	// Client management
	clients      map[uuid.UUID]*ChatClient // userID -> client
	clientsMutex sync.RWMutex

	// Channels
	register   chan *ChatClient
	unregister chan *ChatClient
	broadcast  chan *ChatMessage

	// Services
	messageService      services.MessageService
	conversationService services.ConversationService
	blockService        services.BlockService
	redisService        *redis.RedisService
	pushService         services.PushService

	// Repositories
	conversationRepo repositories.ConversationRepository
	followRepo       repositories.FollowRepository

	// Context
	ctx    context.Context
	cancel context.CancelFunc
}

// ChatClient represents a connected chat user
type ChatClient struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *ChatHub
	Ready  chan bool // Signal when WritePump is ready
}

// ChatMessage represents a WebSocket message
type ChatMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload,omitempty"`
	Error   *ChatError             `json:"error,omitempty"`
}

// ChatError represents an error message
type ChatError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewChatHub creates a new ChatHub
func NewChatHub(
	messageService services.MessageService,
	conversationService services.ConversationService,
	blockService services.BlockService,
	redisService *redis.RedisService,
	conversationRepo repositories.ConversationRepository,
	followRepo repositories.FollowRepository,
	pushService services.PushService,
) *ChatHub {
	ctx, cancel := context.WithCancel(context.Background())

	return &ChatHub{
		clients:             make(map[uuid.UUID]*ChatClient),
		register:            make(chan *ChatClient, 10),
		unregister:          make(chan *ChatClient, 10),
		broadcast:           make(chan *ChatMessage, 256),
		messageService:      messageService,
		conversationService: conversationService,
		blockService:        blockService,
		redisService:        redisService,
		conversationRepo:    conversationRepo,
		followRepo:          followRepo,
		pushService:         pushService,
		ctx:                 ctx,
		cancel:              cancel,
	}
}

// Run starts the hub's main loop
func (h *ChatHub) Run() {
	log.Println("🚀 ChatHub started")

	// Start Redis Pub/Sub listener
	go h.listenRedisPubSub()

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)

		case <-h.ctx.Done():
			log.Println("🛑 ChatHub stopped")
			return
		}
	}
}

// RegisterClient handles client registration (exported for handler)
func (h *ChatHub) RegisterClient(client *ChatClient) {
	h.register <- client
}

// SendToUser sends message to specific user (exported for external use)
func (h *ChatHub) SendToUser(userID uuid.UUID, message *ChatMessage) {
	h.sendToUser(userID, message)
}

// registerClient is internal handler for client registration
func (h *ChatHub) registerClient(client *ChatClient) {
	h.clientsMutex.Lock()
	h.clients[client.UserID] = client
	h.clientsMutex.Unlock()

	// Set user online in Redis
	if err := h.redisService.SetUserOnline(h.ctx, client.UserID); err != nil {
		log.Printf("Failed to set user online in Redis: %v", err)
	}

	log.Printf("✅ Chat client registered: UserID=%s, Total clients=%d", client.UserID, len(h.clients))

	// Send connection success message
	h.sendToClient(client, &ChatMessage{
		Type: "connection.success",
		Payload: map[string]interface{}{
			"userId":      client.UserID.String(),
			"connectedAt": time.Now().Format(time.RFC3339),
		},
	})

	// Send initial online status of relevant users
	go h.sendInitialOnlineStatus(client)

	// Broadcast online status to friends
	go h.broadcastOnlineStatus(client.UserID, true)
}

// unregisterClient handles client disconnection
func (h *ChatHub) unregisterClient(client *ChatClient) {
	h.clientsMutex.Lock()
	if _, ok := h.clients[client.UserID]; ok {
		delete(h.clients, client.UserID)
		close(client.Send)
	}
	h.clientsMutex.Unlock()

	// Set user offline in Redis
	if err := h.redisService.SetUserOffline(h.ctx, client.UserID); err != nil {
		log.Printf("Failed to set user offline in Redis: %v", err)
	}

	log.Printf("❌ Chat client unregistered: UserID=%s, Total clients=%d", client.UserID, len(h.clients))

	// Broadcast offline status to friends
	go h.broadcastOnlineStatus(client.UserID, false)
}

// broadcastMessage sends message to target clients
func (h *ChatHub) broadcastMessage(message *ChatMessage) {
	// Extract target userID if specified
	if userIDStr, ok := message.Payload["targetUserId"].(string); ok {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			h.sendToUser(userID, message)
			return
		}
	}

	// If no target, send to all clients (for system messages)
	h.clientsMutex.RLock()
	defer h.clientsMutex.RUnlock()

	for _, client := range h.clients {
		h.sendToClient(client, message)
	}
}

// sendToUser sends message to specific user
func (h *ChatHub) sendToUser(userID uuid.UUID, message *ChatMessage) {
	h.clientsMutex.RLock()
	client, exists := h.clients[userID]
	h.clientsMutex.RUnlock()

	if exists {
		// User connected to this server - send directly
		h.sendToClient(client, message)
	} else {
		// User might be on another server - publish to Redis
		if err := h.redisService.PublishToUser(h.ctx, userID, message); err != nil {
			log.Printf("Failed to publish to Redis: %v", err)
		}
	}
}

// sendToClient sends message to client's send channel
func (h *ChatHub) sendToClient(client *ChatClient, message *ChatMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	select {
	case client.Send <- data:
		// Message sent successfully
	default:
		// Client's send buffer is full, close connection
		log.Printf("Client send buffer full, closing connection: %s", client.UserID)
		h.unregister <- client
	}
}

// listenRedisPubSub listens for messages from Redis Pub/Sub
func (h *ChatHub) listenRedisPubSub() {
	// Subscribe to all user channels for this server
	// This is a simplified version - in production you'd subscribe to specific channels
	log.Println("🔴 Redis Pub/Sub listener started")

	// For now, we'll implement user-specific subscriptions when clients connect
	// Full implementation will be in Phase 2
}

// Stop gracefully shuts down the hub
func (h *ChatHub) Stop() {
	h.cancel()

	// Close all client connections
	h.clientsMutex.Lock()
	for _, client := range h.clients {
		close(client.Send)
	}
	h.clientsMutex.Unlock()

	log.Println("✓ ChatHub stopped gracefully")
}

// IsUserOnline checks if user is connected to this server
func (h *ChatHub) IsUserOnline(userID uuid.UUID) bool {
	h.clientsMutex.RLock()
	defer h.clientsMutex.RUnlock()

	_, exists := h.clients[userID]
	return exists
}

// GetOnlineCount returns total online clients
func (h *ChatHub) GetOnlineCount() int {
	h.clientsMutex.RLock()
	defer h.clientsMutex.RUnlock()

	return len(h.clients)
}

// broadcastOnlineStatus broadcasts user online/offline status to relevant users
func (h *ChatHub) broadcastOnlineStatus(userID uuid.UUID, isOnline bool) {
	ctx := context.Background()

	// Use map to collect unique user IDs who should receive the status
	recipientMap := make(map[uuid.UUID]bool)

	// 1. Get user's followers (people who follow this user)
	followers, err := h.followRepo.GetFollowers(ctx, userID, 0, 10000)
	if err != nil {
		log.Printf("Failed to get followers for online status broadcast: %v", err)
	} else {
		for _, follower := range followers {
			recipientMap[follower.ID] = true
		}
	}

	// 2. Get user's following (people this user follows)
	following, err := h.followRepo.GetFollowing(ctx, userID, 0, 10000)
	if err != nil {
		log.Printf("Failed to get following for online status broadcast: %v", err)
	} else {
		for _, follow := range following {
			recipientMap[follow.ID] = true
		}
	}

	// 3. Get conversation partners (people who have active conversations with this user)
	conversations, err := h.conversationRepo.ListByUser(ctx, userID, nil, 1000)
	if err != nil {
		log.Printf("Failed to get conversations for online status broadcast: %v", err)
	} else {
		for _, conv := range conversations {
			// Add the other participant in the conversation
			if conv.User1ID == userID {
				recipientMap[conv.User2ID] = true
			} else {
				recipientMap[conv.User1ID] = true
			}
		}
	}

	// Convert map to slice
	recipients := []uuid.UUID{}
	for userID := range recipientMap {
		recipients = append(recipients, userID)
	}

	// Determine event type
	eventType := "user.offline"
	if isOnline {
		eventType = "user.online"
	}

	// Get last seen timestamp
	_, lastSeen, _ := h.redisService.IsUserOnline(ctx, userID)

	// Broadcast to all recipients (followers + following + conversation partners)
	for _, recipientID := range recipients {
		h.sendToUser(recipientID, &ChatMessage{
			Type: eventType,
			Payload: map[string]interface{}{
				"userId":   userID.String(),
				"isOnline": isOnline,
				"lastSeen": lastSeen.Format(time.RFC3339),
			},
		})
	}

	log.Printf("📡 Broadcasted %s for user %s to %d users (followers + following + chat partners)", eventType, userID, len(recipients))
}

// sendInitialOnlineStatus sends the current online status of relevant users to a newly connected client
func (h *ChatHub) sendInitialOnlineStatus(client *ChatClient) {
	ctx := context.Background()

	// Use map to collect unique user IDs
	relevantUserMap := make(map[uuid.UUID]bool)

	// 1. Get followers
	followers, err := h.followRepo.GetFollowers(ctx, client.UserID, 0, 10000)
	if err != nil {
		log.Printf("Failed to get followers for initial status: %v", err)
	} else {
		for _, follower := range followers {
			relevantUserMap[follower.ID] = true
		}
	}

	// 2. Get following
	following, err := h.followRepo.GetFollowing(ctx, client.UserID, 0, 10000)
	if err != nil {
		log.Printf("Failed to get following for initial status: %v", err)
	} else {
		for _, follow := range following {
			relevantUserMap[follow.ID] = true
		}
	}

	// 3. Get conversation partners
	conversations, err := h.conversationRepo.ListByUser(ctx, client.UserID, nil, 1000)
	if err != nil {
		log.Printf("Failed to get conversations for initial status: %v", err)
	} else {
		for _, conv := range conversations {
			if conv.User1ID == client.UserID {
				relevantUserMap[conv.User2ID] = true
			} else {
				relevantUserMap[conv.User1ID] = true
			}
		}
	}

	// Check which users are online
	onlineUsers := []map[string]interface{}{}
	for userID := range relevantUserMap {
		// Check if user is online (in local clients or Redis)
		isOnline, lastSeen, _ := h.redisService.IsUserOnline(ctx, userID)

		onlineUsers = append(onlineUsers, map[string]interface{}{
			"userId":   userID.String(),
			"isOnline": isOnline,
			"lastSeen": lastSeen.Format(time.RFC3339),
		})
	}

	// Send initial status to the client
	h.sendToClient(client, &ChatMessage{
		Type: "initial.online.status",
		Payload: map[string]interface{}{
			"users": onlineUsers,
			"total": len(onlineUsers),
		},
	})

	log.Printf("📤 Sent initial online status to user %s: %d users", client.UserID, len(onlineUsers))
}
