package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupChatWebSocketRoutes(app *fiber.App, h *handlers.Handlers) {
	// Chat WebSocket endpoint with JWT authentication from query parameter
	app.Use("/chat/ws", middleware.WebSocketProtected())
	app.Use("/chat/ws", h.ChatWSHandler.WebSocketUpgrade)
	app.Get("/chat/ws", websocket.New(h.ChatWSHandler.HandleChatWebSocket))
}
