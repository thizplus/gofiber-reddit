package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupVoteRoutes(api fiber.Router, h *handlers.Handlers) {
	votes := api.Group("/votes")

	// Public routes (with optional authentication)
	votes.Get("/:targetType/:targetId", h.VoteHandler.GetVote)
	votes.Get("/:targetType/:targetId/count", h.VoteHandler.GetVoteCount)

	// Protected routes (require authentication)
	votes.Use(middleware.Protected())
	votes.Post("/", h.VoteHandler.Vote)
	votes.Delete("/:targetType/:targetId", h.VoteHandler.Unvote)
	votes.Get("/user", h.VoteHandler.GetUserVotes)
}
