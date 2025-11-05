package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupPushRoutes(api fiber.Router, h *handlers.Handlers) {
	push := api.Group("/push")

	// Public endpoint to get VAPID public key
	push.Get("/public-key", h.PushHandler.GetPublicKey)

	// Protected routes (require authentication)
	push.Use(middleware.Protected())
	push.Post("/subscribe", h.PushHandler.Subscribe)
	push.Post("/unsubscribe", h.PushHandler.Unsubscribe)
}
