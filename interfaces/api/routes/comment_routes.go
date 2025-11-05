package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupCommentRoutes(api fiber.Router, h *handlers.Handlers) {
	comments := api.Group("/comments")

	// Public routes (with optional authentication)
	comments.Get("/:id", middleware.Optional(), h.CommentHandler.GetComment)
	comments.Get("/post/:postId", middleware.Optional(), h.CommentHandler.ListCommentsByPost)
	comments.Get("/post/:postId/tree", middleware.Optional(), h.CommentHandler.GetCommentTree)
	comments.Get("/author/:authorId", middleware.Optional(), h.CommentHandler.ListCommentsByAuthor)
	comments.Get("/:id/replies", middleware.Optional(), h.CommentHandler.ListReplies)
	comments.Get("/:id/parent-chain", middleware.Optional(), h.CommentHandler.GetParentChain)

	// Protected routes (require authentication)
	comments.Use(middleware.Protected())
	comments.Post("/", h.CommentHandler.CreateComment)
	comments.Put("/:id", h.CommentHandler.UpdateComment)
	comments.Delete("/:id", h.CommentHandler.DeleteComment)
}
