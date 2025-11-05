package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupPostRoutes(api fiber.Router, h *handlers.Handlers) {
	posts := api.Group("/posts")

	// Public routes (with optional authentication)
	posts.Get("/", middleware.Optional(), h.PostHandler.ListPosts)
	posts.Get("/:id", middleware.Optional(), h.PostHandler.GetPost)
	posts.Get("/author/:authorId", middleware.Optional(), h.PostHandler.ListPostsByAuthor)
	posts.Get("/tag/:tagName", middleware.Optional(), h.PostHandler.ListPostsByTag)
	posts.Get("/tag-id/:tagId", middleware.Optional(), h.PostHandler.ListPostsByTagID)
	posts.Get("/search", middleware.Optional(), h.PostHandler.SearchPosts)
	posts.Get("/:id/crossposts", middleware.Optional(), h.PostHandler.GetCrossposts)

	// Protected routes (require authentication)
	posts.Use(middleware.Protected())
	posts.Post("/", h.PostHandler.CreatePost)
	posts.Put("/:id", h.PostHandler.UpdatePost)
	posts.Delete("/:id", h.PostHandler.DeletePost)
	posts.Post("/:id/crosspost", h.PostHandler.CreateCrosspost)
	posts.Get("/feed", h.PostHandler.GetFeed)
}
