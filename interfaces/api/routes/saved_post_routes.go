package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupSavedPostRoutes(api fiber.Router, h *handlers.Handlers) {
	saved := api.Group("/saved")
	saved.Use(middleware.Protected())

	saved.Post("/posts/:postId", h.SavedPostHandler.SavePost)
	saved.Delete("/posts/:postId", h.SavedPostHandler.UnsavePost)
	saved.Get("/posts/:postId/status", h.SavedPostHandler.IsSaved)
	saved.Get("/posts", h.SavedPostHandler.GetSavedPosts)
}
