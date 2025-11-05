package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupSearchRoutes(api fiber.Router, h *handlers.Handlers) {
	search := api.Group("/search")

	// Public search (with optional authentication)
	search.Get("/", h.SearchHandler.Search)
	search.Get("/popular", h.SearchHandler.GetPopularSearches)

	// Protected routes (require authentication)
	search.Use(middleware.Protected())
	search.Get("/history", h.SearchHandler.GetSearchHistory)
	search.Delete("/history", h.SearchHandler.ClearSearchHistory)
	search.Delete("/history/:id", h.SearchHandler.DeleteSearchHistoryItem)
}
