package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupUploadRoutes(api fiber.Router, h *handlers.Handlers) {
	upload := api.Group("/upload")

	// Protected routes (require authentication)
	upload.Use(middleware.Protected())
	upload.Post("/file", h.FileUploadHandler.UploadFile)
}
