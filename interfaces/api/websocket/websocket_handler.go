package websocket

import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"gofiber-template/pkg/utils"
	websocketManager "gofiber-template/infrastructure/websocket"
)

type WebSocketHandler struct{}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{}
}

func (h *WebSocketHandler) WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *WebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	var userID uuid.UUID
	var roomID string

	// Try to get user from context (set by Optional middleware)
	if userContext := c.Locals("user"); userContext != nil {
		if user, ok := userContext.(*utils.UserContext); ok {
			userID = user.ID
		}
	}

	// If no user from middleware, try token from query parameter
	if userID == uuid.Nil {
		token := c.Query("token")
		if token != "" {
			jwtSecret := os.Getenv("JWT_SECRET")
			userCtx, err := utils.ValidateTokenStringToUUID(token, jwtSecret)
			if err == nil {
				userID = userCtx.ID
				log.Printf("✅ WebSocket: Token validated from query param for user: %s (%s)", userCtx.Email, userCtx.ID)
			} else {
				log.Printf("⚠️  WebSocket: Invalid token from query param: %v", err)
			}
		}
	}

	// If no user context, generate anonymous user ID
	if userID == uuid.Nil {
		userID = uuid.New()
		log.Printf("WebSocket: Anonymous user connected with ID: %s", userID.String())
	} else {
		log.Printf("WebSocket: Authenticated user connected: %s", userID.String())
	}

	roomID = c.Query("room", "")

	websocketManager.Manager.RegisterClient(c, userID, roomID)

	defer func() {
		websocketManager.Manager.UnregisterClient(c)
	}()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		websocketManager.HandleWebSocketMessage(c, messageType, message)
	}
}