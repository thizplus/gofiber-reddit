package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupProfileRoutes(api fiber.Router, h *handlers.Handlers) {
	profiles := api.Group("/profiles")

	// Public route with optional authentication
	profiles.Get("/:username", middleware.Optional(), h.ProfileHandler.GetPublicProfile)
}
