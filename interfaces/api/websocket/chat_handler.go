package websocket

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	chatWebsocket "gofiber-template/infrastructure/websocket"
)

type ChatWebSocketHandler struct {
	chatHub *chatWebsocket.ChatHub
}

func NewChatWebSocketHandler(chatHub *chatWebsocket.ChatHub) *ChatWebSocketHandler {
	return &ChatWebSocketHandler{
		chatHub: chatHub,
	}
}

// WebSocketUpgrade checks if the request is a WebSocket upgrade
func (h *ChatWebSocketHandler) WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// HandleChatWebSocket handles WebSocket connections for chat
func (h *ChatWebSocketHandler) HandleChatWebSocket(c *websocket.Conn) {
	var userID uuid.UUID

	// Get user from locals (set by Protected middleware)
	if userContext := c.Locals("userID"); userContext != nil {
		if id, ok := userContext.(uuid.UUID); ok {
			userID = id
		}
	}

	// If no user from middleware, deny connection
	if userID == uuid.Nil {
		log.Printf("❌ Chat WebSocket: Unauthorized connection attempt")
		c.WriteJSON(map[string]interface{}{
			"type":  "error",
			"error": map[string]string{
				"code":    "unauthorized",
				"message": "Authentication required",
			},
		})
		c.Close()
		return
	}

	log.Printf("✅ Chat WebSocket: User connected: %s", userID)

	// Create client with ready channel for synchronization
	client := &chatWebsocket.ChatClient{
		UserID: userID,
		Conn:   c,
		Send:   make(chan []byte, 256),
		Hub:    h.chatHub,
		Ready:  make(chan bool),
	}

	// Start WritePump first (it will signal when ready)
	go client.WritePump()

	// Wait for WritePump to be ready (with timeout)
	select {
	case <-client.Ready:
		log.Printf("✅ WritePump ready for user: %s", userID)
	case <-time.After(5 * time.Second):
		log.Printf("⚠️ WritePump timeout for user: %s", userID)
		c.Close()
		return
	}

	// Now it's safe to register client (will send connection.success)
	h.chatHub.RegisterClient(client)

	// Start ReadPump - this will BLOCK until connection closes
	// This is important: the handler must not return while connection is active
	client.ReadPump()

	// When ReadPump returns, the connection has been closed
	log.Printf("📴 Chat WebSocket: Connection closed for user: %s", userID)
}
