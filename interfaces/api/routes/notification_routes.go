package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupNotificationRoutes(api fiber.Router, h *handlers.Handlers) {
	notifications := api.Group("/notifications")
	notifications.Use(middleware.Protected())

	// Settings (must be before /:id to avoid route conflict)
	notifications.Get("/settings", h.NotificationHandler.GetSettings)
	notifications.Put("/settings", h.NotificationHandler.UpdateSettings)

	// Get notifications
	notifications.Get("/", h.NotificationHandler.GetNotifications)
	notifications.Get("/unread", h.NotificationHandler.GetUnreadNotifications)
	notifications.Get("/unread/count", h.NotificationHandler.GetUnreadCount)
	notifications.Get("/:id", h.NotificationHandler.GetNotification)

	// Mark as read
	notifications.Put("/:id/read", h.NotificationHandler.MarkAsRead)
	notifications.Put("/read-all", h.NotificationHandler.MarkAllAsRead)

	// Delete notifications
	notifications.Delete("/:id", h.NotificationHandler.DeleteNotification)
	notifications.Delete("/", h.NotificationHandler.DeleteAllNotifications)
}
