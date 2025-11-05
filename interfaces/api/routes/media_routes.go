package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupMediaRoutes(api fiber.Router, h *handlers.Handlers) {
	media := api.Group("/media")

	// Public route
	media.Get("/:id", h.MediaHandler.GetMedia)
	media.Get("/user/:userId", h.MediaHandler.GetUserMedia)

	// Protected routes (require authentication)
	media.Use(middleware.Protected())
	media.Post("/upload/image", h.MediaHandler.UploadImage)
	media.Post("/upload/video", h.MediaHandler.UploadVideo)
	media.Delete("/:id", h.MediaHandler.DeleteMedia)
}
