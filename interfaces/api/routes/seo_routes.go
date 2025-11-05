package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
)

func SetupSEORoutes(app *fiber.App, h *handlers.Handlers) {
	// SEO routes (no /api prefix)
	app.Get("/sitemap.xml", h.SEOHandler.GetSitemap)
	app.Get("/robots.txt", h.SEOHandler.GetRobotsTxt)
	app.Get("/rss.xml", h.SEOHandler.GetRSSFeed)

	// Open Graph API (with /api prefix)
	api := app.Group("/api/v1")
	og := api.Group("/og")
	og.Get("/post/:id", h.SEOHandler.GetOpenGraphData)
}
