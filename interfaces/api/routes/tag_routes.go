package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
)

func SetupTagRoutes(api fiber.Router, h *handlers.Handlers) {
	tags := api.Group("/tags")

	// All tag routes are public
	tags.Get("/", h.TagHandler.ListTags)
	tags.Get("/popular", h.TagHandler.GetPopularTags)
	tags.Get("/search", h.TagHandler.SearchTags)
	tags.Get("/:id", h.TagHandler.GetTag)
	tags.Get("/name/:name", h.TagHandler.GetTagByName)
}
