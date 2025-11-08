package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
)

func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
	// Setup health and root routes
	SetupHealthRoutes(app)

	// Setup SEO routes (needs app, not api group)
	SetupSEORoutes(app, h)

	// API version group
	api := app.Group("/api/v1")

	// Setup authentication routes
	SetupAuthRoutes(api, h)
	SetupUserRoutes(api, h)
	SetupProfileRoutes(api, h)

	// Setup social media routes
	SetupPostRoutes(api, h)
	SetupCommentRoutes(api, h)
	SetupVoteRoutes(api, h)
	SetupFollowRoutes(api, h)
	SetupSavedPostRoutes(api, h)
	SetupNotificationRoutes(api, h)
	SetupTagRoutes(api, h)
	SetupSearchRoutes(api, h)
	SetupMediaRoutes(api, h)
	SetupPushRoutes(api, h)

	// Setup chat routes
	SetupChatRoutes(api, h)

	// Setup upload routes
	SetupUploadRoutes(api, h)

	// Setup legacy routes (can be removed if not needed)
	SetupTaskRoutes(api, h)
	SetupFileRoutes(api, h)
	SetupJobRoutes(api, h)

	// Setup WebSocket routes (needs app, not api group)
	SetupWebSocketRoutes(app)
	SetupChatWebSocketRoutes(app, h)
}