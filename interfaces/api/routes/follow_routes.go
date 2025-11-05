package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupFollowRoutes(api fiber.Router, h *handlers.Handlers) {
	follows := api.Group("/follows")

	// Public routes (with optional authentication)
	follows.Get("/user/:userId/followers", h.FollowHandler.GetFollowers)
	follows.Get("/user/:userId/following", h.FollowHandler.GetFollowing)
	follows.Get("/user/:userId/status", h.FollowHandler.IsFollowing)

	// Protected routes (require authentication)
	follows.Use(middleware.Protected())
	follows.Post("/user/:userId", h.FollowHandler.Follow)
	follows.Delete("/user/:userId", h.FollowHandler.Unfollow)
	follows.Get("/mutuals", h.FollowHandler.GetMutualFollows)
}
