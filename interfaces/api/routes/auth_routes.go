package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
)

func SetupAuthRoutes(api fiber.Router, h *handlers.Handlers) {
	auth := api.Group("/auth")

	// Standard authentication
	auth.Post("/register", h.UserHandler.Register)
	auth.Post("/login", h.UserHandler.Login)

	// OAuth authentication
	auth.Get("/google", h.OAuthHandler.GetGoogleAuthURL)
	auth.Get("/google/callback", h.OAuthHandler.GoogleCallback)
	auth.Post("/exchange", h.OAuthHandler.ExchangeCodeForToken)
}